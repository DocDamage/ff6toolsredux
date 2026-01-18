package marketplace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Registry manages local plugin installation records and metadata
type Registry struct {
	path             string
	mu               sync.RWMutex
	installedPlugins map[string]*InstallRecord
	ratings          map[string][]*PluginRating
	updateChecks     map[string]time.Time
}

// InstallRecord tracks a plugin installation
type InstallRecord struct {
	PluginID    string    `json:"pluginId"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Enabled     bool      `json:"enabled"`
	AutoUpdate  bool      `json:"autoUpdate"`
	Location    string    `json:"location"`
}

// RegistryData represents the persisted registry structure
type RegistryData struct {
	Plugins map[string]*InstallRecord `json:"plugins"`
	Ratings map[string][]*PluginRating `json:"ratings"`
}

// NewRegistry creates a new plugin registry
func NewRegistry(path string) (*Registry, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create registry directory: %w", err)
	}

	registry := &Registry{
		path:             path,
		installedPlugins: make(map[string]*InstallRecord),
		ratings:          make(map[string][]*PluginRating),
		updateChecks:     make(map[string]time.Time),
	}

	// Load existing data
	if err := registry.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	return registry, nil
}

// TrackInstallation records a plugin installation
func (r *Registry) TrackInstallation(pluginID, version string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.installedPlugins[pluginID] = &InstallRecord{
		PluginID:    pluginID,
		Version:     version,
		InstalledAt: now,
		UpdatedAt:   now,
		Enabled:     true,
		AutoUpdate:  true,
		Location:    "",
	}

	return r.save()
}

// UninstallPlugin removes a plugin from the registry
func (r *Registry) UninstallPlugin(pluginID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.installedPlugins, pluginID)
	delete(r.ratings, pluginID)
	delete(r.updateChecks, pluginID)

	return r.save()
}

// GetInstalledPlugins returns a list of all installed plugins
func (r *Registry) GetInstalledPlugins() []InstallRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]InstallRecord, 0, len(r.installedPlugins))
	for _, record := range r.installedPlugins {
		plugins = append(plugins, *record)
	}
	return plugins
}

// GetPlugin returns a specific installed plugin
func (r *Registry) GetPlugin(pluginID string) (*InstallRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, ok := r.installedPlugins[pluginID]
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	return record, nil
}

// UpdatePlugin updates a plugin to a new version
func (r *Registry) UpdatePlugin(pluginID, newVersion string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.installedPlugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	record.Version = newVersion
	record.UpdatedAt = time.Now()

	return r.save()
}

// SetPluginEnabled enables or disables a plugin
func (r *Registry) SetPluginEnabled(pluginID string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.installedPlugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	record.Enabled = enabled
	return r.save()
}

// SetAutoUpdate sets the auto-update flag for a plugin
func (r *Registry) SetAutoUpdate(pluginID string, autoUpdate bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.installedPlugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	record.AutoUpdate = autoUpdate
	return r.save()
}

// SaveRatings saves plugin ratings/reviews
func (r *Registry) SaveRatings(ratings []*PluginRating) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rating := range ratings {
		if r.ratings[rating.PluginID] == nil {
			r.ratings[rating.PluginID] = make([]*PluginRating, 0)
		}
		r.ratings[rating.PluginID] = append(r.ratings[rating.PluginID], rating)
	}

	return r.save()
}

// LoadRatings retrieves stored ratings for a plugin
func (r *Registry) LoadRatings(pluginID string) ([]*PluginRating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ratings, ok := r.ratings[pluginID]
	if !ok {
		return []*PluginRating{}, nil
	}
	return ratings, nil
}

// GetCachedPlugins returns cached plugin list
func (r *Registry) GetCachedPlugins() ([]RemotePlugin, error) {
	// TODO: Implement cached plugin list retrieval
	return []RemotePlugin{}, nil
}

// CachePluginList caches a list of remote plugins
func (r *Registry) CachePluginList(plugins []RemotePlugin) error {
	// TODO: Implement plugin list caching
	return nil
}

// CheckUpdates checks for available updates
func (r *Registry) CheckUpdates(plugins []string) (map[string]*UpdateInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	updates := make(map[string]*UpdateInfo)
	// TODO: Implement update checking
	return updates, nil
}

// MarkPluginUpdated records when a plugin was updated
func (r *Registry) MarkPluginUpdated(pluginID, newVersion string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.installedPlugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	record.Version = newVersion
	record.UpdatedAt = time.Now()

	return r.save()
}

// RecordUpdateCheck records when an update check was performed
func (r *Registry) RecordUpdateCheck(pluginID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.updateChecks[pluginID] = time.Now()
	return nil
}

// GetLastUpdateCheck gets the last update check time for a plugin
func (r *Registry) GetLastUpdateCheck(pluginID string) time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if t, ok := r.updateChecks[pluginID]; ok {
		return t
	}
	return time.Time{}
}

// Private methods

// load loads the registry from disk
func (r *Registry) load() error {
	registryPath := filepath.Join(r.path, "registry.json")

	data, err := os.ReadFile(registryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Registry doesn't exist yet, return nil
		}
		return err
	}

	var regData RegistryData
	if err := json.Unmarshal(data, &regData); err != nil {
		return fmt.Errorf("failed to unmarshal registry: %w", err)
	}

	r.installedPlugins = regData.Plugins
	r.ratings = regData.Ratings

	return nil
}

// save persists the registry to disk
func (r *Registry) save() error {
	regData := RegistryData{
		Plugins: r.installedPlugins,
		Ratings: r.ratings,
	}

	data, err := json.MarshalIndent(regData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	registryPath := filepath.Join(r.path, "registry.json")
	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// Close closes the registry and performs cleanup
func (r *Registry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Persist any pending changes
	return r.save()
}
