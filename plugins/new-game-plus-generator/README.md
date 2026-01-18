# New Game Plus Generator Plugin

Generate customizable New Game Plus save files with configurable carryover content and difficulty scaling.

## Overview

This plugin creates New Game Plus (NG+) configurations for FF6, allowing you to start fresh playthroughs while carrying over specific content. Configure what carries over (items, levels, spells), set enemy difficulty scaling, choose your starting world state, and generate NG+ profiles for endless replay value.

## Features

### Carryover Configuration
- **Items**: Carry over consumables, key items, or start fresh
- **Equipment**: Keep your weapons/armor or start with nothing
- **Levels**: Maintain character levels or reset to base
- **Spells**: Preserve learned magic or start clean
- **Espers**: Keep collected espers or discover them again

### Difficulty Scaling
- **Enemy Stat Multipliers**: Scale enemy HP/Attack/Defense from 1x to 10x
- **Preset Difficulty Tiers**: Easy, Normal, Hard, Nightmare
- **Custom Scaling**: Fine-tune difficulty to your preference
- **Balanced Progression**: Higher scaling requires strategic planning

### World State Selection
- **World of Balance (WoB)**: Start at the beginning
- **World of Ruin (WoR)**: Skip to the second half
- **Character Unlocking**: Optional "all characters unlocked" mode
- **Fresh Storylines**: Experience different narrative progressions

### Profile Management
- **Multiple Profiles**: Create unlimited NG+ configurations
- **Named Profiles**: Give each configuration a memorable name
- **Preset Templates**: Use Easy/Normal/Hard/Nightmare presets
- **Custom Profiles**: Configure every aspect individually
- **Profile Library**: Browse all created profiles

### Export & Sharing
- **Configuration Export**: Export NG+ profiles to text files
- **Timestamped Exports**: All exports include date/time stamps
- **Shareable Configs**: Share your NG+ setups with other players
- **Generation Logs**: Track what gets carried over

## Usage

### Creating a Custom NG+ Profile

1. Run the plugin
2. Select "Create Custom Profile"
3. Enter a profile name
4. Configure carryover options:
   - Items: Yes/No
   - Equipment: Yes/No
   - Levels: Yes/No
   - Spells: Yes/No
   - Espers: Yes/No
5. Set enemy stat scaling (1x-10x)
6. Choose starting world state (WoB/WoR)
7. Choose to unlock all characters or not
8. Profile is saved

### Creating from Preset

1. Run the plugin
2. Select "Create from Preset"
3. Choose preset:
   - **Easy**: Keep everything, 1x scaling
   - **Normal**: Keep levels + espers, 1.5x scaling
   - **Hard**: Keep espers only, 2x scaling
   - **Nightmare**: Nothing carries over, 3x scaling
4. Profile is created with timestamp

### Generating NG+

1. Run the plugin
2. Select "Generate NG+"
3. Choose a profile from your library
4. Configuration file is created
5. Follow manual steps to apply settings

## Difficulty Presets

### Easy Mode
**Philosophy**: Experience the story again with all your power

- **Carryover**: Everything (items, equipment, levels, spells, espers)
- **Scaling**: 1.0x (normal enemy stats)
- **Best For**: Story replays, casual runs, experimentation
- **Challenge**: Minimal - you'll be overpowered

### Normal Mode
**Philosophy**: Balanced replay with some challenge

- **Carryover**: Levels + Espers
- **Scaling**: 1.5x (enemies have 50% more stats)
- **Best For**: Standard NG+ experience
- **Challenge**: Moderate - enemies hit harder but you're still strong

### Hard Mode
**Philosophy**: Serious challenge for veterans

- **Carryover**: Espers only
- **Scaling**: 2.0x (double enemy stats)
- **Best For**: Experienced players wanting difficulty
- **Challenge**: High - strategic planning required

### Nightmare Mode
**Philosophy**: Brutal difficulty for masochists

- **Carryover**: Nothing
- **Scaling**: 3.0x (triple enemy stats)
- **Best For**: Ultimate challenge seekers
- **Challenge**: Extreme - requires perfect execution

## Technical Details

### Profile Data Structure

```lua
{
  name = "My NG+ Run",
  carryover = {
    items = true,
    equipment = false,
    levels = true,
    spells = false,
    espers = true
  },
  scaling = 2.0,
  world_state = "WoB",
  unlock_all = false,
  created = "2024-01-15 14:30:00"
}
```

### Configuration Storage

Profiles stored in `new_game_plus/profiles.lua`:
```lua
{
  profiles = {
    { name = "Profile 1", ... },
    { name = "Profile 2", ... }
  }
}
```

### Generation Process (Conceptual)

1. **Load Profile**: Read selected NG+ configuration
2. **Base Save Creation**: Create fresh game save at selected world state
3. **Carryover Application**: 
   - If `carryover.items`: Copy inventory from old save
   - If `carryover.equipment`: Copy equipment from old save
   - If `carryover.levels`: Copy character levels from old save
   - If `carryover.spells`: Copy learned spells from old save
   - If `carryover.espers`: Copy esper collection from old save
4. **Scaling Application**: Multiply all enemy stats by scaling factor
5. **Character Unlocking**: If `unlock_all`, unlock all party members
6. **Save Generation**: Write new save file

### Export Format

```
=== FF6 New Game Plus Profiles ===
Export Date: 2024-01-15 18:45:30

Profile 1: Nightmare Challenge
  Created: 2024-01-15 14:30:00
  Preset: Nightmare
  Enemy Scaling: 3x
  World State: WoB
  All Characters: No
  Carryover:
    items: No
    equipment: No
    levels: No
    spells: No
    espers: No
```

## API Requirements

**Current Status**: Configuration system works, actual save generation requires APIs not yet available.

### Working Features

1. **Profile Creation**: ✅ Configure and save NG+ profiles
2. **Profile Management**: ✅ View, browse, select profiles
3. **Export System**: ✅ Export profiles to text files
4. **Configuration Persistence**: ✅ Save/load via Lua serialization

### Conceptual Features

1. **Save File Creation**: ⚠️ Requires save file creation API
2. **World State Initialization**: ⚠️ Requires game state flag modification
3. **Enemy Stat Scaling**: ⚠️ Requires enemy data modification API
4. **Character Unlocking**: ⚠️ Requires party roster manipulation API
5. **Carryover Application**: ⚠️ Requires selective data copying between saves

### Required APIs (Not Yet Available)

1. `save.create(world_state)` - Create new save at specified game state
2. `save.copy_inventory(source, dest, types)` - Selective inventory transfer
3. `save.copy_character_data(source, dest, fields)` - Character data transfer
4. `enemy.modify_stats(multiplier)` - Global enemy stat scaling
5. `party.unlock_all_characters()` - Unlock full roster
6. `worldstate.set_flags(state)` - Set world progression flags

## Current Implementation

### What Works

- ✅ **Profile Creation**: Create unlimited custom NG+ profiles
- ✅ **Preset System**: Use difficulty presets for quick setup
- ✅ **Profile Library**: Browse all created profiles
- ✅ **Configuration Export**: Export profiles for backup/sharing
- ✅ **Configuration Files**: Generate .config files for manual application

### What's Conceptual

- ⚠️ **Automatic Save Generation**: Requires save creation APIs
- ⚠️ **Carryover Application**: Requires data transfer APIs
- ⚠️ **Enemy Scaling**: Requires enemy data modification
- ⚠️ **Character Unlocking**: Requires party roster APIs

### Manual Application

Until APIs are available, use generated .config files to manually:
1. Start a new game or load a save at desired world state
2. Use other plugins/editor features to apply carryover settings
3. Manually adjust enemy encounters (community ROM hack)
4. Unlock characters via save editing

## Examples

### Story Replay Build
```
Profile: "Easy Story Replay"
Carryover: All (items, equipment, levels, spells, espers)
Scaling: 1x
World: WoB
All Characters: No

Use Case: Re-experience the story without grinding
```

### Challenge Run Build
```
Profile: "Hard Mode Gauntlet"
Carryover: Espers only
Scaling: 2.5x
World: WoB
All Characters: No

Use Case: Balanced challenge - you have esper bonuses, enemies are tough
```

### Speedrun Practice Build
```
Profile: "Speedrun Practice"
Carryover: Levels only
Scaling: 1x
World: WoR
All Characters: Yes

Use Case: Practice WoR speedrun routes with appropriate levels
```

### Ultimate Challenge Build
```
Profile: "Nightmare Permadeath"
Carryover: Nothing
Scaling: 5x
World: WoB
All Characters: No

Use Case: Combine with Iron Man rules for extreme challenge
```

## FAQ

**Q: Can I modify a profile after creation?**  
A: Not currently - create a new profile with updated settings.

**Q: What's the highest recommended scaling for first NG+ run?**  
A: 1.5x-2x is comfortable for most players. 3x+ is expert territory.

**Q: If I carry over levels, do esper bonuses carry over too?**  
A: Conceptually yes - level carryover includes all permanent stat bonuses.

**Q: Can I start in WoR with all characters unlocked?**  
A: Yes! Perfect for testing endgame strategies or practicing specific fights.

**Q: Do preset profiles get special bonuses?**  
A: No, presets are just pre-configured settings. You can replicate them manually.

**Q: Why does generation create .config files instead of .sav files?**  
A: Save file creation APIs aren't available yet. Configs document your desired settings.

## Tips & Strategies

### First NG+ Run
- Start with Normal preset (1.5x scaling, levels + espers)
- Test the waters before going harder
- Export your profile before generating

### Maximum Difficulty
- Use Nightmare preset (3x scaling, no carryover)
- Combine with challenge plugins (Iron Man, NMB)
- WoB start for full experience

### Character Experimentation
- Use Easy preset with all characters unlocked
- WoR start to test endgame builds
- Perfect for testing party compositions

### Speedrun Practice
- Carry over levels only
- WoR start with all characters
- 1x scaling for consistency

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## License

Part of the FF6 Save Editor plugin system. See main LICENSE file.
