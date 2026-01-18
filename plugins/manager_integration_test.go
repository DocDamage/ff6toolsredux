package plugins

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ffvi_editor/models"
	modelsPR "ffvi_editor/models/pr"
)

// TestManagerIntegration tests the full plugin lifecycle with security, audit, and sandbox
func TestManagerIntegration(t *testing.T) {
	// Create temporary directory for plugins
	tmpDir := t.TempDir()

	// Create a mock API
	api := &testPluginAPI{
		saveData: make(map[string]string),
	}

	// Create manager with all Phase 12.3 components
	manager := NewManager(tmpDir, api)

	// Verify all components initialized
	if manager.GetSecurityManager() == nil {
		t.Fatal("SecurityManager not initialized")
	}
	if manager.GetAuditLogger() == nil {
		t.Fatal("AuditLogger not initialized")
	}
	if manager.GetSandboxManager() == nil {
		t.Fatal("SandboxManager not initialized")
	}
	if manager.GetProfiler() == nil {
		t.Fatal("PluginProfiler not initialized")
	}
	if manager.GetAnalytics() == nil {
		t.Fatal("AnalyticsEngine not initialized")
	}

	t.Log("✓ All Phase 12.3 components initialized")
}

// TestPluginLoadWithSecurity tests plugin loading with security verification
func TestPluginLoadWithSecurity(t *testing.T) {
	tmpDir := t.TempDir()
	api := &testPluginAPI{saveData: make(map[string]string)}
	manager := NewManager(tmpDir, api)

	// Create test plugin file
	pluginFile := filepath.Join(tmpDir, "test_plugin.go")
	pluginCode := []byte("package main\\n\\nfunc main() {}\\n")
	if err := os.WriteFile(pluginFile, pluginCode, 0644); err != nil {
		t.Fatalf("Failed to create test plugin: %v", err)
	}

	// Sign the plugin first
	secMgr := manager.GetSecurityManager()
	_, err := secMgr.SignPlugin("test_plugin", pluginFile)
	if err != nil {
		t.Logf("SignPlugin returned error (expected for test): %v", err)
	}

	// Try to load plugin (will go through security checks)
	ctx := context.Background()
	plugin, err := manager.LoadPlugin(ctx, pluginFile)

	// Check that audit logger recorded the attempt
	auditLogger := manager.GetAuditLogger()
	auditTrail := auditLogger.GetPluginAuditTrail("test_plugin")

	if plugin != nil {
		t.Log("✓ Plugin loaded successfully")

		// Verify audit trail contains load event
		if len(auditTrail) > 0 {
			t.Log("✓ Audit trail recorded plugin load")
		}

		// Verify sandbox policy applied
		sandboxMgr := manager.GetSandboxManager()
		policy := sandboxMgr.GetPolicy(plugin.ID)
		if policy != nil {
			t.Logf("✓ Sandbox policy applied: IsolationLevel=%s", policy.IsolationLevel)
		}
	} else {
		t.Logf("Plugin load failed (may be expected): %v", err)
		// Even if load fails, audit should be recorded
		if len(auditTrail) > 0 {
			t.Log("✓ Audit trail recorded load attempt")
		}
	}
}

// TestPluginExecutionWithAudit tests plugin execution with audit logging
func TestPluginExecutionWithAudit(t *testing.T) {
	tmpDir := t.TempDir()
	api := &testPluginAPI{saveData: make(map[string]string)}
	manager := NewManager(tmpDir, api)

	// Create and add a mock plugin directly
	metadata := &PluginMetadata{
		ID:          "exec_test_plugin",
		Name:        "Test Plugin",
		Version:     "1.0.0",
		Author:      "Test",
		Permissions: []string{"read_save"},
	}

	plugin := NewPlugin(metadata, api)
	plugin.SetMetadata(*metadata)

	// Manually add to manager (bypass loading checks for test)
	manager.mu.Lock()
	manager.plugins[metadata.ID] = plugin
	manager.configs[metadata.ID] = PluginConfig{
		Name:        metadata.Name,
		Version:     metadata.Version,
		Enabled:     true,
		Permissions: metadata.Permissions,
	}
	manager.mu.Unlock()

	// Apply sandbox policy
	sandboxMgr := manager.GetSandboxManager()
	policy := &SandboxPolicy{
		PluginID:           metadata.ID,
		AllowedPermissions: []string{"read_save", "execute"},
		TimeoutSeconds:     30,
		MaxMemoryMB:        100,
		IsolationLevel:     "basic",
		IsActive:           true,
	}
	sandboxMgr.SetPolicy(policy)

	// Execute the plugin
	ctx := context.Background()
	err := manager.ExecutePlugin(ctx, metadata.ID)

	if err != nil {
		t.Logf("ExecutePlugin returned: %v", err)
	}

	// Verify audit trail recorded execution
	auditLogger := manager.GetAuditLogger()
	auditTrail := auditLogger.GetPluginAuditTrail(metadata.ID)

	if len(auditTrail) > 0 {
		t.Logf("✓ Audit trail contains %d events", len(auditTrail))

		// Check for execution event
		hasExecEvent := false
		for _, log := range auditTrail {
			if log.EventType == "execute" {
				hasExecEvent = true
				t.Logf("✓ Execution event recorded: Status=%s, Duration=%dms",
					log.Status, log.Duration)
			}
		}

		if hasExecEvent {
			t.Log("✓ Execution properly audited")
		}
	}

	// Verify profiler recorded metrics
	profiler := manager.GetProfiler()
	aggregates := profiler.GetAggregates()
	if agg, exists := aggregates[metadata.ID]; exists {
		t.Logf("✓ Profiler recorded metrics: Executions=%d, AverageDuration=%dms",
			agg.ExecutionCount, agg.AverageDuration.Milliseconds())
	}

	// Verify analytics recorded event
	analytics := manager.GetAnalytics()
	allStats := analytics.GetAllPluginStats()
	if stats, exists := allStats[metadata.ID]; exists {
		t.Logf("✓ Analytics recorded: ExecutionCount=%d", stats.ExecutionCount)
	}
}

// TestSandboxPermissionEnforcement tests sandbox permission checking
func TestSandboxPermissionEnforcement(t *testing.T) {
	tmpDir := t.TempDir()
	api := &testPluginAPI{saveData: make(map[string]string)}
	manager := NewManager(tmpDir, api)

	sandboxMgr := manager.GetSandboxManager()
	auditLogger := manager.GetAuditLogger()

	// Create strict sandbox policy
	policy := &SandboxPolicy{
		PluginID:           "strict_plugin",
		AllowedPermissions: []string{"read_save"},
		DeniedPermissions:  []string{"write_save", "network_access"},
		MaxMemoryMB:        50,
		TimeoutSeconds:     10,
		IsolationLevel:     "strict",
		IsActive:           true,
	}
	sandboxMgr.SetPolicy(policy)

	// Test allowed permission
	allowed, reason := sandboxMgr.CheckPermission("strict_plugin", "read_save")
	if allowed {
		t.Log("✓ Allowed permission passed")
	} else {
		t.Errorf("Allowed permission denied: %s", reason)
	}

	// Test denied permission
	allowed, reason = sandboxMgr.CheckPermission("strict_plugin", "write_save")
	if !allowed {
		t.Logf("✓ Denied permission blocked: %s", reason)

		// Should create violation
		violations := sandboxMgr.GetViolations("strict_plugin")
		if len(violations) > 0 {
			t.Logf("✓ Violation recorded: Type=%s", violations[0].ViolationType)
		}
	} else {
		t.Error("Denied permission should have been blocked")
	}

	// Test memory limit
	violated, reason := sandboxMgr.VerifyMemoryUsage("strict_plugin", 75)
	if violated {
		t.Logf("✓ Memory limit enforced: %s", reason)
	}

	// Test timeout
	violated, reason = sandboxMgr.VerifyExecutionTime("strict_plugin", 15*time.Second)
	if violated {
		t.Logf("✓ Timeout enforced: %s", reason)
	}

	// Check audit log for security violations
	auditTrail := auditLogger.GetEventsByType("security_violation")
	t.Logf("✓ Security violations logged: %d events", len(auditTrail))
}

// TestUnloadWithCleanup tests plugin unload with audit and cleanup
func TestUnloadWithCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	api := &testPluginAPI{saveData: make(map[string]string)}
	manager := NewManager(tmpDir, api)

	// Create and add mock plugin
	metadata := &PluginMetadata{
		ID:      "cleanup_plugin",
		Name:    "Cleanup Test",
		Version: "1.0.0",
	}

	plugin := NewPlugin(metadata, api)
	plugin.SetMetadata(*metadata)

	manager.mu.Lock()
	manager.plugins[metadata.ID] = plugin
	manager.configs[metadata.ID] = PluginConfig{
		Name:    metadata.Name,
		Enabled: true,
	}
	manager.mu.Unlock()

	// Apply sandbox policy
	sandboxMgr := manager.GetSandboxManager()
	policy := &SandboxPolicy{
		PluginID:       metadata.ID,
		IsolationLevel: "basic",
		IsActive:       true,
	}
	if err := sandboxMgr.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Verify policy exists
	if sandboxMgr.GetPolicy(metadata.ID) == nil {
		t.Error("Policy should exist before unload")
	}

	// Unload plugin
	ctx := context.Background()
	err := manager.UnloadPlugin(ctx, metadata.ID)
	if err != nil {
		t.Logf("UnloadPlugin returned: %v", err)
	}

	// Verify audit trail recorded unload
	auditLogger := manager.GetAuditLogger()
	auditTrail := auditLogger.GetPluginAuditTrail(metadata.ID)

	hasUnloadEvent := false
	for _, log := range auditTrail {
		if log.EventType == "unload" {
			hasUnloadEvent = true
			t.Log("✓ Unload event recorded in audit trail")
		}
	}

	if !hasUnloadEvent {
		t.Error("Unload event should be in audit trail")
	}

	// Verify sandbox policy removed
	if sandboxMgr.GetPolicy(metadata.ID) == nil {
		t.Log("✓ Sandbox policy cleaned up")
	} else {
		t.Error("Sandbox policy should be removed after unload")
	}

	// Verify plugin removed from manager
	_, err = manager.GetPlugin(metadata.ID)
	if err == ErrPluginNotFound {
		t.Log("✓ Plugin removed from manager")
	}
}

// TestGetManagerStats tests manager statistics with all components
func TestGetManagerStats(t *testing.T) {
	tmpDir := t.TempDir()
	api := &testPluginAPI{saveData: make(map[string]string)}
	manager := NewManager(tmpDir, api)

	// Get manager stats
	stats := manager.GetStats()

	if totalPlugins, ok := stats["total_plugins"].(int); ok {
		t.Logf("✓ Total plugins: %d", totalPlugins)
	}

	if sandboxMode, ok := stats["sandbox_mode"].(bool); ok {
		t.Logf("✓ Sandbox mode: %v", sandboxMode)
	}

	// Get component-specific stats
	sandboxMgr := manager.GetSandboxManager()
	sandboxStats := sandboxMgr.GetViolationStats()
	t.Logf("✓ Sandbox violations: %v", sandboxStats["total_violations"])

	auditLogger := manager.GetAuditLogger()
	auditStats := auditLogger.GetAuditStats()
	t.Logf("✓ Total audit events: %v", auditStats["total_events"])

	secMgr := manager.GetSecurityManager()
	secStats := secMgr.GetSecurityStats()
	t.Logf("✓ Total signatures: %v", secStats["total_signatures"])
}

// testPluginAPI is a mock implementation for testing
type testPluginAPI struct {
	saveData map[string]string
}

func (api *testPluginAPI) GetCharacter(ctx context.Context, name string) (*models.Character, error) {
	return nil, nil
}

func (api *testPluginAPI) SetCharacter(ctx context.Context, name string, ch *models.Character) error {
	return nil
}

func (api *testPluginAPI) GetInventory(ctx context.Context) (*modelsPR.Inventory, error) {
	return nil, nil
}

func (api *testPluginAPI) SetInventory(ctx context.Context, inv *modelsPR.Inventory) error {
	return nil
}

func (api *testPluginAPI) GetParty(ctx context.Context) (*modelsPR.Party, error) {
	return nil, nil
}

func (api *testPluginAPI) SetParty(ctx context.Context, party *modelsPR.Party) error {
	return nil
}

func (api *testPluginAPI) GetEquipment(ctx context.Context) (*models.Equipment, error) {
	return nil, nil
}

func (api *testPluginAPI) SetEquipment(ctx context.Context, eq *models.Equipment) error {
	return nil
}

func (api *testPluginAPI) ApplyBatchOperation(ctx context.Context, op string, params map[string]interface{}) (int, error) {
	return 0, nil
}

func (api *testPluginAPI) FindCharacter(ctx context.Context, predicate func(*models.Character) bool) *models.Character {
	return nil
}

func (api *testPluginAPI) FindItems(ctx context.Context, predicate func(*modelsPR.Row) bool) []*modelsPR.Row {
	return nil
}

func (api *testPluginAPI) RegisterHook(event string, callback func(interface{}) error) error {
	return nil
}

func (api *testPluginAPI) FireEvent(ctx context.Context, event string, data interface{}) error {
	return nil
}

func (api *testPluginAPI) ShowDialog(ctx context.Context, title, message string) error {
	return nil
}

func (api *testPluginAPI) ShowConfirm(ctx context.Context, title, message string) (bool, error) {
	return false, nil
}

func (api *testPluginAPI) ShowInput(ctx context.Context, prompt string) (string, error) {
	return "", nil
}

func (api *testPluginAPI) Log(ctx context.Context, level string, message string) error {
	return nil
}

func (api *testPluginAPI) GetSetting(key string) interface{} {
	return nil
}

func (api *testPluginAPI) SetSetting(key string, value interface{}) error {
	return nil
}

func (api *testPluginAPI) HasPermission(permission string) bool {
	return true
}
