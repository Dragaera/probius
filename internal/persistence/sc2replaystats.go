package persistence

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type SC2ReplayStatsUser struct {
	ID            int
	DiscordUserID int
	APIKey        string
	LastReplayID  int
	CreatedAt     time.Time
}

func GetSC2ReplayStatsUser(db *pgxpool.Pool, discordId string) (SC2ReplayStatsUser, error) {
	u := SC2ReplayStatsUser{}

	du, err := GetOrCreateDiscordUser(db, discordId)
	if err != nil {
		return u, fmt.Errorf("Unable to get SC2Replaystats user: %v", err)
	}

	err = db.QueryRow(
		context.Background(),
		"SELECT id, discord_user_id, api_key, last_replay_id, created_at FROM sc2replaystats_users WHERE discord_user_id=$1",
		du.ID,
	).Scan(&u.ID, &u.DiscordUserID, &u.APIKey, &u.LastReplayID, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one, so GetOrCreate can check for it
		return u, err
	} else if err != nil {
		return u, fmt.Errorf("Unable to get SC2ReplayStats user: %v", err)
	}

	return u, nil
}

func GetOrCreateSC2ReplayStatsUser(db *pgxpool.Pool, discordId string, apiKey string) (SC2ReplayStatsUser, error) {
	u, err := GetSC2ReplayStatsUser(db, discordId)

	if err == pgx.ErrNoRows {
		err = CreateSC2ReplayStatsUser(db, discordId, apiKey)
	}

	// Either the getting failed (with an error other than ErrNoRows), or creation failed
	if err != nil {
		return u, fmt.Errorf("Unable to get/create SC2ReplayStats user: %v", err)
	}

	return GetSC2ReplayStatsUser(db, discordId)
}

func CreateSC2ReplayStatsUser(db *pgxpool.Pool, discordId string, apiKey string) error {
	du, err := GetOrCreateDiscordUser(db, discordId)
	if err != nil {
		return fmt.Errorf("Unable to create SC2Replaystats user: %v", err)
	}

	_, err = db.Exec(
		context.Background(),
		"INSERT INTO sc2replaystats_users (discord_user_id, api_key) VALUES ($1, $2)",
		du.ID,
		apiKey,
	)

	if err != nil {
		return fmt.Errorf("Unable to create SC2Replaystats user: %v", err)
	}
	return nil
}

func (user *SC2ReplayStatsUser) UpdateAPIKey(db *pgxpool.Pool, apiKey string) error {
	_, err := db.Exec(
		context.Background(),
		"UPDATE sc2replaystats_users SET api_key = $1 WHERE id = $2",
		apiKey,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("Unable to update API key: %v", err)
	}
	return nil
}
