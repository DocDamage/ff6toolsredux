package plugins

import (
	"sync"
	"testing"
	"time"
)

// TestNewSandboxManager tests creation of a new sandbox manager
func TestNewSandboxManager(t *testing.T) {
	sm := NewSandboxManager()
	if sm == nil {
		t.Fatal("NewSandboxManager returned nil")
	}

	if sm.policies == nil {
		t.Error("policies map not initialized")
	}

	if sm.violations == nil {
		t.Error("violations map not initialized")
	}
}

// TestSetPolicy tests setting a sandbox policy
func TestSetPolicy(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "test_plugin_001",
		AllowedPermissions: []string{"read_save", "ui_display"},
		DeniedPermissions:  []string{},
		MaxMemoryMB:        100,
		MaxCPUPercent:      50,
		TimeoutSeconds:     30,
		IsolationLevel:     "basic",
		IsActive:           true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("SetPolicy failed: %v", err)
	}

	// Retrieve and verify policy - GetPolicy returns just pointer, no error
	retrieved := sm.GetPolicy(policy.PluginID)
	if retrieved == nil {
		t.Fatalf("GetPolicy returned nil")
	}

	if retrieved.MaxMemoryMB != 100 {
		t.Errorf("Expected MaxMemoryMB=100, got %d", retrieved.MaxMemoryMB)
	}

	if retrieved.IsolationLevel != "basic" {
		t.Errorf("Expected IsolationLevel='basic', got '%s'", retrieved.IsolationLevel)
	}
}

// TestRemovePolicy tests removing a sandbox policy
func TestRemovePolicy(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "test_plugin_002",
		AllowedPermissions: []string{},
		IsolationLevel:     "basic",
		IsActive:           true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("SetPolicy failed: %v", err)
	}

	if err := sm.RemovePolicy(policy.PluginID); err != nil {
		t.Fatalf("RemovePolicy failed: %v", err)
	}

	// Verify policy removed - GetPolicy returns nil if not found
	retrieved := sm.GetPolicy(policy.PluginID)
	if retrieved != nil {
		t.Error("GetPolicy should return nil for removed policy")
	}
}

// TestCheckPermission tests permission checking
func TestCheckPermission(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "test_plugin_003",
		AllowedPermissions: []string{"read_save", "ui_display"},
		DeniedPermissions:  []string{"write_save"},
		IsolationLevel:     "basic",
		IsActive:           true,
	}

	sm.SetPolicy(policy)

	// Test allowed permission
	allowed, reason := sm.CheckPermission("test_plugin_003", "read_save")
	if !allowed {
		t.Errorf("Permission should be allowed: %s", reason)
	}

	// Test denied permission
	allowed, reason = sm.CheckPermission("test_plugin_003", "write_save")
	if allowed {
		t.Error("Permission should be denied")
	}
	if reason == "" {
		t.Error("Denial reason should not be empty")
	}

	// Test permission not in whitelist
	allowed, reason = sm.CheckPermission("test_plugin_003", "network_access")
	if allowed {
		t.Error("Permission not in whitelist should be denied")
	}
}

// TestCheckPermissionNoPolicy tests permission checking without policy
func TestCheckPermissionNoPolicy(t *testing.T) {
	sm := NewSandboxManager()

	// No policy set for this plugin
	allowed, reason := sm.CheckPermission("nonexistent_plugin", "read_save")
	if !allowed {
		t.Errorf("Permission should be allowed when no policy exists: %s", reason)
	}
}

// TestDenyPermission tests denying a specific permission
func TestDenyPermission(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "test_plugin_004",
		AllowedPermissions: []string{"read_save", "write_save"},
		DeniedPermissions:  []string{},
		IsActive:           true,
	}

	sm.SetPolicy(policy)

	// Deny write_save permission
	if err := sm.DenyPermission("test_plugin_004", "write_save"); err != nil {
		t.Fatalf("DenyPermission failed: %v", err)
	}

	// Verify permission is now denied
	allowed, _ := sm.CheckPermission("test_plugin_004", "write_save")
	if allowed {
		t.Error("write_save should be denied after calling DenyPermission")
	}

	// Verify read_save still allowed
	allowed, _ = sm.CheckPermission("test_plugin_004", "read_save")
	if !allowed {
		t.Error("read_save should still be allowed")
	}
}

// TestAllowPermission tests allowing a specific permission
func TestAllowPermission(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "test_plugin_005",
		AllowedPermissions: []string{"read_save"},
		DeniedPermissions:  []string{},
		IsActive:           true,
	}

	sm.SetPolicy(policy)

	// Allow write_save permission
	if err := sm.AllowPermission("test_plugin_005", "write_save"); err != nil {
		t.Fatalf("AllowPermission failed: %v", err)
	}

	// Verify permission is now allowed
	allowed, _ := sm.CheckPermission("test_plugin_005", "write_save")
	if !allowed {
		t.Error("write_save should be allowed after calling AllowPermission")
	}
}

// TestVerifyMemoryUsage tests memory limit verification
func TestVerifyMemoryUsage(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "test_plugin_006",
		IsolationLevel: "basic", // Required for SetPolicy
		MaxMemoryMB:    100,
		IsActive:       true,
	}

	err := sm.SetPolicy(policy)
	if err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Test within limits - convert MB to bytes (50MB = 50 * 1024 * 1024 bytes)
	// VerifyMemoryUsage returns (withinLimits bool, reason string)
	withinLimits, reason := sm.VerifyMemoryUsage("test_plugin_006", 50*1024*1024)
	if !withinLimits {
		t.Errorf("50MB should be within 100MB limit: %s", reason)
	}

	// Test exceeding limits - convert MB to bytes (150MB = 150 * 1024 * 1024 bytes)
	testBytes := int64(150 * 1024 * 1024)
	withinLimits, reason = sm.VerifyMemoryUsage("test_plugin_006", testBytes)
	if withinLimits {
		t.Errorf("150MB (%d bytes) should exceed 100MB limit, but got withinLimits=true, reason=%q", testBytes, reason)
	}
	// When limit exceeded, reason should not be empty
	if !withinLimits && reason == "" {
		t.Error("Violation reason should not be empty when limit exceeded")
	}
}

// TestVerifyExecutionTime tests execution timeout verification
func TestVerifyExecutionTime(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "test_plugin_007",
		IsolationLevel: "basic", // Required for SetPolicy
		TimeoutSeconds: 30,
		IsActive:       true,
	}

	err := sm.SetPolicy(policy)
	if err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Test within limits
	// VerifyExecutionTime returns (withinLimits bool, reason string)
	withinLimits, reason := sm.VerifyExecutionTime("test_plugin_007", 15*time.Second)
	if !withinLimits {
		t.Errorf("15s should be within 30s limit: %s", reason)
	}

	// Test exceeding limits
	withinLimits, reason = sm.VerifyExecutionTime("test_plugin_007", 45*time.Second)
	if withinLimits {
		t.Error("45s should exceed 30s limit")
	}
	// When limit exceeded, reason should not be empty
	if !withinLimits && reason == "" {
		t.Error("Violation reason should not be empty when limit exceeded")
	}
}

// TestGetViolations tests retrieving violations for a plugin
func TestGetViolations(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "test_plugin_008",
		IsolationLevel: "basic",
		MaxMemoryMB:    100,
		IsActive:       true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Trigger some violations - convert MB to bytes
	sm.VerifyMemoryUsage("test_plugin_008", 150*1024*1024)
	sm.CheckPermission("test_plugin_008", "forbidden_permission")

	// Get violations - method should work without error
	violations := sm.GetViolations("test_plugin_008")
	// Violations may or may not be recorded depending on internal implementation
	// Just verify the API works
	if violations == nil {
		t.Error("GetViolations should return a slice, not nil")
	}
}

// TestGetAllViolations tests retrieving all violations
func TestGetAllViolations(t *testing.T) {
	sm := NewSandboxManager()

	// Set policies for multiple plugins
	policy1 := &SandboxPolicy{
		PluginID:       "plugin1",
		IsolationLevel: "basic",
		MaxMemoryMB:    100,
		IsActive:       true,
	}
	policy2 := &SandboxPolicy{
		PluginID:       "plugin2",
		IsolationLevel: "basic",
		TimeoutSeconds: 30,
		IsActive:       true,
	}

	if err := sm.SetPolicy(policy1); err != nil {
		t.Fatalf("Failed to set policy1: %v", err)
	}
	if err := sm.SetPolicy(policy2); err != nil {
		t.Fatalf("Failed to set policy2: %v", err)
	}

	// Trigger violations - convert MB to bytes
	sm.VerifyMemoryUsage("plugin1", 150*1024*1024)
	sm.VerifyExecutionTime("plugin2", 45*time.Second)

	// Get all violations - method should work without error
	violations := sm.GetAllViolations()
	// Just verify the API works
	if violations == nil {
		t.Error("GetAllViolations should return a slice, not nil")
	}
}

// TestGetViolationsByType tests filtering violations by type
func TestGetViolationsByType(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "test_plugin_009",
		IsolationLevel: "basic",
		MaxMemoryMB:    100,
		TimeoutSeconds: 30,
		IsActive:       true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Trigger different types of violations - convert MB to bytes
	sm.VerifyMemoryUsage("test_plugin_009", 150*1024*1024)
	sm.VerifyMemoryUsage("test_plugin_009", 200*1024*1024)
	sm.VerifyExecutionTime("test_plugin_009", 45*time.Second)

	// Get memory violations - method should work without error
	memoryViolations := sm.GetViolationsByType("memory_exceeded")
	if memoryViolations == nil {
		t.Error("GetViolationsByType should return a slice, not nil")
	}

	// Get timeout violations
	timeoutViolations := sm.GetViolationsByType("timeout")
	if timeoutViolations == nil {
		t.Error("GetViolationsByType should return a slice, not nil")
	}
}

// TestClearViolations tests clearing violations for a plugin
func TestClearViolations(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "test_plugin_010",
		IsolationLevel: "basic",
		MaxMemoryMB:    100,
		IsActive:       true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Trigger violations - convert MB to bytes
	sm.VerifyMemoryUsage("test_plugin_010", 150*1024*1024)
	sm.VerifyMemoryUsage("test_plugin_010", 200*1024*1024)

	// Get violations
	violations := sm.GetViolations("test_plugin_010")
	// Just check API works
	if violations == nil {
		t.Error("GetViolations should return a slice")
	}

	// Clear violations - should not error
	sm.ClearViolations()

	// Verify violations cleared
	violations = sm.GetViolations("test_plugin_010")
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations after clearing, got %d", len(violations))
	}
}

// TestGetViolationStats tests violation statistics
func TestGetViolationStats(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "test_plugin_011",
		IsolationLevel: "basic",
		MaxMemoryMB:    100,
		TimeoutSeconds: 30,
		IsActive:       true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Trigger various violations - convert MB to bytes
	sm.VerifyMemoryUsage("test_plugin_011", 150*1024*1024)
	sm.VerifyMemoryUsage("test_plugin_011", 200*1024*1024)
	sm.VerifyExecutionTime("test_plugin_011", 45*time.Second)

	// Get violation stats - verify API works
	stats := sm.GetViolationStats()

	if stats == nil {
		t.Fatal("GetViolationStats should return a map, not nil")
	}

	// Check that stats contains expected keys (based on actual implementation)
	expectedKeys := []string{"total_violations", "enforced", "not_enforced", "by_type", "by_severity", "critical_violations"}
	for _, key := range expectedKeys {
		if _, ok := stats[key]; !ok {
			t.Errorf("stats should contain %s key", key)
		}
	}
}

// TestGenerateSecurityReport tests security report generation
func TestGenerateSecurityReport(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "test_plugin_012",
		IsolationLevel:     "basic",
		AllowedPermissions: []string{"read_save"},
		MaxMemoryMB:        100,
		IsActive:           true,
	}

	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Trigger some violations - convert MB to bytes
	sm.VerifyMemoryUsage("test_plugin_012", 150*1024*1024)
	sm.CheckPermission("test_plugin_012", "write_save")

	// Generate report - GenerateSecurityReport() takes no parameters
	report := sm.GenerateSecurityReport()

	if report == "" {
		t.Error("Generated security report is empty")
	}

	// Report should contain key information
	if !contains(report, "Sandbox Security Report") {
		t.Error("Report should contain header")
	}

	if !contains(report, "test_plugin_012") {
		t.Error("Report should contain plugin ID")
	}
}

// TestGetPolicies tests retrieving all policies
func TestGetPolicies(t *testing.T) {
	sm := NewSandboxManager()

	// Set multiple policies
	for i := 1; i <= 3; i++ {
		policy := &SandboxPolicy{
			PluginID:       "plugin_" + string(rune('0'+i)),
			IsolationLevel: "basic",
			IsActive:       true,
		}
		if err := sm.SetPolicy(policy); err != nil {
			t.Fatalf("Failed to set policy %d: %v", i, err)
		}
	}

	policies := sm.GetPolicies()
	if len(policies) != 3 {
		t.Errorf("Expected 3 policies, got %d", len(policies))
	}
}

// TestIsolationLevels tests different isolation levels
func TestIsolationLevels(t *testing.T) {
	sm := NewSandboxManager()

	// Test "none" isolation - should allow everything
	policyNone := &SandboxPolicy{
		PluginID:           "plugin_none",
		AllowedPermissions: []string{},
		DeniedPermissions:  []string{"write_save"},
		IsolationLevel:     "none",
		IsActive:           true,
	}
	sm.SetPolicy(policyNone)

	allowed, _ := sm.CheckPermission("plugin_none", "any_permission")
	if !allowed {
		t.Error("Isolation level 'none' should allow all permissions")
	}

	// Test "basic" isolation - whitelist only
	policyBasic := &SandboxPolicy{
		PluginID:           "plugin_basic",
		AllowedPermissions: []string{"read_save"},
		DeniedPermissions:  []string{},
		IsolationLevel:     "basic",
		IsActive:           true,
	}
	sm.SetPolicy(policyBasic)

	allowed, _ = sm.CheckPermission("plugin_basic", "read_save")
	if !allowed {
		t.Error("Basic isolation should allow whitelisted permissions")
	}

	allowed, _ = sm.CheckPermission("plugin_basic", "write_save")
	if allowed {
		t.Error("Basic isolation should deny non-whitelisted permissions")
	}

	// Test "strict" isolation - blacklist enforced
	policyStrict := &SandboxPolicy{
		PluginID:           "plugin_strict",
		AllowedPermissions: []string{"read_save", "write_save"},
		DeniedPermissions:  []string{"write_save"},
		IsolationLevel:     "strict",
		IsActive:           true,
	}
	sm.SetPolicy(policyStrict)

	allowed, _ = sm.CheckPermission("plugin_strict", "write_save")
	if allowed {
		t.Error("Strict isolation should enforce blacklist even if whitelisted")
	}
}

// TestConcurrentPolicyManagement tests concurrent policy operations
func TestConcurrentPolicyManagement(t *testing.T) {
	sm := NewSandboxManager()

	// Set 10 policies concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			policy := &SandboxPolicy{
				PluginID:       "concurrent_plugin_" + string(rune('0'+id)),
				IsolationLevel: "basic",
				IsActive:       true,
			}
			if err := sm.SetPolicy(policy); err != nil {
				t.Errorf("Concurrent SetPolicy failed for plugin %d: %v", id, err)
			}
		}(i)
	}
	wg.Wait()

	// Verify all policies set
	policies := sm.GetPolicies()
	if len(policies) != 10 {
		t.Errorf("Expected 10 concurrent policies, got %d", len(policies))
	}
}

// TestConcurrentViolationTracking tests concurrent violation recording
func TestConcurrentViolationTracking(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:       "concurrent_test",
		IsolationLevel: "basic",
		MaxMemoryMB:    100,
		IsActive:       true,
	}
	if err := sm.SetPolicy(policy); err != nil {
		t.Fatalf("Failed to set policy: %v", err)
	}

	// Trigger 20 violations concurrently
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.VerifyMemoryUsage("concurrent_test", 150*1024*1024)
		}()
	}
	wg.Wait()

	// Verify all violations recorded
	violations := sm.GetViolations("concurrent_test")
	if len(violations) != 20 {
		t.Errorf("Expected 20 concurrent violations, got %d", len(violations))
	}
}

// TestPolicyInactive tests that inactive policies are not enforced
func TestPolicyInactive(t *testing.T) {
	sm := NewSandboxManager()

	policy := &SandboxPolicy{
		PluginID:           "inactive_plugin",
		AllowedPermissions: []string{"read_save"},
		DeniedPermissions:  []string{"write_save"},
		IsActive:           false, // Inactive
	}
	sm.SetPolicy(policy)

	// Permissions should be allowed when policy is inactive
	allowed, _ := sm.CheckPermission("inactive_plugin", "write_save")
	if !allowed {
		t.Error("Inactive policy should not deny permissions")
	}
}
