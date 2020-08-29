package main

import (
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/persistence"
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

	orm, err := persistence.InitializeORM(cfg.DB)
	if err != nil {
		log.Fatal("Error while initializing ORM persistence layer: ", err)
	}

	log.Print("Starting automigration.")
	orm.AutoMigrate(
		&persistence.DiscordUser{},
		&persistence.DiscordGuild{},
		&persistence.DiscordChannel{},

		&persistence.SC2ReplayStatsUser{},
		&persistence.Subscription{},
	)
}
