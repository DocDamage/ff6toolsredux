# Location Database Plugin - Changelog

## [1.0.0] - 2026-01-16
### Added
- Location lookup by ID, name, region, and type
- Location catalog management (add locations, summary by region/type)
- Visit tracking and pattern analysis
- Phase 11 integrations:
  - Analytics Engine for visit pattern analysis
  - Import/Export Manager for location export (JSON/CSV/XML)
  - Integration Hub sync for cross-plugin map data
  - Backup/Restore snapshots for location catalog
  - API Gateway REST endpoints (/api/locations/:id, /api/locations/search, /api/locations/summary)
- Database Persistence Layer integration for persistent storage
- Operation logging with configurable log size

### Configuration
- Support for 256 locations (CONFIG.MAX_LOCATIONS)
- Location types: town, dungeon, overworld, tower, cave, ruin
- Regions: World of Balance, World of Ruin
- Operation log capped at 50 entries
