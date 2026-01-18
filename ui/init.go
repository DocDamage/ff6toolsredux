package ui

import (
	"fmt"

	"ffvi_editor/achievements"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// initializeSettings loads settings and applies them to the GUI
func (g *gui) initializeSettings() error {
	// Load settings from file
	if err := g.settingsManager.Load(); err != nil {
		// If load fails, use defaults
		g.settingsManager.Reset()
	}

	// Apply settings to app
	g.applySettings()

	// Start auto-save if enabled
	s := g.settingsManager.Get()
	if s != nil && s.AutoSave {
		// TODO: Start auto-save timer
	}

	return nil
}

// applySettings applies loaded settings to the GUI
func (g *gui) applySettings() {
	// Get the settings object
	s := g.settingsManager.Get()
	if s == nil {
		return
	}

	// Apply theme - simplified for now
	// In production, would check s.Theme and apply accordingly

	// Apply window size if not maximized
	if s.WindowWidth > 0 && s.WindowHeight > 0 {
		g.window.Resize(fyne.NewSize(float32(s.WindowWidth), float32(s.WindowHeight)))
	}

	// Apply font size
	if s.EditorFontSize > 0 {
		// TODO: Apply font size to editor
	}
}

// initializeAchievements sets up achievement tracking
func (g *gui) initializeAchievements() {
	// Set unlock callback for notifications
	g.achievementTracker.SetUnlockCallback(func(achievement *achievements.Achievement) {
		g.showAchievementUnlock(achievement)
	})

	// Initialize some achievements as examples
	// In real usage, these would be triggered by actual events

	// Register progress tracking for editor achievements
	g.registerEditorAchievements()
}

// registerEditorAchievements registers callbacks for editor-related achievements
func (g *gui) registerEditorAchievements() {
	// First edit achievement
	// This would be triggered when user first edits a character

	// Save 10 achievements
	// This would be triggered each time user saves

	// Undo master
	// This would be triggered when user performs 50 undos

	// Template creator
	// This would be triggered when user creates a template

	// Batch operator
	// This would be triggered when user performs batch operations

	// Scripter
	// This would be triggered when user runs Lua scripts

	// Plugin user
	// This would be triggered when user enables plugins
}

// showAchievementUnlock displays an achievement unlock notification
func (g *gui) showAchievementUnlock(achievement *achievements.Achievement) {
	message := fmt.Sprintf(
		"🏆 %s\n\n%s\n\n+%d points",
		achievement.Name,
		achievement.Description,
		achievement.Points,
	)

	dialog.ShowInformation(
		"Achievement Unlocked!",
		message,
		g.window,
	)
}

// saveSettings persists current settings to file
func (g *gui) saveSettings() error {
	// TODO: Implement settings persistence
	return nil
}

// initializeCloudSync initializes cloud synchronization if enabled
func (g *gui) initializeCloudSync() {
	// TODO: Initialize cloud sync if enabled
}

// initializePlugins initializes the plugin system if enabled
func (g *gui) initializePlugins() {
	// TODO: Initialize plugins if enabled
}
