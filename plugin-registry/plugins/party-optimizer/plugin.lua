--[[
    Party Optimizer Plugin
    Version: 1.0.0
    
    Analyzes characters and provides party composition recommendations.
    Features character stat comparison, auto-select best party members,
    and party balance analysis.
    
    Permissions Required:
    - read_save: Read character and party data
    - write_save: Modify party composition
    - ui_display: Show analysis results and recommendations
]]

-- Configuration
local config = {
    partySize = 4,              -- Standard party size
    showTopN = 3,               -- Show top N characters per category
    autoBalance = true,         -- Consider party balance in recommendations
    preferNaturalStats = true,  -- Prefer natural stats over equipment bonuses
    minLevel = 1,               -- Minimum level for consideration
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

-- Calculate combat power rating
local function calculateCombatPower(char)
    if not char then return 0 end
    
    local attack = char.attack or 0
    local magic = char.magic or 0
    local defense = char.defense or 0
    local magicDefense = char.magicDefense or 0
    local speed = char.speed or 0
    local stamina = char.stamina or 0
    
    -- Weighted combat power formula
    local power = (attack * 2.0) +
                  (magic * 1.8) +
                  (defense * 1.5) +
                  (magicDefense * 1.5) +
                  (speed * 1.2) +
                  (stamina * 1.0)
    
    return math.floor(power)
end

-- Calculate magic rating
local function calculateMagicRating(char)
    if not char then return 0 end
    
    local magic = char.magic or 0
    local magicDefense = char.magicDefense or 0
    local mp = char.maxMP or 0
    
    -- Weighted magic rating formula
    local rating = (magic * 3.0) +
                   (magicDefense * 2.0) +
                   (mp * 0.5)
    
    return math.floor(rating)
end

-- Calculate physical rating
local function calculatePhysicalRating(char)
    if not char then return 0 end
    
    local attack = char.attack or 0
    local defense = char.defense or 0
    local stamina = char.stamina or 0
    local hp = char.maxHP or 0
    
    -- Weighted physical rating formula
    local rating = (attack * 3.0) +
                   (defense * 2.0) +
                   (stamina * 1.5) +
                   (hp * 0.3)
    
    return math.floor(rating)
end

-- Calculate tank rating
local function calculateTankRating(char)
    if not char then return 0 end
    
    local defense = char.defense or 0
    local magicDefense = char.magicDefense or 0
    local stamina = char.stamina or 0
    local hp = char.maxHP or 0
    
    -- Weighted tank rating formula
    local rating = (defense * 3.5) +
                   (magicDefense * 3.5) +
                   (stamina * 2.0) +
                   (hp * 0.8)
    
    return math.floor(rating)
end

-- Calculate speed rating
local function calculateSpeedRating(char)
    if not char then return 0 end
    
    local speed = char.speed or 0
    local stamina = char.stamina or 0
    
    -- Weighted speed rating formula
    local rating = (speed * 4.0) +
                   (stamina * 1.0)
    
    return math.floor(rating)
end

-- Get character class/role based on stats
local function getCharacterRole(char)
    if not char then return "Unknown" end
    
    local physical = calculatePhysicalRating(char)
    local magic = calculateMagicRating(char)
    local tank = calculateTankRating(char)
    local speed = calculateSpeedRating(char)
    
    local maxStat = math.max(physical, magic, tank, speed)
    
    if maxStat == tank then
        return "Tank"
    elseif maxStat == physical then
        return "Physical DPS"
    elseif maxStat == magic then
        return "Mage"
    elseif maxStat == speed then
        return "Fast Attacker"
    else
        return "Balanced"
    end
end

-- Format character info string
local function formatCharacter(char, index)
    if not char then return "Character #" .. index .. ": Not Available" end
    
    local name = char.name or ("Character #" .. index)
    local level = char.level or 1
    local hp = char.currentHP or 0
    local maxHP = char.maxHP or 0
    local role = getCharacterRole(char)
    local power = calculateCombatPower(char)
    
    return string.format("%s (Lv%d) - %s - Power: %d - HP: %d/%d",
                         name, level, role, power, hp, maxHP)
end

-- Sort characters by rating function
local function sortCharactersByRating(characters, ratingFunc)
    local sorted = {}
    
    for i, char in pairs(characters) do
        if char and (char.level or 0) >= config.minLevel then
            table.insert(sorted, {
                index = i,
                character = char,
                rating = ratingFunc(char)
            })
        end
    end
    
    table.sort(sorted, function(a, b)
        return a.rating > b.rating
    end)
    
    return sorted
end

-- Get party balance score
local function getPartyBalance(party)
    if not party or #party == 0 then return 0 end
    
    local roles = {}
    local totalPower = 0
    
    for _, member in ipairs(party) do
        local role = getCharacterRole(member)
        roles[role] = (roles[role] or 0) + 1
        totalPower = totalPower + calculateCombatPower(member)
    end
    
    -- Calculate role diversity score (0-100)
    local roleCount = 0
    for _ in pairs(roles) do
        roleCount = roleCount + 1
    end
    
    local diversity = (roleCount / 4) * 100
    local avgPower = totalPower / #party
    
    return math.floor((diversity * 0.4) + (avgPower * 0.001))
end

-- ============================================================================
-- OPERATION FUNCTIONS
-- ============================================================================

-- Show character rankings
local function showCharacterRankings()
    -- Get all characters
    local characters = {}
    for i = 0, 15 do  -- Support up to 16 characters
        local char = safeCall("GetCharacter", i)
        if char then
            characters[i] = char
        end
    end
    
    if not next(characters) then
        API:ShowError("No characters found in save file")
        return false
    end
    
    -- Calculate rankings by category
    local byCombatPower = sortCharactersByRating(characters, calculateCombatPower)
    local byMagic = sortCharactersByRating(characters, calculateMagicRating)
    local byPhysical = sortCharactersByRating(characters, calculatePhysicalRating)
    local byTank = sortCharactersByRating(characters, calculateTankRating)
    local bySpeed = sortCharactersByRating(characters, calculateSpeedRating)
    
    -- Build ranking display
    local output = "=== CHARACTER RANKINGS ===\n\n"
    
    output = output .. "TOP OVERALL COMBAT POWER:\n"
    for i = 1, math.min(config.showTopN, #byCombatPower) do
        local entry = byCombatPower[i]
        output = output .. string.format("%d. %s (Power: %d)\n",
                                         i,
                                         formatCharacter(entry.character, entry.index),
                                         entry.rating)
    end
    
    output = output .. "\nTOP MAGIC USERS:\n"
    for i = 1, math.min(config.showTopN, #byMagic) do
        local entry = byMagic[i]
        output = output .. string.format("%d. %s (Magic: %d)\n",
                                         i,
                                         formatCharacter(entry.character, entry.index),
                                         entry.rating)
    end
    
    output = output .. "\nTOP PHYSICAL ATTACKERS:\n"
    for i = 1, math.min(config.showTopN, #byPhysical) do
        local entry = byPhysical[i]
        output = output .. string.format("%d. %s (Physical: %d)\n",
                                         i,
                                         formatCharacter(entry.character, entry.index),
                                         entry.rating)
    end
    
    output = output .. "\nTOP TANKS:\n"
    for i = 1, math.min(config.showTopN, #byTank) do
        local entry = byTank[i]
        output = output .. string.format("%d. %s (Tank: %d)\n",
                                         i,
                                         formatCharacter(entry.character, entry.index),
                                         entry.rating)
    end
    
    output = output .. "\nFASTEST CHARACTERS:\n"
    for i = 1, math.min(config.showTopN, #bySpeed) do
        local entry = bySpeed[i]
        output = output .. string.format("%d. %s (Speed: %d)\n",
                                         i,
                                         formatCharacter(entry.character, entry.index),
                                         entry.rating)
    end
    
    API:ShowInfo(output)
    return true
end

-- Recommend optimal party
local function recommendOptimalParty()
    -- Get all characters
    local characters = {}
    for i = 0, 15 do
        local char = safeCall("GetCharacter", i)
        if char then
            characters[i] = char
        end
    end
    
    if not next(characters) then
        API:ShowError("No characters found in save file")
        return false
    end
    
    -- Sort by combat power
    local byCombatPower = sortCharactersByRating(characters, calculateCombatPower)
    
    if #byCombatPower < config.partySize then
        API:ShowError("Not enough characters available for party size " .. config.partySize)
        return false
    end
    
    local recommended = {}
    
    if config.autoBalance then
        -- Try to balance party composition
        local roles = {
            Tank = 0,
            ["Physical DPS"] = 0,
            Mage = 0,
            ["Fast Attacker"] = 0,
            Balanced = 0
        }
        
        for _, entry in ipairs(byCombatPower) do
            if #recommended < config.partySize then
                local role = getCharacterRole(entry.character)
                
                -- Prefer diverse roles, but take best if already have 2 of same role
                if not roles[role] or roles[role] < 2 then
                    table.insert(recommended, entry)
                    roles[role] = roles[role] + 1
                end
            end
        end
        
        -- Fill remaining slots with best available
        for _, entry in ipairs(byCombatPower) do
            if #recommended < config.partySize then
                local alreadyAdded = false
                for _, rec in ipairs(recommended) do
                    if rec.index == entry.index then
                        alreadyAdded = true
                        break
                    end
                end
                
                if not alreadyAdded then
                    table.insert(recommended, entry)
                end
            end
        end
    else
        -- Just take top N by combat power
        for i = 1, config.partySize do
            table.insert(recommended, byCombatPower[i])
        end
    end
    
    -- Build recommendation display
    local output = "=== RECOMMENDED OPTIMAL PARTY ===\n\n"
    output = output .. string.format("Party Size: %d\n", config.partySize)
    output = output .. string.format("Balance Mode: %s\n\n",
                                     config.autoBalance and "Enabled" or "Disabled")
    
    for i, entry in ipairs(recommended) do
        output = output .. string.format("%d. %s\n",
                                         i,
                                         formatCharacter(entry.character, entry.index))
    end
    
    -- Calculate party stats
    local totalPower = 0
    local avgLevel = 0
    for _, entry in ipairs(recommended) do
        totalPower = totalPower + entry.rating
        avgLevel = avgLevel + (entry.character.level or 0)
    end
    avgLevel = math.floor(avgLevel / #recommended)
    
    output = output .. string.format("\nTotal Combat Power: %d\n", totalPower)
    output = output .. string.format("Average Level: %d\n", avgLevel)
    output = output .. string.format("Balance Score: %d/100\n",
                                     getPartyBalance(
                                         table.unpack(
                                             (function()
                                                 local chars = {}
                                                 for _, e in ipairs(recommended) do
                                                     table.insert(chars, e.character)
                                                 end
                                                 return chars
                                             end)()
                                         )
                                     ))
    
    API:ShowInfo(output)
    
    -- Ask if user wants to apply
    local apply = API:ShowConfirm(
        "Apply this party composition?\n\n" ..
        "This will update the active party in your save file."
    )
    
    if not apply then
        API:ShowInfo("Party composition not applied")
        return false
    end
    
    -- Apply party composition
    local party = safeCall("GetParty")
    if not party then
        API:ShowError("Failed to get current party")
        return false
    end
    
    -- Update party members
    for i, entry in ipairs(recommended) do
        party[i-1] = entry.index  -- Party uses 0-based indexing
    end
    
    local success = safeCall("SetParty", party)
    if not success then
        API:ShowError("Failed to update party composition")
        return false
    end
    
    API:ShowInfo("✓ Party composition updated successfully!")
    return true
end

-- Show current party analysis
local function analyzeCurrentParty()
    local party = safeCall("GetParty")
    if not party then
        API:ShowError("Failed to get current party")
        return false
    end
    
    -- Get party member characters
    local members = {}
    for i = 0, config.partySize - 1 do
        local charIndex = party[i]
        if charIndex then
            local char = safeCall("GetCharacter", charIndex)
            if char then
                table.insert(members, {
                    index = charIndex,
                    character = char
                })
            end
        end
    end
    
    if #members == 0 then
        API:ShowError("No party members found")
        return false
    end
    
    -- Build analysis display
    local output = "=== CURRENT PARTY ANALYSIS ===\n\n"
    output = output .. string.format("Party Members: %d/%d\n\n",
                                     #members, config.partySize)
    
    local totalPower = 0
    local totalHP = 0
    local totalMaxHP = 0
    local avgLevel = 0
    local roles = {}
    
    for i, member in ipairs(members) do
        local char = member.character
        local role = getCharacterRole(char)
        local power = calculateCombatPower(char)
        
        output = output .. string.format("%d. %s\n",
                                         i,
                                         formatCharacter(char, member.index))
        
        totalPower = totalPower + power
        totalHP = totalHP + (char.currentHP or 0)
        totalMaxHP = totalMaxHP + (char.maxHP or 0)
        avgLevel = avgLevel + (char.level or 0)
        roles[role] = (roles[role] or 0) + 1
    end
    
    avgLevel = math.floor(avgLevel / #members)
    
    output = output .. "\n=== PARTY STATISTICS ===\n"
    output = output .. string.format("Total Combat Power: %d\n", totalPower)
    output = output .. string.format("Average Level: %d\n", avgLevel)
    output = output .. string.format("Total HP: %d/%d (%.1f%%)\n",
                                     totalHP, totalMaxHP,
                                     (totalHP / totalMaxHP) * 100)
    
    output = output .. "\nRole Distribution:\n"
    for role, count in pairs(roles) do
        output = output .. string.format("  %s: %d\n", role, count)
    end
    
    -- Calculate balance score
    local partyChars = {}
    for _, m in ipairs(members) do
        table.insert(partyChars, m.character)
    end
    local balanceScore = getPartyBalance(partyChars)
    
    output = output .. string.format("\nBalance Score: %d/100\n", balanceScore)
    
    if balanceScore < 40 then
        output = output .. "⚠️ Party may be unbalanced - consider diversifying roles\n"
    elseif balanceScore < 70 then
        output = output .. "✓ Party composition is decent\n"
    else
        output = output .. "✓✓ Party composition is excellent!\n"
    end
    
    API:ShowInfo(output)
    return true
end

-- Compare two characters
local function compareCharacters()
    local char1Index = API:PromptNumber("Enter first character index (0-15):", 0)
    if not char1Index then return false end
    
    local char2Index = API:PromptNumber("Enter second character index (0-15):", 1)
    if not char2Index then return false end
    
    local char1 = safeCall("GetCharacter", char1Index)
    local char2 = safeCall("GetCharacter", char2Index)
    
    if not char1 or not char2 then
        API:ShowError("One or both characters not found")
        return false
    end
    
    local output = "=== CHARACTER COMPARISON ===\n\n"
    
    output = output .. formatCharacter(char1, char1Index) .. "\n"
    output = output .. "vs\n"
    output = output .. formatCharacter(char2, char2Index) .. "\n\n"
    
    -- Compare stats
    local stats = {
        {name = "Level", key = "level"},
        {name = "Max HP", key = "maxHP"},
        {name = "Max MP", key = "maxMP"},
        {name = "Attack", key = "attack"},
        {name = "Defense", key = "defense"},
        {name = "Magic", key = "magic"},
        {name = "Mag Def", key = "magicDefense"},
        {name = "Speed", key = "speed"},
        {name = "Stamina", key = "stamina"}
    }
    
    output = output .. "STAT COMPARISON:\n"
    for _, stat in ipairs(stats) do
        local val1 = char1[stat.key] or 0
        local val2 = char2[stat.key] or 0
        local diff = val1 - val2
        local winner = ""
        
        if diff > 0 then
            winner = " ← BETTER"
        elseif diff < 0 then
            winner = " → BETTER"
        else
            winner = " (TIED)"
        end
        
        output = output .. string.format("%s: %d vs %d%s\n",
                                         stat.name, val1, val2, winner)
    end
    
    -- Compare ratings
    output = output .. "\nRATING COMPARISON:\n"
    local ratings = {
        {name = "Combat Power", func = calculateCombatPower},
        {name = "Magic Rating", func = calculateMagicRating},
        {name = "Physical Rating", func = calculatePhysicalRating},
        {name = "Tank Rating", func = calculateTankRating},
        {name = "Speed Rating", func = calculateSpeedRating}
    }
    
    for _, rating in ipairs(ratings) do
        local val1 = rating.func(char1)
        local val2 = rating.func(char2)
        local diff = val1 - val2
        local winner = ""
        
        if diff > 0 then
            winner = " ← BETTER"
        elseif diff < 0 then
            winner = " → BETTER"
        else
            winner = " (TIED)"
        end
        
        output = output .. string.format("%s: %d vs %d%s\n",
                                         rating.name, val1, val2, winner)
    end
    
    API:ShowInfo(output)
    return true
end

-- ============================================================================
-- MAIN PLUGIN ENTRY POINT
-- ============================================================================

local function main()
    -- Verify permissions
    if not API.GetCharacter or not API.GetParty then
        API:ShowError("Missing required API methods. Please update the editor.")
        return
    end
    
    -- Show main menu
    local menu = [[
=== PARTY OPTIMIZER ===

Choose an operation:

1. Show Character Rankings
2. Recommend Optimal Party
3. Analyze Current Party
4. Compare Two Characters
5. Exit

Enter choice (1-5):]]
    
    local choice = API:PromptNumber(menu, 1)
    
    if not choice then
        API:ShowInfo("Operation canceled")
        return
    end
    
    if choice == 1 then
        showCharacterRankings()
    elseif choice == 2 then
        recommendOptimalParty()
    elseif choice == 3 then
        analyzeCurrentParty()
    elseif choice == 4 then
        compareCharacters()
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
