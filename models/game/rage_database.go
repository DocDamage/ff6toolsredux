package game

// RageEntry represents a single rage that Gau can learn
type RageEntry struct {
	ID            int
	Name          string
	Description   string
	Enemy         string
	Location      string
	Level         uint8
	HP            uint16
	Speed         uint8
	Attack        string
	SecondAttack  string
	ElementType   string
	Status        []string
	Difficulty    int // 1-5, difficulty to learn
	Notes         string
}

// RageDatabase stores all 254 rages
var RageDatabase = map[int]*RageEntry{
	1: {
		ID:           1,
		Name:         "Goblin",
		Enemy:        "Goblin",
		Location:     "South Figaro Cave",
		Level:        3,
		HP:           8,
		Speed:        5,
		Attack:       "Knife",
		ElementType:  "Physical",
		Difficulty:   1,
		Notes:        "Very weak, learn early",
	},
	2: {
		ID:           2,
		Name:         "Flan",
		Enemy:        "Flan",
		Location:     "Zozo",
		Level:        8,
		HP:           30,
		Speed:        3,
		Attack:       "Mute",
		ElementType:  "Magic",
		Difficulty:   1,
		Notes:        "Casts Mute",
	},
	3: {
		ID:           3,
		Name:         "Guard",
		Enemy:        "Guard",
		Location:     "Narshes Shrine",
		Level:        10,
		HP:           50,
		Speed:        6,
		Attack:       "Hit",
		ElementType:  "Physical",
		Difficulty:   2,
		Notes:        "Strong physical attacker",
	},
	4: {
		ID:           4,
		Name:         "Spikey Tiger",
		Enemy:        "Spikey Tiger",
		Location:     "Sabin's scenario",
		Level:        15,
		HP:           80,
		Speed:        7,
		Attack:       "Maul",
		ElementType:  "Physical",
		Difficulty:   2,
		Notes:        "Good physical rage",
	},
	5: {
		ID:           5,
		Name:         "Bounty Man",
		Enemy:        "Bounty Man",
		Location:     "Figaro Castle",
		Level:        12,
		HP:           40,
		Speed:        5,
		Attack:       "Scratch",
		ElementType:  "Physical",
		Difficulty:   2,
		Notes:        "Mid-game rage",
	},
	// Adding more rages up to 254 would follow similar pattern
	// For brevity, showing structure with key rages throughout game
	50: {
		ID:           50,
		Name:         "Phantom",
		Enemy:        "Phantom",
		Location:     "Opera House",
		Level:        30,
		HP:           200,
		Speed:        8,
		Attack:       "Curse",
		ElementType:  "Magic",
		Difficulty:   3,
		Notes:        "Mid-game magic rage",
	},
	100: {
		ID:           100,
		Name:         "Seraph",
		Enemy:        "Seraph",
		Location:     "Esper World",
		Level:        60,
		HP:           500,
		Speed:        9,
		Attack:       "Holy",
		ElementType:  "Holy",
		Difficulty:   4,
		Notes:        "Late-game powerful rage",
	},
	150: {
		ID:           150,
		Name:         "Duke Lich",
		Enemy:        "Duke Lich",
		Location:     "Ancient Castle",
		Level:        75,
		HP:           800,
		Speed:        10,
		Attack:       "Death",
		ElementType:  "Death",
		Difficulty:   5,
		Notes:        "Very dangerous late-game rage",
	},
	200: {
		ID:           200,
		Name:         "Behemoth",
		Enemy:        "Behemoth",
		Location:     "Final Dungeon",
		Level:        80,
		HP:           1000,
		Speed:        11,
		Attack:       "Meteor",
		ElementType:  "Magic",
		Difficulty:   5,
		Notes:        "End-game powerhouse",
	},
	254: {
		ID:           254,
		Name:         "Kefka",
		Enemy:        "Kefka",
		Location:     "Final Boss",
		Level:        99,
		HP:           9999,
		Speed:        20,
		Attack:       "Ultima",
		ElementType:  "Magic",
		Difficulty:   5,
		Notes:        "Absolute final rage, nearly impossible",
	},
}

// RageInfo contains stats about the rage database
type RageInfo struct {
	TotalRages      int
	LearnedCount    int
	UnlearnedCount  int
	ProgressPercent float64
}

// GetRage retrieves a rage by ID
func GetRage(id int) *RageEntry {
	if rage, exists := RageDatabase[id]; exists {
		return rage
	}
	return nil
}

// GetRageByName retrieves a rage by name
func GetRageByName(name string) *RageEntry {
	for _, rage := range RageDatabase {
		if rage.Name == name {
			return rage
		}
	}
	return nil
}

// GetAllRages returns all rages sorted by ID
func GetAllRages() []*RageEntry {
	rages := make([]*RageEntry, 0, len(RageDatabase))
	for i := 1; i <= 254; i++ {
		if rage, exists := RageDatabase[i]; exists {
			rages = append(rages, rage)
		}
	}
	return rages
}

// GetRagesByDifficulty returns rages filtered by difficulty
func GetRagesByDifficulty(difficulty int) []*RageEntry {
	var rages []*RageEntry
	for _, rage := range RageDatabase {
		if rage.Difficulty == difficulty {
			rages = append(rages, rage)
		}
	}
	return rages
}

// GetRagesByLocation returns rages found in a specific location
func GetRagesByLocation(location string) []*RageEntry {
	var rages []*RageEntry
	for _, rage := range RageDatabase {
		if rage.Location == location {
			rages = append(rages, rage)
		}
	}
	return rages
}

// GetRageInfo returns information about rage learning progress
func GetRageInfo(learned []bool) RageInfo {
	totalRages := len(RageDatabase)
	learnedCount := 0

	for _, isLearned := range learned {
		if isLearned {
			learnedCount++
		}
	}

	percent := float64(learnedCount) / float64(totalRages) * 100

	return RageInfo{
		TotalRages:      totalRages,
		LearnedCount:    learnedCount,
		UnlearnedCount:  totalRages - learnedCount,
		ProgressPercent: percent,
	}
}

// init initializes the rage database with comprehensive entries
// This would normally load from a data file, but for now we have key entries
func init() {
	// Ensure minimum entries exist
	if len(RageDatabase) == 0 {
		RageDatabase = make(map[int]*RageEntry)
	}

	// Fill in entries to reach 254 rages
	// This is a simplified version; production would have all entries
	for i := 1; i <= 254; i++ {
		if _, exists := RageDatabase[i]; !exists {
			// Create placeholder entry
			RageDatabase[i] = &RageEntry{
				ID:         i,
				Name:       "Unknown Rage " + string(rune(i)),
				Enemy:      "Unknown Enemy",
				Location:   "Unknown Location",
				Level:      uint8((i % 99) + 1),
				HP:         uint16(i * 10),
				Speed:      uint8((i % 20) + 1),
				Attack:     "Attack",
				Difficulty: (i % 5) + 1,
			}
		}
	}
}
