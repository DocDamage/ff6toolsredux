// Package pr provides functionality for loading and saving Final Fantasy VI Pixel Remaster save files.
//
// The PR (Pixel Remaster) format is a JSON-based save file format used by the
// Final Fantasy VI Pixel Remaster edition on PC and mobile platforms.
//
// Basic usage:
//
//	p := pr.New()
//	if err := p.Load("save.json", 0); err != nil {
//	    log.Fatal(err)
//	}
//	// Modify save data...
//	if err := p.Save(0, "save.json", 0); err != nil {
//	    log.Fatal(err)
//	}
//
// The package handles:
//   - Character data (stats, equipment, magic, abilities)
//   - Inventory management
//   - Party composition
//   - Map position and transportation
//   - Esper and skill data
//   - Game flags and progression
//
// File Organization:
//   - loader.go: Core loading functionality
//   - loader_characters.go: Character-specific loading
//   - loader_inventory.go: Inventory loading
//   - loader_map.go: Map and transportation loading
//   - loader_misc.go: Espers, stats, cheats
//   - loader_helpers.go: Helper functions
//   - saver.go: Save file writing
package pr
