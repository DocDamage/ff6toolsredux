package plugins

import (
	"encoding/json"
	"ffvi_editor/models"
	modelsPR "ffvi_editor/models/pr"
	"fmt"
	jo "gitlab.com/c0b/go-ordered-json"
)

// Update helper methods for SetCharacter

// updateParameter updates the parameter nested object with character stats
func (a *APIImpl) updateParameter(charMap *jo.OrderedMap, char *models.Character, baseOffset *modelsPR.CharacterBase) error {
	// Get existing parameter or create new one
	paramMap := jo.NewOrderedMap()
	if paramVal := charMap.Get("parameter"); paramVal != nil {
		if paramStr, ok := paramVal.(string); ok {
			// Unmarshal existing parameter
			if err := paramMap.UnmarshalJSON([]byte(paramStr)); err != nil {
				return fmt.Errorf("failed to unmarshal parameter: %w", err)
			}
		}
	}

	// Update HP
	paramMap.Set("currentHP", json.Number(fmt.Sprintf("%d", char.HP.Current)))

	// Calculate additional max HP (subtract base from max)
	additionalMaxHP := char.HP.Max
	if baseOffset != nil {
		additionalMaxHP -= baseOffset.HPBase
	}
	paramMap.Set("addtionalMaxHp", json.Number(fmt.Sprintf("%d", additionalMaxHP)))

	// Update MP
	paramMap.Set("currentMP", json.Number(fmt.Sprintf("%d", char.MP.Current)))

	// Calculate additional max MP (subtract base from max)
	additionalMaxMP := char.MP.Max
	if baseOffset != nil {
		additionalMaxMP -= baseOffset.MPBase
	}
	paramMap.Set("addtionalMaxMp", json.Number(fmt.Sprintf("%d", additionalMaxMP)))

	// Update Level
	paramMap.Set("addtionalLevel", json.Number(fmt.Sprintf("%d", char.Level)))

	// Update Vigor (Power)
	paramMap.Set("addtionalPower", json.Number(fmt.Sprintf("%d", char.Vigor)))

	// Update Stamina (Vitality)
	paramMap.Set("addtionalVitality", json.Number(fmt.Sprintf("%d", char.Stamina)))

	// Update Speed (Agility)
	paramMap.Set("addtionalAgility", json.Number(fmt.Sprintf("%d", char.Speed)))

	// Update Magic
	paramMap.Set("addtionMagic", json.Number(fmt.Sprintf("%d", char.Magic)))

	// Marshal parameter back to JSON string
	paramJSON, err := paramMap.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal parameter: %w", err)
	}

	// Set parameter as JSON string
	charMap.Set("parameter", string(paramJSON))

	return nil
}

// updateEquipment updates equipment slots in character data
func (a *APIImpl) updateEquipment(charMap *jo.OrderedMap, char *models.Character) error {
	// Create equipment structure
	eqMap := jo.NewOrderedMap()

	// Create values array with equipment slots
	values := make([]interface{}, 6)

	// Equipment slots: 0=Weapon, 1=Shield, 2=Armor, 3=Helmet, 4=Relic1, 5=Relic2
	equipmentIDs := []int{
		char.Equipment.WeaponID,
		char.Equipment.ShieldID,
		char.Equipment.ArmorID,
		char.Equipment.HelmetID,
		char.Equipment.Relic1ID,
		char.Equipment.Relic2ID,
	}

	// Build equipment array
	for i, equipID := range equipmentIDs {
		itemMap := make(map[string]interface{})
		itemMap["contentId"] = json.Number(fmt.Sprintf("%d", equipID))
		itemMap["count"] = json.Number("1")
		values[i] = itemMap
	}

	eqMap.Set("values", values)

	// Marshal equipment to JSON string
	eqJSON, err := eqMap.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal equipment: %w", err)
	}

	// Set equipmentList as JSON string
	charMap.Set("equipmentList", string(eqJSON))

	return nil
}

// updateSpells updates spell list in character data
func (a *APIImpl) updateSpells(charMap *jo.OrderedMap, char *models.Character) error {
	// Create ability list structure
	abilityMap := jo.NewOrderedMap()

	// Build array of spell JSONs
	var values []interface{}

	// Iterate through learned spells (filter by spell ID range 1-54)
	for spellID := 1; spellID <= 54; spellID++ {
		if spell, ok := char.SpellsByID[spellID]; ok && spell.Value > 0 {
			// Create ability object for this spell
			abilityObj := jo.NewOrderedMap()
			abilityObj.Set("abilityId", json.Number(fmt.Sprintf("%d", spellID)))
			abilityObj.Set("skillLevel", json.Number(fmt.Sprintf("%d", spell.Value)))

			// Marshal to JSON string
			abilityJSON, err := abilityObj.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal ability: %w", err)
			}

			values = append(values, string(abilityJSON))
		}
	}

	abilityMap.Set("values", values)

	// Marshal ability list to JSON string
	abilityListJSON, err := abilityMap.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal ability list: %w", err)
	}

	// Set abilityList as JSON string
	charMap.Set("abilityList", string(abilityListJSON))

	return nil
}

// updateCommands updates command list in character data
func (a *APIImpl) updateCommands(charMap *jo.OrderedMap, char *models.Character) error {
	// Build command ID array
	commandIDs := make([]interface{}, len(char.Commands))
	for i, cmd := range char.Commands {
		if cmd != nil {
			commandIDs[i] = json.Number(fmt.Sprintf("%d", cmd.Value))
		}
	}

	// Set commandList as array
	charMap.Set("commandList", commandIDs)

	return nil
}
