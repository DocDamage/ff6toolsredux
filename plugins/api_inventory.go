package plugins

import (
	"context"
	"encoding/json"
	modelsPR "ffvi_editor/models/pr"
)

// GetInventory retrieves the current inventory
func (a *APIImpl) GetInventory(ctx context.Context) (*modelsPR.Inventory, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	inv := modelsPR.GetInventory()
	if inv == nil {
		return nil, ErrNilPRData
	}

	return inv, nil
}

// SetInventory updates the inventory
func (a *APIImpl) SetInventory(ctx context.Context, inv *modelsPR.Inventory) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil || a.prData.UserData == nil {
		return ErrNilPRData
	}

	// Get ownedItemList from UserData
	if ownedItemListVal := a.prData.UserData.Get("ownedItemList"); ownedItemListVal != nil {
		if ownedItemListStr, ok := ownedItemListVal.(string); ok {
			// Parse the JSON
			var itemList map[string]interface{}
			if err := json.Unmarshal([]byte(ownedItemListStr), &itemList); err == nil {
				// Update values array
				if inv != nil && inv.Rows != nil {
					values := make([]map[string]interface{}, 0, inv.Size)
					for _, row := range inv.Rows {
						if row != nil && row.ItemID > 0 {
							values = append(values, map[string]interface{}{
								"contentId": row.ItemID,
								"count":     row.Count,
							})
						}
					}
					itemList["values"] = values

					// Marshal back to JSON
					if updatedJSON, err := json.Marshal(itemList); err == nil {
						a.prData.UserData.Set("ownedItemList", string(updatedJSON))
					}
				}
			}
		}
	}

	return nil
}

// FindItems finds inventory items matching a predicate
func (a *APIImpl) FindItems(ctx context.Context, predicate func(*modelsPR.Row) bool) []*modelsPR.Row {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil
	}

	inv := modelsPR.GetInventory()
	if inv == nil {
		return nil
	}

	var results []*modelsPR.Row
	for _, row := range inv.GetRows() {
		if predicate(row) {
			results = append(results, row)
		}
	}
	return results
}
