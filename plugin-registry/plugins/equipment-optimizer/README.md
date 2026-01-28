# Equipment Optimizer

An intelligent equipment optimization plugin for the FF6 Save Editor that analyzes your available equipment and recommends the best loadout for each character based on configurable optimization goals.

## Features

- **5 Optimization Goals:** Attack, Defense, Magic, Speed, or Balanced
- **Equipment Analysis:** Evaluates all weapons, shields, helms, and armor
- **Stat Comparison:** Shows before/after stat changes
- **Party-Wide Optimization:** Optimize entire party with one command
- **Equipment Comparison:** Compare top options for each equipment slot
- **Preview Mode:** See changes before applying
- **One-Click Apply:** Automatically equip recommended gear

## Installation

1. Open FF6 Save Editor (version 3.4.0 or higher)
2. Go to `Tools → Marketplace` or `Community → Marketplace`
3. Search for "Equipment Optimizer"
4. Click "Install"

## Usage

### Basic Usage

1. Load a save file in the FF6 Save Editor
2. Go to `Tools → Plugin Manager`
3. Find "Equipment Optimizer" in the installed plugins list
4. Click "Run"
5. Select an operation from the menu

### Available Operations

#### 1. Analyze Single Character

Analyzes one character's equipment and recommends optimal loadout.

**Steps:**
1. Select option 1 from menu
2. Enter character index (0-15)
3. Review current equipment and recommendations
4. Check stat improvements
5. Confirm to apply changes

**Example Output:**
```
=== EQUIPMENT ANALYSIS ===

Terra (Lv45)
Current Equipment:
  Weapon: Mythril Pike
  Shield: Mythril Shield
  Helm: Leather Hat
  Armor: Leather Armor

Optimization Goal: BALANCED

Recommended Equipment:
  Weapon: Illumina (Atk+110)
  Shield: Paladin Shield (Def+60, MDef+40)
  Helm: Genji Helm (Def+40, MDef+20)
  Armor: Force Armor (Def+65, MDef+35)

Stat Changes:
  Attack: 28 → 110 (+82)
  Defense: 53 → 165 (+112)
  Magic Def: 15 → 95 (+80)
```

#### 2. Optimize Party Equipment

Optimizes equipment for all current party members at once.

**Steps:**
1. Select option 2 from menu
2. Review party members to be optimized
3. Confirm to proceed
4. Plugin analyzes and applies optimal equipment for each party member

**Use Case:** Quickly gear up your party before major boss battles or difficult dungeons.

#### 3. Change Optimization Goal

Switch between different optimization strategies.

**Available Goals:**
- **Attack:** Maximize offensive power (best for physical attackers)
- **Defense:** Maximize survivability (best for tanks)
- **Magic:** Maximize magic power (best for mages)
- **Speed:** Maximize turn frequency (best for support/fast attackers)
- **Balanced:** All-around optimization (best for versatile characters)

**Steps:**
1. Select option 3 from menu
2. Choose desired goal (1-5)
3. New goal will be used for all subsequent analyses

#### 4. Compare Equipment Options

View top 3 equipment options for each slot based on current optimization goal.

**Steps:**
1. Select option 4 from menu
2. Enter character index
3. Review top options for weapons, shields, helms, and armor

**Example Output:**
```
=== EQUIPMENT COMPARISON ===

Celes (Lv40)

Weapons:
  1. Illumina (Score: 565)
  2. Excalibur (Score: 495)
  3. Ragnarok (Score: 530)

Shields:
  1. Paladin Shield (Score: 500)
  2. Force Shield (Score: 445)
  3. Crystal Shield (Score: 350)
...
```

## Configuration

The plugin supports several configuration options (modify at top of plugin.lua):

```lua
local config = {
    optimizationGoal = "balanced",  -- attack, defense, balanced, speed, magic
    showAllOptions = false,         -- Show all equipment options or just top recommendation
    considerRoleMatch = true,       -- Factor in character's primary role
    preserveRelics = true,          -- Don't change relic assignments
    previewChanges = true,          -- Show preview before applying
}
```

### Configuration Options

- **optimizationGoal:** Default optimization strategy
- **showAllOptions:** If true, shows all viable equipment options
- **considerRoleMatch:** Factors in character's natural role (not yet implemented)
- **preserveRelics:** If true, doesn't modify relic slots
- **previewChanges:** If true, shows confirmation dialog before applying

## Optimization Formulas

### Attack Goal
```
Score = (Attack × 5) + (Defense × 1) + (Speed × 2)
```
Prioritizes offensive power with some consideration for speed.

### Defense Goal
```
Score = (Defense × 5) + (Magic Defense × 4) + (Attack × 1)
```
Maximizes survivability with balanced physical and magic defense.

### Magic Goal
```
Score = (Magic × 5) + (Magic Defense × 3) + (Defense × 1)
```
Optimizes for spellcasting power and magic survivability.

### Speed Goal
```
Score = (Speed × 5) + (Attack × 2) + (Defense × 1)
```
Maximizes turn frequency while maintaining some offense.

### Balanced Goal
```
Score = (Attack × 2) + (Defense × 2) + (Magic × 2) + (Magic Defense × 2) + (Speed × 1)
```
Provides well-rounded stats across all categories.

## Equipment Database

The plugin includes a simplified equipment database covering:
- **13 Weapons:** From Dirk to Atma Weapon
- **8 Shields:** From Buckler to Paladin Shield
- **8 Helms:** From Leather Hat to Red Cap
- **9 Armor:** From Leather Armor to Minerva Bustier

**Note:** The database is simplified for demonstration. A full implementation would include all ~250+ equipment items from the game.

## Permissions

This plugin requires the following permissions:

- **read_save:** Required to read character stats and inventory
- **write_save:** Required to apply equipment recommendations
- **ui_display:** Required to show analysis and menus

**Note:** This plugin MODIFIES save data (equipment). Always backup your save before use!

## Use Cases

### For Casual Players
- Quickly find best equipment without manual checking
- Optimize party before difficult battles
- Ensure you're using your best gear

### For Power Users
- Min-max character builds
- Theory-craft different optimization strategies
- Compare equipment effectiveness

### For Speedrunners
- Quickly optimize party for specific boss encounters
- Switch between attack/defense builds efficiently
- Verify optimal gear usage

### For Completionists
- Ensure best equipment is properly equipped
- Track equipment effectiveness across characters
- Optimize inventory usage

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered
- **Platform:** Windows, macOS, Linux (wherever the editor runs)

## Troubleshooting

### "Missing required API methods"
**Cause:** Editor version too old  
**Solution:** Update to FF6 Save Editor 3.4.0 or higher

### Recommendations seem suboptimal
**Cause:** Simplified equipment database may not include all items  
**Solution:** Ensure you have the latest plugin version

### Changes not visible
**Cause:** Save file not reloaded  
**Solution:** Reload save file to see equipment changes

### "Unable to read inventory"
**Cause:** Inventory data structure incompatible  
**Solution:** Ensure save file is properly loaded and compatible

## Known Limitations

- Simplified equipment database (~40 items vs 250+ in game)
- Doesn't consider character-specific equipment restrictions
- Relic optimization not implemented (preserveRelics=true by default)
- Doesn't factor in equipment special effects (elemental resistances, status immunities)
- Inventory availability checking is simplified

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Future Enhancements

Planned features for future versions:
- Complete equipment database with all 250+ items
- Character-specific equipment compatibility checking
- Relic optimization recommendations
- Equipment special effects consideration (elemental, status)
- Set bonus detection (e.g., Genji set)
- Boss-specific equipment recommendations
- Equipment enchantment tracking
- Multi-character equipment comparison
- Equipment swap optimization (minimize inventory changes)
- Export/import equipment loadouts

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or feature requests:
- **GitHub Issues:** [Create an issue](https://github.com/ff6-save-editor/plugin-registry/issues)
- **Discord:** Join the FF6 Editor community Discord
- **Email:** plugins@ff6editor.dev

## Credits

**Author:** FF6 Editor Team  
**Repository:** https://github.com/ff6-save-editor/plugin-registry  
**Editor:** https://github.com/ff6-save-editor/editor
