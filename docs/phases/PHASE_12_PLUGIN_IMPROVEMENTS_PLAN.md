# Phase 12: Advanced Plugin Improvements — Strategic Roadmap

**Objective:** Transform plugin system from functional to production-grade with developer tools, performance monitoring, security, and community features.

**Duration:** 8-10 weeks (divided into 4 phases)  
**Priority:** High (enhances ecosystem adoption)

---

## Executive Summary

Current plugin system is solid but lacks:
- Developer experience (no hot-reload, profiling, debugging)
- Performance visibility (no metrics, bottleneck detection)
- Security enforcement (no code signing, audit logging)
- User experience (manual config, no reviews/ratings)
- Operational tooling (no analytics, health monitoring)

This plan addresses all gaps across 4 strategic phases.

---

## Architecture Overview

```
┌────────────────────────────────────────────────────────────┐
│                    Plugin Infrastructure                   │
├────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────┐  ┌──────────────────────────────┐  │
│  │ Phase 12.1       │  │ Phase 12.2                   │  │
│  │ Hot-Reload       │  │ Performance Profiling        │  │
│  │ + Dependency Mgmt│  │ + Analytics Dashboard        │  │
│  └──────────────────┘  └──────────────────────────────┘  │
│         │                          │                       │
│         └──────────┬───────────────┘                       │
│                    ▼                                        │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Phase 12.3: Security + Audit                       │   │
│  │ - Code signing/verification (SHA256)              │   │
│  │ - Permission audit logs                           │   │
│  │ - Runtime sandboxing enforcement                  │   │
│  └────────────────────────────────────────────────────┘   │
│         │                                                   │
│         ▼                                                   │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Phase 12.4: Developer & Community Tools           │   │
│  │ - Plugin template generator                       │   │
│  │ - Testing framework & debug mode                  │   │
│  │ - Config UI auto-generator                        │   │
│  │ - Community review system                         │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
└────────────────────────────────────────────────────────────┘
```

---

## Phase 12.1: Hot-Reload System + Dependency Management

### Duration: 2 weeks

### Features

#### 12.1.1 Plugin Hot-Reload
**Capability:** Reload plugins without restarting application

**Implementation:**
```go
// plugins/hot_reload.go (NEW)
type HotReloadManager struct {
    fileWatcher    *fsnotify.Watcher
    pluginManager  *Manager
    reloadCh       chan ReloadEvent
    stateSnapshot  map[string]PluginState
}

type ReloadEvent struct {
    PluginID      string
    OldVersion    string
    NewVersion    string
    Timestamp     time.Time
    Status        ReloadStatus // Success, Failed, Rolled Back
}

// Methods:
func (h *HotReloadManager) WatchPlugin(pluginID string) error
func (h *HotReloadManager) ReloadPlugin(ctx context.Context, pluginID string) (*ReloadEvent, error)
func (h *HotReloadManager) SnapshotState(pluginID string) (PluginState, error)
func (h *HotReloadManager) RollbackPlugin(pluginID, previousVersion string) error
func (h *HotReloadManager) ValidateReload(old, new *Plugin) error
```

**Files to Create:**
- `plugins/hot_reload.go` (~300 lines)
- `plugins/reload_state.go` (~200 lines)

**Testing:**
- Unit: Reload same plugin version 10x without data loss
- Integration: Reload with dependent plugins, verify deps still work
- E2E: Hot-reload Combat Pack, verify Lua state preserved

#### 12.1.2 Dependency Resolution
**Capability:** Automatic version compatibility checking with transitive deps

**Implementation:**
```go
// plugins/dependency_resolver.go (NEW)
type DependencyResolver struct {
    pluginGraph  *PluginDependencyGraph
    versionCache map[string][]string // pluginID -> versions
}

type PluginDependencyGraph struct {
    nodes map[string]*PluginNode
    edges map[string][]string // pluginID -> dependent pluginIDs
}

type VersionConstraint struct {
    Min      string // "1.0.0"
    Max      string // "2.0.0"
    Operator string // "==", ">=", "<=", "~", "^"
}

// Methods:
func (d *DependencyResolver) ResolveDependencies(pluginID, version string) (map[string]string, error)
func (d *DependencyResolver) DetectConflicts(plugins map[string]string) []ConflictInfo
func (d *DependencyResolver) GetTransitiveDeps(pluginID string) map[string]string
func (d *DependencyResolver) ValidateVersionConstraint(constraint VersionConstraint, version string) bool
```

**Files to Create:**
- `plugins/dependency_resolver.go` (~400 lines)
- `plugins/version_constraint.go` (~150 lines)

**Data Model:**
```json
{
  "plugin_id": "speedrun-tools",
  "version": "2.1.0",
  "dependencies": {
    "ui-framework": "^1.5.0",
    "save-parser": ">=2.0.0,<3.0.0",
    "lua-runtime": "1.1.0"
  }
}
```

**Testing:**
- Semver resolution (caret, tilde operators)
- Circular dependency detection
- Multi-level transitive deps
- Version range merging

---

## Phase 12.2: Performance Profiling + Analytics

### Duration: 2.5 weeks

### Features

#### 12.2.1 Plugin Profiler
**Capability:** Real-time execution metrics for each plugin

**Implementation:**
```go
// plugins/profiler.go (NEW)
type PluginProfiler struct {
    metrics        map[string]*PluginMetrics
    mu             sync.RWMutex
    recordInterval time.Duration
    maxRecords     int
}

type PluginMetrics struct {
    ExecutionCount   int64
    TotalDuration    time.Duration
    AvgDuration      time.Duration
    PeakDuration     time.Duration
    MemoryPeak       uint64
    MemoryAvg        uint64
    LastExecuted     time.Time
    ErrorCount       int64
    WarningCount     int64
}

type ExecutionTrace struct {
    PluginID      string
    FunctionName  string
    StartTime     time.Time
    Duration      time.Duration
    MemoryBefore  uint64
    MemoryAfter   uint64
    Status        string // "success", "error", "timeout"
    Error         string
}

// Methods:
func (p *PluginProfiler) StartProfiling(pluginID string)
func (p *PluginProfiler) StopProfiling(pluginID string) *PluginMetrics
func (p *PluginProfiler) GetMetrics(pluginID string) *PluginMetrics
func (p *PluginProfiler) ExportMetrics(format string) ([]byte, error) // json, csv
func (p *PluginProfiler) DetectBottlenecks(threshold time.Duration) []string
func (p *PluginProfiler) GetExecutionTrace(pluginID string) []*ExecutionTrace
```

**Files to Create:**
- `plugins/profiler.go` (~350 lines)
- `plugins/metrics.go` (~200 lines)

#### 12.2.2 Performance Dashboard UI
**Capability:** Visual dashboard showing plugin performance metrics

**Implementation:**
```go
// ui/forms/plugin_performance_dashboard.go (NEW, ~400 lines)
type PluginPerformanceDashboard struct {
    profiler    *plugins.PluginProfiler
    window      fyne.Window
    refreshRate time.Duration
}

// UI Components:
// - Overall system health (CPU, memory, plugin count)
// - Per-plugin metrics (execution time, memory, errors)
// - Performance trends (graph over time)
// - Bottleneck alerts
// - Performance export (CSV/JSON)

func NewPluginPerformanceDashboard(win fyne.Window, profiler *plugins.PluginProfiler) *PluginPerformanceDashboard
func (d *PluginPerformanceDashboard) Show()
func (d *PluginPerformanceDashboard) RefreshMetrics()
func (d *PluginPerformanceDashboard) ExportMetrics(format string) error
```

**Files to Create:**
- `ui/forms/plugin_performance_dashboard.go` (~400 lines)

#### 12.2.3 Analytics Collection
**Capability:** Aggregate usage statistics (anonymized)

**Implementation:**
```go
// plugins/analytics.go (NEW)
type AnalyticsCollector struct {
    db               map[string]*PluginAnalytics
    persistencePath  string
    collectionRate   time.Duration
}

type PluginAnalytics struct {
    PluginID           string
    InstallCount       int64
    ActiveUsers        int64
    ExecutionCount     int64
    AverageRating      float32
    CommunityScore     float32 // Based on usage + ratings
    CompatibilityScore float32 // Version compatibility tracking
    LastUpdated        time.Time
}

// Methods:
func (a *AnalyticsCollector) RecordExecution(pluginID string, duration time.Duration, success bool)
func (a *AnalyticsCollector) GetAnalytics(pluginID string) *PluginAnalytics
func (a *AnalyticsCollector) GetTrending(limit int) []*PluginAnalytics
func (a *AnalyticsCollector) ExportAnalytics() ([]byte, error)
```

**Files to Create:**
- `plugins/analytics.go` (~250 lines)

**Data Collection (Anonymized):**
- Plugin execution counts
- Average execution time
- Error rates
- Compatibility with app versions
- User retention

---

## Phase 12.3: Security & Audit System

### Duration: 2 weeks

### Features

#### 12.3.1 Code Signing & Verification
**Capability:** Cryptographically verify plugin integrity

**Implementation:**
```go
// plugins/security/signer.go (NEW)
type PluginSigner struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
}

type PluginSignature struct {
    PluginID      string
    Version       string
    Checksum      string // SHA256
    Signature     []byte // RSA-PSS signature
    SignedAt      time.Time
    SignedBy      string // Developer ID
    CertChain     []string
}

// Methods:
func (s *PluginSigner) SignPlugin(path string) (*PluginSignature, error)
func (s *PluginSigner) VerifyPlugin(path string, sig *PluginSignature) (bool, error)
func (s *PluginSigner) ValidateCertChain(sig *PluginSignature) (bool, error)
func (s *PluginSigner) CheckRevocation(pluginID string) (bool, error)
```

**Files to Create:**
- `plugins/security/signer.go` (~300 lines)
- `plugins/security/verifier.go` (~250 lines)

#### 12.3.2 Permission Audit Logging
**Capability:** Track all plugin permission usage

**Implementation:**
```go
// plugins/security/audit_log.go (NEW)
type AuditLog struct {
    records    []*AuditRecord
    mu         sync.RWMutex
    maxRecords int
}

type AuditRecord struct {
    Timestamp     time.Time
    PluginID      string
    Action        string // "ReadSave", "WriteSave", "UIDisplay", etc.
    Resource      string // "character.level", "save.gil", etc.
    Allowed       bool
    Reason        string // If denied
    Context       map[string]interface{}
}

// Methods:
func (a *AuditLog) LogAction(pluginID, action, resource string, allowed bool, reason string)
func (a *AuditLog) GetPluginActions(pluginID string, limit int) []*AuditRecord
func (a *AuditLog) ExportAudit(startTime, endTime time.Time) ([]byte, error)
func (a *AuditLog) GetPermissionViolations(pluginID string) []*AuditRecord
```

**Files to Create:**
- `plugins/security/audit_log.go` (~250 lines)

#### 12.3.3 Runtime Sandboxing Enforcement
**Capability:** Strict enforcement of plugin permissions during execution

**Implementation:**
```go
// plugins/security/sandbox_enforcer.go (NEW)
type SandboxEnforcer struct {
    permissions   map[string][]string // pluginID -> allowed actions
    auditLog      *AuditLog
}

type SandboxViolation struct {
    PluginID    string
    Action      string
    Timestamp   time.Time
    Severity    string // "warning", "error", "critical"
}

// Methods:
func (s *SandboxEnforcer) EnforcePermission(pluginID, action string) error
func (s *SandboxEnforcer) TrackResourceAccess(pluginID, resource string) error
func (s *SandboxEnforcer) GetViolations(pluginID string) []*SandboxViolation
func (s *SandboxEnforcer) RevokePermission(pluginID, action string) error
```

**Files to Create:**
- `plugins/security/sandbox_enforcer.go` (~200 lines)

**UI for Security:**
```go
// ui/forms/plugin_security_panel.go (NEW, ~300 lines)
// Components:
// - Permission grant/revoke interface
// - Audit log viewer
// - Code signature verification status
// - Trust level indicator
// - Security warnings
```

---

## Phase 12.4: Developer Tools + Community Features

### Duration: 2.5 weeks

### Features

#### 12.4.1 Plugin Template Generator
**Capability:** Scaffold new plugins with best practices

**Implementation:**
```go
// plugins/dev/scaffolder.go (NEW)
type PluginScaffolder struct {
    templateDir string
}

type ScaffoldOptions struct {
    PluginID    string
    Name        string
    Author      string
    Description string
    Permissions []string
    Dependencies map[string]string
    UseTypeScript bool // Future: compile TS to Lua
}

// Methods:
func (s *PluginScaffolder) GeneratePlugin(opts ScaffoldOptions, outputPath string) error
func (s *PluginScaffolder) ListTemplates() []PluginTemplate
func (s *PluginScaffolder) GetTemplate(name string) (*PluginTemplate, error)
```

**Template Structure:**
```
my-plugin/
├── plugin.yaml (metadata)
├── manifest.lua (entry point)
├── src/
│   ├── main.lua
│   ├── config.lua
│   └── handlers.lua
├── test/
│   ├── test_main.lua
│   └── fixtures.lua
├── docs/
│   ├── README.md
│   └── API.md
└── package.json (if TypeScript)
```

**Files to Create:**
- `plugins/dev/scaffolder.go` (~300 lines)
- Template files (~500 lines combined)

#### 12.4.2 Plugin Testing Framework
**Capability:** Built-in testing utilities for plugins

**Implementation:**
```go
// plugins/dev/testing.go (NEW)
type PluginTestRunner struct {
    pluginManager *Manager
    fixtures      map[string]*TestFixture
}

type TestFixture struct {
    Name   string
    Setup  func() error
    Data   map[string]interface{}
}

// Methods:
func (t *PluginTestRunner) RunTests(pluginID string) (*TestResults, error)
func (t *PluginTestRunner) RunTest(pluginID, testName string) error
func (t *PluginTestRunner) MockSaveFile() *SaveFileMock
func (t *PluginTestRunner) AssertEqual(expected, actual interface{}) error
func (t *PluginTestRunner) AssertError(err error, expectedMsg string) error
```

**Lua Test Helper Library:**
```lua
-- In plugins, tests would look like:
local test = require("test_framework")

test.case("should increase character level", function()
    local save = test.mock_save()
    local char = save:getCharacter(0)
    
    plugin.increaseLevel(char, 10)
    
    test.assert_equal(char:getLevel(), 10)
end)
```

**Files to Create:**
- `plugins/dev/testing.go` (~350 lines)
- `plugins/dev/lua_testing_lib.go` (~200 lines)

#### 12.4.3 Plugin Config Auto-Generator UI
**Capability:** Auto-generate config UI from JSON schema

**Implementation:**
```go
// ui/forms/plugin_config_builder.go (NEW)
type PluginConfigBuilder struct {
    schema *ConfigSchema
}

type ConfigSchema struct {
    Title       string
    Description string
    Properties  map[string]*ConfigProperty
}

type ConfigProperty struct {
    Type        string // "string", "number", "boolean", "enum", "array"
    Description string
    Default     interface{}
    Validation  *ValidationRule
    UIType      string // "text", "number", "toggle", "dropdown", "color"
}

// Methods:
func (c *PluginConfigBuilder) GenerateUI(schema *ConfigSchema) fyne.CanvasObject
func (c *PluginConfigBuilder) ValidateConfig(config map[string]interface{}, schema *ConfigSchema) error
func (c *PluginConfigBuilder) ExportConfig(config map[string]interface{}) ([]byte, error)
func (c *PluginConfigBuilder) ImportConfig(data []byte) (map[string]interface{}, error)
```

**Example Plugin Config Schema:**
```yaml
title: "Combat Pack Configuration"
properties:
  difficulty:
    type: "enum"
    enum: ["normal", "hard", "extreme"]
    default: "normal"
    description: "Combat difficulty level"
    
  enable_rare_encounters:
    type: "boolean"
    default: false
    description: "Enable rare encounter spawning"
    
  encounter_rate:
    type: "number"
    min: 0.1
    max: 5.0
    default: 1.0
    description: "Encounter rate multiplier"
```

**Files to Create:**
- `ui/forms/plugin_config_builder.go` (~400 lines)
- `plugins/dev/schema_validator.go` (~200 lines)

#### 12.4.4 Plugin Debug Mode
**Capability:** Step-through debugging for plugins

**Implementation:**
```go
// plugins/dev/debugger.go (NEW)
type PluginDebugger struct {
    luaVM          *lua.LState
    breakpoints    map[string][]int
    callStack      []*DebugFrame
    isRunning      bool
    isPaused       bool
}

type DebugFrame struct {
    PluginID     string
    FunctionName string
    LineNumber   int
    LocalVars    map[string]interface{}
}

// Methods:
func (d *PluginDebugger) SetBreakpoint(filename string, lineNumber int) error
func (d *PluginDebugger) Continue() error
func (d *PluginDebugger) Step() error
func (d *PluginDebugger) StepOver() error
func (d *PluginDebugger) GetCallStack() []*DebugFrame
func (d *PluginDebugger) EvaluateExpression(expr string) (interface{}, error)
func (d *PluginDebugger) GetLocalVariables() map[string]interface{}
```

**Files to Create:**
- `plugins/dev/debugger.go` (~300 lines)

#### 12.4.5 Community Review System
**Capability:** In-app ratings and reviews for plugins

**Implementation:**
```go
// ui/forms/plugin_reviews_panel.go (NEW)
type PluginReviewsPanel struct {
    marketplace *marketplace.Client
    window      fyne.Window
}

// UI Components:
// - Star rating (1-5)
// - Text review
// - Helpful votes
// - Version-specific reviews
// - Security audit status badge
// - Compatibility warnings

func NewPluginReviewsPanel(win fyne.Window, marketplace *marketplace.Client) *PluginReviewsPanel
func (p *PluginReviewsPanel) ShowReviews(pluginID string)
func (p *PluginReviewsPanel) SubmitReview(rating int, text string) error
func (p *PluginReviewsPanel) GetTrustScore(pluginID string) float32
```

**Files to Create:**
- `ui/forms/plugin_reviews_panel.go` (~350 lines)

---

## Implementation Timeline

```
Week 1-2:   Phase 12.1 (Hot-Reload + Deps)
Week 3-4:   Phase 12.2 (Profiling + Analytics)  
Week 5-6:   Phase 12.3 (Security + Audit)
Week 7-9:   Phase 12.4 (Dev Tools + Community)
Week 10:    Integration Testing + Polish
```

---

## Dependencies Between Phases

```
12.1 (Hot-Reload)
  ↓
12.2 (Profiling) - depends on hot-reload for metrics
  ↓
12.3 (Security) - independent but benefits from metrics
  ↓
12.4 (Dev Tools) - depends on all previous for complete ecosystem
```

---

## File Count Summary

**Phase 12.1:** 4 new files (~850 lines)  
**Phase 12.2:** 3 new files (~950 lines)  
**Phase 12.3:** 4 new files (~1,000 lines)  
**Phase 12.4:** 7 new files (~2,200 lines)  

**Total:** 18 new files, ~5,000 lines of code

---

## Testing Strategy

### Unit Tests
- Hot-reload state snapshot/restore
- Dependency resolution (semver, conflicts)
- Permission audit logging
- Config schema validation

### Integration Tests
- Hot-reload with dependent plugins
- Profiler with real plugin execution
- Security sandbox enforcements
- Config auto-generation with various schemas

### E2E Tests
- Full plugin lifecycle with hot-reload
- Performance dashboard data accuracy
- Security audit trail completeness
- Developer workflow (scaffold → test → deploy)

---

## Success Metrics

✅ **Developer Experience:**
- Hot-reload time < 500ms
- Plugin scaffold ready in < 30 seconds
- Test runner provides feedback in < 2 seconds

✅ **Performance:**
- Profiling overhead < 2% CPU
- Dashboard refresh rate 1Hz
- Memory tracking accurate to ±5%

✅ **Security:**
- 100% of permissions audit-logged
- Code signature verification < 100ms
- Sandbox violations detected in real-time

✅ **Community:**
- Plugin review submission < 10 seconds
- Trust score calculation < 500ms
- Dependencies resolve < 1 second

---

## Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| Hot-reload breaks plugin state | Snapshot/restore with rollback capability |
| Profiling overhead too high | Optional profiling, batched metrics collection |
| Security sandbox too restrictive | Whitelist-based model, can grant permissions |
| Complexity scares developers | Templates, scaffolder, comprehensive docs |

---

## Next Steps

1. **Approve Phase 12.1** to start hot-reload system
2. **Design plugin template** structure for scaffolder
3. **Define security model** for signature verification
4. **Plan marketplace integration** for reviews/ratings

---

**Status:** Ready for implementation  
**Estimated Effort:** 80-100 developer hours  
**Expected Completion:** 8-10 weeks from start
