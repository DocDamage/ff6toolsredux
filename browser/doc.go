// Package browser provides web browser integration and update checking functionality.
//
// The browser package handles:
//   - Opening URLs in the system's default web browser
//   - Checking for application updates from remote sources
//   - Version comparison and update notification
//
// This package is used for:
//   - Opening help documentation in the user's browser
//   - Checking for new versions of the save editor
//   - Opening marketplace links
//
// Example usage:
//
//	// Open a URL in the default browser
//	browser.OpenURL("https://example.com/help")
//
//	// Check for updates
//	updateChecker := browser.NewUpdateChecker()
//	version, err := updateChecker.CheckForUpdate()
//	if err != nil {
//	    log.Printf("Update check failed: %v", err)
//	}
package browser
