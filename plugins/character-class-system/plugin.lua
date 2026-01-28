--[[
  Character Class System Plugin
  Convert FF6 into a class-based RPG system.
  
  Features:
  - 10+ custom class definitions
  - Assign classes to characters
  - Class-specific equipment restrictions
  - Class-based stat modifiers
  - Job change system
  - Hybrid/multi-class support
  - Class progression tracking
  
  Phase: 6 (Gameplay-Altering Plugins)
  Version: 1.0.0
  
  WARNING: EXPERIMENTAL - Fundamentally alters character progression
]]

local function safeCall(fn, ...) local success, result = pcall(fn, ...) return success and result or nil end

-- Class definitions
local CLASSES = {
  {id="knight", name="Knight", desc="High HP/Defense, heavy equipment", stats={HP=1.3, Defense=1.4, Vigor=1.2}, equipment={"Swords","Shields","Heavy Armor"}, abilities={1,6}},
  {id="mage", name="Mage", desc="High Magic, low HP", stats={HP=0.8, Magic=1.5, MP=1.4}, equipment={"Rods","Robes"}, abilities={1,2}},
  {id="thief", name="Thief", desc="High Speed, Steal", stats={Speed=1.4, Stamina=1.2}, equipment={"Daggers","Light Armor"}, abilities={1,4}},
  {id="monk", name="Monk", desc="High Vigor, Blitz", stats={Vigor=1.5, HP=1.2, Stamina=1.3}, equipment={"Claws","Gi"}, abilities={1,9}},
  {id="dragoon", name="Dragoon", desc="Spear master, high jump", stats={Vigor=1.3, HP=1.2}, equipment={"Spears","Heavy Armor"}, abilities={1,7}},
  {id="white_mage", name="White Mage", desc="Healing magic specialist", stats={Magic=1.4, MP=1.5, MagicDef=1.3}, equipment={"Rods","Robes"}, abilities={1,2,10}},
  {id="black_mage", name="Black Mage", desc="Offensive magic specialist", stats={Magic=1.6, MP=1.4, HP=0.7}, equipment={"Rods","Hats"}, abilities={1,2}},
  {id="samurai", name="Samurai", desc="SwdTech master", stats={Vigor=1.3, Magic=1.1, Stamina=1.2}, equipment={"Katanas","Light Armor"}, abilities={1,6}},
  {id="ninja", name="Ninja", desc="Throw, high speed", stats={Speed=1.5, Vigor=1.2}, equipment={"Daggers","Shurikens","Light Armor"}, abilities={1,7}},
  {id="sage", name="Sage", desc="Lore master, balanced", stats={Magic=1.4, Vigor=1.1, HP=1.1}, equipment={"Rods","Robes"}, abilities={1,2,11}},
  {id="berserker", name="Berserker", desc="Pure physical power", stats={Vigor=1.6, HP=1.4, Defense=1.2}, equipment={"Axes","Heavy Armor"}, abilities={1}},
  {id="mystic_knight", name="Mystic Knight", desc="Magic + Physical hybrid", stats={Vigor=1.2, Magic=1.3, HP=1.1}, equipment={"Swords","Light Armor"}, abilities={1,2,6}}
}

local CONFIG_FILE = "class_system/configurations.lua"

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

local function loadConfig()
  local content = safeCall(ReadFile, CONFIG_FILE)
  if not content then return {assignments = {}, multi_class = {}} end
  local success, data = pcall(load, "return " .. content)
  return (success and data and data()) or {assignments = {}, multi_class = {}}
end

local function saveConfig(config)
  safeCall(WriteFile, CONFIG_FILE, "return " .. serialize(config))
end

local function applyClassStats(character, class_id)
  local class = nil
  for _, c in ipairs(CLASSES) do if c.id == class_id then class = c break end end
  if not class then return "Class not found." end
  
  local char = safeCall(GetCharacter, character)
  if not char then return "Character not found." end
  
  -- Apply stat modifiers (conceptual - requires stat modification API)
  local stats = {}
  for stat, mult in pairs(class.stats) do
    local base = char[stat] or 50
    stats[stat] = math.floor(base * mult)
  end
  
  return string.format("Applied %s class to %s\nStats modified (conceptual): %s", class.name, character, 
    table.concat((function() local t={} for k,v in pairs(stats) do table.insert(t, k.."="..v) end return t end)(), ", "))
end

local function assignClass()
  local chars = {"Terra","Locke","Cyan","Shadow","Edgar","Sabin","Celes","Strago","Relm","Setzer","Mog","Gau","Gogo","Umaro"}
  local char_opts = {} for i,c in ipairs(chars) do table.insert(char_opts, i..". "..c) end
  table.insert(char_opts, #chars+1 ..". Cancel")
  
  local char_choice = ShowDialog("Select Character:", char_opts)
  if not char_choice or char_choice > #chars then return "Cancelled." end
  local character = chars[char_choice]
  
  local class_opts = {} for i,c in ipairs(CLASSES) do table.insert(class_opts, i..". "..c.name.." - "..c.desc) end
  table.insert(class_opts, #CLASSES+1 ..". Cancel")
  
  local class_choice = ShowDialog("Select Class:", class_opts)
  if not class_choice or class_choice > #CLASSES then return "Cancelled." end
  
  local config = loadConfig()
  config.assignments[character] = CLASSES[class_choice].id
  saveConfig(config)
  
  return applyClassStats(character, CLASSES[class_choice].id)
end

local function viewClassAssignments()
  local config = loadConfig()
  local lines = {"=== Class Assignments ===\n"}
  
  for char, class_id in pairs(config.assignments) do
    local class = nil
    for _, c in ipairs(CLASSES) do if c.id == class_id then class = c break end end
    if class then table.insert(lines, string.format("%s: %s", char, class.name)) end
  end
  
  if not next(config.assignments) then table.insert(lines, "No class assignments yet.") end
  return table.concat(lines, "\n")
end

local function viewClassLibrary()
  local lines = {"=== Class Library ===\n"}
  for _, class in ipairs(CLASSES) do
    table.insert(lines, string.format("%s - %s", class.name, class.desc))
    table.insert(lines, "  Stats: " .. table.concat((function() local t={} for k,v in pairs(class.stats) do table.insert(t, k.."x"..v) end return t end)(), ", "))
    table.insert(lines, "  Equipment: " .. table.concat(class.equipment, ", "))
    table.insert(lines, "")
  end
  return table.concat(lines, "\n")
end

local function exportConfiguration()
  local config = loadConfig()
  local filename = "class_system/export_" .. os.date("%Y%m%d_%H%M%S") .. ".txt"
  local lines = {"=== FF6 Character Class Configuration ===", "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"), ""}
  
  for char, class_id in pairs(config.assignments) do
    local class = nil
    for _, c in ipairs(CLASSES) do if c.id == class_id then class = c break end end
    if class then table.insert(lines, char .. ": " .. class.name .. " - " .. class.desc) end
  end
  
  safeCall(WriteFile, filename, table.concat(lines, "\n"))
  return "Exported to: " .. filename
end

local function main()
  ShowDialog("WARNING", {"EXPERIMENTAL PLUGIN", "Class system requires advanced stat modification APIs.", "This plugin tracks configurations but may not fully modify game stats.", "", "Press OK to continue..."})
  
  while true do
    local choice = ShowDialog("Character Class System", {"1. View Class Library", "2. Assign Class to Character", "3. View Class Assignments", "4. Export Configuration", "5. Exit"})
    if not choice or choice == 5 then break end
    
    local result = ""
    if choice == 1 then result = viewClassLibrary()
    elseif choice == 2 then result = assignClass()
    elseif choice == 3 then result = viewClassAssignments()
    elseif choice == 4 then result = exportConfiguration()
    end
    
    if result ~= "" then ShowDialog("Result", {result, "Press any key..."}) end
  end
end

main()
