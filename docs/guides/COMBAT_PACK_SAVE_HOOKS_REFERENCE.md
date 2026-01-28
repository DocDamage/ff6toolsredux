# Combat Depth Pack - Save Hooks Reference

## Overview
Combat Depth Pack now includes full save data manipulation capabilities through Lua bindings, allowing scripts to read and modify FF6 save files directly.

## Implementation Details

### Save Bindings (scripting/runner.go)
The `RunSnippetWithSave` function registers a global `save` table in Lua with these functions:

#### Character Functions
- **`save.getCharacterCount()`** → Returns number of characters in save
- **`save.getCharacterName(index)`** → Returns character name by index (0-based)
- **`save.setCharacterLevel(index, level)`** → Sets character level
- **`save.setCharacterHP(index, hp)`** → Sets character current HP  
- **`save.setCharacterMP(index, mp)`** → Sets character current MP

#### Inventory/Economy Functions
- **`save.getGil()`** → Returns current gil amount
- **`save.setGil(amount)`** → Sets gil amount

#### Utility Functions
- **`save.log(message)`** → Logs message to console with [LUA] prefix

### UI Integration (ui/forms/combat_depth_pack_dialog.go)
- Dialog constructor now accepts `*pr.PR` save parameter
- All button handlers use `RunSnippetWithSave` instead of `RunSnippet`
- Save data is passed from main window's loaded PR instance

### CLI Integration (cli/commands_stub.go)
Added optional `--file` flag to all combat-pack modes:
```bash
# Run with save manipulation
ffvi_editor combat-pack --mode encounter --zone "Mt. Kolts" --rate 1.2 --file mysave.json

# Run without save (dry-run mode)
ffvi_editor combat-pack --mode smoke
```

When `--file` is provided:
1. Save file is loaded via `LoadSaveFile`
2. Scripts execute with save bindings enabled
3. Changes are made in-memory (caller must save)

## Usage Examples

### Example 1: Boost Party via UI
1. Open save file in editor
2. Tools → Combat Depth Pack
3. Click any action button - script has full save access
4. Save file to persist changes

### Example 2: Boost Party via CLI
```bash
ffvi_editor combat-pack --mode smoke --file mysave.json
```

### Example 3: Custom Lua Script
```lua
local examples = require('plugins.combat_depth_pack_save_examples')

-- Heal all party members
local result = examples.emergency_heal()
-- result = {success=true, healed_count=4, characters={"Terra", "Locke", ...}}

-- Give starting bonus
local bonus = examples.give_starting_bonus()
-- bonus = {success=true, previous_gil=5000, bonus=50000, new_gil=55000}
```

### Example 4: Integration with Combat Pack
```lua
-- Prepare for boss fight (heal + AI recommendations)
local prep = require('plugins.combat_depth_pack_save_examples')
local result = prep.prepare_for_boss_fight()
-- Combines: party healing + boss remix + AI companion recommendations
```

## Save Data Structure Notes

### PR.Characters Array
- Fixed array of 40 OrderedMap entries
- Index 0-15: Playable characters
- Use `Get(key)` to read, `Set(key, value)` to write
- Common keys: `name`, `level`, `current_hp`, `current_mp`, `vigor`, `speed`, etc.

### PR.UserData OrderedMap
- Contains global save state
- Common keys: `gil`, `playtime`, `location`, etc.

### Safety
- All bindings validate indices/ranges
- Sandbox prevents file I/O, OS calls, code loading
- 3-second timeout per script execution
- Changes are in-memory until save explicitly called

## Extending Save Bindings

To add new bindings (in `registerSaveBindings`):
```go
L.SetField(saveTable, "yourFunctionName", L.NewFunction(func(L *lua.LState) int {
    // Get arguments: L.CheckNumber(1), L.CheckString(2), etc.
    // Manipulate save: save.Characters[idx].Set(key, value)
    // Return values: L.Push(lua.LNumber(42)); return 1
}))
```

## Testing
- Smoke tests run without save: `combat-pack --mode smoke`
- Smoke tests with save: `combat-pack --mode smoke --file test.json`
- Unit test: `go test ./cli -run TestCombatPackSmokeCLI`
- Manual test: Load save in UI, run Combat Pack actions

## Files Modified
- `scripting/runner.go` - Added RunSnippetWithSave, registerSaveBindings
- `ui/forms/combat_depth_pack_dialog.go` - Added save parameter, use RunSnippetWithSave
- `ui/window.go` - Pass g.pr to dialog constructor
- `cli/commands_stub.go` - Added --file flag, load save, use RunSnippetWithSave
- `plugins/combat_depth_pack_save_examples.lua` - Example scripts demonstrating save manipulation

## Next Steps
1. Add more save bindings (inventory items, magic, espers, party composition)
2. Add save persistence option to CLI (--save flag)
3. Create UI toggle for "auto-save changes"
4. Extend examples with more complex manipulations
