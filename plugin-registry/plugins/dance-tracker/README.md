# Mog's Dance Tracker

**Version:** 1.0.0  
**Category:** Tracker  
**Author:** FF6 Editor Team  
**Permissions:** read_save, ui_display

---

## Overview

The **Dance Tracker** plugin provides comprehensive tracking of all 8 dances that Mog can learn throughout Final Fantasy VI. Each dance is learned by visiting specific locations and provides Mog with four random abilities that trigger during battle. This plugin helps you collect all of Mog's unique dance abilities by showing which dances you've learned, where to find remaining dances, and what each dance does.

### Key Features

- ✅ **Complete Dance Checklist** - Track all 8 dances with learned/unlearned status
- 📍 **Location Guides** - Detailed maps showing where to learn each dance
- 🎭 **Dance Effects** - View all four abilities for each dance
- 🌍 **World Availability** - Know which dances are WoB vs WoR
- 📊 **Progress Tracking** - View completion percentage and statistics
- 🗺️ **Interactive Location Guide** - Find where to step to learn each dance
- 🎨 **Clean Interface** - Easy-to-read display with organized information

---

## Installation

1. Copy the `dance-tracker` folder to your `plugin-registry/plugins/` directory
2. The plugin will automatically appear in the FF6 Save Editor plugin menu
3. No additional configuration required

### Requirements

- FF6 Save Editor v3.4.0 or higher
- A save file with Mog recruited
- Read save and UI display permissions (automatically granted)

---

## Usage Guide

### Starting the Plugin

1. Open your FF6 save file in the editor
2. Navigate to the **Plugins** menu
3. Select **Mog's Dance Tracker**
4. The main menu will display with your current progress

### Main Menu Options

#### 1. View All Dances (Complete List)

Displays all 8 dances with:
- Learned/unlearned status
- Location information
- World availability (WoB/WoR)
- All four dance abilities

**Example Output:**
```
==============================================================
MOG'S DANCE TRACKER - Dance Abilities
==============================================================

Progress: 5 / 8 Dances learned (62.5%)

Dance List:
------------------------------------------------------------
[✓] 1. Wind Song
       Location: Narshe (outside, where Mog is found)
       Available: WoB and WoR
       Effects (random selection):
       • Sonic Boom - Wind damage to one enemy
       • Plasma - Lightning damage to all enemies
       • Harvester - Death to one enemy
       • Wombat - Poison damage to all enemies

[ ] 2. Forest Suite
       Location: Phantom Forest
       Available: WoB and WoR
       Effects (random selection):
       • Will o' Wisp - Fire damage to random enemies
       • Apparition - Drain HP from one enemy
       • Poltergeist - Image status to party
       • Specter - Invisible status to party
...
```

#### 2. View Learned Dances Only

Shows only dances Mog has already learned. Useful for:
- Reviewing available abilities
- Planning battle strategies
- Checking your collection

#### 3. View Unlearned Dances Only

Displays only dances you haven't learned yet. Essential for:
- Finding where to go next
- Completing your dance collection
- Quick reference for remaining dances

#### 4. Location Guide (Where to Learn)

Comprehensive guide showing:
- Exact location for each dance
- Map area names
- World availability
- Special notes for tricky dances

**Example Location Guide:**
```
==============================================================
DANCE LOCATION GUIDE
==============================================================

To learn a dance, Mog must step on the location where
the dance can be learned. The dance is learned immediately.

------------------------------------------------------------

1. Wind Song
   Location: Narshe (outside, where Mog is found)
   Map Area: Narshe Exterior
   Available: WoB and WoR

2. Forest Suite
   Location: Phantom Forest
   Map Area: Phantom Forest
   Available: WoB and WoR

...
```

#### 5. View Dance Statistics

Displays comprehensive statistics:
- Total dances learned
- WoB vs WoR availability breakdown
- Dance category analysis
- Completion status

---

## Complete Dance Reference

### All 8 Dances

| ID | Dance Name | Location | World | Category |
|----|------------|----------|-------|----------|
| 1 | Wind Song | Narshe Exterior | WoB/WoR | Elemental |
| 2 | Forest Suite | Phantom Forest | WoB/WoR | Elemental |
| 3 | Desert Aria | Figaro Castle area | WoB/WoR | Elemental |
| 4 | Love Sonata | Darill's Tomb | WoR | Support |
| 5 | Earth Blues | Mt. Kolts | WoB/WoR | Elemental |
| 6 | Water Rondo | Serpent Trench | WoB/WoR | Elemental |
| 7 | Dusk Requiem | Zozo (evening) | WoB/WoR | Status |
| 8 | Snowman Rondo | Narshe Mines | WoB/WoR | Elemental |

### Dance Details

#### 1. Wind Song
- **Location:** Narshe (outside, where you first find Mog)
- **Abilities:** Sonic Boom, Plasma, Harvester, Wombat
- **Strategy:** Good all-around offensive dance with instant death
- **Element Focus:** Wind and Lightning

#### 2. Forest Suite
- **Location:** Phantom Forest (where you meet the phantom train)
- **Abilities:** Will o' Wisp, Apparition, Poltergeist, Specter
- **Strategy:** Fire damage with party support (Image, Invisible)
- **Element Focus:** Fire with utility

#### 3. Desert Aria
- **Location:** Figaro Castle (WoB) or Kohlingen Desert (WoR)
- **Abilities:** Sand Storm, Antlion, Dissolve, Kitty
- **Strategy:** Contains Kitty (reduces enemy HP to single digits!)
- **Element Focus:** Wind and Earth with powerful HP manipulation

#### 4. Love Sonata ⭐ (WoR Only)
- **Location:** Darill's Tomb (World of Ruin only)
- **Abilities:** Aurora, Serenade, Minuet, Love Token
- **Strategy:** Best healing/support dance - Aurora heals, Minuet gives Haste
- **Element Focus:** Non-elemental support
- **Note:** Only dance that's WoR exclusive

#### 5. Earth Blues
- **Location:** Mt. Kolts (where you first meet Sabin)
- **Abilities:** Avalanche, Cave In, Landslide, Sonic Boom
- **Strategy:** All earth-based damage, great vs. earth-weak enemies
- **Element Focus:** Earth

#### 6. Water Rondo
- **Location:** Serpent Trench or Lethe River
- **Abilities:** El Niño, Specter, Plasma, Apparition
- **Strategy:** Water damage with HP drain
- **Element Focus:** Water and Lightning

#### 7. Dusk Requiem
- **Location:** Zozo or any town at night
- **Abilities:** Pois. Frog, Evil Toot, Sour Mouth, Snare
- **Strategy:** Best status ailment dance - Confusion, Poison, Stop, Level reduction
- **Element Focus:** Status effects
- **Note:** Must be evening hours in-game

#### 8. Snowman Rondo
- **Location:** Narshe Mines or Umaro's Cave
- **Abilities:** Absolute Zero, Surge, Avalanche, Snare
- **Strategy:** Ice and earth damage with stop
- **Element Focus:** Ice and Earth

---

## How Dance Learning Works

### Learning Mechanics

1. **Mog must be in the party** when visiting the location
2. **Step on the specific location** where the dance can be learned
3. **The dance is learned immediately** - no battle required
4. **Each dance has a specific trigger location**

### Dance Trigger Mechanics in Battle

- **Random Selection:** Each time Dance is used, one of four abilities is randomly chosen
- **25% Chance Each:** Each ability has equal probability
- **Can't Control:** You cannot choose which ability triggers
- **Multiple Uses:** Can use Dance multiple times per battle for different effects

---

## Dance Learning Guide

### Optimal Learning Route (WoB)

1. **Wind Song** - Narshe (where you recruit Mog)
2. **Earth Blues** - Mt. Kolts (early game area)
3. **Forest Suite** - Phantom Forest (Phantom Train quest)
4. **Water Rondo** - Serpent Trench (early WoB)
5. **Desert Aria** - Figaro Castle desert
6. **Snowman Rondo** - Narshe Mines
7. **Dusk Requiem** - Zozo (go at night)

### World of Ruin Additions

8. **Love Sonata** - Darill's Tomb (WoR only)

### Tips for Learning Dances

#### General Tips

- **Visit early:** Learn dances as soon as you recruit Mog
- **Explore thoroughly:** Some locations are easy to miss
- **Check both worlds:** Most dances work in both WoB and WoR
- **Evening visits:** Dusk Requiem requires nighttime

#### Specific Location Tips

**Wind Song:** The exact spot where you first find Mog in Narshe

**Forest Suite:** Anywhere in the Phantom Forest area

**Desert Aria:** 
- WoB: Near Figaro Castle in the desert
- WoR: Kohlingen desert area

**Love Sonata:** 
- Inside Darill's Tomb
- WoR exclusive - don't miss it!

**Earth Blues:** Inside Mt. Kolts caves

**Water Rondo:**
- Serpent Trench underwater area
- Also Lethe River

**Dusk Requiem:**
- Zozo town
- Or any town during evening hours
- Game time must be evening/night

**Snowman Rondo:**
- Narshe underground mines
- Umaro's Cave area

---

## Battle Strategy Guide

### Best Dances for Different Situations

#### **Boss Battles**
1. **Desert Aria** - Kitty can reduce boss HP dramatically
2. **Love Sonata** - Healing and Haste support
3. **Wind Song** - Harvester for instant death attempts

#### **Regular Encounters**
1. **Dusk Requiem** - Status ailments to control enemies
2. **Water Rondo** - Damage with HP drain
3. **Snowman Rondo** - Ice damage with stop

#### **Multi-Enemy Groups**
1. **Wind Song** - Plasma hits all enemies
2. **Earth Blues** - Multiple earth attacks hit all
3. **Dusk Requiem** - Mass status ailments

#### **Support Needs**
1. **Love Sonata** - Best healing and buffs
2. **Forest Suite** - Image/Invisible for defense

---

## Configuration Options

```lua
local config = {
    showLearnedOnly = false,        -- Show only learned dances
    showUnlearnedOnly = false,      -- Show only unlearned dances
    showEffects = true,             -- Show dance abilities
    sortByName = true,              -- Sort by name (false = by ID)
}
```

Customize these in `plugin.lua` to change default behavior.

---

## Technical Details

### How It Works

1. Searches for Mog in character roster
2. Checks dance learning status (implementation-specific)
3. Displays dance information from built-in database
4. Shows locations and effects

### Dance Detection

Dance learning is tracked via character-specific flags in save data. The exact implementation depends on save file structure.

---

## Troubleshooting

### "Mog not found in save file"

**Cause:** Mog hasn't been recruited yet

**Solutions:**
- Recruit Mog from Narshe (after Mobliz events)
- Ensure save file is valid

### All Dances Show as Unlearned

**Cause:** Dance tracking not yet implemented in save structure

**Solutions:**
- This is a known limitation in v1.0.0
- Plugin will track properly once dance API is implemented

---

## Use Cases

### For Completionists
- Collect all 8 dances
- Verify each location
- Track progress with statistics

### For New Players
- Learn where to find each dance
- Understand dance effects
- Plan dance learning route

### For Battle Strategists
- Review available dance abilities
- Plan which dances to use
- Understand random effect probabilities

---

## Known Limitations

1. **Read-Only:** Cannot modify dance data
2. **Detection:** Dance learning detection may need API enhancement
3. **Static Database:** Location data is hardcoded

---

## Version History

See [CHANGELOG.md](CHANGELOG.md)

### v1.0.0 (January 16, 2026)
- Initial release
- Complete 8 dance database
- Location guides
- Effect descriptions
- Statistics tracking

---

## Credits

**Plugin Development:** FF6 Editor Team  
**Dance Database:** FF6 community resources

---

## License

MIT License

---

**Happy Dancing!** 🎭✨
