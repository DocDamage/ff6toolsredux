-- ============================================================================
-- Sabin's Blitz Input Trainer Plugin
-- ============================================================================
-- @id: blitz-trainer
-- @name: Sabin's Blitz Input Trainer
-- @version: 1.0.0
-- @author: FF6 Editor Team
-- @description: Reference and training guide for Sabin's 8 Blitz commands
-- @permissions: read_save, ui_display

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    showDamageFormulas = true,      -- Show damage calculation formulas
    showInputDiagrams = true,       -- Show ASCII input diagrams
    estimateDamage = true,          -- Calculate estimated damage
}

-- ============================================================================
-- Blitz Database
-- ============================================================================
-- Complete database of all 8 Blitzes with input sequences and effects

local BLITZ_DATABASE = {
    {
        id = 1,
        name = "Raging Fist",
        input = "Left, Right",
        inputSymbols = "← →",
        difficulty = "Easy",
        learnLevel = 1,
        damageType = "Physical",
        power = 110,
        target = "Single Enemy",
        element = "None",
        formula = "Damage = 110 + (Level * Level * Vigor) / 256",
        description = "Powerful single-target physical attack. Double damage compared to normal attack.",
        effects = "Physical damage to one enemy, ignores defense",
        mpCost = 0,
        hitRate = 100,
        strategy = "Your starting Blitz. Great for early game damage. Use against single tough enemies.",
        tips = "Easy input makes this reliable for any situation."
    },
    {
        id = 2,
        name = "Aura Cannon",
        input = "Down, Down-Left, Left",
        inputSymbols = "↓ ↙ ←",
        difficulty = "Easy",
        learnLevel = 6,
        damageType = "Physical",
        power = 130,
        target = "Single Enemy",
        element = "None",
        formula = "Damage = 130 + (Level * Level * Vigor) / 256",
        description = "Enhanced single-target attack, stronger than Raging Fist.",
        effects = "Physical damage to one enemy",
        mpCost = 0,
        hitRate = 100,
        strategy = "Upgrade to Raging Fist learned at level 6. Better damage output.",
        tips = "Diagonal input (Down-Left) can be tricky. Practice the smooth motion."
    },
    {
        id = 3,
        name = "Meteor Strike",
        input = "Left, Left-Up, Up, Right-Up, Right",
        inputSymbols = "← ↖ ↑ ↗ →",
        difficulty = "Hard",
        learnLevel = 10,
        damageType = "Physical",
        power = 150,
        target = "Single Enemy",
        element = "None",
        formula = "Damage = 150 + (Level * Level * Vigor) / 256 * 4",
        description = "Four-hit combo attack dealing massive damage to single target.",
        effects = "Physical damage to one enemy (hits 4 times)",
        mpCost = 0,
        hitRate = 100,
        strategy = "Hardest input but incredible damage. Best single-target Blitz.",
        tips = "Practice the smooth quarter-circle motion. Take your time with each direction."
    },
    {
        id = 4,
        name = "Suplex",
        input = "Up, Up, Down, Down, Left, Right",
        inputSymbols = "↑ ↑ ↓ ↓ ← →",
        difficulty = "Medium",
        learnLevel = 15,
        damageType = "Physical",
        power = 120,
        target = "Single Enemy",
        element = "None",
        formula = "Damage = 120 + (Level * Level * Vigor) / 256",
        description = "Grab and slam enemy. Works on most non-floating enemies.",
        effects = "Physical damage, can work on some bosses!",
        mpCost = 0,
        hitRate = 100,
        strategy = "Famous for working on the Phantom Train. Try it on everything!",
        tips = "Long input but all cardinal directions. No diagonals makes it manageable."
    },
    {
        id = 5,
        name = "Fire Dance",
        input = "Left, Left-Down, Down, Right-Down, Right",
        inputSymbols = "← ↙ ↓ ↘ →",
        difficulty = "Hard",
        learnLevel = 23,
        damageType = "Magical",
        power = 106,
        target = "All Enemies",
        element = "Fire",
        formula = "Damage = (106 * Level / 2) + Magic Power",
        description = "Fire-elemental AoE attack hitting all enemies.",
        effects = "Fire damage to all enemies",
        mpCost = 0,
        hitRate = 100,
        strategy = "First multi-target Blitz. Great against groups, especially ice-weak enemies.",
        tips = "Quarter-circle motion with all diagonals. Practice smooth execution."
    },
    {
        id = 6,
        name = "Mantra",
        input = "Right, Left, Right, Left, Right, Left",
        inputSymbols = "→ ← → ← → ←",
        difficulty = "Medium",
        learnLevel = 30,
        damageType = "Healing",
        power = 0,
        target = "All Allies",
        element = "None",
        formula = "Healing = (Sabin's Max HP + Sabin's Current HP) / 16",
        description = "Restore HP to all party members based on Sabin's HP.",
        effects = "Heals entire party",
        mpCost = 0,
        hitRate = 100,
        strategy = "Free party-wide healing! Keep Sabin healthy for better healing output.",
        tips = "Long input but simple pattern. Alternate right-left six times."
    },
    {
        id = 7,
        name = "Air Blade",
        input = "Up, Up-Right, Right, Right-Down, Down, Down-Left, Left",
        inputSymbols = "↑ ↗ → ↘ ↓ ↙ ←",
        difficulty = "Very Hard",
        learnLevel = 42,
        damageType = "Magical",
        power = 127,
        target = "All Enemies",
        element = "Wind",
        formula = "Damage = (127 * Level / 2) + Magic Power",
        description = "Wind-elemental AoE attack, more powerful than Fire Dance.",
        effects = "Wind damage to all enemies",
        mpCost = 0,
        hitRate = 100,
        strategy = "Strongest multi-target Blitz. Excellent for random encounters.",
        tips = "Three-quarter circle motion. Most complex input in the game."
    },
    {
        id = 8,
        name = "Phantom Rush",
        input = "Left, Up, Right, Down, Left-Down",
        inputSymbols = "← ↑ → ↓ ↙",
        difficulty = "Hard",
        learnLevel = 70,
        damageType = "Physical",
        power = 255,
        target = "Random Enemies",
        element = "None",
        formula = "Damage = Max(9999) per hit, 4 hits to random enemies",
        description = "Ultimate Blitz! Four hits dealing maximum damage to random enemies.",
        effects = "Physical damage (up to 9999 per hit, 4 times)",
        mpCost = 0,
        hitRate = 100,
        strategy = "Endgame destroyer. Can deal 39,996 total damage. Overkill for most fights.",
        tips = "Box pattern with diagonal end. Practice makes perfect!"
    }
}

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Find Sabin's character index
local function findSabin(api)
    for i = 0, 15 do
        local char = api:GetCharacter(i)
        if char and char.Name then
            if char.Name == "Sabin" or char.Name == "SABIN" or 
               char.Name:lower():find("sabin") then
                return i, char
            end
        end
    end
    return nil, nil
end

-- Estimate damage based on Sabin's stats
local function estimateBlitzDamage(blitz, sabinChar)
    if not sabinChar or not config.estimateDamage then
        return "N/A"
    end
    
    local level = sabinChar.Level or 30
    local vigor = sabinChar.Vigor or 50
    local magic = sabinChar.Magic or 35
    
    if blitz.damageType == "Physical" then
        local baseDamage = blitz.power + (level * level * vigor) / 256
        
        if blitz.name == "Meteor Strike" then
            baseDamage = baseDamage * 4  -- 4 hits
        end
        
        return string.format("~%.0f damage", baseDamage)
        
    elseif blitz.damageType == "Magical" then
        local baseDamage = (blitz.power * level / 2) + magic
        return string.format("~%.0f damage per enemy", baseDamage)
        
    elseif blitz.damageType == "Healing" then
        local maxHP = sabinChar.MaxHP or 1000
        local currentHP = sabinChar.CurrentHP or 800
        local healing = (maxHP + currentHP) / 16
        return string.format("~%.0f HP to all allies", healing)
    end
    
    return "N/A"
end

-- Format input diagram
local function formatInputDiagram(blitz)
    local lines = {}
    
    -- ASCII art representation
    table.insert(lines, "   Input Sequence:")
    table.insert(lines, "   " .. blitz.inputSymbols)
    table.insert(lines, "")
    table.insert(lines, "   Button Sequence:")
    table.insert(lines, "   " .. blitz.input)
    
    return table.concat(lines, "\n")
end

-- Format single Blitz entry
local function formatBlitzEntry(blitz, sabinChar)
    local sections = {}
    
    table.insert(sections, string.format("═══ %s ═══", blitz.name))
    table.insert(sections, "")
    
    -- Basic info
    table.insert(sections, string.format("Learn Level: %d", blitz.learnLevel))
    table.insert(sections, string.format("Difficulty: %s", blitz.difficulty))
    table.insert(sections, string.format("Target: %s", blitz.target))
    table.insert(sections, string.format("Type: %s", blitz.damageType))
    
    if blitz.element ~= "None" then
        table.insert(sections, string.format("Element: %s", blitz.element))
    end
    
    table.insert(sections, "")
    
    -- Input diagram
    if config.showInputDiagrams then
        table.insert(sections, formatInputDiagram(blitz))
        table.insert(sections, "")
    end
    
    -- Description
    table.insert(sections, "Description:")
    table.insert(sections, "  " .. blitz.description)
    table.insert(sections, "")
    
    -- Damage estimate
    if sabinChar and config.estimateDamage then
        table.insert(sections, "Estimated Damage (Current Stats):")
        table.insert(sections, "  " .. estimateBlitzDamage(blitz, sabinChar))
        table.insert(sections, "")
    end
    
    -- Formula
    if config.showDamageFormulas then
        table.insert(sections, "Damage Formula:")
        table.insert(sections, "  " .. blitz.formula)
        table.insert(sections, "")
    end
    
    -- Strategy
    table.insert(sections, "Strategy:")
    table.insert(sections, "  " .. blitz.strategy)
    table.insert(sections, "")
    
    -- Tips
    table.insert(sections, "Input Tips:")
    table.insert(sections, "  " .. blitz.tips)
    
    return table.concat(sections, "\n")
end

-- Build quick reference card
local function buildQuickReference(sabinChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "BLITZ QUICK REFERENCE CARD")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    
    for _, blitz in ipairs(BLITZ_DATABASE) do
        local damage = estimateBlitzDamage(blitz, sabinChar)
        table.insert(sections, string.format("%-17s  %-25s  %s", 
            blitz.name, blitz.inputSymbols, blitz.target))
        table.insert(sections, string.format("  Learn: Lv%2d  |  %s  |  Est: %s", 
            blitz.learnLevel, blitz.difficulty, damage))
        table.insert(sections, "")
    end
    
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n")
end

-- Build all Blitzes list
local function buildAllBlitzes(sabinChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "ALL BLITZ TECHNIQUES")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    
    for _, blitz in ipairs(BLITZ_DATABASE) do
        table.insert(sections, formatBlitzEntry(blitz, sabinChar))
        table.insert(sections, "")
        table.insert(sections, string.rep("-", 70))
        table.insert(sections, "")
    end
    
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n")
end

-- Build input guide
local function buildInputGuide()
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "BLITZ INPUT SYSTEM GUIDE")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    table.insert(sections, "How to Perform Blitzes:")
    table.insert(sections, "")
    table.insert(sections, "1. Select Sabin in battle")
    table.insert(sections, "2. Choose 'Blitz' command")
    table.insert(sections, "3. Input the button sequence (you have about 5 seconds)")
    table.insert(sections, "4. Press A/Confirm to execute")
    table.insert(sections, "")
    table.insert(sections, "Direction Inputs:")
    table.insert(sections, "  ↑ = Up        ↗ = Up-Right      → = Right")
    table.insert(sections, "  ↘ = Down-Right ↓ = Down         ↙ = Down-Left")
    table.insert(sections, "  ← = Left      ↖ = Up-Left")
    table.insert(sections, "")
    table.insert(sections, "Input Tips:")
    table.insert(sections, "  • Take your time - you have ~5 seconds")
    table.insert(sections, "  • Smooth motions work better than button mashing")
    table.insert(sections, "  • Diagonal inputs require precision")
    table.insert(sections, "  • Practice makes perfect!")
    table.insert(sections, "  • If you mess up, the command fails (wastes turn)")
    table.insert(sections, "")
    table.insert(sections, "Learning Blitzes:")
    table.insert(sections, "  • Sabin learns new Blitzes as he levels up")
    table.insert(sections, "  • Once learned, they're available forever")
    table.insert(sections, "  • No MP cost for any Blitz!")
    table.insert(sections, "  • All Blitzes ignore defense/magic defense")
    table.insert(sections, "")
    table.insert(sections, "Common Input Patterns:")
    table.insert(sections, "  • Horizontal: ← → or → ←")
    table.insert(sections, "  • Quarter-circle: ↓ ↙ ← or ← ↙ ↓ ↘ →")
    table.insert(sections, "  • Semi-circle: ← ↖ ↑ ↗ →")
    table.insert(sections, "  • Three-quarter: ↑ ↗ → ↘ ↓ ↙ ←")
    table.insert(sections, "  • Box: ← ↑ → ↓")
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n")
end

-- Build damage calculator
local function buildDamageCalculator(sabinChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "BLITZ DAMAGE CALCULATOR")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    
    if not sabinChar then
        table.insert(sections, "Sabin not found. Cannot calculate damage.")
    else
        local level = sabinChar.Level or 30
        local vigor = sabinChar.Vigor or 50
        local magic = sabinChar.Magic or 35
        local maxHP = sabinChar.MaxHP or 1000
        local currentHP = sabinChar.CurrentHP or 800
        
        table.insert(sections, "Sabin's Current Stats:")
        table.insert(sections, string.format("  Level: %d", level))
        table.insert(sections, string.format("  Vigor: %d (Physical Power)", vigor))
        table.insert(sections, string.format("  Magic: %d (Magical Power)", magic))
        table.insert(sections, string.format("  HP: %d / %d", currentHP, maxHP))
        table.insert(sections, "")
        table.insert(sections, "Estimated Blitz Damage:")
        table.insert(sections, string.rep("-", 70))
        
        for _, blitz in ipairs(BLITZ_DATABASE) do
            local damage = estimateBlitzDamage(blitz, sabinChar)
            table.insert(sections, string.format("%-17s: %s", blitz.name, damage))
        end
    end
    
    table.insert(sections, "")
    table.insert(sections, "Note: Actual damage varies by enemy defense and modifiers")
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n")
end

-- Build main menu
local function buildMainMenu(sabinChar)
    local level = sabinChar and sabinChar.Level or "?"
    
    local menu = "=" .. string.rep("=", 70) .. "\n"
    menu = menu .. "SABIN'S BLITZ INPUT TRAINER - Main Menu\n"
    menu = menu .. "=" .. string.rep("=", 70) .. "\n\n"
    menu = menu .. string.format("Sabin's Level: %s\n\n", level)
    menu = menu .. "Select an option:\n\n"
    menu = menu .. "1. Quick Reference Card (All Blitzes Overview)\n"
    menu = menu .. "2. View Individual Blitz (Choose Specific)\n"
    menu = menu .. "3. View All Blitzes (Complete Details)\n"
    menu = menu .. "4. Input System Guide (How Blitzes Work)\n"
    menu = menu .. "5. Damage Calculator (Current Stats)\n"
    menu = menu .. "6. Blitz Recommendations (Best for Each Situation)\n"
    menu = menu .. "0. Exit\n\n"
    menu = menu .. "Enter your choice (0-6):"
    
    return menu
end

-- Build Blitz selection menu
local function buildBlitzSelection()
    local menu = "=" .. string.rep("=", 70) .. "\n"
    menu = menu .. "SELECT A BLITZ\n"
    menu = menu .. "=" .. string.rep("=", 70) .. "\n\n"
    
    for i, blitz in ipairs(BLITZ_DATABASE) do
        menu = menu .. string.format("%d. %-17s (Learn: Lv%2d, %s)\n", 
            i, blitz.name, blitz.learnLevel, blitz.difficulty)
    end
    
    menu = menu .. "0. Back to Main Menu\n\n"
    menu = menu .. "Enter your choice (0-8):"
    
    return menu
end

-- Build recommendations
local function buildRecommendations(sabinChar)
    local sections = {}
    
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "BLITZ RECOMMENDATIONS")
    table.insert(sections, "=" .. string.rep("=", 70))
    table.insert(sections, "")
    table.insert(sections, "Best Blitz for Each Situation:")
    table.insert(sections, "")
    table.insert(sections, "★ Early Game (Levels 1-15)")
    table.insert(sections, "  Raging Fist - Simple input, reliable damage")
    table.insert(sections, "  Aura Cannon - Better damage at level 6")
    table.insert(sections, "")
    table.insert(sections, "★ Mid Game (Levels 16-40)")
    table.insert(sections, "  Suplex - Fun and effective (try it on everything!)")
    table.insert(sections, "  Fire Dance - First AoE attack for groups")
    table.insert(sections, "  Mantra - Free healing for the whole party")
    table.insert(sections, "")
    table.insert(sections, "★ Late Game (Levels 41+)")
    table.insert(sections, "  Meteor Strike - Highest single-target damage")
    table.insert(sections, "  Air Blade - Best AoE for random encounters")
    table.insert(sections, "  Phantom Rush - Ultimate attack (overkill for most)")
    table.insert(sections, "")
    table.insert(sections, "★ Boss Battles")
    table.insert(sections, "  Best: Meteor Strike (massive single-target damage)")
    table.insert(sections, "  Alternative: Suplex (yes, it works on some bosses!)")
    table.insert(sections, "")
    table.insert(sections, "★ Random Encounters")
    table.insert(sections, "  Best: Air Blade (wind AoE damage)")
    table.insert(sections, "  Alternative: Fire Dance (fire AoE damage)")
    table.insert(sections, "")
    table.insert(sections, "★ Emergency Healing")
    table.insert(sections, "  Mantra - Free party-wide healing")
    table.insert(sections, "  (Keep Sabin's HP high for better healing)")
    table.insert(sections, "")
    table.insert(sections, "★ Easiest Inputs (Learning Stage)")
    table.insert(sections, "  1. Raging Fist (← →)")
    table.insert(sections, "  2. Aura Cannon (↓ ↙ ←)")
    table.insert(sections, "  3. Suplex (↑ ↑ ↓ ↓ ← →)")
    table.insert(sections, "")
    table.insert(sections, "★ Hardest Inputs (Master Stage)")
    table.insert(sections, "  1. Air Blade (↑ ↗ → ↘ ↓ ↙ ←)")
    table.insert(sections, "  2. Meteor Strike (← ↖ ↑ ↗ →)")
    table.insert(sections, "  3. Fire Dance (← ↙ ↓ ↘ →)")
    table.insert(sections, "")
    table.insert(sections, "=" .. string.rep("=", 70))
    
    return table.concat(sections, "\n")
end

-- ============================================================================
-- Main Entry Point
-- ============================================================================
function main(api)
    -- Verify permissions
    if not api:HasPermission("read_save") then
        api:Log("error", "Missing required permission: read_save")
        api:ShowDialog("Permission Error", "This plugin requires 'read_save' permission.")
        return false
    end
    
    if not api:HasPermission("ui_display") then
        api:Log("error", "Missing required permission: ui_display")
        return false
    end
    
    api:Log("info", "Blitz Input Trainer plugin started")
    
    -- Find Sabin (optional - plugin works as reference without him)
    local sabinIndex, sabinChar = findSabin(api)
    if sabinIndex then
        api:Log("info", string.format("Found Sabin at character index: %d", sabinIndex))
    else
        api:Log("info", "Sabin not found - running in reference mode")
    end
    
    -- Main menu loop
    while true do
        local menu = buildMainMenu(sabinChar)
        local choice = api:ShowInput("Blitz Trainer - Main Menu", menu, "1")
        
        if not choice or choice == "0" then
            api:Log("info", "User exited Blitz Trainer")
            return true
        end
        
        if choice == "1" then
            -- Quick reference
            local ref = buildQuickReference(sabinChar)
            api:ShowDialog("Blitz Quick Reference", ref)
            
        elseif choice == "2" then
            -- Individual Blitz
            local selMenu = buildBlitzSelection()
            local blitzChoice = api:ShowInput("Select Blitz", selMenu, "1")
            
            if blitzChoice and blitzChoice ~= "0" then
                local blitzNum = tonumber(blitzChoice)
                if blitzNum and blitzNum >= 1 and blitzNum <= #BLITZ_DATABASE then
                    local blitz = BLITZ_DATABASE[blitzNum]
                    local detail = formatBlitzEntry(blitz, sabinChar)
                    api:ShowDialog(blitz.name, detail)
                end
            end
            
        elseif choice == "3" then
            -- All Blitzes
            local allBlitzes = buildAllBlitzes(sabinChar)
            api:ShowDialog("All Blitz Techniques", allBlitzes)
            
        elseif choice == "4" then
            -- Input guide
            local guide = buildInputGuide()
            api:ShowDialog("Blitz Input System", guide)
            
        elseif choice == "5" then
            -- Damage calculator
            local calc = buildDamageCalculator(sabinChar)
            api:ShowDialog("Damage Calculator", calc)
            
        elseif choice == "6" then
            -- Recommendations
            local recs = buildRecommendations(sabinChar)
            api:ShowDialog("Blitz Recommendations", recs)
            
        else
            api:ShowDialog("Invalid Choice", "Please enter a number between 0 and 6.")
        end
    end
    
    return true
end

-- ============================================================================
-- Error Handling Wrapper
-- ============================================================================
local function safeMain(api)
    local success, result = pcall(main, api)
    if not success then
        local errorMsg = tostring(result)
        api:Log("error", "Plugin error: " .. errorMsg)
        api:ShowDialog("Plugin Error", 
            "An error occurred while running the plugin:\n\n" .. errorMsg .. 
            "\n\nCheck the plugin log for details.")
        return false
    end
    return result
end

-- Execute with error handling
return safeMain(api)
