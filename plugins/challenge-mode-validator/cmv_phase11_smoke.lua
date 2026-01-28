--[[
  Challenge Mode Validator - Phase 11 Smoke Test
  Validates all Phase 11 integration hooks load and execute without errors
  
  Date: 2026-01-16
  Phase: 11 Tier 1 - Challenge Mode Validator
]]

-- Setup package path for module loading
package.path = package.path .. ";./?.lua;./?/init.lua;./plugins/?.lua;./plugins/?/?.lua;./plugins/?/init.lua"

-- Stub editor APIs (override file I/O and UI)
ReadFile = function(path)
  return nil
end

WriteFile = function(path, content)
  return true
end

ShowDialog = function(title, lines)
  print("[UI] " .. title)
  return 6  -- Exit
end

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

-- Load challenge validator plugin
print("=" .. string.rep("=", 78))
print("PHASE 11 SMOKE TEST: Challenge Mode Validator v1.2.0")
print("=" .. string.rep("=", 78))
print()

-- Load plugin
test("Load Challenge Validator plugin", function()
  assert(loadfile("plugins/challenge-mode-validator/plugin.lua"))()
  return true
end)

print()
print("Testing Phase 11 Integration Hooks:")
print("-" .. string.rep("-", 78))
print()

-- Test each Phase 11 hook
test("analyzeChallengeTrends()", function()
  local result = analyzeChallengeTrends()
  assert(result, "Result should not be nil")
  assert(result.counts, "Result should have counts")
  assert(result.trend, "Result should have trend")
  return result
end)

test("generateComplianceDashboard()", function()
  local result = generateComplianceDashboard()
  assert(result, "Result should not be nil")
  assert(result.dashboard, "Result should have dashboard")
  assert(result.chart, "Result should have chart")
  assert(result.analysis, "Result should have analysis")
  return result
end)

test("exportComplianceReport('json')", function()
  local result = exportComplianceReport("json", "cmv_compliance_test.json")
  assert(result, "Result should not be nil")
  assert(result.name or result.format, "Result should have name or format")
  return result
end)

test("exportChallengeTemplate({id='test'})", function()
  local result = exportChallengeTemplate({id = "test"}, "json", "cmv_template_test.json")
  assert(result, "Result should not be nil")
  assert(result.path, "Result should have path")
  assert(result.format, "Result should have format")
  return result
end)

test("importChallengeTemplate('cmv_template_test.json')", function()
  local result = importChallengeTemplate("cmv_template_test.json", "json")
  assert(result, "Result should not be nil")
  return result
end)

test("snapshotChallengeState('smoke')", function()
  local result = snapshotChallengeState("smoke")
  assert(result, "Result should not be nil")
  assert(result.snapshot_id or result.name, "Result should have snapshot_id or name")
  return result
end)

test("restoreChallengeState()", function()
  local result = restoreChallengeState("local")
  assert(result, "Result should not be nil")
  return result
end)

test("scheduleChallengeValidation('*/30 * * * *')", function()
  local result = scheduleChallengeValidation("*/30 * * * *")
  assert(result, "Result should not be nil")
  assert(result.scheduled or result.cron, "Result should have scheduled or cron")
  return result
end)

test("syncChallengeData({'party-optimizer', 'equipment-optimizer'})", function()
  local result = syncChallengeData({"party-optimizer", "equipment-optimizer"})
  assert(result, "Result should not be nil")
  return result
end)

test("registerChallengeAPI()", function()
  local result = registerChallengeAPI()
  -- This may return nil if API Gateway unavailable (expected fallback)
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
  print()
  for _, result in ipairs(test_results) do
    if result.status == "PASS" then
      print(string.format("  ✓ %s", result.label))
    end
  end
  os.exit(0)
else
  print("✗ Some tests failed:")
  print()
  for _, result in ipairs(test_results) do
    if result.status == "FAIL" then
      print(string.format("  ✗ %s", result.label))
      print(string.format("    Error: %s", result.error))
    end
  end
  os.exit(1)
end
