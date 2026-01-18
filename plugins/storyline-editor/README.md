# Storyline Editor v1.0.0 - Complete Documentation

**Edit FF6's narrative, dialogue, quests, story branches, and character interactions with precision**

## Overview

The Storyline Editor is an advanced narrative modification system for Final Fantasy VI that allows you to craft custom stories, edit dialogue, modify quest objectives, create alternate story paths, and control story progression with unprecedented flexibility. Perfect for creating fan fiction, story mods, test playthroughs, and narrative experimentation.

## Quick Start

```lua
-- Simple dialogue editing
editDialogue(42, "Welcome, brave adventurers!")
getDialogue(42)

-- Modify quest
editQuestObjective(5, "Find the legendary sword")
editQuestReward(5, "experience", 5000)

-- Skip to story point
skipToStoryPoint("wor_mid")

-- View story status
displayStoryStatus()
```

## Table of Contents

1. [Features](#features)
2. [Story System Overview](#story-system-overview)
3. [Dialogue Editing](#dialogue-editing)
4. [Story Events](#story-events)
5. [Quest Modification](#quest-modification)
6. [Character Interactions](#character-interactions)
7. [Story Branches](#story-branches)
8. [Cutscenes](#cutscenes)
9. [Story Progression](#story-progression)
10. [Story Templates](#story-templates)
11. [Analysis & Export](#analysis--export)
12. [Backup & Restore](#backup--restore)
13. [Advanced Usage](#advanced-usage)
14. [Use Cases](#use-cases)
15. [API Reference](#api-reference)

## Features

### Core Capabilities

✅ **Dialogue Editing** - Modify ~500 NPC and story dialogue entries
✅ **Story Events** - Create and modify 50+ story-triggering events
✅ **Quest System** - Edit 30 quest objectives, descriptions, and rewards
✅ **Character Interactions** - Define relationships and interactions between 14 characters
✅ **Story Branches** - Unlock alternate story paths (6 predefined + custom)
✅ **Cutscene Control** - Toggle cutscene flags and replay sequences
✅ **Story Templates** - Apply preset narrative configurations
✅ **Progression Control** - Skip, rewind, or fast-forward story points
✅ **Batch Operations** - Edit multiple dialogues/quests simultaneously
✅ **Full Backup/Restore** - Create and restore complete story snapshots
✅ **Export/Import** - Share and load story modifications
✅ **Operation Logging** - Track all changes with timestamps

### Advanced Features

- **Custom Story Events** - Create unique story triggers beyond predefined events
- **Event Chaining** - Link events so one triggers another automatically
- **Relationship System** - Set character relationships (friend, rival, romance, neutral)
- **Alternate Paths** - Create completely new story branches
- **Story Analysis** - View detailed statistics about your modifications
- **Template System** - Quick-apply themed story modifications (dark timeline, romance, etc.)

## Story System Overview

### FF6 Story Structure

Final Fantasy VI contains:

- **~500 Dialogue Entries** - NPC greetings, quest text, story conversations, item descriptions
- **50+ Story Events** - Combat events, character recruitment, major plot points, side quests
- **14 Main Characters** - Each with backstory, relationships, and development arcs
- **30 Quests** - Optional and mandatory objectives throughout WoB/WoR
- **6 Story Branches** - WoB Early/Mid/Late and WoR Early/Mid/Late progression points
- **Cutscene System** - Story-driven animated sequences with trigger flags

### Story Progression Phases

The game progresses through these primary story branches:

| Branch | Timeframe | Key Events | Story State |
|--------|-----------|-----------|------------|
| **wob_early** | Game start | Character recruitment, introduction | Getting to know everyone |
| **wob_mid** | Mid World of Balance | Party grows, major challenges | Becoming stronger |
| **wob_late** | Late World of Balance | Empire escalation, preparation | Building for confrontation |
| **wor_early** | World of Ruin start | Cataclysm happens, exploration | Adapting to devastation |
| **wor_mid** | Mid World of Ruin | Reunion and regrouping | Reassembling allies |
| **wor_late** | Final preparations | Gathering power, endgame | Ultimate confrontation |

## Dialogue Editing

### Basic Dialogue Functions

The dialogue system manages ~500 dialogue entries covering:
- NPC dialogue and greetings
- Character-specific conversation lines
- Quest descriptions and objectives
- Story exposition and plot advancement
- Item and magic descriptions

### Edit a Single Dialogue

```lua
-- Edit dialogue by ID
editDialogue(42, "Welcome, brave warriors! The empire grows stronger each day.")

-- Retrieve current dialogue
local current = getDialogue(42)
print("Current dialogue: " .. current)
```

**Parameters:**
- `dialogue_id` (number): 0-499
- `new_text` (string): Replacement text (max 255 characters recommended)

**Returns:** true on success, false on error

**Example:**
```lua
editDialogue(15, "You seem troubled. Can I help?")
editDialogue(16, "There's an ancient legend about a power hidden in the mountains...")
editDialogue(17, "Will you help me find it?")
```

### Batch Edit Dialogues

```lua
-- Edit multiple dialogues at once
local changes = {
    [42] = "Welcome, travelers!",
    [43] = "It's been a long journey for you, hasn't it?",
    [44] = "You must meet the village elder.",
    [45] = "Follow me to his house.",
}

batchEditDialogues(changes)
```

This is efficient for creating complete conversation sequences without multiple individual edits.

### List Dialogues by Type

```lua
-- Get all dialogues of a certain type
local dialogues = listDialoguesByType("npc")

-- Usage in dialogue creation
for _, dialogue in ipairs(dialogues) do
    print(string.format("Dialogue %d: %s", dialogue.id, dialogue.text))
end
```

### Dialogue Editing Use Cases

**Create alternative dialogue:**
```lua
-- Original confrontation scene
editDialogue(100, "You are my enemy!")

-- Change to pacifist option
editDialogue(100, "Perhaps we can talk about this...")
```

**Add character personality:**
```lua
editDialogue(201, "The situation is... dire.")  -- Celes
editDialogue(202, "We've faced worse!")         -- Edgar
editDialogue(203, "Hmph. Bring it on.")         -- Sabin
```

**Write fan fiction dialogue:**
```lua
editDialogue(50, "Locke, I... I've always cared about you.")
editDialogue(51, "Rachel, I thought I'd lost you forever...")
editDialogue(52, "We have a second chance. Let's make it count.")
```

## Story Events

### Story Event System

Story Events are the core narrative triggers that advance FF6's plot:

- **Main Story Events** - Critical plot progression (recruitment, major battles, cutscenes)
- **Side Story Events** - Optional content (side quests, hidden scenes)
- **Character Events** - Personal character development moments
- **World Events** - Environmental and location-specific occurrences

### Modify Story Events

```lua
-- Modify an existing story event
local event_data = {
    name = "Shadow Joins",
    description = "Shadow becomes available as a recruitable character",
    trigger = "wob_mid",
    dialogue_id = 123,
    reward_exp = 1000,
}

modifyStoryEvent(10, event_data)
```

### Create Custom Story Events

```lua
-- Create a new story event
local custom_event = createCustomStoryEvent(
    "Mystery Dungeon Discovery",
    "The party discovers a hidden dungeon",
    "character_level >= 50"
)
```

This creates an event structure you can use in your custom story flow.

### Chain Story Events

```lua
-- Make one event trigger another
chainStoryEvents(5, 6)  -- Event 6 triggers after Event 5

-- Create event sequences
chainStoryEvents(10, 11)
chainStoryEvents(11, 12)
chainStoryEvents(12, 13)

-- Results in: Event 10 → 11 → 12 → 13 (cascade)
```

### Story Event Use Cases

**Create boss progression:**
```lua
-- First boss defeated
modifyStoryEvent(1, {
    name = "Duel with Commander",
    reward_exp = 2000,
})

-- Boss 1 triggers Boss 2
chainStoryEvents(1, 2)

-- Boss 2 progression
modifyStoryEvent(2, {
    name = "Encounter with Emperor",
    reward_exp = 5000,
})
```

**Control recruitment order:**
```lua
modifyStoryEvent(5, {name = "Terra Joins", trigger = "wob_early"})
modifyStoryEvent(6, {name = "Locke Joins", trigger = "wob_early"})
modifyStoryEvent(7, {name = "Edgar Joins", trigger = "wob_mid"})
```

## Quest Modification

### Quest System

FF6 contains ~30 quests covering:
- Main story objectives
- Optional side quests
- Hidden challenges
- Character-specific goals
- Treasure and reward quests

### Modify Quest Data

```lua
-- Complete quest modification
modifyQuest(5, {
    name = "Rescue the Princess",
    description = "The empire has taken the princess. You must save her.",
    difficulty = "medium",
    reward_gold = 10000,
    reward_exp = 5000,
})
```

### Edit Quest Objectives

```lua
-- Change what the player needs to do
editQuestObjective(5, "Infiltrate the imperial base and free the hostages")
editQuestObjective(6, "Defeat the imperial commander")
editQuestObjective(7, "Escape and return to the hideout")
```

### Edit Quest Rewards

```lua
-- Change quest completion rewards
editQuestReward(5, "experience", 10000)  -- XP reward
editQuestReward(5, "gold", 5000)         -- Gold reward
editQuestReward(5, "item", "sword")      -- Item reward
editQuestReward(5, "esper", "shiva")     -- Esper reward
```

### Quest Modification Examples

**Create a side quest:**
```lua
-- Original optional task
editQuestObjective(20, "Optional: Explore the eastern forest")
editQuestReward(20, "experience", 2000)
editQuestReward(20, "item", "potion")

-- Make it mandatory later
chainStoryEvents(8, 20)  -- Story event 8 unlocks this quest
```

**Adjust difficulty progression:**
```lua
-- Early game easier
editQuestObjective(2, "Defeat the commander (Easy)")
editQuestReward(2, "experience", 1000)

-- Late game harder
editQuestObjective(30, "Defeat Kefka's Divine Form (Extreme)")
editQuestReward(30, "experience", 50000)
```

## Character Interactions

### Relationship System

The character system allows you to define how 14 main characters relate to each other:

- **Relationship Types:** friend, rival, romance, neutral
- **Interaction Types:** dialogue, action, reaction, confrontation

### Modify Character Relationships

```lua
-- Define relationships between characters
modifyCharacterRelationship(0, 1, "friend")      -- Terra and Locke
modifyCharacterRelationship(1, 2, "rival")      -- Locke and Edgar
modifyCharacterRelationship(3, 4, "romance")    -- Sabin and ??? (secret)
```

Character IDs:
- 0: Terra
- 1: Locke  
- 2: Edgar
- 3: Sabin
- 4: Shadow
- 5: Cyan
- 6: Gau
- 7: Setzer
- 8: Strago
- 9: Relm
- 10: Mog
- 11: Umaro
- 12: Kefka
- 13: Gestahl

### Add Character Interactions

```lua
-- Add specific interactions between characters
addCharacterInteraction(0, 1, "Terra confides in Locke about her Esper heritage")
addCharacterInteraction(1, 2, "Locke and Edgar discuss the empire's expansion")
addCharacterInteraction(3, 5, "Sabin and Cyan train together in martial arts")
```

### Character Development Arcs

```lua
-- Create a character romance arc
addCharacterInteraction(0, 1, "Act 1: First meeting")
addCharacterInteraction(0, 1, "Act 2: Growing connection")
addCharacterInteraction(0, 1, "Act 3: Deep feelings revealed")
addCharacterInteraction(0, 1, "Act 4: Confession and commitment")
```

### Relationship Configuration Examples

**Create love triangle:**
```lua
modifyCharacterRelationship(0, 1, "romance")    -- Terra loves Locke
modifyCharacterRelationship(0, 2, "romance")    -- Terra loves Edgar
modifyCharacterRelationship(1, 2, "rival")      -- Locke and Edgar rivalry
```

**Build party dynamics:**
```lua
-- Battle-hardened unit
modifyCharacterRelationship(5, 3, "friend")    -- Cyan & Sabin bond
addCharacterInteraction(5, 3, "Cyan and Sabin share stories of martial honor")

-- Rivalry and respect
modifyCharacterRelationship(2, 3, "rival")    -- Edgar vs Sabin
addCharacterInteraction(2, 3, "The brothers compete to prove their strength")
```

## Story Branches

### Branch System Overview

Story Branches represent major progression points in the narrative. FF6 has 6 primary branches:

**World of Balance (WoB):**
- `wob_early` - Introductions and recruitment (hours 0-10)
- `wob_mid` - Growing crisis (hours 10-20)
- `wob_late` - Empire escalation (hours 20-30)

**World of Ruin (WoR):**
- `wor_early` - Post-cataclysm exploration (hours 30-40)
- `wor_mid` - Reunion and regrouping (hours 40-50)
- `wor_late` - Final preparations (hours 50-60)

### Unlock Story Branches

```lua
-- Make a story branch accessible
unlockStoryBranch("wor_early")

-- Unlock all branches at once
unlockStoryBranch("wob_early")
unlockStoryBranch("wob_mid")
unlockStoryBranch("wob_late")
unlockStoryBranch("wor_early")
unlockStoryBranch("wor_mid")
unlockStoryBranch("wor_late")
```

### Set Active Story Branch

```lua
-- Jump to a specific story point
setActiveStoryBranch("wor_mid")

-- This enables:
-- - New dialogue options
-- - Updated NPC locations
-- - Access to new areas
-- - Different quests available
```

### Create Alternate Story Paths

```lua
-- Create a completely new narrative path
createAlternateStoryPath(
    "secret_island_arc",
    "Discover a hidden island with its own questline and story"
)

-- Create multiple branches
createAlternateStoryPath(
    "time_travel_arc",
    "A mysterious portal allows you to revisit past moments"
)

createAlternateStoryPath(
    "esper_origin_story",
    "Learn the true history of the Espers"
)
```

### Story Branch Examples

**Create time-traveling story:**
```lua
-- Normal progression
setActiveStoryBranch("wob_mid")

-- Unlock secret branch
createAlternateStoryPath("time_anomaly", "Time itself is breaking down")

-- Trigger special events
chainStoryEvents(15, 16)  -- Time anomaly detected
chainStoryEvents(16, 17)  -- Travel to past
chainStoryEvents(17, 18)  -- Prevent catastrophe
```

**Create difficulty-based branching:**
```lua
-- Normal difficulty
if difficulty == "normal" then
    setActiveStoryBranch("wob_mid")
elseif difficulty == "hard" then
    -- Skip to harder challenges
    setActiveStoryBranch("wob_late")
    createAlternateStoryPath("ironman", "Permadeath mode - no second chances")
end
```

## Cutscenes

### Cutscene System

Cutscenes are story-driven animated sequences with trigger conditions:

- **Character Cutscenes** - Introduction and development sequences
- **Story Cutscenes** - Major plot progression scenes
- **Action Cutscenes** - Battle-related animated sequences
- **Ending Cutscenes** - Conclusion scenes

### Toggle Cutscene Flags

```lua
-- Enable a cutscene
toggleCutsceneFlag(5, true)

-- Disable/skip a cutscene
toggleCutsceneFlag(5, false)

-- Shorthand: skip specific cutscene
skipCutscene(10)
```

### Replay Cutscenes

```lua
-- Mark cutscene for replay
replayCutscene(5)

-- Allows watching story scenes again
replayCutscene(12)
replayCutscene(25)
```

### Cutscene Control Examples

**Create endless story moments:**
```lua
-- Make key cutscenes replayable
replayCutscene(3)   -- Terra's origin reveal
replayCutscene(8)   -- Character recruitment scenes
replayCutscene(15)  -- Key plot point
replayCutscene(20)  -- Major confrontation
```

**Skip traumatic sequences (player choice):**
```lua
-- Some players may want to skip dark moments
skipCutscene(12)  -- Particularly dark scene

-- But keep important plot scenes
-- (don't skip story-critical cutscene 15)
```

## Story Progression

### Progression Control

Three functions allow precise control over story advancement:

### Skip to Story Point

```lua
-- Jump directly to a story phase
skipToStoryPoint("wor_mid")

-- This will:
-- 1. Unlock all preceding branches
-- 2. Set wor_mid as active
-- 3. Enable all relevant dialogue/NPCs
-- 4. Make appropriate quests available
```

### Rewind Story Progression

```lua
-- Go back in the story
rewindStoryProgression(5)  -- Rewind 5 story steps

-- Use cases:
-- - Undo a major decision
-- - Replay earlier content
-- - Experience alternate dialogue
```

### Fast-Forward Story Progression

```lua
-- Advance the story
fastForwardStoryProgression(10)  -- Skip ahead 10 steps

-- Use cases:
-- - Rush through early game
-- - Focus on endgame content
-- - Speed up pacing
```

### Progression Control Examples

**Create speedrun setup:**
```lua
-- Skip directly to endgame
skipToStoryPoint("wor_late")

-- Player can now:
-- - Access all areas
-- - Recruit all characters
-- - Face final challenges
```

**Create story branch exploration:**
```lua
-- Let player explore each branch
setActiveStoryBranch("wob_early")
-- Player explores...

setActiveStoryBranch("wob_mid")
-- Player explores...

setActiveStoryBranch("wor_early")
-- Player explores with new perspective
```

## Story Templates

### Template System

Story Templates are preset configurations that quickly apply thematic story modifications:

- **alternate_ending** - Changes final story conclusion
- **dark_timeline** - Darker tone, grittier story
- **heroic_journey** - Emphasize hero arcs
- **romance_focused** - Emphasize character relationships
- **political_intrigue** - Focus on political conflict

### Apply Story Template

```lua
-- Apply a preset story configuration
applyStoryTemplate("alternate_ending")

-- This will modify:
-- - Final dialogue
-- - Ending events
-- - Resolution text
-- - Character fates
```

### List Available Templates

```lua
-- See all templates
listStoryTemplates()

-- Returns:
-- alternate_ending: Alternate story ending
-- dark_timeline: Darker story tone and events
-- heroic_journey: Emphasize hero character arcs
-- romance_focused: Emphasize romantic relationships
-- political_intrigue: Emphasize political conflict
```

### Template Examples

**Create fan fiction themes:**
```lua
-- Emphasize romance storyline
applyStoryTemplate("romance_focused")

-- Make story darker
applyStoryTemplate("dark_timeline")

-- Focus on heroic journeys
applyStoryTemplate("heroic_journey")
```

**Combine templates:**
```lua
-- First apply romance focus
applyStoryTemplate("romance_focused")

-- Then make darker
applyStoryTemplate("dark_timeline")

-- Result: Dark romance storyline
```

## Analysis & Export

### Get Story Status

```lua
-- View current storyline state
local status = getStoryStatus()

-- Returns: {
--   dialogues_edited = 15,
--   events_modified = 8,
--   quests_modified = 12,
--   branches_unlocked = 4,
--   current_branch = "wor_mid"
-- }
```

### Display Story Status

```lua
-- Print formatted story status
displayStoryStatus()

-- Output:
-- === Story Status ===
-- Dialogues Edited: 15
-- Story Events Modified: 8
-- Quests Modified: 12
-- Branches Unlocked: 4 / 6
-- Current Branch: wor_mid
```

### Export Storyline Configuration

```lua
-- Export all modifications
exportStorylineConfig()

-- Creates summary including:
-- - Export timestamp
-- - All edited dialogues
-- - Modified events
-- - Changed quests
-- - Unlocked branches
-- - Complete operation log
```

## Backup & Restore

### Automatic Backups

The plugin automatically creates backups before major modifications:

```lua
-- Backups are created before:
editDialogue(42, "New text")           -- Auto-backup
modifyStoryEvent(5, event_data)        -- Auto-backup
editQuestObjective(10, "New objective") -- Auto-backup
```

### Manual Backup

```lua
-- Create a full snapshot
create_backup()

-- Stores:
-- - All dialogue texts
-- - Story event configurations
-- - Quest modifications
-- - Story branch states
```

### Restore from Backup

```lua
-- Restore entire story to last backup point
restoreBackup()

-- This will:
-- - Undo all modifications since backup
-- - Restore dialogue texts
-- - Reset story events
-- - Revert quest changes
-- - Clear branch modifications
```

## Advanced Usage

### Creating Complex Story Arcs

```lua
-- 5-act story structure
local acts = {
    exposition = "wob_early",
    rising_action = "wob_mid",
    climax = "wob_late",
    falling_action = "wor_mid",
    resolution = "wor_late"
}

for phase, branch in pairs(acts) do
    setActiveStoryBranch(branch)
    -- Add dialogue and events for this phase
end
```

### Building Character Development

```lua
-- Character arc progression
local character_arc = {
    act1 = {objective = "Escape", reward = 500},
    act2 = {objective = "Learn truth", reward = 1000},
    act3 = {objective = "Gain power", reward = 2000},
    act4 = {objective = "Confront past", reward = 5000},
    act5 = {objective = "Make choice", reward = 10000},
}

for act, data in pairs(character_arc) do
    local quest_id = tonumber(act:match("%d+"))
    if quest_id then
        editQuestObjective(quest_id, data.objective)
        editQuestReward(quest_id, "experience", data.reward)
    end
end
```

### Conditional Story Paths

```lua
-- Create branching narrative
function advanceStory(player_choice)
    if player_choice == "love" then
        applyStoryTemplate("romance_focused")
        setActiveStoryBranch("wor_mid")
    elseif player_choice == "power" then
        applyStoryTemplate("dark_timeline")
        skipToStoryPoint("wor_late")
    elseif player_choice == "peace" then
        applyStoryTemplate("heroic_journey")
        setActiveStoryBranch("wor_mid")
    end
end
```

## Use Cases

### 1. Fan Fiction Creation

Create your own FF6 storyline with custom dialogue and character interactions:

```lua
-- Write your story
editDialogue(100, "But what if the empire never rose to power?")
editDialogue(101, "Our world would be vastly different...")

-- Create alternate history
createAlternateStoryPath("divergence", "What if everything changed?")

-- Build character relationships
modifyCharacterRelationship(0, 1, "romance")
addCharacterInteraction(0, 1, "In this timeline, destiny found them together")
```

### 2. Story Testing & QA

Test different story branches and dialogue variations:

```lua
-- Quickly jump between story points
skipToStoryPoint("wob_early")
-- Test WoB early dialogue

skipToStoryPoint("wor_late")
-- Test WoR late endgame content

skipToStoryPoint("wob_late")
-- Test late WoB before cataclysm
```

### 3. Narrative Experimentation

Experiment with story themes and tone:

```lua
-- Try darker tone
applyStoryTemplate("dark_timeline")

-- Compare with heroic tone
applyStoryTemplate("heroic_journey")

-- Export both versions
exportStorylineConfig()
```

### 4. Educational Storytelling

Use for writing classes and story structure analysis:

```lua
-- Display 3-act structure
setActiveStoryBranch("wob_mid")   -- Act 1
setActiveStoryBranch("wob_late")  -- Act 2
setActiveStoryBranch("wor_late")  -- Act 3

-- Modify for specific lessons
editDialogue(10, "Each character had their own story...")
editDialogue(11, "Which wove together into the greater narrative...")
```

### 5. Speedrun/Challenge Runs

Modify story pacing for different playstyle:

```lua
-- Speedrun configuration
skipToStoryPoint("wor_late")
-- Allows quick access to all content

-- Challenge run
applyStoryTemplate("dark_timeline")
-- Darker, grittier story for atmosphere
```

### 6. Localization Testing

Test dialogue translations and space constraints:

```lua
-- Original English
editDialogue(50, "Welcome, brave warriors!")

-- Spanish test
editDialogue(50, "¡Bienvenidos, valientes guerreros!")

-- German test
editDialogue(50, "Willkommen, tapfere Krieger!")

-- Check rendering and text box fit
displayStoryStatus()
```

## API Reference

### Core Functions

#### Dialogue Editing
```lua
editDialogue(dialogue_id: number, new_text: string) -> boolean
getDialogue(dialogue_id: number) -> string
batchEditDialogues(dialogue_map: table) -> number
listDialoguesByType(dialogue_type: string) -> table
```

#### Story Events
```lua
modifyStoryEvent(event_id: number, event_data: table) -> boolean
createCustomStoryEvent(name: string, description: string, trigger: string) -> table
chainStoryEvents(event_id_1: number, event_id_2: number) -> boolean
```

#### Quest Modification
```lua
modifyQuest(quest_id: number, quest_data: table) -> boolean
editQuestObjective(quest_id: number, new_objective: string) -> boolean
editQuestReward(quest_id: number, reward_type: string, reward_value) -> boolean
```

#### Character Interactions
```lua
modifyCharacterRelationship(char_id_1: number, char_id_2: number, relationship_type: string) -> boolean
addCharacterInteraction(char_id_1: number, char_id_2: number, interaction_text: string) -> boolean
```

#### Story Branches
```lua
unlockStoryBranch(branch_name: string) -> boolean
setActiveStoryBranch(branch_name: string) -> boolean
createAlternateStoryPath(path_name: string, description: string) -> table
```

#### Cutscenes
```lua
toggleCutsceneFlag(cutscene_id: number, enable: boolean) -> boolean
skipCutscene(cutscene_id: number) -> boolean
replayCutscene(cutscene_id: number) -> boolean
```

#### Story Progression
```lua
skipToStoryPoint(target_branch: string) -> boolean
rewindStoryProgression(steps: number) -> boolean
fastForwardStoryProgression(steps: number) -> boolean
```

#### Templates
```lua
applyStoryTemplate(template_name: string) -> boolean
listStoryTemplates() -> table
```

#### Analysis & Export
```lua
getStoryStatus() -> table
displayStoryStatus() -> table
exportStorylineConfig() -> string
```

#### Backup & Restore
```lua
restoreBackup() -> boolean
```

## Performance Considerations

- **Dialogue Cache:** ~500 entries, minimal memory impact
- **Story Events:** ~50 events, negligible performance cost
- **Backups:** Create before major operations, stored in memory
- **Export:** Fast operation, suitable for frequent use

## Tips & Tricks

1. **Create story outline first** - Plan your story before editing dialogue
2. **Use batch editing** - Faster than individual edits for multiple changes
3. **Test with templates** - Templates provide quick story tone experiments
4. **Regular backups** - The system auto-backs up, but you can add manual checkpoints
5. **Export variations** - Try different paths and export each version
6. **Character relationships** - Define relationships before writing interactions
7. **Event chaining** - Use chains to create natural story flow
8. **Story templates** - Layer multiple templates for unique themes

## Examples & Recipes

### Recipe 1: Time Travel Story

```lua
createAlternateStoryPath("time_anomaly", "A rift in time has opened")
unlockStoryBranch("wob_early")
setActiveStoryBranch("wob_early")

editDialogue(100, "We've been pulled into the past!")
editDialogue(101, "We must be careful not to change history...")

replayCutscene(5)  -- Replay scene from past

skipToStoryPoint("wob_late")  -- Return to present

editDialogue(102, "The timeline has stabilized... for now.")
```

### Recipe 2: Character-Focused Story

```lua
-- Pick a character
local protagonist = 0  -- Terra

-- Build relationship web
modifyCharacterRelationship(protagonist, 1, "friend")
modifyCharacterRelationship(protagonist, 2, "rival")
modifyCharacterRelationship(protagonist, 3, "romance")

-- Add interactions
addCharacterInteraction(protagonist, 1, "I trust you above all others")
addCharacterInteraction(protagonist, 2, "We're more alike than you know")
addCharacterInteraction(protagonist, 3, "In another world, perhaps we could...")

-- Create character arc
for quest_id = 10, 15 do
    editQuestObjective(quest_id, string.format("Terra's Quest Part %d", quest_id - 9))
end

displayStoryStatus()
```

### Recipe 3: Difficulty-Adjusted Story

```lua
local difficulty = "hard"

if difficulty == "easy" then
    applyStoryTemplate("heroic_journey")
    editQuestReward(1, "experience", 1000)
elseif difficulty == "hard" then
    applyStoryTemplate("dark_timeline")
    editQuestReward(1, "experience", 5000)
    createAlternateStoryPath("nightmare_mode", "The stakes have never been higher")
end

displayStoryStatus()
```

## Troubleshooting

**Issue:** Dialogue not appearing in-game
- Check dialogue ID is valid (0-499)
- Verify text encoding is correct
- Test with simpler dialogue first

**Issue:** Story branch not unlocking
- Ensure branch name matches configuration (wob_early, wor_mid, etc.)
- Check previous branches are unlocked first
- Use skipToStoryPoint() instead for guaranteed unlock

**Issue:** Can't restore backup
- Verify backup was created (check last operation)
- Ensure restoreBackup() called while plugin active
- Create new backup and try different story state

## Support & Extension

This plugin is designed to be extensible. Advanced users can:
- Write custom story event triggers
- Create new story templates
- Build story generation systems
- Integrate with other narrative tools

The plugin provides safe API wrappers (pcall) for error handling, ensuring your FF6 save remains protected even with complex modifications.

## Version History

**v1.0.0 (January 16, 2026)** - Initial Release
- Complete dialogue editing system (~500 entries)
- Story event modification and chaining
- Quest objective and reward editing
- Character relationship system (14 characters)
- Story branch unlocking and control (6 branches)
- Cutscene flag management
- Story progression control
- Template system (5 preset themes)
- Full backup/restore functionality
- Export/analysis tools
- Operation logging with timestamps

## Author Notes

The Storyline Editor represents the culmination of advanced FF6 plugin design. It provides narrative designers, fan fiction writers, and story enthusiasts with professional-grade tools for crafting custom Final Fantasy VI stories. Whether you're creating alternate timelines, exploring character relationships, or building complete narrative mods, this plugin provides the flexibility and power to bring your story to life.

The system respects the FF6 lore while allowing unlimited creative freedom. All modifications are reversible through the backup system, encouraging experimentation without fear of permanence.

---

**Storyline Editor v1.0.0** - *Edit Destiny. Create History.*
