# Infinite Resources Mode

**Version:** 1.0.0  
**Category:** Experimental / Quality of Life  
**Phase:** 6, Batch 4  
**Author:** FF6 Editor Plugin System

---

## Overview

The **Infinite Resources Mode** plugin removes all resource management constraints from Final Fantasy VI, allowing you to focus entirely on story and exploration without worrying about running out of items, Gil, or MP. This plugin is perfect for casual playthroughs, testing scenarios, or simply enjoying the game without resource grinding.

### Key Features

- ✅ **Max All Items** - Set every item to 99 quantity
- ✅ **Max Gil** - Set Gil to 9,999,999
- ✅ **Infinite MP** - Max out MP for all characters
- ✅ **Selective Mode** - Choose which resources to make infinite
- ✅ **Auto-Replenish** - Automatically restore consumed resources
- ✅ **One-Click Activation** - Apply full infinite mode instantly
- ✅ **Backup & Restore** - Undo changes if needed
- ✅ **Operation Logging** - Track all modifications

⚠️ **WARNING:** This plugin is **EXPERIMENTAL** and removes all resource management challenge from the game. It's designed for casual play, testing, and experimentation - not recommended for challenge runs or competitive play.

---

## Installation

1. Copy the `infinite-resources-mode` folder to your FF6 Editor `plugins/` directory
2. Restart FF6 Editor or reload plugins
3. The plugin will appear in the **Experimental** category

**Requirements:**
- FF6 Editor v3.4.0 or higher
- Valid FF6 save file
- Permissions: `read_save`, `write_save`, `ui_display`

---

## Usage Guide

### Full Infinite Mode (One-Click)

Apply complete infinite resources with a single command:

```lua
InfiniteResourcesMode.applyFullInfinite()
```

**This will:**
- Set Gil to 9,999,999
- Set all items to 99 quantity
- Max MP for all characters (9,999 max MP)
- Create automatic backup for restoration

**Use Cases:**
- Casual story-focused playthroughs
- Testing specific builds or strategies
- Removing resource grind completely
- Screenshot/video creation with perfect resources

---

### Selective Infinite Mode

Choose exactly which resources to make infinite:

```lua
InfiniteResourcesMode.applySelectiveInfinite({
    gil = true,           -- Max Gil
    consumables = true,   -- Max consumable items only
    equipment = false,    -- Leave equipment as-is
    all_items = false,    -- Max all items (overrides above)
    mp = true,            -- Infinite MP
})
```

**Examples:**

**Infinite Consumables Only (Challenge Lite):**
```lua
InfiniteResourcesMode.applySelectiveInfinite({
    gil = false,
    consumables = true,
    mp = false,
})
```

**Infinite MP Only (Magic Focus):**
```lua
InfiniteResourcesMode.applySelectiveInfinite({
    gil = false,
    consumables = false,
    mp = true,
})
```

**Everything Except Gil (Partial Challenge):**
```lua
InfiniteResourcesMode.applySelectiveInfinite({
    gil = false,
    all_items = true,
    mp = true,
})
```

---

### Auto-Replenish Mode

Enable automatic resource replenishment - consumed items automatically restore to 99:

```lua
-- Enable auto-replenish
InfiniteResourcesMode.enableAutoReplenish()

-- Disable auto-replenish
InfiniteResourcesMode.disableAutoReplenish()

-- Check replenishment status (called automatically)
InfiniteResourcesMode.checkReplenish()
```

**How Auto-Replenish Works:**
- Checks resources every 60 seconds (configurable)
- Restores items that dropped below 99
- Restores MP to maximum for all characters
- Logs all replenishment operations

**When to Use:**
- Long play sessions without micromanagement
- Testing scenarios that consume resources
- Ensuring resources never run out mid-battle

---

### Backup & Restore

The plugin automatically creates backups before making changes:

```lua
-- Restore from backup (undo infinite resources)
InfiniteResourcesMode.restoreBackup()

-- Check if backup exists
local status = InfiniteResourcesMode.getStatus()
if status.has_backup then
    print("Backup available for restoration")
end
```

**What Gets Backed Up:**
- Current Gil amount
- All item quantities (only non-zero items)
- Current and Max MP for all characters
- Timestamp of backup

---

## Status & Monitoring

### Check Current Status

```lua
local status = InfiniteResourcesMode.getStatus()
```

**Returns:**
```lua
{
    enabled = true,              -- Infinite mode active?
    auto_replenish = true,       -- Auto-replenish active?
    has_backup = true,           -- Backup available?
    operations_count = 5,        -- Number of logged operations
}
```

### Display Resources

View current resource state:

```lua
InfiniteResourcesMode.displayResources()
```

**Output Example:**
```
=== Current Resources ===
Gil: 9999999
Items: 180 types owned, 180 at max (99)

Character MP:
  Character 0: 9999 / 9999 MP
  Character 1: 9999 / 9999 MP
  Character 2: 9999 / 9999 MP
  ...
========================
```

### Operation Log

View history of all modifications:

```lua
-- Get last 10 operations
local log = InfiniteResourcesMode.getLog(10)

-- Get last 50 operations
local log = InfiniteResourcesMode.getLog(50)
```

**Log Entry Format:**
```lua
{
    timestamp = 1704672000,
    operation = "full_infinite_applied",
    details = { ... },
}
```

---

## Configuration

Edit `plugin.lua` to customize behavior:

```lua
local CONFIG = {
    MAX_GIL = 9999999,                 -- Maximum Gil value
    MAX_ITEM_QUANTITY = 99,            -- Max items per slot
    MAX_MP_BONUS = 9999,               -- Max MP to set
    
    -- Auto-replenish settings
    AUTO_REPLENISH_ENABLED = false,    -- Enable by default?
    REPLENISH_CHECK_INTERVAL = 60,     -- Check every 60 seconds
    
    -- Logging
    MAX_LOG_ENTRIES = 100,             -- Max log entries to keep
    LOG_OPERATIONS = true,             -- Enable operation logging
}
```

---

## Use Cases & Scenarios

### 1. **Casual Story Playthrough**
**Goal:** Experience the story without resource management.

```lua
InfiniteResourcesMode.applyFullInfinite()
```

**Result:** Focus entirely on narrative, never worry about resources.

---

### 2. **Build Testing**
**Goal:** Test specific character builds with unlimited resources.

```lua
InfiniteResourcesMode.applySelectiveInfinite({
    gil = true,
    all_items = true,
    mp = false,  -- Keep MP limited for testing
})
```

**Result:** Buy any equipment, use any items, test builds freely.

---

### 3. **Magic-Only Playthrough**
**Goal:** Infinite MP for spellcasting focus.

```lua
InfiniteResourcesMode.applySelectiveInfinite({
    mp = true,  -- Only make MP infinite
})
```

**Result:** Cast spells freely without MP management.

---

### 4. **Screenshot/Video Creation**
**Goal:** Perfect resources for content creation.

```lua
InfiniteResourcesMode.applyFullInfinite()
InfiniteResourcesMode.enableAutoReplenish()
```

**Result:** All resources maxed and automatically maintained.

---

### 5. **Challenge Run with QoL**
**Goal:** Challenge difficulty but without item grinding.

```lua
InfiniteResourcesMode.applySelectiveInfinite({
    gil = false,           -- Keep Gil limited
    consumables = true,    -- Infinite potions/ethers
    equipment = false,     -- Must find equipment
    mp = false,            -- Keep MP limited
})
```

**Result:** Tactical challenge remains but no grinding for consumables.

---

## Technical Details

### Resource Modification Process

1. **Backup Creation:** Saves current state before any changes
2. **Confirmation:** Prompts user to confirm (optional)
3. **Gil Modification:** Sets Gil to 9,999,999
4. **Item Modification:** Iterates through all 255 item IDs, sets to 99
5. **MP Modification:** Sets max MP to 9,999 for all 14 characters
6. **Logging:** Records all operations with timestamps
7. **Completion:** Reports success and modified counts

### Auto-Replenish System

- **Check Interval:** 60 seconds (default)
- **Detection:** Scans for items < 99 quantity
- **Restoration:** Sets items back to 99
- **MP Restore:** Sets current MP to max MP
- **Logging:** Records replenishment events

### API Usage

The plugin uses these FF6 Editor API calls:

- `API.getGil()` - Read current Gil
- `API.setGil(amount)` - Set Gil amount
- `API.getItemQuantity(itemId)` - Read item quantity
- `API.setItemQuantity(itemId, quantity)` - Set item quantity
- `API.getCharacterCurrentMP(charId)` - Read current MP
- `API.getCharacterMaxMP(charId)` - Read max MP
- `API.setCharacterCurrentMP(charId, amount)` - Set current MP
- `API.setCharacterMaxMP(charId, amount)` - Set max MP

All API calls are wrapped in `pcall()` for error safety.

---

## Troubleshooting

### Issue: "No backup available to restore"

**Cause:** Trying to restore without applying infinite mode first.

**Solution:** Apply infinite mode (which creates backup) before attempting restoration.

---

### Issue: Items not maxing to 99

**Cause:** Some item IDs may not exist or be invalid.

**Solution:** Check operation log for specific failed item IDs. Some slots may be unused in FF6.

---

### Issue: MP not becoming infinite

**Cause:** Character may not be unlocked yet.

**Solution:** Only unlocked characters can have MP modified. Unlock characters first.

---

### Issue: Auto-replenish not working

**Cause:** Auto-replenish requires periodic calling of `checkReplenish()`.

**Solution:** Ensure the plugin's update loop is calling `checkReplenish()` regularly.

---

### Issue: Changes not persisting after save

**Cause:** Changes need to be saved to the save file.

**Solution:** Save your game after applying infinite resources mode.

---

## Frequently Asked Questions

### Q: Does this work with challenge modes?
**A:** Technically yes, but it defeats the purpose of challenge modes. Not recommended for competitive play or achievements.

### Q: Can I restore my original resources?
**A:** Yes! Use `restoreBackup()` to undo all changes. Backup is created automatically before applying infinite mode.

### Q: Does this affect game balance?
**A:** Absolutely. This removes all resource management challenge. Use only for casual play or testing.

### Q: Can I make only specific items infinite?
**A:** Yes! Use selective infinite mode and specify consumables, equipment, or individual item categories.

### Q: Will this corrupt my save?
**A:** No. The plugin uses safe API calls and creates backups. However, always back up your save file externally as a precaution.

### Q: Can I use this with other plugins?
**A:** Yes! This plugin is compatible with most other plugins. Combine with progression pacing or difficulty modifiers for custom experiences.

### Q: Does auto-replenish work in battle?
**A:** Auto-replenish checks every 60 seconds. It may not activate instantly during battle, but will restore resources shortly after.

### Q: Can I change the max item quantity?
**A:** Yes! Edit `CONFIG.MAX_ITEM_QUANTITY` in `plugin.lua`. 99 is the in-game maximum for FF6.

### Q: What happens if I apply infinite mode twice?
**A:** It's safe but redundant. The plugin checks current state and only modifies what's needed.

### Q: Can I export my infinite resources configuration?
**A:** Not currently implemented, but you can view your operation log to see what was applied.

---

## Known Limitations

1. **Item Slots:** Only existing item slots can be modified. Empty slots remain empty.
2. **Character Availability:** MP can only be modified for unlocked characters.
3. **In-Battle Resources:** Changes apply outside battle. In-battle resources update on next save load.
4. **Auto-Replenish Timing:** 60-second intervals may not catch instant consumption.
5. **Equipment Durability:** FF6 doesn't have equipment durability, so infinite equipment is just quantity.

---

## Safety & Warnings

⚠️ **ALWAYS BACK UP YOUR SAVE FILE** before using this plugin.

⚠️ **NOT RECOMMENDED FOR:**
- Challenge runs (defeats the purpose)
- Speedruns (considered cheating)
- First playthroughs (removes intended experience)
- Competitive play or achievements

✅ **RECOMMENDED FOR:**
- Casual replays focused on story
- Testing builds and strategies
- Removing grind from subsequent playthroughs
- Content creation (screenshots/videos)
- Sandbox experimentation

---

## Changelog

### Version 1.0.0 (January 16, 2026)
- Initial release
- Full infinite mode implementation
- Selective infinite mode
- Auto-replenish system
- Backup and restore functionality
- Operation logging
- Status monitoring

---

## Support & Feedback

For issues, suggestions, or feedback:
- Check the operation log: `InfiniteResourcesMode.getLog()`
- Verify plugin status: `InfiniteResourcesMode.getStatus()`
- Review this README's troubleshooting section
- Report bugs with full operation log attached

---

## Credits

**Developed by:** FF6 Editor Plugin System  
**Plugin Category:** Experimental / Quality of Life  
**Phase 6, Batch 4** - Part of the 44-plugin expansion

**Special Thanks:**
- FF6 modding community
- FF6 Editor development team
- Plugin testing volunteers

---

## License

This plugin is part of the FF6 Editor Plugin System and follows the main project's license terms.

---

**Enjoy infinite resources and stress-free gaming! 🎮✨**
