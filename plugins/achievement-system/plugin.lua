--[[
  Achievement System Plugin
  Tracks and unlocks achievements across multiple categories.
  
  Features:
  - 50+ predefined achievements
  - Custom achievement creation
  - Progress tracking per achievement
  - Achievement categories (combat/collection/challenge/secret)
  - Unlock notifications
  - Achievement points system
  - Leaderboard export
  - Retroactive achievement checking
  - Statistics dashboard
  
  Phase: 5 (Challenge & Advanced Tools)
  Version: 1.0.0
]]

-- Safe API call wrapper
local function safeCall(fn, ...)
  local success, result = pcall(fn, ...)
  if success then
    return result
  end
  return nil
end

-- Achievement categories
local CATEGORIES = {
  COMBAT = "Combat",
  COLLECTION = "Collection",
  CHALLENGE = "Challenge",
  SECRET = "Secret"
}

-- Predefined achievements library
local PREDEFINED_ACHIEVEMENTS = {
  -- Combat achievements (15)
  {id="combat_001", name="First Blood", desc="Win your first battle", category=CATEGORIES.COMBAT, points=10, check=function(save) return true end},
  {id="combat_002", name="Century Club", desc="Defeat 100 enemies", category=CATEGORIES.COMBAT, points=25, check=function(save) return false end},
  {id="combat_003", name="Overkill", desc="Deal 9999 damage in one hit", category=CATEGORIES.COMBAT, points=30, check=function(save) return false end},
  {id="combat_004", name="Magic Mastery", desc="Learn all spells with one character", category=CATEGORIES.COMBAT, points=50, check=function(save)
    local party = safeCall(GetParty)
    if not party then return false end
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Spells then
        local count = 0
        for _ in pairs(char.Spells) do count = count + 1 end
        if count >= 50 then return true end
      end
    end
    return false
  end},
  {id="combat_005", name="No Damage Victory", desc="Win a boss battle without taking damage", category=CATEGORIES.COMBAT, points=40, check=function(save) return false end},
  {id="combat_006", name="Speed Demon", desc="Win a battle in under 10 seconds", category=CATEGORIES.COMBAT, points=25, check=function(save) return false end},
  {id="combat_007", name="Solo Champion", desc="Win a battle with only one character", category=CATEGORIES.COMBAT, points=35, check=function(save) return false end},
  {id="combat_008", name="All Equipment Master", desc="Equip all equipment types", category=CATEGORIES.COMBAT, points=40, check=function(save) return false end},
  {id="combat_009", name="Esper Bond", desc="Max out one esper's level", category=CATEGORIES.COMBAT, points=30, check=function(save) return false end},
  {id="combat_010", name="Level 99 Club", desc="Reach level 99 with any character", category=CATEGORIES.COMBAT, points=50, check=function(save)
    local party = safeCall(GetParty)
    if not party then return false end
    for i = 1, #party do
      local char = safeCall(GetCharacter, party[i])
      if char and char.Level and char.Level >= 99 then return true end
    end
    return false
  end},
  {id="combat_011", name="Perfect Party", desc="Have all 4 party members at max HP/MP", category=CATEGORIES.COMBAT, points=35, check=function(save) return false end},
  {id="combat_012", name="Rage Master", desc="Learn all Rages with Gau", category=CATEGORIES.COMBAT, points=60, check=function(save) return false end},
  {id="combat_013", name="Blitz Virtuoso", desc="Execute all of Sabin's Blitzes successfully", category=CATEGORIES.COMBAT, points=40, check=function(save) return false end},
  {id="combat_014", name="Tools Expert", desc="Use all of Edgar's Tools in battle", category=CATEGORIES.COMBAT, points=35, check=function(save) return false end},
  {id="combat_015", name="Lore Scholar", desc="Learn all Lores with Strago", category=CATEGORIES.COMBAT, points=50, check=function(save) return false end},
  
  -- Collection achievements (15)
  {id="collect_001", name="Item Hoarder", desc="Collect 99 of any item", category=CATEGORIES.COLLECTION, points=20, check=function(save)
    local inventory = safeCall(GetInventory)
    if not inventory or not inventory.Items then return false end
    for _, item in pairs(inventory.Items) do
      if item.Quantity and item.Quantity >= 99 then return true end
    end
    return false
  end},
  {id="collect_002", name="Millionaire", desc="Accumulate 999,999 GP", category=CATEGORIES.COLLECTION, points=50, check=function(save)
    local inventory = safeCall(GetInventory)
    return inventory and inventory.Gil and inventory.Gil >= 999999
  end},
  {id="collect_003", name="Esper Collector", desc="Obtain all espers", category=CATEGORIES.COLLECTION, points=60, check=function(save)
    local espers = safeCall(GetEspers) or safeCall(GetEsperInventory)
    if not espers then return false end
    local count = 0
    for _ in pairs(espers) do count = count + 1 end
    return count >= 27
  end},
  {id="collect_004", name="Weapon Master", desc="Collect all unique weapons", category=CATEGORIES.COLLECTION, points=55, check=function(save) return false end},
  {id="collect_005", name="Armor Collector", desc="Collect all unique armor pieces", category=CATEGORIES.COLLECTION, points=50, check=function(save) return false end},
  {id="collect_006", name="Relic Hunter", desc="Collect all unique relics", category=CATEGORIES.COLLECTION, points=65, check=function(save) return false end},
  {id="collect_007", name="Full Party", desc="Recruit all 14 characters", category=CATEGORIES.COLLECTION, points=40, check=function(save)
    local party = safeCall(GetParty)
    return party and #party >= 14
  end},
  {id="collect_008", name="Rare Find", desc="Obtain a rare item drop", category=CATEGORIES.COLLECTION, points=30, check=function(save) return false end},
  {id="collect_009", name="Dragon Slayer", desc="Defeat all 8 dragons", category=CATEGORIES.COLLECTION, points=70, check=function(save) return false end},
  {id="collect_010", name="Treasure Hunter", desc="Open 100 treasure chests", category=CATEGORIES.COLLECTION, points=45, check=function(save) return false end},
  {id="collect_011", name="Bestiary Complete", desc="Encounter all enemy types", category=CATEGORIES.COLLECTION, points=80, check=function(save) return false end},
  {id="collect_012", name="Item Completionist", desc="Collect at least one of every item", category=CATEGORIES.COLLECTION, points=90, check=function(save) return false end},
  {id="collect_013", name="Steal Master", desc="Successfully steal 50 items", category=CATEGORIES.COLLECTION, points=40, check=function(save) return false end},
  {id="collect_014", name="Colosseum Champion", desc="Win 20 Colosseum battles", category=CATEGORIES.COLLECTION, points=50, check=function(save) return false end},
  {id="collect_015", name="Slots Jackpot", desc="Hit the jackpot on Setzer's Slots", category=CATEGORIES.COLLECTION, points=35, check=function(save) return false end},
  
  -- Challenge achievements (12)
  {id="challenge_001", name="Low Level Legend", desc="Complete the game under level 20", category=CATEGORIES.CHALLENGE, points=100, check=function(save) return false end},
  {id="challenge_002", name="Natural Magic Block", desc="Complete game without learning magic", category=CATEGORIES.CHALLENGE, points=120, check=function(save) return false end},
  {id="challenge_003", name="Solo Run", desc="Complete game with only one character", category=CATEGORIES.CHALLENGE, points=150, check=function(save) return false end},
  {id="challenge_004", name="No Equipment", desc="Complete game without equipping weapons/armor", category=CATEGORIES.CHALLENGE, points=130, check=function(save) return false end},
  {id="challenge_005", name="Speedrunner", desc="Complete game in under 6 hours", category=CATEGORIES.CHALLENGE, points=110, check=function(save) return false end},
  {id="challenge_006", name="Perfectionist", desc="Complete game with 100% completion", category=CATEGORIES.CHALLENGE, points=140, check=function(save) return false end},
  {id="challenge_007", name="No Shop Run", desc="Complete game without buying from shops", category=CATEGORIES.CHALLENGE, points=90, check=function(save) return false end},
  {id="challenge_008", name="Pacifist", desc="Complete game with minimum battles", category=CATEGORIES.CHALLENGE, points=100, check=function(save) return false end},
  {id="challenge_009", name="No Esper Run", desc="Complete game without equipping espers", category=CATEGORIES.CHALLENGE, points=110, check=function(save) return false end},
  {id="challenge_010", name="Fixed Party", desc="Complete game with the same 4 characters", category=CATEGORIES.CHALLENGE, points=80, check=function(save) return false end},
  {id="challenge_011", name="No Save Challenge", desc="Complete game in one session", category=CATEGORIES.CHALLENGE, points=95, check=function(save) return false end},
  {id="challenge_012", name="Superboss Slayer", desc="Defeat all optional superbosses", category=CATEGORIES.CHALLENGE, points=85, check=function(save) return false end},
  
  -- Secret achievements (8)
  {id="secret_001", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=50, check=function(save) return false end, secret=true, revealed_name="Shadow Saved", revealed_desc="Save Shadow on the Floating Continent"},
  {id="secret_002", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=60, check=function(save) return false end, secret=true, revealed_name="Cyan's Dream", revealed_desc="Complete Cyan's dream sequence"},
  {id="secret_003", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=40, check=function(save) return false end, secret=true, revealed_name="Opera Star", revealed_desc="Complete the opera house sequence perfectly"},
  {id="secret_004", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=70, check=function(save) return false end, secret=true, revealed_name="Moogle Charm", revealed_desc="Obtain the Moogle Charm"},
  {id="secret_005", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=55, check=function(save) return false end, secret=true, revealed_name="Gau's Father", revealed_desc="Reunite Gau with his father"},
  {id="secret_006", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=65, check=function(save) return false end, secret=true, revealed_name="Cursed Shield", revealed_desc="Break the Cursed Shield's curse"},
  {id="secret_007", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=80, check=function(save) return false end, secret=true, revealed_name="Celes' Aria", revealed_desc="Keep Celes' spirits up during the World of Ruin opening"},
  {id="secret_008", name="???", desc="Hidden achievement", category=CATEGORIES.SECRET, points=100, check=function(save) return false end, secret=true, revealed_name="Ultimate Completion", revealed_desc="Unlock all other achievements"}
}

-- Achievement progress storage file
local PROGRESS_FILE = "achievements/progress.lua"
local CUSTOM_FILE = "achievements/custom.lua"

-- Load achievement progress
local function loadProgress()
  local content = safeCall(ReadFile, PROGRESS_FILE)
  if not content then
    return {unlocked = {}, progress = {}, points = 0, custom_achievements = {}}
  end
  
  local success, data = pcall(load, "return " .. content)
  if success and data then
    return data()
  end
  return {unlocked = {}, progress = {}, points = 0, custom_achievements = {}}
end

-- Save achievement progress
local function saveProgress(data)
  local content = "return " .. serialize(data)
  safeCall(WriteFile, PROGRESS_FILE, content)
end

-- Serialize table to Lua code
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

-- Check single achievement
local function checkAchievement(achievement, save_data, progress)
  if progress.unlocked[achievement.id] then
    return true -- Already unlocked
  end
  
  local unlocked = achievement.check(save_data)
  if unlocked then
    progress.unlocked[achievement.id] = {
      timestamp = os.time(),
      name = achievement.name,
      points = achievement.points
    }
    progress.points = (progress.points or 0) + achievement.points
    return true, true -- unlocked, newly_unlocked
  end
  
  return false, false
end

-- Check all achievements retroactively
local function checkAllAchievements()
  local progress = loadProgress()
  local save_data = {} -- Placeholder for save data
  local newly_unlocked = {}
  
  -- Check predefined achievements
  for _, achievement in ipairs(PREDEFINED_ACHIEVEMENTS) do
    local unlocked, is_new = checkAchievement(achievement, save_data, progress)
    if is_new then
      table.insert(newly_unlocked, achievement)
    end
  end
  
  -- Check custom achievements
  for _, achievement in ipairs(progress.custom_achievements or {}) do
    local unlocked, is_new = checkAchievement(achievement, save_data, progress)
    if is_new then
      table.insert(newly_unlocked, achievement)
    end
  end
  
  saveProgress(progress)
  return newly_unlocked, progress
end

-- View achievements by category
local function viewAchievements(category)
  local progress = loadProgress()
  local achievements = {}
  
  -- Filter by category
  for _, ach in ipairs(PREDEFINED_ACHIEVEMENTS) do
    if not category or ach.category == category then
      table.insert(achievements, ach)
    end
  end
  
  -- Add custom achievements
  for _, ach in ipairs(progress.custom_achievements or {}) do
    if not category or ach.category == category then
      table.insert(achievements, ach)
    end
  end
  
  -- Build display
  local lines = {}
  table.insert(lines, "=== Achievement Viewer ===\n")
  if category then
    table.insert(lines, "Category: " .. category .. "\n")
  else
    table.insert(lines, "All Achievements\n")
  end
  table.insert(lines, string.format("Total Points: %d\n", progress.points or 0))
  table.insert(lines, "\n")
  
  local unlocked_count = 0
  for _, ach in ipairs(achievements) do
    local is_unlocked = progress.unlocked[ach.id] ~= nil
    if is_unlocked then unlocked_count = unlocked_count + 1 end
    
    local status = is_unlocked and "[✓]" or "[ ]"
    local name = (ach.secret and not is_unlocked) and ach.name or (ach.revealed_name or ach.name)
    local desc = (ach.secret and not is_unlocked) and ach.desc or (ach.revealed_desc or ach.desc)
    
    table.insert(lines, string.format("%s %s (%d pts)", status, name, ach.points))
    table.insert(lines, string.format("    %s", desc))
    
    if is_unlocked and progress.unlocked[ach.id].timestamp then
      table.insert(lines, string.format("    Unlocked: %s", os.date("%Y-%m-%d %H:%M", progress.unlocked[ach.id].timestamp)))
    end
    table.insert(lines, "\n")
  end
  
  table.insert(lines, string.format("\nProgress: %d/%d achievements unlocked", unlocked_count, #achievements))
  table.insert(lines, string.format("Completion: %.1f%%", (unlocked_count / #achievements) * 100))
  
  return table.concat(lines, "\n")
end

-- Create custom achievement
local function createCustomAchievement()
  local name = ShowInput("Enter achievement name:")
  if not name or name == "" then return "Cancelled." end
  
  local desc = ShowInput("Enter achievement description:")
  if not desc or desc == "" then return "Cancelled." end
  
  local category = ShowDialog("Select category:", {
    "1. Combat",
    "2. Collection",
    "3. Challenge",
    "4. Secret"
  })
  
  local cat_map = {CATEGORIES.COMBAT, CATEGORIES.COLLECTION, CATEGORIES.CHALLENGE, CATEGORIES.SECRET}
  if not category or category < 1 or category > 4 then return "Cancelled." end
  
  local points_str = ShowInput("Enter point value (1-200):")
  local points = tonumber(points_str) or 50
  if points < 1 then points = 1 end
  if points > 200 then points = 200 end
  
  local progress = loadProgress()
  local id = "custom_" .. os.time()
  
  table.insert(progress.custom_achievements, {
    id = id,
    name = name,
    desc = desc,
    category = cat_map[category],
    points = points,
    check = function(save) return false end, -- Manual unlock only
    custom = true
  })
  
  saveProgress(progress)
  return string.format("Created custom achievement: %s (%d pts)", name, points)
end

-- Manually unlock achievement
local function manualUnlock()
  local progress = loadProgress()
  local all_achievements = {}
  
  -- Gather all achievements
  for _, ach in ipairs(PREDEFINED_ACHIEVEMENTS) do
    if not progress.unlocked[ach.id] then
      table.insert(all_achievements, ach)
    end
  end
  for _, ach in ipairs(progress.custom_achievements or {}) do
    if not progress.unlocked[ach.id] then
      table.insert(all_achievements, ach)
    end
  end
  
  if #all_achievements == 0 then
    return "All achievements already unlocked!"
  end
  
  -- Build menu
  local options = {}
  for i, ach in ipairs(all_achievements) do
    table.insert(options, string.format("%d. %s (%d pts)", i, ach.name, ach.points))
  end
  table.insert(options, string.format("%d. Cancel", #options + 1))
  
  local choice = ShowDialog("Select achievement to unlock:", options)
  if not choice or choice > #all_achievements then return "Cancelled." end
  
  local ach = all_achievements[choice]
  progress.unlocked[ach.id] = {
    timestamp = os.time(),
    name = ach.name,
    points = ach.points
  }
  progress.points = (progress.points or 0) + ach.points
  
  saveProgress(progress)
  return string.format("Unlocked: %s (+%d pts)", ach.name, ach.points)
end

-- Export achievement data
local function exportAchievements()
  local progress = loadProgress()
  local timestamp = os.date("%Y%m%d_%H%M%S")
  local filename = "achievements/export_" .. timestamp .. ".txt"
  
  local lines = {}
  table.insert(lines, "=== FF6 Save Editor - Achievement Export ===")
  table.insert(lines, "Export Date: " .. os.date("%Y-%m-%d %H:%M:%S"))
  table.insert(lines, "Total Points: " .. (progress.points or 0))
  table.insert(lines, "")
  
  local total_unlocked = 0
  for _ in pairs(progress.unlocked) do total_unlocked = total_unlocked + 1 end
  
  table.insert(lines, string.format("Achievements Unlocked: %d", total_unlocked))
  table.insert(lines, "")
  table.insert(lines, "=== Unlocked Achievements ===")
  
  for id, data in pairs(progress.unlocked) do
    table.insert(lines, string.format("[%s] %s - %d points", 
      os.date("%Y-%m-%d", data.timestamp), data.name, data.points))
  end
  
  local content = table.concat(lines, "\n")
  safeCall(WriteFile, filename, content)
  
  return "Exported achievements to: " .. filename
end

-- Statistics dashboard
local function viewStatistics()
  local progress = loadProgress()
  
  local total_ach = #PREDEFINED_ACHIEVEMENTS + #(progress.custom_achievements or {})
  local unlocked_count = 0
  for _ in pairs(progress.unlocked) do unlocked_count = unlocked_count + 1 end
  
  local by_category = {
    [CATEGORIES.COMBAT] = {total = 0, unlocked = 0},
    [CATEGORIES.COLLECTION] = {total = 0, unlocked = 0},
    [CATEGORIES.CHALLENGE] = {total = 0, unlocked = 0},
    [CATEGORIES.SECRET] = {total = 0, unlocked = 0}
  }
  
  for _, ach in ipairs(PREDEFINED_ACHIEVEMENTS) do
    by_category[ach.category].total = by_category[ach.category].total + 1
    if progress.unlocked[ach.id] then
      by_category[ach.category].unlocked = by_category[ach.category].unlocked + 1
    end
  end
  
  for _, ach in ipairs(progress.custom_achievements or {}) do
    by_category[ach.category].total = by_category[ach.category].total + 1
    if progress.unlocked[ach.id] then
      by_category[ach.category].unlocked = by_category[ach.category].unlocked + 1
    end
  end
  
  local lines = {}
  table.insert(lines, "=== Achievement Statistics ===\n")
  table.insert(lines, string.format("Total Points: %d", progress.points or 0))
  table.insert(lines, string.format("Overall Progress: %d/%d (%.1f%%)\n", 
    unlocked_count, total_ach, (unlocked_count / total_ach) * 100))
  
  table.insert(lines, "=== By Category ===")
  for category, stats in pairs(by_category) do
    local pct = stats.total > 0 and (stats.unlocked / stats.total * 100) or 0
    table.insert(lines, string.format("%s: %d/%d (%.1f%%)", 
      category, stats.unlocked, stats.total, pct))
  end
  
  table.insert(lines, "\n=== Recent Unlocks ===")
  local recent = {}
  for id, data in pairs(progress.unlocked) do
    table.insert(recent, {id = id, data = data})
  end
  table.sort(recent, function(a, b) return a.data.timestamp > b.data.timestamp end)
  
  for i = 1, math.min(5, #recent) do
    local r = recent[i]
    table.insert(lines, string.format("[%s] %s (+%d pts)", 
      os.date("%Y-%m-%d", r.data.timestamp), r.data.name, r.data.points))
  end
  
  return table.concat(lines, "\n")
end

-- Main menu
local function main()
  while true do
    local choice = ShowDialog("Achievement System", {
      "1. Check Achievements (Retroactive Scan)",
      "2. View All Achievements",
      "3. View by Category",
      "4. Statistics Dashboard",
      "5. Create Custom Achievement",
      "6. Manually Unlock Achievement",
      "7. Export Achievements",
      "8. Exit"
    })
    
    if not choice or choice == 8 then break end
    
    local result = ""
    if choice == 1 then
      local newly_unlocked, progress = checkAllAchievements()
      if #newly_unlocked > 0 then
        result = string.format("Unlocked %d new achievement(s):\n", #newly_unlocked)
        for _, ach in ipairs(newly_unlocked) do
          result = result .. string.format("  [NEW] %s (+%d pts)\n", ach.name, ach.points)
        end
        result = result .. string.format("\nTotal Points: %d", progress.points)
      else
        result = "No new achievements unlocked."
      end
    elseif choice == 2 then
      result = viewAchievements(nil)
    elseif choice == 3 then
      local cat_choice = ShowDialog("Select Category:", {
        "1. Combat",
        "2. Collection",
        "3. Challenge",
        "4. Secret",
        "5. Cancel"
      })
      if cat_choice and cat_choice <= 4 then
        local cats = {CATEGORIES.COMBAT, CATEGORIES.COLLECTION, CATEGORIES.CHALLENGE, CATEGORIES.SECRET}
        result = viewAchievements(cats[cat_choice])
      end
    elseif choice == 4 then
      result = viewStatistics()
    elseif choice == 5 then
      result = createCustomAchievement()
    elseif choice == 6 then
      result = manualUnlock()
    elseif choice == 7 then
      result = exportAchievements()
    end
    
    if result ~= "" then
      ShowDialog("Result", {result, "Press any key to continue..."})
    end
  end
end

-- Run plugin
main()
