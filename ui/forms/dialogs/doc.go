// Package dialogs provides modal dialog implementations.
//
// The dialogs package contains:
//   - Backup manager dialog
//   - Help dialog
//   - Preferences dialog
//   - Cloud settings dialog
//   - Validation dialogs
//
// These are reusable modal dialogs used throughout
// the application.
//
// Example usage:
//
//	// Show backup manager
//	dialog := dialogs.NewBackupManager()
//	dialog.Show()
//
//	// Show help
//	help := dialogs.NewHelpDialog("getting_started")
//	help.Show()
package dialogs
