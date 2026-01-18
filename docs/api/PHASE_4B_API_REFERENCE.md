# Phase 4B Plugin API Reference

**Version:** 1.0  
**Complete API Documentation**

---

## API Overview

The Plugin API (`editor` global) provides safe access to FF6 Save Editor functionality.

| Category | Functions |
|----------|-----------|
| **Characters** | getCharacter, setCharacter |
| **Inventory** | getInventory, setInventory |
| **Party** | getParty, setParty |
| **Equipment** | getEquipment, setEquipment |
| **Queries** | findCharacter, findItems |
| **Batch Ops** | applyBatchOperation |
| **Events** | registerHook, fireEvent |
| **UI** | showDialog, showConfirm, showInput |
| **Logging** | log |
| **Settings** | getSetting, setSetting |
| **Permissions** | hasPermission |

---

## Character Functions

### getCharacter(name: string) → Character | nil

Gets a character by name.

**Parameters:**
- `name` (string) - Character name

**Returns:** Character object or nil if not found

**Permissions:** `read_save`

**Example:**
```lua
local terra = api.getCharacter("Terra")
if terra then
    print("Level: " .. terra.level)
end
```

**Character Object:**
```lua
{
    name = "Terra",
    level = 65,
    exp = 50000,
    hp = 450,
    mp = 120,
    vigor = 28,
    speed = 12,
    stamina = 24,
    magic = 31,
    defense = 8,
    magic_defense = 10,
    -- ... more fields
}
```

---

### setCharacter(name: string, character: Character) → error | nil

Sets a character's data.

**Parameters:**
- `name` (string) - Character name
- `character` (Character) - Updated character object

**Returns:** error message or nil on success

**Permissions:** `write_save`

**Example:**
```lua
local terra = api.getCharacter("Terra")
terra.level = 99
api.setCharacter("Terra", terra)
```

---

## Inventory Functions

### getInventory() → Inventory

Gets the current inventory.

**Returns:** Inventory object

**Permissions:** `read_save`

**Example:**
```lua
local inv = api.getInventory()
print("Items: " .. #inv.items)
```

**Inventory Object:**
```lua
{
    items = {
        {name = "Antidote", count = 5},
        {name = "Potion", count = 10},
        -- ...
    },
    equipment = {...},
    relics = {...}
}
```

---

### setInventory(inventory: Inventory) → error | nil

Sets the inventory.

**Parameters:**
- `inventory` (Inventory) - Updated inventory object

**Returns:** error message or nil on success

**Permissions:** `write_save`

**Example:**
```lua
local inv = api.getInventory()
inv.items[1].count = 99
api.setInventory(inv)
```

---

## Party Functions

### getParty() → Party

Gets the current party.

**Returns:** Party object

**Permissions:** `read_save`

**Example:**
```lua
local party = api.getParty()
print("Party size: " .. #party.members)
```

---

### setParty(party: Party) → error | nil

Sets the party.

**Parameters:**
- `party` (Party) - Updated party object

**Returns:** error message or nil on success

**Permissions:** `write_save`

---

## Equipment Functions

### getEquipment() → Equipment

Gets equipped items.

**Returns:** Equipment object

**Permissions:** `read_save`

---

### setEquipment(equipment: Equipment) → error | nil

Sets equipped items.

**Parameters:**
- `equipment` (Equipment) - Updated equipment object

**Returns:** error message or nil on success

**Permissions:** `write_save`

---

## Query Functions

### findCharacter(predicate: function) → Character | nil

Finds first character matching predicate.

**Parameters:**
- `predicate` (function) - Function(character) → boolean

**Returns:** Character or nil

**Permissions:** `read_save`

**Example:**
```lua
local ch = api.findCharacter(function(c)
    return c.name == "Terra" and c.level > 50
end)
```

---

### findItems(predicate: function) → Item[]

Finds all items matching predicate.

**Parameters:**
- `predicate` (function) - Function(item) → boolean

**Returns:** Array of items

**Permissions:** `read_save`

**Example:**
```lua
local potions = api.findItems(function(item)
    return item.name:find("Potion") ~= nil
end)
```

---

## Batch Operations

### applyBatchOperation(operation: string, params: table) → count, error

Applies a batch operation.

**Parameters:**
- `operation` (string) - Operation name
- `params` (table) - Operation parameters

**Returns:** Number of affected items, error message

**Permissions:** `write_save`

**Supported Operations:**
- `max_all_stats`
- `max_level`
- `learn_all_magic`
- `learn_all_skills`
- And more...

**Example:**
```lua
local count, err = api.applyBatchOperation("max_level", {})
if err then
    api.log("error", err)
else
    api.log("info", count .. " characters leveled")
end
```

---

## Event Functions

### registerHook(event: string, callback: function) → error | nil

Registers for an event.

**Parameters:**
- `event` (string) - Event name
- `callback` (function) - Callback function

**Returns:** error message or nil

**Permissions:** `events`

**Available Events:**
- `on_save` - Before save
- `on_load` - After load
- `on_export` - Before export
- `on_sync` - Before sync

**Example:**
```lua
api.registerHook("on_save", function(data)
    api.log("info", "Saving: " .. data.filename)
end)
```

---

### fireEvent(event: string, data: any) → error | nil

Fires an event.

**Parameters:**
- `event` (string) - Event name
- `data` (any) - Event data

**Returns:** error message or nil

**Permissions:** `events`

---

## UI Functions

### showDialog(title: string, message: string) → error | nil

Shows a dialog box.

**Parameters:**
- `title` (string) - Dialog title
- `message` (string) - Dialog message

**Returns:** error message or nil

**Permissions:** `ui_display`

**Example:**
```lua
api.showDialog("Success", "Operation completed!")
```

---

### showConfirm(title: string, message: string) → boolean

Shows a confirmation dialog.

**Parameters:**
- `title` (string) - Dialog title
- `message` (string) - Confirmation message

**Returns:** true if confirmed, false if cancelled

**Permissions:** `ui_display`

**Example:**
```lua
if api.showConfirm("Confirm", "Continue operation?") then
    -- User clicked OK
else
    -- User clicked Cancel
end
```

---

### showInput(prompt: string) → string

Shows an input dialog.

**Parameters:**
- `prompt` (string) - Prompt text

**Returns:** User input string

**Permissions:** `ui_display`

**Example:**
```lua
local name = api.showInput("Enter character name:")
if name and name ~= "" then
    api.log("info", "Name: " .. name)
end
```

---

## Logging Functions

### log(level: string, message: string) → error | nil

Logs a message.

**Parameters:**
- `level` (string) - Log level: "debug", "info", "warn", "error"
- `message` (string) - Message text

**Returns:** error message or nil

**Example:**
```lua
api.log("info", "Plugin started")
api.log("error", "Something went wrong")
```

---

## Settings Functions

### getSetting(key: string) → any

Gets a plugin setting.

**Parameters:**
- `key` (string) - Setting key

**Returns:** Setting value

**Example:**
```lua
local debug = api.getSetting("debug_mode")
```

---

### setSetting(key: string, value: any) → error | nil

Sets a plugin setting.

**Parameters:**
- `key` (string) - Setting key
- `value` (any) - Setting value

**Returns:** error message or nil

**Example:**
```lua
api.setSetting("last_run", os.time())
```

---

## Permission Functions

### hasPermission(permission: string) → boolean

Checks if plugin has permission.

**Parameters:**
- `permission` (string) - Permission name

**Returns:** true if has permission

**Permissions:**
- `read_save` - Read save data
- `write_save` - Write save data
- `ui_display` - Show UI dialogs
- `events` - Register event hooks

**Example:**
```lua
if api.hasPermission("write_save") then
    -- Modify save
else
    api.log("error", "No write permission")
end
```

---

## Standard Library

Lua standard library support is limited to safe modules:

### Table Library
- `table.insert(table, value)`
- `table.remove(table, [pos])`
- `table.concat(table, [sep])`
- `table.sort(table, [comp])`

### String Library
- `string.sub(s, i, [j])`
- `string.len(s)`
- `string.upper(s)`
- `string.lower(s)`
- `string.find(s, pattern, [init], [plain])`
- `string.gsub(s, pattern, repl, [n])`

### Math Library
- `math.abs(x)`
- `math.floor(x)`
- `math.ceil(x)`
- `math.min(x, ...)`
- `math.max(x, ...)`
- `math.random([m], [n])`

### UTF8 Library
- `utf8.len(s)`
- `utf8.codepoint(s, i, [j])`

---

## Error Handling

All API functions return error messages in their return values. Use `pcall` for safety:

```lua
local ok, result = pcall(function()
    return api.getCharacter("Terra")
end)

if not ok then
    api.log("error", "Failed: " .. tostring(result))
else
    -- Use result
end
```

---

## Best Practices

1. **Always check permissions**
```lua
if not api.hasPermission("write_save") then
    return
end
```

2. **Handle errors gracefully**
```lua
local ch = api.getCharacter("Terra")
if not ch then
    api.log("error", "Character not found")
    return
end
```

3. **Log important operations**
```lua
api.log("info", "Starting operation")
-- Do work
api.log("info", "Operation completed")
```

4. **Use pcall for robustness**
```lua
pcall(function()
    -- Safe code
end)
```

5. **Minimize API calls**
```lua
-- Cache results
local terra = api.getCharacter("Terra")
-- Use terra multiple times
```

---

## Constants

**Permission Constants:**
```lua
"read_save"    -- Read character/inventory data
"write_save"   -- Modify character/inventory data
"ui_display"   -- Show UI dialogs
"events"       -- Register event hooks
```

**Hook Constants:**
```lua
"on_save"      -- Before save operation
"on_load"      -- After load operation
"on_export"    -- Before export operation
"on_sync"      -- Before cloud sync
```

**Log Levels:**
```lua
"debug"        -- Debug messages
"info"         -- Information
"warn"         -- Warnings
"error"        -- Errors
```

---

## Limitations

- **Timeout:** 30 seconds per execution
- **Memory:** 50MB per plugin
- **Files:** No file system access
- **Network:** No network access
- **System:** No OS module access

---

## API Version

- **API Version:** 1.0
- **Last Updated:** January 16, 2026
- **Status:** Stable

For more examples, see PHASE_4B_PLUGIN_EXAMPLES.md
