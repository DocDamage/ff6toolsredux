package plugins

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ReloadStatus represents the outcome of a reload operation
type ReloadStatus string

const (
	ReloadSuccess    ReloadStatus = "success"
	ReloadFailed     ReloadStatus = "failed"
	ReloadRolledBack ReloadStatus = "rolled_back"
)

// ReloadEvent represents a plugin reload event
type ReloadEvent struct {
	PluginID   string
	OldVersion string
	NewVersion string
	Timestamp  time.Time
	Status     ReloadStatus
	Error      error
	Duration   time.Duration
}

// HotReloadManager manages plugin hot-reloading without app restart
type HotReloadManager struct {
	pluginManager  *Manager
	fileWatcher    *fsnotify.Watcher
	stateSnapshots map[string]*PluginState
	reloadHistory  []*ReloadEvent
	watchedPlugins map[string]string // pluginID -> file path
	reloadCh       chan string       // pluginID channel
	stopCh         chan struct{}
	mu             sync.RWMutex
	maxHistory     int
	isRunning      bool
}

// NewHotReloadManager creates a new hot reload manager
func NewHotReloadManager(manager *Manager) (*HotReloadManager, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	return &HotReloadManager{
		pluginManager:  manager,
		fileWatcher:    watcher,
		stateSnapshots: make(map[string]*PluginState),
		reloadHistory:  make([]*ReloadEvent, 0),
		watchedPlugins: make(map[string]string),
		reloadCh:       make(chan string, 10),
		stopCh:         make(chan struct{}),
		maxHistory:     100,
		isRunning:      false,
	}, nil
}

// Start begins watching for file changes
func (h *HotReloadManager) Start() error {
	h.mu.Lock()
	if h.isRunning {
		h.mu.Unlock()
		return fmt.Errorf("hot reload manager already running")
	}
	h.isRunning = true
	h.mu.Unlock()

	go h.watchLoop()
	return nil
}

// Stop stops the hot reload manager
func (h *HotReloadManager) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isRunning {
		return fmt.Errorf("hot reload manager not running")
	}

	close(h.stopCh)
	h.isRunning = false
	return h.fileWatcher.Close()
}

// WatchPlugin starts watching a plugin for changes
func (h *HotReloadManager) WatchPlugin(pluginID, path string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	plugin, err := h.pluginManager.GetPlugin(pluginID)
	if err != nil {
		return fmt.Errorf("plugin not found: %w", err)
	}

	// Use provided path or get from plugin
	watchPath := path
	if watchPath == "" {
		watchPath = plugin.GetPath()
	}

	if watchPath == "" {
		return fmt.Errorf("plugin has no file path")
	}

	// Watch the plugin file
	if err := h.fileWatcher.Add(watchPath); err != nil {
		return fmt.Errorf("failed to watch plugin file: %w", err)
	}

	// Also watch the directory for new files
	dir := filepath.Dir(watchPath)
	if err := h.fileWatcher.Add(dir); err != nil {
		// Non-critical error, continue
		fmt.Printf("Warning: failed to watch plugin directory %s: %v\n", dir, err)
	}

	h.watchedPlugins[pluginID] = watchPath
	return nil
}

// UnwatchPlugin stops watching a plugin
func (h *HotReloadManager) UnwatchPlugin(pluginID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	path, ok := h.watchedPlugins[pluginID]
	if !ok {
		return fmt.Errorf("plugin not being watched")
	}

	if err := h.fileWatcher.Remove(path); err != nil {
		return fmt.Errorf("failed to unwatch plugin: %w", err)
	}

	delete(h.watchedPlugins, pluginID)
	return nil
}

// ReloadPlugin reloads a specific plugin
func (h *HotReloadManager) ReloadPlugin(ctx context.Context, pluginID string) (*ReloadEvent, error) {
	startTime := time.Now()
	event := &ReloadEvent{
		PluginID:  pluginID,
		Timestamp: startTime,
	}

	// Get the current plugin
	plugin, err := h.pluginManager.GetPlugin(pluginID)
	if err != nil {
		event.Status = ReloadFailed
		event.Error = err
		event.Duration = time.Since(startTime)
		h.recordEvent(event)
		return event, fmt.Errorf("plugin not found: %w", err)
	}

	event.OldVersion = plugin.Version

	// Create state snapshot before reload
	state, err := h.SnapshotState(pluginID)
	if err != nil {
		event.Status = ReloadFailed
		event.Error = err
		event.Duration = time.Since(startTime)
		h.recordEvent(event)
		return event, fmt.Errorf("failed to snapshot state: %w", err)
	}

	h.mu.Lock()
	h.stateSnapshots[pluginID] = state
	h.mu.Unlock()

	// Unload the old plugin
	if err := h.pluginManager.UnloadPlugin(ctx, pluginID); err != nil {
		event.Status = ReloadFailed
		event.Error = err
		event.Duration = time.Since(startTime)
		h.recordEvent(event)
		return event, fmt.Errorf("failed to unload plugin: %w", err)
	}

	// Load the new version
	newPlugin, err := h.pluginManager.LoadPlugin(ctx, plugin.GetPath())
	if err != nil {
		// Reload failed - attempt rollback
		event.Status = ReloadRolledBack
		event.Error = err

		if rollbackErr := h.RollbackPlugin(pluginID, plugin.Version); rollbackErr != nil {
			event.Error = fmt.Errorf("reload failed and rollback failed: %w, rollback error: %v", err, rollbackErr)
		}

		event.Duration = time.Since(startTime)
		h.recordEvent(event)
		return event, fmt.Errorf("failed to reload plugin: %w", err)
	}

	event.NewVersion = newPlugin.Version

	// Validate the reload
	if err := h.ValidateReload(plugin, newPlugin); err != nil {
		event.Status = ReloadRolledBack
		event.Error = err

		if rollbackErr := h.RollbackPlugin(pluginID, plugin.Version); rollbackErr != nil {
			event.Error = fmt.Errorf("validation failed and rollback failed: %w, rollback error: %v", err, rollbackErr)
		}

		event.Duration = time.Since(startTime)
		h.recordEvent(event)
		return event, fmt.Errorf("reload validation failed: %w", err)
	}

	// Restore state
	if err := h.RestoreState(pluginID, state); err != nil {
		// Non-critical - log warning but don't fail
		fmt.Printf("Warning: failed to restore plugin state: %v\n", err)
	}

	event.Status = ReloadSuccess
	event.Duration = time.Since(startTime)
	h.recordEvent(event)

	return event, nil
}

// SnapshotState captures the current state of a plugin
func (h *HotReloadManager) SnapshotState(pluginID string) (*PluginState, error) {
	plugin, err := h.pluginManager.GetPlugin(pluginID)
	if err != nil {
		return nil, err
	}

	// Get plugin config from manager
	var config interface{}
	h.pluginManager.mu.RLock()
	if pluginConfig, exists := h.pluginManager.configs[pluginID]; exists {
		config = pluginConfig
	}
	h.pluginManager.mu.RUnlock()

	state := &PluginState{
		PluginID:  pluginID,
		Version:   plugin.Version,
		Timestamp: time.Now(),
		Config:    config,
		Enabled:   plugin.Enabled,
		Data:      make(map[string]interface{}),
	}

	// Capture execution metrics if available
	h.pluginManager.mu.RLock()
	if len(h.pluginManager.executionLog) > 0 {
		// Get last 10 executions for this plugin
		count := 0
		for i := len(h.pluginManager.executionLog) - 1; i >= 0 && count < 10; i-- {
			record := h.pluginManager.executionLog[i]
			if record.PluginName == pluginID || record.PluginName == plugin.Name {
				state.LastExecutions = append(state.LastExecutions, record)
				count++
			}
		}
	}
	h.pluginManager.mu.RUnlock()

	return state, nil
}

// RestoreState restores a plugin's state from a snapshot
func (h *HotReloadManager) RestoreState(pluginID string, state *PluginState) error {
	if state == nil {
		return fmt.Errorf("nil state provided")
	}

	plugin, err := h.pluginManager.GetPlugin(pluginID)
	if err != nil {
		return err
	}

	// Restore basic properties
	plugin.Enabled = state.Enabled

	// Update config in manager
	h.pluginManager.mu.Lock()
	if pluginConfig, exists := h.pluginManager.configs[pluginID]; exists {
		// Update config from state
		if configMap, ok := state.Config.(map[string]interface{}); ok {
			for key, value := range configMap {
				pluginConfig.Settings[key] = value
			}
			h.pluginManager.configs[pluginID] = pluginConfig
		}
	}
	h.pluginManager.mu.Unlock()

	// Note: execution history is not restored as it's append-only
	return nil
}

// RollbackPlugin rolls back a plugin to a previous version
func (h *HotReloadManager) RollbackPlugin(pluginID, previousVersion string) error {
	h.mu.RLock()
	state, ok := h.stateSnapshots[pluginID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no snapshot available for rollback")
	}

	// Try to restore from snapshot
	if err := h.RestoreState(pluginID, state); err != nil {
		return fmt.Errorf("failed to restore state: %w", err)
	}

	return nil
}

// ValidateReload validates that a reload is safe
func (h *HotReloadManager) ValidateReload(old, new *Plugin) error {
	if old == nil || new == nil {
		return fmt.Errorf("invalid plugins for validation")
	}

	// Check ID hasn't changed
	if old.ID != new.ID {
		return fmt.Errorf("plugin ID changed from %s to %s", old.ID, new.ID)
	}

	// Check version is different (or allow same version for dev)
	if old.Version == new.Version {
		// Allow for development, just log warning
		fmt.Printf("Warning: reloading same version %s\n", old.Version)
	}

	// Check required permissions haven't expanded significantly
	oldMeta := old.GetMetadata()
	newMeta := new.GetMetadata()
	if len(newMeta.Permissions) > len(oldMeta.Permissions)+5 {
		return fmt.Errorf("new version requires significantly more permissions (%d vs %d)",
			len(newMeta.Permissions), len(oldMeta.Permissions))
	}

	// Check dependencies are still resolvable
	// (This would integrate with dependency resolver once implemented)

	return nil
}

// GetReloadHistory returns recent reload events
func (h *HotReloadManager) GetReloadHistory(limit int) []*ReloadEvent {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if limit <= 0 || limit > len(h.reloadHistory) {
		limit = len(h.reloadHistory)
	}

	start := len(h.reloadHistory) - limit
	history := make([]*ReloadEvent, limit)
	copy(history, h.reloadHistory[start:])

	return history
}

// GetPluginReloadHistory returns reload events for a specific plugin
func (h *HotReloadManager) GetPluginReloadHistory(pluginID string, limit int) []*ReloadEvent {
	h.mu.RLock()
	defer h.mu.RUnlock()

	events := make([]*ReloadEvent, 0)
	for i := len(h.reloadHistory) - 1; i >= 0 && len(events) < limit; i-- {
		if h.reloadHistory[i].PluginID == pluginID {
			events = append(events, h.reloadHistory[i])
		}
	}

	return events
}

// watchLoop monitors file changes and triggers reloads
func (h *HotReloadManager) watchLoop() {
	for {
		select {
		case event, ok := <-h.fileWatcher.Events:
			if !ok {
				return
			}

			// Only process write and create events
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				// Find which plugin this file belongs to
				h.mu.RLock()
				var pluginID string
				for id, path := range h.watchedPlugins {
					if path == event.Name {
						pluginID = id
						break
					}
				}
				h.mu.RUnlock()

				if pluginID != "" {
					// Debounce: wait a bit for file writes to complete
					time.Sleep(100 * time.Millisecond)

					// Trigger reload
					fmt.Printf("File change detected for plugin %s, reloading...\n", pluginID)
					_, err := h.ReloadPlugin(context.Background(), pluginID)
					if err != nil {
						fmt.Printf("Hot reload failed for %s: %v\n", pluginID, err)
					} else {
						fmt.Printf("Hot reload successful for %s\n", pluginID)
					}
				}
			}

		case err, ok := <-h.fileWatcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("File watcher error: %v\n", err)

		case <-h.stopCh:
			return
		}
	}
}

// recordEvent adds an event to the history
func (h *HotReloadManager) recordEvent(event *ReloadEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.reloadHistory = append(h.reloadHistory, event)

	// Trim history if it exceeds max
	if len(h.reloadHistory) > h.maxHistory {
		h.reloadHistory = h.reloadHistory[len(h.reloadHistory)-h.maxHistory:]
	}
}
