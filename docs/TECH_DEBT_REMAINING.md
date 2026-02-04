# Remaining Technical Debt Items (as of January 2026)

## High-Impact Items
- **Large binaries tracked risk**
  - Edge: history bloat; GitHub rejects >100MB.  
  - Option: keep `.gitignore`, add `pre-push` hook to block binaries; BFG cleanup if already pushed.  
  - Best: hook + periodic audit; run BFG if upstream contains them.

## Quick Wins (1–2 days)
- Embed resources/*.png and fallback to solid color
- Wrap errors with context and show user-friendly messages

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
