--[[
    Bestiary
    Version: 1.0.0

    Comprehensive enemy reference with encounter tracking, weaknesses, drops/steals,
    and multi-criteria search. Degrades gracefully when encounter flags are missing
    by surfacing "Unknown" instead of guessing.

    Permissions:
    - read_save: read encounter flags
    - ui_display: render dialogs and menus
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    pageSize = 15,                     -- Entries per page
    assumeUnknownAsNotSeen = false,    -- Treat unknown encounter state as not seen
    showDropsAndSteals = true,
    showMorphAndSketch = true,
    showRageMoves = true,
    showLocations = true,
}

-- ============================================================================
-- Enemy Database (representative subset; expand with full dataset when available)
-- ============================================================================
-- Fields: id, name, level, hp, mp, type, elements(weak,resist,absorb,null), statusImmune, drop, steal,
-- morph, sketch, control, rageMove, location, boss
local ENEMIES = {
    {
        id = 1, name = "Guard", level = 6, hp = 45, mp = 0, type = "Humanoid", boss = false,
        elements = {weak = {"Poison"}, resist = {}, absorb = {}, null = {}},
        statusImmune = {"Sleep"},
        drop = {item = "Potion", chance = 0.35},
        steal = {common = "Potion", rare = "Antidote"},
        morph = {item = "Hi-Potion", chance = 0.40},
        sketch = "Attack",
        control = {"Attack", "Potion"},
        rageMove = "Fight",
        location = "Narshe Mines (intro)"
    },
    {
        id = 27, name = "Templar", level = 13, hp = 330, mp = 120, type = "Humanoid", boss = false,
        elements = {weak = {"Ice"}, resist = {"Fire"}, absorb = {}, null = {}},
        statusImmune = {"Stop"},
        drop = {item = "Hi-Potion", chance = 0.25},
        steal = {common = "Antidote", rare = "Hi-Potion"},
        morph = {item = "Ether", chance = 0.40},
        sketch = "Attack",
        control = {"Attack", "Cure"},
        rageMove = "Fire",
        location = "Imperial Camp / Vector"
    },
    {
        id = 52, name = "Vaporite", level = 19, hp = 320, mp = 180, type = "Elemental", boss = false,
        elements = {weak = {"Fire"}, resist = {"Ice"}, absorb = {"Lightning"}, null = {}},
        statusImmune = {"Poison"},
        drop = {item = "Ether", chance = 0.20},
        steal = {common = "Echo Screen", rare = "Ether"},
        morph = {item = "Hi-Ether", chance = 0.30},
        sketch = "Bolt",
        control = {"Attack", "Bolt"},
        rageMove = "Bolt",
        location = "Magitek Factory"
    },
    {
        id = 90, name = "Chimera", level = 28, hp = 2500, mp = 200, type = "Beast", boss = false,
        elements = {weak = {"Ice"}, resist = {"Fire"}, absorb = {}, null = {"Earth"}},
        statusImmune = {"Poison", "Sap"},
        drop = {item = "Gold Armor", chance = 0.12},
        steal = {common = "Remedy", rare = "Gold Helmet"},
        morph = {item = "Mithril Vest", chance = 0.25},
        sketch = "Aqua Rake",
        control = {"Aqua Rake", "Fireball"},
        rageMove = "Aqua Rake",
        location = "Esper Cave / WoR plains"
    },
    {
        id = 119, name = "Tyranosaur", level = 35, hp = 6700, mp = 600, type = "Dinosaur", boss = false,
        elements = {weak = {"Ice"}, resist = {"Fire"}, absorb = {}, null = {}},
        statusImmune = {"Instant Death", "Sleep"},
        drop = {item = "T-Rexaur Claw", chance = 0.12},
        steal = {common = "X-Potion", rare = "Ribbon"},
        morph = {item = "Megalixir", chance = 0.12},
        sketch = "Meteor",
        control = {"Meteor", "Fireball"},
        rageMove = "Meteor",
        location = "Dinosaur Forest (WoR)"
    },
    {
        id = 130, name = "Intangir", level = 39, hp = 32000, mp = 5000, type = "Special", boss = false,
        elements = {weak = {}, resist = {}, absorb = {}, null = {"Fire", "Ice", "Lightning", "Pearl", "Earth"}},
        statusImmune = {"Death", "Petrify", "Stop", "Sleep"},
        drop = {item = "Elixir", chance = 1.0},
        steal = {common = "Elixir", rare = "Elixir"},
        morph = {item = "Megalixir", chance = 0.25},
        sketch = "Blowfish",
        control = {"Blowfish"},
        rageMove = "Blowfish",
        location = "Triangle Island (WoB only)"
    },
    {
        id = 160, name = "Phunbaba", level = 41, hp = 28000, mp = 10000, type = "Boss", boss = true,
        elements = {weak = {}, resist = {"Lightning"}, absorb = {"Wind"}, null = {}},
        statusImmune = {"Poison", "Sap", "Stop", "Slow"},
        drop = {item = "None", chance = 0.0},
        steal = {common = "None", rare = "None"},
        morph = nil,
        sketch = "Baba Breath",
        control = nil,
        rageMove = "Baba Breath",
        location = "Mobliz (WoR)"
    },
    {
        id = 200, name = "Brontaur", level = 35, hp = 4200, mp = 350, type = "Beast", boss = false,
        elements = {weak = {"Ice"}, resist = {"Poison"}, absorb = {}, null = {}},
        statusImmune = {"Stop"},
        drop = {item = "Mythril Claw", chance = 0.20},
        steal = {common = "Green Cherry", rare = "Gold Shield"},
        morph = {item = "Gold Hairpin", chance = 0.18},
        sketch = "Snowstorm",
        control = {"Snowstorm", "Bite"},
        rageMove = "Snowstorm",
        location = "Narshe Cliffs (WoB)"
    },
    {
        id = 270, name = "Dragon", level = 57, hp = 18000, mp = 12000, type = "Dragon", boss = true,
        elements = {weak = {"Ice"}, resist = {"Poison"}, absorb = {}, null = {"Earth"}},
        statusImmune = {"Stop", "Slow", "Death"},
        drop = {item = "Crystal Helm", chance = 1.0},
        steal = {common = "Elixir", rare = "Elixir"},
        morph = nil,
        sketch = "Attack",
        control = nil,
        rageMove = "Attack",
        location = "Various Dragon battles"
    },
    {
        id = 310, name = "Kefka (Narshe)", level = 18, hp = 3000, mp = 300, type = "Boss", boss = true,
        elements = {weak = {"Poison"}, resist = {}, absorb = {}, null = {}},
        statusImmune = {"Sleep", "Stop", "Slow"},
        drop = {item = "None", chance = 0.0},
        steal = {common = "None", rare = "None"},
        morph = nil,
        sketch = "Ice",
        control = nil,
        rageMove = "Ice",
        location = "Narshe Defense (WoB)"
    },
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

local function hasValue(list, value)
    if not list or not value then return false end
    local v = toLower(value)
    for _, item in ipairs(list) do
        if toLower(item) == v then return true end
    end
    return false
end

local function isEncountered(api, enemyId)
    local seen = safeCall(api, "GetEncountered", enemyId)
    if seen ~= nil then return seen end
    local flag = safeCall(api, "GetFlag", "bestiary_" .. tostring(enemyId))
    if flag ~= nil then return flag end
    return nil
end

local function statusIcon(seen)
    if seen == true then return "[SEEN]" end
    if seen == false then return "[UNSEEN]" end
    return "[?]"
end

local function considerStatus(seen)
    if seen == nil and config.assumeUnknownAsNotSeen then
        return false
    end
    return seen
end

local function formatElements(el)
    local parts = {}
    if el.weak and #el.weak > 0 then table.insert(parts, "Weak: " .. table.concat(el.weak, ", ")) end
    if el.resist and #el.resist > 0 then table.insert(parts, "Resist: " .. table.concat(el.resist, ", ")) end
    if el.absorb and #el.absorb > 0 then table.insert(parts, "Absorb: " .. table.concat(el.absorb, ", ")) end
    if el.null and #el.null > 0 then table.insert(parts, "Null: " .. table.concat(el.null, ", ")) end
    return (#parts > 0) and table.concat(parts, " | ") or "None"
end

local function formatDrops(enemy)
    local lines = {}
    if config.showDropsAndSteals then
        if enemy.drop and enemy.drop.item then
            table.insert(lines, string.format("Drop: %s (%.0f%%)", enemy.drop.item, (enemy.drop.chance or 0) * 100))
        else
            table.insert(lines, "Drop: None")
        end
        if enemy.steal then
            table.insert(lines, string.format("Steal: Common %s / Rare %s", enemy.steal.common or "-", enemy.steal.rare or "-"))
        else
            table.insert(lines, "Steal: None")
        end
    end
    if config.showMorphAndSketch then
        if enemy.morph then
            table.insert(lines, string.format("Morph: %s (%.0f%%)", enemy.morph.item or "-", (enemy.morph.chance or 0) * 100))
        else
            table.insert(lines, "Morph: None")
        end
        table.insert(lines, "Sketch: " .. (enemy.sketch or "-"))
        table.insert(lines, "Control: " .. ((enemy.control and table.concat(enemy.control, ", ")) or "-"))
    end
    if config.showRageMoves then
        table.insert(lines, "Rage: " .. (enemy.rageMove or "-"))
    end
    return lines
end

local function formatEnemy(enemy, seen)
    local lines = {}
    table.insert(lines, string.format("%s #%d %s (Lv %d, HP %d, MP %d) [%s]", statusIcon(seen), enemy.id, enemy.name, enemy.level, enemy.hp, enemy.mp, enemy.type))
    if config.showLocations and enemy.location then
        table.insert(lines, "Location: " .. enemy.location)
    end
    table.insert(lines, "Elements: " .. formatElements(enemy.elements or {}))
    if enemy.statusImmune and #enemy.statusImmune > 0 then
        table.insert(lines, "Status Immune: " .. table.concat(enemy.statusImmune, ", "))
    end
    for _, l in ipairs(formatDrops(enemy)) do table.insert(lines, l) end
    return table.concat(lines, "\n")
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

-- ============================================================================
-- Filters
-- ============================================================================

local function filterEnemies(api, opts)
    opts = opts or {}
    local results = {}
    local typeFilter = opts.type and toLower(opts.type)
    local elementFilter = opts.element and toLower(opts.element)
    local search = opts.search and toLower(opts.search)
    local bossesOnly = opts.bossesOnly or false
    local unseenOnly = opts.unseenOnly or false

    for _, e in ipairs(ENEMIES) do
        if bossesOnly and not e.boss then goto continue end
        if typeFilter and toLower(e.type) ~= typeFilter then goto continue end
        if elementFilter then
            local el = e.elements or {weak = {}, resist = {}, absorb = {}, null = {}}
            local hit = hasValue(el.weak, elementFilter) or hasValue(el.resist, elementFilter) or hasValue(el.absorb, elementFilter) or hasValue(el.null, elementFilter)
            if not hit then goto continue end
        end
        if search then
            local hay = table.concat({e.name, e.type, e.location or "", table.concat(e.statusImmune or {}, " ")}, " ")
            if not string.find(toLower(hay), search, 1, true) then goto continue end
        end
        local seen = considerStatus(isEncountered(api, e.id))
        if unseenOnly and seen == true then goto continue end
        table.insert(results, {enemy = e, seen = seen})
        ::continue::
    end
    return results
end

-- ============================================================================
-- Views
-- ============================================================================

local function renderList(api, title, entries)
    if #entries == 0 then
        api:ShowDialog(title, "No entries match your filters.")
        return
    end
    local pages = paginate(entries, config.pageSize)
    for i, page in ipairs(pages) do
        local lines = {}
        for _, entry in ipairs(page) do
            table.insert(lines, formatEnemy(entry.enemy, entry.seen))
            table.insert(lines, "")
        end
        api:ShowDialog(string.format("%s (Page %d/%d)", title, i, #pages), table.concat(lines, "\n"))
    end
end

local function viewSummary(api)
    local total = #ENEMIES
    local seenCount, unseenCount, unknownCount = 0, 0, 0
    for _, e in ipairs(ENEMIES) do
        local s = isEncountered(api, e.id)
        if s == true then seenCount = seenCount + 1
        elseif s == false then unseenCount = unseenCount + 1
        else unknownCount = unknownCount + 1 end
    end
    local lines = {
        string.format("Enemies tracked: %d", total),
        string.format("Seen: %d", seenCount),
        string.format("Unseen: %d", unseenCount),
        string.format("Unknown: %d", unknownCount),
        "Use filters to find specific weaknesses, bosses, or unseen enemies.",
    }
    api:ShowDialog("Bestiary Summary", table.concat(lines, "\n"))
end

local function viewAll(api)
    local entries = filterEnemies(api, {})
    renderList(api, "All Enemies", entries)
end

local function viewUnseen(api)
    local entries = filterEnemies(api, {unseenOnly = true})
    renderList(api, "Unseen Enemies", entries)
end

local function viewByType(api)
    local typeInput = api:ShowInput("Filter by Type", "Enter type (e.g., Humanoid, Beast, Dragon, Boss):", "Humanoid")
    if not typeInput then return end
    local entries = filterEnemies(api, {type = typeInput})
    renderList(api, "Type: " .. typeInput, entries)
end

local function viewByElement(api)
    local el = api:ShowInput("Filter by Element", "Enter element to match (weak/resist/absorb/null):", "Ice")
    if not el then return end
    local entries = filterEnemies(api, {element = el})
    renderList(api, "Element match: " .. el, entries)
end

local function viewSearch(api)
    local term = api:ShowInput("Search", "Keyword across name/type/location/status:", "")
    if not term or term == "" then return end
    local entries = filterEnemies(api, {search = term})
    renderList(api, "Search: " .. term, entries)
end

local function viewBosses(api)
    local entries = filterEnemies(api, {bossesOnly = true})
    renderList(api, "Bosses", entries)
end

local function viewExport(api)
    local entries = filterEnemies(api, {})
    local lines = {"Bestiary Export"}
    for _, entry in ipairs(entries) do
        table.insert(lines, formatEnemy(entry.enemy, entry.seen))
        table.insert(lines, "")
    end
    api:ShowDialog("Export (copyable)", table.concat(lines, "\n"))
end

local function viewHelp(api)
    local lines = {
        "Bestiary (read-only)",
        "- Summary: Counts of seen/unseen/unknown.",
        "- All: Full list with status icons.",
        "- Unseen: Focus on missing encounters (Unknown optional).",
        "- Type: Match on enemy type (Humanoid, Beast, Dragon, Boss, etc.).",
        "- Element: Match if enemy is weak/resist/absorb/null to the element.",
        "- Bosses: Boss-only view.",
        "- Search: Keyword across name/type/location/status.",
        "- Export: Copyable text of all entries.",
        "Status icons: [SEEN], [UNSEEN], [?]=unknown (API missing).",
        "Config: see top of plugin.lua for page size and unknown handling.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu Loop
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Bestiary",
            "Select an option:\n1) Summary\n2) All\n3) Unseen\n4) By Type\n5) By Element\n6) Bosses\n7) Search\n8) Export\n9) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            viewSummary(api)
        elseif sel == 2 then
            viewAll(api)
        elseif sel == 3 then
            viewUnseen(api)
        elseif sel == 4 then
            viewByType(api)
        elseif sel == 5 then
            viewByElement(api)
        elseif sel == 6 then
            viewBosses(api)
        elseif sel == 7 then
            viewSearch(api)
        elseif sel == 8 then
            viewExport(api)
        elseif sel == 9 then
            viewHelp(api)
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
        api:ShowDialog("Error", "Bestiary failed: " .. tostring(err))
    end
end

return {
    id = "bestiary",
    name = "Bestiary",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
