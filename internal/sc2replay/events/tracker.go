package events

import (
	"encoding/json"
	"fmt"
)

// Used only for composition
type BaseEvent struct {
	ID   int `json:"id"`
	Loop int `json:"loop"`
	// EventType EventType `json:"evtTypeName"`
}

// 1 = true, 0 = false. Because we're in 1995
type IntBool bool

func (intBool *IntBool) UnmarshalJSON(b []byte) error {
	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	switch i {
	case 0:
		*intBool = false
	case 1:
		*intBool = true
	default:
		return fmt.Errorf("Invalid error for boolean: %v", i)
	}

	return nil
}

// type EventType struct {
// 	// reflect.Type is an interface, and as such cannot be the receiver of
// 	// any functions. So to define UnmarshalJSON on it, we'll have it as
// 	// member of a struct.
// 	// Embedding might have bee a more natural choice, but I couldn't get
// 	// it working right away.
// 	Type reflect.Type
// }
//
// func (typ *EventType) UnmarshalJSON(b []byte) error {
// 	var eventTypeString string
// 	if err := json.Unmarshal(b, &eventTypeString); err != nil {
// 		return err
// 	}
//
// 	switch eventTypeString {
// 	case "PlayerSetup":
// 		typ.Type = reflect.TypeOf(PlayerSetup{})
//
// 	case "Upgrade":
// 		typ.Type = reflect.TypeOf(Upgrade{})
//
// 	case "UnitBorn":
// 		typ.Type = reflect.TypeOf(UnitBorn{})
// 	case "UnitInit":
// 		typ.Type = reflect.TypeOf(UnitInit{})
// 	case "UnitDone":
// 		typ.Type = reflect.TypeOf(UnitDone{})
// 	case "UnitTypeChange":
// 		typ.Type = reflect.TypeOf(UnitTypeChange{})
// 	case "UnitOwnerChange":
// 		typ.Type = reflect.TypeOf(UnitOwnerChange{})
// 	case "UnitPositions":
// 		typ.Type = reflect.TypeOf(UnitPositions{})
// 	case "UnitDied":
// 		typ.Type = reflect.TypeOf(UnitDied{})
//
// 	case "PlayerStats":
// 		typ.Type = reflect.TypeOf(PlayerStats{})
//
// 	default:
// 		return fmt.Errorf("Invalid event type: %v", eventTypeString)
// 	}
//
// 	return nil
// }

// Used only for composition
type WithPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Used only for composition
type WithPlayerOwned struct {
	ControlPlayerID int `json:"controlPlayerId"`
	UpkeepPlayerID  int `json:"upkeepPlayerId"`
}

// Used only for composition
type WithUnitTag struct {
	UnitTagIndex   int     `json:"unitTagIndex"`
	UnitTagRecycle IntBool `json:"unitTagRecycle"`
}

// Used only for composition
type WithUnitName struct {
	UnitTypeName string `json:"unitTypeName"`
}

type PlayerSetup struct {
	BaseEvent
	PlayerID int `json:"playerId"`
	UserID   int `json:"userId"`
	SlotID   int `json:"slotId"`
	Type     int `json:"type"`
}

type Upgrade struct {
	BaseEvent
	PlayerID        int    `json:"playerId"`
	Count           int    `json:"count"`
	UpgradeTypeName string `json:"upgradeTypeName"`
}

type UnitBorn struct {
	BaseEvent
	WithPosition
	WithPlayerOwned
	WithUnitTag
	WithUnitName

	CreatorAbilityName    *string  `json:"creatorAbilityName"`
	CreatorUnitTagInded   *int     `json:"creatorUnitTagIndex"`
	CreatorUnitTagRecycle *IntBool `json:"creatorUnitTagRecycle"`
}

type UnitInit struct {
	BaseEvent
	WithPosition
	WithPlayerOwned
	WithUnitTag
	WithUnitName
}

type UnitDone struct {
	BaseEvent
	WithUnitTag
}

type UnitTypeChange struct {
	BaseEvent
	WithUnitTag
	WithUnitName
}

type UnitPositions struct {
	BaseEvent

	FirstUnitIndex int   `json:"firstUnitIndex"`
	Items          []int `json:"items"`
}

type UnitOwnerChange struct {
	BaseEvent
	WithUnitTag
	WithPlayerOwned
}

type UnitDied struct {
	BaseEvent
	WithPosition
	WithUnitTag

	KillerPlayerID       *int     `json:"killerPlayerId"`
	KillerUnitTagIndex   *int     `json:"killerUnitTagIndex"`
	KillerUnitTagRecycle *IntBool `json:"killerUnitTagRecycle"`
}

type PlayerStats struct {
	BaseEvent
	PlayerID int `json:"playerId"`

	Stats Stats
}

type Stats struct {
	FoodMade int `json:"scoreValueFoodMade"`
	FoodUsed int `json:"scoreValueFoodUsed"`

	MineralsCollectionRate int `json:"scoreValueMineralsCollectionRate"`
	MineralsCurrent        int `json:"scoreValueMineralsCurrent"`

	MineralsFriendlyFireArmy       int `json:"scoreValueMineralsFriendlyFireArmy"`
	MineralsFriendlyFireEconomy    int `json:"scoreValueMineralsFriendlyFireEconomy"`
	MineralsFriendlyFireTechnology int `json:"scoreValueMineralsFriendlyFireEconomy"`

	MineralsKilledArmy       int `json:"scoreValueMineralsKilledArmy"`
	MineralsKilledEconomy    int `json:"scoreValueMineralsKilledEconomy"`
	MineralsKilledTechnology int `json:"scoreValueMineralsKilledTechnology"`

	MineralsLostArmy       int `json:"scoreValueMineralsLostArmy"`
	MineralsLostEconomy    int `json:"scoreValueMineralsLostEconomy"`
	MineralsLostTechnology int `json:"scoreValueMineralsLostTechnology"`

	MineralsUsedActiveForces      int `json:"scoreValueMineralsUsedActiveForces"`
	MineralsUsedCurrentArmy       int `json:"scoreValueMineralsUsedCurrentArmy"`
	MineralsUsedCurrentEconomy    int `json:"scoreValueMineralsUsedCurrentEconomy"`
	MineralsUsedCurrentTechnology int `json:"scoreValueMineralsUsedCurrentTechnology"`

	MineralsUsedInProgressArmy       int `json:"scoreValueMineralsUsedInProgressArmy"`
	MineralsUsedInProgressEconomy    int `json:"scoreValueMineralsUsedInProgressEconomy"`
	MineralsUsedInProgressTechnology int `json:"scoreValueMineralsUsedInProgressTechnology"`

	VespeneCollectionRate int `json:"scoreValueVespeneCollectionRate"`
	VespeneCurrent        int `json:"scoreValueVespeneCurrent"`

	VespeneFriendlyFireArmy       int `json:"scoreValueVespeneFriendlyFireArmy"`
	VespeneFriendlyFireEconomy    int `json:"scoreValueVespeneFriendlyFireEconomy"`
	VespeneFriendlyFireTechnology int `json:"scoreValueVespeneFriendlyFireEconomy"`

	VespeneKilledArmy       int `json:"scoreValueVespeneKilledArmy"`
	VespeneKilledEconomy    int `json:"scoreValueVespeneKilledEconomy"`
	VespeneKilledTechnology int `json:"scoreValueVespeneKilledTechnology"`

	VespeneLostArmy       int `json:"scoreValueVespeneLostArmy"`
	VespeneLostEconomy    int `json:"scoreValueVespeneLostEconomy"`
	VespeneLostTechnology int `json:"scoreValueVespeneLostTechnology"`

	VespeneUsedActiveForces      int `json:"scoreValueVespeneUsedActiveForces"`
	VespeneUsedCurrentArmy       int `json:"scoreValueVespeneUsedCurrentArmy"`
	VespeneUsedCurrentEconomy    int `json:"scoreValueVespeneUsedCurrentEconomy"`
	VespeneUsedCurrentTechnology int `json:"scoreValueVespeneUsedCurrentTechnology"`

	VespeneUsedInProgressArmy       int `json:"scoreValueVespeneUsedInProgressArmy"`
	VespeneUsedInProgressEconomy    int `json:"scoreValueVespeneUsedInProgressEconomy"`
	VespeneUsedInProgressTechnology int `json:"scoreValueVespeneUsedInProgressTechnology"`

	WorkersActiveCount int `json:"scoreValueWorkersActiveCount"`
}
