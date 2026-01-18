# Changelog - Alternate Start Generator

## [1.0.0] - 2026-01-16

### Added - Initial Release

#### Core Features
- Skip to World of Ruin
  - Begin at post-apocalypse with Celes and available team
  - Level 30, 5,000 gil
  - All WoR story flags set

- Choose starting characters
  - Select any combination from 14 characters
  - Character availability per world state
  - Roster configuration UI

- Start at specific story events
  - 8 major story event points
  - From game start to final boss
  - Automatic flag configuration

- Custom starting inventory/levels
  - Configurable starting level (1-99)
  - Custom starting gil amount
  - Initial equipment setup
  - Item quantity configuration

- Scenario selection
  - 8 preset templates
  - Pre-configured for different playstyles
  - One-click application

#### Preset Scenarios
- Skip to World of Ruin (WoR, level 30)
- Speedrun - Any% (WoB, level 1, Sabin route)
- Speedrun - 100% (WoB, level 10)
- Celes Solo Run (WoB, level 1, solo character)
- Three-Character Challenge (WoB, level 1, 3 chars)
- Boss Rush - WoB (WoB, level 20, 8 characters)
- Boss Rush - WoR (WoR, level 50, 10 characters)
- Balanced Challenge Party (WoB, level 15, mixed)

#### Advanced Features
- Boss rush mode (fight only major bosses)
- Character-specific scenarios (Celes solo, etc.)
- Speedrun practice starts (optimized routes)
- Event flag manipulation (story progression)
- Export start configurations (shareable)
- Custom scenario creation and saving
- World state toggle (WoB/WoR)

#### Story Events
- Game Start (Terra flashback, WoB start)
- Figaro Castle (after siege, WoB)
- Narshe (first visit, WoB)
- South Figaro (after liberation, WoB)
- Mt. Kolts (after passage, WoB)
- World of Ruin (transition point, WoR)
- Floating Island (final prep, WoR)
- Kefka's Tower (final boss, WoR)

#### Configuration Management
- World state selection (WoB or WoR)
- Story event selection (8 events)
- Character roster configuration
- Starting level customization
- Starting gil customization
- Starting inventory setup
- Event flag configuration
- Validation and safety checks

#### Safety Features
- Scenario compatibility checking
- Dangerous operation warnings
- Invalid combination detection
- Backup reminders
- Configuration export

#### UI and UX
- Scenario list display
- Story event selection
- Character roster UI
- Configuration review
- Export interface

#### Documentation
- Comprehensive README (~5,000+ words)
- Installation instructions
- Usage guide with examples
- Scenario reference table
- Story events listing
- Technical specifications
- Troubleshooting guide
- FAQ section
- Known limitations

### Technical Implementation
- ~750 lines of Lua code
- Preset scenario system
- Custom scenario builder
- World state management
- Event flag system
- Character roster control
- Validation and clamping
- Configuration export/import
- Plugin interface pattern
- Safe API wrappers

### Quality Assurance
- Code tested for syntax
- Scenario validation verified
- World state transitions tested
- Character roster configuration validated
- Event flag application confirmed
- Documentation accuracy checked
- Example scenarios provided

### Compatibility
- Phase 6 (Gameplay-Altering Plugins)
- Batch 3 implementation
- Compatible with other Phase 6 plugins
- Works with any FF6 save file
- Requires read_save, write_save, ui_display permissions

### Known Issues - None
Stable initial release with no known issues.

### Future Roadmap
- v1.1.0: Additional story events (20+ total)
- v1.2.0: Event chain configuration
- v1.3.0: Character placement control
- v2.0.0: Interactive scenario builder

## Versioning

- **v1.x.x** - Phase 6 Batch 3 core alternate starts
- **v2.x.x** - Advanced scenario system
- **v3.x.x** - Integration with other plugins

## Support

For issues:
1. Check README Troubleshooting
2. Review FAQ section
3. Verify scenario compatibility
4. Report with specific scenario details

## Credits

**Phase:** 6 (Gameplay-Altering Plugins)  
**Batch:** 3  
**Week:** 13, Day 4  
**Developer:** FF6 Save Editor Team  
**Date:** January 16, 2026

---

**Current Version:** 1.0.0 (Stable)  
**Last Updated:** 2026-01-16  
**Status:** Active Maintenance
