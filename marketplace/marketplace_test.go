package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// MockRegistry implements test registry
type mockHTTPServer struct {
	server *httptest.Server
	client *Client
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	client := NewClient("https://api.example.com", "test-key")
	if client.baseURL != "https://api.example.com" {
		t.Errorf("Expected baseURL to be https://api.example.com, got %s", client.baseURL)
	}
	if client.apiKey != "test-key" {
		t.Errorf("Expected apiKey to be test-key, got %s", client.apiKey)
	}
}

// TestRegistry creates a test registry
func TestNewRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	if registry.path != tmpDir {
		t.Errorf("Expected path %s, got %s", tmpDir, registry.path)
	}

	plugins := registry.GetInstalledPlugins()
	if len(plugins) != 0 {
		t.Errorf("Expected empty registry, got %d plugins", len(plugins))
	}
}

// TestTrackInstallation tests plugin installation tracking
func TestTrackInstallation(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	pluginID := "test-plugin"
	version := "1.0.0"

	err = registry.TrackInstallation(pluginID, version)
	if err != nil {
		t.Fatalf("Failed to track installation: %v", err)
	}

	record, err := registry.GetPlugin(pluginID)
	if err != nil {
		t.Fatalf("Failed to get plugin: %v", err)
	}

	if record.PluginID != pluginID {
		t.Errorf("Expected plugin ID %s, got %s", pluginID, record.PluginID)
	}
	if record.Version != version {
		t.Errorf("Expected version %s, got %s", version, record.Version)
	}
	if !record.Enabled {
		t.Error("Expected plugin to be enabled")
	}
}

// TestUninstallPlugin tests plugin uninstallation
func TestUninstallPlugin(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	pluginID := "test-plugin"
	_ = registry.TrackInstallation(pluginID, "1.0.0")

	err = registry.UninstallPlugin(pluginID)
	if err != nil {
		t.Fatalf("Failed to uninstall plugin: %v", err)
	}

	_, err = registry.GetPlugin(pluginID)
	if err == nil {
		t.Error("Expected error getting uninstalled plugin")
	}
}

// TestUpdatePlugin tests plugin updates
func TestUpdatePlugin(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	pluginID := "test-plugin"
	_ = registry.TrackInstallation(pluginID, "1.0.0")

	newVersion := "2.0.0"
	err = registry.UpdatePlugin(pluginID, newVersion)
	if err != nil {
		t.Fatalf("Failed to update plugin: %v", err)
	}

	record, _ := registry.GetPlugin(pluginID)
	if record.Version != newVersion {
		t.Errorf("Expected version %s, got %s", newVersion, record.Version)
	}
}

// TestSetPluginEnabled tests enabling/disabling plugins
func TestSetPluginEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	pluginID := "test-plugin"
	_ = registry.TrackInstallation(pluginID, "1.0.0")

	err = registry.SetPluginEnabled(pluginID, false)
	if err != nil {
		t.Fatalf("Failed to disable plugin: %v", err)
	}

	record, _ := registry.GetPlugin(pluginID)
	if record.Enabled {
		t.Error("Expected plugin to be disabled")
	}
}

// TestSetAutoUpdate tests auto-update flag
func TestSetAutoUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	pluginID := "test-plugin"
	_ = registry.TrackInstallation(pluginID, "1.0.0")

	err = registry.SetAutoUpdate(pluginID, false)
	if err != nil {
		t.Fatalf("Failed to set auto-update: %v", err)
	}

	record, _ := registry.GetPlugin(pluginID)
	if record.AutoUpdate {
		t.Error("Expected auto-update to be disabled")
	}
}

// TestRegistryPersistence tests registry save/load
func TestRegistryPersistence(t *testing.T) {
	tmpDir := t.TempDir()

	// Create registry and add plugin
	registry1, _ := NewRegistry(tmpDir)
	pluginID := "test-plugin"
	registry1.TrackInstallation(pluginID, "1.0.0")
	registry1.Close()

	// Load registry from disk
	registry2, _ := NewRegistry(tmpDir)
	defer registry2.Close()

	record, err := registry2.GetPlugin(pluginID)
	if err != nil {
		t.Fatalf("Failed to get plugin from persistent registry: %v", err)
	}

	if record.PluginID != pluginID || record.Version != "1.0.0" {
		t.Error("Plugin data not persisted correctly")
	}
}

// TestGetInstalledPlugins tests listing installed plugins
func TestGetInstalledPlugins(t *testing.T) {
	tmpDir := t.TempDir()
	registry, _ := NewRegistry(tmpDir)
	defer registry.Close()

	plugins := []string{"plugin1", "plugin2", "plugin3"}
	for _, plugin := range plugins {
		registry.TrackInstallation(plugin, "1.0.0")
	}

	installed := registry.GetInstalledPlugins()
	if len(installed) != len(plugins) {
		t.Errorf("Expected %d plugins, got %d", len(plugins), len(installed))
	}
}

// TestSaveLoadRatings tests rating persistence
func TestSaveLoadRatings(t *testing.T) {
	tmpDir := t.TempDir()
	registry, _ := NewRegistry(tmpDir)
	defer registry.Close()

	rating := &PluginRating{
		ID:       "rating1",
		PluginID: "test-plugin",
		UserID:   "user123",
		Rating:   4.5,
		Review:   "Great plugin!",
	}

	err := registry.SaveRatings([]*PluginRating{rating})
	if err != nil {
		t.Fatalf("Failed to save rating: %v", err)
	}

	ratings, err := registry.LoadRatings("test-plugin")
	if err != nil {
		t.Fatalf("Failed to load ratings: %v", err)
	}

	if len(ratings) == 0 {
		t.Error("No ratings found")
	}
	if ratings[0].Review != "Great plugin!" {
		t.Errorf("Expected review 'Great plugin!', got %s", ratings[0].Review)
	}
}

// TestClientListPlugins tests plugin listing
func TestClientListPlugins(t *testing.T) {
	server := startMockServer(t)
	defer server.Close()

	client := NewClient(server.URL, "")

	opts := &ListOptions{Limit: 20}
	plugins, err := client.ListPlugins(context.Background(), opts)
	if err != nil {
		t.Fatalf("Failed to list plugins: %v", err)
	}

	if len(plugins) == 0 {
		t.Error("Expected plugins from server")
	}
}

// TestClientSearchPlugins tests plugin search
func TestClientSearchPlugins(t *testing.T) {
	server := startMockServer(t)
	defer server.Close()

	client := NewClient(server.URL, "")

	plugins, err := client.SearchPlugins(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to search plugins: %v", err)
	}

	if len(plugins) == 0 {
		t.Error("Expected search results")
	}
}

// TestClientCache tests response caching
func TestClientCache(t *testing.T) {
	server := startMockServer(t)
	defer server.Close()

	client := NewClient(server.URL, "")

	// First search
	_, err1 := client.SearchPlugins(context.Background(), "cached")
	// Second search (should use cache)
	_, err2 := client.SearchPlugins(context.Background(), "cached")

	if err1 != nil || err2 != nil {
		t.Fatalf("Search failed: %v, %v", err1, err2)
	}

	// Verify cache is working
	if _, ok := client.cache["search:cached"]; !ok {
		t.Error("Expected search result in cache")
	}
}

// TestClientClearCache tests cache clearing
func TestClientClearCache(t *testing.T) {
	client := NewClient("https://example.com", "")
	client.setCache("test", "value")

	client.ClearCache()

	if len(client.cache) != 0 {
		t.Error("Expected cache to be empty after clear")
	}
}

// TestClientSetCacheTTL tests cache TTL setting
func TestClientSetCacheTTL(t *testing.T) {
	client := NewClient("https://example.com", "")
	ttl := 1 * time.Hour
	client.SetCacheTTL(ttl)

	if client.cacheTTL != ttl {
		t.Errorf("Expected TTL %v, got %v", ttl, client.cacheTTL)
	}
}

// TestClientSubmitRating tests rating submission
func TestClientSubmitRating(t *testing.T) {
	server := startMockServerWithRatings(t)
	defer server.Close()

	tmpDir := t.TempDir()
	client, _ := NewClientWithRegistry(server.URL, "", tmpDir)
	defer client.Close()

	rating := &PluginRating{
		PluginID: "test-plugin",
		Rating:   4.5,
		Review:   "Excellent!",
	}

	err := client.SubmitRating(context.Background(), rating)
	if err != nil {
		t.Fatalf("Failed to submit rating: %v", err)
	}
}

// TestClientGetPluginRatings tests rating retrieval
func TestClientGetPluginRatings(t *testing.T) {
	server := startMockServerWithRatings(t)
	defer server.Close()

	client := NewClient(server.URL, "")

	ratings, err := client.GetPluginRatings(context.Background(), "test-plugin")
	if err != nil {
		t.Fatalf("Failed to get ratings: %v", err)
	}

	if len(ratings) == 0 {
		t.Error("Expected ratings from server")
	}
}

// Helper functions

func startMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/plugins":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"plugins": []RemotePlugin{
					{
						ID:       "plugin1",
						Name:     "Test Plugin 1",
						Author:   "Author1",
						Version:  "1.0.0",
						Category: "Tools",
						Rating:   4.5,
					},
				},
				"total": 1,
			})
		case "/plugins/search":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"plugins": []RemotePlugin{
					{
						ID:       "plugin1",
						Name:     "Test Plugin 1",
						Author:   "Author1",
						Version:  "1.0.0",
						Category: "Tools",
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
}

func startMockServerWithRatings(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/plugins/test-plugin/ratings":
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"ratings": []PluginRating{
						{
							PluginID: "test-plugin",
							Rating:   4.5,
							Review:   "Great!",
						},
					},
				})
			} else if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{"id": "rating1"})
			}
		default:
			http.NotFound(w, r)
		}
	}))
}

// Integration tests

// TestEndToEndInstallation tests complete installation workflow
func TestEndToEndInstallation(t *testing.T) {
	tmpDir := t.TempDir()
	registry, _ := NewRegistry(tmpDir)
	defer registry.Close()

	pluginID := "awesome-plugin"
	version := "1.2.3"

	// Track installation
	err := registry.TrackInstallation(pluginID, version)
	if err != nil {
		t.Fatalf("Installation failed: %v", err)
	}

	// Verify installation
	record, _ := registry.GetPlugin(pluginID)
	if record.PluginID != pluginID || record.Version != version {
		t.Error("Installation verification failed")
	}

	// Update plugin
	err = registry.UpdatePlugin(pluginID, "1.2.4")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	updated, _ := registry.GetPlugin(pluginID)
	if updated.Version != "1.2.4" {
		t.Error("Update verification failed")
	}
}

// TestEndToEndRating tests complete rating workflow
func TestEndToEndRating(t *testing.T) {
	tmpDir := t.TempDir()
	registry, _ := NewRegistry(tmpDir)
	defer registry.Close()

	// Save multiple ratings
	ratings := []*PluginRating{
		{
			PluginID: "plugin1",
			Rating:   5.0,
			Review:   "Perfect!",
		},
		{
			PluginID: "plugin1",
			Rating:   4.0,
			Review:   "Good",
		},
	}

	err := registry.SaveRatings(ratings)
	if err != nil {
		t.Fatalf("Failed to save ratings: %v", err)
	}

	// Load ratings
	loaded, _ := registry.LoadRatings("plugin1")
	if len(loaded) != 2 {
		t.Errorf("Expected 2 ratings, got %d", len(loaded))
	}
}

// BenchmarkRegistryTrackInstallation benchmarks installation tracking
func BenchmarkRegistryTrackInstallation(b *testing.B) {
	tmpDir := os.TempDir()
	registry, _ := NewRegistry(filepath.Join(tmpDir, "bench-registry"))
	defer registry.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pluginID := fmt.Sprintf("plugin-%d", i)
		registry.TrackInstallation(pluginID, "1.0.0")
	}
}

// BenchmarkClientSearch benchmarks search operations
func BenchmarkClientSearch(b *testing.B) {
	server := startMockServer(&testing.T{})
	defer server.Close()

	client := NewClient(server.URL, "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.SearchPlugins(context.Background(), "test")
	}
}
