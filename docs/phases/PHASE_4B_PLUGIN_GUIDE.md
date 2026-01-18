# Phase 4B Plugin System - User Guide

**Version:** 1.0  
**Last Updated:** January 16, 2026  
**Status:** Complete

---

## Introduction

The FF6 Save Editor Plugin System enables users to extend the editor with custom functionality through Lua scripting. Plugins can manipulate save data, add UI elements, and automate complex tasks.

---

## Getting Started

### Installation

1. Create a `plugins` directory in your FF6 Editor home directory
   ```
   ~/.ff6editor/plugins/
   ```

2. Place your plugin files there (`.lua` files)

3. Restart the editor or use Plugin Manager to reload

### Your First Plugin

Create `HelloWorld.lua`:

```lua
-- @name Hello World
-- @version 1.0.0
-- @author You
-- @description A simple hello world plugin

local plugin = {
    name = "Hello World",
    version = "1.0.0",
    author = "You",
    description = "A simple hello world plugin"
}

function plugin.execute(api)
    api.showDialog("Hello", "Welcome to FF6 Plugins!")
    api.log("info", "Hello World plugin executed")
end

return plugin
```

---

## Plugin Structure

Every plugin must export a table with required fields:

```lua
local plugin = {
    -- Required metadata
    name = "Plugin Name",
    version = "1.0.0",
    author = "Your Name",
    description = "What this plugin does",
    
    -- Optional settings
    permissions = {"read_save", "write_save"},
    hooks = {"on_save", "on_load"},
    settings = {
        debug_mode = false,
        max_retries = 3
    }
}

-- Required function
function plugin.execute(api)
    -- Plugin code here
end

-- Optional lifecycle functions
function plugin.initialize()
    -- Called on plugin load
end

function plugin.shutdown()
    -- Called on plugin unload
end

return plugin
```

---

## Plugin API

### Character Functions

```lua
-- Get character by name
local terra = api.getCharacter("Terra")
if terra then
    api.log("info", "Terra Level: " .. terra.level)
end

-- Set character
terra.level = 99
api.setCharacter("Terra", terra)
```

### Inventory Functions

```lua
-- Get inventory
local inv = api.getInventory()
api.log("info", "Item count: " .. #inv.items)

-- Set inventory
api.setInventory(inv)
```

### UI Functions

```lua
-- Show dialog
api.showDialog("Title", "Message text")

-- Show confirmation
if api.showConfirm("Confirm", "Continue?") then
    api.log("info", "User confirmed")
end

-- Get user input
local name = api.showInput("Enter name:")
api.log("info", "Name entered: " .. name)
```

### Logging

```lua
api.log("debug", "Debug message")
api.log("info", "Information")
api.log("warn", "Warning")
api.log("error", "Error occurred")
```

### Permissions Check

```lua
if api.hasPermission("write_save") then
    -- Modify save
else
    api.log("error", "No write permission")
end
```

### Settings

```lua
-- Get setting
local value = api.getSetting("plugin_setting")

-- Set setting  
api.setSetting("plugin_setting", 42)
```

---

## Plugin Manager

Access via **Tools → Plugin Manager**

### Tabs

**Installed Plugins**
- List of loaded plugins
- Enable/disable toggles
- Settings button
- Remove button

**Available Plugins**
- Browse plugin registry
- Search and filter
- Install button
- Rating and downloads

**Plugin Settings**
- Auto-load plugins on startup
- Sandbox security level
- Max concurrent plugins
- Execution timeout

**Plugin Output**
- Real-time execution log
- Filter by level
- Search capability
- Clear log

---

## Example Plugins

### Character Max Stats

```lua
local plugin = {
    name = "Max Stats",
    version = "1.0.0",
    author = "Examples",
    description = "Maximize character stats"
}

function plugin.execute(api)
    local characters = {"Terra", "Locke", "Edgar", "Sabin", "Cyan", "Gau", "Setzer", "Strago", "Relm", "Shadow", "Kefka"}
    
    for _, name in ipairs(characters) do
        local char = api.getCharacter(name)
        if char then
            char.level = 99
            char.hp = 9999
            char.mp = 999
            api.setCharacter(name, char)
        end
    end
    
    api.showDialog("Success", "All characters maxed!")
end

return plugin
```

### Item Duplicator

```lua
local plugin = {
    name = "Item Duplicator",
    version = "1.0.0",
    author = "Examples",
    description = "Duplicate items"
}

function plugin.execute(api)
    if not api.hasPermission("write_save") then
        api.showDialog("Error", "Write permission required")
        return
    end
    
    local inv = api.getInventory()
    local count = api.showInput("Duplicate count (1-99):")
    count = tonumber(count) or 1
    
    for i, item in ipairs(inv.items) do
        inv.items[i].count = math.min(count, 99)
    end
    
    api.setInventory(inv)
    api.showDialog("Done", "Items duplicated")
end

return plugin
```

### Save Backup

```lua
local plugin = {
    name = "Auto Backup",
    version = "1.0.0",
    author = "Examples",
    description = "Auto-backup on save",
    hooks = {"on_save"}
}

function plugin.execute(api)
    api.log("info", "Creating backup...")
    api.setSetting("last_backup", os.time())
    api.log("info", "Backup completed")
end

return plugin
```

---

## Best Practices

### Performance
- Keep plugins small and fast
- Avoid heavy loops
- Cache API results
- Use efficient data structures

### Error Handling
```lua
function plugin.execute(api)
    local success, err = pcall(function()
        -- Your code here
    end)
    
    if not success then
        api.log("error", "Plugin error: " .. tostring(err))
        api.showDialog("Error", err)
    end
end
```

### Permissions
- Always check permissions
- Request minimal permissions needed
- Handle permission denials gracefully

### Logging
- Log important operations
- Include debug info
- Use appropriate levels
- Don't spam logs

---

## Troubleshooting

### Plugin Not Loading
- Check file is in `~/.ff6editor/plugins/`
- Verify `.lua` extension
- Check for Lua syntax errors
- Look in Plugin Output log

### Plugin Disabled
- Check Plugin Manager settings
- Verify permissions are granted
- Review execution log for errors

### Slow Plugin
- Check for infinite loops
- Minimize API calls
- Profile with debug logging
- Consider caching results

### Permission Denied
- Check plugin requests required permissions
- Verify user has granted permissions
- Review Plugin Manager settings

---

## Advanced Topics

### Custom Settings

```lua
local plugin = {
    name = "Advanced",
    settings = {
        debug = false,
        max_items = 99,
        features = {"feature1", "feature2"}
    }
}

function plugin.execute(api)
    local debug = api.getSetting("debug")
    if debug then api.log("debug", "Debug mode enabled") end
end
```

### Event Hooks

```lua
-- Plugins can register for events
-- Available hooks: on_save, on_load, on_export, on_sync

function plugin.initialize()
    -- Register for save event
    api.registerHook("on_save", function(data)
        api.log("info", "File saved: " .. data.filename)
    end)
end
```

### Error Recovery

```lua
function plugin.execute(api)
    -- Wrap in pcall for safety
    local ok, result = pcall(function()
        return api.getCharacter("Terra")
    end)
    
    if not ok then
        api.log("error", "Failed to get character")
        return
    end
    
    if result then
        -- Process result
    end
end
```

---

## Common Patterns

### Batch Character Updates
```lua
local names = {"Terra", "Locke", "Edgar"}
for _, name in ipairs(names) do
    local ch = api.getCharacter(name)
    if ch then
        ch.level = 99
        api.setCharacter(name, ch)
    end
end
```

### Conditional Logic
```lua
local inv = api.getInventory()
if #inv.items > 50 then
    api.showDialog("Info", "Inventory full!")
end
```

### User Prompts
```lua
local response = api.showInput("Enter value:")
if response and response ~= "" then
    api.log("info", "User entered: " .. response)
end
```

---

## Security Notes

- Plugins run in a sandboxed Lua environment
- File system access is blocked
- Network access is restricted
- Memory and timeout limits enforced
- Only safe Lua libraries available

Permissions required for different operations:
- `read_save` - Read character/inventory data
- `write_save` - Modify character/inventory data  
- `ui_display` - Show dialogs and UI elements
- `events` - Register event hooks

---

## Support & Resources

- **API Reference:** See PHASE_4B_API_REFERENCE.md
- **Examples:** See PHASE_4B_PLUGIN_EXAMPLES.md
- **Plugin List:** Available in Plugin Manager
- **Documentation:** Check in-game help

---

## Creating Plugins

### Step 1: Planning
- Define what your plugin does
- Identify required API functions
- List needed permissions

### Step 2: Development
- Create `.lua` file in plugins directory
- Write plugin code with metadata
- Test with Plugin Manager
- Review execution log

### Step 3: Refinement
- Handle edge cases
- Add error handling
- Optimize performance
- Document for users

### Step 4: Distribution
- Share plugin file with community
- Document requirements
- Provide examples
- Get user feedback

---

**Happy Coding!**

For more information, see the API reference and examples documentation.
