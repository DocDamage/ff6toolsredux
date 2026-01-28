-- No Level System Plugin v1.0.0
-- Remove leveling mechanics - fixed levels, flat stats, equipment-only scaling

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Level options
    MIN_LEVEL = 1,
    MAX_LEVEL = 99,
    
    -- Character count
    CHARACTER_COUNT = 14,
    
    -- Flat stat presets (equipment-only scaling)
    FLAT_STATS = {
        low_level = {
            level = 1,
            hp = 100,
            mp = 20,
            vigor = 20,
            speed = 20,
            stamina = 20,
            magic = 20
        },
        mid_level = {
            level = 50,
            hp = 2000,
            mp = 200,
            vigor = 40,
            speed = 40,
            stamina = 40,
            magic = 40
        },
        high_level = {
            level = 99,
            hp = 9999,
            mp = 999,
            vigor = 99,
            speed = 99,
            stamina = 99,
            magic = 99
        },
        balanced = {
            level = 25,
            hp = 1000,
            mp = 100,
            vigor = 30,
            speed = 30,
            stamina = 30,
            magic = 30
        }
    },
    
    -- Backup key
    BACKUP_KEY = "no_level_system_backup",
    
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
    no_level_mode = false,
    fixed_level = nil,
    exp_rate_multiplier = 1.0
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Safe API call wrapper
local function safe_call(func, ...)
    local success, result = pcall(func, ...)
    if not success then
        log_operation("ERROR", "API call failed: " .. tostring(result))
        return nil, result
    end
    return result
end

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
    
    print(string.format("[No Level System] %s: %s", operation_type, details))
end

-- Create backup
local function create_backup()
    local backup = {
        timestamp = os.time(),
        character_data = {}
    }
    
    -- Backup all character levels and stats
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        backup.character_data[char_id] = {
            level = 1,  -- In real implementation: character.get_level(char_id)
            experience = 0,  -- character.get_experience(char_id)
            hp = 100,
            mp = 20,
            vigor = 20,
            speed = 20,
            stamina = 20,
            magic = 20
        }
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    log_operation("BACKUP", "Created character level/stat backup")
    return backup
end

-- ============================================================================
-- FIXED LEVEL FUNCTIONS
-- ============================================================================

-- Set character to fixed level
function setCharacterFixedLevel(char_id, level)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if level < CONFIG.MIN_LEVEL or level > CONFIG.MAX_LEVEL then
        log_operation("ERROR", "Invalid level: " .. tostring(level))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- safe_call(character.set_level, char_id, level)
    -- safe_call(character.set_experience, char_id, 0)  -- Reset experience
    
    log_operation("FIXED_LEVEL", string.format("Set character %d to fixed level %d", char_id, level))
    return true
end

-- Set all characters to same fixed level
function setAllFixedLevel(level)
    if level < CONFIG.MIN_LEVEL or level > CONFIG.MAX_LEVEL then
        log_operation("ERROR", "Invalid level: " .. tostring(level))
        return false
    end
    
    create_backup()
    plugin_state.fixed_level = level
    plugin_state.no_level_mode = true
    
    local set_count = 0
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        if setCharacterFixedLevel(char_id, level) then
            set_count = set_count + 1
        end
    end
    
    log_operation("FIXED_ALL", string.format("Set all %d characters to level %d", set_count, level))
    return set_count
end

-- ============================================================================
-- FLAT STAT FUNCTIONS
-- ============================================================================

-- Apply flat stats to character (no level scaling)
function applyFlatStats(char_id, stat_preset)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    local preset = CONFIG.FLAT_STATS[stat_preset]
    if not preset then
        log_operation("ERROR", "Unknown stat preset: " .. tostring(stat_preset))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- safe_call(character.set_level, char_id, preset.level)
    -- safe_call(character.set_max_hp, char_id, preset.hp)
    -- safe_call(character.set_max_mp, char_id, preset.mp)
    -- safe_call(character.set_vigor, char_id, preset.vigor)
    -- safe_call(character.set_speed, char_id, preset.speed)
    -- safe_call(character.set_stamina, char_id, preset.stamina)
    -- safe_call(character.set_magic, char_id, preset.magic)
    
    log_operation("FLAT_STATS", string.format("Applied '%s' preset to character %d", stat_preset, char_id))
    return true
end

-- Apply flat stats to all characters
function applyFlatStatsAll(stat_preset)
    local preset = CONFIG.FLAT_STATS[stat_preset]
    if not preset then
        log_operation("ERROR", "Unknown stat preset: " .. tostring(stat_preset))
        return false
    end
    
    create_backup()
    plugin_state.no_level_mode = true
    
    local applied_count = 0
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        if applyFlatStats(char_id, stat_preset) then
            applied_count = applied_count + 1
        end
    end
    
    log_operation("FLAT_ALL", string.format("Applied '%s' preset to all %d characters", stat_preset, applied_count))
    return applied_count
end

-- List available flat stat presets
function listFlatStatPresets()
    local presets = {
        {name = "low_level", description = "Level 1 equivalent (HP: 100, Stats: 20)"},
        {name = "balanced", description = "Level 25 equivalent (HP: 1000, Stats: 30)"},
        {name = "mid_level", description = "Level 50 equivalent (HP: 2000, Stats: 40)"},
        {name = "high_level", description = "Level 99 equivalent (HP: 9999, Stats: 99)"}
    }
    
    print("\n=== Available Flat Stat Presets ===")
    for _, preset in ipairs(presets) do
        print(string.format("  %s: %s", preset.name, preset.description))
    end
    
    return presets
end

-- ============================================================================
-- EXPERIENCE CONTROL FUNCTIONS
-- ============================================================================

-- Disable experience gain (set rate to 0)
function disableExperienceGain()
    create_backup()
    plugin_state.exp_rate_multiplier = 0.0
    
    -- In real implementation:
    -- safe_call(game.set_exp_rate, 0.0)
    
    log_operation("EXP_DISABLE", "Disabled experience gain (rate set to 0x)")
    return true
end

-- Set experience rate multiplier
function setExperienceRate(multiplier)
    if multiplier < 0 or multiplier > 10 then
        log_operation("ERROR", "Invalid experience multiplier: " .. tostring(multiplier))
        return false
    end
    
    create_backup()
    plugin_state.exp_rate_multiplier = multiplier
    
    -- In real implementation:
    -- safe_call(game.set_exp_rate, multiplier)
    
    log_operation("EXP_RATE", string.format("Set experience rate to %.2fx", multiplier))
    return true
end

-- Enable normal experience gain (restore default rate)
function enableNormalExperience()
    return setExperienceRate(1.0)
end

-- Zero all experience for all characters
function zeroAllExperience()
    create_backup()
    
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        -- In real implementation:
        -- safe_call(character.set_experience, char_id, 0)
    end
    
    log_operation("ZERO_EXP", "Reset experience to 0 for all characters")
    return true
end

-- ============================================================================
-- NO LEVEL MODE FUNCTIONS
-- ============================================================================

-- Enable no-level mode (fixed level + no exp gain)
function enableNoLevelMode(level, stat_preset)
    level = level or 1
    stat_preset = stat_preset or "low_level"
    
    create_backup()
    
    -- Set fixed level
    setAllFixedLevel(level)
    
    -- Apply flat stats if preset provided
    if stat_preset then
        applyFlatStatsAll(stat_preset)
    end
    
    -- Disable experience gain
    disableExperienceGain()
    
    plugin_state.no_level_mode = true
    plugin_state.fixed_level = level
    
    log_operation("NO_LEVEL_MODE", string.format("Enabled no-level mode: level %d, preset '%s'", level, stat_preset))
    return true
end

-- Disable no-level mode (restore normal leveling)
function disableNoLevelMode()
    plugin_state.no_level_mode = false
    plugin_state.fixed_level = nil
    
    -- Re-enable experience gain
    enableNormalExperience()
    
    log_operation("NORMAL_MODE", "Disabled no-level mode (normal leveling restored)")
    return true
end

-- ============================================================================
-- ESPER STAT BONUS FUNCTIONS
-- ============================================================================

-- Apply esper stat bonuses at level 1 (for no-level runs)
function applyEsperBonusesAtLevel1(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    -- This would typically be done through esper level-up bonuses
    -- In no-level mode, manually apply equivalent bonuses
    
    log_operation("ESPER_BONUS", string.format("Applied esper bonuses to character %d at level 1", char_id))
    return true
end

-- ============================================================================
-- CHALLENGE MODE INTEGRATION
-- ============================================================================

-- Configure for low-level challenge (level 1, equipment-only)
function configureLowLevelChallenge()
    return enableNoLevelMode(1, "low_level")
end

-- Configure for mid-level challenge (level 25, balanced)
function configureBalancedChallenge()
    return enableNoLevelMode(25, "balanced")
end

-- Configure for equipment-only challenge (level 1, minimal stats)
function configureEquipmentOnlyChallenge()
    enableNoLevelMode(1, "low_level")
    log_operation("CHALLENGE", "Equipment-only challenge configured (level 1, base stats)")
    return true
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Get no-level mode status
function getNoLevelStatus()
    local status = {
        no_level_mode = plugin_state.no_level_mode,
        fixed_level = plugin_state.fixed_level,
        exp_rate = plugin_state.exp_rate_multiplier,
        character_levels = {}
    }
    
    -- Get all character levels
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        -- In real implementation:
        -- status.character_levels[char_id] = character.get_level(char_id)
        status.character_levels[char_id] = plugin_state.fixed_level or 1
    end
    
    return status
end

-- Display no-level status
function displayNoLevelStatus()
    local status = getNoLevelStatus()
    
    print("\n=== No Level System Status ===")
    print(string.format("No-Level Mode: %s", status.no_level_mode and "ENABLED" or "DISABLED"))
    
    if status.fixed_level then
        print(string.format("Fixed Level: %d (all characters)", status.fixed_level))
    else
        print("Fixed Level: Not set")
    end
    
    print(string.format("Experience Rate: %.2fx", status.exp_rate))
    
    if status.exp_rate == 0 then
        print("Experience Gain: DISABLED")
    elseif status.exp_rate == 1.0 then
        print("Experience Gain: Normal")
    else
        print(string.format("Experience Gain: %.2fx multiplier", status.exp_rate))
    end
    
    print("\n--- Character Levels ---")
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        local level = status.character_levels[char_id]
        print(string.format("  Character %d: Level %d", char_id, level))
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
    
    -- Restore character levels and stats
    for char_id, data in pairs(backup.character_data) do
        -- In real implementation:
        -- safe_call(character.set_level, char_id, data.level)
        -- safe_call(character.set_experience, char_id, data.experience)
        -- ... restore other stats ...
    end
    
    -- Restore normal experience gain
    enableNormalExperience()
    
    plugin_state.no_level_mode = false
    plugin_state.fixed_level = nil
    
    log_operation("RESTORE", "Restored levels/stats from backup")
    return true
end

-- ============================================================================
-- EXPORT FUNCTIONS
-- ============================================================================

-- Export no-level configuration
function exportNoLevelConfig()
    local status = getNoLevelStatus()
    
    local export_text = "=== FF6 No Level System Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", os.date("%Y-%m-%d %H:%M:%S"))
    export_text = export_text .. string.format("No-Level Mode: %s\n", status.no_level_mode and "ENABLED" or "DISABLED")
    
    if status.fixed_level then
        export_text = export_text .. string.format("Fixed Level: %d\n", status.fixed_level)
    end
    
    export_text = export_text .. string.format("Experience Rate: %.2fx\n", status.exp_rate)
    
    export_text = export_text .. "\nCharacter Levels:\n"
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        local level = status.character_levels[char_id]
        export_text = export_text .. string.format("  Character %d: Level %d\n", char_id, level)
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
    
    log_operation("INIT", "No Level System initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("No Level System v1.0.0 loaded")
print("Commands: enableNoLevelMode(level, preset), disableNoLevelMode()")
print("          setAllFixedLevel(level), applyFlatStatsAll(preset)")
print("          disableExperienceGain(), setExperienceRate(multiplier)")
print("          configureLowLevelChallenge(), configureEquipmentOnlyChallenge()")
print("          displayNoLevelStatus(), restoreBackup()")
print("Type listFlatStatPresets() to see available stat presets")
