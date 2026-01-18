-- Sprite Swapper Plugin v1.0.0
-- Swap character sprites with each other, NPCs, enemies, or custom imports

-- ============================================================================
-- CONFIGURATION
-- ============================================================================

local CONFIG = {
    -- Character count
    CHARACTER_COUNT = 14,
    
    -- NPC count
    NPC_COUNT = 50,
    
    -- Enemy count
    ENEMY_COUNT = 384,
    
    -- Custom sprite storage
    MAX_CUSTOM_SPRITES = 100,
    
    -- Sprite types
    SPRITE_TYPES = {
        CHARACTER = "character",
        NPC = "npc",
        ENEMY = "enemy",
        CUSTOM = "custom"
    },
    
    -- Sprite categories
    SPRITE_CATEGORIES = {
        BATTLE = "battle",        -- Battle sprite
        FIELD = "field",          -- Overworld/field sprite
        BOTH = "both"             -- Both battle and field
    },
    
    -- Character names for reference
    CHARACTER_NAMES = {
        [0] = "Terra",
        [1] = "Locke",
        [2] = "Edgar",
        [3] = "Sabin",
        [4] = "Shadow",
        [5] = "Cyan",
        [6] = "Gau",
        [7] = "Setzer",
        [8] = "Strago",
        [9] = "Relm",
        [10] = "Mog",
        [11] = "Umaro",
        [12] = "Kefka",
        [13] = "Gestahl"
    },
    
    -- Backup key
    BACKUP_KEY = "sprite_swapper_backup",
    
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
    sprite_cache = {},
    custom_sprites = {},
    sprite_swatches = {},
    current_swaps = {}
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
    
    print(string.format("[Sprite Swapper] %s: %s", operation_type, details))
end

-- Create backup
local function create_backup()
    local backup = {
        timestamp = os.time(),
        sprite_swaps = {},
        custom_sprites = {},
        original_sprites = {}
    }
    
    -- Backup current sprite configuration
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        backup.sprite_swaps[char_id] = plugin_state.current_swaps[char_id] or {
            type = CONFIG.SPRITE_TYPES.CHARACTER,
            original_id = char_id,
            current_id = char_id,
            battle = char_id,
            field = char_id
        }
    end
    
    -- Backup custom sprites metadata
    for sprite_id, sprite_data in pairs(plugin_state.custom_sprites) do
        backup.custom_sprites[sprite_id] = sprite_data
    end
    
    plugin_state.backups[CONFIG.BACKUP_KEY] = backup
    log_operation("BACKUP", "Created sprite configuration backup")
    return backup
end

-- Validate character ID
local function validate_char_id(char_id)
    return char_id >= 0 and char_id < CONFIG.CHARACTER_COUNT
end

-- Validate NPC ID
local function validate_npc_id(npc_id)
    return npc_id >= 0 and npc_id < CONFIG.NPC_COUNT
end

-- Validate enemy ID
local function validate_enemy_id(enemy_id)
    return enemy_id >= 0 and enemy_id < CONFIG.ENEMY_COUNT
end

-- Get character name
local function get_char_name(char_id)
    return CONFIG.CHARACTER_NAMES[char_id] or "Unknown"
end

-- ============================================================================
-- SPRITE SWAPPING FUNCTIONS
-- ============================================================================

-- Swap character sprites
function swapCharacterSprites(char_id_1, char_id_2)
    if not validate_char_id(char_id_1) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id_1))
        return false
    end
    
    if not validate_char_id(char_id_2) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id_2))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- game.swap_character_sprites(char_id_1, char_id_2)
    
    -- Store swap information
    plugin_state.current_swaps[char_id_1] = {
        type = CONFIG.SPRITE_TYPES.CHARACTER,
        original_id = char_id_1,
        current_id = char_id_2,
        timestamp = os.time()
    }
    
    log_operation("SWAP_CHAR", string.format("%s ↔ %s", 
        get_char_name(char_id_1), get_char_name(char_id_2)))
    
    return true
end

-- Swap with NPC sprite
function swapWithNPCSprite(char_id, npc_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if not validate_npc_id(npc_id) then
        log_operation("ERROR", "Invalid NPC ID: " .. tostring(npc_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- game.set_character_sprite(char_id, npc_id, "npc")
    
    plugin_state.current_swaps[char_id] = {
        type = CONFIG.SPRITE_TYPES.NPC,
        original_id = char_id,
        current_id = npc_id,
        timestamp = os.time()
    }
    
    log_operation("SWAP_NPC", string.format("%s → NPC #%d", 
        get_char_name(char_id), npc_id))
    
    return true
end

-- Swap with enemy sprite
function swapWithEnemySprite(char_id, enemy_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if not validate_enemy_id(enemy_id) then
        log_operation("ERROR", "Invalid enemy ID: " .. tostring(enemy_id))
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- game.set_character_sprite(char_id, enemy_id, "enemy")
    
    plugin_state.current_swaps[char_id] = {
        type = CONFIG.SPRITE_TYPES.ENEMY,
        original_id = char_id,
        current_id = enemy_id,
        timestamp = os.time()
    }
    
    log_operation("SWAP_ENEMY", string.format("%s → Enemy #%d", 
        get_char_name(char_id), enemy_id))
    
    return true
end

-- ============================================================================
-- CUSTOM SPRITE IMPORT FUNCTIONS
-- ============================================================================

-- Import custom sprite file
function importCustomSprite(sprite_name, file_path, sprite_type)
    if not sprite_name or sprite_name == "" then
        log_operation("ERROR", "Sprite name required")
        return false
    end
    
    if not file_path or file_path == "" then
        log_operation("ERROR", "File path required")
        return false
    end
    
    -- Validate file type (PNG, BMP)
    local extension = file_path:match("%.([^%.]+)$"):lower()
    if extension ~= "png" and extension ~= "bmp" then
        log_operation("ERROR", "Only PNG and BMP sprites supported")
        return false
    end
    
    sprite_type = sprite_type or CONFIG.SPRITE_CATEGORIES.BOTH
    
    create_backup()
    
    -- Generate custom sprite ID
    local custom_id = "custom_" .. tostring(os.time()) .. "_" .. sprite_name
    
    local custom_sprite = {
        id = custom_id,
        name = sprite_name,
        file_path = file_path,
        file_type = extension,
        sprite_type = sprite_type,
        import_date = os.time(),
        width = 0,  -- In real implementation, read from file
        height = 0, -- In real implementation, read from file
        size_kb = 0 -- In real implementation, get file size
    }
    
    -- In real implementation:
    -- Load sprite file, validate dimensions, compress if needed
    
    plugin_state.custom_sprites[custom_id] = custom_sprite
    
    log_operation("IMPORT", string.format("Imported custom sprite: %s (%s)", 
        sprite_name, extension:upper()))
    
    return true, custom_id
end

-- List custom sprites
function listCustomSprites()
    local sprites = {}
    
    for sprite_id, sprite_data in pairs(plugin_state.custom_sprites) do
        table.insert(sprites, {
            id = sprite_id,
            name = sprite_data.name,
            file_type = sprite_data.file_type,
            import_date = sprite_data.import_date,
            size_kb = sprite_data.size_kb
        })
    end
    
    return sprites
end

-- Delete custom sprite
function deleteCustomSprite(custom_id)
    if not plugin_state.custom_sprites[custom_id] then
        log_operation("ERROR", "Custom sprite not found: " .. custom_id)
        return false
    end
    
    create_backup()
    
    local sprite_name = plugin_state.custom_sprites[custom_id].name
    plugin_state.custom_sprites[custom_id] = nil
    
    log_operation("DELETE", string.format("Deleted custom sprite: %s", sprite_name))
    return true
end

-- Apply custom sprite to character
function applyCustomSprite(char_id, custom_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if not plugin_state.custom_sprites[custom_id] then
        log_operation("ERROR", "Custom sprite not found: " .. custom_id)
        return false
    end
    
    create_backup()
    
    -- In real implementation:
    -- game.set_character_sprite(char_id, custom_id, "custom")
    
    plugin_state.current_swaps[char_id] = {
        type = CONFIG.SPRITE_TYPES.CUSTOM,
        original_id = char_id,
        current_id = custom_id,
        timestamp = os.time()
    }
    
    local sprite_name = plugin_state.custom_sprites[custom_id].name
    log_operation("APPLY_CUSTOM", string.format("%s → Custom: %s", 
        get_char_name(char_id), sprite_name))
    
    return true
end

-- ============================================================================
-- SPRITE CATEGORY FUNCTIONS
-- ============================================================================

-- Swap battle sprite only
function swapBattleSprite(char_id, source_char_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if not validate_char_id(source_char_id) then
        log_operation("ERROR", "Invalid source character ID: " .. tostring(source_char_id))
        return false
    end
    
    create_backup()
    
    if not plugin_state.current_swaps[char_id] then
        plugin_state.current_swaps[char_id] = {
            original_id = char_id,
            battle = char_id,
            field = char_id
        }
    end
    
    plugin_state.current_swaps[char_id].battle = source_char_id
    
    log_operation("BATTLE_SPRITE", string.format("%s battle sprite → %s", 
        get_char_name(char_id), get_char_name(source_char_id)))
    
    return true
end

-- Swap field sprite only
function swapFieldSprite(char_id, source_char_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    if not validate_char_id(source_char_id) then
        log_operation("ERROR", "Invalid source character ID: " .. tostring(source_char_id))
        return false
    end
    
    create_backup()
    
    if not plugin_state.current_swaps[char_id] then
        plugin_state.current_swaps[char_id] = {
            original_id = char_id,
            battle = char_id,
            field = char_id
        }
    end
    
    plugin_state.current_swaps[char_id].field = source_char_id
    
    log_operation("FIELD_SPRITE", string.format("%s field sprite → %s", 
        get_char_name(char_id), get_char_name(source_char_id)))
    
    return true
end

-- ============================================================================
-- RESTORE & RESET FUNCTIONS
-- ============================================================================

-- Restore single character sprite
function restoreCharacterSprite(char_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return false
    end
    
    create_backup()
    
    plugin_state.current_swaps[char_id] = nil
    
    log_operation("RESTORE", string.format("Restored %s sprite", get_char_name(char_id)))
    
    return true
end

-- Restore all character sprites
function restoreAllSprites()
    create_backup()
    
    plugin_state.current_swaps = {}
    
    log_operation("RESTORE_ALL", "Restored all character sprites to original")
    
    return true
end

-- ============================================================================
-- PRESET MANAGEMENT (SWATCHES)
-- ============================================================================

-- Save sprite swatch (preset)
function saveSpriteSwatch(swatch_name)
    if not swatch_name or swatch_name == "" then
        log_operation("ERROR", "Swatch name required")
        return false
    end
    
    local swatch = {
        name = swatch_name,
        created_at = os.time(),
        swaps = {}
    }
    
    -- Copy current swaps
    for char_id, swap_data in pairs(plugin_state.current_swaps) do
        swatch.swaps[char_id] = {
            type = swap_data.type,
            current_id = swap_data.current_id
        }
    end
    
    plugin_state.sprite_swatches[swatch_name] = swatch
    
    log_operation("SAVE_SWATCH", string.format("Saved sprite swatch: %s", swatch_name))
    
    return true
end

-- Load sprite swatch
function loadSpriteSwatch(swatch_name)
    if not plugin_state.sprite_swatches[swatch_name] then
        log_operation("ERROR", "Swatch not found: " .. swatch_name)
        return false
    end
    
    create_backup()
    
    local swatch = plugin_state.sprite_swatches[swatch_name]
    plugin_state.current_swaps = {}
    
    for char_id, swap_data in pairs(swatch.swaps) do
        plugin_state.current_swaps[tonumber(char_id)] = swap_data
    end
    
    log_operation("LOAD_SWATCH", string.format("Loaded sprite swatch: %s", swatch_name))
    
    return true
end

-- List sprite swatches
function listSpriteSwatches()
    local swatches = {}
    
    for swatch_name, swatch_data in pairs(plugin_state.sprite_swatches) do
        table.insert(swatches, {
            name = swatch_name,
            created_at = swatch_data.created_at,
            swap_count = table.getn(swatch_data.swaps)
        })
    end
    
    return swatches
end

-- Delete sprite swatch
function deleteSpriteSwatch(swatch_name)
    if not plugin_state.sprite_swatches[swatch_name] then
        log_operation("ERROR", "Swatch not found: " .. swatch_name)
        return false
    end
    
    plugin_state.sprite_swatches[swatch_name] = nil
    
    log_operation("DELETE_SWATCH", string.format("Deleted sprite swatch: %s", swatch_name))
    
    return true
end

-- ============================================================================
-- INFORMATION & QUERY FUNCTIONS
-- ============================================================================

-- Get current character sprite info
function getCharacterSprite(char_id)
    if not validate_char_id(char_id) then
        log_operation("ERROR", "Invalid character ID: " .. tostring(char_id))
        return nil
    end
    
    local swap_data = plugin_state.current_swaps[char_id]
    
    if not swap_data then
        return {
            character = get_char_name(char_id),
            current_sprite = get_char_name(char_id),
            type = CONFIG.SPRITE_TYPES.CHARACTER,
            status = "original"
        }
    end
    
    return {
        character = get_char_name(char_id),
        current_sprite = swap_data.current_id,
        type = swap_data.type,
        status = "swapped",
        timestamp = swap_data.timestamp
    }
end

-- List available sprites by type
function listAvailableSprites(sprite_type)
    local sprites = {}
    
    if sprite_type == CONFIG.SPRITE_TYPES.CHARACTER or sprite_type == nil then
        for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
            table.insert(sprites, {
                id = char_id,
                name = get_char_name(char_id),
                type = CONFIG.SPRITE_TYPES.CHARACTER
            })
        end
    end
    
    if sprite_type == CONFIG.SPRITE_TYPES.CUSTOM or sprite_type == nil then
        for custom_id, sprite_data in pairs(plugin_state.custom_sprites) do
            table.insert(sprites, {
                id = custom_id,
                name = sprite_data.name,
                type = CONFIG.SPRITE_TYPES.CUSTOM
            })
        end
    end
    
    return sprites
end

-- ============================================================================
-- FUN OPERATIONS
-- ============================================================================

-- Apply random sprites to party
function applyRandomSprites(party_size)
    party_size = party_size or 4
    
    if party_size < 1 or party_size > CONFIG.CHARACTER_COUNT then
        log_operation("ERROR", "Invalid party size")
        return false
    end
    
    create_backup()
    
    local random_config = {}
    for slot = 0, party_size - 1 do
        local random_char = math.random(0, CONFIG.CHARACTER_COUNT - 1)
        table.insert(random_config, {slot = slot, target = random_char})
    end
    
    for slot, config in ipairs(random_config) do
        if validate_char_id(config.slot) then
            swapCharacterSprites(config.slot, config.target)
        end
    end
    
    log_operation("RANDOM", string.format("Applied random sprites to %d party members", party_size))
    
    return true
end

-- ============================================================================
-- EXPORT & ANALYSIS
-- ============================================================================

-- Export sprite configuration
function exportSpriteConfig()
    local status = {
        exported_at = os.date("%Y-%m-%d %H:%M:%S"),
        total_swaps = 0,
        custom_sprites_count = table.getn(plugin_state.custom_sprites),
        swatches_count = table.getn(plugin_state.sprite_swatches)
    }
    
    for _ in pairs(plugin_state.current_swaps) do
        status.total_swaps = status.total_swaps + 1
    end
    
    local export_text = "=== FF6 Sprite Swapper Configuration ===\n\n"
    export_text = export_text .. string.format("Export Date: %s\n", status.exported_at)
    export_text = export_text .. string.format("Total Sprite Swaps: %d\n", status.total_swaps)
    export_text = export_text .. string.format("Custom Sprites: %d\n", status.custom_sprites_count)
    export_text = export_text .. string.format("Saved Swatches: %d\n\n", status.swatches_count)
    
    export_text = export_text .. "=== Current Sprite Assignments ===\n"
    for char_id = 0, CONFIG.CHARACTER_COUNT - 1 do
        local sprite_info = getCharacterSprite(char_id)
        export_text = export_text .. string.format("[%d] %s: %s (%s)\n",
            char_id, sprite_info.character, sprite_info.current_sprite, sprite_info.status)
    end
    
    export_text = export_text .. "\n=== Custom Sprites ===\n"
    for sprite_id, sprite_data in pairs(plugin_state.custom_sprites) do
        export_text = export_text .. string.format("  %s: %s (%s)\n",
            sprite_id, sprite_data.name, sprite_data.file_type:upper())
    end
    
    export_text = export_text .. "\n=== Saved Swatches ===\n"
    for swatch_name, swatch_data in pairs(plugin_state.sprite_swatches) do
        export_text = export_text .. string.format("  %s (Created: %s, %d swaps)\n",
            swatch_name, os.date("%Y-%m-%d %H:%M:%S", swatch_data.created_at),
            table.getn(swatch_data.swaps))
    end
    
    export_text = export_text .. "\n=== Recent Operations (Last 10) ===\n"
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
-- BACKUP & RESTORE
-- ============================================================================

-- Restore from backup
function restoreBackup()
    local backup = plugin_state.backups[CONFIG.BACKUP_KEY]
    if not backup then
        log_operation("ERROR", "No backup found to restore")
        return false
    end
    
    plugin_state.current_swaps = backup.sprite_swaps
    plugin_state.custom_sprites = backup.custom_sprites
    
    log_operation("RESTORE", "Restored sprite configuration from backup")
    return true
end

-- ============================================================================
-- INITIALIZATION
-- ============================================================================

function initialize()
    if plugin_state.initialized then
        return true
    end
    
    log_operation("INIT", "Sprite Swapper initialized")
    plugin_state.initialized = true
    
    return true
end

-- Auto-initialize
initialize()

-- ============================================================================
-- PLUGIN INFO
-- ============================================================================

print("Sprite Swapper v1.0.0 loaded")
print("Commands: swapCharacterSprites(c1, c2), swapWithNPCSprite(char, npc)")
print("          importCustomSprite(name, path, type), applyCustomSprite(char, id)")
print("          swapBattleSprite(char, source), swapFieldSprite(char, source)")
print("          restoreCharacterSprite(char), saveSpriteSwatch(name)")
print("          applyRandomSprites(count), exportSpriteConfig()")
print("Type listAvailableSprites() to see all available sprites")
print("Type listCustomSprites() to see imported custom sprites")
