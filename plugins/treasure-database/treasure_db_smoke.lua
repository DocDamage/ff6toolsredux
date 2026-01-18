--[[
  Treasure Database Plugin - Smoke Test Suite
  Tests treasure lookup, map management, and Phase 11 integrations
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

local treasure_db = require("plugins.treasure-database.plugin")

print("\n=== TREASURE DATABASE SMOKE TESTS ===\n")

print("--- Core Treasure Lookup Tests ---")

test("Get treasure by valid ID (Narshe Mines)", function()
    local treasure = treasure_db.getTreasureById(0)
    assert(treasure ~= nil, "Treasure should not be nil")
    assert(treasure.contents == "Potion", "Contents should be Potion")
end)

test("Get treasure by invalid ID returns nil", function()
    local treasure = treasure_db.getTreasureById(9999)
    assert(treasure == nil, "Treasure should be nil for invalid ID")
end)

test("Get treasure by location (Phoenix Cave)", function()
    local results = treasure_db.getTreasureByLocation("Phoenix Cave")
    assert(#results >= 1, "Should find Phoenix Cave treasure")
end)

test("Get treasure by type (hidden)", function()
    local results = treasure_db.getTreasureByType("hidden")
    assert(#results >= 2, "Should have hidden treasures")
end)

test("Get treasure by rarity (legendary)", function()
    local results = treasure_db.getTreasureByRarity(5)
    assert(#results >= 2, "Should have legendary treasures")
end)

print("\n--- Treasure Management Tests ---")

test("Add new treasure", function()
    local new_treasure = {
        id = 300,
        location = "Fanatics Tower",
        type = "event",
        rarity = 4,
        contents = "Genji Armor",
        opened = false
    }
    local success = treasure_db.addTreasure(new_treasure)
    assert(success == true, "Add treasure should succeed")
    local t = treasure_db.getTreasureById(300)
    assert(t ~= nil and t.contents == "Genji Armor", "Treasure should be retrievable")
end)

test("Mark treasure opened", function()
    local success = treasure_db.markTreasureOpened(0)
    assert(success == true, "Mark opened should succeed")
    local t = treasure_db.getTreasureById(0)
    assert(t.opened == true, "Treasure should be marked opened")
end)

test("Get treasure summary", function()
    local summary = treasure_db.getTreasureSummary()
    assert(summary.total_treasure > 0, "Total treasure should be > 0")
    assert(summary.by_rarity[5] ~= nil, "Should count legendary treasures")
end)

test("Record collection", function()
    treasure_db.recordCollection(120)
    treasure_db.recordCollection(120)
    treasure_db.recordCollection(240)
    assert(true, "Recording collection should not fail")
end)

print("\n--- Phase 11 Integration Tests ---")

test("Analyze collection progress", function()
    local analysis = treasure_db.analyzeCollectionProgress()
    assert(analysis.total_treasure > 0, "Total treasure should be > 0")
end)

test("Export treasure map (JSON)", function()
    local result = treasure_db.exportTreasureMap("json", "treasure_test.json")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Export treasure map (CSV)", function()
    local result = treasure_db.exportTreasureMap("csv", "treasure_test.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync treasure map to Integration Hub", function()
    local result = treasure_db.syncTreasureMapToHub()
    assert(result ~= nil, "Sync result should not be nil")
end)

test("Create treasure snapshot", function()
    local snap = treasure_db.createTreasureSnapshot("test_snapshot")
    assert(snap ~= nil, "Snapshot should not be nil")
end)

test("Register treasure REST API", function()
    local result = treasure_db.registerTreasureDatabaseAPI()
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
