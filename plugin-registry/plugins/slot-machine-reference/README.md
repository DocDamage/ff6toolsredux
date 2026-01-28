# Slot Machine Reference

Read-only reference for Setzer's Slots outcomes with approximate probabilities, effects, and rough damage estimates based on Setzer's stats when available.

## Features
- Lists core slot outcomes with reels, effects, and notes.
- Approximate probability labels for quick risk/reward reading.
- Rough damage estimates using Setzer's current level/magic/vigor when accessible; falls back to defaults.
- Search by name, effect, type, or notes.
- Safe pcall wrapping with read-only permissions (`read_save`, `ui_display`).

## Usage
1. Open the plugin menu.
2. Choose:
   - `All Outcomes` to view the main list.
   - `Search` to filter by keyword.
   - `Help` for guidance and caveats.
3. Exit with `0`.

### Notes
- Probabilities are labeled qualitatively (Very Low/Low/Medium/High).
- Timing and rigging nuances are not simulated; this is a quick lookup and sanity check tool.

## Implementation Notes
- Data is a representative core set of outcomes; extend `SLOTS` as needed.
- Damage uses simple FF6-like physical/magic approximations with safe defaults.
- Menu loop guards errors via `pcall`; dialogs are read-only.
