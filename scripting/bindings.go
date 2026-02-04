package scripting

import (
	"context"
	"fmt"
	"strings"

	modelsPR "ffvi_editor/models/pr"
	"ffvi_editor/plugins"
)

// Bindings provides Go-to-Lua function bindings for the plugin API
type Bindings struct {
	api plugins.PluginAPI
	vm  *VM
}

// NewBindings creates new bindings for a VM
func NewBindings(api plugins.PluginAPI, vm *VM) *Bindings {
	return &Bindings{
		api: api,
		vm:  vm,
	}
}

// Register registers all API functions as Lua globals
func (b *Bindings) Register(ctx context.Context) error {
	if b.api == nil {
		return fmt.Errorf("API is nil")
	}
	if b.vm == nil {
		return fmt.Errorf("VM is nil")
	}

	// Register editor API functions
	// Note: Real implementation would use gopher-lua to register these functions

	// Character functions
	// b.vm.RegisterFunction("editor.getCharacter", func(name string) ...)
	// b.vm.RegisterFunction("editor.setCharacter", func(name string, ch ...) ...)

	// Inventory functions
	// b.vm.RegisterFunction("editor.getInventory", func() ...)
	// b.vm.RegisterFunction("editor.setInventory", func(inv ...) ...)

	// Party functions
	// b.vm.RegisterFunction("editor.getParty", func() ...)
	// b.vm.RegisterFunction("editor.setParty", func(party ...) ...)

	// Batch operations
	// b.vm.RegisterFunction("editor.applyBatch", func(op string, params ...) ...)

	// Query functions
	// b.vm.RegisterFunction("editor.findCharacter", func(predicate ...) ...)
	// b.vm.RegisterFunction("editor.findItems", func(predicate ...) ...)

	// Event functions
	// b.vm.RegisterFunction("editor.registerHook", func(event string, callback ...) ...)
	// b.vm.RegisterFunction("editor.fireEvent", func(event string, data ...) ...)

	// UI functions
	// b.vm.RegisterFunction("editor.showDialog", func(title, message string) ...)
	// b.vm.RegisterFunction("editor.showConfirm", func(title, message string) bool ...)
	// b.vm.RegisterFunction("editor.showInput", func(prompt string) string ...)

	// Logging
	// b.vm.RegisterFunction("editor.log", func(level, message string) ...)

	// Settings
	// b.vm.RegisterFunction("editor.getSetting", func(key string) any ...)
	// b.vm.RegisterFunction("editor.setSetting", func(key string, value any) ...)

	// Permissions
	// b.vm.RegisterFunction("editor.hasPermission", func(perm string) bool ...)

	return nil
}

// BindGetCharacter binds the GetCharacter API function
func (b *Bindings) BindGetCharacter(ctx context.Context) error {
	// Integrate decodeString for type-safe string assertion
	decodeString := func(s string) string { return s } // Placeholder for actual decodeString implementation

	// Real implementation would create a Lua function that calls:
	// return b.api.GetCharacter(ctx, name)
	return b.vm.RegisterFunction("editor.getCharacter", func(name string) interface{} {
		ch, err := b.api.GetCharacter(ctx, decodeString(name))
		if err != nil {
			return nil
		}
		return ch
	})
}

// BindSetCharacter binds the SetCharacter API function
func (b *Bindings) BindSetCharacter(ctx context.Context) error {
	// Integrate decodeString for type-safe string assertion
	decodeString := func(s string) string { return s } // Placeholder for actual decodeString implementation

	// Real implementation would create a Lua function that calls:
	// return b.api.SetCharacter(ctx, name, ch)
	return b.vm.RegisterFunction("editor.setCharacter", func(name string, ch interface{}) interface{} {
		// Would deserialize ch and call b.api.SetCharacter
		name = decodeString(name)
		return nil
	})
}

// BindGetInventory binds the GetInventory API function
func (b *Bindings) BindGetInventory(ctx context.Context) error {
	return b.vm.RegisterFunction("editor.getInventory", func() interface{} {
		inv, err := b.api.GetInventory(ctx)
		if err != nil {
			return nil
		}
		return inv
	})
}

// BindSetInventory binds the SetInventory API function
func (b *Bindings) BindSetInventory(ctx context.Context) error {
	return b.vm.RegisterFunction("editor.setInventory", func(inv interface{}) interface{} {
		// Would deserialize inv and call b.api.SetInventory
		return nil
	})
}

// BindLog binds the Log function
func (b *Bindings) BindLog(ctx context.Context) error {
	return b.vm.RegisterFunction("editor.log", func(level, message string) error {
		return b.api.Log(ctx, level, message)
	})
}

// BindShowDialog binds the ShowDialog function
func (b *Bindings) BindShowDialog(ctx context.Context) error {
	return b.vm.RegisterFunction("editor.showDialog", func(title, message string) error {
		return b.api.ShowDialog(ctx, title, message)
	})
}

// BindShowConfirm binds the ShowConfirm function
func (b *Bindings) BindShowConfirm(ctx context.Context) error {
	return b.vm.RegisterFunction("editor.showConfirm", func(title, message string) bool {
		result, err := b.api.ShowConfirm(ctx, title, message)
		if err != nil {
			return false
		}
		return result
	})
}

// BindShowInput binds the ShowInput function
func (b *Bindings) BindShowInput(ctx context.Context) error {
	return b.vm.RegisterFunction("editor.showInput", func(prompt string) string {
		result, err := b.api.ShowInput(ctx, prompt)
		if err != nil {
			return ""
		}
		return result
	})
}

// BindHasPermission binds the HasPermission function
func (b *Bindings) BindHasPermission() error {
	return b.vm.RegisterFunction("editor.hasPermission", func(perm string) bool {
		return b.api.HasPermission(perm)
	})
}

// BindGetSetting binds the GetSetting function
func (b *Bindings) BindGetSetting() error {
	return b.vm.RegisterFunction("editor.getSetting", func(key string) interface{} {
		decodeString := func(s string) string { return s } // Placeholder for actual decodeString implementation
		return b.api.GetSetting(decodeString(key))
	})
}

// BindSetSetting binds the SetSetting function
func (b *Bindings) BindSetSetting() error {
	return b.vm.RegisterFunction("editor.setSetting", func(key string, value interface{}) error {
		decodeString := func(s string) string { return s } // Placeholder for actual decodeString implementation
		return b.api.SetSetting(decodeString(key), value)
	})
}

// Character functions

func (b *Bindings) getCharacter(charID int) (interface{}, error) {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return nil, fmt.Errorf("character with ID %d not found", charID)
	}
	// Use the API to get the full character data from the save
	return b.api.GetCharacter(context.Background(), char.Name)
}

func (b *Bindings) setCharacterLevel(charID, level int) error {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return fmt.Errorf("character with ID %d not found", charID)
	}
	// Get current character data from save
	ch, err := b.api.GetCharacter(context.Background(), char.Name)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}
	// Update level
	ch.Level = level
	// Save back
	return b.api.SetCharacter(context.Background(), char.Name, ch)
}

func (b *Bindings) setCharacterHP(charID, hp int) error {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return fmt.Errorf("character with ID %d not found", charID)
	}
	// Get current character data from save
	ch, err := b.api.GetCharacter(context.Background(), char.Name)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}
	// Update HP (both current and max)
	ch.HP.Current = hp
	ch.HP.Max = hp
	// Save back
	return b.api.SetCharacter(context.Background(), char.Name, ch)
}

func (b *Bindings) setCharacterMP(charID, mp int) error {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return fmt.Errorf("character with ID %d not found", charID)
	}
	// Get current character data from save
	ch, err := b.api.GetCharacter(context.Background(), char.Name)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}
	// Update MP (both current and max)
	ch.MP.Current = mp
	ch.MP.Max = mp
	// Save back
	return b.api.SetCharacter(context.Background(), char.Name, ch)
}

func (b *Bindings) setCharacterStat(charID int, stat string, value int) error {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return fmt.Errorf("character with ID %d not found", charID)
	}
	// Get current character data from save
	ch, err := b.api.GetCharacter(context.Background(), char.Name)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}
	// Map stat string to character field
	switch strings.ToLower(stat) {
	case "vigor", "power":
		ch.Vigor = value
	case "speed", "agility":
		ch.Speed = value
	case "stamina", "vitality":
		ch.Stamina = value
	case "magic":
		ch.Magic = value
	default:
		return fmt.Errorf("unknown stat: %s (expected: vigor, speed, stamina, or magic)", stat)
	}
	// Save back
	return b.api.SetCharacter(context.Background(), char.Name, ch)
}

// Inventory functions

func (b *Bindings) getItemCount(itemID int) (int, error) {
	// Get inventory from model
	inv := modelsPR.GetInventory()
	if inv == nil {
		return 0, fmt.Errorf("inventory not available")
	}
	// Find item by ID
	for _, row := range inv.Rows {
		if row != nil && row.ItemID == itemID {
			return row.Count, nil
		}
	}
	return 0, nil
}

func (b *Bindings) setItemCount(itemID, count int) error {
	// Get inventory from model
	inv := modelsPR.GetInventory()
	if inv == nil {
		return fmt.Errorf("inventory not available")
	}
	// Find existing item and update count
	for _, row := range inv.Rows {
		if row != nil && row.ItemID == itemID {
			row.Count = count
			if row.Count <= 0 {
				row.ItemID = 0
				row.Count = 0
			}
			return nil
		}
	}
	// Item not found, add to empty slot if count > 0
	if count > 0 {
		for _, row := range inv.Rows {
			if row != nil && row.ItemID == 0 {
				row.ItemID = itemID
				row.Count = count
				return nil
			}
		}
		return fmt.Errorf("inventory is full")
	}
	return nil
}

func (b *Bindings) addItem(itemID, count int) error {
	// Get inventory from model
	inv := modelsPR.GetInventory()
	if inv == nil {
		return fmt.Errorf("inventory not available")
	}
	// Find existing item and add to count
	for _, row := range inv.Rows {
		if row != nil && row.ItemID == itemID {
			row.Count += count
			if row.Count <= 0 {
				row.ItemID = 0
				row.Count = 0
			}
			return nil
		}
	}
	// Item not found, add to empty slot if count > 0
	if count > 0 {
		for _, row := range inv.Rows {
			if row != nil && row.ItemID == 0 {
				row.ItemID = itemID
				row.Count = count
				return nil
			}
		}
		return fmt.Errorf("inventory is full")
	}
	return nil
}

// Party functions

func (b *Bindings) getPartyMembers() ([]int, error) {
	// Get party from model
	party := modelsPR.GetParty()
	if party == nil {
		return nil, fmt.Errorf("party not available")
	}
	// Extract character IDs from party members
	memberIDs := make([]int, 0, 4)
	for _, member := range party.Members {
		if member != nil && member.CharacterID > 0 {
			memberIDs = append(memberIDs, member.CharacterID)
		}
	}
	return memberIDs, nil
}

func (b *Bindings) setPartyMembers(members []int) error {
	// Get party from model
	party := modelsPR.GetParty()
	if party == nil {
		return fmt.Errorf("party not available")
	}
	// Validate member count (max 4)
	if len(members) > 4 {
		return fmt.Errorf("party cannot have more than 4 members")
	}
	// Clear existing members
	for i := range party.Members {
		party.Members[i] = modelsPR.EmptyPartyMember
	}
	// Set new members by ID
	for i, charID := range members {
		if err := party.SetMemberByID(i, charID); err != nil {
			return fmt.Errorf("failed to set party member %d: %w", i, err)
		}
	}
	return nil
}

// Magic functions

func (b *Bindings) hasMagic(charID, spellID int) (bool, error) {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return false, fmt.Errorf("character with ID %d not found", charID)
	}
	// Get current character data from save
	ch, err := b.api.GetCharacter(context.Background(), char.Name)
	if err != nil {
		return false, fmt.Errorf("failed to get character: %w", err)
	}
	// Check if spell exists in character's spell list
	if spell, ok := ch.SpellsByID[spellID]; ok && spell != nil {
		// Spell exists - check if learned (Value > 0 typically indicates learned)
		return spell.Value > 0, nil
	}
	return false, nil
}

func (b *Bindings) learnMagic(charID, spellID int) error {
	// Get character model by ID
	char := modelsPR.GetCharacterByID(charID)
	if char == nil {
		return fmt.Errorf("character with ID %d not found", charID)
	}
	// Get current character data from save
	ch, err := b.api.GetCharacter(context.Background(), char.Name)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}
	// Check if spell exists in character's spell list
	if spell, ok := ch.SpellsByID[spellID]; ok && spell != nil {
		// Mark spell as learned by setting a proficiency value
		// Value 1 indicates learned (higher values indicate proficiency level)
		if spell.Value == 0 {
			spell.Value = 1
		}
		// Save back
		return b.api.SetCharacter(context.Background(), char.Name, ch)
	}
	return fmt.Errorf("spell ID %d not available for character", spellID)
}

// Utility functions

func (b *Bindings) print(message string) {
	// Print to console/log
	println(message)
}

func (b *Bindings) log(message string) {
	// Log message
	println("[Script]", message)
}
