// Package shortcuts provides keyboard shortcut management.
//
// The shortcuts package handles:
//   - Global keyboard shortcuts
//   - Shortcut configuration
//   - Key mapping
//   - Shortcut conflicts detection
//
// Default shortcuts:
//   - Ctrl+S: Save
//   - Ctrl+Z: Undo
//   - Ctrl+Y: Redo
//   - Ctrl+F: Find
//   - Ctrl+O: Open
//   - Ctrl+N: New
//
// Example usage:
//
//	// Register shortcut
//	shortcuts.Register("save", "Ctrl+S", func() {
//	    saveFile()
//	})
//
//	// Load custom keymap
//	keymap := shortcuts.LoadKeymap("custom.json")
//	shortcuts.ApplyKeymap(keymap)
package shortcuts
