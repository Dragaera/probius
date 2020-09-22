package sc2replay

import (
	"encoding/json"
	"fmt"
	"github.com/dragaera/probius/internal/sc2replay/events"
	"github.com/icza/s2prot"
	"github.com/icza/s2prot/rep"
	"math"
)

type Replay struct {
	rep *rep.Rep
}

func FromFile(path string) (Replay, error) {
	var replay Replay
	rep, err := rep.NewFromFile(path)
	if err != nil {
		return replay, fmt.Errorf("Failed to open replay file: %v", err)
	}
	replay.rep = rep

	return replay, nil
}

func (replay *Replay) Dump() {
	for _, player := range replay.rep.Details.Players() {
		fmt.Println(player)
	}
}

func (replay *Replay) Close() error {
	return replay.rep.Close()
}

func (replay *Replay) TicksPerSecond() (float64, error) {
	switch replay.rep.Details.GameSpeed() {
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
	for _, evt := range replay.rep.GameEvts {
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

func (replay *Replay) OwnerPlayerID() (int64, error) {
	userID, err := replay.OwnerID()
	fmt.Println(userID)
	if err != nil {
		return 0, err
	}

	for playerID, playerDesc := range replay.rep.TrackerEvts.PIDPlayerDescMap {
		if playerDesc.UserID == userID {
			return playerID, nil
		}
	}

	return 0, fmt.Errorf("No player with User ID %d found", userID)
}

func (replay *Replay) UnitsAt(ticks int64, playerID int64) map[string]int {
	fmt.Printf("Calculating units at %d ticks by player %d\n", ticks, playerID)

	stats := UnitStats{PlayerID: playerID, Ticks: ticks}
	stats.Units = make(map[int64]IngameUnit)

	for _, evt := range replay.rep.TrackerEvts.Evts {
		// We handle this here to allow an early exit
		if evt.Loop() > ticks {
			fmt.Printf("Reached %d ticks, stopping\n", ticks)
			break
		}
		if err := stats.handleEvent(evt); err != nil {
			fmt.Printf("Error while handling event: %v\n", err)
			fmt.Printf("%+v\n", evt)
		}
	}

	stats.prune()
	return stats.summarize()
}

type UnitStats struct {
	PlayerID int64
	Ticks    int64
	Units    map[int64]IngameUnit
}

type IngameUnit struct {
	Index   int64
	Recycle int64
	Name    string
	// The ID of the one towards whose *upkeep* it counts. Ie we don't care
	// about neuralled units etc.
	OwnerID int64
}

func (stats *UnitStats) handleEvent(evt s2prot.Event) error {
	if evt.Loop() > stats.Ticks {
		return nil
	}

	switch eventType := evt.EvtType.Name; eventType {
	case "UnitBorn":
		if err := stats.trackUnitBorn(evt); err != nil {
			return err
		}
	case "UnitInit":
		if err := stats.trackUnitInit(evt); err != nil {
			return err
		}
	case "UnitDone":
		// UnitDone is for eg:
		// - A unit finishing warpin
		// - A building finishing morphing
		// As we already add units to the unit list when they start
		// building, we have no need for this.
	case "UnitTypeChange":
		// UnitTypeChange is for eg:
		// - Buildings transforming (eg gateway => warpgate)
		// - Creep tumors burrowing
		// - Larva transforming into eggs
		// - Hellions transforming into Hellbats
		if err := stats.trackUnitTypeChange(evt); err != nil {
			return err
		}
	case "UnitDied":
		if err := stats.trackUnitDied(evt); err != nil {
			return err
		}
	default:
		// fmt.Printf("[%d]: %s by %d\n", evt.Loop(), eventType, evt.UserID())
	}

	return nil
}

func (stats *UnitStats) trackUnitBorn(evt s2prot.Event) error {
	event := events.UnitBorn{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitBorn event: %v", err)
	}

	if err := stats.addUnit(event.UnitTagIndex, event.UnitTagRecycle, event.UnitTypeName, event.UpkeepPlayerID); err != nil {
		return err
	}

	return nil
}

func (stats *UnitStats) trackUnitInit(evt s2prot.Event) error {
	event := events.UnitInit{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitInit event: %v", err)
	}

	if err := stats.addUnit(event.UnitTagIndex, event.UnitTagRecycle, event.UnitTypeName, event.UpkeepPlayerID); err != nil {
		return err
	}

	return nil
}

func (stats *UnitStats) trackUnitTypeChange(evt s2prot.Event) error {
	event := events.UnitTypeChange{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitTypeChange event: %v", err)
	}

	if err := stats.replaceUnit(event.UnitTagIndex, event.UnitTagRecycle, event.UnitTypeName); err != nil {
		return err
	}

	return nil
}

func (stats *UnitStats) trackUnitDied(evt s2prot.Event) error {
	event := events.UnitDied{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitDied event: %v", err)
	}

	if err := stats.removeUnit(event.UnitTagIndex, event.UnitTagRecycle); err != nil {
		return err
	}

	return nil
}

func (stats *UnitStats) addUnit(index int64, recycle int64, name string, ownerID int64) error {
	tag := unitTag(index, recycle)

	if existing, ok := stats.Units[tag]; ok {
		// Unit with given tag exists already => That's a mistake
		return fmt.Errorf("Unit tag %d reused. Existing: %s, new: %s", tag, existing, name)
	}
	stats.Units[tag] = IngameUnit{Index: index, Recycle: recycle, Name: name, OwnerID: ownerID}

	return nil
}

func (stats *UnitStats) replaceUnit(index int64, recycle int64, name string) error {
	tag := unitTag(index, recycle)

	if _, ok := stats.Units[tag]; !ok {
		// Trying to replace a nonexistant unit
		return fmt.Errorf("Tried to replace unit with tag %d but does not exist", tag)
	}
	// Cannot change struct fields in maps
	existing := stats.Units[tag]
	existing.Name = name
	stats.Units[tag] = existing

	return nil
}

func (stats *UnitStats) removeUnit(index int64, recycle int64) error {
	tag := unitTag(index, recycle)

	if _, ok := stats.Units[tag]; !ok {
		// Trying to remove a nonexistant unit
		return fmt.Errorf("Tried to remove unit tag %d but does not exist", tag)
	}
	delete(stats.Units, tag)

	return nil
}

// Remove all units owned (in terms of supply) other than stats.PlayerID
func (stats *UnitStats) prune() {
	for tag, unit := range stats.Units {
		if unit.OwnerID != stats.PlayerID {
			delete(stats.Units, tag)
		}
	}
}

func (stats *UnitStats) summarize() map[string]int {
	out := make(map[string]int)

	for _, unit := range stats.Units {
		out[unit.Name] += 1
	}

	return out
}

func unitTag(unitTagIndex int64, unitTagRecycle int64) int64 {
	// Ripped from https://github.com/Blizzard/s2protocol, search `func
	// unit_tag`. Whoever thought of this system must've been drunk.
	return (unitTagIndex << 18) + unitTagRecycle
}
