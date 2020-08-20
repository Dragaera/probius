package persistence

import (
	"time"
)

type Tracking struct {
	ID                   uint `gorm:"primary_key"`
	DiscordChannelID     uint `gorm:"not null"`
	SC2ReplayStatsUserID uint `gorm:"not null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
