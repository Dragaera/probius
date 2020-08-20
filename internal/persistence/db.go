package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return dbpool, fmt.Errorf("Unable to connect to database:", err)
	}

	return dbpool, nil
}

func InitializeORM(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("Unable to connect to database:", err)
	}

	return db, nil
}
