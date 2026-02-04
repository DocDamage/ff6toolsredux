package settings

import (
	"encoding/json"
	"ffvi_editor/global"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Settings represents application settings
type Settings struct {
	// General
	Theme         string `json:"theme"`         // "dark" or "light"
	Language      string `json:"language"`      // "en", "ja", etc.
	AutoSave      bool   `json:"autoSave"`      // Auto-save on edit
	AutoSaveDelay int    `json:"autoSaveDelay"` // Delay in seconds

	// Backup
	AutoBackup     bool   `json:"autoBackup"`
	BackupsToKeep  int    `json:"backupsToKeep"`
	BackupLocation string `json:"backupLocation"`

	// Validation
	ValidationLevel string `json:"validationLevel"` // "strict", "normal", "permissive"
	AutoFix         bool   `json:"autoFix"`         // Auto-fix issues
	WarnOnSave      bool   `json:"warnOnSave"`      // Warn before saving invalid data

	// Cloud Sync
	CloudEnabled   bool   `json:"cloudEnabled"`
	CloudProvider  string `json:"cloudProvider"` // "gdrive", "dropbox"
	CloudEncrypted bool   `json:"cloudEncrypted"`
	CloudInterval  int    `json:"cloudInterval"` // Sync interval in minutes
	DropboxEnabled bool   `json:"dropboxEnabled,omitempty"`

	// Cloud Provider Credentials
	GoogleDriveClientID     string `json:"googleDriveClientID,omitempty"`
	GoogleDriveClientSecret string `json:"googleDriveClientSecret,omitempty"`
	DropboxAppKey           string `json:"dropboxAppKey,omitempty"`
	DropboxAppSecret        string `json:"dropboxAppSecret,omitempty"`

	// Cloud Sync Settings (for cloud_settings.go compatibility)
	AutoSync            bool   `json:"autoSync,omitempty"`
	SyncIntervalMinutes int    `json:"syncIntervalMinutes,omitempty"`
	EncryptionEnabled   bool   `json:"encryptionEnabled,omitempty"`
	ConflictStrategy    string `json:"conflictStrategy,omitempty"`
	VerifyHashes        bool   `json:"verifyHashes,omitempty"`
	BackupFolderPath    string `json:"backupFolderPath,omitempty"`
	TemplatesFolderPath string `json:"templatesFolderPath,omitempty"`

	// UI
	ShowLineNumbers bool `json:"showLineNumbers"`
	EditorFontSize  int  `json:"editorFontSize"`
	WindowWidth     int  `json:"windowWidth"`
	WindowHeight    int  `json:"windowHeight"`
	WindowMaximized bool `json:"windowMaximized"`

	// Editor Behavior
	ConfirmOnClose  bool `json:"confirmOnClose"`
	ShowUndoHistory bool `json:"showUndoHistory"`
	MaxUndoSteps    int  `json:"maxUndoSteps"`
	EnableDragDrop  bool `json:"enableDragDrop"`
	ShowTooltips    bool `json:"showTooltips"`

	// Shortcuts
	Shortcuts map[string]string `json:"shortcuts"` // action -> key combo

	// Recent Files
	RecentFiles    []string `json:"recentFiles"`
	MaxRecentFiles int      `json:"maxRecentFiles"`

	// ROM Settings
	ROMPath string `json:"romPath"` // Path to FF6 ROM file for sprite loading

	// Advanced
	EnablePlugins      bool   `json:"enablePlugins"`
	EnableScripting    bool   `json:"enableScripting"`
	EnableAchievements bool   `json:"enableAchievements"`
	ShowDebugInfo      bool   `json:"showDebugInfo"`
	LogLevel           string `json:"logLevel"` // "debug", "info", "warn", "error"

	// Legacy config fields
	WindowX             float32                        `json:"windowX,omitempty"`
	WindowY             float32                        `json:"windowY,omitempty"`
	SaveDir             string                         `json:"saveDir,omitempty"`
	AutoEnableCmd       bool                           `json:"autoEnableCmd,omitempty"`
	EnablePlayStation   bool                           `json:"enablePlayStation,omitempty"`
	WorldMapPoints      map[int]map[string]interface{} `json:"worldMapPoints,omitempty"`
	WorldMapLocations   map[int][]string               `json:"worldMapLocations,omitempty"`
	MarketplaceSettings interface{}                    `json:"marketplaceSettings,omitempty"`
}

// Manager manages application settings
type Manager struct {
	settings *Settings
	filePath string
	mu       sync.RWMutex
	logger   *log.Logger
}

// New creates a new settings manager with default settings and empty file path
func New() *Manager {
	return &Manager{
		settings: DefaultSettings(),
		filePath: "",
		logger:   log.New(log.Writer(), "[settings] ", log.LstdFlags),
	}
}

// NewManager creates a new settings manager
func NewManager(configPath string) *Manager {
	return &Manager{
		settings: DefaultSettings(),
		filePath: configPath,
		logger:   log.New(log.Writer(), "[settings] ", log.LstdFlags),
	}
}

// MigrateLegacyConfig migrates ff6editor.config to unified settings if needed
func (m *Manager) MigrateLegacyConfig() error {
	legacyPath := filepath.Join(global.PWD, "ff6editor.config")
	if _, err := os.Stat(legacyPath); err == nil {
		// Read legacy config
		legacyData, err := os.ReadFile(legacyPath)
		if err != nil {
			return err
		}
		// Define legacy struct inline
		type legacyConfig struct {
			WindowX             float32                        `json:"width"`
			WindowY             float32                        `json:"height"`
			SaveDir             string                         `json:"dir"`
			AutoEnableCmd       bool                           `json:"autoEnableCmd"`
			EnablePlayStation   bool                           `json:"ps"`
			WorldMapPoints      map[int]map[string]interface{} `json:"worldMapPoints,omitempty"`
			WorldMapLocations   map[int][]string               `json:"worldMapLocations,omitempty"`
			MarketplaceSettings interface{}                    `json:"marketplaceSettings,omitempty"`
		}
		var legacy legacyConfig
		if err := json.Unmarshal(legacyData, &legacy); err == nil {
			// Map legacy fields to unified settings
			m.settings.WindowX = legacy.WindowX
			m.settings.WindowY = legacy.WindowY
			m.settings.SaveDir = legacy.SaveDir
			m.settings.AutoEnableCmd = legacy.AutoEnableCmd
			m.settings.EnablePlayStation = legacy.EnablePlayStation
			m.settings.WorldMapPoints = legacy.WorldMapPoints
			m.settings.WorldMapLocations = legacy.WorldMapLocations
			m.settings.MarketplaceSettings = legacy.MarketplaceSettings
			// Save unified settings
			if err := m.Save(); err != nil {
				return err
			}
			// Optionally archive or remove legacy config
			_ = os.Rename(legacyPath, legacyPath+".bak")
		}
	}
	return nil
}

// DefaultSettings returns default settings
func DefaultSettings() *Settings {
	return &Settings{
		// General
		Theme:         "dark",
		Language:      "en",
		AutoSave:      false,
		AutoSaveDelay: 30,

		// Backup
		AutoBackup:     true,
		BackupsToKeep:  10,
		BackupLocation: "",

		// Validation
		ValidationLevel: "normal",
		AutoFix:         false,
		WarnOnSave:      true,

		// Cloud Sync
		CloudEnabled:   false,
		CloudProvider:  "gdrive",
		CloudEncrypted: true,
		CloudInterval:  30,

		// UI
		ShowLineNumbers: true,
		EditorFontSize:  12,
		WindowWidth:     1200,
		WindowHeight:    800,
		WindowMaximized: false,

		// Editor Behavior
		ConfirmOnClose:  true,
		ShowUndoHistory: true,
		MaxUndoSteps:    100,
		EnableDragDrop:  true,
		ShowTooltips:    true,

		// Shortcuts
		Shortcuts: map[string]string{
			"save":        "Ctrl+S",
			"undo":        "Ctrl+Z",
			"redo":        "Ctrl+Y",
			"find":        "Ctrl+F",
			"palette":     "Ctrl+Shift+P",
			"close":       "Ctrl+W",
			"new":         "Ctrl+N",
			"open":        "Ctrl+O",
			"backup":      "Ctrl+B",
			"validate":    "Ctrl+Shift+V",
			"preferences": "Ctrl+,",
		},

		// Recent Files
		RecentFiles:    []string{},
		MaxRecentFiles: 10,

		// Advanced
		EnablePlugins:      true,
		EnableScripting:    true,
		EnableAchievements: true,
		ShowDebugInfo:      false,
		LogLevel:           "info",
	}
}

// Load loads settings from file
func (m *Manager) Load() error {
	// Check if file exists (without holding lock to avoid deadlock with MigrateLegacyConfig)
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		// Try to migrate legacy config first
		if err := m.MigrateLegacyConfig(); err != nil {
			return err
		}
		// If still no settings file, use defaults
		if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
			return nil
		}
	}

	// Read file
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	// Parse JSON (with lock)
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := json.Unmarshal(data, m.settings); err != nil {
		return err
	}

	return nil
}

// Save saves settings to file
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(m.settings, "", "  ")
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(m.filePath, data, 0644)
}

// Get returns a copy of the settings
func (m *Manager) Get() *Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return copy
	settingsCopy := *m.settings
	return &settingsCopy
}

// Set updates all settings
func (m *Manager) Set(settings *Settings) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings = settings
}

// GetValue retrieves a specific setting
func (m *Manager) GetValue(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Use reflection or switch on key to return value
	// For simplicity, returning nil for now
	return nil
}

// SetValue sets a specific setting
func (m *Manager) SetValue(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use reflection or switch on key to set value
	// For simplicity, doing nothing for now
}

// AddRecentFile adds a file to recent files
func (m *Manager) AddRecentFile(filePath string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove if already exists
	for i, f := range m.settings.RecentFiles {
		if f == filePath {
			m.settings.RecentFiles = append(
				m.settings.RecentFiles[:i],
				m.settings.RecentFiles[i+1:]...,
			)
			break
		}
	}

	// Add to front
	m.settings.RecentFiles = append([]string{filePath}, m.settings.RecentFiles...)

	// Trim to max
	if len(m.settings.RecentFiles) > m.settings.MaxRecentFiles {
		m.settings.RecentFiles = m.settings.RecentFiles[:m.settings.MaxRecentFiles]
	}
}

// GetRecentFiles returns recent files
func (m *Manager) GetRecentFiles() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string{}, m.settings.RecentFiles...)
}

// ClearRecentFiles clears recent files
func (m *Manager) ClearRecentFiles() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings.RecentFiles = []string{}
}

// Reset resets settings to defaults
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings = DefaultSettings()
}

// --- Map Data Management ---

// GetMapLocations returns user-added map location names for a world
func (m *Manager) GetMapLocations(world int) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.settings.WorldMapLocations == nil {
		return nil
	}
	return append([]string{}, m.settings.WorldMapLocations[world]...)
}

// AddMapLocation adds a location name to a world
func (m *Manager) AddMapLocation(world int, name string) {
	m.mu.Lock()
	if m.settings.WorldMapLocations == nil {
		m.settings.WorldMapLocations = make(map[int][]string)
	}
	for _, n := range m.settings.WorldMapLocations[world] {
		if n == name {
			m.mu.Unlock()
			return
		}
	}
	m.settings.WorldMapLocations[world] = append(m.settings.WorldMapLocations[world], name)
	m.mu.Unlock()
	if m.filePath != "" {
		if err := m.Save(); err != nil && m.logger != nil {
			m.logger.Printf("failed to save settings after adding map location: %v", err)
		}
	}
}

// GetAllMapPoints returns all map points for a world
func (m *Manager) GetAllMapPoints(world int) map[string]map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]map[string]float64)
	if m.settings.WorldMapPoints == nil {
		return result
	}
	if points, ok := m.settings.WorldMapPoints[world]; ok {
		for name, v := range points {
			mp, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			x, xok := mp["X"].(float64)
			y, yok := mp["Y"].(float64)
			if !xok || !yok {
				continue
			}
			result[name] = map[string]float64{"X": x, "Y": y}
		}
	}
	return result
}

// GetMapPoint returns a map point for a world and name
func (m *Manager) GetMapPoint(world int, name string) (x, y float64, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.settings.WorldMapPoints == nil {
		return 0, 0, false
	}
	if points, ok := m.settings.WorldMapPoints[world]; ok {
		if v, ok := points[name]; ok {
			mp, ok := v.(map[string]interface{})
			if !ok {
				return 0, 0, false
			}
			x, xok := mp["X"].(float64)
			y, yok := mp["Y"].(float64)
			if !xok || !yok {
				return 0, 0, false
			}
			return x, y, true
		}
	}
	return 0, 0, false
}

// SetMapPoint sets a map point for a world and name
func (m *Manager) SetMapPoint(world int, name string, x, y float64) {
	m.mu.Lock()
	if m.settings.WorldMapPoints == nil {
		m.settings.WorldMapPoints = make(map[int]map[string]interface{})
	}
	if m.settings.WorldMapPoints[world] == nil {
		m.settings.WorldMapPoints[world] = make(map[string]interface{})
	}
	m.settings.WorldMapPoints[world][name] = map[string]interface{}{"X": x, "Y": y}
	m.mu.Unlock()
	if m.filePath != "" {
		if err := m.Save(); err != nil && m.logger != nil {
			m.logger.Printf("failed to save settings after setting map point: %v", err)
		}
	}
}

// ClearMapPoint removes a map point for a world and name
func (m *Manager) ClearMapPoint(world int, name string) {
	m.mu.Lock()
	if m.settings.WorldMapPoints == nil {
		m.mu.Unlock()
		return
	}
	if m.settings.WorldMapPoints[world] == nil {
		m.mu.Unlock()
		return
	}
	delete(m.settings.WorldMapPoints[world], name)
	m.mu.Unlock()
	if m.filePath != "" {
		if err := m.Save(); err != nil && m.logger != nil {
			m.logger.Printf("failed to save settings after clearing map point: %v", err)
		}
	}
}

// ClearAllMapPoints removes all map points for a world
func (m *Manager) ClearAllMapPoints(world int) {
	m.mu.Lock()
	if m.settings.WorldMapPoints == nil {
		m.mu.Unlock()
		return
	}
	delete(m.settings.WorldMapPoints, world)
	m.mu.Unlock()
	if m.filePath != "" {
		if err := m.Save(); err != nil && m.logger != nil {
			m.logger.Printf("failed to save settings after clearing all map points: %v", err)
		}
	}
}

// RemoveMapLocation removes a location from a world
func (m *Manager) RemoveMapLocation(world int, name string) {
	m.mu.Lock()
	if m.settings.WorldMapLocations == nil {
		m.mu.Unlock()
		return
	}
	locations := m.settings.WorldMapLocations[world]
	for i, loc := range locations {
		if loc == name {
			m.settings.WorldMapLocations[world] = append(locations[:i], locations[i+1:]...)
			break
		}
	}
	m.mu.Unlock()
	if m.filePath != "" {
		if err := m.Save(); err != nil && m.logger != nil {
			m.logger.Printf("failed to save settings after removing map location: %v", err)
		}
	}
}

// ClearMapLocations removes all locations for a world
func (m *Manager) ClearMapLocations(world int) {
	m.mu.Lock()
	if m.settings.WorldMapLocations == nil {
		m.mu.Unlock()
		return
	}
	delete(m.settings.WorldMapLocations, world)
	m.mu.Unlock()
	if m.filePath != "" {
		if err := m.Save(); err != nil && m.logger != nil {
			m.logger.Printf("failed to save settings after clearing map locations: %v", err)
		}
	}
}
