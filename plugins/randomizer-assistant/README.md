# Randomizer Assistant Plugin

Support for randomized game seeds with spoiler log import, location tracking, and logic validation for Final Fantasy VI randomizers.

## Features

- **Spoiler Log Import**: Parse JSON and CSV spoiler logs
- **Item Location Tracker**: Check off locations as you find items
- **Character Location Tracker**: Track where characters are recruited
- **Boss Shuffle Tracking**: Monitor boss replacements and defeats
- **Ability Randomization Display**: View randomized character abilities
- **Logic Validation**: Check seed beatability and progression
- **Accessibility Checker**: Verify reachable locations
- **Progression Tracking**: Track overall completion percentage
- **Seed Metadata Display**: View seed settings and information
- **Statistics Export**: Export progress reports

## Supported Formats

### JSON Spoiler Logs
```json
{
  "items": {
    "Narshe Weapon Shop": "Elixir",
    "Figaro Castle Throne": "Ragnarok"
  },
  "characters": {
    "Phantom Train": "Sabin",
    "Opera House": "Locke"
  },
  "bosses": {
    "Narshe": "Ultima Weapon",
    "Mt. Kolts": "Kefka"
  },
  "abilities": {
    "Terra": "Blitz",
    "Sabin": "Magic"
  }
}
```

### CSV Spoiler Logs
```csv
Category,Location,Item/Entity
item,Narshe Weapon Shop,Elixir
item,Figaro Castle Throne,Ragnarok
character,Phantom Train,Sabin
boss,Narshe,Ultima Weapon
ability,Terra,Blitz
```

## Usage

### Importing Spoiler Log
1. Place spoiler log in `randomizer/spoiler_logs/` directory
2. Select "Import Spoiler Log"
3. Enter filename (e.g., `seed_12345.json`)
4. Plugin parses and loads seed data

### Tracking Progress
- **View Item Locations**: See all item placements
- **Mark Location Checked**: Check off locations as you explore
- **View Characters**: Track character recruitment progress
- **View Boss Shuffle**: Monitor which bosses are where
- **View Abilities**: See randomized command abilities

### Logic Checking
- Validates presence of key progression items
- Checks seed beatability
- Identifies potential softlocks
- Simplified logic analysis

## Progression Tracking

The plugin automatically tracks:
- **Locations Checked**: Items obtained from locations
- **Characters Recruited**: Party member recruitment
- **Bosses Defeated**: Boss encounters completed
- **Overall Completion**: Percentage-based progress

## File Structure

```
plugins/randomizer-assistant/
├── metadata.json          # Plugin metadata
├── plugin.lua             # Main plugin code
├── README.md              # This file
└── CHANGELOG.md           # Version history

randomizer/
├── spoiler_logs/          # Place spoiler logs here
│   ├── seed_12345.json
│   └── seed_67890.csv
├── seed_data.lua          # Loaded seed information
├── progress.lua           # Progression tracking
└── stats_*.txt            # Exported statistics
```

## Item Location Tracker

Check off locations as you find items:
- [✓] Checked locations (green checkmark)
- [ ] Unchecked locations (empty box)
- Progress percentage display
- Sorted by location name

## Character Location Tracker

Track character recruitment:
- View where each character is located in this seed
- Check off recruited characters
- Recruitment progress percentage

## Boss Shuffle Tracker

Monitor boss replacements:
- See which boss is at each location
- Track defeated bosses
- Boss completion percentage

## Logic Validation

Simplified beatability check:
- Validates key item presence (Airship, Ragnarok, etc.)
- Checks progression item accessibility
- Warns if critical items missing
- Basic softlock detection

## API Requirements

- `ReadFile()`, `WriteFile()`: Spoiler log import and progress storage

## Permissions

- `read_save`: Read game state for progression checking (future)
- `ui_display`: Display tracker interface
- `file_io`: Read spoiler logs, persist progress

## Notes

- Spoiler logs must be placed in `randomizer/spoiler_logs/` directory
- JSON parser is simplified; complex nested structures may not parse
- CSV format expects specific column order (Category, Location, Item/Entity)
- Logic checking is simplified and may not catch all softlocks
- Progression tracking is manual (mark locations as checked)

## Randomizer Compatibility

Designed for FF6 randomizers that generate spoiler logs:
- Beyond Chaos
- Worlds Collide
- Custom randomizers with JSON/CSV export

## Export Statistics

Export includes:
- Seed filename and import date
- Progress summary (locations/characters/bosses)
- Completion percentages
- Timestamp

## Future Enhancements

- Auto-detection of obtained items from save file
- Advanced logic solver with progression chains
- Hint system integration
- Community seed sharing
- Seed generation support
- Real-time accessibility calculation
- Visual map with location highlights

## Version

**Version**: 1.0.0  
**Phase**: 5 (Challenge & Advanced Tools)  
**Author**: FF6 Save Editor Plugin System
