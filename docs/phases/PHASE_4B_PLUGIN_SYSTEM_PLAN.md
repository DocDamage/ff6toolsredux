# Phase 4B: Plugin System Architecture - Implementation Plan

**Date:** January 16, 2026  
**Status:** Ready to Start  
**Estimated Duration:** 1 comprehensive session  
**Complexity:** High (core architecture + Go plugin system)

---

## Overview

Phase 4B implements a dynamic plugin system allowing users to extend the FF6 Save Editor with custom functionality through Lua scripting. This phase creates the infrastructure for loading, managing, and executing plugins with safe sandboxing and comprehensive APIs.

---

## Objectives

### Primary Goals
1. ✅ Create plugin interface and manager
2. ✅ Implement Lua VM with sandboxed execution
3. ✅ Build comprehensive plugin API for save data manipulation
4. ✅ Create plugin discovery and loading system
5. ✅ Implement plugin lifecycle management
6. ✅ Build plugin repository/marketplace backend
7. ✅ Create plugin UI (browser, installer, manager)
8. ✅ Write comprehensive documentation and examples

### Deliverables
- **Core Plugin System:** 400+ lines (plugin interface, manager, loader)
- **Lua VM Integration:** 300+ lines (bindings, sandboxing, execution)
- **Plugin API:** 400+ lines (save data access, UI hooks, events)
- **Plugin Manager UI:** 350+ lines (discovery, install, settings)
- **Unit Tests:** 250+ lines (plugin loading, execution, API)
- **Documentation:** 800+ lines (user guide, API reference, examples)

---

## Architecture Design

### 1. Plugin System (plugins/plugin.go)

**Plugin Interface**
```go
type Plugin interface {
    // Metadata
    Name() string                    // Plugin name
    Version() string                 // Semantic version
    Author() string                  // Creator
    Description() string             // What it does
    
    // Lifecycle
    Initialize(ctx context.Context, api PluginAPI) error
    Execute(ctx context.Context) error
    Shutdown(ctx context.Context) error
    
    // Configuration
    GetConfig() PluginConfig
    SetConfig(config PluginConfig) error
    ValidateConfig() error
}
```

**PluginConfig Structure**
```go
type PluginConfig struct {
    Name        string
    Version     string
    Author      string
    Description string
    Enabled     bool
    Settings    map[string]interface{}
    Hooks       []string  // UI hooks: "on_save", "on_load", "on_export"
    Permissions []string  // API permissions: "read_save", "write_save", "ui_display"
}
```

**PluginMetadata (Discovery)**
```go
type PluginMetadata struct {
    ID            string
    Name          string
    Version       string
    Author        string
    Description   string
    Repository    string
    License       string
    Tags          []string
    Dependencies  []string
    MinEditorVer  string
    MaxEditorVer  string
    DownloadURL   string
    LastUpdated   time.Time
    Rating        float64
    Downloads     int
}
```

---

### 2. Plugin API (plugins/api.go)

**PluginAPI Interface** - Safe access to save editor functionality

```go
type PluginAPI interface {
    // Save Data Access
    GetCharacter(name string) (*Character, error)
    SetCharacter(name string, ch *Character) error
    GetInventory() (*Inventory, error)
    SetInventory(inv *Inventory) error
    GetParty() (*Party, error)
    SetParty(party *Party) error
    
    // Batch Operations
    ApplyBatchOperation(op string, params map[string]interface{}) (int, error)
    
    // Queries
    FindCharacter(predicate func(*Character) bool) *Character
    FindItems(predicate func(*Item) bool) []*Item
    
    // Events
    RegisterHook(event string, callback func(interface{}) error) error
    FireEvent(event string, data interface{}) error
    
    // UI
    ShowDialog(title, message string) error
    ShowConfirm(title, message string) bool
    ShowInput(prompt string) (string, error)
    
    // Logging
    Log(level string, message string) error
    
    // Settings
    GetSetting(key string) interface{}
    SetSetting(key string, value interface{}) error
}
```

---

### 3. Plugin Manager (plugins/manager.go)

**Manager Responsibilities**
- Load plugins from disk
- Validate plugin configurations
- Manage plugin lifecycle
- Coordinate API access
- Handle plugin execution
- Track errors and metrics

```go
type Manager struct {
    plugins       map[string]Plugin
    configs       map[string]PluginConfig
    api           PluginAPI
    pluginDir     string
    mu            sync.RWMutex
    executionLog  []ExecutionRecord
    maxPlugins    int
    sandbox       bool  // Enable sandboxing
}

// Key Methods:
// LoadPlugin(path string) error
// UnloadPlugin(name string) error
// ExecutePlugin(name string, ctx context.Context) error
// ListPlugins() []PluginMetadata
// GetPluginConfig(name string) PluginConfig
// SetPluginConfig(name string, cfg PluginConfig) error
// ValidatePlugin(path string) error
```

---

### 4. Lua VM Integration (scripting/vm.go)

**Lua VM Wrapper** - Sandboxed Lua execution with safe bindings

```go
type VM struct {
    lua           *lua.LState
    api           PluginAPI
    timeout       time.Duration
    maxMemory     int
    modules       map[string]bool  // Allowed modules: table, string, math, utf8
}

// Key Methods:
// New(api PluginAPI) *VM
// Execute(code string) error
// ExecuteFile(path string) error
// Call(function string, args ...interface{}) (interface{}, error)
// SetGlobal(name string, value interface{}) error
// GetGlobal(name string) interface{}
// Close()
```

**Lua Standard Library** - Restricted access

```lua
-- Available modules:
-- table.*  (insert, remove, concat, sort, unpack)
-- string.* (sub, len, find, format, upper, lower, gsub)
-- math.*   (floor, ceil, abs, min, max, random, sin, cos)
-- utf8.*   (len, codes, codepoint)

-- API exposure in global scope:
-- editor.getCharacter(name)
-- editor.setCharacter(name, data)
-- editor.getInventory()
-- editor.applyBatch(op, params)
-- editor.showDialog(title, message)
-- editor.log(level, message)
```

---

### 5. Plugin Discovery & Registry

**Local Plugin Discovery**
- Location: `~/.ff6editor/plugins/` and `./plugins/`
- Format: Lua files (.lua) with metadata header
- Structure:
  ```
  local plugin = {
      name = "Character Max Stats",
      version = "1.0.0",
      author = "Community",
      description = "Quick maxing tool",
      execute = function(api)
          -- Plugin code
      end
  }
  ```

**Plugin Registry** (marketplace backend - Phase 4C)
- Planned structure for Phase 4C
- Will include: ratings, downloads, reviews, hosting

---

### 6. Plugin UI (ui/forms/plugin_manager.go)

**Plugin Manager Dialog**
- 4-tab interface

**Tab 1: Installed Plugins**
- List of loaded plugins
- Enable/disable toggles
- Settings button per plugin
- Remove button
- Status indicator

**Tab 2: Available Plugins**
- Browse all plugins (from registry)
- Search and filter
- Install button
- Rating and download count
- Description preview

**Tab 3: Plugin Settings**
- Global plugin configuration
- Auto-load on startup toggle
- Sandbox security level
- Max concurrent plugins
- Max execution time per plugin

**Tab 4: Plugin Output**
- Real-time plugin execution log
- Filter by level (DEBUG, INFO, WARN, ERROR)
- Search log
- Clear log button
- Export log

---

## Implementation Phases

### Phase 4B.1: Core Plugin System (2-3 hours)

**Files to Create:**
1. `plugins/plugin.go` - Plugin interface and config structures
2. `plugins/manager.go` - Plugin manager with loader
3. `plugins/api.go` - Plugin API interface and implementation

**Key Components:**
- Plugin interface (8 methods)
- PluginConfig structure (6 fields)
- PluginAPI interface (18 methods)
- Manager with LoadPlugin, UnloadPlugin, ExecutePlugin
- Error handling and validation

**Tests:**
- Plugin interface implementation
- Config validation
- Error cases

---

### Phase 4B.2: Lua VM Integration (1-2 hours)

**Files to Create:**
1. `scripting/vm.go` - Lua VM wrapper
2. `scripting/bindings.go` - Go↔Lua bindings
3. `scripting/stdlib.go` - Safe standard library

**Key Components:**
- Lua state creation and cleanup
- Safe binding of PluginAPI methods to Lua globals
- Restricted standard library exposure
- Error message translation
- Memory/timeout limits

**Tests:**
- Lua code execution
- Binding correctness
- Error handling
- Timeout verification

---

### Phase 4B.3: Plugin Discovery & Loading (1 hour)

**Files to Create/Modify:**
1. `plugins/loader.go` - File-based plugin loader
2. Update `plugins/manager.go` - Add discovery

**Key Components:**
- Scan `~/.ff6editor/plugins/` directory
- Parse plugin metadata from files
- Validate syntax before loading
- Handle dependencies
- Error recovery

---

### Phase 4B.4: Plugin Manager UI (1.5-2 hours)

**Files to Create:**
1. `ui/forms/plugin_manager.go` - Full dialog implementation
2. Update `ui/window.go` - Add Tools menu item

**Key Components:**
- 4-tab interface
- Plugin list with controls
- Registry browser
- Settings configuration
- Execution log display

---

### Phase 4B.5: Documentation & Examples (1.5 hours)

**Files to Create:**
1. `PHASE_4B_PLUGIN_GUIDE.md` - User guide (400+ lines)
2. `PHASE_4B_API_REFERENCE.md` - Complete API docs (300+ lines)
3. `PHASE_4B_PLUGIN_EXAMPLES.md` - 5-10 example plugins (200+ lines)
4. `PHASE_4B_COMPLETION_SUMMARY.md` - Technical summary (200+ lines)

---

## Data Structures

### Plugin Metadata
```go
type PluginMetadata struct {
    ID            string                 // Unique identifier
    Name          string                 // Display name
    Version       string                 // Semantic version
    Author        string                 // Creator
    Description   string                 // Purpose
    Repository    string                 // Source URL
    License       string                 // License type
    Tags          []string               // Categories
    Dependencies  map[string]string      // Required plugins
    MinEditorVer  string                 // Minimum version
    MaxEditorVer  string                 // Maximum version
    DownloadURL   string                 // Download location
    LastUpdated   time.Time              // Last modification
    Rating        float64                // 1.0-5.0 stars
    Downloads     int                    // Download count
}
```

### Execution Record
```go
type ExecutionRecord struct {
    PluginName  string
    StartTime   time.Time
    Duration    time.Duration
    Status      string  // "success", "error", "timeout"
    Output      string  // Log messages
    Error       string  // Error message if any
}
```

---

## Example Plugin (Lua)

```lua
-- Character Max Stats Plugin
local plugin = {
    name = "Character Max Stats",
    version = "1.0.0",
    author = "Community",
    description = "Quickly max out all character stats"
}

function plugin.execute(api)
    -- Get first character
    local char = api.getCharacter("Terra")
    if not char then
        api.log("error", "Character not found")
        return
    end
    
    -- Max stats
    char.HP = 9999
    char.MP = 999
    char.Level = 99
    
    -- Set it back
    api.setCharacter("Terra", char)
    
    api.log("info", "Character stats maximized")
    api.showDialog("Success", "Terra's stats maximized!")
end

return plugin
```

---

## Security Considerations

### Sandboxing
- ✅ Restrict Lua modules (no file I/O, no shell access, no os module)
- ✅ Limit memory usage per plugin
- ✅ Set execution timeout
- ✅ API permission system (read-only, write, UI)
- ✅ No direct file system access
- ✅ No network access from plugins

### Validation
- ✅ Validate plugin before loading
- ✅ Check dependencies
- ✅ Verify version compatibility
- ✅ Signature verification (Phase 4C)

### Permissions
```go
type Permission string

const (
    PermReadSave    Permission = "read_save"
    PermWriteSave   Permission = "write_save"
    PermUIDisplay   Permission = "ui_display"
    PermEvents      Permission = "events"
)
```

---

## Integration Points

### Main Window Integration
```go
// In ui/window.go:
func (w *Window) setupToolsMenu() {
    // ... existing code ...
    
    pluginItem := fyne.NewMenuItem("Plugin Manager", func() {
        dialog := NewPluginManagerDialog(w.window, w.cloudManager, w.pluginManager)
        dialog.Show()
    })
    w.toolsMenu.Items = append(w.toolsMenu.Items, pluginItem)
}
```

### Save Hook
```go
// When saving:
func (w *Window) saveFile() error {
    // ... existing save code ...
    
    // Fire "on_save" hook for plugins
    w.pluginManager.FireHook("on_save", map[string]interface{}{
        "filename": filename,
        "pr": w.prData,
    })
    
    return nil
}
```

---

## Testing Strategy

### Unit Tests
- Plugin interface implementation
- Config validation
- Lua code execution
- API method calls
- Error handling
- Timeout enforcement

### Integration Tests
- Plugin loading and execution
- API calls from Lua
- Event firing
- Hook registration

### Example Tests
- "Max Stats" plugin runs and modifies character
- "Batch Rename" plugin renames items
- "Report Generator" plugin creates export

---

## Deliverables Checklist

### Code Files (3 new, 2 modified)
- [ ] `plugins/plugin.go` - Plugin interface and structures
- [ ] `plugins/manager.go` - Plugin manager
- [ ] `plugins/api.go` - Plugin API implementation
- [ ] `scripting/vm.go` - Lua VM wrapper
- [ ] `scripting/bindings.go` - Go-Lua bindings
- [ ] `scripting/stdlib.go` - Safe standard library
- [ ] `ui/forms/plugin_manager.go` - Plugin UI dialog
- [ ] `plugins/plugins_test.go` - Unit tests
- [ ] Update `ui/window.go` - Add plugin menu
- [ ] Update `scripting/vm.go` - Add to existing scripting module

### Documentation Files (4 new)
- [ ] `PHASE_4B_PLUGIN_GUIDE.md` - User guide
- [ ] `PHASE_4B_API_REFERENCE.md` - API documentation
- [ ] `PHASE_4B_PLUGIN_EXAMPLES.md` - Example plugins
- [ ] `PHASE_4B_COMPLETION_SUMMARY.md` - Technical summary

### Tests (1 new)
- [ ] `plugins/plugins_test.go` - 250+ lines
- [ ] Benchmarks for plugin loading and execution

---

## Success Criteria

- ✅ All plugin interfaces defined and testable
- ✅ Lua VM executes safely with sandboxing
- ✅ Plugin API provides comprehensive save data access
- ✅ Plugins can be discovered and loaded from disk
- ✅ UI dialog shows plugin browser and settings
- ✅ Example plugins work correctly
- ✅ Documentation complete and clear
- ✅ Zero security vulnerabilities
- ✅ All unit tests pass
- ✅ Plugin execution logged and monitorable

---

## Known Considerations

### OAuth2 Requirement
- Plugins cannot access cloud credentials (security)
- Cloud sync hooks will be available as "on_sync" event
- Marketplace integration (Phase 4C) will handle publishing

### Marketplace Preview
- Phase 4C will add cloud-based registry
- Plugin publishing and versioning
- User ratings and reviews
- Automated plugin updates

### Phase 4D (CLI) Preview
- CLI will support `ff6editor plugin list`, `plugin install`, `plugin run`
- Headless plugin execution support

### Phase 4E (Advanced Scripting) Preview
- More Lua libraries (JSON parsing, pattern matching)
- Hot-reload support
- Debug mode with breakpoints

---

## Estimated Timeline

| Phase | Duration | Status |
|-------|----------|--------|
| 4B.1 - Core Plugin System | 2-3 hours | Ready to start |
| 4B.2 - Lua VM Integration | 1-2 hours | After 4B.1 |
| 4B.3 - Plugin Discovery | 1 hour | After 4B.2 |
| 4B.4 - Plugin Manager UI | 1.5-2 hours | After 4B.3 |
| 4B.5 - Documentation | 1.5 hours | Parallel/Final |
| **Total Phase 4B** | **7-9 hours** | **Single Session** |

---

## Success Definition

Phase 4B is complete when:
1. ✅ Plugin interface fully defined with metadata
2. ✅ Lua VM executes plugins safely
3. ✅ Plugin API exposes save data for modification
4. ✅ Plugin UI shows browser and manager
5. ✅ Example plugins demonstrate all features
6. ✅ Documentation explains usage and API
7. ✅ All unit tests pass
8. ✅ Zero compilation errors

**Status:** Ready to commence Phase 4B implementation

---

## Next Phases Preview

### Phase 4C: Marketplace System
- Plugin registry backend
- Publishing system
- User ratings and reviews
- Automated plugin updates

### Phase 4D: CLI Tools
- Command-line interface
- Batch save editing
- Plugin execution from CLI
- Script file support

### Phase 4E: Advanced Scripting
- Extended Lua libraries
- Hot-reload support
- Debug mode
- Performance profiling

---

**Prepared by:** Development Team  
**Date:** January 16, 2026  
**Ready for Implementation:** Yes ✅
