// Package cli provides command-line interface functionality.
//
// The CLI supports running the save editor in headless mode for:
//   - Batch processing save files
//   - Exporting save data to various formats
//   - Importing data into saves
//   - Running scripts on save files
//   - Validating save file integrity
//   - Creating backups
//
// Commands:
//
//	combat-pack  - Combat Depth Pack helpers (FULLY IMPLEMENTED)
//	edit         - Edit character stats (EXPERIMENTAL)
//	export       - Export save data (EXPERIMENTAL)
//	import       - Import data to save (EXPERIMENTAL)
//	batch        - Batch process saves (EXPERIMENTAL)
//	script       - Run Lua script (EXPERIMENTAL)
//	validate     - Validate save file (EXPERIMENTAL)
//	backup       - Create backup (EXPERIMENTAL)
//
// Usage:
//
//	ffvi_editor combat-pack --mode smoke --file save.json
//
// Experimental Commands:
//
// Commands marked as EXPERIMENTAL are defined in the CLI structure but
// not yet fully implemented. They will return an error indicating their
// experimental status. The combat-pack command is fully functional and
// ready for production use.
//
// For the graphical interface, run ffvi_editor without any arguments.
package cli
