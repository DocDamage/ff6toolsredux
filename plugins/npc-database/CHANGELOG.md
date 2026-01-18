# NPC Database Plugin - Changelog

## [1.0.0] - 2026-01-16
### Added
- NPC lookup by ID, name, location, and role
- NPC registry management (add NPCs, summary by role/location)
- Interaction tracking and pattern analysis
- Phase 11 integrations:
  - Analytics Engine for interaction analysis
  - Import/Export Manager for NPC registry export (JSON/CSV/XML)
  - Integration Hub sync for cross-plugin NPC data
  - Backup/Restore snapshots for NPC registry
  - API Gateway REST endpoints (/api/npcs/:id, /api/npcs/search, /api/npcs/summary)
- Database Persistence Layer integration for persistent storage
- Operation logging with configurable log size

### Configuration
- Support for 256 NPCs (CONFIG.MAX_NPCS)
- Roles: merchant, quest_giver, ally, antagonist, generic
- Alignment: friendly, neutral, hostile
- Operation log capped at 50 entries
