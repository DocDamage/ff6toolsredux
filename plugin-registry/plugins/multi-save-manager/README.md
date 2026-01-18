# Multi-Save Manager

Manage multiple playthroughs and save slots with custom names, notes, statistics tracking, and metadata persistence.

> Version: 1.0.0  
> Permissions: `read_save`, `ui_display`, `file_io`  
> Scope: Multi-save metadata and statistics management

---

## Features
- List all save files in managed directory with custom names.
- Add custom names and notes/journal entries per save.
- View save details (playtime, gil, character count, completion estimate).
- Metadata persistence in Lua table format.
- Read-only save analysis (no modifications).

---

## Quick Start
1) Place save files in `saves/` directory.
2) Run the plugin → choose `List Saves` to see all slots.
3) Choose `Add Metadata` to set custom name and notes.
4) Choose `View Details` to see playtime, gil, character count.

---

## Metadata Storage
Metadata stored in `saves_metadata.lua` as Lua table:
```lua
{
  ["mysave.sav"] = {
    name = "Main Playthrough",
    notes = "Just reached World of Ruin",
    playtime = 12345,
    completion = 50,
  },
}
```

---

## Statistics
Automatically extracted when available via `LoadSaveFile`:
- **Playtime**: Total play time in game units.
- **Gil**: Current gil amount.
- **Characters**: Count of enabled characters.
- **Completion**: Rough estimate (characters/14 * 100).

---

## API Requirements
- `ListFiles(dir)`: list saves directory.
- `ReadFile(path)`, `WriteFile(path, content)`: metadata persistence.
- `LoadSaveFile(path)`: optional for statistics extraction.

---

## Limitations
- No quick-switch functionality (requires host integration).
- No archive/unarchive completed playthroughs.
- No save duplication.
- No search/filter by metadata.
- No export playthrough summary.

---

## Roadmap
- Add quick-switch to load selected save in editor.
- Add archive/unarchive for completed playthroughs.
- Add save duplication/cloning.
- Add search and filter by tags/notes.
- Add export playthrough summary to text.
- Add save slot comparison.

---

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for release history.
