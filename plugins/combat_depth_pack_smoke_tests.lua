--[[
  Combat Depth Pack Smoke Tests
  Features covered: EncounterTuner, BossRemix, CompanionDirector
  Total tests: 9 (3 per module)
]]

local tests_passed, tests_failed = 0, 0
local results = {}

local function assert_test(condition, name, message)
  if condition then
    tests_passed = tests_passed + 1
    table.insert(results, {name = name, status = "PASS", message = ""})
  else
    tests_failed = tests_failed + 1
    table.insert(results, {name = name, status = "FAIL", message = message or ""})
  end
end

local function summarize()
  print("==== COMBAT DEPTH PACK SMOKE TESTS ====")
  for _, r in ipairs(results) do
    print(string.format("[%s] %s", r.status, r.name))
    if r.message ~= "" then
      print("   -> " .. r.message)
    end
  end
  print(string.rep("=", 36))
  print(string.format("Summary: %d passed | %d failed | %d total", tests_passed, tests_failed, tests_passed + tests_failed))
  return {passed = tests_passed, failed = tests_failed, total = tests_passed + tests_failed}
end

-- ============================================================================
-- Load module
-- ============================================================================

local pack = require("plugins.combat-depth-pack.v1_0_core")
local EncounterTuner = pack.EncounterTuner
local BossRemix = pack.BossRemix
local CompanionDirector = pack.CompanionDirector

-- ============================================================================
-- EncounterTuner tests
-- ============================================================================

local function test_encounter_configure()
  local applied = EncounterTuner.configureEncounterRates({zone = "Mt. Kolts", encounter_rate = 1.2, elite_chance = 0.15})
  assert_test(applied and applied.success ~= false and applied.zone == "Mt. Kolts", "EncounterTuner.configureEncounterRates", "Failed to configure zone rates")
end

local function test_encounter_generate()
  local encounter = EncounterTuner.generateDynamicEncounter({zone = "Narshe", base_xp = 220, xp_multiplier = 1.1}, {avg_level = 24})
  assert_test(encounter and encounter.enemies and encounter.xp_reward and encounter.danger_rating, "EncounterTuner.generateDynamicEncounter", "Failed to generate encounter")
end

local function test_encounter_preset()
  local preset = EncounterTuner.setPreset("story")
  assert_test(preset and preset.success == true and preset.preset == "story", "EncounterTuner.setPreset", "Failed to set preset")
end

-- ============================================================================
-- BossRemix tests
-- ============================================================================

local function test_boss_apply_affixes()
  local remixed = BossRemix.applyAffixes({name = "Atma", hp = 50000, attack = 1200}, {"enraged", "glass_cannon"})
  assert_test(remixed and remixed.hp and remixed.attack and #remixed.affixes == 2, "BossRemix.applyAffixes", "Failed to apply affixes")
end

local function test_boss_generate_plan()
  local plan = BossRemix.generateRemix({name = "Ultros", hp = 30000, attack = 800}, {affixes = {"arcane_shield"}, phases = {{phase = 1, triggers = {hp_percent = 60}, behavior = "Tricky"}}})
  assert_test(plan and plan.phases and #plan.phases >= 1 and plan.loot_bonus, "BossRemix.generateRemix", "Failed to generate remix plan")
end

local function test_boss_presets()
  local presets = BossRemix.listAffixPresets()
  assert_test(presets and presets.starter and #presets.starter >= 1, "BossRemix.listAffixPresets", "Failed to list presets")
end

-- ============================================================================
-- CompanionDirector tests
-- ============================================================================

local function test_companion_define()
  local profile = CompanionDirector.defineProfile("aggressive", {priorities = {"damage", "survival"}, risk_tolerance = "aggressive"})
  assert_test(profile and profile.success == true and profile.risk_tolerance == "aggressive", "CompanionDirector.defineProfile", "Failed to define profile")
end

local function test_companion_evaluate()
  local snapshot = CompanionDirector.evaluateBattleState({hp_status = "critical", threat_level = "high"}, "aggressive")
  assert_test(snapshot and snapshot.recommendations and #snapshot.recommendations >= 1, "CompanionDirector.evaluateBattleState", "Failed to evaluate battle state")
end

local function test_companion_action()
  local action = CompanionDirector.recommendAction({hp_status = "critical", primary_target = "boss"}, "aggressive")
  assert_test(action and action.type and action.target, "CompanionDirector.recommendAction", "Failed to recommend action")
end

-- ============================================================================
-- Runner
-- ============================================================================

local function run_all()
  test_encounter_configure()
  test_encounter_generate()
  test_encounter_preset()
  test_boss_apply_affixes()
  test_boss_generate_plan()
  test_boss_presets()
  test_companion_define()
  test_companion_evaluate()
  test_companion_action()
  return summarize()
end

return {run_all = run_all}
