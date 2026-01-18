-- ============================================================================
-- Strago's Lore Tracker Plugin
-- ============================================================================
-- @id: lore-tracker
-- @name: Strago's Lore Tracker
-- @version: 1.0.0
-- @author: FF6 Editor Team
-- @description: Track Strago's 24 Blue Magic spells (Lores) with enemy sources and locations
-- @permissions: read_save, ui_display

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    showLearnedOnly = false,        -- Show only learned Lores
    showUnlearnedOnly = false,      -- Show only unlearned Lores
    showMissableOnly = false,       -- Show only missable Lores
    highlightMissable = true,       -- Highlight missable Lores
    sortByName = true,              -- Sort by name (false = by Lore ID)
    showEnemyLocations = true,      -- Show where to find enemies
}

-- ============================================================================
-- Lore Database
-- ============================================================================
-- Complete database of all 24 Lores with enemy sources and locations
local LORE_DATABASE = {
    -- Lore ID, Name, Enemy Sources, Location, Missable (WoB only)
    {id = 1, name = "Condemn", enemies = {"Critic", "Death Machine"}, 
     location = "Kefka's Tower", missable = false, boss = false},
    
    {id = 2, name = "Roulette", enemies = {"Doom Drgn", "Didalos"}, 
     location = "Phoenix Cave, Kefka's Tower", missable = false, boss = false},
    
    {id = 3, name = "CleanSweep", enemies = {"Abolisher", "Siegfried (boss)"}, 
     location = "Colosseum, Floating Continent", missable = false, boss = true},
    
    {id = 4, name = "Aqua Rake", enemies = {"Actaneon", "Rizopas"}, 
     location = "Fanatics' Tower, Thamasa Area", missable = false, boss = false},
    
    {id = 5, name = "Aero", enemies = {"Sprinter", "Innoc"}, 
     location = "Mt. Zozo, Veldt", missable = false, boss = false},
    
    {id = 6, name = "Blow Fish", enemies = {"Anguiform", "Nautiloid"}, 
     location = "Serpent Trench, Ocean (WoR)", missable = false, boss = false},
    
    {id = 7, name = "Big Guard", enemies = {"Dark Force", "Galypdes (boss)"}, 
     location = "Kefka's Tower, WoB Sealed Gate", missable = false, boss = true},
    
    {id = 8, name = "Revenge", enemies = {"Misfit", "Pm Stalker"}, 
     location = "Phoenix Cave, Mt. Zozo", missable = false, boss = false},
    
    {id = 9, name = "Pearl Wind", enemies = {"Flan"}, 
     location = "Cave to South Figaro, Figaro Castle", missable = false, boss = false},
    
    {id = 10, name = "L.5 Doom", enemies = {"Trapper"}, 
     location = "Fanatics' Tower", missable = false, boss = false},
    
    {id = 11, name = "L.4 Flare", enemies = {"Apokryphos", "Trapper"}, 
     location = "Fanatics' Tower, Kefka's Tower", missable = false, boss = false},
    
    {id = 12, name = "L.3 Muddle", enemies = {"Apokryphos", "Abolisher"}, 
     location = "Fanatics' Tower, Floating Continent", missable = false, boss = false},
    
    {id = 13, name = "Reflect???", enemies = {"Reach Frog"}, 
     location = "Triangle Island (WoR)", missable = false, boss = false},
    
    {id = 14, name = "L? Pearl", enemies = {"Magic Urn (boss)"}, 
     location = "Owzer's Mansion (boss fight)", missable = true, boss = true},
    
    {id = 15, name = "Step Mine", enemies = {"Pm Stalker", "Sky Base"}, 
     location = "Mt. Zozo, Floating Continent", missable = false, boss = false},
    
    {id = 16, name = "Force Field", enemies = {"Anemone", "Aquila"}, 
     location = "Fanatics' Tower, Narshe Caves", missable = false, boss = false},
    
    {id = 17, name = "Dischord", enemies = {"Cephaler", "Trixter"}, 
     location = "Ancient Castle, Colosseum", missable = false, boss = false},
    
    {id = 18, name = "Sour Mouth", enemies = {"Grenade", "Bug"}, 
     location = "Magitek Research Facility, Cave in the Veldt", missable = false, boss = false},
    
    {id = 19, name = "Pep Up", enemies = {"Muus"}, 
     location = "Triangle Island (WoR)", missable = false, boss = false},
    
    {id = 20, name = "Rippler", enemies = {"Skull Drgn (boss)", "Galypdes (boss)"}, 
     location = "Kefka's Tower (boss), WoB Sealed Gate (boss)", missable = false, boss = true},
    
    {id = 21, name = "Stone", enemies = {"Peepers", "Mag Roader (brown)"}, 
     location = "Serpent Trench, Magitek Research Facility", missable = false, boss = false},
    
    {id = 22, name = "Quasar", enemies = {"Dark Force", "Guardian (boss)"}, 
     location = "Kefka's Tower, Narshe Cliffs (boss)", missable = true, boss = true},
    
    {id = 23, name = "Grand Train", enemies = {"Dullahan (boss)", "Hidon (boss)"}, 
     location = "Colosseum (boss), Ebot's Rock (boss)", missable = false, boss = true},
    
    {id = 24, name = "Exploder", enemies = {"Bomb", "Grenade"}, 
     location = "Magitek Research Facility, Kefka's Tower", missable = false, boss = false},
}

-- Lore IDs mapped to spell IDs in the game (these might need adjustment based on actual spell data)
-- Note: This is a simplified mapping. In a real implementation, you'd need to map
-- Lore spell IDs to the actual spell system IDs used by the game
local LORE_SPELL_IDS = {
    -- Mapping Lore ID to potential spell ID in character's spell list
    -- This may need to be adjusted based on actual game data structure
    [1] = 38,   -- Condemn
    [2] = 39,   -- Roulette
    [3] = 40,   -- CleanSweep
    [4] = 41,   -- Aqua Rake
    [5] = 42,   -- Aero
    [6] = 43,   -- Blow Fish
    [7] = 44,   -- Big Guard
    [8] = 45,   -- Revenge
    [9] = 46,   -- Pearl Wind
    [10] = 47,  -- L.5 Doom
    [11] = 48,  -- L.4 Flare
    [12] = 49,  -- L.3 Muddle
    [13] = 50,  -- Reflect???
    [14] = 51,  -- L? Pearl
    [15] = 52,  -- Step Mine
    [16] = 53,  -- Force Field
    [17] = 54,  -- Dischord
    [18] = 37,  -- Sour Mouth
    [19] = 36,  -- Pep Up
    [20] = 35,  -- Rippler
    [21] = 34,  -- Stone
    [22] = 33,  -- Quasar
    [23] = 32,  -- Grand Train
    [24] = 31,  -- Exploder
}

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Find Strago's character index
local function findStrago(api)
    -- Try to find Strago by name
    for i = 0, 15 do
        local char = api:GetCharacter(i)
        if char and char.Name then
            -- Strago's name variations
            if char.Name == "Strago" or char.Name == "STRAGO" or 
               char.Name:lower():find("strago") then
                return i, char
            end
        end
    end
    return nil, nil
end

-- Check if a Lore is learned
local function isLoreLearned(stragoChar, loreID)
    if not stragoChar or not stragoChar.Spells then
        return false
    end
    
    -- Get the spell ID for this Lore
    local spellID = LORE_SPELL_IDS[loreID]
    if not spellID then
        return false
    end
    
    -- Check if Strago has this spell with proficiency > 0
    if stragoChar.Spells[spellID] then
        return stragoChar.Spells[spellID] > 0
    end
    
    return false
end

-- Count learned Lores
local function countLearnedLores(stragoChar)
    local count = 0
    for loreID = 1, 24 do
        if isLoreLearned(stragoChar, loreID) then
            count = count + 1
        end
    end
    return count
end

-- Format enemy list
local function formatEnemies(enemies, isBoss)
    local enemyStr = table.concat(enemies, ", ")
    if isBoss then
        enemyStr = enemyStr .. " [BOSS]"
    end
    return enemyStr
end

-- Format Lore entry
local function formatLoreEntry(lore, learned, config)
    local status = learned and "[✓]" or "[ ]"
    local missableTag = (config.highlightMissable and lore.missable) and " ⚠️ MISSABLE" or ""
    
    local line = string.format("%s %2d. %-15s", status, lore.id, lore.name)
    
    if config.showEnemyLocations then
        line = line .. "\n       Enemies: " .. formatEnemies(lore.enemies, lore.boss)
        line = line .. "\n       Location: " .. lore.location .. missableTag
    end
    
    return line
end

-- Filter Lores based on configuration
local function filterLores(stragoChar, config)
    local filtered = {}
    
    for _, lore in ipairs(LORE_DATABASE) do
        local learned = isLoreLearned(stragoChar, lore.id)
        
        -- Apply filters
        local include = true
        
        if config.showLearnedOnly and not learned then
            include = false
        end
        
        if config.showUnlearnedOnly and learned then
            include = false
        end
        
        if config.showMissableOnly and not lore.missable then
            include = false
        end
        
        if include then
            table.insert(filtered, {lore = lore, learned = learned})
        end
    end
    
    return filtered
end

-- Sort Lores
local function sortLores(lores, config)
    if config.sortByName then
        table.sort(lores, function(a, b)
            return a.lore.name < b.lore.name
        end)
    else
        table.sort(lores, function(a, b)
            return a.lore.id < b.lore.id
        end)
    end
end

-- Build Lore tracker display
local function buildLoreDisplay(stragoChar, config)
    local sections = {}
    
    -- Header
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "STRAGO'S LORE TRACKER - Blue Magic Spells")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    -- Statistics
    local learnedCount = countLearnedLores(stragoChar)
    local totalCount = 24
    local percentage = (learnedCount / totalCount) * 100
    
    table.insert(sections, "")
    table.insert(sections, string.format("Progress: %d / %d Lores learned (%.1f%%)", 
        learnedCount, totalCount, percentage))
    table.insert(sections, "")
    
    -- Missable warning
    local missableCount = 0
    for _, lore in ipairs(LORE_DATABASE) do
        if lore.missable and not isLoreLearned(stragoChar, lore.id) then
            missableCount = missableCount + 1
        end
    end
    
    if missableCount > 0 then
        table.insert(sections, string.format("⚠️  WARNING: %d missable Lore(s) not yet learned!", missableCount))
        table.insert(sections, "")
    end
    
    -- Get and sort Lores
    local filteredLores = filterLores(stragoChar, config)
    sortLores(filteredLores, config)
    
    if #filteredLores == 0 then
        table.insert(sections, "No Lores match the current filter settings.")
    else
        -- Display Lores
        table.insert(sections, "Lore List:")
        table.insert(sections, string.rep("-", 60))
        
        for _, entry in ipairs(filteredLores) do
            table.insert(sections, formatLoreEntry(entry.lore, entry.learned, config))
            table.insert(sections, "")
        end
    end
    
    -- Legend
    table.insert(sections, string.rep("-", 60))
    table.insert(sections, "Legend:")
    table.insert(sections, "  [✓] = Learned    [ ] = Not Learned")
    table.insert(sections, "  [BOSS] = Only available from boss enemies")
    table.insert(sections, "  ⚠️ MISSABLE = Can only be learned in World of Balance")
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    return table.concat(sections, "\n")
end

-- Build main menu
local function buildMainMenu(stragoChar)
    local learnedCount = countLearnedLores(stragoChar)
    local menu = "=" .. string.rep("=", 59) .. "\n"
    menu = menu .. "STRAGO'S LORE TRACKER - Main Menu\n"
    menu = menu .. "=" .. string.rep("=", 59) .. "\n\n"
    menu = menu .. string.format("Current Progress: %d / 24 Lores learned\n\n", learnedCount)
    menu = menu .. "Select an option:\n\n"
    menu = menu .. "1. View All Lores (Complete List)\n"
    menu = menu .. "2. View Learned Lores Only\n"
    menu = menu .. "3. View Unlearned Lores Only\n"
    menu = menu .. "4. View Missable Lores Only\n"
    menu = menu .. "5. View Lore Statistics\n"
    menu = menu .. "0. Exit\n\n"
    menu = menu .. "Enter your choice (0-5):"
    
    return menu
end

-- Build statistics display
local function buildStatistics(stragoChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "LORE STATISTICS")
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "")
    
    -- Overall stats
    local learnedCount = countLearnedLores(stragoChar)
    local percentage = (learnedCount / 24) * 100
    
    table.insert(sections, string.format("Total Lores Learned: %d / 24 (%.1f%%)", learnedCount, percentage))
    table.insert(sections, "")
    
    -- Missable Lores
    local missableTotal = 0
    local missableLearned = 0
    for _, lore in ipairs(LORE_DATABASE) do
        if lore.missable then
            missableTotal = missableTotal + 1
            if isLoreLearned(stragoChar, lore.id) then
                missableLearned = missableLearned + 1
            end
        end
    end
    
    table.insert(sections, string.format("Missable Lores: %d / %d learned", missableLearned, missableTotal))
    if missableLearned < missableTotal then
        table.insert(sections, "  ⚠️  Some missable Lores not yet obtained!")
    end
    table.insert(sections, "")
    
    -- Boss-only Lores
    local bossTotal = 0
    local bossLearned = 0
    for _, lore in ipairs(LORE_DATABASE) do
        if lore.boss then
            bossTotal = bossTotal + 1
            if isLoreLearned(stragoChar, lore.id) then
                bossLearned = bossLearned + 1
            end
        end
    end
    
    table.insert(sections, string.format("Boss-Only Lores: %d / %d learned", bossLearned, bossTotal))
    table.insert(sections, "")
    
    -- Completion status
    if learnedCount == 24 then
        table.insert(sections, "🎉 CONGRATULATIONS! All Lores learned! 🎉")
    else
        table.insert(sections, string.format("Lores Remaining: %d", 24 - learnedCount))
    end
    
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    return table.concat(sections, "\n")
end

-- ============================================================================
-- Main Entry Point
-- ============================================================================
function main(api)
    -- Verify permissions
    if not api:HasPermission("read_save") then
        api:Log("error", "Missing required permission: read_save")
        api:ShowDialog("Permission Error", "This plugin requires 'read_save' permission.")
        return false
    end
    
    if not api:HasPermission("ui_display") then
        api:Log("error", "Missing required permission: ui_display")
        return false
    end
    
    api:Log("info", "Lore Tracker plugin started")
    
    -- Find Strago
    local stragoIndex, stragoChar = findStrago(api)
    if not stragoIndex then
        api:ShowDialog("Error", "Strago not found in save file. Make sure you have recruited Strago.")
        api:Log("error", "Strago character not found")
        return false
    end
    
    api:Log("info", string.format("Found Strago at character index: %d", stragoIndex))
    
    -- Main menu loop
    while true do
        local menu = buildMainMenu(stragoChar)
        local choice = api:ShowInput("Lore Tracker - Main Menu", menu, "1")
        
        if not choice or choice == "0" then
            api:Log("info", "User exited Lore Tracker")
            return true
        end
        
        if choice == "1" then
            -- View all Lores
            config.showLearnedOnly = false
            config.showUnlearnedOnly = false
            config.showMissableOnly = false
            local display = buildLoreDisplay(stragoChar, config)
            api:ShowDialog("All Lores", display)
            
        elseif choice == "2" then
            -- View learned only
            config.showLearnedOnly = true
            config.showUnlearnedOnly = false
            config.showMissableOnly = false
            local display = buildLoreDisplay(stragoChar, config)
            api:ShowDialog("Learned Lores", display)
            
        elseif choice == "3" then
            -- View unlearned only
            config.showLearnedOnly = false
            config.showUnlearnedOnly = true
            config.showMissableOnly = false
            local display = buildLoreDisplay(stragoChar, config)
            api:ShowDialog("Unlearned Lores", display)
            
        elseif choice == "4" then
            -- View missable only
            config.showLearnedOnly = false
            config.showUnlearnedOnly = false
            config.showMissableOnly = true
            local display = buildLoreDisplay(stragoChar, config)
            api:ShowDialog("Missable Lores", display)
            
        elseif choice == "5" then
            -- View statistics
            local stats = buildStatistics(stragoChar)
            api:ShowDialog("Lore Statistics", stats)
            
        else
            api:ShowDialog("Invalid Choice", "Please enter a number between 0 and 5.")
        end
    end
    
    return true
end

-- ============================================================================
-- Error Handling Wrapper
-- ============================================================================
local function safeMain(api)
    local success, result = pcall(main, api)
    if not success then
        local errorMsg = tostring(result)
        api:Log("error", "Plugin error: " .. errorMsg)
        api:ShowDialog("Plugin Error", 
            "An error occurred while running the plugin:\n\n" .. errorMsg .. 
            "\n\nCheck the plugin log for details.")
        return false
    end
    return result
end

-- Execute with error handling
return safeMain(api)
