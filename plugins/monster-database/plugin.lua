--[[
  Monster Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Provides searchable monster bestiary with analytics, export, and cross-plugin integration
  
  Features:
  - Monster lookup by ID, name, type, or location
  - Monster stats and weaknesses tracking
  - Analytics for encounter patterns
  - Export/Import bestiary
  - Integration Hub sync for cross-plugin monster data
  - Backup/Restore for bestiary
  
  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Monster types
    MONSTER_TYPES = {
        BEAST = "beast",
        UNDEAD = "undead",
        HUMAN = "human",
        DRAGON = "dragon",
        ESPER = "esper",
        MACHINE = "machine",
        BOSS = "boss"
    },
    
    -- Element types
    ELEMENTS = {
        FIRE = "fire",
        ICE = "ice",
        LIGHTNING = "lightning",
        WATER = "water",
        EARTH = "earth",
        WIND = "wind",
        HOLY = "holy",
        POISON = "poison"
    },
    
    MAX_MONSTERS = 384,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    bestiary = {},
    encounter_history = {},
    weakness_cache = {},
    operation_log = {}
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

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
    
    print(string.format("[Monster Database] %s: %s", operation_type, details))
end

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

-- Phase 11 dependency handles (lazy-loaded)
local dependencies = {
    analytics = nil,
    import_export = nil,
    backup_restore = nil,
    integration_hub = nil,
    api_gateway = nil
}

local function load_phase11_dependencies()
    dependencies.analytics = dependencies.analytics or safe_require("plugins.advanced-analytics-engine.v1_0_core")
    dependencies.import_export = dependencies.import_export or safe_require("plugins.import-export-manager.v1_0_core")
    dependencies.backup_restore = dependencies.backup_restore or safe_require("plugins.backup-restore-system.v1_0_core")
    dependencies.integration_hub = dependencies.integration_hub or safe_require("plugins.integration-hub.v1_0_core")
    dependencies.api_gateway = dependencies.api_gateway or safe_require("plugins.api-gateway.v1_0_core")
    return dependencies
end

-- ============================================================================
-- UTILITY HELPER
-- ============================================================================

local function table_count(t)
    local count = 0
    for _ in pairs(t) do count = count + 1 end
    return count
end

-- ============================================================================
-- MONSTER DATABASE INITIALIZATION
-- ============================================================================

local function initialize_bestiary()
    if plugin_state.initialized then return end
    
    -- Sample bestiary (in real implementation, load from game data)
    plugin_state.bestiary = {
        -- Normal monsters
        [0] = {
            id = 0, name = "Guard", type = CONFIG.MONSTER_TYPES.HUMAN,
            hp = 50, mp = 10, attack = 10, defense = 10,
            weaknesses = {CONFIG.ELEMENTS.POISON}, location = "Narshe",
            exp = 10, gp = 20
        },
        [1] = {
            id = 1, name = "Leaf Bunny", type = CONFIG.MONSTER_TYPES.BEAST,
            hp = 25, mp = 5, attack = 5, defense = 5,
            weaknesses = {CONFIG.ELEMENTS.FIRE}, location = "Narshe",
            exp = 5, gp = 10
        },
        
        -- Mid-tier monsters
        [100] = {
            id = 100, name = "Fossil Fang", type = CONFIG.MONSTER_TYPES.UNDEAD,
            hp = 1000, mp = 100, attack = 40, defense = 35,
            weaknesses = {CONFIG.ELEMENTS.FIRE, CONFIG.ELEMENTS.HOLY},
            location = "Mt. Kolts",
            exp = 150, gp = 300
        },
        
        -- Dragons
        [200] = {
            id = 200, name = "Ice Dragon", type = CONFIG.MONSTER_TYPES.DRAGON,
            hp = 24400, mp = 5000, attack = 13, defense = 110,
            weaknesses = {CONFIG.ELEMENTS.FIRE}, absorbs = {CONFIG.ELEMENTS.ICE},
            location = "Narshe Cliffs",
            exp = 0, gp = 0, boss = true
        },
        
        -- Espers
        [250] = {
            id = 250, name = "Tritoch", type = CONFIG.MONSTER_TYPES.ESPER,
            hp = 30000, mp = 50000, attack = 13, defense = 125,
            weaknesses = {}, absorbs = {CONFIG.ELEMENTS.ICE},
            location = "Narshe",
            exp = 0, gp = 0, boss = true
        },
        
        -- Machines
        [300] = {
            id = 300, name = "Guardian", type = CONFIG.MONSTER_TYPES.MACHINE,
            hp = 60000, mp = 10000, attack = 30, defense = 200,
            weaknesses = {CONFIG.ELEMENTS.LIGHTNING}, location = "Kefka's Tower",
            exp = 0, gp = 0, boss = true
        },
        
        -- Final boss
        [383] = {
            id = 383, name = "Kefka (Final)", type = CONFIG.MONSTER_TYPES.BOSS,
            hp = 62000, mp = 50000, attack = 25, defense = 230,
            weaknesses = {}, location = "Kefka's Tower",
            exp = 0, gp = 0, boss = true, final_boss = true
        }
    }
    
    plugin_state.initialized = true
    log_operation("INIT", string.format("Monster database initialized with %d monsters", table_count(plugin_state.bestiary)))
end

-- ============================================================================
-- CORE MONSTER LOOKUP FUNCTIONS
-- ============================================================================

function getMonsterById(monster_id)
    initialize_bestiary()
    
    if not monster_id or monster_id < 0 or monster_id >= CONFIG.MAX_MONSTERS then
        log_operation("ERROR", "Invalid monster ID: " .. tostring(monster_id))
        return nil
    end
    
    local monster = plugin_state.bestiary[monster_id]
    if monster then
        log_operation("LOOKUP", string.format("Retrieved monster %d: %s", monster_id, monster.name))
    end
    
    return monster
end

function searchMonstersByName(name_query)
    initialize_bestiary()
    
    local results = {}
    local query_lower = string.lower(name_query or "")
    
    for monster_id, monster in pairs(plugin_state.bestiary) do
        if string.find(string.lower(monster.name), query_lower, 1, true) then
            table.insert(results, monster)
        end
    end
    
    log_operation("SEARCH", string.format("Found %d monsters matching '%s'", #results, name_query))
    return results
end

function getMonstersByType(monster_type)
    initialize_bestiary()
    
    local results = {}
    
    for monster_id, monster in pairs(plugin_state.bestiary) do
        if monster.type == monster_type then
            table.insert(results, monster)
        end
    end
    
    log_operation("FILTER", string.format("Found %d monsters of type '%s'", #results, monster_type))
    return results
end

function getMonstersByLocation(location)
    initialize_bestiary()
    
    local results = {}
    
    for monster_id, monster in pairs(plugin_state.bestiary) do
        if monster.location and string.find(string.lower(monster.location), string.lower(location), 1, true) then
            table.insert(results, monster)
        end
    end
    
    log_operation("FILTER", string.format("Found %d monsters in location '%s'", #results, location))
    return results
end

function getMonsterWeaknesses(monster_id)
    initialize_bestiary()
    
    local monster = getMonsterById(monster_id)
    if not monster then return nil end
    
    return {
        weaknesses = monster.weaknesses or {},
        absorbs = monster.absorbs or {},
        immunities = monster.immunities or {}
    }
end

-- ============================================================================
-- BESTIARY MANAGEMENT
-- ============================================================================

function addMonsterToBestiary(monster_data)
    initialize_bestiary()
    
    if not monster_data or not monster_data.id or not monster_data.name then
        log_operation("ERROR", "Invalid monster data")
        return false
    end
    
    plugin_state.bestiary[monster_data.id] = monster_data
    
    -- Persist to database layer
    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("bestiary", plugin_state.bestiary)
    end
    
    log_operation("ADD", string.format("Added monster %d: %s", monster_data.id, monster_data.name))
    return true
end

function getBestiarySummary()
    initialize_bestiary()
    
    local summary = {
        total_monsters = table_count(plugin_state.bestiary),
        by_type = {},
        bosses = 0,
        total_hp = 0,
        avg_hp = 0
    }
    
    local hp_sum = 0
    
    -- Count by type and calculate stats
    for _, monster in pairs(plugin_state.bestiary) do
        summary.by_type[monster.type] = (summary.by_type[monster.type] or 0) + 1
        if monster.boss then summary.bosses = summary.bosses + 1 end
        hp_sum = hp_sum + (monster.hp or 0)
    end
    
    summary.total_hp = hp_sum
    summary.avg_hp = summary.total_monsters > 0 and (hp_sum / summary.total_monsters) or 0
    
    return summary
end

function recordEncounter(monster_id)
    initialize_bestiary()
    
    plugin_state.encounter_history[monster_id] = (plugin_state.encounter_history[monster_id] or 0) + 1
    
    log_operation("ENCOUNTER", string.format("Recorded encounter with monster %d", monster_id))
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

-- Analyze encounter patterns via Analytics Engine
function analyzeEncounterPatterns()
    load_phase11_dependencies()
    initialize_bestiary()
    
    local analytics = dependencies.analytics
    
    -- Collect encounter data
    local encounter_data = {}
    for monster_id, count in pairs(plugin_state.encounter_history) do
        local monster = plugin_state.bestiary[monster_id]
        if monster then
            table.insert(encounter_data, {
                monster_id = monster_id,
                name = monster.name,
                count = count,
                location = monster.location
            })
        end
    end
    
    local analysis = {
        total_monsters = table_count(plugin_state.bestiary),
        encounters_tracked = #encounter_data,
        patterns = {}
    }
    
    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(encounter_data)
        analysis.patterns = patterns
    end
    
    log_operation("ANALYTICS", string.format("Analyzed encounter patterns for %d monsters", analysis.encounters_tracked))
    return analysis
end

-- Export bestiary via Import/Export Manager
function exportBestiary(format, path)
    load_phase11_dependencies()
    initialize_bestiary()
    
    local exporter = dependencies.import_export and dependencies.import_export.DataExporter
    
    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        monster_count = table_count(plugin_state.bestiary),
        bestiary = plugin_state.bestiary
    }
    
    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("bestiary_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)
        
        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end
        
        log_operation("EXPORT", string.format("Exported bestiary to %s", output))
        return {path = output, format = fmt}
    end
    
    return {success = false, error = "Import/Export unavailable"}
end

-- Sync bestiary to Integration Hub
function syncBestiaryToHub()
    load_phase11_dependencies()
    initialize_bestiary()
    
    local hub = dependencies.integration_hub
    
    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("bestiary_sync", {
            bestiary = plugin_state.bestiary,
            summary = getBestiarySummary(),
            timestamp = os.time()
        })
        
        log_operation("SYNC_HUB", "Synced bestiary to Integration Hub")
        return result or {success = true}
    end
    
    return {success = false, error = "Integration Hub unavailable"}
end

-- Create bestiary snapshot via Backup/Restore System
function createBestiarySnapshot(label)
    load_phase11_dependencies()
    initialize_bestiary()
    
    local backup = dependencies.backup_restore
    
    local snapshot = {
        label = label or "bestiary_snapshot",
        timestamp = os.time(),
        bestiary = plugin_state.bestiary,
        summary = getBestiarySummary()
    }
    
    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created bestiary snapshot: %s", label))
        return snap_id
    end
    
    -- Persist to database layer as fallback
    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("snapshot_" .. label, snapshot)
        return {snapshot_id = "local_" .. label}
    end
    
    return {success = false}
end

-- Register REST API endpoints for monster database
function registerMonsterDatabaseAPI()
    load_phase11_dependencies()
    initialize_bestiary()
    
    local api = dependencies.api_gateway
    
    if not api or not api.RESTInterface then return {success = false} end
    
    -- Monster lookup endpoint
    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/monsters/:id", function(params)
        return getMonsterById(tonumber(params.id))
    end)
    
    -- Monster search endpoint
    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/monsters/search", function(params)
        return searchMonstersByName(params.q or "")
    end)
    
    -- Bestiary summary endpoint
    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/monsters/summary", function()
        return getBestiarySummary()
    end)
    
    -- Weakness lookup endpoint
    local weakness_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/monsters/:id/weaknesses", function(params)
        return getMonsterWeaknesses(tonumber(params.id))
    end)
    
    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end
    
    log_operation("API", "Registered monster database REST endpoints")
    return {
        lookup = lookup_endpoint or {registered = true}, 
        search = search_endpoint or {registered = true}, 
        summary = summary_endpoint or {registered = true},
        weakness = weakness_endpoint or {registered = true}
    }
end

-- Initialize on load
initialize_bestiary()
log_operation("LOAD", "Monster Database Plugin v1.0.0 loaded")

-- Export public API
return {
    -- Core lookup
    getMonsterById = getMonsterById,
    searchMonstersByName = searchMonstersByName,
    getMonstersByType = getMonstersByType,
    getMonstersByLocation = getMonstersByLocation,
    getMonsterWeaknesses = getMonsterWeaknesses,
    
    -- Bestiary management
    addMonsterToBestiary = addMonsterToBestiary,
    getBestiarySummary = getBestiarySummary,
    recordEncounter = recordEncounter,
    
    -- Phase 11 integrations
    analyzeEncounterPatterns = analyzeEncounterPatterns,
    exportBestiary = exportBestiary,
    syncBestiaryToHub = syncBestiaryToHub,
    createBestiarySnapshot = createBestiarySnapshot,
    registerMonsterDatabaseAPI = registerMonsterDatabaseAPI
}
