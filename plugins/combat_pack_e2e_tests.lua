-- Test script to verify save manipulation works end-to-end
-- Run via CLI: go run . combat-pack --mode smoke --file test_save.json

-- This extends the smoke tests to verify save bindings
local base_tests = require('plugins.combat_depth_pack_smoke_tests')
local results = {}

-- Test save bindings availability
local function test_save_bindings()
  if not save then
    return {success = false, error = "save table not available"}
  end
  
  local tests = {
    "getCharacterCount",
    "getCharacterName", 
    "setCharacterLevel",
    "setCharacterHP",
    "setCharacterMP",
    "getGil",
    "setGil",
    "log"
  }
  
  local missing = {}
  for _, fn_name in ipairs(tests) do
    if not save[fn_name] then
      table.insert(missing, fn_name)
    end
  end
  
  if #missing > 0 then
    return {success = false, missing_functions = missing}
  end
  
  return {success = true, bindings_count = #tests}
end

-- Test reading character data
local function test_read_characters()
  if not save then
    return {success = false, error = "save not available"}
  end
  
  local count = save.getCharacterCount()
  local chars = {}
  
  for i = 0, math.min(count - 1, 3) do
    local name = save.getCharacterName(i)
    if name and name ~= "Unknown" then
      table.insert(chars, name)
    end
  end
  
  return {
    success = true,
    total_characters = count,
    named_characters = chars,
    found_count = #chars
  }
end

-- Test gil manipulation
local function test_gil_operations()
  if not save then
    return {success = false, error = "save not available"}
  end
  
  local original = save.getGil()
  local test_amount = 99999
  
  save.setGil(test_amount)
  local after_set = save.getGil()
  
  -- Restore original
  save.setGil(original)
  local restored = save.getGil()
  
  return {
    success = (after_set == test_amount and restored == original),
    original_gil = original,
    test_amount = test_amount,
    after_set = after_set,
    restored = restored
  }
end

-- Test character modification
local function test_character_modification()
  if not save then
    return {success = false, error = "save not available"}
  end
  
  local count = save.getCharacterCount()
  if count == 0 then
    return {success = false, error = "no characters in save"}
  end
  
  -- Try to set level for first character
  local success, err = save.setCharacterLevel(0, 99)
  if not success then
    return {success = false, error = err or "setCharacterLevel failed"}
  end
  
  -- Try HP/MP
  success = save.setCharacterHP(0, 9999)
  if not success then
    return {success = false, error = "setCharacterHP failed"}
  end
  
  success = save.setCharacterMP(0, 999)
  if not success then
    return {success = false, error = "setCharacterMP failed"}
  end
  
  return {
    success = true,
    modified_character = 0,
    operations = {"setLevel(99)", "setHP(9999)", "setMP(999)"}
  }
end

-- Run enhanced test suite
local function run_all_with_save()
  local summary = {tests_run = 0, tests_passed = 0, tests_failed = 0}
  
  -- Run base smoke tests
  local base_result = base_tests.run_all()
  summary.tests_run = base_result.total or 0
  summary.tests_passed = base_result.passed or 0
  summary.tests_failed = base_result.failed or 0
  
  -- Run save binding tests
  local binding_test = test_save_bindings()
  summary.tests_run = summary.tests_run + 1
  if binding_test.success then
    summary.tests_passed = summary.tests_passed + 1
    table.insert(results, {name = "Save Bindings Available", status = "PASS"})
  else
    summary.tests_failed = summary.tests_failed + 1
    table.insert(results, {name = "Save Bindings Available", status = "FAIL", error = binding_test.error})
  end
  
  -- Only run data tests if save is available
  if save then
    local read_test = test_read_characters()
    summary.tests_run = summary.tests_run + 1
    if read_test.success then
      summary.tests_passed = summary.tests_passed + 1
      table.insert(results, {
        name = "Read Characters", 
        status = "PASS",
        details = read_test
      })
    else
      summary.tests_failed = summary.tests_failed + 1
      table.insert(results, {name = "Read Characters", status = "FAIL", error = read_test.error})
    end
    
    local gil_test = test_gil_operations()
    summary.tests_run = summary.tests_run + 1
    if gil_test.success then
      summary.tests_passed = summary.tests_passed + 1
      table.insert(results, {name = "Gil Operations", status = "PASS", details = gil_test})
    else
      summary.tests_failed = summary.tests_failed + 1
      table.insert(results, {name = "Gil Operations", status = "FAIL", error = gil_test.error})
    end
    
    local mod_test = test_character_modification()
    summary.tests_run = summary.tests_run + 1
    if mod_test.success then
      summary.tests_passed = summary.tests_passed + 1
      table.insert(results, {name = "Character Modification", status = "PASS", details = mod_test})
    else
      summary.tests_failed = summary.tests_failed + 1
      table.insert(results, {name = "Character Modification", status = "FAIL", error = mod_test.error})
    end
  end
  
  summary.save_hooks_enabled = (save ~= nil)
  summary.test_results = results
  
  return summary
end

-- Export
return {
  run_all = run_all_with_save,
  test_save_bindings = test_save_bindings,
  test_read_characters = test_read_characters,
  test_gil_operations = test_gil_operations,
  test_character_modification = test_character_modification
}
