--[[
    Multi-Save Manager
    Version: 1.0.0

    Manage multiple playthroughs and save slots with metadata, notes,
    statistics, and quick switching between saves.

    Permissions: read_save, ui_display, file_io
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    savesDirectory = "saves",
    metadataFile = "saves_metadata.lua",
}

-- ============================================================================
-- Helpers
-- ============================================================================

local function safeCall(api, method, ...)
    if not api or not method then return nil end
    local ok, result = pcall(function(...)
        local fn = api[method]
        if type(fn) == "function" then
            return fn(api, ...)
        end
        return nil
    end, ...)
    if not ok then return nil end
    return result
end

local function loadMetadata(api)
    local path = config.metadataFile
    local content = safeCall(api, "ReadFile", path)
    if not content then return {} end
    
    local ok, data = pcall(function()
        return load("return " .. content)()
    end)
    
    if not ok or type(data) ~= "table" then return {} end
    return data
end

local function saveMetadata(api, metadata)
    local lines = {"{"}
    for filename, meta in pairs(metadata) do
        table.insert(lines, string.format('  ["%s"] = {', filename))
        table.insert(lines, string.format('    name = "%s",', meta.name or ""))
        table.insert(lines, string.format('    notes = "%s",', meta.notes or ""))
        table.insert(lines, string.format('    playtime = %d,', meta.playtime or 0))
        table.insert(lines, string.format('    completion = %d,', meta.completion or 0))
        table.insert(lines, "  },")
    end
    table.insert(lines, "}")
    
    local content = table.concat(lines, "\n")
    return safeCall(api, "WriteFile", config.metadataFile, content) ~= nil
end

local function listSaveFiles(api)
    local files = safeCall(api, "ListFiles", config.savesDirectory)
    if not files then return {} end
    
    local saves = {}
    for _, file in ipairs(files) do
        if file:match("%.sav$") or file:match("%.srm$") then
            table.insert(saves, file)
        end
    end
    
    table.sort(saves)
    return saves
end

local function getSaveStats(api, filename)
    local path = string.format("%s/%s", config.savesDirectory, filename)
    local data = safeCall(api, "LoadSaveFile", path)
    if not data then return nil end
    
    local stats = {
        playtime = data.PlayTime or data.playtime or 0,
        gil = data.Gil or data.gil or 0,
        characters = 0,
    }
    
    if data.Characters then
        for _, c in pairs(data.Characters) do
            if c and (c.Enabled or c.enabled) then
                stats.characters = stats.characters + 1
            end
        end
    end
    
    return stats
end

-- ============================================================================
-- Views
-- ============================================================================

local function listSavesView(api)
    local saves = listSaveFiles(api)
    if #saves == 0 then
        api:ShowDialog("Save Library", "No saves found in " .. config.savesDirectory)
        return
    end
    
    local metadata = loadMetadata(api)
    local lines = {string.format("%d save slots:", #saves)}
    
    for i, file in ipairs(saves) do
        local meta = metadata[file] or {}
        local name = meta.name or file
        local notes = meta.notes and (#meta.notes > 30 and (meta.notes:sub(1, 30) .. "...") or meta.notes) or ""
        table.insert(lines, string.format("%d) %s %s", i, name, notes ~= "" and ("- " .. notes) or ""))
    end
    
    api:ShowDialog("Save Library", table.concat(lines, "\n"))
end

local function addSaveMetadataView(api)
    local saves = listSaveFiles(api)
    if #saves == 0 then
        api:ShowDialog("Add Metadata", "No saves found.")
        return
    end
    
    local list = {}
    for i, file in ipairs(saves) do
        table.insert(list, string.format("%d) %s", i, file))
    end
    
    local choice = api:ShowInput("Add Metadata", "Select save:\n" .. table.concat(list, "\n") .. "\n\nEnter number:", "1")
    if not choice then return end
    local idx = tonumber(choice)
    if not idx or idx < 1 or idx > #saves then
        api:ShowDialog("Error", "Invalid selection.")
        return
    end
    
    local filename = saves[idx]
    local metadata = loadMetadata(api)
    local meta = metadata[filename] or {}
    
    local name = api:ShowInput("Add Metadata", "Custom name:", meta.name or filename)
    if not name then return end
    meta.name = name
    
    local notes = api:ShowInput("Add Metadata", "Notes/journal:", meta.notes or "")
    if notes then meta.notes = notes end
    
    -- Try to get stats from file
    local stats = getSaveStats(api, filename)
    if stats then
        meta.playtime = stats.playtime
        meta.completion = math.floor((stats.characters / 14) * 100)
    end
    
    metadata[filename] = meta
    
    if saveMetadata(api, metadata) then
        api:ShowDialog("Success", string.format("Metadata saved for '%s'.", filename))
    else
        api:ShowDialog("Error", "Failed to save metadata.")
    end
end

local function viewSaveDetailsView(api)
    local saves = listSaveFiles(api)
    if #saves == 0 then
        api:ShowDialog("Save Details", "No saves found.")
        return
    end
    
    local list = {}
    for i, file in ipairs(saves) do
        table.insert(list, string.format("%d) %s", i, file))
    end
    
    local choice = api:ShowInput("Save Details", "Select save:\n" .. table.concat(list, "\n") .. "\n\nEnter number:", "1")
    if not choice then return end
    local idx = tonumber(choice)
    if not idx or idx < 1 or idx > #saves then
        api:ShowDialog("Error", "Invalid selection.")
        return
    end
    
    local filename = saves[idx]
    local metadata = loadMetadata(api)
    local meta = metadata[filename] or {}
    local stats = getSaveStats(api, filename)
    
    local lines = {
        string.format("File: %s", filename),
        string.format("Name: %s", meta.name or "(none)"),
        string.format("Notes: %s", meta.notes or "(none)"),
    }
    
    if stats then
        table.insert(lines, string.format("Playtime: %d", stats.playtime))
        table.insert(lines, string.format("Gil: %d", stats.gil))
        table.insert(lines, string.format("Characters: %d", stats.characters))
        table.insert(lines, string.format("Completion: ~%d%%", meta.completion or 0))
    else
        table.insert(lines, "Stats: (unable to load)")
    end
    
    api:ShowDialog("Save Details", table.concat(lines, "\n"))
end

local function helpView(api)
    local lines = {
        "Multi-Save Manager",
        "- List: View all saves with custom names and notes.",
        "- Add Metadata: Set custom name, notes for a save slot.",
        "- View Details: See playtime, gil, completion stats.",
        "Requires: read_save, file_io permissions.",
        "Saves directory: " .. config.savesDirectory,
        "Metadata file: " .. config.metadataFile,
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Multi-Save Manager",
            "Select an option:\n1) List Saves\n2) Add Metadata\n3) View Details\n4) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            listSavesView(api)
        elseif sel == 2 then
            addSaveMetadataView(api)
        elseif sel == 3 then
            viewSaveDetailsView(api)
        elseif sel == 4 then
            helpView(api)
        elseif sel == 0 then
            return
        else
            api:ShowDialog("Invalid", "Please choose a valid option.")
        end
    end
end

local function safeMain(api)
    local ok, err = pcall(function()
        mainMenu(api)
    end)
    if not ok and api and api.ShowDialog then
        api:ShowDialog("Error", "Multi-Save Manager failed: " .. tostring(err))
    end
end

return {
    id = "multi-save-manager",
    name = "Multi-Save Manager",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display", "file_io"},
    run = safeMain
}
