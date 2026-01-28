# Save Backup Manager

Automated backup management with timestamped backups, restore functionality, and automatic backup rotation to maintain a safety net.

> Version: 1.0.0  
> Permissions: `read_save`, `write_save`, `ui_display`, `file_io`  
> Scope: Backup creation, restoration, and library management

---

## Features
- Create timestamped backups with optional labels.
- List all backups sorted by date (latest first).
- Restore backups to any target path.
- Delete old backups manually.
- Automatic rotation (keeps last 10 backups by default).
- File-based persistence in `backups/` directory.

---

## Quick Start
1) Run the plugin → choose `Create Backup`.
2) Enter source save file path.
3) Enter optional label (default: "manual").
4) Backup created with timestamp (e.g., `manual_20260116_143052.sav`).
5) To restore: choose `Restore Backup` → select from library → enter target path.

---

## Backup Naming
Format: `<label>_<timestamp>.sav`
- `<label>`: User-provided label (default: "manual").
- `<timestamp>`: YYYYMMDD_HHMMSS format.

Example: `manual_20260116_143052.sav`

---

## Auto-Rotation
After each backup creation, oldest backups are deleted if total exceeds `maxBackups` (10 by default). Configure in `config.maxBackups` at top of plugin.lua.

---

## API Requirements
- `ReadFile(path)`: read save file content.
- `WriteFile(path, content)`: write backup/restore.
- `ListFiles(dir)`: list backup directory.
- `DeleteFile(path)`: remove old backups.

---

## Limitations
- No automatic backup before plugin modifications (requires host integration).
- No backup compression.
- No backup metadata (size, description, playtime snapshot).
- No scheduled/periodic auto-backup.

---

## Roadmap
- Add automatic backup before write operations.
- Add backup metadata (file size, description, save summary).
- Add backup compression.
- Add scheduled auto-backup with configurable intervals.
- Add backup comparison to current save.

---

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for release history.
