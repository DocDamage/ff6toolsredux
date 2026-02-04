// Package json provides JSON import and export functionality for save files.
//
// The json package handles:
//   - Exporting save data to JSON format
//   - Importing save data from JSON format
//   - Data transformation between save formats
//
// This allows users to:
//   - Backup saves in a human-readable format
//   - Edit saves with external tools
//   - Share save data between different versions
//
// Example usage:
//
//	// Export save to JSON
//	jsonData, err := json.Export(save)
//	if err != nil {
//	    return err
//	}
//
//	// Import from JSON
//	save, err := json.Import(jsonData)
//	if err != nil {
//	    return err
//	}
package json
