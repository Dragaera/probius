package persistence

import (
	"fmt"
	sc2r "github.com/dragaera/probius/internal/sc2replaystats"
	"gorm.io/gorm"
	"log"
	"time"
)

type SC2ReplayStatsUser struct {
	ID                uint        `gorm:"primaryKey"`
	DiscordUserID     uint        `gorm:"not null"`
	DiscordUser       DiscordUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	APIKey            string
	LastReplayID      int
	LastCheckedAt     time.Time
	UpdateScheduledAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (user *SC2ReplayStatsUser) GetSubscriptions(db *gorm.DB) ([]Subscription, error) {
	subscriptions := make([]Subscription, 10)
	err := db.
		Where(Subscription{SC2ReplayStatsUserID: user.ID}).
		Preload("DiscordChannel.DiscordGuild").
		Find(&subscriptions).
		Error
	if err != nil {
		err = fmt.Errorf("Unable to retrieve subscriptions for user %v: %v", user.ID, err)
	}

	return subscriptions, err
}

func (user *SC2ReplayStatsUser) API() sc2r.API {
	return sc2r.API{APIKey: user.APIKey}
}

func (user *SC2ReplayStatsUser) FetchLastReplay() (sc2r.Replay, error) {
	api := user.API()
	replay, err := api.LastReplay()
	if err != nil {
		return replay, fmt.Errorf("Unable to retrieve last replay: %v", err)
	}

	return replay, nil
}

func (user *SC2ReplayStatsUser) UpdateLastReplay(orm *gorm.DB) (sc2r.Replay, bool, error) {
	defer user.TouchLastCheckedAt(orm)
	defer user.UnlockForUpdate(orm)

	replay, err := user.FetchLastReplay()

	if err != nil {
		return replay, false, err
	}

	replayChanged := replay.ReplayID != user.LastReplayID
	if !replayChanged {
		return replay, replayChanged, err
	}

	err = orm.Model(&user).Update("last_replay_id", replay.ReplayID).Error
	if err != nil {
		return replay, replayChanged, fmt.Errorf("Unable to update last replay: %v", err)
	}

	return replay, replayChanged, err
}

func (user *SC2ReplayStatsUser) LockForUpdate(orm *gorm.DB) error {
	return orm.
		Model(&user).
		Update("update_scheduled_at", time.Now()).
		Error
}

func (user *SC2ReplayStatsUser) UnlockForUpdate(orm *gorm.DB) error {
	return orm.
		Model(&user).
		Update("update_scheduled_at", nil).
		Error
}

func (user *SC2ReplayStatsUser) TouchLastCheckedAt(orm *gorm.DB) error {
	return orm.
		Model(&user).
		Update("last_checked_at", time.Now()).
		Error
}

func SC2ReplayStatsUsersWithStaleData(orm *gorm.DB, updateInterval int) ([]SC2ReplayStatsUser, error) {
	// Timestamp where players which haven't been updated since then are considered stale
	ts := time.Now().Add(time.Second * time.Duration(-updateInterval))

	users := make([]SC2ReplayStatsUser, 10)
	err := orm.
		Where("(update_scheduled_at is NULL AND last_checked_at <= ?) OR last_checked_at is NULL", ts).
		Find(&users).
		Error
	if err != nil {
		err = fmt.Errorf("Unable to retrieve players with stale data: %v", err)
	}

	return users, err
}

func ClearStaleSC2ReplayStatsUpdateLocks(orm *gorm.DB, ttl int) error {
	ts := time.Now().Add(time.Second * time.Duration(-ttl))
	err := orm.
		Model(&SC2ReplayStatsUser{}).
		Where("update_scheduled_at <= ?", ts).
		Update("update_scheduled_at", nil).
		Error
	if err != nil {
		err = fmt.Errorf("Unable to clear stale locks: %v", err)
		log.Print(">>>>>", err)
	}

	return err
}
