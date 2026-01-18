# Changelog - Colosseum Guide

All notable changes to this plugin will be documented here. This project follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16
### Added
- Initial release of the Colosseum Guide (read-only).
- Bet → opponent → reward mapping with strategy notes and reward notes.
- Filters: All, Owned (inventory-aware), By Tier, Search, Export, Help.
- Enemy hint overlay without requiring full bestiary import.
- Configurable handling for unknown inventory (treat as missing or show as unknown).
- Pagination and lightweight data model for easy extension.

### Known Limitations
- Inventory API may be absent; Owned view will show Unknown unless configured otherwise.
- Database is a representative core; extend `BETS` for full coverage.
- No write operations; cannot mark bets as completed.

---

## [Unreleased]
### Planned
- Full bet table coverage with rarity tags.
- Sorting by reward value and difficulty.
- Optional write-enabled manual tracking.
- Cross-link to Bestiary for full enemy stats.
