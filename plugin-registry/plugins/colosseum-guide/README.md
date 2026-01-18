# Colosseum Guide

Interactive betting reference that maps bet item → opponent → reward, with enemy hints, strategy notes, filters, and an exportable checklist. The plugin is **read-only** and surfaces Unknown when inventory/bestiary APIs are missing.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: Colosseum bets/rewards, opponent hints, routing aid

---

## Why Use It

- Quickly see what you get for betting an item, who you fight, and how to prepare.
- Filter to items you own (or unknown) so you do not scroll a huge list.
- Search by bet/reward/opponent/strategy to answer questions fast mid-run.
- Export to a text block for routing, streaming overlays, or co-op notes.
- Safety-first: Unknown inventory stays visible—no guessing.

---

## Features

- **All Bets view**: Full mapping bet → opponent → reward.
- **Owned view**: Filters to bets currently in your inventory (or Unknown if the inventory API is unavailable; configurable).
- **Tier filter**: weapon / armor / shield / relic / misc.
- **Search**: Keyword across bet, reward, opponent name, and strategy notes.
- **Enemy hints**: Lightweight weakness/notes shown when known (no full bestiary import required).
- **Export**: Copyable plaintext of the current mapping for quick sharing.
- **Help**: Quick reminder of controls and safety behavior.

---

## Quick Start

1) Load a save.  
2) Run the plugin → `1) Summary` to confirm counts.  
3) `3) Bets You Own` to see actionable options.  
4) `4) By Tier` or `5) Search` to narrow the list.  
5) `6) Export` if you want a copyable text list.

---

## Data Model

Each entry includes:
- `bet`: Item you wager.
- `reward`: Item you receive if you win.
- `opponent`: Enemy you face.
- `enemyId` (optional): Used for hint lookup.
- `strategy`: One-liner tactics.
- `rewardNotes`: Why the reward matters.
- `tier`: weapon / armor / shield / relic / misc.

The included database is a **representative core** (high-value and common bets). Extend `BETS` in `plugin.lua` as needed.

---

## Status and Safety

- Inventory check uses `GetInventory()`; if missing, bets are treated as Unknown (or missing if `treatUnknownInventoryAsMissing` is true).
- The plugin does **not** modify inventory or flags.
- Enemy hints are lightweight overrides; if you want full stats, open the Bestiary plugin in parallel.

---

## Menu Options

- **1) Summary**: Counts and reminders.  
- **2) All Bets**: Full list.  
- **3) Bets You Own**: Filters to owned (or Unknown) bets.  
- **4) By Tier**: weapon/armor/shield/relic/misc.  
- **5) Search**: Keyword across bet/reward/opponent/strategy.  
- **6) Export**: Copyable text block.  
- **7) Help**: Quick reference.  
- **0) Exit**.

---

## Configuration

Adjust in `config` at the top of `plugin.lua`:
- `pageSize` (default 15): Pagination size for list dialogs.
- `showEnemyStats` (default true): Show enemy hint lines when available.
- `showStrategyNotes` (default true): Show the strategy note line.
- `treatUnknownInventoryAsMissing` (default false): If true, Unknown inventory is treated as missing for the Owned filter; if false, Unknown stays visible.

---

## Extending the Database

Add rows to `BETS` in `plugin.lua`. Keep fields consistent: `bet`, `reward`, `opponent`, `enemyId` (optional), `strategy`, `rewardNotes`, `tier`. Add enemy hints to `ENEMY_HINTS` keyed by `enemyId` for weakness/notes.

---

## Troubleshooting

- **No items show in Owned view**: Inventory API might be missing; either set `treatUnknownInventoryAsMissing = false` (default) to still show Unknown, or wait for API support.
- **Opponent hint missing**: Add an entry to `ENEMY_HINTS` with the `enemyId` used in the bet row.
- **Too many entries per dialog**: Lower `pageSize`.

---

## Roadmap

- Expand to full bet table.
- Add per-reward rarity tags and sorting.
- Optional write-enabled variant for manual tracking.
- Cross-links to Bestiary for full enemy stats.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history and planned enhancements.
