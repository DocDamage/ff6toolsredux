--[[
    Treasure Chest Tracker
    Version: 1.0.0

    Tracks opened vs. unopened treasure chests across both worlds with filtering,
    missable warnings, search, and exportable checklists. Designed to degrade
    gracefully when chest flag APIs are missing: unknown statuses are surfaced
    plainly instead of guessed.

    Permissions required:
    - read_save: read chest flags and world state
    - ui_display: render dialogs and menus
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    pageSize = 20,                    -- Items per page in list views
    assumeUnknownAsUnopened = false,  -- If true, treat unknown as unopened for planning
    showValueEstimates = true,        -- Show gil value estimates when available
    includeWorldFilterInSummary = true,
    defaultWorldFallback = "World of Balance", -- Used when world state cannot be detected
    highlightMissables = true,        -- Prefix missable chests with ! marker
}

-- ============================================================================
-- Chest Database (representative subset; expand as APIs allow)
-- ============================================================================
-- Fields: id (string/number), world ("WoB"|"WoR"|"Both"), region, area, item, value, missable, note
local CHESTS = {
    {id = "NARSHE-CLIFF-01", world = "WoB", region = "Narshe", area = "Cliffs", item = "Elixir", value = 5000, missable = true, note = "Pre-raid only; vanishes after WoR"},
    {id = "NARSHE-MINE-02", world = "WoB", region = "Narshe", area = "Mines", item = "Phoenix Down", value = 300, missable = false, note = "Reusable path later"},
    {id = "FIGARO-CASTLE-01", world = "Both", region = "Figaro", area = "Castle Basement", item = "Ether", value = 800, missable = false, note = "Accessible after submerge"},
    {id = "SOUTH-FIGARO-03", world = "WoB", region = "South Figaro", area = "Rich Man's House", item = "Sprint Shoes", value = 1200, missable = false, note = "During infiltration"},
    {id = "RETURNERS-01", world = "WoB", region = "Returner Hideout", area = "Storage", item = "Green Cherry", value = 150, missable = false, note = "Easy pickup"},
    {id = "LETE-RIVER-01", world = "WoB", region = "Lete River", area = "River Path", item = "White Cape", value = 1000, missable = true, note = "One-time river route"},
    {id = "OPERA-01", world = "WoB", region = "Opera House", area = "Backstage", item = "X-Potion", value = 1500, missable = true, note = "Before Ultros event"},
    {id = "VECTOR-ARMORY-01", world = "WoB", region = "Vector", area = "Armory", item = "Flame Sabre", value = 3000, missable = true, note = "Before Magitek Factory"},
    {id = "MAGITEK-FACTORY-02", world = "WoB", region = "Magitek Facility", area = "Rails", item = "Thunder Blade", value = 3200, missable = true, note = "Factory escape route"},
    {id = "MAGITEK-FACTORY-03", world = "WoB", region = "Magitek Facility", area = "Rails", item = "Poison Rod", value = 1800, missable = true, note = "Factory escape route"},
    {id = "SEALED-GATE-01", world = "WoB", region = "Sealed Gate", area = "Caves", item = "Chocobo Suit", value = 2500, missable = true, note = "Before WoR collapse"},
    {id = "SEALED-GATE-02", world = "WoB", region = "Sealed Gate", area = "Caves", item = "Barrier Ring", value = 1200, missable = true, note = "Before WoR collapse"},
    {id = "FLOATING-CONT-01", world = "WoB", region = "Floating Continent", area = "Ridge", item = "Genji Glove", value = 6000, missable = true, note = "Point-of-no-return"},
    {id = "FLOATING-CONT-02", world = "WoB", region = "Floating Continent", area = "Ridge", item = "Murasame", value = 5200, missable = true, note = "Point-of-no-return"},
    {id = "THAMASA-01", world = "WoB", region = "Thamasa", area = "Houses", item = "Earring", value = 1500, missable = false, note = "Pre-Esper attack"},
    {id = "PHOENIX-CAVE-01", world = "WoR", region = "Phoenix Cave", area = "Upper Path", item = "Ribbon", value = 8000, missable = false, note = "Requires split parties"},
    {id = "PHOENIX-CAVE-02", world = "WoR", region = "Phoenix Cave", area = "Lava Path", item = "Flame Shield", value = 7000, missable = false, note = "Split path"},
    {id = "ANCIENT-CASTLE-01", world = "WoR", region = "Ancient Castle", area = "Throne", item = "Offering", value = 10000, missable = false, note = "Hidden passage"},
    {id = "ANCIENT-CASTLE-02", world = "WoR", region = "Ancient Castle", area = "Basement", item = "Enhancer", value = 7500, missable = false, note = "Hidden switch"},
    {id = "PHOENIX-TOWER-01", world = "WoR", region = "Kefka's Tower", area = "Path A", item = "Atma Weapon", value = 12000, missable = false, note = "Final dungeon"},
    {id = "PHOENIX-TOWER-02", world = "WoR", region = "Kefka's Tower", area = "Path B", item = "Ragnarok Sword", value = 15000, missable = false, note = "Final dungeon"},
    {id = "PHOENIX-TOWER-03", world = "WoR", region = "Kefka's Tower", area = "Path C", item = "Minerva", value = 11000, missable = false, note = "Final dungeon"},
    {id = "VECTORSIDE-RUINS-01", world = "WoR", region = "Vector Ruins", area = "Factory Ruins", item = "Force Armor", value = 9000, missable = false, note = "WoR revisit"},
    {id = "DARRYL-TOMB-01", world = "WoR", region = "Darryl's Tomb", area = "Main", item = "Exp. Egg", value = 9999, missable = false, note = "Switch puzzle"},
    {id = "DARRYL-TOMB-02", world = "WoR", region = "Darryl's Tomb", area = "Side", item = "Crystal Mail", value = 6500, missable = false, note = "Optional wing"},
    {id = "OWZER-MANSION-01", world = "WoR", region = "Owzer's Mansion", area = "Gallery", item = "Soul Sabre", value = 4000, missable = false, note = "Pre-painting battle"},
    {id = "OWZER-MANSION-02", world = "WoR", region = "Owzer's Mansion", area = "Gallery", item = "Thornlet", value = 3500, missable = false, note = "Pre-painting battle"},
    {id = "ESPER-CAVE-01", world = "WoR", region = "Esper Cave", area = "Entrance", item = "Ice Shield", value = 7000, missable = false, note = "Side path"},
    {id = "ESPER-CAVE-02", world = "WoR", region = "Esper Cave", area = "Deep", item = "Ribbon", value = 8000, missable = false, note = "Side path"},
    {id = "COLLOSEUM-01", world = "WoR", region = "Colosseum", area = "Arena", item = "X-Ether", value = 2000, missable = false, note = "Before betting"},
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

local function getWorld(api)
    local world = safeCall(api, "GetWorldState")
    if world then return world end
    local flagWorld = safeCall(api, "GetFlag", "world_state")
    if flagWorld then return flagWorld end
    return config.defaultWorldFallback or "Unknown"
end

local function isChestOpened(api, chestId)
    -- Try dedicated chest API
    local opened = safeCall(api, "GetChestOpened", chestId)
    if opened ~= nil then return opened end
    -- Fallback: generic flag lookup
    local flagVal = safeCall(api, "GetFlag", chestId)
    if flagVal ~= nil then return flagVal end
    return nil -- Unknown
end

local function formatValue(val)
    if not val or not config.showValueEstimates then return "" end
    return string.format(" (est. %d gil)", val)
end

local function statusIcon(status, missable)
    local base = status == true and "[X]" or (status == false and "[ ]" or "[?]")
    if missable and config.highlightMissables then
        return "!" .. base
    end
    return base
end

local function filteredChests(api, opts)
    opts = opts or {}
    local world = opts.world or nil
    local region = opts.region and toLower(opts.region) or nil
    local unopenedOnly = opts.unopenedOnly or false
    local missablesOnly = opts.missablesOnly or false
    local search = opts.search and toLower(opts.search) or nil

    local results = {}
    for _, chest in ipairs(CHESTS) do
        if (not world or chest.world == world or chest.world == "Both") and
           (not region or string.find(toLower(chest.region), region, 1, true) or string.find(toLower(chest.area), region, 1, true)) then
            local status = isChestOpened(api, chest.id)
            local consideredStatus = status
            if consideredStatus == nil and config.assumeUnknownAsUnopened then
                consideredStatus = false
            end

            if missablesOnly and not chest.missable then
                goto continue
            end
            if unopenedOnly and consideredStatus == true then
                goto continue
            end
            if search then
                local haystack = table.concat({chest.region, chest.area, chest.item, chest.note or ""}, " ")
                if not string.find(toLower(haystack), search, 1, true) then
                    goto continue
                end
            end

            table.insert(results, {
                chest = chest,
                status = status,
                consideredStatus = consideredStatus,
            })
        end
        ::continue::
    end
    return results
end

local function paginate(list, size)
    local pages = {}
    local buf = {}
    size = size or config.pageSize
    for i, entry in ipairs(list) do
        table.insert(buf, entry)
        if #buf >= size or i == #list then
            table.insert(pages, buf)
            buf = {}
        end
    end
    return pages
end

local function formatChestLine(entry)
    local c = entry.chest
    local s = entry.consideredStatus
    local icon = statusIcon(s, c.missable)
    local note = c.note and (" - " .. c.note) or ""
    local val = formatValue(c.value)
    return string.format("%s %s / %s / %s%s%s", icon, c.region, c.area, c.item, val, note)
end

local function summarize(api)
    local world = getWorld(api)
    local total = #CHESTS
    local opened, unopened, unknown = 0, 0, 0
    for _, chest in ipairs(CHESTS) do
        local status = isChestOpened(api, chest.id)
        if status == true then opened = opened + 1
        elseif status == false then unopened = unopened + 1
        else unknown = unknown + 1 end
    end
    local lines = {
        string.format("World detected: %s", world),
        string.format("Chests tracked: %d", total),
        string.format("Opened: %d", opened),
        string.format("Unopened: %d", unopened),
        string.format("Unknown: %d", unknown),
        "Use Unopened/Missables views to drill down. Unknown means the API was not available for that chest ID.",
    }
    return table.concat(lines, "\n")
end

local function renderList(api, title, entries)
    if #entries == 0 then
        api:ShowDialog(title, "No entries match your filters.")
        return
    end
    local pages = paginate(entries, config.pageSize)
    for i, page in ipairs(pages) do
        local lines = {}
        for _, entry in ipairs(page) do
            table.insert(lines, formatChestLine(entry))
        end
        api:ShowDialog(string.format("%s (Page %d/%d)", title, i, #pages), table.concat(lines, "\n"))
    end
end

local function menuBrowseByRegion(api)
    -- Build region list
    local regions = {}
    local seen = {}
    for _, chest in ipairs(CHESTS) do
        if not seen[chest.region] then
            table.insert(regions, chest.region)
            seen[chest.region] = true
        end
    end
    table.sort(regions)

    local prompt = {"Select a region:"}
    for i, r in ipairs(regions) do
        table.insert(prompt, string.format("%d) %s", i, r))
    end
    table.insert(prompt, "0) Cancel")

    local choice = api:ShowInput("Browse by Region", table.concat(prompt, "\n"), "1")
    if not choice then return end
    local sel = tonumber(choice)
    if sel == nil or sel < 0 or sel > #regions then
        api:ShowDialog("Invalid", "Please choose a valid option.")
        return
    end
    if sel == 0 then return end

    local region = regions[sel]
    local entries = filteredChests(api, {region = region})
    renderList(api, "Region: " .. region, entries)
end

local function menuSearch(api)
    local term = api:ShowInput("Search", "Enter item, area, or region keyword:", "")
    if not term or term == "" then return end
    local entries = filteredChests(api, {search = term})
    renderList(api, "Search: " .. term, entries)
end

local function menuUnopened(api)
    local entries = filteredChests(api, {unopenedOnly = true})
    renderList(api, "Unopened Chests", entries)
end

local function menuMissables(api)
    local entries = filteredChests(api, {missablesOnly = true})
    renderList(api, "Missable Chests (reference)", entries)
end

local function menuWorldFilter(api)
    local world = getWorld(api)
    local entries = filteredChests(api, {world = world})
    renderList(api, string.format("Chests in %s", world), entries)
end

local function menuExport(api)
    local entries = filteredChests(api, {})
    local lines = {"Treasure Chest Tracker Export"}
    table.insert(lines, string.format("World: %s", getWorld(api)))
    table.insert(lines, "")
    for _, entry in ipairs(entries) do
        table.insert(lines, formatChestLine(entry))
    end
    api:ShowDialog("Export (copyable)", table.concat(lines, "\n"))
end

local function menuHelp(api)
    local lines = {
        "Treasure Chest Tracker",
        "- Summary: Shows counts of opened/unopened/unknown.",
        "- Unopened: Filters to chests not yet opened (Unknown can be treated as unopened via config).",
        "- Missables: Flags critical WoB-only or point-of-no-return chests.",
        "- World Filter: Uses detected world to narrow the list (WoB/WoR/Both).",
        "- Region Browse: Pick a region to list its chests.",
        "- Search: Keyword match on item/area/region/note.",
        "- Export: Copyable text of all tracked chests.",
        "Notes: Unknown means the API did not return a value for that chest ID. The plugin will never guess.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu Loop
-- ============================================================================
local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Treasure Chest Tracker",
            "Select an option:\n1) Summary\n2) Unopened\n3) Missables\n4) World Filter\n5) Browse by Region\n6) Search\n7) Export\n8) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            api:ShowDialog("Summary", summarize(api))
        elseif sel == 2 then
            menuUnopened(api)
        elseif sel == 3 then
            menuMissables(api)
        elseif sel == 4 then
            menuWorldFilter(api)
        elseif sel == 5 then
            menuBrowseByRegion(api)
        elseif sel == 6 then
            menuSearch(api)
        elseif sel == 7 then
            menuExport(api)
        elseif sel == 8 then
            menuHelp(api)
        elseif sel == 0 then
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
    if not ok and api and api.ShowDialog then
        api:ShowDialog("Error", "Treasure Chest Tracker failed: " .. tostring(err))
    end
end

return {
    id = "treasure-chest-tracker",
    name = "Treasure Chest Tracker",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
