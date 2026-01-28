--[[
  Party Optimizer Plugin - v1.1+ Upgrade Extension
  Equipment integration, growth prediction, esper optimization, and boss strategies
  
  Phase: 7D (High-Impact Tracking)
  Version: 1.1+ (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: EQUIPMENT INTEGRATION (~200 LOC)
-- ============================================================================

local EquipmentIntegration = {}

---Analyze equipment bonuses in party calculations
---@param character table Character to analyze
---@param equipment table Equipped gear
---@return table analysis Equipment impact analysis
function EquipmentIntegration.analyzeEquipmentBonuses(character, equipment)
  if not character or not equipment then return {} end
  
  local analysis = {
    baseStats = character.stats or {},
    equipmentBonuses = {},
    totalStats = {},
    bonusPercentage = 0
  }
  
  -- Calculate bonuses per equipment slot
  local totalBonus = 0
  for slotName, item in pairs(equipment) do
    if item and item.bonuses then
      analysis.equipmentBonuses[slotName] = {
        item = item.name,
        bonuses = item.bonuses,
        impact = {}
      }
      
      for stat, bonus in pairs(item.bonuses) do
        totalBonus = totalBonus + bonus
        if not analysis.totalStats[stat] then
          analysis.totalStats[stat] = (character.stats and character.stats[stat]) or 0
        end
        analysis.totalStats[stat] = analysis.totalStats[stat] + bonus
      end
    end
  end
  
  return analysis
end

---Suggest optimal equipment per character
---@param character table Character to equip
---@param availableEquipment table Available gear
---@return table suggestions Recommended equipment
function EquipmentIntegration.suggestOptimalEquipment(character, availableEquipment)
  if not character or not availableEquipment then return {} end
  
  local suggestions = {
    character = character.name,
    recommended = {},
    alternates = {},
    totalBonus = 0
  }
  
  -- Analyze role and suggest gear
  local role = character.role or "warrior"
  
  for slotName, availableGear in pairs(availableEquipment) do
    local bestItem = nil
    local bestScore = 0
    
    for _, item in ipairs(availableGear) do
      local score = EquipmentIntegration._scoreEquipmentForRole(item, role)
      if score > bestScore then
        bestScore = score
        bestItem = item
      end
    end
    
    if bestItem then
      suggestions.recommended[slotName] = bestItem
    end
  end
  
  return suggestions
end

---Analyze equipment synergies
---@param equipment table All equipped items
---@return table synergies Equipment synergy bonuses
function EquipmentIntegration.equipmentSynergyAnalysis(equipment)
  if not equipment then return {} end
  
  local synergies = {
    combos = {},
    setBonuses = {},
    totalSynergyBonus = 0
  }
  
  -- Check for equipment set bonuses
  local elementalItems = {}
  for _, item in pairs(equipment) do
    if item and item.element then
      elementalItems[item.element] = (elementalItems[item.element] or 0) + 1
    end
  end
  
  for element, count in pairs(elementalItems) do
    if count >= 2 then
      table.insert(synergies.setBonuses, {
        element = element,
        count = count,
        bonus = "Elemental resistance +20%"
      })
      synergies.totalSynergyBonus = synergies.totalSynergyBonus + 20
    end
  end
  
  return synergies
end

---Estimate equipment impact on stats
---@param character table Character to analyze
---@param newEquipment table Equipment to try
---@return table impact Impact projection
function EquipmentIntegration.estimateEquipmentImpact(character, newEquipment)
  if not character or not newEquipment then return {} end
  
  local impact = {
    originalStats = character.stats or {},
    withEquipment = {},
    improvements = {},
    totalImprovement = 0
  }
  
  -- Calculate new stats
  for stat, value in pairs(character.stats or {}) do
    impact.withEquipment[stat] = value
    
    for _, item in pairs(newEquipment) do
      if item and item.bonuses and item.bonuses[stat] then
        impact.withEquipment[stat] = impact.withEquipment[stat] + item.bonuses[stat]
      end
    end
    
    impact.improvements[stat] = impact.withEquipment[stat] - value
    impact.totalImprovement = impact.totalImprovement + impact.improvements[stat]
  end
  
  return impact
end

-- Helper function
function EquipmentIntegration._scoreEquipmentForRole(item, role)
  local score = 0
  
  if role == "warrior" and item.strength then
    score = score + item.strength * 2
  elseif role == "mage" and item.magic then
    score = score + item.magic * 2
  end
  
  if item.defense then score = score + item.defense end
  
  return score
end

-- ============================================================================
-- FEATURE 2: GROWTH PREDICTION (~250 LOC)
-- ============================================================================

local GrowthPrediction = {}

---Predict stat growth at future levels
---@param character table Character to project
---@param targetLevel number Level to project to
---@return table projection Growth projection
function GrowthPrediction.predictStatGrowth(character, targetLevel)
  if not character or not targetLevel then return {} end
  
  local currentLevel = character.level or 1
  local projection = {
    character = character.name,
    currentLevel = currentLevel,
    targetLevel = targetLevel,
    levelRange = targetLevel - currentLevel,
    projectedStats = {},
    levelUpDetails = {}
  }
  
  -- Project each level
  for level = currentLevel + 1, targetLevel do
    local levelStats = {
      level = level,
      statGains = {}
    }
    
    -- Estimate gains
    for stat, value in pairs(character.stats or {}) do
      local growthRate = GrowthPrediction._getGrowthRate(character, stat)
      levelStats.statGains[stat] = math.floor(growthRate)
    end
    
    table.insert(projection.levelUpDetails, levelStats)
  end
  
  -- Calculate final stats
  for stat, value in pairs(character.stats or {}) do
    local finalValue = value
    for _, levelStats in ipairs(projection.levelUpDetails) do
      finalValue = finalValue + levelStats.statGains[stat]
    end
    projection.projectedStats[stat] = finalValue
  end
  
  return projection
end

---Predict when new abilities will unlock
---@param character table Character to analyze
---@param currentLevel number Current level
---@return table abilitiesTimeline Ability unlock timeline
function GrowthPrediction.predictAbilityGain(character, currentLevel)
  if not character then return {} end
  
  currentLevel = currentLevel or character.level or 1
  
  local timeline = {
    character = character.name,
    currentLevel = currentLevel,
    unlockedAbilities = {},
    futureAbilities = {}
  }
  
  -- Standard ability unlock levels (can be customized per character)
  local abilityLevels = {
    [10] = "Basic Ability 1",
    [20] = "Intermediate Ability",
    [30] = "Advanced Ability",
    [40] = "Rare Ability",
    [50] = "Ultimate Ability"
  }
  
  for levelReq, ability in pairs(abilityLevels) do
    if levelReq <= currentLevel then
      table.insert(timeline.unlockedAbilities, {
        level = levelReq,
        ability = ability
      })
    else
      table.insert(timeline.futureAbilities, {
        level = levelReq,
        ability = ability,
        levelsAway = levelReq - currentLevel
      })
    end
  end
  
  return timeline
end

---Analyze growth potential for each stat
---@param character table Character to analyze
---@return table potential Growth potential ranking
function GrowthPrediction.analyzeGrowthPotential(character)
  if not character then return {} end
  
  local potential = {
    character = character.name,
    stats = {},
    highestGrowth = nil,
    lowestGrowth = nil
  }
  
  local maxGrowth = 0
  local minGrowth = 999
  
  for stat, value in pairs(character.stats or {}) do
    local growthRate = GrowthPrediction._getGrowthRate(character, stat)
    
    table.insert(potential.stats, {
      stat = stat,
      currentValue = value,
      growthRate = growthRate,
      levelTo99 = math.ceil((99 - value) / growthRate)
    })
    
    if growthRate > maxGrowth then
      maxGrowth = growthRate
      potential.highestGrowth = stat
    end
    if growthRate < minGrowth then
      minGrowth = growthRate
      potential.lowestGrowth = stat
    end
  end
  
  return potential
end

---Optimize character for late game
---@param character table Character to optimize
---@param targetLevel number Target level
---@return table strategy Late-game strategy
function GrowthPrediction.optimizeForLateGame(character, targetLevel)
  if not character or not targetLevel then return {} end
  
  local strategy = {
    character = character.name,
    targetLevel = targetLevel,
    recommendations = {},
    buildStrategy = ""
  }
  
  local potential = GrowthPrediction.analyzeGrowthPotential(character)
  
  -- Suggest focusing on highest growth stats
  strategy.recommendations = {
    "Focus training on " .. (potential.highestGrowth or "balanced") .. " stat",
    "Save esper points for late-game abilities",
    "Equip items that boost " .. (potential.highestGrowth or "core") .. " stats"
  }
  
  return strategy
end

-- Helper function
function GrowthPrediction._getGrowthRate(character, stat)
  local rates = {
    strength = 0.8,
    speed = 0.6,
    stamina = 0.7,
    magic = 0.9,
    defense = 0.5
  }
  return rates[stat] or 0.5
end

-- ============================================================================
-- FEATURE 3: ESPER OPTIMIZATION (~180 LOC)
-- ============================================================================

local EsperOptimization = {}

---Suggest optimal esper per character
---@param character table Character to optimize
---@param availableEspers table Available espers
---@return table suggestions Esper recommendations
function EsperOptimization.suggestOptimalEspers(character, availableEspers)
  if not character or not availableEspers then return {} end
  
  local suggestions = {
    character = character.name,
    primary = nil,
    secondary = nil,
    recommendations = {}
  }
  
  local role = character.role or "warrior"
  
  for _, esper in ipairs(availableEspers) do
    local score = EsperOptimization._scoreEsperForRole(esper, role)
    
    if not suggestions.primary or score > suggestions.primary.score then
      suggestions.secondary = suggestions.primary
      suggestions.primary = {
        esper = esper.name,
        score = score,
        statBonus = esper.statBonus or {}
      }
    elseif not suggestions.secondary or score > suggestions.secondary.score then
      suggestions.secondary = {
        esper = esper.name,
        score = score,
        statBonus = esper.statBonus or {}
      }
    end
  end
  
  return suggestions
end

---Analyze esper synergies
---@param espers table Equipped espers
---@return table synergies Synergy analysis
function EsperOptimization.esperSynergyAnalysis(espers)
  if not espers then return {} end
  
  local synergies = {
    combos = {},
    powerLevel = 0
  }
  
  -- Check for element combinations
  local elements = {}
  for _, esper in ipairs(espers) do
    if esper.element then
      elements[esper.element] = (elements[esper.element] or 0) + 1
    end
  end
  
  for element, count in pairs(elements) do
    if count >= 2 then
      table.insert(synergies.combos, {
        element = element,
        count = count,
        bonus = "Elemental affinity bonus +15%"
      })
    end
  end
  
  return synergies
end

---Predict stat bonuses from espers
---@param character table Character
---@param espers table Equipped espers
---@return table predictions Stat predictions
function EsperOptimization.predictEsperStats(character, espers)
  if not character or not espers then return {} end
  
  local predictions = {
    baseStats = character.stats or {},
    esperBonuses = {},
    totalStats = {}
  }
  
  for stat, value in pairs(character.stats or {}) do
    predictions.totalStats[stat] = value
    
    for _, esper in ipairs(espers) do
      if esper.statBonus and esper.statBonus[stat] then
        local bonus = esper.statBonus[stat]
        predictions.esperBonuses[stat] = (predictions.esperBonuses[stat] or 0) + bonus
        predictions.totalStats[stat] = predictions.totalStats[stat] + bonus
      end
    end
  end
  
  return predictions
end

-- Helper function
function EsperOptimization._scoreEsperForRole(esper, role)
  local score = 0
  
  if role == "warrior" and esper.defense then
    score = score + esper.defense * 1.5
  elseif role == "mage" and esper.magic then
    score = score + esper.magic * 1.5
  end
  
  if esper.level then score = score + esper.level end
  
  return score
end

-- ============================================================================
-- FEATURE 4: BOSS-SPECIFIC STRATEGIES (~220 LOC)
-- ============================================================================

local BossStrategies = {}

---Analyze party for specific boss fight
---@param party table Current party
---@param bossData table Boss information
---@return table analysis Boss-specific analysis
function BossStrategies.analyzeForBoss(party, bossData)
  if not party or not bossData then return {} end
  
  local analysis = {
    boss = bossData.name,
    difficulty = BossStrategies._calculateDifficulty(party, bossData),
    strengths = {},
    weaknesses = {},
    recommendations = {}
  }
  
  -- Analyze party strengths vs boss
  for _, member in ipairs(party) do
    if member.stats then
      if member.stats.magic and bossData.weakTo and bossData.weakTo.magic then
        table.insert(analysis.strengths, member.name .. " has high magic")
      end
    end
  end
  
  -- Identify weaknesses
  if bossData.strongAgainst then
    for weakness, _ in pairs(bossData.strongAgainst) do
      for _, member in ipairs(party) do
        if member.stats and member.stats[weakness] then
          table.insert(analysis.weaknesses, 
            member.name .. " is weak to " .. weakness)
        end
      end
    end
  end
  
  return analysis
end

---Suggest best party for boss fight
---@param characters table Available characters
---@param bossData table Boss information
---@return table partyRecommendation Suggested party
function BossStrategies.suggestBossBattleParty(characters, bossData)
  if not characters or not bossData then return {} end
  
  local recommendation = {
    boss = bossData.name,
    suggestedParty = {},
    reasoning = {}
  }
  
  -- Score each character for this fight
  local scored = {}
  for _, character in ipairs(characters) do
    local score = BossStrategies._scoreCharacterForBoss(character, bossData)
    table.insert(scored, {
      character = character,
      score = score
    })
  end
  
  -- Sort by score and select top 4
  table.sort(scored, function(a, b) return a.score > b.score end)
  
  for i = 1, math.min(4, #scored) do
    table.insert(recommendation.suggestedParty, scored[i].character)
  end
  
  return recommendation
end

---Predict boss battle difficulty
---@param party table Current party
---@param bossData table Boss information
---@return number difficulty Difficulty rating 1-10
function BossStrategies.predictBossDifficulty(party, bossData)
  if not party or not bossData then return 5 end
  
  local partyPower = 0
  for _, member in ipairs(party) do
    partyPower = partyPower + (member.level or 1) * 10
  end
  
  local bossPower = (bossData.level or 50) * 15
  
  local difficulty = 5 + math.floor((bossPower - partyPower) / 50)
  
  return math.max(1, math.min(10, difficulty))
end

-- Helper functions
function BossStrategies._calculateDifficulty(party, boss)
  return BossStrategies.predictBossDifficulty(party, boss)
end

function BossStrategies._scoreCharacterForBoss(character, boss)
  local score = character.level or 1
  
  if character.stats and boss.weakTo then
    for weakness, multiplier in pairs(boss.weakTo) do
      if character.stats[weakness] and character.stats[weakness] > 70 then
        score = score + 50
      end
    end
  end
  
  return score
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1+",
  EquipmentIntegration = EquipmentIntegration,
  GrowthPrediction = GrowthPrediction,
  EsperOptimization = EsperOptimization,
  BossStrategies = BossStrategies,
  
  -- Feature completion
  features = {
    equipmentIntegration = true,
    growthPrediction = true,
    esperOptimization = true,
    bossStrategies = true
  }
}
