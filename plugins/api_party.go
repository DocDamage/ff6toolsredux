package plugins

import (
	"context"
	"encoding/json"
	modelsPR "ffvi_editor/models/pr"
)

// GetParty retrieves the current party composition
func (a *APIImpl) GetParty(ctx context.Context) (*modelsPR.Party, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	party := modelsPR.GetParty()
	if party == nil {
		return nil, ErrNilPRData
	}

	return party, nil
}

// SetParty updates the party composition
func (a *APIImpl) SetParty(ctx context.Context, party *modelsPR.Party) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil || a.prData.UserData == nil {
		return ErrNilPRData
	}

	// Get ownedCharacterList from UserData
	if ownedCharListVal := a.prData.UserData.Get("ownedCharacterList"); ownedCharListVal != nil {
		if ownedCharListStr, ok := ownedCharListVal.(string); ok {
			// Parse the JSON
			var charList map[string]interface{}
			if err := json.Unmarshal([]byte(ownedCharListStr), &charList); err == nil {
				// Update values array with party members
				if party != nil {
					values := make([]map[string]interface{}, 0, len(party.Members))
					for _, member := range party.Members {
						if member != nil {
							values = append(values, map[string]interface{}{
								"contentId": member.CharacterID,
								"count":     1,
							})
						}
					}
					charList["values"] = values

					// Marshal back to JSON
					if updatedJSON, err := json.Marshal(charList); err == nil {
						a.prData.UserData.Set("ownedCharacterList", string(updatedJSON))
					}
				}
			}
		}
	}

	return nil
}
