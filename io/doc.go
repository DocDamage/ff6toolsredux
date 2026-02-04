// Package io provides input/output operations for save files and related data.
//
// The io package and its subpackages handle:
//   - Reading and writing save files in various formats
//   - Backup management
//   - Configuration file I/O
//   - JSON import/export
//   - Sprite and animation data
//   - Palette editing
//
// Subpackages:
//   - backup: Save file backup management
//   - config: Application configuration I/O
//   - file: Low-level file operations
//   - json: JSON format import/export
//   - presets: Party preset I/O
//   - pr: Pixel Remastered save format
//   - templates: Template I/O
//   - validation: Save file validation
//
// Example usage:
//
//	// Load a save file
//	data, trimmed, err := file.LoadFile("save.sav", global.PR)
//	if err != nil {
//	    return err
//	}
//
//	// Parse the save data
//	save, err := pr.Load(data)
//	if err != nil {
//	    return err
//	}
package io
