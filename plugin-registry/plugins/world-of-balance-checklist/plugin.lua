--[[
    World of Balance Checklist
    Version: 1.0.0

    Tracks missable content before the World of Ruin transition: character recruitment,
    espers, missable items/events, and key flags. Designed as a read-only planning tool
    with graceful fallbacks when APIs are unavailable.

    Permissions required:
    - read_save: read roster, inventory, espers, and flags
    - ui_display: render dialogs and menus
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    assumeWorldOfBalanceWhenUnknown = true,   -- If world state cannot be detected, assume WoB for warnings
    showCompletionPercentages = true,         -- Display percentages per category
    warnPointOfNoReturn = true,               -- Show floating continent warning when still in WoB
    showMissableGuides = true,                -- Include short guidance blurbs per missable entry
    showEventFlagDetails = true,              -- Include event flag descriptions in event view
    maxListPerPage = 20,                      -- Pagination size for large lists
}

-- ============================================================================
-- Data Tables (reference-oriented; detection is dynamic where possible)
-- ============================================================================

local WOB_CHARACTERS = {
    {name = "Terra", category = "Story", note = "Joins during Narshe battle; mandatory"},
    {name = "Locke", category = "Story", note = "Joins early; ensure recruited before WoR"},
    {name = "Edgar", category = "Story", note = "Figaro king; available from start"},
    {name = "Sabin", category = "Story", note = "Rejoins after Lete River; bring him to Narshe battle"},
    {name = "Celes", category = "Story", note = "Mandatory in Narshe defense"},
    {name = "Cyan", category = "Story", note = "Recruited in the Imperial Camp / Phantom Train arc"},
    {name = "Gau", category = "Optional", note = "Leap on the Veldt; ensure recruited before WoR"},
    {name = "Setzer", category = "Story", note = "Recruited with the Blackjack before the Floating Continent"},
    {name = "Shadow", category = "Optional", note = "Ensure he survives Floating Continent; Coliseum in WoR"},
    {name = "Mog", category = "Optional", note = "Save him in Narshe to unlock dances"},
    {name = "Strago", category = "Optional", note = "Recruited in Thamasa"},
    {name = "Relm", category = "Optional", note = "Joins in Thamasa / Esper attack"},
    {name = "Umaro", category = "WoR", note = "Not available in WoB; listed for completeness"},
    {name = "Gogo", category = "WoR", note = "Not available in WoB; listed for completeness"},
}

-- Expected WoB-accessible espers (reference list; actual detection is dynamic)
local ESPER_TARGETS = {
    "Ramuh", "Kirin", "Siren", "Stray", "Ifrit", "Shiva", "Unicorn", "Maduin",
    "Carbuncle", "Bismarck", "Phantom", "Seraphim", "Golem", "ZoneSeek", "Shoat",
    "Tritoch", "Fenrir", "Terrato", "Starlet"
}

-- Missable content buckets; status is usually manual unless APIs exist
local MISSABLES = {
    {
        category = "Missable Items / Rewards",
        entries = {
            {name = "Charm Bangle", location = "South Figaro Cave (Merchant path)", note = "Get before WoR; boosts encounter control"},
            {name = "Tintinabar", location = "Returner's Hideout", note = "Talk to Arvis after Banon events"},
            {name = "Memento Ring", location = "Thamasa burning house", note = "Check Relm's house after the fire"},
            {name = "Hero Ring (Banquet)", location = "Imperial Banquet rewards", note = "Earn by optimal banquet responses"},
            {name = "Charm Bangle (Banquet)", location = "Imperial Banquet rewards", note = "Second reward path"},
        }
    },
    {
        category = "Rages to catch before WoR",
        entries = {
            {name = "Tyranosaur", location = "Dinosaur Forest (rare)", note = "Harder to encounter later"},
            {name = "Intangir", location = "Triangle Island", note = "Unavailable post-WoB"},
            {name = "Brontaur", location = "Narshe cliffs", note = "Great early Rage"},
        }
    },
    {
        category = "Dances to learn before WoR",
        entries = {
            {name = "Wind Rhapsody", location = "Any overworld grass", note = "Very early"},
            {name = "Forest Suite", location = "Forest tiles", note = "Learn before WoR for completeness"},
            {name = "Desert Aria", location = "Overworld desert", note = "Easier before WoR"},
            {name = "Snowman Jazz", location = "Narshe snowfield", note = "Use before the WoR freeze"},
        }
    },
    {
        category = "Key Events",
        entries = {
            {name = "Opera Success", location = "Opera House", note = "Ensure the show succeeds"},
            {name = "Imperial Banquet", location = "Vector", note = "Complete with optimal answers"},
            {name = "Magitek Factory Cleared", location = "Vector", note = "Finish the espers escape"},
        }
    }
}

local EVENT_FLAGS = {
    {name = "Opera Completed", flag = "event_opera_complete", description = "Indicates the opera sequence was completed successfully."},
    {name = "Banquet Outcome", flag = "event_banquet_score", description = "Score from the imperial banquet responses."},
    {name = "Floating Continent Reached", flag = "event_floating_continent", description = "Whether the party reached the Floating Continent."},
    {name = "WoR Transition", flag = "world_state", description = "Set when entering the World of Ruin."},
    {name = "Magitek Facility Cleared", flag = "event_magitek_factory", description = "Set after defeating Number 024 and leaving the facility."},
}

-- ============================================================================
-- Helper functions
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
    if not ok then
        return nil
    end
    return result
end

local function toLower(str)
    if not str then return "" end
    return string.lower(tostring(str))
end

local function getWorldState(api)
    local world = safeCall(api, "GetWorldState")
    if world then return world end
    local flagState = safeCall(api, "GetFlag", "world_state")
    if flagState then return flagState end
    return config.assumeWorldOfBalanceWhenUnknown and "World of Balance" or "Unknown"
end

local function fetchCharacters(api)
    local chars = {}
    for i = 0, 15 do
        local data = safeCall(api, "GetCharacter", i)
        if data then
            table.insert(chars, {index = i, data = data})
        end
    end
    return chars
end

local function characterName(char)
    return (char and (char.Name or char.name)) or "Unknown"
end

local function isCharacterRecruited(char)
    if not char then return false end
    if char.Enabled ~= nil then return char.Enabled end
    if char.enabled ~= nil then return char.enabled end
    if char.Available ~= nil then return char.Available end
    return false
end

local function findCharacter(characters, targetName)
    local targetLower = toLower(targetName)
    for _, entry in ipairs(characters) do
        local c = entry.data
        if c then
            local nameLower = toLower(characterName(c))
            if nameLower == targetLower or string.find(nameLower, targetLower, 1, true) then
                return c
            end
        end
    end
    return nil
end

local function fetchEspers(api)
    local espers = safeCall(api, "GetEspers")
    if not espers then return {} end
    return espers
end

local function esperName(esper)
    return (esper and (esper.Name or esper.name)) or "Esper"
end

local function esperOwned(esper)
    if not esper then return false end
    if esper.Owned ~= nil then return esper.Owned end
    if esper.owned ~= nil then return esper.owned end
    if esper.Acquired ~= nil then return esper.Acquired end
    return false
end

local function isEsperAcquired(esperList, targetName)
    local targetLower = toLower(targetName)
    for _, e in ipairs(esperList) do
        local nameLower = toLower(esperName(e))
        if nameLower == targetLower or string.find(nameLower, targetLower, 1, true) then
            return esperOwned(e)
        end
    end
    return false
end

local function formatPercent(current, total)
    if not total or total == 0 then return "0%" end
    local pct = (current / total) * 100
    return string.format("%.1f%%", pct)
end

local function paginate(list, pageSize)
    local pages = {}
    local size = pageSize or config.maxListPerPage
    local current = {}
    for i, item in ipairs(list) do
        table.insert(current, item)
        if #current >= size or i == #list then
            table.insert(pages, table.concat(current, "\n"))
            current = {}
        end
    end
    return pages
end

-- ============================================================================
-- Renderers
-- ============================================================================

local function renderCharacterChecklist(api)
    local chars = fetchCharacters(api)
    local lines = {}
    local recruited = 0
    for _, entry in ipairs(WOB_CHARACTERS) do
        local charData = findCharacter(chars, entry.name)
        local status = isCharacterRecruited(charData)
        if status then recruited = recruited + 1 end
        local availability = (entry.category == "WoR") and "(WoR)" or ""
        table.insert(lines, string.format("- [%s] %s %s - %s", status and "X" or " ", entry.name, availability, entry.note))
    end
    local header = string.format("World of Balance Character Checklist\nRecruited: %d/%d%s", recruited, #WOB_CHARACTERS, config.showCompletionPercentages and (" (" .. formatPercent(recruited, #WOB_CHARACTERS) .. ")") or "")
    local pages = paginate(lines, config.maxListPerPage)
    for i, page in ipairs(pages) do
        api:ShowDialog("WoB Characters (Page " .. i .. "/" .. #pages .. ")", header .. "\n\n" .. page)
    end
end

local function renderEsperChecklist(api)
    local espers = fetchEspers(api)
    local lines = {}
    local owned = 0
    for _, name in ipairs(ESPER_TARGETS) do
        local has = isEsperAcquired(espers, name)
        if has then owned = owned + 1 end
        table.insert(lines, string.format("- [%s] %s", has and "X" or " ", name))
    end
    local header = string.format("Esper Targets (WoB focus)\nAcquired: %d/%d%s", owned, #ESPER_TARGETS, config.showCompletionPercentages and (" (" .. formatPercent(owned, #ESPER_TARGETS) .. ")") or "")
    local pages = paginate(lines, config.maxListPerPage)
    for i, page in ipairs(pages) do
        api:ShowDialog("Esper Checklist (Page " .. i .. "/" .. #pages .. ")", header .. "\n\n" .. page)
    end
end

local function renderMissables(api)
    local sections = {}
    for _, bucket in ipairs(MISSABLES) do
        table.insert(sections, "[" .. bucket.category .. "]")
        for _, entry in ipairs(bucket.entries) do
            local status = "?"
            local line = string.format("- [%s] %s (Location: %s)", status, entry.name, entry.location)
            if config.showMissableGuides and entry.note then
                line = line .. "\n  Tip: " .. entry.note
            end
            table.insert(sections, line)
        end
        table.insert(sections, "")
    end
    local pages = paginate(sections, config.maxListPerPage)
    for i, page in ipairs(pages) do
        api:ShowDialog("Missable Content (Page " .. i .. "/" .. #pages .. ")", page)
    end
end

local function renderEvents(api)
    local lines = {}
    for _, ev in ipairs(EVENT_FLAGS) do
        local val = safeCall(api, "GetFlag", ev.flag)
        local status = (val == nil) and "Unknown" or tostring(val)
        local line = string.format("- %s: %s", ev.name, status)
        if config.showEventFlagDetails and ev.description then
            line = line .. "\n  " .. ev.description
        end
        table.insert(lines, line)
    end
    local pages = paginate(lines, config.maxListPerPage)
    for i, page in ipairs(pages) do
        api:ShowDialog("Event Flags (Page " .. i .. "/" .. #pages .. ")", page)
    end
end

local function renderSummary(api)
    local world = getWorldState(api)
    local chars = fetchCharacters(api)
    local espers = fetchEspers(api)

    local recruited = 0
    for _, entry in ipairs(WOB_CHARACTERS) do
        local charData = findCharacter(chars, entry.name)
        if isCharacterRecruited(charData) then recruited = recruited + 1 end
    end

    local owned = 0
    for _, name in ipairs(ESPER_TARGETS) do
        if isEsperAcquired(espers, name) then owned = owned + 1 end
    end

    local summary = {
        string.format("World state: %s", world or "Unknown"),
        string.format("Characters: %d/%d%s", recruited, #WOB_CHARACTERS, config.showCompletionPercentages and (" (" .. formatPercent(recruited, #WOB_CHARACTERS) .. ")") or ""),
        string.format("Espers: %d/%d%s", owned, #ESPER_TARGETS, config.showCompletionPercentages and (" (" .. formatPercent(owned, #ESPER_TARGETS) .. ")") or ""),
        "Missables: manual status (see missable view)",
    }

    if config.warnPointOfNoReturn and (world == "World of Balance" or world == "Unknown") then
        table.insert(summary, "\nWarning: Ensure all WoB missables are collected before progressing past the Floating Continent.")
    end

    api:ShowDialog("WoB Checklist Summary", table.concat(summary, "\n"))
end

local function renderExport(api)
    local world = getWorldState(api)
    local chars = fetchCharacters(api)
    local espers = fetchEspers(api)

    local lines = {}
    table.insert(lines, "World of Balance Checklist Export")
    table.insert(lines, "World: " .. (world or "Unknown"))
    table.insert(lines, "")
    table.insert(lines, "Characters:")
    for _, entry in ipairs(WOB_CHARACTERS) do
        local c = findCharacter(chars, entry.name)
        local status = isCharacterRecruited(c)
        table.insert(lines, string.format("- [%s] %s", status and "X" or " ", entry.name))
    end
    table.insert(lines, "")
    table.insert(lines, "Espers:")
    for _, name in ipairs(ESPER_TARGETS) do
        local has = isEsperAcquired(espers, name)
        table.insert(lines, string.format("- [%s] %s", has and "X" or " ", name))
    end
    table.insert(lines, "")
    table.insert(lines, "Missables:")
    for _, bucket in ipairs(MISSABLES) do
        table.insert(lines, "[" .. bucket.category .. "]")
        for _, entry in ipairs(bucket.entries) do
            table.insert(lines, string.format("- [?] %s (Location: %s)", entry.name, entry.location))
        end
    end

    api:ShowDialog("Export (copyable)", table.concat(lines, "\n"))
end

local function renderHelp(api)
    local help = {
        "World of Balance Checklist",
        "- Summary: Quick dashboard for characters/espers and a floating continent warning.",
        "- Characters: Shows recruitment status using your save data when available.",
        "- Espers: Marks expected WoB-accessible espers as acquired if found in your esper list.",
        "- Missables: Reference-only list for items/rages/dances/events before WoR.",
        "- Events: Reads flag values when available; otherwise shows Unknown.",
        "Notes:",
        "- If world state cannot be detected, the plugin assumes World of Balance (configurable).",
        "- No write operations are performed; this plugin is read-only.",
    }
    api:ShowDialog("Help", table.concat(help, "\n"))
end

-- ============================================================================
-- Menu Loop
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "World of Balance Checklist",
            "Select an option:\n1) Summary\n2) Characters\n3) Espers\n4) Missables\n5) Event Flags\n6) Export (copyable)\n7) Help\n0) Exit",
            "1"
        )

        if not choice then return end
        local selection = tonumber(choice)
        if selection == 1 then
            renderSummary(api)
        elseif selection == 2 then
            renderCharacterChecklist(api)
        elseif selection == 3 then
            renderEsperChecklist(api)
        elseif selection == 4 then
            renderMissables(api)
        elseif selection == 5 then
            renderEvents(api)
        elseif selection == 6 then
            renderExport(api)
        elseif selection == 7 then
            renderHelp(api)
        elseif selection == 0 then
            return
        else
            api:ShowDialog("Invalid", "Please choose a valid option.")
        end
    end
end

-- ============================================================================
-- Entry Point with safety wrapper
-- ============================================================================

local function safeMain(api)
    local ok, err = pcall(function()
        mainMenu(api)
    end)
    if not ok then
        if api and api.ShowDialog then
            api:ShowDialog("Error", "World of Balance Checklist failed: " .. tostring(err))
        end
    end
end

return {
    id = "world-of-balance-checklist",
    name = "World of Balance Checklist",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
