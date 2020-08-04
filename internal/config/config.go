package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	DB      DBConfig
	Discord DiscordConfig
}

type DiscordConfig struct {
	ClientID string
	Token    string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func (cfg *DBConfig) DBURL() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func ConfigFromEnv() Config {
	log.Print("Loading configuration from environment")

	discordCfg := DiscordConfig{}
	discordCfg.ClientID = fromEnv("DISCORD_CLIENT_ID")
	discordCfg.Token = fromEnv("DISCORD_TOKEN")

	dbCfg := DBConfig{}
	dbCfg.User = fromEnv("DB_USER")
	dbCfg.Password = fromEnv("DB_PASSWORD")
	dbCfg.Host = fromEnvWithDefault("DB_HOST", "127.0.0.1")
	dbCfg.Port = fromEnvWithDefault("DB_PORT", "5432")
	dbCfg.Database = fromEnvWithDefault("DB_DATABASE", "probius")

	cfg := Config{
		DB:      dbCfg,
		Discord: discordCfg,
	}

	return cfg
}

func fromEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Missing env variable: %v", key)
	}

	return val
}

func fromEnvWithDefault(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Using default value for key: %v = %v", key, fallback)
		val = fallback
	}

	return val
}
