// Package presets provides party preset management.
//
// The presets package handles:
//   - Saving and loading party configurations
//   - Managing preset collections
//   - Sharing presets between saves
//
// Features:
//   - Named party presets
//   - Character, equipment, and ability snapshots
//   - Import/export functionality
//
// Example usage:
//
//	// Create party manager
//	pm := presets.NewPartyManager(presetDir)
//
//	// Save current party as preset
//	if err := pm.SavePreset("My Dream Team", party); err != nil {
//	    return err
//	}
//
//	// Load a preset
//	party, err := pm.LoadPreset("My Dream Team")
//	if err != nil {
//	    return err
//	}
package presets
