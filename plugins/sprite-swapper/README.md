# Sprite Swapper v1.0.0 - Complete Documentation

**Swap, customize, and import sprites to transform character appearance in Final Fantasy VI**

## Overview

The Sprite Swapper is a comprehensive cosmetic customization system that lets you swap character sprites with each other, replace them with NPC or enemy sprites, or import your own custom sprite files. Perfect for creating unique party compositions, cosmetic mods, visual challenges, and personal favorite combinations.

## Quick Start

```lua
-- Simple sprite swap
swapCharacterSprites(0, 1)  -- Terra becomes Locke

-- Import custom sprite
importCustomSprite("CustomTerra", "sprites/terra_custom.png", "both")
applyCustomSprite(0, "custom_...")  -- Apply to Terra

-- Save and load configurations
saveSpriteSwatch("MyCustomParty")
loadSpriteSwatch("MyCustomParty")

-- View current state
getCharacterSprite(0)
exportSpriteConfig()
```

## Table of Contents

1. [Features](#features)
2. [Character Reference](#character-reference)
3. [Basic Sprite Swapping](#basic-sprite-swapping)
4. [Custom Sprite Import](#custom-sprite-import)
5. [Sprite Categories](#sprite-categories)
6. [Preset Management](#preset-management)
7. [Advanced Operations](#advanced-operations)
8. [Use Cases](#use-cases)
9. [API Reference](#api-reference)
10. [Troubleshooting](#troubleshooting)

## Features

### Core Capabilities

✅ **Character-to-Character Swaps** - Replace any character with another  
✅ **NPC Sprites** - Use NPC appearances for characters (~50 NPCs available)  
✅ **Enemy Sprites** - Transform characters into enemies (384 enemy sprites)  
✅ **Custom Sprite Import** - Load your own PNG/BMP sprite files  
✅ **Custom Sprite Management** - Store, organize, and manage up to 100 custom sprites  
✅ **Battle vs. Field** - Swap sprites independently for battle and overworld  
✅ **Sprite Presets** - Save and load sprite configurations (swatches)  
✅ **Random Sprites** - Apply random sprite combinations  
✅ **Full Restore** - Revert all changes to original sprites  
✅ **Export/Analysis** - Export sprite configuration for sharing  
✅ **Operation Logging** - Track all sprite modifications  

### Advanced Features

- Import multiple custom sprites and build library
- Create themed party configurations
- Swap only battle OR field sprites (keep one normal)
- Delete unused custom sprites to save space
- Save multiple sprite presets for different playstyles
- Preview current sprite assignments
- List all available sprite sources

## Character Reference

FF6 has 14 main playable characters:

| ID | Name | Default Battle Sprite | Default Field Sprite |
|----|------|----------------------|----------------------|
| 0 | **Terra** | Esper form / Magic-user | Mage robes |
| 1 | **Locke** | Thief gear | Casual outfit |
| 2 | **Edgar** | Royal blue armor | Mechanical suit |
| 3 | **Sabin** | Martial arts gi | Fighting clothes |
| 4 | **Shadow** | Assassin garb | Dark cloak |
| 5 | **Cyan** | Doma knight armor | Samurai suit |
| 6 | **Gau** | Feral rage form | Barbarian outfit |
| 7 | **Setzer** | Gambler outfit | Elegant suit |
| 8 | **Strago** | Mage robes | Adventurer clothing |
| 9 | **Relm** | Elemental robes | Artist garb |
| 10 | **Mog** | Moogle body | Moogle appearance |
| 11 | **Umaro** | Yeti form | Yeti body |
| 12 | **Kefka** | Emperor robes | Emperor appearance |
| 13 | **Gestahl** | Emperor armor | Emperor suit |

## Basic Sprite Swapping

### Swap Two Characters

```lua
-- Terra and Locke swap sprites
swapCharacterSprites(0, 1)

-- Terra now looks like Locke, Locke looks like Terra
-- In battle and on field
```

The most straightforward customization - simply exchange sprite data between two characters.

**Common Combinations:**
```lua
-- All female characters look similar (fun mod)
swapCharacterSprites(0, 9)   -- Terra ↔ Relm
swapCharacterSprites(1, 1)   -- (Locke already)

-- All male warriors look stronger
swapCharacterSprites(2, 5)   -- Edgar ↔ Cyan
swapCharacterSprites(3, 11)  -- Sabin ↔ Umaro

-- Gender-swapped party
swapCharacterSprites(0, 2)   -- Terra ↔ Edgar
swapCharacterSprites(1, 9)   -- Locke ↔ Relm
```

### Replace with NPC Sprite

```lua
-- Make Terra look like a specific NPC
swapWithNPCSprite(0, 15)  -- NPC #15 (example)

-- Now Terra uses NPC sprite #15
```

There are ~50 NPC sprites available, including:
- Shopkeepers
- Innkeepers
- Townspeople
- Soldiers
- Scientists
- Magitek knights
- Espers

**Use Cases:**
```lua
-- Party of scientists
swapWithNPCSprite(0, 20)  -- NPC scientist sprite
swapWithNPCSprite(1, 20)
swapWithNPCSprite(2, 20)

-- Mix soldiers and civilians
swapWithNPCSprite(0, 10)  -- NPC soldier
swapWithNPCSprite(1, 12)  -- NPC civilian
```

### Replace with Enemy Sprite

```lua
-- Make characters look like monsters/enemies
swapWithEnemySprite(0, 50)   -- First tier enemy
swapWithEnemySprite(1, 150)  -- Mid-tier enemy
swapWithEnemySprite(2, 300)  -- Boss-tier enemy
```

FF6 has 384 enemy sprites, from basic monsters to powerful bosses:
- Common enemies (Soldier, Flan, etc.) - IDs 0-100
- Mid-tier enemies (Undead, Machines) - IDs 100-250
- Powerful enemies (Dragons, Undead bosses) - IDs 250-350
- Mega bosses (Atma Weapon, Kefka, etc.) - IDs 350-384

**Monster Party:**
```lua
-- Create a monster party
swapWithEnemySprite(0, 50)   -- Flan
swapWithEnemySprite(1, 75)   -- Bomb
swapWithEnemySprite(2, 100)  -- Ice Dragon
swapWithEnemySprite(3, 120)  -- Undead
```

## Custom Sprite Import

### Import Sprite File

```lua
-- Import custom sprite
importCustomSprite(
    "CustomTerra",                    -- Name for reference
    "C:/Sprites/terra_custom.png",   -- File path
    "both"                            -- Type: "battle", "field", or "both"
)

-- Returns: true, custom_sprite_id
```

**Supported Formats:**
- PNG (Portable Network Graphics) - Recommended
- BMP (Bitmap) - Also supported

**Sprite Dimensions:**
- Battle sprite: 72x72 pixels (standard)
- Field sprite: 24x32 pixels (standard)
- Can be slightly larger, will be scaled

**File Organization Tips:**
```
My Sprites/
├── terra_custom.png
├── locke_custom.png
├── edgar_custom.png
├── sabin_custom.png
├── shadow_custom.png
├── cyan_custom.png
├── gau_custom.png
├── setzer_custom.png
├── strago_custom.png
├── relm_custom.png
├── mog_custom.png
├── umaro_custom.png
├── kefka_custom.png
└── gestahl_custom.png
```

### Batch Import Sprites

```lua
-- Import multiple custom sprites
local sprites_to_import = {
    {name = "Pixel Terra", path = "sprites/terra_pixel.png"},
    {name = "Chibi Terra", path = "sprites/terra_chibi.png"},
    {name = "Anime Terra", path = "sprites/terra_anime.png"},
}

for _, sprite in ipairs(sprites_to_import) do
    importCustomSprite(sprite.name, sprite.path, "both")
end
```

### List Custom Sprites

```lua
-- View all imported custom sprites
local custom_sprites = listCustomSprites()

for _, sprite in ipairs(custom_sprites) do
    print(string.format("%s (%s): %s", 
        sprite.name, sprite.file_type:upper(), sprite.id))
end
```

Output:
```
Pixel Terra (PNG): custom_1234567_PixelTerra
Chibi Terra (PNG): custom_1234568_ChibiTerra
Anime Terra (PNG): custom_1234569_AnimaTerra
```

### Apply Custom Sprite

```lua
-- Apply custom sprite to character
applyCustomSprite(0, "custom_1234567_PixelTerra")

-- Terra now uses the custom pixel sprite
```

### Delete Custom Sprite

```lua
-- Remove custom sprite from library
deleteCustomSprite("custom_1234567_PixelTerra")

-- Frees up space for new imports
```

### Example: Complete Custom Sprite Workflow

```lua
-- Step 1: Import custom sprites
importCustomSprite("ChibiTerra", "sprites/terra_chibi.png", "both")
importCustomSprite("ChibiLocke", "sprites/locke_chibi.png", "both")
importCustomSprite("ChibiEdgar", "sprites/edgar_chibi.png", "both")
importCustomSprite("ChibiSabin", "sprites/sabin_chibi.png", "both")

-- Step 2: Apply to party
local custom_ids = {}
for sprite_id in listCustomSprites() do
    table.insert(custom_ids, sprite_id.id)
end

applyCustomSprite(0, custom_ids[1])
applyCustomSprite(1, custom_ids[2])
applyCustomSprite(2, custom_ids[3])
applyCustomSprite(3, custom_ids[4])

-- Step 3: Save configuration
saveSpriteSwatch("ChibiParty")

-- Step 4: Later, load it back
loadSpriteSwatch("ChibiParty")
```

## Sprite Categories

### Battle vs. Field Sprites

FF6 stores two separate sprite graphics for each character:
- **Battle Sprite** - Used during combat
- **Field Sprite** - Used on the overworld and indoors

By default, swapping both changes both. You can also swap them independently:

### Swap Battle Sprite Only

```lua
-- Change battle appearance only
swapBattleSprite(0, 2)

-- Terra's battle sprite becomes Edgar's
-- Terra's field sprite stays the same
```

**Use Case:**
```lua
-- Create a mismatch for comedy
-- Field: Terra looks normal (blue mage)
-- Battle: Terra looks like Kefka (for surprise)
swapBattleSprite(0, 12)
```

### Swap Field Sprite Only

```lua
-- Change overworld appearance only
swapFieldSprite(0, 1)

-- Terra's field sprite becomes Locke's
-- Terra's battle sprite stays the same
```

**Use Case:**
```lua
-- Create visual storytelling
-- Field: Party members look like normal people
-- Battle: Reveal their true monster forms
swapFieldSprite(0, 50)   -- Keep normal field look
swapBattleSprite(0, 100) -- But battle as monster
```

### Independent Category Example

```lua
-- Create unique character
-- Field: Looks like Locke (civilian appearance)
-- Battle: Looks like Cyan (warrior form)

swapFieldSprite(0, 1)    -- Locke field sprite
swapBattleSprite(0, 5)   -- Cyan battle sprite

-- Result: Character looks like Locke outside, Cyan in combat
```

## Preset Management

### Save Sprite Configuration

```lua
-- Save current sprite setup as preset
saveSpriteSwatch("FemaleParty")

-- Stores all current sprite assignments
-- Can load it again later
```

### Load Sprite Configuration

```lua
-- Load saved preset
loadSpriteSwatch("FemaleParty")

-- Restores all sprites to saved configuration
```

### List Saved Swatches

```lua
-- View all saved presets
local swatches = listSpriteSwatches()

for _, swatch in ipairs(swatches) do
    print(string.format("%s (%d swaps, created %s)",
        swatch.name, swatch.swap_count, os.date("%Y-%m-%d", swatch.created_at)))
end
```

Output:
```
FemaleParty (3 swaps, created 2026-01-16)
MonsterParty (6 swaps, created 2026-01-16)
RandomTheme (4 swaps, created 2026-01-16)
```

### Delete Saved Swatch

```lua
-- Remove preset
deleteSpriteSwatch("OldConfiguration")
```

### Swatch Management Examples

```lua
-- Create themed parties
saveSpriteSwatch("AllFemales")
saveSpriteSwatch("AllMales")
saveSpriteSwatch("PixelArtTheme")
saveSpriteSwatch("RetroTheme")
saveSpriteSwatch("CosplayParty")

-- Switch between them
loadSpriteSwatch("AllFemales")    -- Speedrun attempt 1
loadSpriteSwatch("MonsterParty")  -- Challenge run
loadSpriteSwatch("RetroTheme")    -- Nostalgia playthrough
```

## Advanced Operations

### Apply Random Sprites

```lua
-- Randomize party appearance
applyRandomSprites(4)  -- Randomize first 4 characters

-- Each character gets random sprite from available pool
```

**Use Cases:**
```lua
-- Challenge mode: Never know who looks like who
applyRandomSprites(4)

-- Replay randomizer
applyRandomSprites(6)  -- Randomize 6-character party

-- Comedy run
applyRandomSprites(14)  -- All 14 characters randomized
```

### Get Current Sprite Info

```lua
-- Check current sprite assignment
local info = getCharacterSprite(0)

-- Returns:
-- {
--   character = "Terra",
--   current_sprite = 1,
--   type = "character",
--   status = "swapped",
--   timestamp = 1673894400
-- }
```

### List Available Sprites

```lua
-- View all available sprite sources
local all_sprites = listAvailableSprites()

-- Filter by type
local char_sprites = listAvailableSprites("character")
local custom_sprites = listAvailableSprites("custom")

-- Results in list of {id, name, type}
```

## Use Cases

### 1. All-Female Party

```lua
-- Swap male characters to look female
swapCharacterSprites(2, 0)   -- Edgar → Terra
swapCharacterSprites(3, 9)   -- Sabin → Relm
swapCharacterSprites(5, 0)   -- Cyan → Terra (duplicate OK)
-- Alternatively: Locke, Strago, Mog, Umaro can also be swapped

-- Save configuration
saveSpriteSwatch("FemaleParty")
```

### 2. All-Male Party

```lua
-- Swap female characters to look male
swapCharacterSprites(0, 2)   -- Terra → Edgar
swapCharacterSprites(9, 5)   -- Relm → Cyan
swapCharacterSprites(8, 3)   -- Strago → Sabin

saveSpriteSwatch("MaleParty")
```

### 3. Custom Mod Appearance

```lua
-- Import custom sprite pack
importCustomSprite("Amano Art - Terra", "sprites/amano_terra.png", "both")
importCustomSprite("Amano Art - Locke", "sprites/amano_locke.png", "both")
importCustomSprite("Amano Art - Edgar", "sprites/amano_edgar.png", "both")
importCustomSprite("Amano Art - Sabin", "sprites/amano_sabin.png", "both")

-- Apply to party
applyCustomSprite(0, "custom_...")  -- Get actual IDs from import
applyCustomSprite(1, "custom_...")
applyCustomSprite(2, "custom_...")
applyCustomSprite(3, "custom_...")

saveSpriteSwatch("AmanoArtTheme")
```

### 4. Monster Party Challenge

```lua
-- Transform entire party into monsters
swapWithEnemySprite(0, 50)   -- Flan
swapWithEnemySprite(1, 75)   -- Bomb
swapWithEnemySprite(2, 100)  -- Ice Dragon
swapWithEnemySprite(3, 120)  -- Undead boss
swapWithEnemySprite(4, 150)  -- Phantom
swapWithEnemySprite(5, 200)  -- Machine
swapWithEnemySprite(6, 100)  -- Ice Dragon
swapWithEnemySprite(7, 120)  -- Undead boss

saveSpriteSwatch("MonsterParty")
```

### 5. NPC Squad Cosplay

```lua
-- Everyone looks like civilians/NPCs
swapWithNPCSprite(0, 10)  -- Innkeeper woman
swapWithNPCSprite(1, 15)  -- Shopkeeper man
swapWithNPCSprite(2, 20)  -- Soldier
swapWithNPCSprite(3, 12)  -- Town civilian

saveSpriteSwatch("CivilianParty")
```

### 6. Storytelling - Secret Forms

```lua
-- On field: look like normal people
swapFieldSprite(0, 1)  -- Terra looks like Locke
swapFieldSprite(1, 2)  -- Locke looks like Edgar
swapFieldSprite(2, 3)  -- Edgar looks like Sabin

-- In battle: reveal true forms
swapBattleSprite(0, 50)   -- Terra is actually monster
swapBattleSprite(1, 100)  -- Locke is actually beast
swapBattleSprite(2, 150)  -- Edgar is actually demon

-- Players won't expect it until combat!
```

### 7. Pixel Art Throwback

```lua
-- Import pixel art style sprites
importCustomSprite("FF1 Style - Terra", "sprites/ff1_terra.png", "both")
importCustomSprite("FF1 Style - Locke", "sprites/ff1_locke.png", "both")
importCustomSprite("FF1 Style - Edgar", "sprites/ff1_edgar.png", "both")
importCustomSprite("FF1 Style - Sabin", "sprites/ff1_sabin.png", "both")

-- Create throwback party
saveSpriteSwatch("FF1StyleParty")
```

## API Reference

### Sprite Swapping Functions

```lua
swapCharacterSprites(char_id_1: number, char_id_2: number) -> boolean
-- Swap two characters' sprites
-- Returns: true on success, false on error

swapWithNPCSprite(char_id: number, npc_id: number) -> boolean
-- Replace character sprite with NPC sprite
-- npc_id: 0-49 (50 NPC sprites available)

swapWithEnemySprite(char_id: number, enemy_id: number) -> boolean
-- Replace character sprite with enemy sprite
-- enemy_id: 0-383 (384 enemy sprites available)
```

### Custom Sprite Functions

```lua
importCustomSprite(name: string, file_path: string, sprite_type: string) -> boolean, custom_id
-- Import custom sprite file (PNG or BMP)
-- sprite_type: "battle", "field", or "both"
-- Returns: success, custom_sprite_id (use for applying)

listCustomSprites() -> table
-- Get all imported custom sprites
-- Returns: array of {id, name, file_type, import_date, size_kb}

deleteCustomSprite(custom_id: string) -> boolean
-- Remove custom sprite from library

applyCustomSprite(char_id: number, custom_id: string) -> boolean
-- Apply custom sprite to character
```

### Category Functions

```lua
swapBattleSprite(char_id: number, source_char_id: number) -> boolean
-- Change only battle sprite

swapFieldSprite(char_id: number, source_char_id: number) -> boolean
-- Change only field sprite
```

### Restore Functions

```lua
restoreCharacterSprite(char_id: number) -> boolean
-- Restore single character to original sprite

restoreAllSprites() -> boolean
-- Restore all characters to original sprites
```

### Preset Functions

```lua
saveSpriteSwatch(swatch_name: string) -> boolean
-- Save current sprite configuration

loadSpriteSwatch(swatch_name: string) -> boolean
-- Load saved sprite configuration

listSpriteSwatches() -> table
-- Get all saved presets
-- Returns: array of {name, created_at, swap_count}

deleteSpriteSwatch(swatch_name: string) -> boolean
-- Delete saved preset
```

### Query Functions

```lua
getCharacterSprite(char_id: number) -> table
-- Get current sprite assignment info
-- Returns: {character, current_sprite, type, status, timestamp}

listAvailableSprites(sprite_type: string) -> table
-- List available sprites
-- sprite_type: "character", "npc", "enemy", "custom", or nil for all
-- Returns: array of {id, name, type}
```

### Fun Functions

```lua
applyRandomSprites(party_size: number) -> boolean
-- Apply random sprites to party
-- party_size: 1-14 (how many characters to randomize)
```

### Export Functions

```lua
exportSpriteConfig() -> string
-- Export sprite configuration for sharing/backup
-- Returns: formatted text with all current assignments
```

## Troubleshooting

**Issue:** Can't import sprite
- Check file is PNG or BMP format
- Verify file path is correct
- Check sprite dimensions are reasonable (24x32 to 72x72 pixels)
- Ensure file is readable and not corrupted

**Issue:** Swapped sprite doesn't appear
- Verify character ID is valid (0-13)
- Check that swap was applied (use getCharacterSprite() to verify)
- Some swaps may not display until next battle/area change

**Issue:** Can't load swatch
- Verify swatch name matches exactly (case-sensitive)
- Check swatch was saved (use listSpriteSwatches())
- Try restoring from manual backup first

**Issue:** Custom sprite file too large
- Recommended: Compress sprite PNG files
- Max file size: typically under 50 KB per sprite
- Consider using indexed PNG for smaller files

## Advanced Tips

1. **Layered Configurations** - Save multiple swatches for different runs (speedrun, challenge, casual)
2. **Thematic Parties** - Create sprite themes (all same artist style, all anime, all pixel art)
3. **Story Mods** - Use swaps to create visual narrative (field vs. battle reveals)
4. **Comedy Runs** - Randomize sprites for chaotic playthroughs
5. **Art Appreciation** - Create fan art party combinations and share
6. **Accessibility** - Swap sprites for visual clarity (use contrasting colors)

## Version History

**v1.0.0 (January 16, 2026)** - Initial Release
- Character-to-character sprite swaps
- NPC sprite replacement (50 NPCs)
- Enemy sprite replacement (384 enemies)
- Custom sprite import (PNG/BMP)
- Custom sprite library management (100 sprites)
- Battle vs. field sprite control
- Preset/swatch system
- Random sprite application
- Full restore functionality
- Export/analysis tools
- 20 core functions, complete documentation

---

**Sprite Swapper v1.0.0** - *Customize Your Party's Appearance*
