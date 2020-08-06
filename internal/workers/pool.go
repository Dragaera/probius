package workers

import (
	"fmt"
	"github.com/dragaera/probius/internal/config"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Pool struct {
	Config *config.Config
	pool   *work.WorkerPool
	Redis  *redis.Pool
	db     *pgxpool.Pool
}

func (pool *Pool) Run(db *pgxpool.Pool) error {
	pool.db = db
	defer pool.db.Close()

	pool.pool.Start()
	log.Print("Worker pool started")
	defer pool.pool.Stop()

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

	workerPool := work.NewWorkerPool(
		JobContext{},
		uint(pool.Config.Worker.Concurrency), // Number of worker processes per pool process
		pool.Config.Worker.Namespace,
		pool.Redis,
	)
	pool.pool = workerPool

	workerPool.Middleware((*JobContext).LogStart)
	// Further middlewares to look up eg Discord user from DB, if needed

	workerPool.Job("check_last_replay", (*JobContext).CheckLastReplay)

	return pool, nil
}

type JobContext struct {
	id int
}

func (c *JobContext) LogStart(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Job starting: %v\n", job.Name)

	return next()
}

func (c *JobContext) CheckLastReplay(job *work.Job) error {
	fmt.Printf("CheckLastReplay: %+v\n", c)

	return nil
}
