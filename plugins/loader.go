package plugins

import (
	"fmt"
	"os"
	"path/filepath"
)

// Loader handles plugin discovery and loading from the filesystem
type Loader struct {
	pluginDir string
}

// NewLoader creates a new plugin loader
func NewLoader(pluginDir string) *Loader {
	return &Loader{
		pluginDir: pluginDir,
	}
}

// DiscoverPlugins discovers plugins in the plugin directory
func (l *Loader) DiscoverPlugins() ([]string, error) {
	// Ensure directory exists
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Read directory
	entries, err := os.ReadDir(l.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory: %w", err)
	}

	var plugins []string

	// Look for Lua plugin files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check for .lua files
		if filepath.Ext(entry.Name()) == ".lua" {
			pluginPath := filepath.Join(l.pluginDir, entry.Name())
			plugins = append(plugins, pluginPath)
		}
	}

	return plugins, nil
}

// ValidatePlugin validates a plugin file
func (l *Loader) ValidatePlugin(pluginPath string) error {
	if pluginPath == "" {
		return fmt.Errorf("plugin path is empty")
	}

	// Check if file exists
	fileInfo, err := os.Stat(pluginPath)
	if err != nil {
		return fmt.Errorf("plugin file not found: %w", err)
	}

	// Check if it's a file
	if fileInfo.IsDir() {
		return fmt.Errorf("plugin path is a directory, not a file")
	}

	// Check file extension
	if filepath.Ext(pluginPath) != ".lua" {
		return fmt.Errorf("plugin must be a .lua file")
	}

	// Check file size (reasonable limit: 10MB)
	if fileInfo.Size() > 10*1024*1024 {
		return fmt.Errorf("plugin file is too large (max 10MB)")
	}

	return nil
}

// LoadPluginCode loads the plugin code from a file
func (l *Loader) LoadPluginCode(pluginPath string) (string, error) {
	if err := l.ValidatePlugin(pluginPath); err != nil {
		return "", err
	}

	// Read file
	code, err := os.ReadFile(pluginPath)
	if err != nil {
		return "", fmt.Errorf("failed to read plugin file: %w", err)
	}

	return string(code), nil
}

// ParsePluginMetadata extracts metadata from plugin code
func (l *Loader) ParsePluginMetadata(code string) (*PluginMetadata, error) {
	if code == "" {
		return nil, fmt.Errorf("plugin code is empty")
	}

	// Placeholder: Real implementation would parse Lua code to extract metadata
	// For now, return a default metadata structure
	metadata := &PluginMetadata{
		ID:      "plugin_default",
		Name:    "Unknown Plugin",
		Version: "1.0.0",
		Author:  "Unknown",
		License: "Unknown",
	}

	// TODO: Parse Lua comment header for metadata:
	// -- @name PluginName
	// -- @version 1.0.0
	// -- @author Author Name
	// -- @description Plugin description
	// -- @license MIT
	// -- @permissions read_save, ui_display
	// -- @hooks on_save, on_load

	return metadata, nil
}

// LoadPlugin loads a plugin from a file
func (l *Loader) LoadPlugin(pluginPath string) (string, *PluginMetadata, error) {
	// Validate plugin file
	if err := l.ValidatePlugin(pluginPath); err != nil {
		return "", nil, err
	}

	// Load plugin code
	code, err := l.LoadPluginCode(pluginPath)
	if err != nil {
		return "", nil, err
	}

	// Parse metadata
	metadata, err := l.ParsePluginMetadata(code)
	if err != nil {
		return "", nil, err
	}

	return code, metadata, nil
}

// LoadAllPlugins loads all plugins from the plugin directory
func (l *Loader) LoadAllPlugins() ([]PluginLoadResult, error) {
	// Discover plugins
	pluginPaths, err := l.DiscoverPlugins()
	if err != nil {
		return nil, err
	}

	var results []PluginLoadResult

	// Load each plugin
	for _, path := range pluginPaths {
		code, metadata, err := l.LoadPlugin(path)
		if err != nil {
			results = append(results, PluginLoadResult{
				Path:  path,
				Error: fmt.Sprintf("failed to load plugin: %v", err),
			})
			continue
		}

		results = append(results, PluginLoadResult{
			Path:     path,
			Code:     code,
			Metadata: metadata,
			Error:    "",
		})
	}

	return results, nil
}

// PluginLoadResult represents the result of loading a plugin
type PluginLoadResult struct {
	Path     string
	Code     string
	Metadata *PluginMetadata
	Error    string
}

// IsError checks if the load result is an error
func (r *PluginLoadResult) IsError() bool {
	return r.Error != ""
}
