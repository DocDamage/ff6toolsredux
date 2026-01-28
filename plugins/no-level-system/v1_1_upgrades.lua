--[[
  No-Level System Plugin - v1.1+ Upgrade Extension
  Level cap enforcement, stat scaling, progression alternatives, challenge modes
  
  Phase: 7F (Difficulty Systems)
  Version: 1.1+ (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: LEVEL CAP ENFORCEMENT (~240 LOC)
-- ============================================================================

local LevelCapEnforcement = {}

---Set level cap for run
---@param maxLevel number Maximum allowed level
---@return table capSettings Level cap configuration
function LevelCapEnforcement.setLevelCap(maxLevel)
  if not maxLevel or maxLevel < 1 then maxLevel = 99 end
  
  local capSettings = {
    enabled = maxLevel < 99,
    maxLevel = maxLevel,
    enforcement = "strict",
    rewards = "scaled",
    penalties = {}
  }
  
  return capSettings
end

---Validate character level against cap
---@param character table Character to check
---@param capLevel number Maximum allowed level
---@return table validation Level validation result
function LevelCapEnforcement.validateLevelCap(character, capLevel)
  if not character or not capLevel then return {} end
  
  local validation = {
    character = character.name,
    currentLevel = character.level or 1,
    capLevel = capLevel,
    violatescap = false,
    excessLevel = 0,
    penalty = 0
  }
  
  if (character.level or 1) > capLevel then
    validation.violatesCapLevel = true
    validation.excessLevel = (character.level or 1) - capLevel
    validation.penalty = validation.excessLevel * 10  -- 10% penalty per level
  end
  
  return validation
end

---Prevent level up at cap
---@param character table Character near cap
---@param capLevel number Level cap
---@return table preventionStatus Prevention result
function LevelCapEnforcement.preventLevelUpAtCap(character, capLevel)
  if not character or not capLevel then return {} end
  
  local preventionStatus = {
    character = character.name,
    currentLevel = character.level or 1,
    capLevel = capLevel,
    canLevelUp = (character.level or 1) < capLevel,
    message = ""
  }
  
  if (character.level or 1) >= capLevel then
    preventionStatus.message = "Cannot level up: Level cap reached (" .. capLevel .. ")"
  else
    local remaining = capLevel - (character.level or 1)
    preventionStatus.message = "Levels remaining: " .. remaining
  end
  
  return preventionStatus
end

---Calculate scaled experience at cap
---@param baseExp number Base experience from battle
---@param capLevel number Current level cap
---@return number scaledExp Adjusted experience value
function LevelCapEnforcement.calculateScaledExperience(baseExp, capLevel)
  if not baseExp or not capLevel then return 0 end
  
  -- At cap, experience is scaled down significantly
  local scaleFactor = 1.0
  
  if capLevel <= 10 then
    scaleFactor = 0.25
  elseif capLevel <= 30 then
    scaleFactor = 0.50
  elseif capLevel <= 50 then
    scaleFactor = 0.75
  else
    scaleFactor = 0.90
  end
  
  return math.floor(baseExp * scaleFactor)
end

---Generate level cap summary
---@param party table Current party
---@param capLevel number Level cap
---@return table summary Cap status for entire party
function LevelCapEnforcement.generateLevelCapSummary(party, capLevel)
  if not party or not capLevel then return {} end
  
  local summary = {
    capLevel = capLevel,
    partyStatus = {},
    charactersAtCap = 0,
    averageLevel = 0,
    totalLevelDeficit = 0
  }
  
  local totalLevel = 0
  
  for _, character in ipairs(party) do
    local level = character.level or 1
    totalLevel = totalLevel + level
    
    if level >= capLevel then
      summary.charactersAtCap = summary.charactersAtCap + 1
    else
      summary.totalLevelDeficit = summary.totalLevelDeficit + (capLevel - level)
    end
    
    table.insert(summary.partyStatus, {
      character = character.name,
      level = level,
      atCap = level >= capLevel
    })
  end
  
  summary.averageLevel = math.floor(totalLevel / #party)
  
  return summary
end

-- ============================================================================
-- FEATURE 2: STAT SCALING SYSTEM (~240 LOC)
-- ============================================================================

local StatScaling = {}

---Calculate scaled stat based on level and cap
---@param baseStat number Base stat value
---@param currentLevel number Character's current level
---@param maxLevel number Character's maximum level
---@return number scaledStat Adjusted stat value
function StatScaling.calculateScaledStat(baseStat, currentLevel, maxLevel)
  if not baseStat or not currentLevel or not maxLevel then return 0 end
  
  if maxLevel == 0 then return baseStat end
  
  local scaleFactor = currentLevel / maxLevel
  local scaledStat = math.floor(baseStat * scaleFactor)
  
  -- Apply minimum value
  return math.max(1, scaledStat)
end

---Apply stat cap scaling to character
---@param character table Character to scale
---@param capLevel number Level cap
---@return table scaledStats Adjusted stat block
function StatScaling.applyStatCapScaling(character, capLevel)
  if not character or not capLevel then return {} end
  
  local scaledStats = {
    character = character.name,
    originalStats = {},
    scaledStats = {},
    scalingFactor = 0
  }
  
  local baseStats = character.stats or {}
  scaledStats.originalStats = baseStats
  
  -- Calculate scaling factor
  local currentLevel = character.level or 1
  scaledStats.scalingFactor = currentLevel / math.max(1, capLevel)
  
  for stat, value in pairs(baseStats) do
    local scaled = StatScaling.calculateScaledStat(value, currentLevel, capLevel)
    scaledStats.scaledStats[stat] = scaled
  end
  
  return scaledStats
end

---Predict stat growth with level cap
---@param character table Character data
---@param capLevel number Level cap
---@param targetLevel number Target level within cap
---@return table growthProjection Stat growth to target
function StatScaling.predictGrowthWithCap(character, capLevel, targetLevel)
  if not character or not capLevel then return {} end
  
  targetLevel = math.min(targetLevel or capLevel, capLevel)
  
  local growthProjection = {
    character = character.name,
    currentLevel = character.level or 1,
    capLevel = capLevel,
    targetLevel = targetLevel,
    levelUpCount = targetLevel - (character.level or 1),
    projectedStats = {}
  }
  
  local baseStats = character.stats or {}
  
  for stat, value in pairs(baseStats) do
    local projectedValue = StatScaling.calculateScaledStat(value, targetLevel, capLevel)
    
    table.insert(growthProjection.projectedStats, {
      stat = stat,
      current = value,
      projected = projectedValue,
      gain = projectedValue - value
    })
  end
  
  return growthProjection
end

---Create stat scaling profile for challenge
---@param capLevel number Level cap
---@param difficulty string Difficulty name
---@return table profile Scaling profile
function StatScaling.createScalingProfile(capLevel, difficulty)
  if not capLevel then capLevel = 50 end
  
  local profile = {
    capLevel = capLevel,
    difficulty = difficulty,
    scalingMultipliers = {},
    recommendations = {}
  }
  
  if capLevel <= 20 then
    profile.scalingMultipliers = {hp = 0.6, attack = 0.5, defense = 0.7}
    profile.recommendations = {"Focus on strategy over raw stats", "Use equipment bonuses"}
  elseif capLevel <= 50 then
    profile.scalingMultipliers = {hp = 0.8, attack = 0.75, defense = 0.85}
    profile.recommendations = {"Balanced stat development", "Smart ability usage"}
  else
    profile.scalingMultipliers = {hp = 0.95, attack = 0.90, defense = 0.95}
    profile.recommendations = {"Near-normal gameplay", "Use advanced tactics"}
  end
  
  return profile
end

-- ============================================================================
-- FEATURE 3: PROGRESSION ALTERNATIVES (~220 LOC)
-- ============================================================================

local ProgressionAlternatives = {}

---Enable alternative progression system
---@param type string Progression type: "mastery", "challenge", "skill"
---@return table progressionSystem Enabled system
function ProgressionAlternatives.enableAlternativeProgression(type)
  if not type then type = "challenge" end
  
  local progressionSystem = {
    type = type,
    enabled = true,
    metrics = {},
    milestones = {}
  }
  
  if type == "mastery" then
    progressionSystem.metrics = {"Abilities learned", "Spells mastered", "Techniques perfected"}
    progressionSystem.milestones = {
      {name = "Novice Mastery", threshold = 10},
      {name = "Adept Mastery", threshold = 25},
      {name = "Expert Mastery", threshold = 50}
    }
  elseif type == "challenge" then
    progressionSystem.metrics = {"Bosses defeated", "Achievements earned", "Challenges completed"}
    progressionSystem.milestones = {
      {name = "Rising Hero", threshold = 3},
      {name = "Legendary Hero", threshold = 10},
      {name = "Ultimate Champion", threshold = 25}
    }
  elseif type == "skill" then
    progressionSystem.metrics = {"Battle efficiency", "Resource conservation", "Strategic complexity"}
    progressionSystem.milestones = {
      {name = "Apprentice Warrior", threshold = 5},
      {name = "Master Warrior", threshold = 15},
      {name = "Grand Master", threshold = 30}
    }
  end
  
  return progressionSystem
end

---Track alternative progression metric
---@param saveData table Save file data
---@param metricType string Type of metric
---@return number progress Current progress value
function ProgressionAlternatives.trackProgressionMetric(saveData, metricType)
  if not saveData or not metricType then return 0 end
  
  local progress = 0
  
  if metricType == "abilitiesLearned" then
    if saveData.characters then
      for _, char in ipairs(saveData.characters) do
        progress = progress + (char.abilities and #char.abilities or 0)
      end
    end
  elseif metricType == "bossesFought" then
    progress = saveData.bossesFought or 0
  elseif metricType == "achievementsEarned" then
    progress = saveData.achievements and #saveData.achievements or 0
  end
  
  return progress
end

---Calculate alternative progression level
---@param metrics table Progress metrics
---@return number level Current progression level
function ProgressionAlternatives.calculateProgressionLevel(metrics)
  if not metrics then return 0 end
  
  local totalProgress = 0
  local metricCount = 0
  
  for _, value in pairs(metrics) do
    if type(value) == "number" then
      totalProgress = totalProgress + value
      metricCount = metricCount + 1
    end
  end
  
  if metricCount == 0 then return 0 end
  
  return math.floor(totalProgress / metricCount)
end

---Recommend progression milestones
---@param currentProgress number Current progress value
---@param progressionType string Type of progression
---@return table recommendations Milestone recommendations
function ProgressionAlternatives.recommendMilestones(currentProgress, progressionType)
  if not currentProgress or not progressionType then return {} end
  
  local recommendations = {
    currentProgress = currentProgress,
    nextMilestone = nil,
    milestoneDistance = 0,
    recommendations = {}
  }
  
  if progressionType == "challenge" then
    if currentProgress < 3 then
      recommendations.nextMilestone = "Rising Hero (3 bosses)"
      recommendations.milestoneDistance = 3 - currentProgress
    elseif currentProgress < 10 then
      recommendations.nextMilestone = "Legendary Hero (10 bosses)"
      recommendations.milestoneDistance = 10 - currentProgress
    else
      recommendations.nextMilestone = "Ultimate Champion (25 bosses)"
      recommendations.milestoneDistance = 25 - currentProgress
    end
  end
  
  return recommendations
end

-- ============================================================================
-- FEATURE 4: CHALLENGE MODES (~220 LOC)
-- ============================================================================

local ChallengeModes = {}

---Define custom challenge mode
---@param modeName string Name of challenge
---@param settings table Challenge settings
---@return table challengeMode Challenge definition
function ChallengeModes.defineCustomChallenge(modeName, settings)
  if not modeName or not settings then return {} end
  
  local challengeMode = {
    name = modeName,
    enabled = true,
    levelCap = settings.levelCap or 99,
    scalingMultiplier = settings.scalingMultiplier or 1.0,
    restrictedItems = settings.restrictedItems or {},
    restrictedAbilities = settings.restrictedAbilities or {},
    bonusMultiplier = settings.bonusMultiplier or 1.0,
    description = settings.description or "Custom challenge mode"
  }
  
  return challengeMode
end

---Create preset challenge configuration
---@param presetName string Preset type: "speedrun", "hardcore", "ironman", "minimal"
---@return table preset Challenge preset
function ChallengeModes.createPresetChallenge(presetName)
  if not presetName then presetName = "standard" end
  
  local presets = {
    speedrun = {
      name = "Speedrun Mode",
      levelCap = 40,
      scalingMultiplier = 0.8,
      bonusMultiplier = 1.5,
      restrictions = {"Long animations disabled"},
      description = "Optimized for speed with level restrictions"
    },
    hardcore = {
      name = "Hardcore Mode",
      levelCap = 50,
      scalingMultiplier = 1.0,
      bonusMultiplier = 1.0,
      restrictions = {"Limited saves", "Permadeath enabled"},
      description = "Challenging run with permanent consequences"
    },
    ironman = {
      name = "Ironman Mode",
      levelCap = 99,
      scalingMultiplier = 0.9,
      bonusMultiplier = 2.0,
      restrictions = {"Single save slot", "No reset"},
      description = "Ultimate single-run challenge"
    },
    minimal = {
      name = "Minimal Stats Mode",
      levelCap = 25,
      scalingMultiplier = 0.5,
      bonusMultiplier = 1.0,
      restrictions = {"No equipment bonuses", "Limited abilities"},
      description = "Pure strategy with minimal stat growth"
    }
  }
  
  return presets[presetName] or presets.speedrun
end

---Verify challenge mode compliance
---@param saveData table Save file data
---@param challenge table Challenge definition
---@return table compliance Compliance check results
function ChallengeModes.verifyChallengCompliance(saveData, challenge)
  if not saveData or not challenge then return {} end
  
  local compliance = {
    challenge = challenge.name,
    compliant = true,
    violations = {},
    warnings = {}
  }
  
  -- Check level cap
  if saveData.characters then
    for _, char in ipairs(saveData.characters) do
      if (char.level or 1) > challenge.levelCap then
        table.insert(compliance.violations, 
          char.name .. " exceeds level cap (" .. challenge.levelCap .. ")")
        compliance.compliant = false
      end
    end
  end
  
  -- Check restricted items
  if saveData.equipment and challenge.restrictedItems then
    for _, item in ipairs(saveData.equipment) do
      for _, restricted in ipairs(challenge.restrictedItems) do
        if item == restricted then
          table.insert(compliance.violations, 
            "Restricted item equipped: " .. item)
          compliance.compliant = false
        end
      end
    end
  end
  
  return compliance
end

---Calculate challenge reward multiplier
---@param challenge table Challenge definition
---@param successFlag boolean Whether challenge completed successfully
---@return number multiplier Reward multiplier
function ChallengeModes.calculateChallengeReward(challenge, successFlag)
  if not challenge then return 1.0 end
  
  local multiplier = challenge.bonusMultiplier or 1.0
  
  if not successFlag then
    multiplier = multiplier * 0.5  -- 50% of bonus if failed
  end
  
  return multiplier
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1+",
  LevelCapEnforcement = LevelCapEnforcement,
  StatScaling = StatScaling,
  ProgressionAlternatives = ProgressionAlternatives,
  ChallengeModes = ChallengeModes,
  
  -- Feature completion
  features = {
    levelCapEnforcement = true,
    statScaling = true,
    progressionAlternatives = true,
    challengeModes = true
  }
}
