# Treasure Database Plugin - Changelog

## [1.0.0] - 2026-01-16
### Added
- Treasure lookup by ID, location, type, and rarity
- Treasure map management (add treasures, mark opened, summary by type/rarity)
- Collection tracking and analytics
- Phase 11 integrations:
  - Analytics Engine for collection progress analysis
  - Import/Export Manager for treasure map export (JSON/CSV/XML)
  - Integration Hub sync for cross-plugin treasure data
  - Backup/Restore snapshots for treasure map
  - API Gateway REST endpoints (/api/treasure/:id, /api/treasure/search, /api/treasure/summary)
- Database Persistence Layer integration for persistent storage
- Operation logging with configurable log size

### Configuration
- Support for 512 treasures (CONFIG.MAX_TREASURE)
- Treasure types: chest, hidden, event, shop
- Rarity levels: common, uncommon, rare, epic, legendary
- Operation log capped at 50 entries
