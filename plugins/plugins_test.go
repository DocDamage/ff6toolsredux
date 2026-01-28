package plugins

import (
	"context"
	"testing"
	"time"
)

// TestPluginCreation tests plugin creation and metadata
func TestPluginCreation(t *testing.T) {
	metadata := &PluginMetadata{
		ID:      "test_plugin",
		Name:    "Test Plugin",
		Version: "1.0.0",
		Author:  "Test Author",
	}

	if err := metadata.Validate(); err != nil {
		t.Fatalf("Valid metadata failed validation: %v", err)
	}
}

// TestPluginConfigValidation tests plugin configuration validation
func TestPluginConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PluginConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: PluginConfig{
				Name:    "Test",
				Version: "1.0.0",
				Author:  "Author",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: PluginConfig{
				Version: "1.0.0",
				Author:  "Author",
			},
			wantErr: true,
		},
		{
			name: "missing version",
			config: PluginConfig{
				Name:   "Test",
				Author: "Author",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePluginConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePluginConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestManagerCreation tests manager creation
func TestManagerCreation(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)

	if manager == nil {
		t.Fatal("Manager creation failed")
	}

	if manager.maxPlugins != 50 {
		t.Errorf("Default max plugins = %d, want 50", manager.maxPlugins)
	}
}

// TestManagerStats tests manager statistics
func TestManagerStats(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)

	stats := manager.GetStats()

	if stats["total_plugins"] != 0 {
		t.Errorf("Initial plugins = %v, want 0", stats["total_plugins"])
	}

	if stats["max_plugins"] != 50 {
		t.Errorf("Max plugins = %v, want 50", stats["max_plugins"])
	}
}

// TestExecutionLog tests execution logging
func TestExecutionLog(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)

	log := manager.GetExecutionLog()
	if len(log) != 0 {
		t.Errorf("Initial log length = %d, want 0", len(log))
	}

	manager.ClearExecutionLog()
	log = manager.GetExecutionLog()
	if len(log) != 0 {
		t.Errorf("After clear log length = %d, want 0", len(log))
	}
}

// TestManagerMaxPlugins tests manager plugin limit
func TestManagerMaxPlugins(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)

	manager.SetMaxPlugins(1)

	if manager.maxPlugins != 1 {
		t.Errorf("Max plugins = %d, want 1", manager.maxPlugins)
	}
}

// TestSandboxMode tests sandbox mode toggling
func TestSandboxMode(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)

	if !manager.sandbox {
		t.Error("Sandbox mode should be enabled by default")
	}

	manager.SetSandboxMode(false)
	if manager.sandbox {
		t.Error("Sandbox mode should be disabled")
	}

	manager.SetSandboxMode(true)
	if !manager.sandbox {
		t.Error("Sandbox mode should be enabled")
	}
}

// TestAPIPermissions tests API permission checking
func TestAPIPermissions(t *testing.T) {
	perms := []string{CommonPermissions.ReadSave, CommonPermissions.UIDisplay}
	api := NewAPIImpl(nil, perms)

	if !api.HasPermission(CommonPermissions.ReadSave) {
		t.Error("Should have read_save permission")
	}

	if !api.HasPermission(CommonPermissions.UIDisplay) {
		t.Error("Should have ui_display permission")
	}

	if api.HasPermission(CommonPermissions.WriteSave) {
		t.Error("Should not have write_save permission")
	}
}

// TestAPILogging tests API logging
func TestAPILogging(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	ctx := context.Background()

	// Should not error even with no logger set
	err := api.Log(ctx, "info", "test message")
	if err != nil {
		t.Errorf("Log() error = %v, want nil", err)
	}

	// Custom logger
	var loggedLevel, loggedMessage string
	api.SetLogger(func(level, msg string) {
		loggedLevel = level
		loggedMessage = msg
	})

	err = api.Log(ctx, "debug", "debug message")
	if err != nil {
		t.Errorf("Log() error = %v, want nil", err)
	}

	if loggedLevel != "debug" || loggedMessage != "debug message" {
		t.Errorf("Logger called with level=%s, msg=%s", loggedLevel, loggedMessage)
	}
}

// TestCommonPermissions tests common permissions constants
func TestCommonPermissions(t *testing.T) {
	if CommonPermissions.ReadSave != "read_save" {
		t.Errorf("ReadSave = %s, want read_save", CommonPermissions.ReadSave)
	}

	if CommonPermissions.WriteSave != "write_save" {
		t.Errorf("WriteSave = %s, want write_save", CommonPermissions.WriteSave)
	}

	if CommonPermissions.UIDisplay != "ui_display" {
		t.Errorf("UIDisplay = %s, want ui_display", CommonPermissions.UIDisplay)
	}

	if CommonPermissions.Events != "events" {
		t.Errorf("Events = %s, want events", CommonPermissions.Events)
	}
}

// TestCommonHooks tests common hooks constants
func TestCommonHooks(t *testing.T) {
	if CommonHooks.OnSave != "on_save" {
		t.Errorf("OnSave = %s, want on_save", CommonHooks.OnSave)
	}

	if CommonHooks.OnLoad != "on_load" {
		t.Errorf("OnLoad = %s, want on_load", CommonHooks.OnLoad)
	}

	if CommonHooks.OnExport != "on_export" {
		t.Errorf("OnExport = %s, want on_export", CommonHooks.OnExport)
	}

	if CommonHooks.OnSync != "on_sync" {
		t.Errorf("OnSync = %s, want on_sync", CommonHooks.OnSync)
	}
}

// BenchmarkPluginMetadataValidation benchmarks metadata validation
func BenchmarkPluginMetadataValidation(b *testing.B) {
	metadata := &PluginMetadata{
		ID:      "bench_plugin",
		Name:    "Benchmark Plugin",
		Version: "1.0.0",
		Author:  "Benchmark",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metadata.Validate()
	}
}

// BenchmarkPluginConfigValidation benchmarks config validation
func BenchmarkPluginConfigValidation(b *testing.B) {
	config := PluginConfig{
		Name:    "Bench",
		Version: "1.0.0",
		Author:  "Bench",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidatePluginConfig(config)
	}
}

// BenchmarkAPIPermissionCheck benchmarks permission checking
func BenchmarkAPIPermissionCheck(b *testing.B) {
	api := NewAPIImpl(nil, []string{CommonPermissions.ReadSave, CommonPermissions.UIDisplay})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = api.HasPermission(CommonPermissions.ReadSave)
	}
}

// TestManagerStartStop tests manager start/stop lifecycle
func TestManagerStartStop(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)
	ctx := context.Background()

	if manager.isRunning {
		t.Error("Manager should not be running initially")
	}

	// Start manager
	err := manager.Start(ctx)
	if err != nil {
		t.Errorf("Start() error = %v, want nil", err)
	}

	if !manager.isRunning {
		t.Error("Manager should be running after Start()")
	}

	// Stop manager
	err = manager.Stop(ctx)
	if err != nil {
		t.Errorf("Stop() error = %v, want nil", err)
	}

	if manager.isRunning {
		t.Error("Manager should not be running after Stop()")
	}
}

// TestExecutionRecord tests execution record creation
func TestExecutionRecord(t *testing.T) {
	now := time.Now()
	record := ExecutionRecord{
		PluginName: "TestPlugin",
		PluginVer:  "1.0.0",
		StartTime:  now,
		Duration:   100 * time.Millisecond,
		Status:     "success",
		Output:     "Test output",
		Error:      "",
	}

	if record.PluginName != "TestPlugin" {
		t.Errorf("PluginName = %s, want TestPlugin", record.PluginName)
	}

	if record.Status != "success" {
		t.Errorf("Status = %s, want success", record.Status)
	}

	if record.Duration != 100*time.Millisecond {
		t.Errorf("Duration = %v, want 100ms", record.Duration)
	}
}

// TestPluginEnableDisable tests enabling and disabling plugins
func TestPluginEnableDisable(t *testing.T) {
	api := NewAPIImpl(nil, []string{})
	manager := NewManager("./plugins", api)

	// Create a plugin metadata
	metadata := &PluginMetadata{
		ID:      "test_enable_disable",
		Name:    "Test Plugin",
		Version: "1.0.0",
		Author:  "Test",
	}

	plugin := NewPlugin(metadata, api)
	manager.plugins["test_enable_disable"] = plugin
	manager.configs["test_enable_disable"] = PluginConfig{
		Name:    "Test Plugin",
		Version: "1.0.0",
		Author:  "Test",
		Enabled: true,
	}

	// Test initial state
	if !plugin.Enabled {
		t.Error("Plugin should be enabled initially")
	}

	// Disable plugin
	err := manager.DisablePlugin("test_enable_disable")
	if err != nil {
		t.Errorf("DisablePlugin() error = %v, want nil", err)
	}

	if plugin.Enabled {
		t.Error("Plugin should be disabled after DisablePlugin()")
	}

	// Enable plugin
	err = manager.EnablePlugin("test_enable_disable")
	if err != nil {
		t.Errorf("EnablePlugin() error = %v, want nil", err)
	}

	if !plugin.Enabled {
		t.Error("Plugin should be enabled after EnablePlugin()")
	}
}
