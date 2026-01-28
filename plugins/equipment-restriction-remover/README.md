# Equipment Restriction Remover Plugin

**Version:** 1.0.0  
**Category:** Experimental  
**Author:** FF6 Plugin Team

## Overview

Break ALL equipment restrictions in FF6. Any character can equip anything, gender restrictions removed, dual-wielding for all, relic stacking enabled, equipment conflicts removed.

## Features

- **Character Restrictions** - Any character equips any weapon/armor/relic
- **Gender Restrictions** - Male characters wear female-only items (Minerva Bustier, etc.)
- **Universal Dual-Wield** - All characters can dual-wield weapons
- **Relic Stacking** - Equip same relic twice (double Economizer!)
- **Conflict Removal** - Equip conflicting relics together (Genji Glove + Offering)
- **Unrestricted Mode** - One command removes ALL restrictions
- **Build System** - Apply preset unrestricted builds
- **Stat Calculation** - Calculate stats with "illegal" equipment

## Quick Start

### Remove All Restrictions
```lua
enableUnrestrictedMode()  -- Remove ALL restrictions at once
```

### Selective Removal
```lua
removeGenderRestrictions()   -- Female-only items for everyone
enableRelicStacking()         -- Stack relics
```

### Force Equip Item
```lua
forceEquipItem(char_id, item_id, slot)  -- Bypass validation
```

## Function Reference

### Restriction Removal
- `enableUnrestrictedMode()` - Remove ALL restrictions
- `removeCharacterRestrictions()` - Any character, any equipment
- `removeGenderRestrictions()` - Remove male/female restrictions
- `enableUniversalDualWield()` - All characters dual-wield
- `enableRelicStacking()` - Same relic in both slots
- `removeEquipmentConflicts()` - Remove relic conflicts
- `restoreNormalRestrictions()` - Restore default restrictions

### Equipment Functions
- `canEquip(char_id, item_id, slot)` - Check if equippable
- `forceEquipItem(char_id, item_id, slot)` - Force equip (bypass validation)
- `applyUnrestrictedBuild(char_id, build_config)` - Apply preset build

### Analysis
- `displayRestrictionStatus()` - View current restrictions
- `calculateUnrestrictedStats(char_id)` - Stats with unrestricted gear
- `listExampleBuilds()` - Show example unrestricted builds

### Backup
- `restoreBackup()` - Undo all equipment changes

## Example Builds

### Double Economizer
```lua
-- Two Economizers = 1 MP all spells
build = {relic1 = 114, relic2 = 114}
applyUnrestrictedBuild(char_id, build)
```

### Universal Genji Glove + Offering
```lua
-- 8x attacks per turn (normally conflicts)
enableUnrestrictedMode()
forceEquipItem(char_id, genji_glove_id, CONFIG.SLOTS.RELIC_1)
forceEquipItem(char_id, offering_id, CONFIG.SLOTS.RELIC_2)
```

### Male Character Minerva Build
```lua
-- Male character wearing female-only armor
removeGenderRestrictions()
forceEquipItem(male_char_id, minerva_bustier_id, CONFIG.SLOTS.ARMOR)
```

## Use Cases

- **Extreme Builds** - Double Economizer, double Offering, etc.
- **Build Testing** - Test any equipment combination
- **Gender Equality** - All characters access all equipment
- **Experimentation** - Try "illegal" combinations

## Warnings

⚠️ **Game Balance** - Breaks intended equipment system  
⚠️ **Overpowered** - Some combinations extremely strong  
⚠️ **Visual Bugs** - Character sprites may not match equipment  
⚠️ **Save Compatibility** - Modified saves may not work in vanilla FF6  

## Restrictions Explained

### Character Restrictions (REMOVED)
- **Before:** Only specific characters can use certain equipment
- **After:** ALL characters can use ALL equipment

### Gender Restrictions (REMOVED)
- **Before:** Female-only items (Minerva Bustier, Red Jacket)
- **After:** Male characters can wear female-only items

### Dual-Wield (ENABLED FOR ALL)
- **Before:** Only with Genji Glove/Offering relic
- **After:** All characters can dual-wield naturally

### Relic Stacking (ENABLED)
- **Before:** Can't equip same relic twice
- **After:** Stack relics (double Ribbon, double Economizer)

### Equipment Conflicts (REMOVED)
- **Before:** Certain relics conflict (Genji Glove + Offering)
- **After:** All relics compatible

## Version History

See [CHANGELOG.md](CHANGELOG.md).

---

**Plugin Version:** 1.0.0  
**LOC:** ~570 LOC  
**Documentation:** ~6,000 words
