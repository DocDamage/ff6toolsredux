# Magic Learning Planner

A comprehensive magic learning tracking and planning tool for the FF6 Save Editor that helps you efficiently learn all spells across all characters by tracking progress and recommending optimal esper assignments.

## Features

- **Spell Progress Tracking:** View learning progress for each character by magic school
- **Esper Recommendations:** AI-driven recommendations for optimal esper assignments
- **AP Calculations:** Calculate exact AP required to master spells
- **Spell Source Lookup:** Find which espers teach specific spells
- **Progress by School:** Group spells by Black/White/Blue/Red/Special magic
- **Learning Rate Display:** Shows learn rates (x1, x2, x5, x10, etc.)
- **Priority System:** Automatically prioritizes partially learned spells
- **Completion Tracking:** Track overall spell learning completion percentage

## Installation

1. Open FF6 Save Editor (version 3.4.0 or higher)
2. Go to `Tools → Marketplace` or `Community → Marketplace`
3. Search for "Magic Learning Planner"
4. Click "Install"

## Usage

### Basic Usage

1. Load a save file in the FF6 Save Editor
2. Go to `Tools → Plugin Manager`
3. Find "Magic Learning Planner" in the installed plugins list
4. Click "Run"
5. Select an operation from the menu

### Available Operations

#### 1. Show Spell Progress (Character)

Displays comprehensive spell learning progress for a specific character.

**Steps:**
1. Select option 1 from menu
2. Enter character index (0-15)
3. View progress by magic school and detailed spell list

**Example Output:**
```
=== SPELL LEARNING PROGRESS ===

Terra (Lv45)

Black Magic: 14/14 (100.0%)
White Magic: 12/12 (100.0%)
Blue Magic: 8/10 (80.0%)
Red Magic: 6/6 (100.0%)
Special: 3/5 (60.0%)

Detailed Progress:

Black Magic:
  [✓] Fire
  [✓] Ice
  [✓] Bolt
  [✓] Fira
  ...
  [✓] Ultima

White Magic:
  [✓] Cure
  [✓] Cura
  [45%] Curaga
  [✗] Holy
  ...
```

**Symbols:**
- **✓** = Spell fully learned (100%)
- **%** = Spell partially learned (shows percentage)
- **✗** = Spell not started (0%)

#### 2. Recommend Espers for Learning

Provides intelligent esper recommendations based on unlearned spells.

**Steps:**
1. Select option 2 from menu
2. Enter character index (0-15)
3. Review top priority spells and teaching espers
4. See AP requirements for each option

**Example Output:**
```
=== ESPER RECOMMENDATIONS ===

Celes - Unlearned Spells: 15

Top 10 Priority Spells:

1. Curaga (45% learned)
   Taught by: Unicorn (20 AP), Alexandr (50 AP)

2. Holy (0% learned)
   Taught by: Alexandr (50 AP), Bahamut (100 AP)

3. Firaga
   Taught by: Ifrit (100 AP), Phoenix (50 AP)

4. Blizzaga
   Taught by: Shiva (100 AP), Valigarmanda (50 AP)
...
```

**Priority System:**
- Spells with highest learning progress are listed first
- Shows all espers that teach each spell
- Displays AP needed from current progress to 100%

#### 3. Show All Esper Details

Displays complete database of espers and the spells they teach.

**Steps:**
1. Select option 3 from menu
2. View all espers with their spell teaching capabilities

**Example Output:**
```
=== ESPER SPELL TEACHING ===

Ramuh:
  Bolt (Rate: x10, AP: 10)
  Thundara (Rate: x5, AP: 20)
  Thundaga (Rate: x2, AP: 50)

Ifrit:
  Fire (Rate: x10, AP: 10)
  Fira (Rate: x5, AP: 20)
  Firaga (Rate: x1, AP: 100)

Shiva:
  Ice (Rate: x10, AP: 10)
  Blizzara (Rate: x5, AP: 20)
  Blizzaga (Rate: x1, AP: 100)

Bahamut:
  Flare (Rate: x2, AP: 50)
  Meteor (Rate: x1, AP: 100)
...
```

**Use Case:** Quick reference for planning long-term esper assignments.

#### 4. Find Spell Teaching Source

Search for a specific spell to see which espers teach it.

**Steps:**
1. Select option 4 from menu
2. Enter spell name (partial matches work)
3. View espers that teach the spell and learning rates

**Example Output:**
```
=== SPELL LEARNING SOURCE ===

Spell: Ultima
School: Black Magic
Power: 100

Taught by:
  Ragnarok (x1 rate, 100 AP to master)
  Crusader (x10 rate, 10 AP to master)
```

**Use Case:** Quickly look up where to get a specific spell.

## Configuration

The plugin supports several configuration options (modify at top of plugin.lua):

```lua
local config = {
    showLearnedOnly = false,     -- Show only learned spells
    sortByProgress = true,        -- Sort by learning progress
    highlightMissing = true,      -- Highlight unlearned spells
    calculateAPRequired = true,   -- Calculate AP needed for mastery
    groupBySchool = true,         -- Group spells by magic school
}
```

### Configuration Options

- **showLearnedOnly:** If true, hides unstarted spells from progress view
- **sortByProgress:** Prioritizes partially learned spells in recommendations
- **highlightMissing:** Emphasizes spells not yet learned (not yet implemented)
- **calculateAPRequired:** Shows AP calculations in recommendations
- **groupBySchool:** Organizes spells by magic school type

## Magic Schools

The plugin organizes spells into five schools:

### Black Magic (14 spells)
Offensive elemental and damage spells
- Fire, Ice, Bolt (tier 1)
- Fira, Blizzara, Thundara (tier 2)
- Firaga, Blizzaga, Thundaga (tier 3)
- Drain, Bio, Poison, Flare, Ultima

### White Magic (12 spells)
Healing and protective spells
- Cure, Cura, Curaga
- Raise, Arise
- Regen, Esuna
- Shell, Protect, Reflect
- Dispel, Holy

### Blue Magic (10 spells)
Status and utility spells
- Scan, Slow, Rasp, Osmose
- Haste, Stop, Berserk
- Slow 2, Haste 2, Vanish

### Red Magic (6 spells)
Status ailment spells
- Mute, Sleep, Confuse
- Break, Death, Banish

### Special (5 spells)
Unique and powerful spells
- Quick, Warp
- Merton, Meteor, Quake

## AP Calculation Formula

The plugin calculates AP required to master spells using this formula:

```
Total AP Needed = 100 / Learn Rate

AP Remaining = Total AP Needed - (Current Progress × Total AP Needed / 100)
```

**Examples:**
- **Learn Rate x1:** 100 AP to master from 0%
- **Learn Rate x2:** 50 AP to master from 0%
- **Learn Rate x5:** 20 AP to master from 0%
- **Learn Rate x10:** 10 AP to master from 0%

**With Progress:**
- **50% learned, x1 rate:** 50 AP remaining
- **75% learned, x2 rate:** 12.5 AP remaining (rounds to 13)
- **25% learned, x5 rate:** 15 AP remaining

## Esper Database

The plugin includes a simplified esper database:

- **Ramuh:** Bolt, Thundara, Thundaga
- **Ifrit:** Fire, Fira, Firaga
- **Shiva:** Ice, Blizzara, Blizzaga
- **Unicorn:** Cure, Cura, Regen
- **Carbuncle:** Reflect, Shell, Protect
- **Bahamut:** Flare, Meteor
- **Alexandr:** Shell, Protect, Holy

**Note:** The database is simplified for demonstration. A full implementation would include all 27 espers with complete spell lists.

## Permissions

This plugin requires the following permissions:

- **read_save:** Required to read character spell data and esper assignments
- **ui_display:** Required to show progress tracking and recommendations

**Note:** This is a READ-ONLY plugin. It does not modify save data.

## Use Cases

### For New Players
- Understand which espers to prioritize
- Track spell learning progress
- Plan which characters should equip which espers

### For Completionists
- Track progress toward learning all spells
- Identify missing spells quickly
- Optimize esper assignments for complete mastery

### For Speedrunners
- Quickly identify fastest spell learning paths
- Calculate exact AP requirements for route planning
- Verify which spells are available at key points

### For Challenge Runs
- Track spell availability under restricted conditions
- Plan esper-less or limited esper strategies
- Verify natural spell learning progress

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered
- **Platform:** Windows, macOS, Linux (wherever the editor runs)

## Troubleshooting

### "Character not found"
**Cause:** Invalid character index or character not recruited  
**Solution:** Enter valid index (0-15) for recruited character

### Progress shows 0% for all spells
**Cause:** Character hasn't equipped any espers yet  
**Solution:** This is expected for characters without esper exposure

### Spell not found in search
**Cause:** Typo in spell name or spell not in database  
**Solution:** Check spelling or try partial name match

### AP calculations seem incorrect
**Cause:** Current progress data may be inaccurate  
**Solution:** Reload save file to refresh spell data

## Known Limitations

- Simplified esper database (~7 espers vs 27 in game)
- Simplified spell database (~50 spells vs 54 in game)
- Doesn't track natural spell learning (character-specific)
- No esper application functionality (read-only plugin)
- Learn rates are estimates (may not match exact game values)
- Doesn't consider equipment that teaches spells

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Future Enhancements

Planned features for future versions:
- Complete esper database with all 27 espers
- Complete spell database with all 54 spells
- Natural spell learning tracking (character-specific)
- Equipment-based spell learning (e.g., Paladin Shield)
- Esper assignment recommendations with write permissions
- Multi-character spell progress comparison
- Visual progress bars and charts
- Export spell learning plan to text/JSON
- Integration with Esper Growth Optimizer
- Spell learning timeline (track history)
- Character-specific spell recommendations
- Boss-specific spell requirements

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
