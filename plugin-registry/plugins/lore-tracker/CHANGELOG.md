# Changelog - Strago's Lore Tracker

All notable changes to the Lore Tracker plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16

### Added
- Initial release of Strago's Lore Tracker plugin
- Complete database of all 24 Lores with detailed information
- Enemy source tracking for each Lore
- Location guides showing where to find each enemy
- Missable Lore warnings (L? Pearl and Quasar)
- Boss enemy indicators for boss-exclusive Lores
- Five viewing modes:
  - View all Lores
  - View learned Lores only
  - View unlearned Lores only
  - View missable Lores only
  - View statistics
- Progress tracking with completion percentage
- Comprehensive statistics display:
  - Total Lores learned
  - Missable Lores status
  - Boss-only Lores status
- Main menu with user-friendly navigation
- Automatic Strago character detection
- Error handling with informative messages
- Detailed README with complete usage guide
- Configuration options for customization

### Features
- **24 Lore Database:** Complete information for all Blue Magic spells
- **Enemy Sources:** Lists which enemies can teach each Lore
- **Location Information:** Shows where to encounter enemies
- **Missable Warnings:** Highlights Lores only available in WoB
- **Boss Indicators:** Marks Lores that require boss battles
- **Flexible Filtering:** Multiple views for different needs
- **Progress Tracking:** Real-time completion statistics
- **User-Friendly UI:** Clean, organized text-based interface

### Technical Details
- Lines of Code: ~450 LOC
- Documentation: ~4,700 words
- Permissions: read_save, ui_display
- Minimum Editor Version: 3.4.0
- Language: Lua 5.1

### Known Limitations
- Read-only plugin (cannot modify save data)
- Lore spell ID mapping may need adjustment for some ROM versions
- No spell effect descriptions in v1.0.0
- Static database (no external data file support)

---

## [Unreleased]

### Planned Features
- Lore effect descriptions (damage formulas, status effects)
- Damage calculator for each Lore
- Export checklist to text file
- Visual map integration
- Lore priority recommendations
- Enhanced Lore detection for different ROM versions
- Enemy encounter rate information
- Best practices guide for learning each Lore

### Potential Enhancements
- Write mode to unlock Lores (requires write_save permission)
- Integration with other character tracker plugins
- Community-contributed strategy guides
- Video tutorial links
- Speedrun-specific Lore recommendations

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
