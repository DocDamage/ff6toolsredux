--[[
  Character Ability Swap Plugin
  Give any character any special command ability.
  
  Features:
  - Swap command abilities between characters
  - 14 character command ability library
  - Mix and match up to 4 commands per character
  - Command preview with descriptions
  - Save/load custom ability sets
  - Reset to default abilities
  - Preset builds
  - Export configurations
  
  Phase: 6 (Gameplay-Altering Plugins)
  Version: 1.0.0
  
  WARNING: This is an EXPERIMENTAL plugin that fundamentally alters gameplay.
  API support for command ability swapping may be limited or unavailable.
  This plugin demonstrates the concept with metadata tracking.
]]

-- Safe API call wrapper
local function safeCall(fn, ...)
  local success, result = pcall(fn, ...)
  if success then
    return result
  end
  return nil
end

-- Command ability library (all 14 character commands)
local COMMAND_ABILITIES = {
  {id = 1, name = "Fight", char = "All", desc = "Basic attack command"},
  {id = 2, name = "Magic", char = "Terra/Celes", desc = "Cast learned magic spells"},
  {id = 3, name = "Morph", char = "Terra", desc = "Transform into Esper form (Terra only)"},
  {id = 4, name = "Steal", char = "Locke", desc = "Steal items from enemies"},
  {id = 5, name = "Capture", char = "Locke", desc = "Capture enemy (Colosseum)"},
  {id = 6, name = "SwdTech", char = "Cyan", desc = "Execute Bushido sword techniques"},
  {id = 7, name = "Throw", char = "Shadow", desc = "Throw weapons and items"},
  {id = 8, name = "Tools", char = "Edgar", desc = "Use Edgar's mechanical tools"},
  {id = 9, name = "Blitz", char = "Sabin", desc = "Execute martial arts techniques"},
  {id = 10, name = "Runic", char = "Celes", desc = "Absorb magic spells as MP"},
  {id = 11, name = "Lore", char = "Strago", desc = "Cast learned enemy Lores"},
  {id = 12, name = "Sketch", char = "Relm", desc = "Mimic enemy appearance and attacks"},
  {id = 13, name = "Control", char = "Relm", desc = "Take control of enemy"},
  {id = 14, name = "Slot", char = "Setzer", desc = "Use slot machine attacks"},
  {id = 15, name = "Rage", char = "Gau", desc = "Use learned enemy Rages"},
  {id = 16, name = "Leap", char = "Gau", desc = "Jump onto Veldt"},
  {id = 17, name = "Mimic", char = "Gogo", desc = "Repeat last action"},
  {id = 18, name = "Dance", char = "Mog", desc = "Perform learned Dances"}
}

-- Default command setups (character defaults)
local DEFAULT_COMMANDS = {
  Terra = {1, 2, 3, 0}, -- Fight, Magic, Morph, (Item)
  Locke = {1, 4, 5, 0}, -- Fight, Steal, Capture, (Item)
  Cyan = {1, 6, 0, 0},  -- Fight, SwdTech, (Item), (empty)
  Shadow = {1, 7, 0, 0}, -- Fight, Throw, (Item), (empty)
  Edgar = {1, 8, 0, 0},  -- Fight, Tools, (Item), (empty)
  Sabin = {1, 9, 0, 0},  -- Fight, Blitz, (Item), (empty)
  Celes = {1, 2, 10, 0}, -- Fight, Magic, Runic, (Item)
  Strago = {1, 2, 11, 0}, -- Fight, Magic, Lore, (Item)
  Relm = {1, 12, 13, 0}, -- Fight, Sketch, Control, (Item)
  Setzer = {1, 14, 0, 0}, -- Fight, Slot, (Item), (empty)
  Mog = {1, 2, 18, 0},   -- Fight, Magic, Dance, (Item)
  Gau = {1, 15, 16, 0},  -- Fight, Rage, Leap, (Item)
  Gogo = {1, 17, 0, 0},  -- Fight, Mimic, (Item), (empty) - can set any
  Umaro = {1, 0, 0, 0}   -- Fight only (auto-battle)
}

-- Preset ability builds
local PRESET_BUILDS = {
  {
    id = "all_magic",
    name = "All-Magic Party",
    desc = "Every character has Magic command",
    setup = {
      Terra = {1, 2, 3, 0},
      Locke = {1, 2, 4, 0},
      Cyan = {1, 2, 6, 0},
      Shadow = {1, 2, 7, 0},
      Edgar = {1, 2, 8, 0},
      Sabin = {1, 2, 9, 0},
      Celes = {1, 2, 10, 0},
      Strago = {1, 2, 11, 0},
      Relm = {1, 2, 12, 0},
      Setzer = {1, 2, 14, 0},
      Mog = {1, 2, 18, 0},
      Gau = {1, 2, 15, 0},
      Gogo = {1, 2, 17, 0},
      Umaro = {1, 2, 0, 0}
    }
  },
  {
    id = "all_physical",
    name = "All-Physical Party",
    desc = "Focus on physical abilities only",
    setup = {
      Terra = {1, 9, 7, 0},  -- Fight, Blitz, Throw
      Locke = {1, 4, 9, 0},  -- Fight, Steal, Blitz
      Cyan = {1, 6, 9, 0},   -- Fight, SwdTech, Blitz
      Shadow = {1, 7, 9, 0}, -- Fight, Throw, Blitz
      Edgar = {1, 8, 9, 0},  -- Fight, Tools, Blitz
      Sabin = {1, 9, 7, 0},  -- Fight, Blitz, Throw
      Celes = {1, 6, 9, 0},  -- Fight, SwdTech, Blitz
      Strago = {1, 9, 7, 0}, -- Fight, Blitz, Throw
      Relm = {1, 7, 9, 0},   -- Fight, Throw, Blitz
      Setzer = {1, 14, 9, 0}, -- Fight, Slot, Blitz
      Mog = {1, 9, 7, 0},    -- Fight, Blitz, Throw
      Gau = {1, 15, 9, 0},   -- Fight, Rage, Blitz
      Gogo = {1, 9, 17, 0},  -- Fight, Blitz, Mimic
      Umaro = {1, 9, 0, 0}   -- Fight, Blitz
    }
  },
  {
    id = "utility",
    name = "Utility Focused",
    desc = "Steal, Control, and support abilities",
    setup = {
      Terra = {1, 2, 4, 0},
      Locke = {1, 4, 13, 0},
      Cyan = {1, 4, 13, 0},
      Shadow = {1, 4, 7, 0},
      Edgar = {1, 8, 4, 0},
      Sabin = {1, 4, 13, 0},
      Celes = {1, 2, 10, 0},
      Strago = {1, 2, 11, 0},
      Relm = {1, 12, 13, 0},
      Setzer = {1, 4, 14, 0},
      Mog = {1, 2, 18, 0},
      Gau = {1, 15, 4, 0},
      Gogo = {1, 17, 4, 0},
      Umaro = {1, 4, 0, 0}
    }
  },
  {
    id = "chaos",
    name = "Chaos Mode",
    desc = "Random unusual ability combinations",
    setup = {
      Terra = {1, 15, 12, 0},  -- Fight, Rage, Sketch
      Locke = {1, 18, 11, 0},  -- Fight, Dance, Lore
      Cyan = {1, 14, 17, 0},   -- Fight, Slot, Mimic
      Shadow = {1, 12, 13, 0}, -- Fight, Sketch, Control
      Edgar = {1, 15, 18, 0},  -- Fight, Rage, Dance
      Sabin = {1, 2, 14, 0},   -- Fight, Magic, Slot
      Celes = {1, 9, 15, 0},   -- Fight, Blitz, Rage
      Strago = {1, 6, 14, 0},  -- Fight, SwdTech, Slot
      Relm = {1, 15, 18, 0},   -- Fight, Rage, Dance
      Setzer = {1, 9, 11, 0},  -- Fight, Blitz, Lore
      Mog = {1, 11, 12, 0},    -- Fight, Lore, Sketch
      Gau = {1, 2, 18, 0},     -- Fight, Magic, Dance
      Gogo = {1, 8, 17, 0},    -- Fight, Tools, Mimic
      Umaro = {1, 15, 0, 0}    -- Fight, Rage
    }
  }
}

-- Storage file
local CONFIG_FILE = "ability_swap/configurations.lua"

-- Serialize table
local function serialize(tbl, indent)
  indent = indent or 0
  local spacing = string.rep("  ", indent)
  local lines = {"{"}
  
  for k, v in pairs(tbl) do
    local key = type(k) == "number" and string.format("[%d]", k) or string.format('["%s"]', k)
    local value
    if type(v) == "table" then
      value = serialize(v, indent + 1)
    elseif type(v) == "string" then
      value = string.format('"%s"', v:gsub('"', '\\"'))
    elseif type(v) == "boolean" then
      value = tostring(v)
    else
      value = tostring(v)
    end
    table.insert(lines, spacing .. "  " .. key .. " = " .. value .. ",")
  end
  
  table.insert(lines, spacing .. "}")
  return table.concat(lines, "\n")
end

-- Load configuration
local function loadConfig()
  local content = safeCall(ReadFile, CONFIG_FILE)
  if not content then return {active_setup = nil, saved_setups = {}} end
  
  local success, data = pcall(load, "return " .. content)
  if success and data then
    return data()
  end
  return {active_setup = nil, saved_setups = {}}
end

-- Save configuration
local function saveConfig(config)
  local content = "return " .. serialize(config)
  safeCall(WriteFile, CONFIG_FILE, content)
end

-- Get ability name by ID
local function getAbilityName(id)
  if id == 0 then return "(None)" end
  for _, cmd in ipairs(COMMAND_ABILITIES) do
    if cmd.id == id then return cmd.name end
  end
  return "Unknown"
end

-- Get ability by ID
local function getAbilityById(id)
  for _, cmd in ipairs(COMMAND_ABILITIES) do
    if cmd.id == id then return cmd end
  end
  return nil
end

-- View current setup
local function viewCurrentSetup()
  local config = loadConfig()
  local setup = config.active_setup or DEFAULT_COMMANDS
  
  local lines = {}
  table.insert(lines, "=== Current Ability Setup ===\n")
  
  if config.active_setup then
    table.insert(lines, "Custom setup active\n")
  else
    table.insert(lines, "Default setup (vanilla game)\n")
  end
  
  local characters = {"Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", 
                      "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"}
  
  for _, char in ipairs(characters) do
    local commands = setup[char] or DEFAULT_COMMANDS[char]
    if commands then
      table.insert(lines, string.format("%s:", char))
      for slot = 1, 4 do
        local cmd_id = commands[slot] or 0
        table.insert(lines, string.format("  Slot %d: %s", slot, getAbilityName(cmd_id)))
      end
      table.insert(lines, "")
    end
  end
  
  return table.concat(lines, "\n")
end

-- Edit character abilities
local function editCharacterAbilities()
  local config = loadConfig()
  local setup = config.active_setup or {}
  
  -- Copy defaults if no active setup
  if not config.active_setup then
    for char, commands in pairs(DEFAULT_COMMANDS) do
      setup[char] = {commands[1], commands[2], commands[3], commands[4]}
    end
  end
  
  -- Select character
  local char_options = {}
  local characters = {"Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", 
                      "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"}
  for i, char in ipairs(characters) do
    table.insert(char_options, string.format("%d. %s", i, char))
  end
  table.insert(char_options, string.format("%d. Cancel", #char_options + 1))
  
  local char_choice = ShowDialog("Select Character:", char_options)
  if not char_choice or char_choice > #characters then return "Cancelled." end
  
  local character = characters[char_choice]
  local commands = setup[character] or {1, 0, 0, 0}
  
  -- Select slot
  local slot_choice = ShowDialog(string.format("Edit %s - Select Slot:", character), {
    string.format("1. Slot 1: %s", getAbilityName(commands[1])),
    string.format("2. Slot 2: %s", getAbilityName(commands[2])),
    string.format("3. Slot 3: %s", getAbilityName(commands[3])),
    string.format("4. Slot 4: %s", getAbilityName(commands[4])),
    "5. Cancel"
  })
  
  if not slot_choice or slot_choice == 5 then return "Cancelled." end
  
  -- Select ability
  local ability_options = {}
  for i, cmd in ipairs(COMMAND_ABILITIES) do
    table.insert(ability_options, string.format("%d. %s (%s) - %s", i, cmd.name, cmd.char, cmd.desc))
  end
  table.insert(ability_options, string.format("%d. (None/Empty)", #ability_options + 1))
  table.insert(ability_options, string.format("%d. Cancel", #ability_options + 2))
  
  local ability_choice = ShowDialog("Select Ability:", ability_options)
  if not ability_choice or ability_choice > #COMMAND_ABILITIES + 1 then return "Cancelled." end
  
  -- Set ability
  local new_cmd_id = ability_choice <= #COMMAND_ABILITIES and COMMAND_ABILITIES[ability_choice].id or 0
  commands[slot_choice] = new_cmd_id
  setup[character] = commands
  
  config.active_setup = setup
  saveConfig(config)
  
  return string.format("Set %s Slot %d to: %s", character, slot_choice, getAbilityName(new_cmd_id))
end

-- Apply preset build
local function applyPresetBuild()
  local config = loadConfig()
  
  local preset_options = {}
  for i, preset in ipairs(PRESET_BUILDS) do
    table.insert(preset_options, string.format("%d. %s - %s", i, preset.name, preset.desc))
  end
  table.insert(preset_options, string.format("%d. Cancel", #preset_options + 1))
  
  local choice = ShowDialog("Select Preset Build:", preset_options)
  if not choice or choice > #PRESET_BUILDS then return "Cancelled." end
  
  local preset = PRESET_BUILDS[choice]
  config.active_setup = preset.setup
  saveConfig(config)
  
  return string.format("Applied preset: %s\n%s", preset.name, preset.desc)
end

-- Reset to defaults
local function resetToDefaults()
  local confirm = ShowDialog("Reset all abilities to defaults?", {
    "1. Yes, reset to vanilla",
    "2. No, keep current setup"
  })
  
  if confirm == 1 then
    local config = loadConfig()
    config.active_setup = nil
    saveConfig(config)
    return "Reset to default abilities."
  end
  
  return "Reset cancelled."
end

-- Save current setup
local function saveCurrentSetup()
  local config = loadConfig()
  
  if not config.active_setup then
    return "No custom setup to save. Edit abilities first."
  end
  
  local name = ShowInput("Enter name for this ability setup:")
  if not name or name == "" then return "Cancelled." end
  
  table.insert(config.saved_setups, {
    name = name,
    setup = config.active_setup,
    saved = os.time()
  })
  
  saveConfig(config)
  return string.format("Saved setup: %s", name)
end

-- Load saved setup
local function loadSavedSetup()
  local config = loadConfig()
  
  if #config.saved_setups == 0 then
    return "No saved setups available."
  end
  
  local options = {}
  for i, saved in ipairs(config.saved_setups) do
    table.insert(options, string.format("%d. %s (saved %s)", i, saved.name, 
      os.date("%Y-%m-%d", saved.saved)))
  end
  table.insert(options, string.format("%d. Cancel", #options + 1))
  
  local choice = ShowDialog("Load Saved Setup:", options)
  if not choice or choice > #config.saved_setups then return "Cancelled." end
  
  local saved = config.saved_setups[choice]
  config.active_setup = saved.setup
  saveConfig(config)
  
  return string.format("Loaded setup: %s", saved.name)
end

-- Export configuration
local function exportConfiguration()
  local config = loadConfig()
  
  if not config.active_setup then
    return "No custom setup to export. Using defaults."
  end
  
  local timestamp = os.date("%Y%m%d_%H%M%S")
  local filename = string.format("ability_swap/export_%s.txt", timestamp)
  
  local lines = {}
  table.insert(lines, "=== FF6 Character Ability Swap Configuration ===")
  table.insert(lines, "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"))
  table.insert(lines, "")
  
  local characters = {"Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", 
                      "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"}
  
  for _, char in ipairs(characters) do
    local commands = config.active_setup[char]
    if commands then
      table.insert(lines, string.format("%s:", char))
      for slot = 1, 4 do
        table.insert(lines, string.format("  Slot %d: %s", slot, getAbilityName(commands[slot] or 0)))
      end
      table.insert(lines, "")
    end
  end
  
  table.insert(lines, "Generated by FF6 Save Editor - Character Ability Swap v1.0.0")
  
  local content = table.concat(lines, "\n")
  safeCall(WriteFile, filename, content)
  
  return "Configuration exported to: " .. filename
end

-- Main menu
local function main()
  ShowDialog("WARNING", {
    "This is an EXPERIMENTAL plugin.",
    "API support for command ability swapping may be unavailable.",
    "This plugin tracks configurations but cannot directly modify",
    "in-game command abilities without advanced API support.",
    "",
    "Press OK to continue..."
  })
  
  while true do
    local choice = ShowDialog("Character Ability Swap", {
      "1. View Current Setup",
      "2. Edit Character Abilities",
      "3. Apply Preset Build",
      "4. Reset to Defaults",
      "5. Save Current Setup",
      "6. Load Saved Setup",
      "7. Export Configuration",
      "8. Exit"
    })
    
    if not choice or choice == 8 then break end
    
    local result = ""
    if choice == 1 then
      result = viewCurrentSetup()
    elseif choice == 2 then
      result = editCharacterAbilities()
    elseif choice == 3 then
      result = applyPresetBuild()
    elseif choice == 4 then
      result = resetToDefaults()
    elseif choice == 5 then
      result = saveCurrentSetup()
    elseif choice == 6 then
      result = loadSavedSetup()
    elseif choice == 7 then
      result = exportConfiguration()
    end
    
    if result ~= "" then
      ShowDialog("Result", {result, "Press any key to continue..."})
    end
  end
end

-- Run plugin
main()
