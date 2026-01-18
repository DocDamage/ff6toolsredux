# No Level System Plugin

**Version:** 1.0.0  
**Category:** Experimental  
**Author:** FF6 Plugin Team

## Overview

Remove FF6's leveling mechanics entirely. Play with fixed levels, flat stat progression, and equipment-only power scaling. Perfect for challenge runs that eliminate traditional grinding.

## Features

- **Fixed Level Mode** - Lock all characters at specific level (1-99)
- **Flat Stat Presets** - 4 preset stat configurations (low/balanced/mid/high)
- **Experience Control** - Disable exp gain or set custom multipliers (0-10x)
- **No-Level Mode** - Combined fixed level + flat stats + no exp
- **Challenge Configurations** - Quick setup for low-level and equipment-only runs
- **Esper Bonuses** - Apply esper stat bonuses at level 1
- **Backup & Restore** - Undo all changes

## Quick Start

### Level 1 Challenge
```lua
enableNoLevelMode(1, "low_level")  -- All characters level 1, base stats
```

### Equipment-Only Run
```lua
configureEquipmentOnlyChallenge()  -- Level 1, minimal stats, no exp
```

### Fixed Level 50
```lua
setAllFixedLevel(50)
applyFlatStatsAll("mid_level")
```

## Function Reference

### Core Functions
- `setAllFixedLevel(level)` - Set all characters to fixed level
- `applyFlatStatsAll(preset)` - Apply flat stat preset to all characters
- `disableExperienceGain()` - Disable exp gain entirely
- `setExperienceRate(multiplier)` - Set exp multiplier (0-10x)
- `enableNoLevelMode(level, preset)` - Enable full no-level mode
- `disableNoLevelMode()` - Restore normal leveling

### Presets
- `"low_level"` - Level 1 (HP: 100, Stats: 20)
- `"balanced"` - Level 25 (HP: 1000, Stats: 30)
- `"mid_level"` - Level 50 (HP: 2000, Stats: 40)
- `"high_level"` - Level 99 (HP: 9999, Stats: 99)

### Challenge Modes
- `configureLowLevelChallenge()` - Level 1 challenge setup
- `configureEquipmentOnlyChallenge()` - Equipment-only challenge

### Analysis
- `displayNoLevelStatus()` - View current configuration
- `restoreBackup()` - Undo all changes

## Use Cases

- **Low-Level Runs** - Beat game at level 1
- **Equipment-Only** - Power from gear alone
- **Challenge Runs** - Remove grinding advantage
- **Speedrunning** - Fixed stats for routing
- **Build Testing** - Test equipment combinations

## Warnings

⚠️ **Game Balance** - No levels drastically changes difficulty  
⚠️ **Boss Scaling** - Some bosses expect specific level ranges  
⚠️ **Esper Bonuses** - Level-up bonuses won't apply naturally  

## Version History

See [CHANGELOG.md](CHANGELOG.md).

---

**Plugin Version:** 1.0.0  
**LOC:** ~520 LOC  
**Documentation:** ~5,500 words
