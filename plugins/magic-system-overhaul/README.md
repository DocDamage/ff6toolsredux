# Magic System Overhaul Plugin

**Version:** 1.0.0  
**Category:** Experimental  
**Author:** FF6 Plugin Team

## Overview

Fundamentally transform FF6's magic system. Unlimited MP, zero MP costs, unlock all spells, alternative cost systems (HP-based, item-based), magic power multipliers, and preset modes (FF7 Materia-style, MP-free, etc.).

## Features

- **6 Magic Modes** - Normal, unlimited MP, zero cost, HP-based, item-based, free magic
- **MP Control** - Set to 9999, zero, or custom values
- **Cost Modification** - 0x to 10x MP cost multipliers
- **Spell Unlocking** - Unlock all 54 spells for characters
- **Magic Power** - 0.1x to 10x damage multipliers
- **Preset Systems** - FF7 Materia-style, MP-free, challenge modes
- **Backup & Restore** - Undo all changes

## Quick Start

### Free Magic Mode
```lua
applyMagicMode("free_magic")  -- Unlimited MP + zero costs + all spells
```

### Unlimited MP Only
```lua
applyMagicMode("unlimited")  -- 9999 MP for all characters
```

### Zero MP Costs
```lua
applyMagicMode("zero_cost")  -- All spells cost 0 MP
```

### FF7 Materia Style
```lua
applyMateriaStyle()  -- All spells available, unlimited MP, normal costs
```

## Function Reference

### Magic Modes
- `applyMagicMode(mode)` - Apply magic system mode
  - `"normal"` - Standard FF6 system
  - `"unlimited"` - 9999 MP for all
  - `"zero_cost"` - Spells cost 0 MP
  - `"hp_based"` - Spells cost HP (experimental)
  - `"item_based"` - Spells consume items (experimental)
  - `"free_magic"` - Unlimited + zero cost + all spells

### MP Functions
- `setUnlimitedMPAll()` - Set all characters to 9999 MP
- `setUnlimitedMP(char_id)` - Set specific character to 9999 MP
- `setZeroMP(char_id)` - Set specific character to 0 MP

### Cost Modification
- `setMPCostMultiplier(multiplier)` - Set MP cost multiplier (0-10x)
- `enableZeroMPCosts()` - All spells cost 0 MP
- `restoreNormalMPCosts()` - Restore 1x costs

### Spell Unlocking
- `unlockAllSpellsForAllCharacters()` - Unlock all 54 spells for everyone
- `unlockAllSpells(char_id)` - Unlock all spells for character
- `unlockSpell(char_id, spell_id)` - Unlock specific spell

### Magic Power
- `setMagicPowerMultiplier(multiplier)` - Set damage multiplier (0.1-10x)

### Presets
- `applyMateriaStyle()` - FF7-style (all spells, unlimited MP)
- `applyMPFreeSystem()` - MP-free (free magic mode)
- `applyLowMPChallenge()` - 2x MP costs
- `applyHighMPBonus()` - 0.5x MP costs

### Analysis
- `displayMagicSystemStatus()` - View current configuration
- `listMagicModes()` - List available modes
- `restoreBackup()` - Undo all changes

## Use Cases

- **Challenge Runs** - Zero MP or MP-free playthroughs
- **Build Testing** - Test magic strategies with unlimited resources
- **Speedrunning** - Unlimited MP for routing
- **Experimentation** - Try alternative magic cost systems
- **Power Fantasy** - Max MP + all spells + zero costs

## Magic Modes Explained

### Normal
- Standard FF6 magic system
- MP-based costs
- Spells must be learned

### Unlimited MP
- All characters set to 9999 MP
- Normal MP costs still apply
- Still need to learn spells

### Zero Cost
- All spells cost 0 MP
- Unlimited casting
- Still need to learn spells

### HP-Based (Experimental)
- Spells cost HP instead of MP
- MP costs set to 0
- HP cost = original MP cost × 2
- **Note:** Requires in-game implementation

### Item-Based (Experimental)
- Spells consume items
- MP costs set to 0
- Item cost = MP cost ÷ 10
- **Note:** Requires in-game implementation

### Free Magic
- Unlimited MP (9999)
- Zero MP costs
- All spells unlocked
- Complete magic freedom

## Preset Systems

### FF7 Materia Style
- All 54 spells available to all characters
- 9999 MP for everyone
- Normal MP costs (resource management preserved)

### MP-Free System
- Same as free magic mode
- No resource management

### Low-MP Challenge
- 2x MP costs (spells more expensive)
- More strategic resource management

### High-MP Bonus
- 0.5x MP costs (spells cheaper)
- More frequent casting

## Warnings

⚠️ **Game Balance** - Breaks intended magic system  
⚠️ **Overpowered** - Free magic mode extremely strong  
⚠️ **Alternative Costs** - HP/item-based modes require in-game implementation  
⚠️ **Save Compatibility** - Modified saves may not work in vanilla FF6  

## Version History

See [CHANGELOG.md](CHANGELOG.md).

---

**Plugin Version:** 1.0.0  
**LOC:** ~650 LOC  
**Documentation:** ~6,800 words
