package templates

import (
	"time"

	"ffvi_editor/models"
)

// CharacterTemplate holds a saved character build that can be applied to other characters
type CharacterTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Character   *models.Character `json:"character"`
	Equipment   *models.Equipment `json:"equipment"`
	Spells      []*models.Spell   `json:"spells"`
	Commands    []*models.Command `json:"commands"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Tags        []string          `json:"tags"`
	Favorite    bool              `json:"favorite"`
	Version     int               `json:"version"`
}

// NewCharacterTemplate creates a new template from an existing character
func NewCharacterTemplate(char *models.Character, name, description string, tags []string) *CharacterTemplate {
	return &CharacterTemplate{
		ID:          generateTemplateID(),
		Name:        name,
		Description: description,
		Character: &models.Character{
			ID:       char.ID,
			RootName: char.RootName,
			Name:     char.Name,
			Level:    char.Level,
			Exp:      char.Exp,
			HP:       char.HP,
			MP:       char.MP,
			Vigor:    char.Vigor,
			Stamina:  char.Stamina,
			Speed:    char.Speed,
			Magic:    char.Magic,
		},
		Equipment: &models.Equipment{
			WeaponID: char.Equipment.WeaponID,
			ShieldID: char.Equipment.ShieldID,
			HelmetID: char.Equipment.HelmetID,
			ArmorID:  char.Equipment.ArmorID,
			Relic1ID: char.Equipment.Relic1ID,
			Relic2ID: char.Equipment.Relic2ID,
		},
		Spells:    copySpells(char.SpellsByIndex),
		Commands:  copyCommands(char.Commands),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      tags,
		Favorite:  false,
		Version:   1,
	}
}

// ApplyToCharacter applies this template to a character with specified mode
func (t *CharacterTemplate) ApplyToCharacter(target *models.Character, mode ApplyMode) error {
	if t.Character == nil {
		return ErrInvalidTemplate
	}

	switch mode {
	case ApplyModeReplace:
		// Replace all character data
		target.Level = t.Character.Level
		target.Exp = t.Character.Exp
		target.HP = t.Character.HP
		target.MP = t.Character.MP
		target.Vigor = t.Character.Vigor
		target.Stamina = t.Character.Stamina
		target.Speed = t.Character.Speed
		target.Magic = t.Character.Magic

		if t.Equipment != nil {
			target.Equipment = *t.Equipment
		}

		if len(t.Spells) > 0 {
			target.SpellsByIndex = copySpells(t.Spells)
			target.SpellsByID = make(map[int]*models.Spell)
			for _, spell := range target.SpellsByIndex {
				target.SpellsByID[spell.Index] = spell
			}
		}

		if len(t.Commands) > 0 {
			target.Commands = copyCommands(t.Commands)
		}

	case ApplyModeMerge:
		// Merge: only update if template has different values
		if t.Character.Level > target.Level {
			target.Level = t.Character.Level
		}
		if t.Character.HP.Max > target.HP.Max {
			target.HP.Max = t.Character.HP.Max
			if target.HP.Current > target.HP.Max {
				target.HP.Current = target.HP.Max
			}
		}
		if t.Character.MP.Max > target.MP.Max {
			target.MP.Max = t.Character.MP.Max
			if target.MP.Current > target.MP.Max {
				target.MP.Current = target.MP.Max
			}
		}
		// For equipment, only equip items that are valid
		if t.Equipment != nil {
			target.Equipment = *t.Equipment
		}

	case ApplyModePreview:
		// Validation only, no changes
		if t.Character == nil {
			return ErrInvalidTemplate
		}
	}

	return nil
}

// Helper functions
func copySpells(spells []*models.Spell) []*models.Spell {
	if spells == nil {
		return nil
	}
	result := make([]*models.Spell, len(spells))
	for i, spell := range spells {
		if spell != nil {
			s := *spell
			result[i] = &s
		}
	}
	return result
}

func copyCommands(commands []*models.Command) []*models.Command {
	if commands == nil {
		return nil
	}
	result := make([]*models.Command, len(commands))
	for i, cmd := range commands {
		if cmd != nil {
			c := *cmd
			result[i] = &c
		}
	}
	return result
}

func generateTemplateID() string {
	return time.Now().Format("20060102_150405_") + generateRandomSuffix()
}

func generateRandomSuffix() string {
	return time.Now().Format("000")
}
