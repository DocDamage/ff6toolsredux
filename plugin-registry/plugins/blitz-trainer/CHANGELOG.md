# Changelog - Sabin's Blitz Input Trainer

All notable changes to the Blitz Input Trainer plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16

### Added
- Initial release of Sabin's Blitz Input Trainer plugin
- Complete database of all 8 Blitzes with detailed information
- Visual input diagrams using arrow symbols (↑ ↓ ← → ↗ ↘ ↙ ↖)
- Six viewing modes:
  - Quick reference card
  - Individual Blitz details
  - All Blitzes compendium
  - Input system guide
  - Damage calculator
  - Strategic recommendations
- Damage calculation based on Sabin's current stats:
  - Physical Blitz formulas
  - Magical Blitz formulas
  - Healing formula (Mantra)
- Difficulty ratings for each Blitz (Easy to Very Hard)
- Learn level information (when each Blitz becomes available)
- Complete input sequences with button-by-button breakdown
- Strategic recommendations:
  - Best Blitzes by game phase
  - Boss battle strategies
  - Random encounter optimization
  - Easiest vs. hardest inputs
- Input tips and tricks for mastering complex sequences
- Common input pattern guide (quarter-circle, semi-circle, etc.)
- Battle usage strategies for each Blitz
- Detailed README with complete usage guide (~3,900 words)
- Configuration options for customization
- Works as reference tool even without save file loaded

### Features
- **Complete Blitz Database:** All 8 Blitzes from level 1 to 70
- **Visual Diagrams:** ASCII arrow representations of inputs
- **Smart Calculator:** Estimates damage using current character stats
- **Difficulty System:** Easy, Medium, Hard, Very Hard ratings
- **Strategic Guide:** Recommendations for different situations
- **Formula Display:** Show/hide damage calculation formulas
- **Reference Mode:** Works without needing Sabin in save file
- **User-Friendly UI:** Clean, organized text-based interface

### Technical Details
- Lines of Code: ~380 LOC
- Documentation: ~3,900 words
- Permissions: read_save, ui_display
- Minimum Editor Version: 3.4.0
- Language: Lua 5.1
- Unique Feature: Works without save file (reference mode)

### Blitz Coverage
1. Raging Fist (Level 1)
2. Aura Cannon (Level 6)
3. Meteor Strike (Level 10)
4. Suplex (Level 15)
5. Fire Dance (Level 23)
6. Mantra (Level 30)
7. Air Blade (Level 42)
8. Phantom Rush (Level 70)

### Known Limitations
- Reference only (cannot execute Blitzes from plugin)
- Damage estimates don't account for enemy defense
- No interactive practice mode for inputs
- Estimates are base calculations without modifiers

---

## [Unreleased]

### Planned Features for v1.1
- Input pattern visual improvements
- More detailed strategy guides per Blitz
- Boss-specific Blitz recommendations
- Enhanced damage formulas with enemy resistance calculations

### Planned Features for v2.0
- Interactive practice mode (if API supports input simulation)
- Success rate tracking
- Input timing recommendations
- Video tutorial links
- Blitz comparison tool

### Potential Enhancements
- Integration with battle simulator (if available)
- Community-contributed strategies
- Speedrun routing recommendations
- Challenge run Blitz restrictions
- Input replay system
- Difficulty scaling recommendations

---

## Version Numbering

- **Major version (X.0.0):** Incompatible API changes or major feature overhauls
- **Minor version (0.X.0):** New features, backwards-compatible
- **Patch version (0.0.X):** Bug fixes and minor improvements

---

## Roadmap

### Short-term (1-3 months)
- Enhanced visual diagrams
- Boss-specific recommendations
- Community strategy integration

### Medium-term (3-6 months)
- Practice mode implementation
- Success rate tracking
- Video tutorials

### Long-term (6-12 months)
- Battle integration (if API permits)
- Input replay/record system
- Community challenges

---

## Links

- [Plugin Repository](https://github.com/ff6-save-editor/plugin-registry)
- [Issue Tracker](https://github.com/ff6-save-editor/plugin-registry/issues)
- [Documentation](README.md)

---

## Notes

This plugin serves as both a reference guide and training tool for Sabin's unique Blitz command system. Unlike character trackers, it functions independently of save file data, making it useful as a pure reference tool for learning inputs and strategies.

The Blitz system is one of FF6's most iconic mechanics, and this plugin aims to make mastering all 8 techniques accessible to players of all skill levels.
