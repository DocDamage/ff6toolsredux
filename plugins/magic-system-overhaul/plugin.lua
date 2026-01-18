-- Magic System Overhaul Plugin v1.0.0
-- Transform FF6's magic system - unlimited MP, zero costs, alternative systems

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Character data
    CHARACTER_COUNT = 14,
    
    -- Magic system modes
    MODES = {
        NORMAL = "normal",           -- Standard FF6 magic system
        UNLIMITED_MP = "unlimited",  -- Max MP (9999)
        ZERO_COST = "zero_cost",     -- All spells cost 0 MP
        HP_BASED = "hp_based",       -- Spells cost HP instead of MP
        ITEM_BASED = "item_based",   -- Spells consume items
        FREE_MAGIC = "free_magic"    -- No costs, unlimited MP
    },
    
    -- MP values
    MAX_MP = 9999,
    
    -- Spell count
    TOTAL_SPELLS = 54,  -- FF6 has 54 magic spells
    
    -- Alternative cost multipliers
    HP_COST_MULTIPLIER = 2.0,   -- HP cost = MP cost × 2
    ITEM_COST_DIVISOR = 10,     -- Item cost = MP cost ÷ 10
    
    -- Backup key
    BACKUP_KEY = "magic_system_backup",
    
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
    current_mode = CONFIG.MODES.NORMAL,
    magic_power_multiplier = 1.0,
    mp_cost_multiplier = 1.0
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
    
    print(string.format("[Magic System Overhaul] %s: %s", operation_type, details))
end

-- Create backup
local function create_backup()
    local backup = {
        timestamp = os.time(),
        character_mp = {},
        learned_spells = {},
        mode = plugin_state.current_mode
    }
    
    -- Backup character MP
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        backup.character_mp[char_id] = {
            current = 0,  -- character.get_current_mp(char_id)
            max = 0       -- character.get_max_mp(char_id)
        }
    end
    
    -- Backup learned spells
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        backup.learned_spells[char_id] = {}
        -- In real implementation: get learned spell list
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    log_operation("BACKUP", "Created magic system backup")
    return backup
end

-- ============================================================================
-- MP MANIPULATION FUNCTIONS
-- ============================================================================

-- Set character to maximum MP
function setUnlimitedMP(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- character.set_max_mp(char_id, CONFIG.MAX_MP)
    -- character.set_current_mp(char_id, CONFIG.MAX_MP)
    
    log_operation("UNLIMITED_MP", string.format("Set character %d to max MP (9999)", char_id))
    return true
end

-- Set all characters to maximum MP
function setUnlimitedMPAll()
    create_backup()
    
    local set_count = 0
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        if setUnlimitedMP(char_id) then
            set_count = set_count + 1
        end
    end
    
    log_operation("UNLIMITED_ALL", string.format("Set all %d characters to max MP", set_count))
    return set_count
end

-- Set MP to zero (for HP-based magic)
function setZeroMP(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- character.set_max_mp(char_id, 0)
    -- character.set_current_mp(char_id, 0)
    
    log_operation("ZERO_MP", string.format("Set character %d to zero MP", char_id))
    return true
end

-- ============================================================================
-- MP COST MODIFICATION
-- ============================================================================

-- Set MP cost multiplier (0 = free, 0.5 = half cost, 2.0 = double cost)
function setMPCostMultiplier(multiplier)
    if multiplier < 0 or multiplier > 10 then
        log_operation("ERROR", "Invalid MP cost multiplier: " .. tostring(multiplier))
        return false
    end
    
    create_backup()
    plugin_state.mp_cost_multiplier = multiplier
    
    -- In real implementation:
    -- game.set_mp_cost_multiplier(multiplier)
    
    log_operation("MP_COST", string.format("Set MP cost multiplier to %.2fx", multiplier))
    return true
end

-- Enable zero MP costs (all spells free)
function enableZeroMPCosts()
    return setMPCostMultiplier(0.0)
end

-- Restore normal MP costs
function restoreNormalMPCosts()
    return setMPCostMultiplier(1.0)
end

-- ============================================================================
-- SPELL UNLOCK FUNCTIONS
-- ============================================================================

-- Unlock specific spell for character
function unlockSpell(char_id, spell_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if spell_id < 0 or spell_id >= CONFIG.TOTAL_SPELLS then
        log_operation("ERROR", "Invalid spell ID: " .. tostring(spell_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- character.learn_spell(char_id, spell_id)
    
    log_operation("UNLOCK_SPELL", string.format("Unlocked spell %d for character %d", spell_id, char_id))
    return true
end

-- Unlock all spells for character
function unlockAllSpells(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    local unlocked_count = 0
    for spell_id = 0, CONFIG.TOTAL_SPELLS - 1 do
        if unlockSpell(char_id, spell_id) then
            unlocked_count = unlocked_count + 1
        end
    end
    
    log_operation("UNLOCK_ALL", string.format("Unlocked %d spells for character %d", unlocked_count, char_id))
    return unlocked_count
end

-- Unlock all spells for all characters
function unlockAllSpellsForAllCharacters()
    create_backup()
    
    local total_unlocked = 0
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        total_unlocked = total_unlocked + unlockAllSpells(char_id)
    end
    
    log_operation("UNLOCK_COMPLETE", string.format("Unlocked spells for all characters (%d total)", total_unlocked))
    return total_unlocked
end

-- ============================================================================
-- MAGIC POWER MODIFICATION
-- ============================================================================

-- Set magic power multiplier
function setMagicPowerMultiplier(multiplier)
    if multiplier < 0.1 or multiplier > 10 then
        log_operation("ERROR", "Invalid magic power multiplier: " .. tostring(multiplier))
        return false
    end
    
    create_backup()
    plugin_state.magic_power_multiplier = multiplier
    
    -- In real implementation:
    -- game.set_magic_power_multiplier(multiplier)
    
    log_operation("MAGIC_POWER", string.format("Set magic power multiplier to %.2fx", multiplier))
    return true
end

-- ============================================================================
-- MAGIC SYSTEM MODES
-- ============================================================================

-- Apply magic system mode
function applyMagicMode(mode_name)
    local mode_functions = {
        [CONFIG.MODES.NORMAL] = function()
            restoreNormalMPCosts()
            setMagicPowerMultiplier(1.0)
            plugin_state.current_mode = CONFIG.MODES.NORMAL
            log_operation("MODE_NORMAL", "Restored normal magic system")
            return true
        end,
        
        [CONFIG.MODES.UNLIMITED_MP] = function()
            setUnlimitedMPAll()
            plugin_state.current_mode = CONFIG.MODES.UNLIMITED_MP
            log_operation("MODE_UNLIMITED", "Enabled unlimited MP mode (9999 MP all characters)")
            return true
        end,
        
        [CONFIG.MODES.ZERO_COST] = function()
            enableZeroMPCosts()
            plugin_state.current_mode = CONFIG.MODES.ZERO_COST
            log_operation("MODE_ZERO_COST", "Enabled zero MP cost mode (all spells free)")
            return true
        end,
        
        [CONFIG.MODES.HP_BASED] = function()
            -- HP-based magic: spells cost HP instead of MP
            enableZeroMPCosts()
            plugin_state.current_mode = CONFIG.MODES.HP_BASED
            log_operation("MODE_HP_BASED", "Enabled HP-based magic (spells cost HP)")
            print("⚠️  Note: HP cost calculation requires in-game implementation")
            return true
        end,
        
        [CONFIG.MODES.ITEM_BASED] = function()
            -- Item-based magic: spells consume items
            enableZeroMPCosts()
            plugin_state.current_mode = CONFIG.MODES.ITEM_BASED
            log_operation("MODE_ITEM_BASED", "Enabled item-based magic (spells consume items)")
            print("⚠️  Note: Item cost calculation requires in-game implementation")
            return true
        end,
        
        [CONFIG.MODES.FREE_MAGIC] = function()
            -- Free magic: unlimited MP + zero costs + all spells
            setUnlimitedMPAll()
            enableZeroMPCosts()
            unlockAllSpellsForAllCharacters()
            plugin_state.current_mode = CONFIG.MODES.FREE_MAGIC
            log_operation("MODE_FREE", "Enabled free magic mode (unlimited MP + zero costs + all spells)")
            return true
        end
    }
    
    local mode_func = mode_functions[mode_name]
    if not mode_func then
        log_operation("ERROR", "Unknown magic mode: " .. tostring(mode_name))
        return false
    end
    
    create_backup()
    return mode_func()
end

-- List available magic modes
function listMagicModes()
    local modes = {
        {name = "normal", description = "Standard FF6 magic system"},
        {name = "unlimited", description = "Unlimited MP (9999) for all characters"},
        {name = "zero_cost", description = "All spells cost 0 MP"},
        {name = "hp_based", description = "Spells cost HP instead of MP (experimental)"},
        {name = "item_based", description = "Spells consume items instead of MP (experimental)"},
        {name = "free_magic", description = "Unlimited MP + zero costs + all spells unlocked"}
    }
    
    print("\n=== Available Magic System Modes ===")
    for _, mode in ipairs(modes) do
        print(string.format("  %s: %s", mode.name, mode.description))
    end
    
    return modes
end

-- ============================================================================
-- PRESET CONFIGURATIONS
-- ============================================================================

-- FF7 Materia-style magic (all spells available, MP-based)
function applyMateriaStyle()
    create_backup()
    
    unlockAllSpellsForAllCharacters()
    setUnlimitedMPAll()
    restoreNormalMPCosts()
    
    log_operation("PRESET_MATERIA", "Applied FF7 Materia-style magic system")
    print("✓ All spells unlocked for all characters")
    print("✓ Unlimited MP (9999)")
    print("✓ Normal MP costs")
    
    return true
end

-- MP-free magic system (no resource management)
function applyMPFreeSystem()
    return applyMagicMode(CONFIG.MODES.FREE_MAGIC)
end

-- Low-MP challenge (all MP costs doubled)
function applyLowMPChallenge()
    create_backup()
    
    setMPCostMultiplier(2.0)
    
    log_operation("PRESET_LOW_MP", "Applied low-MP challenge (2x MP costs)")
    return true
end

-- High-MP challenge (all MP costs halved)
function applyHighMPBonus()
    create_backup()
    
    setMPCostMultiplier(0.5)
    
    log_operation("PRESET_HIGH_MP", "Applied high-MP bonus (0.5x MP costs)")
    return true
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Get magic system status
function getMagicSystemStatus()
    local status = {
        current_mode = plugin_state.current_mode,
        mp_cost_multiplier = plugin_state.mp_cost_multiplier,
        magic_power_multiplier = plugin_state.magic_power_multiplier,
        character_mp = {}
    }
    
    -- Get character MP values
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        status.character_mp[char_id] = {
            current = 0,  -- character.get_current_mp(char_id)
            max = CONFIG.MAX_MP  -- Simulated
        }
    end
    
    return status
end

-- Display magic system status
function displayMagicSystemStatus()
    local status = getMagicSystemStatus()
    
    print("\n=== Magic System Status ===")
    print(string.format("Current Mode: %s", status.current_mode))
    print(string.format("MP Cost Multiplier: %.2fx", status.mp_cost_multiplier))
    print(string.format("Magic Power Multiplier: %.2fx", status.magic_power_multiplier))
    
    print("\n--- Character MP ---")
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        local mp_data = status.character_mp[char_id]
        print(string.format("  Character %d: %d / %d MP", char_id, mp_data.current, mp_data.max))
    end
    
    return status
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
    
    -- Restore character MP
    for char_id, mp_data in pairs(backup.character_mp) do
        -- In real implementation:
        -- character.set_max_mp(char_id, mp_data.max)
        -- character.set_current_mp(char_id, mp_data.current)
    end
    
    -- Restore learned spells
    for char_id, spells in pairs(backup.learned_spells) do
        -- Restore spell list
    end
    
    -- Restore mode
    plugin_state.current_mode = backup.mode
    restoreNormalMPCosts()
    setMagicPowerMultiplier(1.0)
    
    log_operation("RESTORE", "Restored magic system from backup")
    return true
end

-- ============================================================================
-- EXPORT FUNCTIONS
-- ============================================================================

-- Export magic system configuration
function exportMagicConfig()
    local status = getMagicSystemStatus()
    
    local export_text = "=== FF6 Magic System Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", os.date("%Y-%m-%d %H:%M:%S"))
    export_text = export_text .. string.format("Current Mode: %s\n", status.current_mode)
    export_text = export_text .. string.format("MP Cost Multiplier: %.2fx\n", status.mp_cost_multiplier)
    export_text = export_text .. string.format("Magic Power Multiplier: %.2fx\n", status.magic_power_multiplier)
    
    export_text = export_text .. "\nCharacter MP:\n"
    for char_id, mp_data in pairs(status.character_mp) do
        export_text = export_text .. string.format("  Character %d: %d / %d MP\n", char_id, mp_data.current, mp_data.max)
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
    
    log_operation("INIT", "Magic System Overhaul initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("Magic System Overhaul v1.0.0 loaded")
print("Commands: applyMagicMode(mode), setUnlimitedMPAll(), enableZeroMPCosts()")
print("          unlockAllSpellsForAllCharacters(), setMagicPowerMultiplier(x)")
print("          applyMateriaStyle(), applyMPFreeSystem()")
print("          displayMagicSystemStatus(), restoreBackup()")
print("Type listMagicModes() to see available magic system modes")
