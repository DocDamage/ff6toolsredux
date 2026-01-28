package game

// MagicSchool represents different magic types
type MagicSchool string

const (
	SchoolBlack    MagicSchool = "Black"
	SchoolWhite    MagicSchool = "White"
	SchoolBlue     MagicSchool = "Blue"
	SchoolRed      MagicSchool = "Red"
	SchoolSpecial  MagicSchool = "Special"
	SchoolEsper    MagicSchool = "Esper"
	SchoolDance    MagicSchool = "Dance"
	SchoolSkill    MagicSchool = "Skill"
	SchoolLore     MagicSchool = "Lore"
)

// MagicEntry represents a single spell
type MagicEntry struct {
	ID              int
	Name            string
	School          MagicSchool
	Cost            uint8
	Power           uint8
	Accuracy        uint8
	Effect          string
	Description     string
	LearnedBy       []string // Character names
	LearnMethod     string   // "natural", "item", "esper", "lore"
	RequiredLevel   uint8
	CanTarget       string // "single", "multi", "party", "area"
	Element         string
	Status          []string // Status effects applied
}

// MagicCatalog stores all spells organized by school
type MagicCatalog struct {
	Black   map[int]*MagicEntry // 69 spells
	White   map[int]*MagicEntry // 52 spells
	Blue    map[int]*MagicEntry // 24 spells
	Red     map[int]*MagicEntry // 8 spells
	Special map[int]*MagicEntry // 127 spells
	Total   int
}

// GlobalMagicCatalog is the package-level magic database
var GlobalMagicCatalog = NewMagicCatalog()

// NewMagicCatalog creates and initializes a new magic catalog
func NewMagicCatalog() *MagicCatalog {
	catalog := &MagicCatalog{
		Black:   make(map[int]*MagicEntry),
		White:   make(map[int]*MagicEntry),
		Blue:    make(map[int]*MagicEntry),
		Red:     make(map[int]*MagicEntry),
		Special: make(map[int]*MagicEntry),
	}

	initializeBlackMagic(catalog.Black)
	initializeWhiteMagic(catalog.White)
	initializeBlueMagic(catalog.Blue)
	initializeRedMagic(catalog.Red)
	initializeSpecialMagic(catalog.Special)

	// Calculate total
	catalog.Total = len(catalog.Black) + len(catalog.White) + len(catalog.Blue) +
		len(catalog.Red) + len(catalog.Special)

	return catalog
}

// initializeBlackMagic populates black magic spells (69 total)
func initializeBlackMagic(black map[int]*MagicEntry) {
	black[1] = &MagicEntry{
		ID:            1,
		Name:          "Fire",
		School:        SchoolBlack,
		Cost:          4,
		Power:         10,
		Accuracy:      100,
		Effect:        "Deals fire damage",
		Description:   "Fire spell dealing moderate damage",
		LearnMethod:   "natural",
		CanTarget:     "single",
		Element:       "Fire",
	}

	black[2] = &MagicEntry{
		ID:            2,
		Name:          "Ice",
		School:        SchoolBlack,
		Cost:          4,
		Power:         10,
		Accuracy:      100,
		Effect:        "Deals ice damage",
		Description:   "Ice spell dealing moderate damage",
		LearnMethod:   "natural",
		CanTarget:     "single",
		Element:       "Ice",
	}

	black[3] = &MagicEntry{
		ID:            3,
		Name:          "Bolt",
		School:        SchoolBlack,
		Cost:          4,
		Power:         10,
		Accuracy:      100,
		Effect:        "Deals lightning damage",
		Description:   "Lightning spell dealing moderate damage",
		LearnMethod:   "natural",
		CanTarget:     "single",
		Element:       "Lightning",
	}

	black[4] = &MagicEntry{
		ID:            4,
		Name:          "Poison",
		School:        SchoolBlack,
		Cost:          6,
		Power:         8,
		Accuracy:      100,
		Effect:        "Poisons enemy",
		Description:   "Causes poison status",
		LearnMethod:   "natural",
		CanTarget:     "single",
		Status:        []string{"Poison"},
	}

	black[5] = &MagicEntry{
		ID:            5,
		Name:          "Drain",
		School:        SchoolBlack,
		Cost:          8,
		Power:         15,
		Accuracy:      100,
		Effect:        "Drain enemy HP",
		Description:   "Drains enemy HP to caster",
		LearnMethod:   "natural",
		CanTarget:     "single",
	}

	// Continue with more black magic entries...
	// For brevity, showing structure with key spells
	// Production would have all 69

	black[69] = &MagicEntry{
		ID:            69,
		Name:          "Ultima",
		School:        SchoolBlack,
		Cost:          99,
		Power:         100,
		Accuracy:      100,
		Effect:        "Massive damage to all enemies",
		Description:   "Ultimate black magic spell",
		LearnMethod:   "esper",
		RequiredLevel: 90,
		CanTarget:     "multi",
		Element:       "Absolute",
	}
}

// initializeWhiteMagic populates white magic spells (52 total)
func initializeWhiteMagic(white map[int]*MagicEntry) {
	white[1] = &MagicEntry{
		ID:            1,
		Name:          "Cure",
		School:        SchoolWhite,
		Cost:          5,
		Power:         0,
		Accuracy:      100,
		Effect:        "Restores HP",
		Description:   "Basic healing spell",
		LearnMethod:   "natural",
		CanTarget:     "single",
	}

	white[2] = &MagicEntry{
		ID:            2,
		Name:          "Cura",
		School:        SchoolWhite,
		Cost:          15,
		Power:         0,
		Accuracy:      100,
		Effect:        "Restores more HP",
		Description:   "Intermediate healing spell",
		LearnMethod:   "natural",
		RequiredLevel: 15,
		CanTarget:     "single",
	}

	white[3] = &MagicEntry{
		ID:            3,
		Name:          "Curaga",
		School:        SchoolWhite,
		Cost:          30,
		Power:         0,
		Accuracy:      100,
		Effect:        "Restores large amount of HP",
		Description:   "Advanced healing spell",
		LearnMethod:   "esper",
		RequiredLevel: 30,
		CanTarget:     "single",
	}

	white[52] = &MagicEntry{
		ID:            52,
		Name:          "Reraise",
		School:        SchoolWhite,
		Cost:          60,
		Power:         0,
		Accuracy:      100,
		Effect:        "Automatic resurrection",
		Description:   "Auto-revive when HP reaches 0",
		LearnMethod:   "esper",
		RequiredLevel: 80,
		CanTarget:     "single",
	}
	// Continue with more white magic entries...
}

// initializeBlueMagic populates blue magic spells (24 total)
func initializeBlueMagic(blue map[int]*MagicEntry) {
	blue[1] = &MagicEntry{
		ID:            1,
		Name:          "Raging Fists",
		School:        SchoolBlue,
		Cost:          8,
		Power:         20,
		Accuracy:      100,
		Effect:        "Physical attack",
		Description:   "Copied from an enemy ability",
		LearnMethod:   "lore",
		CanTarget:     "single",
	}

	blue[2] = &MagicEntry{
		ID:            2,
		Name:          "Aqualung",
		School:        SchoolBlue,
		Cost:          12,
		Power:         25,
		Accuracy:      100,
		Effect:        "Water damage",
		Description:   "Learned as Blue Magic",
		LearnMethod:   "lore",
		CanTarget:     "multi",
		Element:       "Water",
	}
	// Continue with more blue magic entries...
}

// initializeRedMagic populates red magic spells (8 total)
func initializeRedMagic(red map[int]*MagicEntry) {
	red[1] = &MagicEntry{
		ID:            1,
		Name:          "Fire",
		School:        SchoolRed,
		Cost:          4,
		Power:         10,
		Accuracy:      100,
		Effect:        "Deals fire damage",
		Description:   "Red magic fire spell",
		LearnMethod:   "natural",
		CanTarget:     "single",
		Element:       "Fire",
	}

	red[8] = &MagicEntry{
		ID:            8,
		Name:          "Meteor",
		School:        SchoolRed,
		Cost:          40,
		Power:         50,
		Accuracy:      80,
		Effect:        "Calls meteors",
		Description:   "Deals heavy non-elemental damage",
		LearnMethod:   "esper",
		CanTarget:     "multi",
		Element:       "None",
	}
	// Continue with more red magic entries...
}

// initializeSpecialMagic populates special magic (127 total)
func initializeSpecialMagic(special map[int]*MagicEntry) {
	special[1] = &MagicEntry{
		ID:            1,
		Name:          "Haste",
		School:        SchoolSpecial,
		Cost:          20,
		Power:         0,
		Accuracy:      100,
		Effect:        "Increases speed",
		Description:   "Doubles speed of party",
		LearnMethod:   "natural",
		CanTarget:     "party",
	}

	special[2] = &MagicEntry{
		ID:            2,
		Name:          "Slow",
		School:        SchoolSpecial,
		Cost:          20,
		Power:         0,
		Accuracy:      100,
		Effect:        "Decreases speed",
		Description:   "Halves enemy speed",
		LearnMethod:   "natural",
		CanTarget:     "multi",
	}

	special[127] = &MagicEntry{
		ID:            127,
		Name:          "Apocalypse",
		School:        SchoolSpecial,
		Cost:          99,
		Power:         100,
		Accuracy:      50,
		Effect:        "Massive damage",
		Description:   "Devastating special magic",
		LearnMethod:   "esper",
		RequiredLevel: 99,
		CanTarget:     "multi",
		Element:       "Chaos",
	}
	// Continue with more special magic entries...
}

// GetMagicBySchool returns all spells from a specific school
func (catalog *MagicCatalog) GetMagicBySchool(school MagicSchool) []*MagicEntry {
	var spells []*MagicEntry

	var sourceMap map[int]*MagicEntry

	switch school {
	case SchoolBlack:
		sourceMap = catalog.Black
	case SchoolWhite:
		sourceMap = catalog.White
	case SchoolBlue:
		sourceMap = catalog.Blue
	case SchoolRed:
		sourceMap = catalog.Red
	case SchoolSpecial:
		sourceMap = catalog.Special
	default:
		return spells
	}

	for i := 1; i <= len(sourceMap); i++ {
		if spell, exists := sourceMap[i]; exists {
			spells = append(spells, spell)
		}
	}

	return spells
}

// GetMagicByName returns a spell by name
func (catalog *MagicCatalog) GetMagicByName(name string) *MagicEntry {
	allSchools := []map[int]*MagicEntry{
		catalog.Black,
		catalog.White,
		catalog.Blue,
		catalog.Red,
		catalog.Special,
	}

	for _, school := range allSchools {
		for _, spell := range school {
			if spell.Name == name {
				return spell
			}
		}
	}

	return nil
}

// GetAllMagic returns all spells
func (catalog *MagicCatalog) GetAllMagic() []*MagicEntry {
	var all []*MagicEntry

	allSchools := []map[int]*MagicEntry{
		catalog.Black,
		catalog.White,
		catalog.Blue,
		catalog.Red,
		catalog.Special,
	}

	for _, school := range allSchools {
		for i := 1; i <= len(school); i++ {
			if spell, exists := school[i]; exists {
				all = append(all, spell)
			}
		}
	}

	return all
}

// GetMagicStats returns statistics about magic learning
type MagicStats struct {
	TotalSpells      int
	LearnedCount     int
	UnlearnedCount   int
	ProgressPercent  float64
	BySchool         map[MagicSchool]*SchoolStats
}

// SchoolStats represents stats for a specific magic school
type SchoolStats struct {
	School       MagicSchool
	Total        int
	Learned      int
	Unlearned    int
	Percent      float64
}

// GetMagicStats calculates magic learning statistics
func (catalog *MagicCatalog) GetMagicStats(learned []bool) MagicStats {
	stats := MagicStats{
		TotalSpells: catalog.Total,
		BySchool:    make(map[MagicSchool]*SchoolStats),
	}

	for _, isLearned := range learned {
		if isLearned {
			stats.LearnedCount++
		} else {
			stats.UnlearnedCount++
		}
	}

	stats.ProgressPercent = float64(stats.LearnedCount) / float64(stats.TotalSpells) * 100

	// Calculate per-school stats
	for school, spells := range map[MagicSchool]map[int]*MagicEntry{
		SchoolBlack:   catalog.Black,
		SchoolWhite:   catalog.White,
		SchoolBlue:    catalog.Blue,
		SchoolRed:     catalog.Red,
		SchoolSpecial: catalog.Special,
	} {
		schoolStats := &SchoolStats{
			School: school,
			Total:  len(spells),
		}

		for i := 1; i <= len(spells); i++ {
			if i < len(learned) && learned[i] {
				schoolStats.Learned++
			}
		}

		schoolStats.Unlearned = schoolStats.Total - schoolStats.Learned
		if schoolStats.Total > 0 {
			schoolStats.Percent = float64(schoolStats.Learned) / float64(schoolStats.Total) * 100
		}

		stats.BySchool[school] = schoolStats
	}

	return stats
}
