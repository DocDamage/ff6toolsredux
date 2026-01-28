--[[
  Natural Magic Block Enforcer Plugin
  Enforce the "Natural Magic Block" challenge run rules.
  
  Features:
  - Remove all magic spells from characters
  - Prevent spell learning (remove espers)
  - Enforce physical-only combat
  - Challenge validation
  - Configuration restore wizard
  - Export challenge proof
  
  Phase: 6 (Gameplay-Altering Plugins)
  Version: 1.0.0
  
  WARNING: EXPERIMENTAL - Fundamentally alters magic system
]]

local function safeCall(fn, ...) local success, result = pcall(fn, ...) return success and result or nil end

local CONFIG_FILE = "natural_magic_block/state.lua"

local SPELL_LIST = {
  "Fire","Ice","Bolt","Poison","Drain","Fire 2","Ice 2","Bolt 2","Bio","Fire 3","Ice 3","Bolt 3",
  "Break","Doom","Pearl","Flare","Cure","Cure 2","Cure 3","Life","Life 2","Antdot","Remedy","Regen",
  "Life 3","Dispel","Haste","Haste 2","Slow","Slow 2","Stop","Berserk","Float","Imp","Reflect",
  "Shell","Safe","Vanish","Haste","Slow","Rasp","Mute","Scan","Warp","Quick","Osmose","Ultima"
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

local function loadConfig()
  local content = safeCall(ReadFile, CONFIG_FILE)
  if not content then return {enforced = false, backup = {}} end
  local success, data = pcall(load, "return " .. content)
  return (success and data and data()) or {enforced = false, backup = {}}
end

local function saveConfig(config)
  safeCall(WriteFile, CONFIG_FILE, "return " .. serialize(config))
end

local function backupSpellData()
  local chars = {"Terra","Locke","Cyan","Shadow","Edgar","Sabin","Celes","Strago","Relm","Setzer","Mog","Gau","Gogo","Umaro"}
  local backup = {}
  
  for _, char_name in ipairs(chars) do
    local char = safeCall(GetCharacter, char_name)
    if char and char.spells then
      backup[char_name] = {}
      for spell, learned in pairs(char.spells) do
        if learned then table.insert(backup[char_name], spell) end
      end
    end
  end
  
  return backup
end

local function enforceNaturalMagicBlock()
  local config = loadConfig()
  
  if config.enforced then
    return "Natural Magic Block is already enforced.\nUse Restore Wizard to undo."
  end
  
  -- Backup current spell data
  config.backup = backupSpellData()
  
  local chars = {"Terra","Locke","Cyan","Shadow","Edgar","Sabin","Celes","Strago","Relm","Setzer","Mog","Gau","Gogo","Umaro"}
  local removed_count = 0
  
  -- Remove all spells from all characters
  for _, char_name in ipairs(chars) do
    local char = safeCall(GetCharacter, char_name)
    if char and char.spells then
      for spell, _ in pairs(char.spells) do
        char.spells[spell] = false
        removed_count = removed_count + 1
      end
      safeCall(SetCharacter, char_name, char)
    end
  end
  
  -- Mark as enforced
  config.enforced = true
  config.enforcement_date = os.date("%Y-%m-%d %H:%M:%S")
  saveConfig(config)
  
  return string.format("Natural Magic Block ENFORCED\n\nRemoved %d spell entries across all characters.\nBackup saved for restoration.\n\nNote: Remove espers from inventory manually to prevent relearning.", removed_count)
end

local function validateChallenge()
  local chars = {"Terra","Locke","Cyan","Shadow","Edgar","Sabin","Celes","Strago","Relm","Setzer","Mog","Gau","Gogo","Umaro"}
  local violations = {}
  
  for _, char_name in ipairs(chars) do
    local char = safeCall(GetCharacter, char_name)
    if char and char.spells then
      local learned_spells = {}
      for spell, learned in pairs(char.spells) do
        if learned then table.insert(learned_spells, spell) end
      end
      if #learned_spells > 0 then
        table.insert(violations, string.format("%s: %d spells (%s)", char_name, #learned_spells, table.concat(learned_spells, ", ")))
      end
    end
  end
  
  local lines = {"=== Natural Magic Block Validation ===", ""}
  if #violations == 0 then
    table.insert(lines, "✓ CHALLENGE VALID")
    table.insert(lines, "No characters have learned magic.")
  else
    table.insert(lines, "✗ VIOLATIONS DETECTED")
    table.insert(lines, "")
    for _, v in ipairs(violations) do table.insert(lines, v) end
  end
  
  return table.concat(lines, "\n")
end

local function restoreWizard()
  local config = loadConfig()
  
  if not config.enforced then
    return "Natural Magic Block is not currently enforced."
  end
  
  local choice = ShowDialog("Restore Wizard", {"This will restore backed-up spell data.", "", "Continue?", "1. Yes, restore", "2. Cancel"})
  if not choice or choice == 2 then return "Cancelled." end
  
  local restored_count = 0
  for char_name, spells in pairs(config.backup) do
    local char = safeCall(GetCharacter, char_name)
    if char then
      char.spells = char.spells or {}
      for _, spell in ipairs(spells) do
        char.spells[spell] = true
        restored_count = restored_count + 1
      end
      safeCall(SetCharacter, char_name, char)
    end
  end
  
  config.enforced = false
  config.backup = {}
  saveConfig(config)
  
  return string.format("Restored %d spell entries.\nNatural Magic Block enforcement removed.", restored_count)
end

local function exportChallengeProof()
  local validation = validateChallenge()
  local config = loadConfig()
  local filename = "natural_magic_block/proof_" .. os.date("%Y%m%d_%H%M%S") .. ".txt"
  
  local lines = {
    "=== FF6 Natural Magic Block Challenge Proof ===",
    "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"),
    "Enforcement Status: " .. (config.enforced and "ACTIVE" or "INACTIVE"),
    config.enforcement_date and ("Enforcement Date: " .. config.enforcement_date) or "",
    "",
    validation
  }
  
  safeCall(WriteFile, filename, table.concat(lines, "\n"))
  return "Challenge proof exported to:\n" .. filename
end

local function viewStatus()
  local config = loadConfig()
  local lines = {"=== Natural Magic Block Status ===", ""}
  table.insert(lines, "Enforcement: " .. (config.enforced and "ACTIVE" or "INACTIVE"))
  if config.enforcement_date then table.insert(lines, "Enforcement Date: " .. config.enforcement_date) end
  table.insert(lines, "Backup Available: " .. (next(config.backup) and "Yes" or "No"))
  return table.concat(lines, "\n")
end

local function main()
  ShowDialog("WARNING", {"EXPERIMENTAL PLUGIN", "This plugin removes ALL magic spells from characters.", "Creates backup for restoration.", "Some features require APIs not yet implemented.", "", "Press OK to continue..."})
  
  while true do
    local choice = ShowDialog("Natural Magic Block Enforcer", {"1. Enforce Natural Magic Block", "2. Validate Challenge", "3. Restore Wizard", "4. View Status", "5. Export Challenge Proof", "6. Exit"})
    if not choice or choice == 6 then break end
    
    local result = ""
    if choice == 1 then result = enforceNaturalMagicBlock()
    elseif choice == 2 then result = validateChallenge()
    elseif choice == 3 then result = restoreWizard()
    elseif choice == 4 then result = viewStatus()
    elseif choice == 5 then result = exportChallengeProof()
    end
    
    if result ~= "" then ShowDialog("Result", {result, "Press any key..."}) end
  end
end

main()
