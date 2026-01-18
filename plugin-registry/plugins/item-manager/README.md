# Item Manager

A powerful batch inventory management plugin for the FF6 Save Editor that provides automated operations for adding, removing, and managing item quantities.

## Features

- **Batch Operations:** Perform operations on multiple items simultaneously
- **Max Consumables:** Set all potions and consumables to maximum quantity (99)
- **Add Items:** Add specific items by ID with custom quantities
- **Remove Items:** Remove specific items from inventory
- **Set Quantities:** Modify existing item quantities
- **Duplicate Items:** Max out existing item quantities instantly
- **Clear All:** Remove all items with safety confirmation
- **Inventory Summary:** View detailed inventory statistics
- **Safe Mode:** Confirmation dialogs before applying changes
- **Preview Changes:** See what will be modified before committing

## Installation

1. Open FF6 Save Editor (version 3.4.0 or higher)
2. Go to `Tools → Marketplace` or `Community → Marketplace`
3. Search for "Item Manager"
4. Click "Install"

## Usage

### Basic Usage

1. Load a save file in the FF6 Save Editor
2. Go to `Tools → Plugin Manager`
3. Find "Item Manager" in the installed plugins list
4. Click "Run"
5. Select an operation from the menu
6. Follow the prompts to complete the operation

### Available Operations

#### 1. Max All Consumables
Automatically sets all consumable items (potions, ethers, etc.) to maximum quantity (99).

**Use Case:** Quickly stock up on healing items for difficult battles or speedruns.

**Steps:**
1. Select option 1 from the menu
2. Review the preview showing how many items will be modified
3. Confirm to apply changes

#### 2. Add Specific Item
Add a new item to your inventory or increase the quantity of an existing item.

**Use Case:** Add specific items you need for testing or gameplay.

**Steps:**
1. Select option 2 from the menu
2. Enter the item ID (1-999)
3. Enter the quantity to add (1-99)
4. Confirm to apply changes

**Example:**
- Item ID: 5 (Phoenix Down)
- Quantity: 20
- Result: Adds 20 Phoenix Downs to inventory

#### 3. Remove Specific Item
Completely remove an item from your inventory.

**Use Case:** Clean up unwanted items or free inventory slots.

**Steps:**
1. Select option 3 from the menu
2. Enter the item ID to remove
3. Confirm to apply changes

#### 4. Set Item Quantity
Change the quantity of an existing item in your inventory.

**Use Case:** Set specific item counts for testing or gameplay scenarios.

**Steps:**
1. Select option 4 from the menu
2. Enter the item ID to modify
3. Enter the new quantity (0-99)
4. Confirm to apply changes

**Note:** Setting quantity to 0 removes the item.

#### 5. Duplicate Item
Instantly max out the quantity of an existing item to 99.

**Use Case:** Quickly max out rare items you already have.

**Steps:**
1. Select option 5 from the menu
2. Enter the item ID to duplicate
3. Confirm to apply changes

#### 6. Clear All Items
Remove ALL items from your inventory.

**⚠️ WARNING:** This is a destructive operation!

**Use Case:** Starting fresh or testing from empty inventory.

**Steps:**
1. Select option 6 from the menu
2. Read the warning carefully
3. Confirm you want to proceed
4. Review changes preview
5. Confirm again to apply

**Safety:** Requires double confirmation to prevent accidental deletion.

#### 7. View Inventory Summary
Display detailed statistics about your current inventory without making changes.

**Information Shown:**
- Total unique items
- Total item quantity
- Top 10 items by ID

**Use Case:** Quick inventory overview before performing operations.

#### 8. Cancel
Exit the plugin without making any changes.

## Configuration

The plugin supports several configuration options (modify at top of plugin.lua):

```lua
local config = {
    maxQuantity = 99,           -- Maximum item quantity
    safeMode = true,            -- Confirm before applying changes
    showPreview = true,         -- Show preview of changes
    defaultCategory = "all"     -- Default item category filter
}
```

### Configuration Options

- **maxQuantity:** Maximum quantity for items (default: 99)
- **safeMode:** If true, shows confirmation dialogs before changes
- **showPreview:** If true, shows preview of changes before applying
- **defaultCategory:** Default category filter (not yet implemented)

## Item Categories

The plugin organizes items into categories:

- **Consumables:** Potions, Ethers, Phoenix Downs, etc. (IDs 1-10)
- **Weapons:** Swords, Spears, Claws, etc. (IDs 100-105)
- **Armor:** Shields, Helmets, Body armor (IDs 200-205)
- **Relics:** Accessories and special items (IDs 250-255)
- **Key Items:** Story-important items (IDs 300-305)

**Note:** These are simplified categories. Actual game item IDs may vary.

## Permissions

This plugin requires the following permissions:

- **read_save:** Required to read current inventory data
- **write_save:** Required to modify inventory items
- **ui_display:** Required to show operation menu and dialogs

**Note:** This plugin MODIFIES save data. Always backup your save before use!

## Use Cases

### For Speedrunners
- Quickly max consumables for marathon runs
- Add specific items needed for route optimization
- Clear unwanted items to reduce menu navigation

### For Testing
- Add items to test specific game scenarios
- Remove items to verify game behavior
- Set precise item quantities for edge case testing

### For Casual Play
- Stock up on healing items for difficult sections
- Add missed items without replaying
- Organize inventory by removing duplicates

### For Modders
- Test custom item implementations
- Verify item ID mappings
- Batch operations for development workflows

## Compatibility

- **Editor Version:** 3.4.0 or higher
- **Game:** Final Fantasy VI Pixel Remastered
- **Platform:** Windows, macOS, Linux (wherever the editor runs)

## Screenshots

![Item Manager Menu](screenshot.png)
*Item Manager operation menu showing all available batch operations*

## Safety Features

### Confirmation Dialogs
All destructive operations show confirmation dialogs before applying changes.

### Preview Mode
When enabled, shows exactly what will change before committing.

### Safe Mode
Double confirmation for dangerous operations like "Clear All Items".

### Error Handling
Comprehensive error checking for invalid item IDs or quantities.

### Logging
All operations are logged for troubleshooting and audit purposes.

## Troubleshooting

### "Permission Error"
**Cause:** Plugin permissions not granted  
**Solution:** Reinstall the plugin through the Marketplace

### "Invalid item ID"
**Cause:** Item ID outside valid range (1-999)  
**Solution:** Enter a valid item ID between 1 and 999

### "Item not found in inventory"
**Cause:** Trying to modify/remove non-existent item  
**Solution:** Use operation 7 to view current inventory first

### Changes not applied
**Cause:** Canceled confirmation dialog  
**Solution:** Run plugin again and confirm changes when prompted

### Inventory appears unchanged
**Cause:** Save file not reloaded  
**Solution:** Reload the save file in the editor to see changes

## Known Limitations

- Item names show as "Item #XXX" (future versions may include name lookup)
- Maximum 255 inventory slots (game limitation)
- No support for equipped items (use character editor instead)
- Category filtering not yet implemented

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Future Enhancements

Planned features for future versions:
- Item name lookups from game database
- Category-based filtering and operations
- Import/export item lists
- Preset configurations (speedrun kits, starter packs)
- Batch add from text list
- Inventory sorting and organization
- Search by item name

## License

MIT License

Copyright (c) 2026 FF6 Editor Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Support

For issues, questions, or feature requests:

- **GitHub Issues:** [Create an issue](https://github.com/ff6-save-editor/plugin-registry/issues)
- **Discord:** Join the FF6 Editor community Discord
- **Email:** plugins@ff6editor.dev

## Credits

**Author:** FF6 Editor Team  
**Repository:** https://github.com/ff6-save-editor/plugin-registry  
**Editor:** https://github.com/ff6-save-editor/editor
