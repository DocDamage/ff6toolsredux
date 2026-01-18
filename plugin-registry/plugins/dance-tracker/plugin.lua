-- ============================================================================
-- Mog's Dance Tracker Plugin
-- ============================================================================
-- @id: dance-tracker
-- @name: Mog's Dance Tracker
-- @version: 1.0.0
-- @author: FF6 Editor Team
-- @description: Track Mog's 8 dances with locations and effects
-- @permissions: read_save, ui_display

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    showLearnedOnly = false,        -- Show only learned dances
    showUnlearnedOnly = false,      -- Show only unlearned dances
    showEffects = true,             -- Show dance abilities
    sortByName = true,              -- Sort by name (false = by Dance ID)
}

-- ============================================================================
-- Dance Database
-- ============================================================================
-- Complete database of all 8 dances with locations and effects
local DANCE_DATABASE = {
    {
        id = 1,
        name = "Wind Song",
        location = "Narshe (outside, where Mog is found)",
        worldAvailability = "WoB and WoR",
        mapArea = "Narshe Exterior",
        effects = {
            "Sonic Boom - Wind damage to one enemy",
            "Plasma - Lightning damage to all enemies",
            "Harvester - Death to one enemy",
            "Wombat - Poison damage to all enemies"
        },
        description = "Random wind and lightning attacks with instant death"
    },
    {
        id = 2,
        name = "Forest Suite",
        location = "Phantom Forest",
        worldAvailability = "WoB and WoR",
        mapArea = "Phantom Forest",
        effects = {
            "Will o' Wisp - Fire damage to random enemies",
            "Apparition - Drain HP from one enemy",
            "Poltergeist - Image status to party",
            "Specter - Invisible status to party"
        },
        description = "Fire attacks and party-enhancing phantom effects"
    },
    {
        id = 3,
        name = "Desert Aria",
        location = "Figaro Castle (WoB) / Kohlingen Desert (WoR)",
        worldAvailability = "WoB and WoR",
        mapArea = "Desert regions near Figaro",
        effects = {
            "Sand Storm - Wind damage to all enemies",
            "Antlion - Earth damage to one enemy",
            "Dissolve - Remove positive statuses from all enemies",
            "Kitty - Set HP to single digits for all enemies"
        },
        description = "Desert-themed attacks with powerful HP manipulation"
    },
    {
        id = 4,
        name = "Love Sonata",
        location = "Darill's Tomb",
        worldAvailability = "WoR only",
        mapArea = "Darill's Tomb interior",
        effects = {
            "Aurora - Restore HP to party",
            "Serenade - Sleep all enemies",
            "Minuet - Haste party",
            "Love Token - Regen to party"
        },
        description = "Healing and support effects for the party"
    },
    {
        id = 5,
        name = "Earth Blues",
        location = "Mt. Kolts",
        worldAvailability = "WoB and WoR",
        mapArea = "Mt. Kolts interior",
        effects = {
            "Avalanche - Earth damage to all enemies",
            "Cave In - Earth damage to all enemies",
            "Landslide - Earth damage to all enemies",
            "Sonic Boom - Wind damage to one enemy"
        },
        description = "Multiple earth-based attacks"
    },
    {
        id = 6,
        name = "Water Rondo",
        location = "Serpent Trench / Lethe River",
        worldAvailability = "WoB and WoR",
        mapArea = "Water areas (Serpent Trench, Lethe River)",
        effects = {
            "El Niño - Water damage to all enemies",
            "Specter - Invisible status to party",
            "Plasma - Lightning damage to all enemies",
            "Apparition - Drain HP from one enemy"
        },
        description = "Water and lightning attacks with drain effects"
    },
    {
        id = 7,
        name = "Dusk Requiem",
        location = "Zozo / Any town at night",
        worldAvailability = "WoB and WoR",
        mapArea = "Zozo or towns (evening hours)",
        effects = {
            "Pois. Frog - Poison all enemies",
            "Evil Toot - Confusion to all enemies",
            "Sour Mouth - Reduce level of all enemies",
            "Snare - Stop one enemy"
        },
        description = "Status ailment attacks for crowd control"
    },
    {
        id = 8,
        name = "Snowman Rondo",
        location = "Narshe Mines / Umaro's Cave",
        worldAvailability = "WoB and WoR",
        mapArea = "Narshe underground areas",
        effects = {
            "Absolute Zero - Ice damage to all enemies",
            "Surge - Ice damage to all enemies",
            "Avalanche - Earth damage to all enemies",
            "Snare - Stop one enemy"
        },
        description = "Ice and earth attacks with stop effect"
    }
}

-- Dance IDs mapped to potential tracking IDs in the game
-- Note: This mapping may need adjustment based on actual game data structure
local DANCE_TRACKING_IDS = {
    [1] = 1,   -- Wind Song
    [2] = 2,   -- Forest Suite
    [3] = 3,   -- Desert Aria
    [4] = 4,   -- Love Sonata
    [5] = 5,   -- Earth Blues
    [6] = 6,   -- Water Rondo
    [7] = 7,   -- Dusk Requiem
    [8] = 8,   -- Snowman Rondo
}

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Find Mog's character index
local function findMog(api)
    -- Try to find Mog by name
    for i = 0, 15 do
        local char = api:GetCharacter(i)
        if char and char.Name then
            -- Mog's name variations
            if char.Name == "Mog" or char.Name == "MOG" or 
               char.Name:lower():find("mog") then
                return i, char
            end
        end
    end
    return nil, nil
end

-- Check if a dance is learned
-- Note: This is a placeholder implementation. The actual method to check
-- dance learning status would depend on the game's data structure.
local function isDanceLearned(mogChar, danceID)
    if not mogChar then
        return false
    end
    
    -- Placeholder: In a real implementation, you would check specific
    -- character flags or data fields that track learned dances
    -- For now, we'll check if Mog has certain abilities or flags
    
    -- This would need to be replaced with actual dance tracking logic
    -- based on the save file structure
    
    -- For demonstration, we'll return false (not learned) for all
    -- A real implementation would check mogChar.Dances or similar field
    return false
end

-- Count learned dances
local function countLearnedDances(mogChar)
    local count = 0
    for danceID = 1, 8 do
        if isDanceLearned(mogChar, danceID) then
            count = count + 1
        end
    end
    return count
end

-- Format dance effects list
local function formatEffects(effects)
    local formatted = {}
    for i, effect in ipairs(effects) do
        table.insert(formatted, "       • " .. effect)
    end
    return table.concat(formatted, "\n")
end

-- Format dance entry
local function formatDanceEntry(dance, learned, config)
    local status = learned and "[✓]" or "[ ]"
    
    local lines = {}
    table.insert(lines, string.format("%s %d. %s", status, dance.id, dance.name))
    table.insert(lines, string.format("       Location: %s", dance.location))
    table.insert(lines, string.format("       Available: %s", dance.worldAvailability))
    
    if config.showEffects then
        table.insert(lines, "       Effects (random selection):")
        table.insert(lines, formatEffects(dance.effects))
    end
    
    return table.concat(lines, "\n")
end

-- Filter dances based on configuration
local function filterDances(mogChar, config)
    local filtered = {}
    
    for _, dance in ipairs(DANCE_DATABASE) do
        local learned = isDanceLearned(mogChar, dance.id)
        
        -- Apply filters
        local include = true
        
        if config.showLearnedOnly and not learned then
            include = false
        end
        
        if config.showUnlearnedOnly and learned then
            include = false
        end
        
        if include then
            table.insert(filtered, {dance = dance, learned = learned})
        end
    end
    
    return filtered
end

-- Sort dances
local function sortDances(dances, config)
    if config.sortByName then
        table.sort(dances, function(a, b)
            return a.dance.name < b.dance.name
        end)
    else
        table.sort(dances, function(a, b)
            return a.dance.id < b.dance.id
        end)
    end
end

-- Build dance tracker display
local function buildDanceDisplay(mogChar, config)
    local sections = {}
    
    -- Header
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "MOG'S DANCE TRACKER - Dance Abilities")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    -- Statistics
    local learnedCount = countLearnedDances(mogChar)
    local totalCount = 8
    local percentage = (learnedCount / totalCount) * 100
    
    table.insert(sections, "")
    table.insert(sections, string.format("Progress: %d / %d Dances learned (%.1f%%)", 
        learnedCount, totalCount, percentage))
    table.insert(sections, "")
    
    -- Get and sort dances
    local filteredDances = filterDances(mogChar, config)
    sortDances(filteredDances, config)
    
    if #filteredDances == 0 then
        table.insert(sections, "No dances match the current filter settings.")
    else
        -- Display dances
        table.insert(sections, "Dance List:")
        table.insert(sections, string.rep("-", 60))
        
        for _, entry in ipairs(filteredDances) do
            table.insert(sections, formatDanceEntry(entry.dance, entry.learned, config))
            table.insert(sections, "")
        end
    end
    
    -- Legend
    table.insert(sections, string.rep("-", 60))
    table.insert(sections, "Legend:")
    table.insert(sections, "  [✓] = Learned    [ ] = Not Learned")
    table.insert(sections, "  WoB = World of Balance    WoR = World of Ruin")
    table.insert(sections, "")
    table.insert(sections, "Note: Each dance randomly selects one of four abilities")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    return table.concat(sections, "\n")
end

-- Build location guide
local function buildLocationGuide()
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "DANCE LOCATION GUIDE")
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "")
    table.insert(sections, "To learn a dance, Mog must step on the location where")
    table.insert(sections, "the dance can be learned. The dance is learned immediately.")
    table.insert(sections, "")
    table.insert(sections, string.rep("-", 60))
    
    for _, dance in ipairs(DANCE_DATABASE) do
        table.insert(sections, "")
        table.insert(sections, string.format("%d. %s", dance.id, dance.name))
        table.insert(sections, string.format("   Location: %s", dance.location))
        table.insert(sections, string.format("   Map Area: %s", dance.mapArea))
        table.insert(sections, string.format("   Available: %s", dance.worldAvailability))
        
        -- Special notes for certain dances
        if dance.id == 4 then
            table.insert(sections, "   Note: WoR only - cannot be learned in WoB")
        elseif dance.id == 7 then
            table.insert(sections, "   Note: Zozo or any town during evening hours")
        end
    end
    
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    return table.concat(sections, "\n")
end

-- Build statistics display
local function buildStatistics(mogChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "DANCE STATISTICS")
    table.insert(sections, "=" .. string.rep("=", 59))
    table.insert(sections, "")
    
    -- Overall stats
    local learnedCount = countLearnedDances(mogChar)
    local percentage = (learnedCount / 8) * 100
    
    table.insert(sections, string.format("Total Dances Learned: %d / 8 (%.1f%%)", learnedCount, percentage))
    table.insert(sections, "")
    
    -- World availability
    local wobCount = 0
    local worOnlyCount = 0
    
    for _, dance in ipairs(DANCE_DATABASE) do
        if dance.worldAvailability == "WoR only" then
            worOnlyCount = worOnlyCount + 1
        else
            wobCount = wobCount + 1
        end
    end
    
    table.insert(sections, string.format("Available in WoB: %d dances", wobCount))
    table.insert(sections, string.format("Available in WoR only: %d dance(s)", worOnlyCount))
    table.insert(sections, "")
    
    -- Dance categories
    table.insert(sections, "Dance Categories:")
    table.insert(sections, "  - Elemental Dances: 6 (Wind, Forest, Desert, Earth, Water, Snow)")
    table.insert(sections, "  - Support Dance: 1 (Love Sonata - healing/buffs)")
    table.insert(sections, "  - Status Dance: 1 (Dusk Requiem - status ailments)")
    table.insert(sections, "")
    
    -- Completion status
    if learnedCount == 8 then
        table.insert(sections, "🎉 PERFECT! All dances learned! 🎉")
    else
        table.insert(sections, string.format("Dances Remaining: %d", 8 - learnedCount))
    end
    
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 59))
    
    return table.concat(sections, "\n")
end

-- Build main menu
local function buildMainMenu(mogChar)
    local learnedCount = countLearnedDances(mogChar)
    local menu = "=" .. string.rep("=", 59) .. "\n"
    menu = menu .. "MOG'S DANCE TRACKER - Main Menu\n"
    menu = menu .. "=" .. string.rep("=", 59) .. "\n\n"
    menu = menu .. string.format("Current Progress: %d / 8 Dances learned\n\n", learnedCount)
    menu = menu .. "Select an option:\n\n"
    menu = menu .. "1. View All Dances (Complete List)\n"
    menu = menu .. "2. View Learned Dances Only\n"
    menu = menu .. "3. View Unlearned Dances Only\n"
    menu = menu .. "4. Location Guide (Where to Learn)\n"
    menu = menu .. "5. View Dance Statistics\n"
    menu = menu .. "0. Exit\n\n"
    menu = menu .. "Enter your choice (0-5):"
    
    return menu
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
    
    api:Log("info", "Dance Tracker plugin started")
    
    -- Find Mog
    local mogIndex, mogChar = findMog(api)
    if not mogIndex then
        api:ShowDialog("Error", "Mog not found in save file. Make sure you have recruited Mog.")
        api:Log("error", "Mog character not found")
        return false
    end
    
    api:Log("info", string.format("Found Mog at character index: %d", mogIndex))
    
    -- Main menu loop
    while true do
        local menu = buildMainMenu(mogChar)
        local choice = api:ShowInput("Dance Tracker - Main Menu", menu, "1")
        
        if not choice or choice == "0" then
            api:Log("info", "User exited Dance Tracker")
            return true
        end
        
        if choice == "1" then
            -- View all dances
            config.showLearnedOnly = false
            config.showUnlearnedOnly = false
            local display = buildDanceDisplay(mogChar, config)
            api:ShowDialog("All Dances", display)
            
        elseif choice == "2" then
            -- View learned only
            config.showLearnedOnly = true
            config.showUnlearnedOnly = false
            local display = buildDanceDisplay(mogChar, config)
            api:ShowDialog("Learned Dances", display)
            
        elseif choice == "3" then
            -- View unlearned only
            config.showLearnedOnly = false
            config.showUnlearnedOnly = true
            local display = buildDanceDisplay(mogChar, config)
            api:ShowDialog("Unlearned Dances", display)
            
        elseif choice == "4" then
            -- Location guide
            local guide = buildLocationGuide()
            api:ShowDialog("Dance Location Guide", guide)
            
        elseif choice == "5" then
            -- View statistics
            local stats = buildStatistics(mogChar)
            api:ShowDialog("Dance Statistics", stats)
            
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
