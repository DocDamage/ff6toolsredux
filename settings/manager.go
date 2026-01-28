package settings

import (
	"encoding/json"
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
}

// Manager manages application settings
type Manager struct {
	settings *Settings
	filePath string
	mu       sync.RWMutex
}

// New creates a new settings manager with default settings and empty file path
func New() *Manager {
	return &Manager{
		settings: DefaultSettings(),
		filePath: "",
	}
}

// NewManager creates a new settings manager
func NewManager(configPath string) *Manager {
	return &Manager{
		settings: DefaultSettings(),
		filePath: configPath,
	}
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
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if file exists
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		// Use defaults
		return nil
	}

	// Read file
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	// Parse JSON
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
