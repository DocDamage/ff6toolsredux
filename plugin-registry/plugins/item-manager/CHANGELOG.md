# Changelog

All notable changes to the Item Manager plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-15

### Added
- Initial release of Item Manager plugin
- Max All Consumables operation: Set all potions/ethers to maximum quantity (99)
- Add Specific Item operation: Add items by ID with custom quantities
- Remove Specific Item operation: Remove items from inventory
- Set Item Quantity operation: Modify existing item quantities
- Duplicate Item operation: Max out existing item quantities instantly
- Clear All Items operation: Remove all items with double confirmation
- View Inventory Summary operation: Display detailed inventory statistics
- Menu-driven UI with 8 operations
- Safe Mode: Confirmation dialogs before applying changes
- Preview Mode: See changes before committing
- Error handling with user-friendly messages
- Permission validation for read_save, write_save, ui_display
- Comprehensive helper functions (formatItem, countItems, findItemIndex)
- Logging for troubleshooting and audit purposes

### Security
- Sandboxed Lua execution environment
- Permission-based access control
- Double confirmation for destructive operations
- Input validation for item IDs (1-999) and quantities (0-99)
- Error boundaries with pcall wrapper

### Documentation
- Complete README.md with usage instructions
- Configuration options documentation
- Troubleshooting guide
- Use cases for different user types
- Safety features explanation

## [Unreleased]

### Planned
- Item name lookups from game database
- Category-based filtering and operations
- Import/export item lists from CSV/JSON
- Preset configurations (speedrun kits, starter packs)
- Batch add from text list
- Inventory sorting and organization
- Search by item name
- Undo/redo functionality
- Inventory comparison between saves
- Auto-backup before operations
