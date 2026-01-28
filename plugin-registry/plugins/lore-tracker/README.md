# Strago's Lore Tracker

**Version:** 1.0.0  
**Category:** Tracker  
**Author:** FF6 Editor Team  
**Permissions:** read_save, ui_display

---

## Overview

The **Lore Tracker** plugin provides comprehensive tracking of all 24 Blue Magic spells (Lores) that Strago can learn throughout Final Fantasy VI. This plugin helps completionists ensure they don't miss any of Strago's unique abilities by providing detailed information about enemy sources, locations, and missable indicators for each Lore.

### Key Features

- ✅ **Complete Lore Checklist** - Track all 24 Lores with learned/unlearned status
- 📍 **Enemy Source Database** - See which enemies can teach each Lore
- 🗺️ **Location Guides** - Find out where to encounter each enemy
- ⚠️ **Missable Warnings** - Identify Lores only available in World of Balance
- 👑 **Boss Indicators** - Know which Lores require boss battles
- 📊 **Progress Tracking** - View completion percentage and statistics
- 🔍 **Flexible Filtering** - View learned, unlearned, or missable Lores only
- 📈 **Detailed Statistics** - Analyze your Lore collection progress

---

## Installation

1. Copy the `lore-tracker` folder to your `plugin-registry/plugins/` directory
2. The plugin will automatically appear in the FF6 Save Editor plugin menu
3. No additional configuration required

### Requirements

- FF6 Save Editor v3.4.0 or higher
- A save file with Strago recruited
- Read save and UI display permissions (automatically granted)

---

## Usage Guide

### Starting the Plugin

1. Open your FF6 save file in the editor
2. Navigate to the **Plugins** menu
3. Select **Strago's Lore Tracker**
4. The main menu will display with your current progress

### Main Menu Options

The plugin provides five main viewing modes:

#### 1. View All Lores (Complete List)

Displays all 24 Lores with:
- Learned/unlearned status markers
- Enemy sources for each Lore
- Location information
- Boss and missable indicators

**Example Output:**
```
==============================================================
STRAGO'S LORE TRACKER - Blue Magic Spells
==============================================================

Progress: 18 / 24 Lores learned (75.0%)

⚠️  WARNING: 1 missable Lore(s) not yet learned!

Lore List:
------------------------------------------------------------
[✓]  1. Condemn
       Enemies: Critic, Death Machine
       Location: Kefka's Tower

[ ]  2. Roulette
       Enemies: Doom Drgn, Didalos
       Location: Phoenix Cave, Kefka's Tower

[✓]  3. CleanSweep
       Enemies: Abolisher, Siegfried (boss) [BOSS]
       Location: Colosseum, Floating Continent

...
```

#### 2. View Learned Lores Only

Shows only the Lores that Strago has already learned. Perfect for:
- Reviewing your current abilities
- Verifying which Lores you have
- Checking proficiency levels

#### 3. View Unlearned Lores Only

Displays only Lores you haven't learned yet. Essential for:
- Planning which Lores to acquire next
- Seeing what you're missing
- Prioritizing rare or missable Lores

#### 4. View Missable Lores Only

Shows the two Lores that can only be obtained in World of Balance:
- **L? Pearl** - From Magic Urn boss in Owzer's Mansion
- **Quasar** - From Guardian boss on Narshe Cliffs

**Important:** If you're in World of Ruin and haven't learned these, they're permanently missable!

#### 5. View Lore Statistics

Displays comprehensive statistics including:
- Total Lores learned (count and percentage)
- Missable Lores status
- Boss-only Lores status
- Completion indicators

**Example Statistics:**
```
==============================================================
LORE STATISTICS
==============================================================

Total Lores Learned: 18 / 24 (75.0%)

Missable Lores: 1 / 2 learned
  ⚠️  Some missable Lores not yet obtained!

Boss-Only Lores: 4 / 6 learned

Lores Remaining: 6

==============================================================
```

---

## Complete Lore Reference

### All 24 Lores

| ID | Lore Name | Enemy Sources | Location | Missable? | Boss? |
|----|-----------|---------------|----------|-----------|-------|
| 1 | Condemn | Critic, Death Machine | Kefka's Tower | No | No |
| 2 | Roulette | Doom Drgn, Didalos | Phoenix Cave, Kefka's Tower | No | No |
| 3 | CleanSweep | Abolisher, Siegfried | Colosseum, Floating Continent | No | Boss |
| 4 | Aqua Rake | Actaneon, Rizopas | Fanatics' Tower, Thamasa Area | No | No |
| 5 | Aero | Sprinter, Innoc | Mt. Zozo, Veldt | No | No |
| 6 | Blow Fish | Anguiform, Nautiloid | Serpent Trench, Ocean (WoR) | No | No |
| 7 | Big Guard | Dark Force, Galypdes | Kefka's Tower, WoB Sealed Gate | No | Boss |
| 8 | Revenge | Misfit, Pm Stalker | Phoenix Cave, Mt. Zozo | No | No |
| 9 | Pearl Wind | Flan | Cave to South Figaro, Figaro Castle | No | No |
| 10 | L.5 Doom | Trapper | Fanatics' Tower | No | No |
| 11 | L.4 Flare | Apokryphos, Trapper | Fanatics' Tower, Kefka's Tower | No | No |
| 12 | L.3 Muddle | Apokryphos, Abolisher | Fanatics' Tower, Floating Continent | No | No |
| 13 | Reflect??? | Reach Frog | Triangle Island (WoR) | No | No |
| 14 | L? Pearl | Magic Urn | Owzer's Mansion (boss fight) | **YES** | Boss |
| 15 | Step Mine | Pm Stalker, Sky Base | Mt. Zozo, Floating Continent | No | No |
| 16 | Force Field | Anemone, Aquila | Fanatics' Tower, Narshe Caves | No | No |
| 17 | Dischord | Cephaler, Trixter | Ancient Castle, Colosseum | No | No |
| 18 | Sour Mouth | Grenade, Bug | Magitek Research Facility, Cave in the Veldt | No | No |
| 19 | Pep Up | Muus | Triangle Island (WoR) | No | No |
| 20 | Rippler | Skull Drgn, Galypdes | Kefka's Tower (boss), WoB Sealed Gate (boss) | No | Boss |
| 21 | Stone | Peepers, Mag Roader | Serpent Trench, Magitek Research Facility | No | No |
| 22 | Quasar | Dark Force, Guardian | Kefka's Tower, Narshe Cliffs (boss) | **YES** | Boss |
| 23 | Grand Train | Dullahan, Hidon | Colosseum (boss), Ebot's Rock (boss) | No | Boss |
| 24 | Exploder | Bomb, Grenade | Magitek Research Facility, Kefka's Tower | No | No |

### Missable Lores (⚠️ Critical!)

Only **2 Lores** are missable if not obtained in World of Balance:

#### L? Pearl (Lore #14)
- **Enemy:** Magic Urn (boss)
- **Location:** Owzer's Mansion (WoB only)
- **How to Get:** Fight Magic Urn during the Owzer's Mansion event
- **Why Missable:** Owzer's Mansion event only happens once in WoB

#### Quasar (Lore #22)
- **Enemy:** Guardian (boss)
- **Location:** Narshe Cliffs (WoB only)
- **How to Get:** Guardian can use Quasar during the Narshe defense
- **Why Missable:** Guardian boss fight only available in WoB
- **Alternative:** Also learnable from Dark Force in Kefka's Tower (WoR)

**Note:** While Quasar has a WoR source (Dark Force), it's much rarer there, so it's best to learn it from Guardian in WoB.

### Boss-Only Lores

Six Lores can only be learned from boss enemies:

1. **CleanSweep** - Siegfried (Colosseum)
2. **Big Guard** - Galypdes (Sealed Gate)
3. **L? Pearl** - Magic Urn (Owzer's Mansion) - MISSABLE
4. **Rippler** - Skull Drgn or Galypdes
5. **Quasar** - Guardian (Narshe) - MISSABLE
6. **Grand Train** - Dullahan (Colosseum) or Hidon (Ebot's Rock)

---

## Learning Lores: Strategy Guide

### How Lore Learning Works

1. **Strago must be in the party** during the battle
2. **The enemy must use the Lore** on your party
3. **Strago must survive** the battle (doesn't need to survive the Lore itself)
4. The Lore is automatically learned after battle

### Tips for Learning Lores

#### General Tips

- **Control/Sketch Strategy:** Use Control or Sketch on enemies to force them to use their Lores
- **Rage Strategy:** Use Gau's Rage with the same enemy to see what abilities it has
- **Defend & Wait:** Many enemies will eventually use their Lores if you wait long enough
- **Reflect Trick:** Some Lores can be reflected back to the enemy and then to you

#### Boss Lore Tips

- **Save before boss fights** if you need a specific Lore
- **Prepare for long battles** - bosses may not use Lores immediately
- **Use Control if available** to force the boss to use their Lore
- **Defeat the boss slowly** to give more chances for the Lore to be used

#### Missable Lore Priority

Complete these ASAP in World of Balance:

1. **L? Pearl** - During Owzer's Mansion event
   - Magic Urn uses it frequently
   - Have Strago in your party for this event!

2. **Quasar** - During Narshe defense (Guardian boss)
   - Guardian may not use it every time
   - Be prepared to defend for several turns
   - Alternative: Wait for Kefka's Tower (less reliable)

---

## Configuration Options

The plugin includes several configuration options (modifiable in plugin.lua):

```lua
local config = {
    showLearnedOnly = false,        -- Show only learned Lores
    showUnlearnedOnly = false,      -- Show only unlearned Lores
    showMissableOnly = false,       -- Show only missable Lores
    highlightMissable = true,       -- Highlight missable Lores
    sortByName = true,              -- Sort by name (false = by ID)
    showEnemyLocations = true,      -- Show enemy locations
}
```

### Customization

Edit these values in `plugin.lua` to change default behavior:

- **highlightMissable:** Set to `false` to disable missable warnings
- **sortByName:** Set to `false` to sort by Lore ID instead of name
- **showEnemyLocations:** Set to `false` for a compact view without locations

---

## Technical Details

### How It Works

The plugin:
1. Searches for Strago in the character roster
2. Reads Strago's learned spell data
3. Maps Lore IDs to spell system IDs
4. Checks proficiency levels to determine learned status
5. Displays information from the built-in Lore database

### Lore Detection Logic

Lores are detected by:
- Checking Strago's spell list for specific spell IDs
- Verifying proficiency > 0 (spell has been learned)
- Cross-referencing with the comprehensive Lore database

### Data Sources

All Lore information is hardcoded in the plugin:
- 24 Lores with complete details
- Enemy source mappings
- Location information
- Boss and missable flags

This means the plugin works offline without external data files.

---

## Troubleshooting

### "Strago not found in save file"

**Cause:** Strago hasn't been recruited yet or save file is corrupted

**Solutions:**
- Progress story until you recruit Strago (available after Thamasa events)
- Verify save file loads correctly in the main editor
- Check that character roster is intact

### Incorrect Lore Count

**Cause:** Lore spell ID mapping may not match your game version

**Solutions:**
- This is a known limitation in v1.0.0
- Future versions will improve Lore detection
- Manual verification recommended for critical Lores

### Missing Location Information

**Cause:** Some enemies appear in multiple locations

**Solutions:**
- The plugin shows primary locations
- Consult external guides for comprehensive enemy locations
- Some enemies are on the Veldt

---

## Use Cases

### For Completionists

- Track every Lore to achieve 100% collection
- Identify missable Lores before entering World of Ruin
- Monitor progress with detailed statistics

### For New Players

- Learn which enemies teach which Lores
- Find enemy locations without external guides
- Understand which Lores require special battles

### For Speedrunners

- Quickly check which essential Lores are learned
- Skip unnecessary Lore collection
- Focus on missable Lores only

### For Challenge Runs

- Verify Lore availability for specific strategies
- Plan party compositions around Lore collection
- Track boss-only Lores for completionist runs

---

## Known Limitations

1. **Read-Only:** Plugin cannot modify Lore data (would require write_save permission)
2. **Detection Accuracy:** Lore spell ID mapping may vary by ROM version
3. **No Spell Effects:** Doesn't show Lore damage/effects (reference guide feature for future)
4. **Static Database:** Enemy location data is hardcoded (accurate as of SNES version)

---

## Future Enhancements

Planned features for future versions:

- **Lore Effect Descriptions:** Show what each Lore does
- **Damage Calculator:** Calculate expected Lore damage based on Strago's stats
- **Better Detection:** Improved Lore learning detection across ROM versions
- **Export Checklist:** Export Lore checklist to text file
- **Visual Maps:** Integration with location maps
- **Lore Recommendations:** Suggest which Lores to prioritize

---

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

### v1.0.0 (January 16, 2026)
- Initial release
- Complete 24 Lore database
- Enemy source tracking
- Location guides
- Missable indicators
- Statistics tracking
- Multiple viewing modes

---

## Credits

**Plugin Development:** FF6 Editor Team  
**Lore Database:** Compiled from FF6 community resources  
**Testing:** FF6 speedrun and completionist communities

### Data Sources

- Final Fantasy VI SNES original game data
- FF6 community wikis and databases
- Speedrun routing guides

---

## Support

- **Issues:** Report bugs on GitHub Issues
- **Feature Requests:** Submit via GitHub Discussions
- **Community:** Join our Discord server

---

## License

MIT License - See [LICENSE](../../LICENSE) file for details

---

## FAQ

### Q: Can this plugin teach Strago new Lores?

**A:** No, this is a read-only tracker. It shows which Lores Strago has learned but cannot modify save data.

### Q: Do I need Strago in my active party to use this plugin?

**A:** No, you only need Strago recruited. The plugin reads his data regardless of party position.

### Q: What if I'm already in World of Ruin and missed a Lore?

**A:** Unfortunately, the two missable Lores (L? Pearl and Quasar) cannot be obtained if you're in WoR and didn't learn them in WoB. However, Quasar has an alternative source in Kefka's Tower.

### Q: Can I export my Lore checklist?

**A:** Not in v1.0.0, but this feature is planned for a future update.

### Q: Does this work with randomizers?

**A:** The plugin reads actual learned Lore data, so it works with randomizers. However, enemy location information may be inaccurate if enemy placements are randomized.

### Q: How accurate is the enemy location data?

**A:** Location data is based on the SNES original version and is 100% accurate for that version. Some locations may differ in GBA/mobile versions.

---

**Happy Lore Hunting!** 🎭✨
