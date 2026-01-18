--[[
    Slot Machine Reference
    Version: 1.0.0

    Reference guide for Setzer's Slots: outcomes, probabilities (approx), and
    strategy notes. Includes simple damage estimates using Setzer's level/magic
    when available; otherwise uses defaults.

    Permissions: read_save, ui_display
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    defaultLevel = 40,
    defaultMagic = 35,
    defaultWeaponPower = 120,
    pageSize = 12,
}

-- ============================================================================
-- Slots Data (representative core)
-- ============================================================================
-- Fields: name, reels, effect, type, basePower, notes, approxChance
local SLOTS = {
    {name = "7-Flush", reels = "7 / 7 / 7", effect = "Fire-elemental nuke", type = "Magic", basePower = 110, approxChance = "Low", notes = "Hits all; boosted by Magic."},
    {name = "Prismatic Flash", reels = "BAR / BAR / BAR", effect = "Non-elemental damage", type = "Magic", basePower = 100, approxChance = "Low", notes = "All-target."},
    {name = "Lagomorph", reels = "Chocobo / Chocobo / Chocobo", effect = "Small heal + removes Sap", type = "Utility", basePower = 0, approxChance = "Medium", notes = "Safe miss."},
    {name = "H-Bomb", reels = "Airship / Airship / Airship", effect = "Fire-elemental bomb", type = "Magic", basePower = 90, approxChance = "Low", notes = "All-target; strong early."},
    {name = "Mega Flare", reels = "Dragon / Dragon / Dragon", effect = "Non-elemental heavy damage", type = "Magic", basePower = 140, approxChance = "Very Low", notes = "Late-game nuke."},
    {name = "Muddle", reels = "Imp / Imp / Imp", effect = "Confuse on all", type = "Status", basePower = 0, approxChance = "Medium", notes = "Crowd control."},
    {name = "Golem", reels = "Golem / Golem / Golem", effect = "Party physical shield", type = "Defense", basePower = 0, approxChance = "Very Low", notes = "Excellent survival."},
    {name = "Setzer Attack", reels = "ANY mix", effect = "Random physical hit", type = "Physical", basePower = 75, approxChance = "High", notes = "Fallback."},
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

local function getSetzer(api)
    -- Try to find Setzer by name
    for i = 0, 15 do
        local c = safeCall(api, "GetCharacter", i)
        if c and (c.Name == "Setzer" or c.name == "Setzer") then
            return c
        end
    end
    return nil
end

local function calcMagicDamage(level, magic, basePower)
    level = level or config.defaultLevel
    magic = magic or config.defaultMagic
    basePower = basePower or 80
    local dmg = ((level + basePower) ^ 2) / 32 + magic * 4
    return math.max(0, math.floor(dmg))
end

local function calcPhysicalDamage(level, vigor, weaponPower)
    level = level or config.defaultLevel
    vigor = vigor or 35
    weaponPower = weaponPower or config.defaultWeaponPower
    local attack = math.floor(((level * level * vigor) / 256) + weaponPower)
    local dmg = attack - 60 -- nominal defense
    return math.max(0, math.floor(dmg))
end

local function estimateOutcome(slot, setzer)
    if slot.type == "Magic" then
        return calcMagicDamage(setzer and setzer.Level, setzer and setzer.Magic, slot.basePower)
    elseif slot.type == "Physical" then
        return calcPhysicalDamage(setzer and setzer.Level, setzer and setzer.Vigor, setzer and setzer.WeaponPower)
    end
    return 0
end

local function formatSlot(slot, setzer)
    local dmg = estimateOutcome(slot, setzer)
    local lines = {
        string.format("%s (%s)", slot.name, slot.reels),
        string.format("Effect: %s", slot.effect),
        string.format("Type: %s", slot.type),
        string.format("Approx Chance: %s", slot.approxChance or "Unknown"),
        string.format("Est. Damage: %s", slot.type == "Utility" or slot.type == "Status" or slot.type == "Defense" and "-" or tostring(dmg)),
    }
    if slot.notes then table.insert(lines, "Notes: " .. slot.notes) end
    return table.concat(lines, "\n")
end

-- ============================================================================
-- Views
-- ============================================================================

local function listAll(api)
    local setzer = getSetzer(api)
    local lines = {}
    for _, slot in ipairs(SLOTS) do
        table.insert(lines, formatSlot(slot, setzer))
        table.insert(lines, "")
    end
    api:ShowDialog("Slot Outcomes", table.concat(lines, "\n"))
end

local function searchView(api)
    local term = api:ShowInput("Search", "Keyword (name/effect/type):", "")
    if not term or term == "" then return end
    term = term:lower()
    local setzer = getSetzer(api)
    local lines = {}
    for _, slot in ipairs(SLOTS) do
        local hay = table.concat({slot.name, slot.effect, slot.type, slot.notes or ""}, " "):lower()
        if hay:find(term, 1, true) then
            table.insert(lines, formatSlot(slot, setzer))
            table.insert(lines, "")
        end
    end
    if #lines == 0 then
        api:ShowDialog("Search", "No matching outcomes.")
        return
    end
    api:ShowDialog("Search Results", table.concat(lines, "\n"))
end

local function helpView(api)
    local lines = {
        "Slot Machine Reference (read-only)",
        "- Lists common outcomes with effects and rough probabilities.",
        "- Estimates damage using Setzer's level/magic if available, otherwise defaults.",
        "- Search by name/effect/type/notes.",
        "Notes: Probabilities are approximate; timing skills and rigging tricks not modeled.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Slot Machine Reference",
            "Select an option:\n1) All Outcomes\n2) Search\n3) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            listAll(api)
        elseif sel == 2 then
            searchView(api)
        elseif sel == 3 then
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
        api:ShowDialog("Error", "Slot Machine Reference failed: " .. tostring(err))
    end
end

return {
    id = "slot-machine-reference",
    name = "Slot Machine Reference",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
