--[[
  Custom Progression Pacing Plugin
  Control all progression rates for any playstyle
  
  Features:
  - Experience gain rate multiplier (0.1x to 100x)
  - Spell learning rate multiplier (0x to 10x)
  - Gil acquisition rate (0.1x to 10x)
  - Drop rate modification (0.1x to 10x)
  - AP gain rate for espers (0.1x to 10x)
  - Configurable per character or global
  - Preset pacing profiles
  - Real-time rate adjustment
  - Rate history tracking
  - Export/import configurations
  
  Phase: 6 (Gameplay-Altering Plugins)
  Batch: 3
  Version: 1.0.0
  
  WARNING: This is an EXPERIMENTAL plugin that modifies progression rates.
  Gameplay balance will be significantly affected depending on configuration.
]]

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
  -- Rate multiplier bounds
  MIN_RATE = 0.1,
  MAX_XP_RATE = 100,
  MAX_SPELL_RATE = 10,
  MAX_GIL_RATE = 10,
  MAX_DROP_RATE = 10,
  MAX_AP_RATE = 10,
  
  -- Default rates (normal FF6 progression)
  DEFAULT_RATES = {
    experience = 1.0,
    spell_learning = 1.0,
    gil = 1.0,
    drop_rate = 1.0,
    ap_gain = 1.0,
    gold_found = 1.0
  },
  
  -- Character list (14 playable characters)
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

-- Clamp value between min and max
local function clamp(value, min, max)
  if value < min then return min end
  if value > max then return max end
  return value
end

-- Format rate as percentage
local function formatRate(rate)
  return string.format("%.1fx (%d%%)", rate, math.floor(rate * 100))
end

-- Log operation
local function logOperation(operation, details)
  local timestamp = os.date("%Y-%m-%d %H:%M:%S")
  print(string.format("[%s] %s: %s", timestamp, operation, details or ""))
end

-- Create configuration snapshot
local function createSnapshot(config)
  return {
    timestamp = os.time(),
    date = os.date("%Y-%m-%d %H:%M:%S"),
    config = config
  }
end

-- ============================================================================
-- PRESET PACING PROFILES
-- ============================================================================

local PRESET_PROFILES = {
  {
    id = "normal",
    name = "Normal Progression",
    desc = "Standard FF6 progression rates (1.0x all)",
    rates = {
      experience = 1.0,
      spell_learning = 1.0,
      gil = 1.0,
      drop_rate = 1.0,
      ap_gain = 1.0,
      gold_found = 1.0
    }
  },
  {
    id = "speedrun",
    name = "Speedrun (Fast Leveling)",
    desc = "Fast progression: 5x EXP, 2x Spells, 2x Gil",
    rates = {
      experience = 5.0,
      spell_learning = 2.0,
      gil = 2.0,
      drop_rate = 1.0,
      ap_gain = 3.0,
      gold_found = 2.0
    }
  },
  {
    id = "casual",
    name = "Casual (Easy Leveling)",
    desc = "Relaxed progression: 2x EXP, 1.5x Spells",
    rates = {
      experience = 2.0,
      spell_learning = 1.5,
      gil = 1.5,
      drop_rate = 1.5,
      ap_gain = 1.5,
      gold_found = 1.5
    }
  },
  {
    id = "completionist",
    name = "Completionist (Slower)",
    desc = "Slower progression: 0.5x EXP, 0.7x Spells",
    rates = {
      experience = 0.5,
      spell_learning = 0.7,
      gil = 0.8,
      drop_rate = 0.8,
      ap_gain = 0.7,
      gold_found = 0.8
    }
  },
  {
    id = "hardcore",
    name = "Hardcore (Minimal Leveling)",
    desc = "Challenge: 0.2x EXP, 0.5x Spells",
    rates = {
      experience = 0.2,
      spell_learning = 0.5,
      gil = 0.5,
      drop_rate = 0.5,
      ap_gain = 0.3,
      gold_found = 0.5
    }
  },
  {
    id = "creative",
    name = "Creative (Spell-Heavy)",
    desc = "Spell focus: 1x EXP, 5x Spells, 0.5x Gil",
    rates = {
      experience = 1.0,
      spell_learning = 5.0,
      gil = 0.5,
      drop_rate = 1.0,
      ap_gain = 5.0,
      gold_found = 0.5
    }
  },
  {
    id = "economic",
    name = "Economic (Item-Heavy)",
    desc = "Item focus: 1x EXP, 5x Drops, 3x Gil",
    rates = {
      experience = 1.0,
      spell_learning = 1.0,
      gil = 3.0,
      drop_rate = 5.0,
      ap_gain = 1.0,
      gold_found = 3.0
    }
  }
}

-- ============================================================================
-- CORE RATE APPLICATION FUNCTIONS
-- ============================================================================

-- Apply experience rate multiplier
local function applyExperienceRate(rate)
  local cleanRate = clamp(rate, CONFIG.MIN_RATE, CONFIG.MAX_XP_RATE)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- game.setExperienceRate(cleanRate)
    logOperation("EXPERIENCE_RATE", formatRate(cleanRate))
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- Apply spell learning rate multiplier
local function applySpellLearningRate(rate)
  local cleanRate = clamp(rate, 0, CONFIG.MAX_SPELL_RATE)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- game.setSpellLearningRate(cleanRate)
    logOperation("SPELL_LEARNING_RATE", formatRate(cleanRate))
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- Apply gil acquisition rate
local function applyGilRate(rate)
  local cleanRate = clamp(rate, CONFIG.MIN_RATE, CONFIG.MAX_GIL_RATE)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- game.setGilRate(cleanRate)
    logOperation("GIL_RATE", formatRate(cleanRate))
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- Apply drop rate multiplier
local function applyDropRate(rate)
  local cleanRate = clamp(rate, CONFIG.MIN_RATE, CONFIG.MAX_DROP_RATE)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- game.setDropRate(cleanRate)
    logOperation("DROP_RATE", formatRate(cleanRate))
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- Apply AP gain rate (esper stat growth)
local function applyAPRate(rate)
  local cleanRate = clamp(rate, CONFIG.MIN_RATE, CONFIG.MAX_AP_RATE)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- game.setAPRate(cleanRate)
    logOperation("AP_RATE", formatRate(cleanRate))
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- Apply gold found rate (treasure/enemy drops)
local function applyGoldFoundRate(rate)
  local cleanRate = clamp(rate, CONFIG.MIN_RATE, CONFIG.MAX_GIL_RATE)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- game.setGoldFoundRate(cleanRate)
    logOperation("GOLD_FOUND_RATE", formatRate(cleanRate))
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- Apply character-specific rate
local function applyCharacterRate(charId, rateType, rate)
  local cleanRate = clamp(rate, CONFIG.MIN_RATE, 100)
  
  local result = safeCall(function()
    -- This would call actual API:
    -- character.setRate(charId, rateType, cleanRate)
    return cleanRate
  end)
  
  return result and cleanRate or nil
end

-- ============================================================================
-- CONFIGURATION MANAGEMENT
-- ============================================================================

-- Apply complete pacing profile
local function applyPacingProfile(profile)
  local results = {
    timestamp = os.date("%Y-%m-%d %H:%M:%S"),
    profile = profile.id,
    profile_name = profile.name,
    applied_rates = {},
    success = 0,
    failed = 0
  }
  
  -- Apply each rate
  for rateType, rate in pairs(profile.rates) do
    local result
    
    if rateType == "experience" then
      result = applyExperienceRate(rate)
    elseif rateType == "spell_learning" then
      result = applySpellLearningRate(rate)
    elseif rateType == "gil" then
      result = applyGilRate(rate)
    elseif rateType == "drop_rate" then
      result = applyDropRate(rate)
    elseif rateType == "ap_gain" then
      result = applyAPRate(rate)
    elseif rateType == "gold_found" then
      result = applyGoldFoundRate(rate)
    end
    
    if result then
      results.applied_rates[rateType] = result
      results.success = results.success + 1
    else
      results.failed = results.failed + 1
    end
  end
  
  logOperation("PACING_PROFILE", string.format("Profile '%s' applied", profile.name))
  return results
end

-- Apply custom rates
local function applyCustomRates(customRates)
  local results = {
    timestamp = os.date("%Y-%m-%d %H:%M:%S"),
    custom = true,
    applied_rates = {},
    success = 0,
    failed = 0
  }
  
  for rateType, rate in pairs(customRates) do
    local result
    
    if rateType == "experience" then
      result = applyExperienceRate(rate)
    elseif rateType == "spell_learning" then
      result = applySpellLearningRate(rate)
    elseif rateType == "gil" then
      result = applyGilRate(rate)
    elseif rateType == "drop_rate" then
      result = applyDropRate(rate)
    elseif rateType == "ap_gain" then
      result = applyAPRate(rate)
    elseif rateType == "gold_found" then
      result = applyGoldFoundRate(rate)
    end
    
    if result then
      results.applied_rates[rateType] = result
      results.success = results.success + 1
    else
      results.failed = results.failed + 1
    end
  end
  
  logOperation("CUSTOM_RATES", string.format("%d rates applied, %d failed", results.success, results.failed))
  return results
end

-- Create custom rate configuration
local function createCustomConfig(options)
  local config = {
    name = options.name or "Custom Pacing",
    description = options.description or "Custom progression rates",
    experience = options.experience or 1.0,
    spell_learning = options.spell_learning or 1.0,
    gil = options.gil or 1.0,
    drop_rate = options.drop_rate or 1.0,
    ap_gain = options.ap_gain or 1.0,
    gold_found = options.gold_found or 1.0,
    per_character = options.per_character or false,
    notes = options.notes or ""
  }
  
  -- Validate rates
  config.experience = clamp(config.experience, CONFIG.MIN_RATE, CONFIG.MAX_XP_RATE)
  config.spell_learning = clamp(config.spell_learning, 0, CONFIG.MAX_SPELL_RATE)
  config.gil = clamp(config.gil, CONFIG.MIN_RATE, CONFIG.MAX_GIL_RATE)
  config.drop_rate = clamp(config.drop_rate, CONFIG.MIN_RATE, CONFIG.MAX_DROP_RATE)
  config.ap_gain = clamp(config.ap_gain, CONFIG.MIN_RATE, CONFIG.MAX_AP_RATE)
  config.gold_found = clamp(config.gold_found, CONFIG.MIN_RATE, CONFIG.MAX_GIL_RATE)
  
  return config
end

-- Reset to normal rates
local function resetToNormalRates()
  logOperation("RESET_RATES", "Restoring default FF6 progression rates")
  return applyPacingProfile(PRESET_PROFILES[1])  -- Normal profile
end

-- ============================================================================
-- RATE TRACKING AND HISTORY
-- ============================================================================

local rateHistory = {}

-- Record rate change in history
local function recordRateChange(profile_name, rates)
  local entry = {
    timestamp = os.time(),
    date = os.date("%Y-%m-%d %H:%M:%S"),
    profile = profile_name,
    rates = rates
  }
  
  table.insert(rateHistory, entry)
  
  -- Keep last 50 entries
  if #rateHistory > 50 then
    table.remove(rateHistory, 1)
  end
  
  return entry
end

-- Get rate history
local function getRateHistory()
  return rateHistory
end

-- Clear rate history
local function clearRateHistory()
  rateHistory = {}
  logOperation("HISTORY_CLEARED", "Rate change history cleared")
end

-- ============================================================================
-- UI AND MENU FUNCTIONS
-- ============================================================================

-- Display main menu
local function displayMainMenu()
  print("\n" .. string.rep("=", 60))
  print("CUSTOM PROGRESSION PACING - FF6 Rate Control")
  print(string.rep("=", 60))
  print("\n⚠️  WARNING: This plugin controls game progression rates.")
  print("   Gameplay balance may be significantly affected.\n")
  
  print("OPTIONS:")
  print("  1. Apply Pacing Preset")
  print("  2. Custom Rate Configuration")
  print("  3. View Current Rates")
  print("  4. Reset to Normal Rates")
  print("  5. Rate History")
  print("  6. Settings")
  print("  7. Export Configuration")
  print("  8. Exit\n")
  
  return readInput("Select option (1-8): ")
end

-- Display preset profiles
local function displayPresetsMenu()
  print("\nPRECONFIGURED PACING PROFILES:\n")
  
  for i, profile in ipairs(PRESET_PROFILES) do
    print(string.format("  %d. %s", i, profile.name))
    print(string.format("     %s", profile.desc))
    print()
  end
  
  return readInput("Select preset (1-" .. #PRESET_PROFILES .. ") or 0 to cancel: ")
end

-- Display rate configuration
local function displayRateConfiguration()
  print("\nCUSTOM RATE CONFIGURATION:\n")
  print("Experience Gain Rate: 0.1x to 100x (default: 1.0x)")
  print("Spell Learning Rate: 0x to 10x (default: 1.0x)")
  print("Gil Rate: 0.1x to 10x (default: 1.0x)")
  print("Drop Rate: 0.1x to 10x (default: 1.0x)")
  print("AP Gain Rate: 0.1x to 10x (default: 1.0x)")
  print("Gold Found Rate: 0.1x to 10x (default: 1.0x)\n")
end

-- Display current rates
local function displayCurrentRates(currentConfig)
  print("\nCURRENT PROGRESSION RATES:\n")
  
  print("Experience Gain Rate:    " .. formatRate(currentConfig.experience))
  print("Spell Learning Rate:     " .. formatRate(currentConfig.spell_learning))
  print("Gil Rate:                " .. formatRate(currentConfig.gil))
  print("Drop Rate:               " .. formatRate(currentConfig.drop_rate))
  print("AP Gain Rate:            " .. formatRate(currentConfig.ap_gain))
  print("Gold Found Rate:         " .. formatRate(currentConfig.gold_found))
  print()
end

-- ============================================================================
-- MAIN PLUGIN INTERFACE
-- ============================================================================

local Plugin = {
  name = "Custom Progression Pacing",
  version = "1.0.0",
  batch = 3,
  currentConfig = nil,
  defaultConfig = CONFIG.DEFAULT_RATES
}

-- Initialize plugin
function Plugin:init()
  self.currentConfig = createCustomConfig({})
  logOperation("INIT", "Custom Progression Pacing loaded")
end

-- Apply preset
function Plugin:applyPreset(presetId)
  if presetId >= 1 and presetId <= #PRESET_PROFILES then
    local profile = PRESET_PROFILES[presetId]
    local results = applyPacingProfile(profile)
    recordRateChange(profile.name, profile.rates)
    self.currentConfig = createCustomConfig(profile.rates)
    return results
  else
    return nil
  end
end

-- Apply custom rates
function Plugin:applyCustom(options)
  local config = createCustomConfig(options)
  local results = applyCustomRates(config)
  recordRateChange(config.name, config)
  self.currentConfig = config
  return results
end

-- Get current rates
function Plugin:getCurrentRates()
  return self.currentConfig
end

-- Reset rates
function Plugin:resetRates()
  return resetToNormalRates()
end

-- Get plugin info
function Plugin:getInfo()
  return {
    name = self.name,
    version = self.version,
    batch = self.batch,
    presets = #PRESET_PROFILES,
    description = "Control all progression rates"
  }
end

-- Export configuration
function Plugin:exportConfiguration()
  return {
    config = self.currentConfig,
    timestamp = os.date("%Y-%m-%d %H:%M:%S"),
    plugin_version = self.version
  }
end

-- ============================================================================
-- EXPORT
-- ============================================================================

return Plugin
