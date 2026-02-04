// Package pr provides data models for Final Fantasy VI save file components.
//
// This package defines the core data structures used throughout the save editor
// to represent game state including characters, inventory, party composition,
// and map data.
//
// Core Types:
//
//	Character - Represents a playable character with stats, equipment, and abilities
//	Party     - The current active party composition (up to 4 members)
//	Inventory - Item storage with normal and important item lists
//	MapData   - Player position and map-related information
//	Veldt     - Rage encounter tracking for Gau
//
// Singleton Access:
//
// Most model types are accessed through singleton functions:
//
//	characters := pr.GetCharacters()
//	party := pr.GetParty()
//	inventory := pr.GetInventory()
//
// Thread Safety:
//
// The models in this package are not thread-safe. Access should be synchronized
// when used from multiple goroutines.
package pr
