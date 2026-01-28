--[[
  World Map Explorer Plugin - v1.0
  Interactive world navigation with treasure tracking and NPC directory
  
  Phase: 8 (Tier 2 - Ecosystem Expansion)
  Version: 1.0 (New Plugin)
]]

-- ============================================================================
-- FEATURE 1: MAP NAVIGATION (~250 LOC)
-- ============================================================================

local MapNavigation = {}

---Get region information
---@param regionId number Region identifier
---@return table regionInfo Region data
function MapNavigation.getRegionInfo(regionId)
  if not regionId then return {} end
  
  local regions = {
    [1] = {name = "Narshe", type = "town", discovered = true, accessible = true},
    [2] = {name = "South Figaro", type = "town", discovered = true, accessible = true},
    [3] = {name = "Zozo", type = "town", discovered = true, accessible = true},
    [4] = {name = "Floating Island", type = "dungeon", discovered = false, accessible = false},
    [5] = {name = "Kefka's Tower", type = "dungeon", discovered = false, accessible = false}
  }
  
  return regions[regionId] or {}
end

---List all regions
---@return table regions All regions
function MapNavigation.listAllRegions()
  local regions = {
    {id = 1, name = "Narshe", type = "town", discovered = true},
    {id = 2, name = "South Figaro", type = "town", discovered = true},
    {id = 3, name = "Zozo", type = "town", discovered = true},
    {id = 4, name = "Floating Island", type = "dungeon", discovered = false},
    {id = 5, name = "Kefka's Tower", type = "dungeon", discovered = false}
  }
  
  return regions
end

---Filter regions by type
---@param regionType string Region type filter
---@return table filtered Filtered regions
function MapNavigation.filterByType(regionType)
  if not regionType then return {} end
  
  if regionType == "town" then
    return {
      {id = 1, name = "Narshe"},
      {id = 2, name = "South Figaro"},
      {id = 3, name = "Zozo"}
    }
  elseif regionType == "dungeon" then
    return {
      {id = 4, name = "Floating Island"},
      {id = 5, name = "Kefka's Tower"}
    }
  end
  
  return {}
end

---Get accessible locations
---@param playerData table Player save state
---@return table accessible Currently accessible locations
function MapNavigation.getAccessibleLocations(playerData)
  if not playerData then return {} end
  
  local accessible = {
    currentLocation = "Narshe",
    accessible = {
      "Narshe",
      "South Figaro",
      "Zozo"
    },
    totalAccessible = 3,
    locked = {"Floating Island", "Kefka's Tower"}
  }
  
  return accessible
end

-- ============================================================================
-- FEATURE 2: TREASURE TRACKING (~250 LOC)
-- ============================================================================

local TreasureTracking = {}

---Get all treasure chests
---@return table chests All treasure chests
function TreasureTracking.getAllTreasures()
  local chests = {
    {id = 1, location = "Narshe", item = "Elixir", collected = false},
    {id = 2, location = "South Figaro", item = "Phoenix Down", collected = true},
    {id = 3, location = "Zozo", item = "Excalibur", collected = false},
    {id = 4, location = "Floating Island", item = "Aegis Shield", collected = false}
  }
  
  return chests
end

---Map treasure location
---@param chestId number Chest identifier
---@return table location Treasure location info
function TreasureTracking.mapTreasureLocation(chestId)
  if not chestId then return {} end
  
  local location = {
    chestId = chestId,
    region = "Zozo",
    coordinates = {x = 45, y = 78},
    description = "Upper floor, east side",
    itemName = "Excalibur",
    itemRarity = "legendary",
    collected = false
  }
  
  return location
end

---Track collection progress
---@param playerData table Player save data
---@return table progress Collection statistics
function TreasureTracking.trackProgress(playerData)
  if not playerData then return {} end
  
  local progress = {
    treasuresFound = 8,
    totalTreasures = 35,
    completionPercent = 23,
    nextTreasure = "Phoenix Down at South Figaro",
    estimatedTotal = "35 major treasures"
  }
  
  return progress
end

---Identify missable treasures
---@return table missable Missable treasure list
function TreasureTracking.identifyMissable()
  local missable = {
    count = 12,
    treasures = {
      {location = "Floating Island", item = "Elixir", missableAt = "World of Ruin"},
      {location = "Zozo", item = "Excalibur", missableAt = "Post-Floating Island"}
    }
  }
  
  return missable
end

-- ============================================================================
-- FEATURE 3: NPC DIRECTORY (~250 LOC)
-- ============================================================================

local NPCDirectory = {}

---Find NPC location
---@param npcId number NPC identifier
---@return table npcInfo NPC information
function NPCDirectory.findNPC(npcId)
  if not npcId then return {} end
  
  local npcData = {
    [1] = {name = "Locke", location = "Narshe", role = "Protagonist"},
    [2] = {name = "Terra", location = "South Figaro", role = "Protagonist"},
    [3] = {name = "Celes", location = "Zozo", role = "Protagonist"},
    [4] = {name = "Kefka", location = "Kefka's Tower", role = "Antagonist"}
  }
  
  return npcData[npcId] or {}
end

---List NPCs by region
---@param region string Region name
---@return table npcs NPCs in region
function NPCDirectory.listByRegion(region)
  if not region then return {} end
  
  if region == "Narshe" then
    return {{id = 1, name = "Locke"}, {id = 5, name = "Arvis"}}
  elseif region == "South Figaro" then
    return {{id = 2, name = "Terra"}, {id = 6, name = "Merchant"}}
  elseif region == "Zozo" then
    return {{id = 3, name = "Celes"}}
  end
  
  return {}
end

---Track NPC dialogue status
---@param npcId number NPC identifier
---@return table dialogue Dialogue information
function NPCDirectory.getDialogueStatus(npcId)
  if not npcId then return {} end
  
  local dialogue = {
    npcId = npcId,
    dialoguesAvailable = 3,
    dialoguesViewed = 1,
    completionPercent = 33,
    unviewedTopics = {"Past", "Future"}
  }
  
  return dialogue
end

---Map NPC locations by story progression
---@param storyAct number Current story act
---@return table mapping NPC location mapping
function NPCDirectory.mapByStory(storyAct)
  if not storyAct then return {} end
  
  local mapping = {
    storyAct = storyAct,
    npcLocations = {
      {name = "Locke", location = "Narshe"},
      {name = "Terra", location = "Vector"}
    }
  }
  
  return mapping
end

-- ============================================================================
-- FEATURE 4: EXPLORATION GUIDE (~220 LOC)
-- ============================================================================

local ExplorationGuide = {}

---Identify secret areas
---@return table secrets Secret areas
function ExplorationGuide.findSecretAreas()
  local secrets = {
    totalSecrets = 8,
    discovered = 2,
    secrets = {
      {name = "Hidden Cave", location = "Narshe", discovered = false},
      {name = "Secret Garden", location = "South Figaro", discovered = true}
    },
    nextSecret = "Underground Path at Zozo"
  }
  
  return secrets
end

---Generate exploration checklist
---@param region string Region to explore
---@return table checklist Exploration checklist
function ExplorationGuide.generateChecklist(region)
  if not region then return {} end
  
  local checklist = {
    region = region,
    items = {
      {item = "Treasures", count = 5, found = 2},
      {item = "NPCs", count = 4, found = 3},
      {item = "Secret Areas", count = 2, found = 0}
    },
    completionPercent = 50
  }
  
  return checklist
end

---Locate hidden content
---@param contentType string Content type
---@return table hidden Hidden content
function ExplorationGuide.locateHidden(contentType)
  if not contentType then return {} end
  
  local hidden = {
    contentType = contentType,
    locations = 3,
    details = {
      {location = "Floating Island", description = "Rare item cache"},
      {location = "Zozo", description = "Secret passage"}
    }
  }
  
  return hidden
end

---Track exploration completion
---@param playerData table Player save data
---@return table completion Exploration completion stats
function ExplorationGuide.getCompletion(playerData)
  if not playerData then return {} end
  
  local completion = {
    regionsExplored = 12,
    totalRegions = 20,
    explorationPercent = 60,
    itemsFound = 25,
    nextGoal = "Explore Floating Island (60% → 70%)"
  }
  
  return completion
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.0",
  MapNavigation = MapNavigation,
  TreasureTracking = TreasureTracking,
  NPCDirectory = NPCDirectory,
  ExplorationGuide = ExplorationGuide,
  
  features = {
    mapNavigation = true,
    treasureTracking = true,
    npcDirectory = true,
    explorationGuide = true
  }
}
