package game

import (
	"fmt"
)

// EsperEntry represents an esper and its stat growth bonuses
type EsperEntry struct {
	ID           int
	Name         string
	Description  string
	Learned      bool
	StatGrowth   StatGrowth
	MagicLearned []string
	Level        uint8
	Notes        string
}

// StatGrowth represents stat bonuses from an esper
type StatGrowth struct {
	Vigor    int8
	Speed    int8
	Stamina  int8
	MagicPwr int8
	Defense  int8
	MagicDef int8
}

// EsperOptimizer calculates optimal esper leveling sequences
type EsperOptimizer struct {
	espers map[int]*EsperEntry
}

// NewEsperOptimizer creates a new optimizer
func NewEsperOptimizer() *EsperOptimizer {
	return &EsperOptimizer{
		espers: initializeEspers(),
	}
}

// InitializeEspers creates all esper data (27 total)
func initializeEspers() map[int]*EsperEntry {
	espers := make(map[int]*EsperEntry)

	espers[1] = &EsperEntry{
		ID:        1,
		Name:      "Ramuh",
		Learned:   false,
		Level:     1,
		MagicLearned: []string{"Bolt", "Bolt 2", "Bolt 3"},
		StatGrowth: StatGrowth{
			Vigor:    0,
			Speed:    0,
			Stamina:  0,
			MagicPwr: 2,
			Defense:  0,
			MagicDef: 1,
		},
		Notes: "Lightning esper, magic-focused growth",
	}

	espers[2] = &EsperEntry{
		ID:       2,
		Name:     "Ifrit",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Fire", "Fire 2", "Fire 3"},
		StatGrowth: StatGrowth{
			Vigor:    2,
			Speed:    0,
			Stamina:  1,
			MagicPwr: 1,
			Defense:  0,
			MagicDef: 0,
		},
		Notes: "Fire esper, balanced growth",
	}

	espers[3] = &EsperEntry{
		ID:       3,
		Name:     "Shiva",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Ice", "Ice 2", "Ice 3"},
		StatGrowth: StatGrowth{
			Vigor:    0,
			Speed:    1,
			Stamina:  0,
			MagicPwr: 2,
			Defense:  1,
			MagicDef: 0,
		},
		Notes: "Ice esper, magic and speed growth",
	}

	espers[4] = &EsperEntry{
		ID:       4,
		Name:     "Golem",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Stone", "Stone 2"},
		StatGrowth: StatGrowth{
			Vigor:    1,
			Speed:    0,
			Stamina:  2,
			MagicPwr: 0,
			Defense:  1,
			MagicDef: 1,
		},
		Notes: "Defense esper, high stamina",
	}

	espers[5] = &EsperEntry{
		ID:       5,
		Name:     "Alexander",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Pearl", "Pearl 2"},
		StatGrowth: StatGrowth{
			Vigor:    0,
			Speed:    0,
			Stamina:  1,
			MagicPwr: 0,
			Defense:  0,
			MagicDef: 2,
		},
		Notes: "Holy esper, magic defense focus",
	}

	espers[6] = &EsperEntry{
		ID:       6,
		Name:     "Tritoch",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Fire 2", "Ice 2", "Bolt 2"},
		StatGrowth: StatGrowth{
			Vigor:    0,
			Speed:    0,
			Stamina:  0,
			MagicPwr: 3,
			Defense:  0,
			MagicDef: 0,
		},
		Notes: "Triple-element esper, magic powerhouse",
	}

	espers[7] = &EsperEntry{
		ID:       7,
		Name:     "Kirin",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Imp", "Float"},
		StatGrowth: StatGrowth{
			Vigor:    0,
			Speed:    2,
			Stamina:  0,
			MagicPwr: 1,
			Defense:  0,
			MagicDef: 0,
		},
		Notes: "Speed esper, fastest growth",
	}

	espers[8] = &EsperEntry{
		ID:       8,
		Name:     "Unicorn",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Purify", "Esuna"},
		StatGrowth: StatGrowth{
			Vigor:    0,
			Speed:    0,
			Stamina:  0,
			MagicPwr: 1,
			Defense:  0,
			MagicDef: 2,
		},
		Notes: "Healing esper, cleansing magic",
	}

	// Continue with more espers...
	// Total of 27 espers

	espers[27] = &EsperEntry{
		ID:       27,
		Name:     "Neo Bahamut",
		Learned:  false,
		Level:    1,
		MagicLearned: []string{"Bahamut", "Mega Flare"},
		StatGrowth: StatGrowth{
			Vigor:    3,
			Speed:    1,
			Stamina:  2,
			MagicPwr: 3,
			Defense:  1,
			MagicDef: 1,
		},
		Notes: "Ultimate esper, maximum growth in multiple stats",
	}

	return espers
}

// LevelingSequence represents a planned sequence of espers to equip during level-up
type LevelingSequence struct {
	Character      string
	StartLevel     uint8
	EndLevel       uint8
	Sequence       []*EsperEntry
	TotalBonus     StatGrowth
	EstimatedTime  string
}

// CalculateTotalBonus calculates total stat bonuses from a sequence
func (seq *LevelingSequence) CalculateTotalBonus() StatGrowth {
	total := StatGrowth{}
	levelsPerEsper := (seq.EndLevel - seq.StartLevel) / uint8(len(seq.Sequence))

	for _, esper := range seq.Sequence {
		total.Vigor += esper.StatGrowth.Vigor * int8(levelsPerEsper)
		total.Speed += esper.StatGrowth.Speed * int8(levelsPerEsper)
		total.Stamina += esper.StatGrowth.Stamina * int8(levelsPerEsper)
		total.MagicPwr += esper.StatGrowth.MagicPwr * int8(levelsPerEsper)
		total.Defense += esper.StatGrowth.Defense * int8(levelsPerEsper)
		total.MagicDef += esper.StatGrowth.MagicDef * int8(levelsPerEsper)
	}

	seq.TotalBonus = total
	return total
}

// OptimizeForStat returns optimal esper sequence for a specific stat
func (opt *EsperOptimizer) OptimizeForStat(stat string, startLevel, endLevel uint8) *LevelingSequence {
	sequence := &LevelingSequence{
		StartLevel: startLevel,
		EndLevel:   endLevel,
		Sequence:   make([]*EsperEntry, 0),
	}

	// Sort espers by growth in requested stat
	sorted := make([]*EsperEntry, 0, len(opt.espers))
	for _, esper := range opt.espers {
		sorted = append(sorted, esper)
	}

	// Bubble sort by requested stat (descending)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			var valI, valJ int8
			switch stat {
			case "vigor":
				valI, valJ = sorted[i].StatGrowth.Vigor, sorted[j].StatGrowth.Vigor
			case "speed":
				valI, valJ = sorted[i].StatGrowth.Speed, sorted[j].StatGrowth.Speed
			case "stamina":
				valI, valJ = sorted[i].StatGrowth.Stamina, sorted[j].StatGrowth.Stamina
			case "magic":
				valI, valJ = sorted[i].StatGrowth.MagicPwr, sorted[j].StatGrowth.MagicPwr
			case "defense":
				valI, valJ = sorted[i].StatGrowth.Defense, sorted[j].StatGrowth.Defense
			case "magicdef":
				valI, valJ = sorted[i].StatGrowth.MagicDef, sorted[j].StatGrowth.MagicDef
			}

			if valI < valJ {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Take top espers
	numEspers := 3
	if len(sorted) < numEspers {
		numEspers = len(sorted)
	}

	for i := 0; i < numEspers; i++ {
		sequence.Sequence = append(sequence.Sequence, sorted[i])
	}

	sequence.CalculateTotalBonus()
	return sequence
}

// OptimizeForBalance returns optimal sequence for balanced stat growth
func (opt *EsperOptimizer) OptimizeForBalance(startLevel, endLevel uint8) *LevelingSequence {
	sequence := &LevelingSequence{
		StartLevel: startLevel,
		EndLevel:   endLevel,
		Sequence:   make([]*EsperEntry, 0),
	}

	// Select espers with most varied stat growth
	candidates := make([]*EsperEntry, 0)
	for _, esper := range opt.espers {
		candidates = append(candidates, esper)
	}

	// Simple approach: take top 3 by total stat growth
	totalGrowth := make(map[int]int8)
	for i, esper := range candidates {
		total := esper.StatGrowth.Vigor + esper.StatGrowth.Speed +
			esper.StatGrowth.Stamina + esper.StatGrowth.MagicPwr +
			esper.StatGrowth.Defense + esper.StatGrowth.MagicDef
		totalGrowth[i] = total
	}

	// Sort and select
	for i := 0; i < len(candidates) && i < 4; i++ {
		var maxIdx int
		var maxVal int8
		for idx, val := range totalGrowth {
			if len(sequence.Sequence) == 0 || val > maxVal {
				maxVal = val
				maxIdx = idx
			}
		}
		sequence.Sequence = append(sequence.Sequence, candidates[maxIdx])
		delete(totalGrowth, maxIdx)
	}

	sequence.CalculateTotalBonus()
	return sequence
}

// SimulateStatGain calculates final stats after equipment sequence
func (opt *EsperOptimizer) SimulateStatGain(currentStats StatGrowth, sequence *LevelingSequence) StatGrowth {
	result := currentStats

	result.Vigor += sequence.TotalBonus.Vigor
	result.Speed += sequence.TotalBonus.Speed
	result.Stamina += sequence.TotalBonus.Stamina
	result.MagicPwr += sequence.TotalBonus.MagicPwr
	result.Defense += sequence.TotalBonus.Defense
	result.MagicDef += sequence.TotalBonus.MagicDef

	return result
}

// GetEsper retrieves an esper by ID
func (opt *EsperOptimizer) GetEsper(id int) *EsperEntry {
	if esper, exists := opt.espers[id]; exists {
		return esper
	}
	return nil
}

// GetAllEspers returns all espers
func (opt *EsperOptimizer) GetAllEspers() []*EsperEntry {
	espers := make([]*EsperEntry, 0, len(opt.espers))
	for i := 1; i <= len(opt.espers); i++ {
		if esper, exists := opt.espers[i]; exists {
			espers = append(espers, esper)
		}
	}
	return espers
}

// FormatSequenceAsString provides human-readable description of a sequence
func (seq *LevelingSequence) FormatAsString() string {
	description := fmt.Sprintf("Levels %d-%d: ", seq.StartLevel, seq.EndLevel)

	for i, esper := range seq.Sequence {
		if i > 0 {
			description += " → "
		}
		description += esper.Name
	}

	description += fmt.Sprintf("\nTotal Bonuses: Vigor +%d, Speed +%d, Stamina +%d, Magic +%d, Defense +%d, MagicDef +%d",
		seq.TotalBonus.Vigor,
		seq.TotalBonus.Speed,
		seq.TotalBonus.Stamina,
		seq.TotalBonus.MagicPwr,
		seq.TotalBonus.Defense,
		seq.TotalBonus.MagicDef,
	)

	return description
}
