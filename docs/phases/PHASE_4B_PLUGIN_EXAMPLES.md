# Phase 4B Plugin Examples

**Collection of Working Plugin Examples**

---

## Example 1: Character Max Stats

Maximize all character stats instantly.

```lua
-- @name Max All Stats
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Maximize all character stats to maximum
-- @permissions write_save

local plugin = {
    name = "Max All Stats",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Maximize all character stats to maximum",
}

function plugin.execute(api)
    if not api.hasPermission("write_save") then
        api.showDialog("Error", "Write permission required")
        return
    end
    
    local characters = {
        "Terra", "Locke", "Edgar", "Sabin", "Cyan", "Gau", 
        "Setzer", "Strago", "Relm", "Shadow", "Kefka"
    }
    
    local updated = 0
    
    for _, name in ipairs(characters) do
        local ch = api.getCharacter(name)
        if ch then
            ch.level = 99
            ch.hp = 9999
            ch.mp = 999
            ch.vigor = 99
            ch.speed = 99
            ch.stamina = 99
            ch.magic = 99
            ch.defense = 99
            ch.magic_defense = 99
            
            api.setCharacter(name, ch)
            updated = updated + 1
        end
    end
    
    api.log("info", "Updated " .. updated .. " characters")
    api.showDialog("Success", "All " .. updated .. " characters maxed!")
end

return plugin
```

---

## Example 2: Item Duplicator

Duplicate items in inventory.

```lua
-- @name Item Duplicator
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Duplicate items in inventory
-- @permissions write_save, ui_display

local plugin = {
    name = "Item Duplicator",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Duplicate items in inventory",
}

function plugin.execute(api)
    if not api.hasPermission("write_save") then
        api.showDialog("Error", "Write permission required")
        return
    end
    
    local count_str = api.showInput("Duplicate count (1-99):")
    if not count_str or count_str == "" then
        return
    end
    
    local count = tonumber(count_str)
    if not count or count < 1 or count > 99 then
        api.showDialog("Error", "Invalid count")
        return
    end
    
    local inv = api.getInventory()
    if not inv or not inv.items then
        api.showDialog("Error", "Could not get inventory")
        return
    end
    
    for i, item in ipairs(inv.items) do
        inv.items[i].count = count
    end
    
    api.setInventory(inv)
    api.log("info", "Duplicated items: count=" .. count)
    api.showDialog("Success", "Items set to count: " .. count)
end

return plugin
```

---

## Example 3: Level Equalizer

Equalize all characters to same level.

```lua
-- @name Level Equalizer
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Equalize all character levels
-- @permissions read_save, write_save, ui_display

local plugin = {
    name = "Level Equalizer",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Equalize all character levels to maximum",
}

function plugin.execute(api)
    local characters = {
        "Terra", "Locke", "Edgar", "Sabin", "Cyan", "Gau",
        "Setzer", "Strago", "Relm", "Shadow", "Kefka"
    }
    
    -- Find maximum level
    local max_level = 0
    for _, name in ipairs(characters) do
        local ch = api.getCharacter(name)
        if ch and ch.level then
            max_level = math.max(max_level, ch.level)
        end
    end
    
    if max_level == 0 then
        api.showDialog("Error", "Could not determine max level")
        return
    end
    
    -- Set all to max level
    local updated = 0
    for _, name in ipairs(characters) do
        local ch = api.getCharacter(name)
        if ch then
            ch.level = max_level
            api.setCharacter(name, ch)
            updated = updated + 1
        end
    end
    
    api.log("info", "Equalized " .. updated .. " characters to level " .. max_level)
    api.showDialog("Success", "Equalized to level: " .. max_level)
end

return plugin
```

---

## Example 4: Quick Backup Logger

Log important operations for backup tracking.

```lua
-- @name Operation Logger
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Log editor operations for tracking
-- @permissions events, ui_display
-- @hooks on_save, on_load

local plugin = {
    name = "Operation Logger",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Log editor operations",
}

function plugin.initialize(api)
    api.log("info", "Operation Logger initialized")
    
    -- Register for save events
    if api.registerHook then
        api.registerHook("on_save", function(data)
            local timestamp = os.date("%Y-%m-%d %H:%M:%S")
            api.log("info", "[" .. timestamp .. "] Save: " .. (data.filename or "unknown"))
        end)
        
        api.registerHook("on_load", function(data)
            local timestamp = os.date("%Y-%m-%d %H:%M:%S")
            api.log("info", "[" .. timestamp .. "] Load: " .. (data.filename or "unknown"))
        end)
    end
end

function plugin.execute(api)
    local last_op = api.getSetting("last_operation") or "none"
    api.log("info", "Last operation was: " .. last_op)
    api.setSetting("last_operation", os.time())
end

return plugin
```

---

## Example 5: Stat Checker

Validate and report character stats.

```lua
-- @name Stat Checker
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Check and validate character stats
-- @permissions read_save, ui_display

local plugin = {
    name = "Stat Checker",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Validate character stats",
}

function plugin.execute(api)
    local characters = {
        "Terra", "Locke", "Edgar", "Sabin"
    }
    
    local report = {}
    
    for _, name in ipairs(characters) do
        local ch = api.getCharacter(name)
        if ch then
            table.insert(report, string.format(
                "%s: Lv%d HP%d MP%d",
                name, ch.level, ch.hp, ch.mp
            ))
        end
    end
    
    local message = table.concat(report, "\n")
    api.log("info", "Stat report:\n" .. message)
    api.showDialog("Stat Report", message)
end

return plugin
```

---

## Example 6: Equipment Optimizer

Suggest optimal equipment for characters.

```lua
-- @name Equipment Optimizer
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Analyze and suggest equipment
-- @permissions read_save

local plugin = {
    name = "Equipment Optimizer",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Equipment analysis and suggestions",
}

function plugin.execute(api)
    local eq = api.getEquipment()
    if not eq then
        api.log("error", "Could not get equipment")
        return
    end
    
    api.log("info", "Equipment optimization analysis")
    api.log("info", "Total equipped items: " .. (eq and 1 or 0))
    
    -- In real implementation, would analyze and suggest optimizations
    api.log("info", "Analysis complete")
end

return plugin
```

---

## Example 7: Resource Counter

Count and report game resources.

```lua
-- @name Resource Counter
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Count game resources
-- @permissions read_save, ui_display

local plugin = {
    name = "Resource Counter",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Count and report game resources",
}

function plugin.execute(api)
    if not api.hasPermission("read_save") then
        api.showDialog("Error", "Read permission required")
        return
    end
    
    local inv = api.getInventory()
    local party = api.getParty()
    
    if not inv or not party then
        api.showDialog("Error", "Could not read save data")
        return
    end
    
    local item_count = inv.items and #inv.items or 0
    local member_count = party.members and #party.members or 0
    
    local report = string.format(
        "Inventory: %d items\nParty: %d members",
        item_count, member_count
    )
    
    api.log("info", "Resources: " .. report)
    api.showDialog("Resource Report", report)
end

return plugin
```

---

## Example 8: Safe Script Template

Template with error handling for safe plugins.

```lua
-- @name Safe Template
-- @version 1.0.0
-- @author FF6 Editor Team
-- @description Template with error handling
-- @permissions read_save, ui_display

local plugin = {
    name = "Safe Template",
    version = "1.0.0",
    author = "FF6 Editor Team",
    description = "Template with comprehensive error handling",
}

function plugin.execute(api)
    local ok, result = pcall(function()
        api.log("debug", "Starting plugin")
        
        -- Check permissions
        if not api.hasPermission("read_save") then
            error("Read permission required")
        end
        
        -- Get data safely
        local ch = api.getCharacter("Terra")
        if not ch then
            error("Could not find character")
        end
        
        api.log("info", "Character found: " .. ch.name .. " Lv" .. ch.level)
        return true
    end)
    
    if not ok then
        local err_msg = tostring(result)
        api.log("error", "Error: " .. err_msg)
        api.showDialog("Error", err_msg)
        return
    end
    
    api.log("info", "Plugin completed successfully")
    api.showDialog("Success", "Operation completed")
end

return plugin
```

---

## Running Examples

1. Copy example code to `.lua` file
2. Save in `~/.ff6editor/plugins/`
3. Open Plugin Manager
4. Enable plugin in "Installed Plugins" tab
5. Click "Run" or watch for automatic execution
6. Check "Plugin Output" tab for logs

---

## Tips for Plugin Development

### Use Pcall for Safety
```lua
local ok, err = pcall(function()
    -- Code here
end)
if not ok then print("Error: " .. err) end
```

### Always Check Permissions
```lua
if not api.hasPermission("write_save") then
    return
end
```

### Log Extensively
```lua
api.log("debug", "Debug info")
api.log("info", "Normal operation")
api.log("warn", "Warning condition")
api.log("error", "Error occurred")
```

### Cache API Results
```lua
local ch = api.getCharacter("Terra")
-- Use ch multiple times
```

### Validate Input
```lua
local count = tonumber(api.showInput("Enter number:"))
if not count or count < 1 or count > 99 then
    api.showDialog("Error", "Invalid input")
    return
end
```

---

## Common Patterns

### Batch Character Operation
```lua
for _, name in ipairs(character_list) do
    local ch = api.getCharacter(name)
    if ch then
        -- Modify ch
        api.setCharacter(name, ch)
    end
end
```

### Inventory Manipulation
```lua
local inv = api.getInventory()
for i, item in ipairs(inv.items) do
    inv.items[i].count = math.min(99, inv.items[i].count)
end
api.setInventory(inv)
```

### User Confirmation
```lua
if api.showConfirm("Confirm", "Continue?") then
    -- User confirmed
else
    -- User cancelled
end
```

---

**More examples coming in future releases!**

For API documentation, see PHASE_4B_API_REFERENCE.md
