# Changelog

All notable changes to the Character Stats Display plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-16

### Added
- Initial release of Character Stats Display plugin
- Complete character stats overview display
- Basic stats display (Level, Experience, Status)
- HP/MP display with percentage calculations
- Combat stats display (Vigor, Stamina, Speed, Magic)
- Equipment slots display (all 6 slots)
- Learned spells list with proficiency levels
- Command assignments display
- Interactive character selection dialog
- Configurable display sections
- Spell sorting by proficiency or ID
- Configurable maximum spells display limit
- Comprehensive error handling
- Permission validation
- Detailed logging

### Technical Details
- Plugin API Version: 1.0
- Minimum Editor Version: 3.4.0
- Permissions: read_save, ui_display
- Lines of Code: ~280
- Configuration Options: 5

### Testing
- Tested with multiple save files
- Tested with all 16 character slots
- Tested with various equipment configurations
- Tested with different spell counts
- Tested error conditions (no characters, invalid selection)
- Tested permission requirements

### Known Limitations
- Equipment shown as "Item #XXX" (no name lookup)
- Spells shown as "Spell #XX" (no name lookup)
- Maximum 20 spells displayed by default

## [Unreleased]

### Planned for v1.1.0
- Add equipment name lookups from game database
- Add spell name lookups from game database
- Add character comparison feature
- Add export to clipboard functionality
- Improve display formatting

### Planned for v2.0.0
- Add graphical character preview
- Add stat visualization (bars/graphs)
- Add equipment recommendations
- Add spell coverage analysis
- Add party-wide stats comparison
