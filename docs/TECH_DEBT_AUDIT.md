# Technical Debt Audit (January 2026)

## High-Impact Items (with options, edge cases, and recommended fixes)
- **Build/toolchain fragility** ✅
  - **Status:** Completed. Windows GitHub Action added: installs MSYS2, GCC, Go; runs `go test ./...` and builds exe. Mirrors user environment and catches CGO/GL issues early.
- **Module hygiene** ✅
  - **Status:** Completed. Go module hygiene workflow added: enforces `go mod tidy && go mod verify` in CI.
- **Settings persistence split** ✅
  - **Status:** Completed. Unified settings struct implemented; migration logic added; all code refactored to use new settings manager; legacy config removed.
- **Loader robustness (`io/pr/loader`)**
  - Edge cases: malformed save data (OrderedMap vs string), nil keys, type asserts panic.  
  - Status: Complete. All loader and decode helpers now use safe type assertions, with table-driven and fuzz tests for all edge cases.  
  - Fix: Added typed decode helpers, table-driven tests, and Go fuzz tests for all loader and decode helper functions. All interface{} assertions are now guarded.
- **Lua sandbox/package path**
  - Edge cases: CWD not repo; unbounded memory/time; plugin escapes via `io/os`.  
  - Options:  
    1) Set `package.path` to embedded/bundled plugins only and block `require` elsewhere.  
    2) Add instruction limit (`lua.State.SetContext` with hook), and cap heap via allocator wrapper.  
    3) Disallow `os`, `io`, `debug` (already partially).  
  - Best: combine (1)+(2)+(3); add tests that change CWD and assert failures.
- **Tests writing to disk**
  - Edge cases: parallel tests colliding on real paths.  
  - Option: enforce `t.TempDir()` wrappers; linters for `os.WriteFile` without temp.  
  - Best: wrap helpers (e.g., `testfs.WriteTempFile(t, data)`).
- **Plugin smoke tests dependency**
  - Edge: running from other dirs misses `./plugins`.  
  - Options: fail fast with clear message if plugin files absent; embed minimal fixtures.  
  - Best: embed tiny fixture copies for smoke tests; still allow full plugin dir when present.
- **Assets/resources**
  - Edge: missing PNG/font crashes; working directory not repo root.  
  - Options: embed via `//go:embed`, or check existence and fallback solid color.  
  - Best: embed key assets; add runtime fallback to avoid crash.
- **Lint/format gap**
  - Options: `golangci-lint` with selective enables; `gofmt -w` check; `staticcheck` later.  
  - Best: start with `golangci-lint` baseline + `gofmt` in CI.
- **Undo/state coverage**
  - Edge: regressions invisible until user data loss.  
  - Option: small unit tests for push/pop/limit/merge; property tests for idempotence.  
  - Best: add table tests now; property tests later.
- **Config error handling**
  - Edge: silent failure to save, user thinks settings persisted.  
  - Option: bubble errors to UI dialog + log; retry with backup file.  
  - Best: bubble + log + keep last-good backup.
- **Large binaries tracked risk**
  - Edge: history bloat; GitHub rejects >100MB.  
  - Option: keep `.gitignore`, add `pre-push` hook to block binaries; BFG cleanup if already pushed.  
  - Best: hook + periodic audit; run BFG if upstream contains them.

## Quick Wins (1–2 days)
- Add GitHub Actions workflow (Windows runner): **Done**
- Hardened `io/pr` loaders with typed decode helpers, table-driven tests, and fuzz tests for inventory/transportation/map data.
- Wrap resource loads (fonts/backgrounds) with safe fallbacks to prevent startup crashes.

## Suggested Follow-ups
- Implement the Windows CI workflow. **Done**
- Add loader/unit tests and stricter Lua sandbox limits.
- Unify settings storage and document the migration path. **Done**

## Additional Findings (second sweep)
- **Assets packaging**
  - Edge: relative paths break when running from another directory or after `go install`.  
  - Solutions:  
    1) `//go:embed resources/*.png` and load via `fyne.NewStaticResource`.  
    2) At startup, check file presence and fallback to solid-color background.  
    Best: (1) embed; (2) as safety net.
- **Error surfaces**
  - Edge: dialog shows raw `os` errors, confusing users.  
  - Solutions: wrap errors with context (`fmt.Errorf("saving %s: %w", file, err)`) and show user-friendly text while logging full stack.  
  - Best: structured log + friendly toast/dialog.
- **Logging**
  - Edge: no audit trail; debugging user reports is hard.  
  - Solutions: minimal logger writing to `%AppData%/ff6editor/log.txt`; level from settings.  
  - Best: lightweight file+stderr logger; rotate by size.
- **Skipped tests**
  - Edge: `*_test.go.skip` in plugins rot silently.  
  - Solutions: CI job listing skipped tests; convert to build tags; document rationale.  
  - Best: add lint that fails if new `.skip` tests appear; convert important ones to normal tests.
- **Security / plugin sandbox**
  - Edge: Lua plugins might read/write host FS or network.  
  - Solutions: limit `package.path`, remove `os/io/debug`, add VFS layer, deny network; document allowed APIs.  
  - Best: VFS-backed sandbox + documented contract; enforce in tests.
- **Docs lint**
  - Edge: many Markdown files; link rot likely.  
  - Solutions: `markdown-link-check` + `codespell` in CI (allowlist repo to avoid noise).  
  - Best: run on PRs; nightly optional to keep noise low.
- **Cross-platform build helper**
  - Edge: new contributors stumble on env vars.  
  - Solutions: `build.sh` + `build.ps1` that detect tools, print actionable hints, and exit non-zero.  
  - Best: add both scripts; reuse in CI.
- **Binary cache pollution**
  - Edge: stray settings/saves/ROMs committed.  
  - Solutions: extend `.gitignore`; add pre-commit hook to block >5MB binaries (except resources); CI job to fail on large git objects.  
  - Best: hook + CI size guard; keep `.gitignore` updated.
