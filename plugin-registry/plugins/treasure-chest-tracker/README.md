# Treasure Chest Tracker

Track opened vs. unopened treasure chests across both worlds, surface missables before the World of Ruin, and export a copyable checklist. The plugin is **read-only** and refuses to guess: if a chest flag is not exposed by your editor build, it shows `?` (Unknown) instead of faking a status.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: WoB/WoR treasure chests, with filters for world, region, status, and missables

---

## Why You Need This

- WoB has numerous point-of-no-return chests (Floating Continent, Sealed Gate, Factory). Missing them can lock you out of items.
- WoR has sprawling dungeons where unopened chests are easy to forget (Phoenix Cave, Ancient Castle, Kefka's Tower).
- Memory is unreliable; a deterministic checklist avoids backtracking and regret.
- Safety: Unknown API coverage is clearly shown, so you never assume a chest was opened when the editor cannot confirm it.

---

## Features

- **Summary dashboard**: Counts of opened, unopened, and unknown statuses plus detected world.
- **Unopened view**: Filters to chests that are definitely unopened (or Unknown if you opt in via config).
- **Missables view**: Highlights WoB-only and point-of-no-return chests.
- **World filter**: Uses detected world (WoB/WoR/Both) to narrow lists.
- **Region browser**: Quick region picker (Narshe, Figaro, Vector, etc.).
- **Search**: Keyword search across item, area, region, and notes.
- **Export**: Copyable plaintext of all tracked chests for routing/stream notes.
- **Graceful fallback**: Unknown status surfaced when the chest flag API is missing.

---

## Quick Start

1) Load a save.  
2) Run the plugin and choose `1) Summary` to see world detection and counts.  
3) Use `2) Unopened` to list what is still available.  
4) Use `3) Missables` before the Floating Continent.  
5) Use `5) Browse by Region` to sweep areas methodically.  
6) Optionally `6) Search` for a specific item or location.  
7) `7) Export` to copy the list into notes or a stream overlay.

---

## Data and Status Logic

- Each chest entry has: `id`, `world` (`WoB`/`WoR`/`Both`), `region`, `area`, `item`, `value`, `missable`, `note`.
- The plugin checks status in this order:
  1. `GetChestOpened(chestId)` if available
  2. `GetFlag(chestId)` fallback
  3. Otherwise `Unknown`
- Config `assumeUnknownAsUnopened` lets planners treat Unknown as unopened (off by default to avoid false negatives).

### World detection
1) `GetWorldState()`  
2) `GetFlag("world_state")`  
3) Fallback to `World of Balance` (configurable: `defaultWorldFallback`).

---

## Menus

- **1) Summary**: World, counts for opened/unopened/unknown. Reminder that Unknown = API not available.
- **2) Unopened**: Filters to unopened (and optionally unknown). Good for final sweeps.
- **3) Missables**: WoB-only and point-of-no-return chests, flagged with `!` when `highlightMissables` is true.
- **4) World Filter**: Shows chests matching the detected world (WoB/WoR/Both).
- **5) Browse by Region**: Pick a region to list all chests there.
- **6) Search**: Keyword search over region/area/item/note.
- **7) Export**: Copyable text block of all tracked chests (status icons included).
- **8) Help**: Quick reminders and safety notes.

Status icons:
- `[X]` opened
- `[ ]` unopened
- `[?]` unknown (API could not confirm)
- `!` prefix when `highlightMissables` is enabled

---

## Configuration

All options are in `config` at the top of `plugin.lua`:

- `pageSize` (default 20): Pagination size for list views.
- `assumeUnknownAsUnopened` (default false): Treat Unknown as unopened to be conservative when planning.
- `showValueEstimates` (default true): Show gil estimates next to items when present.
- `includeWorldFilterInSummary` (default true): Note world context in the summary.
- `defaultWorldFallback` (default "World of Balance"): Used if world detection fails.
- `highlightMissables` (default true): Adds a `!` prefix to missable entries.

---

## Missables Covered (subset example)

- **Floating Continent**: Genji Glove, Murasame (point-of-no-return).  
- **Sealed Gate**: Chocobo Suit, Barrier Ring (pre-collapse).  
- **Magitek Factory**: Flame/Thunder gear along the escape rails.  
- **Lete River / Opera / Vector Armory**: One-time sequences.  
- These show `!` and appear in the Missables view; status may be `?` if the API is missing.

---

## Regions Covered (sample)

Narshe, Figaro, South Figaro, Returner Hideout, Lete River, Opera House, Vector, Magitek Facility, Sealed Gate, Floating Continent, Thamasa, Phoenix Cave, Ancient Castle, Kefka's Tower, Vector Ruins, Darryl's Tomb, Owzer's Mansion, Esper Cave, Colosseum.

> The database is a representative core set. Expand it as chest IDs become available from the API or data dumps.

---

## Export Format

`7) Export` produces a text block with status icons, region/area/item, value (if enabled), and notes. Copy it into your route notes, a tracker overlay, or a shared doc with co-op partners.

---

## Limitations and Gaps

- Chest flag API may not exist in your build; affected entries show `?` until the API ships.
- Database is a curated subset (high-value and missable spots). Extend `CHESTS` as more IDs become known.
- Gil values are rough estimates for prioritization only.
- No write operations; this plugin cannot toggle chest states.

---

## Troubleshooting

- **Everything is Unknown**: Your build may lack chest APIs; confirm `GetChestOpened` or `GetFlag` exposes IDs. If not, wait for the new API or use manual planning.
- **A chest shows opened but I never looted it**: Some versions repurpose chest flags; verify the ID mapping before relying on it. Adjust the `CHESTS` table if needed.
- **I want Unknown treated as unopened**: Set `assumeUnknownAsUnopened = true` in `config`.
- **Pagination too small/large**: Adjust `pageSize` in `config`.

---

## How to Extend

- Add new entries to `CHESTS` with `id`, `world`, `region`, `area`, `item`, `value`, `missable`, and `note`.
- If your build exposes different API names, update `isChestOpened()` to match.
- Add world variants by setting `world = "Both"` for chests that persist into WoR.

---

## Roadmap (Phase 2 alignment)

- Integrate full chest database (~400+) once the chest flag API is finalized.
- Add per-dungeon progress summaries and value totals.
- Cross-link with World of Balance Checklist for pre-WoR missable coverage.
- Optional manual checkboxes (write-enabled variant) in a future release.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and planned enhancements.
