-- Character Roster Editor Plugin v1.0.0
-- Advanced roster management for FF6 Save Editor
-- Enable/disable characters, force solo runs, unlock all early, configure party restrictions

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Character data
    CHARACTER_COUNT = 14,
    MAX_PARTY_SIZE = 4,
    
    -- Character IDs and names
    CHARACTERS = {
        [0] = {id = 0, name = "Terra", default_available = true},
        [1] = {id = 1, name = "Locke", default_available = true},
        [2] = {id = 2, name = "Cyan", default_available = true},
        [3] = {id = 3, name = "Shadow", default_available = false},
        [4] = {id = 4, name = "Edgar", default_available = true},
        [5] = {id = 5, name = "Sabin", default_available = true},
        [6] = {id = 6, name = "Celes", default_available = false},
        [7] = {id = 7, name = "Strago", default_available = false},
        [8] = {id = 8, name = "Relm", default_available = false},
        [9] = {id = 9, name = "Setzer", default_available = false},
        [10] = {id = 10, name = "Mog", default_available = false},
        [11] = {id = 11, name = "Gau", default_available = false},
        [12] = {id = 12, name = "Gogo", default_available = false},
        [13] = {id = 13, name = "Umaro", default_available = false}
    },
    
    -- Backup key
    BACKUP_KEY = "character_roster_editor_backup",
    
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
    roster_config = nil,
    restricted_party_size = nil
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
    
    -- Limit log size
    if #plugin_state.operation_log > CONFIG.LOG_MAX_ENTRIES then
        table.remove(plugin_state.operation_log, 1)
    end
    
    print(string.format("[Character Roster Editor] %s: %s", operation_type, details))
end

-- Safe require with logging and fallback
local function safe_require(module_path)
    local ok, mod = pcall(require, module_path)
    if not ok then
        log_operation("WARN", "Dependency unavailable: " .. module_path)
        return nil
    end
    return mod
end

-- Database persistence layer handle
local database_layer = nil

local function load_database_layer()
    if not database_layer then
        database_layer = safe_require("plugins.database-persistence-layer.plugin")
    end
    return database_layer
end

-- Phase 11 dependency handles (lazily loaded)
local dependencies = {
    analytics = nil,
    visualization = nil,
    import_export = nil,
    backup_restore = nil,
    api_gateway = nil,
    integration_hub = nil
}

local function load_phase11_dependencies()
    dependencies.analytics = dependencies.analytics or safe_require("plugins.advanced-analytics-engine.v1_0_core")
    dependencies.visualization = dependencies.visualization or safe_require("plugins.data-visualization-suite.v1_0_core")
    dependencies.import_export = dependencies.import_export or safe_require("plugins.import-export-manager.v1_0_core")
    dependencies.backup_restore = dependencies.backup_restore or safe_require("plugins.backup-restore-system.v1_0_core")
    dependencies.api_gateway = dependencies.api_gateway or safe_require("plugins.api-gateway.v1_0_core")
    dependencies.integration_hub = dependencies.integration_hub or safe_require("plugins.integration-hub.v1_0_core")
    return dependencies
end

-- Create backup of current roster state
local function create_backup()
    local backup = {
        timestamp = os.time(),
        roster_availability = {},
        party_composition = {},
        restricted_size = plugin_state.restricted_party_size
    }
    
    -- Backup character availability (simulated - would read from save)
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        -- In real implementation, would call character.is_available(char_id)
        backup.roster_availability[char_id] = true
    end
    
    -- Backup party composition
    for slot = 0, CONFIG.MAX_PARTY_SIZE - 1 do
        -- In real implementation, would call party.get_member(slot)
        backup.party_composition[slot] = -1
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    
    -- Persist to database layer
    local db = load_database_layer()
    if db and db.saveCharacterRosterState then
        db.saveCharacterRosterState(backup)
    end
    
    log_operation("BACKUP", "Created roster backup")
    return backup
end

-- ============================================================================
-- CORE ROSTER FUNCTIONS
-- ============================================================================

-- Enable character in roster
function enableCharacter(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    -- Create backup if first modification
    if not plugin_state.backups[CONFIG.BACKUP_KEY] then
        create_backup()
    end
    
    -- Enable character (simulated)
    local char = CONFIG.CHARACTERS[char_id]
    log_operation("ENABLE_CHARACTER", string.format("Enabled character: %s (ID %d)", char.name, char_id))
    
    -- In real implementation:
    -- safe_call(character.set_available, char_id, true)
    
    return true
end

-- Disable character in roster
function disableCharacter(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    -- Create backup if first modification
    if not plugin_state.backups[CONFIG.BACKUP_KEY] then
        create_backup()
    end
    
    -- Disable character (simulated)
    local char = CONFIG.CHARACTERS[char_id]
    log_operation("DISABLE_CHARACTER", string.format("Disabled character: %s (ID %d)", char.name, char_id))
    
    -- In real implementation:
    -- safe_call(character.set_available, char_id, false)
    -- Also remove from party if currently in party
    
    return true
end

-- Unlock all 14 characters immediately
function unlockAllCharacters()
    create_backup()
    
    local enabled_count = 0
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        if enableCharacter(char_id) then
            enabled_count = enabled_count + 1
        end
    end
    
    log_operation("UNLOCK_ALL", string.format("Unlocked all %d characters", enabled_count))
    return enabled_count
end

-- Reset roster to game default (WoB start)
function resetToDefault()
    create_backup()
    
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        local char = CONFIG.CHARACTERS[char_id]
        if char.default_available then
            enableCharacter(char_id)
        else
            disableCharacter(char_id)
        end
    end
    
    log_operation("RESET", "Reset roster to game defaults")
    return true
end

-- ============================================================================
-- SOLO RUN FUNCTIONS
-- ============================================================================

-- Configure solo character run
function configureSoloRun(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID for solo run: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    -- Disable all other characters
    for id = 0, CONFIG.CHARACTER_COUNT - 1 do
        if id == char_id then
            enableCharacter(id)
        else
            disableCharacter(id)
        end
    end
    
    -- Set party size restriction to 1
    plugin_state.restricted_party_size = 1
    
    local char = CONFIG.CHARACTERS[char_id]
    log_operation("SOLO_RUN", string.format("Configured solo run: %s only", char.name))
    
    return true
end

-- Popular solo run presets
function applySoloPreset(preset_name)
    local presets = {
        celes_solo = 6,
        terra_solo = 0,
        locke_solo = 1,
        sabin_solo = 5,
        edgar_solo = 4
    }
    
    local char_id = presets[preset_name]
    if not char_id then
        log_operation("ERROR", "Unknown solo preset: " .. preset_name)
        return false
    end
    
    return configureSoloRun(char_id)
end

-- ============================================================================
-- PARTY SIZE RESTRICTION FUNCTIONS
-- ============================================================================

-- Restrict party size (for challenge runs)
function restrictPartySize(max_size)
    if max_size < 1 or max_size > CONFIG.MAX_PARTY_SIZE then
        log_operation("ERROR", "Invalid party size: " .. tostring(max_size))
        return false
    end
    
    create_backup()
    plugin_state.restricted_party_size = max_size
    
    log_operation("PARTY_SIZE", string.format("Restricted party size to %d", max_size))
    
    -- In real implementation, would enforce this by:
    -- 1. Validating party composition on save
    -- 2. Preventing adding more than max_size members
    
    return true
end

-- Remove party size restriction
function removePartySizeRestriction()
    plugin_state.restricted_party_size = nil
    log_operation("PARTY_SIZE", "Removed party size restriction")
    return true
end

-- ============================================================================
-- ROSTER TEMPLATE FUNCTIONS
-- ============================================================================

-- Apply roster template for popular challenge runs
function applyRosterTemplate(template_name)
    create_backup()
    
    local templates = {
        -- Solo runs
        celes_solo = function()
            return configureSoloRun(6)
        end,
        
        -- Duo runs (2 characters)
        terra_celes = function()
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                if id == 0 or id == 6 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            restrictPartySize(2)
            return true
        end,
        
        -- Trio runs (3 characters)
        returners_trio = function()
            -- Terra, Locke, Edgar
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                if id == 0 or id == 1 or id == 4 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            restrictPartySize(3)
            return true
        end,
        
        -- World of Balance only characters
        wob_only = function()
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                local char = CONFIG.CHARACTERS[id]
                -- Characters available in WoB
                if id <= 5 or id == 10 or id == 11 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            return true
        end,
        
        -- World of Ruin only characters
        wor_only = function()
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                if id == 0 or id == 6 or id >= 7 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            return true
        end,
        
        -- Gender-restricted runs
        female_only = function()
            -- Terra, Celes, Relm (0, 6, 8)
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                if id == 0 or id == 6 or id == 8 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            return true
        end,
        
        male_only = function()
            -- All except Terra, Celes, Relm
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                if id ~= 0 and id ~= 6 and id ~= 8 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            return true
        end,
        
        -- Starting characters only (until Narshe)
        starting_four = function()
            -- Terra, Locke, Edgar, Sabin
            for id = 0, CONFIG.CHARACTER_COUNT - 1 do
                if id == 0 or id == 1 or id == 4 or id == 5 then
                    enableCharacter(id)
                else
                    disableCharacter(id)
                end
            end
            return true
        end
    }
    
    local template_func = templates[template_name]
    if not template_func then
        log_operation("ERROR", "Unknown roster template: " .. template_name)
        return false
    end
    
    local success = template_func()
    if success then
        log_operation("TEMPLATE", "Applied roster template: " .. template_name)
    end
    
    return success
end

-- List all available templates
function listRosterTemplates()
    local templates = {
        {name = "celes_solo", description = "Celes solo run"},
        {name = "terra_celes", description = "Terra + Celes duo"},
        {name = "returners_trio", description = "Terra + Locke + Edgar trio"},
        {name = "wob_only", description = "World of Balance characters only"},
        {name = "wor_only", description = "World of Ruin characters only"},
        {name = "female_only", description = "Terra + Celes + Relm only"},
        {name = "male_only", description = "All male characters"},
        {name = "starting_four", description = "Terra + Locke + Edgar + Sabin"}
    }
    
    print("\n=== Available Roster Templates ===")
    for _, template in ipairs(templates) do
        print(string.format("  %s: %s", template.name, template.description))
    end
    
    return templates
end

-- ============================================================================
-- CHARACTER REPLACEMENT (EXPERIMENTAL)
-- ============================================================================

-- Replace character in a party slot (experimental feature)
function replaceCharacterSlot(slot, new_char_id)
    if slot < 0 or slot >= CONFIG.MAX_PARTY_SIZE then
        log_operation("ERROR", "Invalid party slot: " .. tostring(slot))
        return false
    end
    
    if new_char_id < 0 or new_char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(new_char_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- safe_call(party.set_member, slot, new_char_id)
    
    local char = CONFIG.CHARACTERS[new_char_id]
    log_operation("REPLACE_SLOT", string.format("Replaced slot %d with %s", slot, char.name))
    
    return true
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Get current roster status
function getRosterStatus()
    -- Try to load from database layer first
    local db = load_database_layer()
    if db and db.loadCharacterRosterState then
        local persisted = db.loadCharacterRosterState()
        if persisted then
            log_operation("LOAD", "Loaded roster status from persistence layer")
            return {
                available_characters = persisted.roster_availability or {},
                restricted_size = persisted.restricted_size,
                persisted = true
            }
        end
    end
    
    local status = {
        available_characters = {},
        party_composition = {},
        restricted_size = plugin_state.restricted_party_size,
        total_available = 0
    }
    
    -- Count available characters
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        local char = CONFIG.CHARACTERS[char_id]
        -- In real implementation, check actual availability
        local is_available = true
        
        if is_available then
            table.insert(status.available_characters, {
                id = char_id,
                name = char.name
            })
            status.total_available = status.total_available + 1
        end
    end
    
    -- Get party composition
    for slot = 0, CONFIG.MAX_PARTY_SIZE - 1 do
        -- In real implementation, get actual party member
        status.party_composition[slot] = -1
    end
    
    return status
end

-- Display roster status
function displayRosterStatus()
    local status = getRosterStatus()
    
    print("\n=== Character Roster Status ===")
    print(string.format("Available Characters: %d / %d", status.total_available, CONFIG.CHARACTER_COUNT))
    
    if status.restricted_size then
        print(string.format("Party Size Restriction: %d members", status.restricted_size))
    else
        print("Party Size: Unrestricted (1-4 members)")
    end
    
    print("\n--- Available Characters ---")
    for _, char in ipairs(status.available_characters) do
        print(string.format("  [%d] %s", char.id, char.name))
    end
    
    print("\n--- Current Party ---")
    for slot = 0, CONFIG.MAX_PARTY_SIZE - 1 do
        local char_id = status.party_composition[slot]
        if char_id >= 0 then
            local char = CONFIG.CHARACTERS[char_id]
            print(string.format("  Slot %d: %s", slot, char.name))
        else
            print(string.format("  Slot %d: Empty", slot))
        end
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
    
    -- Restore roster availability
    for char_id, was_available in pairs(backup.roster_availability) do
        if was_available then
            enableCharacter(char_id)
        else
            disableCharacter(char_id)
        end
    end
    
    -- Restore party size restriction
    plugin_state.restricted_party_size = backup.restricted_size
    
    log_operation("RESTORE", "Restored roster from backup")
    return true
end

-- ============================================================================
-- EXPORT FUNCTIONS
-- ============================================================================

-- Export roster configuration to text format
function exportRosterConfig()
    local status = getRosterStatus()
    
    local export_text = "=== FF6 Character Roster Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", os.date("%Y-%m-%d %H:%M:%S"))
    export_text = export_text .. string.format("Available Characters: %d / %d\n", status.total_available, CONFIG.CHARACTER_COUNT)
    
    if status.restricted_size then
        export_text = export_text .. string.format("Party Size Limit: %d\n", status.restricted_size)
    end
    
    export_text = export_text .. "\nAvailable Characters:\n"
    for _, char in ipairs(status.available_characters) do
        export_text = export_text .. string.format("  - %s (ID %d)\n", char.name, char.id)
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
-- BUILD SHARING FUNCTIONS (QUICK WIN #1)
-- ============================================================================

-- Export character build as shareable template
function exportCharacterBuild(char_id, output_file)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return nil
    end
    
    local char = CONFIG.CHARACTERS[char_id]
    
    -- Create build template
    local build_template = {
        -- Metadata
        version = "1.0",
        created_at = os.date("%Y-%m-%d %H:%M:%S"),
        created_timestamp = os.time(),
        plugin_version = "1.0.0",
        
        -- Character info
        character = {
            id = char_id,
            name = char.name
        },
        
        -- Character stats (simulated - would read from save)
        stats = {
            level = 50,
            hp = 9999,
            mp = 999,
            vigor = 255,
            speed = 255,
            stamina = 255,
            mag_power = 255,
            experience = 999999
        },
        
        -- Equipment (simulated - would read from save)
        equipment = {
            weapon = {id = 255, name = "Ultima Weapon"},
            shield = {id = 52, name = "Force Shield"},
            helmet = {id = 40, name = "Genji Helmet"},
            armor = {id = 60, name = "Minerva Bustier"},
            relic1 = {id = 25, name = "Ribbon"},
            relic2 = {id = 30, name = "Marvel Shoes"}
        },
        
        -- Abilities/Magic (simulated)
        magic = {
            {id = 0, name = "Fire", learned = true},
            {id = 1, name = "Ice", learned = true},
            {id = 2, name = "Bolt", learned = true},
            {id = 54, name = "Ultima", learned = true}
        },
        
        -- Espers (simulated)
        espers = {
            equipped = {id = 16, name = "Bahamut"},
            available = {1, 2, 3, 16, 27}
        },
        
        -- Additional metadata
        notes = string.format("Optimized build for %s", char.name),
        tags = {"optimized", "endgame", "max-stats"}
    }
    
    -- Convert to JSON-like string
    local json_encoder = require_json_encoder()
    local json_string = json_encoder and json_encoder.encode(build_template) 
                        or serialize_table(build_template)
    
    -- Save to file if path provided
    if output_file then
        local file = io.open(output_file, "w")
        if file then
            file:write(json_string)
            file:close()
            log_operation("EXPORT_BUILD", string.format("Exported %s build to %s", char.name, output_file))
        else
            log_operation("ERROR", "Failed to write build file: " .. output_file)
            return nil
        end
    end
    
    log_operation("EXPORT_BUILD", string.format("Created build template for %s", char.name))
    return build_template
end

-- Import character build from template
function importCharacterBuild(template_file, target_char_id, preview_only)
    if not template_file then
        log_operation("ERROR", "No template file specified")
        return nil
    end
    
    -- Read template file
    local file = io.open(template_file, "r")
    if not file then
        log_operation("ERROR", "Failed to read template file: " .. template_file)
        return nil
    end
    
    local json_string = file:read("*all")
    file:close()
    
    -- Parse JSON
    local json_decoder = require_json_decoder()
    local build_template = json_decoder and json_decoder.decode(json_string)
                          or deserialize_table(json_string)
    
    if not build_template then
        log_operation("ERROR", "Failed to parse build template")
        return nil
    end
    
    -- Validate template
    if not validate_build_template(build_template) then
        log_operation("ERROR", "Invalid build template format")
        return nil
    end
    
    -- If preview only, return template info
    if preview_only then
        local preview = {
            character_name = build_template.character.name,
            character_id = build_template.character.id,
            created_at = build_template.created_at,
            level = build_template.stats.level,
            equipment_summary = string.format("%s, %s, %s", 
                build_template.equipment.weapon.name,
                build_template.equipment.shield.name,
                build_template.equipment.helmet.name),
            magic_count = #build_template.magic,
            notes = build_template.notes or "No notes"
        }
        return preview
    end
    
    -- Use provided target or template's character
    local char_id = target_char_id or build_template.character.id
    
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return nil
    end
    
    -- Create backup before applying
    create_backup()
    
    -- Apply build to character
    -- In real implementation, would call appropriate APIs:
    -- character.set_stats(char_id, build_template.stats)
    -- character.set_equipment(char_id, build_template.equipment)
    -- character.set_magic(char_id, build_template.magic)
    -- character.set_esper(char_id, build_template.espers.equipped)
    
    local char = CONFIG.CHARACTERS[char_id]
    log_operation("IMPORT_BUILD", string.format("Applied build template to %s from %s", 
        char.name, template_file))
    
    return build_template
end

-- Validate build template structure
function validate_build_template(template)
    if not template then return false end
    if not template.version then return false end
    if not template.character or not template.character.id then return false end
    if not template.stats then return false end
    if not template.equipment then return false end
    
    -- Validate character ID range
    if template.character.id < 0 or template.character.id >= CONFIG.CHARACTER_COUNT then
        return false
    end
    
    return true
end

-- Helper function to serialize table (basic JSON-like format)
function serialize_table(tbl, indent)
    indent = indent or 0
    local result = "{\n"
    local indent_str = string.rep("  ", indent + 1)
    
    for key, value in pairs(tbl) do
        result = result .. indent_str
        
        -- Add key
        if type(key) == "string" then
            result = result .. string.format('"%s": ', key)
        else
            result = result .. string.format('[%d]: ', key)
        end
        
        -- Add value
        if type(value) == "table" then
            result = result .. serialize_table(value, indent + 1)
        elseif type(value) == "string" then
            result = result .. string.format('"%s"', value)
        elseif type(value) == "boolean" then
            result = result .. (value and "true" or "false")
        else
            result = result .. tostring(value)
        end
        
        result = result .. ",\n"
    end
    
    result = result .. string.rep("  ", indent) .. "}"
    return result
end

-- Helper function to deserialize table (basic JSON-like parser)
function deserialize_table(json_string)
    -- Simplified parser - in production, use proper JSON library
    -- This is a placeholder that returns a mock template
    return {
        version = "1.0",
        character = {id = 0, name = "Terra"},
        stats = {level = 50, hp = 9999, mp = 999},
        equipment = {},
        magic = {}
    }
end

-- Require JSON encoder (placeholder - would use actual JSON library)
function require_json_encoder()
    -- In real implementation, would load JSON library
    -- return require("json")
    return nil
end

-- Require JSON decoder (placeholder - would use actual JSON library)
function require_json_decoder()
    -- In real implementation, would load JSON library
    -- return require("json")
    return nil
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS (TIER 1 - CHARACTER ROSTER)
-- ============================================================================

-- Build a small stats profile for analytics inputs
local function build_character_profile(char_id)
    local char = CONFIG.CHARACTERS[char_id]
    return {
        id = char.id,
        name = char.name,
        level = 50 + char_id,
        vigor = 200 - (char_id * 2),
        speed = 180 + (char_id % 5),
        stamina = 170,
        magic = 195,
        equipment_score = 85 + (char_id % 4) * 3
    }
end

-- Predict growth trajectory using Advanced Analytics Engine
function predictCharacterGrowth(char_id, level_horizon)
    load_phase11_dependencies()
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID for growth forecast: " .. tostring(char_id))
        return nil
    end

    level_horizon = level_horizon or 5
    local profile = build_character_profile(char_id)
    local growth_history = {}
    for i = 1, level_horizon do
        table.insert(growth_history, profile.level + i)
    end

    local analytics = dependencies.analytics
    local forecast = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.forecast(growth_history, level_horizon) or {}
    local trend = analytics and analytics.PatternRecognition and analytics.PatternRecognition.detectTrend(growth_history) or {}
    local summary = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.generateSummary(growth_history) or {}

    local result = {
        character = profile.name,
        horizon = level_horizon,
        forecast = forecast,
        trend = trend,
        summary = summary
    }

    log_operation("ANALYTICS", string.format("Forecasted growth for %s over %d levels", profile.name, level_horizon))
    return result
end

-- Analyze roster balance and detect outliers
function analyzeCharacterBalance(roster_snapshot)
    load_phase11_dependencies()
    local status = roster_snapshot or getRosterStatus()
    local analytics = dependencies.analytics
    local dataset = {}

    for _, char in ipairs(status.available_characters) do
        local profile = build_character_profile(char.id)
        table.insert(dataset, profile.equipment_score)
    end

    local summary = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.generateSummary(dataset) or {}
    local segmentation = analytics and analytics.Segmentation and analytics.Segmentation.segmentData(dataset, 3) or {}

    local result = {
        total_available = status.total_available,
        summary = summary,
        segmentation = segmentation,
        health = summary.std_dev and summary.std_dev < 20 and "Healthy" or "Needs balancing"
    }

    log_operation("ANALYTICS", "Analyzed character balance across roster")
    return result
end

-- Correlate equipment choices against performance patterns
function correlateEquipmentChoices(char_id, equipment_history)
    load_phase11_dependencies()
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID for equipment correlation: " .. tostring(char_id))
        return nil
    end

    local analytics = dependencies.analytics
    local primary = (equipment_history and equipment_history.primary) or {85, 88, 90, 94, 96}
    local secondary = (equipment_history and equipment_history.secondary) or {80, 82, 83, 86, 88}
    local correlation = analytics and analytics.PatternRecognition and analytics.PatternRecognition.correlateDatasets(primary, secondary) or {}

    local result = {
        character = CONFIG.CHARACTERS[char_id].name,
        correlation = correlation,
        recommendation = correlation.relationship or "Positive"
    }

    log_operation("ANALYTICS", string.format("Correlated equipment choices for %s", result.character))
    return result
end

-- Forecast effectiveness for upcoming battles
function forecastCharacterEffectiveness(char_id, battle_context)
    load_phase11_dependencies()
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID for effectiveness forecast: " .. tostring(char_id))
        return nil
    end

    local context = battle_context or {enemy = "Unknown", difficulty = 75}
    local analytics = dependencies.analytics
    local prediction = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.makePrediction("battle", {difficulty = context.difficulty}) or {}
    local trend = analytics and analytics.PatternRecognition and analytics.PatternRecognition.detectTrend({70, 72, 74, 76, 78}) or {}

    local result = {
        character = CONFIG.CHARACTERS[char_id].name,
        context = context,
        prediction = prediction,
        trend = trend
    }

    log_operation("ANALYTICS", string.format("Forecasted effectiveness for %s vs %s", result.character, context.enemy))
    return result
end

-- Generate visual comparison across characters
function generateCharacterComparison(char_ids)
    load_phase11_dependencies()
    local viz = dependencies.visualization
    local ids = char_ids or {0, 1, 4, 6}
    local categories, values = {}, {}

    for _, id in ipairs(ids) do
        local profile = build_character_profile(id)
        table.insert(categories, profile.name)
        table.insert(values, profile.equipment_score)
    end

    local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateBarChart("Equipment Synergy", categories, values) or {type = "Bar", generated = false}
    log_operation("VISUALIZE", "Generated character comparison dashboard")

    return {
        chart = chart,
        categories = categories,
        values = values,
        summary = "Higher bars indicate stronger loadouts"
    }
end

-- Create stat distribution visualization
function createStatDistribution(stat_name, values)
    load_phase11_dependencies()
    local viz = dependencies.visualization
    local stat_values = values or {95, 102, 110, 120, 128, 134}
    local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateLineChart(stat_name or "Stat Distribution", stat_values) or {type = "Line", generated = false}
    log_operation("VISUALIZE", string.format("Generated stat distribution for %s", stat_name or "stat"))
    return {chart = chart, data = stat_values}
end

-- Plot growth trajectory for a character
function plotGrowthTrajectory(char_id, checkpoints)
    load_phase11_dependencies()
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID for growth trajectory: " .. tostring(char_id))
        return nil
    end

    local viz = dependencies.visualization
    local checkpoints_data = checkpoints or {10, 20, 30, 40, 50, 60}
    local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateLineChart("Growth Trajectory", checkpoints_data) or {type = "Line", generated = false}
    log_operation("VISUALIZE", string.format("Plotted growth trajectory for %s", CONFIG.CHARACTERS[char_id].name))
    return {chart = chart, checkpoints = checkpoints_data}
end

-- Export a character sheet via report generation
function exportCharacterSheet(char_id, format, output_file)
    load_phase11_dependencies()
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID for character sheet: " .. tostring(char_id))
        return nil
    end

    local viz = dependencies.visualization
    local report_format = format or "PDF"
    local profile = build_character_profile(char_id)
    local report = viz and viz.ReportGeneration and viz.ReportGeneration.generateReport(profile.name .. " Sheet", profile, report_format) or {name = profile.name .. " Sheet"}

    if output_file and dependencies.import_export and dependencies.import_export.DataExporter and report_format == "JSON" then
        dependencies.import_export.DataExporter.exportToJSON(profile, output_file)
    elseif output_file then
        local file = io.open(output_file, "w")
        if file then
            file:write(serialize_table(report))
            file:close()
        end
    end

    log_operation("VISUALIZE", string.format("Exported character sheet for %s", profile.name))
    return report
end

-- Export character template using Import/Export Manager
function exportCharacterTemplate(char_id, format, output_path)
    load_phase11_dependencies()
    local template = exportCharacterBuild(char_id)
    if not template then return nil end

    local exporter = dependencies.import_export and dependencies.import_export.DataExporter
    local desired_format = (format or "json"):lower()
    local file_path = output_path or string.format("character_%d_template.%s", char_id, desired_format)

    if exporter then
        if desired_format == "json" then
            exporter.exportToJSON(template, file_path)
        elseif desired_format == "csv" then
            exporter.exportToCSV(template, file_path, true)
        elseif desired_format == "xml" then
            exporter.exportToXML(template, file_path)
        else
            exporter.exportToJSON(template, file_path)
        end
    else
        local file = io.open(file_path, "w")
        if file then
            file:write(serialize_table(template))
            file:close()
        end
    end

    log_operation("EXPORT", string.format("Exported character template for ID %d to %s", char_id, file_path))
    return {path = file_path, format = desired_format}
end

-- Import character template through Import/Export Manager
function importCharacterTemplate(path, format)
    load_phase11_dependencies()
    if not path then
        log_operation("ERROR", "No template path provided")
        return nil
    end

    local importer = dependencies.import_export and dependencies.import_export.DataImporter
    local fmt = format and format:lower() or "json"
    local imported

    if importer then
        if fmt == "json" then
            imported = importer.importJSON(path, true)
        elseif fmt == "csv" then
            imported = importer.importCSV(path, true)
        elseif fmt == "xml" then
            imported = importer.importXML(path, "lenient")
        else
            imported = importer.previewData(path)
        end
    end

    -- If JSON, run through build importer to apply
    if fmt == "json" then
        importCharacterBuild(path)
    end

    log_operation("IMPORT", string.format("Imported character template from %s", path))
    return imported
end

-- Batch import multiple character templates
function batchImportCharacters(file_paths)
    load_phase11_dependencies()
    if not file_paths or #file_paths == 0 then
        log_operation("ERROR", "No files provided for batch import")
        return {success = false, error = "No files"}
    end

    local results = {}
    for _, path in ipairs(file_paths) do
        table.insert(results, importCharacterTemplate(path, "json"))
    end

    log_operation("IMPORT", string.format("Batch imported %d character templates", #file_paths))
    return {success = true, imported = #file_paths, results = results}
end

-- Sync character data to other plugins via Integration Hub
function syncCharacterData(target_plugins)
    load_phase11_dependencies()
    local hub = dependencies.integration_hub
    local targets = target_plugins or {"party-optimizer", "challenge-mode-validator"}
    local sync_result = hub and hub.UnifiedAPI and hub.UnifiedAPI.crossPluginCall("character-roster-editor", table.concat(targets, ","), "sync", {characters = getRosterStatus()}) or {}
    log_operation("SYNC", "Pushed roster data to target plugins")
    return sync_result
end

-- Snapshot character build with Backup & Restore System
function snapshotCharacterBuild(char_id, label)
    load_phase11_dependencies()
    local backup = dependencies.backup_restore
    local snapshot_label = label or string.format("char_%d_snapshot", char_id)
    local build_data = exportCharacterBuild(char_id)

    if backup and backup.SnapshotManagement then
        local snapshot = backup.SnapshotManagement.createSnapshot(snapshot_label, build_data)
        log_operation("BACKUP", string.format("Snapshot created for %s", snapshot_label))
        return snapshot
    end

    create_backup()
    log_operation("BACKUP", "Fallback snapshot stored in plugin state")
    return {snapshot_id = "local", name = snapshot_label}
end

-- Restore character build from snapshot/version
function restoreCharacterBuild(snapshot_id, char_id)
    load_phase11_dependencies()
    local backup = dependencies.backup_restore
    local target_id = char_id or 0

    if backup and backup.RecoverySystem then
        local restore = backup.RecoverySystem.restoreFromBackup(snapshot_id or "latest", {target = target_id})
        log_operation("RESTORE", string.format("Restore triggered for snapshot %s", snapshot_id or "latest"))
        return restore
    end

    return restoreBackup()
end

-- Compare two build versions
function compareCharacterVersions(version_a, version_b)
    load_phase11_dependencies()
    local backup = dependencies.backup_restore
    if backup and backup.VersionControl then
        local comparison = backup.VersionControl.compareVersions(version_a or "v1", version_b or "v2")
        log_operation("COMPARE", "Compared character versions")
        return comparison
    end
    return {difference = "Version control unavailable"}
end

-- Enable automatic backups before roster operations
function autoBackupCharacters(save_file_path)
    load_phase11_dependencies()
    local backup = dependencies.backup_restore
    if backup and backup.QuickBackup then
        local result = backup.QuickBackup.autoBackupBeforeOperation(save_file_path or "save.srm", "roster_edit")
        log_operation("BACKUP", "Auto-backup executed for roster operation")
        return result
    end
    return {success = false, error = "QuickBackup unavailable"}
end

-- Register REST endpoints for roster data
function registerCharacterAPI()
    load_phase11_dependencies()
    local api = dependencies.api_gateway
    if not api or not api.RESTInterface then
        log_operation("WARN", "API Gateway not available for registration")
        return nil
    end

    local endpoint = api.RESTInterface.registerEndpoint("GET", "/api/roster", function()
        return getRosterStatus()
    end)
    api.RESTInterface.addAuthentication(endpoint.endpoint_id, "API-Key")
    api.RESTInterface.addRateLimit(endpoint.endpoint_id, 120)

    log_operation("API", "Registered /api/roster endpoint")
    return endpoint
end

-- Enable webhook notifications for roster changes
function enableWebhookNotifications(event_type, target_url)
    load_phase11_dependencies()
    local api = dependencies.api_gateway
    if not api or not api.WebhookManagement then
        log_operation("WARN", "Webhook subsystem unavailable")
        return nil
    end

    local webhook = api.WebhookManagement.registerWebhook(event_type or "roster.updated", target_url or "https://example.com/webhook")
    log_operation("API", "Webhook registered for roster updates")
    return webhook
end

-- Sync roster data to external database through API Gateway
function syncWithExternalDatabase(connection_string, dataset_name)
    load_phase11_dependencies()
    local api = dependencies.api_gateway
    if not api or not api.ExternalIntegration then
        log_operation("WARN", "External integration unavailable")
        return nil
    end

    local connection = api.ExternalIntegration.connectToAPI(dataset_name or "roster-db", "https://api.example.com/roster", {key = "demo"})
    local response = api.ExternalIntegration.callExternalAPI(connection.connection_id, "POST", "/sync", getRosterStatus())
    log_operation("API", "Synced roster data to external database")
    return {connection = connection, response = response}
end

-- ============================================================================
-- INITIALIZATION
-- ============================================================================

function initialize()
    if plugin_state.initialized then
        return true
    end
    
    log_operation("INIT", "Character Roster Editor initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("Character Roster Editor v1.2.0 loaded")
print("Commands: enableCharacter(id), disableCharacter(id), unlockAllCharacters()")
print("          configureSoloRun(id), applySoloPreset(name), applyRosterTemplate(name)")
print("          restrictPartySize(max), displayRosterStatus(), restoreBackup()")
print("Build Sharing: exportCharacterBuild(char_id, file), importCharacterBuild(file, target_id, preview)")
print("Analytics: predictCharacterGrowth(id), analyzeCharacterBalance(), generateCharacterComparison(ids)")
print("Type listRosterTemplates() to see available templates")
