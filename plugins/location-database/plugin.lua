--[[
  Location Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Tracks world locations with search, export, analytics, and integration

  Features:
  - Location lookup by ID, name, region, or type
  - Region and accessibility metadata
  - Analytics for visitation patterns
  - Export/Import location catalog
  - Integration Hub sync for cross-plugin map data
  - Backup/Restore snapshots

  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    LOCATION_TYPES = {
        TOWN = "town",
        DUNGEON = "dungeon",
        OVERWORLD = "overworld",
        TOWER = "tower",
        CAVE = "cave",
        RUIN = "ruin"
    },

    REGIONS = {
        WOB = "World of Balance",
        WOR = "World of Ruin"
    },

    MAX_LOCATIONS = 256,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    location_catalog = {},
    visit_history = {},
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

    print(string.format("[Location Database] %s: %s", operation_type, details))
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
-- LOCATION INITIALIZATION
-- ============================================================================

local function initialize_location_catalog()
    if plugin_state.initialized then return end

    plugin_state.location_catalog = {
        [0] = {id = 0, name = "Narshe", region = CONFIG.REGIONS.WOB, type = CONFIG.LOCATION_TYPES.TOWN, accessible = true},
        [5] = {id = 5, name = "Figaro Castle", region = CONFIG.REGIONS.WOB, type = CONFIG.LOCATION_TYPES.TOWN, accessible = true},
        [20] = {id = 20, name = "Mt. Kolts", region = CONFIG.REGIONS.WOB, type = CONFIG.LOCATION_TYPES.CAVE, accessible = true},
        [50] = {id = 50, name = "Opera House", region = CONFIG.REGIONS.WOB, type = CONFIG.LOCATION_TYPES.TOWN, accessible = true},
        [80] = {id = 80, name = "Vector", region = CONFIG.REGIONS.WOB, type = CONFIG.LOCATION_TYPES.TOWN, accessible = false},
        [120] = {id = 120, name = "Solitary Island", region = CONFIG.REGIONS.WOR, type = CONFIG.LOCATION_TYPES.TOWN, accessible = true},
        [140] = {id = 140, name = "Phoenix Cave", region = CONFIG.REGIONS.WOR, type = CONFIG.LOCATION_TYPES.CAVE, accessible = true},
        [160] = {id = 160, name = "Kefka's Tower", region = CONFIG.REGIONS.WOR, type = CONFIG.LOCATION_TYPES.TOWER, accessible = true}
    }

    plugin_state.initialized = true
    log_operation("INIT", string.format("Location database initialized with %d locations", table_count(plugin_state.location_catalog)))
end

-- ============================================================================
-- CORE LOCATION LOOKUP FUNCTIONS
-- ============================================================================

function getLocationById(location_id)
    initialize_location_catalog()

    if not location_id or location_id < 0 or location_id >= CONFIG.MAX_LOCATIONS then
        log_operation("ERROR", "Invalid location ID: " .. tostring(location_id))
        return nil
    end

    local location = plugin_state.location_catalog[location_id]
    if location then
        log_operation("LOOKUP", string.format("Retrieved location %d: %s", location_id, location.name))
    end

    return location
end

function searchLocationsByName(name_query)
    initialize_location_catalog()

    local results = {}
    local query_lower = string.lower(name_query or "")

    for location_id, location in pairs(plugin_state.location_catalog) do
        if string.find(string.lower(location.name), query_lower, 1, true) then
            table.insert(results, location)
        end
    end

    log_operation("SEARCH", string.format("Found %d locations matching '%s'", #results, name_query))
    return results
end

function getLocationsByRegion(region)
    initialize_location_catalog()

    local results = {}

    for location_id, location in pairs(plugin_state.location_catalog) do
        if location.region == region then
            table.insert(results, location)
        end
    end

    log_operation("FILTER", string.format("Found %d locations in region '%s'", #results, region))
    return results
end

function getLocationsByType(location_type)
    initialize_location_catalog()

    local results = {}

    for location_id, location in pairs(plugin_state.location_catalog) do
        if location.type == location_type then
            table.insert(results, location)
        end
    end

    log_operation("FILTER", string.format("Found %d locations of type '%s'", #results, location_type))
    return results
end

-- ============================================================================
-- LOCATION MANAGEMENT
-- ============================================================================

function addLocationToCatalog(location_data)
    initialize_location_catalog()

    if not location_data or not location_data.id or not location_data.name then
        log_operation("ERROR", "Invalid location data")
        return false
    end

    plugin_state.location_catalog[location_data.id] = location_data

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("location_catalog", plugin_state.location_catalog)
    end

    log_operation("ADD", string.format("Added location %d: %s", location_data.id, location_data.name))
    return true
end

function getLocationCatalogSummary()
    initialize_location_catalog()

    local summary = {
        total_locations = table_count(plugin_state.location_catalog),
        by_type = {},
        by_region = {}
    }

    for _, location in pairs(plugin_state.location_catalog) do
        summary.by_type[location.type] = (summary.by_type[location.type] or 0) + 1
        summary.by_region[location.region] = (summary.by_region[location.region] or 0) + 1
    end

    return summary
end

function recordVisit(location_id)
    initialize_location_catalog()

    plugin_state.visit_history[location_id] = (plugin_state.visit_history[location_id] or 0) + 1
    log_operation("VISIT", string.format("Recorded visit to location %d", location_id))
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

function analyzeVisitPatterns()
    load_phase11_dependencies()
    initialize_location_catalog()

    local analytics = dependencies.analytics

    local visit_data = {}
    for location_id, count in pairs(plugin_state.visit_history) do
        local location = plugin_state.location_catalog[location_id]
        if location then
            table.insert(visit_data, {
                location_id = location_id,
                name = location.name,
                count = count,
                region = location.region,
                type = location.type
            })
        end
    end

    local analysis = {
        total_locations = table_count(plugin_state.location_catalog),
        visits_tracked = #visit_data,
        patterns = {}
    }

    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(visit_data)
        analysis.patterns = patterns
    end

    log_operation("ANALYTICS", string.format("Analyzed visit patterns for %d locations", analysis.visits_tracked))
    return analysis
end

function exportLocationCatalog(format, path)
    load_phase11_dependencies()
    initialize_location_catalog()

    local exporter = dependencies.import_export and dependencies.import_export.DataExporter

    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        location_count = table_count(plugin_state.location_catalog),
        locations = plugin_state.location_catalog
    }

    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("location_catalog_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)

        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end

        log_operation("EXPORT", string.format("Exported location catalog to %s", output))
        return {path = output, format = fmt}
    end

    return {success = false, error = "Import/Export unavailable"}
end

function syncLocationCatalogToHub()
    load_phase11_dependencies()
    initialize_location_catalog()

    local hub = dependencies.integration_hub

    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("location_catalog_sync", {
            catalog = plugin_state.location_catalog,
            summary = getLocationCatalogSummary(),
            timestamp = os.time()
        })

        log_operation("SYNC_HUB", "Synced location catalog to Integration Hub")
        return result or {success = true}
    end

    return {success = false, error = "Integration Hub unavailable"}
end

function createLocationCatalogSnapshot(label)
    load_phase11_dependencies()
    initialize_location_catalog()

    local backup = dependencies.backup_restore

    local snapshot = {
        label = label or "location_catalog_snapshot",
        timestamp = os.time(),
        catalog = plugin_state.location_catalog,
        summary = getLocationCatalogSummary()
    }

    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created location catalog snapshot: %s", label))
        return snap_id
    end

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("snapshot_" .. label, snapshot)
        return {snapshot_id = "local_" .. label}
    end

    return {success = false}
end

function registerLocationDatabaseAPI()
    load_phase11_dependencies()
    initialize_location_catalog()

    local api = dependencies.api_gateway

    if not api or not api.RESTInterface then return {success = false} end

    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/locations/:id", function(params)
        return getLocationById(tonumber(params.id))
    end)

    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/locations/search", function(params)
        return searchLocationsByName(params.q or "")
    end)

    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/locations/summary", function()
        return getLocationCatalogSummary()
    end)

    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end

    log_operation("API", "Registered location database REST endpoints")
    return {lookup = lookup_endpoint or {registered = true}, search = search_endpoint or {registered = true}, summary = summary_endpoint or {registered = true}}
end

-- Initialize on load
initialize_location_catalog()
log_operation("LOAD", "Location Database Plugin v1.0.0 loaded")

return {
    getLocationById = getLocationById,
    searchLocationsByName = searchLocationsByName,
    getLocationsByRegion = getLocationsByRegion,
    getLocationsByType = getLocationsByType,
    addLocationToCatalog = addLocationToCatalog,
    getLocationCatalogSummary = getLocationCatalogSummary,
    recordVisit = recordVisit,
    analyzeVisitPatterns = analyzeVisitPatterns,
    exportLocationCatalog = exportLocationCatalog,
    syncLocationCatalogToHub = syncLocationCatalogToHub,
    createLocationCatalogSnapshot = createLocationCatalogSnapshot,
    registerLocationDatabaseAPI = registerLocationDatabaseAPI
}
