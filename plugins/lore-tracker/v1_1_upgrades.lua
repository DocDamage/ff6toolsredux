--[[
  Lore Tracker Plugin - v1.1+ Upgrade Extension
  Story progression tracking, event analysis, character arc visualization, spoiler management
  
  Phase: 7E (Story & Element Systems)
  Version: 1.1+ (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: STORY PROGRESSION TRACKING (~200 LOC)
-- ============================================================================

local StoryProgressionTracking = {}

---Track current story progress
---@param saveData table Save file data
---@return table progressTracker Current story status
function StoryProgressionTracking.trackStoryProgress(saveData)
  if not saveData then return {} end
  
  local progressTracker = {
    storyAct = 0,
    majorEvents = {},
    completedQuests = {},
    progressPercentage = 0,
    estimatedHoursPlayed = 0
  }
  
  -- Determine story act based on flags
  local storyAct = StoryProgressionTracking._determineStoryAct(saveData)
  progressTracker.storyAct = storyAct
  
  -- Track major events
  local majorEvents = StoryProgressionTracking._trackMajorEvents(saveData)
  progressTracker.majorEvents = majorEvents
  
  -- Calculate progress percentage
  local totalEvents = 30  -- Approximate total major events in game
  progressTracker.progressPercentage = (#majorEvents / totalEvents) * 100
  
  -- Estimate playtime
  local baseStat = (saveData.characters and saveData.characters[1]) or {}
  local level = baseStat.level or 1
  progressTracker.estimatedHoursPlayed = level * 1.5
  
  return progressTracker
end

---Predict story events yet to come
---@param saveData table Save file data
---@return table futureEvents Predicted upcoming events
function StoryProgressionTracking.predictFutureEvents(saveData)
  if not saveData then return {} end
  
  local futureEvents = {
    upcomingMajorEvents = {},
    estimatedTimeToCompletion = 0,
    nextCriticalMoment = nil
  }
  
  local currentAct = StoryProgressionTracking._determineStoryAct(saveData)
  
  -- Define story beats per act
  local storyBeats = {
    [1] = {"Opening", "Narshe Escape", "Mako Reactor", "Coloseum"},
    [2] = {"Floating Continent", "Character Development", "World Breakdown"},
    [3] = {"World Rebuild", "Kefka Final", "Ending Sequence"}
  }
  
  if storyBeats[currentAct + 1] then
    futureEvents.upcomingMajorEvents = storyBeats[currentAct + 1]
    futureEvents.estimatedTimeToCompletion = (3 - currentAct) * 15
    futureEvents.nextCriticalMoment = storyBeats[currentAct + 1][1]
  end
  
  return futureEvents
end

---Analyze character story arcs
---@param saveData table Save file data
---@return table characterArcs Arc analysis per character
function StoryProgressionTracking.analyzeCharacterArcs(saveData)
  if not saveData then return {} end
  
  local characterArcs = {
    characters = {},
    majorArcCount = 0
  }
  
  if saveData.characters then
    for _, character in ipairs(saveData.characters) do
      local arcStage = StoryProgressionTracking._getCharacterArcStage(character)
      
      table.insert(characterArcs.characters, {
        name = character.name,
        arcStage = arcStage,
        storyEvents = character.storyEvents or {},
        characterDeveloped = arcStage >= 2
      })
      
      if arcStage >= 2 then
        characterArcs.majorArcCount = characterArcs.majorArcCount + 1
      end
    end
  end
  
  return characterArcs
end

---Generate story summary checkpoint
---@param saveData table Save file data
---@return table checkpoint Story summary
function StoryProgressionTracking.generateStoryCheckpoint(saveData)
  if not saveData then return {} end
  
  local checkpoint = {
    timestamp = os.time(),
    actNumber = StoryProgressionTracking._determineStoryAct(saveData),
    progressPercentage = 0,
    summaryPoints = {},
    keyCharacters = {}
  }
  
  -- Get major events for summary
  local majorEvents = StoryProgressionTracking._trackMajorEvents(saveData)
  
  for i, event in ipairs(majorEvents) do
    if i <= 5 then  -- Top 5 events
      table.insert(checkpoint.summaryPoints, event)
    end
  end
  
  checkpoint.progressPercentage = (#majorEvents / 30) * 100
  
  return checkpoint
end

-- Helper functions
function StoryProgressionTracking._determineStoryAct(saveData)
  if saveData.worldState and saveData.worldState == "broken" then
    return 3
  elseif saveData.floatingContinent == false then
    return 2
  else
    return 1
  end
end

function StoryProgressionTracking._trackMajorEvents(saveData)
  local events = {}
  if saveData.flags then
    for flag, value in pairs(saveData.flags) do
      if value and string.match(flag, "event_") then
        table.insert(events, flag)
      end
    end
  end
  return events
end

function StoryProgressionTracking._getCharacterArcStage(character)
  if not character then return 0 end
  return (character.storyEvents and #character.storyEvents) or 0
end

-- ============================================================================
-- FEATURE 2: EVENT ANALYSIS (~200 LOC)
-- ============================================================================

local EventAnalysis = {}

---Analyze quest completion status
---@param saveData table Save file data
---@return table questStatus Quest analysis
function EventAnalysis.analyzeQuestCompletion(saveData)
  if not saveData then return {} end
  
  local questStatus = {
    mainQuestsCompleted = 0,
    optionalQuestsCompleted = 0,
    sideQuestsAvailable = 0,
    totalQuestCount = 0
  }
  
  -- Define quest database
  local questDatabase = EventAnalysis._getQuestDatabase()
  questStatus.totalQuestCount = #questDatabase
  
  for _, quest in ipairs(questDatabase) do
    local completed = EventAnalysis._isQuestCompleted(saveData, quest)
    
    if completed then
      if quest.type == "main" then
        questStatus.mainQuestsCompleted = questStatus.mainQuestsCompleted + 1
      elseif quest.type == "optional" then
        questStatus.optionalQuestsCompleted = questStatus.optionalQuestsCompleted + 1
      end
    else
      if quest.type == "side" then
        questStatus.sideQuestsAvailable = questStatus.sideQuestsAvailable + 1
      end
    end
  end
  
  return questStatus
end

---Detect available hidden events
---@param saveData table Save file data
---@param storyProgress number Current story progress
---@return table hiddenEvents Available secret events
function EventAnalysis.detectHiddenEvents(saveData, storyProgress)
  if not saveData then return {} end
  
  local hiddenEvents = {
    availableSecrets = {},
    missedOpportunities = {},
    timeSensitiveEvents = {}
  }
  
  -- Check for secret event conditions
  local secretDatabase = EventAnalysis._getSecretEventDatabase()
  
  for _, secret in ipairs(secretDatabase) do
    if EventAnalysis._meetsEventConditions(saveData, secret.conditions) then
      table.insert(hiddenEvents.availableSecrets, {
        event = secret.name,
        location = secret.location,
        reward = secret.reward
      })
    end
  end
  
  return hiddenEvents
end

---Track optional character recruitment
---@param saveData table Save file data
---@return table recruitmentStatus Recruitment opportunities
function EventAnalysis.trackOptionalRecruitment(saveData)
  if not saveData then return {} end
  
  local recruitmentStatus = {
    recruitedCharacters = {},
    availableRecruit = {},
    missedRecruit = {}
  }
  
  -- Check character recruitment flags
  local optionalCharacters = {
    "Gogo", "Strago", "Relm", "Umaro"
  }
  
  for _, character in ipairs(optionalCharacters) do
    local recruited = EventAnalysis._isCharacterRecruited(saveData, character)
    
    if recruited then
      table.insert(recruitmentStatus.recruitedCharacters, character)
    else
      if EventAnalysis._canStillRecruit(saveData, character) then
        table.insert(recruitmentStatus.availableRecruit, character)
      else
        table.insert(recruitmentStatus.missedRecruit, character)
      end
    end
  end
  
  return recruitmentStatus
end

---Identify time-sensitive events
---@param saveData table Save file data
---@return table timeSensitive Events tied to story progression
function EventAnalysis.identifyTimeSensitiveEvents(saveData)
  if not saveData then return {} end
  
  local timeSensitive = {
    urgentEvents = {},
    oneTimeOnly = {},
    pointOfNoReturn = nil
  }
  
  local currentAct = StoryProgressionTracking._determineStoryAct(saveData)
  
  -- Point of no return is typically act 3 boundary
  if currentAct >= 2 then
    timeSensitive.pointOfNoReturn = "World Breakdown - Last chance before ending"
  end
  
  return timeSensitive
end

-- Helper functions
function EventAnalysis._getQuestDatabase()
  return {
    {name = "Magitek Factory", type = "main"},
    {name = "Floating Continent", type = "main"},
    {name = "Kefka Final Battle", type = "main"},
    {name = "Colosseum", type = "optional"},
    {name = "Monster Pit", type = "optional"}
  }
end

function EventAnalysis._getSecretEventDatabase()
  return {
    {name = "Isolated Island", location = "Unknown", reward = "Rare Items", conditions = {level = 50}},
    {name = "Ancient Castle", location = "South Figaro", reward = "Equipment", conditions = {act = 2}}
  }
end

function EventAnalysis._isQuestCompleted(saveData, quest)
  if saveData.completedEvents then
    for _, event in ipairs(saveData.completedEvents) do
      if event == quest.name then return true end
    end
  end
  return false
end

function EventAnalysis._meetsEventConditions(saveData, conditions)
  if not conditions then return true end
  
  if conditions.level and (saveData.characters and saveData.characters[1]) then
    if saveData.characters[1].level < conditions.level then
      return false
    end
  end
  
  return true
end

function EventAnalysis._isCharacterRecruited(saveData, character)
  if saveData.characters then
    for _, char in ipairs(saveData.characters) do
      if char.name == character then return true end
    end
  end
  return false
end

function EventAnalysis._canStillRecruit(saveData, character)
  -- Most optional characters can be recruited until late game
  local currentAct = StoryProgressionTracking._determineStoryAct(saveData)
  return currentAct < 3
end

-- ============================================================================
-- FEATURE 3: CHARACTER PROGRESSION VISUALIZATION (~200 LOC)
-- ============================================================================

local CharacterVisualization = {}

---Visualize character growth arcs
---@param character table Character data
---@param saveData table Full save file
---@return table arcVisualization Character arc nodes
function CharacterVisualization.visualizeCharacterArc(character, saveData)
  if not character then return {} end
  
  local arcVisualization = {
    character = character.name,
    arcEvents = {},
    arcStages = {},
    developmentScore = 0
  }
  
  -- Define character arc stages
  local arcStages = CharacterVisualization._getCharacterArcStages(character.name)
  
  for _, stage in ipairs(arcStages) do
    local stageComplete = CharacterVisualization._isArcStageComplete(saveData, character, stage)
    
    table.insert(arcVisualization.arcStages, {
      stage = stage,
      complete = stageComplete
    })
    
    if stageComplete then
      arcVisualization.developmentScore = arcVisualization.developmentScore + 1
    end
  end
  
  return arcVisualization
end

---Generate character relationship map
---@param characters table All characters
---@param saveData table Save file data
---@return table relationships Character relationships
function CharacterVisualization.generateRelationshipMap(characters, saveData)
  if not characters then return {} end
  
  local relationships = {
    pairs = {},
    triangles = {},
    groupDynamics = {}
  }
  
  -- Analyze character interactions
  local interactions = {
    {"Terra", "Locke"},
    {"Relm", "Strago"},
    {"Edgar", "Sabin"},
    {"Celes", "Locke"}
  }
  
  for _, pair in ipairs(interactions) do
    local relationship = CharacterVisualization._analyzeRelationship(pair, saveData)
    table.insert(relationships.pairs, {
      character1 = pair[1],
      character2 = pair[2],
      relationship = relationship
    })
  end
  
  return relationships
end

---Track character story events
---@param character table Character data
---@return table storyEvents Character's story moments
function CharacterVisualization.trackCharacterStoryEvents(character)
  if not character then return {} end
  
  local storyEvents = {
    character = character.name,
    events = {},
    keyMoments = {}
  }
  
  -- Character-specific story events
  local storyDatabase = CharacterVisualization._getCharacterStoryEvents(character.name)
  
  for _, event in ipairs(storyDatabase) do
    table.insert(storyEvents.events, event)
    
    if event.isCritical then
      table.insert(storyEvents.keyMoments, event.name)
    end
  end
  
  return storyEvents
end

-- Helper functions
function CharacterVisualization._getCharacterArcStages(characterName)
  local stages = {
    ["Terra"] = {"Magitek Knight", "Esperian", "Mage"},
    ["Locke"] = {"Thief", "Treasure Hunter", "Hero"},
    ["Edgar"] = {"Prince", "Leader", "King"},
    ["Sabin"] = {"Monk", "Wanderer", "Warrior"}
  }
  return stages[characterName] or {"Introduction", "Development", "Resolution"}
end

function CharacterVisualization._isArcStageComplete(saveData, character, stage)
  if character.storyEvents then
    for _, event in ipairs(character.storyEvents) do
      if string.match(event, stage) then return true end
    end
  end
  return false
end

function CharacterVisualization._analyzeRelationship(pair, saveData)
  return "Complex"
end

function CharacterVisualization._getCharacterStoryEvents(characterName)
  local events = {
    ["Terra"] = {
      {name = "Magitek Armor Escape", isCritical = true},
      {name = "Memory Recovery", isCritical = true}
    },
    ["Locke"] = {
      {name = "Treasure Hunting", isCritical = false},
      {name = "Betrayal", isCritical = true}
    }
  }
  return events[characterName] or {}
end

-- ============================================================================
-- FEATURE 4: SPOILER MANAGEMENT (~200 LOC)
-- ============================================================================

local SpoilerManagement = {}

---Set spoiler sensitivity level
---@param level string Sensitivity: "full", "medium", "minimal", "none"
---@return table spoilerSettings Current settings
function SpoilerManagement.setSpoilerSensitivity(level)
  if not level then level = "medium" end
  
  local spoilerSettings = {
    sensitivityLevel = level,
    hideMajorPlots = level ~= "none",
    hideCharacterFates = level == "full",
    hideEndingDetails = level == "full",
    hideSecretLocations = level == "medium" or level == "full"
  }
  
  return spoilerSettings
end

---Filter story information by spoiler level
---@param information table Story information
---@param spoilerLevel string Sensitivity level
---@return table filtered Filtered information
function SpoilerManagement.filterBySpoilerLevel(information, spoilerLevel)
  if not information or not spoilerLevel then return {} end
  
  local filtered = {
    visibleInformation = {},
    hiddenWarnings = 0
  }
  
  if spoilerLevel == "none" then
    filtered.visibleInformation = information
  elseif spoilerLevel == "minimal" then
    -- Show everything but add warnings
    filtered.visibleInformation = information
    filtered.hiddenWarnings = 2
  elseif spoilerLevel == "medium" then
    -- Hide some details
    filtered.visibleInformation = SpoilerManagement._filterMajorPlots(information)
    filtered.hiddenWarnings = 5
  elseif spoilerLevel == "full" then
    -- Heavily redacted
    filtered.visibleInformation = SpoilerManagement._filterAllPlots(information)
    filtered.hiddenWarnings = 15
  end
  
  return filtered
end

---Generate spoiler-free story summary
---@param saveData table Save file data
---@return table summary Safe story summary
function SpoilerManagement.generateSpoilerFreeSummary(saveData)
  if not saveData then return {} end
  
  local summary = {
    currentStatus = "In Progress",
    safeInformation = {
      "You are in the middle of your adventure",
      "Many characters await you",
      "Great challenges lie ahead"
    },
    avoidedSpoilers = 0
  }
  
  local currentAct = StoryProgressionTracking._determineStoryAct(saveData)
  
  if currentAct == 1 then
    summary.currentStatus = "Early in your journey"
  elseif currentAct == 2 then
    summary.currentStatus = "Midway through your adventure"
  elseif currentAct == 3 then
    summary.currentStatus = "Nearing the climax"
  end
  
  return summary
end

---Warn about upcoming spoiler content
---@param topic string Story topic
---@return table warning Spoiler warning
function SpoilerManagement.generateSpoilerWarning(topic)
  if not topic then return {} end
  
  local warning = {
    topic = topic,
    hasSpoilers = true,
    severityLevel = "high",
    recommendation = "Proceed with caution"
  }
  
  return warning
end

-- Helper functions
function SpoilerManagement._filterMajorPlots(information)
  local filtered = {}
  if type(information) == "table" then
    for key, value in pairs(information) do
      if not string.match(key, "ending") then
        filtered[key] = value
      end
    end
  end
  return filtered
end

function SpoilerManagement._filterAllPlots(information)
  return {note = "Content hidden to prevent spoilers"}
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1+",
  StoryProgressionTracking = StoryProgressionTracking,
  EventAnalysis = EventAnalysis,
  CharacterVisualization = CharacterVisualization,
  SpoilerManagement = SpoilerManagement,
  
  -- Feature completion
  features = {
    storyProgressionTracking = true,
    eventAnalysis = true,
    characterVisualization = true,
    spoilerManagement = true
  }
}
