package plugins

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHotReloadManager_Creation(t *testing.T) {
	api := NewMockAPI()
	manager := NewManager("testdata/plugins", api)

	if manager.hotReloadManager == nil {
		t.Fatal("HotReloadManager should be initialized")
	}
}

func TestHotReloadManager_WatchPlugin(t *testing.T) {
	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, "test.lua")

	// Create test plugin file
	if err := os.WriteFile(pluginPath, []byte("-- test plugin"), 0644); err != nil {
		t.Fatalf("Failed to create test plugin: %v", err)
	}

	api := NewMockAPI()
	manager := NewManager(tmpDir, api)

	if err := manager.hotReloadManager.WatchPlugin("test-plugin", pluginPath); err != nil {
		t.Fatalf("Failed to watch plugin: %v", err)
	}

	// Verify plugin is being watched
	manager.hotReloadManager.mu.RLock()
	_, exists := manager.hotReloadManager.watchedPlugins["test-plugin"]
	manager.hotReloadManager.mu.RUnlock()

	if !exists {
		t.Error("Plugin should be in watched list")
	}
}

func TestHotReloadManager_UnwatchPlugin(t *testing.T) {
	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, "test.lua")

	if err := os.WriteFile(pluginPath, []byte("-- test plugin"), 0644); err != nil {
		t.Fatalf("Failed to create test plugin: %v", err)
	}

	api := NewMockAPI()
	manager := NewManager(tmpDir, api)

	// Watch then unwatch
	if err := manager.hotReloadManager.WatchPlugin("test-plugin", pluginPath); err != nil {
		t.Fatalf("Failed to watch plugin: %v", err)
	}

	if err := manager.hotReloadManager.UnwatchPlugin("test-plugin"); err != nil {
		t.Fatalf("Failed to unwatch plugin: %v", err)
	}

	// Verify plugin is no longer watched
	manager.hotReloadManager.mu.RLock()
	_, exists := manager.hotReloadManager.watchedPlugins["test-plugin"]
	manager.hotReloadManager.mu.RUnlock()

	if exists {
		t.Error("Plugin should not be in watched list")
	}
}

func TestHotReloadManager_StateSnapshot(t *testing.T) {
	api := NewMockAPI()
	manager := NewManager("testdata/plugins", api)

	// Create test plugin
	plugin := &Plugin{
		ID:      "test-plugin",
		Name:    "Test Plugin",
		Version: "1.0.0",
		Enabled: true,
	}

	manager.mu.Lock()
	manager.plugins["test-plugin"] = plugin
	manager.mu.Unlock()

	// Create snapshot
	_, err := manager.hotReloadManager.SnapshotState("test-plugin")
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	// Verify snapshot exists
	manager.hotReloadManager.mu.RLock()
	_, exists := manager.hotReloadManager.stateSnapshots["test-plugin"]
	manager.hotReloadManager.mu.RUnlock()

	if !exists {
		t.Error("Snapshot should exist")
	}
}

func TestHotReloadManager_RollbackPlugin(t *testing.T) {
	api := NewMockAPI()
	manager := NewManager("testdata/plugins", api)

	// Create test plugin
	originalPlugin := &Plugin{
		ID:      "test-plugin",
		Name:    "Test Plugin",
		Version: "1.0.0",
		Enabled: true,
	}

	manager.mu.Lock()
	manager.plugins["test-plugin"] = originalPlugin
	manager.mu.Unlock()

	// Create snapshot
	if _, err := manager.hotReloadManager.SnapshotState("test-plugin"); err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	// Modify plugin
	manager.mu.Lock()
	manager.plugins["test-plugin"].Version = "2.0.0"
	manager.mu.Unlock()

	// Rollback
	if err := manager.hotReloadManager.RollbackPlugin("test-plugin", "1.0.0"); err != nil {
		t.Fatalf("Failed to rollback plugin: %v", err)
	}

	// Verify rollback worked
	plugin, err := manager.GetPlugin("test-plugin")
	if err != nil {
		t.Fatalf("Failed to get plugin: %v", err)
	}

	if plugin.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", plugin.Version)
	}
}

func TestHotReloadManager_ValidateReload(t *testing.T) {
	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, "test.lua")

	// Create valid plugin
	validCode := `
-- metadata
return {
	id = "test-plugin",
	name = "Test Plugin",
	version = "1.0.0"
}
`
	if err := os.WriteFile(pluginPath, []byte(validCode), 0644); err != nil {
		t.Fatalf("Failed to create test plugin: %v", err)
	}

	api := NewMockAPI()
	manager := NewManager(tmpDir, api)

	oldPlugin := &Plugin{ID: "test-plugin", Name: "Test Plugin", Version: "1.0.0"}
	newPlugin := &Plugin{ID: "test-plugin", Name: "Test Plugin", Version: "1.0.0"}
	if err := manager.hotReloadManager.ValidateReload(oldPlugin, newPlugin); err != nil {
		t.Errorf("Valid plugin should pass validation: %v", err)
	}

	// Test invalid plugin
	invalidCode := `
-- invalid lua syntax
return {
	id = "test-plugin"
	name "missing equals"
}
`
	if err := os.WriteFile(pluginPath, []byte(invalidCode), 0644); err != nil {
		t.Fatalf("Failed to update test plugin: %v", err)
	}

	if err := manager.hotReloadManager.ValidateReload(oldPlugin, newPlugin); err == nil {
		t.Error("Invalid plugin should fail validation")
	}
}

func TestHotReloadManager_GetReloadHistory(t *testing.T) {
	api := NewMockAPI()
	manager := NewManager("testdata/plugins", api)

	// Create test plugin
	plugin := &Plugin{
		ID:      "test-plugin",
		Name:    "Test Plugin",
		Version: "1.0.0",
		Enabled: true,
	}

	manager.mu.Lock()
	manager.plugins["test-plugin"] = plugin
	manager.mu.Unlock()

	// Simulate reload
	manager.hotReloadManager.mu.Lock()
	manager.hotReloadManager.reloadHistory = append(manager.hotReloadManager.reloadHistory, &ReloadEvent{
		PluginID:   "test-plugin",
		Timestamp:  time.Now(),
		Status:     ReloadSuccess,
		NewVersion: "1.0.0",
	})
	manager.hotReloadManager.mu.Unlock()

	// Get history
	history, err := manager.GetReloadHistory("test-plugin")
	if err != nil {
		t.Fatalf("Failed to get reload history: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("Expected 1 history record, got %d", len(history))
	}

	if history[0].PluginID != "test-plugin" {
		t.Errorf("Expected plugin ID 'test-plugin', got %s", history[0].PluginID)
	}
}

func TestHotReloadManager_StartStop(t *testing.T) {
	api := NewMockAPI()
	manager := NewManager("testdata/plugins", api)

	// Start hot-reload
	if err := manager.hotReloadManager.Start(); err != nil {
		t.Fatalf("Failed to start hot-reload: %v", err)
	}

	// Verify running
	manager.hotReloadManager.mu.RLock()
	isRunning := manager.hotReloadManager.isRunning
	manager.hotReloadManager.mu.RUnlock()

	if !isRunning {
		t.Error("Hot-reload should be running")
	}

	// Stop hot-reload
	if err := manager.hotReloadManager.Stop(); err != nil {
		t.Fatalf("Failed to stop hot-reload: %v", err)
	}

	// Verify stopped
	manager.hotReloadManager.mu.RLock()
	isRunning = manager.hotReloadManager.isRunning
	manager.hotReloadManager.mu.RUnlock()

	if isRunning {
		t.Error("Hot-reload should be stopped")
	}
}

func TestHotReloadManager_FileChangeDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping file watching test in short mode")
	}

	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, "test.lua")

	// Create initial plugin
	initialCode := `-- version 1`
	if err := os.WriteFile(pluginPath, []byte(initialCode), 0644); err != nil {
		t.Fatalf("Failed to create test plugin: %v", err)
	}

	api := NewMockAPI()
	manager := NewManager(tmpDir, api)

	// Load plugin
	metadata := PluginMetadata{
		ID:      "test-plugin",
		Name:    "Test Plugin",
		Version: "1.0.0",
	}

	plugin := &Plugin{
		ID:       metadata.ID,
		Name:     metadata.Name,
		Version:  metadata.Version,
		path:     pluginPath,
		Enabled:  true,
		metadata: metadata,
	}

	manager.mu.Lock()
	manager.plugins["test-plugin"] = plugin
	manager.mu.Unlock()

	// Start watching
	if err := manager.hotReloadManager.Start(); err != nil {
		t.Fatalf("Failed to start hot-reload: %v", err)
	}

	if err := manager.hotReloadManager.WatchPlugin("test-plugin", pluginPath); err != nil {
		t.Fatalf("Failed to watch plugin: %v", err)
	}

	// Wait for watcher to initialize
	time.Sleep(200 * time.Millisecond)

	// Record initial history count
	history, _ := manager.GetReloadHistory("test-plugin")
	initialHistoryCount := len(history)

	// Modify plugin file
	updatedCode := `-- version 2`
	if err := os.WriteFile(pluginPath, []byte(updatedCode), 0644); err != nil {
		t.Fatalf("Failed to update test plugin: %v", err)
	}

	// Wait for change detection and reload
	time.Sleep(500 * time.Millisecond)

	// Check if reload was triggered
	newHistory, _ := manager.GetReloadHistory("test-plugin")
	newHistoryCount := len(newHistory)

	if newHistoryCount <= initialHistoryCount {
		t.Error("File change should trigger reload")
	}

	// Stop watching
	if err := manager.hotReloadManager.Stop(); err != nil {
		t.Fatalf("Failed to stop hot-reload: %v", err)
	}
}
