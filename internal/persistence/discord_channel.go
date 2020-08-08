package persistence

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type DiscordChannel struct {
	ID        int
	DiscordID string
	Name      string
	IsDM      bool
	CreatedAt time.Time
}

func GetDiscordChannel(db *pgxpool.Pool, discordID string) (DiscordChannel, error) {
	c := DiscordChannel{}
	err := db.QueryRow(
		context.Background(),
		"SELECT id, discord_id, name, is_dm, created_at FROM discord_channels WHERE discord_id=$1",
		discordID,
	).Scan(&c.ID, &c.DiscordID, &c.Name, &c.IsDM, &c.CreatedAt)

	if err == pgx.ErrNoRows {
		// We preserve this one
		return c, err
	} else if err != nil {
		return c, fmt.Errorf("Unable to get Discord channel: %v", err)
	}

	return c, nil
}

func CreateDiscordChannel(db *pgxpool.Pool, discordID string, guildID int, name string, isDM bool) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO discord_channels (discord_id, guild_id, name, is_dm) VALUES ($1, $2, $3, $4)",
		discordID,
		guildID,
		name,
		isDM,
	)

	if err != nil {
		return fmt.Errorf("Unable to create discord channel: %v", err)
	}
	return nil
}

// TODO: Update channel if its details changed
func DiscordChannelFromDgo(db *pgxpool.Pool, data *discordgo.Channel) (DiscordChannel, error) {
	c, err := GetDiscordChannel(db, data.ID)
	if err == nil {
		return c, nil
	} else if err == pgx.ErrNoRows {
		var guild DiscordGuild
		if len(data.GuildID) == 0 {
			// DM guild
			guild, err = GetDiscordGuild(db, "0")
			if err != nil {
				return c, fmt.Errorf("Unable to retrieve DM guild from DB: %v", err)
			}
		} else {
			guild, err = GetDiscordGuild(db, data.GuildID)
			if err != nil {
				return c, fmt.Errorf("Unable to retrieve guild of channel: %v", err)
			}
		}

		err = CreateDiscordChannel(
			db,
			data.ID,
			guild.ID,
			data.Name,
			data.Type == discordgo.ChannelTypeDM,
		)
	}

	if err != nil {
		// Either getting failed (with an error other than ErrNoRows), or creation failed.
		return c, fmt.Errorf("Unable to get/create Discord channel: %v", err)
	}

	return GetDiscordChannel(db, data.ID)
}
