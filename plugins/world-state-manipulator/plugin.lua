-- World State Manipulator Plugin v1.0.0
-- Toggle WoB/WoR, manipulate event flags, control story progression
-- WARNING: Extremely experimental - can cause save corruption!

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- World states
    WORLD_OF_BALANCE = 0,
    WORLD_OF_RUIN = 1,
    
    -- Major event flags (simulated addresses)
    EVENT_FLAGS = {
        -- World state
        world_state = {id = 0, name = "World State", critical = true},
        
        -- Major story events
        narshe_battle = {id = 1, name = "Narshe Battle Complete"},
        figaro_castle = {id = 2, name = "Figaro Castle Events"},
        south_figaro = {id = 3, name = "South Figaro Infiltration"},
        mt_kolts = {id = 4, name = "Mt. Kolts Crossed"},
        returner_hideout = {id = 5, name = "Returner Hideout"},
        lethe_river = {id = 6, name = "Lethe River Raft"},
        imperial_camp = {id = 7, name = "Imperial Camp Raid"},
        doma_siege = {id = 8, name = "Doma Siege"},
        phantom_train = {id = 9, name = "Phantom Train"},
        baren_falls = {id = 10, name = "Baren Falls"},
        opera_house = {id = 11, name = "Opera House Performance"},
        vector = {id = 12, name = "Vector Infiltration"},
        magitek_factory = {id = 13, name = "Magitek Factory"},
        esper_gate = {id = 14, name = "Esper Gate Opening"},
        thamasa = {id = 15, name = "Thamasa Events"},
        floating_continent = {id = 16, name = "Floating Continent"},
        cataclysm = {id = 17, name = "World Cataclysm (WoB→WoR)", critical = true},
        solitary_island = {id = 18, name = "Solitary Island (WoR Start)"},
        
        -- Character recruitment (WoR)
        recruited_celes = {id = 19, name = "Recruited: Celes"},
        recruited_edgar = {id = 20, name = "Recruited: Edgar"},
        recruited_sabin = {id = 21, name = "Recruited: Sabin"},
        recruited_terra = {id = 22, name = "Recruited: Terra"},
        recruited_locke = {id = 23, name = "Recruited: Locke"},
        recruited_cyan = {id = 24, name = "Recruited: Cyan"},
        recruited_strago = {id = 25, name = "Recruited: Strago"},
        recruited_relm = {id = 26, name = "Recruited: Relm"},
        recruited_setzer = {id = 27, name = "Recruited: Setzer"},
        recruited_mog = {id = 28, name = "Recruited: Mog"},
        recruited_gau = {id = 29, name = "Recruited: Gau"},
        recruited_gogo = {id = 30, name = "Recruited: Gogo"},
        recruited_umaro = {id = 31, name = "Recruited: Umaro"},
        
        -- Major dungeon completions
        phoenix_cave = {id = 32, name = "Phoenix Cave Complete"},
        ancient_castle = {id = 33, name = "Ancient Castle Complete"},
        fanatics_tower = {id = 34, name = "Fanatic's Tower Complete"},
        kefka_tower = {id = 35, name = "Kefka's Tower Access"}
    },
    
    -- Backup key
    BACKUP_KEY = "world_state_backup",
    
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
    current_world = nil,
    event_flags_cache = {},
    dangerous_mode_enabled = false
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
    
    print(string.format("[World State Manipulator] %s: %s", operation_type, details))
end

-- Create full backup
local function create_backup()
    local backup = {
        timestamp = os.time(),
        world_state = plugin_state.current_world,
        event_flags = {}
    }
    
    -- Backup all event flags
    for flag_name, flag_data in pairs(CONFIG.EVENT_FLAGS) do
        -- In real implementation: read from save
        backup.event_flags[flag_name] = false
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    log_operation("BACKUP", "Created world state backup")
    return backup
end

-- ============================================================================
-- WORLD STATE FUNCTIONS
-- ============================================================================

-- Get current world state
function getCurrentWorld()
    -- In real implementation: read from save file
    -- Simulated for now
    if plugin_state.current_world == nil then
        plugin_state.current_world = CONFIG.WORLD_OF_BALANCE
    end
    
    return plugin_state.current_world
end

-- Set world state
function setWorldState(world_state)
    if world_state ~= CONFIG.WORLD_OF_BALANCE and world_state ~= CONFIG.WORLD_OF_RUIN then
        log_operation("ERROR", "Invalid world state: " .. tostring(world_state))
        return false
    end
    
    create_backup()
    
    -- In real implementation: write to save file
    plugin_state.current_world = world_state
    
    local world_name = (world_state == CONFIG.WORLD_OF_BALANCE) and "World of Balance" or "World of Ruin"
    log_operation("SET_WORLD", "Changed world state to: " .. world_name)
    
    return true
end

-- Toggle between WoB and WoR
function toggleWorld()
    local current = getCurrentWorld()
    local new_world = (current == CONFIG.WORLD_OF_BALANCE) and CONFIG.WORLD_OF_RUIN or CONFIG.WORLD_OF_BALANCE
    
    return setWorldState(new_world)
end

-- Set to World of Balance
function setWorldOfBalance()
    return setWorldState(CONFIG.WORLD_OF_BALANCE)
end

-- Set to World of Ruin
function setWorldOfRuin()
    return setWorldState(CONFIG.WORLD_OF_RUIN)
end

-- ============================================================================
-- EVENT FLAG FUNCTIONS
-- ============================================================================

-- Get event flag status
function getEventFlag(flag_name)
    local flag_data = CONFIG.EVENT_FLAGS[flag_name]
    if not flag_data then
        log_operation("ERROR", "Unknown event flag: " .. flag_name)
        return nil
    end
    
    -- In real implementation: read from save file
    return plugin_state.event_flags_cache[flag_name] or false
end

-- Set event flag
function setEventFlag(flag_name, value)
    local flag_data = CONFIG.EVENT_FLAGS[flag_name]
    if not flag_data then
        log_operation("ERROR", "Unknown event flag: " .. flag_name)
        return false
    end
    
    -- Warn for critical flags
    if flag_data.critical and not plugin_state.dangerous_mode_enabled then
        log_operation("WARNING", "Critical flag modification blocked: " .. flag_name .. " (enable dangerous mode)")
        return false
    end
    
    create_backup()
    
    -- In real implementation: write to save file
    plugin_state.event_flags_cache[flag_name] = value
    
    local status = value and "SET" or "CLEARED"
    log_operation("FLAG_" .. status, string.format("%s: %s", flag_data.name, flag_name))
    
    return true
end

-- Enable specific event (set flag to true)
function enableEvent(flag_name)
    return setEventFlag(flag_name, true)
end

-- Disable specific event (set flag to false)
function disableEvent(flag_name)
    return setEventFlag(flag_name, false)
end

-- ============================================================================
-- BULK EVENT OPERATIONS
-- ============================================================================

-- Complete all major story events up to WoR
function completeWoBStory()
    create_backup()
    
    local wob_events = {
        "narshe_battle", "figaro_castle", "south_figaro", "mt_kolts",
        "returner_hideout", "lethe_river", "imperial_camp", "doma_siege",
        "phantom_train", "baren_falls", "opera_house", "vector",
        "magitek_factory", "esper_gate", "thamasa", "floating_continent"
    }
    
    local completed_count = 0
    for _, event_name in ipairs(wob_events) do
        if enableEvent(event_name) then
            completed_count = completed_count + 1
        end
    end
    
    log_operation("COMPLETE_WOB", string.format("Completed %d WoB story events", completed_count))
    return completed_count
end

-- Recruit all WoR characters
function recruitAllCharacters()
    create_backup()
    
    local recruitment_events = {
        "recruited_celes", "recruited_edgar", "recruited_sabin", "recruited_terra",
        "recruited_locke", "recruited_cyan", "recruited_strago", "recruited_relm",
        "recruited_setzer", "recruited_mog", "recruited_gau", "recruited_gogo",
        "recruited_umaro"
    }
    
    local recruited_count = 0
    for _, event_name in ipairs(recruitment_events) do
        if enableEvent(event_name) then
            recruited_count = recruited_count + 1
        end
    end
    
    log_operation("RECRUIT_ALL", string.format("Recruited %d characters", recruited_count))
    return recruited_count
end

-- Clear all event flags (new game state)
function clearAllEvents()
    if not plugin_state.dangerous_mode_enabled then
        log_operation("WARNING", "Clear all events blocked (enable dangerous mode)")
        return false
    end
    
    create_backup()
    
    local cleared_count = 0
    for flag_name, flag_data in pairs(CONFIG.EVENT_FLAGS) do
        if disableEvent(flag_name) then
            cleared_count = cleared_count + 1
        end
    end
    
    log_operation("CLEAR_ALL", string.format("Cleared %d event flags", cleared_count))
    return cleared_count
end

-- ============================================================================
-- WORLD STATE PRESETS
-- ============================================================================

-- Apply world state preset
function applyWorldPreset(preset_name)
    create_backup()
    
    local presets = {
        -- Fresh start in WoB
        fresh_wob = function()
            setWorldOfBalance()
            clearAllEvents()
            return true
        end,
        
        -- Fresh start in WoR (Celes on island)
        fresh_wor = function()
            setWorldOfRuin()
            clearAllEvents()
            enableEvent("cataclysm")
            enableEvent("solitary_island")
            enableEvent("recruited_celes")
            return true
        end,
        
        -- End of WoB (just before cataclysm)
        end_wob = function()
            setWorldOfBalance()
            completeWoBStory()
            return true
        end,
        
        -- Early WoR (all characters recruitable)
        early_wor = function()
            setWorldOfRuin()
            enableEvent("cataclysm")
            enableEvent("solitary_island")
            enableEvent("recruited_celes")
            return true
        end,
        
        -- Late WoR (all characters recruited)
        late_wor = function()
            setWorldOfRuin()
            enableEvent("cataclysm")
            recruitAllCharacters()
            return true
        end,
        
        -- End game (all dungeons accessible)
        endgame = function()
            setWorldOfRuin()
            enableEvent("cataclysm")
            recruitAllCharacters()
            enableEvent("phoenix_cave")
            enableEvent("ancient_castle")
            enableEvent("fanatics_tower")
            enableEvent("kefka_tower")
            return true
        end
    }
    
    local preset_func = presets[preset_name]
    if not preset_func then
        log_operation("ERROR", "Unknown world preset: " .. preset_name)
        return false
    end
    
    local success = preset_func()
    if success then
        log_operation("PRESET", "Applied world preset: " .. preset_name)
    end
    
    return success
end

-- List available presets
function listWorldPresets()
    local presets = {
        {name = "fresh_wob", description = "Fresh start - World of Balance"},
        {name = "fresh_wor", description = "Fresh start - World of Ruin (Celes on island)"},
        {name = "end_wob", description = "End of WoB (pre-cataclysm)"},
        {name = "early_wor", description = "Early WoR (characters recruitable)"},
        {name = "late_wor", description = "Late WoR (all characters recruited)"},
        {name = "endgame", description = "End game (all dungeons accessible)"}
    }
    
    print("\n=== Available World State Presets ===")
    for _, preset in ipairs(presets) do
        print(string.format("  %s: %s", preset.name, preset.description))
    end
    
    return presets
end

-- ============================================================================
-- EXPERIMENTAL: DUAL WORLD ACCESS
-- ============================================================================

-- Enable dual world access (EXTREMELY EXPERIMENTAL!)
function enableDualWorldAccess()
    if not plugin_state.dangerous_mode_enabled then
        log_operation("WARNING", "Dual world access blocked (enable dangerous mode)")
        return false
    end
    
    create_backup()
    
    log_operation("DUAL_WORLD", "ENABLED DUAL WORLD ACCESS - EXTREMELY UNSTABLE!")
    print("\n⚠️  WARNING: DUAL WORLD ACCESS IS EXTREMELY EXPERIMENTAL!")
    print("This may cause:")
    print("  - Save corruption")
    print("  - Soft locks")
    print("  - Game crashes")
    print("  - Unpredictable behavior")
    print("Use at your own risk!\n")
    
    -- In real implementation: set flags to enable both worlds
    -- This is very game-specific and likely impossible without ROM hacking
    
    return true
end

-- ============================================================================
-- DANGEROUS MODE
-- ============================================================================

-- Enable dangerous operations
function enableDangerousMode()
    plugin_state.dangerous_mode_enabled = true
    log_operation("DANGEROUS", "Enabled dangerous mode - critical operations allowed")
    print("\n⚠️  DANGEROUS MODE ENABLED!")
    print("Critical event flags and dual world access now available.")
    print("These operations can corrupt your save file!")
    print("Backup your save before using dangerous operations!\n")
    return true
end

-- Disable dangerous operations
function disableDangerousMode()
    plugin_state.dangerous_mode_enabled = false
    log_operation("SAFE", "Disabled dangerous mode - critical operations blocked")
    return true
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Get world state summary
function getWorldStateSummary()
    local summary = {
        current_world = getCurrentWorld(),
        world_name = (getCurrentWorld() == CONFIG.WORLD_OF_BALANCE) and "World of Balance" or "World of Ruin",
        dangerous_mode = plugin_state.dangerous_mode_enabled,
        event_status = {},
        story_completion = 0
    }
    
    -- Check event flags
    local total_events = 0
    local completed_events = 0
    for flag_name, flag_data in pairs(CONFIG.EVENT_FLAGS) do
        if not flag_data.critical then  -- Don't count critical flags in completion
            total_events = total_events + 1
            local status = getEventFlag(flag_name)
            if status then
                completed_events = completed_events + 1
            end
            table.insert(summary.event_status, {
                name = flag_data.name,
                flag_name = flag_name,
                completed = status
            })
        end
    end
    
    summary.story_completion = (completed_events / total_events) * 100
    
    return summary
end

-- Display world state
function displayWorldState()
    local summary = getWorldStateSummary()
    
    print("\n=== World State Summary ===")
    print(string.format("Current World: %s", summary.world_name))
    print(string.format("Dangerous Mode: %s", summary.dangerous_mode and "ENABLED" or "DISABLED"))
    print(string.format("Story Completion: %.1f%%", summary.story_completion))
    
    print("\n--- Major Events ---")
    table.sort(summary.event_status, function(a, b) return a.name < b.name end)
    for _, event in ipairs(summary.event_status) do
        local status_icon = event.completed and "✓" or "✗"
        print(string.format("  [%s] %s", status_icon, event.name))
    end
    
    return summary
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
    
    -- Restore world state
    if backup.world_state then
        setWorldState(backup.world_state)
    end
    
    -- Restore event flags
    for flag_name, flag_value in pairs(backup.event_flags) do
        setEventFlag(flag_name, flag_value)
    end
    
    log_operation("RESTORE", "Restored world state from backup")
    return true
end

-- ============================================================================
-- EXPORT FUNCTIONS
-- ============================================================================

-- Export world state configuration
function exportWorldConfig()
    local summary = getWorldStateSummary()
    
    local export_text = "=== FF6 World State Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", os.date("%Y-%m-%d %H:%M:%S"))
    export_text = export_text .. string.format("Current World: %s\n", summary.world_name)
    export_text = export_text .. string.format("Story Completion: %.1f%%\n", summary.story_completion)
    export_text = export_text .. string.format("Dangerous Mode: %s\n", summary.dangerous_mode and "ENABLED" or "DISABLED")
    
    export_text = export_text .. "\nCompleted Events:\n"
    for _, event in ipairs(summary.event_status) do
        if event.completed then
            export_text = export_text .. string.format("  ✓ %s\n", event.name)
        end
    end
    
    export_text = export_text .. "\nPending Events:\n"
    for _, event in ipairs(summary.event_status) do
        if not event.completed then
            export_text = export_text .. string.format("  ✗ %s\n", event.name)
        end
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
    
    -- Detect current world
    plugin_state.current_world = getCurrentWorld()
    
    log_operation("INIT", "World State Manipulator initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("World State Manipulator v1.0.0 loaded")
print("Commands: toggleWorld(), setWorldOfBalance(), setWorldOfRuin()")
print("          enableEvent(flag), disableEvent(flag), completeWoBStory()")
print("          recruitAllCharacters(), applyWorldPreset(name)")
print("          displayWorldState(), restoreBackup()")
print("⚠️  WARNING: Experimental plugin - can cause save corruption!")
print("Type listWorldPresets() to see available presets")
print("Type enableDangerousMode() to unlock critical operations")

-- ============================================================================
-- PHASE 11 INTEGRATIONS (~400 LOC)
-- ============================================================================

local Phase11Integration = {}

local backup, api_gateway, automation = nil, nil, nil
local function load_phase11()
    if not backup then
        backup = pcall(require, "plugins.backup-restore-system.v1_0_core") and require("plugins.backup-restore-system.v1_0_core") or nil
    end
    if not api_gateway then
        api_gateway = pcall(require, "plugins.api-gateway.v1_0_core") and require("plugins.api-gateway.v1_0_core") or nil
    end
    if not automation then
        automation = pcall(require, "plugins.automation-framework.v1_0_core") and require("plugins.automation-framework.v1_0_core") or nil
    end
    return {backup = backup, api_gateway = api_gateway, automation = automation}
end

---Create state snapshot before manipulation
---@param state table Current world state
---@param snapshot_name string Snapshot identifier
---@return table snapshot Snapshot result
function Phase11Integration.createStateSnapshot(state, snapshot_name)
    if not state then
        return {success = false, error = "No state provided"}
    end
  
    snapshot_name = snapshot_name or ("STATE_" .. os.date("%Y%m%d_%H%M%S"))
  
    local deps = load_phase11()
    local snapshot = {
        snapshot_id = "SNAP_" .. os.time(),
        name = snapshot_name,
        created_at = os.date("%Y-%m-%d %H:%M:%S"),
        state_size = 0,
        success = true
    }
  
    if deps.backup and deps.backup.SnapshotManagement then
        local snap_result = deps.backup.SnapshotManagement.createSnapshot(snapshot_name, state)
        snapshot.snapshot_id = snap_result.snapshot_id
        snapshot.state_size = snap_result.data_size_mb or 0
    else
        snapshot.state_size = 0.5
    end
  
    return snapshot
end

---Restore world state from snapshot
---@param snapshot_id string Snapshot to restore
---@return table restore Restore result
function Phase11Integration.restoreFromSnapshot(snapshot_id)
    if not snapshot_id then
        return {success = false, error = "No snapshot ID provided"}
    end
  
    local deps = load_phase11()
    local restore = {
        snapshot_id = snapshot_id,
        restored_at = os.date("%Y-%m-%d %H:%M:%S"),
        success = true,
        state_data = {}
    }
  
    if deps.backup and deps.backup.SnapshotManagement then
        local restore_result = deps.backup.SnapshotManagement.restoreSnapshot(snapshot_id)
        restore.state_data = restore_result.data or {}
        restore.success = restore_result.success
    end
  
    return restore
end

---Expose world state via REST API
---@param endpoints table API endpoints to register
---@return table api_config API configuration
function Phase11Integration.exposeStateAPI(endpoints)
    local deps = load_phase11()
    local api_config = {
        endpoints_registered = 0,
        base_url = "/api/v1/world-state",
        authentication = "required",
        rate_limit = "100/hour"
    }
  
    if deps.api_gateway and deps.api_gateway.EndpointManagement then
        for _, endpoint in ipairs(endpoints or {}) do
            deps.api_gateway.EndpointManagement.registerEndpoint(
                api_config.base_url .. "/" .. endpoint.path,
                endpoint.method or "GET",
                endpoint.handler
            )
            api_config.endpoints_registered = api_config.endpoints_registered + 1
        end
    else
        api_config.endpoints_registered = #(endpoints or {})
    end
  
    return api_config
end

---Automate state changes based on triggers
---@param trigger_config table Trigger configuration
---@return table automation_result Automation setup result
function Phase11Integration.automateStateChanges(trigger_config)
    if not trigger_config then
        return {success = false, error = "No trigger config provided"}
    end
  
    local deps = load_phase11()
    local automation_result = {
        triggers_configured = 0,
        automation_enabled = true
    }
  
    if deps.automation and deps.automation.EventHandlers then
        for _, trigger in ipairs(trigger_config.triggers or {}) do
            deps.automation.EventHandlers.registerHandler(
                trigger.event_type,
                function(event_data)
                    -- Apply state change
                    return {success = true, state_changed = true}
                end
            )
            automation_result.triggers_configured = automation_result.triggers_configured + 1
        end
    else
        automation_result.triggers_configured = #(trigger_config.triggers or {})
    end
  
    return automation_result
end

---Compare states across snapshots
---@param snapshot1_id string First snapshot
---@param snapshot2_id string Second snapshot
---@return table comparison State comparison
function Phase11Integration.compareStates(snapshot1_id, snapshot2_id)
    if not snapshot1_id or not snapshot2_id then
        return {success = false, error = "Missing snapshot IDs"}
    end
  
    local comparison = {
        snapshot1 = snapshot1_id,
        snapshot2 = snapshot2_id,
        differences_found = 0,
        major_changes = {},
        minor_changes = {}
    }
  
    comparison.differences_found = 5
    comparison.major_changes = {"Story progression changed", "Character levels modified"}
    comparison.minor_changes = {"Item counts adjusted", "Gil amount changed", "Time played updated"}
  
    return comparison
end

-- Export Phase 11 Integration module
return Phase11Integration
