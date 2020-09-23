package units

type Critter struct {
	Name string
}

var Critters = map[string]Critter{
	"KarakMale":   Critter{"Karak"},
	"KarakFemale": Critter{"Karak"},
}
