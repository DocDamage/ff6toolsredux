# Changelog - Treasure Chest Tracker

All notable changes to this plugin are documented here. This project follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16
### Added
- Initial release of the Treasure Chest Tracker (read-only).
- Summary dashboard with opened/unopened/unknown counts and world detection.
- Unopened view with optional treatment of Unknown as unopened (configurable).
- Missables view highlighting WoB-only and point-of-no-return chests.
- World filter, region browser, keyword search, and exportable checklist.
- Safety-first status handling: unknown flags surface as `?`, never guessed.
- Configurable pagination, missable highlighting, value estimates, and world fallback.

### Known Limitations
- Chest flag API may be unavailable in some builds; affected entries show `?`.
- Chest database is a curated subset; expand as more IDs and flags become available.
- Gil values are estimates for prioritization and may vary by version.

---

## [Unreleased]
### Planned
- Full chest database (~400+) when the flag API is finalized.
- Per-dungeon progress summaries and total value aggregation.
- Optional manual toggles (write-enabled variant) for users without chest APIs.
- Cross-links to World of Balance Checklist and Bestiary for route planning.
