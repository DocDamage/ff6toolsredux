// marketplace_config.go - Marketplace configuration helpers

package config

import (
	"path/filepath"
	"time"
)

// GetMarketplaceSettings returns the marketplace configuration
func GetMarketplaceSettings() MarketplaceConfig {
	mu.RLock()
	defer mu.RUnlock()

	// Apply defaults if not configured
	settings := data.MarketplaceSettings
	if settings.RegistryURL == "" {
		settings.RegistryURL = "https://raw.githubusercontent.com/ff6-marketplace/registry/main"
	}
	if settings.CacheTTLHours == 0 {
		settings.CacheTTLHours = 24
	}
	if settings.UpdateCheckIntervalHours == 0 {
		settings.UpdateCheckIntervalHours = 12
	}
	if settings.MaxConcurrentDownloads == 0 {
		settings.MaxConcurrentDownloads = 3
	}
	if settings.CachePath == "" {
		settings.CachePath = filepath.Join(SaveDir(), "marketplace", "cache")
	}

	// Default to enabled if not explicitly disabled
	if !settings.Enabled && settings.RegistryURL != "" {
		settings.Enabled = true
	}

	return settings
}

// SetMarketplaceSettings updates marketplace configuration and persists to file
func SetMarketplaceSettings(settings MarketplaceConfig) error {
	mu.Lock()
	defer mu.Unlock()

	data.MarketplaceSettings = settings
	return save()
}

// UpdateMarketplaceRegistryURL updates the marketplace registry URL
func UpdateMarketplaceRegistryURL(url string) error {
	mu.Lock()
	defer mu.Unlock()

	data.MarketplaceSettings.RegistryURL = url
	return save()
}

// UpdateLastMarketplaceCheckTime records when marketplace was last checked for updates
func UpdateLastMarketplaceCheckTime() error {
	mu.Lock()
	defer mu.Unlock()

	data.MarketplaceSettings.LastUpdateCheckTime = time.Now().Format(time.RFC3339)
	return save()
}

// SetMarketplaceEnabled toggles the marketplace feature
func SetMarketplaceEnabled(enabled bool) error {
	mu.Lock()
	defer mu.Unlock()

	data.MarketplaceSettings.Enabled = enabled
	return save()
}

// SetAutoCheckUpdates configures automatic update checking
func SetAutoCheckUpdates(enabled bool, intervalHours int) error {
	mu.Lock()
	defer mu.Unlock()

	data.MarketplaceSettings.AutoCheckUpdates = enabled
	if intervalHours > 0 {
		data.MarketplaceSettings.UpdateCheckIntervalHours = intervalHours
	}
	return save()
}

// SetMarketplaceCacheTTL sets cache expiry time in hours
func SetMarketplaceCacheTTL(hours int) error {
	mu.Lock()
	defer mu.Unlock()

	if hours > 0 {
		data.MarketplaceSettings.CacheTTLHours = hours
	}
	return save()
}

// SetMaxConcurrentDownloads sets the maximum number of concurrent plugin downloads
func SetMaxConcurrentDownloads(max int) error {
	mu.Lock()
	defer mu.Unlock()

	if max > 0 {
		data.MarketplaceSettings.MaxConcurrentDownloads = max
	}
	return save()
}

// GetMarketplaceCachePath returns the path where marketplace cache is stored
func GetMarketplaceCachePath() string {
	settings := GetMarketplaceSettings()
	if settings.CachePath == "" {
		return filepath.Join(SaveDir(), "marketplace", "cache")
	}
	return settings.CachePath
}

// IsMarketplaceEnabled checks if marketplace is enabled
func IsMarketplaceEnabled() bool {
	settings := GetMarketplaceSettings()
	return settings.Enabled
}

// ShouldCheckForUpdates determines if marketplace should check for plugin updates
func ShouldCheckForUpdates() bool {
	settings := GetMarketplaceSettings()

	if !settings.AutoCheckUpdates {
		return false
	}

	// Check if enough time has passed since last check
	if settings.LastUpdateCheckTime != "" {
		lastCheck, err := time.Parse(time.RFC3339, settings.LastUpdateCheckTime)
		if err == nil {
			nextCheck := lastCheck.Add(time.Duration(settings.UpdateCheckIntervalHours) * time.Hour)
			if time.Now().Before(nextCheck) {
				return false
			}
		}
	}

	return true
}

// GetMarketplaceRegistryURL returns the configured registry URL
func GetMarketplaceRegistryURL() string {
	settings := GetMarketplaceSettings()
	return settings.RegistryURL
}

// ResetMarketplaceSettings resets marketplace configuration to defaults
func ResetMarketplaceSettings() error {
	mu.Lock()
	defer mu.Unlock()

	data.MarketplaceSettings = MarketplaceConfig{
		Enabled:                true,
		RegistryURL:            "https://raw.githubusercontent.com/ff6-marketplace/registry/main",
		CachePath:              filepath.Join(SaveDir(), "marketplace", "cache"),
		CacheTTLHours:          24,
		AutoCheckUpdates:       true,
		UpdateCheckIntervalHours: 12,
		MaxConcurrentDownloads: 3,
		EnableRatings:          true,
		DefaultCategory:        "All Categories",
	}
	return save()
}
