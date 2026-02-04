package plugins

import (
	"context"
	"encoding/json"
	"ffvi_editor/models"
	modelsPR "ffvi_editor/models/pr"
)

// GetEquipment retrieves equipment data
func (a *APIImpl) GetEquipment(ctx context.Context) (*models.Equipment, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	// Get first character's equipment as representative
	for _, charMap := range a.prData.Characters {
		if charMap == nil {
			continue
		}

		// Get character name
		nameVal := charMap.Get("name")
		if nameVal == nil {
			continue
		}
		name, ok := nameVal.(string)
		if !ok {
			continue
		}

		// Get character
		char := modelsPR.GetCharacter(name)
		if char == nil {
			continue
		}

		// Extract equipment
		a.extractEquipment(charMap, char)
		return &char.Equipment, nil
	}

	return nil, ErrCharacterNotFound
}

// SetEquipment updates equipment data
func (a *APIImpl) SetEquipment(ctx context.Context, eq *models.Equipment) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil {
		return ErrNilPRData
	}

	// Find first character and update their equipment
	for _, charMap := range a.prData.Characters {
		if charMap == nil {
			continue
		}

		// Get equipmentList
		if eqListVal := charMap.Get("equipmentList"); eqListVal != nil {
			if eqListStr, ok := eqListVal.(string); ok {
				var eqList map[string]interface{}
				if err := json.Unmarshal([]byte(eqListStr), &eqList); err == nil {
					// Build values array from equipment
					values := make([]map[string]interface{}, 6)
					values[0] = map[string]interface{}{"contentId": eq.WeaponID, "count": 1}
					values[1] = map[string]interface{}{"contentId": eq.ShieldID, "count": 1}
					values[2] = map[string]interface{}{"contentId": eq.ArmorID, "count": 1}
					values[3] = map[string]interface{}{"contentId": eq.HelmetID, "count": 1}
					values[4] = map[string]interface{}{"contentId": eq.Relic1ID, "count": 1}
					values[5] = map[string]interface{}{"contentId": eq.Relic2ID, "count": 1}

					eqList["values"] = values

					// Marshal back to JSON
					if updatedJSON, err := json.Marshal(eqList); err == nil {
						charMap.Set("equipmentList", string(updatedJSON))
					}
				}
			}
		}

		return nil
	}

	return ErrCharacterNotFound
}

// ApplyBatchOperation applies a batch operation
func (a *APIImpl) ApplyBatchOperation(ctx context.Context, op string, params map[string]interface{}) (int, error) {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return 0, ErrInsufficientPermissions
	}

	// TODO: Implement batch operations
	return 0, nil
}
