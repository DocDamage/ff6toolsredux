-- ============================================================================
-- Character Stats Display Plugin
-- ============================================================================
-- @id: stats-display
-- @name: Character Stats Display
-- @version: 1.0.0
-- @author: FF6 Editor Team
-- @description: Display comprehensive character stats including HP, MP, equipment, and spells
-- @permissions: read_save, ui_display

-- ============================================================================
-- Configuration
-- ============================================================================
local config = {
    showEquipment = true,      -- Show equipment section
    showSpells = true,          -- Show learned spells
    showCommands = true,        -- Show command assignments
    maxSpellsDisplay = 20,      -- Maximum spells to display
    sortSpellsByProficiency = true  -- Sort spells by proficiency level
}

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Format HP/MP with percentage
local function formatHPMP(current, max)
    if max == 0 then return "0/0 (0%)" end
    local percent = (current / max) * 100
    return string.format("%d/%d (%.1f%%)", current, max, percent)
end

-- Get equipment name from ID
local function getEquipmentName(equipID)
    -- Equipment ID mappings (simplified - would normally use game data)
    local emptyEquipment = {
        [93] = "Empty",   -- Weapon/Shield empty
        [198] = "Empty",  -- Helmet empty
        [199] = "Empty",  -- Armor empty
        [200] = "Empty"   -- Relic empty
    }
    
    if emptyEquipment[equipID] then
        return emptyEquipment[equipID]
    end
    
    return string.format("Item #%d", equipID)
end

-- Format spell list
local function formatSpells(spells, config)
    if not spells or #spells == 0 then
        return "No spells learned"
    end
    
    -- Convert spells map to sortable array
    local spellArray = {}
    for spellID, proficiency in pairs(spells) do
        table.insert(spellArray, {id = spellID, proficiency = proficiency})
    end
    
    -- Sort by proficiency if configured
    if config.sortSpellsByProficiency then
        table.sort(spellArray, function(a, b)
            return a.proficiency > b.proficiency
        end)
    else
        table.sort(spellArray, function(a, b)
            return a.id < b.id
        end)
    end
    
    -- Build spell list string
    local lines = {}
    local count = 0
    for _, spell in ipairs(spellArray) do
        if count >= config.maxSpellsDisplay then
            table.insert(lines, string.format("... and %d more", #spellArray - count))
            break
        end
        table.insert(lines, string.format("  Spell #%d (Proficiency: %d%%)", spell.id, spell.proficiency))
        count = count + 1
    end
    
    return table.concat(lines, "\n")
end

-- Format command list
local function formatCommands(commands)
    if not commands or #commands == 0 then
        return "No commands assigned"
    end
    
    local lines = {}
    for i, cmd in ipairs(commands) do
        table.insert(lines, string.format("  %d. %s", i, cmd.Name or "Unknown"))
    end
    
    return table.concat(lines, "\n")
end

-- Build character info display
local function buildCharacterDisplay(char, config)
    local sections = {}
    
    -- Header
    table.insert(sections, "=".rep(60))
    table.insert(sections, string.format("CHARACTER: %s", char.Name or "Unknown"))
    table.insert(sections, "=".rep(60))
    
    -- Basic Stats
    table.insert(sections, "\n[BASIC STATS]")
    table.insert(sections, string.format("  Status: %s", char.Enabled and "Active" or "Inactive"))
    table.insert(sections, string.format("  Level: %d", char.Level or 0))
    table.insert(sections, string.format("  Experience: %d", char.Experience or 0))
    
    -- HP/MP
    table.insert(sections, "\n[HP / MP]")
    table.insert(sections, string.format("  HP: %s", formatHPMP(char.CurrentHP or 0, char.MaxHP or 1)))
    table.insert(sections, string.format("  MP: %s", formatHPMP(char.CurrentMP or 0, char.MaxMP or 1)))
    
    -- Combat Stats
    table.insert(sections, "\n[COMBAT STATS]")
    table.insert(sections, string.format("  Vigor (Physical Attack): %d", char.Vigor or 0))
    table.insert(sections, string.format("  Stamina (Physical Defense): %d", char.Stamina or 0))
    table.insert(sections, string.format("  Speed (Agility): %d", char.Speed or 0))
    table.insert(sections, string.format("  Magic (Magical Power): %d", char.Magic or 0))
    
    -- Equipment
    if config.showEquipment then
        table.insert(sections, "\n[EQUIPMENT]")
        table.insert(sections, string.format("  Weapon: %s", getEquipmentName(char.Weapon or 93)))
        table.insert(sections, string.format("  Shield: %s", getEquipmentName(char.Shield or 93)))
        table.insert(sections, string.format("  Armor: %s", getEquipmentName(char.Armor or 199)))
        table.insert(sections, string.format("  Helmet: %s", getEquipmentName(char.Helmet or 198)))
        table.insert(sections, string.format("  Relic 1: %s", getEquipmentName(char.Relic1 or 200)))
        table.insert(sections, string.format("  Relic 2: %s", getEquipmentName(char.Relic2 or 200)))
    end
    
    -- Spells
    if config.showSpells and char.Spells then
        local spellCount = 0
        for _ in pairs(char.Spells) do spellCount = spellCount + 1 end
        
        table.insert(sections, string.format("\n[LEARNED SPELLS] (%d total)", spellCount))
        table.insert(sections, formatSpells(char.Spells, config))
    end
    
    -- Commands
    if config.showCommands and char.Commands then
        table.insert(sections, string.format("\n[COMMANDS] (%d assigned)", #char.Commands))
        table.insert(sections, formatCommands(char.Commands))
    end
    
    -- Footer
    table.insert(sections, "\n" .. "=".rep(60))
    
    return table.concat(sections, "\n")
end

-- ============================================================================
-- Character Selection Dialog
-- ============================================================================
local function selectCharacter(api)
    -- Get all characters
    local characters = {}
    for i = 0, 15 do
        local char = api:GetCharacter(i)
        if char and char.Name then
            table.insert(characters, {index = i, name = char.Name, enabled = char.Enabled})
        end
    end
    
    if #characters == 0 then
        api:ShowDialog("Error", "No characters found in save file.")
        return nil
    end
    
    -- Build character selection prompt
    local prompt = "Select a character:\n\n"
    for i, char in ipairs(characters) do
        local status = char.enabled and "Active" or "Inactive"
        prompt = prompt .. string.format("%d. %s (%s)\n", i, char.name, status)
    end
    prompt = prompt .. "\nEnter character number (1-" .. #characters .. "):"
    
    -- Get user selection
    local input = api:ShowInput("Select Character", prompt, "1")
    if not input then
        return nil  -- User canceled
    end
    
    local selection = tonumber(input)
    if not selection or selection < 1 or selection > #characters then
        api:ShowDialog("Error", "Invalid selection. Please enter a number between 1 and " .. #characters)
        return nil
    end
    
    return characters[selection].index
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
    
    api:Log("info", "Character Stats Display plugin started")
    
    -- Select character
    local charIndex = selectCharacter(api)
    if not charIndex then
        api:Log("info", "User canceled character selection")
        return false
    end
    
    -- Get character data
    local char = api:GetCharacter(charIndex)
    if not char then
        api:ShowDialog("Error", "Failed to retrieve character data.")
        api:Log("error", "GetCharacter returned nil for index: " .. charIndex)
        return false
    end
    
    -- Build and display character info
    local display = buildCharacterDisplay(char, config)
    api:ShowDialog("Character Stats: " .. (char.Name or "Unknown"), display)
    
    api:Log("info", string.format("Displayed stats for character: %s (ID: %d)", char.Name, char.ID))
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
