package game

// SketchMove represents a move that Relm can sketch from an enemy
type SketchMove struct {
	ID          int
	Name        string
	Enemy       string
	Location    string
	Power       uint8
	Accuracy    uint8
	Type        string // "Physical", "Magic", "Special"
	Element     string
	Effect      string
	Description string
	Difficulty int // 1-5, how rare/hard to sketch
}

// SketchDatabase contains all sketch moves
var SketchDatabase = map[int]*SketchMove{
	1: {
		ID:          1,
		Name:        "Guard",
		Enemy:       "Guard",
		Location:    "Narshes Shrine",
		Power:       5,
		Accuracy:    100,
		Type:        "Physical",
		Effect:      "Defend",
		Description: "Basic defend sketch",
		Difficulty:  1,
	},
	2: {
		ID:          2,
		Name:        "Poison",
		Enemy:       "Flan",
		Location:    "Zozo",
		Power:       0,
		Accuracy:    100,
		Type:        "Magic",
		Element:     "Poison",
		Effect:      "Status effect",
		Description: "Inflicts poison",
		Difficulty:  1,
	},
	3: {
		ID:          3,
		Name:        "Net",
		Enemy:       "Exocyte",
		Location:    "Floating Continent",
		Power:       8,
		Accuracy:    80,
		Type:        "Special",
		Effect:      "Trap",
		Description: "Catches and damages",
		Difficulty:  2,
	},
	4: {
		ID:          4,
		Name:        "Fire Dance",
		Enemy:       "Dancing Zombie",
		Location:    "Doma Castle",
		Power:       15,
		Accuracy:    100,
		Type:        "Magic",
		Element:     "Fire",
		Effect:      "Multi-target fire",
		Description: "Fire attack on group",
		Difficulty:  2,
	},
	5: {
		ID:          5,
		Name:        "Meteor",
		Enemy:       "Atma",
		Location:    "Floating Continent",
		Power:       50,
		Accuracy:    80,
		Type:        "Magic",
		Element:     "None",
		Effect:      "Meteor strike",
		Description: "Powerful random strike",
		Difficulty:  4,
	},
	6: {
		ID:          6,
		Name:        "Ultima",
		Enemy:       "Kefka",
		Location:    "Final Boss",
		Power:       100,
		Accuracy:    100,
		Type:        "Magic",
		Element:     "Absolute",
		Effect:      "Ultimate attack",
		Description: "Most powerful sketch",
		Difficulty:  5,
	},
	// Add more sketch moves...
	// Total typically ranges from 30-50
}

// SketchInfo contains stats about sketches
type SketchInfo struct {
	TotalSketches int
	SketchedCount int
	UnSketchedCount int
	ProgressPercent float64
}

// GetSketch retrieves a sketch move by ID
func GetSketch(id int) *SketchMove {
	if sketch, exists := SketchDatabase[id]; exists {
		return sketch
	}
	return nil
}

// GetSketchByName retrieves a sketch by name
func GetSketchByName(name string) *SketchMove {
	for _, sketch := range SketchDatabase {
		if sketch.Name == name {
			return sketch
		}
	}
	return nil
}

// GetSketchByEnemy retrieves sketches that can be sketched from an enemy
func GetSketchByEnemy(enemy string) []*SketchMove {
	var sketches []*SketchMove
	for _, sketch := range SketchDatabase {
		if sketch.Enemy == enemy {
			sketches = append(sketches, sketch)
		}
	}
	return sketches
}

// GetAllSketches returns all available sketch moves
func GetAllSketches() []*SketchMove {
	sketches := make([]*SketchMove, 0, len(SketchDatabase))
	for i := 1; i <= len(SketchDatabase); i++ {
		if sketch, exists := SketchDatabase[i]; exists {
			sketches = append(sketches, sketch)
		}
	}
	return sketches
}

// GetSketchesByType returns sketches filtered by type
func GetSketchesByType(sketchType string) []*SketchMove {
	var sketches []*SketchMove
	for _, sketch := range SketchDatabase {
		if sketch.Type == sketchType {
			sketches = append(sketches, sketch)
		}
	}
	return sketches
}

// GetSketchesByLocation returns sketches available in a location
func GetSketchesByLocation(location string) []*SketchMove {
	var sketches []*SketchMove
	for _, sketch := range SketchDatabase {
		if sketch.Location == location {
			sketches = append(sketches, sketch)
		}
	}
	return sketches
}

// GetSketchesByDifficulty returns sketches filtered by difficulty
func GetSketchesByDifficulty(difficulty int) []*SketchMove {
	var sketches []*SketchMove
	for _, sketch := range SketchDatabase {
		if sketch.Difficulty == difficulty {
			sketches = append(sketches, sketch)
		}
	}
	return sketches
}

// GetSketchInfo returns statistics about sketch completion
func GetSketchInfo(sketched []bool) SketchInfo {
	totalSketches := len(SketchDatabase)
	sketchedCount := 0

	for _, isSketched := range sketched {
		if isSketched {
			sketchedCount++
		}
	}

	percent := float64(sketchedCount) / float64(totalSketches) * 100

	return SketchInfo{
		TotalSketches:      totalSketches,
		SketchedCount:      sketchedCount,
		UnSketchedCount:    totalSketches - sketchedCount,
		ProgressPercent:    percent,
	}
}
