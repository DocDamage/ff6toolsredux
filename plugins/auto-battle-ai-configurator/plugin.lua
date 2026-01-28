-- ============================================================================
-- AUTO-BATTLE AI CONFIGURATOR PLUGIN
-- ============================================================================
-- Purpose: Plan and configure custom AI battle strategies (planning tool)
-- Category: Experimental / Strategy Planning
-- Phase: 6, Batch: 4
-- Complexity: Advanced
-- Lines of Code: ~750
-- ============================================================================

-- ⚠️ IMPORTANT: PLANNING TOOL ONLY ⚠️
-- This plugin is a STRATEGY PLANNING and ANALYSIS tool.
-- It does NOT modify in-game AI behavior (FF6 limitation).
-- Use this tool to:
--   - Design battle strategies
--   - Plan conditional actions
--   - Calculate optimal rotations
--   - Export strategies for manual play
--   - Analyze character abilities and stats
-- 
-- Strategies created here are for planning and reference only.
-- You must manually execute strategies in-game.

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    MAX_STRATEGIES = 50,
    MAX_CONDITIONS_PER_STRATEGY = 20,
    MAX_ACTIONS_PER_PRIORITY = 10,
    
    -- Priority levels
    PRIORITY = {
        CRITICAL = 1,    -- Life-threatening situations
        HIGH = 2,        -- Important actions
        MEDIUM = 3,      -- Normal actions
        LOW = 4,         -- Filler actions
        IDLE = 5,        -- Nothing else to do
    },
    
    -- Condition operators
    OPERATORS = {
        EQUAL = "==",
        NOT_EQUAL = "~=",
        GREATER = ">",
        LESS = "<",
        GREATER_EQUAL = ">=",
        LESS_EQUAL = "<=",
    },
    
    -- Target types
    TARGETS = {
        SELF = "self",
        ALLY = "ally",
        ENEMY = "enemy",
        ALL_ALLIES = "all_allies",
        ALL_ENEMIES = "all_enemies",
        LOWEST_HP_ALLY = "lowest_hp_ally",
        HIGHEST_HP_ENEMY = "highest_hp_enemy",
        RANDOM = "random",
    },
    
    -- Action types
    ACTION_TYPES = {
        ATTACK = "attack",
        MAGIC = "magic",
        ITEM = "item",
        SKILL = "skill",
        DEFEND = "defend",
        ROW_CHANGE = "row_change",
    },
    
    TOTAL_CHARACTERS = 14,
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local AutoBattleAI = {
    version = "1.0.0",
    strategies = {},
    strategy_count = 0,
    performance_stats = {},
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

--- Safe API call wrapper
local function safeCall(func, ...)
    local success, result = pcall(func, ...)
    if not success then
        print("[ERROR] API call failed: " .. tostring(result))
        return nil, result
    end
    return result
end

--- Generate unique strategy ID
local function generateStrategyId()
    AutoBattleAI.strategy_count = AutoBattleAI.strategy_count + 1
    return string.format("strategy_%d_%d", os.time(), AutoBattleAI.strategy_count)
end

--- Validate condition
local function validateCondition(condition)
    if not condition.variable then
        return false, "Condition missing variable"
    end
    if not condition.operator then
        return false, "Condition missing operator"
    end
    if condition.value == nil then
        return false, "Condition missing value"
    end
    return true
end

--- Evaluate condition (simulation)
local function evaluateCondition(condition, context)
    -- Simulated evaluation based on context
    local variable_value = context[condition.variable]
    if not variable_value then
        return false
    end
    
    local op = condition.operator
    local target = condition.value
    
    if op == "==" then
        return variable_value == target
    elseif op == "~=" then
        return variable_value ~= target
    elseif op == ">" then
        return variable_value > target
    elseif op == "<" then
        return variable_value < target
    elseif op == ">=" then
        return variable_value >= target
    elseif op == "<=" then
        return variable_value <= target
    end
    
    return false
end

-- ============================================================================
-- STRATEGY DEFINITION
-- ============================================================================

--- Create new AI strategy
function AutoBattleAI.createStrategy(config)
    if not config or not config.name then
        print("[ERROR] Strategy requires a name")
        return nil
    end
    
    if AutoBattleAI.strategy_count >= CONFIG.MAX_STRATEGIES then
        print("[ERROR] Maximum strategies reached")
        return nil
    end
    
    local strategy = {
        id = generateStrategyId(),
        name = config.name,
        description = config.description or "",
        character_id = config.character_id,
        conditions = {},
        actions_by_priority = {
            [CONFIG.PRIORITY.CRITICAL] = {},
            [CONFIG.PRIORITY.HIGH] = {},
            [CONFIG.PRIORITY.MEDIUM] = {},
            [CONFIG.PRIORITY.LOW] = {},
            [CONFIG.PRIORITY.IDLE] = {},
        },
        mp_conservation = config.mp_conservation or false,
        mp_threshold = config.mp_threshold or 50,
        item_usage_enabled = config.item_usage_enabled or false,
        created_at = os.time(),
    }
    
    AutoBattleAI.strategies[strategy.id] = strategy
    
    print(string.format("[AI Configurator] Strategy '%s' created (ID: %s)", strategy.name, strategy.id))
    return strategy.id
end

--- Add condition to strategy
function AutoBattleAI.addCondition(strategy_id, condition)
    local strategy = AutoBattleAI.strategies[strategy_id]
    if not strategy then
        print("[ERROR] Strategy not found")
        return false
    end
    
    local valid, error = validateCondition(condition)
    if not valid then
        print("[ERROR] Invalid condition: " .. tostring(error))
        return false
    end
    
    if #strategy.conditions >= CONFIG.MAX_CONDITIONS_PER_STRATEGY then
        print("[ERROR] Maximum conditions reached for strategy")
        return false
    end
    
    table.insert(strategy.conditions, condition)
    return true
end

--- Add action to strategy at priority level
function AutoBattleAI.addAction(strategy_id, priority, action)
    local strategy = AutoBattleAI.strategies[strategy_id]
    if not strategy then
        print("[ERROR] Strategy not found")
        return false
    end
    
    if not action.type or not action.name then
        print("[ERROR] Action requires type and name")
        return false
    end
    
    local actions = strategy.actions_by_priority[priority]
    if not actions then
        print("[ERROR] Invalid priority level")
        return false
    end
    
    if #actions >= CONFIG.MAX_ACTIONS_PER_PRIORITY then
        print("[ERROR] Maximum actions reached for priority level")
        return false
    end
    
    action.target_type = action.target_type or CONFIG.TARGETS.ENEMY
    table.insert(actions, action)
    
    return true
end

-- ============================================================================
-- PRE-CONFIGURED STRATEGY TEMPLATES
-- ============================================================================

--- Create healer strategy template
function AutoBattleAI.createHealerStrategy(character_id)
    local strategy_id = AutoBattleAI.createStrategy({
        name = "Healer",
        description = "Dedicated healing and support role",
        character_id = character_id,
        mp_conservation = true,
        mp_threshold = 30,
        item_usage_enabled = true,
    })
    
    if not strategy_id then return nil end
    
    -- CRITICAL: Ally below 20% HP
    AutoBattleAI.addCondition(strategy_id, {
        variable = "ally_hp_percent",
        operator = "<",
        value = 20,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.CRITICAL, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Cure3",
        target_type = CONFIG.TARGETS.LOWEST_HP_ALLY,
    })
    
    -- HIGH: Ally below 50% HP
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.HIGH, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Cure2",
        target_type = CONFIG.TARGETS.LOWEST_HP_ALLY,
    })
    
    -- MEDIUM: Apply status protection
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Protect",
        target_type = CONFIG.TARGETS.ALL_ALLIES,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Shell",
        target_type = CONFIG.TARGETS.ALL_ALLIES,
    })
    
    -- LOW: Damage dealing
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.LOW, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Fire2",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    -- IDLE: Attack
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.IDLE, {
        type = CONFIG.ACTION_TYPES.ATTACK,
        name = "Attack",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    print(string.format("[AI Configurator] Healer strategy created for character %d", character_id))
    return strategy_id
end

--- Create attacker strategy template
function AutoBattleAI.createAttackerStrategy(character_id)
    local strategy_id = AutoBattleAI.createStrategy({
        name = "Physical Attacker",
        description = "High damage physical attacks",
        character_id = character_id,
        mp_conservation = true,
        mp_threshold = 20,
    })
    
    if not strategy_id then return nil end
    
    -- CRITICAL: Self below 30% HP
    AutoBattleAI.addCondition(strategy_id, {
        variable = "self_hp_percent",
        operator = "<",
        value = 30,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.CRITICAL, {
        type = CONFIG.ACTION_TYPES.ITEM,
        name = "X-Potion",
        target_type = CONFIG.TARGETS.SELF,
    })
    
    -- HIGH: Boss/strong enemy
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.HIGH, {
        type = CONFIG.ACTION_TYPES.SKILL,
        name = "Jump",
        target_type = CONFIG.TARGETS.HIGHEST_HP_ENEMY,
    })
    
    -- MEDIUM: Normal attack
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.ATTACK,
        name = "Attack",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    -- LOW: Switch to back row if low HP
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.LOW, {
        type = CONFIG.ACTION_TYPES.ROW_CHANGE,
        name = "BackRow",
        target_type = CONFIG.TARGETS.SELF,
    })
    
    print(string.format("[AI Configurator] Attacker strategy created for character %d", character_id))
    return strategy_id
end

--- Create mage strategy template
function AutoBattleAI.createMageStrategy(character_id)
    local strategy_id = AutoBattleAI.createStrategy({
        name = "Offensive Mage",
        description = "Magic-based damage dealer",
        character_id = character_id,
        mp_conservation = true,
        mp_threshold = 40,
    })
    
    if not strategy_id then return nil end
    
    -- CRITICAL: Low MP
    AutoBattleAI.addCondition(strategy_id, {
        variable = "self_mp_percent",
        operator = "<",
        value = 25,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.CRITICAL, {
        type = CONFIG.ACTION_TYPES.ITEM,
        name = "Ether",
        target_type = CONFIG.TARGETS.SELF,
    })
    
    -- HIGH: Exploit weakness
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.HIGH, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Ultima",
        target_type = CONFIG.TARGETS.ALL_ENEMIES,
    })
    
    -- MEDIUM: Elemental magic
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Fire3",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Ice3",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.MAGIC,
        name = "Bolt3",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    -- LOW: Basic attack if MP conservation needed
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.LOW, {
        type = CONFIG.ACTION_TYPES.ATTACK,
        name = "Attack",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    print(string.format("[AI Configurator] Mage strategy created for character %d", character_id))
    return strategy_id
end

--- Create tank strategy template
function AutoBattleAI.createTankStrategy(character_id)
    local strategy_id = AutoBattleAI.createStrategy({
        name = "Tank/Defender",
        description = "Absorb damage and protect allies",
        character_id = character_id,
        mp_conservation = false,
        item_usage_enabled = true,
    })
    
    if not strategy_id then return nil end
    
    -- CRITICAL: Self very low HP
    AutoBattleAI.addCondition(strategy_id, {
        variable = "self_hp_percent",
        operator = "<",
        value = 25,
    })
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.CRITICAL, {
        type = CONFIG.ACTION_TYPES.DEFEND,
        name = "Defend",
        target_type = CONFIG.TARGETS.SELF,
    })
    
    -- HIGH: Taunt enemies
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.HIGH, {
        type = CONFIG.ACTION_TYPES.SKILL,
        name = "Provoke",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    -- MEDIUM: Counter attack
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.MEDIUM, {
        type = CONFIG.ACTION_TYPES.ATTACK,
        name = "Attack",
        target_type = CONFIG.TARGETS.ENEMY,
    })
    
    -- LOW: Heal self
    AutoBattleAI.addAction(strategy_id, CONFIG.PRIORITY.LOW, {
        type = CONFIG.ACTION_TYPES.ITEM,
        name = "Hi-Potion",
        target_type = CONFIG.TARGETS.SELF,
    })
    
    print(string.format("[AI Configurator] Tank strategy created for character %d", character_id))
    return strategy_id
end

-- ============================================================================
-- STRATEGY ANALYSIS
-- ============================================================================

--- Analyze strategy effectiveness (simulation)
function AutoBattleAI.analyzeStrategy(strategy_id, simulation_context)
    local strategy = AutoBattleAI.strategies[strategy_id]
    if not strategy then
        print("[ERROR] Strategy not found")
        return nil
    end
    
    local analysis = {
        strategy_id = strategy_id,
        strategy_name = strategy.name,
        simulated_turns = 0,
        action_distribution = {
            [CONFIG.PRIORITY.CRITICAL] = 0,
            [CONFIG.PRIORITY.HIGH] = 0,
            [CONFIG.PRIORITY.MEDIUM] = 0,
            [CONFIG.PRIORITY.LOW] = 0,
            [CONFIG.PRIORITY.IDLE] = 0,
        },
        estimated_dps = 0,
        estimated_hps = 0,
        mp_efficiency = 0,
    }
    
    -- Simulate 100 turns
    for turn = 1, 100 do
        analysis.simulated_turns = analysis.simulated_turns + 1
        
        -- Check conditions and select action
        local selected_priority = CONFIG.PRIORITY.IDLE
        for _, condition in ipairs(strategy.conditions) do
            if evaluateCondition(condition, simulation_context or {}) then
                selected_priority = CONFIG.PRIORITY.CRITICAL
                break
            end
        end
        
        analysis.action_distribution[selected_priority] = analysis.action_distribution[selected_priority] + 1
    end
    
    -- Calculate metrics
    analysis.mp_efficiency = (strategy.mp_conservation and 80 or 50)
    analysis.estimated_dps = analysis.action_distribution[CONFIG.PRIORITY.HIGH] * 500
    analysis.estimated_hps = analysis.action_distribution[CONFIG.PRIORITY.CRITICAL] * 300
    
    return analysis
end

--- Get strategy statistics
function AutoBattleAI.getStrategyStats(strategy_id)
    local strategy = AutoBattleAI.strategies[strategy_id]
    if not strategy then
        return nil
    end
    
    local stats = {
        name = strategy.name,
        conditions_count = #strategy.conditions,
        actions_count = 0,
        mp_conservation = strategy.mp_conservation,
        mp_threshold = strategy.mp_threshold,
        created_at = strategy.created_at,
    }
    
    for priority, actions in pairs(strategy.actions_by_priority) do
        stats.actions_count = stats.actions_count + #actions
    end
    
    return stats
end

-- ============================================================================
-- STRATEGY MANAGEMENT
-- ============================================================================

--- List all strategies
function AutoBattleAI.listStrategies()
    local list = {}
    for id, strategy in pairs(AutoBattleAI.strategies) do
        table.insert(list, {
            id = id,
            name = strategy.name,
            description = strategy.description,
            character_id = strategy.character_id,
        })
    end
    return list
end

--- Get strategy details
function AutoBattleAI.getStrategy(strategy_id)
    return AutoBattleAI.strategies[strategy_id]
end

--- Delete strategy
function AutoBattleAI.deleteStrategy(strategy_id)
    if not AutoBattleAI.strategies[strategy_id] then
        print("[ERROR] Strategy not found")
        return false
    end
    
    AutoBattleAI.strategies[strategy_id] = nil
    print(string.format("[AI Configurator] Strategy %s deleted", strategy_id))
    return true
end

--- Clone strategy
function AutoBattleAI.cloneStrategy(strategy_id, new_name)
    local original = AutoBattleAI.strategies[strategy_id]
    if not original then
        print("[ERROR] Strategy not found")
        return nil
    end
    
    local clone_id = AutoBattleAI.createStrategy({
        name = new_name or (original.name .. " (Copy)"),
        description = original.description,
        character_id = original.character_id,
        mp_conservation = original.mp_conservation,
        mp_threshold = original.mp_threshold,
        item_usage_enabled = original.item_usage_enabled,
    })
    
    if not clone_id then return nil end
    
    local clone = AutoBattleAI.strategies[clone_id]
    clone.conditions = {}
    for _, condition in ipairs(original.conditions) do
        table.insert(clone.conditions, {
            variable = condition.variable,
            operator = condition.operator,
            value = condition.value,
        })
    end
    
    for priority, actions in pairs(original.actions_by_priority) do
        clone.actions_by_priority[priority] = {}
        for _, action in ipairs(actions) do
            table.insert(clone.actions_by_priority[priority], {
                type = action.type,
                name = action.name,
                target_type = action.target_type,
            })
        end
    end
    
    print(string.format("[AI Configurator] Strategy cloned: %s -> %s", strategy_id, clone_id))
    return clone_id
end

-- ============================================================================
-- EXPORT & IMPORT
-- ============================================================================

--- Export strategy to text format
function AutoBattleAI.exportStrategy(strategy_id)
    local strategy = AutoBattleAI.strategies[strategy_id]
    if not strategy then
        print("[ERROR] Strategy not found")
        return nil
    end
    
    local export_text = string.format([[
==============================================
AI STRATEGY: %s
==============================================
Description: %s
Character ID: %s
MP Conservation: %s (Threshold: %d%%)
Item Usage: %s

--- CONDITIONS ---
]], strategy.name, strategy.description, tostring(strategy.character_id), 
    tostring(strategy.mp_conservation), strategy.mp_threshold,
    tostring(strategy.item_usage_enabled))
    
    for i, condition in ipairs(strategy.conditions) do
        export_text = export_text .. string.format("%d. IF %s %s %s\n",
            i, condition.variable, condition.operator, tostring(condition.value))
    end
    
    export_text = export_text .. "\n--- ACTIONS BY PRIORITY ---\n"
    
    for priority = 1, 5 do
        local priority_name = ({"CRITICAL", "HIGH", "MEDIUM", "LOW", "IDLE"})[priority]
        local actions = strategy.actions_by_priority[priority]
        if #actions > 0 then
            export_text = export_text .. string.format("\n[%s PRIORITY]\n", priority_name)
            for i, action in ipairs(actions) do
                export_text = export_text .. string.format("  %d. %s: %s (Target: %s)\n",
                    i, action.type, action.name, action.target_type)
            end
        end
    end
    
    export_text = export_text .. "\n==============================================\n"
    
    return export_text
end

--- Export all strategies
function AutoBattleAI.exportAllStrategies()
    local export_text = "AUTO-BATTLE AI CONFIGURATOR - ALL STRATEGIES\n"
    export_text = export_text .. string.format("Exported: %s\n\n", os.date("%Y-%m-%d %H:%M:%S"))
    
    for id, strategy in pairs(AutoBattleAI.strategies) do
        local strategy_export = AutoBattleAI.exportStrategy(id)
        if strategy_export then
            export_text = export_text .. strategy_export .. "\n\n"
        end
    end
    
    return export_text
end

-- ============================================================================
-- CHARACTER ANALYSIS
-- ============================================================================

--- Analyze character for optimal strategy
function AutoBattleAI.analyzeCharacter(character_id)
    print(string.format("\n[AI Configurator] Analyzing character %d...", character_id))
    
    -- Read character stats
    local level = safeCall(API.getCharacterLevel, character_id)
    local vigor = safeCall(API.getCharacterVigor, character_id)
    local magic = safeCall(API.getCharacterMagic, character_id)
    local speed = safeCall(API.getCharacterSpeed, character_id)
    local hp = safeCall(API.getCharacterMaxHP, character_id)
    local mp = safeCall(API.getCharacterMaxMP, character_id)
    
    local analysis = {
        character_id = character_id,
        level = level,
        stats = {
            vigor = vigor,
            magic = magic,
            speed = speed,
            hp = hp,
            mp = mp,
        },
        recommended_role = "balanced",
        recommended_strategy = nil,
    }
    
    -- Determine recommended role
    if magic and vigor then
        if magic > vigor + 20 then
            analysis.recommended_role = "mage"
        elseif vigor > magic + 20 then
            analysis.recommended_role = "physical"
        elseif hp and hp > 5000 then
            analysis.recommended_role = "tank"
        elseif mp and mp > 500 then
            analysis.recommended_role = "healer"
        end
    end
    
    print(string.format("Level %d | Vigor: %d | Magic: %d | HP: %d | MP: %d",
        level or 0, vigor or 0, magic or 0, hp or 0, mp or 0))
    print(string.format("Recommended Role: %s", analysis.recommended_role))
    
    return analysis
end

-- ============================================================================
-- PLUGIN INTERFACE
-- ============================================================================

return AutoBattleAI
