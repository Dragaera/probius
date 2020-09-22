package events

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
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

// Used only for composition
type WithPlayerOwned struct {
	ControlPlayerID int64 `json:"controlPlayerId"`
	UpkeepPlayerID  int64 `json:"upkeepPlayerId"`
}

// Used only for composition
type WithUnitTag struct {
	UnitTagIndex   int64 `json:"unitTagIndex"`
	UnitTagRecycle int64 `json:"unitTagRecycle"`
}

// Used only for composition
type WithUnitName struct {
	UnitTypeName string `json:"unitTypeName"`
}

type PlayerSetup struct {
	BaseEvent
	PlayerID int64 `json:"playerId"`
	UserID   int64 `json:"userId"`
	SlotID   int64 `json:"slotId"`
	Type     int64 `json:"type"`
}

type Upgrade struct {
	BaseEvent
	PlayerID        int64  `json:"playerId"`
	Count           int64  `json:"count"`
	UpgradeTypeName string `json:"upgradeTypeName"`
}

type UnitBorn struct {
	BaseEvent
	WithPosition
	WithPlayerOwned
	WithUnitTag
	WithUnitName

	CreatorAbilityName    *string `json:"creatorAbilityName"`
	CreatorUnitTagIndex   *int64  `json:"creatorUnitTagIndex"`
	CreatorUnitTagRecycle *int64  `json:"creatorUnitTagRecycle"`
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

	FirstUnitIndex int64   `json:"firstUnitIndex"`
	Items          []int64 `json:"items"`
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

	KillerPlayerID       *int64 `json:"killerPlayerId"`
	KillerUnitTagIndex   *int64 `json:"killerUnitTagIndex"`
	KillerUnitTagRecycle *int64 `json:"killerUnitTagRecycle"`
}

type PlayerStats struct {
	BaseEvent
	PlayerID int64 `json:"playerId"`

	Stats Stats
}

type Stats struct {
	FoodMade int64 `json:"scoreValueFoodMade"`
	FoodUsed int64 `json:"scoreValueFoodUsed"`

	MineralsCollectionRate int64 `json:"scoreValueMineralsCollectionRate"`
	MineralsCurrent        int64 `json:"scoreValueMineralsCurrent"`

	MineralsFriendlyFireArmy       int64 `json:"scoreValueMineralsFriendlyFireArmy"`
	MineralsFriendlyFireEconomy    int64 `json:"scoreValueMineralsFriendlyFireEconomy"`
	MineralsFriendlyFireTechnology int64 `json:"scoreValueMineralsFriendlyFireEconomy"`

	MineralsKilledArmy       int64 `json:"scoreValueMineralsKilledArmy"`
	MineralsKilledEconomy    int64 `json:"scoreValueMineralsKilledEconomy"`
	MineralsKilledTechnology int64 `json:"scoreValueMineralsKilledTechnology"`

	MineralsLostArmy       int64 `json:"scoreValueMineralsLostArmy"`
	MineralsLostEconomy    int64 `json:"scoreValueMineralsLostEconomy"`
	MineralsLostTechnology int64 `json:"scoreValueMineralsLostTechnology"`

	MineralsUsedActiveForces      int64 `json:"scoreValueMineralsUsedActiveForces"`
	MineralsUsedCurrentArmy       int64 `json:"scoreValueMineralsUsedCurrentArmy"`
	MineralsUsedCurrentEconomy    int64 `json:"scoreValueMineralsUsedCurrentEconomy"`
	MineralsUsedCurrentTechnology int64 `json:"scoreValueMineralsUsedCurrentTechnology"`

	MineralsUsedInProgressArmy       int64 `json:"scoreValueMineralsUsedInProgressArmy"`
	MineralsUsedInProgressEconomy    int64 `json:"scoreValueMineralsUsedInProgressEconomy"`
	MineralsUsedInProgressTechnology int64 `json:"scoreValueMineralsUsedInProgressTechnology"`

	VespeneCollectionRate int64 `json:"scoreValueVespeneCollectionRate"`
	VespeneCurrent        int64 `json:"scoreValueVespeneCurrent"`

	VespeneFriendlyFireArmy       int64 `json:"scoreValueVespeneFriendlyFireArmy"`
	VespeneFriendlyFireEconomy    int64 `json:"scoreValueVespeneFriendlyFireEconomy"`
	VespeneFriendlyFireTechnology int64 `json:"scoreValueVespeneFriendlyFireEconomy"`

	VespeneKilledArmy       int64 `json:"scoreValueVespeneKilledArmy"`
	VespeneKilledEconomy    int64 `json:"scoreValueVespeneKilledEconomy"`
	VespeneKilledTechnology int64 `json:"scoreValueVespeneKilledTechnology"`

	VespeneLostArmy       int64 `json:"scoreValueVespeneLostArmy"`
	VespeneLostEconomy    int64 `json:"scoreValueVespeneLostEconomy"`
	VespeneLostTechnology int64 `json:"scoreValueVespeneLostTechnology"`

	VespeneUsedActiveForces      int64 `json:"scoreValueVespeneUsedActiveForces"`
	VespeneUsedCurrentArmy       int64 `json:"scoreValueVespeneUsedCurrentArmy"`
	VespeneUsedCurrentEconomy    int64 `json:"scoreValueVespeneUsedCurrentEconomy"`
	VespeneUsedCurrentTechnology int64 `json:"scoreValueVespeneUsedCurrentTechnology"`

	VespeneUsedInProgressArmy       int64 `json:"scoreValueVespeneUsedInProgressArmy"`
	VespeneUsedInProgressEconomy    int64 `json:"scoreValueVespeneUsedInProgressEconomy"`
	VespeneUsedInProgressTechnology int64 `json:"scoreValueVespeneUsedInProgressTechnology"`

	WorkersActiveCount int64 `json:"scoreValueWorkersActiveCount"`
}
