package main

import (
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/persistence"
	"github.com/dragaera/probius/internal/workers"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Print("Error loading .env file: ", err)
	}

	// Will `log.Fatal()` if an env variable is missing
	cfg := config.ConfigFromEnv()

	// TODO: Error here is always nil, maybe redis library has a way to perform a
	// health check which InitializeRedis() could use?
	redis, _ := persistence.InitializeRedis(
		cfg.Redis.Host,
		cfg.Redis.Port,
	)

	orm, err := persistence.InitializeORM(cfg.DB)
	if err != nil {
		log.Fatal("Error while initializing persistence layer: ", err)
	}

	pool, err := workers.Create(
		&workers.Pool{
			Config: &cfg,
			Redis:  redis,
			DB:     orm,
		},
	)
	if err != nil {
		log.Fatal("Error while creating worker pool: ", err)
	}

	log.Print("Starting worker pool.")

	err = pool.Run()
	if err != nil {
		log.Fatal("Error while starting worker pool: ", err)
	}

	log.Print("Worker pool shut down.")
}
