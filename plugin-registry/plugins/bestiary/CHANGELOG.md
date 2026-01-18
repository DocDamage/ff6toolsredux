# Changelog - Bestiary

All notable changes to this plugin will be documented here. This project follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16
### Added
- Initial release of the Bestiary (read-only).
- Summary view with seen/unseen/unknown counts.
- Filters: All, Unseen, By Type, By Element, Bosses, Search, Export, Help.
- Status icons with conservative Unknown handling (no guessing when APIs are missing).
- Element, status immunity, drop/steal, morph, sketch/control, rage, and location data for representative core enemies.
- Config toggles for unknown handling, pagination, and which data sections to show.

### Known Limitations
- Encounter flag API may be absent; entries show `?` until provided.
- Database is a curated subset; structured for full 384-entry import later.
- No write operations; cannot mark encounters manually in this version.

---

## [Unreleased]
### Planned
- Full enemy dataset import once available.
- Per-element and per-type statistics in Summary.
- Optional write-enabled manual toggles for encounter flags.
- Cross-links to Rage Tracker and Treasure Chest Tracker for routing.
