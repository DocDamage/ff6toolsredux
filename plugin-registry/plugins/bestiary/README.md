# Bestiary

Complete enemy reference with encounter tracking, weaknesses, drops, steals, sketch/control data, and multi-criteria search. The plugin is **read-only** and intentionally honest: if the encounter flag API is missing, it marks entries as `?` (Unknown) rather than guessing.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: Enemies, bosses, dragons, special encounters; encounter tracking if API is available

---

## Why Use This

- Know which enemies you have seen for completion goals (Rages, Sketch, Control, Coliseum planning).
- Quickly find weaknesses/resistances without leaving the editor.
- Plan drops/steals/morphs for farming and speedruns.
- Identify bosses and special encounters you may have missed.
- Stay safe: the plugin surfaces Unknown when it cannot verify a flag.

---

## Features

- **Summary**: Counts of seen, unseen, and unknown enemies.
- **All/Unseen views**: Browse everything or focus on what you have not encountered (Unknown optionally treated as unseen via config).
- **Filters**: By type (Humanoid, Beast, Dragon, Boss, etc.), element match (weak/resist/absorb/null), bosses only, keyword search (name/location/type/status).
- **Drops/Steals/Morph/Sketch/Control/Rage**: Displayed when present.
- **Locations**: Shown when known for quick hunting.
- **Export**: Copyable plaintext for routing notes.

---

## Quick Start

1) Load a save.  
2) Run the plugin → `1) Summary` to see counts.  
3) Use `3) Unseen` to focus on missing encounters.  
4) Use `4) By Type` (e.g., Dragon) or `5) By Element` (e.g., Ice weak) to target what you need.  
5) `7) Search` for a name/location/status keyword.  
6) `8) Export` to copy the list into notes or a tracker overlay.

---

## Data Model and Status Logic

- Each enemy entry includes: `id`, `name`, `level`, `hp`, `mp`, `type`, `elements` (weak/resist/absorb/null), `statusImmune`, `drop`, `steal`, `morph`, `sketch`, `control`, `rageMove`, `location`, `boss`.
- Encounter state is determined in this order:
  1. `GetEncountered(enemyId)` if provided by the API
  2. `GetFlag("bestiary_<id>")` fallback
  3. Otherwise `Unknown`
- Config `assumeUnknownAsNotSeen` lets you treat Unknown as unseen for conservative planning (off by default).

Status icons:
- `[SEEN]` encountered
- `[UNSEEN]` not encountered
- `[?]` unknown (API missing)

---

## Menu Options

- **1) Summary**: Counts of seen/unseen/unknown. Quick health check.
- **2) All**: Full list with status icons.
- **3) Unseen**: Filters to unseen (and optionally unknown) entries.
- **4) By Type**: Exact match on `type` (case-insensitive): Humanoid, Beast, Dragon, Boss, etc.
- **5) By Element**: Matches any enemy that is weak/resist/absorb/null to the given element.
- **6) Bosses**: Boss-only view.
- **7) Search**: Keyword across name, type, location, or status immunity text.
- **8) Export**: Copyable plaintext of all entries.
- **9) Help**: Quick reference.

---

## Configuration

Edit `config` at the top of `plugin.lua`:

- `pageSize` (default 15): Pagination size for list dialogs.
- `assumeUnknownAsNotSeen` (default false): Treat Unknown as unseen to be conservative.
- `showDropsAndSteals` (default true): Toggle drop/steal lines.
- `showMorphAndSketch` (default true): Toggle morph/sketch/control lines.
- `showRageMoves` (default true): Show Rage move name.
- `showLocations` (default true): Show location text when available.

---

## Element and Status Coverage

- Elements are grouped into `weak`, `resist`, `absorb`, `null`. Searches match any bucket.
- Status immunities are listed per enemy; use the search function to find enemies immune to a status by keyword.

---

## Notes on Data Completeness

- The database included here is a **representative core** covering key enemies/bosses/specials (e.g., Intangir, Tyranosaur, Dragons). It is structured for a full 384-entry import when available.
- Extend the `ENEMIES` table with additional rows; the UI and filters work automatically.
- If your build exposes different API names for encounter flags, update `isEncountered()` accordingly.

---

## Troubleshooting

- **Everything is Unknown**: Your build may not expose `GetEncountered` or the `bestiary_*` flags. Wait for API updates or treat Unknown as unseen via config.
- **Mismatched names/types**: Update the relevant entries in `ENEMIES` to align with your localization/version.
- **Pagination too small/large**: Adjust `pageSize`.

---

## Roadmap

- Import full enemy dataset (~384) once finalized.
- Add per-element and per-type statistics to Summary.
- Cross-link with Rage Tracker and Treasure Chest Tracker for route planning.
- Optional write-enabled variant for manual toggles (future, opt-in).

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history and planned enhancements.
