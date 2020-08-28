package persistence

import (
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
