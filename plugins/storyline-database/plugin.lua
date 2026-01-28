--[[
  Storyline Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Tracks story events, chapters, and flags with analytics, export, and integration

  Features:
  - Story event lookup by ID, name, chapter, or status
  - Story flags and progression metadata
  - Analytics for storyline coverage
  - Export/Import storyline data
  - Integration Hub sync for cross-plugin story state
  - Backup/Restore snapshots

  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    EVENT_TYPES = {
        MAIN = "main",
        SIDE = "side",
        OPTIONAL = "optional"
    },

    STATUS = {
        PENDING = "pending",
        IN_PROGRESS = "in_progress",
        COMPLETED = "completed",
        BLOCKED = "blocked"
    },

    MAX_EVENTS = 256,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    events = {},
    flags = {},
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

    print(string.format("[Storyline Database] %s: %s", operation_type, details))
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
-- STORYLINE INITIALIZATION
-- ============================================================================

local function initialize_storyline()
    if plugin_state.initialized then return end

    plugin_state.events = {
        [0] = {id = 0, name = "Narshe Defense", chapter = "Opening", type = CONFIG.EVENT_TYPES.MAIN, status = CONFIG.STATUS.COMPLETED, location = "Narshe"},
        [10] = {id = 10, name = "Figaro Escape", chapter = "Opening", type = CONFIG.EVENT_TYPES.MAIN, status = CONFIG.STATUS.COMPLETED, location = "Figaro"},
        [50] = {id = 50, name = "Opera House", chapter = "Midgame", type = CONFIG.EVENT_TYPES.MAIN, status = CONFIG.STATUS.COMPLETED, location = "Jidoor"},
        [80] = {id = 80, name = "Floating Continent", chapter = "Midgame", type = CONFIG.EVENT_TYPES.MAIN, status = CONFIG.STATUS.COMPLETED, location = "Floating Continent"},
        [120] = {id = 120, name = "World of Ruin Awakening", chapter = "WOR", type = CONFIG.EVENT_TYPES.MAIN, status = CONFIG.STATUS.IN_PROGRESS, location = "Solitary Island"},
        [140] = {id = 140, name = "Phoenix Cave", chapter = "WOR", type = CONFIG.EVENT_TYPES.SIDE, status = CONFIG.STATUS.PENDING, location = "Phoenix Cave"},
        [160] = {id = 160, name = "Ancient Castle", chapter = "WOR", type = CONFIG.EVENT_TYPES.OPTIONAL, status = CONFIG.STATUS.PENDING, location = "Ancient Castle"},
        [200] = {id = 200, name = "Kefka's Tower", chapter = "Finale", type = CONFIG.EVENT_TYPES.MAIN, status = CONFIG.STATUS.PENDING, location = "Kefka's Tower"}
    }

    plugin_state.flags = {
        world_state = "WOR",
        opera_completed = true,
        floating_continent_cleared = true,
        party_split = false
    }

    plugin_state.initialized = true
    log_operation("INIT", string.format("Storyline database initialized with %d events", table_count(plugin_state.events)))
end

-- ============================================================================
-- CORE STORYLINE LOOKUP FUNCTIONS
-- ============================================================================

function getStoryEventById(event_id)
    initialize_storyline()

    if not event_id or event_id < 0 or event_id >= CONFIG.MAX_EVENTS then
        log_operation("ERROR", "Invalid event ID: " .. tostring(event_id))
        return nil
    end

    local event = plugin_state.events[event_id]
    if event then
        log_operation("LOOKUP", string.format("Retrieved story event %d: %s", event_id, event.name))
    end

    return event
end

function searchStoryEventsByName(name_query)
    initialize_storyline()

    local results = {}
    local query_lower = string.lower(name_query or "")

    for event_id, event in pairs(plugin_state.events) do
        if string.find(string.lower(event.name), query_lower, 1, true) then
            table.insert(results, event)
        end
    end

    log_operation("SEARCH", string.format("Found %d story events matching '%s'", #results, name_query))
    return results
end

function getEventsByChapter(chapter)
    initialize_storyline()

    local results = {}

    for event_id, event in pairs(plugin_state.events) do
        if event.chapter == chapter then
            table.insert(results, event)
        end
    end

    log_operation("FILTER", string.format("Found %d story events in chapter '%s'", #results, chapter))
    return results
end

function getEventsByStatus(status)
    initialize_storyline()

    local results = {}

    for event_id, event in pairs(plugin_state.events) do
        if event.status == status then
            table.insert(results, event)
        end
    end

    log_operation("FILTER", string.format("Found %d story events with status '%s'", #results, status))
    return results
end

function getEventsByType(event_type)
    initialize_storyline()

    local results = {}

    for event_id, event in pairs(plugin_state.events) do
        if event.type == event_type then
            table.insert(results, event)
        end
    end

    log_operation("FILTER", string.format("Found %d story events of type '%s'", #results, event_type))
    return results
end

-- ============================================================================
-- STORYLINE MANAGEMENT
-- ============================================================================

function addStoryEvent(event_data)
    initialize_storyline()

    if not event_data or not event_data.id or not event_data.name then
        log_operation("ERROR", "Invalid story event data")
        return false
    end

    plugin_state.events[event_data.id] = event_data

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("storyline_events", plugin_state.events)
    end

    log_operation("ADD", string.format("Added story event %d: %s", event_data.id, event_data.name))
    return true
end

function updateStoryFlag(flag_key, value)
    initialize_storyline()

    if not flag_key then
        log_operation("ERROR", "Invalid flag key")
        return false
    end

    plugin_state.flags[flag_key] = value

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("storyline_flags", plugin_state.flags)
    end

    log_operation("FLAG", string.format("Updated flag '%s'", flag_key))
    return true
end

function getStoryProgressSummary()
    initialize_storyline()

    local summary = {
        total_events = table_count(plugin_state.events),
        by_status = {},
        by_type = {},
        flags = plugin_state.flags
    }

    for _, event in pairs(plugin_state.events) do
        summary.by_status[event.status] = (summary.by_status[event.status] or 0) + 1
        summary.by_type[event.type] = (summary.by_type[event.type] or 0) + 1
    end

    return summary
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

function analyzeStorylineCoverage()
    load_phase11_dependencies()
    initialize_storyline()

    local analytics = dependencies.analytics

    local coverage_data = {}
    for event_id, event in pairs(plugin_state.events) do
        table.insert(coverage_data, {
            event_id = event_id,
            status = event.status,
            type = event.type,
            chapter = event.chapter
        })
    end

    local analysis = {
        total_events = table_count(plugin_state.events),
        patterns = {}
    }

    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(coverage_data)
        analysis.patterns = patterns
    end

    log_operation("ANALYTICS", string.format("Analyzed storyline coverage for %d events", analysis.total_events))
    return analysis
end

function exportStorylineData(format, path)
    load_phase11_dependencies()
    initialize_storyline()

    local exporter = dependencies.import_export and dependencies.import_export.DataExporter

    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        event_count = table_count(plugin_state.events),
        events = plugin_state.events,
        flags = plugin_state.flags
    }

    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("storyline_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)

        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end

        log_operation("EXPORT", string.format("Exported storyline data to %s", output))
        return {path = output, format = fmt}
    end

    return {success = false, error = "Import/Export unavailable"}
end

function syncStorylineToHub()
    load_phase11_dependencies()
    initialize_storyline()

    local hub = dependencies.integration_hub

    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("storyline_sync", {
            events = plugin_state.events,
            flags = plugin_state.flags,
            summary = getStoryProgressSummary(),
            timestamp = os.time()
        })

        log_operation("SYNC_HUB", "Synced storyline data to Integration Hub")
        return result or {success = true}
    end

    return {success = false, error = "Integration Hub unavailable"}
end

function createStorylineSnapshot(label)
    load_phase11_dependencies()
    initialize_storyline()

    local backup = dependencies.backup_restore

    local snapshot = {
        label = label or "storyline_snapshot",
        timestamp = os.time(),
        events = plugin_state.events,
        flags = plugin_state.flags,
        summary = getStoryProgressSummary()
    }

    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created storyline snapshot: %s", label))
        return snap_id
    end

    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("snapshot_" .. label, snapshot)
        return {snapshot_id = "local_" .. label}
    end

    return {success = false}
end

function registerStorylineDatabaseAPI()
    load_phase11_dependencies()
    initialize_storyline()

    local api = dependencies.api_gateway

    if not api or not api.RESTInterface then return {success = false} end

    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/storyline/events/:id", function(params)
        return getStoryEventById(tonumber(params.id))
    end)

    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/storyline/events/search", function(params)
        return searchStoryEventsByName(params.q or "")
    end)

    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/storyline/summary", function()
        return getStoryProgressSummary()
    end)

    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end

    log_operation("API", "Registered storyline database REST endpoints")
    return {lookup = lookup_endpoint or {registered = true}, search = search_endpoint or {registered = true}, summary = summary_endpoint or {registered = true}}
end

-- Initialize on load
initialize_storyline()
log_operation("LOAD", "Storyline Database Plugin v1.0.0 loaded")

return {
    getStoryEventById = getStoryEventById,
    searchStoryEventsByName = searchStoryEventsByName,
    getEventsByChapter = getEventsByChapter,
    getEventsByStatus = getEventsByStatus,
    getEventsByType = getEventsByType,
    addStoryEvent = addStoryEvent,
    updateStoryFlag = updateStoryFlag,
    getStoryProgressSummary = getStoryProgressSummary,
    analyzeStorylineCoverage = analyzeStorylineCoverage,
    exportStorylineData = exportStorylineData,
    syncStorylineToHub = syncStorylineToHub,
    createStorylineSnapshot = createStorylineSnapshot,
    registerStorylineDatabaseAPI = registerStorylineDatabaseAPI
}
