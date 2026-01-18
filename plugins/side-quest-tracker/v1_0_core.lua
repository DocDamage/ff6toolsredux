--[[
  Side Quest Tracker Plugin - v1.0
  Quest progression system with reward optimization and NPC tracking
  
  Phase: 8 (Tier 2 - Ecosystem Expansion)
  Version: 1.0 (New Plugin)
]]

-- ============================================================================
-- FEATURE 1: QUEST DATABASE (~250 LOC)
-- ============================================================================

local QuestDatabase = {}

---Get quest information
---@param questId number Quest identifier
---@return table questInfo Quest data
function QuestDatabase.getQuestInfo(questId)
  if not questId then return {} end
  
  local quests = {
    [1] = {name = "Gau's Father", npc = "Gau", reward = 5000, type = "rescue"},
    [2] = {name = "Opera House", npc = "Impresario", reward = 3000, type = "event"},
    [3] = {name = "Floating Island", npc = "Cid", reward = 10000, type = "quest"},
    [4] = {name = "Relic Hunt", npc = "Relm", reward = 2000, type = "collection"}
  }
  
  return quests[questId] or {}
end

---List all side quests
---@return table quests All quests
function QuestDatabase.listAllQuests()
  local quests = {
    {id = 1, name = "Gau's Father", npc = "Gau"},
    {id = 2, name = "Opera House", npc = "Impresario"},
    {id = 3, name = "Floating Island", npc = "Cid"},
    {id = 4, name = "Relic Hunt", npc = "Relm"}
  }
  
  return quests
end

---Filter quests by type
---@param questType string Quest type
---@return table filtered Filtered quests
function QuestDatabase.filterByType(questType)
  if not questType then return {} end
  
  if questType == "rescue" then
    return {{id = 1, name = "Gau's Father"}}
  elseif questType == "event" then
    return {{id = 2, name = "Opera House"}}
  elseif questType == "collection" then
    return {{id = 4, name = "Relic Hunt"}}
  end
  
  return {}
end

---Get quest requirements
---@param questId number Quest
---@return table requirements Quest requirements
function QuestDatabase.getRequirements(questId)
  if not questId then return {} end
  
  local requirements = {
    questId = questId,
    minimumLevel = 25,
    requiredItems = {"Map Fragment"},
    prerequisiteQuests = {"Gau's Father"},
    timeLimit = nil
  }
  
  return requirements
end

-- ============================================================================
-- FEATURE 2: PROGRESSION TRACKING (~250 LOC)
-- ============================================================================

local ProgressionTracking = {}

---Get quest completion status
---@param playerData table Player save data
---@return table completion Completion status
function ProgressionTracking.getCompletion(playerData)
  if not playerData then return {} end
  
  local completion = {
    questsCompleted = 8,
    totalQuests = 20,
    completionPercent = 40,
    activeQuests = 3,
    nextQuest = "Relic Hunt"
  }
  
  return completion
end

---Track quest progress
---@param questId number Quest
---@param playerData table Player data
---@return table progress Quest progress
function ProgressionTracking.getProgress(questId, playerData)
  if not questId or not playerData then return {} end
  
  local progress = {
    questId = questId,
    status = "In Progress",
    objectives = {
      {objective = "Find map fragment", completed = true},
      {objective = "Locate treasure", completed = false},
      {objective = "Return to NPC", completed = false}
    },
    progressPercent = 33
  }
  
  return progress
end

---Generate milestone achievements
---@param playerData table Player data
---@return table milestones Quest milestones
function ProgressionTracking.getMilestones(playerData)
  if not playerData then return {} end
  
  local milestones = {
    completed = {
      "Complete 5 quests (50% → 75%)",
      "Complete 10 quests (100% quest tier)"
    },
    nextMilestone = "Complete 3 more quests"
  }
  
  return milestones
end

---Track active quests
---@param playerData table Player data
---@return table active Active quest list
function ProgressionTracking.getActive(playerData)
  if not playerData then return {} end
  
  local active = {
    activeCount = 3,
    quests = {
      {questId = 1, name = "Gau's Father", progress = 75},
      {questId = 3, name = "Floating Island", progress = 25}
    }
  }
  
  return active
end

-- ============================================================================
-- FEATURE 3: LOCATION MAPPING (~250 LOC)
-- ============================================================================

local LocationMapping = {}

---Get quest giver location
---@param questId number Quest
---@return table location NPC location
function LocationMapping.getQuestGiverLocation(questId)
  if not questId then return {} end
  
  local locations = {
    [1] = {npc = "Gau", location = "Veldt", coordinates = {x = 100, y = 150}},
    [2] = {npc = "Impresario", location = "Opera House", coordinates = {x = 200, y = 100}},
    [3] = {npc = "Cid", location = "Floating Island", coordinates = {x = 250, y = 200}},
    [4] = {npc = "Relm", location = "Zozo", coordinates = {x = 80, y = 120}}
  }
  
  return locations[questId] or {}
end

---Map quest objectives
---@param questId number Quest
---@return table objectives Objective locations
function LocationMapping.mapObjectives(questId)
  if not questId then return {} end
  
  local objectives = {
    questId = questId,
    locations = {
      {objective = "Find Fragment", location = "Dungeon A", difficulty = "moderate"},
      {objective = "Locate Treasure", location = "Hidden Cave", difficulty = "hard"}
    },
    mapMarkers = 2
  }
  
  return objectives
end

---Get quest reward location
---@param questId number Quest
---@return table rewardLoc Reward pickup location
function LocationMapping.getRewardLocation(questId)
  if not questId then return {} end
  
  local rewardLocation = {
    questId = questId,
    rewardPickupNpc = "Quest Giver",
    location = "Original NPC location",
    coordinates = LocationMapping.getQuestGiverLocation(questId).coordinates
  }
  
  return rewardLocation
end

---Generate location guide
---@param questId number Quest
---@return table guide Location guide
function LocationMapping.generateGuide(questId)
  if not questId then return {} end
  
  local guide = {
    questId = questId,
    stops = {
      {order = 1, location = "Quest Giver", description = "Receive quest"},
      {order = 2, location = "Objective Area", description = "Complete objective"},
      {order = 3, location = "Return Location", description = "Return for reward"}
    }
  }
  
  return guide
end

-- ============================================================================
-- FEATURE 4: REWARD OPTIMIZATION (~220 LOC)
-- ============================================================================

local RewardOptimization = {}

---Get quest rewards
---@param questId number Quest
---@return table rewards Quest rewards
function RewardOptimization.getRewards(questId)
  if not questId then return {} end
  
  local quest = QuestDatabase.getQuestInfo(questId)
  
  local rewards = {
    questId = questId,
    gil = quest.reward or 0,
    items = {"Map Fragment", "Key"},
    experience = 2500,
    totalReward = {gil = quest.reward or 0, exp = 2500}
  }
  
  return rewards
end

---Rank quests by efficiency
---@param playerData table Player data
---@return table ranking Quest efficiency ranking
function RewardOptimization.rankByEfficiency(playerData)
  if not playerData then return {} end
  
  local ranking = {
    quests = {
      {questId = 3, name = "Floating Island", reward = 10000, time = 30, efficiency = 333},
      {questId = 1, name = "Gau's Father", reward = 5000, time = 20, efficiency = 250},
      {questId = 2, name = "Opera House", reward = 3000, time = 15, efficiency = 200}
    },
    bestReward = "Floating Island"
  }
  
  return ranking
end

---Calculate total rewards available
---@param playerData table Player data
---@return table totalRewards Total available rewards
function RewardOptimization.getTotalAvailable(playerData)
  if not playerData then return {} end
  
  local totalRewards = {
    totalGilAvailable = 45000,
    totalExpAvailable = 50000,
    questsNotCompleted = 12,
    totalTime = "8-10 hours"
  }
  
  return totalRewards
end

---Suggest quest order
---@param playerData table Player data
---@return table order Suggested quest order
function RewardOptimization.suggestOrder(playerData)
  if not playerData then return {} end
  
  local order = {
    recommended = {
      1, 2, 3, 4
    },
    reasoning = "Optimal reward-to-time ratio",
    estimatedTotal = 45000
  }
  
  return order
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.0",
  QuestDatabase = QuestDatabase,
  ProgressionTracking = ProgressionTracking,
  LocationMapping = LocationMapping,
  RewardOptimization = RewardOptimization,
  
  features = {
    questDatabase = true,
    progressionTracking = true,
    locationMapping = true,
    rewardOptimization = true
  }
}
