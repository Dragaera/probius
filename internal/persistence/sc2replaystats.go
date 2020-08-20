package persistence

import (
	"fmt"
	sc2r "github.com/dragaera/probius/internal/sc2replaystats"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type SC2ReplayStatsUser struct {
	ID            uint `gorm:"primary_key"`
	DiscordUserID uint
	APIKey        string
	LastReplayID  int
	CreatedAt     time.Time
	UpdatedAt     time.Time
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

func (user *SC2ReplayStatsUser) UpdateLastReplay(db *pgxpool.Pool) (sc2r.Replay, bool, error) {
	replay, err := user.FetchLastReplay()
	// TODO
	return replay, true, err
}
