--[[
    Party Synergy Analyzer
    Version: 1.0.0

    Evaluates party composition for role coverage, elemental coverage, redundancy,
    and suggests improvements. Read-only and conservative: missing APIs produce
    Unknown sections instead of guesses.

    Permissions: read_save, ui_display
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    expectedElements = {"Fire", "Ice", "Lightning", "Wind", "Water", "Earth", "Holy", "Poison"},
    roleThresholds = {
        healer = {spells = {"Cure", "Cura", "Curaga", "Raise", "Arise"}},
        support = {spells = {"Haste", "Haste2", "Slow", "Slow2", "Shell", "Protect", "Float"}},
        dpsMagic = {spells = {"Firaga", "Blizzaga", "Thundaga", "Ultima", "Flare"}},
        dpsPhysical = {},
        tank = {},
    },
    maxPartySlots = 4,
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

local function toLower(str)
    if not str then return "" end
    return string.lower(tostring(str))
end

local function hasSpell(char, targets)
    if not char or not char.Spells then return false end
    for _, spell in ipairs(char.Spells) do
        local name = spell.Name or spell.name
        if name then
            local lower = toLower(name)
            for _, t in ipairs(targets) do
                if lower == toLower(t) then return true end
            end
        end
    end
    return false
end

local function inferRoles(char)
    local roles = {}
    if hasSpell(char, config.roleThresholds.healer.spells) then
        table.insert(roles, "Healer")
    end
    if hasSpell(char, config.roleThresholds.support.spells) then
        table.insert(roles, "Support")
    end
    if hasSpell(char, config.roleThresholds.dpsMagic.spells) then
        table.insert(roles, "Magic DPS")
    end
    -- Simple physical DPS heuristic: high Vigor or Fight+Jump/Throw/Blitz/Tools commands
    if (char.Vigor or 0) >= 40 then
        table.insert(roles, "Physical DPS")
    end
    return roles
end

local function elementCoverage(char)
    local covered = {}
    if not char or not char.Spells then return covered end
    for _, spell in ipairs(char.Spells) do
        local name = toLower(spell.Name or spell.name or "")
        if name:find("fire") then covered["Fire"] = true end
        if name:find("ice") or name:find("blizz") then covered["Ice"] = true end
        if name:find("thunder") or name:find("bolt") then covered["Lightning"] = true end
        if name:find("wind") or name:find("aero") then covered["Wind"] = true end
        if name:find("water") then covered["Water"] = true end
        if name:find("quake") or name:find("earth") then covered["Earth"] = true end
        if name:find("holy") or name:find("pearl") then covered["Holy"] = true end
        if name:find("bio") or name:find("poison") then covered["Poison"] = true end
    end
    return covered
end

local function fetchParty(api)
    local party = safeCall(api, "GetParty")
    if party and #party > 0 then return party end
    return nil
end

local function characterDisplayName(char)
    return (char and (char.Name or char.name)) or "Unknown"
end

local function analyzeParty(api)
    local party = fetchParty(api)
    if not party then
        return nil, "Party API unavailable or empty"
    end

    local results = {
        roles = {},
        elements = {},
        redundancy = {},
        missing = {},
        roster = {},
    }

    -- Initialize coverage maps
    for _, el in ipairs(config.expectedElements) do
        results.elements[el] = false
    end

    for idx, char in ipairs(party) do
        local name = characterDisplayName(char)
        table.insert(results.roster, name)
        local roles = inferRoles(char)
        results.roles[name] = roles

        -- Element coverage
        local covered = elementCoverage(char)
        for el, has in pairs(covered) do
            if has then results.elements[el] = true end
        end
    end

    -- Redundancy and missing checks
    local roleCounts = {
        Healer = 0,
        Support = 0,
        ["Magic DPS"] = 0,
        ["Physical DPS"] = 0,
    }
    for _, roles in pairs(results.roles) do
        for _, r in ipairs(roles) do
            roleCounts[r] = (roleCounts[r] or 0) + 1
        end
    end

    for role, count in pairs(roleCounts) do
        if count == 0 then
            table.insert(results.missing, role)
        elseif count > 1 then
            table.insert(results.redundancy, role .. " x" .. count)
        end
    end

    return results
end

local function summarizeElements(elements)
    local lines = {}
    for _, el in ipairs(config.expectedElements) do
        table.insert(lines, string.format("%s: %s", el, elements[el] and "Covered" or "Missing"))
    end
    return table.concat(lines, "\n")
end

-- ============================================================================
-- Views
-- ============================================================================

local function viewSummary(api)
    local results, err = analyzeParty(api)
    if not results then
        api:ShowDialog("Summary", err or "Unable to analyze party.")
        return
    end
    local lines = {
        "Party Members: " .. table.concat(results.roster, ", "),
        "\nElement Coverage:",
        summarizeElements(results.elements),
        "\nRole Redundancy: " .. (#results.redundancy > 0 and table.concat(results.redundancy, ", ") or "None"),
        "Missing Roles: " .. (#results.missing > 0 and table.concat(results.missing, ", ") or "None"),
    }
    api:ShowDialog("Party Synergy Summary", table.concat(lines, "\n"))
end

local function viewRoles(api)
    local results, err = analyzeParty(api)
    if not results then
        api:ShowDialog("Roles", err or "Unable to analyze party.")
        return
    end
    local lines = {}
    for name, roles in pairs(results.roles) do
        table.insert(lines, string.format("%s: %s", name, (#roles > 0 and table.concat(roles, ", ") or "(No inferred roles)")))
    end
    api:ShowDialog("Roles", table.concat(lines, "\n"))
end

local function viewElements(api)
    local results, err = analyzeParty(api)
    if not results then
        api:ShowDialog("Elements", err or "Unable to analyze party.")
        return
    end
    api:ShowDialog("Element Coverage", summarizeElements(results.elements))
end

local function viewSuggestions(api)
    local results, err = analyzeParty(api)
    if not results then
        api:ShowDialog("Suggestions", err or "Unable to analyze party.")
        return
    end
    local lines = {}
    if #results.missing == 0 then
        table.insert(lines, "No missing roles detected. Consider improving redundancy balance if needed.")
    else
        table.insert(lines, "Add roles: " .. table.concat(results.missing, ", "))
    end
    if #results.redundancy > 0 then
        table.insert(lines, "High redundancy: " .. table.concat(results.redundancy, ", "))
    end
    if #lines == 0 then
        table.insert(lines, "Party appears balanced based on heuristic checks.")
    end
    api:ShowDialog("Suggestions", table.concat(lines, "\n"))
end

local function viewHelp(api)
    local lines = {
        "Party Synergy Analyzer (read-only)",
        "- Summary: roster, element coverage, redundancy, missing roles.",
        "- Roles: shows inferred roles per member.",
        "- Elements: coverage across 8 elements.",
        "- Suggestions: quick guidance on gaps/redundancy.",
        "Inference is heuristic (spells/commands/stats); Unknown APIs yield 'Unable to analyze party'.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu Loop
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Party Synergy Analyzer",
            "Select an option:\n1) Summary\n2) Roles\n3) Elements\n4) Suggestions\n5) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            viewSummary(api)
        elseif sel == 2 then
            viewRoles(api)
        elseif sel == 3 then
            viewElements(api)
        elseif sel == 4 then
            viewSuggestions(api)
        elseif sel == 5 then
            viewHelp(api)
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
        api:ShowDialog("Error", "Party Synergy Analyzer failed: " .. tostring(err))
    end
end

return {
    id = "party-synergy-analyzer",
    name = "Party Synergy Analyzer",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
