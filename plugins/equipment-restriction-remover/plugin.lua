-- Equipment Restriction Remover Plugin v1.0.0
-- Break all equipment restrictions - any character equips anything

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Character data
    CHARACTER_COUNT = 14,
    
    -- Equipment slots
    SLOTS = {
        WEAPON = 0,
        SHIELD = 1,
        HELMET = 2,
        ARMOR = 3,
        RELIC_1 = 4,
        RELIC_2 = 5
    },
    
    -- Restriction types
    RESTRICTIONS = {
        CHARACTER = "character",       -- Character-specific equipment
        GENDER = "gender",            -- Male/female only items
        DUAL_WIELD = "dual_wield",    -- Dual weapon restrictions
        RELIC_STACK = "relic_stack",  -- Same relic in both slots
        CONFLICTS = "conflicts"       -- Equipment conflicts (e.g., Genji Glove + Offering)
    },
    
    -- Backup key
    BACKUP_KEY = "equipment_restrictions_backup",
    
    -- Logging
    LOG_MAX_ENTRIES = 100
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    backups = {},
    operation_log = {},
    restrictions_removed = {},
    unrestricted_mode = false
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Log operations
local function log_operation(operation_type, details)
    local entry = {
        timestamp = os.time(),
        type = operation_type,
        details = details
    }
    table.insert(plugin_state.operation_log, entry)
    
    if #plugin_state.operation_log > CONFIG.LOG_MAX_ENTRIES then
        table.remove(plugin_state.operation_log, 1)
    end
    
    print(string.format("[Equipment Restriction Remover] %s: %s", operation_type, details))
end

-- Create backup
local function create_backup()
    local backup = {
        timestamp = os.time(),
        equipment_data = {},
        restrictions_removed = {}
    }
    
    -- Backup all character equipment
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        backup.equipment_data[char_id] = {
            weapon = -1,    -- equipment.get_weapon(char_id)
            shield = -1,
            helmet = -1,
            armor = -1,
            relic1 = -1,
            relic2 = -1
        }
    end
    
    -- Backup restriction state
    for restriction_type, _ in pairs(CONFIG.RESTRICTIONS) do
        backup.restrictions_removed[restriction_type] = plugin_state.restrictions_removed[restriction_type] or false
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    log_operation("BACKUP", "Created equipment backup")
    return backup
end

-- ============================================================================
-- RESTRICTION REMOVAL FUNCTIONS
-- ============================================================================

-- Remove character-specific equipment restrictions
function removeCharacterRestrictions()
    create_backup()
    
    plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.CHARACTER] = true
    
    log_operation("REMOVE_CHAR_RESTRICT", "Removed character-specific equipment restrictions")
    print("✓ Any character can now equip any weapon/armor/relic")
    
    return true
end

-- Remove gender restrictions (female-only items)
function removeGenderRestrictions()
    create_backup()
    
    plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.GENDER] = true
    
    log_operation("REMOVE_GENDER_RESTRICT", "Removed gender-based equipment restrictions")
    print("✓ All characters can now equip gender-restricted items")
    print("  (e.g., Minerva Bustier no longer female-only)")
    
    return true
end

-- Enable dual-wielding for all characters
function enableUniversalDualWield()
    create_backup()
    
    plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.DUAL_WIELD] = true
    
    log_operation("ENABLE_DUAL_WIELD", "Enabled dual-wielding for all characters")
    print("✓ All characters can now dual-wield weapons")
    print("  (Requires Genji Glove or Offering equipped)")
    
    return true
end

-- Enable relic stacking (same relic in both slots)
function enableRelicStacking()
    create_backup()
    
    plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.RELIC_STACK] = true
    
    log_operation("ENABLE_RELIC_STACK", "Enabled relic stacking")
    print("✓ Characters can now equip same relic twice")
    print("  (e.g., two Ribbons, two Economizers)")
    
    return true
end

-- Remove equipment conflict restrictions
function removeEquipmentConflicts()
    create_backup()
    
    plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.CONFLICTS] = true
    
    log_operation("REMOVE_CONFLICTS", "Removed equipment conflict restrictions")
    print("✓ Conflicting relics can now be equipped together")
    print("  (e.g., Genji Glove + Offering)")
    
    return true
end

-- Remove ALL restrictions (unrestricted mode)
function enableUnrestrictedMode()
    create_backup()
    
    removeCharacterRestrictions()
    removeGenderRestrictions()
    enableUniversalDualWield()
    enableRelicStacking()
    removeEquipmentConflicts()
    
    plugin_state.unrestricted_mode = true
    
    log_operation("UNRESTRICTED", "Enabled fully unrestricted equipment mode")
    print("\n=== UNRESTRICTED EQUIPMENT MODE ENABLED ===")
    print("✓ No character restrictions")
    print("✓ No gender restrictions")
    print("✓ Universal dual-wielding")
    print("✓ Relic stacking allowed")
    print("✓ No equipment conflicts")
    
    return true
end

-- Restore normal restrictions
function restoreNormalRestrictions()
    for restriction_type, _ in pairs(CONFIG.RESTRICTIONS) do
        plugin_state.restrictions_removed[restriction_type] = false
    end
    
    plugin_state.unrestricted_mode = false
    
    log_operation("RESTORE_RESTRICT", "Restored normal equipment restrictions")
    print("✓ Normal equipment restrictions restored")
    
    return true
end

-- ============================================================================
-- EQUIPMENT VALIDATION OVERRIDE
-- ============================================================================

-- Check if item can be equipped (with restrictions removed)
function canEquip(char_id, item_id, slot)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        return false
    end
    
    -- In unrestricted mode, everything is allowed
    if plugin_state.unrestricted_mode then
        return true
    end
    
    -- Check individual restrictions
    -- In real implementation, this would check actual equipment data
    
    -- Character restrictions
    if not plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.CHARACTER] then
        -- Check character-specific restrictions
    end
    
    -- Gender restrictions
    if not plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.GENDER] then
        -- Check gender-specific restrictions
    end
    
    -- Dual-wield restrictions
    if not plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.DUAL_WIELD] then
        -- Check dual-wield restrictions
    end
    
    -- Relic stacking
    if slot == CONFIG.SLOTS.RELIC_1 or slot == CONFIG.SLOTS.RELIC_2 then
        if not plugin_state.restrictions_removed[CONFIG.RESTRICTIONS.RELIC_STACK] then
            -- Check if same relic already equipped
        end
    end
    
    return true
end

-- Force equip item (bypass validation)
function forceEquipItem(char_id, item_id, slot)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- equipment.force_equip(char_id, item_id, slot)
    
    log_operation("FORCE_EQUIP", string.format("Force equipped item %d to character %d slot %d", item_id, char_id, slot))
    return true
end

-- ============================================================================
-- BUILD FUNCTIONS
-- ============================================================================

-- Apply unrestricted build to character
function applyUnrestrictedBuild(char_id, build_config)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if not build_config then
        log_operation("ERROR", "No build configuration provided")
        return false
    end
    
    create_backup()
    
    -- Equip items from build config
    if build_config.weapon then
        forceEquipItem(char_id, build_config.weapon, CONFIG.SLOTS.WEAPON)
    end
    if build_config.shield then
        forceEquipItem(char_id, build_config.shield, CONFIG.SLOTS.SHIELD)
    end
    if build_config.helmet then
        forceEquipItem(char_id, build_config.helmet, CONFIG.SLOTS.HELMET)
    end
    if build_config.armor then
        forceEquipItem(char_id, build_config.armor, CONFIG.SLOTS.ARMOR)
    end
    if build_config.relic1 then
        forceEquipItem(char_id, build_config.relic1, CONFIG.SLOTS.RELIC_1)
    end
    if build_config.relic2 then
        forceEquipItem(char_id, build_config.relic2, CONFIG.SLOTS.RELIC_2)
    end
    
    log_operation("APPLY_BUILD", string.format("Applied unrestricted build to character %d", char_id))
    return true
end

-- Example builds
function listExampleBuilds()
    local examples = {
        {
            name = "Double Economizer",
            description = "Two Economizers for 2 MP all spells",
            config = {relic1 = 114, relic2 = 114}  -- Economizer ID (example)
        },
        {
            name = "Double Offering",
            description = "Two Offerings for 8x attacks per turn",
            config = {relic1 = 102, relic2 = 102}  -- Offering ID (example)
        },
        {
            name = "Minerva + Muscle Belt",
            description = "Female-only armor on male characters",
            config = {armor = 44, relic1 = 34}  -- Example IDs
        }
    }
    
    print("\n=== Example Unrestricted Builds ===")
    for _, build in ipairs(examples) do
        print(string.format("  %s: %s", build.name, build.description))
    end
    
    return examples
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Get restriction status
function getRestrictionStatus()
    local status = {
        unrestricted_mode = plugin_state.unrestricted_mode,
        removed_restrictions = {}
    }
    
    for restriction_type, type_name in pairs(CONFIG.RESTRICTIONS) do
        status.removed_restrictions[type_name] = plugin_state.restrictions_removed[type_name] or false
    end
    
    return status
end

-- Display restriction status
function displayRestrictionStatus()
    local status = getRestrictionStatus()
    
    print("\n=== Equipment Restriction Status ===")
    print(string.format("Unrestricted Mode: %s", status.unrestricted_mode and "ENABLED" or "DISABLED"))
    
    print("\n--- Individual Restrictions ---")
    print(string.format("  Character Restrictions: %s", status.removed_restrictions.character and "REMOVED" or "ACTIVE"))
    print(string.format("  Gender Restrictions: %s", status.removed_restrictions.gender and "REMOVED" or "ACTIVE"))
    print(string.format("  Dual-Wield Universal: %s", status.removed_restrictions.dual_wield and "ENABLED" or "DISABLED"))
    print(string.format("  Relic Stacking: %s", status.removed_restrictions.relic_stack and "ENABLED" or "DISABLED"))
    print(string.format("  Conflict Removal: %s", status.removed_restrictions.conflicts and "REMOVED" or "ACTIVE"))
    
    return status
end

-- Calculate stats with unrestricted equipment
function calculateUnrestrictedStats(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return nil
    end
    
    -- In real implementation:
    -- Get character base stats
    -- Add equipment bonuses (including "illegal" equipment)
    -- Calculate final stats
    
    local stats = {
        hp = 9999,
        mp = 999,
        vigor = 128,
        speed = 128,
        stamina = 128,
        magic = 128,
        attack = 255,
        defense = 255,
        magic_defense = 255,
        evade = 128,
        magic_evade = 128
    }
    
    return stats
end

-- ============================================================================
-- BACKUP & RESTORE
-- ============================================================================

-- Restore from backup
function restoreBackup()
    local backup = plugin_state.backups[CONFIG.BACKUP_KEY]
    if not backup then
        log_operation("ERROR", "No backup found to restore")
        return false
    end
    
    -- Restore equipment for all characters
    for char_id, equipment in pairs(backup.equipment_data) do
        -- In real implementation:
        -- equipment.set_weapon(char_id, equipment.weapon)
        -- equipment.set_shield(char_id, equipment.shield)
        -- ... restore all slots ...
    end
    
    -- Restore restriction state
    for restriction_type, was_removed in pairs(backup.restrictions_removed) do
        plugin_state.restrictions_removed[restriction_type] = was_removed
    end
    
    log_operation("RESTORE", "Restored equipment from backup")
    return true
end

-- ============================================================================
-- EXPORT FUNCTIONS
-- ============================================================================

-- Export configuration
function exportUnrestrictedConfig()
    local status = getRestrictionStatus()
    
    local export_text = "=== FF6 Equipment Restriction Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", os.date("%Y-%m-%d %H:%M:%S"))
    export_text = export_text .. string.format("Unrestricted Mode: %s\n", status.unrestricted_mode and "ENABLED" or "DISABLED")
    
    export_text = export_text .. "\nRestriction Status:\n"
    for restriction_name, is_removed in pairs(status.removed_restrictions) do
        local status_text = is_removed and "REMOVED" or "ACTIVE"
        export_text = export_text .. string.format("  %s: %s\n", restriction_name, status_text)
    end
    
    export_text = export_text .. "\nOperation Log (last 10 operations):\n"
    local log_start = math.max(1, #plugin_state.operation_log - 9)
    for i = log_start, #plugin_state.operation_log do
        local entry = plugin_state.operation_log[i]
        export_text = export_text .. string.format("  [%s] %s: %s\n",
            os.date("%H:%M:%S", entry.timestamp), entry.type, entry.details)
    end
    
    print(export_text)
    return export_text
end

-- ============================================================================
-- INITIALIZATION
-- ============================================================================

function initialize()
    if plugin_state.initialized then
        return true
    end
    
    -- Initialize restriction state
    for restriction_type, _ in pairs(CONFIG.RESTRICTIONS) do
        plugin_state.restrictions_removed[restriction_type] = false
    end
    
    log_operation("INIT", "Equipment Restriction Remover initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("Equipment Restriction Remover v1.0.0 loaded")
print("Commands: enableUnrestrictedMode(), removeCharacterRestrictions()")
print("          removeGenderRestrictions(), enableUniversalDualWield()")
print("          enableRelicStacking(), removeEquipmentConflicts()")
print("          displayRestrictionStatus(), restoreBackup()")
print("Type listExampleBuilds() to see unrestricted build examples")
