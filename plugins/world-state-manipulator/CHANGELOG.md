# Changelog - World State Manipulator Plugin

All notable changes to the World State Manipulator plugin will be documented in this file.

⚠️ **EXPERIMENTAL PLUGIN - USE WITH EXTREME CAUTION**

## [1.0.0] - 2026-01-16

### Added
- Initial release of World State Manipulator plugin
- World state control (toggle WoB ↔ WoR, set specific world)
- 36+ event flag system (story events, character recruitment, dungeons)
- Event flag manipulation (enable/disable individual flags)
- Bulk event operations:
  - `completeWoBStory()` - Complete all 16 WoB events
  - `recruitAllCharacters()` - Flag all 13 WoR characters as recruited
  - `clearAllEvents()` - Reset to new game state (dangerous)
- 6 world state presets:
  - `fresh_wob` - New game in World of Balance
  - `fresh_wor` - New game in World of Ruin (Celes on island)
  - `end_wob` - End of WoB (pre-cataclysm)
  - `early_wor` - Early WoR (characters recruitable)
  - `late_wor` - Late WoR (all characters recruited)
  - `endgame` - All dungeons accessible
- Dangerous mode protection for critical operations
- Dual world access (extremely experimental, not recommended)
- World state analysis and display
- Configuration export to text format
- Automatic backup before modifications
- Backup restoration system
- Operation logging (last 100 operations)

### Features
- **36 Event Flags:** Major story events, character recruitment, dungeon completions
- **6 Presets:** Quick access to major story points
- **Safety System:** Dangerous mode lock prevents accidental corruption
- **Export System:** Save world state details to text
- **Logging:** Track all state modifications

### Technical
- ~680 lines of Lua code
- ~7,000 words of documentation
- Safe API call wrappers with error handling
- Operation logging system
- State management for backups and caches

### Event Flags Included
**Major Story Events (18 flags):**
- Narshe Battle, Figaro Castle, South Figaro, Mt. Kolts
- Returner Hideout, Lethe River, Imperial Camp, Doma Siege
- Phantom Train, Baren Falls, Opera House, Vector
- Magitek Factory, Esper Gate, Thamasa, Floating Continent
- Cataclysm (critical), Solitary Island

**Character Recruitment (13 flags):**
- All 13 recruitable characters in WoR

**Dungeons (4 flags):**
- Phoenix Cave, Ancient Castle, Fanatic's Tower, Kefka's Tower

### Use Cases
- Speedrun testing (jump to specific story points)
- Character recruitment testing (skip to WoR)
- Story sequence skipping (bypass early game)
- Event flag debugging (identify blocking events)
- World state experiments (test WoB/WoR interactions)

### Safety Features
- Automatic backups before modifications
- Dangerous mode protection for critical flags
- Operation logging (100 entries)
- Restore system (undo all changes)
- Warning system for risky operations

### Known Limitations
- Save corruption risk with improper flag combinations
- Softlock potential with wrong world state + event flags
- Vanilla FF6 incompatibility (modified saves won't load)
- Dual world access extremely unstable (crashes/corruption)
- Event dependencies not enforced (user must know prerequisites)

### Warnings
⚠️ **High-Risk Operations:**
- Toggling world mid-game
- Clearing all events
- Dual world access (don't use this!)
- Modifying critical flags without dangerous mode

⚠️ **Corruption Risks:**
- Wrong world + event flag combinations
- Clearing events at wrong story point
- Enabling dual world access

## Future Enhancements (Potential)

### Version 1.1.0
- [ ] Event dependency validation (enforce prerequisites)
- [ ] World state validation (detect invalid combinations)
- [ ] Safe world toggle (auto-set required events)
- [ ] More event flags (minor story events, optional content)
- [ ] Preset customization (create custom presets)

### Version 1.2.0
- [ ] Event timeline visualization
- [ ] Story progression graph
- [ ] Suggested next events (based on current state)
- [ ] Integration with Character Roster Editor
- [ ] Auto-fix invalid state combinations

### Version 2.0.0
- [ ] Full event dependency system
- [ ] Visual world state editor
- [ ] Story branching support
- [ ] Save state snapshots (multiple backups)
- [ ] Integration with Challenge Mode Validator
- [ ] Community preset sharing

## Changelog Format

This changelog follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.

### Change Categories
- **Added:** New features
- **Changed:** Changes to existing functionality
- **Deprecated:** Soon-to-be removed features
- **Removed:** Removed features
- **Fixed:** Bug fixes
- **Security:** Security fixes

---

**Current Version:** 1.0.0  
**Release Date:** 2026-01-16  
**Plugin Status:** Experimental (High Risk)

⚠️ **ALWAYS BACKUP YOUR SAVE BEFORE USING THIS PLUGIN!**
