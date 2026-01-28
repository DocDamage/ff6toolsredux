# Changelog

All notable changes to the Equipment Optimizer plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-16

### Added
- Initial release of Equipment Optimizer plugin
- Equipment analysis for individual characters
- Party-wide equipment optimization
- 5 optimization goals: Attack, Defense, Magic, Speed, Balanced
- Stat comparison (before/after preview)
- Equipment comparison tool showing top 3 options per slot
- Simplified equipment database (13 weapons, 8 shields, 8 helms, 9 armor)
- Scoring formulas for each optimization goal
- One-click equipment application
- Preview mode with confirmation dialogs
- Current equipment display
- Menu-driven UI with 5 operations
- Comprehensive error handling
- Permission validation (read_save, write_save, ui_display)

### Optimization Goals
- **Attack:** Maximizes offensive power (Attack × 5 + Defense + Speed × 2)
- **Defense:** Maximizes survivability (Defense × 5 + Magic Def × 4 + Attack)
- **Magic:** Maximizes magic power (Magic × 5 + Magic Def × 3 + Defense)
- **Speed:** Maximizes turn frequency (Speed × 5 + Attack × 2 + Defense)
- **Balanced:** All-around optimization (Attack × 2 + Defense × 2 + Magic × 2 + Magic Def × 2 + Speed)

### Configuration
- configurable optimization goal (default: balanced)
- preserveRelics option (default: true)
- previewChanges option (default: true)
- showAllOptions flag for verbose output
- considerRoleMatch flag for future role-based optimization

### Security
- Sandboxed Lua execution environment
- Permission-based access control
- Confirmation dialogs before modifications
- Input validation for character indices
- Error boundaries with pcall wrapper

### Documentation
- Complete README.md with usage instructions
- Optimization formula explanations
- Configuration options documentation
- Troubleshooting guide
- Use cases for different player types

## [Unreleased]

### Planned
- Complete equipment database with all 250+ items from FF6
- Character-specific equipment compatibility checking
- Relic optimization (currently preserved)
- Equipment special effects consideration (elemental, status immunity)
- Set bonus detection (Genji set, etc.)
- Boss-specific equipment recommendations
- Equipment enchantment tracking
- Multi-character equipment comparison side-by-side
- Equipment swap optimization (minimize inventory shuffling)
- Export/import equipment loadouts to JSON
- Role-based recommendations (tank, mage, physical DPS)
- Equipment rarity/value consideration
- Inventory management integration
- Character natural role detection and matching
