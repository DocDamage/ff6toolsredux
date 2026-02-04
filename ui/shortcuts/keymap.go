package shortcuts

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
)

// Shortcut represents a keyboard shortcut binding
type Shortcut struct {
	ID           string
	Description  string
	Keys         fyne.KeyName
	Modifiers    fyne.KeyModifier
	Action       func() // Function to call when shortcut is triggered
	Customizable bool
}

// KeyMap holds all registered keyboard shortcuts
type KeyMap struct {
	shortcuts map[string]*Shortcut
}

// NewKeyMap creates a new keyboard shortcut map with defaults
func NewKeyMap() *KeyMap {
	km := &KeyMap{
		shortcuts: make(map[string]*Shortcut),
	}

	km.registerDefaultShortcuts()
	return km
}

// registerDefaultShortcuts registers all default shortcuts
func (km *KeyMap) registerDefaultShortcuts() {
	// File Operations
	km.shortcuts["open"] = &Shortcut{
		ID:           "open",
		Description:  "Open save file",
		Keys:         fyne.KeyO,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}

	km.shortcuts["save"] = &Shortcut{
		ID:           "save",
		Description:  "Save changes",
		Keys:         fyne.KeyS,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}

	km.shortcuts["save_as"] = &Shortcut{
		ID:           "save_as",
		Description:  "Save as",
		Keys:         fyne.KeyS,
		Modifiers:    fyne.KeyModifierControl | fyne.KeyModifierShift,
		Customizable: true,
	}

	// Editing
	km.shortcuts["undo"] = &Shortcut{
		ID:           "undo",
		Description:  "Undo",
		Keys:         fyne.KeyZ,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}

	km.shortcuts["redo"] = &Shortcut{
		ID:           "redo",
		Description:  "Redo",
		Keys:         fyne.KeyY,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}

	km.shortcuts["search"] = &Shortcut{
		ID:           "search",
		Description:  "Search",
		Keys:         fyne.KeyF,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}

	// Tools
	km.shortcuts["backup"] = &Shortcut{
		ID:           "backup",
		Description:  "Backup manager",
		Keys:         fyne.KeyB,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}

	km.shortcuts["compare"] = &Shortcut{
		ID:           "compare",
		Description:  "Compare saves",
		Keys:         fyne.KeyC,
		Modifiers:    fyne.KeyModifierControl | fyne.KeyModifierShift,
		Customizable: true,
	}

	km.shortcuts["batch"] = &Shortcut{
		ID:           "batch",
		Description:  "Batch operations",
		Keys:         fyne.KeyB,
		Modifiers:    fyne.KeyModifierControl | fyne.KeyModifierShift,
		Customizable: true,
	}

	km.shortcuts["templates"] = &Shortcut{
		ID:           "templates",
		Description:  "Template manager",
		Keys:         fyne.KeyP,
		Modifiers:    fyne.KeyModifierControl | fyne.KeyModifierShift,
		Customizable: true,
	}

	km.shortcuts["command_palette"] = &Shortcut{
		ID:           "command_palette",
		Description:  "Command palette",
		Keys:         fyne.KeyP,
		Modifiers:    fyne.KeyModifierControl | fyne.KeyModifierShift,
		Customizable: false,
	}

	// Views
	km.shortcuts["help"] = &Shortcut{
		ID:           "help",
		Description:  "Help",
		Keys:         fyne.KeyF1,
		Modifiers:    0,
		Customizable: true,
	}

	km.shortcuts["settings"] = &Shortcut{
		ID:           "settings",
		Description:  "Settings",
		Keys:         fyne.KeyComma,
		Modifiers:    fyne.KeyModifierControl,
		Customizable: true,
	}
}

// GetShortcut retrieves a shortcut by ID
func (km *KeyMap) GetShortcut(id string) *Shortcut {
	return km.shortcuts[id]
}

// SetShortcut updates a shortcut binding
func (km *KeyMap) SetShortcut(id string, keys fyne.KeyName, modifiers fyne.KeyModifier) error {
	shortcut, exists := km.shortcuts[id]
	if !exists {
		return fmt.Errorf("shortcut not found: %s", id)
	}

	if !shortcut.Customizable {
		return fmt.Errorf("shortcut cannot be customized: %s", id)
	}

	// Check for conflicts
	for _, sc := range km.shortcuts {
		if sc.ID != id && sc.Keys == keys && sc.Modifiers == modifiers {
			return fmt.Errorf("shortcut already in use by: %s", sc.ID)
		}
	}

	shortcut.Keys = keys
	shortcut.Modifiers = modifiers
	return nil
}

// ResetToDefaults resets all customizable shortcuts to defaults
func (km *KeyMap) ResetToDefaults() {
	km.shortcuts = make(map[string]*Shortcut)
	km.registerDefaultShortcuts()
}

// GetAllShortcuts returns all registered shortcuts
func (km *KeyMap) GetAllShortcuts() []*Shortcut {
	result := make([]*Shortcut, 0, len(km.shortcuts))
	for _, sc := range km.shortcuts {
		result = append(result, sc)
	}
	return result
}

// GetShortcutsByCategory returns shortcuts for a specific category
func (km *KeyMap) GetShortcutsByCategory(category string) []*Shortcut {
	result := make([]*Shortcut, 0)

	var categoryPrefix string
	switch category {
	case "File":
		categoryPrefix = "save"
	case "Edit":
		categoryPrefix = "undo"
	case "Tools":
		categoryPrefix = "backup"
	case "View":
		categoryPrefix = "help"
	}

	for _, sc := range km.shortcuts {
		if strings.Contains(sc.ID, categoryPrefix) {
			result = append(result, sc)
		}
	}

	return result
}

// FormatKeyCombo returns a human-readable key combination
func (km *KeyMap) FormatKeyCombo(shortcut *Shortcut) string {
	var parts []string

	if shortcut.Modifiers&fyne.KeyModifierControl != 0 {
		parts = append(parts, "Ctrl")
	}
	if shortcut.Modifiers&fyne.KeyModifierShift != 0 {
		parts = append(parts, "Shift")
	}
	if shortcut.Modifiers&fyne.KeyModifierAlt != 0 {
		parts = append(parts, "Alt")
	}
	if shortcut.Modifiers&fyne.KeyModifierSuper != 0 {
		parts = append(parts, "Super")
	}

	// Format key name
	keyName := string(shortcut.Keys)
	if keyName != "" {
		parts = append(parts, keyName)
	}

	return strings.Join(parts, "+")
}

// ListAllShortcuts returns a formatted list of all shortcuts
func (km *KeyMap) ListAllShortcuts() string {
	var output strings.Builder
	output.WriteString("Keyboard Shortcuts:\n")
	output.WriteString("==================\n\n")

	categories := []string{"File", "Edit", "Tools", "View"}
	for _, category := range categories {
		output.WriteString(fmt.Sprintf("%s Operations:\n", category))
		shortcuts := km.GetShortcutsByCategory(category)
		for _, sc := range shortcuts {
			keyCombo := km.FormatKeyCombo(sc)
			output.WriteString(fmt.Sprintf("  %-30s %s\n", sc.Description, keyCombo))
		}
		output.WriteString("\n")
	}

	return output.String()
}
