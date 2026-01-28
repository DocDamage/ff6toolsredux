package speedrun

// SpeedrunConfig represents a speedrun preset configuration
type SpeedrunConfig struct {
	ID          string
	Name        string
	Description string
	Category    string // "any%", "100%", "low_level", "solo", "glitchless", "pacifist"
	
	// Character configuration
	Characters []CharacterConfig
	
	// Party configuration
	PartyMembers []string // Character names
	
	// Equipment priority
	EquipmentPriority map[string][]string // Character -> [equipment in priority order]
	
	// Magic learning order
	MagicOrder []string
	
	// Item acquisition priority
	ItemPriority []string
	
	// Level caps and restrictions
	LevelCap      uint8
	AllowSequenceBreaks bool
	AllowMysteryTour bool
	
	// Notes and tips
	Tips []string
}

// CharacterConfig represents a speedrun character configuration
type CharacterConfig struct {
	Name      string
	TargetLevel uint8
	Priority  int // Lower = higher priority
	Equipment []string
	Magic     []string
	Commands  []string
}

// Config package-level variables
var (
	// RegisteredConfigs stores all available speedrun configs
	RegisteredConfigs = map[string]*SpeedrunConfig{}
)

func init() {
	registerAllConfigs()
}

// registerAllConfigs registers all built-in speedrun configs
func registerAllConfigs() {
	RegisteredConfigs["any_percent"] = NewAnyPercentConfig()
	RegisteredConfigs["one_hundred_percent"] = NewOneHundredPercentConfig()
	RegisteredConfigs["low_level"] = NewLowLevelConfig()
	RegisteredConfigs["solo_character"] = NewSoloCharacterConfig()
	RegisteredConfigs["glitchless"] = NewGlitchlessConfig()
	RegisteredConfigs["pacifist"] = NewPacifistConfig()
}

// NewAnyPercentConfig creates the Any% speedrun config
func NewAnyPercentConfig() *SpeedrunConfig {
	return &SpeedrunConfig{
		ID:          "any_percent",
		Name:        "Any% Speedrun",
		Description: "Fastest completion, any methods allowed",
		Category:    "any%",
		Characters: []CharacterConfig{
			{
				Name:        "Terra",
				TargetLevel: 50,
				Priority:    1,
			},
			{
				Name:        "Locke",
				TargetLevel: 45,
				Priority:    2,
			},
			{
				Name:        "Edgar",
				TargetLevel: 45,
				Priority:    3,
			},
			{
				Name:        "Sabin",
				TargetLevel: 45,
				Priority:    4,
			},
		},
		PartyMembers: []string{"Terra", "Locke", "Edgar", "Sabin"},
		EquipmentPriority: map[string][]string{
			"Terra": {"Mithril Rod", "Tiara", "Cotton Robe", "Tincture Ring"},
			"Locke": {"Thief Knife", "Iron Helm", "Leather Armor", "Barrier Ring"},
			"Edgar": {"Steel Lance", "Iron Helm", "Iron Armor", "Barrier Ring"},
			"Sabin": {"Iron Knuckles", "Iron Helm", "Gi", "Power Source"},
		},
		LevelCap:           99,
		AllowSequenceBreaks: true,
		AllowMysteryTour:   true,
		Tips: []string{
			"Skip non-essential boss fights",
			"Use sequence breaks to skip story sections",
			"Equip characters based on next dungeon",
			"Level up only when necessary",
			"Take the fastest route through dungeons",
		},
	}
}

// NewOneHundredPercentConfig creates the 100% speedrun config
func NewOneHundredPercentConfig() *SpeedrunConfig {
	return &SpeedrunConfig{
		ID:          "one_hundred_percent",
		Name:        "100% Speedrun",
		Description: "All espers, all spells, all items acquired",
		Category:    "100%",
		Characters: []CharacterConfig{
			{
				Name:        "Terra",
				TargetLevel: 99,
				Priority:    1,
			},
			{
				Name:        "Locke",
				TargetLevel: 99,
				Priority:    2,
			},
			{
				Name:        "Edgar",
				TargetLevel: 99,
				Priority:    3,
			},
			{
				Name:        "Celes",
				TargetLevel: 99,
				Priority:    4,
			},
		},
		PartyMembers: []string{"Terra", "Locke", "Edgar", "Celes"},
		EquipmentPriority: map[string][]string{
			"Terra": {"Ultima Rod", "Circlet", "Celestriad", "Ring of Power"},
			"Locke": {"Ultima Weapon", "Coronet", "Gaia Gear", "Genji Glove"},
			"Edgar": {"Chainsaw", "Coronet", "Dueling Armor", "Genji Glove"},
			"Celes": {"Ultima Weapon", "Coronet", "Celestriad", "Genji Glove"},
		},
		LevelCap:           99,
		AllowSequenceBreaks: false,
		AllowMysteryTour:   true,
		Tips: []string{
			"Collect all espers throughout the game",
			"Learn all available magic for all characters",
			"Equip all ultimate weapons",
			"Complete all side quests",
			"Max out all stats before final boss",
		},
	}
}

// NewLowLevelConfig creates the Low Level challenge config
func NewLowLevelConfig() *SpeedrunConfig {
	return &SpeedrunConfig{
		ID:          "low_level",
		Name:        "Low Level Challenge",
		Description: "Complete game with characters capped at Level 30",
		Category:    "low_level",
		Characters: []CharacterConfig{
			{
				Name:        "Terra",
				TargetLevel: 30,
				Priority:    1,
			},
			{
				Name:        "Locke",
				TargetLevel: 30,
				Priority:    2,
			},
			{
				Name:        "Edgar",
				TargetLevel: 30,
				Priority:    3,
			},
			{
				Name:        "Sabin",
				TargetLevel: 30,
				Priority:    4,
			},
		},
		PartyMembers: []string{"Terra", "Locke", "Edgar", "Sabin"},
		LevelCap:    30,
		AllowSequenceBreaks: false,
		Tips: []string{
			"Level cap strictly enforced at 30",
			"Rely on equipment and magic rather than raw stats",
			"Strategy is critical - proper party composition essential",
			"Stock healing items for boss fights",
			"Use status effects strategically",
		},
	}
}

// NewSoloCharacterConfig creates the Solo Character run config
func NewSoloCharacterConfig() *SpeedrunConfig {
	return &SpeedrunConfig{
		ID:          "solo_character",
		Name:        "Solo Character Run",
		Description: "Complete game with only one character",
		Category:    "solo",
		Characters: []CharacterConfig{
			{
				Name:        "Terra",
				TargetLevel: 99,
				Priority:    1,
			},
		},
		PartyMembers: []string{"Terra"},
		EquipmentPriority: map[string][]string{
			"Terra": {"Ultima Rod", "Circlet", "Celestriad", "Ring of Power"},
		},
		LevelCap:    99,
		AllowSequenceBreaks: false,
		Tips: []string{
			"Only one character allowed in party at all times",
			"Max out all stats for the solo character",
			"Equip best possible gear",
			"Learn all magic spells",
			"Boss fights will be challenging - prepare potions",
		},
	}
}

// NewGlitchlessConfig creates the Glitchless run config
func NewGlitchlessConfig() *SpeedrunConfig {
	return &SpeedrunConfig{
		ID:          "glitchless",
		Name:        "Glitchless Run",
		Description: "Complete game without using exploits or glitches",
		Category:    "glitchless",
		Characters: []CharacterConfig{
			{
				Name:        "Terra",
				TargetLevel: 99,
				Priority:    1,
			},
			{
				Name:        "Locke",
				TargetLevel: 99,
				Priority:    2,
			},
			{
				Name:        "Edgar",
				TargetLevel: 99,
				Priority:    3,
			},
			{
				Name:        "Celes",
				TargetLevel: 99,
				Priority:    4,
			},
		},
		PartyMembers: []string{"Terra", "Locke", "Edgar", "Celes"},
		LevelCap:    99,
		AllowSequenceBreaks: false,
		AllowMysteryTour: false,
		Tips: []string{
			"No glitches or exploits allowed",
			"Follow the intended game path",
			"Sequence breaks are not permitted",
			"Use only legitimate methods to obtain items",
			"Follow intended boss strategies",
		},
	}
}

// NewPacifistConfig creates the Pacifist run config
func NewPacifistConfig() *SpeedrunConfig {
	return &SpeedrunConfig{
		ID:          "pacifist",
		Name:        "Pacifist Run",
		Description: "Complete game with minimal combat",
		Category:    "pacifist",
		Characters: []CharacterConfig{
			{
				Name:        "Terra",
				TargetLevel: 30,
				Priority:    1,
			},
			{
				Name:        "Locke",
				TargetLevel: 30,
				Priority:    2,
			},
			{
				Name:        "Edgar",
				TargetLevel: 30,
				Priority:    3,
			},
			{
				Name:        "Celes",
				TargetLevel: 30,
				Priority:    4,
			},
		},
		PartyMembers: []string{"Terra", "Locke", "Edgar", "Celes"},
		LevelCap:    30,
		AllowSequenceBreaks: false,
		Tips: []string{
			"Avoid all non-essential battles",
			"Use escape/flee commands when possible",
			"Equip armor for defense, not offense",
			"Stock healing items heavily",
			"Focus on status healing and defensive spells",
		},
	}
}

// GetConfig retrieves a speedrun config by ID
func GetConfig(id string) *SpeedrunConfig {
	if config, exists := RegisteredConfigs[id]; exists {
		return config
	}
	return nil
}

// GetAllConfigs returns all registered speedrun configs
func GetAllConfigs() []*SpeedrunConfig {
	configs := make([]*SpeedrunConfig, 0, len(RegisteredConfigs))
	for _, config := range RegisteredConfigs {
		configs = append(configs, config)
	}
	return configs
}

// ApplyConfigToSave applies a speedrun config to save data
// TODO: This function needs refactoring to work with the actual PR structure
// For now, it's a placeholder. The PR struct uses OrderedMap for characters.
/*
func ApplyConfigToSave(save *pr.Save, config *SpeedrunConfig) error {
	if save == nil || config == nil {
		return nil
	}

	// Apply character configurations
	for _, charConfig := range config.Characters {
		for idx, char := range save.Characters {
			if char.Name == charConfig.Name {
				// Set target level
				if charConfig.TargetLevel > 0 {
					save.Characters[idx].Level = charConfig.TargetLevel
				}

				// Set up equipment if specified
				if len(charConfig.Equipment) > 0 {
					if save.Characters[idx].Equipment == nil {
						save.Characters[idx].Equipment = &models.Equipment{}
					}
					// Equipment would be applied by name lookup
				}

				break
			}
		}
	}

	// Set party members
	if len(config.PartyMembers) > 0 {
		newParty := make([]uint8, 0, len(config.PartyMembers))
		for _, memberName := range config.PartyMembers {
			for idx, char := range save.Characters {
				if char.Name == memberName {
					newParty = append(newParty, uint8(idx))
					break
				}
			}
		}
		save.Party = newParty
	}

	return nil
}
*/
func ApplyConfigToSave(config *SpeedrunConfig) error {
	if config == nil {
		return nil
	}
	// TODO: Implement configuration application
	return nil
}
