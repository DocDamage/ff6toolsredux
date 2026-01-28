# Changelog - Party Synergy Analyzer

All notable changes to this plugin will be documented here. This project follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16
### Added
- Initial release (read-only) with role inference, elemental coverage, redundancy, and missing-role detection.
- Views: Summary, Roles, Elements, Suggestions, Help.
- Heuristic spell-based role detection and spell-name-based elemental parsing.
- Safe handling when party/spell APIs are unavailable (explicit error, no guessing).

### Known Limitations
- Heuristic only; equipment/commands not yet considered.
- Localization differences may require updating spell name keys.
- No export view in this version.

---

## [Unreleased]
### Planned
- Equipment/command-aware role inference.
- Exportable analysis report.
- Optional scoring system and party comparison.
