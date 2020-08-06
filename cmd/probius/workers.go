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

	pool, err := workers.Create(
		&workers.Pool{
			Config: &cfg,
			Redis:  redis,
		},
	)
	if err != nil {
		log.Fatal("Error while creating worker pool: ", err)
	}

	db, err := persistence.InitializeDB(cfg.DB.DBURL())
	if err != nil {
		log.Fatal("Error while initializing persistence layer: ", err)
	}

	log.Print("Starting worker pool.")

	err = pool.Run(db)
	if err != nil {
		log.Fatal("Error while starting worker pool: ", err)
	}

	log.Print("Worker pool shut down.")
}
