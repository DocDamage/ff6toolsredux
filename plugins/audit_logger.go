package plugins

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// AuditLog represents a single audit event
type AuditLog struct {
	EventID      string                 // Unique event identifier
	PluginID     string                 // ID of plugin involved
	EventType    string                 // "load", "execute", "error", "permission_used", "security_violation"
	Action       string                 // "READ_SAVE", "WRITE_SAVE", "EXECUTE_COMMAND", etc.
	PermissionID string                 // Permission being used/denied
	Status       string                 // "success", "denied", "error"
	Error        string                 // Error message if any
	Timestamp    time.Time              // When event occurred
	Duration     int64                  // Duration in microseconds (if applicable)
	Details      map[string]interface{} // Additional context
}

// AuditLogger tracks all plugin operations and permissions
type AuditLogger struct {
	logs           []*AuditLog
	permissionUse  map[string]int         // permission -> use count
	errorLog       map[string]int         // error message -> count
	pluginEventLog map[string][]*AuditLog // pluginID -> events
	mu             sync.RWMutex
	maxLogs        int
	autoCleanup    bool
	retentionDays  int
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(maxLogs int) *AuditLogger {
	return &AuditLogger{
		logs:           make([]*AuditLog, 0, maxLogs),
		permissionUse:  make(map[string]int),
		errorLog:       make(map[string]int),
		pluginEventLog: make(map[string][]*AuditLog),
		maxLogs:        maxLogs,
		autoCleanup:    true,
		retentionDays:  90, // Default 90-day retention
	}
}

// LogEvent records an audit event
func (al *AuditLogger) LogEvent(pluginID, eventType, action, permissionID, status, errMsg string, duration int64, details map[string]interface{}) {
	al.mu.Lock()
	defer al.mu.Unlock()

	event := &AuditLog{
		EventID:      fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		PluginID:     pluginID,
		EventType:    eventType,
		Action:       action,
		PermissionID: permissionID,
		Status:       status,
		Error:        errMsg,
		Timestamp:    time.Now(),
		Duration:     duration,
		Details:      details,
	}

	al.logs = append(al.logs, event)

	// Track permission usage
	if permissionID != "" {
		al.permissionUse[permissionID]++
	}

	// Track errors
	if errMsg != "" {
		al.errorLog[errMsg]++
	}

	// Track per-plugin events
	al.pluginEventLog[pluginID] = append(al.pluginEventLog[pluginID], event)

	// Enforce max size
	if len(al.logs) > al.maxLogs {
		al.logs = al.logs[len(al.logs)-al.maxLogs:]
	}
}

// LogPermissionUsed logs permission usage
func (al *AuditLogger) LogPermissionUsed(pluginID, permission string) {
	al.LogEvent(pluginID, "permission_used", "READ", permission, "success", "", 0, map[string]interface{}{
		"permission_type": "access",
	})
}

// LogPermissionDenied logs denied permission attempt
func (al *AuditLogger) LogPermissionDenied(pluginID, permission, reason string) {
	al.LogEvent(pluginID, "permission_denied", "DENIED", permission, "denied", reason, 0, map[string]interface{}{
		"permission_type": "access",
		"reason":          reason,
	})
}

// LogPluginLoad logs plugin loading
func (al *AuditLogger) LogPluginLoad(pluginID string) {
	al.LogEvent(pluginID, "load", "LOAD", "", "success", "", 0, nil)
}

// LogPluginUnload logs plugin unloading
func (al *AuditLogger) LogPluginUnload(pluginID string) {
	al.LogEvent(pluginID, "unload", "UNLOAD", "", "success", "", 0, nil)
}

// LogPluginExecution logs plugin execution with metrics
func (al *AuditLogger) LogPluginExecution(pluginID string, duration int64, success bool, err string) {
	status := "success"
	if !success {
		status = "error"
	}

	al.LogEvent(pluginID, "execute", "EXECUTE", "", status, err, duration, map[string]interface{}{
		"duration_microseconds": duration,
	})
}

// LogSecurityViolation logs security violations
func (al *AuditLogger) LogSecurityViolation(pluginID, violationType, description string) {
	al.LogEvent(pluginID, "security_violation", "SECURITY", "", "denied", description, 0, map[string]interface{}{
		"violation_type": violationType,
	})
}

// LogError logs plugin errors
func (al *AuditLogger) LogError(pluginID, errMsg string) {
	al.LogEvent(pluginID, "error", "ERROR", "", "error", errMsg, 0, map[string]interface{}{
		"error_category": "plugin_error",
	})
}

// GetPluginAuditTrail returns all events for a plugin
func (al *AuditLogger) GetPluginAuditTrail(pluginID string) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	events := al.pluginEventLog[pluginID]
	eventsCopy := make([]*AuditLog, len(events))
	copy(eventsCopy, events)
	return eventsCopy
}

// GetAuditLog returns events within a time range
func (al *AuditLogger) GetAuditLog(startTime, endTime time.Time) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var filtered []*AuditLog
	for _, log := range al.logs {
		if log.Timestamp.After(startTime) && log.Timestamp.Before(endTime) {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// GetEventsByType returns events of a specific type
func (al *AuditLogger) GetEventsByType(eventType string) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var filtered []*AuditLog
	for _, log := range al.logs {
		if log.EventType == eventType {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// GetEventsByStatus returns events with a specific status
func (al *AuditLogger) GetEventsByStatus(status string) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var filtered []*AuditLog
	for _, log := range al.logs {
		if log.Status == status {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// GetPermissionUsage returns usage count for each permission
func (al *AuditLogger) GetPermissionUsage() map[string]int {
	al.mu.RLock()
	defer al.mu.RUnlock()

	usage := make(map[string]int)
	for k, v := range al.permissionUse {
		usage[k] = v
	}
	return usage
}

// GetTopPermissions returns N most-used permissions
func (al *AuditLogger) GetTopPermissions(count int) []string {
	al.mu.RLock()
	defer al.mu.RUnlock()

	type permStat struct {
		name  string
		count int
	}

	var stats []permStat
	for perm, cnt := range al.permissionUse {
		stats = append(stats, permStat{perm, cnt})
	}

	// Sort by count descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})

	// Extract top N names
	if len(stats) > count {
		stats = stats[:count]
	}

	var results []string
	for _, s := range stats {
		results = append(results, s.name)
	}
	return results
}

// GetTopErrors returns N most common errors
func (al *AuditLogger) GetTopErrors(count int) []string {
	al.mu.RLock()
	defer al.mu.RUnlock()

	type errStat struct {
		msg   string
		count int
	}

	var stats []errStat
	for msg, cnt := range al.errorLog {
		stats = append(stats, errStat{msg, cnt})
	}

	// Sort by count descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})

	// Extract top N
	if len(stats) > count {
		stats = stats[:count]
	}

	var results []string
	for _, s := range stats {
		results = append(results, s.msg)
	}
	return results
}

// GetAuditStats returns comprehensive audit statistics
func (al *AuditLogger) GetAuditStats() map[string]interface{} {
	al.mu.RLock()
	defer al.mu.RUnlock()

	stats := map[string]interface{}{
		"total_events":          len(al.logs),
		"total_plugins":         len(al.pluginEventLog),
		"total_permissions":     len(al.permissionUse),
		"total_errors":          len(al.errorLog),
		"events_by_type":        make(map[string]int),
		"events_by_status":      make(map[string]int),
		"permission_violations": 0,
		"security_violations":   0,
	}

	eventsByType := make(map[string]int)
	eventsByStatus := make(map[string]int)
	permViolations := 0
	secViolations := 0

	for _, log := range al.logs {
		eventsByType[log.EventType]++
		eventsByStatus[log.Status]++

		if log.Status == "denied" {
			permViolations++
		}
		if log.EventType == "security_violation" {
			secViolations++
		}
	}

	stats["events_by_type"] = eventsByType
	stats["events_by_status"] = eventsByStatus
	stats["permission_violations"] = permViolations
	stats["security_violations"] = secViolations

	return stats
}

// GenerateAuditReport generates a text report
func (al *AuditLogger) GenerateAuditReport() string {
	al.mu.RLock()
	defer al.mu.RUnlock()

	report := fmt.Sprintf("=== Audit Log Report ===\n")
	report += fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339))

	stats := al.GetAuditStats()
	report += fmt.Sprintf("Total Events: %d\n", stats["total_events"])
	report += fmt.Sprintf("Total Plugins: %d\n", stats["total_plugins"])
	report += fmt.Sprintf("Permission Violations: %d\n", stats["permission_violations"])
	report += fmt.Sprintf("Security Violations: %d\n", stats["security_violations"])
	report += fmt.Sprintf("\n--- Event Summary ---\n")

	eventsByType := stats["events_by_type"].(map[string]int)
	for eventType, count := range eventsByType {
		report += fmt.Sprintf("%s: %d\n", eventType, count)
	}

	return report
}

// ExportAuditLog exports logs to CSV file
func (al *AuditLogger) ExportAuditLog(filePath string) error {
	al.mu.RLock()
	defer al.mu.RUnlock()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"EventID", "PluginID", "EventType", "Action", "PermissionID",
		"Status", "Error", "Timestamp", "DurationUS",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Write events
	for _, log := range al.logs {
		record := []string{
			log.EventID,
			log.PluginID,
			log.EventType,
			log.Action,
			log.PermissionID,
			log.Status,
			log.Error,
			log.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%d", log.Duration),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// ExportAuditJSON exports logs to JSON file
func (al *AuditLogger) ExportAuditJSON(filePath string) error {
	al.mu.RLock()
	defer al.mu.RUnlock()

	type exportData struct {
		Timestamp time.Time
		Events    []*AuditLog
		Stats     map[string]interface{}
	}

	stats := al.GetAuditStats()
	data := exportData{
		Timestamp: time.Now(),
		Events:    al.logs,
		Stats:     stats,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ClearOldLogs removes logs older than specified duration
func (al *AuditLogger) ClearOldLogs(olderThan time.Duration) {
	al.mu.Lock()
	defer al.mu.Unlock()

	cutoffTime := time.Now().Add(-olderThan)

	var retained []*AuditLog
	for _, log := range al.logs {
		if log.Timestamp.After(cutoffTime) {
			retained = append(retained, log)
		}
	}

	al.logs = retained

	// Rebuild plugin event log
	al.pluginEventLog = make(map[string][]*AuditLog)
	for _, log := range al.logs {
		al.pluginEventLog[log.PluginID] = append(al.pluginEventLog[log.PluginID], log)
	}
}

// ClearAllLogs clears all audit logs
func (al *AuditLogger) ClearAllLogs() {
	al.mu.Lock()
	defer al.mu.Unlock()

	al.logs = make([]*AuditLog, 0, al.maxLogs)
	al.permissionUse = make(map[string]int)
	al.errorLog = make(map[string]int)
	al.pluginEventLog = make(map[string][]*AuditLog)
}

// GetAllLogs returns all audit logs
func (al *AuditLogger) GetAllLogs() []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	logsCopy := make([]*AuditLog, len(al.logs))
	copy(logsCopy, al.logs)
	return logsCopy
}

// GetLogCount returns total number of logs
func (al *AuditLogger) GetLogCount() int {
	al.mu.RLock()
	defer al.mu.RUnlock()

	return len(al.logs)
}

// GetEventTimeline returns events in chronological order
func (al *AuditLogger) GetEventTimeline(startTime, endTime time.Time, limit int) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var filtered []*AuditLog
	for _, log := range al.logs {
		if log.Timestamp.After(startTime) && log.Timestamp.Before(endTime) {
			filtered = append(filtered, log)
		}
	}

	// Already in chronological order (appended in order)
	if len(filtered) > limit {
		filtered = filtered[len(filtered)-limit:]
	}

	return filtered
}
