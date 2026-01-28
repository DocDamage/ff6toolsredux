--[[
  Ability Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Searchable ability catalog with analytics, export, and cross-plugin integration

  Features:
  - Ability lookup by ID, name, type, or element
  - MP cost and power metadata
  - Usage analytics
  - Export/Import ability catalog
  - Integration Hub sync for cross-plugin ability data
  - Backup/Restore snapshots

  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    ABILITY_TYPES = {
        MAGIC = "magic",
        BLITZ = "blitz",
        BUSHIDO = "bushido",
        LORE = "lore",
        TOOL = "tool",
        ESPER = "esper",
        RAGE = "rage"
    },

    ELEMENTS = {
        FIRE = "fire",
        ICE = "ice",
        LIGHTNING = "lightning",
        HOLY = "holy",
        EARTH = "earth",
        WIND = "wind",
        POISON = "poison",
        WATER = "water",
        NEUTRAL = "neutral"
    },

    MAX_ABILITIES = 256,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    ability_catalog = {},
    usage_stats = {},
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

    print(string.format("[Ability Database] %s: %s", operation_type, details))
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
-- ABILITY DATABASE INITIALIZATION
-- ============================================================================

local function initialize_ability_database()
    if plugin_state.initialized then return end

    plugin_state.ability_catalog = {
        [0] = {id = 0, name = "Fire", type = CONFIG.ABILITY_TYPES.MAGIC, element = CONFIG.ELEMENTS.FIRE, mp_cost = 4, power = 22},
        [1] = {id = 1, name = "Cura", type = CONFIG.ABILITY_TYPES.MAGIC, element = CONFIG.ELEMENTS.HOLY, mp_cost = 6, power = 45, effect = "Heal"},
        [10] = {id = 10, name = "Thundara", type = CONFIG.ABILITY_TYPES.MAGIC, element = CONFIG.ELEMENTS.LIGHTNING, mp_cost = 12, power = 55},
        [40] = {id = 40, name = "Suplex", type = CONFIG.ABILITY_TYPES.BLITZ, element = CONFIG.ELEMENTS.EARTH, mp_cost = 0, power = 70},
        [50] = {id = 50, name = "Bum Rush", type = CONFIG.ABILITY_TYPES.BLITZ, element = CONFIG.ELEMENTS.WIND, mp_cost = 0, power = 110},
        [75] = {id = 75, name = "Shock", type = CONFIG.ABILITY_TYPES.ESPER, element = CONFIG.ELEMENTS.LIGHTNING, mp_cost = 25, power = 90},
        [90] = {id = 90, name = "Chainsaw", type = CONFIG.ABILITY_TYPES.TOOL, element = CONFIG.ELEMENTS.NEUTRAL, mp_cost = 0, power = 85},
        [120] = {id = 120, name = "Aqua Rake", type = CONFIG.ABILITY_TYPES.LORE, element = CONFIG.ELEMENTS.WATER, mp_cost = 20, power = 75},
        [200] = {id = 200, name = "Ultima", type = CONFIG.ABILITY_TYPES.MAGIC, element = CONFIG.ELEMENTS.NEUTRAL, mp_cost = 80, power = 255},
        [220] = {id = 220, name = "Ragnarok", type = CONFIG.ABILITY_TYPES.ESPER, element = CONFIG.ELEMENTS.HOLY, mp_cost = 70, power = 180}
    }

    plugin_state.initialized = true
    log_operation("INIT", string.format("Ability database initialized with %d abilities", table_count(plugin_state.ability_catalog)))
end

-- ============================================================================
-- CORE ABILITY LOOKUP FUNCTIONS
-- ============================================================================

function getAbilityById(ability_id)
    initialize_ability_database()

    if not ability_id or ability_id < 0 or ability_id >= CONFIG.MAX_ABILITIES then
        log_operation("ERROR", "Invalid ability ID: " .. tostring(ability_id))
        return nil
    end

    local ability = plugin_state.ability_catalog[ability_id]
    if ability then
        log_operation("LOOKUP", string.format("Retrieved ability %d: %s", ability_id, ability.name))
    end

    return ability
end

function searchAbilitiesByName(name_query)
    initialize_ability_database()

    local results = {}
    local query_lower = string.lower(name_query or "")

    for ability_id, ability in pairs(plugin_state.ability_catalog) do
        if string.find(string.lower(ability.name), query_lower, 1, true) then
            table.insert(results, ability)
        end
    end

    log_operation("SEARCH", string.format("Found %d abilities matching '%s'", #results, name_query))
    return results
end

function getAbilitiesByType(ability_type)
    initialize_ability_database()

    local results = {}

    for ability_id, ability in pairs(plugin_state.ability_catalog) do
        if ability.type == ability_type then
            table.insert(results, ability)
        end
    end

    log_operation("FILTER", string.format("Found %d abilities of type '%s'", #results, ability_type))
    return results
end

function getAbilitiesByElement(element)
    initialize_ability_database()

    local results = {}

    for ability_id, ability in pairs(plugin_state.ability_catalog) do
        if ability.element == element then
            table.insert(results, ability)
        end
    end

    log_operation("FILTER", string.format("Found %d abilities with element '%s'", #results, element))
    return results
end

function getAbilitiesByCostRange(min_cost, max_cost)
    initialize_ability_database()

    local results = {}
    local minc = min_cost or 0
    local maxc = max_cost or 999

    for ability_id, ability in pairs(plugin_state.ability_catalog) do
        local cost = ability.mp_cost or 0
        if cost >= minc and cost <= maxc then
            table.insert(results, ability)
        end
    end

    log_operation("FILTER", string.format("Found %d abilities within cost %d-%d", #results, minc, maxc))
    return results
end

-- ============================================================================
-- ABILITY CATALOG MANAGEMENT
-- ============================================================================

function addAbilityToCatalog(ability_data)
    initialize_ability_database()

    if not ability_data or not ability_data.id or not ability_data.name then
        log_operation("ERROR", "Invalid ability data")
        return false
    end

    plugin_state.ability_catalog[ability_data.id] = ability_data

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("ability_catalog", plugin_state.ability_catalog)
    end

    log_operation("ADD", string.format("Added ability %d: %s", ability_data.id, ability_data.name))
    return true
end

function getAbilityCatalogSummary()
    initialize_ability_database()

    local summary = {
        total_abilities = table_count(plugin_state.ability_catalog),
        by_type = {},
        by_element = {}
    }

    for _, ability in pairs(plugin_state.ability_catalog) do
        summary.by_type[ability.type] = (summary.by_type[ability.type] or 0) + 1
        summary.by_element[ability.element] = (summary.by_element[ability.element] or 0) + 1
    end

    return summary
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

function analyzeAbilityUsagePatterns()
    load_phase11_dependencies()
    initialize_ability_database()

    local analytics = dependencies.analytics

    local usage_data = {}
    for ability_id, count in pairs(plugin_state.usage_stats) do
        table.insert(usage_data, {ability_id = ability_id, count = count})
    end

    local analysis = {
        total_abilities = table_count(plugin_state.ability_catalog),
        usage_tracked = #usage_data,
        patterns = {}
    }

    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(usage_data)
        analysis.patterns = patterns
    end

    log_operation("ANALYTICS", string.format("Analyzed usage patterns for %d abilities", analysis.usage_tracked))
    return analysis
end

function exportAbilityCatalog(format, path)
    load_phase11_dependencies()
    initialize_ability_database()

    local exporter = dependencies.import_export and dependencies.import_export.DataExporter

    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        ability_count = table_count(plugin_state.ability_catalog),
        abilities = plugin_state.ability_catalog
    }

    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("ability_catalog_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)

        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end

        log_operation("EXPORT", string.format("Exported ability catalog to %s", output))
        return {path = output, format = fmt}
    end

    return {success = false, error = "Import/Export unavailable"}
end

function syncAbilityCatalogToHub()
    load_phase11_dependencies()
    initialize_ability_database()

    local hub = dependencies.integration_hub

    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("ability_catalog_sync", {
            catalog = plugin_state.ability_catalog,
            summary = getAbilityCatalogSummary(),
            timestamp = os.time()
        })

        log_operation("SYNC_HUB", "Synced ability catalog to Integration Hub")
        return result or {success = true}
    end

    return {success = false, error = "Integration Hub unavailable"}
end

function createAbilityCatalogSnapshot(label)
    load_phase11_dependencies()
    initialize_ability_database()

    local backup = dependencies.backup_restore

    local snapshot = {
        label = label or "ability_catalog_snapshot",
        timestamp = os.time(),
        catalog = plugin_state.ability_catalog,
        summary = getAbilityCatalogSummary()
    }

    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created ability catalog snapshot: %s", label))
        return snap_id
    end

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("snapshot_" .. label, snapshot)
        return {snapshot_id = "local_" .. label}
    end

    return {success = false}
end

function registerAbilityDatabaseAPI()
    load_phase11_dependencies()
    initialize_ability_database()

    local api = dependencies.api_gateway

    if not api or not api.RESTInterface then return {success = false} end

    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/abilities/:id", function(params)
        return getAbilityById(tonumber(params.id))
    end)

    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/abilities/search", function(params)
        return searchAbilitiesByName(params.q or "")
    end)

    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/abilities/summary", function()
        return getAbilityCatalogSummary()
    end)

    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end

    log_operation("API", "Registered ability database REST endpoints")
    return {lookup = lookup_endpoint or {registered = true}, search = search_endpoint or {registered = true}, summary = summary_endpoint or {registered = true}}
end

-- Initialize on load
initialize_ability_database()
log_operation("LOAD", "Ability Database Plugin v1.0.0 loaded")

return {
    getAbilityById = getAbilityById,
    searchAbilitiesByName = searchAbilitiesByName,
    getAbilitiesByType = getAbilitiesByType,
    getAbilitiesByElement = getAbilitiesByElement,
    getAbilitiesByCostRange = getAbilitiesByCostRange,
    addAbilityToCatalog = addAbilityToCatalog,
    getAbilityCatalogSummary = getAbilityCatalogSummary,
    analyzeAbilityUsagePatterns = analyzeAbilityUsagePatterns,
    exportAbilityCatalog = exportAbilityCatalog,
    syncAbilityCatalogToHub = syncAbilityCatalogToHub,
    createAbilityCatalogSnapshot = createAbilityCatalogSnapshot,
    registerAbilityDatabaseAPI = registerAbilityDatabaseAPI
}
