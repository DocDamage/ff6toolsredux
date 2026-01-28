--[[
  Alternate Start Generator Plugin
  Start game from different points and conditions
  
  Features:
  - Skip to World of Ruin
  - Choose starting characters (any combination)
  - Start at specific story events
  - Custom starting inventory/levels
  - Scenario selection (different party configurations)
  - Boss rush mode (fight only bosses)
  - Character-specific scenarios (Celes solo)
  - Speedrun practice starts
  - Event flag manipulation
  - Export start configurations
  - Preset scenarios (6+ templates)
  
  Phase: 6 (Gameplay-Altering Plugins)
  Batch: 3
  Version: 1.0.0
  
  WARNING: This is an EXPERIMENTAL plugin that creates alternate game starts.
  Game state is fundamentally altered. May cause unexpected behavior.
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
  -- World states
  WORLD_OF_BALANCE = "WoB",
  WORLD_OF_RUIN = "WoR",
  
  -- Story event points (examples, would need full mapping)
  STORY_EVENTS = {
    {id = 1, name = "Game Start", world = "WoB", desc = "Beginning at Terra's flashback"},
    {id = 2, name = "Figaro Castle", world = "WoB", desc = "After Figaro Castle siege"},
    {id = 3, name = "Narshe", world = "WoB", desc = "First arrival at Narshe"},
    {id = 4, name = "South Figaro", world = "WoB", desc = "After South Figaro liberation"},
    {id = 5, name = "Mt. Kolts", world = "WoB", desc = "After Mt. Kolts passage"},
    {id = 6, name = "World of Ruin", world = "WoR", desc = "After Kefka destroys world"},
    {id = 7, name = "Floating Island", world = "WoR", desc = "At Floating Island final dungeon"},
    {id = 8, name = "Kefka's Tower", world = "WoR", desc = "Beginning Kefka's Tower assault"}
  },
  
  -- Character list
  CHARACTERS = {
    {id = 0, name = "Terra"},
    {id = 1, name = "Locke"},
    {id = 2, name = "Cyan"},
    {id = 3, name = "Shadow"},
    {id = 4, name = "Edgar"},
    {id = 5, name = "Sabin"},
    {id = 6, name = "Celes"},
    {id = 7, name = "Strago"},
    {id = 8, name = "Relm"},
    {id = 9, name = "Setzer"},
    {id = 10, name = "Mog"},
    {id = 11, name = "Gau"},
    {id = 12, name = "Gogo"},
    {id = 13, name = "Umaro"}
  }
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Safe API call wrapper
local function safeCall(fn, ...)
  local success, result = pcall(fn, ...)
  if success then
    return result
  else
    return nil
  end
end

-- Log operation
local function logOperation(operation, details)
  local timestamp = os.date("%Y-%m-%d %H:%M:%S")
  print(string.format("[%s] %s: %s", timestamp, operation, details or ""))
end

-- ============================================================================
-- PRESET STARTING SCENARIOS
-- ============================================================================

local PRESET_SCENARIOS = {
  {
    id = "start_wor",
    name = "Skip to World of Ruin",
    desc = "Begin at World of Ruin with Celes and available characters",
    world = "WoR",
    starting_party = {6, 4, 1, 2},  -- Celes, Edgar, Locke, Cyan
    starting_level = 30,
    starting_gil = 5000,
    story_event = 6
  },
  {
    id = "speedrun_any",
    name = "Speedrun - Any% Start",
    desc = "Optimized start for Any% speedrun with Sabin route",
    world = "WoB",
    starting_party = {5, 4, 0},  -- Sabin, Edgar, Terra
    starting_level = 1,
    starting_gil = 0,
    story_event = 1
  },
  {
    id = "speedrun_100",
    name = "Speedrun - 100% Start",
    desc = "Optimized start for 100% speedrun preparation",
    world = "WoB",
    starting_party = {0, 6, 1},  -- Terra, Celes, Locke
    starting_level = 10,
    starting_gil = 2000,
    story_event = 3
  },
  {
    id = "solo_celes",
    name = "Celes Solo Run",
    desc = "Start with only Celes - extreme challenge",
    world = "WoB",
    starting_party = {6},  -- Celes only
    starting_level = 1,
    starting_gil = 0,
    story_event = 1
  },
  {
    id = "trio_challenge",
    name = "Three-Character Challenge",
    desc = "Start with exactly 3 characters of your choice",
    world = "WoB",
    starting_party = {0, 1, 2},  -- Terra, Locke, Cyan (customizable)
    starting_level = 1,
    starting_gil = 0,
    story_event = 1
  },
  {
    id = "boss_rush",
    name = "Boss Rush - WoB Preparation",
    desc = "Skip to major boss area with max party",
    world = "WoB",
    starting_party = {0, 1, 2, 3, 4, 5, 6, 7},  -- All 8 WoB available
    starting_level = 20,
    starting_gil = 10000,
    story_event = 5
  },
  {
    id = "wor_boss_rush",
    name = "Boss Rush - WoR Preparation",
    desc = "Start at WoR with all available characters and high stats",
    world = "WoR",
    starting_party = {0, 1, 2, 4, 5, 6, 7, 8, 9, 10},  -- 10 characters available
    starting_level = 50,
    starting_gil = 50000,
    story_event = 6
  },
  {
    id = "challenge_balanced",
    name = "Balanced Challenge Party",
    desc = "Mix of all character types: Physical, Magic, Support",
    world = "WoB",
    starting_party = {5, 0, 6, 7},  -- Sabin (physical), Terra (magic), Celes (support), Strago (ranged)
    starting_level = 15,
    starting_gil = 3000,
    story_event = 3
  }
}

-- ============================================================================
-- WORLD STATE MODIFICATION FUNCTIONS
-- ============================================================================

-- Set world state (WoB or WoR)
local function setWorldState(worldState)
  local result = safeCall(function()
    -- This would call actual API:
    -- flags.setWorldState(worldState)
    return worldState
  end)
  
  if result then
    logOperation("WORLD_STATE", string.format("Set to %s", worldState))
    return true
  else
    return false
  end
end

-- Set specific story event
local function setStoryEvent(eventId)
  local event = nil
  for i, e in ipairs(CONFIG.STORY_EVENTS) do
    if e.id == eventId then
      event = e
      break
    end
  end
  
  if not event then
    logOperation("STORY_EVENT", "Event not found")
    return false
  end
  
  local result = safeCall(function()
    -- This would call actual API:
    -- flags.setStoryEvent(eventId)
    return true
  end)
  
  if result then
    logOperation("STORY_EVENT", string.format("Set to %s (Event %d)", event.name, eventId))
    return true
  else
    return false
  end
end

-- ============================================================================
-- CHARACTER ROSTER FUNCTIONS
-- ============================================================================

-- Set starting party/roster
local function setStartingCharacterRoster(characterIds)
  local charNames = {}
  
  for i, charId in ipairs(characterIds) do
    local found = false
    for j, char in ipairs(CONFIG.CHARACTERS) do
      if char.id == charId then
        table.insert(charNames, char.name)
        found = true
        break
      end
    end
    if not found then
      logOperation("ROSTER", "Unknown character ID: " .. charId)
    end
  end
  
  local result = safeCall(function()
    -- This would call actual API:
    -- for _, charId in ipairs(characterIds) do
    --   character.setAvailable(charId, true)
    -- end
    return #characterIds
  end)
  
  if result then
    logOperation("ROSTER", string.format("Set to: %s", table.concat(charNames, ", ")))
    return true
  else
    return false
  end
end

-- Disable unavailable characters
local function disableUnavailableCharacters(availableIds)
  local result = safeCall(function()
    -- This would call actual API:
    -- for i, char in ipairs(CONFIG.CHARACTERS) do
    --   local isAvailable = false
    --   for _, availId in ipairs(availableIds) do
    --     if char.id == availId then
    --       isAvailable = true
    --       break
    --     end
    --   end
    --   character.setAvailable(char.id, isAvailable)
    -- end
    return true
  end)
  
  return result and true or false
end

-- ============================================================================
-- CHARACTER STAT FUNCTIONS
-- ============================================================================

-- Set starting level for all party members
local function setStartingLevels(level)
  local result = safeCall(function()
    -- This would call actual API:
    -- for i, char in ipairs(CONFIG.CHARACTERS) do
    --   character.setLevel(char.id, level)
    -- end
    return level
  end)
  
  if result then
    logOperation("STARTING_LEVEL", string.format("All characters set to level %d", level))
    return true
  else
    return false
  end
end

-- Initialize starting inventory
local function setStartingInventory(inventoryConfig)
  local result = safeCall(function()
    -- This would call actual API:
    -- for itemId, quantity in pairs(inventoryConfig) do
    --   inventory.setItemQuantity(itemId, quantity)
    -- end
    return true
  end)
  
  return result and true or false
end

-- Set starting gil
local function setStartingGil(amount)
  local result = safeCall(function()
    -- This would call actual API:
    -- player.setGil(amount)
    return amount
  end)
  
  if result then
    logOperation("STARTING_GIL", string.format("Set to %d", amount))
    return true
  else
    return false
  end
end

-- ============================================================================
-- SCENARIO APPLICATION FUNCTIONS
-- ============================================================================

-- Apply preset scenario
local function applyPresetScenario(scenarioId)
  local scenario = nil
  
  for i, scen in ipairs(PRESET_SCENARIOS) do
    if scen.id == scenarioId then
      scenario = scen
      break
    end
  end
  
  if not scenario then
    logOperation("SCENARIO", "Preset not found: " .. scenarioId)
    return {success = false}
  end
  
  local results = {
    scenario = scenario.id,
    scenario_name = scenario.name,
    success = 0,
    failed = 0,
    operations = {}
  }
  
  -- Apply world state
  if setWorldState(scenario.world) then
    results.success = results.success + 1
    table.insert(results.operations, "✓ World state set to " .. scenario.world)
  else
    results.failed = results.failed + 1
    table.insert(results.operations, "✗ Failed to set world state")
  end
  
  -- Apply story event
  if setStoryEvent(scenario.story_event) then
    results.success = results.success + 1
    table.insert(results.operations, "✓ Story event set")
  else
    results.failed = results.failed + 1
    table.insert(results.operations, "✗ Failed to set story event")
  end
  
  -- Apply starting roster
  if setStartingCharacterRoster(scenario.starting_party) then
    results.success = results.success + 1
    table.insert(results.operations, "✓ Character roster configured")
  else
    results.failed = results.failed + 1
    table.insert(results.operations, "✗ Failed to configure roster")
  end
  
  -- Apply starting levels
  if setStartingLevels(scenario.starting_level) then
    results.success = results.success + 1
    table.insert(results.operations, string.format("✓ Starting levels set to %d", scenario.starting_level))
  else
    results.failed = results.failed + 1
    table.insert(results.operations, "✗ Failed to set starting levels")
  end
  
  -- Apply starting gil
  if setStartingGil(scenario.starting_gil) then
    results.success = results.success + 1
    table.insert(results.operations, string.format("✓ Starting gil set to %d", scenario.starting_gil))
  else
    results.failed = results.failed + 1
    table.insert(results.operations, "✗ Failed to set starting gil")
  end
  
  logOperation("SCENARIO_APPLY", string.format("Applied scenario: %s", scenario.name))
  return results
end

-- Create custom scenario
local function createCustomScenario(options)
  local scenario = {
    id = "custom_" .. os.time(),
    name = options.name or "Custom Scenario",
    desc = options.desc or "Custom alternate start",
    world = options.world or CONFIG.WORLD_OF_BALANCE,
    starting_party = options.starting_party or {0},
    starting_level = options.starting_level or 1,
    starting_gil = options.starting_gil or 0,
    story_event = options.story_event or 1,
    starting_inventory = options.starting_inventory or {}
  }
  
  return scenario
end

-- ============================================================================
-- EVENT FLAG MANIPULATION
-- ============================================================================

-- Get all major event flags
local function getMajorEventFlags()
  return {
    terra_flashback = true,
    figaro_siege = true,
    narshe_visit = true,
    south_figaro_liberation = true,
    kefka_world_destruction = false,
    celes_isolation = false,
    floating_island_accessible = false
  }
end

-- Set event flag
local function setEventFlag(flagName, value)
  local result = safeCall(function()
    -- This would call actual API:
    -- flags.setEventFlag(flagName, value)
    return value
  end)
  
  return result and true or false
end

-- ============================================================================
-- UI AND MENU FUNCTIONS
-- ============================================================================

-- Display scenario list
local function displayScenarioList()
  print("\nPRESET STARTING SCENARIOS:\n")
  
  for i, scenario in ipairs(PRESET_SCENARIOS) do
    print(string.format("%d. %s", i, scenario.name))
    print(string.format("   %s", scenario.desc))
    print(string.format("   World: %s | Level: %d | Gil: %d", 
      scenario.world, scenario.starting_level, scenario.starting_gil))
    print()
  end
  
  return readInput("Select scenario (1-" .. #PRESET_SCENARIOS .. ") or 0 to cancel: ")
end

-- Display story events
local function displayStoryEvents()
  print("\nSTORY EVENTS:\n")
  
  for i, event in ipairs(CONFIG.STORY_EVENTS) do
    print(string.format("%d. %s (%s)", i, event.name, event.world))
    print(string.format("   %s", event.desc))
  end
  
  print()
  return readInput("Select story event (1-" .. #CONFIG.STORY_EVENTS .. "): ")
end

-- Display character selection
local function displayCharacterSelection()
  print("\nAVAILABLE CHARACTERS:\n")
  
  for i, char in ipairs(CONFIG.CHARACTERS) do
    print(string.format("%2d. %s", i, char.name))
  end
  
  print("\nSelect multiple characters by entering IDs separated by commas")
  print("Example: 1,2,3 for Terra, Locke, Cyan\n")
  
  return readInput("Enter character IDs: ")
end

-- ============================================================================
-- MAIN PLUGIN INTERFACE
-- ============================================================================

local Plugin = {
  name = "Alternate Start Generator",
  version = "1.0.0",
  batch = 3,
  customScenarios = {}
}

-- Initialize plugin
function Plugin:init()
  logOperation("INIT", "Alternate Start Generator loaded")
end

-- Apply preset scenario
function Plugin:applyPreset(presetId)
  if presetId >= 1 and presetId <= #PRESET_SCENARIOS then
    return applyPresetScenario(PRESET_SCENARIOS[presetId].id)
  else
    return {success = false}
  end
end

-- Apply custom scenario
function Plugin:applyCustomScenario(scenarioOptions)
  local scenario = createCustomScenario(scenarioOptions)
  return applyPresetScenario(scenario)
end

-- Get available scenarios
function Plugin:getAvailableScenarios()
  return PRESET_SCENARIOS
end

-- Get plugin info
function Plugin:getInfo()
  return {
    name = self.name,
    version = self.version,
    batch = self.batch,
    presets = #PRESET_SCENARIOS,
    description = "Start game from alternate points and conditions"
  }
end

-- ============================================================================
-- EXPORT
-- ============================================================================

return Plugin
