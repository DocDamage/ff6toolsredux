--[[
  Phase C (Tier 2 Phase 2) Plugin Smoke Tests
  Tests Phase 11 integrations for 6 analytics plugins
  
  Plugins Tested:
  1. Build Optimizer (v1.1)
  2. Strategy Library (v1.1)
  3. Performance Profiler (v1.1)
  4. Randomizer Assistant (v1.1)
  5. World State Manipulator (v1.1)
  6. Custom Report Generator (v1.1)
  
  Total Test Count: 30 tests (5 per plugin)
]]

local PhaseC_Tests = {}

-- Test tracking
local tests_passed = 0
local tests_failed = 0
local test_results = {}

---Test assertion helper
---@param condition boolean Test condition
---@param test_name string Test name
---@param error_message string Error description
local function assert_test(condition, test_name, error_message)
  if condition then
    tests_passed = tests_passed + 1
    table.insert(test_results, {name = test_name, status = "✓ PASS", message = ""})
  else
    tests_failed = tests_failed + 1
    table.insert(test_results, {name = test_name, status = "✗ FAIL", message = error_message or "Assertion failed"})
  end
end

-- ============================================================================
-- BUILD OPTIMIZER TESTS (5 tests)
-- ============================================================================

function PhaseC_Tests.test_build_optimizer_ml_prediction()
  local build_optimizer = require("plugins.build-optimizer.v1_0_core")
  
  local prediction = build_optimizer.Phase11Integration.predictBuildPerformance({
    character = "Terra",
    equipment = {"Illumina", "Crystal Shield", "Minerva"},
    stats = {hp = 8500, mp = 600, strength = 220, magic = 255}
  })
  
  assert_test(
    prediction and prediction.predicted_score and prediction.predicted_score > 0,
    "Build Optimizer: ML Performance Prediction",
    "Failed to predict build performance"
  )
end

function PhaseC_Tests.test_build_optimizer_visualization()
  local build_optimizer = require("plugins.build-optimizer.v1_0_core")
  
  local dashboard = build_optimizer.Phase11Integration.createProgressionDashboard({
    character = "Edgar",
    progression_data = {{level = 10, power = 120}, {level = 20, power = 240}}
  })
  
  assert_test(
    dashboard and dashboard.dashboard_id and #dashboard.widgets > 0,
    "Build Optimizer: Visualization Dashboard",
    "Failed to create progression dashboard"
  )
end

function PhaseC_Tests.test_build_optimizer_import_export()
  local build_optimizer = require("plugins.build-optimizer.v1_0_core")
  
  local export_result = build_optimizer.Phase11Integration.exportBuildTemplate({
    build_id = "TEST_BUILD_001",
    character = "Celes",
    equipment = {"Lightbringer", "Genji Glove"}
  })
  
  assert_test(
    export_result and export_result.success and export_result.filepath,
    "Build Optimizer: Build Template Export",
    "Failed to export build template"
  )
end

function PhaseC_Tests.test_build_optimizer_pattern_analysis()
  local build_optimizer = require("plugins.build-optimizer.v1_0_core")
  
  local patterns = build_optimizer.Phase11Integration.analyzeSuccessPatterns({
    {build_id = "B1", success_rate = 85},
    {build_id = "B2", success_rate = 92},
    {build_id = "B3", success_rate = 78}
  })
  
  assert_test(
    patterns and patterns.patterns_found and patterns.patterns_found > 0,
    "Build Optimizer: Success Pattern Analysis",
    "Failed to analyze success patterns"
  )
end

function PhaseC_Tests.test_build_optimizer_automation()
  local build_optimizer = require("plugins.build-optimizer.v1_0_core")
  
  local automation = build_optimizer.Phase11Integration.autoOptimizeBuild({
    character = "Locke",
    goal = "maximize_damage"
  })
  
  assert_test(
    automation and automation.optimizations_applied and automation.optimizations_applied > 0,
    "Build Optimizer: Auto-Optimization",
    "Failed to apply build optimizations"
  )
end

-- ============================================================================
-- STRATEGY LIBRARY TESTS (5 tests)
-- ============================================================================

function PhaseC_Tests.test_strategy_library_effectiveness()
  local strategy_library = require("plugins.strategy-library.v1_0_core")
  
  local analysis = strategy_library.Phase11Integration.analyzeStrategyEffectiveness({
    strategy_id = "STRAT_001",
    usage_count = 45,
    success_rate = 88,
    boss = "Kefka"
  })
  
  assert_test(
    analysis and analysis.effectiveness_score and analysis.effectiveness_score > 0,
    "Strategy Library: Effectiveness Analysis",
    "Failed to analyze strategy effectiveness"
  )
end

function PhaseC_Tests.test_strategy_library_prediction()
  local strategy_library = require("plugins.strategy-library.v1_0_core")
  
  local prediction = strategy_library.Phase11Integration.predictStrategySuccess({
    strategy = "Quick Sketch + Relic Bug",
    boss = "Wrexsoul",
    party_level = 45
  })
  
  assert_test(
    prediction and prediction.predicted_success and prediction.confidence > 0,
    "Strategy Library: Success Prediction",
    "Failed to predict strategy success"
  )
end

function PhaseC_Tests.test_strategy_library_similarity()
  local strategy_library = require("plugins.strategy-library.v1_0_core")
  
  local similar = strategy_library.Phase11Integration.findSimilarStrategies({
    strategy_id = "STRAT_BERSERK",
    tags = {"offensive", "status_effect"}
  }, {
    {strategy_id = "STRAT_RAGE", tags = {"offensive", "berserk"}},
    {strategy_id = "STRAT_QUICK", tags = {"offensive", "haste"}}
  })
  
  assert_test(
    similar and similar.similar_count and similar.similar_count > 0,
    "Strategy Library: Similarity Search",
    "Failed to find similar strategies"
  )
end

function PhaseC_Tests.test_strategy_library_export()
  local strategy_library = require("plugins.strategy-library.v1_0_core")
  
  local export_result = strategy_library.Phase11Integration.exportStrategies({
    {strategy_id = "S1", name = "Vanish-Doom"},
    {strategy_id = "S2", name = "Quick Ultima"}
  }, "json")
  
  assert_test(
    export_result and export_result.success and export_result.strategies_exported > 0,
    "Strategy Library: Strategy Export",
    "Failed to export strategies"
  )
end

function PhaseC_Tests.test_strategy_library_auto_suggest()
  local strategy_library = require("plugins.strategy-library.v1_0_core")
  
  local suggestions = strategy_library.Phase11Integration.autoSuggestStrategies({
    boss = "Atma Weapon",
    party_composition = {"Terra", "Celes", "Edgar", "Sabin"}
  })
  
  assert_test(
    suggestions and suggestions.suggestions_count and suggestions.suggestions_count > 0,
    "Strategy Library: Auto-Suggest Strategies",
    "Failed to suggest strategies"
  )
end

-- ============================================================================
-- PERFORMANCE PROFILER TESTS (5 tests)
-- ============================================================================

function PhaseC_Tests.test_performance_profiler_ml_prediction()
  local performance_profiler = require("plugins.performance-profiler.v1_0_core")
  
  local prediction = performance_profiler.Phase11Integration.predictDegradation({
    {timestamp = 1, performance_score = 95},
    {timestamp = 2, performance_score = 92},
    {timestamp = 3, performance_score = 90}
  })
  
  assert_test(
    prediction and prediction.predicted_future and #prediction.predicted_future > 0,
    "Performance Profiler: Degradation Prediction",
    "Failed to predict performance degradation"
  )
end

function PhaseC_Tests.test_performance_profiler_dashboard()
  local performance_profiler = require("plugins.performance-profiler.v1_0_core")
  
  local dashboard = performance_profiler.Phase11Integration.createPerformanceDashboard({
    {action = "Battle", execution_time = 150},
    {action = "Menu", execution_time = 50},
    {action = "Save", execution_time = 200}
  })
  
  assert_test(
    dashboard and dashboard.dashboard_id and #dashboard.widgets >= 3,
    "Performance Profiler: Performance Dashboard",
    "Failed to create performance dashboard"
  )
end

function PhaseC_Tests.test_performance_profiler_monitoring()
  local performance_profiler = require("plugins.performance-profiler.v1_0_core")
  
  local monitoring = performance_profiler.Phase11Integration.enableRealTimeMonitoring("SESSION_TEST")
  
  assert_test(
    monitoring and monitoring.enabled and #monitoring.metrics > 0,
    "Performance Profiler: Real-Time Monitoring",
    "Failed to enable real-time monitoring"
  )
end

function PhaseC_Tests.test_performance_profiler_bottleneck_detection()
  local performance_profiler = require("plugins.performance-profiler.v1_0_core")
  
  -- Test existing bottleneck detection
  local bottlenecks = performance_profiler.BottleneckDetection.identifyBottlenecks({
    {action = "DatabaseQuery", duration = 500},
    {action = "Render", duration = 800},
    {action = "Compute", duration = 1200}
  })
  
  assert_test(
    bottlenecks and #bottlenecks > 0,
    "Performance Profiler: Bottleneck Detection",
    "Failed to identify bottlenecks"
  )
end

function PhaseC_Tests.test_performance_profiler_benchmarking()
  local performance_profiler = require("plugins.performance-profiler.v1_0_core")
  
  local benchmark = performance_profiler.BenchmarkComparison.runBenchmark({
    name = "Load Save Test",
    iterations = 10
  })
  
  assert_test(
    benchmark and benchmark.average_time and benchmark.average_time > 0,
    "Performance Profiler: Benchmark Execution",
    "Failed to run benchmark"
  )
end

-- ============================================================================
-- RANDOMIZER ASSISTANT TESTS (5 tests)
-- ============================================================================

function PhaseC_Tests.test_randomizer_seed_difficulty()
  local Phase11Integration = require("plugins.randomizer-assistant.plugin")
  
  local analysis = Phase11Integration.analyzeSeedDifficulty({
    seed_id = "ABC123",
    flags = {enemy_scaling = 1.5, randomization = 7, items = 2}
  })
  
  assert_test(
    analysis and analysis.predicted_difficulty and analysis.predicted_difficulty > 0,
    "Randomizer Assistant: Seed Difficulty Analysis",
    "Failed to analyze seed difficulty"
  )
end

function PhaseC_Tests.test_randomizer_similarity_search()
  local Phase11Integration = require("plugins.randomizer-assistant.plugin")
  
  local similar = Phase11Integration.findSimilarSeeds(
    {seed_id = "TEST_SEED"},
    {
      {seed_id = "SEED_A", difficulty = 7},
      {seed_id = "SEED_B", difficulty = 8}
    }
  )
  
  assert_test(
    similar and similar.similar_count and similar.similar_count > 0,
    "Randomizer Assistant: Similar Seed Search",
    "Failed to find similar seeds"
  )
end

function PhaseC_Tests.test_randomizer_community_sharing()
  local Phase11Integration = require("plugins.randomizer-assistant.plugin")
  
  local share_result = Phase11Integration.shareSeedToCommunity({
    seed_id = "SHARE_TEST",
    flags = {randomization = 5}
  }, false)
  
  assert_test(
    share_result and share_result.success and share_result.share_url,
    "Randomizer Assistant: Community Sharing",
    "Failed to share seed to community"
  )
end

function PhaseC_Tests.test_randomizer_visualization()
  local Phase11Integration = require("plugins.randomizer-assistant.plugin")
  
  local visualization = Phase11Integration.visualizeSeedCharacteristics({
    seed_id = "VIZ_TEST"
  })
  
  assert_test(
    visualization and visualization.type == "radar" and #visualization.axes > 0,
    "Randomizer Assistant: Seed Visualization",
    "Failed to visualize seed characteristics"
  )
end

function PhaseC_Tests.test_randomizer_confidence_scoring()
  local Phase11Integration = require("plugins.randomizer-assistant.plugin")
  
  local analysis = Phase11Integration.analyzeSeedDifficulty({
    seed_id = "CONF_TEST",
    flags = {}
  })
  
  assert_test(
    analysis and analysis.confidence and analysis.confidence > 0,
    "Randomizer Assistant: Confidence Scoring",
    "Failed to generate confidence score"
  )
end

-- ============================================================================
-- WORLD STATE MANIPULATOR TESTS (5 tests)
-- ============================================================================

function PhaseC_Tests.test_world_state_snapshot_creation()
  local Phase11Integration = require("plugins.world-state-manipulator.plugin")
  
  local snapshot = Phase11Integration.createStateSnapshot({
    world = "WoR",
    characters_recruited = 10
  }, "TEST_SNAPSHOT")
  
  assert_test(
    snapshot and snapshot.success and snapshot.snapshot_id,
    "World State Manipulator: Snapshot Creation",
    "Failed to create state snapshot"
  )
end

function PhaseC_Tests.test_world_state_restoration()
  local Phase11Integration = require("plugins.world-state-manipulator.plugin")
  
  local restore = Phase11Integration.restoreFromSnapshot("SNAP_12345")
  
  assert_test(
    restore and restore.success and restore.snapshot_id,
    "World State Manipulator: State Restoration",
    "Failed to restore from snapshot"
  )
end

function PhaseC_Tests.test_world_state_api_exposure()
  local Phase11Integration = require("plugins.world-state-manipulator.plugin")
  
  local api_config = Phase11Integration.exposeStateAPI({
    {path = "getState", method = "GET"},
    {path = "updateState", method = "POST"}
  })
  
  assert_test(
    api_config and api_config.endpoints_registered >= 2,
    "World State Manipulator: API Exposure",
    "Failed to expose state API"
  )
end

function PhaseC_Tests.test_world_state_automation()
  local Phase11Integration = require("plugins.world-state-manipulator.plugin")
  
  local automation = Phase11Integration.automateStateChanges({
    triggers = {
      {event_type = "character_recruited"},
      {event_type = "boss_defeated"}
    }
  })
  
  assert_test(
    automation and automation.triggers_configured >= 2,
    "World State Manipulator: State Automation",
    "Failed to configure automation"
  )
end

function PhaseC_Tests.test_world_state_comparison()
  local Phase11Integration = require("plugins.world-state-manipulator.plugin")
  
  local comparison = Phase11Integration.compareStates("SNAP_001", "SNAP_002")
  
  assert_test(
    comparison and comparison.differences_found >= 0,
    "World State Manipulator: State Comparison",
    "Failed to compare states"
  )
end

-- ============================================================================
-- CUSTOM REPORT GENERATOR TESTS (5 tests)
-- ============================================================================

function PhaseC_Tests.test_report_visual_generation()
  local report_generator = require("plugins.custom-report-generator.v1_0_core")
  
  local visual_report = report_generator.Phase11Integration.generateVisualReport({
    title = "Test Report",
    data_series = {{name = "Series 1", values = {10, 20, 30}}}
  }, {"bar", "line"})
  
  assert_test(
    visual_report and visual_report.report_id and #visual_report.charts >= 2,
    "Custom Report Generator: Visual Report Generation",
    "Failed to generate visual report"
  )
end

function PhaseC_Tests.test_report_interactive_dashboard()
  local report_generator = require("plugins.custom-report-generator.v1_0_core")
  
  local dashboard = report_generator.Phase11Integration.createInteractiveDashboard({
    title = "Test Dashboard",
    layout = "grid",
    widgets = {
      {type = "chart", title = "Widget 1"},
      {type = "gauge", title = "Widget 2"}
    }
  })
  
  assert_test(
    dashboard and dashboard.interactive and #dashboard.widgets >= 2,
    "Custom Report Generator: Interactive Dashboard",
    "Failed to create interactive dashboard"
  )
end

function PhaseC_Tests.test_report_multi_format_export()
  local report_generator = require("plugins.custom-report-generator.v1_0_core")
  
  local export_result = report_generator.Phase11Integration.exportMultiFormat({
    report_id = "REPORT_001",
    title = "Test Report"
  }, {"json", "csv"}, "./test_reports")
  
  assert_test(
    export_result and export_result.total_exports >= 2,
    "Custom Report Generator: Multi-Format Export",
    "Failed to export report in multiple formats"
  )
end

function PhaseC_Tests.test_report_ai_insights()
  local report_generator = require("plugins.custom-report-generator.v1_0_core")
  
  local insights = report_generator.Phase11Integration.generateAIInsights({
    data = {
      {metric = "performance", value = 85},
      {metric = "efficiency", value = 92}
    }
  })
  
  assert_test(
    insights and insights.insights_count > 0 and #insights.key_findings > 0,
    "Custom Report Generator: AI-Powered Insights",
    "Failed to generate AI insights"
  )
end

function PhaseC_Tests.test_report_scheduling()
  local report_generator = require("plugins.custom-report-generator.v1_0_core")
  
  local scheduled = report_generator.Phase11Integration.scheduleReports({
    frequency = "weekly",
    recipients = {"user@example.com"}
  })
  
  assert_test(
    scheduled and scheduled.schedule_id and scheduled.enabled,
    "Custom Report Generator: Report Scheduling",
    "Failed to schedule reports"
  )
end

-- ============================================================================
-- TEST RUNNER
-- ============================================================================

function PhaseC_Tests.run_all()
  print("=" .. string.rep("=", 79))
  print("  PHASE C (TIER 2 PHASE 2) SMOKE TESTS - 6 PLUGINS")
  print("=" .. string.rep("=", 79))
  print()
  
  -- Build Optimizer Tests
  print("[1/6] Build Optimizer v1.1 Tests...")
  PhaseC_Tests.test_build_optimizer_ml_prediction()
  PhaseC_Tests.test_build_optimizer_visualization()
  PhaseC_Tests.test_build_optimizer_import_export()
  PhaseC_Tests.test_build_optimizer_pattern_analysis()
  PhaseC_Tests.test_build_optimizer_automation()
  print()
  
  -- Strategy Library Tests
  print("[2/6] Strategy Library v1.1 Tests...")
  PhaseC_Tests.test_strategy_library_effectiveness()
  PhaseC_Tests.test_strategy_library_prediction()
  PhaseC_Tests.test_strategy_library_similarity()
  PhaseC_Tests.test_strategy_library_export()
  PhaseC_Tests.test_strategy_library_auto_suggest()
  print()
  
  -- Performance Profiler Tests
  print("[3/6] Performance Profiler v1.1 Tests...")
  PhaseC_Tests.test_performance_profiler_ml_prediction()
  PhaseC_Tests.test_performance_profiler_dashboard()
  PhaseC_Tests.test_performance_profiler_monitoring()
  PhaseC_Tests.test_performance_profiler_bottleneck_detection()
  PhaseC_Tests.test_performance_profiler_benchmarking()
  print()
  
  -- Randomizer Assistant Tests
  print("[4/6] Randomizer Assistant v1.1 Tests...")
  PhaseC_Tests.test_randomizer_seed_difficulty()
  PhaseC_Tests.test_randomizer_similarity_search()
  PhaseC_Tests.test_randomizer_community_sharing()
  PhaseC_Tests.test_randomizer_visualization()
  PhaseC_Tests.test_randomizer_confidence_scoring()
  print()
  
  -- World State Manipulator Tests
  print("[5/6] World State Manipulator v1.1 Tests...")
  PhaseC_Tests.test_world_state_snapshot_creation()
  PhaseC_Tests.test_world_state_restoration()
  PhaseC_Tests.test_world_state_api_exposure()
  PhaseC_Tests.test_world_state_automation()
  PhaseC_Tests.test_world_state_comparison()
  print()
  
  -- Custom Report Generator Tests
  print("[6/6] Custom Report Generator v1.1 Tests...")
  PhaseC_Tests.test_report_visual_generation()
  PhaseC_Tests.test_report_interactive_dashboard()
  PhaseC_Tests.test_report_multi_format_export()
  PhaseC_Tests.test_report_ai_insights()
  PhaseC_Tests.test_report_scheduling()
  print()
  
  -- Display Results
  print("=" .. string.rep("=", 79))
  print("  TEST RESULTS")
  print("=" .. string.rep("=", 79))
  print()
  
  for _, result in ipairs(test_results) do
    print(string.format("%s %s", result.status, result.name))
    if result.message ~= "" then
      print(string.format("    → %s", result.message))
    end
  end
  
  print()
  print("=" .. string.rep("=", 79))
  print(string.format("  SUMMARY: %d PASSED | %d FAILED | %d TOTAL", 
    tests_passed, tests_failed, tests_passed + tests_failed))
  print("=" .. string.rep("=", 79))
  
  if tests_failed == 0 then
    print()
    print("✓ ALL PHASE C SMOKE TESTS PASSED!")
    print("  All 6 plugins have functional Phase 11 integrations")
    print("  Ready for production deployment")
  else
    print()
    print("⚠ SOME TESTS FAILED - REVIEW REQUIRED")
  end
  
  return {
    passed = tests_passed,
    failed = tests_failed,
    total = tests_passed + tests_failed,
    success_rate = (tests_passed / (tests_passed + tests_failed)) * 100
  }
end

return PhaseC_Tests
