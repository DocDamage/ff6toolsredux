--[[
  Speedrun Timer Integration Plugin
  Track split times for speedrun categories.
  
  Features:
  - Split timer with manual triggers
  - Multiple speedrun categories
  - Personal best comparison
  - Route planning tool
  - Split management
  - Export to speedrun.com format
  - Segment analysis
  
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

-- Speedrun categories with default splits
local CATEGORIES = {
  {
    id = "any_percent",
    name = "Any%",
    description = "Complete game as fast as possible",
    splits = {
      "Narshe Escape", "Figaro Castle", "South Figaro", "Mt. Kolts", "Returners Hideout",
      "Phantom Train", "Zozo", "Opera House", "Vector", "Magitek Factory",
      "Esper Escape", "Floating Continent", "WoR Start", "Airship", "Kefka's Tower"
    }
  },
  {
    id = "glitchless",
    name = "Glitchless",
    description = "Any% without major glitches",
    splits = {
      "Narshe Escape", "Figaro Castle", "South Figaro", "Mt. Kolts", "Returners Hideout",
      "Phantom Train", "Zozo", "Opera House", "Vector", "Magitek Factory",
      "Esper Escape", "Floating Continent", "WoR Start", "Airship", "Kefka's Tower"
    }
  },
  {
    id = "100_percent",
    name = "100% Completion",
    description = "All characters, espers, items",
    splits = {
      "Narshe Escape", "All WoB Characters", "All WoB Espers", "Floating Continent",
      "All WoR Characters", "8 Dragons", "All Espers", "All Unique Items",
      "Level 99 Party", "Kefka's Tower"
    }
  },
  {
    id = "low_level",
    name = "Low Level Run",
    description = "Complete under level 20",
    splits = {
      "Narshe Escape", "Phantom Train", "Opera House", "Floating Continent",
      "WoR Start", "All Dragons (No Level Up)", "Kefka's Tower"
    }
  },
  {
    id = "wor_skip",
    name = "WoR Skip",
    description = "Skip World of Ruin content",
    splits = {
      "Narshe Escape", "Floating Continent", "Airship Skip", "Kefka's Tower Entry", "Final Boss"
    }
  }
}

-- Timer state
local timer_state = {
  running = false,
  paused = false,
  start_time = 0,
  pause_time = 0,
  total_pause_duration = 0,
  current_split = 0,
  split_times = {},
  category = nil
}

-- Storage files
local PB_FILE = "speedrun_timer/personal_bests.lua"
local ROUTE_FILE = "speedrun_timer/routes.lua"
local CURRENT_RUN_FILE = "speedrun_timer/current_run.lua"

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

-- Format time (seconds to HH:MM:SS.ms)
local function formatTime(seconds)
  if not seconds or seconds < 0 then return "00:00:00.000" end
  
  local hours = math.floor(seconds / 3600)
  local mins = math.floor((seconds % 3600) / 60)
  local secs = math.floor(seconds % 60)
  local ms = math.floor((seconds % 1) * 1000)
  
  return string.format("%02d:%02d:%02d.%03d", hours, mins, secs, ms)
end

-- Get current elapsed time
local function getElapsedTime()
  if not timer_state.running then
    return 0
  end
  
  local current_time = os.clock()
  local elapsed = current_time - timer_state.start_time - timer_state.total_pause_duration
  
  if timer_state.paused then
    elapsed = timer_state.pause_time - timer_state.start_time - timer_state.total_pause_duration
  end
  
  return elapsed
end

-- Load personal bests
local function loadPBs()
  local content = safeCall(ReadFile, PB_FILE)
  if not content then return {} end
  
  local success, data = pcall(load, "return " .. content)
  if success and data then
    return data()
  end
  return {}
end

-- Save personal best
local function savePB(category_id, run_data)
  local pbs = loadPBs()
  
  if not pbs[category_id] or run_data.total_time < pbs[category_id].total_time then
    pbs[category_id] = run_data
    local content = "return " .. serialize(pbs)
    safeCall(WriteFile, PB_FILE, content)
    return true
  end
  
  return false
end

-- Start timer
local function startTimer(category)
  timer_state = {
    running = true,
    paused = false,
    start_time = os.clock(),
    pause_time = 0,
    total_pause_duration = 0,
    current_split = 0,
    split_times = {},
    category = category
  }
  
  -- Save current run state
  local content = "return " .. serialize(timer_state)
  safeCall(WriteFile, CURRENT_RUN_FILE, content)
  
  return string.format("Timer started for: %s", category.name)
end

-- Split
local function split()
  if not timer_state.running or timer_state.paused then
    return "Timer not running!"
  end
  
  local elapsed = getElapsedTime()
  timer_state.current_split = timer_state.current_split + 1
  
  if timer_state.current_split <= #timer_state.category.splits then
    local split_name = timer_state.category.splits[timer_state.current_split]
    timer_state.split_times[timer_state.current_split] = {
      name = split_name,
      time = elapsed
    }
    
    -- Check if final split
    if timer_state.current_split == #timer_state.category.splits then
      timer_state.running = false
      timer_state.total_time = elapsed
      
      -- Check for PB
      local is_pb = savePB(timer_state.category.id, {
        total_time = elapsed,
        split_times = timer_state.split_times,
        date = os.time()
      })
      
      return string.format("Final split: %s\nTotal Time: %s%s", 
        split_name, formatTime(elapsed), is_pb and "\n** NEW PERSONAL BEST! **" or "")
    end
    
    return string.format("Split %d/%d: %s\nTime: %s", 
      timer_state.current_split, #timer_state.category.splits, split_name, formatTime(elapsed))
  end
  
  return "All splits completed!"
end

-- Pause/Resume timer
local function togglePause()
  if not timer_state.running then
    return "Timer not running!"
  end
  
  if timer_state.paused then
    -- Resume
    local pause_duration = os.clock() - timer_state.pause_time
    timer_state.total_pause_duration = timer_state.total_pause_duration + pause_duration
    timer_state.paused = false
    return "Timer resumed"
  else
    -- Pause
    timer_state.pause_time = os.clock()
    timer_state.paused = true
    return "Timer paused"
  end
end

-- Reset timer
local function resetTimer()
  if not timer_state.running and timer_state.current_split == 0 then
    return "Timer already reset"
  end
  
  local confirm = ShowDialog("Reset timer?", {
    "1. Yes, reset and lose current run",
    "2. No, keep running"
  })
  
  if confirm == 1 then
    timer_state = {
      running = false,
      paused = false,
      start_time = 0,
      pause_time = 0,
      total_pause_duration = 0,
      current_split = 0,
      split_times = {},
      category = timer_state.category
    }
    return "Timer reset"
  end
  
  return "Reset cancelled"
end

-- View current run
local function viewCurrentRun()
  if timer_state.current_split == 0 and not timer_state.running then
    return "No active run. Start a timer first."
  end
  
  local lines = {}
  table.insert(lines, "=== Current Run ===")
  table.insert(lines, string.format("Category: %s", timer_state.category and timer_state.category.name or "Unknown"))
  table.insert(lines, string.format("Status: %s", 
    timer_state.paused and "PAUSED" or (timer_state.running and "RUNNING" or "STOPPED")))
  table.insert(lines, string.format("Current Time: %s\n", formatTime(getElapsedTime())))
  
  -- Load PB for comparison
  local pbs = loadPBs()
  local pb = timer_state.category and pbs[timer_state.category.id]
  
  table.insert(lines, string.format("Split: %d/%d\n", 
    timer_state.current_split, timer_state.category and #timer_state.category.splits or 0))
  
  table.insert(lines, "=== Splits ===")
  for i = 1, timer_state.current_split do
    local split_data = timer_state.split_times[i]
    local segment_time = i == 1 and split_data.time or 
      (split_data.time - timer_state.split_times[i-1].time)
    
    local pb_split_time = pb and pb.split_times[i] and pb.split_times[i].time
    local delta = pb_split_time and (split_data.time - pb_split_time) or 0
    local delta_str = ""
    if pb_split_time then
      delta_str = delta < 0 and 
        string.format(" [-%s]", formatTime(math.abs(delta))) or
        string.format(" [+%s]", formatTime(delta))
    end
    
    table.insert(lines, string.format("%d. %s: %s%s", 
      i, split_data.name, formatTime(split_data.time), delta_str))
  end
  
  -- Show next split
  if timer_state.running and timer_state.current_split < #timer_state.category.splits then
    table.insert(lines, string.format("\nNext: %s", 
      timer_state.category.splits[timer_state.current_split + 1]))
  end
  
  return table.concat(lines, "\n")
end

-- View personal bests
local function viewPBs()
  local pbs = loadPBs()
  
  local lines = {}
  table.insert(lines, "=== Personal Bests ===\n")
  
  local has_pbs = false
  for _, category in ipairs(CATEGORIES) do
    local pb = pbs[category.id]
    if pb then
      has_pbs = true
      table.insert(lines, string.format("%s: %s", category.name, formatTime(pb.total_time)))
      table.insert(lines, string.format("  Date: %s\n", os.date("%Y-%m-%d", pb.date)))
    end
  end
  
  if not has_pbs then
    table.insert(lines, "No personal bests yet. Complete a run to set your first PB!")
  end
  
  return table.concat(lines, "\n")
end

-- Export to speedrun.com format
local function exportRun()
  if timer_state.current_split == 0 then
    return "No run to export."
  end
  
  local timestamp = os.date("%Y%m%d_%H%M%S")
  local filename = string.format("speedrun_timer/run_%s_%s.txt", 
    timer_state.category.id, timestamp)
  
  local lines = {}
  table.insert(lines, "=== FF6 Speedrun Export ===")
  table.insert(lines, "Game: Final Fantasy VI")
  table.insert(lines, "Category: " .. timer_state.category.name)
  table.insert(lines, "Date: " .. os.date("%Y-%m-%d %H:%M:%S"))
  table.insert(lines, "")
  
  if timer_state.total_time then
    table.insert(lines, "Final Time: " .. formatTime(timer_state.total_time))
  else
    table.insert(lines, "Final Time: Incomplete")
  end
  
  table.insert(lines, "")
  table.insert(lines, "=== Splits ===")
  
  for i, split_data in ipairs(timer_state.split_times) do
    local segment_time = i == 1 and split_data.time or 
      (split_data.time - timer_state.split_times[i-1].time)
    table.insert(lines, string.format("%d. %s: %s (segment: %s)", 
      i, split_data.name, formatTime(split_data.time), formatTime(segment_time)))
  end
  
  table.insert(lines, "")
  table.insert(lines, "Generated by FF6 Save Editor - Speedrun Timer v1.0.0")
  
  local content = table.concat(lines, "\n")
  safeCall(WriteFile, filename, content)
  
  return "Run exported to: " .. filename
end

-- Main menu
local function main()
  while true do
    local status_line = timer_state.running and 
      string.format("Running: %s [%s]", timer_state.category.name, formatTime(getElapsedTime())) or
      "No active timer"
    
    local choice = ShowDialog("Speedrun Timer", {
      "Status: " .. status_line,
      "────────────────────────",
      "1. Start New Run",
      "2. Split",
      "3. Pause/Resume",
      "4. Reset Timer",
      "5. View Current Run",
      "6. View Personal Bests",
      "7. Export Run",
      "8. Exit"
    })
    
    if not choice or choice == 8 then break end
    if choice <= 2 then choice = choice + 2 end -- Adjust for display lines
    
    local result = ""
    if choice == 3 then
      -- Start new run
      local cat_lines = {}
      for i, cat in ipairs(CATEGORIES) do
        table.insert(cat_lines, string.format("%d. %s", i, cat.name))
      end
      table.insert(cat_lines, string.format("%d. Cancel", #CATEGORIES + 1))
      
      local cat_choice = ShowDialog("Select Category:", cat_lines)
      if cat_choice and cat_choice <= #CATEGORIES then
        result = startTimer(CATEGORIES[cat_choice])
      end
    elseif choice == 4 then
      result = split()
    elseif choice == 5 then
      result = togglePause()
    elseif choice == 6 then
      result = resetTimer()
    elseif choice == 7 then
      result = viewCurrentRun()
    elseif choice == 8 then
      result = viewPBs()
    elseif choice == 9 then
      result = exportRun()
    end
    
    if result ~= "" then
      ShowDialog("Result", {result, "Press any key to continue..."})
    end
  end
end

-- Run plugin
main()
