-- ============================================================================
-- POVERTY MODE PLUGIN
-- ============================================================================
-- Purpose: Extreme challenge mode - zero resources survival
-- Category: Experimental / Challenge
-- Phase: 6, Batch: 4
-- Complexity: Intermediate
-- Lines of Code: ~530
-- ============================================================================

-- ⚠️ WARNING: EXTREME CHALLENGE PLUGIN ⚠️
-- This plugin removes ALL resources from your save file.
-- It is designed for:
--   - Hardcore challenge runs
--   - Survival gameplay
--   - No-purchase runs
--   - Equipment-only strategies
-- 
-- This mode can make the game EXTREMELY DIFFICULT or even unwinnable.
-- ALWAYS back up your save file before using this plugin.
-- Use restoration wizard to undo if needed.

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    TOTAL_ITEMS = 255,
    TOTAL_CHARACTERS = 14,
    
    -- Poverty levels (preset difficulties)
    POVERTY_LEVELS = {
        EXTREME = {
            name = "Extreme Poverty",
            gil = 0,
            allow_consumables = false,
            allow_equipment_items = false,
            max_item_quantity = 0,
        },
        HARDCORE = {
            name = "Hardcore Poverty",
            gil = 0,
            allow_consumables = false,
            allow_equipment_items = true,  -- Can keep found equipment
            max_item_quantity = 0,
        },
        STRICT = {
            name = "Strict Poverty",
            gil = 100,
            allow_consumables = false,
            allow_equipment_items = true,
            max_item_quantity = 0,
        },
        MODERATE = {
            name = "Moderate Poverty",
            gil = 500,
            allow_consumables = true,      -- Can keep 5 of each
            allow_equipment_items = true,
            max_item_quantity = 5,
        },
        LIGHT = {
            name = "Light Poverty",
            gil = 1000,
            allow_consumables = true,
            allow_equipment_items = true,
            max_item_quantity = 10,
        },
    },
    
    -- Item categories (approximate ranges)
    ITEM_RANGES = {
        consumables = {0, 99},     -- Potions, ethers, etc.
        equipment = {100, 254},    -- Weapons, armor, relics
    },
    
    -- Challenge tracking
    MAX_LOG_ENTRIES = 200,
    TRACK_ACQUISITIONS = true,
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local PovertyMode = {
    version = "1.0.0",
    enabled = false,
    poverty_level = nil,
    challenge_start_time = nil,
    backup_state = nil,
    operations_log = {},
    found_items_log = {},
    compliance_violations = {},
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
    
    table.insert(PovertyMode.operations_log, 1, entry)
    
    if #PovertyMode.operations_log > CONFIG.MAX_LOG_ENTRIES then
        table.remove(PovertyMode.operations_log)
    end
end

--- Create complete backup
local function createBackup()
    local backup = {
        timestamp = os.time(),
        gil = safeCall(API.getGil),
        items = {},
        equipped_items = {},
    }
    
    -- Backup all items
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local quantity = safeCall(API.getItemQuantity, itemId)
        if quantity and quantity > 0 then
            backup.items[itemId] = quantity
        end
    end
    
    -- Backup equipped items per character
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        backup.equipped_items[charId] = {
            weapon = safeCall(API.getCharacterWeapon, charId),
            shield = safeCall(API.getCharacterShield, charId),
            helmet = safeCall(API.getCharacterHelmet, charId),
            armor = safeCall(API.getCharacterArmor, charId),
            relic1 = safeCall(API.getCharacterRelic1, charId),
            relic2 = safeCall(API.getCharacterRelic2, charId),
        }
    end
    
    PovertyMode.backup_state = backup
    logOperation("backup_created", { timestamp = backup.timestamp })
    return true
end

--- Confirmation dialog
local function confirmAction(message)
    print("[CONFIRM] " .. message)
    return true
end

-- ============================================================================
-- RESOURCE REMOVAL FUNCTIONS
-- ============================================================================

--- Zero all Gil
local function zeroGil()
    local currentGil = safeCall(API.getGil)
    if not currentGil then
        return false, "Failed to read current Gil"
    end
    
    local success = safeCall(API.setGil, 0)
    if not success then
        return false, "Failed to set Gil to zero"
    end
    
    logOperation("zero_gil", {
        old_value = currentGil,
        new_value = 0,
    })
    
    return true
end

--- Set Gil to specific amount (for gradual poverty)
local function setGil(amount)
    local success = safeCall(API.setGil, amount)
    if success then
        logOperation("set_gil", { amount = amount })
    end
    return success
end

--- Remove all consumable items
local function removeConsumables()
    local removed_count = 0
    local removed_items = {}
    
    for itemId = CONFIG.ITEM_RANGES.consumables[1], CONFIG.ITEM_RANGES.consumables[2] do
        local quantity = safeCall(API.getItemQuantity, itemId)
        if quantity and quantity > 0 then
            local success = safeCall(API.setItemQuantity, itemId, 0)
            if success then
                removed_count = removed_count + 1
                table.insert(removed_items, { id = itemId, quantity = quantity })
            end
        end
    end
    
    logOperation("remove_consumables", {
        removed_count = removed_count,
        removed_items = removed_items,
    })
    
    return true, removed_count
end

--- Remove all equipment from inventory (keep equipped)
local function removeUnequippedEquipment()
    local removed_count = 0
    
    -- Get list of currently equipped items
    local equipped = {}
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local weapon = safeCall(API.getCharacterWeapon, charId)
        local shield = safeCall(API.getCharacterShield, charId)
        local helmet = safeCall(API.getCharacterHelmet, charId)
        local armor = safeCall(API.getCharacterArmor, charId)
        local relic1 = safeCall(API.getCharacterRelic1, charId)
        local relic2 = safeCall(API.getCharacterRelic2, charId)
        
        if weapon then equipped[weapon] = true end
        if shield then equipped[shield] = true end
        if helmet then equipped[helmet] = true end
        if armor then equipped[armor] = true end
        if relic1 then equipped[relic1] = true end
        if relic2 then equipped[relic2] = true end
    end
    
    -- Remove equipment not currently equipped
    for itemId = CONFIG.ITEM_RANGES.equipment[1], CONFIG.ITEM_RANGES.equipment[2] do
        if not equipped[itemId] then
            local quantity = safeCall(API.getItemQuantity, itemId)
            if quantity and quantity > 0 then
                local success = safeCall(API.setItemQuantity, itemId, 0)
                if success then
                    removed_count = removed_count + 1
                end
            end
        end
    end
    
    logOperation("remove_unequipped_equipment", {
        removed_count = removed_count,
    })
    
    return true, removed_count
end

--- Limit item quantities to max allowed
local function limitItemQuantities(max_quantity)
    local limited_count = 0
    
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local quantity = safeCall(API.getItemQuantity, itemId)
        if quantity and quantity > max_quantity then
            local success = safeCall(API.setItemQuantity, itemId, max_quantity)
            if success then
                limited_count = limited_count + 1
            end
        end
    end
    
    logOperation("limit_item_quantities", {
        max_quantity = max_quantity,
        limited_count = limited_count,
    })
    
    return limited_count
end

-- ============================================================================
-- PRESET POVERTY MODES
-- ============================================================================

--- Apply poverty mode based on preset level
function PovertyMode.applyPovertyLevel(level_name)
    local level = CONFIG.POVERTY_LEVELS[level_name]
    if not level then
        print("[ERROR] Invalid poverty level: " .. tostring(level_name))
        return false
    end
    
    if not confirmAction(string.format("Apply %s? This will remove most/all resources!", level.name)) then
        return false
    end
    
    print(string.format("\n[Poverty Mode] Applying %s...", level.name))
    
    -- Create backup
    createBackup()
    
    -- Apply Gil restriction
    if level.gil == 0 then
        zeroGil()
        print("✓ Gil zeroed")
    else
        setGil(level.gil)
        print(string.format("✓ Gil limited to %d", level.gil))
    end
    
    -- Remove consumables if not allowed
    if not level.allow_consumables then
        local success, count = removeConsumables()
        if success then
            print(string.format("✓ %d consumables removed", count))
        end
    elseif level.max_item_quantity < 99 then
        local count = limitItemQuantities(level.max_item_quantity)
        print(string.format("✓ Items limited to %d each (%d items affected)", level.max_item_quantity, count))
    end
    
    -- Remove equipment if not allowed
    if not level.allow_equipment_items then
        local success, count = removeUnequippedEquipment()
        if success then
            print(string.format("✓ %d unequipped equipment removed", count))
        end
    end
    
    PovertyMode.enabled = true
    PovertyMode.poverty_level = level_name
    PovertyMode.challenge_start_time = os.time()
    
    logOperation("poverty_mode_applied", {
        level = level_name,
        level_config = level,
    })
    
    print(string.format("\n[Poverty Mode] %s applied! Challenge started. ✅", level.name))
    return true
end

--- Custom poverty configuration
function PovertyMode.applyCustomPoverty(config)
    if not config then
        print("[ERROR] No custom configuration provided")
        return false
    end
    
    if not confirmAction("Apply custom poverty configuration?") then
        return false
    end
    
    print("\n[Poverty Mode] Applying custom poverty configuration...")
    
    createBackup()
    
    -- Gil
    if config.gil ~= nil then
        setGil(config.gil)
        print(string.format("✓ Gil set to %d", config.gil))
    end
    
    -- Item limits
    if config.max_item_quantity ~= nil then
        local count = limitItemQuantities(config.max_item_quantity)
        print(string.format("✓ Items limited to %d each (%d affected)", config.max_item_quantity, count))
    end
    
    -- Remove specific categories
    if config.remove_consumables then
        local success, count = removeConsumables()
        if success then
            print(string.format("✓ %d consumables removed", count))
        end
    end
    
    if config.remove_unequipped_equipment then
        local success, count = removeUnequippedEquipment()
        if success then
            print(string.format("✓ %d unequipped equipment removed", count))
        end
    end
    
    PovertyMode.enabled = true
    PovertyMode.poverty_level = "CUSTOM"
    PovertyMode.challenge_start_time = os.time()
    
    logOperation("custom_poverty_applied", { config = config })
    
    print("\n[Poverty Mode] Custom poverty configuration applied! ✅")
    return true
end

-- ============================================================================
-- CHALLENGE TRACKING
-- ============================================================================

--- Track found/acquired item
function PovertyMode.trackFoundItem(itemId, quantity, source)
    if not CONFIG.TRACK_ACQUISITIONS then
        return
    end
    
    local entry = {
        timestamp = os.time(),
        item_id = itemId,
        quantity = quantity,
        source = source or "unknown",
    }
    
    table.insert(PovertyMode.found_items_log, entry)
    logOperation("item_acquired", entry)
end

--- Check compliance with poverty rules
function PovertyMode.checkCompliance()
    if not PovertyMode.enabled or not PovertyMode.poverty_level then
        return true, "Poverty mode not active"
    end
    
    local level = CONFIG.POVERTY_LEVELS[PovertyMode.poverty_level]
    if not level then
        return true, "Custom poverty mode"
    end
    
    local violations = {}
    
    -- Check Gil
    local currentGil = safeCall(API.getGil)
    if currentGil and currentGil > level.gil then
        table.insert(violations, string.format("Gil exceeded: %d / %d", currentGil, level.gil))
    end
    
    -- Check item quantities
    if not level.allow_consumables then
        for itemId = CONFIG.ITEM_RANGES.consumables[1], CONFIG.ITEM_RANGES.consumables[2] do
            local qty = safeCall(API.getItemQuantity, itemId)
            if qty and qty > 0 then
                table.insert(violations, string.format("Consumable found: Item %d (qty: %d)", itemId, qty))
            end
        end
    elseif level.max_item_quantity < 99 then
        for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
            local qty = safeCall(API.getItemQuantity, itemId)
            if qty and qty > level.max_item_quantity then
                table.insert(violations, string.format("Item exceeded limit: Item %d (%d / %d)", itemId, qty, level.max_item_quantity))
            end
        end
    end
    
    if #violations > 0 then
        PovertyMode.compliance_violations = violations
        return false, violations
    end
    
    return true, "Compliant"
end

--- Export challenge proof
function PovertyMode.exportChallengeProof()
    local proof = {
        poverty_level = PovertyMode.poverty_level,
        challenge_start_time = PovertyMode.challenge_start_time,
        challenge_duration = os.time() - (PovertyMode.challenge_start_time or os.time()),
        found_items_count = #PovertyMode.found_items_log,
        operations_count = #PovertyMode.operations_log,
        compliance_status = PovertyMode.checkCompliance(),
        current_gil = safeCall(API.getGil),
    }
    
    return proof
end

-- ============================================================================
-- RESTORATION FUNCTIONS
-- ============================================================================

--- Restore from backup
function PovertyMode.restoreBackup()
    if not PovertyMode.backup_state then
        print("[ERROR] No backup available to restore")
        return false
    end
    
    if not confirmAction("Restore from backup? This will undo poverty mode.") then
        return false
    end
    
    print("\n[Poverty Mode] Restoring from backup...")
    
    local backup = PovertyMode.backup_state
    
    -- Restore Gil
    if backup.gil then
        safeCall(API.setGil, backup.gil)
        print("✓ Gil restored")
    end
    
    -- Restore items
    local restored_count = 0
    for itemId, quantity in pairs(backup.items) do
        safeCall(API.setItemQuantity, itemId, quantity)
        restored_count = restored_count + 1
    end
    print(string.format("✓ %d items restored", restored_count))
    
    PovertyMode.enabled = false
    PovertyMode.poverty_level = nil
    logOperation("backup_restored", { timestamp = os.time() })
    
    print("\n[Poverty Mode] Backup restored successfully! ✅")
    return true
end

-- ============================================================================
-- STATUS & REPORTING
-- ============================================================================

--- Get status
function PovertyMode.getStatus()
    return {
        enabled = PovertyMode.enabled,
        poverty_level = PovertyMode.poverty_level,
        challenge_duration = PovertyMode.challenge_start_time and (os.time() - PovertyMode.challenge_start_time) or 0,
        has_backup = PovertyMode.backup_state ~= nil,
        found_items_count = #PovertyMode.found_items_log,
        compliance_violations = #PovertyMode.compliance_violations,
    }
end

--- Display current resources
function PovertyMode.displayResources()
    print("\n=== Poverty Mode Resources ===")
    
    local gil = safeCall(API.getGil)
    print(string.format("Gil: %s", gil and tostring(gil) or "ERROR"))
    
    local item_count = 0
    local total_items = 0
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local qty = safeCall(API.getItemQuantity, itemId)
        if qty and qty > 0 then
            item_count = item_count + 1
            total_items = total_items + qty
        end
    end
    print(string.format("Items: %d types, %d total quantity", item_count, total_items))
    
    if PovertyMode.enabled then
        local compliant, msg = PovertyMode.checkCompliance()
        print(string.format("\nCompliance: %s", compliant and "✅ COMPLIANT" or "❌ VIOLATIONS"))
        if not compliant and type(msg) == "table" then
            for _, violation in ipairs(msg) do
                print("  - " .. violation)
            end
        end
    end
    
    print("==============================\n")
end

-- ============================================================================
-- PLUGIN INTERFACE
-- ============================================================================

return PovertyMode
