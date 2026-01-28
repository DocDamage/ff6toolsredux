--[[
    Magic Learning Planner Plugin
    Version: 1.0.0
    
    Tracks spell learning progress, recommends optimal esper assignments,
    and calculates AP requirements for complete spell mastery.
    
    Permissions Required:
    - read_save: Read character spell data and esper assignments
    - write_save: Apply esper recommendations
    - ui_display: Show progress tracking and recommendations
]]

-- Configuration
local config = {
    showLearnedOnly = false,     -- Show only learned spells
    sortByProgress = true,        -- Sort by learning progress
    highlightMissing = true,      -- Highlight unlearned spells
    calculateAPRequired = true,   -- Calculate AP needed for mastery
    groupBySchool = true,         -- Group spells by magic school
}

-- ============================================================================
-- SPELL AND ESPER DATABASE
-- ============================================================================

-- Spell schools
local SCHOOL_BLACK = "Black Magic"
local SCHOOL_WHITE = "White Magic"
local SCHOOL_BLUE = "Blue Magic"
local SCHOOL_RED = "Red Magic"
local SCHOOL_SPECIAL = "Special"

-- Spell database (simplified - real implementation would load from game data)
local spellDB = {
    -- Black Magic
    {id = 0, name = "Fire", school = SCHOOL_BLACK, power = 20},
    {id = 1, name = "Ice", school = SCHOOL_BLACK, power = 20},
    {id = 2, name = "Bolt", school = SCHOOL_BLACK, power = 20},
    {id = 3, name = "Poison", school = SCHOOL_BLACK, power = 10},
    {id = 4, name = "Drain", school = SCHOOL_BLACK, power = 25},
    {id = 5, name = "Fira", school = SCHOOL_BLACK, power = 40},
    {id = 6, name = "Blizzara", school = SCHOOL_BLACK, power = 40},
    {id = 7, name = "Thundara", school = SCHOOL_BLACK, power = 40},
    {id = 8, name = "Bio", school = SCHOOL_BLACK, power = 30},
    {id = 9, name = "Firaga", school = SCHOOL_BLACK, power = 60},
    {id = 10, name = "Blizzaga", school = SCHOOL_BLACK, power = 60},
    {id = 11, name = "Thundaga", school = SCHOOL_BLACK, power = 60},
    {id = 12, name = "Flare", school = SCHOOL_BLACK, power = 80},
    {id = 13, name = "Ultima", school = SCHOOL_BLACK, power = 100},
    
    -- White Magic
    {id = 20, name = "Cure", school = SCHOOL_WHITE, power = 20},
    {id = 21, name = "Cura", school = SCHOOL_WHITE, power = 40},
    {id = 22, name = "Curaga", school = SCHOOL_WHITE, power = 60},
    {id = 23, name = "Raise", school = SCHOOL_WHITE, power = 0},
    {id = 24, name = "Arise", school = SCHOOL_WHITE, power = 0},
    {id = 25, name = "Regen", school = SCHOOL_WHITE, power = 0},
    {id = 26, name = "Esuna", school = SCHOOL_WHITE, power = 0},
    {id = 27, name = "Dispel", school = SCHOOL_WHITE, power = 0},
    {id = 28, name = "Shell", school = SCHOOL_WHITE, power = 0},
    {id = 29, name = "Protect", school = SCHOOL_WHITE, power = 0},
    {id = 30, name = "Reflect", school = SCHOOL_WHITE, power = 0},
    {id = 31, name = "Holy", school = SCHOOL_WHITE, power = 100},
    
    -- Blue Magic
    {id = 40, name = "Scan", school = SCHOOL_BLUE, power = 0},
    {id = 41, name = "Slow", school = SCHOOL_BLUE, power = 0},
    {id = 42, name = "Rasp", school = SCHOOL_BLUE, power = 15},
    {id = 43, name = "Osmose", school = SCHOOL_BLUE, power = 0},
    {id = 44, name = "Haste", school = SCHOOL_BLUE, power = 0},
    {id = 45, name = "Stop", school = SCHOOL_BLUE, power = 0},
    {id = 46, name = "Berserk", school = SCHOOL_BLUE, power = 0},
    {id = 47, name = "Slow 2", school = SCHOOL_BLUE, power = 0},
    {id = 48, name = "Haste 2", school = SCHOOL_BLUE, power = 0},
    {id = 49, name = "Vanish", school = SCHOOL_BLUE, power = 0},
    
    -- Red Magic (Effect spells)
    {id = 60, name = "Mute", school = SCHOOL_RED, power = 0},
    {id = 61, name = "Sleep", school = SCHOOL_RED, power = 0},
    {id = 62, name = "Confuse", school = SCHOOL_RED, power = 0},
    {id = 63, name = "Break", school = SCHOOL_RED, power = 0},
    {id = 64, name = "Death", school = SCHOOL_RED, power = 0},
    {id = 65, name = "Banish", school = SCHOOL_RED, power = 0},
    
    -- Special
    {id = 80, name = "Quick", school = SCHOOL_SPECIAL, power = 0},
    {id = 81, name = "Warp", school = SCHOOL_SPECIAL, power = 0},
    {id = 82, name = "Merton", school = SCHOOL_SPECIAL, power = 90},
    {id = 83, name = "Meteor", school = SCHOOL_SPECIAL, power = 85},
    {id = 84, name = "Quake", school = SCHOOL_SPECIAL, power = 70},
}

-- Esper database with spells they teach
local esperDB = {
    {id = 0, name = "Ramuh", spells = {
        {spellID = 2, learnRate = 10},  -- Bolt
        {spellID = 7, learnRate = 5},   -- Thundara
        {spellID = 11, learnRate = 2},  -- Thundaga
    }},
    {id = 1, name = "Ifrit", spells = {
        {spellID = 0, learnRate = 10},  -- Fire
        {spellID = 5, learnRate = 5},   -- Fira
        {spellID = 9, learnRate = 1},   -- Firaga
    }},
    {id = 2, name = "Shiva", spells = {
        {spellID = 1, learnRate = 10},  -- Ice
        {spellID = 6, learnRate = 5},   -- Blizzara
        {spellID = 10, learnRate = 1},  -- Blizzaga
    }},
    {id = 3, name = "Unicorn", spells = {
        {spellID = 20, learnRate = 5},  -- Cure
        {spellID = 21, learnRate = 3},  -- Cura
        {spellID = 25, learnRate = 2},  -- Regen
    }},
    {id = 4, name = "Carbuncle", spells = {
        {spellID = 30, learnRate = 5},  -- Reflect
        {spellID = 28, learnRate = 3},  -- Shell
        {spellID = 29, learnRate = 3},  -- Protect
    }},
    {id = 5, name = "Bahamut", spells = {
        {spellID = 12, learnRate = 2},  -- Flare
        {spellID = 83, learnRate = 1},  -- Meteor
    }},
    {id = 6, name = "Alexandr", spells = {
        {spellID = 28, learnRate = 10},  -- Shell
        {spellID = 29, learnRate = 10},  -- Protect
        {spellID = 31, learnRate = 2},   -- Holy
    }},
}

-- ============================================================================
-- HELPER FUNCTIONS
-- ============================================================================

-- Safe API call wrapper
local function safeCall(funcName, ...)
    local success, result = pcall(API[funcName], ...)
    if not success then
        API:ShowError("Error calling " .. funcName .. ": " .. tostring(result))
        return nil
    end
    return result
end

-- Get spell by ID
local function getSpellByID(spellID)
    for _, spell in ipairs(spellDB) do
        if spell.id == spellID then
            return spell
        end
    end
    return nil
end

-- Get esper by ID
local function getEsperByID(esperID)
    for _, esper in ipairs(esperDB) do
        if esper.id == esperID then
            return esper
        end
    end
    return nil
end

-- Check if character has learned spell
local function hasLearnedSpell(char, spellID)
    if not char or not char.spells then return false end
    
    for _, learnedSpell in ipairs(char.spells) do
        if learnedSpell.id == spellID and learnedSpell.learned then
            return true
        end
    end
    return false
end

-- Get spell learning progress for character
local function getSpellProgress(char, spellID)
    if not char or not char.spells then return 0 end
    
    for _, learnedSpell in ipairs(char.spells) do
        if learnedSpell.id == spellID then
            return learnedSpell.progress or 0
        end
    end
    return 0
end

-- Calculate AP required to learn spell
local function calculateAPRequired(learnRate, currentProgress)
    if not learnRate or learnRate == 0 then
        return 999  -- Can't learn from this esper
    end
    
    local totalAPNeeded = 100 / learnRate
    local currentAP = currentProgress * totalAPNeeded / 100
    return math.ceil(totalAPNeeded - currentAP)
end

-- Get character's spell learning progress by school
local function getSpellProgressBySchool(char)
    local progress = {
        [SCHOOL_BLACK] = {learned = 0, total = 0},
        [SCHOOL_WHITE] = {learned = 0, total = 0},
        [SCHOOL_BLUE] = {learned = 0, total = 0},
        [SCHOOL_RED] = {learned = 0, total = 0},
        [SCHOOL_SPECIAL] = {learned = 0, total = 0},
    }
    
    for _, spell in ipairs(spellDB) do
        local school = spell.school
        progress[school].total = progress[school].total + 1
        
        if hasLearnedSpell(char, spell.id) then
            progress[school].learned = progress[school].learned + 1
        end
    end
    
    return progress
end

-- Find espers that teach a spell
local function findEspersThatTeachSpell(spellID)
    local espers = {}
    
    for _, esper in ipairs(esperDB) do
        for _, spell in ipairs(esper.spells) do
            if spell.spellID == spellID then
                table.insert(espers, {
                    esper = esper,
                    learnRate = spell.learnRate
                })
            end
        end
    end
    
    return espers
end

-- ============================================================================
-- OPERATION FUNCTIONS
-- ============================================================================

-- Show spell progress for character
local function showSpellProgress(charIndex)
    local char = safeCall("GetCharacter", charIndex)
    if not char then
        API:ShowError("Character #" .. charIndex .. " not found")
        return false
    end
    
    local name = char.name or ("Character #" .. charIndex)
    local output = "=== SPELL LEARNING PROGRESS ===\n\n"
    output = output .. string.format("%s (Lv%d)\n\n", name, char.level or 1)
    
    if config.groupBySchool then
        -- Show progress by school
        local progressBySchool = getSpellProgressBySchool(char)
        
        for school, data in pairs(progressBySchool) do
            local percentage = data.total > 0 and (data.learned / data.total * 100) or 0
            output = output .. string.format("%s: %d/%d (%.1f%%)\n",
                school, data.learned, data.total, percentage)
        end
        
        output = output .. "\nDetailed Progress:\n\n"
        
        -- Show spells grouped by school
        for _, school in ipairs({SCHOOL_BLACK, SCHOOL_WHITE, SCHOOL_BLUE, SCHOOL_RED, SCHOOL_SPECIAL}) do
            output = output .. school .. ":\n"
            
            for _, spell in ipairs(spellDB) do
                if spell.school == school then
                    local learned = hasLearnedSpell(char, spell.id)
                    local progress = getSpellProgress(char, spell.id)
                    local status = learned and "✓" or (progress > 0 and string.format("%d%%", progress) or "✗")
                    
                    if config.showLearnedOnly and not learned and progress == 0 then
                        -- Skip unlearned spells if filter is on
                    else
                        output = output .. string.format("  [%s] %s\n", status, spell.name)
                    end
                end
            end
            output = output .. "\n"
        end
    else
        -- Show all spells in list format
        local totalSpells = #spellDB
        local learnedCount = 0
        
        for _, spell in ipairs(spellDB) do
            if hasLearnedSpell(char, spell.id) then
                learnedCount = learnedCount + 1
            end
        end
        
        output = output .. string.format("Overall Progress: %d/%d (%.1f%%)\n\n",
            learnedCount, totalSpells, (learnedCount / totalSpells * 100))
        
        for _, spell in ipairs(spellDB) do
            local learned = hasLearnedSpell(char, spell.id)
            local progress = getSpellProgress(char, spell.id)
            local status = learned and "✓" or (progress > 0 and string.format("%d%%", progress) or "✗")
            
            if config.showLearnedOnly and not learned and progress == 0 then
                -- Skip
            else
                output = output .. string.format("[%s] %s (%s)\n", status, spell.name, spell.school)
            end
        end
    end
    
    API:ShowInfo(output)
    return true
end

-- Recommend espers for spell learning
local function recommendEspers(charIndex)
    local char = safeCall("GetCharacter", charIndex)
    if not char then
        API:ShowError("Character #" .. charIndex .. " not found")
        return false
    end
    
    -- Find spells character hasn't learned yet
    local unlearnedSpells = {}
    for _, spell in ipairs(spellDB) do
        if not hasLearnedSpell(char, spell.id) then
            local progress = getSpellProgress(char, spell.id)
            table.insert(unlearnedSpells, {
                spell = spell,
                progress = progress
            })
        end
    end
    
    if #unlearnedSpells == 0 then
        API:ShowInfo(string.format("%s has learned all spells!", char.name or "Character"))
        return true
    end
    
    -- Sort by progress (highest first) if configured
    if config.sortByProgress then
        table.sort(unlearnedSpells, function(a, b)
            return a.progress > b.progress
        end)
    end
    
    local output = "=== ESPER RECOMMENDATIONS ===\n\n"
    output = output .. string.format("%s - Unlearned Spells: %d\n\n", 
        char.name or "Character", #unlearnedSpells)
    
    -- Show top 10 unlearned spells and which espers teach them
    local count = math.min(10, #unlearnedSpells)
    output = output .. string.format("Top %d Priority Spells:\n\n", count)
    
    for i = 1, count do
        local entry = unlearnedSpells[i]
        local spell = entry.spell
        
        output = output .. string.format("%d. %s", i, spell.name)
        if entry.progress > 0 then
            output = output .. string.format(" (%d%% learned)", entry.progress)
        end
        output = output .. "\n"
        
        -- Find espers that teach this spell
        local espers = findEspersThatTeachSpell(spell.id)
        if #espers > 0 then
            output = output .. "   Taught by: "
            for j, esperInfo in ipairs(espers) do
                output = output .. esperInfo.esper.name
                if config.calculateAPRequired then
                    local apRequired = calculateAPRequired(esperInfo.learnRate, entry.progress)
                    output = output .. string.format(" (%d AP)", apRequired)
                end
                if j < #espers then
                    output = output .. ", "
                end
            end
            output = output .. "\n"
        else
            output = output .. "   No espers teach this spell\n"
        end
        output = output .. "\n"
    end
    
    API:ShowInfo(output)
    return true
end

-- Show esper teaching details
local function showEsperDetails()
    local output = "=== ESPER SPELL TEACHING ===\n\n"
    
    for _, esper in ipairs(esperDB) do
        output = output .. esper.name .. ":\n"
        
        for _, spellInfo in ipairs(esper.spells) do
            local spell = getSpellByID(spellInfo.spellID)
            if spell then
                local apRequired = calculateAPRequired(spellInfo.learnRate, 0)
                output = output .. string.format("  %s (Rate: x%d, AP: %d)\n",
                    spell.name, spellInfo.learnRate, apRequired)
            end
        end
        output = output .. "\n"
    end
    
    API:ShowInfo(output)
    return true
end

-- Find spell teaching source
local function findSpellSource()
    local spellName = API:PromptText("Enter spell name to search:", "")
    if not spellName or spellName == "" then
        return false
    end
    
    -- Find spell
    local foundSpell = nil
    for _, spell in ipairs(spellDB) do
        if spell.name:lower():find(spellName:lower()) then
            foundSpell = spell
            break
        end
    end
    
    if not foundSpell then
        API:ShowError("Spell not found: " .. spellName)
        return false
    end
    
    local output = "=== SPELL LEARNING SOURCE ===\n\n"
    output = output .. string.format("Spell: %s\n", foundSpell.name)
    output = output .. string.format("School: %s\n", foundSpell.school)
    output = output .. string.format("Power: %d\n\n", foundSpell.power or 0)
    
    -- Find espers
    local espers = findEspersThatTeachSpell(foundSpell.id)
    
    if #espers > 0 then
        output = output .. "Taught by:\n"
        for _, esperInfo in ipairs(espers) do
            local apRequired = calculateAPRequired(esperInfo.learnRate, 0)
            output = output .. string.format("  %s (x%d rate, %d AP to master)\n",
                esperInfo.esper.name, esperInfo.learnRate, apRequired)
        end
    else
        output = output .. "No espers teach this spell.\n"
        output = output .. "Spell may be learned naturally or from equipment.\n"
    end
    
    API:ShowInfo(output)
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
=== MAGIC LEARNING PLANNER ===

Choose an operation:

1. Show Spell Progress (Character)
2. Recommend Espers for Learning
3. Show All Esper Details
4. Find Spell Teaching Source
5. Exit

Enter choice (1-5):]]
    
    local choice = API:PromptNumber(menu, 1)
    
    if not choice then
        API:ShowInfo("Operation canceled")
        return
    end
    
    if choice == 1 then
        local charIndex = API:PromptNumber("Enter character index (0-15):", 0)
        if charIndex then
            showSpellProgress(charIndex)
        end
    elseif choice == 2 then
        local charIndex = API:PromptNumber("Enter character index (0-15):", 0)
        if charIndex then
            recommendEspers(charIndex)
        end
    elseif choice == 3 then
        showEsperDetails()
    elseif choice == 4 then
        findSpellSource()
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
