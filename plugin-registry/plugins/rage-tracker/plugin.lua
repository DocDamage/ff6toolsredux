-- ============================================================================
-- Gau's Rage Tracker Plugin
-- ============================================================================
-- @id: rage-tracker
-- @name: Gau's Rage Tracker
-- @version: 1.0.0
-- @author: FF6 Editor Team
-- @description: Track all 384 Rages with filtering, search, and location guides
-- @permissions: read_save, ui_display

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    showLearnedOnly = false,
    showUnlearnedOnly = false,
    itemsPerPage = 20,              -- Pagination for large lists
    showAbilities = true,            -- Show Rage abilities
    sortByName = true,               -- Sort by name (false = by ID)
}

-- ============================================================================
-- Enemy Type & Location Constants
-- ============================================================================
local ENEMY_TYPES = {
    NORMAL = "Normal",
    BOSS = "Boss",
    MAGIC = "Magic",
    UNDEAD = "Undead",
    DRAGON = "Dragon",
    MECHANICAL = "Mechanical",
    HUMANOID = "Humanoid",
    FLYING = "Flying",
    AQUATIC = "Aquatic"
}

local LOCATIONS = {
    "World Map (General)",
    "Narshe Area",
    "Figaro Area",
    "South Figaro Area",
    "Mt. Kolts",
    "Returner's Hideout",
    "Phantom Forest",
    "Barren Falls",
    "Veldt (Common)",
    "Veldt (Rare)",
    "Serpent Trench",
    "Zozo",
    "Opera House",
    "Southern Continent",
    "Crescent Mountain",
    "Doma Castle",
    "Mt. Zozo",
    "Floating Continent",
    "World of Ruin (General)",
    "Kefka's Tower",
    "Phoenix Cave",
    "Fanatics' Tower",
    "Ancient Castle",
    "Other"
}

-- ============================================================================
-- Rage Database (Sample - Full database would be 384 entries)
-- ============================================================================
-- Note: This is a representative sample. A production version would include
-- all 384 enemies. For this implementation, we'll include a variety of
-- enemies to demonstrate the system.

local RAGE_DATABASE = {
    -- Early Game Enemies
    {id = 1, name = "Guard", type = ENEMY_TYPES.HUMANOID, location = "Narshe Area", 
     ability = "Fight", power = "Low", veldt = true},
    {id = 2, name = "Lobo", type = ENEMY_TYPES.NORMAL, location = "Figaro Area", 
     ability = "Scratch", power = "Low", veldt = true},
    {id = 3, name = "Marshal", type = ENEMY_TYPES.HUMANOID, location = "South Figaro Area", 
     ability = "Fire", power = "Low", veldt = true},
    {id = 4, name = "Hornet", type = ENEMY_TYPES.FLYING, location = "Figaro Area", 
     ability = "Sting", power = "Low", veldt = true},
    {id = 5, name = "Leafer", type = ENEMY_TYPES.NORMAL, location = "Phantom Forest", 
     ability = "Acid Rain", power = "Low", veldt = true},
    
    -- Mid Game Enemies
    {id = 50, name = "Aspik", type = ENEMY_TYPES.NORMAL, location = "Serpent Trench", 
     ability = "Lode Stone", power = "Medium", veldt = true},
    {id = 51, name = "Actaneon", type = ENEMY_TYPES.AQUATIC, location = "Serpent Trench", 
     ability = "Aqua Rake", power = "Medium", veldt = true},
    {id = 52, name = "Anguiform", type = ENEMY_TYPES.AQUATIC, location = "Serpent Trench", 
     ability = "Blow Fish", power = "Medium", veldt = true},
    
    -- Zozo Enemies
    {id = 100, name = "Vaporite", type = ENEMY_TYPES.FLYING, location = "Zozo", 
     ability = "Acid Rain", power = "Medium", veldt = true},
    {id = 101, name = "Brawler", type = ENEMY_TYPES.HUMANOID, location = "Zozo", 
     ability = "Pummel", power = "Medium", veldt = true},
    
    -- Mt. Zozo
    {id = 150, name = "Harpy", type = ENEMY_TYPES.FLYING, location = "Mt. Zozo", 
     ability = "Aero", power = "Medium", veldt = true},
    {id = 151, name = "Pm Stalker", type = ENEMY_TYPES.MAGIC, location = "Mt. Zozo", 
     ability = "Revenge", power = "Medium", veldt = true},
    
    -- Floating Continent
    {id = 200, name = "Apocrypha", type = ENEMY_TYPES.MAGIC, location = "Floating Continent", 
     ability = "L.3 Muddle", power = "High", veldt = true},
    {id = 201, name = "Behemoth", type = ENEMY_TYPES.NORMAL, location = "Floating Continent", 
     ability = "Meteor", power = "High", veldt = true},
    {id = 202, name = "Dragon", type = ENEMY_TYPES.DRAGON, location = "Floating Continent", 
     ability = "Flare Star", power = "High", veldt = true},
    
    -- World of Ruin
    {id = 250, name = "Chaos Dragon", type = ENEMY_TYPES.DRAGON, location = "Kefka's Tower", 
     ability = "Fallen One", power = "Very High", veldt = true},
    {id = 251, name = "Wizard", type = ENEMY_TYPES.MAGIC, location = "Fanatics' Tower", 
     ability = "Flare", power = "Very High", veldt = true},
    {id = 252, name = "Dark Force", type = ENEMY_TYPES.MAGIC, location = "Kefka's Tower", 
     ability = "Quasar", power = "Very High", veldt = true},
    
    -- Special/Rare Enemies
    {id = 300, name = "Intangir", type = ENEMY_TYPES.MAGIC, location = "Triangle Island", 
     ability = "Meteor", power = "Very High", veldt = true},
    {id = 301, name = "Cactrot", type = ENEMY_TYPES.NORMAL, location = "World Map (Desert)", 
     ability = "1000 Needles", power = "Fixed", veldt = true},
    {id = 302, name = "Muus", type = ENEMY_TYPES.NORMAL, location = "Triangle Island", 
     ability = "Pep Up", power = "Support", veldt = true},
    
    -- Mechanical Enemies
    {id = 350, name = "Mag Roader", type = ENEMY_TYPES.MECHANICAL, location = "Magitek Facility", 
     ability = "Wheel", power = "Medium", veldt = true},
    {id = 351, name = "Spit Fire", type = ENEMY_TYPES.MECHANICAL, location = "Magitek Facility", 
     ability = "Fireball", power = "Medium", veldt = true},
    
    -- Undead Enemies
    {id = 370, name = "Ghost", type = ENEMY_TYPES.UNDEAD, location = "Phantom Forest", 
     ability = "Possess", power = "Special", veldt = true},
    {id = 371, name = "Zombie", type = ENEMY_TYPES.UNDEAD, location = "Doma Castle", 
     ability = "Zombie Touch", power = "Low", veldt = true},
}

-- Note: In production, this would contain all 384 enemies with complete data

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Find Gau's character index
local function findGau(api)
    for i = 0, 15 do
        local char = api:GetCharacter(i)
        if char and char.Name then
            if char.Name == "Gau" or char.Name == "GAU" or 
               char.Name:lower():find("gau") then
                return i, char
            end
        end
    end
    return nil, nil
end

-- Check if a Rage is learned (placeholder - needs actual implementation)
local function isRageLearned(gauChar, rageID)
    -- In actual implementation, this would check a bitfield or array
    -- that tracks which of the 384 Rages are learned
    -- For now, return false for all
    return false
end

-- Count learned Rages
local function countLearnedRages(gauChar)
    local count = 0
    for _, enemy in ipairs(RAGE_DATABASE) do
        if isRageLearned(gauChar, enemy.id) then
            count = count + 1
        end
    end
    return count
end

-- Get total Rages in database
local function getTotalRages()
    return #RAGE_DATABASE
end

-- Filter Rages by various criteria
local function filterRages(gauChar, filterType, filterValue)
    local filtered = {}
    
    for _, enemy in ipairs(RAGE_DATABASE) do
        local learned = isRageLearned(gauChar, enemy.id)
        local include = true
        
        -- Apply status filters
        if config.showLearnedOnly and not learned then
            include = false
        end
        
        if config.showUnlearnedOnly and learned then
            include = false
        end
        
        -- Apply type/location filters
        if filterType == "type" and filterValue then
            if enemy.type ~= filterValue then
                include = false
            end
        end
        
        if filterType == "location" and filterValue then
            if not enemy.location:find(filterValue) then
                include = false
            end
        end
        
        if filterType == "search" and filterValue then
            local searchLower = filterValue:lower()
            if not (enemy.name:lower():find(searchLower) or 
                    enemy.ability:lower():find(searchLower)) then
                include = false
            end
        end
        
        if include then
            table.insert(filtered, {enemy = enemy, learned = learned})
        end
    end
    
    return filtered
end

-- Sort Rages
local function sortRages(rages, config)
    if config.sortByName then
        table.sort(rages, function(a, b)
            return a.enemy.name < b.enemy.name
        end)
    else
        table.sort(rages, function(a, b)
            return a.enemy.id < b.enemy.id
        end)
    end
end

-- Search Rages
local function searchRages(gauChar, searchTerm)
    return filterRages(gauChar, "search", searchTerm)
end

-- Paginate results
local function paginateResults(items, page, itemsPerPage)
    local totalPages = math.ceil(#items / itemsPerPage)
    local startIdx = (page - 1) * itemsPerPage + 1
    local endIdx = math.min(page * itemsPerPage, #items)
    
    local pageItems = {}
    for i = startIdx, endIdx do
        table.insert(pageItems, items[i])
    end
    
    return pageItems, totalPages, startIdx
end

-- Format Rage entry
local function formatRageEntry(rage, learned, index)
    local status = learned and "[✓]" or "[ ]"
    local line = string.format("%s %3d. %-20s", status, index or rage.id, rage.name)
    
    if config.showAbilities then
        line = line .. string.format("\n       Type: %-12s  Ability: %s (Power: %s)", 
            rage.type, rage.ability, rage.power)
        line = line .. string.format("\n       Location: %s", rage.location)
    end
    
    return line
end

-- ============================================================================
-- Display Functions
-- ============================================================================

-- Build main display
local function buildRageDisplay(gauChar, rages, page, title)
    local sections = {}
    
    -- Header
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "GAU'S RAGE TRACKER - " .. title)
    table.insert(sections, "=" .. string.rep("=", 70))
    
    -- Statistics
    local learnedCount = countLearnedRages(gauChar)
    local totalCount = getTotalRages()
    local percentage = (learnedCount / totalCount) * 100
    
    table.insert(sections, "")
    table.insert(sections, string.format("Overall Progress: %d / %d Rages learned (%.1f%%)", 
        learnedCount, totalCount, percentage))
    
    -- Current filter stats
    if #rages < totalCount then
        table.insert(sections, string.format("Filtered Results: %d Rages", #rages))
    end
    table.insert(sections, "")
    
    -- Paginate
    local pageItems, totalPages, startIdx = paginateResults(rages, page, config.itemsPerPage)
    
    if #pageItems == 0 then
        table.insert(sections, "No Rages match the current filter.")
    else
        table.insert(sections, string.format("Showing %d-%d of %d (Page %d/%d)", 
            startIdx, startIdx + #pageItems - 1, #rages, page, totalPages))
        table.insert(sections, string.rep("-", 70))
        
        for i, entry in ipairs(pageItems) do
            table.insert(sections, formatRageEntry(entry.enemy, entry.learned, startIdx + i - 1))
            table.insert(sections, "")
        end
        
        if totalPages > 1 then
            table.insert(sections, string.rep("-", 70))
            table.insert(sections, string.format("Page %d of %d - Enter page number to view", page, totalPages))
        end
    end
    
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n"), totalPages
end

-- Build statistics
local function buildStatistics(gauChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "RAGE STATISTICS")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    
    local learnedCount = countLearnedRages(gauChar)
    local totalCount = getTotalRages()
    local percentage = (learnedCount / totalCount) * 100
    
    table.insert(sections, string.format("Total Rages Learned: %d / %d (%.1f%%)", 
        learnedCount, totalCount, percentage))
    table.insert(sections, "")
    
    -- By enemy type
    table.insert(sections, "Rages by Enemy Type:")
    local typeCounts = {}
    for _, enemy in ipairs(RAGE_DATABASE) do
        typeCounts[enemy.type] = (typeCounts[enemy.type] or 0) + 1
    end
    
    for enemyType, count in pairs(typeCounts) do
        table.insert(sections, string.format("  %-15s: %3d Rages", enemyType, count))
    end
    
    table.insert(sections, "")
    
    -- Completion
    if learnedCount == totalCount then
        table.insert(sections, "🎉 MASTER OF RAGES! All 384 Rages learned! 🎉")
    else
        table.insert(sections, string.format("Rages Remaining: %d", totalCount - learnedCount))
    end
    
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n")
end

-- Build filter menu
local function buildFilterMenu()
    local menu = "=" .. string.rep("=", 70) .. "\n"
    menu = menu .. "RAGE TRACKER - Filter Menu\n"
    menu = menu .. "=" .. string.rep("=", 70) .. "\n\n"
    menu = menu .. "Filter by Enemy Type:\n\n"
    menu = menu .. "1. Normal Enemies\n"
    menu = menu .. "2. Magic Enemies\n"
    menu = menu .. "3. Mechanical Enemies\n"
    menu = menu .. "4. Undead Enemies\n"
    menu = menu .. "5. Dragon Enemies\n"
    menu = menu .. "6. Flying Enemies\n"
    menu = menu .. "7. Aquatic Enemies\n"
    menu = menu .. "8. Humanoid Enemies\n"
    menu = menu .. "0. Back to Main Menu\n\n"
    menu = menu .. "Enter your choice (0-8):"
    
    return menu
end

-- Build main menu
local function buildMainMenu(gauChar)
    local learnedCount = countLearnedRages(gauChar)
    local totalCount = getTotalRages()
    
    local menu = "=" .. string.rep("=", 70) .. "\n"
    menu = menu .. "GAU'S RAGE TRACKER - Main Menu\n"
    menu = menu .. "=" .. string.rep("=", 70) .. "\n\n"
    menu = menu .. string.format("Current Progress: %d / %d Rages learned\n\n", 
        learnedCount, totalCount)
    menu = menu .. "Select an option:\n\n"
    menu = menu .. "1. View All Rages\n"
    menu = menu .. "2. View Learned Rages Only\n"
    menu = menu .. "3. View Unlearned Rages Only\n"
    menu = menu .. "4. Filter by Enemy Type\n"
    menu = menu .. "5. Filter by Location\n"
    menu = menu .. "6. Search Rages\n"
    menu = menu .. "7. View Statistics\n"
    menu = menu .. "8. Veldt Formation Guide\n"
    menu = menu .. "0. Exit\n\n"
    menu = menu .. "Enter your choice (0-8):"
    
    return menu
end

-- Build location filter menu
local function buildLocationMenu()
    local menu = "=" .. string.rep("=", 70) .. "\n"
    menu = menu .. "RAGE TRACKER - Location Filter\n"
    menu = menu .. "=" .. string.rep("=", 70) .. "\n\n"
    menu = menu .. "Select location (showing first 10 for demo):\n\n"
    
    for i = 1, math.min(10, #LOCATIONS) do
        menu = menu .. string.format("%d. %s\n", i, LOCATIONS[i])
    end
    
    menu = menu .. "0. Back to Main Menu\n\n"
    menu = menu .. "Enter your choice:"
    
    return menu
end

-- Build Veldt formation guide
local function buildVeldtGuide()
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "VELDT FORMATION GUIDE")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    table.insert(sections, "The Veldt has multiple formations where different enemies appear.")
    table.insert(sections, "Enemies are added to the Veldt after you encounter them elsewhere.")
    table.insert(sections, "")
    table.insert(sections, "Key Facts:")
    table.insert(sections, "  • 384 total Rages available")
    table.insert(sections, "  • Must encounter enemy at least once before it appears on Veldt")
    table.insert(sections, "  • Some enemies never appear on Veldt (boss-exclusive)")
    table.insert(sections, "  • Formations have 1-6 enemies")
    table.insert(sections, "  • Some enemies are rare (low encounter rate)")
    table.insert(sections, "")
    table.insert(sections, "Strategy Tips:")
    table.insert(sections, "  • Use the Veldt as your training ground")
    table.insert(sections, "  • Gau learns Rage after successfully Leaping an enemy")
    table.insert(sections, "  • Check your unlearned list to see what you're missing")
    table.insert(sections, "  • Some Rages are WoB only - get them early!")
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 70))
    
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
    
    api:Log("info", "Rage Tracker plugin started")
    
    -- Find Gau
    local gauIndex, gauChar = findGau(api)
    if not gauIndex then
        api:ShowDialog("Error", "Gau not found in save file. Make sure you have recruited Gau.")
        api:Log("error", "Gau character not found")
        return false
    end
    
    api:Log("info", string.format("Found Gau at character index: %d", gauIndex))
    
    -- Main menu loop
    while true do
        local menu = buildMainMenu(gauChar)
        local choice = api:ShowInput("Rage Tracker - Main Menu", menu, "1")
        
        if not choice or choice == "0" then
            api:Log("info", "User exited Rage Tracker")
            return true
        end
        
        if choice == "1" then
            -- View all Rages
            config.showLearnedOnly = false
            config.showUnlearnedOnly = false
            local rages = filterRages(gauChar, nil, nil)
            sortRages(rages, config)
            local display = buildRageDisplay(gauChar, rages, 1, "All Rages")
            api:ShowDialog("All Rages", display)
            
        elseif choice == "2" then
            -- View learned only
            config.showLearnedOnly = true
            config.showUnlearnedOnly = false
            local rages = filterRages(gauChar, nil, nil)
            sortRages(rages, config)
            local display = buildRageDisplay(gauChar, rages, 1, "Learned Rages")
            api:ShowDialog("Learned Rages", display)
            config.showLearnedOnly = false
            
        elseif choice == "3" then
            -- View unlearned only
            config.showLearnedOnly = false
            config.showUnlearnedOnly = true
            local rages = filterRages(gauChar, nil, nil)
            sortRages(rages, config)
            local display = buildRageDisplay(gauChar, rages, 1, "Unlearned Rages")
            api:ShowDialog("Unlearned Rages", display)
            config.showUnlearnedOnly = false
            
        elseif choice == "4" then
            -- Filter by type
            local filterMenu = buildFilterMenu()
            local typeChoice = api:ShowInput("Filter by Type", filterMenu, "1")
            
            if typeChoice and typeChoice ~= "0" then
                local typeMap = {
                    ["1"] = ENEMY_TYPES.NORMAL,
                    ["2"] = ENEMY_TYPES.MAGIC,
                    ["3"] = ENEMY_TYPES.MECHANICAL,
                    ["4"] = ENEMY_TYPES.UNDEAD,
                    ["5"] = ENEMY_TYPES.DRAGON,
                    ["6"] = ENEMY_TYPES.FLYING,
                    ["7"] = ENEMY_TYPES.AQUATIC,
                    ["8"] = ENEMY_TYPES.HUMANOID,
                }
                
                local selectedType = typeMap[typeChoice]
                if selectedType then
                    local rages = filterRages(gauChar, "type", selectedType)
                    sortRages(rages, config)
                    local display = buildRageDisplay(gauChar, rages, 1, selectedType .. " Rages")
                    api:ShowDialog(selectedType .. " Rages", display)
                end
            end
            
        elseif choice == "5" then
            -- Filter by location
            local locMenu = buildLocationMenu()
            local locChoice = api:ShowInput("Filter by Location", locMenu, "1")
            
            if locChoice and locChoice ~= "0" then
                local locNum = tonumber(locChoice)
                if locNum and locNum >= 1 and locNum <= #LOCATIONS then
                    local location = LOCATIONS[locNum]
                    local rages = filterRages(gauChar, "location", location)
                    sortRages(rages, config)
                    local display = buildRageDisplay(gauChar, rages, 1, location)
                    api:ShowDialog(location, display)
                end
            end
            
        elseif choice == "6" then
            -- Search
            local searchTerm = api:ShowInput("Search Rages", 
                "Enter search term (enemy name or ability):", "")
            
            if searchTerm and searchTerm ~= "" then
                local rages = searchRages(gauChar, searchTerm)
                sortRages(rages, config)
                local display = buildRageDisplay(gauChar, rages, 1, "Search: " .. searchTerm)
                api:ShowDialog("Search Results", display)
            end
            
        elseif choice == "7" then
            -- Statistics
            local stats = buildStatistics(gauChar)
            api:ShowDialog("Rage Statistics", stats)
            
        elseif choice == "8" then
            -- Veldt guide
            local guide = buildVeldtGuide()
            api:ShowDialog("Veldt Formation Guide", guide)
            
        else
            api:ShowDialog("Invalid Choice", "Please enter a number between 0 and 8.")
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
