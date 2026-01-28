# Character Ability Swap Plugin

⚠️ **EXPERIMENTAL PLUGIN** - Fundamentally alters gameplay mechanics

Give any character any special command ability. Mix and match up to 4 commands per character.

## Features

- **Command Ability Swapping**: Assign any of the 14 character commands to any character
- **14 Command Ability Library**: Complete database of all character special commands
- **Slot Customization**: Configure up to 4 command slots per character
- **Command Preview**: View descriptions before assigning
- **Save/Load Setups**: Store and recall custom ability configurations
- **Reset to Defaults**: Restore vanilla command setups
- **Preset Builds**: 4 pre-configured party setups
- **Export Functionality**: Export configurations to text files

## Command Ability Library

The plugin includes all 14 character-specific commands:

1. **Fight** - Basic attack (All characters)
2. **Magic** - Cast learned spells (Terra/Celes default)
3. **Morph** - Esper transformation (Terra only)
4. **Steal** - Steal items from enemies (Locke)
5. **Capture** - Capture for Colosseum (Locke)
6. **SwdTech** - Bushido techniques (Cyan)
7. **Throw** - Throw weapons/items (Shadow)
8. **Tools** - Mechanical tools (Edgar)
9. **Blitz** - Martial arts (Sabin)
10. **Runic** - Absorb magic as MP (Celes)
11. **Lore** - Cast enemy Lores (Strago)
12. **Sketch** - Mimic enemy (Relm)
13. **Control** - Control enemy (Relm)
14. **Slot** - Slot machine attacks (Setzer)
15. **Rage** - Enemy Rages (Gau)
16. **Leap** - Jump to Veldt (Gau)
17. **Mimic** - Repeat last action (Gogo)
18. **Dance** - Perform Dances (Mog)

## Preset Builds

### 1. All-Magic Party
Every character has the Magic command for full spell-casting party.

### 2. All-Physical Party
Focus on physical abilities: Blitz, Throw, SwdTech, Tools combinations.

### 3. Utility Focused
Emphasize Steal, Control, and support abilities across the party.

### 4. Chaos Mode
Random unusual combinations for experimental gameplay (Rage/Dance/Lore mixes).

## Usage

### Editing Character Abilities
1. Select "Edit Character Abilities"
2. Choose character to modify
3. Select command slot (1-4)
4. Pick new ability from library
5. Changes saved automatically

### Applying Presets
1. Select "Apply Preset Build"
2. Choose from 4 preset configurations
3. Entire party updated instantly

### Saving Custom Setups
1. Configure abilities to your preference
2. Select "Save Current Setup"
3. Enter descriptive name
4. Setup stored for future use

### Loading Saved Setups
1. Select "Load Saved Setup"
2. Choose from saved configurations
3. Setup applied immediately

## File Structure

```
plugins/character-ability-swap/
├── metadata.json          # Plugin metadata
├── plugin.lua             # Main plugin code (~800 LOC)
├── README.md              # This file
└── CHANGELOG.md           # Version history

ability_swap/
├── configurations.lua     # Active and saved setups
└── export_*.txt           # Exported configurations
```

## Configuration Storage

Configurations stored in `ability_swap/configurations.lua`:
```lua
{
  active_setup = {
    Terra = {1, 2, 3, 0},  -- Fight, Magic, Morph, (empty)
    Locke = {1, 4, 9, 0},  -- Fight, Steal, Blitz, (empty)
    -- ... all 14 characters
  },
  saved_setups = {
    {name = "My Setup", setup = {...}, saved = timestamp}
  }
}
```

## API Requirements

**IMPORTANT**: This plugin requires NEW APIs that may not be available:
- `character.get_command_abilities(char_id)` - Read current commands
- `character.set_command_ability(char_id, slot, ability_id)` - Write commands

Without these APIs, the plugin:
- ✅ Tracks configurations and setups
- ✅ Provides planning/documentation
- ✅ Exports configurations
- ❌ Cannot directly modify in-game abilities

## Permissions

- `read_save`: Read character data
- `write_save`: Write ability changes (if API available)
- `ui_display`: Display interface

## Known Limitations

- **API Dependency**: Command swapping requires advanced game state APIs
- Without API support, plugin acts as configuration planner only
- Some abilities may have character-specific implementations (Morph, Leap)
- Ability effectiveness depends on character stats and equipment
- Balance considerations: some combinations may be overpowered or ineffective

## Experimental Warning

This plugin fundamentally alters core gameplay mechanics. Effects include:

- **Game Balance**: Certain combinations may trivialize challenges
- **Character Identity**: Removes character uniqueness
- **Save Compatibility**: Modified saves may behave unexpectedly
- **Story Integration**: Some abilities tied to character-specific events

**Recommendation**: Test on backup saves before extensive use.

## Example Configurations

### Magic-Heavy Party
- Terra: Fight, Magic, Morph, Lore
- Celes: Fight, Magic, Runic, Lore
- Strago: Fight, Magic, Lore, Control
- Relm: Fight, Magic, Sketch, Control

### Physical Powerhouse
- Sabin: Fight, Blitz, SwdTech, Rage
- Edgar: Fight, Tools, Blitz, Throw
- Cyan: Fight, SwdTech, Blitz, Rage
- Gau: Fight, Rage, Blitz, Throw

### Utility Squad
- Locke: Fight, Steal, Control, Magic
- Shadow: Fight, Throw, Steal, Magic
- Setzer: Fight, Slot, Steal, Magic
- Gogo: Fight, Mimic, Steal, Magic

## Future Enhancements

- Visual ability icons
- Compatibility validation (warn about problematic combinations)
- Stat-based ability recommendations
- Ability synergy calculator
- Community preset sharing
- In-game ability swapping (with API support)
- Preset import/export format

## Version

**Version**: 1.0.0  
**Phase**: 6 (Gameplay-Altering Plugins)  
**Author**: FF6 Save Editor Plugin System  
**Priority**: HIGH-EXPERIMENTAL
