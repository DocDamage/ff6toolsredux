--[[
  Ability Database Plugin - Smoke Test Suite
  Tests core ability lookup, catalog management, and Phase 11 integrations
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

local ability_db = require("plugins.ability-database.plugin")

print("\n=== ABILITY DATABASE SMOKE TESTS ===\n")

print("--- Core Ability Lookup Tests ---")

test("Get ability by valid ID (Fire)", function()
    local ability = ability_db.getAbilityById(0)
    assert(ability ~= nil, "Ability should not be nil")
    assert(ability.name == "Fire", "Ability name should be 'Fire'")
    assert(ability.type == "magic", "Ability type should be 'magic'")
end)

test("Get ability by valid ID (Ultima)", function()
    local ability = ability_db.getAbilityById(200)
    assert(ability ~= nil, "Ability should not be nil")
    assert(ability.name == "Ultima", "Ability name should be 'Ultima'")
    assert(ability.mp_cost == 80, "MP cost should be 80")
end)

test("Get ability by invalid ID returns nil", function()
    local ability = ability_db.getAbilityById(999)
    assert(ability == nil, "Ability should be nil for invalid ID")
end)

test("Search abilities by name (partial match)", function()
    local results = ability_db.searchAbilitiesByName("ra")
    assert(results ~= nil, "Results should not be nil")
    assert(#results >= 1, "Should find at least one ability matching 'ra'")
end)

test("Get abilities by type (blitz)", function()
    local results = ability_db.getAbilitiesByType("blitz")
    assert(#results >= 2, "Should have at least two blitz abilities")
end)

test("Get abilities by element (lightning)", function()
    local results = ability_db.getAbilitiesByElement("lightning")
    assert(#results >= 2, "Should have at least two lightning abilities")
end)

test("Get abilities by cost range", function()
    local results = ability_db.getAbilitiesByCostRange(0, 10)
    assert(#results >= 3, "Should have at least three abilities in cost range")
end)

print("\n--- Catalog Management Tests ---")

test("Add new ability to catalog", function()
    local new_ability = {
        id = 30,
        name = "Fira",
        type = "magic",
        element = "fire",
        mp_cost = 14,
        power = 55
    }
    local success = ability_db.addAbilityToCatalog(new_ability)
    assert(success == true, "Add ability should succeed")
    local retrieved = ability_db.getAbilityById(30)
    assert(retrieved ~= nil and retrieved.name == "Fira", "New ability should be retrievable")
end)

test("Get ability catalog summary", function()
    local summary = ability_db.getAbilityCatalogSummary()
    assert(summary.total_abilities > 0, "Total abilities should be > 0")
    assert(summary.by_type["magic"] ~= nil, "Should count magic abilities")
end)

print("\n--- Phase 11 Integration Tests ---")

test("Analyze ability usage patterns", function()
    local analysis = ability_db.analyzeAbilityUsagePatterns()
    assert(analysis.total_abilities > 0, "Total abilities should be > 0")
end)

test("Export ability catalog (JSON)", function()
    local result = ability_db.exportAbilityCatalog("json", "test_ability.json")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Export ability catalog (CSV)", function()
    local result = ability_db.exportAbilityCatalog("csv", "test_ability.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync ability catalog to Integration Hub", function()
    local result = ability_db.syncAbilityCatalogToHub()
    assert(result ~= nil, "Sync result should not be nil")
end)

test("Create ability catalog snapshot", function()
    local snap = ability_db.createAbilityCatalogSnapshot("test_snapshot")
    assert(snap ~= nil, "Snapshot should not be nil")
end)

test("Register ability database REST API endpoints", function()
    local result = ability_db.registerAbilityDatabaseAPI()
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
