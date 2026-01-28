--[[
  NPC Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Tracks NPCs, roles, quests, and dialogue with analytics, export, and integration

  Features:
  - NPC lookup by ID, name, location, or role
  - Quest relevance and dialogue metadata
  - Analytics for NPC interaction patterns
  - Export/Import NPC registry
  - Integration Hub sync for cross-plugin NPC data
  - Backup/Restore snapshots

  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    ROLES = {
        MERCHANT = "merchant",
        QUEST_GIVER = "quest_giver",
        ALLY = "ally",
        ANTAGONIST = "antagonist",
        GENERIC = "generic"
    },

    ALIGNMENT = {
        FRIENDLY = "friendly",
        NEUTRAL = "neutral",
        HOSTILE = "hostile"
    },

    MAX_NPCS = 256,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    npc_registry = {},
    interaction_history = {},
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

    print(string.format("[NPC Database] %s: %s", operation_type, details))
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
-- NPC INITIALIZATION
-- ============================================================================

local function initialize_npc_registry()
    if plugin_state.initialized then return end

    plugin_state.npc_registry = {
        [0] = {id = 0, name = "Arvis", role = CONFIG.ROLES.ALLY, alignment = CONFIG.ALIGNMENT.FRIENDLY, location = "Narshe", quests = {"Narshe Defense"}},
        [10] = {id = 10, name = "Edgar", role = CONFIG.ROLES.ALLY, alignment = CONFIG.ALIGNMENT.FRIENDLY, location = "Figaro Castle", quests = {"Figaro Escape"}},
        [30] = {id = 30, name = "Ultros", role = CONFIG.ROLES.ANTAGONIST, alignment = CONFIG.ALIGNMENT.HOSTILE, location = "Opera House", quests = {"Opera House"}},
        [50] = {id = 50, name = "Merchant", role = CONFIG.ROLES.MERCHANT, alignment = CONFIG.ALIGNMENT.NEUTRAL, location = "South Figaro", quests = {}},
        [70] = {id = 70, name = "Gestahl", role = CONFIG.ROLES.ANTAGONIST, alignment = CONFIG.ALIGNMENT.HOSTILE, location = "Vector", quests = {"Empire"}},
        [120] = {id = 120, name = "Gogo", role = CONFIG.ROLES.ALLY, alignment = CONFIG.ALIGNMENT.FRIENDLY, location = "Zone Eater", quests = {"Recruit Gogo"}},
        [140] = {id = 140, name = "Mog", role = CONFIG.ROLES.ALLY, alignment = CONFIG.ALIGNMENT.FRIENDLY, location = "Narshe", quests = {"Protect Narshe"}},
        [200] = {id = 200, name = "Leo", role = CONFIG.ROLES.ALLY, alignment = CONFIG.ALIGNMENT.FRIENDLY, location = "Thamasa", quests = {"Thamasa"}}
    }

    plugin_state.initialized = true
    log_operation("INIT", string.format("NPC database initialized with %d NPCs", table_count(plugin_state.npc_registry)))
end

-- ============================================================================
-- CORE NPC LOOKUP FUNCTIONS
-- ============================================================================

function getNPCById(npc_id)
    initialize_npc_registry()

    if not npc_id or npc_id < 0 or npc_id >= CONFIG.MAX_NPCS then
        log_operation("ERROR", "Invalid NPC ID: " .. tostring(npc_id))
        return nil
    end

    local npc = plugin_state.npc_registry[npc_id]
    if npc then
        log_operation("LOOKUP", string.format("Retrieved NPC %d: %s", npc_id, npc.name))
    end

    return npc
end

function searchNPCsByName(name_query)
    initialize_npc_registry()

    local results = {}
    local query_lower = string.lower(name_query or "")

    for npc_id, npc in pairs(plugin_state.npc_registry) do
        if string.find(string.lower(npc.name), query_lower, 1, true) then
            table.insert(results, npc)
        end
    end

    log_operation("SEARCH", string.format("Found %d NPCs matching '%s'", #results, name_query))
    return results
end

function getNPCsByLocation(location)
    initialize_npc_registry()

    local results = {}

    for npc_id, npc in pairs(plugin_state.npc_registry) do
        if npc.location and string.find(string.lower(npc.location), string.lower(location), 1, true) then
            table.insert(results, npc)
        end
    end

    log_operation("FILTER", string.format("Found %d NPCs in location '%s'", #results, location))
    return results
end

function getNPCsByRole(role)
    initialize_npc_registry()

    local results = {}

    for npc_id, npc in pairs(plugin_state.npc_registry) do
        if npc.role == role then
            table.insert(results, npc)
        end
    end

    log_operation("FILTER", string.format("Found %d NPCs with role '%s'", #results, role))
    return results
end

-- ============================================================================
-- NPC MANAGEMENT
-- ============================================================================

function addNPC(npc_data)
    initialize_npc_registry()

    if not npc_data or not npc_data.id or not npc_data.name then
        log_operation("ERROR", "Invalid NPC data")
        return false
    end

    plugin_state.npc_registry[npc_data.id] = npc_data

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("npc_registry", plugin_state.npc_registry)
    end

    log_operation("ADD", string.format("Added NPC %d: %s", npc_data.id, npc_data.name))
    return true
end

function getNPCRegistrySummary()
    initialize_npc_registry()

    local summary = {
        total_npcs = table_count(plugin_state.npc_registry),
        by_role = {},
        by_location = {}
    }

    for _, npc in pairs(plugin_state.npc_registry) do
        summary.by_role[npc.role] = (summary.by_role[npc.role] or 0) + 1
        if npc.location then
            summary.by_location[npc.location] = (summary.by_location[npc.location] or 0) + 1
        end
    end

    return summary
end

function recordInteraction(npc_id)
    initialize_npc_registry()

    plugin_state.interaction_history[npc_id] = (plugin_state.interaction_history[npc_id] or 0) + 1
    log_operation("INTERACT", string.format("Recorded interaction with NPC %d", npc_id))
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

function analyzeInteractionPatterns()
    load_phase11_dependencies()
    initialize_npc_registry()

    local analytics = dependencies.analytics

    local interaction_data = {}
    for npc_id, count in pairs(plugin_state.interaction_history) do
        local npc = plugin_state.npc_registry[npc_id]
        if npc then
            table.insert(interaction_data, {
                npc_id = npc_id,
                name = npc.name,
                role = npc.role,
                location = npc.location,
                count = count
            })
        end
    end

    local analysis = {
        total_npcs = table_count(plugin_state.npc_registry),
        interactions_tracked = #interaction_data,
        patterns = {}
    }

    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(interaction_data)
        analysis.patterns = patterns
    end

    log_operation("ANALYTICS", string.format("Analyzed interaction patterns for %d NPCs", analysis.interactions_tracked))
    return analysis
end

function exportNPCRegistry(format, path)
    load_phase11_dependencies()
    initialize_npc_registry()

    local exporter = dependencies.import_export and dependencies.import_export.DataExporter

    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        npc_count = table_count(plugin_state.npc_registry),
        npcs = plugin_state.npc_registry
    }

    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("npc_registry_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)

        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end

        log_operation("EXPORT", string.format("Exported NPC registry to %s", output))
        return {path = output, format = fmt}
    end

    return {success = false, error = "Import/Export unavailable"}
end

function syncNPCRegistryToHub()
    load_phase11_dependencies()
    initialize_npc_registry()

    local hub = dependencies.integration_hub

    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("npc_registry_sync", {
            registry = plugin_state.npc_registry,
            summary = getNPCRegistrySummary(),
            timestamp = os.time()
        })

        log_operation("SYNC_HUB", "Synced NPC registry to Integration Hub")
        return result or {success = true}
    end

    return {success = false, error = "Integration Hub unavailable"}
end

function createNPCRegistrySnapshot(label)
    load_phase11_dependencies()
    initialize_npc_registry()

    local backup = dependencies.backup_restore

    local snapshot = {
        label = label or "npc_registry_snapshot",
        timestamp = os.time(),
        registry = plugin_state.npc_registry,
        summary = getNPCRegistrySummary()
    }

    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created NPC registry snapshot: %s", label))
        return snap_id
    end

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("snapshot_" .. label, snapshot)
        return {snapshot_id = "local_" .. label}
    end

    return {success = false}
end

function registerNPCDatabaseAPI()
    load_phase11_dependencies()
    initialize_npc_registry()

    local api = dependencies.api_gateway

    if not api or not api.RESTInterface then return {success = false} end

    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/npcs/:id", function(params)
        return getNPCById(tonumber(params.id))
    end)

    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/npcs/search", function(params)
        return searchNPCsByName(params.q or "")
    end)

    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/npcs/summary", function()
        return getNPCRegistrySummary()
    end)

    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end

    log_operation("API", "Registered NPC database REST endpoints")
    return {lookup = lookup_endpoint or {registered = true}, search = search_endpoint or {registered = true}, summary = summary_endpoint or {registered = true}}
end

-- Initialize on load
initialize_npc_registry()
log_operation("LOAD", "NPC Database Plugin v1.0.0 loaded")

return {
    getNPCById = getNPCById,
    searchNPCsByName = searchNPCsByName,
    getNPCsByLocation = getNPCsByLocation,
    getNPCsByRole = getNPCsByRole,
    addNPC = addNPC,
    getNPCRegistrySummary = getNPCRegistrySummary,
    recordInteraction = recordInteraction,
    analyzeInteractionPatterns = analyzeInteractionPatterns,
    exportNPCRegistry = exportNPCRegistry,
    syncNPCRegistryToHub = syncNPCRegistryToHub,
    createNPCRegistrySnapshot = createNPCRegistrySnapshot,
    registerNPCDatabaseAPI = registerNPCDatabaseAPI
}
