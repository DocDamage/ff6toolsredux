--[[
  New Game Plus Generator Plugin
  Generate NG+ save files with configurable carryover content.
  
  Features:
  - Configurable carryover (items, equipment, levels, spells, espers)
  - Enemy stat scaling (1x-10x multiplier)
  - All characters unlocked option
  - World state selection (WoB/WoR)
  - Multiple NG+ profiles
  - Preset difficulty tiers
  - Export NG+ configuration
  
  Phase: 6 (Gameplay-Altering Plugins)
  Version: 1.0.0
  
  WARNING: EXPERIMENTAL - Creates modified save files
]]

local function safeCall(fn, ...) local success, result = pcall(fn, ...) return success and result or nil end

local CONFIG_FILE = "new_game_plus/profiles.lua"

local DIFFICULTY_PRESETS = {
  {name="Easy", desc="Keep everything", carryover={items=true,equipment=true,levels=true,spells=true,espers=true}, scaling=1.0},
  {name="Normal", desc="Keep levels + espers", carryover={items=false,equipment=false,levels=true,spells=false,espers=true}, scaling=1.5},
  {name="Hard", desc="Keep espers only", carryover={items=false,equipment=false,levels=false,spells=false,espers=true}, scaling=2.0},
  {name="Nightmare", desc="Nothing carries over", carryover={items=false,equipment=false,levels=false,spells=false,espers=false}, scaling=3.0}
}

local function serialize(tbl, indent)
  indent = indent or 0
  local spacing = string.rep("  ", indent)
  local lines = {"{"}
  for k, v in pairs(tbl) do
    local key = type(k) == "number" and string.format("[%d]", k) or string.format('["%s"]', k)
    local value = type(v) == "table" and serialize(v, indent + 1) or (type(v) == "string" and string.format('"%s"', v:gsub('"', '\\"')) or tostring(v))
    table.insert(lines, spacing .. "  " .. key .. " = " .. value .. ",")
  end
  table.insert(lines, spacing .. "}")
  return table.concat(lines, "\n")
end

local function loadProfiles()
  local content = safeCall(ReadFile, CONFIG_FILE)
  if not content then return {profiles = {}} end
  local success, data = pcall(load, "return " .. content)
  return (success and data and data()) or {profiles = {}}
end

local function saveProfiles(data)
  safeCall(WriteFile, CONFIG_FILE, "return " .. serialize(data))
end

local function createNGPlusProfile()
  local profile_name = ShowDialog("Enter Profile Name:", {"Enter a name for this NG+ profile:"})
  if not profile_name or profile_name == "" then return "Cancelled." end
  
  -- Carryover configuration
  local lines = {"Configure what carries over to NG+:"}
  local carryover = {}
  
  local items_choice = ShowDialog("Carry over Items?", {"1. Yes", "2. No"})
  carryover.items = (items_choice == 1)
  
  local equipment_choice = ShowDialog("Carry over Equipment?", {"1. Yes", "2. No"})
  carryover.equipment = (equipment_choice == 1)
  
  local levels_choice = ShowDialog("Carry over Levels?", {"1. Yes", "2. No"})
  carryover.levels = (levels_choice == 1)
  
  local spells_choice = ShowDialog("Carry over Spells?", {"1. Yes", "2. No"})
  carryover.spells = (spells_choice == 1)
  
  local espers_choice = ShowDialog("Carry over Espers?", {"1. Yes", "2. No"})
  carryover.espers = (espers_choice == 1)
  
  -- Enemy scaling
  local scaling_opts = {}
  for i = 1, 10 do table.insert(scaling_opts, i .. ". " .. i .. "x scaling") end
  local scaling_choice = ShowDialog("Enemy Stat Scaling:", scaling_opts)
  local scaling = scaling_choice or 1
  
  -- World state
  local world_choice = ShowDialog("Starting World State:", {"1. World of Balance (WoB)", "2. World of Ruin (WoR)"})
  local world_state = (world_choice == 1) and "WoB" or "WoR"
  
  -- All characters unlocked
  local unlock_choice = ShowDialog("Unlock All Characters?", {"1. Yes", "2. No"})
  local unlock_all = (unlock_choice == 1)
  
  -- Save profile
  local data = loadProfiles()
  table.insert(data.profiles, {
    name = profile_name,
    carryover = carryover,
    scaling = scaling,
    world_state = world_state,
    unlock_all = unlock_all,
    created = os.date("%Y-%m-%d %H:%M:%S")
  })
  saveProfiles(data)
  
  return string.format("NG+ Profile Created: %s\n\nCarryover:\n  Items: %s\n  Equipment: %s\n  Levels: %s\n  Spells: %s\n  Espers: %s\n\nScaling: %dx\nWorld: %s\nAll Characters: %s",
    profile_name,
    carryover.items and "Yes" or "No",
    carryover.equipment and "Yes" or "No",
    carryover.levels and "Yes" or "No",
    carryover.spells and "Yes" or "No",
    carryover.espers and "Yes" or "No",
    scaling,
    world_state,
    unlock_all and "Yes" or "No")
end

local function createFromPreset()
  local preset_opts = {}
  for i, p in ipairs(DIFFICULTY_PRESETS) do
    table.insert(preset_opts, i .. ". " .. p.name .. " - " .. p.desc)
  end
  table.insert(preset_opts, #DIFFICULTY_PRESETS + 1 .. ". Cancel")
  
  local choice = ShowDialog("Select Preset:", preset_opts)
  if not choice or choice > #DIFFICULTY_PRESETS then return "Cancelled." end
  
  local preset = DIFFICULTY_PRESETS[choice]
  local profile_name = preset.name .. " NG+ " .. os.date("%Y%m%d_%H%M%S")
  
  local data = loadProfiles()
  table.insert(data.profiles, {
    name = profile_name,
    carryover = preset.carryover,
    scaling = preset.scaling,
    world_state = "WoB",
    unlock_all = false,
    created = os.date("%Y-%m-%d %H:%M:%S"),
    preset = preset.name
  })
  saveProfiles(data)
  
  return string.format("Created %s profile: %s", preset.name, profile_name)
end

local function viewProfiles()
  local data = loadProfiles()
  local lines = {"=== New Game Plus Profiles ===", ""}
  
  if #data.profiles == 0 then
    table.insert(lines, "No profiles created yet.")
  else
    for i, profile in ipairs(data.profiles) do
      table.insert(lines, string.format("%d. %s", i, profile.name))
      table.insert(lines, string.format("   Created: %s", profile.created))
      table.insert(lines, string.format("   Scaling: %dx | World: %s", profile.scaling, profile.world_state))
      local carry_parts = {}
      for k, v in pairs(profile.carryover) do if v then table.insert(carry_parts, k) end end
      table.insert(lines, string.format("   Carryover: %s", #carry_parts > 0 and table.concat(carry_parts, ", ") or "none"))
      table.insert(lines, "")
    end
  end
  
  return table.concat(lines, "\n")
end

local function generateNGPlus()
  local data = loadProfiles()
  if #data.profiles == 0 then return "No profiles created. Create a profile first." end
  
  local profile_opts = {}
  for i, p in ipairs(data.profiles) do
    table.insert(profile_opts, i .. ". " .. p.name)
  end
  table.insert(profile_opts, #data.profiles + 1 .. ". Cancel")
  
  local choice = ShowDialog("Select Profile to Generate:", profile_opts)
  if not choice or choice > #data.profiles then return "Cancelled." end
  
  local profile = data.profiles[choice]
  
  -- This is a conceptual implementation - actual save generation requires advanced APIs
  local filename = "new_game_plus/ngplus_" .. os.date("%Y%m%d_%H%M%S") .. ".sav"
  local config_text = serialize(profile)
  
  safeCall(WriteFile, filename .. ".config", config_text)
  
  return string.format("NG+ Generation Request Created\n\nProfile: %s\nConfig saved to: %s\n\nNote: Actual save generation requires save file creation APIs.\nUse config file to manually apply settings.", profile.name, filename .. ".config")
end

local function exportProfiles()
  local data = loadProfiles()
  local filename = "new_game_plus/export_" .. os.date("%Y%m%d_%H%M%S") .. ".txt"
  
  local lines = {"=== FF6 New Game Plus Profiles ===", "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"), ""}
  
  for i, profile in ipairs(data.profiles) do
    table.insert(lines, string.format("Profile %d: %s", i, profile.name))
    table.insert(lines, string.format("  Created: %s", profile.created))
    if profile.preset then table.insert(lines, string.format("  Preset: %s", profile.preset)) end
    table.insert(lines, string.format("  Enemy Scaling: %dx", profile.scaling))
    table.insert(lines, string.format("  World State: %s", profile.world_state))
    table.insert(lines, string.format("  All Characters: %s", profile.unlock_all and "Yes" or "No"))
    table.insert(lines, "  Carryover:")
    for k, v in pairs(profile.carryover) do
      table.insert(lines, string.format("    %s: %s", k, v and "Yes" or "No"))
    end
    table.insert(lines, "")
  end
  
  safeCall(WriteFile, filename, table.concat(lines, "\n"))
  return "Exported profiles to:\n" .. filename
end

local function main()
  ShowDialog("WARNING", {"EXPERIMENTAL PLUGIN", "NG+ generation requires advanced save file creation APIs.", "This plugin creates configuration files for NG+ profiles.", "Actual save generation may require manual steps.", "", "Press OK to continue..."})
  
  while true do
    local choice = ShowDialog("New Game Plus Generator", {"1. Create Custom Profile", "2. Create from Preset", "3. View Profiles", "4. Generate NG+", "5. Export Profiles", "6. Exit"})
    if not choice or choice == 6 then break end
    
    local result = ""
    if choice == 1 then result = createNGPlusProfile()
    elseif choice == 2 then result = createFromPreset()
    elseif choice == 3 then result = viewProfiles()
    elseif choice == 4 then result = generateNGPlus()
    elseif choice == 5 then result = exportProfiles()
    end
    
    if result ~= "" then ShowDialog("Result", {result, "Press any key..."}) end
  end
end

main()
