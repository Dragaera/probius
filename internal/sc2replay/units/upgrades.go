package units

type Upgrade struct {
	Name string
}

var Upgrades = map[string]Upgrade{
	// Protoss
	"ProtossGroundWeaponsLevel1": Upgrade{"Ground Weapons 1"},
	"ProtossGroundWeaponsLevel2": Upgrade{"Ground Weapons 2"},
	"ProtossGroundWeaponsLevel3": Upgrade{"Ground Weapons 3"},
	"ProtossAirWeaponsLevel1":    Upgrade{"Air Weapons 1"},
	"ProtossAirWeaponsLevel2":    Upgrade{"Air Weapons 2"},
	"ProtossAirWeaponsLevel3":    Upgrade{"Air Weapons 3"},
	"ProtossGroundArmorsLevel1":  Upgrade{"Ground Armor 1"},
	"ProtossGroundArmorsLevel2":  Upgrade{"Ground Armor 2"},
	"ProtossGroundArmorsLevel3":  Upgrade{"Ground Armor 3"},
	"ProtossAirArmorsLevel1":     Upgrade{"Air Armor 1"},
	"ProtossAirArmorsLevel2":     Upgrade{"Air Armor 2"},
	"ProtossAirArmorsLevel3":     Upgrade{"Air Armor 3"},
	"ProtossShieldsLevel1":       Upgrade{"Shields 1"},
	"ProtossShieldsLevel2":       Upgrade{"Shields 2"},
	"ProtossShieldsLevel3":       Upgrade{"Shields 3"},
	"TempestGroundAttackUpgrade": Upgrade{"Tectonic Destabilizers"},
	"Charge":                     Upgrade{"Charge"},
	"ObserverGraviticBooster":    Upgrade{"Gravitic Boosters"},
	"GraviticDrive":              Upgrade{"Gravitic Drive"},
	"VoidRaySpeedUpgrade":        Upgrade{"Flux Vanes"},
	"AdeptPiercingAttack":        Upgrade{"Resonating Glaives"},
	"PhoenixRangeUpgrade":        Upgrade{"Anion Pulse-Crystals"},
	"ExtendedThermalLance":       Upgrade{"Extended Thermal Lance"},
	"PsiStormTech":               Upgrade{"Psionic Storm"},
	"BlinkTech":                  Upgrade{"Blink"},
	"DarkTemplarBlinkUpgrade":    Upgrade{"Shadow Stride"},
	"WarpGateResearch":           Upgrade{"Warp Gate"},
	// Terran
	"TerranInfantryWeaponsLevel1":        Upgrade{"Infantry Weapons 1"},
	"TerranInfantryWeaponsLevel2":        Upgrade{"Infantry Weapons 2"},
	"TerranInfantryWeaponsLevel3":        Upgrade{"Infantry Weapons 3"},
	"TerranVehicleWeaponsLevel1":         Upgrade{"Vehicle Weapons 1"},
	"TerranVehicleWeaponsLevel2":         Upgrade{"Vehicle Weapons 2"},
	"TerranVehicleWeaponsLevel3":         Upgrade{"Vehicle Weapons 3"},
	"TerranShipWeaponsLevel1":            Upgrade{"Ship Weapons 1"},
	"TerranShipWeaponsLevel2":            Upgrade{"Ship Weapons 2"},
	"TerranShipWeaponsLevel3":            Upgrade{"Ship Weapons 3"},
	"TerranInfantryArmorsLevel1":         Upgrade{"Infantry Armor 1"},
	"TerranInfantryArmorsLevel2":         Upgrade{"Infantry Armor 2"},
	"TerranInfantryArmorsLevel3":         Upgrade{"Infantry Armor 3"},
	"TerranVehicleAndShipArmorsLevel1":   Upgrade{"Vehicle and Ship Plating 1"},
	"TerranVehicleAndShipArmorsLevel2":   Upgrade{"Vehicle and Ship Plating 2"},
	"TerranVehicleAndShipArmorsLevel3":   Upgrade{"Vehicle and Ship Plating 3"},
	"BansheeSpeed":                       Upgrade{"Hyperflight Rotors"},
	"MedivacIncreaseSpeedBoost":          Upgrade{"Rapid Reignition System"},
	"SmartServos":                        Upgrade{"Smart Servos"},
	"LiberatorAGRangeUpgrade":            Upgrade{"Advanced Ballistics"},
	"EnhancedShockwaves":                 Upgrade{"Enhanced Shockwaves"},
	"HiSecAutoTracking":                  Upgrade{"Hi-Sec Auto Tracking"},
	"CycloneLockOnDamageUpgrade":         Upgrade{"Mag-Field Accelerator"},
	"BansheeCloak":                       Upgrade{"Cloaking Field"},
	"RavenCorvidReactor":                 Upgrade{"Corvid Reactor"},
	"PunisherGrenades":                   Upgrade{"Concussive Shells"},
	"PersonalCloaking":                   Upgrade{"Personal Cloaking"},
	"Stimpack":                           Upgrade{"Stimpack"},
	"BattlecruiserEnableSpecializations": Upgrade{"Weapon Refit"},
	"DrillClaws":                         Upgrade{"Drilling Claws"},
	"ShieldWall":                         Upgrade{"Combat Shield"},
	"HighCapacityBarrels":                Upgrade{"Infernal Pre-Igniter"},
	"TerranBuildingArmor":                Upgrade{"Neosteel Armor"},
	// Zerg
	"ZergMeleeWeaponsLevel1":   Upgrade{"Melee Attacks 1"},
	"ZergMeleeWeaponsLevel2":   Upgrade{"Melee Attacks 2"},
	"ZergMeleeWeaponsLevel3":   Upgrade{"Melee Attacks 3"},
	"ZergMissileWeaponsLevel1": Upgrade{"Missile Attacks 1"},
	"ZergMissileWeaponsLevel2": Upgrade{"Missile Attacks 2"},
	"ZergMissileWeaponsLevel3": Upgrade{"Missile Attacks 3"},
	"ZergFlyerWeaponsLevel1":   Upgrade{"Flyer Attacks 1"},
	"ZergFlyerWeaponsLevel2":   Upgrade{"Flyer Attacks 2"},
	"ZergFlyerWeaponsLevel3":   Upgrade{"Flyer Attacks 3"},
	"ZergGroundArmorsLevel1":   Upgrade{"Ground Carapace 1"},
	"ZergGroundArmorsLevel2":   Upgrade{"Ground Carapace 2"},
	"ZergGroundArmorsLevel3":   Upgrade{"Ground Carapace 3"},
	"ZergFlyerArmorsLevel1":    Upgrade{"Flyer Carapace 1"},
	"ZergFlyerArmorsLevel2":    Upgrade{"Flyer Carapace 2"},
	"ZergFlyerArmorsLevel3":    Upgrade{"Flyer Carapace 3"},
	"ChitinousPlating":         Upgrade{"Chitinous Plating"},
	"DiggingClaws":             Upgrade{"Adaptive Talons"},
	"AnabolicSynthesis":        Upgrade{"Anabolic Synthesis"},
	"CentrificalHooks":         Upgrade{"Centrifugal Hooks"},
	"GlialReconstitution":      Upgrade{"Glial Reconstitution"},
	"zerglingmovementspeed":    Upgrade{"Metabolic Boost"},
	"overlordspeed":            Upgrade{"Pneumatized Carapace"},
	"EvolveMuscularAugments":   Upgrade{"Muscular Augments"},
	"EvolveGroovedSpines":      Upgrade{"Grooved Spines"},
	"LurkerRange":              Upgrade{"Seismic Spines"},
	"Burrow":                   Upgrade{"Burrow"},
	"NeuralParasite":           Upgrade{"NeuralParasite"},
	"InfestorEnergyUpgrade":    Upgrade{"Pathogen Glands"},
	"zerglingattackspeed":      Upgrade{"Adrenal Glands"},
	"TunnelingClaws":           Upgrade{"Tunneling Claws"},
}