--[[
  NPC Database Plugin - Smoke Test Suite
  Tests NPC lookup, registry management, and Phase 11 integrations
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

local npc_db = require("plugins.npc-database.plugin")

print("\n=== NPC DATABASE SMOKE TESTS ===\n")

print("--- Core NPC Lookup Tests ---")

test("Get NPC by valid ID (Arvis)", function()
    local npc = npc_db.getNPCById(0)
    assert(npc ~= nil, "NPC should not be nil")
    assert(npc.name == "Arvis", "NPC name should be Arvis")
end)

test("Get NPC by invalid ID returns nil", function()
    local npc = npc_db.getNPCById(999)
    assert(npc == nil, "NPC should be nil for invalid ID")
end)

test("Search NPCs by name (Mog)", function()
    local results = npc_db.searchNPCsByName("Mog")
    assert(#results >= 1, "Should find Mog")
end)

test("Get NPCs by location (Narshe)", function()
    local results = npc_db.getNPCsByLocation("Narshe")
    assert(#results >= 2, "Should find NPCs in Narshe")
end)

test("Get NPCs by role (ally)", function()
    local results = npc_db.getNPCsByRole("ally")
    assert(#results >= 4, "Should have ally NPCs")
end)

print("\n--- NPC Management Tests ---")

test("Add new NPC", function()
    local new_npc = {
        id = 210,
        name = "Relm",
        role = "ally",
        alignment = "friendly",
        location = "Thamasa",
        quests = {"Thamasa"}
    }
    local success = npc_db.addNPC(new_npc)
    assert(success == true, "Add NPC should succeed")
    local npc = npc_db.getNPCById(210)
    assert(npc ~= nil and npc.name == "Relm", "NPC should be retrievable")
end)

test("Record interaction", function()
    npc_db.recordInteraction(0)
    npc_db.recordInteraction(0)
    npc_db.recordInteraction(30)
    assert(true, "Recording interactions should not fail")
end)

test("Get NPC registry summary", function()
    local summary = npc_db.getNPCRegistrySummary()
    assert(summary.total_npcs > 0, "Total NPCs should be > 0")
    assert(summary.by_role["ally"] ~= nil, "Should count allies")
end)

print("\n--- Phase 11 Integration Tests ---")

test("Analyze interaction patterns", function()
    local analysis = npc_db.analyzeInteractionPatterns()
    assert(analysis.total_npcs > 0, "Total NPCs should be > 0")
end)

test("Export NPC registry (JSON)", function()
    local result = npc_db.exportNPCRegistry("json", "npc_test.json")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Export NPC registry (CSV)", function()
    local result = npc_db.exportNPCRegistry("csv", "npc_test.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync NPC registry to Integration Hub", function()
    local result = npc_db.syncNPCRegistryToHub()
    assert(result ~= nil, "Sync result should not be nil")
end)

test("Create NPC registry snapshot", function()
    local snap = npc_db.createNPCRegistrySnapshot("test_snapshot")
    assert(snap ~= nil, "Snapshot should not be nil")
end)

test("Register NPC REST API", function()
    local result = npc_db.registerNPCDatabaseAPI()
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
