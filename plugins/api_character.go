package plugins

import (
	"context"
	"encoding/json"
	"ffvi_editor/models"
	modelsPR "ffvi_editor/models/pr"
	"fmt"

	jo "gitlab.com/c0b/go-ordered-json"
)

// extractIntFromMap extracts an integer value from an OrderedMap by key
func extractIntFromMap(m *jo.OrderedMap, key string) (int, bool) {
	if val := m.Get(key); val != nil {
		if num, ok := val.(json.Number); ok {
			if n, err := num.Int64(); err == nil {
				return int(n), true
			}
		}
	}
	return 0, false
}

// extractBoolFromMap extracts a boolean value from an OrderedMap by key
func extractBoolFromMap(m *jo.OrderedMap, key string) (bool, bool) {
	if val := m.Get(key); val != nil {
		if b, ok := val.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// extractCharacterStats extracts HP, MP, Level, and base stats from paramMap
func extractCharacterStats(paramMap *jo.OrderedMap, char *models.Character, baseOffset *modelsPR.CharacterBase) {
	// Extract HP
	if hp, ok := extractIntFromMap(paramMap, "currentHP"); ok {
		char.HP.Current = hp
	}
	if maxHP, ok := extractIntFromMap(paramMap, "addtionalMaxHp"); ok {
		char.HP.Max = maxHP
		if baseOffset != nil {
			char.HP.Max += baseOffset.HPBase
		}
	}

	// Extract MP
	if mp, ok := extractIntFromMap(paramMap, "currentMP"); ok {
		char.MP.Current = mp
	}
	if maxMP, ok := extractIntFromMap(paramMap, "addtionalMaxMp"); ok {
		char.MP.Max = maxMP
		if baseOffset != nil {
			char.MP.Max += baseOffset.MPBase
		}
	}

	// Extract Level
	if level, ok := extractIntFromMap(paramMap, "addtionalLevel"); ok {
		char.Level = level
	}

	// Extract Vigor (Power)
	if vigor, ok := extractIntFromMap(paramMap, "addtionalPower"); ok {
		char.Vigor = vigor
	}

	// Extract Stamina (Vitality)
	if stamina, ok := extractIntFromMap(paramMap, "addtionalVitality"); ok {
		char.Stamina = stamina
	}

	// Extract Speed (Agility)
	if speed, ok := extractIntFromMap(paramMap, "addtionalAgility"); ok {
		char.Speed = speed
	}

	// Extract Magic
	if magic, ok := extractIntFromMap(paramMap, "addtionMagic"); ok {
		char.Magic = magic
	}
}

// findCharacterMap finds a character OrderedMap by name
func (a *APIImpl) findCharacterMap(name string) *jo.OrderedMap {
	for _, charMap := range a.prData.Characters {
		if charMap == nil {
			continue
		}
		charNameVal := charMap.Get("name")
		if charNameVal == nil {
			continue
		}
		if charName, ok := charNameVal.(string); ok && charName == name {
			return charMap
		}
	}
	return nil
}

// GetCharacter retrieves a character by name
func (a *APIImpl) GetCharacter(ctx context.Context, name string) (*models.Character, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	if a.prData == nil {
		return nil, ErrNilPRData
	}

	charMap := a.findCharacterMap(name)
	if charMap == nil {
		return nil, ErrCharacterNotFound
	}

	char := modelsPR.GetCharacter(name)
	if char == nil {
		return nil, ErrCharacterNotFound
	}

	// Extract enabled status
	if enabled, ok := extractBoolFromMap(charMap, "isEnableCorps"); ok {
		char.IsEnabled = enabled
	}

	// Extract experience
	if exp, ok := extractIntFromMap(charMap, "currentExp"); ok {
		char.Exp = exp
	}

	// Extract character ID and job ID for base stats lookup
	charID, _ := extractIntFromMap(charMap, "id")
	char.ID = charID
	jobID, _ := extractIntFromMap(charMap, "jobId")

	baseOffset, _ := modelsPR.GetCharacterBaseOffset(charID, jobID)

	// Extract parameter nested object for stats
	if paramVal := charMap.Get("parameter"); paramVal != nil {
		if paramStr, ok := paramVal.(string); ok {
			paramMap := jo.NewOrderedMap()
			if err := paramMap.UnmarshalJSON([]byte(paramStr)); err == nil {
				extractCharacterStats(paramMap, char, baseOffset)
			}
		}
	}

	// Extract equipment, spells, and commands
	a.extractEquipment(charMap, char)
	a.extractSpells(charMap, char)
	a.extractCommands(charMap, char)

	return char, nil
}

// getCharacterIDs extracts character ID and job ID from character map
func getCharacterIDs(charMap *jo.OrderedMap) (charID, jobID int) {
	if idVal := charMap.Get("id"); idVal != nil {
		if idNum, ok := idVal.(json.Number); ok {
			if id, err := idNum.Int64(); err == nil {
				charID = int(id)
			}
		}
	}
	if jobIDVal := charMap.Get("jobId"); jobIDVal != nil {
		if jobIDNum, ok := jobIDVal.(json.Number); ok {
			if id, err := jobIDNum.Int64(); err == nil {
				jobID = int(id)
			}
		}
	}
	return
}

// updateCharacterParams updates the parameter object with character stats
func updateCharacterParams(paramMap *jo.OrderedMap, ch *models.Character, baseOffset *modelsPR.CharacterBase) {
	// Update HP
	if ch.HP.Current >= 0 {
		paramMap.Set("currentHP", json.Number(fmt.Sprintf("%d", ch.HP.Current)))
	}
	if ch.HP.Max >= 0 && baseOffset != nil {
		additionalMaxHp := ch.HP.Max - baseOffset.HPBase
		if additionalMaxHp < 0 {
			additionalMaxHp = 0
		}
		paramMap.Set("addtionalMaxHp", json.Number(fmt.Sprintf("%d", additionalMaxHp)))
	}

	// Update MP
	if ch.MP.Current >= 0 {
		paramMap.Set("currentMP", json.Number(fmt.Sprintf("%d", ch.MP.Current)))
	}
	if ch.MP.Max >= 0 && baseOffset != nil {
		additionalMaxMp := ch.MP.Max - baseOffset.MPBase
		if additionalMaxMp < 0 {
			additionalMaxMp = 0
		}
		paramMap.Set("addtionalMaxMp", json.Number(fmt.Sprintf("%d", additionalMaxMp)))
	}

	// Update Level
	if ch.Level > 0 {
		paramMap.Set("addtionalLevel", json.Number(fmt.Sprintf("%d", ch.Level)))
	}

	// Update stats
	if ch.Vigor >= 0 {
		paramMap.Set("addtionalPower", json.Number(fmt.Sprintf("%d", ch.Vigor)))
	}
	if ch.Stamina >= 0 {
		paramMap.Set("addtionalVitality", json.Number(fmt.Sprintf("%d", ch.Stamina)))
	}
	if ch.Speed >= 0 {
		paramMap.Set("addtionalAgility", json.Number(fmt.Sprintf("%d", ch.Speed)))
	}
	if ch.Magic >= 0 {
		paramMap.Set("addtionMagic", json.Number(fmt.Sprintf("%d", ch.Magic)))
	}
}

// SetCharacter updates a character with full data modification
func (a *APIImpl) SetCharacter(ctx context.Context, name string, ch *models.Character) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil {
		return ErrNilPRData
	}

	// Find character in PR data
	charMap := a.findCharacterMap(name)
	if charMap == nil {
		return ErrCharacterNotFound
	}

	// Update enabled status
	if ch.IsEnabled {
		charMap.Set("isEnableCorps", true)
	}

	// Update experience
	if ch.Exp > 0 {
		charMap.Set("currentExp", json.Number(fmt.Sprintf("%d", ch.Exp)))
	}

	// Get character ID and job ID for base stats
	charID, jobID := getCharacterIDs(charMap)
	baseOffset, _ := modelsPR.GetCharacterBaseOffset(charID, jobID)

	// Update parameter object
	if paramVal := charMap.Get("parameter"); paramVal != nil {
		if paramStr, ok := paramVal.(string); ok {
			paramMap := jo.NewOrderedMap()
			if err := paramMap.UnmarshalJSON([]byte(paramStr)); err == nil {
				updateCharacterParams(paramMap, ch, baseOffset)

				// Marshal back to string
				if paramBytes, err := paramMap.MarshalJSON(); err == nil {
					charMap.Set("parameter", string(paramBytes))
				}
			}
		}
	}

	return nil
}

// FindCharacter finds a character matching a predicate
func (a *APIImpl) FindCharacter(ctx context.Context, predicate func(*models.Character) bool) *models.Character {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil
	}

	for _, char := range a.prData.Characters {
		if char == nil {
			continue
		}
		// Get character name
		nameVal := char.Get("name")
		if nameVal == nil {
			continue
		}
		name, ok := nameVal.(string)
		if !ok {
			continue
		}

		// Get full character data
		character, err := a.GetCharacter(ctx, name)
		if err != nil {
			continue
		}

		if predicate(character) {
			return character
		}
	}
	return nil
}

// SetCharacterStat sets a character stat
func (a *APIImpl) SetCharacterStat(charID int, stat string, value int) error {
	// TODO: Set stat in PR save
	return nil
}
