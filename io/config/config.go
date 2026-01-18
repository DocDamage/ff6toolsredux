package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"ffvi_editor/global"
)

const (
	file = "ff6editor.config"
)

var (
	data         d
	mu           sync.RWMutex // Thread-safe access to config
	lastWriteErr error        // Tracks last write error for debugging
)

type (
	d struct {
		WindowX             float32           `json:"width"`
		WindowY             float32           `json:"height"`
		SaveDir             string            `json:"dir"`
		AutoEnableCmd       bool              `json:"autoEnableCmd"`
		EnablePlayStation   bool              `json:"ps"`
		CloudSettings       CloudConfig       `json:"cloud"`
		MarketplaceSettings MarketplaceConfig `json:"marketplace"`
	}

	// CloudConfig contains all cloud backup settings
	CloudConfig struct {
		GoogleDriveEnabled      bool   `json:"googleDriveEnabled"`
		DropboxEnabled          bool   `json:"dropboxEnabled"`
		GoogleDriveClientID     string `json:"googleDriveClientID"`
		GoogleDriveClientSecret string `json:"googleDriveClientSecret"`
		DropboxAppKey           string `json:"dropboxAppKey"`
		DropboxAppSecret        string `json:"dropboxAppSecret"`
		AutoSync                bool   `json:"autoSync"`
		SyncIntervalMinutes     int    `json:"syncIntervalMinutes"`
		ConflictStrategy        string `json:"conflictStrategy"` // "local", "remote", "newest", "both"
		EncryptionEnabled       bool   `json:"encryptionEnabled"`
		BackupFolderPath        string `json:"backupFolderPath"`
		TemplatesFolderPath     string `json:"templatesFolderPath"`
		MaxRetries              int    `json:"maxRetries"`
		VerifyHashes            bool   `json:"verifyHashes"`
	}

	// MarketplaceConfig contains all marketplace settings
	MarketplaceConfig struct {
		Enabled                  bool   `json:"enabled"`
		RegistryURL              string `json:"registryURL"`
		CachePath                string `json:"cachePath"`
		CacheTTLHours            int    `json:"cacheTTLHours"`
		AutoCheckUpdates         bool   `json:"autoCheckUpdates"`
		UpdateCheckIntervalHours int    `json:"updateCheckIntervalHours"`
		MaxConcurrentDownloads   int    `json:"maxConcurrentDownloads"`
		EnableRatings            bool   `json:"enableRatings"`
		DefaultCategory          string `json:"defaultCategory"`
		LastUpdateCheckTime      string `json:"lastUpdateCheckTime"`
	}
)

// init loads the configuration from file with proper error handling
func init() {
	mu.Lock()
	defer mu.Unlock()

	configPath := filepath.Join(global.PWD, file)
	b, err := os.ReadFile(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			// Log error for debugging, but don't fail - use defaults
			lastWriteErr = fmt.Errorf("failed to read config file: %w", err)
		}
		// Set defaults
		data.WindowX = global.WindowWidth
		data.WindowY = global.WindowHeight
		return
	}

	// Unmarshal with error checking
	if err := json.Unmarshal(b, &data); err != nil {
		lastWriteErr = fmt.Errorf("failed to unmarshal config: %w", err)
		data.WindowX = global.WindowWidth
		data.WindowY = global.WindowHeight
		return
	}

	// Apply defaults for zero values
	if data.WindowX == 0 {
		data.WindowX = global.WindowWidth
	}
	if data.WindowY == 0 {
		data.WindowY = global.WindowHeight
	}
}

// WindowSize returns the saved window dimensions
func WindowSize() (x, y float32) {
	mu.RLock()
	defer mu.RUnlock()

	x = data.WindowX
	y = data.WindowY
	return
}

// SaveDir returns the configured save directory
func SaveDir() string {
	mu.RLock()
	defer mu.RUnlock()

	return data.SaveDir
}

// AutoEnableCmd returns whether auto-enable commands is enabled
func AutoEnableCmd() bool {
	mu.RLock()
	defer mu.RUnlock()

	return data.AutoEnableCmd
}

// EnablePlayStation returns whether PlayStation format support is enabled
func EnablePlayStation() bool {
	mu.RLock()
	defer mu.RUnlock()

	return data.EnablePlayStation
}

// SetWindowSize updates and persists the window size
func SetWindowSize(x, y float32) error {
	if x <= 0 || y <= 0 {
		return fmt.Errorf("invalid window size: %fx%f", x, y)
	}

	mu.Lock()
	data.WindowX = x
	data.WindowY = y
	mu.Unlock()

	return save()
}

// SetSaveDir updates and persists the save directory
func SetSaveDir(dir string) error {
	if dir == "" {
		return errors.New("save directory cannot be empty")
	}

	mu.Lock()
	data.SaveDir = dir
	mu.Unlock()

	return save()
}

// SetAutoEnableCmd updates and persists the auto-enable commands setting
func SetAutoEnableCmd(v bool) error {
	mu.Lock()
	data.AutoEnableCmd = v
	mu.Unlock()

	return save()
}

// SetEnablePlayStation updates and persists the PlayStation support setting
func SetEnablePlayStation(v bool) error {
	mu.Lock()
	data.EnablePlayStation = v
	mu.Unlock()

	return save()
}

// GetCloudSettings returns a copy of the cloud configuration
func GetCloudSettings() CloudConfig {
	mu.RLock()
	defer mu.RUnlock()

	return data.CloudSettings
}

// SetCloudSettings updates and persists the cloud configuration
func SetCloudSettings(cfg CloudConfig) error {
	if cfg.SyncIntervalMinutes < 1 {
		cfg.SyncIntervalMinutes = 30 // Default to 30 minutes
	}
	if cfg.MaxRetries < 1 {
		cfg.MaxRetries = 3 // Default to 3 retries
	}

	mu.Lock()
	data.CloudSettings = cfg
	mu.Unlock()

	return save()
}

// UpdateCloudProvider updates settings for a specific cloud provider
func UpdateCloudProvider(providerName string, enabled bool, clientID, clientSecret string) error {
	mu.Lock()
	defer mu.Unlock()

	switch providerName {
	case "google_drive":
		data.CloudSettings.GoogleDriveEnabled = enabled
		if clientID != "" {
			data.CloudSettings.GoogleDriveClientID = clientID
		}
		if clientSecret != "" {
			data.CloudSettings.GoogleDriveClientSecret = clientSecret
		}

	case "dropbox":
		data.CloudSettings.DropboxEnabled = enabled
		if clientID != "" {
			data.CloudSettings.DropboxAppKey = clientID
		}
		if clientSecret != "" {
			data.CloudSettings.DropboxAppSecret = clientSecret
		}

	default:
		return fmt.Errorf("unknown cloud provider: %s", providerName)
	}

	return save()
}

// save persists the configuration to disk with proper error handling
func save() error {
	mu.RLock()
	defer mu.RUnlock()

	// Apply defaults if not set
	if data.WindowX == 0 {
		data.WindowX = global.WindowWidth
	}
	if data.WindowY == 0 {
		data.WindowY = global.WindowHeight
	}

	// Marshal the configuration
	b, err := json.Marshal(&data)
	if err != nil {
		lastWriteErr = fmt.Errorf("failed to marshal config: %w", err)
		return lastWriteErr
	}

	if len(b) == 0 {
		lastWriteErr = errors.New("marshaled config is empty")
		return lastWriteErr
	}

	configPath := filepath.Join(global.PWD, file)
	if err := os.WriteFile(configPath, b, 0644); err != nil {
		lastWriteErr = fmt.Errorf("failed to write config file: %w", err)
		return lastWriteErr
	}

	lastWriteErr = nil
	return nil
}

// GetLastError returns the last error encountered during config operations
// Useful for debugging configuration issues
func GetLastError() error {
	mu.RLock()
	defer mu.RUnlock()

	return lastWriteErr
}

// ValidateConfig checks the current configuration for validity
func ValidateConfig() error {
	mu.RLock()
	defer mu.RUnlock()

	if data.WindowX <= 0 || data.WindowY <= 0 {
		return fmt.Errorf("invalid window size: %fx%f", data.WindowX, data.WindowY)
	}

	return nil
}
