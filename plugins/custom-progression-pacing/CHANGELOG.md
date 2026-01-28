# Changelog - Custom Progression Pacing

All notable changes to the Custom Progression Pacing plugin are documented in this file.

## [1.0.0] - 2026-01-16

### Added - Initial Release

#### Core Features
- Experience gain rate multiplier (0.1x to 100x)
  - Control character leveling speed
  - From extreme slowdown to instant leveling
  - 0.1x = 10x slower, 10x = 10x faster, 100x = 100x faster

- Spell learning rate multiplier (0x to 10x)
  - Control magic ability acquisition speed
  - 0x = spells never learned
  - 10x = spells learned 10x faster

- Gil acquisition rate (0.1x to 10x)
  - Control currency earning speed
  - 0.1x = 10% of normal rewards
  - 10x = 10x the normal rewards

- Drop rate modification (0.1x to 10x)
  - Control item drop frequency
  - 0.1x = 10% drop frequency
  - 10x = 10x drop frequency

- AP gain rate for espers (0.1x to 10x)
  - Control esper stat growth speed
  - Affects overall stat distribution
  - From minimal growth to rapid growth

- Gold found rate (0.1x to 10x)
  - Control treasure chest values
  - Independent from combat rewards
  - Affects item purchasing power

#### Preset Pacing Profiles
- Normal Progression (1.0x all rates)
  - Standard FF6 progression
  - Recommended for balanced gameplay

- Speedrun Profile (5x EXP, 2x Spells, 2x Gil)
  - Fast progression for quick playthroughs
  - Good for speedrun practice

- Casual Profile (2x EXP, 1.5x Spells, 1.5x Items)
  - Relaxed progression
  - Story-focused gameplay

- Completionist Profile (0.5x EXP, 0.7x Spells, 0.8x Items)
  - Slower progression
  - Extended playthrough
  - More grinding required

- Hardcore Profile (0.2x EXP, 0.5x Spells, 0.5x Items)
  - Extreme challenge
  - Minimal progression
  - For experienced players

- Creative Profile (1x EXP, 5x Spells, 0.5x Gil)
  - Spell-focused gameplay
  - All magic available quickly
  - Less combat requirement

- Economic Profile (1x EXP, 5x Drops, 3x Gil)
  - Item-heavy progression
  - Frequent equipment upgrades
  - Treasure-focused

#### Configuration Management
- Custom rate configuration
  - Define rates for each progression system
  - Per-rate validation
  - Clamping to valid ranges

- Configurable per character or global
  - Global application in v1.0.0
  - Per-character planned for future

- Real-time rate adjustment
  - Rates apply to new actions
  - No need to restart game
  - Can switch profiles mid-game

- Rate validation
  - Automatic clamping to min/max
  - Error handling for invalid rates
  - Feedback on actual applied rates

#### Rate Tracking
- Rate history tracking (50 entries)
  - Timestamps for all changes
  - Profile names recorded
  - Applied rates stored

- View current active rates
  - Formatted as percentage (e.g., "5.0x (500%)")
  - All 6 systems displayed
  - Easy comparison to defaults

- Rate change history
  - Last 50 modifications recorded
  - Chronological order
  - Clear timestamps

- Clear history option
  - Reset rate change tracking
  - Clean history management

#### Configuration Export
- Export current configuration
  - JSON format support
  - Text format support
  - Shareable with other players

- Import custom profiles
  - Load saved configurations
  - Apply previous setups
  - Community profile sharing

#### UI and UX
- Clean menu interface
- Preset profile selection
- Custom rate input
- Current rate display
- History viewer
- Export functionality
- Settings menu

#### Documentation
- Comprehensive README (~6,000+ words)
- Installation instructions
- Usage guide with examples
- Technical specifications
- Troubleshooting guide
- FAQ section
- Rate combination recommendations
- Performance specifications
- Known limitations

### Technical Implementation
- ~620 lines of Lua code
- Modular function design
- Preset profile system
- Rate validation and clamping
- History management (last 50 entries)
- Configuration export/import
- Plugin interface pattern
- Safe API call wrappers
- Comprehensive logging

### Quality Assurance
- Code tested for syntax
- Rate clamping verified
- Profile application tested
- History recording validated
- Export functionality confirmed
- Documentation accuracy checked
- Example scenarios provided

### Compatibility
- Phase 6 (Gameplay-Altering Plugins)
- Batch 3 implementation
- Compatible with existing Phase 6 plugins
- Works with any FF6 save file
- Requires read_save, write_save, ui_display permissions

### Known Issues - None
This is a stable initial release with no known issues.

### Future Roadmap
- v1.1.0: Per-character rate configuration
- v1.2.0: Rate curve customization
- v1.3.0: Conditional rate changes
- v2.0.0: Advanced profile builder

## Versioning

- **v1.x.x** - Phase 6 Batch 3 core rate control
- **v2.x.x** - Advanced profile system
- **v3.x.x** - Integration with other plugins

## Support

For issues or questions:
1. Check README Troubleshooting section
2. Review FAQ section
3. Check GitHub Issues
4. Create new issue with specific rates and effects

## Credits

**Phase:** 6 (Gameplay-Altering Plugins)  
**Batch:** 3  
**Week:** 13, Day 3  
**Developer:** FF6 Save Editor Team  
**Date:** January 16, 2026

---

**Current Version:** 1.0.0 (Stable)  
**Last Updated:** 2026-01-16  
**Status:** Active Maintenance
