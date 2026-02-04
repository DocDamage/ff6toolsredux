// Package settings provides persistent configuration management for the save editor.
//
// Settings are stored as JSON and include user preferences for:
//   - UI theme, language, and appearance
//   - Auto-save and backup options
//   - Cloud sync configuration
//   - Window size and position
//   - Recent files list
//   - Keyboard shortcuts
//   - Plugin and scripting preferences
//
// Usage:
//
//	manager := settings.NewManager("settings.json")
//	if err := manager.Load(); err != nil {
//	    // Use defaults
//	}
//	settings := manager.Get()
//	settings.Theme = "dark"
//	manager.Set(settings)
//	manager.Save()
//
// Legacy Migration:
//
// The package includes MigrateLegacyConfig() to convert old ff6editor.config
// files to the new unified settings format.
package settings
