# Damage Calculator

Approximate physical and magical damage planner with elemental modifiers, row/back attack penalties, criticals, multi-target penalties, and quick physical-vs-magic comparisons. Designed to be **read-only** and to surface Unknown when APIs are missing.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: Physical/Magic damage estimates; enemy stats optional

---

## What It Does

- Physical calculator: level, Vigor, weapon power, enemy defense, row/back, critical.
- Magic calculator: level, Magic, spell power, enemy magic defense, element (weak/resist/absorb/null), multi-target penalty.
- Comparison mode: side-by-side physical vs. magic outputs using shared inputs.
- Elemental multipliers (configurable): weak 1.5x, resist 0.5x, absorb -1x, null 0x.
- Defaults provided when enemy stats are unknown; you can override interactively.

---

## Quick Start

1) Run the plugin → choose Physical or Magic.  
2) Enter level, stats, power, and enemy defenses (or accept defaults).  
3) For Magic, pick the elemental relationship and multi-target flag.  
4) Use Compare to decide whether to cast or swing.

---

## Formulas (Approximate FF6-Style)

- **Physical**: `(level^2 * Vigor) / 256 + weaponPower - (enemyDef / 2)`, then modifiers for row (0.5 back row), back attack (0.5), and critical (2.0).
- **Magic**: `((level + spellPower)^2)/32 + Magic*4 - enemyMagDef`, then elemental multiplier and multi-target penalty (0.75).

These are planner-grade approximations, not frame-perfect emulator formulas.

---

## Configuration (top of plugin.lua)

- `defaultEnemyDefense` (120), `defaultEnemyMagDefense` (150)
- `defaultWeaponPower` (120), `defaultSpellPower` (50), `defaultLevel` (50)
- `elementalWeakMultiplier` (1.5), `elementalResistMultiplier` (0.5)
- `backAttackPenalty` (0.5), `rowPenalty` (0.5), `criticalMultiplier` (2.0)
- `multiTargetPenalty` (0.75), `pageSize` (15)

---

## Safety and Fallbacks

- If enemy data is unavailable, the plugin prompts for values or uses defaults.
- No write operations; entirely read-only.
- Unknown inputs are never guessed silently; prompts allow manual override.

---

## Limitations

- Uses simplified formulas; certain edge-case modifiers (relic/weapon quirks) are not modeled.
- Does not account for defense ignoring, split damage rounding, or status-driven multipliers beyond the basics listed.
- Enemy fetch uses `GetBestiaryEntry(id)` if available; otherwise defaults are used.

---

## Roadmap

- Add weapon/relic flag handling (ignore defense, random variance, vigor caps).
- Add spell-specific quirks (Pearl proc, fixed damage, unblockable effects).
- Add batch scenario save/load for repeated comparisons.
- Add party equipment picker integration when API exposure allows.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and planned improvements.
