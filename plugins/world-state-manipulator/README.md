# World State Manipulator Plugin

**Version:** 1.0.0  
**Category:** Experimental  
**Author:** FF6 Plugin Team

⚠️ **EXTREMELY EXPERIMENTAL - HIGH RISK OF SAVE CORRUPTION**

## Overview

The **World State Manipulator** is an advanced plugin that gives you complete control over FF6's world state and story progression. Toggle between World of Balance and World of Ruin, manipulate event flags, skip story sequences, or create custom game states for testing and experimentation.

## Features

### World State Control
- **Toggle WoB ↔ WoR** - Switch between worlds instantly
- **Set Specific World** - Force World of Balance or World of Ruin
- **World Detection** - Automatically detect current world state

### Event Flag System
- **36+ Event Flags** - Major story events, character recruitment, dungeon completions
- **Enable/Disable Events** - Set or clear individual event flags
- **Bulk Operations** - Complete entire story arcs at once
- **Dangerous Mode** - Protection for critical flags that can corrupt saves

### Story Progression Control
- **Complete WoB Story** - Mark all WoB events as complete
- **Recruit All Characters** - Flag all 13 WoR characters as recruited
- **Clear All Events** - Reset to new game state
- **Dungeon Completions** - Mark major dungeons as cleared

### World State Presets
- **6 Pre-configured States** - Quick jump to major story points
- **Fresh Starts** - New game in WoB or WoR
- **Mid-game States** - Skip to specific story points
- **Endgame Access** - All dungeons and characters available

### Safety Features
- **Automatic Backups** - Save state before modifications
- **Dangerous Mode Lock** - Critical operations require explicit unlock
- **Operation Logging** - Track all state changes
- **Restore System** - Undo all modifications

### Analysis & Export
- **World State Display** - View current world and event completion
- **Story Completion %** - Track overall progress
- **Configuration Export** - Save state details to text

## Installation

1. Copy `world-state-manipulator` folder to `plugins` directory
2. Restart FF6 Save Editor
3. Load save file
4. Access via Plugin Menu → World State Manipulator

⚠️ **BACKUP YOUR SAVE FIRST!** This plugin can corrupt save files!

## Quick Start Guide

### Example 1: Toggle Between Worlds
```lua
-- Switch from WoB to WoR (or vice versa)
toggleWorld()
```

### Example 2: Jump to World of Ruin Start
```lua
-- Fresh WoR start (Celes on Solitary Island)
applyWorldPreset("fresh_wor")
```

### Example 3: Skip to End Game
```lua
-- All characters recruited, all dungeons accessible
applyWorldPreset("endgame")
```

### Example 4: Complete WoB Story
```lua
-- Mark all WoB events as complete
completeWoBStory()
```

### Example 5: Enable Specific Event
```lua
-- Enable Opera House event
enableEvent("opera_house")
```

## Function Reference

### World State Functions

#### `getCurrentWorld()`
Get current world state (0 = WoB, 1 = WoR).

**Returns:** `number` - World state code

**Example:**
```lua
local world = getCurrentWorld()
if world == 0 then
    print("World of Balance")
else
    print("World of Ruin")
end
```

#### `setWorldState(world_state)`
Set world state explicitly.

**Parameters:**
- `world_state` (number) - 0 for WoB, 1 for WoR

**Returns:** `boolean` - Success status

**Example:**
```lua
setWorldState(0)  -- Force World of Balance
setWorldState(1)  -- Force World of Ruin
```

#### `toggleWorld()`
Toggle between WoB and WoR.

**Returns:** `boolean` - Success status

**Example:**
```lua
toggleWorld()  -- WoB→WoR or WoR→WoB
```

#### `setWorldOfBalance()`
Set world to WoB (convenience function).

**Returns:** `boolean` - Success status

**Example:**
```lua
setWorldOfBalance()
```

#### `setWorldOfRuin()`
Set world to WoR (convenience function).

**Returns:** `boolean` - Success status

**Example:**
```lua
setWorldOfRuin()
```

---

### Event Flag Functions

#### `getEventFlag(flag_name)`
Get status of an event flag.

**Parameters:**
- `flag_name` (string) - Event flag identifier

**Returns:** `boolean` - Flag status (true = completed)

**Example:**
```lua
local completed = getEventFlag("opera_house")
if completed then
    print("Opera House event complete")
end
```

#### `setEventFlag(flag_name, value)`
Set event flag to specific value.

**Parameters:**
- `flag_name` (string) - Event flag identifier
- `value` (boolean) - New flag value

**Returns:** `boolean` - Success status

**Example:**
```lua
setEventFlag("opera_house", true)   -- Mark as complete
setEventFlag("opera_house", false)  -- Mark as incomplete
```

#### `enableEvent(flag_name)`
Enable event (set flag to true).

**Parameters:**
- `flag_name` (string) - Event flag identifier

**Returns:** `boolean` - Success status

**Example:**
```lua
enableEvent("magitek_factory")  -- Mark Magitek Factory as complete
```

#### `disableEvent(flag_name)`
Disable event (set flag to false).

**Parameters:**
- `flag_name` (string) - Event flag identifier

**Returns:** `boolean` - Success status

**Example:**
```lua
disableEvent("opera_house")  -- Reset Opera House event
```

---

### Bulk Event Operations

#### `completeWoBStory()`
Mark all major WoB story events as complete (16 events).

**Returns:** `number` - Count of events completed

**Example:**
```lua
local count = completeWoBStory()
print("Completed " .. count .. " WoB events")
```

#### `recruitAllCharacters()`
Flag all 13 WoR characters as recruited.

**Returns:** `number` - Count of characters recruited

**Example:**
```lua
local count = recruitAllCharacters()
print("Recruited " .. count .. " characters")
```

#### `clearAllEvents()`
Clear all event flags (reset to new game state). **Requires dangerous mode.**

**Returns:** `number` - Count of events cleared (or `false` if blocked)

**Example:**
```lua
enableDangerousMode()
clearAllEvents()  -- Reset all events
```

---

### World State Presets

#### `applyWorldPreset(preset_name)`
Apply a pre-configured world state preset.

**Parameters:**
- `preset_name` (string) - Preset identifier

**Available Presets:**
- `"fresh_wob"` - Fresh start in World of Balance (new game)
- `"fresh_wor"` - Fresh start in World of Ruin (Celes on island)
- `"end_wob"` - End of WoB (all events complete, pre-cataclysm)
- `"early_wor"` - Early WoR (Celes found, characters recruitable)
- `"late_wor"` - Late WoR (all characters recruited)
- `"endgame"` - End game (all characters + dungeons accessible)

**Returns:** `boolean` - Success status

**Example:**
```lua
applyWorldPreset("late_wor")  -- Skip to late WoR with all characters
```

#### `listWorldPresets()`
List all available world state presets.

**Returns:** `table` - Array of preset information

**Example:**
```lua
listWorldPresets()
-- Prints: Available World State Presets with descriptions
```

---

### Dangerous Mode Functions

#### `enableDangerousMode()`
Enable dangerous operations (critical flag modifications, dual world access).

**Returns:** `boolean` - Success status

**Example:**
```lua
enableDangerousMode()
-- Now clearAllEvents() and enableDualWorldAccess() are available
```

#### `disableDangerousMode()`
Disable dangerous operations (restore safety).

**Returns:** `boolean` - Success status

**Example:**
```lua
disableDangerousMode()  -- Block critical operations
```

---

### Experimental Functions

#### `enableDualWorldAccess()`
Enable simultaneous access to WoB and WoR. **EXTREMELY EXPERIMENTAL!** **Requires dangerous mode.**

**Returns:** `boolean` - Success status (or `false` if blocked)

**Example:**
```lua
enableDangerousMode()
enableDualWorldAccess()
-- WARNING: High risk of save corruption and crashes!
```

---

### Analysis Functions

#### `getWorldStateSummary()`
Get detailed world state information.

**Returns:** `table` - World state data
- `current_world` - World state code (0 or 1)
- `world_name` - "World of Balance" or "World of Ruin"
- `dangerous_mode` - Dangerous operations enabled
- `event_status` - Array of event completion status
- `story_completion` - Percentage of events complete (0-100)

**Example:**
```lua
local summary = getWorldStateSummary()
print("Current world: " .. summary.world_name)
print("Story completion: " .. summary.story_completion .. "%")
```

#### `displayWorldState()`
Display formatted world state and event status.

**Returns:** `table` - World state summary (same as `getWorldStateSummary()`)

**Example:**
```lua
displayWorldState()
-- Output:
-- === World State Summary ===
-- Current World: World of Ruin
-- Dangerous Mode: DISABLED
-- Story Completion: 65.2%
-- --- Major Events ---
--   [✓] Opera House Performance
--   [✗] Floating Continent
--   ...
```

---

### Backup & Restore Functions

#### `restoreBackup()`
Restore world state from backup (undo all modifications).

**Returns:** `boolean` - Success status

**Example:**
```lua
restoreBackup()  -- Undo all world state changes
```

---

### Export Functions

#### `exportWorldConfig()`
Export world state configuration to formatted text.

**Returns:** `string` - Formatted configuration text

**Example:**
```lua
local config = exportWorldConfig()
-- Exports: current world, completed events, pending events, operation log
```

---

## Event Flag Reference

### World State
- `world_state` - Overall world state (CRITICAL)

### Major Story Events
- `narshe_battle` - Narshe Battle Complete
- `figaro_castle` - Figaro Castle Events
- `south_figaro` - South Figaro Infiltration
- `mt_kolts` - Mt. Kolts Crossed
- `returner_hideout` - Returner Hideout
- `lethe_river` - Lethe River Raft
- `imperial_camp` - Imperial Camp Raid
- `doma_siege` - Doma Siege
- `phantom_train` - Phantom Train
- `baren_falls` - Baren Falls
- `opera_house` - Opera House Performance
- `vector` - Vector Infiltration
- `magitek_factory` - Magitek Factory
- `esper_gate` - Esper Gate Opening
- `thamasa` - Thamasa Events
- `floating_continent` - Floating Continent
- `cataclysm` - World Cataclysm WoB→WoR (CRITICAL)
- `solitary_island` - Solitary Island (WoR Start)

### Character Recruitment (WoR)
- `recruited_celes` - Recruited: Celes
- `recruited_edgar` - Recruited: Edgar
- `recruited_sabin` - Recruited: Sabin
- `recruited_terra` - Recruited: Terra
- `recruited_locke` - Recruited: Locke
- `recruited_cyan` - Recruited: Cyan
- `recruited_strago` - Recruited: Strago
- `recruited_relm` - Recruited: Relm
- `recruited_setzer` - Recruited: Setzer
- `recruited_mog` - Recruited: Mog
- `recruited_gau` - Recruited: Gau
- `recruited_gogo` - Recruited: Gogo
- `recruited_umaro` - Recruited: Umaro

### Major Dungeons
- `phoenix_cave` - Phoenix Cave Complete
- `ancient_castle` - Ancient Castle Complete
- `fanatics_tower` - Fanatic's Tower Complete
- `kefka_tower` - Kefka's Tower Access

## World State Preset Details

### Fresh Starts
- **fresh_wob** - New game in WoB (all events cleared)
- **fresh_wor** - New game in WoR (Celes on island, cataclysm complete)

### Story Progression Points
- **end_wob** - All WoB events complete, ready for Floating Continent
- **early_wor** - Just entered WoR, Celes found, characters recruitable
- **late_wor** - All 13 characters recruited, dungeons still locked
- **endgame** - All characters + all major dungeons accessible

## Use Cases

### Speedrunning & Testing
```lua
-- Jump to specific story points for route testing
applyWorldPreset("end_wob")  -- Test Floating Continent strategies
```

### Character Recruitment Testing
```lua
-- Test character-specific content without full playthrough
setWorldOfRuin()
recruitAllCharacters()
```

### Skip Story Sequences
```lua
-- Skip early game, start in mid-WoB
setWorldOfBalance()
completeWoBStory()
```

### Event Flag Debugging
```lua
-- Check which events are blocking progress
displayWorldState()
-- Manually enable missing events
enableEvent("opera_house")
```

### World State Experiments
```lua
-- Try playing WoB areas in WoR state (experimental)
setWorldOfRuin()
-- Visit WoB locations (may cause glitches)
```

## Advanced Techniques

### Custom Story Progression
```lua
-- Skip to specific event without completing earlier ones
enableEvent("thamasa")
enableEvent("recruited_strago")
enableEvent("recruited_relm")
```

### Conditional Event Setting
```lua
-- Enable events based on conditions
local world = getCurrentWorld()
if world == 1 then
    recruitAllCharacters()  -- Only in WoR
end
```

### Event Flag Combinations
```lua
-- Set up specific game states
setWorldOfBalance()
enableEvent("opera_house")
enableEvent("magitek_factory")
-- Now opera is complete but later events aren't
```

## Safety Features

### Automatic Backups
- Backup created before first modification
- Stores: world state, all event flags
- One backup per session

### Dangerous Mode Protection
Critical operations blocked unless dangerous mode enabled:
- Modifying `world_state` flag
- Modifying `cataclysm` flag
- `clearAllEvents()` function
- `enableDualWorldAccess()` function

### Enable Dangerous Operations
```lua
-- Explicitly unlock dangerous operations
enableDangerousMode()
-- Now critical operations available
clearAllEvents()  -- This now works
disableDangerousMode()  -- Re-enable safety
```

### Operation Logging
- All modifications logged with timestamps
- Last 100 operations tracked
- Export logs for debugging

## Warnings & Limitations

⚠️ **CRITICAL WARNINGS**

1. **Save Corruption Risk** - Improper flag combinations can corrupt saves
2. **Softlock Potential** - Wrong world state + event flags = game-breaking bugs
3. **Vanilla Incompatibility** - Modified saves won't work in vanilla FF6
4. **Dual World Access** - Extremely unstable, guaranteed instability
5. **Event Dependencies** - Some events require others to be set first

### Known Issues

**Issue:** Toggling world mid-game causes glitches  
**Why:** Game expects specific event flags for each world  
**Solution:** Use presets instead of raw world toggling

**Issue:** Character recruitment flags don't add characters  
**Why:** Flags only track recruitment status, don't modify roster  
**Solution:** Use Character Roster Editor plugin in combination

**Issue:** Dual world access crashes game  
**Why:** Game engine not designed for this  
**Solution:** Don't use this feature (seriously!)

**Issue:** Clearing events causes softlock  
**Why:** Game needs certain events to progress  
**Solution:** Use `fresh_wob` or `fresh_wor` presets instead

### Dangerous Operations

Operations that **require dangerous mode**:
- Modifying critical flags (`world_state`, `cataclysm`)
- `clearAllEvents()` - Can softlock game
- `enableDualWorldAccess()` - Can crash game

Operations that are **generally safe**:
- Toggling between worlds (if at appropriate story point)
- Enabling events (moving forward in story)
- Disabling non-critical events
- Using presets

## Best Practices

1. **Always backup save** before using plugin
2. **Use presets** instead of manual flag setting
3. **Enable dangerous mode** only when necessary
4. **Export configuration** before risky operations
5. **Test in separate save** before modifying main save
6. **Don't use dual world access** (seriously, don't)

## Troubleshooting

**Problem:** World toggle causes black screen  
**Solution:** Use preset to properly set world + required events

**Problem:** Characters missing after world toggle  
**Solution:** Enable recruitment flags or use Character Roster Editor

**Problem:** Can't access certain areas  
**Solution:** Enable prerequisite events (check event flag list)

**Problem:** Dangerous operation blocked  
**Solution:** Call `enableDangerousMode()` first (and backup your save!)

## Technical Notes

- **Event Storage:** Flags stored as individual bits in save file
- **World State:** Single value (0 or 1) determines world
- **Flag Dependencies:** Many events depend on earlier events being complete
- **Backup System:** Stores complete snapshot of world state and flags

## Version History

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Support

- **Documentation:** See README.md (this file)
- **Presets:** Use `listWorldPresets()` for available configurations
- **Flags:** See Event Flag Reference section
- **Safety:** Always enable dangerous mode explicitly for critical operations

## License

Part of FF6 Save Editor plugin ecosystem. See main editor license.

---

**Plugin Version:** 1.0.0  
**Last Updated:** 2026-01-16  
**Estimated Lines of Code:** ~680 LOC  
**Documentation:** ~7,000 words

⚠️ **USE AT YOUR OWN RISK - BACKUP YOUR SAVE FILES!**
