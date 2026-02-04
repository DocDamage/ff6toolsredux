// Package scripting provides Lua scripting capabilities for the save editor.
//
// The scripting system allows users to automate save file modifications
// using the Lua programming language.
//
// API Bindings:
//
// Scripts have access to editor functions through the editor global:
//
//	editor.getCharacter(name)      - Get character data
//	editor.setCharacter(name, ch)  - Update character
//	editor.getInventory()          - Get inventory
//	editor.log(level, message)     - Logging
//
// Script Execution:
//
// Scripts can be run standalone or with a save file loaded:
//
//	vm := scripting.NewVM()
//	result, err := vm.Execute(ctx, script)
//
// Combat Depth Pack:
//
// The package includes pre-built scripts for the Combat Depth Pack:
//   - Encounter tuning
//   - Boss remixing
//   - Companion director
//
// Note: Core save editing bindings are partially implemented.
// See TECHNICAL_DEBT_AUDIT_COMPLETE.md for details.
package scripting
