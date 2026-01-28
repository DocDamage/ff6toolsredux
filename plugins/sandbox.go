package plugins

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// SandboxPolicy defines runtime constraints for a plugin
type SandboxPolicy struct {
	PluginID           string
	AllowedPermissions []string // Whitelist of permitted actions
	DeniedPermissions  []string // Blacklist of forbidden actions
	MaxMemoryMB        int      // Memory limit in MB (0 = unlimited)
	MaxCPUPercent      int      // CPU usage limit as percentage (0 = unlimited)
	TimeoutSeconds     int      // Execution timeout in seconds (0 = unlimited)
	IsolationLevel     string   // "none", "basic", "strict"
	CreatedAt          time.Time
	ModifiedAt         time.Time
	IsActive           bool
}

// SecurityViolation records a sandbox policy violation
type SecurityViolation struct {
	ViolationID   string
	PluginID      string
	ViolationType string // "permission_denied", "memory_exceeded", "cpu_exceeded", "timeout", "isolation_breach"
	Message       string
	Timestamp     time.Time
	Severity      string // "CRITICAL", "WARNING", "INFO"
	Enforced      bool   // Whether violation was prevented
	Details       map[string]interface{}
}

// SandboxManager enforces sandbox policies
type SandboxManager struct {
	policies      map[string]*SandboxPolicy // pluginID -> policy
	violations    []*SecurityViolation
	mu            sync.RWMutex
	maxViolations int
}

// NewSandboxManager creates a new sandbox manager
func NewSandboxManager() *SandboxManager {
	return &SandboxManager{
		policies:      make(map[string]*SandboxPolicy),
		violations:    make([]*SecurityViolation, 0, 10000),
		maxViolations: 10000,
	}
}

// SetPolicy defines sandbox policy for a plugin
func (sm *SandboxManager) SetPolicy(policy *SandboxPolicy) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if policy.PluginID == "" {
		return fmt.Errorf("plugin ID is required")
	}

	// Validate isolation level
	validLevels := map[string]bool{
		"none":   true,
		"basic":  true,
		"strict": true,
	}
	if !validLevels[policy.IsolationLevel] {
		return fmt.Errorf("invalid isolation level: %s", policy.IsolationLevel)
	}

	policy.CreatedAt = time.Now()
	policy.ModifiedAt = time.Now()
	policy.IsActive = true

	sm.policies[policy.PluginID] = policy

	return nil
}

// GetPolicy retrieves sandbox policy for a plugin
func (sm *SandboxManager) GetPolicy(pluginID string) *SandboxPolicy {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.policies[pluginID]
}

// RemovePolicy removes sandbox policy for a plugin
func (sm *SandboxManager) RemovePolicy(pluginID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.policies[pluginID]; !exists {
		return fmt.Errorf("policy not found for plugin: %s", pluginID)
	}

	delete(sm.policies, pluginID)
	return nil
}

// CheckPermission verifies if plugin is allowed to perform action
func (sm *SandboxManager) CheckPermission(pluginID, permission string) (bool, string) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	policy := sm.policies[pluginID]
	if policy == nil {
		// No policy = allow by default (permissive mode)
		return true, ""
	}

	// Check denied permissions (blacklist)
	for _, denied := range policy.DeniedPermissions {
		if denied == permission || denied == "*" {
			sm.mu.RUnlock()
			sm.logViolation(pluginID, "permission_denied",
				fmt.Sprintf("Permission '%s' is denied", permission), "CRITICAL", true)
			sm.mu.RLock()
			return false, fmt.Sprintf("permission '%s' is denied by policy", permission)
		}
	}

	// If allowed list is specified, check it (whitelist)
	if len(policy.AllowedPermissions) > 0 {
		allowed := false
		for _, allowed_perm := range policy.AllowedPermissions {
			if allowed_perm == permission || allowed_perm == "*" {
				allowed = true
				break
			}
		}
		if !allowed {
			sm.mu.RUnlock()
			sm.logViolation(pluginID, "permission_denied",
				fmt.Sprintf("Permission '%s' not in allowed list", permission), "WARNING", true)
			sm.mu.RLock()
			return false, fmt.Sprintf("permission '%s' not allowed by policy", permission)
		}
	}

	return true, ""
}

// VerifyMemoryUsage checks if memory usage is within limits
func (sm *SandboxManager) VerifyMemoryUsage(pluginID string, memoryBytes int64) (bool, string) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	policy := sm.policies[pluginID]
	if policy == nil || policy.MaxMemoryMB == 0 {
		// No limit or no policy
		return true, ""
	}

	memoryMB := int(memoryBytes / 1024 / 1024)
	if memoryMB > policy.MaxMemoryMB {
		sm.mu.RUnlock()
		sm.logViolation(pluginID, "memory_exceeded",
			fmt.Sprintf("Memory usage %dMB exceeds limit %dMB", memoryMB, policy.MaxMemoryMB),
			"CRITICAL", false) // Not enforced yet
		sm.mu.RLock()
		return false, fmt.Sprintf("memory limit exceeded: %dMB > %dMB", memoryMB, policy.MaxMemoryMB)
	}

	return true, ""
}

// VerifyExecutionTime checks if execution time is within timeout
func (sm *SandboxManager) VerifyExecutionTime(pluginID string, duration time.Duration) (bool, string) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	policy := sm.policies[pluginID]
	if policy == nil || policy.TimeoutSeconds == 0 {
		// No timeout
		return true, ""
	}

	timeoutDuration := time.Duration(policy.TimeoutSeconds) * time.Second
	if duration > timeoutDuration {
		sm.mu.RUnlock()
		sm.logViolation(pluginID, "timeout",
			fmt.Sprintf("Execution time %s exceeds timeout %s", duration, timeoutDuration),
			"CRITICAL", false)
		sm.mu.RLock()
		return false, fmt.Sprintf("execution timeout: %s > %s", duration, timeoutDuration)
	}

	return true, ""
}

// DenyPermission explicitly denies a permission
func (sm *SandboxManager) DenyPermission(pluginID, permission string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	policy := sm.policies[pluginID]
	if policy == nil {
		policy = &SandboxPolicy{
			PluginID:           pluginID,
			IsolationLevel:     "basic",
			AllowedPermissions: []string{},
			DeniedPermissions:  []string{permission},
		}
		sm.policies[pluginID] = policy
	} else {
		// Add to denied list if not already there
		found := false
		for _, d := range policy.DeniedPermissions {
			if d == permission {
				found = true
				break
			}
		}
		if !found {
			policy.DeniedPermissions = append(policy.DeniedPermissions, permission)
		}
		policy.ModifiedAt = time.Now()
	}

	return nil
}

// AllowPermission explicitly allows a permission
func (sm *SandboxManager) AllowPermission(pluginID, permission string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	policy := sm.policies[pluginID]
	if policy == nil {
		policy = &SandboxPolicy{
			PluginID:           pluginID,
			IsolationLevel:     "basic",
			AllowedPermissions: []string{permission},
			DeniedPermissions:  []string{},
		}
		sm.policies[pluginID] = policy
	} else {
		// Add to allowed list if not already there
		found := false
		for _, a := range policy.AllowedPermissions {
			if a == permission {
				found = true
				break
			}
		}
		if !found {
			policy.AllowedPermissions = append(policy.AllowedPermissions, permission)
		}
		policy.ModifiedAt = time.Now()
	}

	return nil
}

// GetViolations returns violations for a plugin
func (sm *SandboxManager) GetViolations(pluginID string) []*SecurityViolation {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var violations []*SecurityViolation
	for _, v := range sm.violations {
		if v.PluginID == pluginID {
			violations = append(violations, v)
		}
	}
	return violations
}

// GetAllViolations returns all violations
func (sm *SandboxManager) GetAllViolations() []*SecurityViolation {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	violationsCopy := make([]*SecurityViolation, len(sm.violations))
	copy(violationsCopy, sm.violations)
	return violationsCopy
}

// GetViolationsByType returns violations of specific type
func (sm *SandboxManager) GetViolationsByType(violationType string) []*SecurityViolation {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var violations []*SecurityViolation
	for _, v := range sm.violations {
		if v.ViolationType == violationType {
			violations = append(violations, v)
		}
	}
	return violations
}

// ClearViolations clears violation history
func (sm *SandboxManager) ClearViolations() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.violations = make([]*SecurityViolation, 0, 10000)
}

// GetViolationStats returns statistics about violations
func (sm *SandboxManager) GetViolationStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := map[string]interface{}{
		"total_violations":    len(sm.violations),
		"enforced":            0,
		"not_enforced":        0,
		"by_type":             make(map[string]int),
		"by_severity":         make(map[string]int),
		"critical_violations": 0,
	}

	enforced := 0
	notEnforced := 0
	byType := make(map[string]int)
	bySeverity := make(map[string]int)
	criticalCount := 0

	for _, v := range sm.violations {
		if v.Enforced {
			enforced++
		} else {
			notEnforced++
		}

		byType[v.ViolationType]++
		bySeverity[v.Severity]++

		if v.Severity == "CRITICAL" {
			criticalCount++
		}
	}

	stats["enforced"] = enforced
	stats["not_enforced"] = notEnforced
	stats["by_type"] = byType
	stats["by_severity"] = bySeverity
	stats["critical_violations"] = criticalCount

	return stats
}

// logViolation adds a security violation record
func (sm *SandboxManager) logViolation(pluginID, violationType, message, severity string, enforced bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	violation := &SecurityViolation{
		ViolationID:   fmt.Sprintf("viol_%d", time.Now().UnixNano()),
		PluginID:      pluginID,
		ViolationType: violationType,
		Message:       message,
		Timestamp:     time.Now(),
		Severity:      severity,
		Enforced:      enforced,
		Details:       make(map[string]interface{}),
	}

	sm.violations = append(sm.violations, violation)

	// Enforce max history size
	if len(sm.violations) > sm.maxViolations {
		sm.violations = sm.violations[len(sm.violations)-sm.maxViolations:]
	}
}

// GenerateSecurityReport generates comprehensive security report
func (sm *SandboxManager) GenerateSecurityReport() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	report := fmt.Sprintf("=== Sandbox Security Report ===\n")
	report += fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339))

	// Policy summary
	report += fmt.Sprintf("Active Policies: %d\n\n", len(sm.policies))

	// Violation summary
	stats := sm.GetViolationStats()
	report += fmt.Sprintf("--- Violations ---\n")
	report += fmt.Sprintf("Total: %d\n", stats["total_violations"])
	report += fmt.Sprintf("Enforced: %d\n", stats["enforced"])
	report += fmt.Sprintf("Not Enforced: %d\n", stats["not_enforced"])
	report += fmt.Sprintf("Critical: %d\n\n", stats["critical_violations"])

	// Most violated plugins
	report += fmt.Sprintf("--- Most Violated Plugins (Top 10) ---\n")

	type pluginViolCount struct {
		id    string
		count int
	}

	pluginViolations := make(map[string]int)
	for _, v := range sm.violations {
		pluginViolations[v.PluginID]++
	}

	var violations []pluginViolCount
	for id, count := range pluginViolations {
		violations = append(violations, pluginViolCount{id, count})
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].count > violations[j].count
	})

	if len(violations) > 10 {
		violations = violations[:10]
	}

	for i, v := range violations {
		report += fmt.Sprintf("%d. %s: %d violations\n", i+1, v.id, v.count)
	}

	return report
}

// GetPolicies returns all active policies
func (sm *SandboxManager) GetPolicies() map[string]*SandboxPolicy {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	policiesCopy := make(map[string]*SandboxPolicy)
	for k, v := range sm.policies {
		policiesCopy[k] = v
	}
	return policiesCopy
}

// GetPolicyCount returns number of active policies
func (sm *SandboxManager) GetPolicyCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.policies)
}
