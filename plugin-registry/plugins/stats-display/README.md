# Character Stats Display

A comprehensive character information viewer for the FF6 Save Editor that displays all character stats in one convenient, formatted view.

## Features

- **Complete Character Overview:** View all character information in one dialog
- **Basic Stats:** Level, Experience, Active/Inactive status
- **HP/MP Display:** Current and maximum values with percentage calculations
- **Combat Stats:** Vigor, Stamina, Speed, and Magic attributes
- **Equipment Slots:** View all 6 equipment slots (Weapon, Shield, Armor, Helmet, 2 Relics)
- **Learned Spells:** List all learned spells with proficiency levels
- **Commands:** Display assigned command abilities
- **Character Selection:** Easy character picker with status indicators
- **Configurable Display:** Customize which sections are shown

## Installation

1. Open FF6 Save Editor (version 3.4.0 or higher)
2. Go to `Tools → Marketplace` or `Community → Marketplace`
3. Search for "Character Stats Display"
4. Click "Install"

## Usage

### Basic Usage

1. Load a save file in the FF6 Save Editor
2. Go to `Tools → Plugin Manager`
3. Find "Character Stats Display" in the installed plugins list
4. Click "Run"
5. Select a character from the list
6. View the comprehensive stats display

### Character Selection

The plugin will present a list of all characters in your save file with their status:
```
1. Terra (Active)
2. Locke (Active)
3. Edgar (Inactive)
...
```

Enter the number corresponding to the character you want to view.

### Display Format

The plugin shows information in organized sections:

```
============================================================
CHARACTER: Terra
============================================================

[BASIC STATS]
  Status: Active
  Level: 50
  Experience: 123456

[HP / MP]
  HP: 850/1000 (85.0%)
  MP: 180/200 (90.0%)

[COMBAT STATS]
  Vigor (Physical Attack): 35
  Stamina (Physical Defense): 30
  Speed (Agility): 40
  Magic (Magical Power): 45

[EQUIPMENT]
  Weapon: Item #100
  Shield: Item #101
  Armor: Item #102
  Helmet: Item #103
  Relic 1: Item #104
  Relic 2: Item #105

[LEARNED SPELLS] (25 total)
  Spell #15 (Proficiency: 100%)
  Spell #8 (Proficiency: 95%)
  Spell #22 (Proficiency: 80%)
  ... and 22 more

[COMMANDS] (4 assigned)
  1. Fight
  2. Magic
  3. Morph
  4. Item

============================================================
```

## Configuration

The plugin supports several configuration options (modify at top of plugin.lua):

```lua
local config = {
    showEquipment = true,           -- Show equipment section
    showSpells = true,              -- Show learned spells
    showCommands = true,            -- Show command assignments
    maxSpellsDisplay = 20,          -- Maximum spells to display
    sortSpellsByProficiency = true  -- Sort spells by proficiency level
}
```

### Configuration Options

- **showEquipment:** Set to `false` to hide the equipment section
- **showSpells:** Set to `false` to hide the learned spells section
- **showCommands:** Set to `false` to hide the commands section
- **maxSpellsDisplay:** Maximum number of spells to show (default: 20)
- **sortSpellsByProficiency:** If `true`, spells are sorted by proficiency (highest first); if `false`, sorted by spell ID

## Permissions

This plugin requires the following permissions:

- **read_save:** Required to read character data from the save file
- **ui_display:** Required to show the character selection dialog and stats display

**Note:** This plugin does NOT modify save data. It only reads and displays information.

## Use Cases

- **Quick Character Review:** Get a complete overview without navigating multiple tabs
- **Build Analysis:** Review character builds and equipment setups
- **Spell Coverage Check:** See which spells each character knows
- **Party Composition:** Compare characters when deciding party makeup
- **Speedrun Planning:** Quickly assess character capabilities

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered
- **Platform:** Windows, macOS, Linux (wherever the editor runs)

## Screenshots

![Stats Display Example](screenshot.png)
*Example of character stats display showing Terra's complete information*

## Known Limitations

- Equipment names show as "Item #XXX" (future versions may include item name lookup)
- Spell names show as "Spell #XX" (future versions may include spell name lookup)
- Maximum of 20 spells displayed by default (configurable)

## Troubleshooting

### "No characters found in save file"
**Cause:** Save file is not loaded or contains no character data  
**Solution:** Ensure a valid save file is loaded before running the plugin

### "Permission Error"
**Cause:** Plugin permissions not granted  
**Solution:** Reinstall the plugin through the Marketplace

### Character selection doesn't work
**Cause:** Invalid number entered  
**Solution:** Enter a number between 1 and the total number of characters shown

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Future Enhancements

Planned features for future versions:
- Item/equipment name lookups
- Spell name lookups
- Character comparison mode
- Export stats to clipboard/file
- Customizable display format
- Quick-access keyboard shortcut

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
