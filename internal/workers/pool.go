package workers

import (
	"fmt"
	"github.com/dragaera/probius/internal/config"
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
	Config *config.Config
	pool   *work.WorkerPool
	Redis  *redis.Pool
	DB     *gorm.DB
}

func (pool *Pool) Run() error {
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

	ctxt := JobContext{
		db: pool.DB,
	}
	workerPool := work.NewWorkerPool(
		ctxt,
		uint(pool.Config.Worker.Concurrency), // Number of worker processes per pool process
		pool.Config.Worker.Namespace,
		pool.Redis,
	)
	pool.pool = workerPool

	workerPool.Middleware(LogStart)
	workerPool.Middleware(pool.InjectDB)

	workerPool.Job("check_last_replay", CheckLastReplay)

	return pool, nil
}

type JobContext struct {
	id int
	db *gorm.DB
}

func LogStart(ctxt *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Job starting: %v\n", job.Name)

	return next()
}

func (pool *Pool) InjectDB(ctxt *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
	ctxt.db = pool.DB

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
		log.Printf("New replay for user %+v found.", user)
		// TODO
	} else {
		log.Printf("No new replay for user %+v found.", user)
	}

	log.Printf("PLACEHOLDER %v", replay)

	return nil
}
