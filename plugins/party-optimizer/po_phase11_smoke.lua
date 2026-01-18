--[[
  Party Optimizer - Phase 11 Smoke Test
  Validates all Phase 11 integration hooks load and execute without errors
  
  Date: 2026-01-16
  Phase: 11 Tier 1 - Party Optimizer
]]

-- Setup package path for module loading
package.path = package.path .. ";./?.lua;./?/init.lua;./plugins/?.lua;./plugins/?/?.lua;./plugins/?/init.lua"

-- Stub editor APIs (override file I/O and UI)
ReadFile = function(path) return nil end
WriteFile = function(path, content) return true end
ShowDialog = function(title, lines) print("[UI] " .. title); return 6 end

-- Initialize test environment
local test_results = {}
local test_count = 0
local pass_count = 0

-- Test helper
local function test(label, fn)
  test_count = test_count + 1
  local success, result = pcall(fn)
  
  if success then
    pass_count = pass_count + 1
    print(string.format("[PASS] %s", label))
    table.insert(test_results, {label = label, status = "PASS", result = result})
    return result
  else
    print(string.format("[FAIL] %s: %s", label, tostring(result)))
    table.insert(test_results, {label = label, status = "FAIL", error = tostring(result)})
    return nil
  end
end

-- Load party optimizer plugin
print("=" .. string.rep("=", 78))
print("PHASE 11 SMOKE TEST: Party Optimizer v1.2.0")
print("=" .. string.rep("=", 78))
print()

-- Load plugin
test("Load Party Optimizer plugin", function()
  assert(loadfile("plugins/party-optimizer/plugin.lua"))()
  return true
end)

print()
print("Testing Core Plugin Functionality:")
print("-" .. string.rep("-", 78))
print()

-- Test core functionality that exists in the plugin
test("optimizePartyComposition({}, 'balanced')", function()
  local result = optimizePartyComposition({}, "balanced")
  return true
end)

test("recommendPartyForScenario({})", function()
  local result = recommendPartyForScenario({})
  return true
end)

test("autoConfigureParty({}, {})", function()
  local result = autoConfigureParty({}, {})
  return true
end)

test("visualizePartySynergy({})", function()
  local result = visualizePartySynergy({})
  return true
end)

test("exportPartyTemplate({}, 'json')", function()
  local result = exportPartyTemplate({}, "json", "po_party.json")
  assert(result, "Result should not be nil")
  return result
end)

test("importPartyTemplate('po_party.json', 'json')", function()
  local result = importPartyTemplate("po_party.json", "json")
  return true
end)

test("snapshotParty({}, 'smoke')", function()
  local result = snapshotParty({}, "smoke")
  return true
end)

test("syncPartyData()", function()
  local result = syncPartyData({"character-roster-editor"})
  return true
end)

-- Summary
print()
print("=" .. string.rep("=", 78))
print(string.format("SMOKE TEST RESULTS: %d/%d TESTS PASSED", pass_count, test_count))
print("=" .. string.rep("=", 78))
print()

if pass_count == test_count then
  print("✓ All Phase 11 hooks loaded and executed successfully!")
  print("✓ Integration is functional with graceful fallbacks.")
  os.exit(0)
else
  print("✗ Some tests failed:")
  for _, result in ipairs(test_results) do
    if result.status == "FAIL" then
      print(string.format("  ✗ %s: %s", result.label, result.error))
    end
  end
  os.exit(1)
end
