# Changelog - Infinite Resources Mode

All notable changes to the Infinite Resources Mode plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16

### Added
- Initial release of Infinite Resources Mode plugin
- **Full Infinite Mode** - One-click max all resources (Gil, items, MP)
- **Selective Infinite Mode** - Choose which resources to make infinite
  - Separate controls for Gil, consumables, equipment, all items, MP
  - Mix and match options for custom configurations
- **Auto-Replenish System** - Automatically restore consumed resources
  - Configurable check interval (default: 60 seconds)
  - Restores items below 99 back to 99
  - Restores MP to maximum for all characters
  - Enable/disable functionality at runtime
- **Backup & Restore** - Automatic backup before modifications
  - Saves Gil, all item quantities, character MP stats
  - One-command restoration to undo changes
  - Timestamped backups for tracking
- **Operation Logging** - Comprehensive modification tracking
  - Logs all resource modifications with timestamps
  - Configurable log size (default: 100 entries)
  - View recent operations for debugging
- **Status Monitoring** - Check current plugin state
  - Enabled/disabled status
  - Auto-replenish active status
  - Backup availability check
  - Operation count tracking
- **Resource Display** - View current resource state
  - Gil display
  - Item count and max item count
  - Character MP display (current/max)
- **Safe API Calls** - Error handling for all API operations
  - pcall wrappers for all FF6 API calls
  - Error logging and reporting
  - Graceful failure handling
- **Configuration Options** - Customize plugin behavior
  - Max Gil value (default: 9,999,999)
  - Max item quantity (default: 99)
  - Max MP bonus (default: 9,999)
  - Auto-replenish interval (default: 60 seconds)
  - Log size limits (default: 100 entries)

### Technical Details
- **Lines of Code:** ~480 LOC
- **Complexity:** Basic-Intermediate
- **API Usage:** Gil, Item, MP read/write operations
- **Permissions:** read_save, write_save, ui_display
- **Phase:** 6, Batch: 4

### Documentation
- Comprehensive README (~5,000 words)
- 10+ use case examples
- Configuration guide
- Troubleshooting section
- FAQ with 10 questions
- Technical details and API usage

### Known Limitations
- Item slots: Only existing item slots can be modified
- Character availability: MP only modified for unlocked characters
- In-battle resources: Changes apply outside battle
- Auto-replenish timing: 60-second intervals (configurable)

### Safety Features
- Automatic backup creation before modifications
- Confirmation dialogs for destructive operations
- Error handling with graceful failures
- Operation logging for audit trail
- Restore functionality to undo changes

---

## Future Plans

### Planned for 1.1.0
- Export/import infinite resource configurations
- Custom item category definitions
- Per-character selective infinite MP
- Battle-synchronous auto-replenish
- Resource usage statistics tracking
- Integration with other economy plugins

### Planned for 1.2.0
- UI integration for graphical control panel
- Preset configurations (Casual, Testing, Magic-Focus, etc.)
- Scheduled resource replenishment
- Conditional infinite resources (if Gil < X, set to Y)
- Resource consumption rate tracking

### Planned for 2.0.0
- Advanced auto-replenish with battle detection
- Resource threshold alerts
- Integration with challenge mode plugins
- Multiplayer-safe resource sharing
- Cloud save integration for configurations

---

## Version History

| Version | Date | LOC | Features | Status |
|---------|------|-----|----------|--------|
| 1.0.0 | 2026-01-16 | 480 | Initial release, full/selective infinite, auto-replenish | ✅ Released |

---

## Breaking Changes

None (initial release)

---

## Migration Guide

Not applicable (initial release)

---

## Deprecations

None (initial release)

---

## Contributors

- FF6 Editor Plugin System - Initial development
- FF6 Modding Community - Testing and feedback

---

**For detailed usage instructions, see README.md**
