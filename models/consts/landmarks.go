package consts

// Landmark represents a key location on the world map for snapping teleport.
type Landmark struct {
	Name  string
	World int     // 1 = Balance, 2 = Ruin
	X     float64 // in-game coordinates (0-256)
	Y     float64 // in-game coordinates (0-256)
}

// Landmarks for both worlds (add more as needed)
var Landmarks = []Landmark{
	// World of Balance locations
	{"Narshe", 1, 83.5, 34.4},
	{"Figaro Castle", 1, 64.6, 79.3},
	{"South Figaro", 1, 84.14, 112.8},
	{"Returners' Hideout", 1, 104.8, 66.0},
	{"Kohlingen", 1, 26.3, 39.6},
	{"Jidoor", 1, 27.3, 130.8},
	{"Zozo", 1, 21.0, 93.0},
	{"Thamasa", 1, 249.4, 127.9},
	{"Vector", 1, 119.9, 187.1},
	{"Doma Castle", 1, 155.8, 84.0},
	{"Crescent Mountain", 1, 212.6, 148.8},
	{"Opera House", 1, 45.0, 154.2},
	{"Maranda", 1, 65.9, 226.2},
	{"Albrook", 1, 138.8, 203.8},
	{"Sealed Gate", 1, 166.4, 194.1},
	{"Tzen", 1, 119.0, 149.2},
	{"Figaro Cave", 1, 73.5, 95.0},
	{"Mt. Kolts", 1, 103.1, 101.3},
	{"Lethe River", 1, 144.0, 72.0},
	{"Espers' Gathering Place", 1, 231.2, 129.8},
	{"House in the Veldt", 1, 164.8, 35.4},
	{"Sabin's Cabin", 1, 90.3, 97.7},
	{"The Veldt", 1, 202.2, 132.9},
	{"Barren Falls", 1, 189.4, 95.6},
	{"Phantom Forest", 1, 178.7, 86.5},
	{"Imperial Base", 1, 179.1, 71.2},
	{"Triangle Island", 1, 230.2, 45.0},
	
	// World of Ruin locations
	{"Narshe", 2, 114.3, 33.3},
	{"Figaro Castle", 2, 53.6, 57.85},
	{"South Figaro", 2, 112.3, 96.1},

	{"Kohlingen", 2, 38.3, 44.1},
	{"Jidoor", 2, 33.8, 155.8},
	{"Zozo", 2, 42.4, 130.8},
	{"Thamasa", 2, 249.5, 230.4},

	{"Doma Castle", 2, 171.0, 74.1},

	{"Solitary Island", 2, 73.5, 237.4},
	{"Albrook", 2, 139.9, 209.4},
	{"Tzen", 2, 129.6, 178.8},
	{"Mobliz", 2, 235.8, 135.4},
	{"Nikeah", 2, 145.5, 75.7},
	{"Figaro Cave", 2, 105.6, 99.2},
	{"Darill's Tomb", 2, 25.1, 53.7},
	{"Maranda", 2, 69.5, 183.8},
	{"Cave on the Veldt", 2, 207.8, 88.5},
	{"Colosseum", 2, 50.0, 17.0},
	{"Phoenix Cave", 2, 116.4, 155.8},
	{"Triangle Island", 2, 235.8, 52.2},
	{"Fanatics' Tower", 2, 159.2, 125.2},
	{"Duncan's House", 2, 138.8, 20.1},
	{"The Ancient Castle", 2, 65.9, 71.7},
	{"Gau's Father's House", 2, 177.1, 44.6},
	{"Opera House", 2, 29.6, 182.9},
	{"Kefka's Tower", 2, 136.3, 197.2},
	{"Ebot's Rock", 2, 248.0, 222.6},
	{"The Veldt", 2, 213.9, 78.8},
}
