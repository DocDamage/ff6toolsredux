package ui

import (
	"context"
	"fmt"
	"time"

	"ffvi_editor/achievements"
	"ffvi_editor/cloud"
	"ffvi_editor/global"
	"ffvi_editor/io/config"
	"ffvi_editor/plugins"

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
		g.startAutoSave(time.Duration(s.AutoSaveDelay) * time.Second)
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
		// Store font size in global config for use by editors
		config.SetFontSize(s.EditorFontSize)
	}
}

// startAutoSave starts the auto-save timer with the specified interval
func (g *gui) startAutoSave(interval time.Duration) {
	if g.autoSaveTicker != nil {
		g.stopAutoSave <- true
		g.autoSaveTicker.Stop()
	}

	g.autoSaveTicker = time.NewTicker(interval)
	g.stopAutoSave = make(chan bool)

	go func() {
		for {
			select {
			case <-g.autoSaveTicker.C:
				g.performAutoSave()
			case <-g.stopAutoSave:
				return
			}
		}
	}()
}

// stopAutoSaveTimer stops the auto-save timer
func (g *gui) stopAutoSaveTimer() {
	if g.autoSaveTicker != nil {
		g.autoSaveTicker.Stop()
		if g.stopAutoSave != nil {
			select {
			case g.stopAutoSave <- true:
			default:
			}
		}
		g.autoSaveTicker = nil
	}
}

// performAutoSave performs the auto-save operation
func (g *gui) performAutoSave() {
	if g.pr == nil {
		return
	}

	// Create auto-save filename
	autoSavePath := global.PWD + "/autosave.json"
	if err := g.pr.Save(0, autoSavePath, 0); err != nil {
		// Silently log error - don't disturb user with dialogs for auto-save
		fmt.Printf("Auto-save failed: %v\n", err)
	} else {
		fmt.Printf("Auto-saved to: %s\n", autoSavePath)
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
	return g.settingsManager.Save()
}

// initializeCloudSync initializes cloud synchronization if enabled
func (g *gui) initializeCloudSync() {
	s := g.settingsManager.Get()
	if s == nil || !s.CloudEnabled {
		return
	}

	g.cloudManager = cloud.New()

	// Register cloud providers based on settings
	if s.DropboxEnabled {
		dropboxProvider := cloud.NewDropboxProvider(s.DropboxAppKey, s.DropboxAppSecret)
		if err := g.cloudManager.RegisterProvider(dropboxProvider); err != nil {
			fmt.Printf("Failed to register Dropbox provider: %v\n", err)
		}
	}

	if s.GoogleDriveClientID != "" {
		gdriveProvider := cloud.NewGoogleDriveProvider(s.GoogleDriveClientID, s.GoogleDriveClientSecret)
		if err := g.cloudManager.RegisterProvider(gdriveProvider); err != nil {
			fmt.Printf("Failed to register Google Drive provider: %v\n", err)
		}
	}

	fmt.Println("Cloud sync initialized")
}

// initializePlugins initializes the plugin system if enabled
func (g *gui) initializePlugins() {
	s := g.settingsManager.Get()
	if s == nil || !s.EnablePlugins {
		return
	}

	// Create a basic plugin API - plugins won't have access to save data
	// until a file is loaded, but they can be initialized
	api := plugins.NewAPIImpl(nil, []string{
		plugins.CommonPermissions.ReadSave,
		plugins.CommonPermissions.UIDisplay,
	})

	// Set up logging
	api.SetLogger(func(level, msg string) {
		fmt.Printf("[Plugin %s] %s\n", level, msg)
	})

	g.pluginManager = plugins.NewManager(global.PWD+"/plugins", api)

	// Start the plugin manager
	ctx := context.Background()
	if err := g.pluginManager.Start(ctx); err != nil {
		fmt.Printf("Failed to start plugin manager: %v\n", err)
		return
	}

	fmt.Println("Plugins initialized")
}

// cleanup performs cleanup before application exit
func (g *gui) cleanup() {
	g.stopAutoSaveTimer()

	if g.cloudManager != nil {
		// Perform final cloud sync if needed
	}

	if g.pluginManager != nil {
		ctx := context.Background()
		_ = g.pluginManager.Stop(ctx)
	}

	// Save settings before exit
	_ = g.saveSettings()
}
