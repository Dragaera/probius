package main

import (
	"fmt"
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/persistence"
	"github.com/gocraft/work"
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

	var enqueuer = work.NewEnqueuer(cfg.Worker.Namespace, redis)
	for i := 0; i < 1; i++ {
		_, err = enqueuer.Enqueue("check_last_replay", work.Q{"id": i + 1})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Job enqueued")
}
