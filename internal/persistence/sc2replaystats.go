package persistence

import (
	"fmt"
	sc2r "github.com/dragaera/probius/internal/sc2replaystats"
	"gorm.io/gorm"
	"time"
)

type SC2ReplayStatsUser struct {
	ID            uint `gorm:"primary_key"`
	DiscordUserID uint
	APIKey        string
	LastReplayID  int
	Trackings     []Tracking
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (user *SC2ReplayStatsUser) GetTrackings(db *gorm.DB) ([]Tracking, error) {
	trackings := make([]Tracking, 10)
	err := db.Model(&user).Association("Trackings").Find(&trackings)

	return trackings, err
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
