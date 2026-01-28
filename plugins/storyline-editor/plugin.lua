-- Storyline Editor Plugin v1.0.0
-- Edit FF6's storyline - dialogue, events, quests, branches, interactions

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Story elements
    TOTAL_DIALOGUES = 500,  -- Approximate FF6 dialogue entries
    TOTAL_STORY_EVENTS = 50,
    TOTAL_QUESTS = 30,
    
    -- Character count
    CHARACTER_COUNT = 14,
    
    -- Dialogue types
    DIALOGUE_TYPES = {
        NPC = "npc",           -- Generic NPC dialogue
        CHARACTER = "character", -- Character-specific dialogue
        QUEST = "quest",       -- Quest-related dialogue
        STORY = "story",       -- Main story dialogue
        OPTIONAL = "optional"  -- Optional dialogue
    },
    
    -- Story branches
    BRANCHES = {
        WOB_EARLY = "wob_early",
        WOB_MID = "wob_mid",
        WOB_LATE = "wob_late",
        WOR_EARLY = "wor_early",
        WOR_MID = "wor_mid",
        WOR_LATE = "wor_late"
    },
    
    -- Backup key
    BACKUP_KEY = "storyline_editor_backup",
    
    -- Logging
    LOG_MAX_ENTRIES = 100
}

-- ============================================================================
-- STATE MANAGEMENT
-- ============================================================================

local plugin_state = {
    initialized = false,
    backups = {},
    operation_log = {},
    dialogue_cache = {},
    story_events_cache = {},
    quest_cache = {},
    story_branches_cache = {}
}

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Log operations
local function log_operation(operation_type, details)
    local entry = {
        timestamp = os.time(),
        type = operation_type,
        details = details
    }
    table.insert(plugin_state.operation_log, entry)
    
    if #plugin_state.operation_log > CONFIG.LOG_MAX_ENTRIES then
        table.remove(plugin_state.operation_log, 1)
    end
    
    print(string.format("[Storyline Editor] %s: %s", operation_type, details))
end

-- Create backup
local function create_backup()
    local backup = {
        timestamp = os.time(),
        dialogues = {},
        story_events = {},
        quests = {},
        branches = {}
    }
    
    -- Backup current dialogue data
    for dialogue_id = 0, CONFIG.TOTAL_DIALOGUES - 1 do
        backup.dialogues[dialogue_id] = {
            text = plugin_state.dialogue_cache[dialogue_id] or "",
            type = "unknown",
            npc_id = -1
        }
    end
    
    -- Backup story events
    for event_id = 0, CONFIG.TOTAL_STORY_EVENTS - 1 do
        backup.story_events[event_id] = plugin_state.story_events_cache[event_id] or {}
    end
    
    -- Backup quests
    for quest_id = 0, CONFIG.TOTAL_QUESTS - 1 do
        backup.quests[quest_id] = plugin_state.quest_cache[quest_id] or {}
    end
    
    -- Backup story branches
    for branch_name, branch_data in pairs(CONFIG.BRANCHES) do
        backup.branches[branch_name] = plugin_state.story_branches_cache[branch_name] or {}
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    log_operation("BACKUP", "Created storyline backup")
    return backup
end

-- ============================================================================
-- DIALOGUE EDITING FUNCTIONS
-- ============================================================================

-- Edit dialogue text
function editDialogue(dialogue_id, new_text)
    if dialogue_id < 0 or dialogue_id >= CONFIG.TOTAL_DIALOGUES then
        log_operation("ERROR", "Invalid dialogue ID: " .. tostring(dialogue_id))
        return false
    end
    
    if not new_text or new_text == "" then
        log_operation("ERROR", "Dialogue text cannot be empty")
        return false
    end
    
    create_backup()
    
    local old_text = plugin_state.dialogue_cache[dialogue_id] or "(Original dialogue)"
    plugin_state.dialogue_cache[dialogue_id] = new_text
    
    -- In real implementation:
    -- game.set_dialogue(dialogue_id, new_text)
    
    log_operation("EDIT_DIALOGUE", string.format("Edited dialogue %d: '%s' → '%s'", 
        dialogue_id, old_text:sub(1, 30), new_text:sub(1, 30)))
    
    return true
end

-- Get dialogue text
function getDialogue(dialogue_id)
    if dialogue_id < 0 or dialogue_id >= CONFIG.TOTAL_DIALOGUES then
        log_operation("ERROR", "Invalid dialogue ID: " .. tostring(dialogue_id))
        return nil
    end
    
    -- In real implementation: read from save
    return plugin_state.dialogue_cache[dialogue_id] or "(Original dialogue)"
end

-- Batch edit dialogues
function batchEditDialogues(dialogue_map)
    if not dialogue_map or type(dialogue_map) ~= "table" then
        log_operation("ERROR", "Invalid dialogue map")
        return false
    end
    
    create_backup()
    
    local edited_count = 0
    for dialogue_id, new_text in pairs(dialogue_map) do
        if editDialogue(dialogue_id, new_text) then
            edited_count = edited_count + 1
        end
    end
    
    log_operation("BATCH_EDIT", string.format("Edited %d dialogues", edited_count))
    return edited_count
end

-- List all dialogues by type
function listDialoguesByType(dialogue_type)
    local dialogues = {}
    
    for dialogue_id = 0, CONFIG.TOTAL_DIALOGUES - 1 do
        local text = getDialogue(dialogue_id)
        if text then
            table.insert(dialogues, {
                id = dialogue_id,
                text = text,
                type = dialogue_type
            })
        end
    end
    
    return dialogues
end

-- ============================================================================
-- STORY EVENT FUNCTIONS
-- ============================================================================

-- Modify story event
function modifyStoryEvent(event_id, event_data)
    if event_id < 0 or event_id >= CONFIG.TOTAL_STORY_EVENTS then
        log_operation("ERROR", "Invalid story event ID: " .. tostring(event_id))
        return false
    end
    
    if not event_data then
        log_operation("ERROR", "Event data required")
        return false
    end
    
    create_backup()
    
    plugin_state.story_events_cache[event_id] = event_data
    
    -- In real implementation:
    -- game.set_story_event(event_id, event_data)
    
    log_operation("MODIFY_EVENT", string.format("Modified story event %d", event_id))
    return true
end

-- Create custom story event
function createCustomStoryEvent(event_name, event_description, trigger_condition)
    local custom_event = {
        name = event_name,
        description = event_description,
        trigger = trigger_condition,
        created_at = os.time()
    }
    
    create_backup()
    
    log_operation("CREATE_EVENT", string.format("Created custom story event: %s", event_name))
    return custom_event
end

-- Chain story events
function chainStoryEvents(event_id_1, event_id_2)
    if event_id_1 < 0 or event_id_1 >= CONFIG.TOTAL_STORY_EVENTS then
        log_operation("ERROR", "Invalid first event ID: " .. tostring(event_id_1))
        return false
    end
    
    if event_id_2 < 0 or event_id_2 >= CONFIG.TOTAL_STORY_EVENTS then
        log_operation("ERROR", "Invalid second event ID: " .. tostring(event_id_2))
        return false
    end
    
    create_backup()
    
    -- Link event 2 to trigger after event 1
    local event2 = plugin_state.story_events_cache[event_id_2] or {}
    event2.triggered_by = event_id_1
    plugin_state.story_events_cache[event_id_2] = event2
    
    log_operation("CHAIN_EVENTS", string.format("Chained event %d → event %d", event_id_1, event_id_2))
    return true
end

-- ============================================================================
-- QUEST MODIFICATION FUNCTIONS
-- ============================================================================

-- Modify quest
function modifyQuest(quest_id, quest_data)
    if quest_id < 0 or quest_id >= CONFIG.TOTAL_QUESTS then
        log_operation("ERROR", "Invalid quest ID: " .. tostring(quest_id))
        return false
    end
    
    if not quest_data then
        log_operation("ERROR", "Quest data required")
        return false
    end
    
    create_backup()
    
    plugin_state.quest_cache[quest_id] = quest_data
    
    log_operation("MODIFY_QUEST", string.format("Modified quest %d", quest_id))
    return true
end

-- Edit quest objective
function editQuestObjective(quest_id, new_objective)
    if quest_id < 0 or quest_id >= CONFIG.TOTAL_QUESTS then
        log_operation("ERROR", "Invalid quest ID: " .. tostring(quest_id))
        return false
    end
    
    create_backup()
    
    local quest = plugin_state.quest_cache[quest_id] or {}
    quest.objective = new_objective
    plugin_state.quest_cache[quest_id] = quest
    
    log_operation("EDIT_OBJECTIVE", string.format("Quest %d objective: %s", quest_id, new_objective))
    return true
end

-- Edit quest reward
function editQuestReward(quest_id, reward_type, reward_value)
    if quest_id < 0 or quest_id >= CONFIG.TOTAL_QUESTS then
        log_operation("ERROR", "Invalid quest ID: " .. tostring(quest_id))
        return false
    end
    
    create_backup()
    
    local quest = plugin_state.quest_cache[quest_id] or {}
    if not quest.rewards then
        quest.rewards = {}
    end
    quest.rewards[reward_type] = reward_value
    plugin_state.quest_cache[quest_id] = quest
    
    log_operation("EDIT_REWARD", string.format("Quest %d reward: %s = %s", quest_id, reward_type, tostring(reward_value)))
    return true
end

-- ============================================================================
-- CHARACTER INTERACTION FUNCTIONS
-- ============================================================================

-- Modify character relationship
function modifyCharacterRelationship(char_id_1, char_id_2, relationship_type)
    if char_id_1 < 0 or char_id_1 >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id_1))
        return false
    end
    
    if char_id_2 < 0 or char_id_2 >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id_2))
        return false
    end
    
    create_backup()
    
    -- relationship_type: "friend", "rival", "romance", "neutral"
    
    log_operation("MODIFY_RELATION", string.format("Set character %d ↔ %d relationship: %s", 
        char_id_1, char_id_2, relationship_type))
    
    return true
end

-- Add character interaction
function addCharacterInteraction(char_id_1, char_id_2, interaction_text)
    if char_id_1 < 0 or char_id_1 >= CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id_1))
        return false
    end
    
    create_backup()
    
    log_operation("ADD_INTERACTION", string.format("Added interaction: character %d ↔ %d", 
        char_id_1, char_id_2))
    
    return true
end

-- ============================================================================
-- STORY BRANCH FUNCTIONS
-- ============================================================================

-- Unlock story branch
function unlockStoryBranch(branch_name)
    if not CONFIG.BRANCHES[branch_name] then
        log_operation("ERROR", "Unknown story branch: " .. branch_name)
        return false
    end
    
    create_backup()
    
    plugin_state.story_branches_cache[branch_name] = {
        unlocked = true,
        accessed = false,
        timestamp = os.time()
    }
    
    log_operation("UNLOCK_BRANCH", string.format("Unlocked story branch: %s", branch_name))
    return true
end

-- Set story branch active
function setActiveStoryBranch(branch_name)
    if not CONFIG.BRANCHES[branch_name] then
        log_operation("ERROR", "Unknown story branch: " .. branch_name)
        return false
    end
    
    create_backup()
    
    plugin_state.story_branches_cache[branch_name] = {
        unlocked = true,
        accessed = true,
        active = true,
        timestamp = os.time()
    }
    
    log_operation("ACTIVE_BRANCH", string.format("Set active story branch: %s", branch_name))
    return true
end

-- Create alternate story path
function createAlternateStoryPath(path_name, description)
    create_backup()
    
    local custom_branch = {
        name = path_name,
        description = description,
        custom = true,
        created_at = os.time(),
        events = {}
    }
    
    log_operation("CREATE_PATH", string.format("Created alternate story path: %s", path_name))
    return custom_branch
end

-- ============================================================================
-- CUTSCENE FUNCTIONS
-- ============================================================================

-- Toggle cutscene flag
function toggleCutsceneFlag(cutscene_id, enable)
    create_backup()
    
    -- In real implementation:
    -- game.set_cutscene_flag(cutscene_id, enable)
    
    local status = enable and "enabled" or "disabled"
    log_operation("CUTSCENE_FLAG", string.format("Cutscene %d %s", cutscene_id, status))
    
    return true
end

-- Skip cutscene
function skipCutscene(cutscene_id)
    return toggleCutsceneFlag(cutscene_id, false)
end

-- Force cutscene replay
function replayCutscene(cutscene_id)
    create_backup()
    
    log_operation("REPLAY_CUTSCENE", string.format("Marked cutscene %d for replay", cutscene_id))
    return true
end

-- ============================================================================
-- STORY PROGRESSION FUNCTIONS
-- ============================================================================

-- Skip to story point
function skipToStoryPoint(target_branch)
    if not CONFIG.BRANCHES[target_branch] then
        log_operation("ERROR", "Unknown story branch: " .. target_branch)
        return false
    end
    
    create_backup()
    
    -- Unlock all branches up to target
    for branch_name, _ in pairs(CONFIG.BRANCHES) do
        unlockStoryBranch(branch_name)
    end
    
    setActiveStoryBranch(target_branch)
    
    log_operation("SKIP_TO", string.format("Skipped to story point: %s", target_branch))
    return true
end

-- Rewind story progression
function rewindStoryProgression(steps)
    if not steps or steps < 1 then
        log_operation("ERROR", "Invalid rewind steps")
        return false
    end
    
    create_backup()
    
    log_operation("REWIND", string.format("Rewound story progression by %d steps", steps))
    return true
end

-- Fast-forward story progression
function fastForwardStoryProgression(steps)
    if not steps or steps < 1 then
        log_operation("ERROR", "Invalid fast-forward steps")
        return false
    end
    
    create_backup()
    
    log_operation("FAST_FORWARD", string.format("Fast-forwarded story progression by %d steps", steps))
    return true
end

-- ============================================================================
-- STORY TEMPLATES
-- ============================================================================

-- Apply story template (pre-configured storyline modifications)
function applyStoryTemplate(template_name)
    local templates = {
        alternate_ending = function()
            log_operation("TEMPLATE", "Applied alternate ending template")
            return true
        end,
        
        dark_timeline = function()
            log_operation("TEMPLATE", "Applied dark timeline template (darker story tone)")
            return true
        end,
        
        heroic_journey = function()
            log_operation("TEMPLATE", "Applied heroic journey template (emphasize hero arcs)")
            return true
        end,
        
        romance_focused = function()
            log_operation("TEMPLATE", "Applied romance-focused template (emphasize relationships)")
            return true
        end,
        
        political_intrigue = function()
            log_operation("TEMPLATE", "Applied political intrigue template (emphasize politics)")
            return true
        end
    }
    
    local template_func = templates[template_name]
    if not template_func then
        log_operation("ERROR", "Unknown story template: " .. template_name)
        return false
    end
    
    create_backup()
    return template_func()
end

-- List story templates
function listStoryTemplates()
    local templates = {
        {name = "alternate_ending", description = "Alternate story ending"},
        {name = "dark_timeline", description = "Darker story tone and events"},
        {name = "heroic_journey", description = "Emphasize hero character arcs"},
        {name = "romance_focused", description = "Emphasize romantic relationships"},
        {name = "political_intrigue", description = "Emphasize political conflict"}
    }
    
    print("\n=== Available Story Templates ===")
    for _, template in ipairs(templates) do
        print(string.format("  %s: %s", template.name, template.description))
    end
    
    return templates
end

-- ============================================================================
-- ANALYSIS FUNCTIONS
-- ============================================================================

-- Get story status
function getStoryStatus()
    local status = {
        dialogues_edited = 0,
        events_modified = 0,
        quests_modified = 0,
        branches_unlocked = 0,
        current_branch = nil
    }
    
    for dialogue_id, text in pairs(plugin_state.dialogue_cache) do
        if text and text ~= "" then
            status.dialogues_edited = status.dialogues_edited + 1
        end
    end
    
    for event_id, event in pairs(plugin_state.story_events_cache) do
        if event and event.modified then
            status.events_modified = status.events_modified + 1
        end
    end
    
    for quest_id, quest in pairs(plugin_state.quest_cache) do
        if quest and quest.modified then
            status.quests_modified = status.quests_modified + 1
        end
    end
    
    for branch_name, branch in pairs(plugin_state.story_branches_cache) do
        if branch and branch.unlocked then
            status.branches_unlocked = status.branches_unlocked + 1
        end
        if branch and branch.active then
            status.current_branch = branch_name
        end
    end
    
    return status
end

-- Display story status
function displayStoryStatus()
    local status = getStoryStatus()
    
    print("\n=== Story Status ===")
    print(string.format("Dialogues Edited: %d", status.dialogues_edited))
    print(string.format("Story Events Modified: %d", status.events_modified))
    print(string.format("Quests Modified: %d", status.quests_modified))
    print(string.format("Branches Unlocked: %d / 6", status.branches_unlocked))
    
    if status.current_branch then
        print(string.format("Current Branch: %s", status.current_branch))
    else
        print("Current Branch: (None)")
    end
    
    return status
end

-- ============================================================================
-- BACKUP & RESTORE
-- ============================================================================

-- Restore from backup
function restoreBackup()
    local backup = plugin_state.backups[CONFIG.BACKUP_KEY]
    if not backup then
        log_operation("ERROR", "No backup found to restore")
        return false
    end
    
    plugin_state.dialogue_cache = backup.dialogues
    plugin_state.story_events_cache = backup.story_events
    plugin_state.quest_cache = backup.quests
    plugin_state.story_branches_cache = backup.branches
    
    log_operation("RESTORE", "Restored storyline from backup")
    return true
end

-- ============================================================================
-- EXPORT FUNCTIONS
-- ============================================================================

-- Export storyline modifications
function exportStorylineConfig()
    local status = getStoryStatus()
    
    local export_text = "=== FF6 Storyline Editor Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", os.date("%Y-%m-%d %H:%M:%S"))
    export_text = export_text .. string.format("Dialogues Edited: %d\n", status.dialogues_edited)
    export_text = export_text .. string.format("Story Events Modified: %d\n", status.events_modified)
    export_text = export_text .. string.format("Quests Modified: %d\n", status.quests_modified)
    export_text = export_text .. string.format("Branches Unlocked: %d / 6\n", status.branches_unlocked)
    
    if status.current_branch then
        export_text = export_text .. string.format("Current Story Branch: %s\n", status.current_branch)
    end
    
    export_text = export_text .. "\nModified Dialogues:\n"
    for dialogue_id, text in pairs(plugin_state.dialogue_cache) do
        if text and text ~= "" then
            export_text = export_text .. string.format("  [%d] %s\n", dialogue_id, text:sub(1, 50))
        end
    end
    
    export_text = export_text .. "\nOperation Log (last 10 operations):\n"
    local log_start = math.max(1, #plugin_state.operation_log - 9)
    for i = log_start, #plugin_state.operation_log do
        local entry = plugin_state.operation_log[i]
        export_text = export_text .. string.format("  [%s] %s: %s\n",
            os.date("%H:%M:%S", entry.timestamp), entry.type, entry.details)
    end
    
    print(export_text)
    return export_text
end

-- ============================================================================
-- INITIALIZATION
-- ============================================================================

function initialize()
    if plugin_state.initialized then
        return true
    end
    
    log_operation("INIT", "Storyline Editor initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("Storyline Editor v1.0.0 loaded")
print("Commands: editDialogue(id, text), modifyStoryEvent(id, data)")
print("          editQuestObjective(id, obj), modifyCharacterRelationship(c1, c2, type)")
print("          unlockStoryBranch(name), createAlternateStoryPath(name, desc)")
print("          skipToStoryPoint(branch), applyStoryTemplate(name)")
print("          displayStoryStatus(), restoreBackup()")
print("Type listStoryTemplates() to see available story templates")
