# Save File Analyzer

A comprehensive save file statistics and analysis tool for the FF6 Save Editor that provides detailed insights into your game progress, completion percentages, and character/inventory metrics.

## Features

- **Overall Completion Percentage:** Weighted scoring across all game categories
- **Character Statistics:** Levels, HP/MP, experience, spells learned
- **Party Analysis:** Current party composition with power ratings
- **Inventory Metrics:** Unique items, total quantities, value estimates
- **Completion Checklist:** Track progress toward 100% completion
- **Visual Progress Bars:** ASCII progress indicators for quick scanning
- **Exportable Reports:** Generate shareable statistics reports
- **Read-Only Analysis:** Safe inspection without modifying save data

## Installation

1. Open FF6 Save Editor (version 3.4.0 or higher)
2. Go to `Tools → Marketplace` or `Community → Marketplace`
3. Search for "Save File Analyzer"
4. Click "Install"

## Usage

### Basic Usage

1. Load a save file in the FF6 Save Editor
2. Go to `Tools → Plugin Manager`
3. Find "Save File Analyzer" in the installed plugins list
4. Click "Run"
5. Select an operation from the menu

### Available Operations

#### 1. Show Comprehensive Statistics

Displays a complete dashboard of save file statistics.

**Includes:**
- Overall completion percentage with visual progress bar
- Character statistics (recruited, levels, HP/MP, spells, experience)
- Current party composition and power rating
- Inventory statistics (unique items, total count, value)

**Example Output:**
```
=== SAVE FILE ANALYSIS ===

Overall Completion: 67%
████████████████░░░░

=== CHARACTERS ===
Recruited: 12/14 (85%)
Average Level: 42
Level Range: 38 - 48
Total HP: 28,450
Total MP: 4,820
Total Spells Learned: 486
Total Experience: 1,248,670

=== CURRENT PARTY ===
Party Size: 4/4
Average Level: 44
Party Power: 18,450
Members:
  1. Terra (Lv48) - HP: 2850/2850
  2. Celes (Lv45) - HP: 2620/2620
  3. Edgar (Lv42) - HP: 2980/2980
  4. Sabin (Lv43) - HP: 2870/2870

=== INVENTORY ===
Unique Items: 142
Total Items: 428
Slots Used: 55%
Estimated Value: 42,800 gil
```

**Use Case:** Get a complete snapshot of your save file at a glance.

#### 2. Show Character Details

Displays detailed information for each recruited character.

**Includes:**
- Character name and level
- Current/Max HP and MP
- Total experience points
- Spells learned count
- Combat stats (Attack, Defense, Magic, Magic Defense)

**Example Output:**
```
=== CHARACTER DETAILS ===

Terra (Lv48)
  HP: 2850/2850  MP: 680/680
  Experience: 145,820
  Spells Learned: 54/54
  Attack: 142  Defense: 128
  Magic: 168  Mag.Def: 145

Locke (Lv45)
  HP: 2420/2420  MP: 480/480
  Experience: 128,450
  Spells Learned: 48/54
  Attack: 155  Defense: 118
  Magic: 98  Mag.Def: 105
...
```

**Use Case:** Quick reference for character capabilities and progression.

#### 3. Show Completion Checklist

Displays a checklist of major completion milestones with progress tracking.

**Includes:**
- All characters recruited (14/14)
- Max level character (Lv99)
- Average level 50+
- All spells learned
- 100+ unique items collected
- Visual progress bars per category

**Example Output:**
```
=== COMPLETION CHECKLIST ===

[✗] All Characters (12/14)
[✗] Max Level Character (Lv48/99)
[✓] Average Level 50+ (Avg: Lv52)
[✗] All Spells Learned (Avg: 48/54 per char)
[✓] 100+ Unique Items (142 items)

Progress Categories:

Characters: [████████░░] 85%
Levels:     [████░░░░░░] 43%
Spells:     [████████░░] 88%
Inventory:  [█████░░░░░] 55%
```

**Use Case:** Track progress toward completionist goals.

#### 4. Export Statistics Report

Generates a text-based statistics report that can be copied and shared.

**Includes:**
- Timestamp of report generation
- Summary statistics
- Character statistics
- Inventory statistics
- Formatted for easy sharing

**Example Output:**
```
=== FF6 SAVE FILE STATISTICS REPORT ===

Generated: 2026-01-16 14:30:00

SUMMARY
Overall Completion: 67%
Characters: 12/14
Average Level: 42
Party Size: 4/4

CHARACTER STATISTICS
  Total HP: 28,450
  Total MP: 4,820
  Total Experience: 1,248,670
  Total Spells Learned: 486

INVENTORY STATISTICS
  Unique Items: 142
  Total Items: 428
  Estimated Value: 42,800 gil

End of Report
```

**Use Case:** Share progress with community, compare with friends, track historical progress.

## Configuration

The plugin supports several configuration options (modify at top of plugin.lua):

```lua
local config = {
    showDetailedStats = true,      -- Show detailed breakdowns
    calculateCompletion = true,    -- Calculate completion percentages
    trackMissables = true,         -- Track missable items/events
    compareToIdeal = true,         -- Compare to ideal/max stats
    groupByCategory = true,        -- Group stats by category
}
```

### Configuration Options

- **showDetailedStats:** If true, shows comprehensive breakdowns
- **calculateCompletion:** Enables percentage calculations
- **trackMissables:** Tracks missable items/events (not yet fully implemented)
- **compareToIdeal:** Compares current stats to maximum possible
- **groupByCategory:** Organizes statistics by category

## Completion Calculation

Overall completion is calculated using weighted scoring:

```
Overall = (Characters × 30%) + (Inventory × 20%) + (Spells × 25%) + (Levels × 25%)
```

**Weight Breakdown:**
- **Characters (30%):** Recruited characters / 14 max
- **Inventory (20%):** Unique items collected / 255 max slots
- **Spells (25%):** Average spells learned / 54 max per character
- **Levels (25%):** Average level / 99 max

**Example:**
- 12/14 characters = 85% × 0.30 = 25.5
- 142 unique items = 55% × 0.20 = 11.0
- 48/54 avg spells = 88% × 0.25 = 22.0
- Level 42/99 avg = 42% × 0.25 = 10.5
- **Total: 69%**

## Statistics Reference

### Character Constants
- **Max Characters:** 14 (Terra through Umaro)
- **Max Level:** 99
- **Max Spells:** 54 per character
- **Max Rages:** 255 (Gau only)

### Inventory Constants
- **Max Inventory Slots:** 255
- **Max Gil:** 9,999,999

### Party Constants
- **Max Party Size:** 4 members
- **Max Parties:** 3 (World of Ruin only)

## Permissions

This plugin requires the following permissions:

- **read_save:** Required to read all save file data
- **ui_display:** Required to show analysis results

**Note:** This is a READ-ONLY plugin. It does not modify any save data.

## Use Cases

### For Casual Players
- Check overall progress at a glance
- Identify which characters need leveling
- Track party composition

### For Completionists
- Monitor progress toward 100% completion
- Identify missing collectibles
- Track all completion milestones

### For Speedrunners
- Quick character/party readiness check
- Verify route progression
- Export run statistics for comparison

### For Community Sharing
- Generate shareable progress reports
- Compare progress with other players
- Document playthrough milestones

### For Challenge Runs
- Track restricted run compliance
- Verify level requirements
- Document unique challenge progress

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered
- **Platform:** Windows, macOS, Linux (wherever the editor runs)
- **Save Format:** Compatible with all FF6PR save formats

## Troubleshooting

### "No characters found"
**Cause:** Save file not loaded or incompatible  
**Solution:** Ensure save file is properly loaded in the editor

### Statistics seem inaccurate
**Cause:** Cached data or partial save load  
**Solution:** Reload the save file completely

### Completion percentage unexpected
**Cause:** Weighted formula may not match personal goals  
**Solution:** Use individual category percentages instead

### "Unable to read inventory"
**Cause:** Inventory data structure incompatible  
**Solution:** Verify save file compatibility with editor version

### Report text truncated
**Cause:** Dialog display limitations  
**Solution:** Copy text immediately before closing dialog

## Known Limitations

- Read-only analysis (no modifications possible)
- Simplified completion formula (doesn't track all collectibles)
- Esper tracking not implemented
- Rage tracking not implemented
- No historical comparison (single snapshot only)
- No multi-save comparison
- Estimated inventory values are rough approximations
- Party power rating is simplified calculation

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Future Enhancements

Planned features for future versions:
- Complete collectibles tracking (espers, rages, etc.)
- Missable items/events tracking
- Historical statistics (track changes over time)
- Multi-save comparison (compare multiple save files)
- Visual charts and graphs
- Detailed boss/enemy defeat tracking
- Step counter and playtime tracking
- Gil spending analysis
- Character usage statistics (battles fought per character)
- Equipment collection tracking
- Colosseum completion tracking
- Dragon tracking (8 legendary dragons)
- Achievement system
- Export to PDF/HTML/JSON formats
- Screenshot generation of statistics
- Integration with cloud save features

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

## Acknowledgments

Special thanks to the FF6 community for feedback on desired statistics tracking and to all completionist players who inspired this tool!
