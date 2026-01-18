# Phase 12.3 Quick Start Guide

## System Overview

Phase 12.3 adds comprehensive security, audit, and monitoring to the plugin system. All components are automatically initialized when creating a Manager.

## Quick Start

### 1. Create Manager (Auto-initializes all Phase 12.3 components)
```go
import "ffvi_editor/plugins"

manager := plugins.NewManager(pluginDir, api)
// ✅ SecurityManager initialized
// ✅ AuditLogger initialized (10,000 events)
// ✅ SandboxManager initialized
// ✅ PluginProfiler initialized (Phase 12.2)
// ✅ AnalyticsEngine initialized (Phase 12.2)
```

### 2. Load Plugin (Automatic security & audit)
```go
ctx := context.Background()
plugin, err := manager.LoadPlugin(ctx, "/path/to/plugin.go")
// Automatically:
// - Verifies plugin signature
// - Logs load event to audit trail
// - Applies default sandbox policy
// - Records analytics event
```

### 3. Execute Plugin (Automatic monitoring)
```go
err := manager.ExecutePlugin(ctx, "plugin_id")
// Automatically:
// - Checks sandbox permissions
// - Records execution metrics
// - Tracks execution time
// - Logs to audit trail
// - Detects timeout violations
```

### 4. Access Components
```go
// Security
secMgr := manager.GetSecurityManager()
signatures := secMgr.GetSignatures()
stats := secMgr.GetSecurityStats()

// Audit
auditLog := manager.GetAuditLogger()
trail := auditLog.GetPluginAuditTrail("plugin_id")
auditLog.ExportAuditLog("audit.csv")

// Sandbox
sandbox := manager.GetSandboxManager()
violations := sandbox.GetViolations("plugin_id")
report := sandbox.GenerateSecurityReport()

// Performance (Phase 12.2)
profiler := manager.GetProfiler()
aggregates := profiler.GetAggregates()

// Analytics (Phase 12.2)
analytics := manager.GetAnalytics()
allStats := analytics.GetAllPluginStats()
```

## Common Tasks

### Sign a Plugin
```go
secMgr := manager.GetSecurityManager()
signature, err := secMgr.SignPlugin("my_plugin", "/path/to/plugin.go")
// Creates RSA-2048 signature with SHA256 hash
// Valid for 1 year
```

### Set Custom Sandbox Policy
```go
sandbox := manager.GetSandboxManager()
policy := &plugins.SandboxPolicy{
    PluginID:           "my_plugin",
    AllowedPermissions: []string{"read_save", "ui_display"},
    DeniedPermissions:  []string{"network_access"},
    MaxMemoryMB:        50,
    MaxCPUPercent:      30,
    TimeoutSeconds:     15,
    IsolationLevel:     "strict",
    IsActive:           true,
}
sandbox.SetPolicy(policy)
```

### Query Audit Trail
```go
auditLog := manager.GetAuditLogger()

// Get all events for a plugin
trail := auditLog.GetPluginAuditTrail("plugin_id")

// Get events by time range
start := time.Now().Add(-24 * time.Hour)
end := time.Now()
logs := auditLog.GetAuditLog(start, end)

// Get events by type
errors := auditLog.GetEventsByType("error")
violations := auditLog.GetEventsByType("security_violation")

// Export to CSV
auditLog.ExportAuditLog("audit_trail.csv")
auditLog.ExportAuditJSON("audit_trail.json")
```

### Check Permissions
```go
sandbox := manager.GetSandboxManager()

allowed, reason := sandbox.CheckPermission("plugin_id", "write_save")
if !allowed {
    fmt.Printf("Permission denied: %s\n", reason)
}
```

### View Performance Metrics
```go
profiler := manager.GetProfiler()
aggregates := profiler.GetAggregates()

for pluginID, metrics := range aggregates {
    fmt.Printf("Plugin: %s\n", pluginID)
    fmt.Printf("  Executions: %d\n", metrics.ExecutionCount)
    fmt.Printf("  Success Rate: %.2f%%\n", (1.0 - metrics.ErrorRate) * 100)
    fmt.Printf("  Avg Duration: %v\n", metrics.AverageDuration)
    fmt.Printf("  Peak Memory: %.2f MB\n", metrics.PeakMemoryMB)
}
```

### Get Analytics Summary
```go
analytics := manager.GetAnalytics()

// Get all plugin stats
allStats := analytics.GetAllPluginStats()
for pluginID, stats := range allStats {
    fmt.Printf("%s: %d loads, %d executions\n", 
        pluginID, stats.LoadCount, stats.ExecutionCount)
}

// Analyze trends
trends := analytics.AnalyzeTrends("plugin_id", 10) // Last 10 samples
for _, trend := range trends {
    fmt.Printf("Metric: %s, Trend: %s, Change: %.2f%%\n",
        trend.MetricName, trend.Trend, trend.ChangePercent)
}
```

## Default Sandbox Policy

When a plugin is loaded, this policy is automatically applied:

```go
AllowedPermissions: []string{"read_save", "ui_display"}
MaxMemoryMB:        100
MaxCPUPercent:      50
TimeoutSeconds:     30
IsolationLevel:     "basic"
```

## Isolation Levels

### none
- All permissions allowed
- No enforcement
- Use for trusted plugins only

### basic (default)
- Whitelist enforcement only
- Memory and CPU limits
- Timeout enforcement

### strict
- Whitelist AND blacklist enforcement
- Blacklist takes priority
- All resource limits enforced

## Audit Event Types

1. **load** - Plugin loaded
2. **unload** - Plugin unloaded
3. **execute** - Plugin executed
4. **error** - Error occurred
5. **permission_used** - Permission granted and used
6. **security_violation** - Security policy violated

## Security Features

### Signature Verification
- Automatic on plugin load
- RSA-2048 encryption
- SHA256 code hashing
- PEM-encoded certificates
- 1-year validity
- Revocation support

### Trusted Keys
```go
secMgr := manager.GetSecurityManager()

// Add trusted key
err := secMgr.AddTrustedKey("key_id", publicKey)

// Check if key is trusted
if secMgr.IsKeyTrusted("key_id") {
    // Key is trusted
}

// Remove trusted key
err := secMgr.RemoveTrustedKey("key_id")
```

### Revoke Signature
```go
secMgr := manager.GetSecurityManager()
err := secMgr.RevokeSignature("plugin_id", "Security issue detected")
// Plugin will fail verification on next load
```

## Violation Tracking

```go
sandbox := manager.GetSandboxManager()

// Get all violations for a plugin
violations := sandbox.GetViolations("plugin_id")

// Get violations by type
memoryViolations := sandbox.GetViolationsByType("memory_exceeded")
timeoutViolations := sandbox.GetViolationsByType("timeout")

// Get violation statistics
stats := sandbox.GetViolationStats()
fmt.Printf("Total violations: %d\n", stats["total_violations"])
fmt.Printf("Memory violations: %d\n", stats["memory_violations"])
fmt.Printf("Timeout violations: %d\n", stats["timeout_violations"])

// Clear violations
sandbox.ClearViolations()
```

## Dashboard UI (Components Ready)

### Performance Dashboard
```go
import "ffvi_editor/ui/forms"

dashboard := forms.NewPluginPerformanceDashboard(profiler)
widget := dashboard.Show() // Returns fyne.CanvasObject
dashboard.RefreshMetrics() // Manual refresh
dashboard.StartAutoRefresh() // Auto-refresh every 5s
```

### Analytics Dashboard
```go
dashboard := forms.NewPluginAnalyticsDashboard(analytics)
widget := dashboard.Show()
dashboard.RefreshAnalytics()
```

### Alert Manager
```go
alertMgr := forms.NewAlertManager()
alertMgr.AddAlert(&forms.BottleneckAlert{
    AlertID:   "alert_001",
    PluginID:  "my_plugin",
    AlertType: "HIGH_ERROR_RATE",
    Severity:  "CRITICAL",
    Message:   "Error rate exceeds 25%",
})
widget := alertMgr.Show()
```

## Exports

### CSV Export
```go
// Audit log
auditLog.ExportAuditLog("audit.csv")

// Performance metrics
dashboard.exportAsCSV("performance.csv")
```

### JSON Export
```go
// Audit log
auditLog.ExportAuditJSON("audit.json")

// Performance metrics
dashboard.exportAsJSON("performance.json")

// Security signatures
secMgr.ExportSignatures("signatures.json")
```

## Testing

### Run Integration Tests
```powershell
# All manager tests
go test ./plugins -run="TestManager" -v

# Specific test
go test ./plugins -run="TestManagerIntegration" -v

# With coverage
go test ./plugins -run="TestManager" -cover
```

### Build Commands
```powershell
# Build plugins package
go build ./plugins

# Build UI forms
go build ./ui/forms

# Build full application
go build .
```

## Troubleshooting

### Plugin Load Fails
1. Check signature: `secMgr.VerifyPlugin(pluginID, path)`
2. Check audit log: `auditLog.GetPluginAuditTrail(pluginID)`
3. Check violations: `sandbox.GetViolations(pluginID)`

### Permission Denied
1. Check policy: `policy := sandbox.GetPolicy(pluginID)`
2. Check allowed permissions: `policy.AllowedPermissions`
3. Update policy if needed: `sandbox.AllowPermission(pluginID, "permission")`

### Timeout Violations
1. Check execution time: Look at profiler metrics
2. Increase timeout: Update sandbox policy `TimeoutSeconds`
3. Optimize plugin code

### Memory Violations
1. Check memory usage: Look at profiler `PeakMemoryMB`
2. Increase limit: Update sandbox policy `MaxMemoryMB`
3. Optimize plugin memory usage

## Best Practices

1. **Always check audit logs** after operations
2. **Use strict isolation** for untrusted plugins
3. **Monitor violation stats** regularly
4. **Export audit logs** periodically for compliance
5. **Set appropriate timeouts** based on plugin complexity
6. **Review security stats** to identify issues early
7. **Use trusted keys** for production plugins
8. **Revoke signatures immediately** when issues detected

## Performance Tips

1. **Batch operations** when possible
2. **Use auto-refresh wisely** - balance freshness vs performance
3. **Export large audit logs** to CSV instead of querying repeatedly
4. **Clear old violations** periodically
5. **Set reasonable retention** on audit logs (default 90 days)
6. **Monitor profiler overhead** (default 10% sampling)

---

**Quick Command Reference**:
```powershell
# Build
go build ./plugins

# Test
go test ./plugins -run="TestManager" -v

# Coverage
go test ./plugins -cover

# Run specific test
go test ./plugins -run="TestManagerIntegration" -v
```

For detailed API documentation, see PHASE_12.3_INTEGRATION_COMPLETE.md
