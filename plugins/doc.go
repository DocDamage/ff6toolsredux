// Package plugins provides a plugin system for extending the save editor functionality.
//
// The plugin system supports:
//   - Lua scripting for custom modifications
//   - Sandboxed execution for security
//   - Hot-reload during development
//   - Permission-based access control
//   - Plugin marketplace integration
//
// Plugin Types:
//
//	Lua Scripts  - Text-based scripts for simple modifications
//	Go Plugins   - Compiled plugins for complex functionality
//
// Security:
//
// Plugins run in a sandboxed environment with configurable permissions:
//   - ReadSave: Access save data
//   - WriteSave: Modify save data
//   - FileSystem: Access files
//   - Network: Network access
//   - UIDisplay: Show UI dialogs
//
// Usage:
//
//	api := plugins.NewAPIImpl(prData, []string{
//	    plugins.CommonPermissions.ReadSave,
//	})
//	manager := plugins.NewManager("plugins/", api)
//	manager.Start(ctx)
package plugins
