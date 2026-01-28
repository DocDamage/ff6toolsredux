--[[
  Location Database Plugin - Smoke Test Suite
  Tests location lookup, management, and Phase 11 integrations
]]

local test_count = 0
local pass_count = 0
local fail_count = 0

local function test(name, fn)
    test_count = test_count + 1
    local status, err = pcall(fn)
    if status then
        pass_count = pass_count + 1
        print(string.format("[PASS] Test %d: %s", test_count, name))
    else
        fail_count = fail_count + 1
        print(string.format("[FAIL] Test %d: %s - %s", test_count, name, tostring(err)))
    end
end

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

local location_db = require("plugins.location-database.plugin")

print("\n=== LOCATION DATABASE SMOKE TESTS ===\n")

print("--- Core Location Lookup Tests ---")

test("Get location by valid ID (Narshe)", function()
    local loc = location_db.getLocationById(0)
    assert(loc ~= nil, "Location should not be nil")
    assert(loc.name == "Narshe", "Location name should be Narshe")
end)

test("Get location by invalid ID returns nil", function()
    local loc = location_db.getLocationById(999)
    assert(loc == nil, "Location should be nil for invalid ID")
end)

test("Search locations by name (Tower)", function()
    local results = location_db.searchLocationsByName("Tower")
    assert(#results >= 1, "Should find Kefka's Tower")
end)

test("Get locations by region (World of Ruin)", function()
    local results = location_db.getLocationsByRegion("World of Ruin")
    assert(#results >= 3, "Should have WOR locations")
end)

test("Get locations by type (town)", function()
    local results = location_db.getLocationsByType("town")
    assert(#results >= 3, "Should have town entries")
end)

print("\n--- Location Management Tests ---")

test("Add new location to catalog", function()
    local new_location = {
        id = 200,
        name = "Fanatics Tower",
        region = "World of Ruin",
        type = "tower",
        accessible = true
    }
    local success = location_db.addLocationToCatalog(new_location)
    assert(success == true, "Add location should succeed")
    local loc = location_db.getLocationById(200)
    assert(loc ~= nil and loc.name == "Fanatics Tower", "Location should be retrievable")
end)

test("Record visit", function()
    location_db.recordVisit(0)
    location_db.recordVisit(0)
    location_db.recordVisit(160)
    assert(true, "Recording visits should not fail")
end)

test("Get location catalog summary", function()
    local summary = location_db.getLocationCatalogSummary()
    assert(summary.total_locations > 0, "Total locations should be > 0")
    assert(summary.by_region["World of Ruin"] ~= nil, "Should count WOR region")
end)

print("\n--- Phase 11 Integration Tests ---")

test("Analyze visit patterns", function()
    local analysis = location_db.analyzeVisitPatterns()
    assert(analysis.total_locations > 0, "Total locations should be > 0")
end)

test("Export location catalog (JSON)", function()
    local result = location_db.exportLocationCatalog("json", "loc_test.json")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Export location catalog (CSV)", function()
    local result = location_db.exportLocationCatalog("csv", "loc_test.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync location catalog to Integration Hub", function()
    local result = location_db.syncLocationCatalogToHub()
    assert(result ~= nil, "Sync result should not be nil")
end)

test("Create location catalog snapshot", function()
    local snap = location_db.createLocationCatalogSnapshot("test_snapshot")
    assert(snap ~= nil, "Snapshot should not be nil")
end)

test("Register location REST API", function()
    local result = location_db.registerLocationDatabaseAPI()
    assert(result.lookup ~= nil, "Lookup endpoint should be registered")
    assert(result.search ~= nil, "Search endpoint should be registered")
    assert(result.summary ~= nil, "Summary endpoint should be registered")
end)

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
