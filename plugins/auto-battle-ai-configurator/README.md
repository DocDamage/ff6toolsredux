# Auto-Battle AI Configurator

**Version:** 1.0.0  
**Category:** Experimental / Strategy Planning  
**Phase:** 6, Batch 4  
**Author:** FF6 Editor Plugin System

---

## Overview

The **Auto-Battle AI Configurator** is a strategy planning and analysis tool that helps you design optimal battle strategies for your Final Fantasy VI characters. Create conditional action systems, priority-based decisions, and ability rotations to plan your battles effectively.

### Key Features

- ✅ **Conditional Action System** - IF-THEN rules for decision making
- ✅ **5 Priority Levels** - Critical, High, Medium, Low, Idle
- ✅ **Target Selection AI** - Choose optimal targets automatically
- ✅ **MP Conservation** - Manage MP usage intelligently
- ✅ **Item Usage Policies** - Control when/how items are used
- ✅ **Ability Rotations** - Plan optimal spell/skill sequences
- ✅ **4 Pre-Made Templates** - Healer, Attacker, Mage, Tank
- ✅ **Strategy Analysis** - Simulate and analyze effectiveness
- ✅ **Export Strategies** - Save plans as text files
- ✅ **Character Analysis** - Determine optimal roles

⚠️ **IMPORTANT:** This is a **PLANNING TOOL ONLY**. It does NOT modify in-game AI behavior (FF6 limitation). Use this to design strategies that you execute manually during gameplay.

---

## Installation

1. Copy the `auto-battle-ai-configurator` folder to your FF6 Editor `plugins/` directory
2. Restart FF6 Editor or reload plugins
3. The plugin will appear in the **Experimental** category

**Requirements:**
- FF6 Editor v3.4.0 or higher
- Valid FF6 save file (for character analysis)
- Permissions: `read_save`, `ui_display` (read-only plugin)

---

## Quick Start

### 1. Analyze Your Character

```lua
local analysis = AutoBattleAI.analyzeCharacter(0)  -- Terra
-- Returns: recommended_role = "mage" (high magic stat)
```

### 2. Create a Strategy from Template

```lua
local strategy_id = AutoBattleAI.createMageStrategy(0)
-- Creates pre-configured mage strategy for Terra
```

### 3. Customize Strategy

```lua
-- Add custom action
AutoBattleAI.addAction(strategy_id, AutoBattleAI.PRIORITY.HIGH, {
    type = "magic",
    name = "Ultima",
    target_type = "all_enemies",
})
```

### 4. Export for Reference

```lua
local export_text = AutoBattleAI.exportStrategy(strategy_id)
print(export_text)  -- Use during gameplay
```

---

## Strategy Creation

### Create Custom Strategy

```lua
local strategy_id = AutoBattleAI.createStrategy({
    name = "Custom Balanced",
    description = "Balanced offense and defense",
    character_id = 0,  -- Terra
    mp_conservation = true,
    mp_threshold = 40,  -- Conserve MP below 40%
    item_usage_enabled = true,
})
```

**Configuration Options:**
- `name` - Strategy display name
- `description` - Strategy purpose/notes
- `character_id` - Which character (0-13)
- `mp_conservation` - Enable MP management?
- `mp_threshold` - MP% threshold for conservation
- `item_usage_enabled` - Allow item usage?

---

## Priority System

Strategies use 5 priority levels for decision-making:

1. **CRITICAL** - Life-threatening situations (HP < 20%)
2. **HIGH** - Important actions (buffs, strong attacks)
3. **MEDIUM** - Normal actions (standard attacks/spells)
4. **LOW** - Low-priority actions (weak attacks, positioning)
5. **IDLE** - Default action when nothing else applies

### Adding Actions by Priority

```lua
-- CRITICAL: Emergency healing
AutoBattleAI.addAction(strategy_id, AutoBattleAI.PRIORITY.CRITICAL, {
    type = "magic",
    name = "Cure3",
    target_type = "lowest_hp_ally",
})

-- HIGH: Strong offensive magic
AutoBattleAI.addAction(strategy_id, AutoBattleAI.PRIORITY.HIGH, {
    type = "magic",
    name = "Ultima",
    target_type = "all_enemies",
})

-- MEDIUM: Standard attack
AutoBattleAI.addAction(strategy_id, AutoBattleAI.PRIORITY.MEDIUM, {
    type = "attack",
    name = "Attack",
    target_type = "enemy",
})
```

---

## Conditional Actions

Add IF-THEN rules to control when actions trigger:

```lua
-- IF ally HP < 30%, prioritize healing
AutoBattleAI.addCondition(strategy_id, {
    variable = "ally_hp_percent",
    operator = "<",
    value = 30,
})

-- IF self MP < 25%, use Ether
AutoBattleAI.addCondition(strategy_id, {
    variable = "self_mp_percent",
    operator = "<",
    value = 25,
})

-- IF enemy count > 3, use AoE
AutoBattleAI.addCondition(strategy_id, {
    variable = "enemy_count",
    operator = ">",
    value = 3,
})
```

**Available Variables:**
- `self_hp_percent` - Your HP percentage
- `self_mp_percent` - Your MP percentage
- `ally_hp_percent` - Lowest ally HP%
- `enemy_count` - Number of enemies
- `boss_battle` - Is this a boss fight? (boolean)

**Available Operators:**
- `==` - Equal to
- `~=` - Not equal to
- `>` - Greater than
- `<` - Less than
- `>=` - Greater than or equal
- `<=` - Less than or equal

---

## Target Types

Specify who actions target:

- `self` - Target yourself
- `ally` - Target random ally
- `enemy` - Target random enemy
- `all_allies` - Target all allies
- `all_enemies` - Target all enemies
- `lowest_hp_ally` - Ally with lowest HP
- `highest_hp_enemy` - Enemy with highest HP
- `random` - Completely random target

---

## Pre-Made Strategy Templates

### Healer Template

```lua
local healer_id = AutoBattleAI.createHealerStrategy(character_id)
```

**Healer Features:**
- CRITICAL: Cure3 on ally < 20% HP
- HIGH: Cure2 on ally < 50% HP
- MEDIUM: Protect/Shell buffs
- LOW: Offensive magic (Fire2)
- IDLE: Physical attack

---

### Physical Attacker Template

```lua
local attacker_id = AutoBattleAI.createAttackerStrategy(character_id)
```

**Attacker Features:**
- CRITICAL: X-Potion on self < 30% HP
- HIGH: Jump on strongest enemy
- MEDIUM: Physical attack
- LOW: Switch to back row if hurt

---

### Offensive Mage Template

```lua
local mage_id = AutoBattleAI.createMageStrategy(character_id)
```

**Mage Features:**
- CRITICAL: Ether on self < 25% MP
- HIGH: Ultima on all enemies
- MEDIUM: Fire3/Ice3/Bolt3 rotation
- LOW: Physical attack (MP conservation)

---

### Tank/Defender Template

```lua
local tank_id = AutoBattleAI.createTankStrategy(character_id)
```

**Tank Features:**
- CRITICAL: Defend when < 25% HP
- HIGH: Provoke enemies
- MEDIUM: Counter attack
- LOW: Hi-Potion self-heal

---

## Strategy Analysis

Simulate strategy effectiveness:

```lua
local analysis = AutoBattleAI.analyzeStrategy(strategy_id, {
    self_hp_percent = 50,
    ally_hp_percent = 80,
    enemy_count = 3,
})
```

**Analysis Returns:**
```lua
{
    simulated_turns = 100,
    action_distribution = {
        CRITICAL = 5,
        HIGH = 20,
        MEDIUM = 40,
        LOW = 25,
        IDLE = 10,
    },
    estimated_dps = 10000,  -- Damage per turn
    estimated_hps = 1500,   -- Healing per turn
    mp_efficiency = 80,     -- MP usage efficiency %
}
```

---

## Strategy Management

### List All Strategies

```lua
local strategies = AutoBattleAI.listStrategies()
for _, strat in ipairs(strategies) do
    print(strat.name .. " - " .. strat.description)
end
```

### Get Strategy Details

```lua
local details = AutoBattleAI.getStrategy(strategy_id)
-- Returns full strategy configuration
```

### Clone Strategy

```lua
local new_id = AutoBattleAI.cloneStrategy(strategy_id, "Modified Healer")
-- Creates copy you can customize
```

### Delete Strategy

```lua
AutoBattleAI.deleteStrategy(strategy_id)
```

---

## Character Analysis

Analyze character stats to determine optimal role:

```lua
local analysis = AutoBattleAI.analyzeCharacter(0)  -- Terra
```

**Returns:**
```lua
{
    character_id = 0,
    level = 35,
    stats = {
        vigor = 32,
        magic = 48,  -- High magic!
        speed = 34,
        hp = 2450,
        mp = 420,
    },
    recommended_role = "mage",  -- Determined by stats
}
```

**Role Determination:**
- **Mage:** Magic > Vigor + 20
- **Physical:** Vigor > Magic + 20
- **Tank:** HP > 5000
- **Healer:** MP > 500
- **Balanced:** Otherwise

---

## Export & Import

### Export Single Strategy

```lua
local export_text = AutoBattleAI.exportStrategy(strategy_id)
-- Save to file or print for reference
```

**Export Format:**
```
==============================================
AI STRATEGY: Healer
==============================================
Description: Dedicated healing and support role
Character ID: 5
MP Conservation: true (Threshold: 30%)
Item Usage: true

--- CONDITIONS ---
1. IF ally_hp_percent < 20

--- ACTIONS BY PRIORITY ---

[CRITICAL PRIORITY]
  1. magic: Cure3 (Target: lowest_hp_ally)

[HIGH PRIORITY]
  1. magic: Cure2 (Target: lowest_hp_ally)
...
==============================================
```

### Export All Strategies

```lua
local all_exports = AutoBattleAI.exportAllStrategies()
-- Exports all strategies in one document
```

---

## Use Cases

### 1. Party Strategy Planning

**Goal:** Plan optimal strategies for 4-person party.

```lua
-- Healer
local healer_id = AutoBattleAI.createHealerStrategy(5)  -- Shadow
AutoBattleAI.exportStrategy(healer_id)

-- Tank
local tank_id = AutoBattleAI.createTankStrategy(2)  -- Sabin
AutoBattleAI.exportStrategy(tank_id)

-- Mage
local mage_id = AutoBattleAI.createMageStrategy(0)  -- Terra
AutoBattleAI.exportStrategy(mage_id)

-- Attacker
local attacker_id = AutoBattleAI.createAttackerStrategy(1)  -- Locke
AutoBattleAI.exportStrategy(attacker_id)

-- Export all for reference during gameplay
print(AutoBattleAI.exportAllStrategies())
```

---

### 2. Boss Strategy Design

**Goal:** Create boss-specific strategies.

```lua
local boss_strategy = AutoBattleAI.createStrategy({
    name = "Kefka Tower Strategy",
    description = "Optimized for final boss",
    character_id = 0,
    mp_conservation = false,  -- No MP conservation needed
})

-- Add boss-specific conditions
AutoBattleAI.addCondition(boss_strategy, {
    variable = "boss_battle",
    operator = "==",
    value = true,
})

-- High-damage ultimates
AutoBattleAI.addAction(boss_strategy, AutoBattleAI.PRIORITY.HIGH, {
    type = "magic",
    name = "Ultima",
    target_type = "all_enemies",
})
```

---

### 3. Challenge Run Planning

**Goal:** Plan strategies for no-item runs.

```lua
local no_item_strategy = AutoBattleAI.createStrategy({
    name = "No-Item Healer",
    character_id = 5,
    item_usage_enabled = false,  -- Disable items
    mp_conservation = true,
})

-- Only use spells for healing
AutoBattleAI.addAction(no_item_strategy, AutoBattleAI.PRIORITY.CRITICAL, {
    type = "magic",
    name = "Cure3",
    target_type = "lowest_hp_ally",
})
```

---

## Troubleshooting

### Issue: "Strategy not found"

**Cause:** Using wrong strategy ID.

**Solution:** Use `listStrategies()` to find correct IDs.

---

### Issue: Cannot add more actions

**Cause:** Maximum actions per priority reached (10).

**Solution:** Organize actions better or remove unnecessary ones.

---

### Issue: Condition not working

**Cause:** Invalid variable or operator.

**Solution:** Check spelling of variables and operators. Use provided constants.

---

## FAQ

### Q: Does this control in-game AI?
**A:** No. This is a planning tool. You manually execute strategies during battles.

### Q: Can I test strategies in battle?
**A:** Not automatically. Export strategies and follow them manually during gameplay.

### Q: How do I share strategies?
**A:** Export to text and share the output file.

### Q: Can I import strategies?
**A:** Not currently. Manual recreation required from exported text.

### Q: Does this work with challenge modes?
**A:** Yes! Design strategies compatible with challenge restrictions.

---

## Known Limitations

1. **Planning Only:** Does not control actual game AI
2. **Manual Execution:** You must follow strategies manually
3. **No Battle Integration:** Cannot auto-execute during battles
4. **Simulation Basic:** Analysis is estimate-based, not exact
5. **No Import:** Exported strategies cannot be re-imported (v1.0)

---

## Credits

**Developed by:** FF6 Editor Plugin System  
**Plugin Category:** Experimental / Strategy Planning  
**Phase 6, Batch 4**

---

**Plan your battles, execute with precision! ⚔️🎯**
