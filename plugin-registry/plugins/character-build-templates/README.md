# Character Build Templates

Save and share character build presets including equipment, espers, and learned spells. Manage a build library with validation and quick-apply functionality.

> Version: 1.0.0  
> Permissions: `read_save`, `write_save`, `ui_display`, `file_io`  
> Scope: Character build persistence and library management

---

## Features
- Save current character builds (equipment IDs, esper, spells with proficiency).
- Load builds and apply to any character slot.
- Build library management (list, delete).
- File-based persistence using `file_io` permission.
- Safe pcall wrappers with fallback mutation when APIs unavailable.

---

## Quick Start
1) Run the plugin → choose `Save Build`.
2) Enter character slot (0-15) to serialize.
3) Enter filename (saved as `builds/<filename>.lua`).
4) To apply: choose `Load Build` → select from library → apply to slot.

---

## Build Format
Builds are serialized as Lua table literals:
- `version`: format version (1.0).
- `charName`, `charID`, `level`: metadata.
- `equipment`: weapon/shield/armor/helmet/relic1/relic2 IDs.
- `esper`: equipped esper ID.
- `spells`: array of `{id, name, value}` (proficiency).

Files stored in `builds/` directory with `.lua` extension.

---

## API Requirements
- `GetCharacter(slot)`: read character data.
- `SaveCharacter(slot)`: persist changes (write_save).
- `SetSpell(slot, id, value)`: optional for spell writes.
- `ReadFile(path)`, `WriteFile(path, content)`, `ListFiles(dir)`, `DeleteFile(path)`: file_io operations.

Fallback: direct mutation of character tables if APIs missing.

---

## Limitations
- No validation of item/esper inventory availability (assumes IDs are valid).
- No build categories/tags/metadata beyond filename.
- Spells applied with best-effort (SetSpell API or direct mutation).
- No import/export to external formats (JSON/clipboard).

---

## Roadmap
- Add build validation against current inventory.
- Add build categories and metadata (author, description, tags).
- Add import/export to JSON for sharing.
- Add build comparison tool.
- Add popular/curated builds database.

---

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for release history.
