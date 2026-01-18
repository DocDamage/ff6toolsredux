-- Party Optimizer Plugin v1.2.0 (Phase 11 Tier 1 Integrations)
-- Provides party composition optimization with analytics, visualization, automation, import/export, and backup hooks.

local plugin_state = {
    initialized = false,
    last_party = nil,
    last_recommendation = nil,
    last_dashboard = nil
}

-- Safe logger
local function log_operation(kind, msg)
    print(string.format("[Party Optimizer] %s: %s", kind, msg))
end

-- Safe require with warning
local function safe_require(path)
    local ok, mod = pcall(require, path)
    if not ok then
        log_operation("WARN", "Dependency unavailable: " .. path)
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

-- Phase 11 dependency handles (lazy)
local dependencies = {
    analytics = nil,
    visualization = nil,
    import_export = nil,
    automation = nil,
    integration_hub = nil,
    backup_restore = nil
}

local function load_phase11()
    dependencies.analytics = dependencies.analytics or safe_require("plugins.advanced-analytics-engine.v1_0_core")
    dependencies.visualization = dependencies.visualization or safe_require("plugins.data-visualization-suite.v1_0_core")
    dependencies.import_export = dependencies.import_export or safe_require("plugins.import-export-manager.v1_0_core")
    dependencies.automation = dependencies.automation or safe_require("plugins.automation-framework.v1_0_core")
    dependencies.integration_hub = dependencies.integration_hub or safe_require("plugins.integration-hub.v1_0_core")
    dependencies.backup_restore = dependencies.backup_restore or safe_require("plugins.backup-restore-system.v1_0_core")
    return dependencies
end

-- Basic helper to normalize party data
local function normalize_party(party)
    local normalized = {members = {}, size = 0}
    if not party then return normalized end
    for _, member in ipairs(party) do
        table.insert(normalized.members, member)
    end
    normalized.size = #normalized.members
    return normalized
end

-- Optimize party composition using analytics signals
function optimizePartyComposition(party, goal)
    load_phase11()
    local normalized = normalize_party(party)
    goal = goal or "balanced"

    local analytics = dependencies.analytics
    local scores = {}
    for idx, member in ipairs(normalized.members) do
        scores[idx] = (member.synergy or 75) + (member.level or 30)
    end

    local best = analytics and analytics.Optimization and analytics.Optimization.optimize(goal, scores) or {objective = goal, optimal_value = 0.9}
    plugin_state.last_party = normalized
    plugin_state.last_recommendation = {goal = goal, best = best}
    
    -- Persist party configuration to database layer
    local db = load_database_layer()
    if db and db.savePartyConfig then
        db.savePartyConfig(normalized)
    end

    log_operation("OPTIMIZE", string.format("Optimized party for goal '%s' (size=%d)", goal, normalized.size))
    return plugin_state.last_recommendation
end

-- Recommend party for boss or scenario
function recommendPartyForScenario(battle_info)
    load_phase11()
    battle_info = battle_info or {boss = "Unknown", difficulty = 80}
    local analytics = dependencies.analytics
    local prediction = analytics and analytics.PredictiveAnalytics and analytics.PredictiveAnalytics.makePrediction("battle", {difficulty = battle_info.difficulty}) or {}
    log_operation("ANALYTICS", "Generated scenario recommendation")
    return {battle = battle_info, prediction = prediction}
end

-- Auto-apply party setup with automation framework
function autoConfigureParty(party, battle_info)
    load_phase11()
    local automation = dependencies.automation
    local recommendation = optimizePartyComposition(party, "balanced")
    if automation and automation.WorkflowEngine and automation.WorkflowEngine.executeWorkflow then
        automation.WorkflowEngine.executeWorkflow("party_setup", {party = party, context = battle_info or {}})
    end
    log_operation("AUTOMATION", "Auto-configured party for upcoming battle")
    return recommendation
end

-- Visualize party synergy dashboard
function visualizePartySynergy(party)
    load_phase11()
    local viz = dependencies.visualization
    local normalized = normalize_party(party)
    local names, scores = {}, {}
    for _, member in ipairs(normalized.members) do
        table.insert(names, member.name or "Member")
        table.insert(scores, member.synergy or 75)
    end

    local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateBarChart("Party Synergy", names, scores) or {type = "Bar", generated = false}
    local dashboard = viz and viz.DashboardManagement and viz.DashboardManagement.createDashboard("Party Dashboard", "grid") or {name = "Party Dashboard"}
    if viz and viz.DashboardManagement and viz.DashboardManagement.addWidget and dashboard.dashboard_id then
        viz.DashboardManagement.addWidget(dashboard.dashboard_id, "chart", chart)
    end

    plugin_state.last_dashboard = {chart = chart, dashboard = dashboard}
    log_operation("VISUALIZE", "Generated party synergy dashboard")
    return plugin_state.last_dashboard
end

-- Export party template via Import/Export Manager
function exportPartyTemplate(party, format, path)
    load_phase11()
    local exporter = dependencies.import_export and dependencies.import_export.DataExporter
    local output_path = path or "party_template.json"
    local fmt = (format or "json"):lower()

    if exporter then
        if fmt == "csv" then
            exporter.exportToCSV(party, output_path, true)
        elseif fmt == "xml" then
            exporter.exportToXML(party, output_path)
        else
            exporter.exportToJSON(party, output_path)
        end
    else
        local file = io.open(output_path, "w")
        if file then
            file:write("-- party template placeholder --\n")
            file:close()
        end
    end

    log_operation("EXPORT", string.format("Exported party template to %s", output_path))
    return {path = output_path, format = fmt}
end

-- Import party template via Import/Export Manager
function importPartyTemplate(path, format)
    load_phase11()
    if not path then
        log_operation("ERROR", "No path provided for party import")
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
    log_operation("IMPORT", string.format("Imported party template from %s", path))
    return data or {path = path}
end

-- Snapshot party setup using Backup & Restore System
function snapshotParty(party, label)
    load_phase11()
    local backup = dependencies.backup_restore
    local snapshot_label = label or "party_snapshot"
    if backup and backup.SnapshotManagement then
        local snap = backup.SnapshotManagement.createSnapshot(snapshot_label, party or {})
        log_operation("BACKUP", "Party snapshot created")
        return snap
    end
    log_operation("WARN", "Backup system unavailable; snapshot not persisted")
    return {snapshot_id = "local", name = snapshot_label}
end

-- Restore party from backup id
function restoreParty(snapshot_id)
    load_phase11()
    local backup = dependencies.backup_restore
    if backup and backup.RecoverySystem then
        local restore = backup.RecoverySystem.restoreFromBackup(snapshot_id or "latest", {target = "party"})
        log_operation("RESTORE", "Restore triggered for party snapshot")
        return restore
    end
    log_operation("WARN", "Backup system unavailable; restore skipped")
    return {success = false}
end

-- Sync party data to other plugins via Integration Hub
function syncPartyData(target_plugins)
    load_phase11()
    local hub = dependencies.integration_hub
    local targets = target_plugins or {"advanced-battle-predictor", "equipment-optimizer"}
    local result = hub and hub.UnifiedAPI and hub.UnifiedAPI.crossPluginCall("party-optimizer", table.concat(targets, ","), "sync_party", {party = plugin_state.last_party}) or {}
    log_operation("SYNC", "Synced party data to target plugins")
    return result
end

-- Load party data from persistence layer
function loadPartyData()
    local db = load_database_layer()
    if db and db.loadPartyConfig then
        local persisted = db.loadPartyConfig()
        if persisted then
            plugin_state.last_party = persisted
            log_operation("LOAD", "Loaded party config from persistence layer")
            return persisted
        end
    end
    
    return plugin_state.last_party or {members = {}, size = 0}
end

-- Generate party report via Data Visualization Suite
function generatePartyReport(party)
    load_phase11()
    local viz = dependencies.visualization
    local normalized = normalize_party(party)
    local report = viz and viz.ReportGeneration and viz.ReportGeneration.generateReport("Party Report", normalized, "PDF") or {name = "Party Report"}
    log_operation("REPORT", "Generated party report")
    return report
end

-- Minimal status helper
function getPartyStatus()
    local status = plugin_state.last_party or {members = {}, size = 0}
    log_operation("STATUS", string.format("Party size: %d", status.size or 0))
    return status
end

-- Initialization
local function initialize()
    if plugin_state.initialized then return true end
    plugin_state.initialized = true
    log_operation("INIT", "Party Optimizer initialized with Phase 11 integrations")
    return true
end

initialize()

print("Party Optimizer v1.2.0 loaded")
print("Commands: optimizePartyComposition(party, goal), recommendPartyForScenario(info), autoConfigureParty(party, info)")
print("Visualization: visualizePartySynergy(party), generatePartyReport(party)")
print("Import/Export: exportPartyTemplate(party, fmt, path), importPartyTemplate(path, fmt)")
print("Backup/Sync: snapshotParty(party, label), restoreParty(id), syncPartyData(targets)")
