# Randomizer Assistant Changelog

## [1.0.0] - 2026-01-16

### Added
- Initial release of Randomizer Assistant plugin
- **Spoiler log import system**:
  - JSON format parser (basic structure support)
  - CSV format parser with header detection
  - Automatic format detection by file extension
  - Import from `randomizer/spoiler_logs/` directory
- **Item location tracker**:
  - Full item placement listing
  - Check-off system for obtained items
  - Progress percentage calculation
  - Sorted display by location name
  - Checkmark indicators ([✓] checked, [ ] unchecked)
- **Character location tracker**:
  - Character recruitment location display
  - Recruitment progress tracking
  - Check-off system for recruited characters
  - Sorted by character name
- **Boss shuffle tracker**:
  - Boss replacement location display
  - Defeated boss tracking
  - Boss completion percentage
  - Sorted by location
- **Ability randomization display**:
  - View randomized character command abilities
  - Sorted by character name
  - Quick reference for ability swaps
- **Logic validation system**:
  - Key item presence checking (Ragnarok, Airship, Sealed Gate Key)
  - Basic beatability validation
  - Warning system for missing critical items
  - Simplified progression logic
- **Seed metadata viewer**:
  - Import date and filename display
  - Format information (JSON/CSV)
  - Seed contents summary (counts)
  - Seed settings display (if present)
- **Progression tracking**:
  - Persistent progress storage
  - Items obtained tracking
  - Characters recruited tracking
  - Bosses defeated tracking
  - Locations checked tracking
- **Statistics export**:
  - Progress summary generation
  - Completion percentages by category
  - Timestamped export files
  - Seed information inclusion
- Safe API call wrappers with pcall error handling
- Lua table serialization for data persistence

### Spoiler Log Format Support

**JSON Structure**:
```lua
{
  items = {location = item},
  characters = {location = character},
  bosses = {location = boss},
  abilities = {character = ability},
  seed_info = {key = value}  -- optional metadata
}
```

**CSV Structure**:
```
Category,Location,Item/Entity
item,Narshe Weapon Shop,Elixir
character,Phantom Train,Sabin
boss,Mt. Kolts,Kefka
ability,Terra,Blitz
```

### Features
- Manual location check-off system
- Progress persistence across sessions
- Multiple tracking categories (items/characters/bosses)
- Basic logic validation
- Export functionality

### Technical Details
- **Lines of Code**: ~870
- **Permissions**: read_save, ui_display, file_io
- **Dependencies**: None (uses Lua standard library)
- **Data Format**: Lua table literal serialization
- **Storage Files**:
  - `randomizer/seed_data.lua` - imported seed information
  - `randomizer/progress.lua` - progression tracking
  - `randomizer/spoiler_logs/` - user-provided spoiler logs
  - `randomizer/stats_*.txt` - exported statistics

### Parsing Logic
- **JSON**: Simplified parser using pcall(load) (requires Lua-compatible JSON structure)
- **CSV**: Manual field parsing with quote handling
- **Format Detection**: Automatic by file extension (.json, .csv)

### Logic Validation
- Checks for key progression items (Airship, Ragnarok, Sealed Gate Key)
- Validates item presence in seed data
- Simplified beatability check (not comprehensive)
- Warning system for potential issues

### Progress Tracking
- **items_obtained**: Dictionary of obtained items
- **characters_recruited**: Dictionary of recruited characters
- **bosses_defeated**: Dictionary of defeated bosses
- **locations_checked**: Dictionary of checked locations

### Known Limitations
- JSON parser is simplified (no support for complex nested structures)
- CSV parsing expects specific column order
- Logic checking is basic (not a full logic solver)
- No automatic progress detection from save file
- Manual check-off system (not integrated with game state)
- No accessibility calculation (which locations are currently reachable)
- No hint system integration

### Notes
- Phase 5, Plugin 4 of 4
- Part of Challenge & Advanced Tools phase
- Designed for Beyond Chaos, Worlds Collide, and similar randomizers
- Spoiler logs must be manually placed in `randomizer/spoiler_logs/` directory
- Progress tracking is persistent but requires manual updates

### Randomizer Compatibility
- **Beyond Chaos**: JSON/CSV spoiler log support
- **Worlds Collide**: JSON/CSV spoiler log support
- **Custom Randomizers**: Any JSON/CSV format with standard structure

### Future API Requirements
- Game state reading for automatic progress detection
- Memory reading for real-time item/character tracking
- Logic solver integration for accessibility calculation
- Hint system API for guided progression
