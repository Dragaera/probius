package persistence

import (
	"context"
	"fmt"
	sc2r "github.com/dragaera/probius/internal/sc2replaystats"
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

func GetSC2ReplayStatsUser(db *pgxpool.Pool, id int) (SC2ReplayStatsUser, error) {
	u := SC2ReplayStatsUser{}

	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_user_id, api_key, last_replay_id, created_at FROM sc2replaystats_users WHERE id=$1",
		id,
	).Scan(&u.ID, &u.DiscordUserID, &u.APIKey, &u.LastReplayID, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one, so GetOrCreate can check for it
		return u, err
	} else if err != nil {
		return u, fmt.Errorf("Unable to get SC2ReplayStats user with ID %v: %v", id, err)
	}

	return u, nil
}

func GetSC2ReplayStatsUserByDiscordUser(db *pgxpool.Pool, discordUser *DiscordUser) (SC2ReplayStatsUser, error) {
	u := SC2ReplayStatsUser{}

	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_user_id, api_key, last_replay_id, created_at FROM sc2replaystats_users WHERE discord_user_id=$1",
		discordUser.ID,
	).Scan(&u.ID, &u.DiscordUserID, &u.APIKey, &u.LastReplayID, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one, so GetOrCreate can check for it
		return u, err
	} else if err != nil {
		return u, fmt.Errorf("Unable to get SC2ReplayStats user: %v", err)
	}

	return u, nil
}

func GetOrCreateSC2ReplayStatsUser(db *pgxpool.Pool, discordUser *DiscordUser, apiKey string) (SC2ReplayStatsUser, error) {
	u, err := GetSC2ReplayStatsUserByDiscordUser(db, discordUser)

	if err == pgx.ErrNoRows {
		err = CreateSC2ReplayStatsUser(db, discordUser, apiKey)
	}

	// Either the getting failed (with an error other than ErrNoRows), or creation failed
	if err != nil {
		return u, fmt.Errorf("Unable to get/create SC2ReplayStats user: %v", err)
	}

	return GetSC2ReplayStatsUserByDiscordUser(db, discordUser)
}

func CreateSC2ReplayStatsUser(db *pgxpool.Pool, discordUser *DiscordUser, apiKey string) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO sc2replaystats_users (discord_user_id, api_key) VALUES ($1, $2)",
		discordUser.ID,
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

func (user *SC2ReplayStatsUser) API() sc2r.API {
	return sc2r.API{APIKey: user.APIKey}
}

func (user *SC2ReplayStatsUser) FetchLastReplay() (sc2r.Replay, error) {
	api := user.API()
	replay, err := api.LastReplay()
	if err != nil {
		return replay, fmt.Errorf("Unable to retrieve last replay: %v", err)
	}

	return replay, nil
}

func (user *SC2ReplayStatsUser) UpdateLastReplay(db *pgxpool.Pool) (sc2r.Replay, bool, error) {
	replay, err := user.FetchLastReplay()
	if err != nil {
		return replay, false, err
	}

	if user.LastReplayID == replay.ReplayID {
		return replay, false, nil
	}

	_, err = db.Exec(
		context.Background(),
		"UPDATE sc2replaystats_users SET last_replay_id = $1 WHERE id = $2",
		replay.ReplayID,
		user.ID,
	)
	if err != nil {
		return replay, false, fmt.Errorf("Unable to update last replay ID: %v", err)
	}

	return replay, true, nil
}
