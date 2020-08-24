package persistence

import (
	"context"
	"fmt"
	"github.com/dragaera/probius/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitializeDB(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return dbpool, fmt.Errorf("Unable to connect to database:", err)
	}

	return dbpool, nil
}

func InitializeORM(cfg config.DBConfig) (*gorm.DB, error) {
	var logMode logger.LogLevel
	if cfg.LogSQL {
		logMode = logger.Info
	} else {
		logMode = logger.Silent
	}

	logger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logMode,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(
		postgres.Open(cfg.DBURL2()),
		&gorm.Config{
			Logger: logger,
		},
	)
	if err != nil {
		return db, fmt.Errorf("Unable to connect to database:", err)
	}

	return db, nil
}
