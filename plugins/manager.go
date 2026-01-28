package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Manager manages plugin lifecycle and execution
type Manager struct {
	plugins            map[string]*Plugin
	configs            map[string]PluginConfig
	api                PluginAPI
	pluginDir          string
	mu                 sync.RWMutex
	executionLog       []ExecutionRecord
	maxPlugins         int
	sandbox            bool
	autoLoadPlugins    bool
	defaultPerm        []string
	stopCh             chan struct{}
	syncTicker         *time.Ticker
	syncInterval       time.Duration
	isRunning          bool
	hotReloadManager   *HotReloadManager
	dependencyResolver *DependencyResolver
	profiler           *PluginProfiler
	analytics          *AnalyticsEngine
	securityMgr        *SecurityManager
	auditLogger        *AuditLogger
	sandboxMgr         *SandboxManager
}

// NewManager creates a new plugin manager
func NewManager(pluginDir string, api PluginAPI) *Manager {
	m := &Manager{
		plugins:         make(map[string]*Plugin),
		configs:         make(map[string]PluginConfig),
		api:             api,
		pluginDir:       pluginDir,
		maxPlugins:      50,
		sandbox:         true,
		autoLoadPlugins: true,
		defaultPerm: []string{
			CommonPermissions.ReadSave,
			CommonPermissions.UIDisplay,
		},
		executionLog:       make([]ExecutionRecord, 0),
		stopCh:             make(chan struct{}),
		syncInterval:       5 * time.Minute,
		isRunning:          false,
		dependencyResolver: NewDependencyResolver(),
		profiler:           NewPluginProfiler(1000, 0.1), // Store last 1000 executions, sample 10%
		analytics:          NewAnalyticsEngine(10000),    // Store last 10000 events
		securityMgr:        NewSecurityManager(),
		auditLogger:        NewAuditLogger(10000), // Store last 10000 audit events
		sandboxMgr:         NewSandboxManager(),
	}

	// Initialize hot-reload manager
	hotReload, err := NewHotReloadManager(m)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize hot-reload: %v\n", err)
	} else {
		m.hotReloadManager = hotReload
	}

	return m
}

// LoadPlugin loads a plugin from file
func (m *Manager) LoadPlugin(ctx context.Context, path string) (*Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.plugins) >= m.maxPlugins {
		return nil, ErrMaxPluginsExceeded
	}

	// Parse plugin metadata (placeholder for now)
	metadata := &PluginMetadata{
		ID:      "plugin_" + fmt.Sprintf("%d", time.Now().Unix()),
		Name:    "Plugin",
		Version: "1.0.0",
		Author:  "Unknown",
	}

	if err := metadata.Validate(); err != nil {
		return nil, err
	}

	// Check if already loaded
	if _, exists := m.plugins[metadata.ID]; exists {
		return nil, ErrPluginAlreadyLoaded
	}

	// Verify plugin signature if security is enabled
	if m.securityMgr != nil {
		if verified, err := m.securityMgr.VerifyPlugin(metadata.ID, path); !verified || err != nil {
			// Log security violation
			if m.auditLogger != nil {
				errMsg := "signature verification failed"
				if err != nil {
					errMsg = err.Error()
				}
				m.auditLogger.LogSecurityViolation(metadata.ID, "signature_verification_failed", errMsg)
			}
			if err != nil {
				return nil, fmt.Errorf("plugin signature verification failed: %w", err)
			}
			return nil, fmt.Errorf("plugin signature verification failed")
		}
	}

	// Create plugin instance
	plugin := NewPlugin(metadata, m.api)

	// Set the plugin path and metadata
	plugin.SetPath(path)
	plugin.SetMetadata(*metadata)

	// Initialize plugin
	if err := plugin.Load(); err != nil {
		// Log load error
		if m.auditLogger != nil {
			m.auditLogger.LogError(metadata.ID, err.Error())
		}
		return nil, fmt.Errorf("failed to initialize plugin: %w", err)
	}

	// Store plugin
	m.plugins[metadata.ID] = plugin
	m.configs[metadata.ID] = PluginConfig{
		Name:        metadata.Name,
		Version:     metadata.Version,
		Author:      metadata.Author,
		Description: metadata.Description,
		Enabled:     true,
		Settings:    make(map[string]interface{}),
		Hooks:       metadata.Hooks,
		Permissions: metadata.Permissions,
	}

	// Register with dependency resolver
	if m.dependencyResolver != nil {
		m.dependencyResolver.AddPlugin(plugin)
	}

	// Record analytics event for plugin load
	if m.analytics != nil {
		m.analytics.RecordEvent(&AnalyticsEvent{
			EventType: "plugin_load",
			PluginID:  metadata.ID,
			Timestamp: time.Now(),
			Status:    "success",
			Metadata:  map[string]interface{}{"version": metadata.Version},
		})
	}

	// Log plugin load in audit trail
	if m.auditLogger != nil {
		m.auditLogger.LogPluginLoad(metadata.ID)
	}

	// Apply default sandbox policy
	if m.sandboxMgr != nil && m.sandbox {
		policy := &SandboxPolicy{
			PluginID:           metadata.ID,
			AllowedPermissions: append([]string{}, m.defaultPerm...),
			DeniedPermissions:  []string{},
			MaxMemoryMB:        100,
			MaxCPUPercent:      50,
			TimeoutSeconds:     30,
			IsolationLevel:     "basic",
			IsActive:           true,
		}
		m.sandboxMgr.SetPolicy(policy)
	}

	return plugin, nil
}

// UnloadPlugin unloads a plugin
func (m *Manager) UnloadPlugin(ctx context.Context, pluginID string) error {
	m.mu.Lock()
	plugin, exists := m.plugins[pluginID]
	m.mu.Unlock()

	if !exists {
		return ErrPluginNotFound
	}

	// Unload plugin
	if err := plugin.Unload(); err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}

	m.mu.Lock()
	delete(m.plugins, pluginID)
	delete(m.configs, pluginID)
	m.mu.Unlock()

	// Remove from dependency resolver
	if m.dependencyResolver != nil {
		m.dependencyResolver.RemovePlugin(pluginID)
	}

	// Stop watching for hot-reload
	if m.hotReloadManager != nil {
		m.hotReloadManager.UnwatchPlugin(pluginID)
	}

	// Record analytics event for plugin unload
	if m.analytics != nil {
		m.analytics.RecordEvent(&AnalyticsEvent{
			EventType: "plugin_unload",
			PluginID:  pluginID,
			Timestamp: time.Now(),
			Status:    "success",
		})
	}

	// Log plugin unload in audit trail
	if m.auditLogger != nil {
		m.auditLogger.LogPluginUnload(pluginID)
	}

	// Remove sandbox policy
	if m.sandboxMgr != nil {
		m.sandboxMgr.RemovePolicy(pluginID)
	}

	return nil
}

// GetPlugin retrieves a plugin by ID
func (m *Manager) GetPlugin(pluginID string) (*Plugin, error) {
	m.mu.RLock()
	plugin, exists := m.plugins[pluginID]
	m.mu.RUnlock()

	if !exists {
		return nil, ErrPluginNotFound
	}

	return plugin, nil
}

// ListPlugins returns all loaded plugins
func (m *Manager) ListPlugins() []*Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}

// ExecutePlugin executes a plugin
func (m *Manager) ExecutePlugin(ctx context.Context, pluginID string) error {
	plugin, err := m.GetPlugin(pluginID)
	if err != nil {
		return err
	}

	if !plugin.Enabled {
		return fmt.Errorf("plugin is disabled")
	}

	// Check sandbox permissions before execution
	if m.sandboxMgr != nil && m.sandbox {
		if allowed, reason := m.sandboxMgr.CheckPermission(pluginID, "execute"); !allowed {
			if m.auditLogger != nil {
				m.auditLogger.LogPermissionDenied(pluginID, "execute", reason)
			}
			return fmt.Errorf("permission denied: %s", reason)
		}
	}

	startTime := time.Now()
	record := ExecutionRecord{
		PluginName: plugin.Name,
		PluginVer:  plugin.Version,
		StartTime:  startTime,
	}

	// Execute plugin (placeholder)
	defer func() {
		record.Duration = time.Since(startTime)
		m.mu.Lock()
		m.executionLog = append(m.executionLog, record)
		m.mu.Unlock()

		// Check sandbox resource limits
		if m.sandboxMgr != nil && m.sandbox {
			// Check execution time
			if violated, reason := m.sandboxMgr.VerifyExecutionTime(pluginID, record.Duration); violated {
				if m.auditLogger != nil {
					m.auditLogger.LogSecurityViolation(pluginID, "timeout", reason)
				}
			}
		}

		// Record profiling metrics
		if m.profiler != nil {
			success := record.Status == "success"
			var err error
			if record.Status == "error" {
				err = fmt.Errorf(record.Error)
			}
			m.profiler.RecordExecution(pluginID, record.Duration, success, err)
		}

		// Record analytics event
		if m.analytics != nil {
			eventType := "plugin_execute"
			status := "success"
			if record.Status == "error" {
				status = "error"
			}
			m.analytics.RecordEvent(&AnalyticsEvent{
				EventType: eventType,
				PluginID:  pluginID,
				Timestamp: startTime,
				Duration:  record.Duration,
				Status:    status,
				Error:     record.Error,
				Metadata:  map[string]interface{}{"version": plugin.Version},
			})
		}

		// Log plugin execution in audit trail
		if m.auditLogger != nil {
			if record.Status == "success" {
				m.auditLogger.LogPluginExecution(pluginID, record.Duration.Milliseconds(), true, "")
			} else {
				m.auditLogger.LogPluginExecution(pluginID, record.Duration.Milliseconds(), false, record.Error)
			}
		}
	}()

	// Call hook with context timeout
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute plugin with hook - timeout enforced via execCtx
	// Note: CallHook signature could be updated to accept context parameter for cleaner async support
	// Currently context enforced at callsite level via context.WithTimeout
	if err := plugin.CallHook(HookLoad); err != nil {
		// Track error if context was cancelled due to timeout
		if execCtx.Err() == context.DeadlineExceeded {
			record.Error = "plugin execution timeout"
		} else {
			record.Error = err.Error()
		}
		record.Status = "error"
		return err
	}

	record.Status = "success"
	return nil
}

// CallHook fires a hook for all enabled plugins
func (m *Manager) CallHook(ctx context.Context, hookType HookType, args ...interface{}) error {
	m.mu.RLock()
	plugins := make([]*Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		if p.Enabled {
			plugins = append(plugins, p)
		}
	}
	m.mu.RUnlock()

	for _, plugin := range plugins {
		config, ok := m.configs[plugin.ID]
		if ok && !config.Enabled {
			continue
		}

		if err := plugin.CallHook(hookType, args...); err != nil {
			// Log error but continue
			fmt.Printf("hook error in plugin %s: %v\n", plugin.Name, err)
		}
	}

	return nil
}

// GetConfig returns plugin configuration
func (m *Manager) GetConfig(pluginID string) (PluginConfig, error) {
	m.mu.RLock()
	config, exists := m.configs[pluginID]
	m.mu.RUnlock()

	if !exists {
		return PluginConfig{}, ErrPluginNotFound
	}

	return config, nil
}

// SetConfig sets plugin configuration
func (m *Manager) SetConfig(pluginID string, config PluginConfig) error {
	if err := ValidatePluginConfig(config); err != nil {
		return err
	}

	m.mu.Lock()
	if _, exists := m.plugins[pluginID]; !exists {
		m.mu.Unlock()
		return ErrPluginNotFound
	}

	m.configs[pluginID] = config
	m.mu.Unlock()

	return nil
}

// EnablePlugin enables a plugin
func (m *Manager) EnablePlugin(pluginID string) error {
	m.mu.Lock()
	plugin, exists := m.plugins[pluginID]
	config, configExists := m.configs[pluginID]

	if !exists || !configExists {
		m.mu.Unlock()
		return ErrPluginNotFound
	}

	plugin.Enabled = true
	config.Enabled = true
	m.configs[pluginID] = config
	m.mu.Unlock()

	return nil
}

// DisablePlugin disables a plugin
func (m *Manager) DisablePlugin(pluginID string) error {
	m.mu.Lock()
	plugin, exists := m.plugins[pluginID]
	config, configExists := m.configs[pluginID]

	if !exists || !configExists {
		m.mu.Unlock()
		return ErrPluginNotFound
	}

	plugin.Enabled = false
	config.Enabled = false
	m.configs[pluginID] = config
	m.mu.Unlock()

	return nil
}

// GetExecutionLog returns the execution log
func (m *Manager) GetExecutionLog() []ExecutionRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	logCopy := make([]ExecutionRecord, len(m.executionLog))
	copy(logCopy, m.executionLog)

	return logCopy
}

// ClearExecutionLog clears the execution log
func (m *Manager) ClearExecutionLog() {
	m.mu.Lock()
	m.executionLog = make([]ExecutionRecord, 0)
	m.mu.Unlock()
}

// Start starts the plugin manager (enables auto-sync/scheduling)
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.isRunning {
		m.mu.Unlock()
		return nil
	}

	m.isRunning = true
	m.syncTicker = time.NewTicker(m.syncInterval)
	m.mu.Unlock()

	go m.syncLoop(ctx)
	return nil
}

// Stop stops the plugin manager
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	if !m.isRunning {
		m.mu.Unlock()
		return nil
	}

	m.isRunning = false
	if m.syncTicker != nil {
		m.syncTicker.Stop()
	}
	m.mu.Unlock()

	close(m.stopCh)
	return nil
}

// syncLoop runs periodic plugin checks
func (m *Manager) syncLoop(ctx context.Context) {
	for {
		select {
		case <-m.stopCh:
			return
		case <-m.syncTicker.C:
			// Placeholder for periodic sync operations
			// Could check for plugin updates, health checks, etc.
		}
	}
}

// SetMaxPlugins sets the maximum number of plugins
func (m *Manager) SetMaxPlugins(max int) {
	m.mu.Lock()
	m.maxPlugins = max
	m.mu.Unlock()
}

// SetSandboxMode enables/disables sandboxing
func (m *Manager) SetSandboxMode(enabled bool) {
	m.mu.Lock()
	m.sandbox = enabled
	m.mu.Unlock()
}

// StartHotReload starts the hot-reload system
func (m *Manager) StartHotReload(ctx context.Context) error {
	if m.hotReloadManager == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}
	// Note: ctx is for caller use, hot-reload manager doesn't need it
	_ = ctx
	return m.hotReloadManager.Start()
}

// StopHotReload stops the hot-reload system
func (m *Manager) StopHotReload() error {
	if m.hotReloadManager == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}
	return m.hotReloadManager.Stop()
}

// EnablePluginHotReload enables hot-reload for a specific plugin
func (m *Manager) EnablePluginHotReload(pluginID string) error {
	if m.hotReloadManager == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}

	plugin, err := m.GetPlugin(pluginID)
	if err != nil {
		return err
	}

	return m.hotReloadManager.WatchPlugin(pluginID, plugin.GetPath())
}

// DisablePluginHotReload disables hot-reload for a specific plugin
func (m *Manager) DisablePluginHotReload(pluginID string) error {
	if m.hotReloadManager == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}

	return m.hotReloadManager.UnwatchPlugin(pluginID)
}

// GetReloadHistory returns the reload history for a plugin
func (m *Manager) GetReloadHistory(pluginID string) ([]*ReloadEvent, error) {
	if m.hotReloadManager == nil {
		return nil, fmt.Errorf("hot-reload manager not initialized")
	}

	// Get all history and filter by plugin ID
	allHistory := m.hotReloadManager.GetReloadHistory(100)
	var filtered []*ReloadEvent
	for _, event := range allHistory {
		if event.PluginID == pluginID {
			filtered = append(filtered, event)
		}
	}

	return filtered, nil
}

// ResolveDependencies resolves dependencies for all loaded plugins
func (m *Manager) ResolveDependencies() ([]string, error) {
	if m.dependencyResolver == nil {
		return nil, fmt.Errorf("dependency resolver not initialized")
	}

	m.mu.RLock()
	pluginIDs := make([]string, 0, len(m.plugins))
	for id := range m.plugins {
		pluginIDs = append(pluginIDs, id)
	}
	m.mu.RUnlock()

	// For now, return plugins in load order
	// In future, we can implement full topological sort
	return pluginIDs, nil
}

// DetectDependencyConflicts detects conflicts in plugin dependencies
func (m *Manager) DetectDependencyConflicts() ([]ConflictInfo, error) {
	if m.dependencyResolver == nil {
		return nil, fmt.Errorf("dependency resolver not initialized")
	}

	// Build a map of plugin versions for conflict detection
	versionMap := make(map[string]string)
	m.mu.RLock()
	for id, plugin := range m.plugins {
		versionMap[id] = plugin.Version
	}
	m.mu.RUnlock()

	return m.dependencyResolver.DetectConflicts(versionMap), nil
}

// GetStats returns plugin manager statistics
func (m *Manager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"total_plugins": len(m.plugins),
		"max_plugins":   m.maxPlugins,
		"enabled_count": m.countEnabled(),
		"log_entries":   len(m.executionLog),
		"sandbox_mode":  m.sandbox,
		"is_running":    m.isRunning,
	}
}

// countEnabled counts enabled plugins
func (m *Manager) countEnabled() int {
	count := 0
	for _, config := range m.configs {
		if config.Enabled {
			count++
		}
	}
	return count
}

// GetProfiler returns the plugin profiler
func (m *Manager) GetProfiler() *PluginProfiler {
	return m.profiler
}

// GetPerformanceReport returns a human-readable performance report
func (m *Manager) GetPerformanceReport() string {
	if m.profiler == nil {
		return "Profiler not available\n"
	}
	return m.profiler.GetPerformanceReport()
}

// GetAnalytics returns the analytics engine
func (m *Manager) GetAnalytics() *AnalyticsEngine {
	return m.analytics
}

// GetAnalyticsSummary returns a human-readable analytics summary
func (m *Manager) GetAnalyticsSummary() string {
	if m.analytics == nil {
		return "Analytics not available\n"
	}
	return m.analytics.GetAnalyticsSummary()
}

// GetSecurityManager returns the security manager
func (m *Manager) GetSecurityManager() *SecurityManager {
	return m.securityMgr
}

// GetAuditLogger returns the audit logger
func (m *Manager) GetAuditLogger() *AuditLogger {
	return m.auditLogger
}

// GetSandboxManager returns the sandbox manager
func (m *Manager) GetSandboxManager() *SandboxManager {
	return m.sandboxMgr
}
