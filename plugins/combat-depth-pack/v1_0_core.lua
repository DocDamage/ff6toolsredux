--[[
  Combat Depth Pack - Pack A
  Features: Dynamic Encounter Tuner, Boss Remix & Affixes, AI Companion Director
  Version: 1.0.0
  Notes: Uses lazy-loaded Phase 11 modules when available (analytics, visualization, automation)
]]

-- ============================================================================
-- DEPENDENCY LOADER (lazy)
-- ============================================================================

local analytics, viz, automation = nil, nil, nil

local function load_phase11()
  if not analytics then
    local success, module = pcall(require, "plugins.advanced-analytics-engine.v1_0_core")
    analytics = success and module or nil
  end
  if not viz then
    local success, module = pcall(require, "plugins.data-visualization-suite.v1_0_core")
    viz = success and module or nil
  end
  if not automation then
    local success, module = pcall(require, "plugins.automation-framework.v1_0_core")
    automation = success and module or nil
  end
  return {analytics = analytics, viz = viz, automation = automation}
end

-- ============================================================================
-- MODULE: DYNAMIC ENCOUNTER TUNER
-- ============================================================================

local EncounterTuner = {
  version = "1.0.0",
  presets = {
    story = {encounter_rate = 0.6, elite_chance = 0.05, flee_allowed = true},
    grind = {encounter_rate = 1.4, elite_chance = 0.15, xp_multiplier = 1.25},
    challenge = {encounter_rate = 1.1, elite_chance = 0.25, ai_aggression = "high"}
  }
}

---Configure per-zone encounter rates and behaviors
---@param zone_config table
---@return table config_applied
function EncounterTuner.configureEncounterRates(zone_config)
  if not zone_config then
    return {success = false, error = "No zone_config provided"}
  end

  local deps = load_phase11()
  local applied = {
    zone = zone_config.zone or "Unknown",
    encounter_rate = zone_config.encounter_rate or 1.0,
    elite_chance = zone_config.elite_chance or 0.1,
    xp_multiplier = zone_config.xp_multiplier or 1.0,
    ai_aggression = zone_config.ai_aggression or "normal",
    success = true
  }

  if deps.analytics and deps.analytics.PatternRecognition and type(deps.analytics.PatternRecognition.detectPatterns) == "function" then
    -- Example: use analytics to validate pacing
    local pacing = deps.analytics.PatternRecognition.detectPatterns({applied}, "encounter_pacing")
    applied.pacing_confidence = pacing.confidence or 75
  end

  return applied
end

---Generate a dynamic encounter based on party and zone state
---@param zone_state table
---@param party_state table
---@return table encounter
function EncounterTuner.generateDynamicEncounter(zone_state, party_state)
  zone_state = zone_state or {}
  party_state = party_state or {}

  local encounter = {
    zone = zone_state.zone or "Unknown",
    enemies = zone_state.enemy_pool or {"Soldier", "Officer", "Beast"},
    elite = (zone_state.elite_roll or 0.1) > 0.2,
    xp_reward = math.floor((zone_state.base_xp or 150) * (zone_state.xp_multiplier or 1.0)),
    gil_reward = math.floor((zone_state.base_gil or 180) * 1.1),
    danger_rating = "Moderate"
  }

  if (party_state.avg_level or 20) < 15 then
    encounter.danger_rating = "Low"
  elseif (party_state.avg_level or 20) > 35 then
    encounter.danger_rating = "High"
  end

  return encounter
end

---Apply a preset to zone configuration
---@param preset_name string
---@return table preset
function EncounterTuner.setPreset(preset_name)
  local preset = EncounterTuner.presets[preset_name]
  if not preset then
    return {success = false, error = "Unknown preset"}
  end
  preset.success = true
  preset.preset = preset_name
  return preset
end

-- ============================================================================
-- MODULE: BOSS REMIX & AFFIXES
-- ============================================================================

local BossRemix = {
  version = "1.0.0",
  affix_catalog = {
    enraged = {name = "Enraged", effect = "Higher offense, lower defense"},
    arcane_shield = {name = "Arcane Shield", effect = "Magic barrier every 3 turns"},
    glass_cannon = {name = "Glass Cannon", effect = "Massive damage, fragile"},
    regenerating = {name = "Regenerating", effect = "Heals 5% HP each turn"},
    stagger_resist = {name = "Stagger Resist", effect = "Reduced stun duration"}
  }
}

---Apply affixes to a boss profile
---@param boss_profile table
---@param affixes table
---@return table remixed
function BossRemix.applyAffixes(boss_profile, affixes)
  if not boss_profile then
    return {success = false, error = "No boss_profile provided"}
  end

  local remixed = {
    boss = boss_profile.name or "Unknown Boss",
    base_hp = boss_profile.hp or 50000,
    base_attack = boss_profile.attack or 1200,
    affixes = {},
    success = true
  }

  for _, affix_key in ipairs(affixes or {}) do
    local affix = BossRemix.affix_catalog[affix_key]
    if affix then
      table.insert(remixed.affixes, affix)
    end
  end

  remixed.hp = math.floor(remixed.base_hp * (1 + (#remixed.affixes * 0.08)))
  remixed.attack = math.floor(remixed.base_attack * (1 + (#remixed.affixes * 0.05)))

  return remixed
end

---Generate a remix plan with phases and affixes
---@param boss_profile table
---@param remix_config table
---@return table plan
function BossRemix.generateRemix(boss_profile, remix_config)
  remix_config = remix_config or {affixes = {"enraged"}}
  local remixed = BossRemix.applyAffixes(boss_profile, remix_config.affixes)

  remixed.phases = remix_config.phases or {
    {phase = 1, triggers = {hp_percent = 75}, behavior = "Balanced"},
    {phase = 2, triggers = {hp_percent = 50}, behavior = "Aggressive"},
    {phase = 3, triggers = {hp_percent = 25}, behavior = "Desperate"}
  }

  remixed.loot_bonus = remix_config.loot_bonus or 1.15
  remixed.experiment_tag = remix_config.tag or "BossRemixV1"

  return remixed
end

---List available affix presets
---@return table presets
function BossRemix.listAffixPresets()
  return {
    starter = {"enraged"},
    defensive = {"arcane_shield", "stagger_resist"},
    glass = {"glass_cannon", "enraged"},
    sustain = {"regenerating", "arcane_shield"}
  }
end

-- ============================================================================
-- MODULE: AI COMPANION DIRECTOR
-- ============================================================================

local CompanionDirector = {
  version = "1.0.0",
  profiles = {}
}

---Define an AI tactics profile
---@param profile_name string
---@param config table
---@return table profile
function CompanionDirector.defineProfile(profile_name, config)
  if not profile_name or not config then
    return {success = false, error = "Missing profile name or config"}
  end

  CompanionDirector.profiles[profile_name] = {
    priorities = config.priorities or {"survival", "damage", "support"},
    triggers = config.triggers or {},
    risk_tolerance = config.risk_tolerance or "normal",
    success = true
  }

  return CompanionDirector.profiles[profile_name]
end

---Evaluate battle state and return a snapshot
---@param battle_state table
---@param profile_name string
---@return table snapshot
function CompanionDirector.evaluateBattleState(battle_state, profile_name)
  battle_state = battle_state or {}
  local profile = CompanionDirector.profiles[profile_name] or {risk_tolerance = "normal"}

  local snapshot = {
    hp_status = battle_state.hp_status or "stable",
    threat_level = battle_state.threat_level or "medium",
    action_bias = profile.priorities or {"damage"},
    risk = profile.risk_tolerance,
    recommendations = {}
  }

  table.insert(snapshot.recommendations, "Prioritize heals below 40% HP")
  table.insert(snapshot.recommendations, "Focus highest threat enemy")

  return snapshot
end

---Recommend a companion action based on profile and state
---@param battle_state table
---@param profile_name string
---@return table action
function CompanionDirector.recommendAction(battle_state, profile_name)
  local snapshot = CompanionDirector.evaluateBattleState(battle_state, profile_name)

  local action = {
    type = "attack",
    target = battle_state.primary_target or "highest_threat",
    rationale = "Default action"
  }

  if snapshot.hp_status == "critical" then
    action.type = "heal"
    action.target = "lowest_hp_ally"
    action.rationale = "Ally HP critical"
  elseif snapshot.threat_level == "high" then
    action.type = "debuff"
    action.target = "highest_threat"
    action.rationale = "Reduce enemy threat"
  elseif snapshot.risk == "aggressive" then
    action.type = "burst"
    action.target = "boss"
    action.rationale = "Aggressive profile"
  end

  return action
end

-- ============================================================================
-- EXPORTS
-- ============================================================================

return {
  version = "1.0.0",
  EncounterTuner = EncounterTuner,
  BossRemix = BossRemix,
  CompanionDirector = CompanionDirector,
  features = {
    dynamicEncounterTuning = true,
    bossRemixAffixes = true,
    companionAIDirector = true
  }
}
