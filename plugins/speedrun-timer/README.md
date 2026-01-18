# Speedrun Timer Integration Plugin

Track split times for speedrun categories with manual/auto splits and personal best comparison for Final Fantasy VI.

## Features

- **Split Timer**: Manual split triggering with precise timing
- **5 Speedrun Categories**: Popular FF6 speedrun categories with preset splits
- **Personal Best Tracking**: Automatic PB detection and comparison
- **Live Timer Display**: Real-time elapsed time in menu
- **Pause/Resume**: Pause timer without losing progress
- **Split Comparison**: Compare current run to personal best
- **Segment Analysis**: View individual segment times
- **Export Functionality**: Export runs to speedrun.com format
- **Reset Handling**: Safe reset with confirmation

## Speedrun Categories

### 1. Any%
Complete game as fast as possible (15 splits)
- Standard route through main story
- Ends at Kefka's Tower completion

### 2. Glitchless
Any% without major glitches (15 splits)
- Same route as Any% but no skip glitches
- Clean play verification

### 3. 100% Completion
All characters, espers, items (10 splits)
- All WoB/WoR characters
- 8 Dragons
- All espers and unique items
- Level 99 party

### 4. Low Level Run
Complete under level 20 (7 splits)
- Optimized route avoiding battles
- Dragon strategy without leveling
- Final boss at low level

### 5. WoR Skip
Skip World of Ruin content (5 splits)
- Airship skip route
- Direct to Kefka's Tower
- Speedrun routing

## Usage

### Starting a Run
1. Select "Start New Run"
2. Choose speedrun category
3. Timer starts immediately

### During Run
- **Split**: Press split after completing each segment
- **Pause/Resume**: Pause timer when needed (pause time excluded)
- **View Current Run**: See splits with PB comparison
- **Reset**: Start over (with confirmation)

### After Run
- Automatic PB detection if new record
- Export run to text file
- View detailed segment analysis

## Split Management

Each category has preset splits for major checkpoints:
- **Any%/Glitchless**: 15 major story checkpoints
- **100%**: Collection milestones (characters, dragons, espers)
- **Low Level**: Reduced splits for optimized routing
- **WoR Skip**: Minimal splits for skip route

## Personal Best Comparison

During runs, see live comparison to PB:
- **Green (-time)**: Ahead of PB
- **Red (+time)**: Behind PB
- Split-by-split delta display

## Timer Features

- **Precision**: Millisecond accuracy using Lua os.clock()
- **Pause Handling**: Pause duration excluded from final time
- **Auto PB Detection**: Automatically saves new personal bests
- **Persistent Storage**: PBs saved across sessions

## File Structure

```
plugins/speedrun-timer/
├── metadata.json          # Plugin metadata
├── plugin.lua             # Main plugin code
├── README.md              # This file
└── CHANGELOG.md           # Version history

speedrun_timer/
├── personal_bests.lua     # PB times by category
├── routes.lua             # Custom route definitions
├── current_run.lua        # Active run state
└── run_*.txt              # Exported run files
```

## Export Format

Exported runs include:
- Game and category information
- Date and final time
- Complete split list with cumulative times
- Segment times for each split
- Generator signature

Compatible with manual entry to speedrun.com.

## API Requirements

- `os.clock()`: High-precision timing (Lua standard library)
- `ReadFile()`, `WriteFile()`: Data persistence

## Permissions

- `read_save`: Read game state for auto-splits (future)
- `ui_display`: Display timer interface
- `file_io`: Persist PBs and export runs

## Notes

- Currently uses manual splits (auto-splits require game state API)
- Timer runs in plugin context (not during actual gameplay)
- Intended for planning and route documentation
- PB times stored per category
- Pause time is automatically excluded from final time

## Timing Guidelines

- Start timer at game start (New Game selected)
- Split at each major checkpoint
- Pause for breaks (pause time not counted)
- Final split on final boss defeat
- Reset if run is abandoned

## Future Enhancements

- Auto-split integration with game state API
- World record comparison database
- Route planning tool with time estimates
- Live timer overlay during gameplay
- Split editing and reordering
- Custom category creation
- Community split presets

## Version

**Version**: 1.0.0  
**Phase**: 5 (Challenge & Advanced Tools)  
**Author**: FF6 Save Editor Plugin System
