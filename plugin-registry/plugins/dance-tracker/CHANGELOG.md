# Changelog - Mog's Dance Tracker

All notable changes to the Dance Tracker plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16

### Added
- Initial release of Mog's Dance Tracker plugin
- Complete database of all 8 dances with detailed information
- Location guides for each dance with map area names
- Dance effect descriptions (all four abilities per dance)
- World availability indicators (WoB/WoR)
- Five viewing modes:
  - View all dances
  - View learned dances only
  - View unlearned dances only
  - Location guide
  - View statistics
- Progress tracking with completion percentage
- Comprehensive statistics display:
  - Total dances learned
  - World availability breakdown
  - Dance category analysis
- Main menu with user-friendly navigation
- Automatic Mog character detection
- Error handling with informative messages
- Detailed README with complete usage guide
- Configuration options for customization

### Features
- **8 Dance Database:** Complete information for all dances
- **Location Information:** Detailed guides for learning each dance
- **Effect Descriptions:** All four random abilities per dance
- **World Indicators:** WoB and WoR availability
- **Flexible Filtering:** Multiple views for different needs
- **Progress Tracking:** Real-time completion statistics
- **User-Friendly UI:** Clean, organized text-based interface

### Technical Details
- Lines of Code: ~350 LOC
- Documentation: ~3,600 words
- Permissions: read_save, ui_display
- Minimum Editor Version: 3.4.0
- Language: Lua 5.1

### Known Limitations
- Read-only plugin (cannot modify save data)
- Dance detection may need API enhancement for accurate tracking
- Static database (no external data file support)

---

## [Unreleased]

### Planned Features
- Dance effect damage/healing formulas
- Battle strategy recommendations per dance
- Enhanced dance learning detection
- Export checklist to text file
- Visual location maps
- Dance usage statistics (if battle data available)

### Potential Enhancements
- Write mode to unlock dances (requires write_save permission)
- Integration with other character tracker plugins
- Video tutorial links for dance locations
- Animation/effect descriptions

---

## Version Numbering

- **Major version (X.0.0):** Incompatible API changes or major feature overhauls
- **Minor version (0.X.0):** New features, backwards-compatible
- **Patch version (0.0.X):** Bug fixes and minor improvements

---

## Links

- [Plugin Repository](https://github.com/ff6-save-editor/plugin-registry)
- [Issue Tracker](https://github.com/ff6-save-editor/plugin-registry/issues)
- [Documentation](README.md)
