package main

import (
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/discord"
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

	bot, err := discord.Create(&discord.Bot{
		Config: cfg,
	})
	if err != nil {
		log.Fatal("Error while creating Discord bot: ", err)
	}

	orm, err := persistence.InitializeORM(cfg.DB)
	if err != nil {
		log.Fatal("Error while initializing ORM persistence layer: ", err)
	}

	log.Print("Starting Discord bot.")

	err = bot.Run(orm)
	if err != nil {
		log.Fatal("Error while starting Discord bot: ", err)
	}

	log.Print("Discord bot shut down.")
}
