package persistence

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"time"
)

type DiscordChannel struct {
	ID             uint         `gorm:"primaryKey"`
	DiscordGuildID uint         `gorm:"not null"`
	DiscordGuild   DiscordGuild `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DiscordID      string
	Name           string
	IsDM           bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (channel *DiscordChannel) UpdateFromDgo(dgoChannel *discordgo.Channel, orm *gorm.DB) error {
	changed := false
	newChannel := DiscordChannel{}

	if channel.Name != dgoChannel.Name {
		newChannel.Name = dgoChannel.Name
		changed = true
	}

	// .Updates with null-struct still caused an update of the `updated_at` field.
	if changed {
		return orm.
			Model(&channel).
			Updates(newChannel).
			Error
	}
	return nil
}
