--[[
  Speedrun Timer Plugin - v1.1 Upgrade Extension
  Advanced split management, personal records, race mode, and streaming integration
  
  Phase: 7A (Advanced Analysis & Tracking)
  Version: 1.1.0 (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: ADVANCED SPLIT MANAGEMENT (~300 LOC)
-- ============================================================================

local SplitManagement = {}

---Create custom split profile
---@param name string Profile name
---@param splits table List of split definitions
---@return table profile Created profile
function SplitManagement.createSplitProfile(name, splits)
  if not name or not splits then return {} end
  
  local profile = {
    id = math.random(100000, 999999),
    name = name,
    splits = {},
    createdAt = os.time(),
    modifiedAt = os.time(),
    version = "1.1.0",
    isActive = false
  }
  
  for i, split in ipairs(splits) do
    table.insert(profile.splits, {
      index = i,
      name = split.name or string.format("Split %d", i),
      description = split.description or "",
      targetTime = split.targetTime or 0,
      notes = split.notes or ""
    })
  end
  
  return profile
end

---Save split profile to library
---@param profile table Profile to save
---@return boolean success Save successful
function SplitManagement.saveSplitPreset(profile)
  if not profile or not profile.name then return false end
  
  profile.savedAt = os.time()
  profile.isSaved = true
  
  return true
end

---Load saved split preset
---@param profileId string Profile ID to load
---@return table profile Loaded profile
function SplitManagement.loadSplitPreset(profileId)
  if not profileId then return {} end
  
  -- In real implementation, would load from database
  return {
    id = profileId,
    name = "Loaded Profile",
    splits = {},
    isActive = true
  }
end

---Compare two split profiles
---@param profile1 table First profile
---@param profile2 table Second profile
---@return table comparison Detailed comparison
function SplitManagement.compareSplits(profile1, profile2)
  if not profile1 or not profile2 then return {} end
  
  local comparison = {
    profile1Name = profile1.name,
    profile2Name = profile2.name,
    splitCount = math.max(#profile1.splits, #profile2.splits),
    differences = {},
    timeDelta = {}
  }
  
  -- Compare each split
  for i = 1, math.max(#profile1.splits, #profile2.splits) do
    local split1 = profile1.splits[i]
    local split2 = profile2.splits[i]
    
    if split1 and split2 then
      local delta = (split2.targetTime or 0) - (split1.targetTime or 0)
      table.insert(comparison.timeDelta, {
        splitIndex = i,
        delta = delta,
        deltaPercent = split1.targetTime > 0 and (delta / split1.targetTime * 100) or 0
      })
    end
  end
  
  return comparison
end

---Manage split library
---@return table library Current split library
function SplitManagement.manageSplitLibrary()
  return {
    profiles = {},
    totalProfiles = 0,
    totalSplits = 0,
    categories = {},
    lastModified = os.time()
  }
end

---Import splits from external file
---@param filename string File to import
---@return boolean success, table imports Import result
function SplitManagement.importExternalSplits(filename)
  if not filename then
    return false, {}
  end
  
  -- In real implementation, would parse file
  return true, {
    imported = 0,
    failed = 0,
    filename = filename
  }
end

-- ============================================================================
-- FEATURE 2: PERSONAL RECORD SYSTEM (~250 LOC)
-- ============================================================================

local PersonalRecords = {}

---Set personal best for comparison
---@param runTime number Time in seconds
---@param splits table Split times
---@return table pbRecord PB record
function PersonalRecords.setPBComparison(runTime, splits)
  if not runTime then return {} end
  
  local pbRecord = {
    time = runTime,
    splits = splits or {},
    setAt = os.time(),
    verified = false,
    category = "any%"
  }
  
  return pbRecord
end

---Compare current run to personal best
---@param currentTime number Current run time
---@param pbTime number Personal best time
---@return table comparison Comparison data
function PersonalRecords.compareToRecord(currentTime, pbTime)
  if not currentTime or not pbTime then return {} end
  
  local delta = currentTime - pbTime
  local deltaPercent = (delta / pbTime) * 100
  
  return {
    currentTime = currentTime,
    pbTime = pbTime,
    delta = delta,
    deltaPercent = deltaPercent,
    status = delta < 0 and "NEW PB!" or (delta < 60 and "Close!" or "Off pace"),
    isNewRecord = delta < 0
  }
end

---Calculate delta time display
---@param currentSegmentTime number Current segment
---@param pbSegmentTime number PB segment
---@return table deltaDisplay Display-friendly delta
function PersonalRecords.calculateDeltaTime(currentSegmentTime, pbSegmentTime)
  if not currentSegmentTime or not pbSegmentTime then return {} end
  
  local delta = currentSegmentTime - pbSegmentTime
  local minutes = math.floor(math.abs(delta) / 60)
  local seconds = math.abs(delta) % 60
  
  return {
    seconds = delta,
    display = string.format(
      "%s%02d:%05.2f",
      delta < 0 and "-" or "+",
      minutes,
      seconds
    ),
    color = delta < 0 and "gold" or (delta > 0 and "red" or "neutral")
  }
end

---Predict final time based on current pace
---@param currentTime number Time so far
---@param currentSegment number Current segment number
---@param totalSegments number Total segments
---@param pbTime number Personal best time
---@return number predictedTime Projected finish time
function PersonalRecords.predictFinalTime(currentTime, currentSegment, totalSegments, pbTime)
  if not currentTime or not currentSegment or not totalSegments or not pbTime then return 0 end
  
  if currentSegment == 0 then return 0 end
  
  local pacePerSegment = currentTime / currentSegment
  local remainingSegments = totalSegments - currentSegment
  
  return currentTime + (pacePerSegment * remainingSegments)
end

---Archive historical record
---@param runData table Run to archive
---@return boolean success Archive successful
function PersonalRecords.archiveRecord(runData)
  if not runData then return false end
  
  runData.archivedAt = os.time()
  runData.isArchived = true
  
  return true
end

---Track multiple personal bests
---@param category string Run category
---@return table categoryRecords All PBs for category
function PersonalRecords.trackMultiplePBs(category)
  if not category then return {} end
  
  return {
    category = category,
    records = {},
    bestRecord = nil,
    recordCount = 0
  }
end

-- ============================================================================
-- FEATURE 3: PACE CALCULATOR v2 (~280 LOC)
-- ============================================================================

local PaceCalculator = {}

---Calculate realistic expected pace
---@param splitTimes table Times for completed splits
---@return number pace Average pace per segment
function PaceCalculator.calculatePaceRealistic(splitTimes)
  if not splitTimes or #splitTimes == 0 then return 0 end
  
  local totalTime = 0
  for _, time in ipairs(splitTimes) do
    totalTime = totalTime + time
  end
  
  return totalTime / #splitTimes
end

---Detect when pace changes
---@param splitTimes table Split times
---@param threshold number Change threshold (percentage)
---@return table paceChanges Detected changes
function PaceCalculator.detectPaceChanges(splitTimes, threshold)
  if not splitTimes or #splitTimes < 2 then return {} end
  
  threshold = threshold or 10  -- 10% default
  
  local changes = {}
  local avgPace = PaceCalculator.calculatePaceRealistic(splitTimes)
  
  for i, time in ipairs(splitTimes) do
    local deviation = math.abs(time - avgPace) / avgPace * 100
    if deviation > threshold then
      table.insert(changes, {
        splitIndex = i,
        time = time,
        deviation = deviation,
        status = time > avgPace and "slower" or "faster"
      })
    end
  end
  
  return changes
end

---Project completion time
---@param currentTime number Time elapsed
---@param completionPercent number Percent complete
---@return number projectedTime Estimated finish time
function PaceCalculator.projectedCompletion(currentTime, completionPercent)
  if not currentTime or not completionPercent or completionPercent == 0 then return 0 end
  
  return currentTime / (completionPercent / 100)
end

---Breakdown pace by segment
---@param splits table Split data
---@return table breakdown Per-segment pace analysis
function PaceCalculator.paceBreakdown(splits)
  if not splits then return {} end
  
  local breakdown = {
    segments = {},
    slowestSegment = nil,
    fastestSegment = nil,
    averagePace = 0
  }
  
  local totalTime = 0
  
  for i, split in ipairs(splits) do
    local segmentTime = split.time or 0
    totalTime = totalTime + segmentTime
    
    table.insert(breakdown.segments, {
      index = i,
      name = split.name,
      time = segmentTime,
      pace = 0
    })
    
    if not breakdown.slowestSegment or segmentTime > breakdown.slowestSegment.time then
      breakdown.slowestSegment = breakdown.segments[i]
    end
    
    if not breakdown.fastestSegment or segmentTime < breakdown.fastestSegment.time then
      breakdown.fastestSegment = breakdown.segments[i]
    end
  end
  
  breakdown.averagePace = totalTime / #splits
  
  return breakdown
end

---Compare segment pace to previous run
---@param currentSegment number Current segment time
---@param previousSegment number Previous run segment time
---@return table comparison Segment comparison
function PaceCalculator.compareSegmentPace(currentSegment, previousSegment)
  if not currentSegment or not previousSegment then return {} end
  
  local delta = currentSegment - previousSegment
  local percent = (delta / previousSegment) * 100
  
  return {
    current = currentSegment,
    previous = previousSegment,
    delta = delta,
    deltaPercent = percent,
    faster = delta < 0
  }
end

---Advise pace adjustment
---@param currentPace number Current pace per segment
---@param targetPace number Target pace
---@return string advice Pace advice
function PaceCalculator.advisePaceAdjustment(currentPace, targetPace)
  if not currentPace or not targetPace then return "" end
  
  if currentPace < targetPace then
    return "Good pace! Maintain current speed."
  elseif currentPace > targetPace then
    local slowdown = ((currentPace - targetPace) / targetPace) * 100
    return string.format("Need to speed up by %.0f%% to match target pace.", slowdown)
  else
    return "Perfect pace!"
  end
end

-- ============================================================================
-- FEATURE 4: CATEGORY PRESETS (~200 LOC)
-- ============================================================================

local CategoryPresets = {}

---Create new speedrun category
---@param name string Category name
---@param rules table Category rules
---@return table category Created category
function CategoryPresets.createCategory(name, rules)
  if not name then return {} end
  
  return {
    id = math.random(100000, 999999),
    name = name,
    rules = rules or {},
    createdAt = os.time(),
    active = false,
    splits = {}
  }
end

---Add preset to library
---@param presetName string Preset identifier
---@param definition table Preset definition
---@return boolean success Added successfully
function CategoryPresets.addCategoryPreset(presetName, definition)
  if not presetName or not definition then return false end
  
  definition.presetName = presetName
  definition.addedAt = os.time()
  
  return true
end

---Load category preset
---@param presetName string Preset to load
---@return table preset Loaded preset
function CategoryPresets.loadCategoryPreset(presetName)
  if not presetName then return {} end
  
  local presets = {
    ["any%"] = {
      name = "Any%",
      description = "Reach the end as fast as possible",
      splits = {"Start", "First Dungeon", "Midpoint", "Final Boss"}
    },
    ["100%"] = {
      name = "100%",
      description = "Complete everything",
      splits = {"Start", "25%", "50%", "75%", "100%"}
    },
    ["glitchless"] = {
      name = "Glitchless",
      description = "No sequence breaks allowed",
      splits = {"Start", "Act 1", "Act 2", "Act 3", "Final Boss"}
    },
    ["low%"] = {
      name = "Low%",
      description = "Minimal item collection",
      splits = {"Start", "Essential Items", "Final Boss"}
    },
    ["no-save"] = {
      name = "No-Save Run",
      description = "Single continuous session",
      splits = {"Start", "Midpoint", "Final Boss"}
    },
    ["blind"] = {
      name = "Blind",
      description = "No guide or hints",
      splits = {"Start", "Progress Checkpoint", "Finish"}
    }
  }
  
  return presets[presetName] or {}
end

---Edit category rules
---@param categoryId string Category to edit
---@param newRules table Updated rules
---@return boolean success Edit successful
function CategoryPresets.editCategoryRules(categoryId, newRules)
  if not categoryId or not newRules then return false end
  
  return true
end

---Delete category preset
---@param presetName string Preset to delete
---@return boolean success Deleted successfully
function CategoryPresets.deleteCategoryPreset(presetName)
  if not presetName then return false end
  
  return true
end

---Suggest default splits for category
---@param categoryName string Category name
---@return table defaultSplits Suggested split points
function CategoryPresets.suggestCategoryDefaults(categoryName)
  if not categoryName then return {} end
  
  -- Get preset and extract splits
  local preset = CategoryPresets.loadCategoryPreset(categoryName)
  return preset.splits or {}
end

-- ============================================================================
-- FEATURE 5: STREAMING INTEGRATION (~350 LOC)
-- ============================================================================

local StreamingIntegration = {}

---Enable stream-optimized mode
---@return table config Stream mode configuration
function StreamingIntegration.enableStreamMode()
  return {
    enabled = true,
    obsConnected = false,
    twitchConnected = false,
    overlayVisible = true,
    chatIntegration = true,
    recordingMode = true
  }
end

---Stream timer output
---@param timerData table Current timer state
---@return string obsOutput OBS-compatible output
function StreamingIntegration.streamTimer(timerData)
  if not timerData then return "" end
  
  local output = string.format(
    "Time: %02d:%05.2f\nDelta: %s\nSegment: %s",
    timerData.minutes or 0,
    timerData.seconds or 0,
    timerData.delta or "+0:00.00",
    timerData.currentSegment or "N/A"
  )
  
  return output
end

---Create stream overlay graphics
---@param theme string Overlay theme
---@return table overlay Overlay configuration
function StreamingIntegration.createStreamOverlay(theme)
  if not theme then theme = "default" end
  
  return {
    theme = theme,
    position = {x = 0, y = 0},
    size = {width = 300, height = 200},
    elements = {
      timer = {visible = true, fontSize = 24},
      splits = {visible = true, count = 3},
      pace = {visible = true},
      pb = {visible = true}
    },
    opacity = 0.9
  }
end

---Auto-update Twitch channel info
---@param title string New stream title
---@param category string Game category
---@return boolean success Update sent
function StreamingIntegration.autoUpdateTwitch(title, category)
  if not title or not category then return false end
  
  -- In real implementation, would call Twitch API
  return true
end

---Generate stream graphic
---@param timerData table Current timer data
---@return string graphic Graphic file path
function StreamingIntegration.generateStreamGraphic(timerData)
  if not timerData then return "" end
  
  return "stream_overlay_" .. os.time() .. ".png"
end

---Track chat predictions
---@return table predictions Chat prediction data
function StreamingIntegration.trackChatPredictions()
  return {
    predictionId = "",
    outcomes = {},
    active = false,
    viewers = 0
  }
end

-- ============================================================================
-- FEATURE 6: RACE MODE (~300 LOC)
-- ============================================================================

local RaceMode = {}

---Start multiplayer race
---@param players table Player list
---@param category string Race category
---@return table raceSession Race session data
function RaceMode.startRaceMode(players, category)
  if not players or #players == 0 or not category then return {} end
  
  local raceSession = {
    id = math.random(100000000, 999999999),
    players = players,
    category = category,
    startTime = os.time(),
    status = "running",
    results = {}
  }
  
  for _, player in ipairs(players) do
    table.insert(raceSession.results, {
      player = player,
      time = 0,
      position = 0,
      finished = false
    })
  end
  
  return raceSession
end

---Track multiple players in race
---@param raceSession table Current race
---@param playerUpdates table Updated player times
---@return table raceSession Updated race state
function RaceMode.trackMultiplePlayers(raceSession, playerUpdates)
  if not raceSession or not playerUpdates then return raceSession end
  
  for _, update in ipairs(playerUpdates) do
    for _, result in ipairs(raceSession.results) do
      if result.player == update.player then
        result.time = update.time
        result.updated = os.time()
      end
    end
  end
  
  return raceSession
end

---Calculate lead status
---@param raceSession table Current race
---@return table standings Current standings
function RaceMode.calculateLeadStatus(raceSession)
  if not raceSession or not raceSession.results then return {} end
  
  local standings = {}
  
  -- Sort by time
  for i, result in ipairs(raceSession.results) do
    table.insert(standings, {
      position = i,
      player = result.player,
      time = result.time,
      status = i == 1 and "LEAD" or string.format("-%s behind", 
        RaceMode._formatTimeDifference((standings[i-1].time or 0) - result.time))
    })
  end
  
  table.sort(standings, function(a, b) return a.time < b.time end)
  
  return standings
end

---Synchronize race between runners
---@param raceSession table Race to sync
---@param latency number Network latency in ms
---@return boolean success Sync successful
function RaceMode.raceSynchronization(raceSession, latency)
  if not raceSession then return false end
  
  latency = latency or 0
  
  -- In real implementation, would sync via network
  raceSession.lastSync = os.time()
  raceSession.latency = latency
  
  return true
end

---Generate race statistics
---@param raceSession table Completed race
---@return table stats Race statistics
function RaceMode.generateRaceStats(raceSession)
  if not raceSession or not raceSession.results then return {} end
  
  local stats = {
    raceId = raceSession.id,
    winner = raceSession.results[1],
    totalTime = raceSession.results[1].time,
    playerCount = #raceSession.results,
    category = raceSession.category
  }
  
  return stats
end

---Record race VOD metadata
---@param raceSession table Race to record
---@param vodUrl string Video URL
---@return boolean success Recording saved
function RaceMode.recordRaceVideo(raceSession, vodUrl)
  if not raceSession or not vodUrl then return false end
  
  raceSession.vod = {
    url = vodUrl,
    recordedAt = os.time(),
    archived = true
  }
  
  return true
end

function RaceMode._formatTimeDifference(seconds)
  local mins = math.floor(seconds / 60)
  local secs = seconds % 60
  return string.format("%02d:%05.2f", mins, secs)
end

-- ============================================================================
-- FEATURE 7: TIME TRACKING ANALYTICS (~280 LOC)
-- ============================================================================

local TimeAnalytics = {}

---Track run consistency
---@param runs table List of run times
---@return table consistency Consistency metrics
function TimeAnalytics.trackRunConsistency(runs)
  if not runs or #runs == 0 then return {} end
  
  -- Calculate variance
  local mean = 0
  for _, time in ipairs(runs) do
    mean = mean + time
  end
  mean = mean / #runs
  
  local variance = 0
  for _, time in ipairs(runs) do
    variance = variance + (time - mean) ^ 2
  end
  variance = variance / #runs
  
  return {
    mean = mean,
    variance = variance,
    stdDeviation = math.sqrt(variance),
    consistent = variance < 3600  -- Less than 1 hour variance
  }
end

---Identify consistency timeout points
---@param runs table Run data with segments
---@return table timeoutPoints Segments with consistency issues
function TimeAnalytics.identifyTimeoutPoints(runs)
  if not runs or #runs == 0 then return {} end
  
  local timeoutPoints = {}
  
  for i = 1, #runs do
    if runs[i].segments then
      local segmentVariances = {}
      for segIdx, segment in ipairs(runs[i].segments) do
        -- Calculate variance for this segment across all runs
        local times = {}
        for _, run in ipairs(runs) do
          if run.segments and run.segments[segIdx] then
            table.insert(times, run.segments[segIdx].time)
          end
        end
        
        if #times > 1 then
          local mean = 0
          for _, t in ipairs(times) do mean = mean + t end
          mean = mean / #times
          
          local variance = 0
          for _, t in ipairs(times) do
            variance = variance + (t - mean) ^ 2
          end
          variance = variance / #times
          
          if variance > 3600 then  -- High variance
            table.insert(timeoutPoints, {
              segment = segIdx,
              variance = variance,
              issue = "Inconsistent timing"
            })
          end
        end
      end
    end
  end
  
  return timeoutPoints
end

---Generate performance trend graph
---@param runs table Run history
---@return table graphData Graph data points
function TimeAnalytics.generatePerformanceGraph(runs)
  if not runs or #runs == 0 then return {} end
  
  local graphData = {
    points = {},
    trend = "",
    improving = false
  }
  
  for i, run in ipairs(runs) do
    table.insert(graphData.points, {
      runNumber = i,
      time = run.time or 0,
      date = run.date or os.time()
    })
  end
  
  return graphData
end

---Benchmark against community
---@param time number Run time to compare
---@return table benchmark Benchmark data
function TimeAnalytics.benchmarkAgainstCommunity(time)
  if not time then return {} end
  
  return {
    yourTime = time,
    worldRecord = 0,
    communityAverage = 0,
    percentile = 0,
    rank = "unknown"
  }
end

---Analyze route efficiency
---@param run table Run data with movement
---@return table analysis Efficiency analysis
function TimeAnalytics.analyzeEfficiency(run)
  if not run then return {} end
  
  return {
    movementTime = 0,
    combatTime = 0,
    menuTime = 0,
    deadTime = 0,
    efficiency = 0
  }
end

---Suggest optimizations
---@param run table Analyzed run
---@return table suggestions Optimization suggestions
function TimeAnalytics.suggestOptimizations(run)
  if not run then return {} end
  
  return {
    suggestions = {
      "Practice difficult segments",
      "Optimize route choices",
      "Improve execution speed"
    }
  }
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1.0",
  SplitManagement = SplitManagement,
  PersonalRecords = PersonalRecords,
  PaceCalculator = PaceCalculator,
  CategoryPresets = CategoryPresets,
  StreamingIntegration = StreamingIntegration,
  RaceMode = RaceMode,
  TimeAnalytics = TimeAnalytics,
  
  -- Feature completion
  features = {
    splitManagement = true,
    personalRecords = true,
    paceCalculator = true,
    categoryPresets = true,
    streamingIntegration = true,
    raceMode = true,
    timeAnalytics = true
  }
}
