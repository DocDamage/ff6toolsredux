--[[
  Combat Depth Pack - Save Manipulation Examples
  Demonstrates how to use the save bindings exposed to Lua
]]

-- Character Manipulation Examples
local examples = {}

-- Example 1: Boost all party members to level 50
function examples.boost_party()
  local count = save.getCharacterCount()
  local boosted = {}
  
  for i = 0, math.min(count - 1, 3) do  -- First 4 characters
    local name = save.getCharacterName(i)
    if name and name ~= "Unknown" then
      local success = save.setCharacterLevel(i, 50)
      if success then
        save.setCharacterHP(i, 9999)
        save.setCharacterMP(i, 999)
        table.insert(boosted, name)
        save.log("Boosted " .. name .. " to level 50")
      end
    end
  end
  
  return {
    success = true,
    characters_boosted = boosted,
    count = #boosted
  }
end

-- Example 2: Give starting gil bonus
function examples.give_starting_bonus()
  local current = save.getGil()
  local bonus = 50000
  save.setGil(current + bonus)
  save.log("Added " .. bonus .. " gil (was " .. current .. ", now " .. (current + bonus) .. ")")
  
  return {
    success = true,
    previous_gil = current,
    bonus = bonus,
    new_gil = current + bonus
  }
end

-- Example 3: Emergency heal for all characters
function examples.emergency_heal()
  local count = save.getCharacterCount()
  local healed = {}
  
  for i = 0, count - 1 do
    local name = save.getCharacterName(i)
    if name and name ~= "Unknown" then
      save.setCharacterHP(i, 9999)
      save.setCharacterMP(i, 999)
      table.insert(healed, name)
    end
  end
  
  return {
    success = true,
    healed_count = #healed,
    characters = healed
  }
end

-- Example 4: Integration with Combat Depth Pack
function examples.prepare_for_boss_fight()
  -- Heal party
  local heal_result = examples.emergency_heal()
  
  -- Load Combat Pack
  local pack = require('plugins.combat-depth-pack.v1_0_core')
  
  -- Generate boss remix
  local boss_plan = pack.BossRemix.generateRemix(
    {name = "Atma", hp = 50000, attack = 1200},
    {affixes = {"enraged", "arcane_shield"}}
  )
  
  -- Get AI recommendations
  local ai_action = pack.CompanionDirector.recommendAction(
    {hp_status = "stable", threat_level = "high", primary_target = "boss"},
    "aggressive"
  )
  
  return {
    success = true,
    party_healed = heal_result.healed_count,
    boss_setup = boss_plan,
    ai_recommendation = ai_action
  }
end

return examples
