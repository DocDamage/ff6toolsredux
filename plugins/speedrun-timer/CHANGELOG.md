# Speedrun Timer Integration Changelog

## [1.0.0] - 2026-01-16

### Added
- Initial release of Speedrun Timer Integration plugin
- **Manual split timer** with millisecond precision using Lua os.clock()
- **5 preset speedrun categories**:
  1. **Any%** (15 splits): Standard completion route
  2. **Glitchless** (15 splits): Any% without major glitches
  3. **100% Completion** (10 splits): All characters/espers/items
  4. **Low Level Run** (7 splits): Complete under level 20
  5. **WoR Skip** (5 splits): Skip World of Ruin content
- **Personal Best (PB) tracking** system:
  - Automatic PB detection on run completion
  - PB comparison during active runs
  - Split-by-split delta display (ahead/behind)
  - Persistent storage per category
- **Timer controls**:
  - Start new run with category selection
  - Manual split triggering
  - Pause/Resume functionality
  - Reset with confirmation dialog
  - Automatic pause time exclusion
- **Live timer display** in menu status line
- **Current run viewer** with:
  - Real-time elapsed time
  - Completed splits with cumulative times
  - Segment time calculation
  - PB comparison deltas
  - Next split preview
- **Personal bests viewer** showing all category PBs with dates
- **Run export** to speedrun.com compatible format:
  - Game and category metadata
  - Final time and date
  - Complete split list
  - Segment times
  - Timestamped filenames
- Time formatting (HH:MM:SS.ms)
- Safe API call wrappers with pcall error handling
- Persistent storage using Lua table serialization

### Features
- Millisecond-accurate timing
- Pause handling (excluded from final time)
- Automatic PB detection and notification
- Split-by-split comparison to personal best
- Segment analysis for performance review
- Export functionality for documentation
- Reset protection with confirmation

### Timer Mechanics
- **Start**: Timer begins on category selection
- **Split**: Records time at each checkpoint, advances to next split
- **Final Split**: Stops timer, checks for PB, saves run
- **Pause/Resume**: Tracks pause duration, excludes from elapsed time
- **Reset**: Clears current run, preserves PB data

### Technical Details
- **Lines of Code**: ~620
- **Permissions**: read_save, ui_display, file_io
- **Dependencies**: Lua standard library (os.clock, os.time, os.date)
- **Data Format**: Lua table literal serialization
- **Storage Files**:
  - `speedrun_timer/personal_bests.lua` - PB times by category
  - `speedrun_timer/routes.lua` - Custom route definitions
  - `speedrun_timer/current_run.lua` - Active run state
  - `speedrun_timer/run_*.txt` - Exported runs

### Split Presets
- **Any%/Glitchless**: Narshe Escape, Figaro Castle, South Figaro, Mt. Kolts, Returners Hideout, Phantom Train, Zozo, Opera House, Vector, Magitek Factory, Esper Escape, Floating Continent, WoR Start, Airship, Kefka's Tower
- **100%**: Narshe Escape, All WoB Characters, All WoB Espers, Floating Continent, All WoR Characters, 8 Dragons, All Espers, All Unique Items, Level 99 Party, Kefka's Tower
- **Low Level**: Narshe Escape, Phantom Train, Opera House, Floating Continent, WoR Start, All Dragons (No Level Up), Kefka's Tower
- **WoR Skip**: Narshe Escape, Floating Continent, Airship Skip, Kefka's Tower Entry, Final Boss

### Export Format
```
=== FF6 Speedrun Export ===
Game: Final Fantasy VI
Category: [Category Name]
Date: [YYYY-MM-DD HH:MM:SS]

Final Time: [HH:MM:SS.ms]

=== Splits ===
1. [Split Name]: [Cumulative] (segment: [Segment Time])
...
```

### Known Limitations
- Manual splits only (no auto-split integration)
- Timer runs in plugin context, not during gameplay
- No world record comparison database
- No custom category creation
- No split editing/reordering
- No route planning tool with time estimates
- No live overlay during gameplay

### Notes
- Phase 5, Plugin 3 of 4
- Part of Challenge & Advanced Tools phase
- Intended for route planning and documentation
- Pause time automatically excluded from final time
- PB notification on new records
- Compatible with speedrun.com manual entry format

### Future API Requirements
- Game state detection for auto-splits
- Event hooks for automatic split triggering
- In-game overlay API for live timer display
- Memory reading for trigger point detection
