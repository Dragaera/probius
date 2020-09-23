package sc2replay

import (
	"encoding/json"
	"fmt"
	"github.com/dragaera/probius/internal/sc2replay/events"
	"github.com/icza/s2prot/rep"
	"math"
)

type Replay struct {
	Rep *rep.Rep
}

func FromFile(path string) (Replay, error) {
	var replay Replay
	rep, err := rep.NewFromFile(path)
	if err != nil {
		return replay, fmt.Errorf("Failed to open replay file: %v", err)
	}
	replay.Rep = rep

	return replay, nil
}

func (replay *Replay) Close() error {
	return replay.Rep.Close()
}

func (replay *Replay) TicksPerSecond() (float64, error) {
	switch replay.Rep.Details.GameSpeed() {
	case rep.GameSpeedSlower:
		// 16 * 0.6
		return 9.6, nil
	case rep.GameSpeedSlow:
		// 16 * 0.8
		return 12.8, nil
	case rep.GameSpeedNormal:
		return 16, nil
	case rep.GameSpeedFast:
		// 16 * 1.2
		return 19.2, nil
	case rep.GameSpeedFaster:
		// 16 * 1.4
		return 22.4, nil
	default:
		return 0, fmt.Errorf("Gamespeed of replay is unknown")
	}
}

func (replay *Replay) TicksUntilSeconds(seconds float64) (int64, error) {
	ticksPerSecond, err := replay.TicksPerSecond()
	if err != nil {
		return 0, err
	}

	return int64(math.Round(ticksPerSecond * seconds)), nil
}

// Return the *User ID* of the replay's owner.
//
// Mind that this is NOT the Player ID, but rather a separate identifier.
func (replay *Replay) OwnerID() (int64, error) {
	var userLeave = events.GameUserLeave{}
	userLeavesFound := 0

	// The way it works: The *last* GameUserLeave event is the owner of the
	// replay. AI does not cause any such events, so in a vs AI game there
	// will be only one.
	for _, evt := range replay.Rep.GameEvts {
		switch eventType := evt.EvtType.Name; eventType {
		case "GameUserLeave":
			err := json.Unmarshal([]byte(evt.String()), &userLeave)
			if err != nil {
				return 0, fmt.Errorf("Unable to parse GameUserLeave: %v", err)
			}
			userLeavesFound += 1
		}
	}

	if userLeavesFound < 1 {
		return 0, fmt.Errorf("No GameUserLeave events found, replay might be corrupt")
	}

	return userLeave.UserID.UserID, nil
}

// Return the player ID corresponding to the *non-AI* player, running as the replay's owner.
func (replay *Replay) OwnerPlayerID() (int64, error) {
	userID, err := replay.OwnerID()
	if err != nil {
		return 0, err
	}

	// If there are AI opponents, there might be multiple players with the
	// same user ID, as AI opponents are run by a human user.
	// In the case of team-games with AI I assume the host is the one
	// running the AIs, but I did not verify that.

	possiblePlayers := make([]*(rep.PlayerDesc), 0)
	for _, playerDesc := range replay.Rep.TrackerEvts.PIDPlayerDescMap {
		if playerDesc.UserID == userID {
			possiblePlayers = append(possiblePlayers, playerDesc)
		}
	}

	switch len(possiblePlayers) {
	case 0:
		return 0, fmt.Errorf("No player with User ID %d found", userID)
	case 1:
		return possiblePlayers[0].PlayerID, nil
	}

	// We're left with multiple possible players owned by the same user.
	// We'll now narrow it down to only *human* players.

	fmt.Println("Multiple possible replay owner player IDs detected, checking which players are human.")
	possibleHumanPlayers := make([]*(rep.PlayerDesc), 0)
	for _, desc := range possiblePlayers {
		if int(desc.SlotID) > len(replay.Rep.InitData.LobbyState.Slots)-1 {
			// Shouldn't ever happen, as slots 0-15 are always populated (even if empty), but who knows...
			fmt.Printf("Player %d has no slot assigned, skipping\n", desc.PlayerID)
			continue
		}
		slot := replay.Rep.InitData.LobbyState.Slots[desc.SlotID]

		if slot.Control() == rep.ControlHuman {
			possibleHumanPlayers = append(possibleHumanPlayers, desc)
		}
	}

	switch len(possibleHumanPlayers) {
	case 0:
		return 0, fmt.Errorf("No possible human replay owners detected")
	case 1:
		return possibleHumanPlayers[0].PlayerID, nil
	default:
		return 0, fmt.Errorf("%d possible human replay owners detected", len(possibleHumanPlayers))
	}
}
