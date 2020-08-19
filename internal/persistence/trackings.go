package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Tracking struct {
	ID                   int
	DiscordChannelID     int
	SC2ReplayStatsUserID int
	CreatedAt            time.Time
}

func CreateTracking(db *pgxpool.Pool, channel *DiscordChannel, user *SC2ReplayStatsUser) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO trackings (discord_channel_id, sc2replaystats_user_id) VALUES ($1, $2)",
		channel.ID,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("Unable to create tracking: %v", err)
	}
	return nil
}

func GetTracking(db *pgxpool.Pool, channel *DiscordChannel, user *SC2ReplayStatsUser) (Tracking, error) {
	t := Tracking{}

	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_channel_id, sc2replaystats_user_id, created_at FROM trackings WHERE discord_channel_id = $1 AND sc2replaystats_user_id = $2",
		channel.ID,
		user.ID,
	).Scan(&t.ID, &t.DiscordChannelID, &t.SC2ReplayStatsUserID, &t.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one, so GetOrCreate can check for it
		return t, err
	} else if err != nil {
		return t, fmt.Errorf("Unable to get tracking with user ID %v, channel ID: %v: %v", user.ID, channel.ID, err)
	}

	return t, nil
}

func DeleteTracking(db *pgxpool.Pool, id int) error {
	_, err := db.Exec(
		context.Background(),
		"DELETE FROM trackings WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("Unable to delete tracking: %v", err)
	}
	return nil
}
