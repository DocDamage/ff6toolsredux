--[[
  Item Database Plugin - Smoke Test Suite
  Tests core item lookup, catalog management, and Phase 11 integrations
]]

-- Test helpers
local test_count = 0
local pass_count = 0
local fail_count = 0

local function test(name, fn)
    test_count = test_count + 1
    local status, err = pcall(fn)
    if status then
        pass_count = pass_count + 1
        print(string.format("[PASS] Test %d: %s", test_count, name))
        return true
    else
        fail_count = fail_count + 1
        print(string.format("[FAIL] Test %d: %s - %s", test_count, name, tostring(err)))
        return false
    end
end

-- Mock dependencies for Phase 11 tests
local mock_analytics = {
    PatternRecognition = {
        analyzePatterns = function(data)
            return {pattern_count = #data, analysis = "mock_patterns"}
        end
    }
}

local mock_import_export = {
    DataExporter = {
        exportToJSON = function(data, path)
            print("[MOCK] Exported to JSON: " .. path)
            return true
        end,
        exportToCSV = function(data, path, headers)
            print("[MOCK] Exported to CSV: " .. path)
            return true
        end,
        exportToXML = function(data, path)
            print("[MOCK] Exported to XML: " .. path)
            return true
        end
    }
}

local mock_backup_restore = {
    SnapshotManagement = {
        createSnapshot = function(label, data)
            return {snapshot_id = "snap_" .. label, created_at = os.time()}
        end
    }
}

local mock_integration_hub = {
    UnifiedAPI = {
        broadcastEvent = function(event_type, data)
            return {broadcasted = true, event = event_type, timestamp = os.time()}
        end
    }
}

local mock_api_gateway = {
    RESTInterface = {
        registerEndpoint = function(method, path, handler)
            return {endpoint_id = "ep_" .. method .. "_" .. path, registered = true}
        end,
        addRateLimit = function(endpoint_id, limit)
            return true
        end
    }
}

_G.mock_analytics = mock_analytics
_G.mock_import_export = mock_import_export
_G.mock_backup_restore = mock_backup_restore
_G.mock_integration_hub = mock_integration_hub
_G.mock_api_gateway = mock_api_gateway

-- Mock database persistence layer
local mock_database_layer = {
    savePersistentData = function(key, data)
        print("[MOCK DB] Saved data: " .. key)
        return true
    end,
    loadPersistentData = function(key)
        print("[MOCK DB] Loaded data: " .. key)
        return nil
    end
}

_G.mock_database_layer = mock_database_layer

-- Safe require with fallback to mock
local original_require = require
_G.require = function(module_path)
    if module_path == "plugins.database-persistence-layer.plugin" then
        return mock_database_layer
    elseif module_path == "plugins.advanced-analytics-engine.v1_0_core" then
        return mock_analytics
    elseif module_path == "plugins.import-export-manager.v1_0_core" then
        return mock_import_export
    elseif module_path == "plugins.backup-restore-system.v1_0_core" then
        return mock_backup_restore
    elseif module_path == "plugins.integration-hub.v1_0_core" then
        return mock_integration_hub
    elseif module_path == "plugins.api-gateway.v1_0_core" then
        return mock_api_gateway
    end
    
    return original_require(module_path)
end

-- Load plugin
local item_db = require("plugins.item-database.plugin")

print("\n=== ITEM DATABASE SMOKE TESTS ===\n")

-- ============================================================================
-- CORE ITEM LOOKUP TESTS
-- ============================================================================

print("--- Core Item Lookup Tests ---")

test("Get item by valid ID (Dirk)", function()
    local item = item_db.getItemById(0)
    assert(item ~= nil, "Item should not be nil")
    assert(item.name == "Dirk", "Item name should be 'Dirk'")
    assert(item.type == "weapon", "Item type should be 'weapon'")
end)

test("Get item by valid ID (Ultima Weapon)", function()
    local item = item_db.getItemById(255)
    assert(item ~= nil, "Item should not be nil")
    assert(item.name == "Ultima Weapon", "Item name should be 'Ultima Weapon'")
    assert(item.rarity == 5, "Item rarity should be 5 (legendary)")
end)

test("Get item by invalid ID returns nil", function()
    local item = item_db.getItemById(999)
    assert(item == nil, "Item should be nil for invalid ID")
end)

test("Search items by name (partial match)", function()
    local results = item_db.searchItemsByName("weapon")
    assert(results ~= nil, "Results should not be nil")
    assert(#results > 0, "Should find at least one item matching 'weapon'")
end)

test("Get items by type (weapons)", function()
    local results = item_db.getItemsByType("weapon")
    assert(results ~= nil, "Results should not be nil")
    assert(#results >= 3, "Should have at least 3 weapons in catalog")
end)

test("Get items by rarity (legendary)", function()
    local results = item_db.getItemsByRarity(5)
    assert(results ~= nil, "Results should not be nil")
    assert(#results >= 2, "Should have at least 2 legendary items")
end)

-- ============================================================================
-- CATALOG MANAGEMENT TESTS
-- ============================================================================

print("\n--- Catalog Management Tests ---")

test("Add new item to catalog", function()
    local new_item = {
        id = 50,
        name = "Test Sword",
        type = "weapon",
        rarity = 2,
        attack = 25,
        price = 500
    }
    local success = item_db.addItemToCatalog(new_item)
    assert(success == true, "Add item should succeed")
    
    local retrieved = item_db.getItemById(50)
    assert(retrieved ~= nil, "Newly added item should be retrievable")
    assert(retrieved.name == "Test Sword", "Retrieved item should match added item")
end)

test("Get catalog summary", function()
    local summary = item_db.getItemCatalogSummary()
    assert(summary ~= nil, "Summary should not be nil")
    assert(summary.total_items > 0, "Total items should be greater than 0")
    assert(summary.by_type ~= nil, "Should have items categorized by type")
    assert(summary.by_rarity ~= nil, "Should have items categorized by rarity")
end)

-- ============================================================================
-- PHASE 11 INTEGRATION TESTS
-- ============================================================================

print("\n--- Phase 11 Integration Tests ---")

test("Analyze item usage patterns", function()
    local analysis = item_db.analyzeItemUsagePatterns()
    assert(analysis ~= nil, "Analysis should not be nil")
    assert(analysis.total_items > 0, "Total items should be greater than 0")
end)

test("Export item catalog (JSON)", function()
    local result = item_db.exportItemCatalog("json", "test_catalog.json")
    assert(result ~= nil, "Export result should not be nil")
    assert(result.path ~= nil or result.success ~= nil, "Export should return path or success")
end)

test("Export item catalog (CSV)", function()
    local result = item_db.exportItemCatalog("csv", "test_catalog.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync item catalog to Integration Hub", function()
    local result = item_db.syncItemCatalogToHub()
    assert(result ~= nil, "Sync result should not be nil")
    assert(result.success ~= false or result.broadcasted == true, "Sync should succeed or indicate broadcast")
end)

test("Create item catalog snapshot", function()
    local snap_id = item_db.createItemCatalogSnapshot("test_snapshot")
    assert(snap_id ~= nil, "Snapshot ID should not be nil")
end)

test("Register item database REST API endpoints", function()
    local result = item_db.registerItemDatabaseAPI()
    assert(result ~= nil, "API registration result should not be nil")
    assert(result.lookup ~= nil, "Lookup endpoint should be registered")
    assert(result.search ~= nil, "Search endpoint should be registered")
    assert(result.summary ~= nil, "Summary endpoint should be registered")
end)

-- ============================================================================
-- TEST SUMMARY
-- ============================================================================

print("\n=== TEST SUMMARY ===")
print(string.format("Total Tests: %d", test_count))
print(string.format("Passed: %d", pass_count))
print(string.format("Failed: %d", fail_count))
print(string.format("Success Rate: %.1f%%", (pass_count / test_count) * 100))

if fail_count == 0 then
    print("\n✓ ALL TESTS PASSED")
    os.exit(0)
else
    print("\n✗ SOME TESTS FAILED")
    os.exit(1)
end
