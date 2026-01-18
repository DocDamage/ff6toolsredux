# Changelog

All notable changes to the Party Optimizer plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-15

### Added
- Initial release of Party Optimizer plugin
- Character ranking system with 5 categories:
  - Overall Combat Power
  - Magic Rating
  - Physical Rating
  - Tank Rating
  - Speed Rating
- Optimal party recommendation engine with configurable options
- Automatic role balancing for diverse party composition
- Current party analysis with detailed statistics
- Balance scoring system (0-100 scale)
- Character comparison tool for head-to-head stat analysis
- Role detection system (Tank, Mage, Physical DPS, Fast Attacker, Balanced)
- Weighted rating formulas for accurate character assessment
- Configurable party size (default: 4)
- Configurable minimum level filter
- One-click party composition application
- Safe Mode: Confirmation before updating party
- Menu-driven UI with 5 operations
- Comprehensive error handling with user-friendly messages
- Permission validation for read_save, write_save, ui_display

### Calculations
- Combat Power: Weighted formula considering all stats
- Magic Rating: Magic × 3.0 + MagicDef × 2.0 + MaxMP × 0.5
- Physical Rating: Attack × 3.0 + Defense × 2.0 + Stamina × 1.5 + MaxHP × 0.3
- Tank Rating: Defense × 3.5 + MagicDef × 3.5 + Stamina × 2.0 + MaxHP × 0.8
- Speed Rating: Speed × 4.0 + Stamina × 1.0
- Balance Score: Role Diversity (40%) + Average Power (60%)

### Security
- Sandboxed Lua execution environment
- Permission-based access control
- Confirmation dialog before modifying party
- Input validation for character indices (0-15)
- Error boundaries with pcall wrapper

### Documentation
- Complete README.md with usage instructions
- Detailed operation descriptions with examples
- Configuration options documentation
- Rating calculation formulas explained
- Troubleshooting guide
- Use cases for different player types

## [Unreleased]

### Planned
- Character name lookups from game database
- Equipment optimization recommendations
- Esper matching suggestions (best espers for each character)
- Growth potential analysis (stat gain predictions)
- Boss-specific party recommendations
- Export party builds to JSON
- Import community party builds
- Historical stat tracking across levels
- Multi-party management (World of Ruin support)
- Character ability scoring (Blitz, SwdTech, Tools, etc.)
- Synergy analysis (character combo effectiveness)
- Equipment set recommendations per character
- Level scaling predictions (what stats will be at max level)
- Party cost analysis (for challenge runs with point systems)
