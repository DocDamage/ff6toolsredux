// Package achievements provides achievement tracking.
//
// The achievements package handles:
//   - Achievement definitions
//   - Progress tracking
//   - Unlock notifications
//   - Statistics collection
//
// Achievement types:
//   - Save-based (edit specific values)
//   - Usage-based (open X files)
//   - Discovery-based (find hidden features)
//
// Example usage:
//
//	// Track achievement progress
//	achievements.Track("master_editor", progress)
//
//	// Unlock achievement
//	achievements.Unlock("first_save")
//
//	// Get achievement status
//	status := achievements.GetStatus("speed_runner")
package achievements
