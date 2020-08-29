package workers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/discord"
	"github.com/dragaera/probius/internal/persistence"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Pool struct {
	Config   *config.Config
	Session  *discordgo.Session
	pool     *work.WorkerPool
	enqueuer *work.Enqueuer
	Redis    *redis.Pool
	DB       *gorm.DB
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

	// Periodic jobs
	// seconds hours minutes day-of-month month week-of-day
	// (as per https://github.com/gocraft/work/)
	workerPool.PeriodicallyEnqueue("0 * * * * *", "check_stale_players")

	return pool, nil
}

type JobContext struct {
	id       int
	db       *gorm.DB
	config   *config.Config
	enqueuer *work.Enqueuer
	session  *discordgo.Session
}

func LogStart(ctxt *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Job starting: %v\n", job.Name)

	return next()
}

func (pool *Pool) EnrichContext(ctxt *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
	ctxt.db = pool.DB
	ctxt.config = pool.Config
	ctxt.enqueuer = pool.enqueuer
	ctxt.session = pool.Session

	return next()
}

func CheckLastReplay(ctxt *JobContext, job *work.Job) error {
	sc2rID := int(job.ArgInt64("id"))
	if err := job.ArgError(); err != nil {
		return fmt.Errorf("Missing SC2ReplayStatsuser ID: %v", err)
	}

	user := persistence.SC2ReplayStatsUser{}
	err := ctxt.db.First(&user, "id = ?", sc2rID).Error
	if err != nil {
		return err
	}

	replay, changed, err := user.UpdateLastReplay(ctxt.db)
	if err != nil {
		return err
	}

	if changed {
		log.Printf("New replay for user %v found.", user.ID)
		subscriptions, err := user.GetSubscriptions(ctxt.db)
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
				ctxt.session,
				sub.DiscordChannel.DiscordID,
				&embed,
			)
			if err != nil {
				log.Print("Error sending replay embed: ", err)
				return err
			}
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
