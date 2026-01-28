# Esper Stat Growth Optimizer

Read-only planner for level-up stat gains via esper bonuses. Targets a stat (Vigor, Magic, Speed, Stamina, Defense, Magic Defense), ranks available espers, and estimates total gain from current to target level.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: Single-stat planning per character

---

## Features
- Ranks espers by per-level bonus for a chosen stat.
- Plans expected total gain from current level to a target level.
- Limits pool to owned espers when inventory APIs are available; otherwise uses the full reference list.
- Quick comparison view (top espers for a stat).
- Safe pcall-wrapped, read-only dialogs.

---

## Quick Start
1) Run the plugin → enter character slot (0-15).  
2) Choose `Plan stat growth`.  
3) Enter target stat (e.g., Magic) and target level (defaults to 99).  
4) Review the recommended esper and backup options.

Use `Compare espers for a stat` to see the ranking table without picking a character.

---

## Data & Assumptions
- Prefers runtime esper data from plugin APIs when available (names, Owned/Acquired flags, stat growth fields). Falls back to the shipped table from `models/game/esper_growth.go` (Vigor, Speed, Stamina, Magic, Defense, MagicDef) if growth data is missing.
- If owned esper inventory is detected (`GetEsperInventory`/`GetEspers`), ranking is limited to owned; otherwise the full list is used.
- Estimates only; intended for quick planning, not exact stat caps.

---

## Configuration
Adjust at the top of `plugin.lua`:
- `defaultTargetLevel` / `defaultCurrentLevel`: planning defaults.
- `ESPERS`: bundled fallback table; only used if runtime growth data is unavailable.

---

## Limitations
- Does not model natural stat caps, equipment, or temporary buffs.
- Uses simplified esper bonuses (flat growth only). HP/MP% growth is not modeled.
- No import/export of plans (reference-only).

---

## Changelog
See [CHANGELOG.md](CHANGELOG.md) for release history.
