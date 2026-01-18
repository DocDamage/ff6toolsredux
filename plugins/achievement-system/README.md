# Achievement System Plugin

Track and unlock achievements across combat, collection, challenge, and secret categories for Final Fantasy VI.

## Features

- **50+ Predefined Achievements**: Comprehensive achievement library across 4 categories
- **Custom Achievement Creation**: Create your own achievement definitions
- **Progress Tracking**: Automatically track achievement progress with timestamps
- **Achievement Categories**:
  - **Combat** (15 achievements): Battle-related feats and mastery
  - **Collection** (15 achievements): Item, character, and bestiary completion
  - **Challenge** (12 achievements): Difficult gameplay challenges and restrictions
  - **Secret** (8 achievements): Hidden achievements revealed upon unlock
- **Points System**: Earn points for each achievement (10-150 points)
- **Retroactive Checking**: Scan your save file to unlock eligible achievements
- **Statistics Dashboard**: View progress by category and recent unlocks
- **Export Functionality**: Export achievement data to text file
- **Manual Unlock**: Unlock custom achievements manually

## Usage

1. Launch the plugin from the plugin manager
2. Use "Check Achievements" to retroactively scan your save file
3. View achievements by category or see all at once
4. Create custom achievements for personal challenges
5. Export achievement data for sharing or backup

## Achievement Examples

### Combat
- **Magic Mastery** (50 pts): Learn all spells with one character
- **Level 99 Club** (50 pts): Reach level 99 with any character
- **Rage Master** (60 pts): Learn all Rages with Gau

### Collection
- **Millionaire** (50 pts): Accumulate 999,999 GP
- **Esper Collector** (60 pts): Obtain all espers
- **Full Party** (40 pts): Recruit all 14 characters

### Challenge
- **Natural Magic Block** (120 pts): Complete game without learning magic
- **Solo Run** (150 pts): Complete game with only one character
- **Low Level Legend** (100 pts): Complete the game under level 20

### Secret
- **???** Hidden achievements revealed upon unlock

## File Structure

```
plugins/achievement-system/
├── metadata.json          # Plugin metadata
├── plugin.lua             # Main plugin code
├── README.md              # This file
└── CHANGELOG.md           # Version history

achievements/
├── progress.lua           # Achievement progress data
├── custom.lua             # Custom achievement definitions
└── export_*.txt           # Exported achievement data
```

## Data Persistence

Achievement progress is saved to `achievements/progress.lua` and includes:
- Unlocked achievement IDs and timestamps
- Total achievement points
- Custom achievement definitions

## API Requirements

- `GetParty()`: Party roster information
- `GetCharacter(id)`: Character stats and spells
- `GetInventory()`: Inventory and GP tracking
- `GetEspers()` / `GetEsperInventory()`: Esper collection
- `ReadFile()`, `WriteFile()`: Achievement data persistence

## Permissions

- `read_save`: Read save file data for achievement validation
- `ui_display`: Display achievement interface
- `file_io`: Persist achievement progress and export data

## Notes

- Many achievements use placeholder validation (return false) pending API availability
- Retroactive scanning checks achievements against current save state
- Secret achievements display "???" until unlocked
- Custom achievements require manual unlock
- Achievement points contribute to total score

## Future Enhancements

- Real-time achievement tracking via event hooks
- Integration with external leaderboards
- Achievement sound effects and visual notifications
- Steam-style pop-up notifications
- Community achievement sharing
- Rare achievement showcase

## Version

**Version**: 1.0.0  
**Phase**: 5 (Challenge & Advanced Tools)  
**Author**: FF6 Save Editor Plugin System
