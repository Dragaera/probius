package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitializeDB(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return dbpool, fmt.Errorf("Unable to connect to database:", err)
	}

	return dbpool, nil
}
