// Package editors provides form-based editor implementations.
//
// The editors package contains:
//   - Character editor (stats, equipment, magic)
//   - Inventory editor
//   - Esper editor
//   - Map data editor
//   - Veldt editor
//   - Command editor
//
// Each editor provides a complete UI for editing
// specific aspects of the save file.
//
// Example usage:
//
//	// Create character editor
//	editor := editors.NewCharacterEditor()
//	editor.SetSave(save)
//
//	// Create inventory editor
//	invEditor := editors.NewInventoryEditor()
//	invEditor.SetInventory(save.Inventory)
package editors
