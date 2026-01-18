--[[
    Equipment Optimizer Plugin
    Version: 1.0.0
    
    Analyzes available equipment and recommends optimal loadouts for each character
    based on configurable optimization goals (attack, defense, balanced, speed).
    
    Permissions Required:
    - read_save: Read character stats and inventory
    - write_save: Apply equipment recommendations
    - ui_display: Show analysis and recommendations
]]

-- Configuration
local config = {
    optimizationGoal = "balanced",  -- attack, defense, balanced, speed, magic
    showAllOptions = false,         -- Show all equipment options or just top recommendation
    considerRoleMatch = true,       -- Factor in character's primary role
    preserveRelics = true,          -- Don't change relic assignments
    previewChanges = true,          -- Show preview before applying
}

-- Equipment slots
local SLOT_WEAPON = 0
local SLOT_SHIELD = 1
local SLOT_HELM = 2
local SLOT_ARMOR = 3
local SLOT_RELIC1 = 4
local SLOT_RELIC2 = 5

-- ============================================================================
-- EQUIPMENT DATABASE (Simplified - Real implementation would load from game data)
-- ============================================================================

-- Equipment types
local equipmentDB = {
    weapons = {
        {id = 0, name = "Dirk", attack = 15, defense = 0, magic = 0, speed = 0},
        {id = 1, name = "Mythril Knife", attack = 20, defense = 0, magic = 0, speed = 2},
        {id = 2, name = "Assassin's Dagger", attack = 35, defense = 0, magic = 0, speed = 5},
        {id = 3, name = "Man-Eater", attack = 45, defense = 0, magic = 0, speed = 0},
        {id = 4, name = "Swordbreaker", attack = 42, defense = 5, magic = 0, speed = 0},
        {id = 10, name = "Mythril Pike", attack = 28, defense = 0, magic = 0, speed = 0},
        {id = 11, name = "Trident", attack = 38, defense = 0, magic = 0, speed = 0},
        {id = 12, name = "Stardust Rod", attack = 25, defense = 0, magic = 15, speed = 0},
        {id = 13, name = "Holy Rod", attack = 30, defense = 0, magic = 20, speed = 0},
        {id = 20, name = "Excalibur", attack = 99, defense = 0, magic = 0, speed = 0},
        {id = 21, name = "Illumina", attack = 110, defense = 0, magic = 10, speed = 3},
        {id = 22, name = "Ragnarok", attack = 105, defense = 0, magic = 5, speed = 0},
        {id = 23, name = "Atma Weapon", attack = 120, defense = 0, magic = 0, speed = 0},
    },
    shields = {
        {id = 0, name = "Buckler", attack = 0, defense = 10, magic = 0, magicDef = 0},
        {id = 1, name = "Heavy Shield", attack = 0, defense = 20, magic = 0, magicDef = 0},
        {id = 2, name = "Mythril Shield", attack = 0, defense = 25, magic = 0, magicDef = 5},
        {id = 3, name = "Gold Shield", attack = 0, defense = 30, magic = 0, magicDef = 10},
        {id = 4, name = "Aegis Shield", attack = 0, defense = 40, magic = 0, magicDef = 15},
        {id = 5, name = "Crystal Shield", attack = 0, defense = 50, magic = 0, magicDef = 20},
        {id = 6, name = "Force Shield", attack = 0, defense = 45, magic = 5, magicDef = 30},
        {id = 7, name = "Paladin Shield", attack = 0, defense = 60, magic = 0, magicDef = 40},
    },
    helms = {
        {id = 0, name = "Leather Hat", attack = 0, defense = 8, magic = 0, magicDef = 5},
        {id = 1, name = "Plumed Hat", attack = 0, defense = 12, magic = 0, magicDef = 8},
        {id = 2, name = "Beret", attack = 0, defense = 15, magic = 2, magicDef = 10},
        {id = 3, name = "Mythril Helm", attack = 0, defense = 20, magic = 0, magicDef = 5},
        {id = 4, name = "Gold Helm", attack = 0, defense = 25, magic = 0, magicDef = 10},
        {id = 5, name = "Crystal Helm", attack = 0, defense = 35, magic = 0, magicDef = 15},
        {id = 6, name = "Genji Helm", attack = 0, defense = 40, magic = 0, magicDef = 20},
        {id = 7, name = "Red Cap", attack = 0, defense = 18, magic = 5, magicDef = 20},
    },
    armor = {
        {id = 0, name = "Leather Armor", attack = 0, defense = 15, magic = 0, magicDef = 5},
        {id = 1, name = "Buckler Suit", attack = 0, defense = 20, magic = 0, magicDef = 8},
        {id = 2, name = "Mythril Vest", attack = 0, defense = 30, magic = 0, magicDef = 10},
        {id = 3, name = "Mythril Mail", attack = 0, defense = 35, magic = 0, magicDef = 5},
        {id = 4, name = "Gold Armor", attack = 0, defense = 45, magic = 0, magicDef = 15},
        {id = 5, name = "Crystal Mail", attack = 0, defense = 60, magic = 0, magicDef = 20},
        {id = 6, name = "Genji Armor", attack = 0, defense = 70, magic = 0, magicDef = 25},
        {id = 7, name = "Force Armor", attack = 0, defense = 65, magic = 10, magicDef = 35},
        {id = 8, name = "Minerva Bustier", attack = 0, defense = 55, magic = 15, magicDef = 40},
    },
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

-- Get equipment by ID from database
local function getEquipmentByID(equipType, id)
    local db = equipmentDB[equipType]
    if not db then return nil end
    
    for _, equip in ipairs(db) do
        if equip.id == id then
            return equip
        end
    end
    return nil
end

-- Get equipment name
local function getEquipmentName(equipType, id)
    if id == 255 or id == -1 then
        return "Empty"
    end
    
    local equip = getEquipmentByID(equipType, id)
    if equip then
        return equip.name
    end
    return "Unknown #" .. id
end

-- Calculate stat bonus from equipment
local function calculateEquipmentBonus(equipType, id, stat)
    if id == 255 or id == -1 then
        return 0
    end
    
    local equip = getEquipmentByID(equipType, id)
    if not equip then return 0 end
    
    return equip[stat] or 0
end

-- Calculate total equipment score based on optimization goal
local function calculateEquipmentScore(equip, goal, charStats)
    local score = 0
    
    if goal == "attack" then
        score = (equip.attack or 0) * 5 + 
                (equip.defense or 0) * 1 +
                (equip.speed or 0) * 2
    elseif goal == "defense" then
        score = (equip.defense or 0) * 5 +
                (equip.magicDef or 0) * 4 +
                (equip.attack or 0) * 1
    elseif goal == "magic" then
        score = (equip.magic or 0) * 5 +
                (equip.magicDef or 0) * 3 +
                (equip.defense or 0) * 1
    elseif goal == "speed" then
        score = (equip.speed or 0) * 5 +
                (equip.attack or 0) * 2 +
                (equip.defense or 0) * 1
    else  -- balanced
        score = (equip.attack or 0) * 2 +
                (equip.defense or 0) * 2 +
                (equip.magic or 0) * 2 +
                (equip.magicDef or 0) * 2 +
                (equip.speed or 0) * 1
    end
    
    return score
end

-- Find best equipment from available inventory
local function findBestEquipment(equipType, goal, charStats, inventory)
    local db = equipmentDB[equipType]
    if not db then return nil end
    
    local bestEquip = nil
    local bestScore = -1
    
    for _, equip in ipairs(db) do
        -- Check if we have this item in inventory
        local hasItem = false
        if inventory then
            for _, item in ipairs(inventory) do
                if item.id == equip.id and item.quantity > 0 then
                    hasItem = true
                    break
                end
            end
        else
            hasItem = true  -- If no inventory check, assume available
        end
        
        if hasItem then
            local score = calculateEquipmentScore(equip, goal, charStats)
            if score > bestScore then
                bestScore = score
                bestEquip = equip
            end
        end
    end
    
    return bestEquip
end

-- Get character's current equipment
local function getCurrentEquipment(char)
    if not char or not char.equipment then
        return {
            weapon = 255,
            shield = 255,
            helm = 255,
            armor = 255,
            relic1 = 255,
            relic2 = 255
        }
    end
    
    return {
        weapon = char.equipment[SLOT_WEAPON] or 255,
        shield = char.equipment[SLOT_SHIELD] or 255,
        helm = char.equipment[SLOT_HELM] or 255,
        armor = char.equipment[SLOT_ARMOR] or 255,
        relic1 = char.equipment[SLOT_RELIC1] or 255,
        relic2 = char.equipment[SLOT_RELIC2] or 255
    }
end

-- Format character equipment display
local function formatEquipment(char, charIndex)
    local equip = getCurrentEquipment(char)
    local name = char.name or ("Character #" .. charIndex)
    
    local output = string.format("%s (Lv%d)\n", name, char.level or 1)
    output = output .. "Current Equipment:\n"
    output = output .. string.format("  Weapon: %s\n", getEquipmentName("weapons", equip.weapon))
    output = output .. string.format("  Shield: %s\n", getEquipmentName("shields", equip.shield))
    output = output .. string.format("  Helm: %s\n", getEquipmentName("helms", equip.helm))
    output = output .. string.format("  Armor: %s\n", getEquipmentName("armor", equip.armor))
    
    if not config.preserveRelics then
        output = output .. string.format("  Relic 1: %s\n", equip.relic1 == 255 and "Empty" or ("Relic #" .. equip.relic1))
        output = output .. string.format("  Relic 2: %s\n", equip.relic2 == 255 and "Empty" or ("Relic #" .. equip.relic2))
    end
    
    return output
end

-- ============================================================================
-- OPERATION FUNCTIONS
-- ============================================================================

-- Analyze single character equipment
local function analyzeCharacterEquipment(charIndex)
    local char = safeCall("GetCharacter", charIndex)
    if not char then
        API:ShowError("Character #" .. charIndex .. " not found")
        return false
    end
    
    local inventory = safeCall("GetInventory")
    local currentEquip = getCurrentEquipment(char)
    
    -- Find optimal equipment
    local recommendations = {
        weapon = findBestEquipment("weapons", config.optimizationGoal, char, inventory),
        shield = findBestEquipment("shields", config.optimizationGoal, char, inventory),
        helm = findBestEquipment("helms", config.optimizationGoal, char, inventory),
        armor = findBestEquipment("armor", config.optimizationGoal, char, inventory),
    }
    
    -- Build analysis output
    local output = "=== EQUIPMENT ANALYSIS ===\n\n"
    output = output .. formatEquipment(char, charIndex) .. "\n"
    
    output = output .. string.format("Optimization Goal: %s\n\n", config.optimizationGoal:upper())
    output = output .. "Recommended Equipment:\n"
    
    if recommendations.weapon then
        output = output .. string.format("  Weapon: %s (Atk+%d)\n", 
            recommendations.weapon.name, recommendations.weapon.attack or 0)
    end
    if recommendations.shield then
        output = output .. string.format("  Shield: %s (Def+%d, MDef+%d)\n",
            recommendations.shield.name, 
            recommendations.shield.defense or 0,
            recommendations.shield.magicDef or 0)
    end
    if recommendations.helm then
        output = output .. string.format("  Helm: %s (Def+%d, MDef+%d)\n",
            recommendations.helm.name,
            recommendations.helm.defense or 0,
            recommendations.helm.magicDef or 0)
    end
    if recommendations.armor then
        output = output .. string.format("  Armor: %s (Def+%d, MDef+%d)\n",
            recommendations.armor.name,
            recommendations.armor.defense or 0,
            recommendations.armor.magicDef or 0)
    end
    
    -- Calculate stat improvements
    local currentStats = {
        attack = calculateEquipmentBonus("weapons", currentEquip.weapon, "attack"),
        defense = calculateEquipmentBonus("shields", currentEquip.shield, "defense") +
                  calculateEquipmentBonus("helms", currentEquip.helm, "defense") +
                  calculateEquipmentBonus("armor", currentEquip.armor, "defense"),
        magicDef = calculateEquipmentBonus("shields", currentEquip.shield, "magicDef") +
                   calculateEquipmentBonus("helms", currentEquip.helm, "magicDef") +
                   calculateEquipmentBonus("armor", currentEquip.armor, "magicDef"),
    }
    
    local newStats = {
        attack = (recommendations.weapon and recommendations.weapon.attack or 0),
        defense = (recommendations.shield and recommendations.shield.defense or 0) +
                  (recommendations.helm and recommendations.helm.defense or 0) +
                  (recommendations.armor and recommendations.armor.defense or 0),
        magicDef = (recommendations.shield and recommendations.shield.magicDef or 0) +
                   (recommendations.helm and recommendations.helm.magicDef or 0) +
                   (recommendations.armor and recommendations.armor.magicDef or 0),
    }
    
    output = output .. "\nStat Changes:\n"
    output = output .. string.format("  Attack: %d → %d (%+d)\n",
        currentStats.attack, newStats.attack, newStats.attack - currentStats.attack)
    output = output .. string.format("  Defense: %d → %d (%+d)\n",
        currentStats.defense, newStats.defense, newStats.defense - currentStats.defense)
    output = output .. string.format("  Magic Def: %d → %d (%+d)\n",
        currentStats.magicDef, newStats.magicDef, newStats.magicDef - currentStats.magicDef)
    
    API:ShowInfo(output)
    
    -- Ask to apply
    if config.previewChanges then
        local apply = API:ShowConfirm("Apply recommended equipment to " .. (char.name or ("Character #" .. charIndex)) .. "?")
        if not apply then
            API:ShowInfo("Equipment not applied")
            return false
        end
    end
    
    -- Apply equipment
    if not char.equipment then
        char.equipment = {}
    end
    
    if recommendations.weapon then
        char.equipment[SLOT_WEAPON] = recommendations.weapon.id
    end
    if recommendations.shield then
        char.equipment[SLOT_SHIELD] = recommendations.shield.id
    end
    if recommendations.helm then
        char.equipment[SLOT_HELM] = recommendations.helm.id
    end
    if recommendations.armor then
        char.equipment[SLOT_ARMOR] = recommendations.armor.id
    end
    
    local success = safeCall("SetCharacter", charIndex, char)
    if not success then
        API:ShowError("Failed to apply equipment")
        return false
    end
    
    API:ShowInfo("✓ Equipment applied successfully!")
    return true
end

-- Optimize all party members
local function optimizePartyEquipment()
    local party = safeCall("GetParty")
    if not party then
        API:ShowError("Failed to get party")
        return false
    end
    
    local count = 0
    local output = "=== PARTY EQUIPMENT OPTIMIZATION ===\n\n"
    output = output .. string.format("Goal: %s\n\n", config.optimizationGoal:upper())
    
    for i = 0, 3 do
        local charIndex = party[i]
        if charIndex then
            local char = safeCall("GetCharacter", charIndex)
            if char then
                output = output .. string.format("%d. %s - Ready for optimization\n",
                    i + 1, char.name or ("Character #" .. charIndex))
                count = count + 1
            end
        end
    end
    
    if count == 0 then
        API:ShowError("No party members found")
        return false
    end
    
    output = output .. string.format("\nOptimize equipment for %d party members?", count)
    
    API:ShowInfo(output)
    
    local confirm = API:ShowConfirm("Proceed with party equipment optimization?")
    if not confirm then
        API:ShowInfo("Operation canceled")
        return false
    end
    
    -- Optimize each party member
    local successCount = 0
    for i = 0, 3 do
        local charIndex = party[i]
        if charIndex then
            if analyzeCharacterEquipment(charIndex) then
                successCount = successCount + 1
            end
        end
    end
    
    API:ShowInfo(string.format("✓ Optimized equipment for %d/%d party members", successCount, count))
    return true
end

-- Change optimization goal
local function changeOptimizationGoal()
    local menu = [[
Select Optimization Goal:

1. Attack (maximize offensive power)
2. Defense (maximize survivability)
3. Magic (maximize magic power)
4. Speed (maximize turn frequency)
5. Balanced (all-around optimization)
6. Cancel

Current: ]] .. config.optimizationGoal:upper() .. "\n\nEnter choice (1-6):"
    
    local choice = API:PromptNumber(menu, 1)
    if not choice then return false end
    
    if choice == 1 then
        config.optimizationGoal = "attack"
    elseif choice == 2 then
        config.optimizationGoal = "defense"
    elseif choice == 3 then
        config.optimizationGoal = "magic"
    elseif choice == 4 then
        config.optimizationGoal = "speed"
    elseif choice == 5 then
        config.optimizationGoal = "balanced"
    elseif choice == 6 then
        return false
    else
        API:ShowError("Invalid choice")
        return false
    end
    
    API:ShowInfo("Optimization goal set to: " .. config.optimizationGoal:upper())
    return true
end

-- Compare equipment options
local function compareEquipmentOptions()
    local charIndex = API:PromptNumber("Enter character index (0-15):", 0)
    if not charIndex then return false end
    
    local char = safeCall("GetCharacter", charIndex)
    if not char then
        API:ShowError("Character not found")
        return false
    end
    
    local output = "=== EQUIPMENT COMPARISON ===\n\n"
    output = output .. string.format("%s (Lv%d)\n\n", char.name or ("Character #" .. charIndex), char.level or 1)
    
    -- Show top 3 options for each slot
    local inventory = safeCall("GetInventory")
    
    for _, slotInfo in ipairs({
        {type = "weapons", name = "Weapons"},
        {type = "shields", name = "Shields"},
        {type = "helms", name = "Helms"},
        {type = "armor", name = "Armor"}
    }) do
        output = output .. slotInfo.name .. ":\n"
        
        local options = {}
        for _, equip in ipairs(equipmentDB[slotInfo.type]) do
            local score = calculateEquipmentScore(equip, config.optimizationGoal, char)
            table.insert(options, {equip = equip, score = score})
        end
        
        table.sort(options, function(a, b) return a.score > b.score end)
        
        for i = 1, math.min(3, #options) do
            local opt = options[i]
            output = output .. string.format("  %d. %s (Score: %d)\n",
                i, opt.equip.name, opt.score)
        end
        output = output .. "\n"
    end
    
    API:ShowInfo(output)
    return true
end

-- ============================================================================
-- MAIN PLUGIN ENTRY POINT
-- ============================================================================

local function main()
    -- Verify permissions
    if not API.GetCharacter or not API.SetCharacter or not API.GetInventory then
        API:ShowError("Missing required API methods. Please update the editor.")
        return
    end
    
    -- Show main menu
    local menu = [[
=== EQUIPMENT OPTIMIZER ===

Choose an operation:

1. Analyze Single Character
2. Optimize Party Equipment
3. Change Optimization Goal
4. Compare Equipment Options
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
            analyzeCharacterEquipment(charIndex)
        end
    elseif choice == 2 then
        optimizePartyEquipment()
    elseif choice == 3 then
        changeOptimizationGoal()
    elseif choice == 4 then
        compareEquipmentOptions()
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
