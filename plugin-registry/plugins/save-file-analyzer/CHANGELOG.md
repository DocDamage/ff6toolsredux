# Changelog

All notable changes to the Save File Analyzer plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-16

### Added
- Initial release of Save File Analyzer plugin
- Comprehensive statistics dashboard
- Overall completion percentage with weighted scoring
- Character statistics analysis:
  - Recruited count (X/14)
  - Level statistics (min, max, average)
  - Total HP/MP across all characters
  - Total experience points
  - Total spells learned
  - HP/MP percentage (current vs max)
- Current party analysis:
  - Party size and composition
  - Average party level
  - Party power rating
  - Individual member HP display
- Inventory statistics:
  - Unique items count
  - Total items count
  - Slots used percentage
  - Estimated inventory value
- Character details view:
  - Per-character HP/MP display
  - Experience points per character
  - Spells learned per character
  - Combat stats (Attack, Defense, Magic, Magic Def)
- Completion checklist:
  - All characters recruited milestone
  - Max level character milestone
  - Average level 50+ milestone
  - All spells learned milestone
  - 100+ unique items milestone
  - Visual progress bars per category
- Statistics export:
  - Text-based report generation
  - Timestamp inclusion
  - Shareable format
- Visual progress indicators (ASCII progress bars)
- Number formatting with comma separators
- Menu-driven UI with 4 operations
- Comprehensive error handling
- Permission validation (read_save, ui_display)

### Completion Calculation
- Weighted formula: Characters (30%) + Inventory (20%) + Spells (25%) + Levels (25%)
- Visual progress bar (20 characters wide)
- Per-category progress tracking
- Milestone detection

### Statistics
- Character count maximum: 14
- Max level: 99
- Max spells per character: 54
- Max inventory slots: 255
- Max gil: 9,999,999
- Party size: 4 members

### Configuration
- showDetailedStats: Toggle detailed breakdowns
- calculateCompletion: Enable/disable percentage calculations
- trackMissables: Track missable items (planned)
- compareToIdeal: Compare to maximum possible stats
- groupByCategory: Organize by category

### Security
- Sandboxed Lua execution environment
- Read-only operations (no write permissions)
- Safe API call wrappers with error handling
- Input validation
- Error boundaries with pcall wrapper

### Documentation
- Complete README.md with usage instructions
- Completion calculation formula documentation
- Statistics reference guide
- Configuration options documentation
- Troubleshooting guide
- Use cases for different player types

## [Unreleased]

### Planned
- Complete collectibles tracking (espers, rages, unique equipment)
- Missable items/events tracking with warnings
- Historical statistics (track changes over time)
- Multi-save comparison (side-by-side analysis)
- Visual charts and graphs
- Detailed boss/enemy defeat tracking
- Step counter and playtime tracking
- Gil spending/earning analysis
- Character usage statistics (battles per character)
- Equipment collection tracking (unique items)
- Colosseum completion tracking
- Dragon tracking (8 legendary dragons)
- Achievement system
- Export to PDF/HTML/JSON formats
- Screenshot generation of statistics dashboards
- Integration with cloud save features
- Custom completion goals (user-defined milestones)
- Playthrough comparison (compare different runs)
- Speed run timer integration
- World of Balance vs World of Ruin separation
