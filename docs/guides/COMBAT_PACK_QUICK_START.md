# Combat Depth Pack - Quick Start Guide

## Overview
Combat Depth Pack adds dynamic encounter tuning, boss remixing, and AI companion recommendations to your FF6 save editor, with full Lua scripting support for save file manipulation.

## Installation
The Combat Pack is built-in. No additional installation required.

## Using the UI

### 1. Load Your Save File
File → Open → Select your FF6 PR save.json

### 2. Open Combat Depth Pack
Tools → Combat Depth Pack...

### 3. Available Actions

#### Encounter Tuning
- **Zone:** Enter zone name (e.g., "Mt. Kolts", "Narshe")
- **Encounter Rate:** Multiplier for encounter frequency (0.5-2.0)
- **Elite Chance:** Probability of elite encounters (0.0-1.0)
- Click **Apply Encounter Tuning**

#### Boss Remix
- **Affixes:** Comma-separated list of modifiers
  - Available: `enraged`, `arcane_shield`, `glass_cannon`, `regenerating`, `stagger_resist`
- Example: `enraged,arcane_shield`
- Click **Generate Remix Plan**

#### AI Companion Director
- **Profile:** Tactics style (`aggressive`, `defensive`, `balanced`, `support`)
- **Risk:** Risk tolerance (`low`, `normal`, `high`, `aggressive`)
- Click **Save Profile** to generate AI recommendations

#### Smoke Tests
- Click **Run Smoke Tests** to verify all modules working

### 4. Save Changes
File → Save to persist any modifications

## Using the CLI

### Basic Syntax
```bash
ffvi_editor combat-pack --mode <mode> [options]
```

### Examples

#### Run Smoke Tests
```bash
# Without save manipulation
ffvi_editor combat-pack --mode smoke

# With save manipulation
ffvi_editor combat-pack --mode smoke --file mysave.json
```

#### Encounter Tuning
```bash
ffvi_editor combat-pack \
  --mode encounter \
  --zone "Mt. Kolts" \
  --rate 1.2 \
  --elite 0.15 \
  --file mysave.json
```

#### Boss Remix
```bash
ffvi_editor combat-pack \
  --mode boss \
  --affixes "enraged,glass_cannon" \
  --file mysave.json
```

#### Companion Director
```bash
ffvi_editor combat-pack \
  --mode companion \
  --profile aggressive \
  --risk high \
  --file mysave.json
```

### CLI Options
- `--mode` - Operation mode: `encounter`, `boss`, `companion`, `smoke`
- `--file` - (Optional) Save file path for manipulation
- `--zone` - Zone name for encounter mode
- `--rate` - Encounter rate multiplier (default: 1.0)
- `--elite` - Elite encounter chance (default: 0.10)
- `--affixes` - Comma-separated affixes for boss mode
- `--profile` - AI profile for companion mode
- `--risk` - Risk tolerance for companion mode (default: normal)

## Custom Lua Scripts

### Basic Save Manipulation
```lua
-- Get character info
local count = save.getCharacterCount()
local name = save.getCharacterName(0)
local gil = save.getGil()

-- Modify save
save.setCharacterLevel(0, 50)
save.setCharacterHP(0, 9999)
save.setCharacterMP(0, 999)
save.setGil(100000)

-- Log output
save.log("Modified " .. name .. " to level 50")

return {success = true, character = name, gil = 100000}
```

### Using Combat Pack Modules
```lua
local pack = require('plugins.combat-depth-pack.v1_0_core')

-- Generate boss encounter
local boss = pack.BossRemix.generateRemix(
  {name = "Atma", hp = 50000, attack = 1200},
  {affixes = {"enraged", "arcane_shield"}}
)

-- Get AI recommendation
local action = pack.CompanionDirector.recommendAction(
  {hp_status = "critical", threat_level = "high"},
  "aggressive"
)

return {boss = boss, ai_action = action}
```

### Example Scripts
See `plugins/combat_depth_pack_save_examples.lua` for ready-to-use examples:
- `boost_party()` - Level up party to 50
- `give_starting_bonus()` - Add 50k gil
- `emergency_heal()` - Full heal all characters
- `prepare_for_boss_fight()` - Complete boss prep

## Available Save Bindings

### Character Operations
```lua
save.getCharacterCount() → number (total slots, usually 40)
save.getCharacterName(index) → string (0-based index)
save.setCharacterLevel(index, level) → boolean
save.setCharacterHP(index, hp) → boolean
save.setCharacterMP(index, mp) → boolean
```

### Economy
```lua
save.getGil() → number
save.setGil(amount) → boolean
```

### Utility
```lua
save.log(message) -- Print to console/log
```

## Safety & Sandbox

### What's Allowed
- ✓ Standard Lua: table, string, math operations
- ✓ Plugin loading from `plugins/` directory
- ✓ Save data read/write via `save` table
- ✓ Logging via `save.log()`

### What's Blocked
- ✗ File I/O (io library)
- ✗ OS commands (os library)
- ✗ Dynamic code loading (loadstring, dofile)
- ✗ Debug hooks (debug library)
- ✗ Execution longer than 3 seconds (timeout)

## Troubleshooting

### Script Times Out
- Reduce complexity or break into smaller operations
- Check for infinite loops
- Default timeout is 3 seconds

### "Character not initialized" Error
- Character slot may be empty
- Use `save.getCharacterName(index)` to check first
- Valid character indices: 0-15 for playable characters

### Changes Don't Persist
- Save changes are in-memory until you save the file
- UI: File → Save
- CLI: Changes are applied but not auto-saved (add explicit save step)

### Module Not Found
- Ensure plugin files are in `plugins/` directory
- Check file naming: `plugins/combat-depth-pack/v1_0_core.lua`
- Verify package.path includes plugin directories

## Tips & Best Practices

1. **Always Backup** - Keep backups of your save files before experimentation
2. **Test First** - Run `--mode smoke` before complex operations
3. **Start Small** - Test on a few characters before batch operations
4. **Check Results** - Verify changes before saving file
5. **Use Logging** - Add `save.log()` calls to track script progress
6. **Read Errors** - Error messages show file/line for Lua errors

## Need Help?

- **API Reference:** See `COMBAT_PACK_SAVE_HOOKS_REFERENCE.md`
- **Test Results:** See `COMBAT_PACK_E2E_TEST_RESULTS.md`
- **Example Scripts:** Check `plugins/combat_depth_pack_save_examples.lua`
- **Smoke Tests:** Run `combat-pack --mode smoke` to verify setup

## Version Info
- **Combat Pack:** v1.0.0
- **Lua Runtime:** gopher-lua v1.1.0
- **Tested:** January 17, 2026
- **Status:** ✅ Production Ready
