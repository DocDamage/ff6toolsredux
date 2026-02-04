package pr

import (
	"fmt"
	"sort"
)

// DiffType represents the type of difference
type DiffType int

const (
	DiffAdded DiffType = iota
	DiffRemoved
	DiffModified
	DiffSame
)

func (d DiffType) String() string {
	switch d {
	case DiffAdded:
		return "Added"
	case DiffRemoved:
		return "Removed"
	case DiffModified:
		return "Modified"
	case DiffSame:
		return "Same"
	default:
		return "Unknown"
	}
}

// Diff represents a single difference
type Diff struct {
	Type     DiffType
	Category string      // e.g., "Character", "Equipment", "Inventory"
	Name     string      // e.g., "Terra", "Item #5"
	Field    string      // e.g., "Level", "HP", "Equipped"
	OldValue interface{} // Previous value
	NewValue interface{} // Current value
}

// DiffReport represents comparison results
type DiffReport struct {
	Diffs      []Diff
	Statistics DiffStatistics
}

// DiffStatistics provides summary of changes
type DiffStatistics struct {
	TotalDiffs    int
	Added         int
	Removed       int
	Modified      int
	CharacterDiff CharacterDiffStats
	EquipmentDiff EquipmentDiffStats
	InventoryDiff InventoryDiffStats
	EsperDiff     EsperDiffStats
}

// CharacterDiffStats tracks character changes
type CharacterDiffStats struct {
	ChangedCount int
	LevelChanges int
	HPChanges    int
	MPChanges    int
	StatChanges  int
}

// EquipmentDiffStats tracks equipment changes
type EquipmentDiffStats struct {
	ChangedCount     int
	WeaponChanges    int
	ArmorChanges     int
	AccessoryChanges int
}

// InventoryDiffStats tracks inventory changes
type InventoryDiffStats struct {
	ItemsAdded   int
	ItemsRemoved int
	ItemsMoved   int
}

// EsperDiffStats tracks esper changes
type EsperDiffStats struct {
	EspersAdded   int
	EspersRemoved int
	EspersLearned int
}

// Comparator compares two save files
type Comparator struct {
	old *PR
	new *PR
}

// NewComparator creates a new comparator
func NewComparator(oldSave, newSave *PR) *Comparator {
	return &Comparator{
		old: oldSave,
		new: newSave,
	}
}

// Compare generates a comprehensive diff report
func (c *Comparator) Compare() DiffReport {
	report := DiffReport{
		Diffs: make([]Diff, 0),
		Statistics: DiffStatistics{
			CharacterDiff: CharacterDiffStats{},
			EquipmentDiff: EquipmentDiffStats{},
			InventoryDiff: InventoryDiffStats{},
			EsperDiff:     EsperDiffStats{},
		},
	}

	// Compare characters
	c.compareCharacters(&report)

	// Compare map data
	c.compareMapData(&report)

	// Update totals
	report.Statistics.TotalDiffs = len(report.Diffs)
	for _, diff := range report.Diffs {
		switch diff.Type {
		case DiffAdded:
			report.Statistics.Added++
		case DiffRemoved:
			report.Statistics.Removed++
		case DiffModified:
			report.Statistics.Modified++
		}
	}

	return report
}

// compareCharacters compares character data
func (c *Comparator) compareCharacters(report *DiffReport) {
	for i := 0; i < len(c.old.Characters) && i < len(c.new.Characters); i++ {
		oldChar := c.old.Characters[i]
		newChar := c.new.Characters[i]

		if oldChar == nil || newChar == nil {
			continue
		}

		changedCount := 0

		// Get character name or fallback to index
		charName := fmt.Sprintf("Character %d", i)
		if oldName := oldChar.Get("Name"); oldName != nil {
			if nameStr, ok := oldName.(string); ok {
				charName = nameStr
			}
		} else if newName := newChar.Get("Name"); newName != nil {
			if nameStr, ok := newName.(string); ok {
				charName = nameStr
			}
		}

		// Compare basic stats
		for _, field := range []string{"Level", "HP", "MaxHP", "MP", "MaxMP"} {
			oldVal := oldChar.Get(field)
			newVal := newChar.Get(field)

			if oldVal != newVal {
				report.Diffs = append(report.Diffs, Diff{
					Type:     DiffModified,
					Category: "Character",
					Name:     charName,
					Field:    field,
					OldValue: oldVal,
					NewValue: newVal,
				})
				changedCount++

				switch field {
				case "Level":
					report.Statistics.CharacterDiff.LevelChanges++
				case "HP":
					report.Statistics.CharacterDiff.HPChanges++
				case "MP":
					report.Statistics.CharacterDiff.MPChanges++
				}
			}
		}

		if changedCount > 0 {
			report.Statistics.CharacterDiff.ChangedCount++
		}
	}
}

// compareMapData compares map data
func (c *Comparator) compareMapData(report *DiffReport) {
	// Compare specific common fields
	commonFields := []string{"MapID", "PlayerPosX", "PlayerPosY", "PlayerFacing", "Airship", "ReliquaryLocations"}

	for _, field := range commonFields {
		oldVal := c.old.MapData.Get(field)
		newVal := c.new.MapData.Get(field)

		if oldVal != newVal && oldVal != nil && newVal != nil {
			report.Diffs = append(report.Diffs, Diff{
				Type:     DiffModified,
				Category: "MapData",
				Name:     "Map",
				Field:    field,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}
}

// GetSortedDiffs returns diffs sorted by category
func (r *DiffReport) GetSortedDiffs() []Diff {
	diffs := make([]Diff, len(r.Diffs))
	copy(diffs, r.Diffs)

	sort.Slice(diffs, func(i, j int) bool {
		if diffs[i].Category != diffs[j].Category {
			return diffs[i].Category < diffs[j].Category
		}
		return diffs[i].Name < diffs[j].Name
	})

	return diffs
}

// GetDiffsByCategory returns diffs filtered by category
func (r *DiffReport) GetDiffsByCategory(category string) []Diff {
	filtered := make([]Diff, 0)
	for _, diff := range r.Diffs {
		if diff.Category == category {
			filtered = append(filtered, diff)
		}
	}
	return filtered
}

// String provides human-readable summary
func (s *DiffStatistics) String() string {
	return fmt.Sprintf(
		"Differences: %d total (%d added, %d removed, %d modified)\n"+
			"Characters: %d changed (%d level, %d HP, %d MP, %d stat)",
		s.TotalDiffs, s.Added, s.Removed, s.Modified,
		s.CharacterDiff.ChangedCount, s.CharacterDiff.LevelChanges, s.CharacterDiff.HPChanges, s.CharacterDiff.MPChanges, s.CharacterDiff.StatChanges,
	)
}
