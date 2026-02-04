// Package ui provides the graphical user interface for the Final Fantasy VI Save Editor.
//
// The UI is built using the Fyne toolkit (fyne.io/fyne) and provides:
//   - File I/O dialogs for loading and saving
//   - Character editors for stats, equipment, and abilities
//   - Inventory management interface
//   - Party composition editor
//   - Map teleportation tools
//   - Settings and preferences
//
// Main Components:
//
//	gui - Main window and application container
//	forms/editors - Data editing widgets
//	forms/selections - Navigation and selection widgets
//	forms/inputs - Custom input components
//
// Initialization:
//
// The UI is initialized through the gui struct which manages settings,
// achievements, cloud sync, and plugin systems during startup.
//
//	gui := ui.New()
//	gui.Load()
//	gui.Run()
package ui
