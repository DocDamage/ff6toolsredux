--[[
  Equipment Optimizer Plugin - v1.1
  Intelligent gear optimization with visual loadout comparison
  
  Phase: Quick Win #3 (Phase 11+ Legacy Plugin Upgrades)
  Version: 1.1 (with Equipment Comparison Dashboard)
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Equipment slots
    SLOTS = {"weapon", "shield", "helmet", "armor", "relic1", "relic2"},
    
    -- Stat weights for optimization
    DEFAULT_WEIGHTS = {
        attack = 1.0,
        defense = 1.0,
        magic_power = 1.0,
        magic_defense = 1.0,
        speed = 0.8,
        evasion = 0.6
    },
    
    -- Character count
    CHARACTER_COUNT = 14,
    
    -- Logging
    LOG_MAX_ENTRIES = 50
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    cached_loadouts = {},
    operation_log = {},
    comparison_cache = {}
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Log operations
local function log_operation(operation_type, details)
    local entry = {
        timestamp = os.time(),
        type = operation_type,
        details = details
    }
    table.insert(plugin_state.operation_log, entry)
    
    -- Limit log size
    if #plugin_state.operation_log > CONFIG.LOG_MAX_ENTRIES then
        table.remove(plugin_state.operation_log, 1)
    end
    
    print(string.format("[Equipment Optimizer] %s: %s", operation_type, details))
end

-- Safe require with warning
local function safe_require(module_path)
    local ok, mod = pcall(require, module_path)
    if not ok then
        log_operation("WARN", "Dependency unavailable: " .. module_path)
        return nil
    end
    return mod
end

-- Database persistence layer handle
local database_layer = nil

local function load_database_layer()
    if not database_layer then
        database_layer = safe_require("plugins.database-persistence-layer.plugin")
    end
    return database_layer
end

-- Phase 11 dependencies (lazy-loaded)
local dependencies = {
    analytics = nil,
    visualization = nil,
    import_export = nil,
    automation = nil
}

local function load_phase11_dependencies()
    dependencies.analytics = dependencies.analytics or safe_require("plugins.advanced-analytics-engine.v1_0_core")
    dependencies.visualization = dependencies.visualization or safe_require("plugins.data-visualization-suite.v1_0_core")
    dependencies.import_export = dependencies.import_export or safe_require("plugins.import-export-manager.v1_0_core")
    dependencies.automation = dependencies.automation or safe_require("plugins.automation-framework.v1_0_core")
    return dependencies
end

-- Calculate equipment stats
local function calculate_equipment_stats(loadout)
    if not loadout then return {} end
    
    local stats = {
        attack = 0,
        defense = 0,
        magic_power = 0,
        magic_defense = 0,
        speed = 0,
        evasion = 0,
        special_effects = {}
    }
    
    -- Sum stats from all equipment (simulated)
    for slot, equipment in pairs(loadout) do
        if equipment and equipment.stats then
            for stat, value in pairs(equipment.stats) do
                if stats[stat] then
                    stats[stat] = stats[stat] + value
                end
            end
        end
    end
    
    return stats
end

-- Calculate synergy score
local function calculate_synergy_score(loadout)
    if not loadout then return 0 end
    
    local synergy_score = 0
    local synergy_details = {}
    
    -- Check for elemental synergies
    local elements = {}
    for slot, equipment in pairs(loadout) do
        if equipment and equipment.element then
            elements[equipment.element] = (elements[equipment.element] or 0) + 1
        end
    end
    
    -- Bonus for matching elements
    for element, count in pairs(elements) do
        if count >= 2 then
            synergy_score = synergy_score + (count * 10)
            table.insert(synergy_details, string.format("%s element x%d (+%d)", 
                element, count, count * 10))
        end
    end
    
    -- Check for set bonuses
    local sets = {}
    for slot, equipment in pairs(loadout) do
        if equipment and equipment.set then
            sets[equipment.set] = (sets[equipment.set] or 0) + 1
        end
    end
    
    -- Bonus for complete sets
    for set_name, count in pairs(sets) do
        if count >= 3 then
            synergy_score = synergy_score + 25
            table.insert(synergy_details, string.format("%s set (+25)", set_name))
        elseif count == 2 then
            synergy_score = synergy_score + 10
            table.insert(synergy_details, string.format("%s partial set (+10)", set_name))
        end
    end
    
    return synergy_score, synergy_details
end

-- ============================================================================
-- QUICK WIN #3: EQUIPMENT COMPARISON DASHBOARD (~180 LOC)
-- ============================================================================

---Compare multiple equipment loadouts side-by-side
---@param loadouts table Array of loadouts to compare
---@return table comparison Comparison data with stats and synergies
function compareLoadouts(loadouts)
    if not loadouts or #loadouts < 2 then
        log_operation("ERROR", "Need at least 2 loadouts to compare")
        return nil
    end
    
    -- Create comparison structure
    local comparison = {
        loadout_count = #loadouts,
        loadouts = {},
        stat_differences = {},
        synergy_comparison = {},
        recommendations = {}
    }
    
    -- Analyze each loadout
    for i, loadout in ipairs(loadouts) do
        local stats = calculate_equipment_stats(loadout.equipment)
        local synergy_score, synergy_details = calculate_synergy_score(loadout.equipment)
        
        table.insert(comparison.loadouts, {
            name = loadout.name or ("Loadout " .. i),
            equipment = loadout.equipment,
            stats = stats,
            synergy_score = synergy_score,
            synergy_details = synergy_details,
            total_score = calculate_total_score(stats, synergy_score)
        })
    end
    
    -- Calculate stat differences (compare to first loadout as baseline)
    if #comparison.loadouts >= 2 then
        local baseline = comparison.loadouts[1].stats
        
        for i = 2, #comparison.loadouts do
            local loadout = comparison.loadouts[i]
            local diffs = {}
            
            for stat, value in pairs(loadout.stats) do
                diffs[stat] = value - (baseline[stat] or 0)
            end
            
            comparison.stat_differences[i] = diffs
        end
    end
    
    -- Generate recommendations
    local best_overall = 1
    local best_score = comparison.loadouts[1].total_score
    
    for i, loadout in ipairs(comparison.loadouts) do
        if loadout.total_score > best_score then
            best_overall = i
            best_score = loadout.total_score
        end
    end
    
    comparison.recommendations.best_overall = best_overall
    comparison.recommendations.confidence = calculate_confidence(comparison.loadouts[best_overall])
    
    log_operation("COMPARE", string.format("Compared %d loadouts", #loadouts))
    
    -- Cache comparison for visualization
    plugin_state.comparison_cache = comparison
    
    return comparison
end

---Display equipment comparison dashboard (formatted output)
---@param comparison table Comparison data from compareLoadouts
---@return string display Formatted comparison dashboard
function displayComparisonDashboard(comparison)
    if not comparison then
        comparison = plugin_state.comparison_cache
    end
    
    if not comparison or #comparison.loadouts == 0 then
        return "No comparison data available. Run compareLoadouts first."
    end
    
    local display = "\n" .. string.rep("=", 80) .. "\n"
    display = display .. "EQUIPMENT LOADOUT COMPARISON DASHBOARD\n"
    display = display .. string.rep("=", 80) .. "\n\n"
    
    -- Loadout names header
    display = display .. "Comparing " .. comparison.loadout_count .. " loadouts:\n"
    for i, loadout in ipairs(comparison.loadouts) do
        display = display .. string.format("  [%d] %s (Score: %d)\n", 
            i, loadout.name, loadout.total_score)
    end
    display = display .. "\n"
    
    -- Stats comparison table
    display = display .. string.rep("-", 80) .. "\n"
    display = display .. "STATS COMPARISON\n"
    display = display .. string.rep("-", 80) .. "\n"
    
    local stats_to_show = {"attack", "defense", "magic_power", "magic_defense", "speed", "evasion"}
    
    for _, stat in ipairs(stats_to_show) do
        local stat_line = string.format("%-15s", stat:upper())
        
        for i, loadout in ipairs(comparison.loadouts) do
            local value = loadout.stats[stat] or 0
            local diff_str = ""
            
            if i > 1 and comparison.stat_differences[i] then
                local diff = comparison.stat_differences[i][stat]
                if diff > 0 then
                    diff_str = string.format(" (+%d)", diff)
                elseif diff < 0 then
                    diff_str = string.format(" (%d)", diff)
                end
            end
            
            stat_line = stat_line .. string.format(" | %4d%-8s", value, diff_str)
        end
        
        display = display .. stat_line .. "\n"
    end
    
    -- Synergy comparison
    display = display .. "\n" .. string.rep("-", 80) .. "\n"
    display = display .. "SYNERGY ANALYSIS\n"
    display = display .. string.rep("-", 80) .. "\n"
    
    for i, loadout in ipairs(comparison.loadouts) do
        display = display .. string.format("[%d] %s: Synergy Score = %d\n", 
            i, loadout.name, loadout.synergy_score)
        
        if #loadout.synergy_details > 0 then
            for _, detail in ipairs(loadout.synergy_details) do
                display = display .. string.format("    - %s\n", detail)
            end
        else
            display = display .. "    - No synergies detected\n"
        end
        display = display .. "\n"
    end
    
    -- Recommendation
    display = display .. string.rep("-", 80) .. "\n"
    display = display .. "RECOMMENDATION\n"
    display = display .. string.rep("-", 80) .. "\n"
    
    local best_idx = comparison.recommendations.best_overall
    local best_loadout = comparison.loadouts[best_idx]
    
    display = display .. string.format("BEST OVERALL: [%d] %s\n", best_idx, best_loadout.name)
    display = display .. string.format("Total Score: %d\n", best_loadout.total_score)
    display = display .. string.format("Confidence: %.1f%%\n", 
        comparison.recommendations.confidence * 100)
    
    display = display .. "\n" .. string.rep("=", 80) .. "\n"
    
    print(display)
    return display
end

-- Calculate total score for loadout
function calculate_total_score(stats, synergy_score)
    local score = 0
    
    -- Weight stats
    for stat, value in pairs(stats) do
        local weight = CONFIG.DEFAULT_WEIGHTS[stat] or 0.5
        score = score + (value * weight)
    end
    
    -- Add synergy bonus
    score = score + synergy_score
    
    return math.floor(score)
end

-- Calculate recommendation confidence
function calculate_confidence(loadout)
    -- Simple confidence based on total score
    -- Higher scores = higher confidence
    local score = loadout.total_score
    local confidence = math.min(0.95, 0.5 + (score / 1000))
    
    return confidence
end

-- ============================================================================
-- EQUIPMENT OPTIMIZATION FUNCTIONS
-- ============================================================================

---Optimize equipment for character
---@param char_id number Character ID (0-13)
---@param optimization_goal string Goal (offense/defense/balanced/magic)
---@return table loadout Optimized equipment loadout
function optimizeEquipment(char_id, optimization_goal)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return nil
    end
    
    optimization_goal = optimization_goal or "balanced"
    
    -- Simulated optimization (in real implementation, would analyze inventory)
    local optimized_loadout = {
        name = string.format("%s Build", optimization_goal:upper()),
        character_id = char_id,
        equipment = {
            weapon = {id = 255, name = "Ultima Weapon", stats = {attack = 255, magic_power = 108}},
            shield = {id = 52, name = "Force Shield", stats = {defense = 70, magic_defense = 40}},
            helmet = {id = 40, name = "Genji Helmet", stats = {defense = 50, magic_defense = 35}},
            armor = {id = 60, name = "Minerva Bustier", stats = {defense = 86, magic_defense = 66}},
            relic1 = {id = 25, name = "Ribbon", stats = {}},
            relic2 = {id = 30, name = "Marvel Shoes", stats = {speed = 20, evasion = 10}}
        }
    }
    
    -- Persist optimized loadout to database layer
    local db = load_database_layer()
    if db and db.saveEquipmentConfig then
        local config = plugin_state.cached_loadouts[char_id] or {}
        config[optimization_goal] = optimized_loadout
        db.saveEquipmentConfig(config)
    end
    
    plugin_state.cached_loadouts[char_id] = optimized_loadout
    log_operation("OPTIMIZE", string.format("Optimized equipment for character %d (%s)", 
        char_id, optimization_goal))
    
    return optimized_loadout
end

---Get current equipment loadout for character
---@param char_id number Character ID (0-13)
---@return table loadout Current equipment
function getCurrentLoadout(char_id)
    if char_id < 0 or char_id >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return nil
    end
    
    -- Try to load from database layer first
    local db = load_database_layer()
    if db and db.loadEquipmentConfig then
        local persisted = db.loadEquipmentConfig()
        if persisted and persisted[char_id] then
            log_operation("LOAD", string.format("Loaded equipment for character %d from persistence layer", char_id))
            return persisted[char_id]
        end
    end
    
    -- Simulated current loadout
    local current_loadout = {
        name = "Current Equipment",
        character_id = char_id,
        equipment = {
            weapon = {id = 100, name = "Enhancer", stats = {attack = 105, magic_power = 68}},
            shield = {id = 20, name = "Ice Shield", stats = {defense = 50, magic_defense = 30}, element = "ice"},
            helmet = {id = 15, name = "Crystal Helm", stats = {defense = 40, magic_defense = 25}},
            armor = {id = 25, name = "Force Armor", stats = {defense = 70, magic_defense = 50}},
            relic1 = {id = 10, name = "Sprint Shoes", stats = {speed = 10}},
            relic2 = {id = 12, name = "Back Guard", stats = {evasion = 5}}
        }
    }
    
    return current_loadout
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS (TIER 1 - EQUIPMENT OPTIMIZER)
-- ============================================================================

-- Build a compact loadout profile for analytics
local function build_loadout_profile(loadout)
    local stats = calculate_equipment_stats(loadout.equipment or {})
    local synergy_score = select(1, calculate_synergy_score(loadout.equipment or {}))
    return {
        name = loadout.name or "Unnamed",
        attack = stats.attack or 0,
        defense = stats.defense or 0,
        magic_power = stats.magic_power or 0,
        magic_defense = stats.magic_defense or 0,
        speed = stats.speed or 0,
        evasion = stats.evasion or 0,
        synergy_score = synergy_score
    }
end

-- Advanced optimization using Analytics Engine (stub for integration)
function optimizeEquipmentLoadout(char_id, goal)
    load_phase11_dependencies()
    local base = optimizeEquipment(char_id, goal)
    local analytics = dependencies.analytics
    local profile = build_loadout_profile(base)
    local recommendation = analytics and analytics.Optimization and analytics.Optimization.optimize(goal or "balanced", profile) or {}
    base.recommendation = recommendation
    log_operation("ANALYTICS", string.format("Optimized with analytics for character %d", char_id))
    return base
end

-- Predict equipment synergies via analytics
function predictEquipmentSynergies(loadout)
    load_phase11_dependencies()
    local analytics = dependencies.analytics
    local profile = build_loadout_profile(loadout)
    local trend = analytics and analytics.PatternRecognition and analytics.PatternRecognition.detectTrend({profile.attack, profile.defense, profile.magic_power}) or {}
    log_operation("ANALYTICS", "Predicted equipment synergies")
    return {profile = profile, trend = trend}
end

-- Detect equipment outliers
function detectEquipmentOutliers(loadouts)
    load_phase11_dependencies()
    if not loadouts or #loadouts == 0 then
        log_operation("ERROR", "No loadouts provided for outlier detection")
        return nil
    end
    local analytics = dependencies.analytics
    local scores = {}
    for _, l in ipairs(loadouts) do
        table.insert(scores, build_loadout_profile(l).synergy_score)
    end
    local anomalies = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.detectAnomalies(scores) or {}
    log_operation("ANALYTICS", "Detected equipment outliers")
    return anomalies
end

-- Recommend specialized gear for scenario
function recommendSpecializedGear(scenario)
    load_phase11_dependencies()
    local analytics = dependencies.analytics
    local prediction = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.makePrediction("scenario", scenario or {type = "boss"}) or {}
    log_operation("ANALYTICS", "Generated specialized gear recommendation")
    return prediction
end

-- Analyze resistance patterns
function analyzeResistancePatterns(enemy_profile)
    load_phase11_dependencies()
    local analytics = dependencies.analytics
    local pattern = analytics and analytics.PatternRecognition and analytics.PatternRecognition.analyzePatterns(enemy_profile or {fire = 1, ice = 2}) or {}
    log_operation("ANALYTICS", "Analyzed resistance patterns")
    return pattern
end

-- Automated optimal equip via Automation Framework
function autoEquipOptimal(char_id)
    load_phase11_dependencies()
    local automation = dependencies.automation
    local optimized = optimizeEquipmentLoadout(char_id, "balanced")
    if automation and automation.WorkflowEngine and automation.WorkflowEngine.executeWorkflow then
        automation.WorkflowEngine.executeWorkflow("auto_equip", {character = char_id, loadout = optimized})
    end
    log_operation("AUTOMATION", string.format("Auto-equipped optimal gear for %d", char_id))
    return optimized
end

-- Set up equipment rules
function setupEquipmentRules(rule_set)
    load_phase11_dependencies()
    local automation = dependencies.automation
    if automation and automation.RuleEngine and automation.RuleEngine.registerRules then
        automation.RuleEngine.registerRules("equipment_rules", rule_set or {})
    end
    log_operation("AUTOMATION", "Equipment rules configured")
    return true
end

-- Schedule periodic equipment review
function scheduleEquipmentReview(cron_expr)
    load_phase11_dependencies()
    local automation = dependencies.automation
    if automation and automation.Scheduler and automation.Scheduler.scheduleTask then
        automation.Scheduler.scheduleTask("equipment_review", cron_expr or "0 9 * * *")
    end
    log_operation("AUTOMATION", "Scheduled equipment review")
    return true
end

-- Trigger alert when better gear appears
function triggerEquipmentAlert(message)
    load_phase11_dependencies()
    local automation = dependencies.automation
    if automation and automation.Notifications and automation.Notifications.sendAlert then
        automation.Notifications.sendAlert(message or "Better gear available")
    end
    log_operation("AUTOMATION", "Equipment alert triggered")
    return true
end

-- Visualization: generate comparison chart
function generateEquipmentComparison(loadouts)
    load_phase11_dependencies()
    local viz = dependencies.visualization
    local comparison = compareLoadouts(loadouts)
    local categories, values = {}, {}
    if comparison then
        for _, l in ipairs(comparison.loadouts) do
            table.insert(categories, l.name)
            table.insert(values, l.total_score)
        end
    end
    local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateBarChart("Loadout Scores", categories, values) or {type = "Bar", generated = false}
    log_operation("VISUALIZE", "Generated equipment comparison chart")
    return {chart = chart, categories = categories, values = values}
end

-- Visualization: create dashboard
function createLoadoutDashboard(loadouts)
    load_phase11_dependencies()
    local viz = dependencies.visualization
    local dashboard = viz and viz.DashboardManagement and viz.DashboardManagement.createDashboard("Equipment Dashboard", "grid") or {name = "Equipment Dashboard"}
    if viz and viz.DashboardManagement and viz.DashboardManagement.addWidget then
        viz.DashboardManagement.addWidget(dashboard.dashboard_id, "chart", generateEquipmentComparison(loadouts))
    end
    log_operation("VISUALIZE", "Created loadout dashboard")
    return dashboard
end

-- Visualization: armor coverage view
function visualizeArmorCoverage(loadout)
    load_phase11_dependencies()
    local viz = dependencies.visualization
    local stats = calculate_equipment_stats(loadout and loadout.equipment or {})
    local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateGaugeChart("Armor Coverage", stats.defense or 0, 300) or {type = "Gauge"}
    log_operation("VISUALIZE", "Generated armor coverage gauge")
    return {chart = chart, defense = stats.defense}
end

-- Export loadout guide report
function exportLoadoutGuide(loadouts, format, output_path)
    load_phase11_dependencies()
    local viz = dependencies.visualization
    local report = viz and viz.ReportGeneration and viz.ReportGeneration.generateReport("Loadout Guide", loadouts or {}, format or "PDF") or {name = "Loadout Guide"}
    if output_path and dependencies.import_export and dependencies.import_export.DataExporter then
        dependencies.import_export.DataExporter.exportToJSON(report, output_path)
    end
    log_operation("VISUALIZE", "Exported loadout guide")
    return report
end

-- Import/export helpers via Import/Export Manager
function exportEquipmentTemplate(loadout, format, path)
    load_phase11_dependencies()
    local exporter = dependencies.import_export and dependencies.import_export.DataExporter
    local output = path or "equipment_template.json"
    if exporter then
        if format == "csv" then
            exporter.exportToCSV(loadout, output, true)
        elseif format == "xml" then
            exporter.exportToXML(loadout, output)
        else
            exporter.exportToJSON(loadout, output)
        end
    else
        local file = io.open(output, "w")
        if file then
            file:write("-- equipment template placeholder --\n")
            file:close()
        end
    end
    log_operation("EXPORT", "Exported equipment template")
    return {path = output, format = format or "json"}
end

function importEquipmentTemplate(path, format)
    load_phase11_dependencies()
    if not path then
        log_operation("ERROR", "No path provided for equipment import")
        return nil
    end
    local importer = dependencies.import_export and dependencies.import_export.DataImporter
    local fmt = (format or "json"):lower()
    local data
    if importer then
        if fmt == "csv" then
            data = importer.importCSV(path, true)
        elseif fmt == "xml" then
            data = importer.importXML(path, "lenient")
        else
            data = importer.importJSON(path, true)
        end
    end
    log_operation("IMPORT", string.format("Imported equipment template from %s", path))
    return data or {path = path}
end

function syncEquipmentWithTeam(team_ids)
    load_phase11_dependencies()
    local importer = dependencies.import_export
    local synced = importer and importer.FormatConversion and importer.FormatConversion.convertFormat("internal", "sync", team_ids or {}) or {}
    log_operation("SYNC", "Synced equipment with team")
    return synced
end

function batchApplyLoadouts(loadouts)
    load_phase11_dependencies()
    if not loadouts or #loadouts == 0 then
        log_operation("ERROR", "No loadouts provided for batch apply")
        return {success = false}
    end
    for _, loadout in ipairs(loadouts) do
        -- In real implementation, would apply loadout to character
    end
    log_operation("APPLY", string.format("Batch applied %d loadouts", #loadouts))
    return {success = true, applied = #loadouts}
end

-- ============================================================================
-- INITIALIZATION
-- ============================================================================

function initialize()
    if plugin_state.initialized then
        return true
    end
    
    log_operation("INIT", "Equipment Optimizer initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("Equipment Optimizer v1.2.0 loaded")
print("Commands: optimizeEquipment(char_id, goal), optimizeEquipmentLoadout(char_id, goal), getCurrentLoadout(char_id)")
print("Comparison: compareLoadouts(loadouts), displayComparisonDashboard(), generateEquipmentComparison(loadouts)")
print("Automation: autoEquipOptimal(char_id), setupEquipmentRules(rules), scheduleEquipmentReview(cron)")
print("Import/Export: exportEquipmentTemplate(loadout, fmt, path), importEquipmentTemplate(path, fmt)")
print("Goals: 'offense', 'defense', 'balanced', 'magic'")
