package scripting

import (
	"context"
	"ffvi_editor/plugins"
	"fmt"
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
	// Real implementation would create a Lua function that calls:
	// return b.api.GetCharacter(ctx, name)
	return b.vm.RegisterFunction("editor.getCharacter", func(name string) interface{} {
		ch, err := b.api.GetCharacter(ctx, name)
		if err != nil {
			return nil
		}
		return ch
	})
}

// BindSetCharacter binds the SetCharacter API function
func (b *Bindings) BindSetCharacter(ctx context.Context) error {
	// Real implementation would create a Lua function that calls:
	// return b.api.SetCharacter(ctx, name, ch)
	return b.vm.RegisterFunction("editor.setCharacter", func(name string, ch interface{}) interface{} {
		// Would deserialize ch and call b.api.SetCharacter
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
		return b.api.GetSetting(key)
	})
}

// BindSetSetting binds the SetSetting function
func (b *Bindings) BindSetSetting() error {
	return b.vm.RegisterFunction("editor.setSetting", func(key string, value interface{}) error {
		return b.api.SetSetting(key, value)
	})
}

// Character functions

func (b *Bindings) getCharacter(charID int) (interface{}, error) {
	// TODO: Extract character from PR save
	return nil, nil
}

func (b *Bindings) setCharacterLevel(charID, level int) error {
	// TODO: Set character level in PR save
	return nil
}

func (b *Bindings) setCharacterHP(charID, hp int) error {
	// TODO: Set character HP in PR save
	return nil
}

func (b *Bindings) setCharacterMP(charID, mp int) error {
	// TODO: Set character MP in PR save
	return nil
}

func (b *Bindings) setCharacterStat(charID int, stat string, value int) error {
	// TODO: Set character stat in PR save
	// stat can be: "vigor", "speed", "stamina", "magic", "defense", "magicdef"
	return nil
}

// Inventory functions

func (b *Bindings) getItemCount(itemID int) (int, error) {
	// TODO: Get item count from PR save
	return 0, nil
}

func (b *Bindings) setItemCount(itemID, count int) error {
	// TODO: Set item count in PR save
	return nil
}

func (b *Bindings) addItem(itemID, count int) error {
	// TODO: Add items to inventory
	return nil
}

// Party functions

func (b *Bindings) getPartyMembers() ([]int, error) {
	// TODO: Get party members from PR save
	return []int{}, nil
}

func (b *Bindings) setPartyMembers(members []int) error {
	// TODO: Set party members in PR save
	return nil
}

// Magic functions

func (b *Bindings) hasMagic(charID, spellID int) (bool, error) {
	// TODO: Check if character has magic
	return false, nil
}

func (b *Bindings) learnMagic(charID, spellID int) error {
	// TODO: Learn magic for character
	return nil
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
