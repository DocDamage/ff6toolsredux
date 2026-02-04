// Package models provides data models for the Final Fantasy VI Save Editor.
//
// The models package contains:
//   - Core data structures for save file representation
//   - Validation logic for game data
//   - Constants and enums for game values
//   - Helper types for UI and editing operations
//
// Subpackages:
//   - batch: Batch operation definitions
//   - consts: Game constants (items, spells, etc.)
//   - game: Game-specific data structures
//   - pr: Pixel Remastered save format models
//   - search: Search indexing
//   - share: Share code generation
//   - speedrun: Speedrun configuration
//   - templates: Character templates
//
// Example usage:
//
//	// Create a new character model
//	char := &models.Character{
//	    Name: "Terra",
//	    Level: 50,
//	}
//
//	// Validate the model
//	if err := char.Validate(); err != nil {
//	    return err
//	}
package models
