package persistence

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"time"
)

type DiscordUser struct {
	ID            uint `gorm:"primaryKey"`
	DiscordID     string
	Name          string
	Discriminator string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (user *DiscordUser) UpdateFromDgo(author *discordgo.User, orm *gorm.DB) error {
	changed := false
	newUser := DiscordUser{}

	if user.Name != author.Username {
		newUser.Name = author.Username
		changed = true
	}

	if user.Discriminator != author.Discriminator {
		newUser.Discriminator = author.Discriminator
		changed = true
	}

	// .Updates with null-struct still caused an update of the `updated_at` field.
	if changed {
		return orm.
			Model(&user).
			Updates(newUser).
			Error
	}
	return nil
}
