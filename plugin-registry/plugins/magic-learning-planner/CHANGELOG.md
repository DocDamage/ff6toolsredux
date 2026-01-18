# Changelog

All notable changes to the Magic Learning Planner plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-16

### Added
- Initial release of Magic Learning Planner plugin
- Spell progress tracking for individual characters
- Progress grouped by magic school (Black, White, Blue, Red, Special)
- Esper recommendations based on unlearned spells
- AP calculation for spell mastery
- Spell source lookup (find which espers teach spells)
- Complete esper teaching database viewer
- Priority system (partially learned spells ranked first)
- Completion percentage tracking per school
- Simplified spell database (~50 spells)
- Simplified esper database (~7 espers)
- Learning rate display (x1, x2, x5, x10)
- Menu-driven UI with 4 operations
- Search functionality for spell sources
- Comprehensive error handling
- Permission validation (read_save, ui_display)

### Spell Organization
- **Black Magic:** 14 spells (Fire through Ultima)
- **White Magic:** 12 spells (Cure through Holy)
- **Blue Magic:** 10 spells (Scan, Haste, Stop, etc.)
- **Red Magic:** 6 spells (Mute, Sleep, Confuse, etc.)
- **Special:** 5 spells (Quick, Meteor, etc.)

### Esper Database
- Ramuh (Bolt line)
- Ifrit (Fire line)
- Shiva (Ice line)
- Unicorn (Cure, Regen)
- Carbuncle (Reflect, Shell, Protect)
- Bahamut (Flare, Meteor)
- Alexandr (Shell, Protect, Holy)

### Features
- Visual progress indicators (✓ learned, % partially learned, ✗ unstarted)
- AP requirement calculations based on learn rate and current progress
- Top 10 priority spell recommendations
- Spell search with partial name matching
- Read-only plugin (no save modifications)

### Configuration
- showLearnedOnly: Filter view to learned spells
- sortByProgress: Prioritize partially learned spells
- highlightMissing: Emphasize unlearned spells (planned)
- calculateAPRequired: Show AP calculations
- groupBySchool: Organize by magic school

### Security
- Sandboxed Lua execution environment
- Read-only operations (no write permissions)
- Input validation for character indices and spell names
- Error boundaries with pcall wrapper

### Documentation
- Complete README.md with usage instructions
- AP calculation formula explanations
- Magic school organization reference
- Configuration options documentation
- Troubleshooting guide
- Use cases for different player types

## [Unreleased]

### Planned
- Complete esper database with all 27 espers
- Complete spell database with all 54 spells
- Natural spell learning tracking (character-specific spells)
- Equipment-based spell learning (Paladin Shield, etc.)
- Esper assignment recommendations with write permissions
- Multi-character spell progress comparison side-by-side
- Visual progress bars and charts
- Export spell learning plan to CSV/JSON
- Integration with Esper Growth Optimizer plugin
- Spell learning timeline (historical tracking)
- Character-specific spell recommendations based on role
- Boss-specific spell requirements checker
- Missable spell warnings
- Optimal esper rotation calculator
- Spell proficiency tracking (usage statistics)
