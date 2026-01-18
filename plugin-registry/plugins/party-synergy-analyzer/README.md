# Party Synergy Analyzer

Heuristic party analysis for role coverage, elemental coverage, redundancy, and quick suggestions. Read-only and conservative: if party/spell APIs are missing, it reports an error instead of guessing.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: 4-slot active party composition analysis

---

## What It Checks

- **Roles**: Healer, Support, Magic DPS, Physical DPS (heuristic via spells and Vigor).
- **Elements**: Fire, Ice, Lightning, Wind, Water, Earth, Holy, Poison coverage via learned spells.
- **Redundancy**: Flags multiple members filling the same role.
- **Missing roles**: Highlights gaps for survivability and damage.

---

## Quick Start

1) Load a save with an active party.  
2) Run the plugin → `1) Summary` for a fast overview.  
3) `2) Roles` to see inferred roles per member.  
4) `3) Elements` to check coverage.  
5) `4) Suggestions` for quick guidance on gaps/redundancy.

---

## Heuristics

- **Roles**: Detected from key spells (e.g., Cura/Curaga for Healer, Haste/Shell for Support, -aga/Ultima/Flare for Magic DPS) and Vigor threshold for Physical DPS.
- **Elements**: Parsed from spell names (fire/ice/bolt/aero/water/quake/holy/bio). Matches are case-insensitive substring checks.
- These heuristics are intentionally lightweight and explainable; adjust the `config` tables near the top of `plugin.lua` to match your build/localization.

---

## Safety and Fallbacks

- If `GetParty()` or spell lists are missing, the plugin surfaces an error dialog and skips guessing.
- No write operations are performed.

---

## Configuration

At the top of `plugin.lua`:
- `expectedElements`: element list checked for coverage.
- `roleThresholds`: spell lists for Healer/Support/Magic DPS detection; Physical DPS uses Vigor >= 40 by default.
- `maxPartySlots`: informational (not enforced).

---

## Limitations

- Heuristic only; does not inspect equipment/commands beyond spells and Vigor.
- Localization differences may require adjusting spell name keys.
- Does not simulate damage; pair with the Damage Calculator for numbers.

---

## Roadmap

- Add equipment/command awareness (e.g., Tools, Blitz, Throw) for finer role inference.
- Add elemental coverage from weapons/relics and Rage/Sketch where APIs allow.
- Add exportable analysis report.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history and planned enhancements.
