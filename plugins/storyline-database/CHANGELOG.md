# Storyline Database Plugin - Changelog

## [1.0.0] - 2026-01-16
### Added
- Story event lookup by ID, name, chapter, status, and type
- Story flags support and progress summary
- Story event management (add events, update flags)
- Phase 11 integrations:
  - Analytics Engine for storyline coverage analysis
  - Import/Export Manager for storyline export (JSON/CSV/XML)
  - Integration Hub sync for cross-plugin story state
  - Backup/Restore snapshots for storyline data
  - API Gateway REST endpoints (/api/storyline/events/:id, /api/storyline/events/search, /api/storyline/summary)
- Database Persistence Layer integration for persistent storage
- Operation logging with configurable log size

### Configuration
- Support for 256 events (CONFIG.MAX_EVENTS)
- Event types: main, side, optional
- Status values: pending, in_progress, completed, blocked
- Operation log capped at 50 entries
