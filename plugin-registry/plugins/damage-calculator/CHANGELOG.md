# Changelog - Damage Calculator

All notable changes to this plugin are documented here, following [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-01-16
### Added
- Initial release (read-only) with physical and magical damage calculators.
- Comparison view (physical vs. magic) with shared inputs.
- Elemental modifiers (weak/resist/absorb/null), row/back attack, critical, and multi-target penalties.
- Default stat fallbacks with user override prompts.
- Help view explaining formulas and safety behavior.

### Known Limitations
- Simplified formulas; weapon/relic quirks and rounding edge cases not modeled.
- Enemy stats fetched only if `GetBestiaryEntry` is available; otherwise defaults.
- No scenario save/load yet.

---

## [Unreleased]
### Planned
- Add weapon/relic flags (ignore defense, vigor caps, random variance).
- Add spell-specific quirks and elemental proc handling.
- Scenario save/load and batch comparisons.
