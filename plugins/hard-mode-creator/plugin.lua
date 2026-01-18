--[[
  Hard Mode Creator Plugin
  Create custom difficulty modifiers for challenge runs.
  
  Features:
  - Enemy stat multipliers (HP, MP, Attack, Defense, Magic)
  - Player stat divisors (handicap mode)
  - Experience/Gil rate modification (0.1x-10x)
  - Item use restrictions
  - Equipment slot restrictions
  - Difficulty presets (Iron Man, Nuzlocke, etc.)
  - Export difficulty profiles
  
  Phase: 6 (Gameplay-Altering Plugins)
  Version: 1.0.0
  
  WARNING: EXPERIMENTAL - Fundamentally alters game difficulty
]]

local function safeCall(fn, ...) local success, result = pcall(fn, ...) return success and result or nil end

local CONFIG_FILE = "hard_mode/profiles.lua"

local DIFFICULTY_PRESETS = {
  {name="Iron Man", desc="Permanent death + restrictions", enemy_mult={HP=1.5,Attack=1.3,Defense=1.2}, player_div={HP=1.0,Attack=1.0}, exp_rate=0.5, gil_rate=0.5, restrictions={"No Phoenix Down","No Life spells"}},
  {name="Nuzlocke", desc="First encounter only + death rules", enemy_mult={HP=1.2,Attack=1.1}, player_div={HP=1.0}, exp_rate=0.75, gil_rate=1.0, restrictions={"Permadeath","First encounter recruitment"}},
  {name="Low Level", desc="Reduced XP, high difficulty", enemy_mult={HP=2.0,Attack=1.5,Defense=1.3,Magic=1.4}, player_div={HP=1.0}, exp_rate=0.1, gil_rate=1.0, restrictions={}},
  {name="Poverty", desc="Minimal gil and items", enemy_mult={HP=1.3,Attack=1.2}, player_div={HP=1.0}, exp_rate=1.0, gil_rate=0.1, restrictions={"No item shops"}},
  {name="Nightmare", desc="Maximum difficulty", enemy_mult={HP=3.0,Attack=2.0,Defense=1.8,Magic=2.0,MagicDef=1.5}, player_div={HP=2.0,Attack=2.0,Magic=2.0}, exp_rate=0.5, gil_rate=0.5, restrictions={"No Save Points","Limited items"}}
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
  if not content then return {profiles = {}, active = nil} end
  local success, data = pcall(load, "return " .. content)
  return (success and data and data()) or {profiles = {}, active = nil}
end

local function saveProfiles(data)
  safeCall(WriteFile, CONFIG_FILE, "return " .. serialize(data))
end

local function createCustomProfile()
  local profile_name = ShowDialog("Enter Profile Name:", {"Enter a name for this difficulty profile:"})
  if not profile_name or profile_name == "" then return "Cancelled." end
  
  -- Enemy multipliers
  local enemy_mult = {}
  local stats = {"HP", "MP", "Attack", "Defense", "Magic", "MagicDef"}
  
  for _, stat in ipairs(stats) do
    local mult_opts = {}
    for i = 1, 20 do table.insert(mult_opts, string.format("%d. %.1fx", i, i * 0.5)) end
    local choice = ShowDialog(string.format("Enemy %s Multiplier:", stat), mult_opts)
    enemy_mult[stat] = (choice or 2) * 0.5
  end
  
  -- Player divisors (handicap)
  local player_div = {}
  local player_stats = {"HP", "MP", "Attack", "Defense", "Magic"}
  
  for _, stat in ipairs(player_stats) do
    local div_opts = {"1. 1.0x (no handicap)", "2. 1.5x handicap", "3. 2.0x handicap", "4. 3.0x handicap"}
    local choice = ShowDialog(string.format("Player %s Divisor:", stat), div_opts)
    local divs = {1.0, 1.5, 2.0, 3.0}
    player_div[stat] = divs[choice or 1]
  end
  
  -- Experience/Gil rates
  local exp_opts = {}
  for i = 1, 20 do table.insert(exp_opts, string.format("%d. %.1fx", i, i * 0.5)) end
  local exp_choice = ShowDialog("Experience Rate:", exp_opts)
  local exp_rate = (exp_choice or 2) * 0.5
  
  local gil_choice = ShowDialog("Gil Rate:", exp_opts)
  local gil_rate = (gil_choice or 2) * 0.5
  
  -- Restrictions (simple list for now)
  local restrictions = {}
  
  -- Save profile
  local data = loadProfiles()
  table.insert(data.profiles, {
    name = profile_name,
    enemy_mult = enemy_mult,
    player_div = player_div,
    exp_rate = exp_rate,
    gil_rate = gil_rate,
    restrictions = restrictions,
    created = os.date("%Y-%m-%d %H:%M:%S")
  })
  saveProfiles(data)
  
  local lines = {string.format("Hard Mode Profile Created: %s", profile_name), ""}
  for stat, mult in pairs(enemy_mult) do
    table.insert(lines, string.format("Enemy %s: %.1fx", stat, mult))
  end
  table.insert(lines, "")
  for stat, div in pairs(player_div) do
    table.insert(lines, string.format("Player %s: %.1fx", stat, 1.0/div))
  end
  table.insert(lines, "")
  table.insert(lines, string.format("Experience: %.1fx", exp_rate))
  table.insert(lines, string.format("Gil: %.1fx", gil_rate))
  
  return table.concat(lines, "\n")
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
  local profile_name = preset.name .. " " .. os.date("%Y%m%d_%H%M%S")
  
  local data = loadProfiles()
  table.insert(data.profiles, {
    name = profile_name,
    enemy_mult = preset.enemy_mult,
    player_div = preset.player_div,
    exp_rate = preset.exp_rate,
    gil_rate = preset.gil_rate,
    restrictions = preset.restrictions,
    created = os.date("%Y-%m-%d %H:%M:%S"),
    preset = preset.name
  })
  saveProfiles(data)
  
  return string.format("Created %s profile: %s\n\nRestrictions: %s", preset.name, profile_name, table.concat(preset.restrictions, ", "))
end

local function viewProfiles()
  local data = loadProfiles()
  local lines = {"=== Hard Mode Profiles ===", ""}
  
  if #data.profiles == 0 then
    table.insert(lines, "No profiles created yet.")
  else
    for i, profile in ipairs(data.profiles) do
      table.insert(lines, string.format("%d. %s%s", i, profile.name, data.active == i and " [ACTIVE]" or ""))
      table.insert(lines, string.format("   Created: %s", profile.created))
      if profile.preset then table.insert(lines, string.format("   Preset: %s", profile.preset)) end
      table.insert(lines, string.format("   Exp: %.1fx | Gil: %.1fx", profile.exp_rate, profile.gil_rate))
      table.insert(lines, "")
    end
  end
  
  return table.concat(lines, "\n")
end

local function activateProfile()
  local data = loadProfiles()
  if #data.profiles == 0 then return "No profiles created. Create a profile first." end
  
  local profile_opts = {}
  for i, p in ipairs(data.profiles) do
    table.insert(profile_opts, i .. ". " .. p.name)
  end
  table.insert(profile_opts, #data.profiles + 1 .. ". Cancel")
  
  local choice = ShowDialog("Select Profile to Activate:", profile_opts)
  if not choice or choice > #data.profiles then return "Cancelled." end
  
  local profile = data.profiles[choice]
  data.active = choice
  saveProfiles(data)
  
  -- This is conceptual - actual stat modification requires game state APIs
  local lines = {string.format("Activated: %s", profile.name), "", "Enemy Multipliers:"}
  for stat, mult in pairs(profile.enemy_mult) do
    table.insert(lines, string.format("  %s: %.1fx", stat, mult))
  end
  table.insert(lines, "")
  table.insert(lines, "Player Handicaps:")
  for stat, div in pairs(profile.player_div) do
    table.insert(lines, string.format("  %s: %.1fx", stat, 1.0/div))
  end
  table.insert(lines, "")
  table.insert(lines, string.format("Experience: %.1fx", profile.exp_rate))
  table.insert(lines, string.format("Gil: %.1fx", profile.gil_rate))
  table.insert(lines, "")
  table.insert(lines, "Note: Actual stat modification requires game engine APIs.")
  
  return table.concat(lines, "\n")
end

local function exportProfiles()
  local data = loadProfiles()
  local filename = "hard_mode/export_" .. os.date("%Y%m%d_%H%M%S") .. ".txt"
  
  local lines = {"=== FF6 Hard Mode Difficulty Profiles ===", "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"), ""}
  
  for i, profile in ipairs(data.profiles) do
    table.insert(lines, string.format("Profile %d: %s", i, profile.name))
    table.insert(lines, string.format("  Created: %s", profile.created))
    if profile.preset then table.insert(lines, string.format("  Preset: %s", profile.preset)) end
    table.insert(lines, "  Enemy Multipliers:")
    for stat, mult in pairs(profile.enemy_mult) do
      table.insert(lines, string.format("    %s: %.1fx", stat, mult))
    end
    table.insert(lines, "  Player Divisors:")
    for stat, div in pairs(profile.player_div) do
      table.insert(lines, string.format("    %s: %.1fx", stat, 1.0/div))
    end
    table.insert(lines, string.format("  Experience Rate: %.1fx", profile.exp_rate))
    table.insert(lines, string.format("  Gil Rate: %.1fx", profile.gil_rate))
    if #profile.restrictions > 0 then
      table.insert(lines, "  Restrictions: " .. table.concat(profile.restrictions, ", "))
    end
    table.insert(lines, "")
  end
  
  safeCall(WriteFile, filename, table.concat(lines, "\n"))
  return "Exported profiles to:\n" .. filename
end

local function main()
  ShowDialog("WARNING", {"EXPERIMENTAL PLUGIN", "Hard mode creation requires game engine stat modification APIs.", "This plugin creates configuration profiles.", "Actual difficulty modification may require manual steps.", "", "Press OK to continue..."})
  
  while true do
    local choice = ShowDialog("Hard Mode Creator", {"1. Create Custom Profile", "2. Create from Preset", "3. View Profiles", "4. Activate Profile", "5. Export Profiles", "6. Exit"})
    if not choice or choice == 6 then break end
    
    local result = ""
    if choice == 1 then result = createCustomProfile()
    elseif choice == 2 then result = createFromPreset()
    elseif choice == 3 then result = viewProfiles()
    elseif choice == 4 then result = activateProfile()
    elseif choice == 5 then result = exportProfiles()
    end
    
    if result ~= "" then ShowDialog("Result", {result, "Press any key..."}) end
  end
end

main()
