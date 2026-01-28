--[[
    Esper Stat Growth Optimizer
    Version: 1.0.0

    Plans esper assignments for level-ups based on target stats. Uses a small
    esper bonus table and estimates total gains from current level to a chosen
    target. Read-only with conservative fallbacks when APIs are missing.

    Permissions: read_save, ui_display
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    defaultTargetLevel = 99,
    defaultCurrentLevel = 50,
}

-- ============================================================================
-- Esper bonuses aligned with models/game/esper_growth.go (fields: Vigor, Speed,
-- Stamina, MagicPwr, Defense, MagicDef). This mirrors the shipped dataset so
-- per-level advice matches in-app growth logic.
local ESPERS = {
    {name = "Ramuh", bonuses = {MagicPwr = 2, MagicDef = 1}, notes = "Lightning, magic growth"},
    {name = "Ifrit", bonuses = {Vigor = 2, Stamina = 1, MagicPwr = 1}, notes = "Balanced fire"},
    {name = "Shiva", bonuses = {Speed = 1, MagicPwr = 2, Defense = 1}, notes = "Ice, speed + magic"},
    {name = "Golem", bonuses = {Vigor = 1, Stamina = 2, Defense = 1, MagicDef = 1}, notes = "Physical shield focus"},
    {name = "Alexander", bonuses = {Stamina = 1, MagicDef = 2}, notes = "Holy, magic defense"},
    {name = "Tritoch", bonuses = {MagicPwr = 3}, notes = "Tri-element caster"},
    {name = "Kirin", bonuses = {Speed = 2, MagicPwr = 1}, notes = "Speed-oriented"},
    {name = "Unicorn", bonuses = {MagicPwr = 1, MagicDef = 2}, notes = "Status resilience"},
    {name = "Neo Bahamut", bonuses = {Vigor = 3, Speed = 1, Stamina = 2, MagicPwr = 3, Defense = 1, MagicDef = 1}, notes = "Late-game all-rounder"},
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

local function normalizeStat(input)
    if not input then return "MagicPwr" end
    local s = string.lower(tostring(input))
    if s == "vig" or s == "vigor" or s == "str" or s == "strength" then return "Vigor" end
    if s == "mag" or s == "magic" then return "MagicPwr" end
    if s == "spd" or s == "speed" then return "Speed" end
    if s == "sta" or s == "stamina" or s == "vit" then return "Stamina" end
    if s == "def" or s == "defense" then return "Defense" end
    if s == "res" or s == "mdef" or s == "magicdef" or s == "mdefense" then return "MagicDef" end
    return "MagicPwr"
end

local function bonusValue(esper, stat)
    if not esper or not esper.bonuses then return 0 end
    return esper.bonuses[stat] or 0
end

local function statLabel(stat)
    if stat == "MagicPwr" then return "Magic" end
    if stat == "MagicDef" then return "Magic Defense" end
    return stat
end

local function bonusesToText(esper)
    local parts = {}
    for k, v in pairs(esper.bonuses or {}) do
        table.insert(parts, string.format("%s %+g", k, v))
    end
    return #parts > 0 and table.concat(parts, ", ") or "None"
end

local function loadRuntimeEspers(api)
    local inv = safeCall(api, "GetEspers") or safeCall(api, "GetEsperInventory")
    if type(inv) ~= "table" then return nil, nil, false, nil end

    local list = {}
    local ownedSet = {}
    local bonusful = 0
    local detectedFields = {statGrowth = false, bonuses = false, owned = false}
    for _, e in ipairs(inv) do
        local name = e.Name or e.name or e.ID or e.id
        if name then
            local stats = e.StatGrowth or e.statGrowth or e.Bonuses or e.bonuses
            if e.StatGrowth or e.statGrowth then detectedFields.statGrowth = true end
            if e.Bonuses or e.bonuses then detectedFields.bonuses = true end
            local function pick(keys)
                for _, k in ipairs(keys) do
                    if stats and stats[k] ~= nil then return stats[k] end
                    if stats and stats[string.lower(k)] ~= nil then return stats[string.lower(k)] end
                    if e[k] ~= nil then return e[k] end
                    if e[string.lower(k)] ~= nil then return e[string.lower(k)] end
                end
                return 0
            end

            local bonuses = {
                Vigor = pick({"Vigor", "Str", "Strength"}),
                Speed = pick({"Speed"}),
                Stamina = pick({"Stamina", "Vit"}),
                MagicPwr = pick({"MagicPwr", "Magic"}),
                Defense = pick({"Defense"}),
                MagicDef = pick({"MagicDef"}),
            }

            if (bonuses.Vigor or 0) ~= 0 or (bonuses.Speed or 0) ~= 0 or (bonuses.Stamina or 0) ~= 0 or
               (bonuses.MagicPwr or 0) ~= 0 or (bonuses.Defense or 0) ~= 0 or (bonuses.MagicDef or 0) ~= 0 then
                bonusful = bonusful + 1
            end

            table.insert(list, {name = tostring(name), bonuses = bonuses, notes = "Runtime esper data"})

            local owned = e.Owned or e.owned or e.Acquired or e.acquired
            if owned then
                ownedSet[string.lower(tostring(name))] = true
                detectedFields.owned = true
            end
        end
    end

    if #list == 0 then return nil, nil, false, nil end
    local reliable = bonusful > 0
    local meta = {
        count = #list,
        withGrowth = bonusful,
        fields = detectedFields,
        ownedCount = 0,
    }
    for _ in pairs(ownedSet) do meta.ownedCount = meta.ownedCount + 1 end
    return list, ownedSet, reliable, meta
end

local function filterEspers(pool, availableSet)
    pool = pool or ESPERS
    if not availableSet then return pool end
    local filtered = {}
    for _, e in ipairs(pool) do
        if availableSet[string.lower(e.name)] then
            table.insert(filtered, e)
        end
    end
    if #filtered == 0 then return pool end
    return filtered
end

local function pickCharacter(api)
    local idxStr = api:ShowInput("Character", "Enter character slot (0-15):", "0")
    if not idxStr then return nil, "No character selected" end
    local idx = tonumber(idxStr)
    if not idx then return nil, "Invalid slot" end
    local char = safeCall(api, "GetCharacter", idx)
    if not char then return nil, "Character API unavailable" end
    return char, nil
end

local function sortByStat(esps, stat)
    table.sort(esps, function(a, b)
        return bonusValue(a, stat) > bonusValue(b, stat)
    end)
end

local function planStat(api)
    local char, err = pickCharacter(api)
    if not char then
        api:ShowDialog("Esper Planner", err or "Unable to fetch character.")
        return
    end

    local statInput = api:ShowInput("Target Stat", "Stat (Vigor/Magic/Speed/Stamina/Defense/MagicDef):", "Magic")
    if not statInput then return end
    local stat = normalizeStat(statInput)

    local currentLevel = char.Level or char.level or config.defaultCurrentLevel
    local targetInput = api:ShowInput("Target Level", string.format("Current L%d. Target level?", currentLevel), tostring(config.defaultTargetLevel))
    if not targetInput then return end
    local targetLevel = tonumber(targetInput) or config.defaultTargetLevel
    if targetLevel < currentLevel then targetLevel = currentLevel end

    local levels = math.max(0, targetLevel - currentLevel)
    if levels == 0 then
        api:ShowDialog("Esper Planner", "No levels remaining to plan.")
        return
    end

    local pool, ownedSet, reliable, meta = loadRuntimeEspers(api)
    if not pool or not reliable then
        pool = ESPERS
        reliable = false
    end
    local poolUsed = filterEspers(pool, ownedSet)
    sortByStat(poolUsed, stat)

    local best = poolUsed[1]
    local second = poolUsed[2]
    local third = poolUsed[3]

    local bestPer = bonusValue(best, stat)
    local totalGain = bestPer * levels

    local lines = {
        string.format("Character: %s (L%d → L%d, %d levels)", char.Name or char.name or "Unknown", currentLevel, targetLevel, levels),
        string.format("Target stat: %s", statLabel(stat)),
        string.format("Recommended esper: %s (per-level %+0.2f; bonuses: %s)", best.name, bestPer, bonusesToText(best)),
        string.format("Estimated total gain: %+0.2f %s", totalGain, statLabel(stat)),
    }
    if second then
        table.insert(lines, string.format("Backup: %s (per-level %+0.2f)", second.name, bonusValue(second, stat)))
    end
    if third then
        table.insert(lines, string.format("Alternates: %s (per-level %+0.2f)", third.name, bonusValue(third, stat)))
    end
    if ownedSet then
        table.insert(lines, string.format("Pool limited to owned espers (%d owned/%d total).", meta and meta.ownedCount or 0, #pool))
    end
    if reliable and meta then
        local fields = {}
        if meta.fields.statGrowth then table.insert(fields, "StatGrowth") end
        if meta.fields.bonuses then table.insert(fields, "Bonuses") end
        if meta.fields.owned then table.insert(fields, "Owned") end
        table.insert(lines, string.format("Source: runtime API (%d espers, %d with growth; detected: %s)", meta.count, meta.withGrowth, table.concat(fields, ", ")))
    else
        table.insert(lines, "Source: bundled esper table (runtime data missing or incomplete).")
    end

    api:ShowDialog("Esper Plan", table.concat(lines, "\n"))
end

local function compareEspers(api)
    local statInput = api:ShowInput("Compare", "Stat to compare (Vigor/Magic/Speed/Stamina/Defense/MagicDef):", "Magic")
    if not statInput then return end
    local stat = normalizeStat(statInput)
    local pool, ownedSet, reliable, meta = loadRuntimeEspers(api)
    if not pool or not reliable then
        pool = ESPERS
        reliable = false
    end
    local poolUsed = filterEspers(pool, ownedSet)
    sortByStat(poolUsed, stat)

    local lines = {string.format("Stat focus: %s", statLabel(stat))}
    for i, e in ipairs(poolUsed) do
        if i > 8 then break end
        table.insert(lines, string.format("%d) %s — %+0.2f per level (%s)", i, e.name, bonusValue(e, stat), bonusesToText(e)))
    end
    if ownedSet then
        table.insert(lines, string.format("Showing owned espers (%d owned/%d total).", meta and meta.ownedCount or 0, #pool))
    end
    if reliable and meta then
        local fields = {}
        if meta.fields.statGrowth then table.insert(fields, "StatGrowth") end
        if meta.fields.bonuses then table.insert(fields, "Bonuses") end
        if meta.fields.owned then table.insert(fields, "Owned") end
        table.insert(lines, string.format("Source: runtime API (%d espers, %d with growth; detected: %s)", meta.count, meta.withGrowth, table.concat(fields, ", ")))
    else
        table.insert(lines, "Source: bundled esper table (runtime data missing or incomplete).")
    end
    api:ShowDialog("Esper Comparison", table.concat(lines, "\n"))
end

local function helpView(api)
    local lines = {
        "Esper Stat Growth Optimizer (read-only)",
        "- Plans level-up esper assignments for one stat from current level to a target level.",
        "- Pulls esper list from runtime APIs when available; otherwise falls back to bundled table aligned to models/game/esper_growth.go (Vigor, Speed, Stamina, Magic, Defense, MagicDef).",
        "- If owned esper inventory is available, rankings are limited to owned espers.",
        "Tips: Stack the highest per-level bonus for your target stat; use balanced espers if you want mixed growth.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Esper Stat Growth Optimizer",
            "Select an option:\n1) Plan stat growth\n2) Compare espers for a stat\n3) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            planStat(api)
        elseif sel == 2 then
            compareEspers(api)
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
        api:ShowDialog("Error", "Esper Stat Growth Optimizer failed: " .. tostring(err))
    end
end

return {
    id = "esper-stat-growth-optimizer",
    name = "Esper Stat Growth Optimizer",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
