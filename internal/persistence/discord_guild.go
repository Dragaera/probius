package persistence

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"time"
)

type DiscordGuild struct {
	ID        uint `gorm:"primaryKey"`
	DiscordID string
	Name      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetDMGuild(orm *gorm.DB) (DiscordGuild, error) {
	guild := DiscordGuild{}

	err := orm.
		Where(DiscordGuild{DiscordID: "0"}).
		Attrs(DiscordGuild{Name: "Direct Message", OwnerID: "0"}).
		FirstOrCreate(&guild).
		Error
	if err != nil {
		return guild, fmt.Errorf("Unable to retrieve DM guild from DB: %v", err)
	}

	return guild, nil
}

func (guild *DiscordGuild) UpdateFromDgo(dgoGuild *discordgo.Guild, orm *gorm.DB) error {
	changed := false
	newUser := DiscordGuild{}

	if guild.Name != dgoGuild.Name {
		newUser.Name = dgoGuild.Name
		changed = true
	}

	if guild.OwnerID != dgoGuild.OwnerID {
		newUser.OwnerID = dgoGuild.OwnerID
		changed = true
	}

	// .Updates with null-struct still caused an update of the `updated_at` field.
	if changed {
		return orm.
			Model(&guild).
			Updates(newUser).
			Error
	}
	return nil
}
