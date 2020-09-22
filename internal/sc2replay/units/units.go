package units

type Unit struct {
	Name   string
	Supply int
}

var Units = map[string]Unit{
	// Protoss
	"Probe":            Unit{"Probe", 1},
	"Zealot":           Unit{"Zealot", 2},
	"Sentry":           Unit{"Sentry", 2},
	"Stalker":          Unit{"Stalker", 2},
	"Adept":            Unit{"Adept", 2},
	"HighTemplar":      Unit{"High Templar", 2},
	"DarkTemplar":      Unit{"Dark Templar", 2},
	"Archon":           Unit{"Archon", 4},
	"Observer":         Unit{"Observer", 1},
	"WarpPrism":        Unit{"Warp Prism", 2},
	"WarpPrismPhasing": Unit{"Warp Prism", 2},
	"Immortal":         Unit{"Immortal", 4},
	"Colossus":         Unit{"Colossus", 6},
	"Disruptor":        Unit{"Disruptor", 3},
	"Phoenix":          Unit{"Phoenix", 2},
	"VoidRay":          Unit{"Void Ray", 4},
	"Oracle":           Unit{"Oracle", 3},
	"Tempest":          Unit{"Tempest", 5},
	"Carrier":          Unit{"Carrier", 6},
	"Mothership":       Unit{"Mothership", 8},
	// Terran
	"SCV":               Unit{"SCV", 1},
	"Marine":            Unit{"Marine", 1},
	"Marauder":          Unit{"Marauder", 2},
	"Reaper":            Unit{"Reaper", 1},
	"GhostAlternate":    Unit{"Ghost", 2},
	"Hellion":           Unit{"Hellion", 2},
	"HellionTank":       Unit{"Hellbat", 2},
	"WidowMine":         Unit{"Widow Mine", 2},
	"WidowMineBurrowed": Unit{"Widow Mine", 2},
	"SiegeTank":         Unit{"Siege Tank", 3},
	"SiegeTankSieged":   Unit{"Siege Tank", 3},
	"Cyclone":           Unit{"Cyclone", 3},
	"Thor":              Unit{"Thor", 6},
	"ThorAP":            Unit{"Thor", 6},
	"VikingFighter":     Unit{"Viking", 2},
	"VikingAssault":     Unit{"Viking", 2},
	"Medivac":           Unit{"Medivac", 2},
	"Liberator":         Unit{"Liberator", 3},
	"LiberatorAG":       Unit{"Liberator", 3},
	"Banshee":           Unit{"Banshee", 3},
	"Raven":             Unit{"Raven", 2},
	"Battlecruiser":     Unit{"Battlecruiser", 6},
	// Zerg
	"Drone": Unit{"Drone", 1},
	"Queen": Unit{"Queen", 2},
	// TODO do we have UnitSpawn / UnitInit for those, or only via UnitTypeChange from an egg?
	"Zergling": Unit{"Zergling", 0.5},
	"Baneling": Unit{"Baneling", 0.5},
	"Roach":    Unit{"Roach", 2},
	// TODO See above
	"Ravager":   Unit{"Ravager", 3},
	"Hydralisk": Unit{"Hydralisk", 2},
	// TODO See above
	"Lurker":     Unit{"Lurker", 3},
	"Infestor":   Unit{"Infestor", 2},
	"Swarm Host": Unit{"Swarm Host", 3},
	"Ultralisk":  Unit{"Ultralisk", 6},
	"Overlord":   Unit{"Overlord", 0},
	"Overseer":   Unit{"Overseer", 0},
	"Mutalisk":   Unit{"Mutalisk", 2},
	"Corruptor":  Unit{"Corruptor", 2},
	"Viper":      Unit{"Viper", 3},
	"BroodLord":  Unit{"Brood Lord", 4},
}
