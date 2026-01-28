--[[
  Storyline Database Plugin - Smoke Test Suite
  Tests story lookup, flag management, and Phase 11 integrations
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

local storyline_db = require("plugins.storyline-database.plugin")

print("\n=== STORYLINE DATABASE SMOKE TESTS ===\n")

print("--- Story Lookup Tests ---")

test("Get story event by valid ID (Narshe Defense)", function()
    local ev = storyline_db.getStoryEventById(0)
    assert(ev ~= nil, "Event should not be nil")
    assert(ev.name == "Narshe Defense", "Name should be Narshe Defense")
end)

test("Get story event by invalid ID returns nil", function()
    local ev = storyline_db.getStoryEventById(999)
    assert(ev == nil, "Event should be nil for invalid ID")
end)

test("Search story events by name (Opera)", function()
    local results = storyline_db.searchStoryEventsByName("Opera")
    assert(#results >= 1, "Should find Opera House")
end)

test("Get events by chapter (WOR)", function()
    local results = storyline_db.getEventsByChapter("WOR")
    assert(#results >= 2, "Should have events in WOR chapter")
end)

test("Get events by status (pending)", function()
    local results = storyline_db.getEventsByStatus("pending")
    assert(#results >= 2, "Should have pending events")
end)

test("Get events by type (main)", function()
    local results = storyline_db.getEventsByType("main")
    assert(#results >= 4, "Should have main events")
end)

print("\n--- Story Management Tests ---")

test("Add new story event", function()
    local new_event = {
        id = 210,
        name = "Fanatics Tower",
        chapter = "WOR",
        type = "optional",
        status = "pending",
        location = "Fanatics Tower"
    }
    local success = storyline_db.addStoryEvent(new_event)
    assert(success == true, "Add event should succeed")
    local ev = storyline_db.getStoryEventById(210)
    assert(ev ~= nil and ev.name == "Fanatics Tower", "Event should be retrievable")
end)

test("Update story flag", function()
    local success = storyline_db.updateStoryFlag("party_split", true)
    assert(success == true, "Flag update should succeed")
    local summary = storyline_db.getStoryProgressSummary()
    assert(summary.flags.party_split == true, "Flag should be true")
end)

test("Get story progress summary", function()
    local summary = storyline_db.getStoryProgressSummary()
    assert(summary.total_events > 0, "Total events should be > 0")
    assert(summary.by_status ~= nil, "Should have by_status counts")
end)

print("\n--- Phase 11 Integration Tests ---")

test("Analyze storyline coverage", function()
    local analysis = storyline_db.analyzeStorylineCoverage()
    assert(analysis.total_events > 0, "Total events should be > 0")
end)

test("Export storyline data (JSON)", function()
    local result = storyline_db.exportStorylineData("json", "story_test.json")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Export storyline data (CSV)", function()
    local result = storyline_db.exportStorylineData("csv", "story_test.csv")
    assert(result ~= nil, "Export result should not be nil")
end)

test("Sync storyline to Integration Hub", function()
    local result = storyline_db.syncStorylineToHub()
    assert(result ~= nil, "Sync result should not be nil")
end)

test("Create storyline snapshot", function()
    local snap = storyline_db.createStorylineSnapshot("test_snapshot")
    assert(snap ~= nil, "Snapshot should not be nil")
end)

test("Register storyline REST API", function()
    local result = storyline_db.registerStorylineDatabaseAPI()
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
