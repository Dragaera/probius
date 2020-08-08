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
	DiscordID     string
	Name          string
	Discriminator string
	CreatedAt     time.Time
}

func GetDiscordUser(db *pgxpool.Pool, discordID string) (DiscordUser, error) {
	u := DiscordUser{}
	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_id, name, discriminator, created_at FROM discord_users WHERE discord_id=$1",
		discordID,
	).Scan(&u.ID, &u.DiscordID, &u.Name, &u.Discriminator, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one
		return u, err
	} else if err != nil {
		return u, fmt.Errorf("Unable to get Discord user: %v", err)
	}

	return u, nil
}

func CreateDiscordUser(db *pgxpool.Pool, discordID string, name string, discriminator string) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO discord_users (discord_id, name, discriminator) VALUES ($1, $2, $3)",
		discordID,
		name,
		discriminator,
	)

	if err != nil {
		return fmt.Errorf("Unable to create discord user: %v", err)
	}
	return nil
}

// TODO: Update user if its details changed
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
