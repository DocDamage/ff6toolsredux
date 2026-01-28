--[[
    Colosseum Guide
    Version: 1.0.0

    Interactive betting reference: maps bet item -> opponent -> reward, with enemy notes,
    filtering, search, and export. Read-only and resilient to missing APIs: unknown
    inventory or bestiary data is surfaced rather than guessed.

    Permissions: read_save, ui_display
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    pageSize = 15,
    showEnemyStats = true,
    showStrategyNotes = true,
    treatUnknownInventoryAsMissing = false, -- When inventory API is missing, treat unknown as missing
}

-- ============================================================================
-- Colosseum Database (representative core; extend as needed)
-- ============================================================================
-- Fields: bet, reward, opponent, rewardNotes, strategy, enemyId (optional), tier
local BETS = {
    {bet = "Dirk", reward = "Aura Lance", opponent = "Borras", enemyId = 2001, tier = "weapon", strategy = "Physical; ice works.", rewardNotes = "High-power lance"},
    {bet = "Aura Lance", reward = "Trump", opponent = "Didalos", enemyId = 2002, tier = "weapon", strategy = "Weak to fire.", rewardNotes = "Throwing weapon"},
    {bet = "Genji Glove", reward = "Master's Scroll", opponent = "Chupon", enemyId = 2003, tier = "relic", strategy = "Beware sneeze; high evasion helps.", rewardNotes = "Dual-wield relic"},
    {bet = "Ragnarok (Sword)", reward = "Lightbringer", opponent = "Didalos", enemyId = 2002, tier = "weapon", strategy = "Haste + Shell; ice/fire both fine.", rewardNotes = "Best sword"},
    {bet = "Minerva", reward = "Force Armor", opponent = "Tyranosaur", enemyId = 119, tier = "armor", strategy = "Ice weak; beware Meteo.", rewardNotes = "High def/mdef"},
    {bet = "Flame Shield", reward = "Ice Shield", opponent = "Chimera", enemyId = 90, tier = "shield", strategy = "Ice weak; null earth.", rewardNotes = "Elemental swap"},
    {bet = "Thunder Shield", reward = "Flame Shield", opponent = "Chimera", enemyId = 90, tier = "shield", strategy = "Ice weak; null earth.", rewardNotes = "Elemental swap"},
    {bet = "Ice Shield", reward = "Aegis Shield", opponent = "Chimera", enemyId = 90, tier = "shield", strategy = "Ice weak; null earth.", rewardNotes = "Evade/MEvade"},
    {bet = "Cat Hood", reward = "Merit Award", opponent = "Coelecite", enemyId = 2004, tier = "relic", strategy = "Low threat; burst it.", rewardNotes = "Equip anything"},
    {bet = "Exp. Egg", reward = "Marvel Shoes", opponent = "Steroidite", enemyId = 2005, tier = "relic", strategy = "Slow + def; bring strong magic.", rewardNotes = "Haste+Regen+Shell+Safe"},
    {bet = "Merit Award", reward = "Cat Hood", opponent = "Coelecite", enemyId = 2004, tier = "relic", strategy = "Low threat; burst it.", rewardNotes = "Gil boost"},
    {bet = "Charm Bangle", reward = "Back Guard", opponent = "Scorpion", enemyId = 2006, tier = "relic", strategy = "Bolt weak.", rewardNotes = "Back attack immunity"},
    {bet = "Tintinabar", reward = "Thornlet", opponent = "Gloomwind", enemyId = 2007, tier = "relic", strategy = "Wind damage; silence helps.", rewardNotes = "HP regen relic"},
    {bet = "Elixir", reward = "Rename Card", opponent = "Cactuar", enemyId = 2008, tier = "misc", strategy = "1 HP; any hit.", rewardNotes = "Utility"},
    {bet = "Rename Card", reward = "Marvel Shoes", opponent = "Steroidite", enemyId = 2005, tier = "relic", strategy = "Slow + def; strong magic.", rewardNotes = "Buff relic"},
    {bet = "Gauntlet", reward = "Genji Glove", opponent = "Steroidite", enemyId = 2005, tier = "relic", strategy = "As above.", rewardNotes = "Dual-wield swap"},
}

-- Enemy hint overrides (lightweight, avoids full bestiary import)
local ENEMY_HINTS = {
    [119] = {name = "Tyranosaur", weak = "Ice", note = "High HP; beware Meteor."},
    [90] = {name = "Chimera", weak = "Ice", note = "Null Earth; uses Aqua Rake."},
    [2001] = {name = "Borras", weak = "Ice", note = "Standard physical."},
    [2002] = {name = "Didalos", weak = "Fire/Ice", note = "Haste + Safe recommended."},
    [2003] = {name = "Chupon", weak = "None", note = "Can Sneeze; Stamina helps."},
    [2004] = {name = "Coelecite", weak = "Fire", note = "Very low threat."},
    [2005] = {name = "Steroidite", weak = "None", note = "High defense; use strong magic."},
    [2006] = {name = "Scorpion", weak = "Bolt", note = "Poison claws; fast burst."},
    [2007] = {name = "Gloomwind", weak = "Fire", note = "Wind attacks; silence to disable."},
    [2008] = {name = "Cactuar", weak = "Any", note = "1 HP; any hit."},
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

local function hasInventoryItem(api, name)
    local inv = safeCall(api, "GetInventory")
    if not inv then
        return config.treatUnknownInventoryAsMissing and false or nil
    end
    local target = toLower(name)
    for _, item in ipairs(inv) do
        local n = item.Name or item.name
        if n and toLower(n) == target and (item.Quantity or item.quantity or 0) > 0 then
            return true
        end
    end
    return false
end

local function getEnemyHint(enemyId)
    return ENEMY_HINTS[enemyId]
end

local function filterBets(opts, api)
    opts = opts or {}
    local tier = opts.tier and toLower(opts.tier)
    local search = opts.search and toLower(opts.search)
    local ownedOnly = opts.ownedOnly or false

    local results = {}
    for _, entry in ipairs(BETS) do
        if tier and toLower(entry.tier or "") ~= tier then goto continue end
        if search then
            local hay = table.concat({entry.bet, entry.reward, entry.opponent or "", entry.strategy or "", entry.rewardNotes or ""}, " ")
            if not string.find(toLower(hay), search, 1, true) then
                goto continue
            end
        end
        if ownedOnly then
            local has = hasInventoryItem(api, entry.bet)
            if has == false then goto continue end
        end
        table.insert(results, entry)
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

local function formatEntry(entry)
    local hint = entry.enemyId and getEnemyHint(entry.enemyId)
    local lines = {}
    table.insert(lines, string.format("Bet: %s -> Reward: %s", entry.bet, entry.reward))
    table.insert(lines, string.format("Opponent: %s", entry.opponent))
    if hint and config.showEnemyStats then
        table.insert(lines, string.format("Enemy Hint: %s (Weak: %s)%s", hint.name or "Enemy", hint.weak or "-", hint.note and (" - " .. hint.note) or ""))
    end
    if entry.strategy and config.showStrategyNotes then
        table.insert(lines, "Strategy: " .. entry.strategy)
    end
    if entry.rewardNotes then
        table.insert(lines, "Reward Notes: " .. entry.rewardNotes)
    end
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
        for _, e in ipairs(page) do
            table.insert(lines, formatEntry(e))
            table.insert(lines, "")
        end
        api:ShowDialog(string.format("%s (Page %d/%d)", title, i, #pages), table.concat(lines, "\n"))
    end
end

-- ============================================================================
-- Views
-- ============================================================================

local function viewSummary(api)
    local total = #BETS
    local lines = {
        string.format("Mapped bets: %d", total),
        "Filters: Owned-only, by tier, search. Export available.",
        "Unknown inventory is shown as Unknown (unless configured as missing).",
    }
    api:ShowDialog("Colosseum Summary", table.concat(lines, "\n"))
end

local function viewAll(api)
    renderList(api, "All Bets", BETS)
end

local function viewOwned(api)
    local entries = filterBets({ownedOnly = true}, api)
    renderList(api, "Bets You Own (or Unknown)", entries)
end

local function viewByTier(api)
    local tier = api:ShowInput("Filter by Tier", "Enter tier (weapon, armor, shield, relic, misc):", "weapon")
    if not tier then return end
    local entries = filterBets({tier = tier}, api)
    renderList(api, "Tier: " .. tier, entries)
end

local function viewSearch(api)
    local term = api:ShowInput("Search", "Keyword across bet/reward/opponent/strategy:", "")
    if not term or term == "" then return end
    local entries = filterBets({search = term}, api)
    renderList(api, "Search: " .. term, entries)
end

local function viewExport(api)
    local entries = filterBets({}, api)
    local lines = {"Colosseum Guide Export"}
    for _, e in ipairs(entries) do
        table.insert(lines, formatEntry(e))
        table.insert(lines, "")
    end
    api:ShowDialog("Export (copyable)", table.concat(lines, "\n"))
end

local function viewHelp(api)
    local lines = {
        "Colosseum Guide (read-only)",
        "- Summary: Counts and reminders.",
        "- All: Full mapping bet -> opponent -> reward.",
        "- Owned: Filters to bets you have (or Unknown if inventory API missing).",
        "- Tier: weapon/armor/shield/relic/misc filters.",
        "- Search: Keyword across bet/reward/opponent/strategy.",
        "- Export: Copyable text for routing notes.",
        "Status: No write operations; Unknown inventory shown if API absent.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu Loop
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Colosseum Guide",
            "Select an option:\n1) Summary\n2) All Bets\n3) Bets You Own\n4) By Tier\n5) Search\n6) Export\n7) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            viewSummary(api)
        elseif sel == 2 then
            viewAll(api)
        elseif sel == 3 then
            viewOwned(api)
        elseif sel == 4 then
            viewByTier(api)
        elseif sel == 5 then
            viewSearch(api)
        elseif sel == 6 then
            viewExport(api)
        elseif sel == 7 then
            viewHelp(api)
        elseif sel == 0 then
            return
        else
            api:ShowDialog("Invalid", "Please choose a valid option.")
        end
    end
end

-- ============================================================================
-- Entry Point
-- ============================================================================

local function safeMain(api)
    local ok, err = pcall(function()
        mainMenu(api)
    end)
    if not ok and api and api.ShowDialog then
        api:ShowDialog("Error", "Colosseum Guide failed: " .. tostring(err))
    end
end

return {
    id = "colosseum-guide",
    name = "Colosseum Guide",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
