package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DB      DBConfig
	Discord DiscordConfig
	Redis   RedisConfig
	Worker  WorkerConfig
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
	SSLMode  string
	LogSQL   bool
}

type RedisConfig struct {
	Host string
	Port string
}

type WorkerConfig struct {
	Concurrency int
	Namespace   string
}

func (cfg *DBConfig) DBURL() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Database,
		cfg.Password,
		cfg.SSLMode,
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
	dbCfg.SSLMode = fromEnvWithDefault("DB_SSL_MODE", "disable")
	dbCfg.LogSQL = boolFromEnvWithDefault("DB_LOG_SQL", false)

	redisCfg := RedisConfig{}
	redisCfg.Host = fromEnvWithDefault("REDIS_HOST", "127.0.0.1")
	redisCfg.Port = fromEnvWithDefault("REDIS_PORT", "6379")

	workerCfg := WorkerConfig{}
	workerCfg.Concurrency = intFromEnvWithDefault("WORKER_CONCURRENCY", 5)
	workerCfg.Namespace = fromEnvWithDefault("WORKER_NAMESPACE", "probius")

	cfg := Config{
		DB:      dbCfg,
		Discord: discordCfg,
		Redis:   redisCfg,
		Worker:  workerCfg,
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

func intFromEnv(key string) int {
	str := fromEnv(key)

	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Unable to convert value of %v to int: %v", key, err)
	}

	return val
}

func intFromEnvWithDefault(key string, fallback int) int {
	str := fromEnvWithDefault(key, strconv.FormatInt(int64(fallback), 10))

	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Unable to convert value of %v to int: %v", key, err)
	}

	return val
}

func boolFromEnv(key string) bool {
	str := fromEnv(key)

	switch str {
	case "true", "t", "yes", "y":
		return true
	case "false", "f", "no", "n":
		return false
	default:
		log.Fatalf("Unable to convert value of %v to bool: %v", key, str)
		panic("This piece of code is unreachable")
	}
}

func boolFromEnvWithDefault(key string, fallback bool) bool {
	if _, ok := os.LookupEnv(key); !ok {
		return fallback
	}

	return boolFromEnv(key)
}
