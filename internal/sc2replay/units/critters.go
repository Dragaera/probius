package units

type Critter struct {
	Name string
}

var Critters = map[string]Critter{
	"KarakMale":   Critter{"Male Karak"},
	"KarakFemale": Critter{"Female Karak"},
	"CarrionBird": Critter{"Urubu"},
	"Ursadon":     Critter{"Ursadon"},
}
