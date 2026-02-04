# Test Coverage Analysis - Final Fantasy VI Save Editor

**Analysis Date:** 2026-01-31  
**Current Coverage:** 21 test files covering 11 of 39 packages (28%)  
**Untested Packages:** 28 of 39 packages (72%)

---

## 📊 Current Coverage Summary

| Package | Coverage | Status | Priority |
|---------|----------|--------|----------|
| **settings** | 74.1% | 🟢 Good | - |
| **io/file** | 50.0% | 🟢 Good | - |
| **ui/state** | 35.2% | 🟡 Medium | - |
| **plugins** | 34.2% | 🟡 Medium | - |
| **marketplace** | 32.6% | 🟡 Medium | - |
| **models/pr** | 28.2% | 🟡 Medium | - |
| **io/pr** | 24.0% | 🟡 Medium | - |
| **cloud** | 14.1% | 🟡 Medium | - |
| **scripting** | 14.5% | 🟡 Medium | - |
| **cli** | 4.2% | 🔴 Low | High |

---

## 🔴 Critical Priority (Untested Core Packages)

### 1. `io/pr/` - Save File I/O (24% coverage, needs expansion)

**Current Tests:**
- `decode_helpers_test.go`
- `decode_helpers_fuzz_test.go`
- `factory_test.go`
- `loader_test.go`
- `saver_test.go`
- `type_safe_extractors_test.go`
- `type_safe_extractors_fuzz_test.go`

**Missing Coverage:**
| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `saver.go` | 813 | 🔴 Critical | Comprehensive save tests |
| `loader_characters.go` | 279 | 🔴 Critical | Character loading tests |
| `loader_map.go` | 163 | 🟡 Medium | Map data loading |
| `loader_misc.go` | 107 | 🟡 Medium | Esper/misc loading |
| `loader_inventory.go` | 84 | 🟡 Medium | Inventory loading |
| `loader_helpers.go` | 190 | 🟡 Medium | Helper functions |
| `encounters.go` | ~100 | 🟢 Low | Veldt encounters |
| `compare.go` | ~80 | 🟢 Low | Save comparison |

**Recommended Tests:**
```go
// TestSaveRoundTrip - Full save/load cycle
// TestCharacterSaveLoad - Character data integrity
// TestInventorySaveLoad - Inventory preservation
// TestMapDataSaveLoad - Map state preservation
// TestCorruptedSaveHandling - Error handling
// TestSaveFileVersioning - Backward compatibility
```

---

### 2. `settings/manager.go` - Configuration (74.1% coverage)

**Status:** Well covered, but could expand edge cases

**Additional Tests Needed:**
- Concurrent access tests
- Invalid config recovery
- Settings migration

---

### 3. `models/pr/` - Data Models (28.2% coverage)

**Current Tests:**
- `characters_test.go`

**Missing Coverage:**
| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `characters.go` | ~200 | 🔴 Critical | Expand existing |
| `inventory.go` | ~150 | 🔴 Critical | Item management |
| `baseOffsets.go` | ~260 | 🟡 Medium | Offset calculations |
| `hpMpCounts.go` | ~100 | 🟡 Medium | HP/MP tables |

**Recommended Tests:**
```go
// TestInventoryValidation - Item limits, duplicates
// TestCharacterStatCalculation - Stat formulas
// TestBaseOffsetIntegrity - Offset table validation
// TestHPMPLookupTables - Table correctness
```

---

## 🟡 High Priority (Important Business Logic)

### 4. `scripting/` - Lua Scripting (14.5% coverage)

**Current Tests:**
- `runner_test.go`

**Missing Coverage:**
| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `bindings.go` | 406 | 🔴 Critical | API bindings |
| `vm.go` | ~300 | 🔴 Critical | VM operations |
| `stdlib.go` | ~200 | 🟡 Medium | Standard library |

**Recommended Tests:**
```go
// TestScriptExecutionTimeout - Timeout handling
// TestScriptSandboxEscape - Security tests
// TestCharacterBindings - Character API
// TestInventoryBindings - Inventory API
// TestMaliciousScriptHandling - Security
```

---

### 5. `plugins/` - Plugin System (34.2% coverage)

**Current Tests:**
- `audit_logger_test.go`
- `dependency_resolver_test.go`
- `manager_integration_test.go`
- `plugins_test.go`
- `sandbox_test.go`
- `security_test.go`

**Missing Coverage:**
| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `api.go` | ~750 | 🔴 Critical | Main API surface |
| `api_character.go` | ~350 | 🔴 Critical | Character API |
| `api_equipment.go` | ~300 | 🔴 Critical | Equipment API |
| `api_inventory.go` | ~200 | 🟡 Medium | Inventory API |
| `api_party.go` | ~150 | 🟡 Medium | Party API |
| `manager.go` | 595 | 🟡 Medium | Plugin lifecycle |

---

### 6. `cloud/` - Cloud Sync (14.1% coverage)

**Current Tests:**
- `cloud_test.go`

**Missing Coverage:**
| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `manager.go` | 400 | 🔴 Critical | Sync orchestration |
| `gdrive.go` | ~200 | 🟡 Medium | Google Drive |
| `dropbox.go` | ~200 | 🟡 Medium | Dropbox |

**Note:** Cloud tests may need mocking for CI

---

### 7. `cli/` - Command Line (4.2% coverage)

**Current Tests:**
- `combat_pack_test.go`

**Missing Coverage:**
| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `commands_stub.go` | ~1000 | 🔴 Critical | All commands |
| `commands_validate.go` | 415 | 🔴 Critical | Validation |
| `commands_core.go` | ~300 | 🟡 Medium | Core CLI |

**Recommended Tests:**
```go
// TestEditCommand - Character editing
// TestExportCommand - JSON export
// TestImportCommand - JSON import
// TestBatchCommand - Batch operations
// TestScriptCommand - Lua script execution
// TestBackupCommand - Backup creation
// TestValidateCommand - Save validation
```

---

## 🟢 Medium Priority (Supporting Packages)

### 8. `io/` Root Package (0% coverage)

| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `palette_editor.go` | 426 | 🟡 Medium | Palette operations |
| `animation_exporter.go` | ~200 | 🟢 Low | Animation export |
| `sprite_converter.go` | ~200 | 🟢 Low | Sprite conversion |
| `frame_sequencer.go` | ~150 | 🟢 Low | Frame sequencing |

---

### 9. `io/config/` (0% coverage)

| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `config.go` | ~150 | 🟡 Medium | Config management |

---

### 10. `models/` Root (0% coverage)

| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `change.go` | ~100 | 🟡 Medium | Change tracking |
| `validation.go` | ~50 | 🟡 Medium | Validation |

---

### 11. `models/batch/` (0% coverage)

| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `operations.go` | ~100 | 🟡 Medium | Batch operations |

---

## 🔵 Low Priority (Data/UI Packages)

### 12. `models/consts/` (0% coverage)

**Status:** Acceptable - mostly constant definitions

### 13. `models/game/` (0% coverage)

| File | Lines | Risk | Test Needed |
|------|-------|------|-------------|
| `rage_database.go` | ~250 | 🟢 Low | Rage data |
| `magic_catalog.go` | ~420 | 🟢 Low | Magic data |
| `sketch_database.go` | ~100 | 🟢 Low | Sketch data |

### 14. `ui/` and Subpackages (0% coverage)

**Status:** Acceptable for now - UI testing is complex

| Package | Files | Approach |
|---------|-------|----------|
| `ui/` | 8 files | Integration tests |
| `ui/forms/` | 35 files | E2E tests |
| `ui/editors/` | 27 files | Manual testing |

---

### 15. `achievements/`, `browser/`, `global/` (0% coverage)

**Status:** Low priority - small packages

---

## 📋 Recommended Test Implementation Plan

### Phase 1: Critical (Week 1-2)

```bash
# 1. Expand io/pr tests (target: 60% coverage)
io/pr/saver_test.go           # Add comprehensive save tests
io/pr/loader_characters_test.go # Character loading
io/pr/loader_integration_test.go # Full round-trip

# 2. Expand models/pr tests (target: 50% coverage)
models/pr/inventory_test.go    # Item operations
models/pr/characters_test.go   # Expand existing
```

### Phase 2: High Priority (Week 3-4)

```bash
# 3. Scripting tests (target: 40% coverage)
scripting/bindings_test.go     # API bindings
scripting/vm_test.go           # VM operations
scripting/security_test.go     # Sandbox (expand)

# 4. CLI tests (target: 50% coverage)
cli/commands_edit_test.go
cli/commands_export_test.go
cli/commands_import_test.go
cli/commands_batch_test.go
cli/commands_validate_test.go
```

### Phase 3: Medium Priority (Week 5-6)

```bash
# 5. Plugin API tests (target: 50% coverage)
plugins/api_character_test.go
plugins/api_equipment_test.go
plugins/api_inventory_test.go

# 6. Cloud tests (with mocks)
cloud/manager_test.go          # Expand existing
```

### Phase 4: Supporting (Week 7-8)

```bash
# 7. I/O utilities
io/palette_editor_test.go
io/config/config_test.go

# 8. Model utilities
models/change_test.go
models/batch/operations_test.go
```

---

## 🎯 Coverage Targets

| Quarter | Target Coverage | Packages to Cover |
|---------|-----------------|-------------------|
| Q1 2026 | 40% overall | io/pr, models/pr, cli |
| Q2 2026 | 60% overall | scripting, plugins |
| Q3 2026 | 75% overall | cloud, io/* |
| Q4 2026 | 85% overall | All core packages |

---

## 💡 Testing Strategy Recommendations

### 1. Unit Tests
- Focus on pure functions in `models/` and `io/pr/`
- Mock file system for I/O tests
- Table-driven tests for validation logic

### 2. Integration Tests
- Save/load round-trip tests
- Plugin loading/unloading
- Cloud sync with mocked providers

### 3. Fuzz Tests (Already Started)
- `io/pr/decode_helpers_fuzz_test.go`
- `io/pr/type_safe_extractors_fuzz_test.go`
- Expand to cover more parsing functions

### 4. E2E Tests
- CLI command execution
- Full workflow: Load → Edit → Save → Validate

---

## 🚨 High-Risk Untested Code

These areas have high complexity and no tests - **test immediately**:

1. **`io/pr/saver.go:813`** - Complex save logic, data corruption risk
2. **`cli/commands_stub.go:~1000`** - CLI entry points, user-facing
3. **`plugins/api.go:~750`** - Plugin API, security boundary
4. **`scripting/bindings.go:406`** - Lua bindings, security critical
5. **`cloud/manager.go:400`** - Data sync, corruption risk

---

## 📈 Expected Impact

Implementing this plan would:
- Increase overall coverage from ~15% to ~60%
- Cover all critical save file operations
- Prevent regressions in core functionality
- Enable safer refactoring
- Improve CI/CD confidence

---

*Analysis completed: 2026-01-31*
