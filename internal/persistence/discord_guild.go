package persistence

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type DiscordGuild struct {
	ID              uint `gorm:"primary_key"`
	DiscordID       string
	Name            string
	OwnerID         string
	DiscordChannels []DiscordChannel
	CreatedAt       time.Time
	UpdatedAt       time.Time
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
