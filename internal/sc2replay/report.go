package sc2replay

import (
	"encoding/json"
	"fmt"
	"github.com/dragaera/probius/internal/sc2replay/events"
	"github.com/dragaera/probius/internal/sc2replay/units"
	"github.com/icza/s2prot"
	"math"
)

type IngameUnit struct {
	Index   int64
	Recycle int64
	Name    string
	// The ID of the one towards whose *upkeep* it counts. Ie we don't care
	// about neuralled units etc.
	OwnerID int64
}

type Report struct {
	PlayerID int64
	Replay   *Replay
	// Map containing everything the game considers a unit. This also
	// includes buildings, mineral patches etc. This also contains units
	// owned by other players.
	IngameUnits map[int64]IngameUnit
	// Map containing enriched units belonging to the specified player ID,
	// and only those for which specific information is available.
	Units map[int64]units.Unit

	// Count of units by ingame name
	UnitCount map[string]int

	// Supply. As there are units with 0.5 supply, this is a float. Use
	// `Report.IngameSupply()` for the integer (rounded) supply as shown in-game.
	Supply float64
}

// Call this to generate the report.
func (rep *Report) At(ticks int64) {
	rep.IngameUnits = make(map[int64]IngameUnit)

	for _, evt := range rep.Replay.Rep.TrackerEvts.Evts {
		// We handle this here to allow an early exit
		if evt.Loop() > int64(ticks) {
			fmt.Printf("Reached %d ticks, stopping\n", ticks)
			break
		}
		if err := rep.handleEvent(evt); err != nil {
			fmt.Printf("Error while handling event: %v\n", err)
			fmt.Printf("%+v\n", evt)
		}
	}

	// Remove units belonging to other players
	rep.prune()
	// Enrich with static information (supply, human-readable name etc)
	rep.enrich()

	rep.calculateUnitCount()
	rep.calculateSupply()
}

// Return rounded supply as shown in-game
func (rep *Report) IngameSupply() int {
	return int(math.Round(rep.Supply))
}

func (rep *Report) handleEvent(evt s2prot.Event) error {
	switch eventType := evt.EvtType.Name; eventType {
	case "UnitBorn":
		if err := rep.trackUnitBorn(evt); err != nil {
			return err
		}
	case "UnitInit":
		if err := rep.trackUnitInit(evt); err != nil {
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
		if err := rep.trackUnitTypeChange(evt); err != nil {
			return err
		}
	case "UnitDied":
		if err := rep.trackUnitDied(evt); err != nil {
			return err
		}
	default:
		// fmt.Printf("[%d]: %s by %d\n", evt.Loop(), eventType, evt.UserID())
	}

	return nil
}

func (rep *Report) trackUnitBorn(evt s2prot.Event) error {
	event := events.UnitBorn{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitBorn event: %v", err)
	}

	if err := rep.addUnit(event.UnitTagIndex, event.UnitTagRecycle, event.UnitTypeName, event.UpkeepPlayerID); err != nil {
		return err
	}

	return nil
}

func (rep *Report) trackUnitInit(evt s2prot.Event) error {
	event := events.UnitInit{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitInit event: %v", err)
	}

	if err := rep.addUnit(event.UnitTagIndex, event.UnitTagRecycle, event.UnitTypeName, event.UpkeepPlayerID); err != nil {
		return err
	}

	return nil
}

func (rep *Report) trackUnitTypeChange(evt s2prot.Event) error {
	event := events.UnitTypeChange{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitTypeChange event: %v", err)
	}

	if err := rep.replaceUnit(event.UnitTagIndex, event.UnitTagRecycle, event.UnitTypeName); err != nil {
		return err
	}

	return nil
}

func (rep *Report) trackUnitDied(evt s2prot.Event) error {
	event := events.UnitDied{}
	if err := json.Unmarshal([]byte(evt.String()), &event); err != nil {
		return fmt.Errorf("Unable to unmarshal UnitDied event: %v", err)
	}

	if err := rep.removeUnit(event.UnitTagIndex, event.UnitTagRecycle); err != nil {
		return err
	}

	return nil
}

func (rep *Report) addUnit(index int64, recycle int64, name string, ownerID int64) error {
	tag := unitTag(index, recycle)

	if existing, ok := rep.IngameUnits[tag]; ok {
		// Unit with given tag exists already => That's a mistake
		return fmt.Errorf("Unit tag %d reused. Existing: %s, new: %s", tag, existing, name)
	}
	rep.IngameUnits[tag] = IngameUnit{Index: index, Recycle: recycle, Name: name, OwnerID: ownerID}

	return nil
}

func (rep *Report) replaceUnit(index int64, recycle int64, name string) error {
	tag := unitTag(index, recycle)

	if _, ok := rep.IngameUnits[tag]; !ok {
		// Trying to replace a nonexistant unit
		return fmt.Errorf("Tried to replace unit with tag %d but does not exist", tag)
	}
	// Cannot change struct fields in maps
	existing := rep.IngameUnits[tag]
	existing.Name = name
	rep.IngameUnits[tag] = existing

	return nil
}

func (rep *Report) removeUnit(index int64, recycle int64) error {
	tag := unitTag(index, recycle)

	if _, ok := rep.IngameUnits[tag]; !ok {
		// Trying to remove a nonexistant unit
		return fmt.Errorf("Tried to remove unit tag %d but does not exist", tag)
	}
	delete(rep.IngameUnits, tag)

	return nil
}

// Remove all units owned (in terms of supply) other than rep.PlayerID
func (rep *Report) prune() {
	for tag, unit := range rep.IngameUnits {
		if unit.OwnerID != rep.PlayerID {
			delete(rep.IngameUnits, tag)
		}
	}
}

// Enrich with static per-unit information such as supply and name.
//
// Enriched information will be in rep.Units, which will only contain those
// units for which enriched information is available.
func (rep *Report) enrich() {
	rep.Units = make(map[int64]units.Unit)

	for tag, unit := range rep.IngameUnits {
		if enrichedUnit, ok := units.Units[unit.Name]; ok {
			rep.Units[tag] = enrichedUnit
		}
	}
}

func (rep *Report) calculateUnitCount() {
	rep.UnitCount = make(map[string]int)

	for _, unit := range rep.Units {
		rep.UnitCount[unit.Name] += 1
	}
}

func (rep *Report) calculateSupply() {
	rep.Supply = 0

	for _, unit := range rep.Units {
		rep.Supply += unit.Supply
	}
}

func unitTag(unitTagIndex int64, unitTagRecycle int64) int64 {
	// Ripped from https://github.com/Blizzard/s2protocol, search `func
	// unit_tag`. Whoever thought of this system must've been drunk.
	return (unitTagIndex << 18) + unitTagRecycle
}
