package plugins

import (
	"fmt"
	"time"
)

// Plugin represents a loaded plugin
type Plugin struct {
	ID          string
	Name        string
	Version     string
	Author      string
	Description string
	Enabled     bool
	LoadedAt    time.Time
	API         PluginAPI

	// Internal fields
	path     string         // Plugin file path
	metadata PluginMetadata // Full plugin metadata

	// Hooks
	OnLoad     func() error
	OnUnload   func() error
	OnSaveOpen func(savePath string) error
	OnSaveSave func(savePath string) error
	OnCharEdit func(charID int) error
}

// PluginMetadata contains plugin metadata from manifest
type PluginMetadata struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	Author        string    `json:"author"`
	Description   string    `json:"description"`
	Homepage      string    `json:"homepage"`
	License       string    `json:"license"`
	Tags          []string  `json:"tags"`
	Dependencies  []string  `json:"dependencies"`
	MinAppVersion string    `json:"minAppVersion"`
	MaxAppVersion string    `json:"maxAppVersion"`
	Repository    string    `json:"repository"`
	DownloadURL   string    `json:"downloadUrl"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Rating        float64   `json:"rating"`
	Downloads     int       `json:"downloads"`
	Permissions   []string  `json:"permissions"`
	Hooks         []string  `json:"hooks"`
}

// PluginConfig contains plugin configuration
type PluginConfig struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Author      string                 `json:"author"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Settings    map[string]interface{} `json:"settings"`
	Hooks       []string               `json:"hooks"`
	Permissions []string               `json:"permissions"`
}

// ExecutionRecord tracks plugin execution history
type ExecutionRecord struct {
	PluginName string
	PluginVer  string
	StartTime  time.Time
	Duration   time.Duration
	Status     string // "success", "error", "timeout"
	Output     string
	Error      string
}

// Hook types for plugins
type HookType string

const (
	HookLoad     HookType = "load"
	HookUnload   HookType = "unload"
	HookSaveOpen HookType = "save:open"
	HookSaveSave HookType = "save:save"
	HookCharEdit HookType = "character:edit"
	HookUIRender HookType = "ui:render"
	HookMenuAdd  HookType = "menu:add"
	HookSync     HookType = "sync"
)

// CommonPermissions defines standard plugin permissions
var CommonPermissions = struct {
	ReadSave  string
	WriteSave string
	UIDisplay string
	Events    string
}{
	ReadSave:  "read_save",
	WriteSave: "write_save",
	UIDisplay: "ui_display",
	Events:    "events",
}

// CommonHooks defines standard plugin hooks
var CommonHooks = struct {
	OnSave   string
	OnLoad   string
	OnExport string
	OnSync   string
}{
	OnSave:   "on_save",
	OnLoad:   "on_load",
	OnExport: "on_export",
	OnSync:   "on_sync",
}

// NewPlugin creates a new plugin instance
func NewPlugin(metadata *PluginMetadata, api PluginAPI) *Plugin {
	return &Plugin{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Version:     metadata.Version,
		Author:      metadata.Author,
		Description: metadata.Description,
		Enabled:     true,
		LoadedAt:    time.Now(),
		API:         api,
	}
}

// Load loads the plugin
func (p *Plugin) Load() error {
	if p.OnLoad != nil {
		return p.OnLoad()
	}
	return nil
}

// Unload unloads the plugin
func (p *Plugin) Unload() error {
	if p.OnUnload != nil {
		return p.OnUnload()
	}
	return nil
}

// CallHook calls a specific hook if it exists
func (p *Plugin) CallHook(hookType HookType, args ...interface{}) error {
	switch hookType {
	case HookLoad:
		return p.Load()
	case HookUnload:
		return p.Unload()
	case HookSaveOpen:
		if p.OnSaveOpen != nil && len(args) > 0 {
			if savePath, ok := args[0].(string); ok {
				return p.OnSaveOpen(savePath)
			}
		}
	case HookSaveSave:
		if p.OnSaveSave != nil && len(args) > 0 {
			if savePath, ok := args[0].(string); ok {
				return p.OnSaveSave(savePath)
			}
		}
	case HookCharEdit:
		if p.OnCharEdit != nil && len(args) > 0 {
			if charID, ok := args[0].(int); ok {
				return p.OnCharEdit(charID)
			}
		}
	}
	return nil
}

// GetPath returns the plugin file path
func (p *Plugin) GetPath() string {
	return p.path
}

// SetPath sets the plugin file path
func (p *Plugin) SetPath(path string) {
	p.path = path
}

// GetMetadata returns the full plugin metadata
func (p *Plugin) GetMetadata() PluginMetadata {
	return p.metadata
}

// SetMetadata sets the plugin metadata
func (p *Plugin) SetMetadata(metadata PluginMetadata) {
	p.metadata = metadata
}

// Validate validates plugin metadata
func (m *PluginMetadata) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("plugin ID is required")
	}
	if m.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	if m.Author == "" {
		return fmt.Errorf("plugin author is required")
	}
	return nil
}

// ValidatePluginConfig validates a plugin configuration
func ValidatePluginConfig(config PluginConfig) error {
	if config.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if config.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	if config.Author == "" {
		return fmt.Errorf("plugin author is required")
	}
	return nil
}
