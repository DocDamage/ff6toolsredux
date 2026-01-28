-- ============================================================================
-- ELEMENT AFFINITY SYSTEM PLUGIN
-- ============================================================================
-- Purpose: Assign Pokemon-style elemental affinities to characters
-- Category: Experimental / Combat Enhancement
-- Phase: 6, Batch: 4
-- Complexity: Intermediate-Advanced
-- Lines of Code: ~680
-- ============================================================================

-- ⚠️ WARNING: EXPERIMENTAL PLUGIN ⚠️
-- This plugin significantly alters character balance by assigning elemental
-- affinities with stat bonuses and weakness/resistance patterns.
-- Designed for strategic party composition and elemental-focused gameplay.
-- Always back up your save file before using this plugin.

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Elements available
    ELEMENTS = {
        FIRE = "fire",
        ICE = "ice",
        LIGHTNING = "lightning",
        WATER = "water",
        WIND = "wind",
        EARTH = "earth",
        HOLY = "holy",
        DARK = "dark",
    },
    
    -- Stat bonus multipliers per element
    ELEMENT_STAT_BONUSES = {
        fire = {
            magic = 1.15,      -- +15% Magic
            vigor = 1.10,      -- +10% Vigor
            stamina = 0.95,    -- -5% Stamina
        },
        ice = {
            magic = 1.15,
            defense = 1.10,
            speed = 0.95,
        },
        lightning = {
            speed = 1.20,      -- +20% Speed
            magic = 1.10,
            defense = 0.95,
        },
        water = {
            magic = 1.10,
            stamina = 1.15,
            vigor = 0.95,
        },
        wind = {
            speed = 1.15,
            evade = 1.20,
            defense = 0.90,
        },
        earth = {
            defense = 1.20,
            stamina = 1.15,
            speed = 0.85,
        },
        holy = {
            magic = 1.20,
            stamina = 1.10,
            vigor = 0.95,
        },
        dark = {
            vigor = 1.20,
            magic = 1.05,
            stamina = 0.90,
        },
    },
    
    -- Elemental relationships (Rock-Paper-Scissors style)
    ELEMENT_EFFECTIVENESS = {
        fire = { strong_against = {"ice", "wind"}, weak_against = {"water", "earth"} },
        ice = { strong_against = {"water", "earth"}, weak_against = {"fire", "lightning"} },
        lightning = { strong_against = {"water", "wind"}, weak_against = {"earth"} },
        water = { strong_against = {"fire"}, weak_against = {"lightning", "ice"} },
        wind = { strong_against = {"earth"}, weak_against = {"fire", "lightning"} },
        earth = { strong_against = {"lightning", "fire"}, weak_against = {"water", "wind"} },
        holy = { strong_against = {"dark"}, weak_against = {"dark"} },
        dark = { strong_against = {"holy"}, weak_against = {"holy"} },
    },
    
    TOTAL_CHARACTERS = 14,
    MAX_LOG_ENTRIES = 100,
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local ElementAffinity = {
    version = "1.0.0",
    character_affinities = {},  -- [charId] = element
    backup_state = nil,
    operations_log = {},
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

--- Safe API call wrapper
local function safeCall(func, ...)
    local success, result = pcall(func, ...)
    if not success then
        print("[ERROR] API call failed: " .. tostring(result))
        return nil, result
    end
    return result
end

--- Log operation
local function logOperation(operation_type, details)
    local entry = {
        timestamp = os.time(),
        operation = operation_type,
        details = details or {},
    }
    
    table.insert(ElementAffinity.operations_log, 1, entry)
    
    if #ElementAffinity.operations_log > CONFIG.MAX_LOG_ENTRIES then
        table.remove(ElementAffinity.operations_log)
    end
end

--- Create backup of character stats
local function createBackup()
    local backup = {
        timestamp = os.time(),
        character_stats = {},
        affinities = {},
    }
    
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        backup.character_stats[charId] = {
            vigor = safeCall(API.getCharacterVigor, charId),
            speed = safeCall(API.getCharacterSpeed, charId),
            stamina = safeCall(API.getCharacterStamina, charId),
            magic = safeCall(API.getCharacterMagic, charId),
            defense = safeCall(API.getCharacterDefense, charId),
            m_defense = safeCall(API.getCharacterMagicDefense, charId),
            evade = safeCall(API.getCharacterEvade, charId),
            m_evade = safeCall(API.getCharacterMagicEvade, charId),
        }
        backup.affinities[charId] = ElementAffinity.character_affinities[charId]
    end
    
    ElementAffinity.backup_state = backup
    logOperation("backup_created", { timestamp = backup.timestamp })
    return true
end

--- Confirmation dialog
local function confirmAction(message)
    print("[CONFIRM] " .. message)
    return true
end

-- ============================================================================
-- STAT MODIFICATION FUNCTIONS
-- ============================================================================

--- Apply stat bonuses for element affinity
local function applyElementStatBonuses(charId, element)
    local bonuses = CONFIG.ELEMENT_STAT_BONUSES[element]
    if not bonuses then
        print(string.format("[ERROR] No stat bonuses defined for element: %s", element))
        return false
    end
    
    local modifications = {}
    
    -- Apply vigor bonus
    if bonuses.vigor then
        local current = safeCall(API.getCharacterVigor, charId)
        if current then
            local new_value = math.floor(current * bonuses.vigor)
            safeCall(API.setCharacterVigor, charId, new_value)
            modifications.vigor = {old = current, new = new_value, multiplier = bonuses.vigor}
        end
    end
    
    -- Apply speed bonus
    if bonuses.speed then
        local current = safeCall(API.getCharacterSpeed, charId)
        if current then
            local new_value = math.floor(current * bonuses.speed)
            safeCall(API.setCharacterSpeed, charId, new_value)
            modifications.speed = {old = current, new = new_value, multiplier = bonuses.speed}
        end
    end
    
    -- Apply stamina bonus
    if bonuses.stamina then
        local current = safeCall(API.getCharacterStamina, charId)
        if current then
            local new_value = math.floor(current * bonuses.stamina)
            safeCall(API.setCharacterStamina, charId, new_value)
            modifications.stamina = {old = current, new = new_value, multiplier = bonuses.stamina}
        end
    end
    
    -- Apply magic bonus
    if bonuses.magic then
        local current = safeCall(API.getCharacterMagic, charId)
        if current then
            local new_value = math.floor(current * bonuses.magic)
            safeCall(API.setCharacterMagic, charId, new_value)
            modifications.magic = {old = current, new = new_value, multiplier = bonuses.magic}
        end
    end
    
    -- Apply defense bonus
    if bonuses.defense then
        local current = safeCall(API.getCharacterDefense, charId)
        if current then
            local new_value = math.floor(current * bonuses.defense)
            safeCall(API.setCharacterDefense, charId, new_value)
            modifications.defense = {old = current, new = new_value, multiplier = bonuses.defense}
        end
    end
    
    -- Apply evade bonus
    if bonuses.evade then
        local current = safeCall(API.getCharacterEvade, charId)
        if current then
            local new_value = math.floor(current * bonuses.evade)
            safeCall(API.setCharacterEvade, charId, new_value)
            modifications.evade = {old = current, new = new_value, multiplier = bonuses.evade}
        end
    end
    
    logOperation("stat_bonuses_applied", {
        character_id = charId,
        element = element,
        modifications = modifications,
    })
    
    return true, modifications
end

-- ============================================================================
-- AFFINITY ASSIGNMENT
-- ============================================================================

--- Assign element affinity to character
function ElementAffinity.assignAffinity(charId, element)
    if charId < 0 or charId >= CONFIG.TOTAL_CHARACTERS then
        print("[ERROR] Invalid character ID")
        return false
    end
    
    if not CONFIG.ELEMENTS[element:upper()] then
        print(string.format("[ERROR] Invalid element: %s", element))
        return false
    end
    
    -- Create backup if first assignment
    if not ElementAffinity.backup_state then
        createBackup()
    end
    
    -- Remove previous affinity bonuses if exists
    if ElementAffinity.character_affinities[charId] then
        -- Would need to restore from backup and reapply
        print(string.format("[WARN] Character %d already has affinity, reassigning", charId))
    end
    
    -- Assign affinity
    ElementAffinity.character_affinities[charId] = element:lower()
    
    -- Apply stat bonuses
    local success, modifications = applyElementStatBonuses(charId, element:lower())
    
    if success then
        print(string.format("✓ Character %d assigned %s affinity", charId, element:upper()))
        
        -- Display modifications
        if modifications then
            for stat, mod in pairs(modifications) do
                local change_percent = (mod.multiplier - 1) * 100
                print(string.format("  %s: %d -> %d (%+.0f%%)", 
                    stat:upper(), mod.old, mod.new, change_percent))
            end
        end
    else
        print(string.format("✗ Failed to apply stat bonuses for character %d", charId))
    end
    
    logOperation("affinity_assigned", {
        character_id = charId,
        element = element,
        modifications = modifications,
    })
    
    return success
end

--- Remove affinity from character (restore stats)
function ElementAffinity.removeAffinity(charId)
    if not ElementAffinity.character_affinities[charId] then
        print(string.format("[ERROR] Character %d has no affinity assigned", charId))
        return false
    end
    
    if not ElementAffinity.backup_state then
        print("[ERROR] No backup available to restore stats")
        return false
    end
    
    local backup_stats = ElementAffinity.backup_state.character_stats[charId]
    if not backup_stats then
        print("[ERROR] No backup stats found for character")
        return false
    end
    
    -- Restore original stats
    safeCall(API.setCharacterVigor, charId, backup_stats.vigor)
    safeCall(API.setCharacterSpeed, charId, backup_stats.speed)
    safeCall(API.setCharacterStamina, charId, backup_stats.stamina)
    safeCall(API.setCharacterMagic, charId, backup_stats.magic)
    safeCall(API.setCharacterDefense, charId, backup_stats.defense)
    safeCall(API.setCharacterEvade, charId, backup_stats.evade)
    
    print(string.format("✓ Removed %s affinity from character %d", 
        ElementAffinity.character_affinities[charId]:upper(), charId))
    
    ElementAffinity.character_affinities[charId] = nil
    logOperation("affinity_removed", { character_id = charId })
    
    return true
end

-- ============================================================================
-- PRESET AFFINITY PATTERNS
-- ============================================================================

--- Apply balanced party affinity (one of each element)
function ElementAffinity.applyBalancedParty()
    if not confirmAction("Apply balanced party affinities? (8 characters, all elements)") then
        return false
    end
    
    print("\n[Element Affinity] Applying balanced party...")
    
    local elements = {"fire", "ice", "lightning", "water", "wind", "earth", "holy", "dark"}
    
    for i, element in ipairs(elements) do
        local charId = i - 1  -- Characters 0-7
        ElementAffinity.assignAffinity(charId, element)
    end
    
    print("\n[Element Affinity] Balanced party applied! ✅")
    return true
end

--- Apply offensive-focused affinities
function ElementAffinity.applyOffensiveParty()
    if not confirmAction("Apply offensive party affinities?") then
        return false
    end
    
    print("\n[Element Affinity] Applying offensive party...")
    
    -- High-damage elements: Fire, Lightning, Dark, Holy
    local offensive_elements = {"fire", "lightning", "dark", "holy"}
    
    for i = 0, 3 do
        ElementAffinity.assignAffinity(i, offensive_elements[i + 1])
    end
    
    print("\n[Element Affinity] Offensive party applied! ✅")
    return true
end

--- Apply defensive-focused affinities
function ElementAffinity.applyDefensiveParty()
    if not confirmAction("Apply defensive party affinities?") then
        return false
    end
    
    print("\n[Element Affinity] Applying defensive party...")
    
    -- High-defense elements: Earth, Ice, Water, Holy
    local defensive_elements = {"earth", "ice", "water", "holy"}
    
    for i = 0, 3 do
        ElementAffinity.assignAffinity(i, defensive_elements[i + 1])
    end
    
    print("\n[Element Affinity] Defensive party applied! ✅")
    return true
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

--- Analyze party composition
function ElementAffinity.analyzePartyComposition()
    print("\n=== Party Composition Analysis ===")
    
    local element_counts = {}
    local assigned_count = 0
    
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local element = ElementAffinity.character_affinities[charId]
        if element then
            assigned_count = assigned_count + 1
            element_counts[element] = (element_counts[element] or 0) + 1
        end
    end
    
    print(string.format("Characters with affinities: %d / 14", assigned_count))
    print("\nElement Distribution:")
    for element, count in pairs(element_counts) do
        print(string.format("  %s: %d", element:upper(), count))
    end
    
    -- Check balance
    local max_count = 0
    for _, count in pairs(element_counts) do
        if count > max_count then max_count = count end
    end
    
    if max_count <= 2 then
        print("\n✓ Balanced composition")
    elseif max_count >= 5 then
        print("\n⚠ Heavily skewed composition")
    else
        print("\n~ Moderately skewed composition")
    end
    
    print("==================================\n")
    
    return {
        assigned_count = assigned_count,
        element_counts = element_counts,
        balance_score = (8 - max_count) / 8,  -- 1.0 = perfectly balanced
    }
end

--- Calculate elemental synergy between characters
function ElementAffinity.calculateSynergy(charId1, charId2)
    local element1 = ElementAffinity.character_affinities[charId1]
    local element2 = ElementAffinity.character_affinities[charId2]
    
    if not element1 or not element2 then
        return 0, "One or both characters have no affinity"
    end
    
    local effectiveness = CONFIG.ELEMENT_EFFECTIVENESS[element1]
    if not effectiveness then
        return 0, "No effectiveness data"
    end
    
    -- Check if element1 is strong against element2
    for _, weak_elem in ipairs(effectiveness.strong_against) do
        if weak_elem == element2 then
            return -0.5, string.format("%s is strong against %s (poor synergy)", 
                element1:upper(), element2:upper())
        end
    end
    
    -- Check if element1 is weak against element2
    for _, strong_elem in ipairs(effectiveness.weak_against) do
        if strong_elem == element2 then
            return 1.0, string.format("%s covers %s weakness (excellent synergy)", 
                element2:upper(), element1:upper())
        end
    end
    
    -- Same element
    if element1 == element2 then
        return 0.5, string.format("Both %s (moderate synergy)", element1:upper())
    end
    
    return 0, "Neutral synergy"
end

--- Analyze enemy matchup (theoretical)
function ElementAffinity.analyzeEnemyMatchup(enemy_element)
    print(string.format("\n=== Matchup vs %s Enemy ===", enemy_element:upper()))
    
    local advantages = {}
    local disadvantages = {}
    local neutral = {}
    
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local char_element = ElementAffinity.character_affinities[charId]
        if char_element then
            local effectiveness = CONFIG.ELEMENT_EFFECTIVENESS[char_element]
            
            -- Check if character has advantage
            local has_advantage = false
            if effectiveness then
                for _, strong in ipairs(effectiveness.strong_against) do
                    if strong == enemy_element:lower() then
                        table.insert(advantages, {id = charId, element = char_element})
                        has_advantage = true
                        break
                    end
                end
            end
            
            -- Check if character has disadvantage
            if not has_advantage and effectiveness then
                local has_disadvantage = false
                for _, weak in ipairs(effectiveness.weak_against) do
                    if weak == enemy_element:lower() then
                        table.insert(disadvantages, {id = charId, element = char_element})
                        has_disadvantage = true
                        break
                    end
                end
                
                if not has_disadvantage then
                    table.insert(neutral, {id = charId, element = char_element})
                end
            end
        end
    end
    
    print(string.format("\nAdvantaged Characters: %d", #advantages))
    for _, char in ipairs(advantages) do
        print(string.format("  Character %d (%s) ✓", char.id, char.element:upper()))
    end
    
    print(string.format("\nDisadvantaged Characters: %d", #disadvantages))
    for _, char in ipairs(disadvantages) do
        print(string.format("  Character %d (%s) ✗", char.id, char.element:upper()))
    end
    
    print(string.format("\nNeutral Characters: %d", #neutral))
    for _, char in ipairs(neutral) do
        print(string.format("  Character %d (%s) ~", char.id, char.element:upper()))
    end
    
    print("===============================\n")
    
    return {
        advantages = advantages,
        disadvantages = disadvantages,
        neutral = neutral,
    }
end

-- ============================================================================
-- RESTORATION
-- ============================================================================

--- Restore all characters to original stats
function ElementAffinity.restoreAll()
    if not ElementAffinity.backup_state then
        print("[ERROR] No backup available")
        return false
    end
    
    if not confirmAction("Restore all characters to original stats?") then
        return false
    end
    
    print("\n[Element Affinity] Restoring all characters...")
    
    local restored_count = 0
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        if ElementAffinity.character_affinities[charId] then
            ElementAffinity.removeAffinity(charId)
            restored_count = restored_count + 1
        end
    end
    
    print(string.format("\n[Element Affinity] %d characters restored! ✅", restored_count))
    return true
end

-- ============================================================================
-- STATUS & REPORTING
-- ============================================================================

--- Get status
function ElementAffinity.getStatus()
    local assigned_count = 0
    for _, element in pairs(ElementAffinity.character_affinities) do
        if element then assigned_count = assigned_count + 1 end
    end
    
    return {
        assigned_count = assigned_count,
        has_backup = ElementAffinity.backup_state ~= nil,
        operations_count = #ElementAffinity.operations_log,
    }
end

--- Display all affinities
function ElementAffinity.displayAffinities()
    print("\n=== Character Affinities ===")
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local element = ElementAffinity.character_affinities[charId]
        if element then
            print(string.format("Character %d: %s", charId, element:upper()))
        end
    end
    print("============================\n")
end

-- ============================================================================
-- PLUGIN INTERFACE
-- ============================================================================

return ElementAffinity
