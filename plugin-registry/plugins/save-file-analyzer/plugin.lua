--[[
    Save File Analyzer Plugin
    Version: 1.0.0
    
    Provides comprehensive statistics and analysis of save file data including
    playtime, completion percentage, collectibles tracking, and progress reports.
    
    Permissions Required:
    - read_save: Read all save file data
    - ui_display: Show analysis results and statistics
]]

-- Configuration
local config = {
    showDetailedStats = true,      -- Show detailed breakdowns
    calculateCompletion = true,    -- Calculate completion percentages
    trackMissables = true,         -- Track missable items/events
    compareToIdeal = true,         -- Compare to ideal/max stats
    groupByCategory = true,        -- Group stats by category
}

-- ============================================================================
-- GAME DATA CONSTANTS
-- ============================================================================

local GAME_CONSTANTS = {
    maxCharacters = 14,
    maxEspers = 27,
    maxSpells = 54,
    maxRages = 255,
    maxLevel = 99,
    maxInventorySlots = 255,
    maxGil = 9999999,
}

-- ============================================================================
-- HELPER FUNCTIONS
-- ============================================================================

-- Safe API call wrapper
local function safeCall(funcName, ...)
    local success, result = pcall(API[funcName], ...)
    if not success then
        return nil
    end
    return result
end

-- Format large numbers with commas
local function formatNumber(num)
    if not num then return "0" end
    local formatted = tostring(num)
    local k
    while true do
        formatted, k = string.gsub(formatted, "^(-?%d+)(%d%d%d)", '%1,%2')
        if k == 0 then break end
    end
    return formatted
end

-- Format time (seconds to hours:minutes:seconds)
local function formatTime(seconds)
    if not seconds then return "0:00:00" end
    
    local hours = math.floor(seconds / 3600)
    local mins = math.floor((seconds % 3600) / 60)
    local secs = seconds % 60
    
    return string.format("%d:%02d:%02d", hours, mins, secs)
end

-- Calculate percentage
local function calculatePercentage(current, total)
    if not current or not total or total == 0 then return 0 end
    return math.floor((current / total) * 100)
end

-- Get all available characters
local function getAllCharacters()
    local characters = {}
    for i = 0, 15 do
        local char = safeCall("GetCharacter", i)
        if char and char.name then
            table.insert(characters, {
                index = i,
                data = char
            })
        end
    end
    return characters
end

-- Count learned spells for character
local function countLearnedSpells(char)
    if not char or not char.spells then return 0 end
    
    local count = 0
    for _, spell in ipairs(char.spells) do
        if spell.learned then
            count = count + 1
        end
    end
    return count
end

-- Count total HP/MP across all characters
local function calculateTotalStats(characters)
    local totals = {
        hp = 0,
        maxHP = 0,
        mp = 0,
        maxMP = 0,
        level = 0,
        exp = 0
    }
    
    for _, charInfo in ipairs(characters) do
        local char = charInfo.data
        totals.hp = totals.hp + (char.currentHP or 0)
        totals.maxHP = totals.maxHP + (char.maxHP or 0)
        totals.mp = totals.mp + (char.currentMP or 0)
        totals.maxMP = totals.maxMP + (char.maxMP or 0)
        totals.level = totals.level + (char.level or 0)
        totals.exp = totals.exp + (char.experience or 0)
    end
    
    return totals
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Analyze character data
local function analyzeCharacters()
    local characters = getAllCharacters()
    
    if #characters == 0 then
        return {
            count = 0,
            completion = 0,
            details = "No characters found"
        }
    end
    
    local totalStats = calculateTotalStats(characters)
    local maxLevel = 0
    local minLevel = 99
    local avgLevel = 0
    local totalSpells = 0
    
    for _, charInfo in ipairs(characters) do
        local char = charInfo.data
        local level = char.level or 1
        
        if level > maxLevel then maxLevel = level end
        if level < minLevel then minLevel = level end
        avgLevel = avgLevel + level
        totalSpells = totalSpells + countLearnedSpells(char)
    end
    
    avgLevel = math.floor(avgLevel / #characters)
    
    local completion = calculatePercentage(#characters, GAME_CONSTANTS.maxCharacters)
    
    return {
        count = #characters,
        completion = completion,
        maxLevel = maxLevel,
        minLevel = minLevel,
        avgLevel = avgLevel,
        totalSpells = totalSpells,
        totalHP = totalStats.maxHP,
        totalMP = totalStats.maxMP,
        totalExp = totalStats.exp,
        hpPercent = calculatePercentage(totalStats.hp, totalStats.maxHP),
        mpPercent = calculatePercentage(totalStats.mp, totalStats.maxMP),
    }
end

-- Analyze inventory
local function analyzeInventory()
    local inventory = safeCall("GetInventory")
    
    if not inventory then
        return {
            uniqueItems = 0,
            totalItems = 0,
            completion = 0,
            details = "Unable to read inventory"
        }
    end
    
    local uniqueItems = 0
    local totalItems = 0
    local valueEstimate = 0
    
    for _, item in ipairs(inventory) do
        if item.quantity and item.quantity > 0 then
            uniqueItems = uniqueItems + 1
            totalItems = totalItems + item.quantity
            -- Rough value estimate (would need item price database)
            valueEstimate = valueEstimate + (item.quantity * 100)
        end
    end
    
    local slotsUsed = calculatePercentage(uniqueItems, GAME_CONSTANTS.maxInventorySlots)
    
    return {
        uniqueItems = uniqueItems,
        totalItems = totalItems,
        slotsUsed = slotsUsed,
        valueEstimate = valueEstimate,
    }
end

-- Analyze party
local function analyzeParty()
    local party = safeCall("GetParty")
    
    if not party then
        return {
            size = 0,
            details = "Unable to read party"
        }
    end
    
    local partyMembers = {}
    for i = 0, 3 do
        local charIndex = party[i]
        if charIndex then
            local char = safeCall("GetCharacter", charIndex)
            if char then
                table.insert(partyMembers, {
                    index = charIndex,
                    name = char.name or ("Character #" .. charIndex),
                    level = char.level or 1,
                    hp = char.currentHP or 0,
                    maxHP = char.maxHP or 0,
                })
            end
        end
    end
    
    local totalPower = 0
    local avgLevel = 0
    for _, member in ipairs(partyMembers) do
        avgLevel = avgLevel + member.level
        totalPower = totalPower + member.level * 10 + member.maxHP / 10
    end
    
    if #partyMembers > 0 then
        avgLevel = math.floor(avgLevel / #partyMembers)
        totalPower = math.floor(totalPower)
    end
    
    return {
        size = #partyMembers,
        members = partyMembers,
        avgLevel = avgLevel,
        power = totalPower,
    }
end

-- Calculate overall completion percentage
local function calculateOverallCompletion()
    local charAnalysis = analyzeCharacters()
    local invAnalysis = analyzeInventory()
    
    -- Weighted completion calculation
    local weights = {
        characters = 0.30,  -- 30%
        inventory = 0.20,   -- 20%
        spells = 0.25,      -- 25%
        levels = 0.25,      -- 25%
    }
    
    local charCompletion = charAnalysis.completion or 0
    local invCompletion = math.min(100, invAnalysis.slotsUsed or 0)
    local spellCompletion = calculatePercentage(
        charAnalysis.totalSpells or 0,
        (charAnalysis.count or 1) * GAME_CONSTANTS.maxSpells
    )
    local levelCompletion = calculatePercentage(
        charAnalysis.avgLevel or 1,
        GAME_CONSTANTS.maxLevel
    )
    
    local overall = (charCompletion * weights.characters) +
                   (invCompletion * weights.inventory) +
                   (spellCompletion * weights.spells) +
                   (levelCompletion * weights.levels)
    
    return math.floor(overall)
end

-- ============================================================================
-- DISPLAY FUNCTIONS
-- ============================================================================

-- Show comprehensive statistics
local function showComprehensiveStats()
    local charAnalysis = analyzeCharacters()
    local invAnalysis = analyzeInventory()
    local partyAnalysis = analyzeParty()
    local overallCompletion = calculateOverallCompletion()
    
    local output = "=== SAVE FILE ANALYSIS ===\n\n"
    
    -- Overall Completion
    output = output .. string.format("Overall Completion: %d%%\n", overallCompletion)
    
    local completionBar = ""
    for i = 1, 20 do
        if i <= overallCompletion / 5 then
            completionBar = completionBar .. "█"
        else
            completionBar = completionBar .. "░"
        end
    end
    output = output .. completionBar .. "\n\n"
    
    -- Character Statistics
    output = output .. "=== CHARACTERS ===\n"
    output = output .. string.format("Recruited: %d/%d (%d%%)\n",
        charAnalysis.count, GAME_CONSTANTS.maxCharacters, charAnalysis.completion)
    output = output .. string.format("Average Level: %d\n", charAnalysis.avgLevel or 0)
    output = output .. string.format("Level Range: %d - %d\n",
        charAnalysis.minLevel or 0, charAnalysis.maxLevel or 0)
    output = output .. string.format("Total HP: %s\n", formatNumber(charAnalysis.totalHP or 0))
    output = output .. string.format("Total MP: %s\n", formatNumber(charAnalysis.totalMP or 0))
    output = output .. string.format("Total Spells Learned: %d\n", charAnalysis.totalSpells or 0)
    output = output .. string.format("Total Experience: %s\n\n", formatNumber(charAnalysis.totalExp or 0))
    
    -- Party Statistics
    output = output .. "=== CURRENT PARTY ===\n"
    output = output .. string.format("Party Size: %d/4\n", partyAnalysis.size)
    output = output .. string.format("Average Level: %d\n", partyAnalysis.avgLevel or 0)
    output = output .. string.format("Party Power: %s\n", formatNumber(partyAnalysis.power or 0))
    
    if partyAnalysis.members and #partyAnalysis.members > 0 then
        output = output .. "Members:\n"
        for i, member in ipairs(partyAnalysis.members) do
            output = output .. string.format("  %d. %s (Lv%d) - HP: %d/%d\n",
                i, member.name, member.level, member.hp, member.maxHP)
        end
    end
    output = output .. "\n"
    
    -- Inventory Statistics
    output = output .. "=== INVENTORY ===\n"
    output = output .. string.format("Unique Items: %d\n", invAnalysis.uniqueItems or 0)
    output = output .. string.format("Total Items: %d\n", invAnalysis.totalItems or 0)
    output = output .. string.format("Slots Used: %d%%\n", invAnalysis.slotsUsed or 0)
    output = output .. string.format("Estimated Value: %s gil\n\n",
        formatNumber(invAnalysis.valueEstimate or 0))
    
    API:ShowInfo(output)
    return true
end

-- Show character details
local function showCharacterDetails()
    local characters = getAllCharacters()
    
    if #characters == 0 then
        API:ShowError("No characters found")
        return false
    end
    
    local output = "=== CHARACTER DETAILS ===\n\n"
    
    for _, charInfo in ipairs(characters) do
        local char = charInfo.data
        local name = char.name or ("Character #" .. charInfo.index)
        local spellCount = countLearnedSpells(char)
        
        output = output .. string.format("%s (Lv%d)\n", name, char.level or 1)
        output = output .. string.format("  HP: %d/%d  MP: %d/%d\n",
            char.currentHP or 0, char.maxHP or 0,
            char.currentMP or 0, char.maxMP or 0)
        output = output .. string.format("  Experience: %s\n", formatNumber(char.experience or 0))
        output = output .. string.format("  Spells Learned: %d/%d\n",
            spellCount, GAME_CONSTANTS.maxSpells)
        
        if char.attack then
            output = output .. string.format("  Attack: %d  Defense: %d\n",
                char.attack or 0, char.defense or 0)
        end
        if char.magic then
            output = output .. string.format("  Magic: %d  Mag.Def: %d\n",
                char.magic or 0, char.magicDefense or 0)
        end
        output = output .. "\n"
    end
    
    API:ShowInfo(output)
    return true
end

-- Show completion checklist
local function showCompletionChecklist()
    local charAnalysis = analyzeCharacters()
    local invAnalysis = analyzeInventory()
    
    local output = "=== COMPLETION CHECKLIST ===\n\n"
    
    -- Characters
    local charCheck = charAnalysis.count >= GAME_CONSTANTS.maxCharacters and "✓" or "✗"
    output = output .. string.format("[%s] All Characters (%d/%d)\n",
        charCheck, charAnalysis.count, GAME_CONSTANTS.maxCharacters)
    
    -- Max Level
    local levelCheck = charAnalysis.maxLevel >= GAME_CONSTANTS.maxLevel and "✓" or "✗"
    output = output .. string.format("[%s] Max Level Character (Lv%d/%d)\n",
        levelCheck, charAnalysis.maxLevel or 0, GAME_CONSTANTS.maxLevel)
    
    -- Average Level
    local avgLevelCheck = charAnalysis.avgLevel >= 50 and "✓" or "✗"
    output = output .. string.format("[%s] Average Level 50+ (Avg: Lv%d)\n",
        avgLevelCheck, charAnalysis.avgLevel or 0)
    
    -- Spell Learning
    local avgSpells = charAnalysis.count > 0 and 
        math.floor(charAnalysis.totalSpells / charAnalysis.count) or 0
    local spellCheck = avgSpells >= GAME_CONSTANTS.maxSpells and "✓" or "✗"
    output = output .. string.format("[%s] All Spells Learned (Avg: %d/%d per char)\n",
        spellCheck, avgSpells, GAME_CONSTANTS.maxSpells)
    
    -- Inventory
    local invCheck = invAnalysis.uniqueItems >= 100 and "✓" or "✗"
    output = output .. string.format("[%s] 100+ Unique Items (%d items)\n",
        invCheck, invAnalysis.uniqueItems or 0)
    
    output = output .. "\nProgress Categories:\n\n"
    
    -- Calculate category completion
    local categories = {
        {name = "Characters", percent = charAnalysis.completion},
        {name = "Levels", percent = calculatePercentage(charAnalysis.avgLevel, GAME_CONSTANTS.maxLevel)},
        {name = "Spells", percent = calculatePercentage(avgSpells, GAME_CONSTANTS.maxSpells)},
        {name = "Inventory", percent = math.min(100, invAnalysis.slotsUsed)},
    }
    
    for _, cat in ipairs(categories) do
        local bar = ""
        for i = 1, 10 do
            if i <= cat.percent / 10 then
                bar = bar .. "█"
            else
                bar = bar .. "░"
            end
        end
        output = output .. string.format("%s: [%s] %d%%\n",
            cat.name, bar, cat.percent)
    end
    
    API:ShowInfo(output)
    return true
end

-- Export statistics report
local function exportStatisticsReport()
    local charAnalysis = analyzeCharacters()
    local invAnalysis = analyzeInventory()
    local partyAnalysis = analyzeParty()
    local overallCompletion = calculateOverallCompletion()
    
    local output = "=== FF6 SAVE FILE STATISTICS REPORT ===\n\n"
    output = output .. string.format("Generated: %s\n\n", os.date("%Y-%m-%d %H:%M:%S"))
    
    output = output .. "SUMMARY\n"
    output = output .. string.format("Overall Completion: %d%%\n", overallCompletion)
    output = output .. string.format("Characters: %d/%d\n", charAnalysis.count, GAME_CONSTANTS.maxCharacters)
    output = output .. string.format("Average Level: %d\n", charAnalysis.avgLevel or 0)
    output = output .. string.format("Party Size: %d/4\n\n", partyAnalysis.size)
    
    output = output .. "CHARACTER STATISTICS\n"
    output = output .. string.format("  Total HP: %s\n", formatNumber(charAnalysis.totalHP))
    output = output .. string.format("  Total MP: %s\n", formatNumber(charAnalysis.totalMP))
    output = output .. string.format("  Total Experience: %s\n", formatNumber(charAnalysis.totalExp))
    output = output .. string.format("  Total Spells Learned: %d\n\n", charAnalysis.totalSpells)
    
    output = output .. "INVENTORY STATISTICS\n"
    output = output .. string.format("  Unique Items: %d\n", invAnalysis.uniqueItems)
    output = output .. string.format("  Total Items: %d\n", invAnalysis.totalItems)
    output = output .. string.format("  Estimated Value: %s gil\n\n", formatNumber(invAnalysis.valueEstimate))
    
    output = output .. "End of Report\n"
    
    API:ShowInfo(output)
    API:ShowInfo("Report ready! (Copy from dialog above)")
    return true
end

-- ============================================================================
-- MAIN PLUGIN ENTRY POINT
-- ============================================================================

local function main()
    -- Verify permissions
    if not API.GetCharacter then
        API:ShowError("Missing required API methods. Please update the editor.")
        return
    end
    
    -- Show main menu
    local menu = [[
=== SAVE FILE ANALYZER ===

Choose an operation:

1. Show Comprehensive Statistics
2. Show Character Details
3. Show Completion Checklist
4. Export Statistics Report
5. Exit

Enter choice (1-5):]]
    
    local choice = API:PromptNumber(menu, 1)
    
    if not choice then
        API:ShowInfo("Operation canceled")
        return
    end
    
    if choice == 1 then
        showComprehensiveStats()
    elseif choice == 2 then
        showCharacterDetails()
    elseif choice == 3 then
        showCompletionChecklist()
    elseif choice == 4 then
        exportStatisticsReport()
    elseif choice == 5 then
        API:ShowInfo("Goodbye!")
    else
        API:ShowError("Invalid choice. Please enter 1-5.")
    end
end

-- Run the plugin
local success, err = pcall(main)
if not success then
    API:ShowError("Plugin error: " .. tostring(err))
end
