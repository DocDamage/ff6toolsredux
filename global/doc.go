// Package global provides application-wide shared state and constants.
//
// The global package contains:
//   - Working directory path (PWD)
//   - Application constants
//   - Shared state that needs to be accessible across packages
//
// ⚠️ Warning: This package uses init() for initialization.
// Consider migrating to explicit initialization in the future.
//
// Example usage:
//
//	// Access the working directory
//	savePath := filepath.Join(global.PWD, "saves", "game.sav")
//
//	// Check save file type
//	if global.IsPR(saveType) {
//	    // Handle Pixel Remastered format
//	}
package global
