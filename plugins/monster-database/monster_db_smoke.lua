--[[
  Monster Database Plugin - Smoke Test Suite
  Tests core monster lookup, bestiary management, and Phase 11 integrations
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
local monster_db = require("plugins.monster-database.plugin")

print("\n=== MONSTER DATABASE SMOKE TESTS ===\n")

-- ============================================================================
-- CORE MONSTER LOOKUP TESTS
-- ============================================================================

print("--- Core Monster Lookup Tests ---")

test("Get monster by valid ID (Guard)", function()
    local monster = monster_db.getMonsterById(0)
    assert(monster ~= nil, "Monster should not be nil")
    assert(monster.name == "Guard", "Monster name should be 'Guard'")
    assert(monster.type == "human", "Monster type should be 'human'")
end)

test("Get monster by valid ID (Kefka Final)", function()
    local monster = monster_db.getMonsterById(383)
    assert(monster ~= nil, "Monster should not be nil")
    assert(monster.name == "Kefka (Final)", "Monster name should be 'Kefka (Final)'")
    assert(monster.final_boss == true, "Monster should be marked as final boss")
end)

test("Get monster by invalid ID returns nil", function()
    local monster = monster_db.getMonsterById(999)
    assert(monster == nil, "Monster should be nil for invalid ID")
end)

test("Search monsters by name (partial match)", function()
    local results = monster_db.searchMonstersByName("dragon")
    assert(results ~= nil, "Results should not be nil")
    assert(#results > 0, "Should find at least one monster matching 'dragon'")
end)

test("Get monsters by type (dragons)", function()
    local results = monster_db.getMonstersByType("dragon")
    assert(results ~= nil, "Results should not be nil")
    assert(#results >= 1, "Should have at least 1 dragon in bestiary")
end)

test("Get monsters by location (Narshe)", function()
    local results = monster_db.getMonstersByLocation("Narshe")
    assert(results ~= nil, "Results should not be nil")
    assert(#results >= 2, "Should have at least 2 monsters in Narshe")
end)

test("Get monster weaknesses", function()
    local weaknesses = monster_db.getMonsterWeaknesses(200) -- Ice Dragon
    assert(weaknesses ~= nil, "Weaknesses should not be nil")
    assert(weaknesses.weaknesses ~= nil, "Should have weaknesses table")
    assert(weaknesses.absorbs ~= nil, "Should have absorbs table")
end)

-- ============================================================================
-- BESTIARY MANAGEMENT TESTS
-- ============================================================================

print("\n--- Bestiary Management Tests ---")

test("Add new monster to bestiary", function()
    local new_monster = {
        id = 50,
        name = "Test Monster",
        type = "beast",
        hp = 500,
        mp = 100,
        attack = 20,
        defense = 15,
        weaknesses = {"fire"},
        location = "Test Area"
    }
    local success = monster_db.addMonsterToBestiary(new_monster)
    assert(success == true, "Add monster should succeed")
    
    local retrieved = monster_db.getMonsterById(50)
    assert(retrieved ~= nil, "Newly added monster should be retrievable")
    assert(retrieved.name == "Test Monster", "Retrieved monster should match added monster")
end)

test("Get bestiary summary", function()
    local summary = monster_db.getBestiarySummary()
    assert(summary ~= nil, "Summary should not be nil")
    assert(summary.total_monsters > 0, "Total monsters should be greater than 0")
    assert(summary.by_type ~= nil, "Should have monsters categorized by type")
    assert(summary.bosses > 0, "Should have boss count")
end)

test("Record encounter", function()
    monster_db.recordEncounter(0) -- Guard
    monster_db.recordEncounter(0) -- Guard again
    monster_db.recordEncounter(1) -- Leaf Bunny
    -- No assertion needed, just verify it doesn't crash
    assert(true, "Record encounter should not crash")
end)

-- ============================================================================
-- PHASE 11 INTEGRATION TESTS
-- ============================================================================

print("\n--- Phase 11 Integration Tests ---")

test("Analyze encounter patterns", function()
    local analysis = monster_db.analyzeEncounterPatterns()
    assert(analysis ~= nil, "Analysis should not be nil")
    assert(analysis.total_monsters > 0, "Total monsters should be greater than 0")
end)

test("Export bestiary (JSON)", function()
    local result = monster_db.exportBestiary("json", "test_bestiary.json")
    assert(result ~= nil, "Export result should not be nil")
    assert(result.path ~= nil or result.success ~= nil, "Export should return path or success")
end)

test("Export bestiary (CSV)", function()
    local result = monster_db.exportBestiary("csv", "test_bestiary.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync bestiary to Integration Hub", function()
    local result = monster_db.syncBestiaryToHub()
    assert(result ~= nil, "Sync result should not be nil")
    assert(result.success ~= false or result.broadcasted == true, "Sync should succeed or indicate broadcast")
end)

test("Create bestiary snapshot", function()
    local snap_id = monster_db.createBestiarySnapshot("test_snapshot")
    assert(snap_id ~= nil, "Snapshot ID should not be nil")
end)

test("Register monster database REST API endpoints", function()
    local result = monster_db.registerMonsterDatabaseAPI()
    assert(result ~= nil, "API registration result should not be nil")
    assert(result.lookup ~= nil, "Lookup endpoint should be registered")
    assert(result.search ~= nil, "Search endpoint should be registered")
    assert(result.summary ~= nil, "Summary endpoint should be registered")
    assert(result.weakness ~= nil, "Weakness endpoint should be registered")
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
