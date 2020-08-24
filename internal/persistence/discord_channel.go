package persistence

import (
	"time"
)

type DiscordChannel struct {
	ID             uint         `gorm:"primary_key"`
	DiscordGuildID uint         `gorm:"not null"`
	DiscordGuild   DiscordGuild `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DiscordID      string
	Name           string
	IsDM           bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
