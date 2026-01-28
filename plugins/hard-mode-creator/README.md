# Hard Mode Creator Plugin

Create custom difficulty modifiers and challenge profiles for extreme FF6 playthroughs.

## Overview

This plugin creates and manages difficulty profiles that modify enemy stats, player stats, experience/gil rates, and enforce gameplay restrictions. Perfect for challenge runs, speedrun practice, or just making the game harder. Build your own extreme difficulty or use presets like Iron Man and Nuzlocke modes.

## Features

### Enemy Stat Multipliers
- **Individual Stat Control**: Adjust HP, MP, Attack, Defense, Magic, Magic Defense independently
- **0.5x-10x Range**: 20 steps of scaling per stat (0.5x, 1x, 1.5x, 2x ... 10x)
- **Flexible Difficulty**: Make enemies tankier without making them hit harder, or vice versa
- **Boss Scaling**: Affects all enemies globally (bosses become true challenges)

### Player Stat Handicaps
- **Stat Divisors**: Reduce player HP, MP, Attack, Defense, Magic
- **Handicap Levels**: 1x (no handicap), 1.5x, 2x, 3x divisors
- **Glass Cannon Builds**: High damage, low survivability
- **Challenge Modes**: Combine with high enemy stats for extreme difficulty

### Experience & Gil Rate Modifiers
- **Experience Rate**: 0.5x-10x multiplier (0.5x = half XP, 10x = 10x XP)
- **Gil Rate**: 0.5x-10x multiplier independently from XP
- **Low Level Runs**: 0.1x XP for extreme challenge
- **Poverty Runs**: 0.1x Gil forces item management
- **Fast Leveling**: 5x-10x XP for testing or quick replays

### Difficulty Presets
- **Iron Man**: Permanent death mechanics + high difficulty
- **Nuzlocke**: First encounter rules + permadeath
- **Low Level**: Minimal XP gain, extreme enemy stats
- **Poverty**: Minimal gil, item restrictions
- **Nightmare**: Maximum difficulty across all stats

### Restriction System
- **Item Use Restrictions**: No Phoenix Downs, No Life spells, etc.
- **Equipment Slot Restrictions**: Limit available equipment slots
- **Shop Restrictions**: Disable item shops for poverty runs
- **Save Point Restrictions**: Limit save opportunities
- **Custom Rules**: Document any house rules for your challenge

### Profile Management
- **Multiple Profiles**: Create unlimited difficulty configurations
- **Named Profiles**: Give each profile a memorable name
- **Active Profile Tracking**: See which profile is currently applied
- **Profile Library**: Browse all created profiles
- **Export Profiles**: Share your difficulty settings

## Usage

### Creating a Custom Profile

1. Run the plugin
2. Select "Create Custom Profile"
3. Enter a profile name
4. Configure enemy stat multipliers (6 stats)
5. Configure player stat divisors (5 stats)
6. Set experience rate multiplier
7. Set gil rate multiplier
8. Profile is saved

### Creating from Preset

1. Run the plugin
2. Select "Create from Preset"
3. Choose preset:
   - Iron Man
   - Nuzlocke
   - Low Level
   - Poverty
   - Nightmare
4. Preset profile is created with timestamp

### Activating a Profile

1. Run the plugin
2. Select "Activate Profile"
3. Choose profile from library
4. Profile becomes active (configuration applied where possible)

### Viewing Profiles

Browse all profiles with:
- Profile name and creation date
- Active status indicator
- Preset source (if from preset)
- Scaling values

## Difficulty Presets

### Iron Man Mode
**Philosophy**: Permanent death + high stakes

- **Enemy Multipliers**: HP 1.5x, Attack 1.3x, Defense 1.2x
- **Player Divisors**: No handicap (1x)
- **Rates**: 0.5x XP, 0.5x Gil
- **Restrictions**: No Phoenix Down, No Life spells
- **Challenge**: Character death is permanent; plan carefully

### Nuzlocke Mode
**Philosophy**: First encounter only + death rules

- **Enemy Multipliers**: HP 1.2x, Attack 1.1x
- **Player Divisors**: No handicap (1x)
- **Rates**: 0.75x XP, 1x Gil
- **Restrictions**: Permadeath, First encounter recruitment
- **Challenge**: Can only recruit first character encountered in each area

### Low Level Mode
**Philosophy**: Minimal leveling, maximum challenge

- **Enemy Multipliers**: HP 2x, Attack 1.5x, Defense 1.3x, Magic 1.4x
- **Player Divisors**: No handicap (1x)
- **Rates**: 0.1x XP, 1x Gil
- **Restrictions**: None (low levels are challenge enough)
- **Challenge**: Stay underleveled throughout entire game

### Poverty Mode
**Philosophy**: Minimal resources, item management hell

- **Enemy Multipliers**: HP 1.3x, Attack 1.2x
- **Player Divisors**: No handicap (1x)
- **Rates**: 1x XP, 0.1x Gil
- **Restrictions**: No item shops
- **Challenge**: Can barely afford anything; treasure hunting essential

### Nightmare Mode
**Philosophy**: Everything is worse, all the time

- **Enemy Multipliers**: HP 3x, Attack 2x, Defense 1.8x, Magic 2x, MagicDef 1.5x
- **Player Divisors**: HP 2x, Attack 2x, Magic 2x (half stats)
- **Rates**: 0.5x XP, 0.5x Gil
- **Restrictions**: No Save Points, Limited items
- **Challenge**: Only for the truly masochistic

## Technical Details

### Profile Data Structure

```lua
{
  name = "My Hard Mode",
  enemy_mult = {
    HP = 2.0,
    MP = 1.5,
    Attack = 1.8,
    Defense = 1.3,
    Magic = 1.5,
    MagicDef = 1.2
  },
  player_div = {
    HP = 1.0,    -- 1.0 = no handicap
    MP = 1.0,
    Attack = 1.5, -- 1.5 = player does 66% damage
    Defense = 2.0, -- 2.0 = player takes 2x damage
    Magic = 1.0
  },
  exp_rate = 0.5,  -- Half experience
  gil_rate = 0.5,  -- Half gil
  restrictions = {"No Phoenix Down", "No Save Points"},
  created = "2024-01-15 14:30:00"
}
```

### Stat Modifier Mechanics

**Enemy Multipliers** (> 1.0 makes enemies stronger):
- HP 2x: Enemies have double health
- Attack 1.5x: Enemies hit 50% harder
- Defense 1.3x: Enemies take 30% less physical damage
- Magic 1.5x: Enemy spells 50% stronger
- MagicDef 1.2x: Enemies resist magic 20% better

**Player Divisors** (> 1.0 makes players weaker):
- HP 2x divisor: Player has half health
- Attack 1.5x divisor: Player does 66% damage
- Defense 2x divisor: Player takes double damage

### Configuration Storage

Profiles stored in `hard_mode/profiles.lua`:
```lua
{
  profiles = {
    { name = "Profile 1", ... },
    { name = "Profile 2", ... }
  },
  active = 1  -- Index of active profile
}
```

### Export Format

```
=== FF6 Hard Mode Difficulty Profiles ===
Export Date: 2024-01-15 18:45:30

Profile 1: Nightmare Challenge
  Created: 2024-01-15 14:30:00
  Preset: Nightmare
  Enemy Multipliers:
    HP: 3.0x
    Attack: 2.0x
    Defense: 1.8x
    Magic: 2.0x
    MagicDef: 1.5x
  Player Divisors:
    HP: 0.5x
    Attack: 0.5x
    Magic: 0.5x
  Experience Rate: 0.5x
  Gil Rate: 0.5x
  Restrictions: No Save Points, Limited items
```

## API Requirements

**Current Status**: Configuration system works, actual stat modification requires APIs not yet available.

### Working Features

1. **Profile Creation**: ✅ Create and save difficulty profiles
2. **Profile Management**: ✅ View, browse, activate profiles
3. **Export System**: ✅ Export profiles to text files
4. **Configuration Persistence**: ✅ Save/load via Lua serialization

### Conceptual Features

1. **Enemy Stat Modification**: ⚠️ Requires enemy data modification API
2. **Player Stat Modification**: ⚠️ Requires character stat override API
3. **XP/Gil Rate Modification**: ⚠️ Requires battle reward modifier API
4. **Restriction Enforcement**: ⚠️ Requires game flag/hook system
5. **Save Point Blocking**: ⚠️ Requires event system hooks

### Required APIs (Not Yet Available)

1. `enemy.modify_stats(multipliers)` - Apply scaling to enemy stats
2. `character.apply_handicap(divisors)` - Reduce player stats
3. `battle.set_reward_rates(exp_mult, gil_mult)` - Modify battle rewards
4. `game.set_restrictions(rules)` - Enforce gameplay restrictions
5. `save.block_save_points(enabled)` - Disable save functionality

## Current Implementation

### What Works

- ✅ **Profile Creation**: Create unlimited custom difficulty profiles
- ✅ **Preset System**: Use 5 difficulty presets for quick setup
- ✅ **Profile Library**: Browse all created profiles
- ✅ **Profile Activation**: Mark profile as active
- ✅ **Configuration Export**: Export profiles for backup/sharing

### What's Conceptual

- ⚠️ **Stat Modification**: Requires enemy/player stat APIs
- ⚠️ **Rate Modification**: Requires battle reward APIs
- ⚠️ **Restriction Enforcement**: Requires game hook system
- ⚠️ **Automatic Application**: Profiles tracked but not auto-applied

### Manual Application

Until APIs are available:
1. Use profile as documentation for house rules
2. Manually track restrictions (honor system)
3. Use ROM hacks for stat scaling (community tools)
4. Use save editor for player stat reduction

## Examples

### Balanced Challenge
```
Profile: "Hard But Fair"
Enemy: HP 1.5x, Attack 1.3x, Defense 1.2x
Player: No handicap
XP: 1x, Gil: 1x
Restrictions: None

Use Case: Enemies are tougher, but you're not crippled
```

### Glass Cannon Build
```
Profile: "Maximum Damage"
Enemy: HP 1x, Attack 1x, Defense 1x
Player: HP 2x handicap, Attack no handicap
XP: 1x, Gil: 1x
Restrictions: None

Use Case: You hit hard, but one mistake kills you
```

### Extreme Low Level
```
Profile: "LLG Challenge"
Enemy: HP 3x, Attack 2x, Defense 2x, Magic 2x
Player: No handicap
XP: 0.1x, Gil: 1x
Restrictions: None

Use Case: Classic FF6 Low Level Game challenge
```

### Speedrun Practice
```
Profile: "Fast Leveling"
Enemy: HP 1x, Attack 1x
Player: No handicap
XP: 10x, Gil: 10x
Restrictions: None

Use Case: Level up fast to practice lategame strats
```

## FAQ

**Q: Can I modify a profile after creation?**  
A: Not currently - create a new profile with updated settings.

**Q: What's the highest enemy multiplier used in community challenges?**  
A: Most stop at 3x-5x. Beyond that becomes nearly unplayable.

**Q: Do player divisors affect magic damage?**  
A: Yes - Magic divisor affects spell damage output.

**Q: Can I have separate XP and Gil rates?**  
A: Yes! Configure them independently (e.g., 0.1x XP, 10x Gil).

**Q: How do I add custom restrictions?**  
A: Currently restrictions are in preset profiles. Custom restriction entry coming in future version.

**Q: Will activating a profile immediately change my save?**  
A: No - activation marks it as active but doesn't modify save (requires APIs).

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## License

Part of the FF6 Save Editor plugin system. See main LICENSE file.
