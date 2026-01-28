// plugin_manager_sync.go - Marketplace synchronization methods for PluginManagerDialog

package forms

import (
	"context"
	"fmt"
	"time"

	"ffvi_editor/marketplace"

	"fyne.io/fyne/v2/dialog"
)

// RefreshInstalledPlugins refreshes the installed plugins list after marketplace installation
func (p *PluginManagerDialog) RefreshInstalledPlugins() error {
	// This would typically re-render the installed plugins tab
	// For now, this is a placeholder for UI refresh logic
	return nil
}

// CheckForUpdates checks if any installed plugins have updates available
func (p *PluginManagerDialog) CheckForUpdates() {
	// Get list of installed plugins from plugin manager
	installedPlugins := p.pluginManager.ListPlugins()

	// Check each plugin for updates
	var updates map[string]*marketplace.UpdateInfo = make(map[string]*marketplace.UpdateInfo)

	for _, plugin := range installedPlugins {
		// Create a mock remote plugin entry to check for updates
		// In production, this would query the actual marketplace registry
		updateInfo := &marketplace.UpdateInfo{
			PluginID:       plugin.ID,
			CurrentVersion: plugin.Version,
			LatestVersion:  "1.0.0", // TODO: Get from marketplace
			ReleaseNotes:   "",
			UpdatedAt:      time.Now(),
		}

		// Check if update is available by comparing versions
		if updateInfo.LatestVersion != updateInfo.CurrentVersion {
			updates[plugin.ID] = updateInfo
		}
	}

	// If updates available, show notification
	if len(updates) > 0 {
		p.showAvailableUpdates(updates)
	} else {
		dialog.ShowInformation("No Updates", "All plugins are up to date", p.window)
	}
}

// showAvailableUpdates displays a dialog showing available plugin updates
func (p *PluginManagerDialog) showAvailableUpdates(updates map[string]*marketplace.UpdateInfo) {
	var updateList string
	updateList = "The following plugins have updates available:\n\n"

	for pluginID, update := range updates {
		updateList += fmt.Sprintf("• %s: v%s → v%s\n",
			pluginID,
			update.CurrentVersion,
			update.LatestVersion,
		)
	}

	dialog.ShowConfirm(
		"Plugin Updates Available",
		updateList+"\nWould you like to update now?",
		func(confirmed bool) {
			if confirmed {
				// In production, iterate through updates and install them
				// For now, show a message
				dialog.ShowInformation("Updates", "Update functionality will be available soon", p.window)
			}
		},
		p.window,
	)
}

// SyncWithMarketplace synchronizes installed plugins with marketplace registry
func (p *PluginManagerDialog) SyncWithMarketplace(marketplace *marketplace.Client, registry *marketplace.Registry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get installed plugins from local manager
	installedPlugins := p.pluginManager.ListPlugins()

	// Sync each installation with registry
	for _, plugin := range installedPlugins {
		// Update registry with installation info
		err := registry.TrackInstallation(plugin.ID, plugin.Version)
		if err != nil {
			// Log but don't fail - registry sync is best-effort
			fmt.Printf("Warning: Failed to sync plugin %s: %v\n", plugin.ID, err)
		}

		// Check for updates in marketplace
		pluginDetails, err := marketplace.GetPluginDetails(ctx, plugin.ID)
		if err == nil && pluginDetails != nil && pluginDetails.Version != plugin.Version {
			// Update available - store in registry
			registry.UpdatePlugin(plugin.ID, pluginDetails.Version)
		}
	}

	return nil
}

// OfferUpdates displays available plugin updates to the user
func (p *PluginManagerDialog) OfferUpdates(updates map[string]*marketplace.UpdateInfo) {
	if len(updates) == 0 {
		return
	}

	p.showAvailableUpdates(updates)
}

// GetInstalledPluginCount returns the number of installed plugins
func (p *PluginManagerDialog) GetInstalledPluginCount() int {
	return len(p.pluginManager.ListPlugins())
}

// GetEnabledPluginCount returns the number of enabled plugins
func (p *PluginManagerDialog) GetEnabledPluginCount() int {
	count := 0
	for _, plugin := range p.pluginManager.ListPlugins() {
		if plugin.Enabled {
			count++
		}
	}
	return count
}

// IsPluginInstalled checks if a specific plugin is installed
func (p *PluginManagerDialog) IsPluginInstalled(pluginID string) bool {
	for _, plugin := range p.pluginManager.ListPlugins() {
		if plugin.ID == pluginID {
			return true
		}
	}
	return false
}

// GetInstalledVersion returns the version of an installed plugin, or empty string if not installed
func (p *PluginManagerDialog) GetInstalledVersion(pluginID string) string {
	for _, plugin := range p.pluginManager.ListPlugins() {
		if plugin.ID == pluginID {
			return plugin.Version
		}
	}
	return ""
}

// NotifyInstallation notifies the plugin manager of a new marketplace installation
func (p *PluginManagerDialog) NotifyInstallation(pluginID, version string) {
	// Trigger UI refresh to show newly installed plugin
	err := p.RefreshInstalledPlugins()
	if err != nil {
		fmt.Printf("Warning: Failed to refresh installed plugins: %v\n", err)
	}

	// Show success message
	dialog.ShowInformation(
		"Installation Complete",
		fmt.Sprintf("Plugin %s v%s has been installed and will be available after restart.",
			pluginID, version),
		p.window,
	)
}

// NotifyUninstallation notifies the plugin manager of a plugin removal
func (p *PluginManagerDialog) NotifyUninstallation(pluginID string) {
	// Trigger UI refresh to remove uninstalled plugin
	err := p.RefreshInstalledPlugins()
	if err != nil {
		fmt.Printf("Warning: Failed to refresh installed plugins: %v\n", err)
	}

	// Show success message
	dialog.ShowInformation(
		"Uninstallation Complete",
		fmt.Sprintf("Plugin %s has been uninstalled.",
			pluginID),
		p.window,
	)
}
