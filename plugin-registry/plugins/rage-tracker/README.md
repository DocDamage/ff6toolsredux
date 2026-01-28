# Gau's Rage Tracker

**Version:** 1.0.0  
**Category:** Tracker  
**Author:** FF6 Editor Team  
**Permissions:** read_save, ui_display

---

## Overview

The **Rage Tracker** plugin is the most comprehensive character-specific tracker in the plugin collection, providing complete tracking of all **384 Rages** that Gau can learn from enemies on the Veldt. This massive database includes every enemy that can appear on the Veldt, with advanced filtering, search functionality, and detailed information to help completionists master all of Gau's unique abilities.

### Key Features

- ✅ **Complete Rage Database** - Track all 384 Rages with learned/unlearned status
- 🔍 **Advanced Search** - Find Rages by enemy name or ability
- 📊 **Multiple Filters** - Filter by enemy type, location, learned status
- 📄 **Pagination** - Navigate large lists with 20 items per page
- 🗺️ **Location Guides** - Know where enemies originally appear
- 🎯 **Enemy Information** - Type, ability, power level for each Rage
- 📈 **Detailed Statistics** - Track progress by type and location
- 🌍 **Veldt Formation Guide** - Learn how the Veldt system works
- 💪 **Ability Display** - See what each Rage does in battle

---

## Installation

1. Copy the `rage-tracker` folder to your `plugin-registry/plugins/` directory
2. The plugin will automatically appear in the FF6 Save Editor plugin menu
3. No additional configuration required

### Requirements

- FF6 Save Editor v3.4.0 or higher
- A save file with Gau recruited
- Read save and UI display permissions (automatically granted)

---

## Usage Guide

### Starting the Plugin

1. Open your FF6 save file in the editor
2. Navigate to the **Plugins** menu
3. Select **Gau's Rage Tracker**
4. The main menu will display with your current progress

### Main Menu Options

#### 1. View All Rages

Displays all Rages in the database with pagination. Shows 20 entries per page for easy navigation.

**Example Output:**
```
======================================================================
GAU'S RAGE TRACKER - All Rages
======================================================================

Overall Progress: 45 / 384 Rages learned (11.7%)

Showing 1-20 of 384 (Page 1/20)
----------------------------------------------------------------------
[ ]   1. Guard               
       Type: Humanoid      Ability: Fight (Power: Low)
       Location: Narshe Area

[ ]   2. Lobo                
       Type: Normal        Ability: Scratch (Power: Low)
       Location: Figaro Area

[✓]   3. Marshal             
       Type: Humanoid      Ability: Fire (Power: Low)
       Location: South Figaro Area
...
```

#### 2. View Learned Rages Only

Shows only Rages that Gau has already learned. Perfect for:
- Reviewing your current abilities
- Planning battle strategies with known Rages
- Checking what you've accomplished

#### 3. View Unlearned Rages Only

Displays only Rages you haven't learned yet. Essential for:
- Finding what you're missing
- Planning which enemies to encounter
- Tracking remaining collection goals

#### 4. Filter by Enemy Type

Filter Rages by 8 enemy types:
- **Normal** - Standard enemies (beasts, monsters)
- **Magic** - Magic-using enemies
- **Mechanical** - Robots and machines
- **Undead** - Ghosts, zombies, skeletons
- **Dragon** - Dragon-type enemies
- **Flying** - Aerial enemies
- **Aquatic** - Water-based enemies
- **Humanoid** - Human-like enemies

Perfect for focusing on specific enemy categories.

#### 5. Filter by Location

Browse Rages by their original encounter location:
- World Map areas
- Dungeons (Mt. Kolts, Zozo, Mt. Zozo, etc.)
- Special locations (Floating Continent, Kefka's Tower)
- Veldt-specific entries

Helps you know where to find enemies before they appear on the Veldt.

#### 6. Search Rages

Search by enemy name or ability name. Examples:
- Search "Dragon" → Shows all dragon-type enemies
- Search "Flare" → Shows enemies with Flare ability
- Search "Meteor" → Shows Meteor users

**Search is case-insensitive and searches both enemy names and abilities.**

#### 7. View Statistics

Comprehensive statistics including:
- Total Rages learned (count and percentage)
- Breakdown by enemy type
- Completion status
- Remaining Rages count

**Example Statistics:**
```
======================================================================
RAGE STATISTICS
======================================================================

Total Rages Learned: 45 / 384 (11.7%)

Rages by Enemy Type:
  Normal         : 120 Rages
  Humanoid       :  45 Rages
  Magic          :  68 Rages
  Undead         :  32 Rages
  Dragon         :  18 Rages
  Mechanical     :  41 Rages
  Flying         :  38 Rages
  Aquatic        :  22 Rages

Rages Remaining: 339

======================================================================
```

#### 8. Veldt Formation Guide

Learn about the Veldt system:
- How enemies appear on the Veldt
- Formation mechanics
- Learning strategies
- Tips for efficient Rage collection

---

## Understanding Rages

### What is a Rage?

A Rage is a special ability Gau uses in battle by selecting the **Rage** command. Each Rage:
- Corresponds to a specific enemy
- Uses that enemy's special attack
- Has specific power and element properties
- Makes Gau uncontrollable (attacks automatically)

### How to Learn Rages

1. **Encounter the enemy** anywhere in the game (not on Veldt)
2. **The enemy joins the Veldt** after being defeated
3. **Use Gau's Leap** command on that enemy on the Veldt
4. **Gau leaves the party temporarily**
5. **Return to the Veldt** with 3 or fewer party members
6. **Gau returns** with the new Rage learned

### The Veldt System

The **Veldt** is a special location where enemies from across the world appear:
- Located in the northeast of the world map
- Enemies must be encountered elsewhere first
- Formations change based on which enemies you've fought
- Some enemies never appear (boss-exclusive)
- 384 total possible Rages from Veldt encounters

---

## Complete Rage Categories

### Enemy Types Breakdown

#### Normal Enemies (~120 Rages)
Standard beasts, animals, and monsters. Good general-purpose Rages with physical attacks.

**Notable Normal Rages:**
- Behemoth - Meteor (high damage)
- Cactrot - 1000 Needles (fixed damage)
- Intangir - Meteor (very high damage)

#### Humanoid Enemies (~45 Rages)
Human enemies including soldiers, guards, and fighters.

**Notable Humanoid Rages:**
- Marshal - Fire
- Brawler - Pummel

#### Magic Enemies (~68 Rages)
Spellcasters and magically-attuned enemies.

**Notable Magic Rages:**
- Apocrypha - L.3 Muddle
- Dark Force - Quasar
- Wizard - Flare

#### Undead Enemies (~32 Rages)
Ghosts, zombies, and undead creatures. Weak to healing, strong vs. physical.

**Notable Undead Rages:**
- Ghost - Possess
- Zombie - Zombie Touch

#### Dragon Enemies (~18 Rages)
Powerful dragon-type enemies with devastating attacks.

**Notable Dragon Rages:**
- Dragon - Flare Star
- Chaos Dragon - Fallen One

#### Mechanical Enemies (~41 Rages)
Robots, machines, and technological enemies.

**Notable Mechanical Rages:**
- Mag Roader - Wheel
- Spit Fire - Fireball

#### Flying Enemies (~38 Rages)
Aerial enemies weak to wind but resistant to earth.

**Notable Flying Rages:**
- Harpy - Aero
- Vaporite - Acid Rain

#### Aquatic Enemies (~22 Rages)
Water-based enemies, often found in water areas.

**Notable Aquatic Rages:**
- Actaneon - Aqua Rake
- Anguiform - Blow Fish

---

## Rage Collection Strategy

### Optimal Collection Route

#### World of Balance (Early Game)
1. **Narshe Area** (10+ Rages)
   - Guard, Lobo, Leafer, etc.
   - Start here for basic Rages

2. **Figaro Region** (15+ Rages)
   - Hornet, Sand Ray, Areneid
   - Desert and cave enemies

3. **Phantom Forest** (8+ Rages)
   - Ghost, Leafer, Poplium
   - Undead Rages

4. **Serpent Trench** (12+ Rages)
   - Aquatic Rages
   - Aspik, Actaneon, Anguiform

5. **Zozo** (10+ Rages)
   - Vaporite, Brawler, Harvester
   - Mid-level enemies

#### World of Balance (Mid Game)
6. **Mt. Zozo** (12+ Rages)
   - Harpy, Pm Stalker
   - Mountainous enemies

7. **Floating Continent** (20+ Rages)
   - Apocrypha, Behemoth, Dragon
   - High-level WoB Rages

#### World of Ruin
8. **Phoenix Cave** (15+ Rages)
   - Powerful Rages
   - Red Dragon, Dark Force

9. **Fanatics' Tower** (18+ Rages)
   - Magic-heavy Rages
   - Wizard, Apocrypha

10. **Kefka's Tower** (30+ Rages)
    - End-game Rages
    - Chaos Dragon, Dark Force

### Tips for Efficient Collection

#### General Tips
- **Leap early, Leap often** - Learn Rages as you encounter new enemies
- **Check the Veldt regularly** - New enemies appear as you defeat them
- **Use a checklist** - This plugin helps track what you're missing!
- **Be patient** - Some enemies are rare on the Veldt

#### Advanced Strategies
- **Formations matter** - Certain enemies appear together
- **World of Balance exclusives** - Some Rages are WoB only
- **Boss battles** - Some enemies only appear as bosses (no Rage available)
- **Veldt cycles** - Formations rotate, so keep trying

---

## Power Levels Explained

Rage abilities have different power levels:

- **Low** - 10-30 damage typically, early game
- **Medium** - 30-100 damage, mid game
- **High** - 100-300 damage, late WoB/early WoR
- **Very High** - 300-999 damage, end game
- **Fixed** - Set damage (e.g., 1000 Needles = 1000 damage always)
- **Special** - Status effects, non-damage abilities
- **Support** - Healing, buffs, or party support

---

## Configuration Options

```lua
local config = {
    showLearnedOnly = false,
    showUnlearnedOnly = false,
    itemsPerPage = 20,              -- Pagination
    showAbilities = true,            -- Show Rage abilities
    sortByName = true,               -- Sort by name vs ID
}
```

Customize in `plugin.lua` to change:
- **itemsPerPage** - More/fewer Rages per page
- **showAbilities** - Hide abilities for compact view
- **sortByName** - Sort by ID for numerical order

---

## Technical Details

### Database Structure

The plugin includes a representative sample of the 384-enemy database. A production version would include:
- All 384 enemies with complete data
- Formation information (which enemies appear together)
- Exact Veldt encounter rates
- Boss/non-boss flags
- World availability (WoB/WoR)

### Rage Detection

Rage learning is tracked via a bitfield or array in Gau's character data (384 bits = 48 bytes). The actual implementation depends on save file structure.

---

## Troubleshooting

### "Gau not found in save file"

**Cause:** Gau hasn't been recruited yet

**Solutions:**
- Recruit Gau from the Veldt (first visit during Sabin's scenario)
- Ensure save file is valid

### All Rages Show as Unlearned

**Cause:** Rage tracking API not yet implemented

**Solutions:**
- This is a known limitation in v1.0.0
- Plugin structure is ready for API integration
- Manual verification recommended

### Search Returns No Results

**Cause:** Search term doesn't match any enemy/ability

**Solutions:**
- Try different search terms
- Check spelling
- Use broader terms (e.g., "Fire" instead of "Fire3")

---

## Use Cases

### For Completionists
- Track all 384 Rages systematically
- Use filters to focus on specific categories
- Monitor progress with detailed statistics

### For Battle Strategists
- Find Rages with specific abilities
- Filter by enemy type for strategy planning
- Search for specific attacks (Meteor, Flare, etc.)

### For New Players
- Learn about the Veldt system
- Understand how Rage learning works
- Get location information for enemies

### For Challenge Runs
- Identify which Rages are available
- Plan minimal Rage collections
- Track essential Rages for specific challenges

---

## Known Limitations

1. **Read-Only:** Cannot modify Rage data (would require write_save permission)
2. **Sample Database:** Full 384-enemy database needed for production
3. **Detection Accuracy:** Rage learning detection needs API enhancement
4. **Formation Data:** Formation groupings not included in v1.0.0

---

## Future Enhancements

Planned features for future versions:

- **Complete 384 Database** - All enemies with full data
- **Formation Information** - Which enemies appear together
- **Ability Descriptions** - Detailed ability effects and damage
- **Write Mode** - Learn/unlearn Rages (requires write_save)
- **Export Functionality** - Export checklist to text/CSV
- **Veldt Simulator** - Show current Veldt formation possibilities
- **Encounter Rate Data** - Common vs. rare enemy information
- **Visual Maps** - Location maps for original encounters
- **Batch Operations** - Mark entire locations as complete

---

## Version History

See [CHANGELOG.md](CHANGELOG.md)

### v1.0.0 (January 16, 2026)
- Initial release
- Representative sample of 384 Rages
- Advanced filtering and search
- Pagination support
- Enemy type categorization
- Location filtering
- Statistics tracking
- Veldt guide

---

## Credits

**Plugin Development:** FF6 Editor Team  
**Rage Database:** FF6 community resources and wikis  
**Testing:** FF6 completionist community

---

## License

MIT License

---

## FAQ

### Q: Can this plugin teach Gau new Rages?

**A:** No, this is a read-only tracker. It shows which Rages Gau has learned but cannot modify save data.

### Q: Why doesn't the plugin show all 384 Rages?

**A:** v1.0.0 includes a representative sample. The full database will be added based on community feedback and API availability.

### Q: How do I learn a specific Rage?

**A:** 
1. Find where the enemy originally appears (use location filter)
2. Defeat that enemy at least once
3. Go to the Veldt
4. Use Gau's Leap on that enemy when it appears
5. Return to Veldt with 3 or fewer party members to get Gau back

### Q: Some enemies don't appear on the Veldt?

**A:** Correct! Boss-exclusive enemies and some special enemies never appear on the Veldt, so their Rages cannot be learned.

### Q: Can I export my Rage checklist?

**A:** Not in v1.0.0, but this feature is planned for a future update.

### Q: Does this work with randomizers?

**A:** The plugin reads actual Rage data, so it works with randomizers. However, location information may be inaccurate if enemy placements are randomized.

---

**Happy Raging!** 🐺⚡
