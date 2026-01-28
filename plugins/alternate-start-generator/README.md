# Alternate Start Generator

**Phase:** 6 (Gameplay-Altering Plugins)  
**Batch:** 3  
**Version:** 1.0.0  
**Status:** ✅ Complete  
**Lines of Code:** ~750  

## Overview

The **Alternate Start Generator** plugin enables you to start FF6 from **any point in the game** under custom conditions. Instead of playing through from the beginning, skip to the World of Ruin, create solo character challenges, or prepare optimized starts for speedruns.

### Key Capabilities

- **Skip to World of Ruin** - Begin at the post-apocalypse with available characters
- **Choose Starting Characters** - Select any character combination
- **Set Story Events** - Position yourself at specific plot points
- **Configure Starting Inventory** - Define initial equipment and items
- **Create Challenge Scenarios** - Solo runs, limited parties, specific difficulty setups

## Features

### 8 Preset Starting Scenarios

1. **Skip to World of Ruin** - Begin WoR with Celes and core team
2. **Speedrun - Any%** - Optimized for Any% speedrun (Sabin route)
3. **Speedrun - 100%** - Optimized for 100% speedrun preparation
4. **Celes Solo Run** - Only Celes available (extreme challenge)
5. **Three-Character Challenge** - Exactly 3 characters for tactical play
6. **Boss Rush - WoB** - All 8 WoB characters, skip grind
7. **Boss Rush - WoR** - All 10 WoR characters, maximum power
8. **Balanced Challenge Party** - Mix of physical, magic, and support

### Customization Options

- **World State Selection** - Choose WoB or WoR
- **Story Event Selection** - Start at 8+ major plot points
- **Character Roster** - Selectable from all 14 characters
- **Starting Levels** - Configurable (1-99)
- **Starting Gil** - Custom currency amount
- **Starting Inventory** - Define initial equipment

### Advanced Features

- Create custom scenarios and save them
- Apply multiple scenario changes
- Event flag manipulation
- Character availability control
- Export scenario configurations
- Complete validation and safety checks

## Installation

1. Extract `alternate-start-generator` to plugins/
2. Restart FF6 Save Editor
3. Available in Phase 6 - Experimental
4. Backup saves before use!

## Usage

### Apply Preset Scenario

1. Open a save file
2. Navigate to Alternate Start Generator
3. Click "Apply Preset Scenario"
4. Choose from 8 available presets
5. Confirm application
6. Your save is now positioned at the selected start

### Create Custom Scenario

1. Click "Custom Scenario"
2. Configure:
   - World state (WoB/WoR)
   - Story event
   - Character roster
   - Starting level
   - Starting gil
3. Save configuration
4. Apply custom scenario

### Scenario Details

| Scenario | World | Party Size | Level | Gil | Best For |
|----------|-------|-----------|-------|-----|----------|
| Skip to WoR | WoR | 4 | 30 | 5K | Experience WoR |
| Any% Speedrun | WoB | 3 | 1 | 0 | Speed runs |
| 100% Speedrun | WoB | 3 | 10 | 2K | Speedrun prep |
| Celes Solo | WoB | 1 | 1 | 0 | Solo challenge |
| 3-Char Challenge | WoB | 3 | 1 | 0 | Team play |
| Boss Rush WoB | WoB | 8 | 20 | 10K | Boss practice |
| Boss Rush WoR | WoR | 10 | 50 | 50K | High-level bosses |
| Balanced Party | WoB | 4 | 15 | 3K | Strategic play |

## Story Events

Available starting points:

1. **Game Start** - Terra's flashback (WoB start)
2. **Figaro Castle** - After siege (WoB)
3. **Narshe** - First arrival (WoB)
4. **South Figaro** - After liberation (WoB)
5. **Mt. Kolts** - After passage (WoB)
6. **World of Ruin** - After Kefka (WoR transition)
7. **Floating Island** - Final dungeon prep (WoR)
8. **Kefka's Tower** - Ultimate boss battle (WoR)

## Examples

### Example 1: Skip to World of Ruin

Goal: Experience WoR content without grinding through WoB

1. Open save file
2. Apply "Skip to World of Ruin" preset
3. Celes, Edgar, Locke, Cyan are available
4. Level 30 with 5,000 gil
5. Immediate access to WoR story

### Example 2: Speedrun Practice

Goal: Practice Any% speedrun strategy

1. Apply "Speedrun - Any%" preset
2. Start as Sabin with Edgar and Terra
3. Level 1, no gil
4. Practice optimal route repeatedly

### Example 3: Solo Celes Challenge

Goal: Play through with only Celes (extreme difficulty)

1. Apply "Celes Solo Run" preset
2. Only Celes available
3. Level 1, minimal resources
4. Test character mastery

### Example 4: Boss Rush Training

Goal: Fight major bosses with unlimited resources

1. Apply "Boss Rush - WoB" preset
2. All 8 characters available
3. Level 20, 10,000 gil
4. Skip random encounters, focus on bosses

## Technical Details

### World States

- **WoB** - World of Balance (first half)
  - Terra, Locke, Cyan, Shadow, Edgar, Sabin, Celes available
  - Story-driven progression
  - Early game mechanics

- **WoR** - World of Ruin (second half)
  - All 14 characters eventually available
  - Open world exploration
  - Post-apocalypse setting

### Story Event Flags

Each story event sets specific game flags:
- Character availability
- Location accessibility
- Item availability
- Enemy types
- Dialogue changes
- Cutscene triggers

### Character Roster

All 14 playable characters can be enabled/disabled:
- Terra, Locke, Cyan, Shadow, Edgar, Sabin
- Celes, Strago, Relm, Setzer, Mog, Gau, Gogo, Umaro

### API Requirements

- `flags.setWorldState(state)` - Set WoB/WoR
- `flags.setStoryEvent(eventId)` - Set story progression
- `character.setAvailable(id, available)` - Enable/disable character
- `character.setLevel(id, level)` - Set starting level
- `player.setGil(amount)` - Set starting currency
- `inventory.setItemQuantity(id, qty)` - Initialize inventory

## Safety Considerations

### Backup Before Use

Creating an alternate start modifies:
- Story progression flags
- Character availability
- Starting resources
- Game state position

Always backup before using!

### Dangerous Operations

⚠️ **Warning:** Some combinations can create impossible game states:
- Solo character with no healing items
- World state mismatch with available characters
- Inaccessible story events
- Missing required items

The plugin warns about dangerous combinations.

### What Gets Modified

✏️ Modified:
- World state (WoB/WoR)
- Story event flags
- Character roster
- Character levels
- Starting gil
- Starting inventory

🔒 Not Modified:
- Existing save data (creates new game start)
- Item quantities (unless customized)
- Character names
- Graphics/audio

## Troubleshooting

### "Invalid Scenario" Error
- Check scenario ID exists
- Verify all settings are valid
- Consult preset list

### Characters Unavailable After Start
- Some character combinations are invalid for story events
- WoB scenarios can't include Relm, Setzer, Mog, Gau, Gogo, Umaro
- Check available roster for your story event

### Game Won't Load Modified Save
- The alternate start may have created an incompatible state
- Restore from backup
- Try a different scenario

### Missing Starting Items
- Inventory customization may be incomplete
- Check item IDs are valid
- Verify quantities are correct

## Known Limitations

### Story Event Constraints

- Not all event combinations are valid
- Some characters are locked out of certain events
- Equipment availability varies by event
- Location accessibility limited by progression

### No Retroactive Changes

- Applies to NEW games only
- Can't modify existing saves past initial point
- Use before playing through the start

### Limited Mid-Game Changes

- World state can't toggle mid-playthrough
- Story events are one-way (mostly)
- Character roster changes are major

## FAQ

### Q: Can I use this mid-game?
**A:** No, designed for new game starts. Use at save creation.

### Q: Will this break my playthrough?
**A:** If compatibility is maintained, no. Mismatched settings can cause issues.

### Q: Can I go back to normal after applying?
**A:** Undo by loading a backup. No in-plugin undo.

### Q: What about story/dialogue?
**A:** Cutscenes may reference "missing" events. Some dialogue may be odd.

### Q: Can I play competitively with this?
**A:** Check speedrun.com rules. Most require unmodified seeds.

### Q: Does this affect emulator compatibility?
**A:** Depends on emulator. Test before competition.

## Version History

### v1.0.0 (January 16, 2026) - Initial Release
- 8 preset scenarios
- Custom scenario creation
- Story event selection (8 events)
- Character roster configuration
- Starting level/gil customization
- World state selection (WoB/WoR)
- Event flag manipulation
- Complete documentation

## Credits

**Phase:** 6 Batch 3  
**Developer:** FF6 Save Editor Team  
**Inspiration:** Speedrunning community requests

---

**Status:** Stable v1.0.0  
**Last Updated:** January 16, 2026
