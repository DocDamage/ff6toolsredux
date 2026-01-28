--[[
  Character Roster Editor - Phase 11 Smoke Test
  Validates all Phase 11 integration hooks load and execute without errors
  
  Date: 2026-01-16
  Phase: 11 Tier 1 - Character Roster Editor
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

-- Load character roster editor plugin
print("=" .. string.rep("=", 78))
print("PHASE 11 SMOKE TEST: Character Roster Editor v1.2.0")
print("=" .. string.rep("=", 78))
print()

-- Load plugin
test("Load Character Roster Editor plugin", function()
  assert(loadfile("plugins/character-roster-editor/plugin.lua"))()
  return true
end)

print()
print("Testing Core Plugin Functionality:")
print("-" .. string.rep("-", 78))
print()

-- Test core functionality that exists in the plugin
test("unlockAllCharacters()", function()
  local result = unlockAllCharacters()
  return true  -- Function exists and executes
end)

test("restrictPartySize(4)", function()
  local result = restrictPartySize(4)
  return true
end)

test("getRosterStatus()", function()
  local result = getRosterStatus()
  assert(result, "Status should not be nil")
  return result
end)

test("exportCharacterBuild(0)", function()
  local result = exportCharacterBuild(0, "test.txt")
  assert(result, "Result should not be nil")
  return result
end)

test("importCharacterBuild('test.txt', 0)", function()
  local result = importCharacterBuild("test.txt", 0)
  return true
end)

test("displayRosterStatus()", function()
  local result = displayRosterStatus()
  assert(result, "Result should not be nil")
  return result
end)

test("listRosterTemplates()", function()
  local result = listRosterTemplates()
  assert(result, "Result should not be nil")
  return result
end)

test("configureSoloRun(0)", function()
  local result = configureSoloRun(0)
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
