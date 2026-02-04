# Full Technical Debt Audit - Final Fantasy VI Save Editor

**Audit Date:** 2026-02-04  
**Auditor:** Codex (GPT-5)  
**Go Version:** 1.22.4 (toolchain); CGO disabled in this environment  
**Total Go Files:** 289  
**Test Files:** 37 (12.8%)  
**Total Packages (dirs with Go files):** 41  
**Total Lines of Go Code:** ~61,000  

---

## Executive Summary

| Category | Status | Notes |
|----------|--------|-------|
| **Build Status** | ⚠️ BLOCKED | `go test ./...` fails due to CGO/OpenGL build constraints (Fyne/GL) |
| **Test Status** | ⚠️ BLOCKED | Coverage and tests cannot run in current Linux env without CGO/OpenGL |
| **go vet** | ⚠️ BLOCKED | Same CGO/OpenGL constraint blocks vet |
| **Code Formatting** | ⚠️ NEEDS ATTENTION | 19 files need `gofmt` |
| **Documentation** | 🟢 GOOD | 36/41 packages have `doc.go` |
| **Test Coverage** | ⚠️ UNKNOWN | Coverage generation blocked; last known analysis in `TEST_COVERAGE_ANALYSIS.md` (2026-01-31) |

---

## Changes Since Last Audit (2026-01-31)

- **Go files:** 231 → 289 (+58)
- **Test files:** 19 → 37 (+18)
- **Go LOC:** ~45k → ~61k (+~16k)
- **Package docs:** 11 → 36 packages with `doc.go`
- **Formatting debt:** ~45 files → 19 files
- **Tooling status:** prior audit showed passing tests; current Linux environment cannot compile Fyne/GL without CGO/OpenGL toolchain.

---

## 🔴 Critical Issues (Immediate Action Required)

### 1. Tooling Blocked by CGO/OpenGL Dependencies

**Impact:** `go test ./...`, `go vet ./...`, and `govulncheck ./...` fail in this Linux environment because CGO is disabled (no GCC/OpenGL dev libs), and Fyne’s GLFW/GL packages require those toolchains. This prevents CI-quality verification outside of Windows.

**Observed Failure:**
- `github.com/go-gl/gl/v2.1/gl`: build constraints exclude all Go files
- Fyne GLFW/GL symbols missing

**Fix Options:**
- Add a Linux CI runner with OpenGL + GCC dev packages.
- Add a headless build tag (or build constraints) for UI packages so non-UI tests can run without GL.
- Keep Windows CI as primary and document that Linux headless testing is unsupported without extra deps.

---

## 🟡 High Priority Issues

### 2. High Cyclomatic Complexity (gocyclo > 15)

Top hotspots:
- `io/pr/saver.go:98` — `(*PR).saveCharacters` (59)
- `ui/forms/editors/mapData.go:76` — `(*MapData).CreateRenderer` (55)
- `plugins/api_character.go:14` — `(*APIImpl).GetCharacter` (52)
- `io/pr/loader_characters.go:15` — `(*PR).loadCharacters` (36)
- `plugins/api_character.go:190` — `(*APIImpl).SetCharacter` (33)
- `ui/forms/editors/mapData.go:487` — `(*MapData).showTeleportDialog` (30)
- `cli/commands_validate.go:155` — `validateCharacters` (24)
- `io/json/importer.go:39` — `(*Importer).ImportFromJSON` (23)
- `io/pr/loader.go:15` — `(*PR).Load` (22)
- `io/pr/loader_map.go:78` — `(*PR).loadTransportation` (22)

**Recommendation:** Extract helpers, reduce nesting, and split responsibilities in these functions.

### 3. Ignored Errors in Production Code (56 instances)

**Impact:** Silent failures in settings persistence, file I/O, cloud sync, or UI can lead to data loss or undiagnosed issues.

**Examples:**
- `settings/manager.go` — repeated `_ = m.Save()` ignores persistence failures
- `io/config/config.go` — `_ = json.Unmarshal(...)` / `_ = os.WriteFile(...)`
- `cloud/manager.go:136` — `_ = m.SyncAll(ctx)`
- `io/file/fileIO.go:78-79` — `_ = zw.Flush()` / `_ = zw.Close()`
- `ui/window.go:84` — `_ = os.Setenv("FYNE_FONT", ...)`

**Recommendation:** Handle errors or document intentional ignores with `//nolint:errcheck` and rationale.

### 4. Unchecked Type Assertions (Potential Panics)

**Examples:**
- `io/pr/loader_characters.go:226` — `k.(json.Number).Int64()` without ok-check
- `settings/manager.go:426-427` — `mp["X"].(float64)` / `mp["Y"].(float64)` with ignored ok

**Recommendation:** Use safe assertions with `ok` checks, or convert via helpers that return errors.

### 5. Skipped Tests in Critical Plugin Security Path

**Files:**
- `plugins/version_constraint_test.go.skip`
- `plugins/security_test.go.skip`
- `plugins/sandbox_test.go.skip`
- `plugins/reload_state_test.go.skip`
- `plugins/analytics_engine_test.go.skip`
- `plugins/plugin_profiler_test.go.skip`
- `plugins/hot_reload_test.go.skip`

**Risk:** Security and lifecycle behavior can regress without detection.

**Recommendation:** Convert to build tags or re-enable with stable fixtures.

---

## 🟢 Medium Priority Issues

### 6. Code Formatting (gofmt)

19 files require formatting, including:
- `models/pr/inventory.go`
- `io/pr/saver_test.go`
- `cli/commands_*_test.go`
- `cloud/manager_test.go`

**Fix:** `gofmt -w .`

### 7. Outstanding TODOs (19 in Go code)

- `scripting/runner.go` — memory cap + network blocking (4)
- `ui/forms/template_manager.go` — template UI (4)
- `plugins/loader.go` — parse Lua metadata
- `plugins/api_character.go` — set stat in PR save
- `plugins/api_equipment.go` — batch operations
- `models/speedrun/configs.go` — refactor + apply config
- `models/pr/inventory.go` — error handling
- `io/palette_editor.go` — undo/redo history
- `ui/forms/party_preset_manager.go` — save preset
- `ui/forms/script_editor.go` — cancellation
- `ui/forms/plugin_manager_sync.go` — marketplace version
- `ui/forms/plugin_browser.go` — username from settings

### 8. `init()` Usage (20 functions)

**Packages:**
- `global/vars.go`, `io/config/config.go`, `models/pr/*`, `models/consts/*`, `models/speedrun/configs.go`, `ui/forms/editors/command.go`

**Recommendation:** Replace with explicit `Initialize()` wiring for clarity and testability.

### 9. Heavy `interface{}` Usage

**Count:** 253 occurrences. Most concentrated in:
- `io/pr/saver.go` (24)
- `io/pr/type_safe_extractors.go` (19)
- `settings/manager.go` (11)
- `plugins/api_equipment.go` (9)

**Recommendation:** Continue replacing with typed helpers or generics where possible.

### 10. Documentation Gaps (5 packages missing `doc.go`)

- Root package (`.`)
- `docs/`
- `models/consts/pr/`
- `tools/`
- `ui/state/`

### 11. Large Files (>500 LOC)

- `ui/forms/editors/text_maps.go` (1306) — data file (acceptable)
- `io/pr/saver.go` (888) — complex save logic
- `marketplace/client.go` (871) — mixed concerns
- `ui/forms/plugin_browser.go` (745) — UI complexity
- `plugins/manager.go` (701) — lifecycle + auditing + execution

---

## 🔵 Low Priority Issues

### 12. `recover()` Usage

- `main.go` — global panic handler
- `ui/forms/plugin_browser.go` — plugin isolation

**Status:** Acceptable if paired with logging and user feedback.

### 13. `exec.Command` Usage

- `browser/update.go` — opening URLs only (expected)
- `io/pr/encounters.go` — commented out Python hook

**Status:** Low risk.

### 14. No `unsafe` / `syscall`

No direct usage detected.

---

## 📊 Code Metrics Summary

| Metric | Value | Target | Status |
|--------|-------|--------|--------|
| Total Go Files | 289 | - | - |
| Test Files | 37 (12.8%) | 25%+ | 🔴 Low |
| Packages with Tests | 15/41 (36.6%) | 100% | 🔴 Low |
| Total TODOs (Go) | 19 | 0 | 🟡 Medium |
| init() Functions | 20 | 0 | 🟡 Medium |
| gofmt Issues | 19 files | 0 | 🟡 Medium |
| `interface{}` Usage | 253 | Reduce | 🟡 Medium |
| go test/vet/vulncheck | Blocked | Pass | 🔴 Blocked |

---

## 🛡️ Security Audit

| Check | Status | Notes |
|-------|--------|-------|
| `unsafe` package | ✅ None | Safe |
| `syscall` usage | ✅ None | Safe |
| `exec.Command` | ✅ Limited | Only URL open |
| Plugin sandbox | 🟡 Needs hardening | TODOs for memory cap + network blocking |
| Error handling | 🟡 Mixed | Many ignored errors |
| CGO/GL | 🔴 Tooling blocked | Tests/vet/vulncheck can’t run here |

---

## 🔧 Dependencies Audit

**Direct Dependencies:**
- `fyne.io/fyne/v2` v2.7.2 — no update detected
- `github.com/fsnotify/fsnotify` v1.9.0 — no update detected
- `github.com/kiamev/ffpr-save-cypher` v1.0.0 — no update detected
- `github.com/sqweek/dialog` v0.0.0-20240226 — **update available** (v0.0.0-20260123)
- `github.com/yuin/gopher-lua` v1.1.0 — **update available** (v1.1.1)
- `gitlab.com/c0b/go-ordered-json` v0.0.0-20201030 — no update detected (stale)

**Note:** `go list -m -u -json all` failed under Go 1.22.4 due to module retractions in dependencies requiring Go 1.24+. Direct module checks succeeded individually.

---

## 📁 Docs + Plugin Ecosystem Review

**Plugins registry:** All entries have `README.md` and `CHANGELOG.md`.  
**Lua TODOs:** none found in `.lua` files.  
**Skipped tests:** 7 `.skip` tests in `plugins/` (see High Priority).  

---

## 📦 Binary / Asset Hygiene

Large binaries tracked (>5MB):
- `ffvi-save-editor.exe`
- `FFVIPR_Save_Editor.exe`
- `ffvi_editor.exe`
- `ffvi_editor_test.exe`
- `test_build.exe`
- `go1.25.6.src.tar.gz`
- `go1.25.6.windows-amd64.msi`

**Risk:** Repository bloat and Git hosting limits. Recommend .gitignore rules + pre-commit size guard.

---

## 🎯 Recommended Action Plan

### Phase 1 (Week 1) — Unblock Tooling + Safety
1. Add headless build tag or UI-optional build path to allow `go test ./...` on Linux CI.
2. Re-enable skipped plugin security tests (convert `.skip` → build tags if needed).
3. Fix critical ignored errors (`settings`, `io/config`, `cloud`).

### Phase 2 (Week 2–3) — Complexity + Quality
4. Refactor top 5 high-complexity functions.
5. Replace unsafe type assertions in `io/pr/loader_characters.go` and `settings/manager.go`.
6. Run `gofmt -w .` across the repo.

### Phase 3 (Week 4) — Coverage + Docs
7. Increase tests in `io/pr`, `plugins`, `scripting`, `settings`.
8. Add `doc.go` to the 5 missing packages.
9. Remove or document large binaries and add size guards.

---

## Summary

The codebase is feature-rich and mature, with improved documentation coverage and expanded test presence since the last audit. The main blocker is **tooling inability to run full tests/vet/vulncheck in a headless Linux environment** due to CGO/OpenGL constraints from Fyne. Addressing this unlocks reliable CI and confidence. Next priorities are complexity hotspots, ignored errors, and re-enabling skipped security tests.

---

*Audit completed: 2026-02-04*  
*Next audit recommended: 2026-03-04*
