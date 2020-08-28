package persistence

import (
	"fmt"
	sc2r "github.com/dragaera/probius/internal/sc2replaystats"
	"gorm.io/gorm"
	"time"
)

type SC2ReplayStatsUser struct {
	ID            uint        `gorm:"primaryKey"`
	DiscordUserID uint        `gorm:"not null"`
	DiscordUser   DiscordUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	APIKey        string
	LastReplayID  int
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
