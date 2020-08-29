package persistence

import (
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	ID                   uint               `gorm:"primaryKey"`
	DiscordChannelID     uint               `gorm:"not null"`
	DiscordChannel       DiscordChannel     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SC2ReplayStatsUserID uint               `gorm:"not null"`
	SC2ReplayStatsUser   SC2ReplayStatsUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (subscription *Subscription) GetSC2ReplayStatsUser(db *gorm.DB) (SC2ReplayStatsUser, error) {
	user := SC2ReplayStatsUser{}
	err := db.First(
		&user,
		"id = ?",
		subscription.SC2ReplayStatsUserID,
	).Error

	return user, err
}
