// Package forms provides UI form implementations.
//
// The forms package contains:
//   - Main editor forms
//   - Dialog implementations
//   - Input components
//   - Selection widgets
//
// Subpackages:
//   - dialogs: Modal dialogs
//   - editors: Character, inventory, map editors
//   - inputs: Specialized input widgets
//   - selections: Selection lists and grids
//
// Example usage:
//
//	// Create character editor form
//	charForm := forms.NewCharacterEditor()
//	charForm.LoadCharacter(char)
//
//	// Show modal dialog
//	dialog := forms.NewBackupManagerDialog()
//	dialog.Show()
package forms
