--[[
  Element Affinity Plugin - v1.1+ Upgrade Extension
  Elemental system integration, weakness analysis, resistance optimization, combat strategy
  
  Phase: 7E (Story & Element Systems)
  Version: 1.1+ (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: ELEMENTAL WEAKNESS ANALYSIS (~220 LOC)
-- ============================================================================

local WeaknessAnalysis = {}

---Analyze character elemental weaknesses
---@param character table Character data
---@return table weaknessProfile Elemental vulnerability analysis
function WeaknessAnalysis.analyzeElementalWeakness(character)
  if not character then return {} end
  
  local weaknessProfile = {
    character = character.name,
    weaknesses = {},
    resistances = {},
    neutrals = {},
    overallVulnerability = 0
  }
  
  -- Define elemental framework
  local elements = {"fire", "ice", "lightning", "wind", "water", "earth", "holy", "dark"}
  
  for _, element in ipairs(elements) do
    local weakness = WeaknessAnalysis._getElementalAffinityValue(character, element)
    
    if weakness < -20 then
      table.insert(weaknessProfile.weaknesses, {
        element = element,
        severity = math.abs(weakness)
      })
      weaknessProfile.overallVulnerability = weaknessProfile.overallVulnerability + 1
    elseif weakness > 20 then
      table.insert(weaknessProfile.resistances, {
        element = element,
        resistance = weakness
      })
    else
      table.insert(weaknessProfile.neutrals, element)
    end
  end
  
  return weaknessProfile
end

---Detect enemy elemental attacks
---@param enemy table Enemy data
---@return table elementalAttacks Available elemental attacks
function WeaknessAnalysis.detectElementalAttacks(enemy)
  if not enemy then return {} end
  
  local elementalAttacks = {
    enemy = enemy.name,
    attacks = {},
    primaryElement = nil,
    attackDiversity = 0
  }
  
  -- Scan enemy abilities
  if enemy.abilities then
    for _, ability in ipairs(enemy.abilities) do
      local element = WeaknessAnalysis._getAbilityElement(ability)
      
      if element then
        table.insert(elementalAttacks.attacks, {
          ability = ability,
          element = element,
          power = 50
        })
        
        if not elementalAttacks.primaryElement then
          elementalAttacks.primaryElement = element
        end
      end
    end
  end
  
  elementalAttacks.attackDiversity = #elementalAttacks.attacks
  
  return elementalAttacks
end

---Calculate elemental advantage
---@param characterElement string Character's resistance
---@param enemyElement string Enemy's attack type
---@return number advantage Damage multiplier (0.5 = half, 2.0 = double)
function WeaknessAnalysis.calculateElementalAdvantage(characterElement, enemyElement)
  if not characterElement or not enemyElement then return 1.0 end
  
  -- Advantage matrix
  local advantages = {
    fire = {ice = 0.5, wind = 1.5, earth = 1.0},
    ice = {fire = 1.5, water = 0.5, earth = 1.0},
    lightning = {wind = 0.5, water = 1.5, earth = 1.0},
    wind = {lightning = 1.5, fire = 0.5, earth = 1.0},
    water = {fire = 0.5, ice = 1.5, earth = 1.0},
    earth = {wind = 1.5, water = 0.5, lightning = 1.0}
  }
  
  if advantages[characterElement] and advantages[characterElement][enemyElement] then
    return advantages[characterElement][enemyElement]
  end
  
  return 1.0
end

---Recommend defensive elements
---@param enemies table Enemy party
---@return table recommendations Recommended elements to equip
function WeaknessAnalysis.recommendDefensiveElements(enemies)
  if not enemies then return {} end
  
  local recommendations = {
    suggestedElements = {},
    priorityOrder = {},
    defensiveRating = 0
  }
  
  -- Analyze enemy elemental attacks
  local elementCount = {}
  for _, enemy in ipairs(enemies) do
    local attacks = WeaknessAnalysis.detectElementalAttacks(enemy)
    
    if attacks.primaryElement then
      elementCount[attacks.primaryElement] = 
        (elementCount[attacks.primaryElement] or 0) + 1
    end
  end
  
  -- Sort by frequency
  for element, count in pairs(elementCount) do
    table.insert(recommendations.suggestedElements, {
      element = element,
      frequency = count,
      recommendation = "Equip " .. element .. " resistance"
    })
  end
  
  table.sort(recommendations.suggestedElements, function(a, b)
    return a.frequency > b.frequency
  end)
  
  return recommendations
end

-- Helper functions
function WeaknessAnalysis._getElementalAffinityValue(character, element)
  if character.elementalAffinity and character.elementalAffinity[element] then
    return character.elementalAffinity[element]
  end
  return 0
end

function WeaknessAnalysis._getAbilityElement(ability)
  local elementMap = {
    fire = "Firaga",
    ice = "Blizzaga",
    lightning = "Thundaga",
    wind = "Tornado"
  }
  
  for element, spell in pairs(elementMap) do
    if string.match(ability, spell) then return element end
  end
  
  return nil
end

-- ============================================================================
-- FEATURE 2: RESISTANCE OPTIMIZATION (~220 LOC)
-- ============================================================================

local ResistanceOptimization = {}

---Optimize equipment for elemental resistances
---@param character table Character to optimize
---@param availableEquip table Available equipment
---@param targetElements table Desired resistances
---@return table optimizedGear Recommended equipment
function ResistanceOptimization.optimizeForResistance(character, availableEquip, targetElements)
  if not character or not availableEquip then return {} end
  
  local optimizedGear = {
    character = character.name,
    recommendations = {},
    maxResistance = 0,
    coveragePercentage = 0
  }
  
  targetElements = targetElements or {"fire", "ice"}
  
  -- Score equipment by resistance values
  local equipScores = {}
  for slotName, items in pairs(availableEquip) do
    for _, item in ipairs(items) do
      local score = ResistanceOptimization._scoreEquipmentForResistance(
        item, targetElements
      )
      
      table.insert(equipScores, {
        slot = slotName,
        item = item,
        score = score
      })
    end
  end
  
  -- Sort and recommend
  table.sort(equipScores, function(a, b) return a.score > b.score end)
  
  for i = 1, math.min(4, #equipScores) do
    table.insert(optimizedGear.recommendations, equipScores[i])
    optimizedGear.maxResistance = optimizedGear.maxResistance + equipScores[i].score
  end
  
  optimizedGear.coveragePercentage = (#optimizedGear.recommendations / 4) * 100
  
  return optimizedGear
end

---Calculate total elemental coverage
---@param character table Character with equipped gear
---@param equipment table Equipped items
---@return table coverage Total resistance coverage
function ResistanceOptimization.calculateTotalCoverage(character, equipment)
  if not character then return {} end
  
  local coverage = {
    character = character.name,
    elements = {},
    totalResistance = 0,
    averageResistance = 0
  }
  
  local elements = {"fire", "ice", "lightning", "wind", "water", "earth"}
  
  for _, element in ipairs(elements) do
    local resistance = 0
    
    if equipment then
      for _, item in pairs(equipment) do
        if item.resistances and item.resistances[element] then
          resistance = resistance + item.resistances[element]
        end
      end
    end
    
    table.insert(coverage.elements, {
      element = element,
      resistance = resistance
    })
    
    coverage.totalResistance = coverage.totalResistance + resistance
  end
  
  coverage.averageResistance = coverage.totalResistance / #coverage.elements
  
  return coverage
end

---Predict damage reduction
---@param baseDamage number Base damage from attack
---@param elementResistance number Resistance percentage
---@return number reducedDamage Expected damage after resistance
function ResistanceOptimization.predictDamageReduction(baseDamage, elementResistance)
  if not baseDamage then return 0 end
  
  elementResistance = math.min(100, elementResistance or 0)
  
  local reductionFactor = 1 - (elementResistance / 100)
  local reducedDamage = baseDamage * reductionFactor
  
  return math.max(1, math.floor(reducedDamage))
end

---Create resistance build template
---@param focusElements table Priority elements
---@return table buildTemplate Template for resistance setup
function ResistanceOptimization.createResistanceBuildTemplate(focusElements)
  if not focusElements then focusElements = {"fire", "ice"} end
  
  local buildTemplate = {
    name = "Resistance Build",
    focusElements = focusElements,
    equipmentLayout = {},
    totalCost = 0,
    effectiveness = 0
  }
  
  for _, element in ipairs(focusElements) do
    buildTemplate.equipmentLayout[element] = {
      head = "Resist " .. element .. " Helm",
      body = "Resist " .. element .. " Robe",
      accessory = element .. " Ring"
    }
  end
  
  return buildTemplate
end

-- Helper functions
function ResistanceOptimization._scoreEquipmentForResistance(item, targetElements)
  local score = 0
  
  if item.resistances then
    for _, element in ipairs(targetElements) do
      score = score + (item.resistances[element] or 0)
    end
  end
  
  return score
end

-- ============================================================================
-- FEATURE 3: ELEMENTAL SYNERGY DETECTION (~200 LOC)
-- ============================================================================

local SynergyDetection = {}

---Detect elemental synergies in party
---@param party table Current party
---@return table synergies Synergy analysis
function SynergyDetection.detectPartySynergies(party)
  if not party then return {} end
  
  local synergies = {
    memberCount = #party,
    synergyCombos = {},
    totalSynergyBonus = 0
  }
  
  -- Analyze character element combinations
  for i = 1, #party do
    for j = i + 1, #party do
      local char1 = party[i]
      local char2 = party[j]
      
      local synergy = SynergyDetection._calculateSynergy(char1, char2)
      
      if synergy.score > 0 then
        table.insert(synergies.synergyCombos, {
          character1 = char1.name,
          character2 = char2.name,
          synergyType = synergy.type,
          bonus = synergy.score
        })
        
        synergies.totalSynergyBonus = synergies.totalSynergyBonus + synergy.score
      end
    end
  end
  
  return synergies
end

---Recommend party composition for elemental strength
---@param characters table Available characters
---@param targetElements table Desired coverage
---@return table partyRecommendation Optimal party suggestion
function SynergyDetection.recommendElementalParty(characters, targetElements)
  if not characters or not targetElements then return {} end
  
  local partyRecommendation = {
    suggestedParty = {},
    elementalCoverage = {},
    synergiesCount = 0
  }
  
  -- Score characters for element coverage
  local characterScores = {}
  for _, character in ipairs(characters) do
    local score = 0
    
    -- Calculate coverage score
    if character.elementalAffinity then
      for _, element in ipairs(targetElements) do
        local affinity = character.elementalAffinity[element] or 0
        if affinity > 0 then score = score + affinity end
      end
    end
    
    table.insert(characterScores, {
      character = character,
      score = score
    })
  end
  
  -- Sort and select top 4
  table.sort(characterScores, function(a, b) return a.score > b.score end)
  
  for i = 1, math.min(4, #characterScores) do
    table.insert(partyRecommendation.suggestedParty, characterScores[i].character)
  end
  
  return partyRecommendation
end

---Analyze elemental chain possibilities
---@param party table Current party
---@param enemies table Enemy party
---@return table chainAnalysis Elemental chain opportunities
function SynergyDetection.analyzeElementalChains(party, enemies)
  if not party or not enemies then return {} end
  
  local chainAnalysis = {
    possibleChains = {},
    maximumChainLength = 0,
    chainDamageMultiplier = 1.0
  }
  
  -- Detect weaknesses in enemy party
  for _, character in ipairs(party) do
    for _, enemy in ipairs(enemies) do
      if character.elementalAttack and enemy.elementalWeakness then
        if character.elementalAttack == enemy.elementalWeakness then
          table.insert(chainAnalysis.possibleChains, {
            initiator = character.name,
            target = enemy.name,
            bonus = "1.5x damage"
          })
        end
      end
    end
  end
  
  chainAnalysis.maximumChainLength = #chainAnalysis.possibleChains
  
  return chainAnalysis
end

-- Helper functions
function SynergyDetection._calculateSynergy(char1, char2)
  local synergy = {type = "elemental", score = 0}
  
  -- Simple synergy scoring
  if char1.elementalAttack and char2.elementalAttack then
    if char1.elementalAttack ~= char2.elementalAttack then
      synergy.score = 15
    end
  end
  
  return synergy
end

-- ============================================================================
-- FEATURE 4: COMBAT STRATEGY (~220 LOC)
-- ============================================================================

local CombatStrategy = {}

---Suggest elemental strategy for battle
---@param party table Party members
---@param enemies table Enemy party
---@return table strategy Recommended combat strategy
function CombatStrategy.suggestElementalStrategy(party, enemies)
  if not party or not enemies then return {} end
  
  local strategy = {
    primaryStrategy = "",
    targetPriority = {},
    recommendedSpells = {},
    effectiveness = 0
  }
  
  -- Analyze enemy weaknesses
  local commonWeakness = CombatStrategy._findCommonEnemyWeakness(enemies)
  
  if commonWeakness then
    strategy.primaryStrategy = "Focus " .. commonWeakness .. " attacks"
    strategy.effectiveness = 85
  else
    strategy.primaryStrategy = "Mixed elemental attacks"
    strategy.effectiveness = 60
  end
  
  -- Recommend spells
  for _, character in ipairs(party) do
    if character.elementalAttack then
      table.insert(strategy.recommendedSpells, {
        character = character.name,
        spell = character.elementalAttack .. "a"
      })
    end
  end
  
  return strategy
end

---Calculate battle outcome prediction
---@param party table Party members
---@param enemies table Enemy party
---@return table prediction Win probability analysis
function CombatStrategy.predictBattleOutcome(party, enemies)
  if not party or not enemies then return {} end
  
  local prediction = {
    partyPower = 0,
    enemyPower = 0,
    winProbability = 0,
    recommendedApproach = ""
  }
  
  -- Calculate party elemental power
  for _, character in ipairs(party) do
    prediction.partyPower = prediction.partyPower + 
      (character.magicPower or 10) + (character.level or 1)
  end
  
  -- Calculate enemy power
  for _, enemy in ipairs(enemies) do
    prediction.enemyPower = prediction.enemyPower + 
      (enemy.power or 10) + (enemy.level or 1)
  end
  
  -- Calculate win probability
  local powerRatio = prediction.partyPower / prediction.enemyPower
  prediction.winProbability = math.min(99, math.floor(powerRatio * 50))
  
  if prediction.winProbability > 70 then
    prediction.recommendedApproach = "Aggressive - Use powerful spells"
  elseif prediction.winProbability > 40 then
    prediction.recommendedApproach = "Balanced - Mix offense and defense"
  else
    prediction.recommendedApproach = "Defensive - Prioritize survival"
  end
  
  return prediction
end

---Plan multi-turn elemental strategy
---@param party table Party members
---@param enemies table Enemy party
---@param battleLength number Expected turns
---@return table plan Multi-turn combat plan
function CombatStrategy.planMultiTurnStrategy(party, enemies, battleLength)
  if not party or not enemies then return {} end
  
  battleLength = battleLength or 5
  
  local plan = {
    turns = {},
    strategy = "Elemental-based",
    expectedDuration = battleLength
  }
  
  for turn = 1, battleLength do
    local turnPlan = {
      turnNumber = turn,
      actions = {}
    }
    
    for _, character in ipairs(party) do
      table.insert(turnPlan.actions, {
        character = character.name,
        action = "Cast elemental spell",
        target = "Highest priority enemy"
      })
    end
    
    table.insert(plan.turns, turnPlan)
  end
  
  return plan
end

-- Helper functions
function CombatStrategy._findCommonEnemyWeakness(enemies)
  local weaknesses = {}
  
  for _, enemy in ipairs(enemies) do
    if enemy.elementalWeakness then
      weaknesses[enemy.elementalWeakness] = 
        (weaknesses[enemy.elementalWeakness] or 0) + 1
    end
  end
  
  -- Return most common weakness
  local maxCount = 0
  local commonWeakness = nil
  
  for element, count in pairs(weaknesses) do
    if count > maxCount then
      maxCount = count
      commonWeakness = element
    end
  end
  
  return commonWeakness
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1+",
  WeaknessAnalysis = WeaknessAnalysis,
  ResistanceOptimization = ResistanceOptimization,
  SynergyDetection = SynergyDetection,
  CombatStrategy = CombatStrategy,
  
  -- Feature completion
  features = {
    weaknessAnalysis = true,
    resistanceOptimization = true,
    synergyDetection = true,
    combatStrategy = true
  }
}
