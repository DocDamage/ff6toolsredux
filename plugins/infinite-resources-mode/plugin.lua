-- ============================================================================
-- INFINITE RESOURCES MODE PLUGIN
-- ============================================================================
-- Purpose: Remove resource management entirely with infinite items, gil, MP
-- Category: Experimental / Quality of Life
-- Phase: 6, Batch: 4
-- Complexity: Basic-Intermediate
-- Lines of Code: ~480
-- ============================================================================

-- ⚠️ WARNING: EXPERIMENTAL PLUGIN ⚠️
-- This plugin removes all resource management challenge from the game.
-- It is designed for:
--   - Testing and debugging
--   - Casual playthroughs focused on story
--   - Removing resource grind
--   - Sandbox experimentation
-- 
-- NOT recommended for challenge runs or competitive play.
-- Always back up your save file before using this plugin.

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    MAX_GIL = 9999999,
    MAX_ITEM_QUANTITY = 99,
    MAX_MP_BONUS = 9999,          -- Additional MP from equipment
    TOTAL_ITEMS = 255,            -- Total items in FF6
    TOTAL_CHARACTERS = 14,        -- All playable characters
    
    -- Auto-replenish settings
    AUTO_REPLENISH_ENABLED = false,
    REPLENISH_CHECK_INTERVAL = 60, -- seconds
    
    -- Operation logging
    MAX_LOG_ENTRIES = 100,
    LOG_OPERATIONS = true,
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local InfiniteResourcesMode = {
    version = "1.0.0",
    enabled = false,
    operations_log = {},
    backup_state = nil,
    auto_replenish_active = false,
    last_replenish_time = 0,
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

--- Safe API call wrapper with error handling
local function safeCall(func, ...)
    local success, result = pcall(func, ...)
    if not success then
        print("[ERROR] API call failed: " .. tostring(result))
        return nil, result
    end
    return result
end

--- Log operation with timestamp
local function logOperation(operation_type, details)
    if not CONFIG.LOG_OPERATIONS then return end
    
    local entry = {
        timestamp = os.time(),
        operation = operation_type,
        details = details or {},
    }
    
    table.insert(InfiniteResourcesMode.operations_log, 1, entry)
    
    -- Keep log size manageable
    if #InfiniteResourcesMode.operations_log > CONFIG.MAX_LOG_ENTRIES then
        table.remove(InfiniteResourcesMode.operations_log)
    end
end

--- Create backup of current state
local function createBackup()
    local backup = {
        timestamp = os.time(),
        gil = safeCall(API.getGil),
        items = {},
        character_mp = {},
    }
    
    -- Backup all items
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local quantity = safeCall(API.getItemQuantity, itemId)
        if quantity and quantity > 0 then
            backup.items[itemId] = quantity
        end
    end
    
    -- Backup character MP
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local currentMP = safeCall(API.getCharacterCurrentMP, charId)
        local maxMP = safeCall(API.getCharacterMaxMP, charId)
        if currentMP and maxMP then
            backup.character_mp[charId] = {
                current = currentMP,
                max = maxMP,
            }
        end
    end
    
    InfiniteResourcesMode.backup_state = backup
    logOperation("backup_created", { timestamp = backup.timestamp })
    return true
end

--- Display confirmation dialog
local function confirmAction(message)
    -- Simulated confirmation (in real implementation, use UI dialog)
    print("[CONFIRM] " .. message)
    return true
end

-- ============================================================================
-- CORE RESOURCE MODIFICATION FUNCTIONS
-- ============================================================================

--- Max all Gil
local function maxGil()
    local currentGil = safeCall(API.getGil)
    if not currentGil then
        return false, "Failed to get current Gil"
    end
    
    local success = safeCall(API.setGil, CONFIG.MAX_GIL)
    if not success then
        return false, "Failed to set Gil"
    end
    
    logOperation("max_gil", {
        old_value = currentGil,
        new_value = CONFIG.MAX_GIL,
    })
    
    return true
end

--- Max all items to 99 quantity
local function maxAllItems()
    local modified_count = 0
    local errors = {}
    
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local currentQty = safeCall(API.getItemQuantity, itemId)
        
        if currentQty and currentQty < CONFIG.MAX_ITEM_QUANTITY then
            local success = safeCall(API.setItemQuantity, itemId, CONFIG.MAX_ITEM_QUANTITY)
            if success then
                modified_count = modified_count + 1
            else
                table.insert(errors, itemId)
            end
        end
    end
    
    logOperation("max_all_items", {
        modified_count = modified_count,
        errors = #errors > 0 and errors or nil,
    })
    
    return true, modified_count, errors
end

--- Max specific item category
local function maxItemCategory(category_ids)
    local modified_count = 0
    
    for _, itemId in ipairs(category_ids) do
        local success = safeCall(API.setItemQuantity, itemId, CONFIG.MAX_ITEM_QUANTITY)
        if success then
            modified_count = modified_count + 1
        end
    end
    
    return modified_count
end

--- Set infinite MP for all characters
local function setInfiniteMP()
    local modified_count = 0
    
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        -- Set current MP to max
        local maxMP = safeCall(API.getCharacterMaxMP, charId)
        if maxMP then
            safeCall(API.setCharacterCurrentMP, charId, maxMP)
        end
        
        -- Boost max MP via equipment bonus
        local success = safeCall(API.setCharacterMaxMP, charId, CONFIG.MAX_MP_BONUS)
        if success then
            modified_count = modified_count + 1
            -- Set current to new max
            safeCall(API.setCharacterCurrentMP, charId, CONFIG.MAX_MP_BONUS)
        end
    end
    
    logOperation("set_infinite_mp", {
        characters_modified = modified_count,
        mp_value = CONFIG.MAX_MP_BONUS,
    })
    
    return true, modified_count
end

--- Restore MP to max for all characters (part of auto-replenish)
local function restoreAllMP()
    local restored_count = 0
    
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local maxMP = safeCall(API.getCharacterMaxMP, charId)
        if maxMP then
            local success = safeCall(API.setCharacterCurrentMP, charId, maxMP)
            if success then
                restored_count = restored_count + 1
            end
        end
    end
    
    return restored_count
end

-- ============================================================================
-- PRESET CONFIGURATIONS
-- ============================================================================

--- Apply full infinite resources (everything maxed)
function InfiniteResourcesMode.applyFullInfinite()
    if not confirmAction("Apply full infinite resources mode? This will max all items, Gil, and MP.") then
        return false
    end
    
    print("\n[Infinite Resources] Applying full infinite mode...")
    
    -- Create backup first
    createBackup()
    
    -- Apply all modifications
    local success, msg
    
    -- 1. Max Gil
    success, msg = maxGil()
    if not success then
        print("[ERROR] Failed to max Gil: " .. tostring(msg))
    else
        print("✓ Gil maxed to 9,999,999")
    end
    
    -- 2. Max all items
    success, modified, errors = maxAllItems()
    if success then
        print(string.format("✓ %d items maxed to 99", modified))
        if #errors > 0 then
            print(string.format("  ⚠ %d items failed", #errors))
        end
    end
    
    -- 3. Infinite MP
    success, modified = setInfiniteMP()
    if success then
        print(string.format("✓ Infinite MP set for %d characters", modified))
    end
    
    InfiniteResourcesMode.enabled = true
    logOperation("full_infinite_applied", { timestamp = os.time() })
    
    print("\n[Infinite Resources] Full infinite mode applied! ✅")
    return true
end

--- Apply selective infinite resources
function InfiniteResourcesMode.applySelectiveInfinite(options)
    if not options then
        print("[ERROR] No options provided for selective infinite mode")
        return false
    end
    
    print("\n[Infinite Resources] Applying selective infinite mode...")
    createBackup()
    
    local applied = {}
    
    -- Gil
    if options.gil then
        local success = maxGil()
        if success then
            table.insert(applied, "Gil")
            print("✓ Gil maxed")
        end
    end
    
    -- Items by category
    if options.consumables then
        -- Item IDs 0-99 are typically consumables
        local count = maxItemCategory({0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
        table.insert(applied, string.format("%d consumables", count))
        print(string.format("✓ %d consumable items maxed", count))
    end
    
    if options.equipment then
        -- Equipment typically starts at item ID 100+
        local equipment_ids = {}
        for i = 100, 200 do
            table.insert(equipment_ids, i)
        end
        local count = maxItemCategory(equipment_ids)
        table.insert(applied, string.format("%d equipment", count))
        print(string.format("✓ %d equipment items maxed", count))
    end
    
    if options.all_items then
        local success, modified = maxAllItems()
        if success then
            table.insert(applied, string.format("%d items", modified))
            print(string.format("✓ %d items maxed", modified))
        end
    end
    
    -- MP
    if options.mp then
        local success, modified = setInfiniteMP()
        if success then
            table.insert(applied, "infinite MP")
            print(string.format("✓ Infinite MP for %d characters", modified))
        end
    end
    
    InfiniteResourcesMode.enabled = true
    logOperation("selective_infinite_applied", {
        options = options,
        applied = applied,
    })
    
    print("\n[Infinite Resources] Selective mode applied: " .. table.concat(applied, ", ") .. " ✅")
    return true
end

-- ============================================================================
-- AUTO-REPLENISH SYSTEM
-- ============================================================================

--- Enable auto-replenish mode
function InfiniteResourcesMode.enableAutoReplenish()
    if InfiniteResourcesMode.auto_replenish_active then
        print("[Infinite Resources] Auto-replenish already active")
        return false
    end
    
    InfiniteResourcesMode.auto_replenish_active = true
    InfiniteResourcesMode.last_replenish_time = os.time()
    
    logOperation("auto_replenish_enabled", { timestamp = os.time() })
    print("[Infinite Resources] Auto-replenish enabled ✅")
    return true
end

--- Disable auto-replenish mode
function InfiniteResourcesMode.disableAutoReplenish()
    InfiniteResourcesMode.auto_replenish_active = false
    logOperation("auto_replenish_disabled", { timestamp = os.time() })
    print("[Infinite Resources] Auto-replenish disabled")
    return true
end

--- Check and replenish resources (called periodically)
function InfiniteResourcesMode.checkReplenish()
    if not InfiniteResourcesMode.auto_replenish_active then
        return
    end
    
    local current_time = os.time()
    if current_time - InfiniteResourcesMode.last_replenish_time < CONFIG.REPLENISH_CHECK_INTERVAL then
        return
    end
    
    -- Replenish items that dropped below 99
    local replenished = 0
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local qty = safeCall(API.getItemQuantity, itemId)
        if qty and qty > 0 and qty < CONFIG.MAX_ITEM_QUANTITY then
            safeCall(API.setItemQuantity, itemId, CONFIG.MAX_ITEM_QUANTITY)
            replenished = replenished + 1
        end
    end
    
    -- Restore MP
    local mp_restored = restoreAllMP()
    
    InfiniteResourcesMode.last_replenish_time = current_time
    
    if replenished > 0 or mp_restored > 0 then
        logOperation("auto_replenish_triggered", {
            items_replenished = replenished,
            mp_restored = mp_restored,
        })
    end
end

-- ============================================================================
-- RESTORATION FUNCTIONS
-- ============================================================================

--- Restore from backup
function InfiniteResourcesMode.restoreBackup()
    if not InfiniteResourcesMode.backup_state then
        print("[ERROR] No backup available to restore")
        return false
    end
    
    if not confirmAction("Restore from backup? This will undo infinite resources modifications.") then
        return false
    end
    
    print("\n[Infinite Resources] Restoring from backup...")
    
    local backup = InfiniteResourcesMode.backup_state
    
    -- Restore Gil
    if backup.gil then
        safeCall(API.setGil, backup.gil)
        print("✓ Gil restored")
    end
    
    -- Restore items
    local restored_items = 0
    for itemId, quantity in pairs(backup.items) do
        safeCall(API.setItemQuantity, itemId, quantity)
        restored_items = restored_items + 1
    end
    print(string.format("✓ %d items restored", restored_items))
    
    -- Restore MP
    local restored_chars = 0
    for charId, mp_data in pairs(backup.character_mp) do
        safeCall(API.setCharacterMaxMP, charId, mp_data.max)
        safeCall(API.setCharacterCurrentMP, charId, mp_data.current)
        restored_chars = restored_chars + 1
    end
    print(string.format("✓ MP restored for %d characters", restored_chars))
    
    InfiniteResourcesMode.enabled = false
    logOperation("backup_restored", { timestamp = os.time() })
    
    print("\n[Infinite Resources] Backup restored successfully ✅")
    return true
end

-- ============================================================================
-- STATUS & REPORTING
-- ============================================================================

--- Get current status
function InfiniteResourcesMode.getStatus()
    return {
        enabled = InfiniteResourcesMode.enabled,
        auto_replenish = InfiniteResourcesMode.auto_replenish_active,
        has_backup = InfiniteResourcesMode.backup_state ~= nil,
        operations_count = #InfiniteResourcesMode.operations_log,
    }
end

--- Get operations log
function InfiniteResourcesMode.getLog(count)
    count = count or 10
    local log = {}
    for i = 1, math.min(count, #InfiniteResourcesMode.operations_log) do
        table.insert(log, InfiniteResourcesMode.operations_log[i])
    end
    return log
end

--- Display current resources
function InfiniteResourcesMode.displayResources()
    print("\n=== Current Resources ===")
    
    local gil = safeCall(API.getGil)
    print(string.format("Gil: %s", gil and tostring(gil) or "ERROR"))
    
    local item_count = 0
    local max_items = 0
    for itemId = 0, CONFIG.TOTAL_ITEMS - 1 do
        local qty = safeCall(API.getItemQuantity, itemId)
        if qty and qty > 0 then
            item_count = item_count + 1
            if qty == CONFIG.MAX_ITEM_QUANTITY then
                max_items = max_items + 1
            end
        end
    end
    print(string.format("Items: %d types owned, %d at max (99)", item_count, max_items))
    
    print("\nCharacter MP:")
    for charId = 0, CONFIG.TOTAL_CHARACTERS - 1 do
        local current = safeCall(API.getCharacterCurrentMP, charId)
        local max = safeCall(API.getCharacterMaxMP, charId)
        if current and max then
            print(string.format("  Character %d: %d / %d MP", charId, current, max))
        end
    end
    
    print("========================\n")
end

-- ============================================================================
-- PLUGIN INTERFACE
-- ============================================================================

return InfiniteResourcesMode
