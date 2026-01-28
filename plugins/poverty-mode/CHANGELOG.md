# Changelog - Poverty Mode

All notable changes to the Poverty Mode plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16

### Added
- Initial release of Poverty Mode plugin
- **5 Poverty Level Presets**
  - Light Poverty (1,000 Gil, 10 items max)
  - Moderate Poverty (500 Gil, 5 items max)
  - Strict Poverty (100 Gil, no consumables)
  - Hardcore Poverty (0 Gil, no consumables, equipped items only)
  - Extreme Poverty (0 Gil, nothing allowed except equipped items)
- **Custom Poverty Configuration**
  - Set custom Gil amounts
  - Set custom item quantity limits
  - Toggle consumable removal
  - Toggle unequipped equipment removal
- **Challenge Tracking**
  - Track found items during poverty run
  - Log item acquisitions with source tracking
  - Found items log (timestamp, item ID, quantity, source)
  - Up to 200 log entries maintained
- **Compliance Checking**
  - Verify save compliance with poverty rules
  - Detect Gil rule violations
  - Detect consumable violations
  - Detect item quantity limit violations
  - Violation reporting with details
- **Challenge Proof Export**
  - Export challenge completion proof
  - Includes poverty level, duration, found items count
  - Includes compliance status and current Gil
  - Shareable for community validation
- **Backup & Restoration**
  - Automatic backup creation before poverty application
  - Backup includes Gil, all items, equipped items
  - One-command restoration to undo poverty mode
  - Timestamped backups for tracking
- **Resource Removal Functions**
  - Zero Gil function
  - Set Gil to custom amount
  - Remove all consumables (ID 0-99)
  - Remove unequipped equipment (ID 100-254)
  - Limit item quantities to max allowed
  - Equipped item detection and preservation
- **Status Monitoring**
  - Check enabled/disabled state
  - View current poverty level
  - Track challenge duration
  - Check backup availability
  - Count found items
  - Count compliance violations
- **Resource Display**
  - Display current Gil
  - Display item count and total quantity
  - Display compliance status
  - List compliance violations
- **Operation Logging**
  - Comprehensive logging of all operations
  - Timestamp all modifications
  - Log size limit (200 entries)
  - Operation details tracking
- **Safe API Calls**
  - Error handling for all API operations
  - pcall wrappers for safety
  - Graceful failure handling

### Technical Details
- **Lines of Code:** ~530 LOC
- **Complexity:** Intermediate
- **API Usage:** Gil, Item, Equipment read/write operations
- **Permissions:** read_save, write_save, ui_display
- **Phase:** 6, Batch: 4

### Documentation
- Comprehensive README (~5,500 words)
- 5 poverty level descriptions
- Custom configuration guide
- Challenge tracking documentation
- 5 use case examples
- Survival strategies section
- Troubleshooting section
- FAQ with 10 questions

### Known Limitations
- Item category ranges are approximate (ID 0-99 consumables, 100-254 equipment)
- Equipped item detection: only currently equipped items preserved
- Story/key items: some cannot be removed (game requirement)
- Shop access: shops accessible but unusable with zero Gil
- Found items tracking: requires manual tracking calls

### Safety Features
- Automatic backup creation before modifications
- Confirmation dialogs for destructive operations
- Error handling with graceful failures
- Operation logging for audit trail
- Restore functionality to undo changes
- Compliance checking to verify challenge validity

### Challenge Features
- 5 difficulty presets from Light to Extreme
- Custom configuration for personalized challenges
- Found items tracking for documentation
- Compliance verification for rule validation
- Challenge proof export for community sharing
- Challenge duration tracking

---

## Future Plans

### Planned for 1.1.0
- Automatic found item detection (via API hook)
- Real-time compliance monitoring
- Challenge milestones and achievements
- Integration with Challenge Timer plugin
- Poverty mode leaderboard export format
- Multiple backup slots

### Planned for 1.2.0
- Advanced compliance rules (time-based, location-based)
- Poverty mode achievements system
- Community challenge templates
- Poverty mode statistics dashboard
- Integration with streaming overlays
- Save file integrity verification

### Planned for 2.0.0
- Multiplayer poverty challenges
- Synchronized poverty runs
- Real-time challenge comparison
- Automated challenge validation
- Cloud-based challenge proof storage
- Community challenge rankings

---

## Version History

| Version | Date | LOC | Features | Status |
|---------|------|-----|----------|--------|
| 1.0.0 | 2026-01-16 | 530 | Initial release, 5 presets, tracking, compliance | ✅ Released |

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
- FF6 Challenge Run Community - Testing and feedback
- No-Shop Run Pioneers - Feature suggestions

---

**For detailed usage instructions, see README.md**
