package forms

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/ui/shortcuts"
)

// CommandPaletteItem represents a single command in the palette
type CommandPaletteItem struct {
	ID          string
	Name        string
	Description string
	Category    string
	Action      func()
	Shortcut    string
}

// CommandPalette provides fuzzy-searchable command access
type CommandPalette struct {
	window         fyne.Window
	commands       []*CommandPaletteItem
	searchEntry    *widget.Entry
	resultList     *widget.List
	keyMap         *shortcuts.KeyMap
	recentCommands []string
	maxRecent      int
}

// NewCommandPalette creates a new command palette
func NewCommandPalette(window fyne.Window, keyMap *shortcuts.KeyMap) *CommandPalette {
	cp := &CommandPalette{
		window:         window,
		commands:       make([]*CommandPaletteItem, 0),
		keyMap:         keyMap,
		recentCommands: make([]string, 0),
		maxRecent:      5,
	}

	cp.registerDefaultCommands()
	return cp
}

// registerDefaultCommands registers built-in commands
func (cp *CommandPalette) registerDefaultCommands() {
	// File commands
	cp.AddCommand(&CommandPaletteItem{
		ID:          "open",
		Name:        "Open Save File",
		Description: "Load a save file from disk",
		Category:    "File",
		Shortcut:    "Ctrl+O",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "save",
		Name:        "Save Changes",
		Description: "Save current modifications",
		Category:    "File",
		Shortcut:    "Ctrl+S",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "save_as",
		Name:        "Save As",
		Description: "Save with a new filename",
		Category:    "File",
		Shortcut:    "Ctrl+Shift+S",
	})

	// Edit commands
	cp.AddCommand(&CommandPaletteItem{
		ID:          "undo",
		Name:        "Undo",
		Description: "Undo last action",
		Category:    "Edit",
		Shortcut:    "Ctrl+Z",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "redo",
		Name:        "Redo",
		Description: "Redo last undone action",
		Category:    "Edit",
		Shortcut:    "Ctrl+Y",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "search",
		Name:        "Search",
		Description: "Search save data",
		Category:    "Edit",
		Shortcut:    "Ctrl+F",
	})

	// Tool commands
	cp.AddCommand(&CommandPaletteItem{
		ID:          "backup",
		Name:        "Manage Backups",
		Description: "Open backup manager",
		Category:    "Tools",
		Shortcut:    "Ctrl+B",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "batch_max_stats",
		Name:        "Batch: Max All Stats",
		Description: "Maximize stats for all characters",
		Category:    "Tools",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "batch_max_level",
		Name:        "Batch: Max All Levels",
		Description: "Set all characters to level 99",
		Category:    "Tools",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "batch_learn_magic",
		Name:        "Batch: Learn All Magic",
		Description: "All characters learn all spells",
		Category:    "Tools",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "templates",
		Name:        "Template Manager",
		Description: "Manage character templates",
		Category:    "Tools",
		Shortcut:    "Ctrl+Shift+P",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "presets",
		Name:        "Party Presets",
		Description: "Manage party configurations",
		Category:    "Tools",
	})

	// View commands
	cp.AddCommand(&CommandPaletteItem{
		ID:          "toggle_theme",
		Name:        "Toggle Theme",
		Description: "Switch between light and dark themes",
		Category:    "View",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "settings",
		Name:        "Settings",
		Description: "Open application settings",
		Category:    "View",
		Shortcut:    "Ctrl+,",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "help",
		Name:        "Help",
		Description: "View help documentation",
		Category:    "View",
		Shortcut:    "F1",
	})

	cp.AddCommand(&CommandPaletteItem{
		ID:          "shortcuts",
		Name:        "Keyboard Shortcuts",
		Description: "View all keyboard shortcuts",
		Category:    "View",
	})
}

// AddCommand registers a new command
func (cp *CommandPalette) AddCommand(cmd *CommandPaletteItem) {
	if cmd != nil {
		cp.commands = append(cp.commands, cmd)
	}
}

// Show displays the command palette dialog
func (cp *CommandPalette) Show() {
	// Build search entry
	cp.searchEntry = widget.NewEntry()
	cp.searchEntry.SetPlaceHolder("Type command name or description...")
	cp.searchEntry.OnChanged = func(s string) {
		cp.updateResults(s)
	}

	// Build results list
	filteredCommands := cp.commands
	cp.resultList = widget.NewList(
		func() int {
			return len(filteredCommands)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Command")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(filteredCommands) {
				cmd := filteredCommands[id]
				label := obj.(*widget.Label)
				shortcutStr := ""
				if cmd.Shortcut != "" {
					shortcutStr = fmt.Sprintf(" (%s)", cmd.Shortcut)
				}
				label.SetText(fmt.Sprintf("%s%s\n%s", cmd.Name, shortcutStr, cmd.Description))
			}
		},
	)

	cp.resultList.OnSelected = func(id widget.ListItemID) {
		if id < len(filteredCommands) {
			cmd := filteredCommands[id]
			cp.executeCommand(cmd)
			cp.addToRecent(cmd.ID)
		}
	}

	// Build layout
	content := container.NewVBox(
		widget.NewCard(
			"Command Palette",
			"Type to search (Ctrl+Shift+P to close)",
			container.NewVBox(
				cp.searchEntry,
				cp.resultList,
			),
		),
	)

	dlg := dialog.NewCustom(
		"Command Palette",
		"Close",
		content,
		cp.window,
	)

	dlg.Resize(fyne.NewSize(600, 400))
	dlg.Show()
}

// updateResults filters commands based on search input
func (cp *CommandPalette) updateResults(query string) {
	query = strings.ToLower(strings.TrimSpace(query))

	// If empty, show recent commands first
	if query == "" {
		// Show recent commands + all others
		cp.resultList.Refresh()
		return
	}

	// Filter commands
	filtered := make([]*CommandPaletteItem, 0)
	for _, cmd := range cp.commands {
		if fuzzyMatch(cmd.Name, query) || fuzzyMatch(cmd.Description, query) {
			filtered = append(filtered, cmd)
		}
	}

	cp.resultList.Refresh()
}

// executeCommand runs a command's action
func (cp *CommandPalette) executeCommand(cmd *CommandPaletteItem) {
	if cmd.Action != nil {
		cmd.Action()
	}
	// Could add logging here for command tracking
}

// addToRecent adds a command to recent history
func (cp *CommandPalette) addToRecent(id string) {
	// Remove if already in list
	for i, item := range cp.recentCommands {
		if item == id {
			cp.recentCommands = append(cp.recentCommands[:i], cp.recentCommands[i+1:]...)
			break
		}
	}

	// Add to front
	cp.recentCommands = append([]string{id}, cp.recentCommands...)

	// Trim to max size
	if len(cp.recentCommands) > cp.maxRecent {
		cp.recentCommands = cp.recentCommands[:cp.maxRecent]
	}
}

// GetRecentCommands returns the recent command list
func (cp *CommandPalette) GetRecentCommands() []*CommandPaletteItem {
	result := make([]*CommandPaletteItem, 0)
	for _, id := range cp.recentCommands {
		for _, cmd := range cp.commands {
			if cmd.ID == id {
				result = append(result, cmd)
				break
			}
		}
	}
	return result
}

// Helper function for fuzzy matching
func fuzzyMatch(text, pattern string) bool {
	text = strings.ToLower(text)
	patternIdx := 0
	for _, char := range text {
		if patternIdx < len(pattern) && char == rune(pattern[patternIdx]) {
			patternIdx++
		}
	}
	return patternIdx == len(pattern)
}
