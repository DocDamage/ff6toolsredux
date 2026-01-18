# Character Class System Plugin

Transform Final Fantasy VI into a class-based RPG with traditional job systems.

## Overview

This plugin adds a comprehensive class system to FF6, allowing you to assign traditional RPG classes to characters. Each class comes with unique stat modifiers, equipment restrictions, and command ability sets, fundamentally changing how characters develop and play.

## Features

### Class Library
- **12 Unique Classes**: Knight, Mage, Thief, Monk, Dragoon, White Mage, Black Mage, Samurai, Ninja, Sage, Berserker, Mystic Knight
- **Class Descriptions**: Each class has unique characteristics and playstyles
- **Equipment Restrictions**: Classes restrict what equipment characters can use
- **Stat Modifiers**: Classes modify character stats (HP, Attack, Defense, Magic, etc.)
- **Command Abilities**: Classes determine available command abilities

### Class Management
- **Assign Classes**: Assign any class to any character
- **View Assignments**: See current class assignments for all characters
- **Class Details**: View detailed information for each class
- **Configuration Persistence**: Class assignments saved between sessions

### Export & Sharing
- **Export Configurations**: Export class assignments to text files
- **Timestamped Exports**: All exports include date/time stamps
- **Shareable Configs**: Share your class configurations with other players

## Class Roster

| Class | Focus | Stat Bonuses | Equipment |
|-------|-------|--------------|-----------|
| **Knight** | Tank | HP+30%, Def+40%, Vig+20% | Swords, Shields, Heavy Armor |
| **Mage** | Magic | Magic+50%, MP+40%, HP-20% | Rods, Robes |
| **Thief** | Speed | Speed+40%, Stam+20% | Daggers, Light Armor |
| **Monk** | Physical | Vigor+50%, HP+20%, Stam+30% | Claws, Gi |
| **Dragoon** | Jump | Vigor+30%, HP+20% | Spears, Heavy Armor |
| **White Mage** | Healing | Magic+40%, MP+50%, MDef+30% | Rods, Robes |
| **Black Mage** | Offense | Magic+60%, MP+40%, HP-30% | Rods, Hats |
| **Samurai** | SwdTech | Vig+30%, Mag+10%, Stam+20% | Katanas, Light Armor |
| **Ninja** | Throw | Speed+50%, Vigor+20% | Daggers, Shurikens, Light Armor |
| **Sage** | Lore | Magic+40%, Vigor+10%, HP+10% | Rods, Robes |
| **Berserker** | Power | Vigor+60%, HP+40%, Def+20% | Axes, Heavy Armor |
| **Mystic Knight** | Hybrid | Vig+20%, Mag+30%, HP+10% | Swords, Light Armor |

## Usage

### Assigning a Class

1. Run the plugin
2. Select "Assign Class to Character"
3. Choose the character to modify
4. Select the desired class
5. Class is applied and configuration saved

### Viewing Class Library

1. Run the plugin
2. Select "View Class Library"
3. Browse all available classes with full details
4. See stat modifiers, equipment restrictions, and command abilities

### Exporting Configuration

1. Run the plugin
2. Select "Export Configuration"
3. Configuration saved to `class_system/export_YYYYMMDD_HHMMSS.txt`
4. Share with other players or keep as backup

## Technical Details

### Stat Modifiers

Stat modifiers are multiplicative:
- Values > 1.0 increase the stat
- Values < 1.0 decrease the stat
- Example: Knight with base 100 HP gets 130 HP (1.3x multiplier)

### Equipment Restrictions

Classes specify allowed equipment types:
- **Weapons**: Swords, Daggers, Spears, Rods, Katanas, Claws, Axes
- **Armor**: Heavy Armor, Light Armor, Robes, Gi, Hats, Shields

Attempting to equip restricted items should be blocked (requires equipment API).

### Command Abilities

Each class specifies command ability IDs:
- 1 = Fight (all classes have this)
- 2 = Magic
- 4 = Steal
- 6 = SwdTech
- 7 = Throw
- 9 = Blitz
- 10 = White Magic (conceptual)
- 11 = Lore

### Configuration Storage

Configurations stored in `class_system/configurations.lua`:
```lua
{
  assignments = {
    ["Terra"] = "mage",
    ["Locke"] = "thief",
    -- etc.
  },
  multi_class = {} -- Reserved for future multi-class support
}
```

## API Requirements

**Current Status**: Conceptual implementation - some features require APIs not yet available.

### Required APIs

1. **Character Stat Modification**
   - `character.modify_stat(name, stat, multiplier)` - Apply class stat modifiers
   - Currently conceptual - plugin tracks but can't modify stats

2. **Equipment Restriction**
   - `character.set_equipment_restrictions(name, allowed_types)` - Enforce class equipment rules
   - Not yet implemented

3. **Command Ability Swapping**
   - `character.set_command_ability(name, slot, ability_id)` - Assign class abilities
   - Not yet implemented (see Character Ability Swap plugin for workaround)

### Current Functionality

**What Works**:
- Class library browsing
- Class assignment tracking
- Configuration save/load
- Export functionality

**What's Conceptual**:
- Actual stat modification (displays calculated values)
- Equipment restrictions (configuration only)
- Command ability changes (tracked but not applied)

## Warnings

⚠️ **EXPERIMENTAL PLUGIN**

This plugin fundamentally alters character progression and balance:
- Class stat modifiers may break game balance
- Equipment restrictions can lock you out of items
- Some characters may become overpowered or underpowered
- Not all features functional without additional APIs

**Recommendations**:
- Test on a backup save first
- Start with balanced classes (Mystic Knight, Sage)
- Avoid extreme classes (Berserker, Black Mage) until you understand the system
- Export your configuration regularly

## Examples

### Traditional Party Setup
```
Terra:     White Mage (healer)
Locke:     Thief (utility)
Edgar:     Knight (tank)
Sabin:     Monk (damage)
```

### Challenge Run Setup
```
Terra:     Black Mage (glass cannon)
Celes:     Mystic Knight (hybrid)
Cyan:      Samurai (SwdTech specialist)
Shadow:    Ninja (speed demon)
```

### Power Gaming Setup
```
Terra:     Sage (balanced magic)
Celes:     Mystic Knight (best hybrid)
Edgar:     Knight (unkillable tank)
Sabin:     Monk (physical powerhouse)
```

## FAQ

**Q: Can I change a character's class later?**  
A: Yes, simply reassign a new class. The plugin will apply the new modifiers.

**Q: Do class stat modifiers stack with espers?**  
A: Conceptually yes, but implementation depends on how APIs calculate stats.

**Q: Can characters multi-class?**  
A: Not yet, but the system has placeholders for future multi-class support.

**Q: Why aren't my class restrictions working in-game?**  
A: Equipment restriction APIs are not yet implemented. The plugin tracks restrictions but can't enforce them.

**Q: What happens if I unassign a class?**  
A: You need to reassign to another class or manually edit the configuration file.

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## License

Part of the FF6 Save Editor plugin system. See main LICENSE file.
