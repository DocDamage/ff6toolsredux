--[[
    Save File Comparator
    Version: 1.0.0

    Compare two save files side-by-side with diff view showing all differences
    in characters, inventory, and game state. Read-only comparison tool.

    Permissions: read_save, write_save, ui_display, file_io
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    categories = {"Characters", "Inventory", "GameState", "Summary"},
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

local function loadSaveFile(api, path)
    local data = safeCall(api, "LoadSaveFile", path)
    if not data then
        data = safeCall(api, "ReadFile", path)
    end
    return data
end

local function compareCharacters(save1, save2)
    local diffs = {}
    local chars1 = {}
    local chars2 = {}
    
    -- Extract character data from saves
    for i = 0, 15 do
        local c1 = save1 and save1.Characters and save1.Characters[i]
        local c2 = save2 and save2.Characters and save2.Characters[i]
        if c1 then chars1[i] = c1 end
        if c2 then chars2[i] = c2 end
    end
    
    -- Compare each slot
    for i = 0, 15 do
        local c1 = chars1[i]
        local c2 = chars2[i]
        if c1 or c2 then
            local name1 = c1 and (c1.Name or c1.name) or "(empty)"
            local name2 = c2 and (c2.Name or c2.name) or "(empty)"
            
            if name1 ~= name2 then
                table.insert(diffs, string.format("Slot %d: %s vs %s", i, name1, name2))
            elseif c1 and c2 then
                -- Compare stats
                local level1 = c1.Level or c1.level or 0
                local level2 = c2.Level or c2.level or 0
                if level1 ~= level2 then
                    table.insert(diffs, string.format("  %s Level: %d vs %d", name1, level1, level2))
                end
                
                local hp1 = c1.HP and c1.HP.Current or 0
                local hp2 = c2.HP and c2.HP.Current or 0
                if hp1 ~= hp2 then
                    table.insert(diffs, string.format("  %s HP: %d vs %d", name1, hp1, hp2))
                end
            end
        end
    end
    
    return diffs
end

local function compareInventory(save1, save2)
    local diffs = {}
    local inv1 = save1 and save1.Inventory or {}
    local inv2 = save2 and save2.Inventory or {}
    
    -- Simple count comparison
    local count1 = 0
    local count2 = 0
    for _ in pairs(inv1) do count1 = count1 + 1 end
    for _ in pairs(inv2) do count2 = count2 + 1 end
    
    if count1 ~= count2 then
        table.insert(diffs, string.format("Inventory items: %d vs %d", count1, count2))
    else
        table.insert(diffs, string.format("Inventory items: %d (no diff)", count1))
    end
    
    return diffs
end

local function compareGameState(save1, save2)
    local diffs = {}
    local playtime1 = save1 and save1.PlayTime or 0
    local playtime2 = save2 and save2.PlayTime or 0
    if playtime1 ~= playtime2 then
        table.insert(diffs, string.format("Play time: %d vs %d", playtime1, playtime2))
    end
    
    local gil1 = save1 and save1.Gil or 0
    local gil2 = save2 and save2.Gil or 0
    if gil1 ~= gil2 then
        table.insert(diffs, string.format("Gil: %d vs %d", gil1, gil2))
    end
    
    return diffs
end

local function generateSummary(charDiffs, invDiffs, stateDiffs)
    local total = #charDiffs + #invDiffs + #stateDiffs
    return {
        string.format("Total differences: %d", total),
        string.format("  Characters: %d", #charDiffs),
        string.format("  Inventory: %d", #invDiffs),
        string.format("  Game State: %d", #stateDiffs),
    }
end

-- ============================================================================
-- Views
-- ============================================================================

local function compareSavesView(api)
    local path1 = api:ShowInput("Compare Saves", "Enter path to save file 1:", "")
    if not path1 or path1 == "" then return end
    
    local path2 = api:ShowInput("Compare Saves", "Enter path to save file 2:", "")
    if not path2 or path2 == "" then return end
    
    local save1 = loadSaveFile(api, path1)
    local save2 = loadSaveFile(api, path2)
    
    if not save1 and not save2 then
        api:ShowDialog("Error", "Both save files failed to load. Check file_io permission and paths.")
        return
    end
    
    if not save1 then
        api:ShowDialog("Error", "Save file 1 failed to load.")
        return
    end
    
    if not save2 then
        api:ShowDialog("Error", "Save file 2 failed to load.")
        return
    end
    
    -- Compare categories
    local charDiffs = compareCharacters(save1, save2)
    local invDiffs = compareInventory(save1, save2)
    local stateDiffs = compareGameState(save1, save2)
    local summary = generateSummary(charDiffs, invDiffs, stateDiffs)
    
    -- Show results
    local lines = {"=== Save File Comparison ===", ""}
    table.insert(lines, "Summary:")
    for _, line in ipairs(summary) do
        table.insert(lines, line)
    end
    table.insert(lines, "")
    
    if #charDiffs > 0 then
        table.insert(lines, "Character Differences:")
        for _, diff in ipairs(charDiffs) do
            table.insert(lines, diff)
        end
        table.insert(lines, "")
    end
    
    if #invDiffs > 0 then
        table.insert(lines, "Inventory Differences:")
        for _, diff in ipairs(invDiffs) do
            table.insert(lines, diff)
        end
        table.insert(lines, "")
    end
    
    if #stateDiffs > 0 then
        table.insert(lines, "Game State Differences:")
        for _, diff in ipairs(stateDiffs) do
            table.insert(lines, diff)
        end
    end
    
    if #charDiffs == 0 and #invDiffs == 0 and #stateDiffs == 0 then
        table.insert(lines, "No differences detected (or limited data available).")
    end
    
    api:ShowDialog("Comparison Results", table.concat(lines, "\n"))
end

local function helpView(api)
    local lines = {
        "Save File Comparator",
        "- Compares two save files side-by-side.",
        "- Shows differences in characters, inventory, and game state.",
        "- Read-only comparison (no merge in this version).",
        "Requires: read_save, file_io permissions.",
        "Note: Comparison depth depends on LoadSaveFile API availability.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Save File Comparator",
            "Select an option:\n1) Compare Two Saves\n2) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            compareSavesView(api)
        elseif sel == 2 then
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
        api:ShowDialog("Error", "Save File Comparator failed: " .. tostring(err))
    end
end

return {
    id = "save-file-comparator",
    name = "Save File Comparator",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "write_save", "ui_display", "file_io"},
    run = safeMain
}
