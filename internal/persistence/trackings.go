package persistence

import (
	"gorm.io/gorm"
	"time"
)

type Tracking struct {
	ID                   uint `gorm:"primary_key"`
	DiscordChannelID     uint `gorm:"not null"`
	SC2ReplayStatsUserID uint `gorm:"not null"`
	SC2ReplayStatsUser   SC2ReplayStatsUser
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (tracking *Tracking) GetSC2ReplayStatsUser(db *gorm.DB) (SC2ReplayStatsUser, error) {
	user := SC2ReplayStatsUser{}
	err := db.First(
		&user,
		"id = ?",
		tracking.SC2ReplayStatsUserID,
	).Error

	return user, err
}
