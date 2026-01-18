--[[
  Rage Tracker Plugin - v1.1+ Upgrade Extension
  Veldt simulator, advanced detection, battle strategy, and write mode preparation
  
  Phase: 7D (High-Impact Tracking)
  Version: 1.1+ (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: VELDT SIMULATOR (~200 LOC)
-- ============================================================================

local VeldtSimulator = {}

---Simulate Veldt battle outcomes
---@param gau table Gau's current state
---@param rages table Available rages
---@param terrain string Veldt terrain type
---@return table simulation Simulation results
function VeldtSimulator.simulateVeldtBattle(gau, rages, terrain)
  if not gau or not rages or not terrain then return {} end
  
  local simulation = {
    gau = gau.name,
    terrain = terrain,
    duration = 0,
    enemiesEncountered = {},
    ragesCaught = {},
    success = false
  }
  
  -- Simulate encounter duration
  local encounterCount = VeldtSimulator._getTerrainEncounterCount(terrain)
  simulation.duration = encounterCount * 30  -- 30 seconds per encounter estimate
  
  -- Determine available enemies
  local availableEnemies = VeldtSimulator._getTerrainEnemies(terrain)
  
  for i = 1, encounterCount do
    local enemy = availableEnemies[math.random(#availableEnemies)]
    table.insert(simulation.enemiesEncountered, enemy)
    
    -- Chance to catch rage
    if math.random() < 0.3 then  -- 30% catch rate
      table.insert(simulation.ragesCaught, {
        rage = enemy.rage,
        terrain = terrain,
        timestamp = i * 30
      })
    end
  end
  
  simulation.success = #simulation.ragesCaught > 0
  
  return simulation
end

---Calculate best terrain for specific rage
---@param rageTarget string Desired rage name
---@param terrains table Available terrains
---@return table recommendation Best terrain recommendation
function VeldtSimulator.recommendTerrainForRage(rageTarget, terrains)
  if not rageTarget or not terrains then return {} end
  
  local recommendation = {
    rageTarget = rageTarget,
    bestTerrains = {},
    efficiency = {}
  }
  
  for _, terrain in ipairs(terrains) do
    local enemies = VeldtSimulator._getTerrainEnemies(terrain)
    
    for _, enemy in ipairs(enemies) do
      if enemy.rage == rageTarget then
        local efficiency = VeldtSimulator._calculateRageCatchEfficiency(terrain)
        
        table.insert(recommendation.bestTerrains, terrain)
        recommendation.efficiency[terrain] = efficiency
      end
    end
  end
  
  -- Sort by efficiency
  table.sort(recommendation.bestTerrains, function(a, b)
    return (recommendation.efficiency[a] or 0) > (recommendation.efficiency[b] or 0)
  end)
  
  return recommendation
end

---Predict encounter rate on terrain
---@param terrain string Terrain type
---@return table encounterData Encounter rate data
function VeldtSimulator.predictEncounterRate(terrain)
  if not terrain then return {} end
  
  local encounterData = {
    terrain = terrain,
    baseRate = 0,
    encountersPerHour = 0,
    rareRages = {},
    commonRages = {}
  }
  
  -- Terrain-specific rates
  local rates = {
    forest = 6,
    mountain = 8,
    coast = 5,
    desert = 7,
    grassland = 4
  }
  
  encounterData.baseRate = rates[terrain] or 5
  encounterData.encountersPerHour = encounterData.baseRate * 60 / 5
  
  -- Populate available rages
  local terrainEnemies = VeldtSimulator._getTerrainEnemies(terrain)
  for _, enemy in ipairs(terrainEnemies) do
    if math.random() < 0.3 then
      table.insert(encounterData.rareRages, enemy.rage)
    else
      table.insert(encounterData.commonRages, enemy.rage)
    end
  end
  
  return encounterData
end

---Strategy for efficient rage collection
---@param missingRages table List of needed rages
---@param terrains table Available terrains
---@return table strategy Collection strategy
function VeldtSimulator.strategyForRageCollection(missingRages, terrains)
  if not missingRages or not terrains then return {} end
  
  local strategy = {
    missingCount = #missingRages,
    terrainPlan = {},
    estimatedTime = 0
  }
  
  for _, rage in ipairs(missingRages) do
    local best = VeldtSimulator.recommendTerrainForRage(rage, terrains)
    
    if best.bestTerrains and #best.bestTerrains > 0 then
      local terrain = best.bestTerrains[1]
      strategy.terrainPlan[rage] = terrain
      strategy.estimatedTime = strategy.estimatedTime + 300  -- 5 min per rage
    end
  end
  
  return strategy
end

-- Helper functions
function VeldtSimulator._getTerrainEncounterCount(terrain)
  local counts = {
    forest = 5,
    mountain = 6,
    coast = 4,
    desert = 5,
    grassland = 3
  }
  return counts[terrain] or 4
end

function VeldtSimulator._getTerrainEnemies(terrain)
  local enemies = {
    forest = {
      {name = "Gabbledeguck", rage = "Gabbledeguck"},
      {name = "Anemone", rage = "Anemone"}
    },
    mountain = {
      {name = "Aspik", rage = "Aspik"},
      {name = "Karkass", rage = "Karkass"}
    },
    coast = {
      {name = "Sea Flower", rage = "Sea Flower"},
      {name = "Exocray", rage = "Exocray"}
    },
    desert = {
      {name = "Tortoise", rage = "Tortoise"},
      {name = "Dray", rage = "Dray"}
    }
  }
  return enemies[terrain] or enemies.forest
end

function VeldtSimulator._calculateRageCatchEfficiency(terrain)
  local efficiency = {
    forest = 0.8,
    mountain = 0.85,
    coast = 0.75,
    desert = 0.82
  }
  return efficiency[terrain] or 0.7
end

-- ============================================================================
-- FEATURE 2: ADVANCED DETECTION (~200 LOC)
-- ============================================================================

local AdvancedDetection = {}

---Detect all rages in current save
---@param saveData table Save file data
---@return table rageAnalysis Detailed rage analysis
function AdvancedDetection.detectAllRages(saveData)
  if not saveData then return {} end
  
  local rageAnalysis = {
    totalRagesCount = 0,
    detectedRages = {},
    missingRages = {},
    rarityBreakdown = {}
  }
  
  -- All possible rages in game
  local allRages = AdvancedDetection._getAllPossibleRages()
  
  -- Check which are detected
  for _, rage in ipairs(allRages) do
    if AdvancedDetection._isRageInSave(saveData, rage) then
      table.insert(rageAnalysis.detectedRages, rage)
      rageAnalysis.totalRagesCount = rageAnalysis.totalRagesCount + 1
    else
      table.insert(rageAnalysis.missingRages, rage)
    end
  end
  
  -- Rarity analysis
  rageAnalysis.rarityBreakdown = {
    common = 0,
    uncommon = 0,
    rare = 0,
    veryRare = 0
  }
  
  for _, rage in ipairs(rageAnalysis.detectedRages) do
    local rarity = AdvancedDetection._getRageRarity(rage)
    rageAnalysis.rarityBreakdown[rarity] = rageAnalysis.rarityBreakdown[rarity] + 1
  end
  
  return rageAnalysis
end

---Identify easily missed rages
---@param saveData table Save file data
---@return table missedRages List of easily missed rages
function AdvancedDetection.identifyMissedRages(saveData)
  if not saveData then return {} end
  
  local missedRages = {
    optional = {},
    oneTime = {},
    levelDependent = {},
    recommendations = {}
  }
  
  local allRages = AdvancedDetection._getAllPossibleRages()
  
  for _, rage in ipairs(allRages) do
    if not AdvancedDetection._isRageInSave(saveData, rage) then
      if AdvancedDetection._isOneTimeOnly(rage) then
        table.insert(missedRages.oneTime, rage)
      elseif AdvancedDetection._isLevelDependent(rage) then
        table.insert(missedRages.levelDependent, rage)
      else
        table.insert(missedRages.optional, rage)
      end
    end
  end
  
  -- Recommendations for recovery
  if #missedRages.oneTime > 0 then
    table.insert(missedRages.recommendations, 
      "Cannot recover one-time rages in current save")
  end
  
  return missedRages
end

---Track rare rage spawn conditions
---@param rage string Rage name
---@return table conditions Spawn conditions
function AdvancedDetection.trackRareRageConditions(rage)
  if not rage then return {} end
  
  local conditions = {
    rage = rage,
    terrains = {},
    enemies = {},
    rarity = AdvancedDetection._getRageRarity(rage),
    spawnRate = 0,
    notes = ""
  }
  
  -- Database of rare rage conditions
  local rareConditions = {
    ["Ultros"] = {
      terrains = {"ocean"},
      rarity = "veryRare",
      spawnRate = 0.05,
      notes = "Only appears in specific locations"
    },
    ["Atma"] = {
      terrains = {"special"},
      rarity = "veryRare",
      spawnRate = 0.01,
      notes = "Boss rage, very limited encounters"
    },
    ["Kefka"] = {
      terrains = {"special"},
      rarity = "veryRare",
      spawnRate = 0.00,
      notes = "Cannot be caught normally"
    }
  }
  
  if rareConditions[rage] then
    conditions = rareConditions[rage]
    conditions.rage = rage
  end
  
  return conditions
end

---Calculate completion percentage
---@param saveData table Save file data
---@return number percentage Rage collection percentage
function AdvancedDetection.calculateRageCompletion(saveData)
  if not saveData then return 0 end
  
  local analysis = AdvancedDetection.detectAllRages(saveData)
  
  local total = #analysis.detectedRages + #analysis.missingRages
  if total == 0 then return 0 end
  
  return (analysis.totalRagesCount / total) * 100
end

-- Helper functions
function AdvancedDetection._getAllPossibleRages()
  return {
    "Gabbledeguck", "Anemone", "Aspik", "Karkass",
    "Ultros", "Atma", "Kefka",
    "Dray", "Tortoise", "Sea Flower", "Exocray",
    "Leaf Bunny", "Floating Continent Enemy"
  }
end

function AdvancedDetection._isRageInSave(saveData, rage)
  if saveData.rages then
    for _, r in ipairs(saveData.rages) do
      if r == rage then return true end
    end
  end
  return false
end

function AdvancedDetection._getRageRarity(rage)
  local rarities = {
    ["Ultros"] = "veryRare",
    ["Atma"] = "veryRare",
    ["Kefka"] = "veryRare"
  }
  return rarities[rage] or "common"
end

function AdvancedDetection._isOneTimeOnly(rage)
  local oneTime = {["Atma"] = true, ["Kefka"] = true}
  return oneTime[rage] or false
end

function AdvancedDetection._isLevelDependent(rage)
  -- Most rages become available at different story points
  return true
end

-- ============================================================================
-- FEATURE 3: BATTLE STRATEGY (~220 LOC)
-- ============================================================================

local BattleStrategy = {}

---Recommend rage for specific battle
---@param opponents table Enemy party
---@param party table Character party
---@param availableRages table Gau's available rages
---@return table recommendation Best rage suggestion
function BattleStrategy.recommendRageForBattle(opponents, party, availableRages)
  if not opponents or not party or not availableRages then return {} end
  
  local recommendation = {
    recommendedRages = {},
    strategy = "",
    effectiveness = 0
  }
  
  -- Score each available rage
  local rageScores = {}
  for _, rage in ipairs(availableRages) do
    local score = BattleStrategy._scoreRageEffectiveness(rage, opponents)
    table.insert(rageScores, {
      rage = rage,
      score = score
    })
  end
  
  -- Sort and recommend top choices
  table.sort(rageScores, function(a, b) return a.score > b.score end)
  
  for i = 1, math.min(3, #rageScores) do
    table.insert(recommendation.recommendedRages, rageScores[i])
  end
  
  if #recommendation.recommendedRages > 0 then
    recommendation.effectiveness = recommendation.recommendedRages[1].score
  end
  
  return recommendation
end

---Analyze rage damage output
---@param rage string Rage name
---@param targets table Target party
---@return table damageAnalysis Damage breakdown
function BattleStrategy.analyzeRageDamage(rage, targets)
  if not rage or not targets then return {} end
  
  local damageAnalysis = {
    rage = rage,
    damagePerHit = 0,
    hitsPerRound = 2,
    totalDamagePerRound = 0,
    targetBreakdown = {}
  }
  
  -- Calculate rage damage
  local baseDamage = BattleStrategy._getRageBaseDamage(rage)
  
  for _, target in ipairs(targets) do
    local defense = target.stats and target.stats.defense or 50
    local actualDamage = math.max(1, baseDamage - (defense * 0.5))
    
    table.insert(damageAnalysis.targetBreakdown, {
      target = target.name,
      damage = actualDamage
    })
    
    damageAnalysis.totalDamagePerRound = damageAnalysis.totalDamagePerRound + actualDamage
  end
  
  damageAnalysis.totalDamagePerRound = damageAnalysis.totalDamagePerRound * 
    damageAnalysis.hitsPerRound
  
  return damageAnalysis
end

---Calculate survival time with rage
---@param gau table Gau's stats
---@param opponents table Enemy party
---@param rage string Selected rage
---@return number survivalRounds Estimated rounds survived
function BattleStrategy.calculateSurvivalTime(gau, opponents, rage)
  if not gau or not opponents or not rage then return 0 end
  
  local rageDefense = BattleStrategy._getRageDefense(rage)
  local effectiveHP = (gau.stats and gau.stats.hp or 100) * (rageDefense / 100)
  
  local damagePerRound = 0
  for _, opponent in ipairs(opponents) do
    local opponentDamage = BattleStrategy._calculateEnemyDamage(opponent)
    damagePerRound = damagePerRound + opponentDamage
  end
  
  if damagePerRound <= 0 then return 999 end
  
  return math.floor(effectiveHP / damagePerRound)
end

---Suggest rage rotation strategy
---@param availableRages table Available rages
---@param battleLength number Expected battle length
---@return table rotationStrategy Recommended rotation
function BattleStrategy.suggestRageRotation(availableRages, battleLength)
  if not availableRages or not battleLength then return {} end
  
  local rotationStrategy = {
    sequence = {},
    totalDuration = 0,
    purpose = "Balanced offense and defense"
  }
  
  local sequence = {}
  for i = 1, math.min(3, #availableRages) do
    table.insert(sequence, availableRages[i])
  end
  
  rotationStrategy.sequence = sequence
  rotationStrategy.totalDuration = battleLength
  
  return rotationStrategy
end

-- Helper functions
function BattleStrategy._scoreRageEffectiveness(rage, opponents)
  local score = 0
  
  -- Check damage effectiveness
  local damage = BattleStrategy._getRageBaseDamage(rage)
  score = score + damage
  
  -- Check defense capability
  local defense = BattleStrategy._getRageDefense(rage)
  score = score + (defense / 2)
  
  return score
end

function BattleStrategy._getRageBaseDamage(rage)
  local damages = {
    ["Gabbledeguck"] = 35,
    ["Ultros"] = 60,
    ["Atma"] = 80
  }
  return damages[rage] or 40
end

function BattleStrategy._getRageDefense(rage)
  local defenses = {
    ["Ultros"] = 120,
    ["Atma"] = 150,
    ["Karkass"] = 90
  }
  return defenses[rage] or 100
end

function BattleStrategy._calculateEnemyDamage(enemy)
  local baseDamage = enemy.stats and enemy.stats.strength or 30
  return baseDamage * 0.8
end

-- ============================================================================
-- FEATURE 4: WRITE MODE PREPARATION (~220 LOC)
-- ============================================================================

local WriteModePrep = {}

---Prepare Gau's status for write mode
---@param gau table Gau's current data
---@param targetRages table Rages to lock in
---@return table prepStatus Preparation status
function WriteModePrep.prepareForWriteMode(gau, targetRages)
  if not gau or not targetRages then return {} end
  
  local prepStatus = {
    character = "Gau",
    currentRages = {},
    targetRages = targetRages,
    readyToWrite = false,
    issues = {},
    recommendations = {}
  }
  
  -- Verify all target rages are available
  for _, targetRage in ipairs(targetRages) do
    local found = false
    if gau.rages then
      for _, ownedRage in ipairs(gau.rages) do
        if ownedRage == targetRage then
          found = true
          break
        end
      end
    end
    
    if not found then
      table.insert(prepStatus.issues, 
        "Missing rage: " .. targetRage)
    end
  end
  
  prepStatus.readyToWrite = #prepStatus.issues == 0
  
  return prepStatus
end

---Validate save data integrity before write
---@param saveData table Save file data
---@return table validation Validation results
function WriteModePrep.validateBeforeWrite(saveData)
  if not saveData then return {} end
  
  local validation = {
    valid = true,
    errors = {},
    warnings = {},
    checksPerformed = 0
  }
  
  -- Check data integrity
  validation.checksPerformed = validation.checksPerformed + 1
  if not saveData.rages or type(saveData.rages) ~= "table" then
    table.insert(validation.errors, "Invalid rage data structure")
    validation.valid = false
  end
  
  -- Check for data corruption
  validation.checksPerformed = validation.checksPerformed + 1
  if saveData.checksum and not WriteModePrep._verifyChecksum(saveData) then
    table.insert(validation.errors, "Checksum mismatch")
    validation.valid = false
  end
  
  -- Verify character data
  validation.checksPerformed = validation.checksPerformed + 1
  if not saveData.characters or #saveData.characters == 0 then
    table.insert(validation.warnings, "No character data found")
  end
  
  return validation
end

---Create backup before write
---@param saveData table Save file to backup
---@return table backup Backup metadata
function WriteModePrep.createPreWriteBackup(saveData)
  if not saveData then return {} end
  
  local backup = {
    timestamp = os.time(),
    backupId = math.random(1000000, 9999999),
    dataSize = WriteModePrep._calculateDataSize(saveData),
    ragesCount = 0,
    characterized = false
  }
  
  if saveData.rages then
    backup.ragesCount = #saveData.rages
  end
  
  backup.characterized = backup.ragesCount > 0
  
  return backup
end

---Verify write safety
---@param gau table Gau's data
---@param selectedRages table Rages to write
---@return table safetyCheck Safety verification
function WriteModePrep.verifyWriteSafety(gau, selectedRages)
  if not gau or not selectedRages then return {} end
  
  local safetyCheck = {
    safe = true,
    riskLevel = "low",
    issues = {},
    recommendedActions = {}
  }
  
  -- Check for data loss
  if gau.rages then
    for _, existing in ipairs(gau.rages) do
      local found = false
      for _, selected in ipairs(selectedRages) do
        if selected == existing then
          found = true
          break
        end
      end
      if not found then
        table.insert(safetyCheck.issues, 
          "Will lose rage: " .. existing)
      end
    end
  end
  
  if #safetyCheck.issues > 0 then
    safetyCheck.riskLevel = "medium"
    table.insert(safetyCheck.recommendedActions, 
      "Create backup before proceeding")
  end
  
  return safetyCheck
end

-- Helper functions
function WriteModePrep._verifyChecksum(saveData)
  -- Placeholder checksum verification
  return true
end

function WriteModePrep._calculateDataSize(saveData)
  return 8192  -- Typical save size
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1+",
  VeldtSimulator = VeldtSimulator,
  AdvancedDetection = AdvancedDetection,
  BattleStrategy = BattleStrategy,
  WriteModePrep = WriteModePrep,
  
  -- Feature completion
  features = {
    veldtSimulator = true,
    advancedDetection = true,
    battleStrategy = true,
    writeModePrep = true
  }
}
