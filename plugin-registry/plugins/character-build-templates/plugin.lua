--[[
    Character Build Templates
    Version: 1.0.0

    Save and share character build presets (equipment, espers, spells). Manage
    a build library with validation and quick-apply. Uses file_io for persistence.

    Permissions: read_save, write_save, ui_display, file_io
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    buildsDirectory = "builds",
    maxBuildsPerPage = 10,
    defaultCategory = "general",
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

local function serializeBuild(char)
    if not char then return nil end
    local build = {
        version = "1.0",
        charName = char.Name or char.name or "Unknown",
        charID = char.ID or char.id or 0,
        level = char.Level or char.level or 1,
        equipment = {
            weapon = char.Equipment and char.Equipment.WeaponID or 0,
            shield = char.Equipment and char.Equipment.ShieldID or 0,
            armor = char.Equipment and char.Equipment.ArmorID or 0,
            helmet = char.Equipment and char.Equipment.HelmetID or 0,
            relic1 = char.Equipment and char.Equipment.Relic1ID or 0,
            relic2 = char.Equipment and char.Equipment.Relic2ID or 0,
        },
        esper = char.EsperID or char.esperID or 0,
        spells = {},
    }
    
    if char.Spells or char.spells then
        local spells = char.Spells or char.spells
        for _, spell in ipairs(spells) do
            if spell.Value and spell.Value > 0 then
                table.insert(build.spells, {
                    id = spell.Index or spell.index or 0,
                    name = spell.Name or spell.name or "",
                    value = spell.Value or spell.value or 0,
                })
            end
        end
    end
    
    return build
end

local function deserializeBuild(data)
    if type(data) ~= "string" then return nil end
    local ok, build = pcall(function()
        -- Simple JSON-like parser (assumes well-formed data)
        return load("return " .. data)()
    end)
    if not ok or type(build) ~= "table" then return nil end
    return build
end

local function buildToString(build)
    if not build then return nil end
    -- Convert build to Lua table literal
    local lines = {"{"}
    table.insert(lines, string.format('  version = "%s",', build.version or "1.0"))
    table.insert(lines, string.format('  charName = "%s",', build.charName or "Unknown"))
    table.insert(lines, string.format('  charID = %d,', build.charID or 0))
    table.insert(lines, string.format('  level = %d,', build.level or 1))
    table.insert(lines, "  equipment = {")
    table.insert(lines, string.format('    weapon = %d,', build.equipment.weapon or 0))
    table.insert(lines, string.format('    shield = %d,', build.equipment.shield or 0))
    table.insert(lines, string.format('    armor = %d,', build.equipment.armor or 0))
    table.insert(lines, string.format('    helmet = %d,', build.equipment.helmet or 0))
    table.insert(lines, string.format('    relic1 = %d,', build.equipment.relic1 or 0))
    table.insert(lines, string.format('    relic2 = %d,', build.equipment.relic2 or 0))
    table.insert(lines, "  },")
    table.insert(lines, string.format('  esper = %d,', build.esper or 0))
    table.insert(lines, "  spells = {")
    for _, spell in ipairs(build.spells or {}) do
        table.insert(lines, string.format('    {id = %d, name = "%s", value = %d},', spell.id, spell.name, spell.value))
    end
    table.insert(lines, "  },")
    table.insert(lines, "}")
    return table.concat(lines, "\n")
end

local function saveBuild(api, build, filename)
    if not build or not filename then return false end
    local path = string.format("%s/%s.lua", config.buildsDirectory, filename)
    local content = buildToString(build)
    if not content then return false end
    local ok = safeCall(api, "WriteFile", path, content)
    return ok ~= nil
end

local function loadBuild(api, filename)
    if not filename then return nil end
    local path = string.format("%s/%s.lua", config.buildsDirectory, filename)
    local content = safeCall(api, "ReadFile", path)
    if not content then return nil end
    return deserializeBuild(content)
end

local function listBuilds(api)
    local files = safeCall(api, "ListFiles", config.buildsDirectory)
    if not files then return {} end
    local builds = {}
    for _, file in ipairs(files) do
        if file:match("%.lua$") then
            local name = file:gsub("%.lua$", "")
            table.insert(builds, name)
        end
    end
    return builds
end

local function deleteBuild(api, filename)
    if not filename then return false end
    local path = string.format("%s/%s.lua", config.buildsDirectory, filename)
    local ok = safeCall(api, "DeleteFile", path)
    return ok ~= nil
end

local function applyBuild(api, build, charSlot)
    if not build or not charSlot then return false end
    local char = safeCall(api, "GetCharacter", charSlot)
    if not char then return false end
    
    -- Apply equipment
    if build.equipment then
        if char.Equipment then
            char.Equipment.WeaponID = build.equipment.weapon or 0
            char.Equipment.ShieldID = build.equipment.shield or 0
            char.Equipment.ArmorID = build.equipment.armor or 0
            char.Equipment.HelmetID = build.equipment.helmet or 0
            char.Equipment.Relic1ID = build.equipment.relic1 or 0
            char.Equipment.Relic2ID = build.equipment.relic2 or 0
        end
    end
    
    -- Apply esper
    if build.esper then
        if char.EsperID ~= nil then
            char.EsperID = build.esper
        elseif char.esperID ~= nil then
            char.esperID = build.esper
        end
    end
    
    -- Apply spells (if write access available)
    if build.spells and (char.Spells or char.spells) then
        for _, spell in ipairs(build.spells) do
            local ok = safeCall(api, "SetSpell", charSlot, spell.id, spell.value)
            if not ok then
                -- Fallback: try direct mutation if API unavailable
                local spells = char.Spells or char.spells
                if spells and spells[spell.id] then
                    spells[spell.id].Value = spell.value
                end
            end
        end
    end
    
    local saveOk = safeCall(api, "SaveCharacter", charSlot)
    return saveOk ~= nil
end

-- ============================================================================
-- Views
-- ============================================================================

local function saveBuildView(api)
    local slotStr = api:ShowInput("Save Build", "Enter character slot (0-15):", "0")
    if not slotStr then return end
    local slot = tonumber(slotStr)
    if not slot then
        api:ShowDialog("Error", "Invalid slot number.")
        return
    end
    
    local char = safeCall(api, "GetCharacter", slot)
    if not char then
        api:ShowDialog("Error", "Unable to fetch character.")
        return
    end
    
    local build = serializeBuild(char)
    if not build then
        api:ShowDialog("Error", "Unable to serialize character build.")
        return
    end
    
    local filename = api:ShowInput("Save Build", "Enter filename (no extension):", char.Name or "build")
    if not filename or filename == "" then return end
    
    local ok = saveBuild(api, build, filename)
    if ok then
        api:ShowDialog("Success", string.format("Build '%s' saved successfully.", filename))
    else
        api:ShowDialog("Error", "Failed to save build. Check file_io permission.")
    end
end

local function loadBuildView(api)
    local builds = listBuilds(api)
    if #builds == 0 then
        api:ShowDialog("Load Build", "No builds found. Save a build first.")
        return
    end
    
    local list = {}
    for i, name in ipairs(builds) do
        table.insert(list, string.format("%d) %s", i, name))
    end
    
    local choice = api:ShowInput("Load Build", "Select build:\n" .. table.concat(list, "\n") .. "\n\nEnter number:", "1")
    if not choice then return end
    local idx = tonumber(choice)
    if not idx or idx < 1 or idx > #builds then
        api:ShowDialog("Error", "Invalid selection.")
        return
    end
    
    local build = loadBuild(api, builds[idx])
    if not build then
        api:ShowDialog("Error", "Failed to load build.")
        return
    end
    
    local slotStr = api:ShowInput("Load Build", string.format("Build: %s\nApply to character slot (0-15):", build.charName), "0")
    if not slotStr then return end
    local slot = tonumber(slotStr)
    if not slot then
        api:ShowDialog("Error", "Invalid slot number.")
        return
    end
    
    local ok = applyBuild(api, build, slot)
    if ok then
        api:ShowDialog("Success", string.format("Build applied to slot %d.", slot))
    else
        api:ShowDialog("Error", "Failed to apply build. Check write_save permission.")
    end
end

local function listBuildsView(api)
    local builds = listBuilds(api)
    if #builds == 0 then
        api:ShowDialog("Build Library", "No builds found.")
        return
    end
    
    local lines = {}
    for i, name in ipairs(builds) do
        table.insert(lines, string.format("%d) %s", i, name))
    end
    
    api:ShowDialog("Build Library", string.format("%d builds:\n%s", #builds, table.concat(lines, "\n")))
end

local function deleteBuildView(api)
    local builds = listBuilds(api)
    if #builds == 0 then
        api:ShowDialog("Delete Build", "No builds found.")
        return
    end
    
    local list = {}
    for i, name in ipairs(builds) do
        table.insert(list, string.format("%d) %s", i, name))
    end
    
    local choice = api:ShowInput("Delete Build", "Select build to delete:\n" .. table.concat(list, "\n") .. "\n\nEnter number:", "1")
    if not choice then return end
    local idx = tonumber(choice)
    if not idx or idx < 1 or idx > #builds then
        api:ShowDialog("Error", "Invalid selection.")
        return
    end
    
    local confirm = api:ShowInput("Delete Build", string.format("Delete '%s'? (y/n)", builds[idx]), "n")
    if not confirm or confirm:lower() ~= "y" then return end
    
    local ok = deleteBuild(api, builds[idx])
    if ok then
        api:ShowDialog("Success", string.format("Build '%s' deleted.", builds[idx]))
    else
        api:ShowDialog("Error", "Failed to delete build.")
    end
end

local function helpView(api)
    local lines = {
        "Character Build Templates",
        "- Save: Serialize character equipment/espers/spells to file.",
        "- Load: Apply saved build to any character slot.",
        "- Library: View all saved builds.",
        "- Delete: Remove builds from library.",
        "Requires: read_save, write_save, ui_display, file_io permissions.",
        "Builds stored in: " .. config.buildsDirectory .. "/ directory.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Character Build Templates",
            "Select an option:\n1) Save Build\n2) Load Build\n3) Build Library\n4) Delete Build\n5) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            saveBuildView(api)
        elseif sel == 2 then
            loadBuildView(api)
        elseif sel == 3 then
            listBuildsView(api)
        elseif sel == 4 then
            deleteBuildView(api)
        elseif sel == 5 then
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
        api:ShowDialog("Error", "Character Build Templates failed: " .. tostring(err))
    end
end

return {
    id = "character-build-templates",
    name = "Character Build Templates",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "write_save", "ui_display", "file_io"},
    run = safeMain
}
