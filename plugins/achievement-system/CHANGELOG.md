# Achievement System Changelog

## [1.0.0] - 2026-01-16

### Added
- Initial release of Achievement System plugin
- 50+ predefined achievements across 4 categories (Combat, Collection, Challenge, Secret)
- Achievement categories with point values (10-150 points per achievement)
- **Combat achievements** (15 total):
  - First Blood, Century Club, Overkill, Magic Mastery, No Damage Victory
  - Speed Demon, Solo Champion, All Equipment Master, Esper Bond, Level 99 Club
  - Perfect Party, Rage Master, Blitz Virtuoso, Tools Expert, Lore Scholar
- **Collection achievements** (15 total):
  - Item Hoarder, Millionaire, Esper Collector, Weapon Master, Armor Collector
  - Relic Hunter, Full Party, Rare Find, Dragon Slayer, Treasure Hunter
  - Bestiary Complete, Item Completionist, Steal Master, Colosseum Champion, Slots Jackpot
- **Challenge achievements** (12 total):
  - Low Level Legend, Natural Magic Block, Solo Run, No Equipment, Speedrunner
  - Perfectionist, No Shop Run, Pacifist, No Esper Run, Fixed Party
  - No Save Challenge, Superboss Slayer
- **Secret achievements** (8 total):
  - Hidden achievements revealed upon unlock (Shadow Saved, Cyan's Dream, Opera Star, etc.)
- Custom achievement creation with name, description, category, and point value
- Retroactive achievement checking to scan save file for eligible unlocks
- Progress tracking with unlock timestamps
- Achievement points system with running total
- Statistics dashboard showing:
  - Overall progress percentage
  - Progress by category
  - Recent unlocks (last 5)
  - Total points earned
- Achievement viewer with filtering by category
- Manual unlock functionality for custom achievements
- Export achievements to timestamped text file
- File persistence in `achievements/` directory
- Safe API call wrappers with pcall error handling
- Secret achievement reveal system (shows "???" until unlocked)

### Features
- View all achievements or filter by category
- Track progress with percentage completion
- Create unlimited custom achievements
- Export achievement data for sharing/backup
- Automatic serialization/deserialization of progress data
- Timestamp tracking for each unlock
- Points-based progression system

### Technical Details
- **Lines of Code**: ~750
- **Permissions**: read_save, ui_display, file_io
- **Dependencies**: GetParty, GetCharacter, GetInventory, GetEspers/GetEsperInventory
- **Data Format**: Lua table literal serialization
- **Storage Files**: 
  - `achievements/progress.lua` - unlock progress and points
  - `achievements/custom.lua` - custom achievement definitions
  - `achievements/export_*.txt` - exported achievement reports

### Known Limitations
- Many achievements use placeholder validation pending API availability
- No real-time event tracking (requires event hook API)
- Custom achievements require manual unlock
- No Steam/platform integration
- No multiplayer/leaderboard functionality

### Notes
- Phase 5, Plugin 1 of 4
- Part of Challenge & Advanced Tools phase
- Provides foundation for achievement-based gameplay
- Extensible framework for future achievement additions
