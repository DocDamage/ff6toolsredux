# Save File Comparator

Compare two save files side-by-side with detailed diff view showing differences in characters, inventory, and game state.

> Version: 1.0.0  
> Permissions: `read_save`, `write_save`, `ui_display`, `file_io`  
> Scope: Read-only save comparison

---

## Features
- Load and compare two save files from disk.
- Category-based comparison (Characters, Inventory, Game State).
- Diff summary with total difference count.
- Character-level comparison (name, level, HP).
- Inventory item count comparison.
- Game state comparison (play time, gil).

---

## Quick Start
1) Run the plugin → choose `Compare Two Saves`.
2) Enter path to save file 1.
3) Enter path to save file 2.
4) Review comparison results dialog.

---

## Comparison Categories
- **Characters**: Name, level, HP differences per slot.
- **Inventory**: Item count comparison.
- **Game State**: Play time, gil differences.
- **Summary**: Aggregated difference counts.

---

## API Requirements
- `LoadSaveFile(path)`: preferred loader for save data.
- `ReadFile(path)`: fallback for raw file access.

Comparison depth depends on save file structure returned by APIs.

---

## Limitations
- No merge functionality (read-only comparison).
- Limited to fields accessible via LoadSaveFile API.
- No detailed item-by-item inventory diff.
- No export of comparison report.

---

## Roadmap
- Add merge capabilities (apply changes from one save to another).
- Add detailed item-by-item inventory comparison.
- Add spell/equipment comparison per character.
- Add export comparison report to text file.
- Add visual diff highlighting.

---

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for release history.
