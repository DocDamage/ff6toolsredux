--[[
  Item Database Plugin v1.0.0 (Tier 2 Phase 1 - Database Suite)
  Provides searchable item database with analytics, export, and cross-plugin integration
  
  Features:
  - Item lookup by ID, name, type, or rarity
  - Item stats and metadata tracking
  - Analytics for item usage patterns
  - Export/Import item catalogs
  - Integration Hub sync for cross-plugin item data
  - Backup/Restore for item collections
  
  Phase: Tier 2 Phase 1 (Database Integration Foundation)
  Version: 1.0.0
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Item categories
    ITEM_TYPES = {
        WEAPON = "weapon",
        ARMOR = "armor",
        RELIC = "relic",
        CONSUMABLE = "consumable",
        KEY_ITEM = "key_item",
        TOOL = "tool"
    },
    
    -- Item rarity
    RARITY = {
        COMMON = 1,
        UNCOMMON = 2,
        RARE = 3,
        EPIC = 4,
        LEGENDARY = 5
    },
    
    MAX_ITEMS = 256,
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    item_catalog = {},
    search_cache = {},
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
    
    print(string.format("[Item Database] %s: %s", operation_type, details))
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
-- ITEM DATABASE INITIALIZATION
-- ============================================================================

local function initialize_item_database()
    if plugin_state.initialized then return end
    
    -- Sample item catalog (in real implementation, load from game data)
    plugin_state.item_catalog = {
        -- Weapons
        [0] = {id = 0, name = "Dirk", type = CONFIG.ITEM_TYPES.WEAPON, rarity = CONFIG.RARITY.COMMON, attack = 12, price = 100},
        [1] = {id = 1, name = "MithrilKnife", type = CONFIG.ITEM_TYPES.WEAPON, rarity = CONFIG.RARITY.UNCOMMON, attack = 30, price = 300},
        [255] = {id = 255, name = "Ultima Weapon", type = CONFIG.ITEM_TYPES.WEAPON, rarity = CONFIG.RARITY.LEGENDARY, attack = 255, magic_power = 108, price = 0},
        
        -- Armor
        [100] = {id = 100, name = "Leather Hat", type = CONFIG.ITEM_TYPES.ARMOR, rarity = CONFIG.RARITY.COMMON, defense = 10, price = 50},
        [150] = {id = 150, name = "Force Armor", type = CONFIG.ITEM_TYPES.ARMOR, rarity = CONFIG.RARITY.EPIC, defense = 70, magic_defense = 50, price = 5000},
        
        -- Relics
        [200] = {id = 200, name = "Sprint Shoes", type = CONFIG.ITEM_TYPES.RELIC, rarity = CONFIG.RARITY.UNCOMMON, effect = "Permanent Haste", price = 1500},
        [225] = {id = 225, name = "Ribbon", type = CONFIG.ITEM_TYPES.RELIC, rarity = CONFIG.RARITY.LEGENDARY, effect = "Immunity to all status", price = 0},
        
        -- Consumables
        [230] = {id = 230, name = "Potion", type = CONFIG.ITEM_TYPES.CONSUMABLE, rarity = CONFIG.RARITY.COMMON, effect = "Restore 50 HP", price = 50},
        [240] = {id = 240, name = "Elixir", type = CONFIG.ITEM_TYPES.CONSUMABLE, rarity = CONFIG.RARITY.RARE, effect = "Restore all HP/MP", price = 5000}
    }
    
    plugin_state.initialized = true
    log_operation("INIT", string.format("Item database initialized with %d items", table_count(plugin_state.item_catalog)))
end

-- ============================================================================
-- CORE ITEM LOOKUP FUNCTIONS
-- ============================================================================

function getItemById(item_id)
    initialize_item_database()
    
    if not item_id or item_id < 0 or item_id >= CONFIG.MAX_ITEMS then
        log_operation("ERROR", "Invalid item ID: " .. tostring(item_id))
        return nil
    end
    
    local item = plugin_state.item_catalog[item_id]
    if item then
        log_operation("LOOKUP", string.format("Retrieved item %d: %s", item_id, item.name))
    end
    
    return item
end

function searchItemsByName(name_query)
    initialize_item_database()
    
    local results = {}
    local query_lower = string.lower(name_query or "")
    
    for item_id, item in pairs(plugin_state.item_catalog) do
        if string.find(string.lower(item.name), query_lower, 1, true) then
            table.insert(results, item)
        end
    end
    
    log_operation("SEARCH", string.format("Found %d items matching '%s'", #results, name_query))
    return results
end

function getItemsByType(item_type)
    initialize_item_database()
    
    local results = {}
    
    for item_id, item in pairs(plugin_state.item_catalog) do
        if item.type == item_type then
            table.insert(results, item)
        end
    end
    
    log_operation("FILTER", string.format("Found %d items of type '%s'", #results, item_type))
    return results
end

function getItemsByRarity(rarity)
    initialize_item_database()
    
    local results = {}
    
    for item_id, item in pairs(plugin_state.item_catalog) do
        if item.rarity == rarity then
            table.insert(results, item)
        end
    end
    
    log_operation("FILTER", string.format("Found %d items with rarity %d", #results, rarity))
    return results
end

-- ============================================================================
-- ITEM CATALOG MANAGEMENT
-- ============================================================================

function addItemToCatalog(item_data)
    initialize_item_database()
    
    if not item_data or not item_data.id or not item_data.name then
        log_operation("ERROR", "Invalid item data")
        return false
    end
    
    plugin_state.item_catalog[item_data.id] = item_data
    
    -- Persist to database layer
    local db = load_database_layer()
    if db and db.savePersistentData then
        db.savePersistentData("item_catalog", plugin_state.item_catalog)
    end
    
    log_operation("ADD", string.format("Added item %d: %s", item_data.id, item_data.name))
    return true
end

function getItemCatalogSummary()
    initialize_item_database()
    
    local summary = {
        total_items = table_count(plugin_state.item_catalog),
        by_type = {},
        by_rarity = {}
    }
    
    -- Count by type
    for _, item in pairs(plugin_state.item_catalog) do
        summary.by_type[item.type] = (summary.by_type[item.type] or 0) + 1
        summary.by_rarity[item.rarity] = (summary.by_rarity[item.rarity] or 0) + 1
    end
    
    return summary
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS
-- ============================================================================

-- Analyze item usage patterns via Analytics Engine
function analyzeItemUsagePatterns()
    load_phase11_dependencies()
    initialize_item_database()
    
    local analytics = dependencies.analytics
    
    -- Collect item usage data
    local usage_data = {}
    for item_id, count in pairs(plugin_state.usage_stats) do
        table.insert(usage_data, {item_id = item_id, count = count})
    end
    
    local analysis = {
        total_items = table_count(plugin_state.item_catalog),
        usage_tracked = #usage_data,
        patterns = {}
    }
    
    if analytics and analytics.PatternRecognition then
        local patterns = analytics.PatternRecognition.analyzePatterns(usage_data)
        analysis.patterns = patterns
    end
    
    log_operation("ANALYTICS", string.format("Analyzed usage patterns for %d items", analysis.usage_tracked))
    return analysis
end

-- Export item catalog via Import/Export Manager
function exportItemCatalog(format, path)
    load_phase11_dependencies()
    initialize_item_database()
    
    local exporter = dependencies.import_export and dependencies.import_export.DataExporter
    
    local export_data = {
        version = "1.0.0",
        timestamp = os.time(),
        item_count = table_count(plugin_state.item_catalog),
        items = plugin_state.item_catalog
    }
    
    if exporter then
        local fmt = (format or "json"):lower()
        local output = path or ("item_catalog_" .. os.date("%Y%m%d_%H%M%S") .. "." .. fmt)
        
        if fmt == "csv" then
            exporter.exportToCSV(export_data, output, true)
        elseif fmt == "xml" then
            exporter.exportToXML(export_data, output)
        else
            exporter.exportToJSON(export_data, output)
        end
        
        log_operation("EXPORT", string.format("Exported item catalog to %s", output))
        return {path = output, format = fmt}
    end
    
    return {success = false, error = "Import/Export unavailable"}
end

-- Sync item catalog to Integration Hub
function syncItemCatalogToHub()
    load_phase11_dependencies()
    initialize_item_database()
    
    local hub = dependencies.integration_hub
    
    if hub and hub.UnifiedAPI then
        local result = hub.UnifiedAPI.broadcastEvent("item_catalog_sync", {
            catalog = plugin_state.item_catalog,
            summary = getItemCatalogSummary(),
            timestamp = os.time()
        })
        
        log_operation("SYNC_HUB", "Synced item catalog to Integration Hub")
        return result or {success = true}
    end
    
    return {success = false, error = "Integration Hub unavailable"}
end

-- Create item catalog snapshot via Backup/Restore System
function createItemCatalogSnapshot(label)
    load_phase11_dependencies()
    initialize_item_database()
    
    local backup = dependencies.backup_restore
    
    local snapshot = {
        label = label or "item_catalog_snapshot",
        timestamp = os.time(),
        catalog = plugin_state.item_catalog,
        summary = getItemCatalogSummary()
    }
    
    if backup and backup.SnapshotManagement then
        local snap_id = backup.SnapshotManagement.createSnapshot(label, snapshot)
        log_operation("SNAPSHOT", string.format("Created item catalog snapshot: %s", label))
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

-- Register REST API endpoints for item database
function registerItemDatabaseAPI()
    load_phase11_dependencies()
    initialize_item_database()
    
    local api = dependencies.api_gateway
    
    if not api or not api.RESTInterface then return {success = false} end
    
    -- Item lookup endpoint
    local lookup_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/items/:id", function(params)
        return getItemById(tonumber(params.id))
    end)
    
    -- Item search endpoint
    local search_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/items/search", function(params)
        return searchItemsByName(params.q or "")
    end)
    
    -- Catalog summary endpoint
    local summary_endpoint = api.RESTInterface.registerEndpoint("GET", "/api/items/summary", function()
        return getItemCatalogSummary()
    end)
    
    if lookup_endpoint then api.RESTInterface.addRateLimit(lookup_endpoint.endpoint_id, 120) end
    if search_endpoint then api.RESTInterface.addRateLimit(search_endpoint.endpoint_id, 60) end
    
    log_operation("API", "Registered item database REST endpoints")
    return {lookup = lookup_endpoint or {registered = true}, search = search_endpoint or {registered = true}, summary = summary_endpoint or {registered = true}}
end

-- Initialize on load
initialize_item_database()
log_operation("LOAD", "Item Database Plugin v1.0.0 loaded")

-- Export public API
return {
    -- Core lookup
    getItemById = getItemById,
    searchItemsByName = searchItemsByName,
    getItemsByType = getItemsByType,
    getItemsByRarity = getItemsByRarity,
    
    -- Catalog management
    addItemToCatalog = addItemToCatalog,
    getItemCatalogSummary = getItemCatalogSummary,
    
    -- Phase 11 integrations
    analyzeItemUsagePatterns = analyzeItemUsagePatterns,
    exportItemCatalog = exportItemCatalog,
    syncItemCatalogToHub = syncItemCatalogToHub,
    createItemCatalogSnapshot = createItemCatalogSnapshot,
    registerItemDatabaseAPI = registerItemDatabaseAPI
}
