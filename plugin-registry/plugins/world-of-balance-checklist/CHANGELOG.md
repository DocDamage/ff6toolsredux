# Changelog - World of Balance Checklist

All notable changes to this plugin will be documented in this file. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16
### Added
- Initial release of the World of Balance Checklist (read-only).
- Summary dashboard with world detection, character, and esper progress.
- Character checklist for all roster slots (WoR-only entries labeled for clarity).
- Esper checklist for expected WoB-accessible magicite with dynamic ownership detection.
- Missable content reference covering items, rages, dances, and key events with quick tips.
- Event flag viewer for opera, banquet, floating continent, world state, and factory completion.
- Export view for copyable progress snapshots.
- Config toggles for world fallback, warnings, pagination, and verbosity.

### Known Limitations
- No chest flag API yet (treasure tracking planned for Phase 2.2).
- Rage/Dance learn flags not exposed; missables shown as reminders with `?` status.
- World detection falls back to WoB when API is unavailable (configurable).

---

## [Unreleased]
### Planned
- Add chest flag integration once the API is available.
- Add optional manual checklist toggles for items without flag coverage.
- Expand esper list variants per localization/version and allow profile selection.
- Cross-link with Bestiary and Treasure Chest Tracker once those plugins are published.
