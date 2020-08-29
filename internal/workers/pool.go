package workers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/discord"
	"github.com/dragaera/probius/internal/persistence"
	"github.com/dragaera/probius/internal/sc2replaystats"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/redigostore"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
)

type Pool struct {
	Config      *config.Config
	Session     *discordgo.Session
	pool        *work.WorkerPool
	enqueuer    *work.Enqueuer
	rateLimiter *throttled.GCRARateLimiter
	Redis       *redis.Pool
	DB          *gorm.DB
}

func (pool *Pool) Run() error {
	err := pool.Session.Open()
	if err != nil {
		return fmt.Errorf("Error connecting to Discord:", err)
	}
	defer pool.Session.Close()
	log.Print("Discord connection established")

	pool.pool.Start()
	defer pool.pool.Stop()
	log.Print("Worker pool started")

	// Terminate on ^c or SIGTERM
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Print("Worker pool terminating")

	return nil
}

func Create(pool *Pool) (*Pool, error) {
	if pool.Redis == nil {
		return pool, fmt.Errorf("Redis pool must not be nil")
	}

	if pool.Config == nil {
		return pool, fmt.Errorf("Config must not be nil")
	}

	if pool.DB == nil {
		return pool, fmt.Errorf("DB pool must not be nil")
	}

	dg, err := discordgo.New("Bot " + pool.Config.Discord.Token)
	if err != nil {
		return pool, fmt.Errorf("Error creating Discord session:", err)
	}
	pool.Session = dg

	store, err := redigostore.New(
		pool.Redis,
		"", // Prefix
		0,  // DB
	)
	if err != nil {
		return pool, fmt.Errorf("Error creating rate limiting Redis store: ", err)
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(pool.Config.SC2ReplayStats.RateLimitAverage),
		MaxBurst: pool.Config.SC2ReplayStats.RateLimitBurst,
	}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return pool, fmt.Errorf("Error creating rate limiter: ", err)
	}
	pool.rateLimiter = rateLimiter

	workerPool := work.NewWorkerPool(
		JobContext{},
		uint(pool.Config.Worker.Concurrency), // Number of worker processes per pool process
		pool.Config.Worker.Namespace,
		pool.Redis,
	)
	pool.pool = workerPool

	pool.enqueuer = work.NewEnqueuer(pool.Config.Worker.Namespace, pool.Redis)

	workerPool.Middleware(LogStart)
	workerPool.Middleware(pool.EnrichContext)

	// Job handlers
	workerPool.Job("check_last_replay", CheckLastReplay)
	workerPool.Job("check_stale_players", CheckStalePlayers)
	workerPool.Job("clear_stale_locks", ClearStaleLocks)

	// Periodic jobs
	// seconds hours minutes day-of-month month week-of-day
	// (as per https://github.com/gocraft/work/)
	workerPool.PeriodicallyEnqueue("0 * * * * *", "check_stale_players")
	workerPool.PeriodicallyEnqueue("0 * * * * *", "clear_stale_locks")

	return pool, nil
}

type JobContext struct {
	id          int
	db          *gorm.DB
	config      *config.Config
	enqueuer    *work.Enqueuer
	rateLimiter *throttled.GCRARateLimiter
	session     *discordgo.Session
}

func LogStart(ctxt *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Job starting: %v\n", job.Name)

	return next()
}

func (pool *Pool) EnrichContext(ctxt *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
	ctxt.db = pool.DB
	ctxt.config = pool.Config
	ctxt.enqueuer = pool.enqueuer
	ctxt.rateLimiter = pool.rateLimiter
	ctxt.session = pool.Session

	return next()
}

func CheckLastReplay(ctxt *JobContext, job *work.Job) error {
	sc2rID := int(job.ArgInt64("id"))
	if err := job.ArgError(); err != nil {
		return fmt.Errorf("Missing SC2ReplayStatsuser ID: %v", err)
	}

	limited, _, err := ctxt.rateLimiter.RateLimit("check_last_replay", 1)
	if err != nil {
		return fmt.Errorf("Unable to query rate limiter: %v", err)
	}

	if limited {
		// Backoff in [5, 59) seconds
		backoff := rand.Int63n(55) + 5
		log.Printf("Warning: Hit rate limit for check_last_replay. Rescheduling in: %vs", backoff)
		// Rate limited, reschedule for later time
		ctxt.enqueuer.EnqueueIn(
			"check_last_replay",
			backoff,
			work.Q{"id": sc2rID},
		)
		return nil
	}

	user := persistence.SC2ReplayStatsUser{}
	err = ctxt.db.First(&user, "id = ?", sc2rID).Error
	if err != nil {
		return err
	}

	replay, changed, err := user.UpdateLastReplay(ctxt.db)
	if err != nil {
		return err
	}

	if changed {
		log.Printf("New replay for user %v found.", user.ID)
		if err := notifySubscriptions(ctxt.db, ctxt.session, user, replay); err != nil {
			return err
		}
	} else {
		log.Printf("No new replay for user %v found.", user.ID)
	}

	return nil
}

func CheckStalePlayers(ctxt *JobContext, job *work.Job) error {
	users, err := persistence.SC2ReplayStatsUsersWithStaleData(
		ctxt.db,
		ctxt.config.SC2ReplayStats.UpdateInterval,
	)
	if err != nil {
		return err
	}

	for _, user := range users {
		err = user.LockForUpdate(ctxt.db)
		if err != nil {
			log.Printf("Error marking player for update: ", err)
			return err
		}
		ctxt.enqueuer.Enqueue("check_last_replay", work.Q{"id": user.ID})
	}

	return nil
}

func ClearStaleLocks(ctxt *JobContext, job *work.Job) error {
	return persistence.ClearStaleSC2ReplayStatsUpdateLocks(
		ctxt.db,
		ctxt.config.SC2ReplayStats.LockTTL,
	)
}

func notifySubscriptions(db *gorm.DB, session *discordgo.Session, user persistence.SC2ReplayStatsUser, replay sc2replaystats.Replay) error {
	subscriptions, err := user.GetSubscriptions(db)
	if err != nil {
		log.Print("Error retrieving user's subscriptions: ", err)
		return err
	}

	for _, sub := range subscriptions {
		embed := discord.BuildReplayEmbed(
			user.API(),
			replay,
		)

		err = sendEmbed(
			session,
			sub.DiscordChannel.DiscordID,
			&embed,
		)
		if err != nil {
			log.Print("Error sending replay embed: ", err)
			return err
		}
	}
	return nil
}

func sendEmbed(sess *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) error {
	_, err := sess.ChannelMessageSendEmbed(
		channelID,
		embed,
	)

	if err != nil {
		return fmt.Errorf("Error sending message: %v", err)
	}
	return nil
}
