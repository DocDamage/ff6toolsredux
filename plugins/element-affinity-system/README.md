# Element Affinity System

**Version:** 1.0.0  
**Category:** Experimental / Combat Enhancement  
**Phase:** 6, Batch 4  
**Author:** FF6 Editor Plugin System

---

## Overview

The **Element Affinity System** plugin adds Pokemon-style elemental types to Final Fantasy VI characters. Each character can be assigned an element (Fire, Ice, Lightning, Water, Wind, Earth, Holy, Dark) that provides stat bonuses and creates strategic weakness/resistance patterns for tactical party composition.

### Key Features

- ✅ **8 Elemental Affinities** - Fire, Ice, Lightning, Water, Wind, Earth, Holy, Dark
- ✅ **Stat Bonuses** - Each element provides unique stat multipliers
- ✅ **Type Effectiveness** - Rock-paper-scissors style matchups
- ✅ **3 Preset Parties** - Balanced, Offensive, Defensive
- ✅ **Synergy Calculator** - Analyze character compatibility
- ✅ **Enemy Matchup Analysis** - Determine advantages/disadvantages
- ✅ **Party Composition Analyzer** - Check element distribution
- ✅ **Backup & Restore** - Undo changes easily
- ✅ **Operation Logging** - Track all modifications

⚠️ **WARNING:** This plugin **significantly alters character balance**. Stat modifications can create powerful or weak builds. Always backup your save file.

---

## Installation

1. Copy the `element-affinity-system` folder to your FF6 Editor `plugins/` directory
2. Restart FF6 Editor or reload plugins
3. The plugin will appear in the **Experimental** category

**Requirements:**
- FF6 Editor v3.4.0 or higher
- Valid FF6 save file
- Permissions: `read_save`, `write_save`, `ui_display`

---

## Elements & Stat Bonuses

### 🔥 Fire
**Stat Bonuses:**
- Magic: +15%
- Vigor: +10%
- Stamina: -5%

**Playstyle:** High offensive magic and physical damage, lower endurance.

---

### ❄️ Ice
**Stat Bonuses:**
- Magic: +15%
- Defense: +10%
- Speed: -5%

**Playstyle:** Defensive mage, slower but tankier.

---

### ⚡ Lightning
**Stat Bonuses:**
- Speed: +20%
- Magic: +10%
- Defense: -5%

**Playstyle:** Fast magic attacker, glass cannon.

---

### 💧 Water
**Stat Bonuses:**
- Magic: +10%
- Stamina: +15%
- Vigor: -5%

**Playstyle:** Sustained magic damage, high endurance.

---

### 🌪️ Wind
**Stat Bonuses:**
- Speed: +15%
- Evade: +20%
- Defense: -10%

**Playstyle:** Evasion tank, hit-and-run tactics.

---

### 🗻 Earth
**Stat Bonuses:**
- Defense: +20%
- Stamina: +15%
- Speed: -15%

**Playstyle:** Ultimate tank, very slow but nearly invulnerable.

---

### ✨ Holy
**Stat Bonuses:**
- Magic: +20%
- Stamina: +10%
- Vigor: -5%

**Playstyle:** Powerful healer/support, balanced offense.

---

### 🌑 Dark
**Stat Bonuses:**
- Vigor: +20%
- Magic: +5%
- Stamina: -10%

**Playstyle:** Physical damage dealer, berserker style.

---

## Type Effectiveness (Weaknesses & Resistances)

### Fire
- **Strong Against:** Ice, Wind
- **Weak Against:** Water, Earth

### Ice
- **Strong Against:** Water, Earth
- **Weak Against:** Fire, Lightning

### Lightning
- **Strong Against:** Water, Wind
- **Weak Against:** Earth

### Water
- **Strong Against:** Fire
- **Weak Against:** Lightning, Ice

### Wind
- **Strong Against:** Earth
- **Weak Against:** Fire, Lightning

### Earth
- **Strong Against:** Lightning, Fire
- **Weak Against:** Water, Wind

### Holy
- **Strong Against:** Dark
- **Weak Against:** Dark

### Dark
- **Strong Against:** Holy
- **Weak Against:** Holy

---

## Usage Guide

### Assign Element to Character

```lua
ElementAffinity.assignAffinity(0, "fire")  -- Terra = Fire
ElementAffinity.assignAffinity(1, "water")  -- Locke = Water
ElementAffinity.assignAffinity(2, "lightning")  -- Cyan = Lightning
```

**Output:**
```
✓ Character 0 assigned FIRE affinity
  MAGIC: 38 -> 44 (+15%)
  VIGOR: 32 -> 35 (+10%)
  STAMINA: 30 -> 29 (-5%)
```

---

### Remove Element from Character

```lua
ElementAffinity.removeAffinity(0)  -- Remove Terra's affinity
```

Restores original stats from backup.

---

### Preset Party Configurations

#### Balanced Party (8 Characters, All Elements)

```lua
ElementAffinity.applyBalancedParty()
```

Assigns:
- Character 0: Fire
- Character 1: Ice
- Character 2: Lightning
- Character 3: Water
- Character 4: Wind
- Character 5: Earth
- Character 6: Holy
- Character 7: Dark

---

#### Offensive Party (High Damage)

```lua
ElementAffinity.applyOffensiveParty()
```

Assigns Fire, Lightning, Dark, Holy to first 4 characters (damage-focused elements).

---

#### Defensive Party (High Survivability)

```lua
ElementAffinity.applyDefensiveParty()
```

Assigns Earth, Ice, Water, Holy to first 4 characters (defense-focused elements).

---

## Analysis Tools

### Party Composition Analysis

```lua
local analysis = ElementAffinity.analyzePartyComposition()
```

**Output:**
```
=== Party Composition Analysis ===
Characters with affinities: 8 / 14

Element Distribution:
  FIRE: 2
  ICE: 1
  LIGHTNING: 2
  WATER: 1
  EARTH: 1
  HOLY: 1

~ Moderately skewed composition
==================================
```

**Returns:**
```lua
{
    assigned_count = 8,
    element_counts = { fire = 2, ice = 1, ... },
    balance_score = 0.75,  -- 1.0 = perfect balance
}
```

---

### Character Synergy Calculator

Check compatibility between two characters:

```lua
local synergy, message = ElementAffinity.calculateSynergy(0, 1)
```

**Example Results:**
- `1.0` - Excellent synergy (covers weaknesses)
- `0.5` - Moderate synergy (same element)
- `0.0` - Neutral synergy (no interaction)
- `-0.5` - Poor synergy (elements oppose)

---

### Enemy Matchup Analysis

Analyze party effectiveness against enemy element:

```lua
ElementAffinity.analyzeEnemyMatchup("fire")
```

**Output:**
```
=== Matchup vs FIRE Enemy ===

Advantaged Characters: 2
  Character 1 (WATER) ✓
  Character 5 (EARTH) ✓

Disadvantaged Characters: 2
  Character 0 (FIRE) ✗
  Character 4 (WIND) ✗

Neutral Characters: 4
  ...
===============================
```

Use this to optimize party for specific bosses!

---

## Restoration

### Restore All Characters

```lua
ElementAffinity.restoreAll()
```

Removes all affinities and restores original stats from backup.

---

### Display Current Affinities

```lua
ElementAffinity.displayAffinities()
```

**Output:**
```
=== Character Affinities ===
Character 0: FIRE
Character 1: WATER
Character 2: LIGHTNING
Character 5: EARTH
============================
```

---

## Use Cases

### 1. Elemental-Themed Playthrough

**Goal:** Play with elemental type advantages.

```lua
ElementAffinity.applyBalancedParty()
-- Now plan battles around type matchups
ElementAffinity.analyzeEnemyMatchup("ice")  -- Check before boss
```

---

### 2. Mono-Element Challenge

**Goal:** Beat game with all characters as same element.

```lua
for charId = 0, 7 do
    ElementAffinity.assignAffinity(charId, "fire")
end
-- All fire party (strong offense, weak to water)
```

---

### 3. Synergy-Optimized Party

**Goal:** Maximize party synergy.

```lua
-- Assign complementary elements
ElementAffinity.assignAffinity(0, "fire")
ElementAffinity.assignAffinity(1, "water")  -- Covers fire's weakness
ElementAffinity.assignAffinity(2, "lightning")
ElementAffinity.assignAffinity(3, "earth")  -- Covers lightning's weakness

-- Check synergies
ElementAffinity.calculateSynergy(0, 1)  -- Fire + Water
ElementAffinity.calculateSynergy(2, 3)  -- Lightning + Earth
```

---

### 4. Boss-Specific Optimization

**Goal:** Optimize party for specific boss element.

```lua
-- Before Ifrit boss (Fire element)
ElementAffinity.analyzeEnemyMatchup("fire")
-- Use recommended Water/Earth characters
```

---

## Advanced Strategies

### Type Coverage Analysis

Create party with maximum type coverage:
- 4 characters with diverse elements
- Each covers another's weakness
- Example: Fire + Water + Lightning + Earth

### Stat Stacking

Maximize specific stats:
- All Holy (4 × +20% Magic = insane magic power)
- All Earth (4 × +20% Defense = wall of steel)
- All Lightning (4 × +20% Speed = strike first always)

### Balanced Approach

Mix offensive and defensive:
- 2 offensive (Fire, Lightning)
- 2 defensive (Earth, Water)
- Perfect balance of damage and survivability

---

## Troubleshooting

### Issue: Stat changes not visible

**Cause:** Stats may need save reload to update in-game.

**Solution:** Save and reload to see stat changes reflected.

---

### Issue: Character becomes too weak

**Cause:** Element has negative stat multipliers.

**Solution:** Remove affinity with `removeAffinity()` or choose different element.

---

### Issue: Cannot remove affinity

**Cause:** No backup available.

**Solution:** Backup is created automatically on first assignment. May need external save restoration.

---

## FAQ

### Q: Do elements affect in-game damage?
**A:** No, this plugin only modifies stats. Actual elemental damage calculation is unchanged.

### Q: Can I assign multiple elements to one character?
**A:** No, only one element per character. Reassigning removes previous element.

### Q: Do NPCs/enemies have elements?
**A:** This plugin doesn't modify enemies, but you can use matchup analysis for planning.

### Q: Which element is strongest?
**A:** Depends on playstyle. Holy/Dark have highest single-stat bonuses. Lightning is fastest. Earth is tankiest.

### Q: Can I create custom elements?
**A:** Not in v1.0, but you can modify CONFIG in plugin.lua to add custom elements.

---

## Known Limitations

1. **Stat Bonuses Only** - Doesn't affect actual elemental damage
2. **Single Element** - One element per character maximum
3. **Visual Indicators** - No in-game visual affinity indicators
4. **Enemy Elements** - Cannot assign elements to enemies
5. **Restoration Requires Backup** - Must have backup to restore stats

---

## Credits

**Developed by:** FF6 Editor Plugin System  
**Plugin Category:** Experimental / Combat Enhancement  
**Phase 6, Batch 4**

**Inspired by:** Pokemon type system

---

**Master the elements and dominate the battlefield! 🔥❄️⚡💧🌪️🗻✨🌑**
