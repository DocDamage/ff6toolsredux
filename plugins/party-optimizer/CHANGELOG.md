# Changelog - Party Optimizer Plugin

## [1.2.1] - 2026-01-16 (TIER 2 PHASE 1 - PERSISTENCE LAYER INTEGRATION)
### Added - Database Persistence Layer
- **Persistent Party Configurations:** Party compositions now persist via Database Persistence Layer v1.2.0
- **Auto-Save Optimization:** `optimizePartyComposition()` automatically persists party configs to database layer
- **Load on Retrieval:** `loadPartyData()` loads persisted party configurations from database
- **Party Durability:** Party configurations survive across editor sessions
- **Graceful Fallback:** If persistence layer unavailable, falls back to in-memory storage

### Changed - Integration Points
- `optimizePartyComposition()` now calls `db.savePartyConfig()` for persistence
- New `loadPartyData()` function calls `db.loadPartyConfig()` to retrieve persisted data

## [1.2.0] - 2026-01-16 (PHASE 11 TIER 1 INTEGRATIONS)
### Added
- Analytics-assisted party optimization: `optimizePartyComposition`, `recommendPartyForScenario`
- Automation hooks: `autoConfigureParty`, workflow/rules ready
- Visualization: `visualizePartySynergy`, `generatePartyReport`
- Import/Export: `exportPartyTemplate`, `importPartyTemplate`
- Backup/Restore: `snapshotParty`, `restoreParty`
- Sync: `syncPartyData` via Integration Hub

### Notes
- Dependencies load lazily and fail gracefully if Phase 11 plugins are missing.
- Designed to remain backward compatible with existing saves and workflows.

## [1.1.0] - 2026-01-16
- (Legacy extension in v1_1_upgrades.lua) Equipment integration, growth prediction, esper optimization, boss strategies

## [1.0.0] - 2026-01-16
- Initial release of Party Optimizer core
