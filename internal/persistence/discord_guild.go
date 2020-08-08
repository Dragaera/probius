package persistence

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type DiscordGuild struct {
	ID        int
	DiscordID string
	Name      string
	OwnerID   string
	CreatedAt time.Time
}

func GetDiscordGuild(db *pgxpool.Pool, discordID string) (DiscordGuild, error) {
	g := DiscordGuild{}
	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_id, name, owner_id, created_at FROM discord_guilds WHERE discord_id=$1",
		discordID,
	).Scan(&g.ID, &g.DiscordID, &g.Name, &g.OwnerID, &g.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one
		return g, err
	} else if err != nil {
		return g, fmt.Errorf("Unable to get Discord guild: %v", err)
	}

	return g, nil
}

func CreateDiscordGuild(db *pgxpool.Pool, discordID string, name string, ownerID string) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO discord_guilds (discord_id, name, owner_id) VALUES ($1, $2, $3)",
		discordID,
		name,
		ownerID,
	)

	return err
	if err != nil {
		return fmt.Errorf("Unable to create discord guild: %v", err)
	}
	return nil
}

// TODO: Update guild if its details changed
func DiscordGuildFromDgo(db *pgxpool.Pool, data *discordgo.Guild) (DiscordGuild, error) {
	g, err := GetDiscordGuild(db, data.ID)
	if err == nil {
		return g, nil
	} else if err == pgx.ErrNoRows {
		err = CreateDiscordGuild(
			db,
			data.ID,
			data.Name,
			data.OwnerID,
		)
	}

	if err != nil {
		// Either getting failed (with an error other than ErrNoRows), or creation failed.
		return g, fmt.Errorf("Unable to get/create Discord guild: %v", err)
	}

	return GetDiscordGuild(db, data.ID)
}
