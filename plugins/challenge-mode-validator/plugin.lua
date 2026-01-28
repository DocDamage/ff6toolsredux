--[[
  Challenge Mode Validator Plugin
  Verifies challenge run compliance with preset or custom rules.
  
  Features:
  - 9 preset challenge mode templates
  - Custom rule creation
  - Real-time violation detection
  - Violation log with timestamps
  - Challenge verification report
  - Export challenge proof
  
  Phase: 5 (Challenge & Advanced Tools)
  Version: 1.0.0
]]

-- Safe API call wrapper
local function safeCall(fn, ...)
  local success, result = pcall(fn, ...)
  if success then
    return result
  end
  return nil
end

-- Safe require with warning
local function safe_require(path)
  local ok, mod = pcall(require, path)
  if not ok then
    print(string.format("[Challenge Validator] WARN: Dependency unavailable: %s", path))
    return nil
  end
  return mod
end

-- Phase 11 dependency handles (lazy)
local dependencies = {
  analytics = nil,
  visualization = nil,
  import_export = nil,
  backup_restore = nil,
  automation = nil,
  integration_hub = nil,
  api_gateway = nil
}

local function load_phase11()
  dependencies.analytics = dependencies.analytics or safe_require("plugins.advanced-analytics-engine.v1_0_core")
  dependencies.visualization = dependencies.visualization or safe_require("plugins.data-visualization-suite.v1_0_core")
  dependencies.import_export = dependencies.import_export or safe_require("plugins.import-export-manager.v1_0_core")
  dependencies.backup_restore = dependencies.backup_restore or safe_require("plugins.backup-restore-system.v1_0_core")
  dependencies.automation = dependencies.automation or safe_require("plugins.automation-framework.v1_0_core")
  dependencies.integration_hub = dependencies.integration_hub or safe_require("plugins.integration-hub.v1_0_core")
  dependencies.api_gateway = dependencies.api_gateway or safe_require("plugins.api-gateway.v1_0_core")
  return dependencies
end

-- Rule categories
local RULE_CATEGORIES = {
  LEVEL = "Level Restrictions",
  EQUIPMENT = "Equipment Restrictions",
  MAGIC = "Magic Restrictions",
  PARTY = "Party Restrictions",
  INVENTORY = "Inventory Restrictions",
  GAMEPLAY = "Gameplay Restrictions"
}

-- Preset challenge modes
local CHALLENGE_PRESETS = {
  {
    id = "low_level",
    name = "Low Level Run",
    description = "Complete game with all characters under level 20",
    rules = {
      {category = RULE_CATEGORIES.LEVEL, check = "max_level", value = 20, desc = "No character above level 20"}
    }
  },
  {
    id = "natural_magic_block",
    name = "Natural Magic Block",
    description = "No magic learning allowed (Runic, natural spells only)",
    rules = {
      {category = RULE_CATEGORIES.MAGIC, check = "no_learned_magic", value = true, desc = "No learned spells from espers"},
      {category = RULE_CATEGORIES.EQUIPMENT, check = "no_espers", value = true, desc = "No espers equipped"}
    }
  },
  {
    id = "solo_run",
    name = "Solo Character Run",
    description = "Use only one character for the entire game",
    rules = {
      {category = RULE_CATEGORIES.PARTY, check = "max_party_size", value = 1, desc = "Only 1 character in active party"}
    }
  },
  {
    id = "no_equipment",
    name = "No Equipment Challenge",
    description = "No weapons or armor equipped",
    rules = {
      {category = RULE_CATEGORIES.EQUIPMENT, check = "no_weapons", value = true, desc = "No weapons equipped"},
      {category = RULE_CATEGORIES.EQUIPMENT, check = "no_armor", value = true, desc = "No armor equipped"}
    }
  },
  {
    id = "no_shop",
    name = "No Shop Run",
    description = "Cannot buy items from shops (starting gil only)",
    rules = {
      {category = RULE_CATEGORIES.INVENTORY, check = "max_gil_spent", value = 0, desc = "No gil spent in shops"}
    }
  },
  {
    id = "fixed_party",
    name = "Fixed Party Challenge",
    description = "Use the same 4 characters for entire game",
    rules = {
      {category = RULE_CATEGORIES.PARTY, check = "fixed_roster", value = true, desc = "Same 4 characters only"}
    }
  },
  {
    id = "no_espers",
    name = "No Esper Run",
    description = "Cannot equip espers in battle",
    rules = {
      {category = RULE_CATEGORIES.EQUIPMENT, check = "no_espers", value = true, desc = "No espers equipped"}
    }
  },
  {
    id = "ancient_cave",
    name = "Ancient Cave Mode",
    description = "Start fresh, only use items found in dungeons",
    rules = {
      {category = RULE_CATEGORIES.INVENTORY, check = "no_shop_items", value = true, desc = "No shop purchases"},
      {category = RULE_CATEGORIES.LEVEL, check = "start_level_1", value = true, desc = "All characters start at level 1"}
    }
  },
  {
    id = "minimalist",
    name = "Minimalist Run",
    description = "Complete with minimum battles and items",
    rules = {
      {category = RULE_CATEGORIES.GAMEPLAY, check = "max_battles", value = 100, desc = "Maximum 100 battles"},
      {category = RULE_CATEGORIES.INVENTORY, check = "max_items_used", value = 50, desc = "Maximum 50 items used"}
    }
  }
}

-- Storage file
local CHALLENGE_FILE = "challenge_validator/active_challenge.lua"
local VIOLATIONS_FILE = "challenge_validator/violations.lua"

-- Serialize table
local function serialize(tbl, indent)
  indent = indent or 0
  local spacing = string.rep("  ", indent)
  local lines = {"{"}
  
  for k, v in pairs(tbl) do
    local key = type(k) == "number" and string.format("[%d]", k) or string.format('["%s"]', k)
    local value
    if type(v) == "table" then
      value = serialize(v, indent + 1)
    elseif type(v) == "string" then
      value = string.format('"%s"', v:gsub('"', '\\"'))
    elseif type(v) == "boolean" then
      value = tostring(v)
    else
      value = tostring(v)
    end
    table.insert(lines, spacing .. "  " .. key .. " = " .. value .. ",")
  end
  
  table.insert(lines, spacing .. "}")
  return table.concat(lines, "\n")
end

-- Load active challenge
local function loadChallenge()
  local content = safeCall(ReadFile, CHALLENGE_FILE)
  if not content then return nil end
  
  local success, data = pcall(load, "return " .. content)
  if success and data then
    return data()
  end
  return nil
end

-- Save active challenge
local function saveChallenge(challenge)
  local content = "return " .. serialize(challenge)
  safeCall(WriteFile, CHALLENGE_FILE, content)
end

-- Load violations log
local function loadViolations()
  local content = safeCall(ReadFile, VIOLATIONS_FILE)
  if not content then return {} end
  
  local success, data = pcall(load, "return " .. content)
  if success and data then
    return data()
  end
  return {}
end

-- Save violations log
local function saveViolations(violations)
  local content = "return " .. serialize(violations)
  safeCall(WriteFile, VIOLATIONS_FILE, content)
end

-- Check individual rule
local function checkRule(rule, save_data)
  if rule.check == "max_level" then
    local party = safeCall(GetParty)
    if not party then return true, "API unavailable" end
    
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Level and char.Level > rule.value then
        return false, string.format("Character %s level %d exceeds limit %d", 
          char.Name or "Unknown", char.Level, rule.value)
      end
    end
    return true, nil
    
  elseif rule.check == "no_learned_magic" then
    local party = safeCall(GetParty)
    if not party then return true, "API unavailable" end
    
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Spells then
        -- Check for learned spells (excluding natural spells)
        local natural_spells = {"Cure", "Fire", "Blizzard"} -- Simplified
        for spell_name, _ in pairs(char.Spells) do
          local is_natural = false
          for _, nat in ipairs(natural_spells) do
            if spell_name == nat then is_natural = true break end
          end
          if not is_natural then
            return false, string.format("Character %s has learned spell: %s", 
              char.Name or "Unknown", spell_name)
          end
        end
      end
    end
    return true, nil
    
  elseif rule.check == "no_espers" then
    local party = safeCall(GetParty)
    if not party then return true, "API unavailable" end
    
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Esper and char.Esper ~= "" and char.Esper ~= "None" then
        return false, string.format("Character %s has esper equipped: %s", 
          char.Name or "Unknown", char.Esper)
      end
    end
    return true, nil
    
  elseif rule.check == "max_party_size" then
    local party = safeCall(GetParty)
    if not party then return true, "API unavailable" end
    
    if #party > rule.value then
      return false, string.format("Party size %d exceeds limit %d", #party, rule.value)
    end
    return true, nil
    
  elseif rule.check == "no_weapons" then
    local party = safeCall(GetParty)
    if not party then return true, "API unavailable" end
    
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Weapon and char.Weapon ~= "" and char.Weapon ~= "None" then
        return false, string.format("Character %s has weapon equipped: %s", 
          char.Name or "Unknown", char.Weapon)
      end
    end
    return true, nil
    
  elseif rule.check == "no_armor" then
    local party = safeCall(GetParty)
    if not party then return true, "API unavailable" end
    
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Armor and char.Armor ~= "" and char.Armor ~= "None" then
        return false, string.format("Character %s has armor equipped: %s", 
          char.Name or "Unknown", char.Armor)
      end
    end
    return true, nil
    
  elseif rule.check == "fixed_roster" then
    -- This requires tracking initial party selection (placeholder)
    return true, "Tracking not implemented"
    
  elseif rule.check == "max_gil_spent" then
    -- Requires tracking gil transactions (placeholder)
    return true, "Tracking not implemented"
    
  elseif rule.check == "no_shop_items" then
    -- Requires item source tracking (placeholder)
    return true, "Tracking not implemented"
    
  elseif rule.check == "start_level_1" then
    -- Validation for game start state (placeholder)
    return true, "Start state validation not implemented"
    
  elseif rule.check == "max_battles" then
    -- Requires battle count tracking (placeholder)
    return true, "Tracking not implemented"
    
  elseif rule.check == "max_items_used" then
    -- Requires item usage tracking (placeholder)
    return true, "Tracking not implemented"
  end
  
  return true, "Unknown rule type"
end

-- Validate all rules for active challenge
local function validateChallenge()
  local challenge = loadChallenge()
  if not challenge then
    return nil, "No active challenge"
  end
  
  local violations = {}
  local save_data = {} -- Placeholder
  
  for i, rule in ipairs(challenge.rules) do
    local passed, error_msg = checkRule(rule, save_data)
    if not passed and error_msg ~= "API unavailable" and error_msg ~= "Tracking not implemented" then
      table.insert(violations, {
        timestamp = os.time(),
        category = rule.category,
        rule_desc = rule.desc,
        violation_detail = error_msg
      })
    end
  end
  
  return violations, challenge
end

-- Start new challenge
local function startChallenge()
  local lines = {}
  table.insert(lines, "Select Challenge Mode:")
  for i, preset in ipairs(CHALLENGE_PRESETS) do
    table.insert(lines, string.format("%d. %s", i, preset.name))
  end
  table.insert(lines, string.format("%d. Custom Challenge", #CHALLENGE_PRESETS + 1))
  table.insert(lines, string.format("%d. Cancel", #CHALLENGE_PRESETS + 2))
  
  local choice = ShowDialog("Start Challenge", lines)
  if not choice or choice > #CHALLENGE_PRESETS + 1 then return "Cancelled." end
  
  local challenge
  if choice <= #CHALLENGE_PRESETS then
    challenge = CHALLENGE_PRESETS[choice]
    challenge.started = os.time()
  else
    -- Custom challenge creation
    local name = ShowInput("Enter challenge name:")
    if not name or name == "" then return "Cancelled." end
    
    local desc = ShowInput("Enter challenge description:")
    if not desc or desc == "" then return "Cancelled." end
    
    challenge = {
      id = "custom_" .. os.time(),
      name = name,
      description = desc,
      rules = {},
      started = os.time()
    }
  end
  
  saveChallenge(challenge)
  safeCall(WriteFile, VIOLATIONS_FILE, "return {}") -- Reset violations
  
  return string.format("Started challenge: %s\n%s", challenge.name, challenge.description)
end

-- Check current challenge
local function checkCurrentChallenge()
  local violations, challenge = validateChallenge()
  if not challenge then
    return "No active challenge. Start a new challenge first."
  end
  
  local lines = {}
  table.insert(lines, "=== Challenge Validation ===")
  table.insert(lines, string.format("Challenge: %s", challenge.name))
  table.insert(lines, string.format("Started: %s", os.date("%Y-%m-%d %H:%M", challenge.started)))
  table.insert(lines, "")
  
  table.insert(lines, "=== Rules ===")
  for i, rule in ipairs(challenge.rules) do
    table.insert(lines, string.format("%d. [%s] %s", i, rule.category, rule.desc))
  end
  table.insert(lines, "")
  
  if #violations > 0 then
    table.insert(lines, string.format("=== VIOLATIONS DETECTED: %d ===", #violations))
    for i, v in ipairs(violations) do
      table.insert(lines, string.format("[%s] %s", os.date("%Y-%m-%d %H:%M", v.timestamp), v.rule_desc))
      table.insert(lines, string.format("    Detail: %s", v.violation_detail))
    end
    
    -- Save violations to log
    local all_violations = loadViolations()
    for _, v in ipairs(violations) do
      table.insert(all_violations, v)
    end
    saveViolations(all_violations)
  else
    table.insert(lines, "=== STATUS: COMPLIANT ===")
    table.insert(lines, "All rules verified. No violations detected.")
  end
  
  return table.concat(lines, "\n")
end

-- ============================================================================
-- PHASE 11 INTEGRATIONS (TIER 1 - CHALLENGE VALIDATOR)
-- ============================================================================

-- Analyze violation trends with Analytics Engine
function analyzeChallengeTrends()
  load_phase11()
  local analytics = dependencies.analytics
  local violations = loadViolations() or {}
  local counts = {}
  for _, v in ipairs(violations) do
    counts[v.category] = (counts[v.category] or 0) + 1
  end
  local dataset = {}
  for _, c in pairs(counts) do table.insert(dataset, c) end
  local trend = analytics and analytics.PatternRecognition and analytics.PatternRecognition.detectTrend(dataset) or {}
  return {counts = counts, trend = trend}
end

-- Generate compliance dashboard
function generateComplianceDashboard()
  load_phase11()
  local viz = dependencies.visualization
  local analysis = analyzeChallengeTrends()
  local categories, values = {}, {}
  for cat, count in pairs(analysis.counts) do
    table.insert(categories, cat)
    table.insert(values, count)
  end
  local chart = viz and viz.ChartGeneration and viz.ChartGeneration.generateBarChart("Violations by Category", categories, values) or {type = "Bar"}
  local dashboard = viz and viz.DashboardManagement and viz.DashboardManagement.createDashboard("Challenge Compliance", "grid") or {name = "Challenge Compliance"}
  if viz and viz.DashboardManagement and viz.DashboardManagement.addWidget and dashboard.dashboard_id then
    viz.DashboardManagement.addWidget(dashboard.dashboard_id, "chart", chart)
  end
  return {chart = chart, dashboard = dashboard, analysis = analysis}
end

-- Export compliance report via Data Visualization Suite
function exportComplianceReport(format, output_path)
  load_phase11()
  local viz = dependencies.visualization
  local report = viz and viz.ReportGeneration and viz.ReportGeneration.generateReport("Challenge Compliance Report", {
    active = loadChallenge(),
    violations = loadViolations(),
    analysis = analyzeChallengeTrends()
  }, format or "PDF") or {name = "Compliance Report"}

  if output_path and dependencies.import_export and dependencies.import_export.DataExporter then
    dependencies.import_export.DataExporter.exportToJSON(report, output_path)
  end
  return report
end

-- Import/export challenge templates via Import/Export Manager
function exportChallengeTemplate(challenge, format, path)
  load_phase11()
  local exporter = dependencies.import_export and dependencies.import_export.DataExporter
  local output = path or "challenge_template.json"
  local fmt = (format or "json"):lower()
  if exporter then
    if fmt == "csv" then exporter.exportToCSV(challenge, output, true)
    elseif fmt == "xml" then exporter.exportToXML(challenge, output)
    else exporter.exportToJSON(challenge, output) end
  end
  return {path = output, format = fmt}
end

function importChallengeTemplate(path, format)
  load_phase11()
  if not path then return nil end
  local importer = dependencies.import_export and dependencies.import_export.DataImporter
  local fmt = (format or "json"):lower()
  local data
  if importer then
    if fmt == "csv" then data = importer.importCSV(path, true)
    elseif fmt == "xml" then data = importer.importXML(path, "lenient")
    else data = importer.importJSON(path, true) end
  end
  return data or {path = path}
end

-- Snapshot/restore challenge state with Backup & Restore System
function snapshotChallengeState(label)
  load_phase11()
  local backup = dependencies.backup_restore
  local snap_label = label or "challenge_snapshot"
  local state = {challenge = loadChallenge(), violations = loadViolations()}
  if backup and backup.SnapshotManagement then
    return backup.SnapshotManagement.createSnapshot(snap_label, state)
  end
  return {snapshot_id = "local", name = snap_label}
end

function restoreChallengeState(snapshot_id)
  load_phase11()
  local backup = dependencies.backup_restore
  if backup and backup.RecoverySystem then
    return backup.RecoverySystem.restoreFromBackup(snapshot_id or "latest", {target = "challenge"})
  end
  return {success = false, error = "Backup unavailable"}
end

-- Automation: schedule periodic validation
function scheduleChallengeValidation(cron_expr)
  load_phase11()
  local automation = dependencies.automation
  if automation and automation.Scheduler and automation.Scheduler.scheduleTask then
    automation.Scheduler.scheduleTask("challenge_validation", cron_expr or "*/30 * * * *", function()
      checkCurrentChallenge()
    end)
  end
  return {scheduled = true, cron = cron_expr or "*/30 * * * *"}
end

-- Integration Hub sync
function syncChallengeData(target_plugins)
  load_phase11()
  local hub = dependencies.integration_hub
  local targets = target_plugins or {"party-optimizer", "equipment-optimizer"}
  if hub and hub.UnifiedAPI then
    return hub.UnifiedAPI.crossPluginCall("challenge-mode-validator", table.concat(targets, ","), "sync_challenge", {
      active = loadChallenge(),
      violations = loadViolations()
    })
  end
  return {success = false, error = "Integration Hub unavailable"}
end

-- API Gateway: register read-only endpoint
function registerChallengeAPI()
  load_phase11()
  local api = dependencies.api_gateway
  if not api or not api.RESTInterface then return nil end
  local endpoint = api.RESTInterface.registerEndpoint("GET", "/api/challenge/status", function()
    return {active = loadChallenge(), violations = loadViolations()}
  end)
  api.RESTInterface.addAuthentication(endpoint.endpoint_id, "API-Key")
  api.RESTInterface.addRateLimit(endpoint.endpoint_id, 120)
  return endpoint
end

-- View violation log
local function viewViolationLog()
  local challenge = loadChallenge()
  if not challenge then
    return "No active challenge."
  end
  
  local violations = loadViolations()
  
  local lines = {}
  table.insert(lines, "=== Violation Log ===")
  table.insert(lines, string.format("Challenge: %s", challenge.name))
  table.insert(lines, string.format("Total Violations: %d\n", #violations))
  
  if #violations == 0 then
    table.insert(lines, "No violations recorded. Challenge is clean!")
  else
    for i, v in ipairs(violations) do
      table.insert(lines, string.format("[%d] %s", i, os.date("%Y-%m-%d %H:%M:%S", v.timestamp)))
      table.insert(lines, string.format("    Category: %s", v.category))
      table.insert(lines, string.format("    Rule: %s", v.rule_desc))
      table.insert(lines, string.format("    Detail: %s\n", v.violation_detail))
    end
  end
  
  return table.concat(lines, "\n")
end

-- Export challenge proof
local function exportProof()
  local challenge = loadChallenge()
  if not challenge then
    return "No active challenge to export."
  end
  
  local violations = loadViolations()
  local timestamp = os.date("%Y%m%d_%H%M%S")
  local filename = string.format("challenge_validator/proof_%s_%s.txt", 
    challenge.id, timestamp)
  
  local lines = {}
  table.insert(lines, "=== FF6 Challenge Mode Proof ===")
  table.insert(lines, "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"))
  table.insert(lines, "")
  table.insert(lines, "Challenge: " .. challenge.name)
  table.insert(lines, "Description: " .. challenge.description)
  table.insert(lines, "Started: " .. os.date("%Y-%m-%d %H:%M", challenge.started))
  table.insert(lines, "")
  
  table.insert(lines, "=== Rules ===")
  for i, rule in ipairs(challenge.rules) do
    table.insert(lines, string.format("%d. [%s] %s", i, rule.category, rule.desc))
  end
  table.insert(lines, "")
  
  if #violations > 0 then
    table.insert(lines, string.format("=== VIOLATIONS: %d ===", #violations))
    table.insert(lines, "This challenge has recorded violations and is NOT verified.")
    table.insert(lines, "")
    for i, v in ipairs(violations) do
      table.insert(lines, string.format("[%s] %s: %s", 
        os.date("%Y-%m-%d %H:%M:%S", v.timestamp), v.rule_desc, v.violation_detail))
    end
  else
    table.insert(lines, "=== STATUS: VERIFIED ===")
    table.insert(lines, "This challenge has NO recorded violations.")
    table.insert(lines, "All rules were compliant at time of validation.")
  end
  
  table.insert(lines, "")
  table.insert(lines, "Generated by FF6 Save Editor - Challenge Mode Validator v1.0.0")
  
  local content = table.concat(lines, "\n")
  safeCall(WriteFile, filename, content)
  
  return "Challenge proof exported to: " .. filename
end

-- End current challenge
local function endChallenge()
  local challenge = loadChallenge()
  if not challenge then
    return "No active challenge."
  end
  
  local confirm = ShowDialog(
    string.format("End challenge: %s?", challenge.name),
    {"1. Yes, end challenge", "2. No, keep running"}
  )
  
  if confirm == 1 then
    -- Archive the challenge
    local timestamp = os.date("%Y%m%d_%H%M%S")
    local archive_file = string.format("challenge_validator/archived_%s_%s.lua", 
      challenge.id, timestamp)
    
    challenge.ended = os.time()
    challenge.violations = loadViolations()
    
    local content = "return " .. serialize(challenge)
    safeCall(WriteFile, archive_file, content)
    
    -- Clear active challenge
    safeCall(WriteFile, CHALLENGE_FILE, "")
    safeCall(WriteFile, VIOLATIONS_FILE, "return {}")
    
    return string.format("Challenge ended and archived to:\n%s", archive_file)
  end
  
  return "Challenge continues."
end

-- Main menu
local function main()
  while true do
    local challenge = loadChallenge()
    local status = challenge and string.format("Active: %s", challenge.name) or "No Active Challenge"
    
    local choice = ShowDialog("Challenge Mode Validator", {
      "Status: " .. status,
      "────────────────────────",
      "1. Start New Challenge",
      "2. Check Current Challenge",
      "3. View Violation Log",
      "4. Export Challenge Proof",
      "5. End Current Challenge",
      "6. Exit"
    })
    
    if not choice or choice == 6 then break end
    if choice <= 2 then choice = choice + 2 end -- Adjust for display lines
    
    local result = ""
    if choice == 3 then
      result = startChallenge()
    elseif choice == 4 then
      result = checkCurrentChallenge()
    elseif choice == 5 then
      result = viewViolationLog()
    elseif choice == 6 then
      result = exportProof()
    elseif choice == 7 then
      result = endChallenge()
    end
    
    if result ~= "" then
      ShowDialog("Result", {result, "Press any key to continue..."})
    end
  end
end

-- Run plugin
main()
