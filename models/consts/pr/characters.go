package pr

var (
	Characters = []string{
		"Celes",
		"Cyan",
		"Edgar",
		"Locke",
		"Gau",
		"Gogo",
		"Mog",
		"Relm",
		"Sabin",
		"Setzer",
		"Shadow",
		"Strago",
		"Terra",
		"Umaro",
		"??????",
		"Banon",
		"Biggs",
		"Leo",
		"Maduin",
		"Wedge",
		"Cosmog",
		"Moggie",
		"Moghan",
		"Moglin",
		"Mogret",
		"Mogsy",
		"Moguel",
		"Mogwin",
		"Molulu",
		"Mugmug",
	}
	CharacterLookup = make(map[string]int)
)

func init() {
	for i, c := range Characters {
		CharacterLookup[c] = i
	}
}

func GetCharacterIndex(name string) int {
	if i, ok := CharacterLookup[name]; ok {
		return i
	}
	return -1
}

/*
func IsMainCharacter(n string) bool {
	return !(n == "Wedge" || n == "Biggs" || n == "Leo" || n == "Banon")
}
*/
