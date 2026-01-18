--[[
  Treasure Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Tracks treasure chests, rewards, and locations with analytics, export, and integration

  Features:
  - Treasure lookup by ID, location, type, or rarity
  - Contents and status metadata
  - Analytics for collection progress
  - Export/Import treasure map
  - Integration Hub sync for cross-plugin treasure data
  - Backup/Restore snapshots

  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    TREASURE_TYPES = {
        CHEST = "chest",
        HIDDEN = "hidden",
        EVENT = "event",
        SHOP = "shop"
    },

    RARITY = {
        COMMON = 1,
        UNCOMMON = 2,
        RARE = 3,
        EPIC = 4,
        LEGENDARY = 5
    },

    MAX_TREASURE = 512,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    treasure_map = {},
    collection_history = {},
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

    print(string.format("[Treasure Database] %s: %s", operation_type, details))
end

local function safe_require(module_path)
    local ok, mod = pcall(require, module_path)
    if not ok then
        log_operation("WARN", "Dependency unavailable: " .. module_path)
        return nil
    end
    return mod
end

local function table_count(t)
    local count = 0
    for _ in pairs(t) do count = count + 1 end
    return count
end

local database_layer = nil

local function load_database_layer()
    if not database_layer then
        database_layer = safe_require("plugins.database-persistence-layer.plugin")
    end
    return database_layer
end

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
-- TREASURE INITIALIZATION
-- ============================================================================

local function initialize_treasure_map()
    if plugin_state.initialized then return end

    plugin_state.treasure_map = {
        [0] = {id = 0, location = "Narshe Mines", type = CONFIG.TREASURE_TYPES.CHEST, rarity = CONFIG.RARITY.COMMON, contents = "Potion", opened = false},
        [5] = {id = 5, location = "South Figaro", type = CONFIG.TREASURE_TYPES.HIDDEN, rarity = CONFIG.RARITY.UNCOMMON, contents = "Phoenix Down", opened = false},
        [20] = {id = 20, location = "Mt. Kolts", type = CONFIG.TREASURE_TYPES.CHEST, rarity = CONFIG.RARITY.UNCOMMON, contents = "Mythril Claw", opened = true},
        [60] = {id = 60, location = "Opera House", type = CONFIG.TREASURE_TYPES.EVENT, rarity = CONFIG.RARITY.RARE, contents = "Earrings", opened = true},
        [120] = {id = 120, location = "Phoenix Cave", type = CONFIG.TREASURE_TYPES.CHEST, rarity = CONFIG.RARITY.EPIC, contents = "Ribbon", opened = false},
        [180] = {id = 180, location = "Ancient Castle", type = CONFIG.TREASURE_TYPES.HIDDEN, rarity = CONFIG.RARITY.LEGENDARY, contents = "Offering", opened = false},
        [240] = {id = 240, location = "Kefka's Tower", type = CONFIG.TREASURE_TYPES.CHEST, rarity = CONFIG.RARITY.LEGENDARY, contents = "Illumina", opened = false}
    }

    plugin_state.initialized = true
    log_operation("INIT", string.format("Treasure database initialized with %d entries", table_count(plugin_state.treasure_map)))
end

-- ============================================================================
-- CORE TREASURE LOOKUP FUNCTIONS
-- ============================================================================

function getTreasureById(treasure_id)
    initialize_treasure_map()

    if not treasure_id or treasure_id < 0 or treasure_id >= CONFIG.MAX_TREASURE then
        log_operation("ERROR", "Invalid treasure ID: " .. tostring(treasure_id))
        return nil
    end

    local treasure = plugin_state.treasure_map[treasure_id]
    if treasure then
        log_operation("LOOKUP", string.format("Retrieved treasure %d at %s", treasure_id, treasure.location))
    end

    return treasure
end

function getTreasureByLocation(location)
    initialize_treasure_map()

    local results = {}
    local query_lower = string.lower(location or "")

    for treasure_id, treasure in pairs(plugin_state.treasure_map) do
        if string.find(string.lower(treasure.location), query_lower, 1, true) then
            table.insert(results, treasure)
        end
    end

    log_operation("FILTER", string.format("Found %d treasures in location '%s'", #results, location))
    return results
end

function getTreasureByType(treasure_type)
    initialize_treasure_map()

    local results = {}

    for treasure_id, treasure in pairs(plugin_state.treasure_map) do
        if treasure.type == treasure_type then
            table.insert(results, treasure)
        end
    end

    log_operation("FILTER", string.format("Found %d treasures of type '%s'", #results, treasure_type))
    return results
end

function getTreasureByRarity(rarity)
    initialize_treasure_map()

    local results = {}

    for treasure_id, treasure in pairs(plugin_state.treasure_map) do
        if treasure.rarity == rarity then
            table.insert(results, treasure)
        end
    end

    log_operation("FILTER", string.format("Found %d treasures with rarity %d", #results, rarity))
    return results
end

-- ============================================================================
-- TREASURE MANAGEMENT
-- ============================================================================

function addTreasure(treasure_data)
    initialize_treasure_map()

    if not treasure_data or not treasure_data.id or not treasure_data.location then
        log_operation("ERROR", "Invalid treasure data")
        return false
    end

    plugin_state.treasure_map[treasure_data.id] = treasure_data

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("treasure_map", plugin_state.treasure_map)
    end

    log_operation("ADD", string.format("Added treasure %d at %s", treasure_data.id, treasure_data.location))
    return true
end

function markTreasureOpened(treasure_id)
    initialize_treasure_map()

    local treasure = plugin_state.treasure_map[treasure_id]
    if not treasure then
        log_operation("ERROR", "Treasure not found for open")
        return false
    end

    treasure.opened = true

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("treasure_map", plugin_state.treasure_map)
    end

    log_operation("UPDATE", string.format("Marked treasure %d as opened", treasure_id))
    return true
end

function getTreasureSummary()
    initialize_treasure_map()

    local summary = {
        total_treasure = table_count(plugin_state.treasure_map),
        by_type = {},
        by_rarity = {},
        opened = 0
    }

    for _, treasure in pairs(plugin_state.treasure_map) do
        summary.by_type[treasure.type] = (summary.by_type[treasure.type] or 0) + 1
        summary.by_rarity[treasure.rarity] = (summary.by_rarity[treasure.rarity] or 0) + 1
        if treasure.opened then summary.opened = summary.opened + 1 end
    end

    return summary
end

function recordCollection(treasure_id)
    initialize_treasure_map()

    plugin_state.collection_history[treasure_id] = (plugin_state.collection_history[treasure_id] or 0) + 1
    log_operation("COLLECT", string.format("Recorded collection for treasure %d", treasure_id))
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

function analyzeCollectionProgress()
    load_phase11_dependencies()
    initialize_treasure_map()

    local analytics = dependencies.analytics

    local collection_data = {}
    for treasure_id, count in pairs(plugin_state.collection_history) do
        local treasure = plugin_state.treasure_map[treasure_id]
        if treasure then
            table.insert(collection_data, {
                treasure_id = treasure_id,
                location = treasure.location,
                type = treasure.type,
                rarity = treasure.rarity,
                count = count
            })
        end
    end

    local analysis = {
        total_treasure = table_count(plugin_state.treasure_map),
        collection_tracked = #collection_data,
        patterns = {}
    }

    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(collection_data)
        analysis.patterns = patterns
    end

    log_operation("ANALYTICS", string.format("Analyzed collection progress for %d treasures", analysis.collection_tracked))
    return analysis
end

function exportTreasureMap(format, path)
    load_phase11_dependencies()
    initialize_treasure_map()

    local exporter = dependencies.import_export and dependencies.import_export.DataExporter

    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        treasure_count = table_count(plugin_state.treasure_map),
        treasure = plugin_state.treasure_map
    }

    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("treasure_map_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)

        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end

        log_operation("EXPORT", string.format("Exported treasure map to %s", output))
        return {path = output, format = fmt}
    end

    return {success = false, error = "Import/Export unavailable"}
end

function syncTreasureMapToHub()
    load_phase11_dependencies()
    initialize_treasure_map()

    local hub = dependencies.integration_hub

    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("treasure_map_sync", {
            treasure = plugin_state.treasure_map,
            summary = getTreasureSummary(),
            timestamp = os.time()
        })

        log_operation("SYNC_HUB", "Synced treasure map to Integration Hub")
        return result or {success = true}
    end

    return {success = false, error = "Integration Hub unavailable"}
end

function createTreasureSnapshot(label)
    load_phase11_dependencies()
    initialize_treasure_map()

    local backup = dependencies.backup_restore

    local snapshot = {
        label = label or "treasure_snapshot",
        timestamp = os.time(),
        treasure = plugin_state.treasure_map,
        summary = getTreasureSummary()
    }

    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created treasure snapshot: %s", label))
        return snap_id
    end

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("snapshot_" .. label, snapshot)
        return {snapshot_id = "local_" .. label}
    end

    return {success = false}
end

function registerTreasureDatabaseAPI()
    load_phase11_dependencies()
    initialize_treasure_map()

    local api = dependencies.api_gateway

    if not api or not api.RESTInterface then return {success = false} end

    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/treasure/:id", function(params)
        return getTreasureById(tonumber(params.id))
    end)

    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/treasure/search", function(params)
        return getTreasureByLocation(params.q or "")
    end)

    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/treasure/summary", function()
        return getTreasureSummary()
    end)

    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end

    log_operation("API", "Registered treasure database REST endpoints")
    return {lookup = lookup_endpoint or {registered = true}, search = search_endpoint or {registered = true}, summary = summary_endpoint or {registered = true}}
end

-- Initialize on load
initialize_treasure_map()
log_operation("LOAD", "Treasure Database Plugin v1.0.0 loaded")

return {
    getTreasureById = getTreasureById,
    getTreasureByLocation = getTreasureByLocation,
    getTreasureByType = getTreasureByType,
    getTreasureByRarity = getTreasureByRarity,
    addTreasure = addTreasure,
    markTreasureOpened = markTreasureOpened,
    getTreasureSummary = getTreasureSummary,
    recordCollection = recordCollection,
    analyzeCollectionProgress = analyzeCollectionProgress,
    exportTreasureMap = exportTreasureMap,
    syncTreasureMapToHub = syncTreasureMapToHub,
    createTreasureSnapshot = createTreasureSnapshot,
    registerTreasureDatabaseAPI = registerTreasureDatabaseAPI
}
