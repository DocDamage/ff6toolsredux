// Package consts provides game constants and enumerations.
//
// The consts package contains:
//   - Item IDs and names
//   - Spell IDs and names
//   - Character and job constants
//   - Map and location data
//   - Experience tables
//
// Subpackages:
//   - pr: Pixel Remastered specific constants
//
// Example usage:
//
//	// Get item name
//	name := consts.ItemName(consts.ItemPotion)
//
//	// Get spell power
//	power := consts.SpellPower(consts.SpellFire)
//
//	// Check job abilities
//	if consts.IsJobWithBushido(jobID) {
//	    // Handle Bushido abilities
//	}
package consts
