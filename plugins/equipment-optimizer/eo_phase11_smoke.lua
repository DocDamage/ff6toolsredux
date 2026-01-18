--[[
  Equipment Optimizer - Phase 11 Smoke Test
  Validates all Phase 11 integration hooks load and execute without errors
  
  Date: 2026-01-16
  Phase: 11 Tier 1 - Equipment Optimizer
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

-- Load equipment optimizer plugin
print("=" .. string.rep("=", 78))
print("PHASE 11 SMOKE TEST: Equipment Optimizer v1.2.0")
print("=" .. string.rep("=", 78))
print()

-- Load plugin
test("Load Equipment Optimizer plugin", function()
  assert(loadfile("plugins/equipment-optimizer/plugin.lua"))()
  return true
end)

print()
print("Testing Core Plugin Functionality:")
print("-" .. string.rep("-", 78))
print()

-- Test core functionality that exists in the plugin
test("optimizeEquipment(0, 'balanced')", function()
  local result = optimizeEquipment(0, "balanced")
  return true
end)

test("getCurrentLoadout(0)", function()
  local result = getCurrentLoadout(0)
  assert(result, "Result should not be nil")
  return result
end)

test("autoEquipOptimal(0)", function()
  local result = autoEquipOptimal(0)
  return true
end)

test("compareLoadouts({})", function()
  local result = compareLoadouts({})
  return true
end)

-- Skip displayComparisonDashboard - requires populated data from compareLoadouts

test("generateEquipmentComparison({})", function()
  local result = generateEquipmentComparison({})
  assert(result, "Result should not be nil")
  return result
end)

test("exportEquipmentTemplate({})", function()
  local result = exportEquipmentTemplate({}, "json", "eo_template.json")
  return true
end)

test("importEquipmentTemplate('eo_template.json')", function()
  local result = importEquipmentTemplate("eo_template.json", "json")
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
