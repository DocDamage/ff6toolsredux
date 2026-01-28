# Party Optimizer

An intelligent party composition analyzer and optimizer for the FF6 Save Editor. Analyzes all available characters, ranks them by multiple metrics, and provides data-driven recommendations for building the optimal party.

## Features

- **Character Rankings:** Rank all characters by combat power, magic, physical, tank, and speed ratings
- **Optimal Party Recommendations:** AI-driven party composition with role balancing
- **Current Party Analysis:** Detailed analysis of your active party with balance scoring
- **Character Comparison:** Head-to-head stat and rating comparisons
- **Role Detection:** Automatic character role classification (Tank, Mage, Physical DPS, etc.)
- **Balance Scoring:** Calculate party composition balance (0-100 scale)
- **Smart Filtering:** Configurable minimum level requirements
- **One-Click Apply:** Automatically update party composition from recommendations

## Installation

1. Open FF6 Save Editor (version 3.4.0 or higher)
2. Go to `Tools → Marketplace` or `Community → Marketplace`
3. Search for "Party Optimizer"
4. Click "Install"

## Usage

### Basic Usage

1. Load a save file in the FF6 Save Editor
2. Go to `Tools → Plugin Manager`
3. Find "Party Optimizer" in the installed plugins list
4. Click "Run"
5. Select an operation from the menu

### Available Operations

#### 1. Show Character Rankings

Displays comprehensive rankings of all available characters across multiple categories.

**Rankings Shown:**
- Top 3 by Overall Combat Power
- Top 3 Magic Users
- Top 3 Physical Attackers
- Top 3 Tanks
- Top 3 Fastest Characters

**Use Case:** Quickly identify strongest characters in each category for strategic team building.

**Example Output:**
```
=== CHARACTER RANKINGS ===

TOP OVERALL COMBAT POWER:
1. Terra (Lv45) - Mage - Power: 2840 - HP: 2450/2450
2. Sabin (Lv43) - Physical DPS - Power: 2720 - HP: 2680/2680
3. Edgar (Lv42) - Balanced - Power: 2650 - HP: 2340/2340

TOP MAGIC USERS:
1. Terra (Lv45) - Mage - Magic: 1950
2. Celes (Lv44) - Mage - Magic: 1880
3. Strago (Lv40) - Mage - Magic: 1750
...
```

#### 2. Recommend Optimal Party

Analyzes all characters and recommends the best 4-member party composition.

**Features:**
- Automatic role balancing (when enabled)
- Combat power optimization
- Diversity scoring
- One-click application to save file

**Balance Mode:**
- **Enabled:** Prefers diverse roles (1 Tank, 1 Mage, 1 Physical, 1 Fast)
- **Disabled:** Simply picks top 4 by combat power

**Use Case:** Building the strongest possible party for endgame content or boss battles.

**Steps:**
1. Select option 2 from menu
2. Review recommended party with stats
3. Check balance score
4. Confirm to apply changes to save file

**Example Output:**
```
=== RECOMMENDED OPTIMAL PARTY ===

Party Size: 4
Balance Mode: Enabled

1. Terra (Lv45) - Mage - Power: 2840 - HP: 2450/2450
2. Sabin (Lv43) - Physical DPS - Power: 2720 - HP: 2680/2680
3. Edgar (Lv42) - Tank - Power: 2650 - HP: 2980/2980
4. Locke (Lv41) - Fast Attacker - Power: 2580 - HP: 2140/2140

Total Combat Power: 10790
Average Level: 42
Balance Score: 85/100
```

#### 3. Analyze Current Party

Provides detailed analysis of your currently active party.

**Analysis Includes:**
- Party member list with roles and stats
- Total combat power
- Average level
- Total HP and percentage
- Role distribution breakdown
- Balance score with recommendations

**Use Case:** Evaluating your current party's strengths and weaknesses before major battles.

**Example Output:**
```
=== CURRENT PARTY ANALYSIS ===

Party Members: 4/4

1. Celes (Lv40) - Mage - Power: 2450 - HP: 2100/2100
2. Locke (Lv39) - Fast Attacker - Power: 2380 - HP: 1980/1980
3. Cyan (Lv38) - Physical DPS - Power: 2290 - HP: 2150/2150
4. Gau (Lv37) - Balanced - Power: 2210 - HP: 2050/2050

=== PARTY STATISTICS ===
Total Combat Power: 9330
Average Level: 38
Total HP: 8280/8280 (100.0%)

Role Distribution:
  Mage: 1
  Fast Attacker: 1
  Physical DPS: 1
  Balanced: 1

Balance Score: 78/100
✓ Party composition is excellent!
```

#### 4. Compare Two Characters

Performs head-to-head comparison of any two characters.

**Comparison Includes:**
- All base stats (Level, HP, MP, Attack, Defense, Magic, Mag Def, Speed, Stamina)
- All calculated ratings (Combat Power, Magic, Physical, Tank, Speed)
- Clear winner indicators for each category

**Use Case:** Deciding between two similar characters for the final party slot.

**Steps:**
1. Select option 4 from menu
2. Enter first character index (0-15)
3. Enter second character index (0-15)
4. Review detailed comparison

**Example Output:**
```
=== CHARACTER COMPARISON ===

Terra (Lv45) - Mage - Power: 2840 - HP: 2450/2450
vs
Celes (Lv44) - Mage - Power: 2750 - HP: 2380/2380

STAT COMPARISON:
Level: 45 vs 44 ← BETTER
Max HP: 2450 vs 2380 ← BETTER
Max MP: 580 vs 550 ← BETTER
Attack: 142 vs 138 ← BETTER
Defense: 125 vs 128 → BETTER
Magic: 165 vs 158 ← BETTER
Mag Def: 135 vs 142 → BETTER
Speed: 48 vs 51 → BETTER
Stamina: 38 vs 35 ← BETTER

RATING COMPARISON:
Combat Power: 2840 vs 2750 ← BETTER
Magic Rating: 1950 vs 1880 ← BETTER
Physical Rating: 1520 vs 1480 ← BETTER
Tank Rating: 1680 vs 1720 → BETTER
Speed Rating: 230 vs 239 → BETTER
```

## Configuration

The plugin supports several configuration options (modify at top of plugin.lua):

```lua
local config = {
    partySize = 4,              -- Standard party size
    showTopN = 3,               -- Show top N characters per category
    autoBalance = true,         -- Consider party balance in recommendations
    preferNaturalStats = true,  -- Prefer natural stats over equipment bonuses
    minLevel = 1,               -- Minimum level for consideration
}
```

### Configuration Options

- **partySize:** Number of party members (default: 4)
- **showTopN:** How many top characters to show per ranking category (default: 3)
- **autoBalance:** If true, recommendations prefer role diversity (default: true)
- **preferNaturalStats:** If true, prioritizes natural stats over equipment bonuses (default: true)
- **minLevel:** Minimum character level to include in analysis (default: 1)

## Rating Calculations

### Combat Power
Weighted formula considering all stats:
- Attack × 2.0
- Magic × 1.8
- Defense × 1.5
- Magic Defense × 1.5
- Speed × 1.2
- Stamina × 1.0

### Magic Rating
Focused on spellcasting ability:
- Magic × 3.0
- Magic Defense × 2.0
- Max MP × 0.5

### Physical Rating
Focused on physical combat:
- Attack × 3.0
- Defense × 2.0
- Stamina × 1.5
- Max HP × 0.3

### Tank Rating
Focused on damage absorption:
- Defense × 3.5
- Magic Defense × 3.5
- Stamina × 2.0
- Max HP × 0.8

### Speed Rating
Focused on turn frequency:
- Speed × 4.0
- Stamina × 1.0

## Role Detection

Characters are automatically classified into roles based on their highest rating:

- **Tank:** Highest defense and HP, absorbs damage
- **Mage:** Highest magic rating, casts spells
- **Physical DPS:** Highest physical rating, deals melee damage
- **Fast Attacker:** Highest speed rating, acts frequently
- **Balanced:** No dominant stat, versatile

## Balance Scoring

Party balance is scored 0-100 based on:
- **Role Diversity (40%):** More unique roles = higher score
- **Average Power (60%):** Higher combat power = higher score

**Score Interpretation:**
- **0-39:** Unbalanced party (consider diversifying roles)
- **40-69:** Decent party composition
- **70-100:** Excellent party composition

## Permissions

This plugin requires the following permissions:

- **read_save:** Required to read character and party data
- **write_save:** Required to update party composition (operation 2 only)
- **ui_display:** Required to show analysis results and menus

**Note:** Operation 2 (Recommend Optimal Party) modifies save data. Always backup your save before use!

## Use Cases

### For New Players
- Identify strongest characters available at current point
- Build balanced party for story progression
- Compare characters to make informed choices

### For Veterans
- Min-max party composition for endgame content
- Theory-craft optimal team builds
- Analyze character potential across playthroughs

### For Speedrunners
- Quickly identify fastest party composition
- Optimize for specific boss encounters
- Validate routing decisions with hard data

### For Challenge Runs
- Analyze low-level party viability
- Identify "worst" characters for difficulty
- Balance parties under self-imposed restrictions

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered
- **Platform:** Windows, macOS, Linux (wherever the editor runs)
- **Save Format:** Compatible with all FF6PR save formats

## Troubleshooting

### "Missing required API methods"
**Cause:** Editor version too old  
**Solution:** Update to FF6 Save Editor 3.4.0 or higher

### "No characters found in save file"
**Cause:** Save file not loaded or corrupted  
**Solution:** Ensure save file is properly loaded in editor

### "Not enough characters available for party size"
**Cause:** Fewer than 4 characters recruited  
**Solution:** Progress further in game to recruit more characters

### Rankings show unexpected results
**Cause:** Equipment bonuses heavily affecting stats  
**Solution:** Set `preferNaturalStats = true` in config to see base potential

### Balance score seems low
**Cause:** Party has duplicate roles or lacks diversity  
**Solution:** Enable `autoBalance = true` in config for recommendations

## Known Limitations

- Character names show as "Character #X" if name data not available
- Equipment bonuses included in calculations (may skew results)
- No support for character-specific abilities (Blitz, SwdTech, etc.)
- Party size fixed at 4 (game limitation)
- Analysis based on current stats (doesn't predict growth)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Future Enhancements

Planned features for future versions:
- Character name lookups from game database
- Equipment optimization recommendations
- Esper matching suggestions (best espers for each character)
- Growth potential analysis (stat gain predictions)
- Boss-specific party recommendations
- Export party builds to JSON
- Import community party builds
- Historical stat tracking across levels
- Multi-party management (for World of Ruin)

## License

MIT License

Copyright (c) 2026 FF6 Editor Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

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

Special thanks to the FF6 community for feedback and testing, and to the original FF6 development team for creating such a rich character system to analyze!
