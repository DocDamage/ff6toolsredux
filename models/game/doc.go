// Package game provides game-specific data structures.
//
// The game package contains:
//   - Rage database (enemy abilities)
//   - Sketch database
//   - Magic catalog
//   - Esper growth data
//   - Monster data
//
// These structures represent game data that is not part of
// the save file but is needed for editing.
//
// Example usage:
//
//	// Get rage abilities for an enemy
//	rages := game.GetRageAbilities(enemyID)
//
//	// Get sketch result
//	sketch := game.GetSketchResult(enemyID)
//
//	// Get magic info
//	magic := game.GetMagicInfo(spellID)
package game
