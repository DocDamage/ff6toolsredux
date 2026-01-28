# Instant Mastery System

**Phase:** 6 (Gameplay-Altering Plugins)  
**Batch:** 3  
**Version:** 1.0.0  
**Status:** ✅ Complete  
**Lines of Code:** ~580  

## Overview

The **Instant Mastery System** is a sandbox mode plugin that enables complete gameplay transformation. Instead of grinding for levels, unlocking spells, or collecting items across a long playthrough, this plugin allows you to unlock **everything instantly** with a single operation or selective choices.

This is the ultimate **"skip the grind"** tool for players who want to:
- Experience FF6's gameplay without progression mechanics
- Test build combinations and strategies instantly
- Run challenge modes with perfect gear
- Create speedrun practice saves
- Experiment with different party compositions

### ⚠️ Warning

This plugin **fundamentally alters your save file** and enables features outside normal FF6 gameplay. It transforms your save into a sandbox environment. **Always backup your save before using!**

## Features

### Full Mastery
The most powerful mode - unlock everything with one click:
- ✅ Max all character levels to 99
- ✅ Max all character stats (HP, MP, Vigor, Speed, Stamina, Magic, Defense, etc.)
- ✅ Unlock all ~50 spells for all characters
- ✅ Learn all 384 Rages (Gau)
- ✅ Learn all 24 Lores (Strago)
- ✅ Learn all 8 Dances (Mog)
- ✅ Learn all 8 Blitzes (Sabin)
- ✅ Obtain all 26 espers
- ✅ Max all items to 99 quantity
- ✅ Set Gil to 9,999,999

### Selective Mastery
Choose exactly which systems to unlock:
- Select individual mastery categories
- Apply only what you need
- Keep some progression mechanics if desired
- Supports partial unlocking

### Quick Presets
Pre-configured mastery profiles for common scenarios:
- **Combat Ready:** Max stats + unlock spells + get espers (for battle testing)
- **Collection Complete:** All abilities + all items + all espers (for item testing)
- **Completionist:** Everything unlocked (full sandbox)

### Safety Features
- Save state backup before any modification
- Confirmation dialogs for all operations
- Undo functionality to restore previous state
- Operation logging with timestamps
- Detailed success/failure reporting

## Installation

1. Extract `instant-mastery-system` folder to your `plugins/` directory
2. Restart the FF6 Save Editor
3. The plugin will appear in the "Phase 6 - Experimental" section
4. Backup your save file first!

## Usage Guide

### Basic Usage (Full Mastery)

1. **Open a save file** in the FF6 Save Editor
2. **Backup your save** (File → Backup or Ctrl+B)
3. **Navigate to Instant Mastery System** in the Phase 6 section
4. **Click "Apply Full Mastery"**
5. **Confirm the operation** in the safety dialog
6. **Wait for completion** - see operation log for progress
7. Your save is now in sandbox mode!

### Selective Mastery

1. **Open a save file** in the FF6 Save Editor
2. **Navigate to Instant Mastery System**
3. **Click "Selective Mastery"**
4. **Choose categories to unlock:**
   - Stats → Max all character levels and stats
   - Spells → Unlock all magic for all characters
   - Rages → Learn all 384 Rages for Gau
   - Lores → Learn all 24 Lores for Strago
   - Dances → Learn all 8 Dances for Mog
   - Blitzes → Learn all 8 Blitzes for Sabin
   - Espers → Obtain all 26 espers
   - Inventory → Max all items to 99
   - Gil → Set to 9,999,999
5. **Confirm selection** and apply
6. Only chosen systems are unlocked

### Using Presets

1. **Navigate to Instant Mastery System**
2. **Click "Quick Presets"**
3. **Choose a preset:**
   - **Combat Ready** - For testing battle mechanics
   - **Collection Complete** - For testing items/equipment
   - **Completionist** - Everything (full sandbox)
4. **Confirm** and apply
5. Preset configuration is instantly applied

### Undoing Changes

1. **Navigate to Instant Mastery System**
2. **Click "Undo Last Operation"**
3. Your save is restored to the state before the last mastery operation
4. Works only for the most recent operation

## Examples

### Example 1: Build Testing
**Goal:** Test a specific character build without grinding

1. Open a new save at the start of the game
2. Apply "Combat Ready" preset
3. All characters are now level 99 with max stats and all spells
4. Equip different gear combinations
5. Test party synergies
6. Go back to normal save when done

### Example 2: Speedrun Practice
**Goal:** Practice speedrun route with perfect gear

1. Apply Full Mastery to a save
2. Jump to a specific story point (World of Ruin)
3. All optimal gear and espers are available
4. Practice the speedrun section repeatedly
5. No time spent on grinding or collecting

### Example 3: Challenge Mode Setup
**Goal:** Create a save for Natural Magic Block challenge

1. Open a save and apply Selective Mastery
2. Unlock only: Stats, Espers, Items, Gil
3. Don't unlock Spells (challenge rule: no magic)
4. Now you have perfect gear but no magic access
5. Run the challenge with proper handicaps

### Example 4: All-Spell Party
**Goal:** Every character can cast all spells

1. Apply Selective Mastery
2. Choose only "Spells" category
3. Every character instantly learns all ~50 spells
4. Enjoy an all-magic focused playthrough

## Technical Details

### Categories and Unlock Amounts

| Category | Amount | System | Notes |
|----------|--------|--------|-------|
| Character Levels | 14 chars × 99 | Progression | All 14 playable characters to max level |
| Hit Points | 14 chars × 9999 | Vitality | Max HP for combat tanking |
| Magic Points | 14 chars × 9999 | Mana | Unlimited spellcasting |
| Stats | 14 chars × 9 stats | Abilities | Vigor, Speed, Stamina, Magic, Defense, M.Defense, Evade, M.Evade |
| Spells | 14 chars × 50 spells | Magic | All magic available to every character |
| Rages | 1 char × 384 rages | Gau | All enemy rages for Gau |
| Lores | 1 char × 24 lores | Strago | All Blue Magic for Strago |
| Dances | 1 char × 8 dances | Mog | All locations dances for Mog |
| Blitzes | 1 char × 8 blitzes | Sabin | All Bushido techniques for Sabin |
| Espers | 26 espers | Equipment | All summoned espers |
| Items | ~256 items × 99 | Inventory | All equipment and consumables maxed |
| Gil | 9,999,999 | Currency | Maximum possible gil |

### API Requirements

This plugin requires the following APIs:
- `character.setLevel(id, level)` - Set character level
- `character.setHP(id, hp)` - Set current HP
- `character.setMP(id, mp)` - Set current MP
- `character.setStat(id, stat_name, value)` - Set character stats
- `character.learnSpell(id, spell_id)` - Learn spell
- `character.learnRage(id, rage_id)` - Learn Rage (Gau)
- `character.learnLore(id, lore_id)` - Learn Lore (Strago)
- `character.learnDance(id, dance_id)` - Learn Dance (Mog)
- `character.learnBlitz(id, blitz_id)` - Learn Blitz (Sabin)
- `inventory.addEsper(id, count)` - Add esper
- `inventory.setItemQuantity(id, quantity)` - Set item quantity
- `player.setGil(amount)` - Set gil amount

### Permissions Required

```json
"permissions": ["read_save", "write_save", "ui_display"]
```

- **read_save** - Read character and inventory data
- **write_save** - Modify all save file data
- **ui_display** - Show plugin UI and menus

### Performance

- **Full Mastery Application:** ~2-5 seconds (depends on system)
- **Selective Mastery:** ~1-3 seconds
- **Memory Usage:** <10MB during operation
- **Backup Creation:** ~500ms
- **Undo Operation:** ~1-2 seconds

## Safety Considerations

### Backup Before Use
Always backup your save file before using Instant Mastery:
1. File → Backup Save (or Ctrl+B)
2. This creates an automatic backup in `backups/` folder
3. You can restore it if something goes wrong

### Restore Original Save
If you want to go back to your original playthrough:
1. Click "Undo Last Operation"
2. Or restore from the backup folder manually

### What Gets Modified

The following save data is modified:
- ✏️ All character levels
- ✏️ All character stats (HP, MP, attributes)
- ✏️ All character spells
- ✏️ All character special abilities
- ✏️ Inventory items (quantities)
- ✏️ Esper collection
- ✏️ Gil amount

**Not Modified:**
- Story/event flags (game progression location stays same)
- Party composition (you choose which characters to use)
- Item locations (they're unlocked, not physically obtained)

### Phase 6 Experimental Features

**Why is this Experimental?**
- This plugin modifies fundamental game systems
- It bypasses normal progression completely
- It enables non-standard gameplay
- API support may be incomplete or unstable
- Gameplay balance is intentionally broken

**Safety Warnings:**
- ⚠️ Your save is transformed to sandbox mode
- ⚠️ Some features may not work correctly
- ⚠️ Game balance is intentionally removed
- ⚠️ Not recommended for first playthroughs
- ⚠️ Always backup before using!

## Troubleshooting

### "Operation Failed" Message
**Problem:** Mastery operation didn't complete
**Solution:**
1. Check that you have a valid save file open
2. Verify your backup is created (see log)
3. Try "Undo Last Operation" first
4. Close and reopen the plugin
5. Report the specific operation that failed

### Stats Didn't Change
**Problem:** Clicked Apply but nothing happened
**Solution:**
1. Check the operation log for errors
2. Verify the save file has write permissions
3. Try closing and reopening the editor
4. Check if API is properly initialized

### Can't Undo Operation
**Problem:** Undo button is disabled
**Solution:**
1. Undo only works for the most recent operation
2. Multiple operations can't be undone in sequence
3. Use your external backup (copy from backups folder)
4. Restore from auto-backup if available

### Game Crashes After Mastery
**Problem:** Game crashes when loading maxed-out save
**Solution:**
1. The game may have limits on stat values
2. Try undoing the operation
3. Use Selective Mastery instead (don't max all stats)
4. Report the specific crash scenario

## Known Limitations

### Maximum Values
Some systems may have internal limits below "max":
- Stats may cap below 99 in calculations
- Item quantities may cap at 99
- Character levels might have gameplay limits
- See your FF6 system documentation

### API Dependency
This plugin's functionality depends on core engine APIs:
- If APIs aren't available, operations fail gracefully
- Some operations may be read-only (no write support)
- Specific character mechanics may not be supported

### Selective Mastery
- Some operations are bundled (can't separate them)
- Some systems depend on others (e.g., high stats need high HP)
- Check operation log for dependency notes

## FAQ

### Q: Will this save work on console emulators?
**A:** Depends on the emulator and save format. The editor uses standard FF6 save format, but some modifications might not be compatible with all emulators. Test first!

### Q: Can I undo multiple operations?
**A:** No, only the most recent operation can be undone. Use external backups if you need to go further back.

### Q: Does this affect story/cutscenes?
**A:** No, story progression is unchanged. You'll still see the same cutscenes. You just have all gear/abilities available from the start.

### Q: Can I partially undo (only remove some unlocks)?
**A:** No, undo restores the entire previous save state. Use Selective Mastery to get partial unlocks instead.

### Q: Is there a way to reset to Normal mode?
**A:** Close without saving and reopen your backup, or use Undo. This plugin doesn't have a "normalize" function.

### Q: Can I use this with other Phase 6 plugins?
**A:** Yes, they're compatible. You can combine Instant Mastery with other gameplay-altering plugins for unique effects.

### Q: What if I accidentally delete a character?
**A:** Use Undo to restore the previous save state, or restore from backup.

## Version History

### v1.0.0 (January 16, 2026) - Initial Release
- Full Mastery mode (unlock everything)
- Selective Mastery (choose categories)
- Quick Presets (pre-configured options)
- Safety backup system
- Undo functionality
- Operation logging
- Complete documentation

## Recommended Use Cases

### ✅ Good Use Cases
- Testing specific builds and strategies
- Practicing speedruns with perfect gear
- Running challenge modes with handicaps
- Experimenting with character combinations
- Playing through story without progression grind
- Testing mod combinations
- Screenshot/video creation

### ❌ Not Recommended For
- First playthrough (use normal mode)
- Learning game mechanics (grind teaches strategy)
- Competitive speedruns (uses sandbox mode)
- Achievements requiring normal progression
- Authentic retro experience

## Permissions and Safety

**Permission Justification:**
- **read_save** - Needs to read all character/inventory data to check current state
- **write_save** - Needs to write all modifications (levels, spells, items, etc.)
- **ui_display** - Needs to show menus, confirmations, and operation logs

**Safety Model:**
- Confirmation dialogs prevent accidental use
- Backup creation before modifications
- Operation logging for audit trail
- Undo functionality for recent changes
- Explicit user confirmation required

## Future Enhancements

Potential improvements for future versions:
- 🔮 Partial stat maxing (e.g., max only Defense)
- 🔮 Experience curves (gradual stat progression)
- 🔮 Preset customization
- 🔮 Batch undo (multiple operations)
- 🔮 Export/import mastery profiles
- 🔮 Balancing profiles (e.g., slight handicap)
- 🔮 Integration with other Phase 6 plugins
- 🔮 GUI preset builder

## Credits

- **Developer:** FF6 Save Editor Team
- **Phase:** 6 Batch 3
- **Inspiration:** Community requests for sandbox testing
- **Compatibility:** FF6 Save Format Specification v1.0

## Support and Feedback

- **Report Bugs:** GitHub Issues with save file details
- **Feature Requests:** GitHub Discussions
- **Plugin Development:** See Plugin Development Guide
- **Community:** Join Discord for discussion

---

**Status:** Stable Release v1.0.0  
**Last Updated:** January 16, 2026  
**Maintenance:** Active  
