package main

import (
	"fmt"
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/persistence"
	"github.com/gocraft/work"
	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
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
	fmt.Println("Initialized enqueuer: ", enqueuer)
	// for i := 0; i < 1; i++ {
	// 	_, err = enqueuer.Enqueue("check_last_replay", work.Q{"id": i + 1})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	orm, _ := persistence.InitializeORM(cfg.DB)
	orm.AutoMigrate(
		&persistence.Subscription{},
		&persistence.SC2ReplayStatsUser{},
		&persistence.DiscordUser{},
		&persistence.DiscordGuild{},
		&persistence.DiscordChannel{},
	)

	channel := persistence.DiscordChannel{}
	err = orm.
		Preload(clause.Associations).
		First(&channel).
		Error
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("Got channel.guild: %+v\n", channel.DiscordGuild)

	guild := persistence.DiscordGuild{}
	fmt.Printf("\n\n")
	err = orm.
		Preload(clause.Associations).
		First(&guild).
		Error
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("Got guild channels: %+v\n", guild.DiscordChannels)

	// tracking := persistence.Tracking{}
	// err = orm.
	// 	Preload(clause.Associations).
	// 	First(&tracking).Error
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }
	// fmt.Printf("Got tracking: %+v\n", tracking)

	// user, err := tracking.GetSC2ReplayStatsUser(orm)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }
	// fmt.Printf("Got attached user: %+v\n", user)

	// trackings, err := user.GetTrackings(orm)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }
	// fmt.Printf("\n\n\nTrackings are: %+v\n", trackings)
}
