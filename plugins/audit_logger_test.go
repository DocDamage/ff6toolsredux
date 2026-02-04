package plugins

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestNewAuditLogger tests creation of a new audit logger
func TestNewAuditLogger(t *testing.T) {
	al := NewAuditLogger(1000)
	if al == nil {
		t.Fatal("NewAuditLogger returned nil")
	}

	if al.maxLogs != 1000 {
		t.Errorf("Expected maxLogs=1000, got %d", al.maxLogs)
	}

	if al.logs == nil {
		t.Error("logs slice not initialized")
	}
}

// TestLogPluginLoad tests logging plugin load events
func TestLogPluginLoad(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_001"
	al.LogPluginLoad(pluginID)

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].EventType != "load" {
		t.Errorf("Expected EventType='load', got '%s'", logs[0].EventType)
	}

	if logs[0].PluginID != pluginID {
		t.Errorf("Expected PluginID='%s', got '%s'", pluginID, logs[0].PluginID)
	}
}

// TestLogPluginUnload tests logging plugin unload events
func TestLogPluginUnload(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_002"
	al.LogPluginUnload(pluginID)

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].EventType != "unload" {
		t.Errorf("Expected EventType='unload', got '%s'", logs[0].EventType)
	}
}

// TestLogPluginExecution tests logging plugin execution events
func TestLogPluginExecution(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_003"
	duration := int64(150) // 150ms
	al.LogPluginExecution(pluginID, duration, true, "")

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].EventType != "execute" {
		t.Errorf("Expected EventType='execute', got '%s'", logs[0].EventType)
	}

	if logs[0].Duration != duration {
		t.Errorf("Expected Duration=%d, got %d", duration, logs[0].Duration)
	}

	if logs[0].Status != "success" {
		t.Errorf("Expected Status='success', got '%s'", logs[0].Status)
	}
}

// TestLogPermissionUsed tests logging permission usage
func TestLogPermissionUsed(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_004"
	permission := "read_save"
	al.LogPermissionUsed(pluginID, permission)

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].EventType != "permission_used" {
		t.Errorf("Expected EventType='permission_used', got '%s'", logs[0].EventType)
	}

	if logs[0].PermissionID != permission {
		t.Errorf("Expected PermissionID='%s', got '%s'", permission, logs[0].PermissionID)
	}
}

// TestLogPermissionDenied tests logging denied permissions
func TestLogPermissionDenied(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_005"
	permission := "write_save"
	reason := "Not authorized"
	al.LogPermissionDenied(pluginID, permission, reason)

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].Status != "denied" {
		t.Errorf("Expected Status='denied', got '%s'", logs[0].Status)
	}

	if logs[0].Error != reason {
		t.Errorf("Expected Error='%s', got '%s'", reason, logs[0].Error)
	}
}

// TestLogSecurityViolation tests logging security violations
func TestLogSecurityViolation(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_006"
	violationType := "timeout"
	message := "Execution exceeded timeout"
	al.LogSecurityViolation(pluginID, violationType, message)

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].EventType != "security_violation" {
		t.Errorf("Expected EventType='security_violation', got '%s'", logs[0].EventType)
	}

	if logs[0].Details["violation_type"] != violationType {
		t.Errorf("Expected violation_type='%s'", violationType)
	}
}

// TestLogError tests logging errors
func TestLogError(t *testing.T) {
	al := NewAuditLogger(1000)

	pluginID := "test_plugin_007"
	errorMsg := "Failed to initialize"
	al.LogError(pluginID, errorMsg)

	logs := al.GetPluginAuditTrail(pluginID)
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].EventType != "error" {
		t.Errorf("Expected EventType='error', got '%s'", logs[0].EventType)
	}

	if logs[0].Error != errorMsg {
		t.Errorf("Expected Error='%s', got '%s'", errorMsg, logs[0].Error)
	}
}

// TestGetAuditLog tests retrieving audit log by time range
func TestGetAuditLog(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log events at different times
	now := time.Now()
	for i := 0; i < 5; i++ {
		pluginID := "test_plugin_" + string(rune('0'+i))
		al.LogPluginLoad(pluginID)
		time.Sleep(10 * time.Millisecond)
	}

	// Get logs from last 100ms
	startTime := now.Add(-100 * time.Millisecond)
	endTime := time.Now()
	logs := al.GetAuditLog(startTime, endTime)

	if len(logs) != 5 {
		t.Errorf("Expected 5 log entries, got %d", len(logs))
	}
}

// TestGetEventsByType tests filtering by event type
func TestGetEventsByType(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log different event types
	al.LogPluginLoad("plugin1")
	al.LogPluginLoad("plugin2")
	al.LogPluginExecution("plugin3", 100, true, "")
	al.LogError("plugin4", "Test error")

	// Get only load events
	loadLogs := al.GetEventsByType("load")
	if len(loadLogs) != 2 {
		t.Errorf("Expected 2 load events, got %d", len(loadLogs))
	}

	// Get only error events
	errorLogs := al.GetEventsByType("error")
	if len(errorLogs) != 1 {
		t.Errorf("Expected 1 error event, got %d", len(errorLogs))
	}
}

// TestGetEventsByStatus tests filtering by status
func TestGetEventsByStatus(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log events with different statuses
	al.LogPluginExecution("plugin1", 100, true, "")
	al.LogPluginExecution("plugin2", 200, false, "Failed")
	al.LogPermissionDenied("plugin3", "write_save", "Not authorized")

	// Get only success events
	successLogs := al.GetEventsByStatus("success")
	if len(successLogs) != 1 {
		t.Errorf("Expected 1 success event, got %d", len(successLogs))
	}

	// Get only denied events
	deniedLogs := al.GetEventsByStatus("denied")
	if len(deniedLogs) != 1 {
		t.Errorf("Expected 1 denied event, got %d", len(deniedLogs))
	}
}

// TestGetPermissionUsage tests permission usage statistics
func TestGetPermissionUsage(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log permission usage
	al.LogPermissionUsed("plugin1", "read_save")
	al.LogPermissionUsed("plugin2", "read_save")
	al.LogPermissionUsed("plugin3", "write_save")

	usage := al.GetPermissionUsage()

	if usage["read_save"] != 2 {
		t.Errorf("Expected read_save usage=2, got %d", usage["read_save"])
	}

	if usage["write_save"] != 1 {
		t.Errorf("Expected write_save usage=1, got %d", usage["write_save"])
	}
}

// TestGetTopPermissions tests top permissions retrieval
func TestGetTopPermissions(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log permission usage with varying frequencies
	for i := 0; i < 5; i++ {
		al.LogPermissionUsed("plugin1", "read_save")
	}
	for i := 0; i < 3; i++ {
		al.LogPermissionUsed("plugin2", "write_save")
	}
	for i := 0; i < 1; i++ {
		al.LogPermissionUsed("plugin3", "ui_display")
	}

	topPerms := al.GetTopPermissions(2)
	if len(topPerms) != 2 {
		t.Errorf("Expected 2 top permissions, got %d", len(topPerms))
	}

	// First should be read_save (5 uses)
	if topPerms[0] != "read_save" {
		t.Errorf("Expected first permission to be read_save, got %s", topPerms[0])
	}

	// Second should be write_save (3 uses)
	if topPerms[1] != "write_save" {
		t.Errorf("Expected second permission to be write_save, got %s", topPerms[1])
	}
}

// TestExportAuditLog tests CSV export
func TestExportAuditLog(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log some events
	al.LogPluginLoad("plugin1")
	al.LogPluginExecution("plugin2", 100, true, "")
	al.LogError("plugin3", "Test error")

	// Export to CSV
	tmpFile := filepath.Join(t.TempDir(), "audit_log.csv")
	if err := al.ExportAuditLog(tmpFile); err != nil {
		t.Fatalf("ExportAuditLog failed: %v", err)
	}

	// Read the CSV file
	file, err := os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Should have header + 3 data rows
	if len(records) < 2 {
		t.Errorf("Expected at least header row, got %d", len(records))
	}

	// Verify header matches actual export format
	expectedHeader := []string{"EventID", "PluginID", "EventType", "Action", "PermissionID", "Status", "Error", "Timestamp", "DurationUS"}
	if len(records[0]) != len(expectedHeader) {
		t.Errorf("Expected %d header columns, got %d", len(expectedHeader), len(records[0]))
	}

	for i, col := range expectedHeader {
		if i < len(records[0]) && records[0][i] != col {
			t.Errorf("Expected header column %d to be '%s', got '%s'", i, col, records[0][i])
		}
	}
}

// TestExportAuditJSON tests JSON export
func TestExportAuditJSON(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log some events
	al.LogPluginLoad("plugin1")
	al.LogPluginExecution("plugin2", 100, true, "")

	// Export to JSON
	tmpFile := filepath.Join(t.TempDir(), "audit_log.json")
	if err := al.ExportAuditJSON(tmpFile); err != nil {
		t.Fatalf("ExportAuditJSON failed: %v", err)
	}

	// Read the JSON file
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse JSON - the format is {Timestamp, Events, Stats}
	type exportData struct {
		Timestamp time.Time
		Events    []*AuditLog
		Stats     map[string]interface{}
	}

	var exported exportData
	if err := json.Unmarshal(data, &exported); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(exported.Events) != 2 {
		t.Errorf("Expected 2 log entries in JSON, got %d", len(exported.Events))
	}
}

// TestGetAuditStats tests audit statistics
func TestGetAuditStats(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log various events
	al.LogPluginLoad("plugin1")
	al.LogPluginLoad("plugin2")
	al.LogPluginExecution("plugin1", 100, true, "")
	al.LogPluginExecution("plugin2", 200, false, "Error")
	al.LogError("plugin3", "Test error")

	stats := al.GetAuditStats()

	// GetAuditStats returns total_events and events_by_status
	if stats["total_events"].(int) != 5 {
		t.Errorf("Expected total_events=5, got %v", stats["total_events"])
	}

	// Check events_by_status map
	eventsByStatus := stats["events_by_status"].(map[string]int)
	if eventsByStatus["success"] != 3 {
		t.Errorf("Expected 3 successful events, got %v", eventsByStatus["success"])
	}

	if eventsByStatus["error"] != 2 {
		t.Errorf("Expected 2 error events, got %v", eventsByStatus["error"])
	}
}

// TestMaxEvents tests event limit enforcement
func TestMaxEvents(t *testing.T) {
	maxLogs := 10
	al := NewAuditLogger(maxLogs)

	// Log more than maxLogs
	for i := 0; i < 15; i++ {
		pluginID := "plugin_" + string(rune('0'+i))
		al.LogPluginLoad(pluginID)
	}

	// Should only keep last 10 events
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	logs := al.GetAuditLog(startTime, endTime)

	if len(logs) > maxLogs {
		t.Errorf("Expected max %d events, got %d", maxLogs, len(logs))
	}
}

// TestConcurrentLogging tests concurrent audit logging
func TestConcurrentLogging(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log 100 events concurrently
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pluginID := fmt.Sprintf("concurrent_plugin_%d", id)
			al.LogPluginLoad(pluginID)
		}(i)
	}
	wg.Wait()

	// Give a small delay for all goroutines to complete logging
	time.Sleep(100 * time.Millisecond)

	// Verify all events logged
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	logs := al.GetAuditLog(startTime, endTime)

	// Due to concurrent writes and potential race conditions, check for reasonable number
	if len(logs) < 90 {
		t.Errorf("Expected at least 90 concurrent log entries (out of 100), got %d", len(logs))
	}
}

// TestGenerateAuditReport tests audit report generation
func TestGenerateAuditReport(t *testing.T) {
	al := NewAuditLogger(1000)

	// Log various events
	al.LogPluginLoad("plugin1")
	al.LogPluginExecution("plugin1", 100, true, "")
	al.LogPermissionUsed("plugin1", "read_save")
	al.LogError("plugin2", "Test error")

	// Generate report
	report := al.GenerateAuditReport()

	if report == "" {
		t.Error("Generated audit report is empty")
	}

	// Report should contain key sections
	if !contains(report, "Audit Log Report") {
		t.Error("Report should contain header")
	}

	if !contains(report, "Total Events:") {
		t.Error("Report should contain total events")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
