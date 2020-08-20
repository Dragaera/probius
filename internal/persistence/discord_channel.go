package persistence

import (
	"time"
)

type DiscordChannel struct {
	ID             uint `gorm:"primary_key"`
	DiscordGuildID uint
	DiscordID      string
	Name           string
	IsDM           bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
