--[[
  Randomizer Mode Plugin
  Randomize various game elements for replayability
  
  Features:
  - Randomize character starting stats (within ranges)
  - Random starting equipment per character
  - Shuffle character command abilities
  - Random esper assignments
  - Random starting inventory
  - Shuffle spell learning (esper to spell assignments)
  - Random character name generator
  - Seed-based randomization (reproducible)
  - Randomization intensity levels (Mild/Moderate/Chaos)
  - Export randomizer settings
  - Share randomizer seeds
  - Validate randomizer balance
  
  Phase: 6 (Gameplay-Altering Plugins)
  Batch: 3
  Version: 1.0.0
  
  WARNING: This is an EXPERIMENTAL plugin that randomizes game elements.
  Creates unique playthroughs for endless replayability.
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
  -- Stat ranges for randomization
  STAT_RANGES = {
    HP = {min = 30, max = 100},
    MP = {min = 0, max = 60},
    VIGOR = {min = 15, max = 35},
    SPEED = {min = 15, max = 35},
    STAMINA = {min = 15, max = 35},
    MAGIC = {min = 15, max = 35},
    DEFENSE = {min = 10, max = 30},
    M_DEFENSE = {min = 10, max = 30},
    EVADE = {min = 0, max = 20},
    M_EVADE = {min = 0, max = 20}
  },
  
  -- Randomization intensity levels
  INTENSITY_LEVELS = {
    {id = "mild", name = "Mild", desc = "Small variations, recognizable gameplay", variance = 0.2},
    {id = "moderate", name = "Moderate", desc = "Significant changes, still balanced", variance = 0.5},
    {id = "chaos", name = "Chaos", desc = "Extreme randomization, wild gameplay", variance = 1.0}
  },
  
  -- Character names for random generation
  CHARACTER_NAMES = {
    "Aria", "Blade", "Cole", "Dane", "Echo", "Flame",
    "Grim", "Hope", "Iris", "Jorn", "Kael", "Lyra",
    "Mage", "Nova", "Omen", "Pike", "Quest", "Rune",
    "Sage", "Thorn", "Unity", "Vale", "Wave", "Xero",
    "Yara", "Zest"
  },
  
  -- Command abilities
  COMMANDS = {
    {id = 1, name = "Fight"},
    {id = 2, name = "Magic"},
    {id = 4, name = "Steal"},
    {id = 6, name = "SwdTech"},
    {id = 7, name = "Throw"},
    {id = 8, name = "Tools"},
    {id = 9, name = "Blitz"},
    {id = 10, name = "Runic"},
    {id = 11, name = "Lore"},
    {id = 12, name = "Sketch"},
    {id = 13, name = "Control"},
    {id = 14, name = "Slot"},
    {id = 15, name = "Rage"},
    {id = 18, name = "Dance"}
  },
  
  -- Characters
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

-- Seed-based random number generator
local function createSeededRandom(seed)
  local mt = {seed = seed or os.time()}
  
  function mt:random(min, max)
    min = min or 0
    max = max or 1
    -- Simple LCG (Linear Congruential Generator)
    self.seed = (self.seed * 1103515245 + 12345) % 2147483648
    local value = self.seed / 2147483648
    return min + (value * (max - min))
  end
  
  function mt:randint(min, max)
    return math.floor(self:random(min, max))
  end
  
  function mt:shuffle(array)
    for i = #array, 2, -1 do
      local j = self:randint(1, i)
      array[i], array[j] = array[j], array[i]
    end
    return array
  end
  
  return mt
end

-- ============================================================================
-- RANDOMIZATION FUNCTIONS
-- ============================================================================

-- Randomize character starting stats
local function randomizeCharacterStats(charId, intensity, rng)
  local variance = intensity.variance
  local randomizedStats = {}
  
  for stat, range in pairs(CONFIG.STAT_RANGES) do
    local mid = (range.min + range.max) / 2
    local halfRange = (range.max - range.min) / 2
    local adjustment = rng:random(-halfRange * variance, halfRange * variance)
    local newValue = math.floor(mid + adjustment)
    randomizedStats[stat] = math.max(range.min, math.min(range.max, newValue))
  end
  
  return randomizedStats
end

-- Shuffle command abilities
local function shuffleCommandAbilities(rng)
  local commands = {}
  for i, cmd in ipairs(CONFIG.COMMANDS) do
    table.insert(commands, cmd)
  end
  
  return rng:shuffle(commands)
end

-- Generate random character names
local function generateRandomNames(count, rng)
  local names = {}
  for i = 1, count do
    local idx = rng:randint(1, #CONFIG.CHARACTER_NAMES)
    table.insert(names, CONFIG.CHARACTER_NAMES[idx])
  end
  return names
end

-- Randomize esper assignments
local function randomizeEsperAssignments(rng)
  local assignments = {}
  
  -- 26 espers randomly assigned to characters
  local esperList = {}
  for i = 1, 26 do
    table.insert(esperList, i)
  end
  
  esperList = rng:shuffle(esperList)
  
  for i, charId in ipairs({0, 1, 2, 4, 5, 6, 7, 8, 10}) do
    -- Assign 2-3 espers per character
    local esperCount = rng:randint(2, 3)
    assignments[charId] = {}
    for j = 1, esperCount do
      if #esperList > 0 then
        table.insert(assignments[charId], table.remove(esperList))
      end
    end
  end
  
  return assignments
end

-- Shuffle spell learning (randomize esper spell assignments)
local function shuffleSpellLearning(rng)
  local spells = {}
  for i = 1, 50 do  -- ~50 spells in FF6
    table.insert(spells, i)
  end
  
  return rng:shuffle(spells)
end

-- Randomize starting equipment
local function randomizeStartingEquipment(rng)
  local equipment = {}
  
  for i, char in ipairs(CONFIG.CHARACTERS) do
    equipment[char.id] = {
      weapon = rng:randint(1, 20),
      shield = rng:randint(0, 10),  -- 0 = none
      helmet = rng:randint(1, 15),
      armor = rng:randint(1, 15),
      relic1 = rng:randint(0, 20),
      relic2 = rng:randint(0, 20)
    }
  end
  
  return equipment
end

-- Randomize starting inventory
local function randomizeStartingInventory(rng)
  local inventory = {}
  
  -- Random consumable items
  for itemId = 1, 20 do  -- Consumable range
    local quantity = rng:randint(0, 5)
    if quantity > 0 then
      inventory[itemId] = quantity
    end
  end
  
  return inventory
end

-- ============================================================================
-- SEED AND EXPORT FUNCTIONS
-- ============================================================================

-- Generate seed from current state
local function generateSeed()
  return os.time()
end

-- Export randomizer configuration
local function exportRandomizerConfig(seed, intensity, config)
  return {
    seed = seed,
    timestamp = os.date("%Y-%m-%d %H:%M:%S"),
    intensity = intensity,
    configuration = config
  }
end

-- Create randomizer export string
local function createExportString(seed, intensity)
  return string.format("FF6_RANDOMIZER_%d_%s_%s", seed, intensity.id, os.date("%Y%m%d"))
end

-- ============================================================================
-- VALIDATION FUNCTIONS
-- ============================================================================

-- Validate randomizer balance
local function validateBalance(config)
  local results = {
    valid = true,
    warnings = {},
    total_stats = 0,
    avg_stats = 0
  }
  
  -- Check for extreme imbalances
  local totalPower = 0
  local charCount = 0
  
  for charId, stats in pairs(config.character_stats or {}) do
    local charPower = 0
    for stat, value in pairs(stats) do
      charPower = charPower + value
    end
    totalPower = totalPower + charPower
    charCount = charCount + 1
  end
  
  if charCount > 0 then
    results.avg_stats = math.floor(totalPower / charCount)
    
    -- Warn if stats are extreme
    if results.avg_stats < 50 then
      table.insert(results.warnings, "WARNING: Average stats very low (may be too easy)")
    elseif results.avg_stats > 300 then
      table.insert(results.warnings, "WARNING: Average stats very high (may be too hard)")
    end
  end
  
  return results
end

-- ============================================================================
-- RANDOMIZER APPLICATION FUNCTIONS
-- ============================================================================

-- Apply randomization with seed
local function applyRandomization(seed, intensityId)
  local intensity = nil
  for i, intLevel in ipairs(CONFIG.INTENSITY_LEVELS) do
    if intLevel.id == intensityId then
      intensity = intLevel
      break
    end
  end
  
  if not intensity then
    logOperation("RANDOMIZER", "Invalid intensity level")
    return {success = false}
  end
  
  -- Create seeded RNG
  local rng = createSeededRandom(seed)
  
  local results = {
    seed = seed,
    intensity = intensity.id,
    timestamp = os.date("%Y-%m-%d %H:%M:%S"),
    randomization_results = {},
    success_count = 0,
    failed_count = 0
  }
  
  -- Apply randomizations
  for i, char in ipairs(CONFIG.CHARACTERS) do
    local result = safeCall(function()
      local stats = randomizeCharacterStats(char.id, intensity, rng)
      table.insert(results.randomization_results, {
        character = char.name,
        stats = stats
      })
      return true
    end)
    
    if result then
      results.success_count = results.success_count + 1
    else
      results.failed_count = results.failed_count + 1
    end
  end
  
  -- Additional randomizations
  local commands = shuffleCommandAbilities(rng)
  local equipment = randomizeStartingEquipment(rng)
  local inventory = randomizeStartingInventory(rng)
  local espers = randomizeEsperAssignments(rng)
  
  results.config = {
    character_stats = results.randomization_results,
    commands = commands,
    equipment = equipment,
    inventory = inventory,
    espers = espers
  }
  
  -- Validate balance
  results.balance_check = validateBalance(results.config)
  
  logOperation("RANDOMIZER", string.format("Applied randomization (seed: %d, intensity: %s)", seed, intensity.id))
  return results
end

-- ============================================================================
-- PRESET INTENSITY PROFILES
-- ============================================================================

-- Get intensity profile
local function getIntensityProfile(intensityId)
  for i, profile in ipairs(CONFIG.INTENSITY_LEVELS) do
    if profile.id == intensityId then
      return profile
    end
  end
  return nil
end

-- ============================================================================
-- UI AND MENU FUNCTIONS
-- ============================================================================

-- Display intensity levels
local function displayIntensityLevels()
  print("\nRANDOMIZATION INTENSITY:\n")
  
  for i, profile in ipairs(CONFIG.INTENSITY_LEVELS) do
    print(string.format("%d. %s (Variance: %d%%)", i, profile.name, math.floor(profile.variance * 100)))
    print(string.format("   %s\n", profile.desc))
  end
  
  return readInput("Select intensity (1-" .. #CONFIG.INTENSITY_LEVELS .. "): ")
end

-- Display randomization options
local function displayRandomizationOptions()
  print("\nRANDOMIZER OPTIONS:\n")
  print("  1. Generate New Randomization (New Seed)")
  print("  2. Use Existing Seed (Share Seeds)")
  print("  3. Randomize Stats Only")
  print("  4. Randomize Equipment Only")
  print("  5. Randomize All Systems")
  print("  6. Validate Current Randomization")
  print("  7. Export Seed/Settings")
  print("  8. Back\n")
end

-- ============================================================================
-- MAIN PLUGIN INTERFACE
-- ============================================================================

local Plugin = {
  name = "Randomizer Mode",
  version = "1.0.0",
  batch = 3,
  lastSeed = nil,
  lastIntensity = nil,
  seedHistory = {}
}

-- Initialize plugin
function Plugin:init()
  logOperation("INIT", "Randomizer Mode loaded")
end

-- Apply randomization with new seed
function Plugin:applyNewRandomization(intensityId)
  local seed = generateSeed()
  self.lastSeed = seed
  self.lastIntensity = intensityId
  
  table.insert(self.seedHistory, {
    seed = seed,
    intensity = intensityId,
    timestamp = os.time()
  })
  
  return applyRandomization(seed, intensityId)
end

-- Apply randomization with existing seed
function Plugin:applyExistingSeed(seed, intensityId)
  self.lastSeed = seed
  self.lastIntensity = intensityId
  return applyRandomization(seed, intensityId)
end

-- Export seed
function Plugin:exportSeed()
  if self.lastSeed then
    return createExportString(self.lastSeed, getIntensityProfile(self.lastIntensity))
  else
    return nil
  end
end

-- Get seed history
function Plugin:getSeedHistory()
  return self.seedHistory
end

-- Get plugin info
function Plugin:getInfo()
  return {
    name = self.name,
    version = self.version,
    batch = self.batch,
    intensities = #CONFIG.INTENSITY_LEVELS,
    description = "Randomize game elements for endless replayability"
  }
end

-- ============================================================================
-- EXPORT
-- ============================================================================

return Plugin
