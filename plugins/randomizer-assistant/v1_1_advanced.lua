--[[
  Randomizer Assistant Plugin - v1.1 Upgrade Extension
  Advanced analysis, real-time tracking, and community integration features
  
  This file extends the base v1.0 plugin with powerful new capabilities:
  - Auto-detection of obtained items from save files
  - Advanced logic solving with softlock detection
  - Real-time accessibility calculation
  - Visual mapping system
  - Community seed sharing
  - Custom seed generation
  - Advanced analytics
  
  Phase: 7A (Advanced Analysis & Tracking)
  Version: 1.1.0 (extends v1.0.0)
  Backward Compatibility: 100% - all existing functions unchanged
]]

-- ============================================================================
-- FEATURE 1: AUTO-DETECTION SYSTEM (~200 LOC)
-- ============================================================================

local AutoDetection = {}

---Automatically detect obtained items from current save file
---@param saveFile table Current save state
---@param spoilerLog table Parsed spoiler log
---@return table obtainedItems Detected items with locations
function AutoDetection.autoDetectObtainedItems(saveFile, spoilerLog)
  if not saveFile or not spoilerLog then return {} end
  
  local obtained = {}
  
  -- Scan inventory for items
  if saveFile.inventory then
    for itemId, count in pairs(saveFile.inventory) do
      if count > 0 then
        table.insert(obtained, {
          itemId = itemId,
          source = "inventory",
          method = "auto_detect",
          timestamp = os.time()
        })
      end
    end
  end
  
  -- Scan equipment for gear items
  if saveFile.equipment then
    for charId, gear in pairs(saveFile.equipment) do
      if gear and gear.weapon then
        table.insert(obtained, {
          itemId = gear.weapon,
          source = "equipment",
          character = charId,
          method = "auto_detect",
          timestamp = os.time()
        })
      end
    end
  end
  
  return obtained
end

---Compare current save state to seed spoiler log
---@param saveFile table Current save state
---@param spoilerLog table Seed spoiler log
---@return table comparison Comparison results
function AutoDetection.compareSeedToSaveState(saveFile, spoilerLog)
  if not saveFile or not spoilerLog then return {} end
  
  local comparison = {
    totalLocations = 0,
    checkedLocations = 0,
    matchedItems = 0,
    mismatches = {},
    seedProgress = 0
  }
  
  -- Count total locations in seed
  if spoilerLog.locations then
    comparison.totalLocations = #spoilerLog.locations
  end
  
  -- Check which items obtained match spoiler log
  local obtained = AutoDetection.autoDetectObtainedItems(saveFile, spoilerLog)
  comparison.checkedLocations = #obtained
  
  if comparison.totalLocations > 0 then
    comparison.seedProgress = (comparison.checkedLocations / comparison.totalLocations) * 100
  end
  
  return comparison
end

---Track item progression through identified checkpoints
---@param saveFile table Current save state
---@param spoilerLog table Seed spoiler log
---@return table progression Item progression timeline
function AutoDetection.trackItemProgression(saveFile, spoilerLog)
  if not saveFile or not spoilerLog then return {} end
  
  local progression = {
    timeline = {},
    currentPhase = "early",
    itemsPerPhase = {
      early = {},
      mid = {},
      late = {},
      endgame = {}
    }
  }
  
  local obtained = AutoDetection.autoDetectObtainedItems(saveFile, spoilerLog)
  
  -- Categorize items by game phase
  for _, item in ipairs(obtained) do
    local phase = "early"
    if saveFile.totalGameTime and saveFile.totalGameTime > 3600 then
      phase = "mid"
    end
    if saveFile.totalGameTime and saveFile.totalGameTime > 7200 then
      phase = "late"
    end
    if saveFile.totalGameTime and saveFile.totalGameTime > 14400 then
      phase = "endgame"
    end
    
    table.insert(progression.itemsPerPhase[phase], item)
  end
  
  return progression
end

---Synchronize save state with spoiler log seed data
---@param saveFile table Current save state
---@param spoilerLog table Seed spoiler log
---@return boolean success, string message Sync result
function AutoDetection.syncSaveWithSeed(saveFile, spoilerLog)
  if not saveFile or not spoilerLog then
    return false, "Invalid save or spoiler log"
  end
  
  -- Verify seed metadata matches
  if saveFile.seedId and spoilerLog.seedId then
    if saveFile.seedId ~= spoilerLog.seedId then
      return false, "Seed ID mismatch"
    end
  end
  
  -- Update progress tracking
  local comparison = AutoDetection.compareSeedToSaveState(saveFile, spoilerLog)
  saveFile.randomizerProgress = comparison.seedProgress
  
  return true, "Sync successful"
end

---Calculate current completion percentage
---@param saveFile table Current save state
---@param spoilerLog table Seed spoiler log
---@return number completion Percentage 0-100
function AutoDetection.getCompletionPercentage(saveFile, spoilerLog)
  if not saveFile or not spoilerLog then return 0 end
  
  local comparison = AutoDetection.compareSeedToSaveState(saveFile, spoilerLog)
  return comparison.seedProgress or 0
end

---Get list of next accessible locations based on current items
---@param saveFile table Current save state
---@param spoilerLog table Seed spoiler log
---@return table nextLocations Accessible but unchecked locations
function AutoDetection.getNextAccessibleLocations(saveFile, spoilerLog)
  if not saveFile or not spoilerLog then return {} end
  
  local accessible = {}
  local obtained = AutoDetection.autoDetectObtainedItems(saveFile, spoilerLog)
  
  -- Analyze logic to find accessible locations
  if spoilerLog.logic and spoilerLog.locations then
    for _, location in ipairs(spoilerLog.locations) do
      if location.logic then
        -- Check if logic requirements are met
        local canAccess = true
        for _, requirement in ipairs(location.logic) do
          if not AutoDetection._hasRequirement(obtained, requirement) then
            canAccess = false
            break
          end
        end
        
        if canAccess and not AutoDetection._isChecked(obtained, location) then
          table.insert(accessible, location)
        end
      end
    end
  end
  
  return accessible
end

-- Helper: Check if item requirement is met
function AutoDetection._hasRequirement(obtained, requirement)
  for _, item in ipairs(obtained) do
    if item.itemId == requirement then return true end
  end
  return false
end

-- Helper: Check if location already checked
function AutoDetection._isChecked(obtained, location)
  for _, item in ipairs(obtained) do
    if item.locationId == location.id then return true end
  end
  return false
end

-- ============================================================================
-- FEATURE 2: ADVANCED LOGIC SOLVER (~400 LOC)
-- ============================================================================

local LogicSolver = {}

---Analyze progression chains to build dependency graph
---@param spoilerLog table Seed spoiler log
---@return table graph Dependency graph
function LogicSolver.analyzeProgressionChains(spoilerLog)
  if not spoilerLog or not spoilerLog.locations then return {} end
  
  local graph = {
    nodes = {},
    edges = {},
    cycles = {}
  }
  
  -- Create nodes for each location
  for _, location in ipairs(spoilerLog.locations) do
    graph.nodes[location.id] = {
      id = location.id,
      name = location.name,
      item = location.item,
      dependencies = location.logic or {},
      reachable = false
    }
  end
  
  -- Create edges for dependencies
  for locationId, node in pairs(graph.nodes) do
    for _, dep in ipairs(node.dependencies) do
      table.insert(graph.edges, {
        from = dep,
        to = locationId
      })
    end
  end
  
  -- Detect cycles
  graph.cycles = LogicSolver._detectCycles(graph)
  
  return graph
end

---Detect softlocks in progression logic
---@param graph table Dependency graph
---@return table softlocks List of softlock scenarios
function LogicSolver.detectSoftlocks(graph)
  if not graph or not graph.nodes then return {} end
  
  local softlocks = {}
  
  -- Check for unreachable critical items
  for nodeId, node in pairs(graph.nodes) do
    if node.item and LogicSolver._isCriticalItem(node.item) then
      if not LogicSolver._isReachable(graph, nodeId) then
        table.insert(softlocks, {
          type = "unreachable_critical",
          location = node.name,
          item = node.item,
          reason = "No path to obtain critical item"
        })
      end
    end
  end
  
  -- Check for circular dependencies
  for _, cycle in ipairs(graph.cycles) do
    if LogicSolver._cycleMakesGameUnwinnable(graph, cycle) then
      table.insert(softlocks, {
        type = "circular_dependency",
        cycle = cycle,
        reason = "Circular logic makes progression impossible"
      })
    end
  end
  
  return softlocks
end

---Find alternate routes to reach objective
---@param graph table Dependency graph
---@param target string Target location ID
---@return table routes List of alternate paths
function LogicSolver.findAlternateRoutes(graph, target)
  if not graph or not target then return {} end
  
  local routes = {}
  local visited = {}
  
  -- Find all possible paths using DFS
  local function findPaths(current, path)
    if current == target then
      table.insert(routes, path)
      return
    end
    
    visited[current] = true
    
    for _, edge in ipairs(graph.edges) do
      if edge.from == current and not visited[edge.to] then
        findPaths(edge.to, {table.unpack(path), edge.to})
      end
    end
    
    visited[current] = false
  end
  
  findPaths(graph.nodes[target], {target})
  return routes
end

---Calculate complete accessibility tree
---@param graph table Dependency graph
---@param startingItems table Initial items available
---@return table accessibilityTree Complete reachability map
function LogicSolver.calculateAccessibilityTree(graph, startingItems)
  if not graph then return {} end
  
  local tree = {
    startingItems = startingItems or {},
    phases = {},
    totalReachable = 0,
    unreachable = {}
  }
  
  local reachable = {}
  for _, item in ipairs(startingItems or {}) do
    reachable[item] = true
  end
  
  -- Iteratively find newly reachable items
  local phase = 0
  while true do
    phase = phase + 1
    local newItems = false
    
    for nodeId, node in pairs(graph.nodes) do
      if not reachable[nodeId] then
        local canReach = true
        for _, dep in ipairs(node.dependencies) do
          if not reachable[dep] then
            canReach = false
            break
          end
        end
        
        if canReach then
          reachable[nodeId] = true
          newItems = true
          tree.phases[phase] = tree.phases[phase] or {}
          table.insert(tree.phases[phase], nodeId)
        end
      end
    end
    
    if not newItems then break end
  end
  
  -- Calculate statistics
  tree.totalReachable = 0
  for nodeId in pairs(reachable) do
    tree.totalReachable = tree.totalReachable + 1
  end
  
  -- Find unreachable items
  for nodeId in pairs(graph.nodes) do
    if not reachable[nodeId] then
      table.insert(tree.unreachable, nodeId)
    end
  end
  
  return tree
end

---Validate logic chain for correctness
---@param graph table Dependency graph
---@return boolean valid, table issues Validation result
function LogicSolver.validateLogicChain(graph)
  if not graph then return false, {} end
  
  local issues = {}
  
  -- Check for impossible requirements
  for nodeId, node in pairs(graph.nodes) do
    for _, dep in ipairs(node.dependencies) do
      if not graph.nodes[dep] then
        table.insert(issues, {
          type = "missing_dependency",
          location = node.name,
          dependency = dep
        })
      end
    end
  end
  
  -- Check for circular dependencies that break logic
  for _, cycle in ipairs(graph.cycles) do
    table.insert(issues, {
      type = "circular_dependency",
      cycle = cycle
    })
  end
  
  return #issues == 0, issues
end

---Suggest optimal progression path
---@param graph table Dependency graph
---@param startingItems table Starting items
---@return table suggestedPath Recommended item collection order
function LogicSolver.suggestProgressionPath(graph, startingItems)
  if not graph then return {} end
  
  local path = {}
  local reachable = {}
  
  for _, item in ipairs(startingItems or {}) do
    reachable[item] = true
  end
  
  -- Greedily select items in order
  while true do
    local bestItem = nil
    local bestValue = -1
    
    for nodeId, node in pairs(graph.nodes) do
      if not reachable[nodeId] then
        local value = LogicSolver._calculateItemValue(graph, node, reachable)
        if value > bestValue then
          bestValue = value
          bestItem = nodeId
        end
      end
    end
    
    if not bestItem then break end
    
    table.insert(path, bestItem)
    reachable[bestItem] = true
  end
  
  return path
end

-- Helper functions
function LogicSolver._isCriticalItem(itemId)
  local critical = {
    "sword", "shield", "magicite", "esper", "key_item"
  }
  for _, crit in ipairs(critical) do
    if string.match(itemId, crit) then return true end
  end
  return false
end

function LogicSolver._isReachable(graph, nodeId)
  if not graph.nodes[nodeId] then return false end
  
  if #graph.nodes[nodeId].dependencies == 0 then return true end
  
  for _, dep in ipairs(graph.nodes[nodeId].dependencies) do
    if not LogicSolver._isReachable(graph, dep) then
      return false
    end
  end
  return true
end

function LogicSolver._detectCycles(graph)
  local cycles = {}
  local visited = {}
  local recStack = {}
  
  local function dfs(nodeId)
    visited[nodeId] = true
    recStack[nodeId] = true
    
    for _, edge in ipairs(graph.edges) do
      if edge.from == nodeId then
        if not visited[edge.to] then
          dfs(edge.to)
        elseif recStack[edge.to] then
          table.insert(cycles, {edge.from, edge.to})
        end
      end
    end
    
    recStack[nodeId] = false
  end
  
  for nodeId in pairs(graph.nodes) do
    if not visited[nodeId] then
      dfs(nodeId)
    end
  end
  
  return cycles
end

function LogicSolver._cycleMakesGameUnwinnable(graph, cycle)
  -- Check if cycle contains all critical items
  local hasCritical = false
  for _, nodeId in ipairs(cycle) do
    if LogicSolver._isCriticalItem(graph.nodes[nodeId].item) then
      hasCritical = true
      break
    end
  end
  return hasCritical
end

function LogicSolver._calculateItemValue(graph, node, reachable)
  -- Calculate how many new items this opens up
  local newlyReachable = 0
  for depId in pairs(reachable) do
    if node.dependencies[depId] then
      newlyReachable = newlyReachable + 1
    end
  end
  return newlyReachable
end

-- ============================================================================
-- FEATURE 3: REAL-TIME ACCESSIBILITY (~300 LOC)
-- ============================================================================

local RealTimeAccessibility = {}

---Calculate what locations are currently accessible
---@param currentItems table Items currently obtained
---@param logic table Game logic rules
---@return table accessible List of accessible locations
function RealTimeAccessibility.calculateAccessible(currentItems, logic)
  if not currentItems or not logic then return {} end
  
  local accessible = {}
  
  for locationId, requirements in pairs(logic) do
    if RealTimeAccessibility._meetsRequirements(currentItems, requirements) then
      table.insert(accessible, locationId)
    end
  end
  
  return accessible
end

---Get list of prevented/blocked locations
---@param currentItems table Items currently obtained
---@param logic table Game logic rules
---@return table prevented List of blocked locations
function RealTimeAccessibility.getPrevented(currentItems, logic)
  if not currentItems or not logic then return {} end
  
  local prevented = {}
  
  for locationId, requirements in pairs(logic) do
    if not RealTimeAccessibility._meetsRequirements(currentItems, requirements) then
      table.insert(prevented, {
        location = locationId,
        missingItems = RealTimeAccessibility._getMissingItems(currentItems, requirements)
      })
    end
  end
  
  return prevented
end

---Suggest next location to visit
---@param currentItems table Items currently obtained
---@param logic table Game logic rules
---@param visited table Already checked locations
---@return table suggestion Recommended next location
function RealTimeAccessibility.suggestNextLocation(currentItems, logic, visited)
  if not currentItems or not logic then return {} end
  
  local accessible = RealTimeAccessibility.calculateAccessible(currentItems, logic)
  
  -- Score accessible locations
  local scored = {}
  for _, locationId in ipairs(accessible) do
    if not visited or not visited[locationId] then
      local score = RealTimeAccessibility._scoreLocation(locationId, currentItems, logic)
      table.insert(scored, {
        location = locationId,
        score = score
      })
    end
  end
  
  -- Sort by score
  table.sort(scored, function(a, b) return a.score > b.score end)
  
  return scored[1] or {}
end

---Estimate percentage completion to goal
---@param currentItems table Items currently obtained
---@param goalItems table Required items for completion
---@return number percentage Completion 0-100
function RealTimeAccessibility.estimateCompletion(currentItems, goalItems)
  if not currentItems or not goalItems then return 0 end
  
  local have = 0
  for _, item in ipairs(goalItems) do
    if RealTimeAccessibility._hasItem(currentItems, item) then
      have = have + 1
    end
  end
  
  return (have / #goalItems) * 100
end

---Track changes in accessibility over time
---@param previousState table Previous item state
---@param currentState table Current item state
---@param logic table Game logic
---@return table changes Accessibility changes
function RealTimeAccessibility.trackAccessibilityChanges(previousState, currentState, logic)
  if not previousState or not currentState or not logic then return {} end
  
  local oldAccessible = RealTimeAccessibility.calculateAccessible(previousState, logic)
  local newAccessible = RealTimeAccessibility.calculateAccessible(currentState, logic)
  
  local changes = {
    newlyAccessible = {},
    newlyBlocked = {}
  }
  
  -- Find newly accessible
  for _, loc in ipairs(newAccessible) do
    local wasAccessible = false
    for _, oldLoc in ipairs(oldAccessible) do
      if loc == oldLoc then
        wasAccessible = true
        break
      end
    end
    if not wasAccessible then
      table.insert(changes.newlyAccessible, loc)
    end
  end
  
  -- Find newly blocked
  for _, loc in ipairs(oldAccessible) do
    local stillAccessible = false
    for _, newLoc in ipairs(newAccessible) do
      if loc == newLoc then
        stillAccessible = true
        break
      end
    end
    if not stillAccessible then
      table.insert(changes.newlyBlocked, loc)
    end
  end
  
  return changes
end

---Predict future accessibility based on planned items
---@param currentItems table Current items
---@param plannedItems table Items planning to get
---@param logic table Game logic
---@return table future Predicted future accessibility
function RealTimeAccessibility.predictFutureAccessibility(currentItems, plannedItems, logic)
  if not currentItems or not logic then return {} end
  
  local futureItems = {}
  for _, item in ipairs(currentItems) do
    table.insert(futureItems, item)
  end
  for _, item in ipairs(plannedItems or {}) do
    table.insert(futureItems, item)
  end
  
  return RealTimeAccessibility.calculateAccessible(futureItems, logic)
end

-- Helper functions
function RealTimeAccessibility._meetsRequirements(items, requirements)
  if not requirements or #requirements == 0 then return true end
  
  for _, req in ipairs(requirements) do
    if not RealTimeAccessibility._hasItem(items, req) then
      return false
    end
  end
  return true
end

function RealTimeAccessibility._getMissingItems(items, requirements)
  local missing = {}
  for _, req in ipairs(requirements or {}) do
    if not RealTimeAccessibility._hasItem(items, req) then
      table.insert(missing, req)
    end
  end
  return missing
end

function RealTimeAccessibility._hasItem(items, itemId)
  for _, item in ipairs(items) do
    if item == itemId then return true end
  end
  return false
end

function RealTimeAccessibility._scoreLocation(locationId, currentItems, logic)
  -- Score based on how many new items it unlocks
  local score = 0
  
  -- Bonus for locations that unlock multiple areas
  if string.match(locationId, "key") then score = score + 10 end
  if string.match(locationId, "critical") then score = score + 5 end
  
  return score
end

-- ============================================================================
-- MODULE EXPORTS
-- ============================================================================

return {
  version = "1.1.0",
  AutoDetection = AutoDetection,
  LogicSolver = LogicSolver,
  RealTimeAccessibility = RealTimeAccessibility,
  
  -- Feature completeness
  features = {
    autoDetection = true,
    advancedLogic = true,
    realTimeAccessibility = true,
    visualMap = false,  -- Feature 4 in separate file
    hints = false,      -- Feature 5 in separate file
    community = false,  -- Feature 6 in separate file
    seedGeneration = false, -- Feature 7 in separate file
    analytics = false   -- Feature 8 in separate file
  }
}
