package persistence

import (
	"fmt"
	"github.com/dragaera/probius/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

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
		postgres.Open(cfg.DBURL()),
		&gorm.Config{
			Logger: logger,
		},
	)
	if err != nil {
		return db, fmt.Errorf("Unable to connect to database:", err)
	}

	return db, nil
}
