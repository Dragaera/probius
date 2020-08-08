package persistence

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type DiscordUser struct {
	ID            int
	DiscordID     uint64
	Name          string
	Discriminator string
	CreatedAt     time.Time
}

func GetDiscordUser(db *pgxpool.Pool, discordId string) (DiscordUser, error) {
	u := DiscordUser{}
	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_id, name, discriminator, created_at FROM discord_users WHERE discord_id=$1",
		discordId,
	).Scan(&u.ID, &u.DiscordID, &u.Name, &u.Discriminator, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one, so GetOrCreate can check for it
		return u, err
	} else if err != nil {
		return u, fmt.Errorf("Unable to get Discord user: %v", err)
	}

	return u, nil
}

// func GetOrCreateDiscordUser(db *pgxpool.Pool, discordId string) (DiscordUser, error) {
// 	u, err := GetDiscordUser(db, discordId)
//
// 	if err == pgx.ErrNoRows {
// 		err = CreateDiscordUser(db, discordId)
// 	}
//
// 	// Either the getting failed (with an error other than ErrNoRows), or creation failed
// 	if err != nil {
// 		return u, fmt.Errorf("Unable to get/create Discord user: %v", err)
// 	}
//
// 	return GetDiscordUser(db, discordId)
// }

func CreateDiscordUser(db *pgxpool.Pool, discordId string, name string, discriminator string) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO discord_users (discord_id, name, discriminator) VALUES ($1, $2, $3)",
		discordId,
		name,
		discriminator,
	)

	if err != nil {
		return err
	}
	return nil
}

func DiscordUserFromDgo(db *pgxpool.Pool, data *discordgo.User) (DiscordUser, error) {
	u, err := GetDiscordUser(db, data.ID)
	if err == nil {
		return u, nil
	} else if err == pgx.ErrNoRows {
		err = CreateDiscordUser(
			db,
			data.ID,
			data.Username,
			data.Discriminator,
		)
	}

	if err != nil {
		// Either getting failed (with an error other than ErrNoRows), or creation failed.
		return u, fmt.Errorf("Unable to get/create Discord user: %v", err)
	}

	return GetDiscordUser(db, data.ID)
}
