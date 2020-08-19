package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/gorm"
	// TODO: This wraps lib/pg which is in maintenance mode. We should use
	// jack/pgx instead.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func InitializeDB(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return dbpool, fmt.Errorf("Unable to connect to database:", err)
	}

	return dbpool, nil
}

func InitializeORM(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return db, fmt.Errorf("Unable to connect to database:", err)
	}

	return db, nil
}
