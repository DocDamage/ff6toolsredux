# Sabin's Blitz Input Trainer

**Version:** 1.0.0  
**Category:** Reference  
**Author:** FF6 Editor Team  
**Permissions:** read_save, ui_display

---

## Overview

The **Blitz Input Trainer** plugin provides a comprehensive reference and practice guide for all 8 of Sabin's Blitz techniques. Blitzes are special moves executed by entering directional button sequences during battle, and this plugin helps you master each one with detailed input diagrams, damage calculations, and strategic recommendations.

### Key Features

- ✅ **Complete Blitz Reference** - All 8 Blitzes with full details
- ⌨️ **Visual Input Diagrams** - Arrow symbols showing button sequences
- 📊 **Damage Calculator** - Estimates based on Sabin's current stats
- 🎯 **Difficulty Ratings** - Know which inputs are easy vs. hard
- 📖 **Quick Reference Card** - One-page overview of all Blitzes
- 💡 **Strategic Recommendations** - Best Blitz for each situation
- 📐 **Damage Formulas** - Understand how damage is calculated
- 🎓 **Input System Guide** - Learn how the Blitz system works

---

## Installation

1. Copy the `blitz-trainer` folder to your `plugin-registry/plugins/` directory
2. The plugin will automatically appear in the FF6 Save Editor plugin menu
3. Works as a reference tool even without Sabin in your save file

### Requirements

- FF6 Save Editor v3.4.0 or higher
- Optional: A save file with Sabin (for damage calculations)
- Read save and UI display permissions (automatically granted)

---

## Usage Guide

### Starting the Plugin

1. Open your FF6 save file in the editor (or use without a save file)
2. Navigate to the **Plugins** menu
3. Select **Sabin's Blitz Input Trainer**
4. Choose from six viewing modes

### Main Menu Options

#### 1. Quick Reference Card

One-page overview showing all 8 Blitzes with:
- Blitz names
- Input sequences (visual)
- Learn levels
- Targets
- Estimated damage

**Perfect for:** Quick lookup during gameplay

#### 2. View Individual Blitz

Detailed view of a specific Blitz including:
- Complete input sequence with diagrams
- Damage formula and estimates
- Strategy recommendations
- Input tips and tricks

**Perfect for:** Learning a specific Blitz

#### 3. View All Blitzes

Complete compendium with every detail for all 8 Blitzes.

**Perfect for:** Comprehensive study session

#### 4. Input System Guide

Learn how the Blitz system works:
- How to perform Blitzes in battle
- Direction input explanations
- Common input patterns
- General tips for success

**Perfect for:** Understanding the mechanics

#### 5. Damage Calculator

See estimated damage for all Blitzes based on Sabin's current stats.

**Perfect for:** Planning battle strategies

#### 6. Blitz Recommendations

Strategic guide showing:
- Best Blitz for each game phase
- Boss battle recommendations
- Random encounter suggestions
- Easiest vs. hardest inputs

**Perfect for:** Optimizing your Blitz usage

---

## Complete Blitz Reference

### All 8 Blitzes

| Blitz | Input | Level | Difficulty | Target | Type |
|-------|-------|-------|------------|--------|------|
| Raging Fist | ← → | 1 | Easy | Single | Physical |
| Aura Cannon | ↓ ↙ ← | 6 | Easy | Single | Physical |
| Meteor Strike | ← ↖ ↑ ↗ → | 10 | Hard | Single | Physical (4x) |
| Suplex | ↑ ↑ ↓ ↓ ← → | 15 | Medium | Single | Physical |
| Fire Dance | ← ↙ ↓ ↘ → | 23 | Hard | All | Magic/Fire |
| Mantra | → ← → ← → ← | 30 | Medium | All Allies | Healing |
| Air Blade | ↑ ↗ → ↘ ↓ ↙ ← | 42 | Very Hard | All | Magic/Wind |
| Phantom Rush | ← ↑ → ↓ ↙ | 70 | Hard | Random | Physical (4x) |

---

## Detailed Blitz Descriptions

### 1. Raging Fist (← →)
**Learn Level:** 1 | **Difficulty:** Easy

Your starting Blitz and the easiest input. Simple left-right motion deals solid physical damage to one enemy.

**Strategy:** Reliable early-game damage dealer. Use until you learn better Blitzes.

**Damage:** 110 + (Level² × Vigor) / 256

---

### 2. Aura Cannon (↓ ↙ ←)
**Learn Level:** 6 | **Difficulty:** Easy

Enhanced version of Raging Fist with better damage. Quarter-circle input introduces diagonals.

**Strategy:** Replace Raging Fist at level 6. Still easy but stronger.

**Damage:** 130 + (Level² × Vigor) / 256

**Tip:** Practice the smooth motion from down to down-left to left.

---

### 3. Meteor Strike (← ↖ ↑ ↗ →)
**Learn Level:** 10 | **Difficulty:** Hard

Semi-circle motion delivers four consecutive hits for massive damage. Best single-target Blitz.

**Strategy:** Your go-to for boss battles. Worth mastering the input.

**Damage:** [150 + (Level² × Vigor) / 256] × 4

**Tip:** Take your time with each direction. Smooth quarter-circle from left to right.

---

### 4. Suplex (↑ ↑ ↓ ↓ ← →)
**Learn Level:** 15 | **Difficulty:** Medium

The legendary Phantom Train Suplex! Long input but all cardinal directions.

**Strategy:** Try it on EVERYTHING. Works on more enemies than you'd think!

**Damage:** 120 + (Level² × Vigor) / 256

**Tip:** No diagonals makes this easier than it looks. Just take your time.

---

### 5. Fire Dance (← ↙ ↓ ↘ →)
**Learn Level:** 23 | **Difficulty:** Hard

First multi-target Blitz! Fire-elemental quarter-circle hits all enemies.

**Strategy:** Great for random encounters and ice-weak enemies.

**Damage:** (106 × Level / 2) + Magic Power [to all enemies]

**Tip:** Smooth quarter-circle with all diagonals. Practice in low-risk battles.

---

### 6. Mantra (→ ← → ← → ←)
**Learn Level:** 30 | **Difficulty:** Medium

Free party-wide healing! Alternating right-left pattern.

**Strategy:** Keep Sabin's HP high for better healing output. No MP cost!

**Healing:** (Sabin's Max HP + Current HP) / 16 [to all allies]

**Tip:** Six inputs total - just alternate right and left.

---

### 7. Air Blade (↑ ↗ → ↘ ↓ ↙ ←)
**Learn Level:** 42 | **Difficulty:** Very Hard

Most complex input in the game. Three-quarter circle hits all enemies with wind.

**Strategy:** Strongest AoE Blitz. Practice until muscle memory kicks in.

**Damage:** (127 × Level / 2) + Magic Power [to all enemies]

**Tip:** Start up, sweep clockwise to down, end left. Seven total inputs.

---

### 8. Phantom Rush (← ↑ → ↓ ↙)
**Learn Level:** 70 | **Difficulty:** Hard

Ultimate Blitz! Box pattern + diagonal delivers four hits of maximum damage.

**Strategy:** Endgame destroyer. Can deal 39,996 total damage (9,999 × 4).

**Damage:** Max 9999 per hit × 4 hits to random enemies

**Tip:** Box (left-up-right-down) then diagonal (down-left).

---

## Battle Strategy Guide

### Early Game (Levels 1-20)
- **Main Blitz:** Aura Cannon
- **Alternative:** Raging Fist (if easier)
- **Boss Strategy:** Meteor Strike (once learned)

### Mid Game (Levels 21-40)
- **Single Target:** Meteor Strike
- **Groups:** Fire Dance
- **Healing:** Mantra
- **Fun:** Suplex everything!

### Late Game (Levels 41+)
- **Single Target:** Meteor Strike or Phantom Rush
- **Groups:** Air Blade
- **Healing:** Mantra
- **Boss Battles:** Meteor Strike

---

## Input Tips & Tricks

### General Input Rules

1. **Timing:** You have about 5 seconds to input the sequence
2. **Confirmation:** Press A/Confirm after completing the sequence
3. **Failure:** Wrong input wastes Sabin's turn with no effect
4. **Practice:** No penalty for practicing in low-risk battles

### Direction Input Basics

**Cardinal Directions (Easy):**
- ↑ = Up
- ↓ = Down
- ← = Left  
- → = Right

**Diagonal Directions (Harder):**
- ↗ = Up-Right
- ↘ = Down-Right
- ↙ = Down-Left
- ↖ = Up-Left

### Common Input Patterns

**Horizontal Line:**
- ← → or → ← 
- Used in: Raging Fist, Mantra

**Quarter-Circle:**
- ↓ ↙ ← or ← ↙ ↓ ↘ →
- Used in: Aura Cannon, Fire Dance

**Semi-Circle:**
- ← ↖ ↑ ↗ →
- Used in: Meteor Strike

**Three-Quarter Circle:**
- ↑ ↗ → ↘ ↓ ↙ ←
- Used in: Air Blade

**Box Pattern:**
- ← ↑ → ↓ or ↑ ↓ ← →
- Used in: Suplex, Phantom Rush

### Mastery Tips

1. **Start Slow:** Perfect the easy Blitzes first
2. **Smooth Motions:** Don't button mash - smooth is better
3. **Muscle Memory:** Repetition builds instinct
4. **No Rush:** You have time - don't panic
5. **Practice Mode:** Use low-risk battles to practice

---

## Damage Calculator

The plugin calculates estimated damage based on Sabin's current stats:

### Physical Blitzes
**Formula:** Power + (Level² × Vigor) / 256

**Affected Blitzes:**
- Raging Fist (Power: 110)
- Aura Cannon (Power: 130)
- Meteor Strike (Power: 150 × 4 hits)
- Suplex (Power: 120)
- Phantom Rush (Power: 255, max 9999 per hit × 4)

### Magical Blitzes
**Formula:** (Power × Level / 2) + Magic Power

**Affected Blitzes:**
- Fire Dance (Power: 106, all enemies)
- Air Blade (Power: 127, all enemies)

### Healing Blitz
**Formula:** (Max HP + Current HP) / 16

**Affected Blitz:**
- Mantra (heals all allies)

---

## Configuration Options

```lua
local config = {
    showDamageFormulas = true,      -- Show calculation formulas
    showInputDiagrams = true,       -- Show ASCII diagrams
    estimateDamage = true,          -- Calculate damage estimates
}
```

---

## Technical Details

### How It Works

The plugin:
1. Optionally loads Sabin's character data
2. Displays Blitz information from built-in database
3. Calculates damage estimates using current stats
4. Provides interactive reference interface

### Works Without Save File

Unlike tracker plugins, the Blitz Trainer works as a reference tool even without loading a save file. Damage calculations require Sabin's stats, but all other features work independently.

---

## Troubleshooting

### "Sabin not found"
**Not an error!** The plugin works as a reference guide without Sabin. Only damage calculations require his stats.

### Damage Estimates Seem Wrong
**Cause:** Estimates don't account for enemy defense or modifiers

**Solution:** These are base damage calculations. Actual battle damage varies.

---

## Use Cases

### For New Players
- Learn all Blitz inputs
- Understand the Blitz system
- Practice with easier Blitzes first

### For Veterans
- Quick reference during gameplay
- Damage optimization
- Refresh muscle memory

### For Speedrunners
- Quick lookup of essential Blitzes
- Damage calculations for routing
- Input practice reference

---

## Known Limitations

1. **Reference Only:** Cannot execute Blitzes from the plugin
2. **Damage Estimates:** Don't account for enemy defense
3. **No Practice Mode:** Can't actually practice inputs in plugin

---

## Future Enhancements

- Interactive practice mode (if API allows input simulation)
- Success rate tracking
- Video tutorials for complex inputs
- Boss-specific Blitz recommendations

---

## Version History

See [CHANGELOG.md](CHANGELOG.md)

### v1.0.0 (January 16, 2026)
- Initial release
- Complete 8-Blitz database
- Input diagrams and formulas
- Damage calculator
- Strategic recommendations

---

## Credits

**Plugin Development:** FF6 Editor Team  
**Blitz Data:** FF6 game mechanics research

---

## License

MIT License

---

**Master the Blitz!** 💪🔥
