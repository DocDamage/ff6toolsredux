--[[
    Damage Calculator
    Version: 1.0.0

    Calculates approximate physical and magical damage using common FF6-style
    formulas, with elemental modifiers, row/back attack, and quick scenario
    comparisons. Designed to be read-only and resilient to missing APIs: if
    stats or enemies are unavailable, it uses manual inputs or shows Unknown.

    Permissions: read_save, ui_display
]]

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    defaultEnemyDefense = 120,
    defaultEnemyMagDefense = 150,
    defaultWeaponPower = 120,
    defaultSpellPower = 50,
    defaultLevel = 50,
    elementalWeakMultiplier = 1.5,
    elementalResistMultiplier = 0.5,
    backAttackPenalty = 0.5,
    rowPenalty = 0.5,
    criticalMultiplier = 2.0,
    multiTargetPenalty = 0.75,
    pageSize = 15,
}

-- ============================================================================
-- Helper functions
-- ============================================================================

local function safeCall(api, method, ...)
    if not api or not method then return nil end
    local ok, result = pcall(function(...)
        local fn = api[method]
        if type(fn) == "function" then
            return fn(api, ...)
        end
        return nil
    end, ...)
    if not ok then return nil end
    return result
end

local function toNumber(v, default)
    local n = tonumber(v)
    if not n then return default end
    return n
end

local function promptNumber(api, title, message, default)
    local input = api:ShowInput(title, message, tostring(default))
    if not input then return nil end
    return toNumber(input, default)
end

local function formatResult(title, value)
    return string.format("%s: %.0f", title, value)
end

-- Basic physical damage approximation
local function calcPhysical(level, vigor, weaponPower, enemyDef, row, backAttack, critical)
    level = level or config.defaultLevel
    vigor = vigor or 40
    weaponPower = weaponPower or config.defaultWeaponPower
    enemyDef = enemyDef or config.defaultEnemyDefense

    local attack = math.floor(((level * level * vigor) / 256) + weaponPower)
    local damage = attack - (enemyDef / 2)
    if damage < 0 then damage = 0 end
    if row == "back" then damage = damage * config.rowPenalty end
    if backAttack then damage = damage * config.backAttackPenalty end
    if critical then damage = damage * config.criticalMultiplier end
    return math.max(0, math.floor(damage))
end

-- Basic magical damage approximation
local function calcMagic(level, magic, spellPower, enemyMDef, elementalMod, multiTarget)
    level = level or config.defaultLevel
    magic = magic or 40
    spellPower = spellPower or config.defaultSpellPower
    enemyMDef = enemyMDef or config.defaultEnemyMagDefense
    elementalMod = elementalMod or 1.0

    local base = ((level + spellPower) ^ 2) / 32 + magic * 4
    local damage = base * elementalMod - enemyMDef
    if multiTarget then damage = damage * config.multiTargetPenalty end
    return math.max(0, math.floor(damage))
end

local function applyElement(choice)
    if choice == "weak" then return config.elementalWeakMultiplier
    elseif choice == "resist" then return config.elementalResistMultiplier
    elseif choice == "absorb" then return -1.0
    elseif choice == "null" then return 0
    end
    return 1.0
end

-- ============================================================================
-- Enemy fetch (lightweight; uses bestiary if present)
-- ============================================================================
local function pickEnemy(api)
    local enemyIdStr = api:ShowInput("Enemy", "Enter enemy ID (optional, for tracking):", "0")
    if not enemyIdStr then return nil end
    local enemyId = tonumber(enemyIdStr)
    if not enemyId or enemyId <= 0 then return nil end
    local enemy = safeCall(api, "GetBestiaryEntry", enemyId)
    return enemy, enemyId
end

local function enemyDefense(enemy)
    if not enemy then return config.defaultEnemyDefense end
    return enemy.Defense or enemy.defense or config.defaultEnemyDefense
end

local function enemyMagDefense(enemy)
    if not enemy then return config.defaultEnemyMagDefense end
    return enemy.MagDefense or enemy.magDefense or enemy.MagicDefense or config.defaultEnemyMagDefense
end

-- ============================================================================
-- Views
-- ============================================================================

local function physicalView(api)
    local enemy, enemyId = pickEnemy(api)
    local level = promptNumber(api, "Level", "Enter attacker level:", config.defaultLevel)
    if not level then return end
    local vigor = promptNumber(api, "Vigor", "Enter Vigor:", 40)
    if not vigor then return end
    local weaponPower = promptNumber(api, "Weapon Power", "Enter weapon power:", config.defaultWeaponPower)
    if not weaponPower then return end
    local def = promptNumber(api, "Enemy Defense", "Enter enemy defense (leave blank to use default/enemy):", enemyDefense(enemy))
    if not def then def = enemyDefense(enemy) end

    local row = api:ShowInput("Row", "Enter row (front/back):", "front") or "front"
    local backAttack = (api:ShowInput("Back Attack", "Is this a back attack? (y/n)", "n") or "n"):lower() == "y"
    local crit = (api:ShowInput("Critical", "Critical hit? (y/n)", "n") or "n"):lower() == "y"

    local dmg = calcPhysical(level, vigor, weaponPower, def, row == "back" and "back" or "front", backAttack, crit)

    local lines = {
        string.format("Enemy ID: %s", enemyId or "(none)"),
        formatResult("Physical Damage", dmg),
        string.format("Inputs: L%d, Vigor %d, WP %d, Def %d, Row %s, Back %s, Crit %s",
            level, vigor, weaponPower, def, row, backAttack and "Y" or "N", crit and "Y" or "N")
    }
    api:ShowDialog("Physical Damage", table.concat(lines, "\n"))
end

local function magicView(api)
    local enemy, enemyId = pickEnemy(api)
    local level = promptNumber(api, "Level", "Enter caster level:", config.defaultLevel)
    if not level then return end
    local magic = promptNumber(api, "Magic", "Enter Magic stat:", 40)
    if not magic then return end
    local spellPower = promptNumber(api, "Spell Power", "Enter spell power:", config.defaultSpellPower)
    if not spellPower then return end
    local mdef = promptNumber(api, "Enemy Mag Def", "Enter enemy magic defense (blank to use default/enemy):", enemyMagDefense(enemy))
    if not mdef then mdef = enemyMagDefense(enemy) end

    local elem = (api:ShowInput("Element", "Target relationship (weak/resist/absorb/null/neutral):", "neutral") or "neutral"):lower()
    local elementalMod = applyElement(elem)
    local multiTarget = (api:ShowInput("Multi-target", "Is this multi-target? (y/n)", "n") or "n"):lower() == "y"

    local dmg = calcMagic(level, magic, spellPower, mdef, elementalMod, multiTarget)

    local lines = {
        string.format("Enemy ID: %s", enemyId or "(none)"),
        formatResult("Magic Damage", dmg),
        string.format("Inputs: L%d, Mag %d, Pow %d, MDef %d, Elem %s, Multi %s",
            level, magic, spellPower, mdef, elem, multiTarget and "Y" or "N")
    }
    api:ShowDialog("Magic Damage", table.concat(lines, "\n"))
end

local function compareView(api)
    local enemy, enemyId = pickEnemy(api)
    local level = promptNumber(api, "Level", "Enter attacker level:", config.defaultLevel)
    if not level then return end
    local vigor = promptNumber(api, "Vigor", "Enter Vigor:", 40)
    if not vigor then return end
    local magic = promptNumber(api, "Magic", "Enter Magic stat:", 40)
    if not magic then return end
    local weaponPower = promptNumber(api, "Weapon Power", "Enter weapon power:", config.defaultWeaponPower)
    if not weaponPower then return end
    local spellPower = promptNumber(api, "Spell Power", "Enter spell power:", config.defaultSpellPower)
    if not spellPower then return end
    local def = enemyDefense(enemy)
    local mdef = enemyMagDefense(enemy)

    local phys = calcPhysical(level, vigor, weaponPower, def, "front", false, false)
    local magNeutral = calcMagic(level, magic, spellPower, mdef, 1.0, false)
    local magWeak = calcMagic(level, magic, spellPower, mdef, config.elementalWeakMultiplier, false)

    local lines = {
        string.format("Enemy ID: %s", enemyId or "(none)"),
        formatResult("Physical (front)", phys),
        formatResult("Magic (neutral)", magNeutral),
        formatResult("Magic (weak)", magWeak),
        string.format("Inputs: L%d, Vig %d, Mag %d, WP %d, Pow %d, Def %d, MDef %d",
            level, vigor, magic, weaponPower, spellPower, def, mdef)
    }
    api:ShowDialog("Compare Physical vs Magic", table.concat(lines, "\n"))
end

local function helpView(api)
    local lines = {
        "Damage Calculator (approximate FF6-style)",
        "- Physical: uses level^2 * vigor / 256 + weapon power - enemy def/2, with row/back/crit modifiers.",
        "- Magic: uses ((level + power)^2)/32 + magic*4 - enemy mdef, with element and multi-target modifiers.",
        "- Element multipliers: weak 1.5x, resist 0.5x, absorb -1x, null 0x (configurable).",
        "- Unknown enemy stats: defaults from config; you can override in prompts.",
        "- Read-only: no save mutations.",
    }
    api:ShowDialog("Help", table.concat(lines, "\n"))
end

-- ============================================================================
-- Menu Loop
-- ============================================================================

local function mainMenu(api)
    while true do
        local choice = api:ShowInput(
            "Damage Calculator",
            "Select an option:\n1) Physical\n2) Magic\n3) Compare Phys vs Magic\n4) Help\n0) Exit",
            "1"
        )
        if not choice then return end
        local sel = tonumber(choice)
        if sel == 1 then
            physicalView(api)
        elseif sel == 2 then
            magicView(api)
        elseif sel == 3 then
            compareView(api)
        elseif sel == 4 then
            helpView(api)
        elseif sel == 0 then
            return
        else
            api:ShowDialog("Invalid", "Please choose a valid option.")
        end
    end
end

local function safeMain(api)
    local ok, err = pcall(function()
        mainMenu(api)
    end)
    if not ok and api and api.ShowDialog then
        api:ShowDialog("Error", "Damage Calculator failed: " .. tostring(err))
    end
end

return {
    id = "damage-calculator",
    name = "Damage Calculator",
    version = "1.0.0",
    author = "FF6 Save Editor Team",
    permissions = {"read_save", "ui_display"},
    run = safeMain
}
