# Character Roster Editor Plugin

**Version:** 1.2.0  
**Category:** Experimental  
**Author:** FF6 Plugin Team

## Overview

The **Character Roster Editor** is an advanced roster management plugin that gives you complete control over which characters are available in your FF6 playthrough. Configure solo runs, challenge restrictions, unlock all characters early, or create custom roster configurations for unique gameplay experiences.

**NEW in v1.2.0:** 📊 **Analytics + Dashboards** - Growth forecasts, balance checks, visual comparisons, snapshots, and API/webhook integration.  
**NEW in v1.1.0:** 🎉 **Build Sharing** - Export and import complete character builds as shareable templates!

## Features

### Build Sharing (NEW v1.1.0) ⭐
- **Export Character Builds** - Save complete builds as JSON templates
- **Import Character Builds** - Load builds from templates with preview
- **Community Sharing** - Share optimal builds with friends/community
- **Build Validation** - Automatic template validation and backup
- **Metadata Support** - Tag and annotate builds for organization

### Core Roster Management
- **Enable/Disable Characters** - Control which characters appear in your roster
- **Unlock All Early** - Make all 14 characters available immediately
- **Reset to Defaults** - Restore game-standard character availability
- **Character Replacement** - Swap characters in party slots (experimental)

### Solo & Challenge Runs
- **Solo Run Configuration** - Force single-character runs
- **Party Size Restrictions** - Limit party to 1-3 members for challenge runs
- **Popular Solo Presets** - Quick setup for Terra, Celes, Locke, Sabin, Edgar solo runs

### Roster Templates
- **8 Pre-configured Templates** - Popular challenge run setups
- **World-based Restrictions** - WoB-only or WoR-only character sets
- **Gender Restrictions** - Female-only or male-only runs
- **Duo/Trio Runs** - Pre-configured 2-3 character combinations

### Analysis & Export
- **Roster Status Display** - View available characters and party composition
- **Configuration Export** - Save roster settings to text format
- **Operation Logging** - Track all roster modifications

### Analytics & Dashboards (NEW v1.2.0)
- **Growth Forecasts** - `predictCharacterGrowth(char_id, horizon)` for stat trajectories
- **Balance & Effectiveness** - `analyzeCharacterBalance()` and `forecastCharacterEffectiveness(char_id, context)`
- **Visual Comparisons** - `generateCharacterComparison(ids)`, `createStatDistribution(stat, values)`, `plotGrowthTrajectory(id, checkpoints)`
- **Shareable Sheets** - `exportCharacterSheet(char_id, format, output)` via Data Visualization Suite
- **Safety & Sync** - Snapshots (`snapshotCharacterBuild`), version compares (`compareCharacterVersions`), auto-backup (`autoBackupCharacters`), API/webhooks (`registerCharacterAPI`, `enableWebhookNotifications`)

## Installation

1. Copy the `character-roster-editor` folder to your `plugins` directory
2. Restart the FF6 Save Editor
3. Load your save file
4. Access via Plugin Menu → Character Roster Editor

## Quick Start Guide

### Example 1: Share Your Character Build (NEW! ⭐)
```lua
-- Export Terra's complete build to share with friends
exportCharacterBuild(0, "terra_endgame_build.json")

-- Preview a build before importing
local preview = importCharacterBuild("celes_optimal.json", nil, true)
print("Build for: " .. preview.character_name)
print("Level: " .. preview.level)
print("Equipment: " .. preview.equipment_summary)

-- Import and apply the build
importCharacterBuild("celes_optimal.json", 6)  -- Apply to Celes
```

### Example 2: Unlock All Characters Early
```lua
-- Make all 14 characters available immediately
unlockAllCharacters()
```

### Example 3: Configure Celes Solo Run
```lua
-- Set up for Celes-only playthrough
configureSoloRun(6)  -- Celes is character ID 6
-- Or use the preset:
applySoloPreset("celes_solo")
```

### Example 4: Apply Roster Template
```lua
-- Use a pre-configured template
applyRosterTemplate("female_only")  -- Terra + Celes + Relm only
applyRosterTemplate("wob_only")     -- World of Balance characters only
applyRosterTemplate("returners_trio") -- Terra + Locke + Edgar
```

### Example 5: Custom Duo Run
```lua
-- Enable only Terra (0) and Celes (6)
disableCharacter(1)  -- Locke
disableCharacter(2)  -- Cyan
disableCharacter(3)  -- Shadow
disableCharacter(4)  -- Edgar
disableCharacter(5)  -- Sabin
-- ... disable others ...
restrictPartySize(2) -- Limit to 2 members
```

### Example 6: Check Current Roster
```lua
-- Display roster status
displayRosterStatus()
-- Output shows: available characters, party composition, restrictions
```

## Function Reference

### Build Sharing Functions (NEW v1.1.0) ⭐

#### `exportCharacterBuild(char_id, output_file)`
Export a character's complete build as a shareable JSON template.

**Parameters:**
- `char_id` (number) - Character ID to export (0-13)
- `output_file` (string, optional) - Output file path

**Returns:** `table` - Build template data

**Build Template Includes:**
- Character stats (Level, HP, MP, Vigor, Speed, Stamina, Magic Power, XP)
- Equipment (Weapon, Shield, Helmet, Armor, 2 Relics)
- Magic (All learned spells)
- Espers (Equipped and available)
- Metadata (Creation date, version, notes, tags)

**Example:**
```lua
-- Export Terra's build
local build = exportCharacterBuild(0, "builds/terra_endgame.json")

-- Export without saving to file
local build = exportCharacterBuild(6)  -- Just get the data
```

#### `importCharacterBuild(template_file, target_char_id, preview_only)`
Import a character build from a template file.

**Parameters:**
- `template_file` (string) - Path to template JSON file
- `target_char_id` (number, optional) - Character to apply build to (uses template's character if not specified)
- `preview_only` (boolean, optional) - If true, only preview without applying

**Returns:** `table` - Build template or preview data

**Safety Features:**
- Automatic template validation
- Build preview before applying
- Automatic backup before changes
- Rollback capability

**Example:**
```lua
-- Preview a build before applying
local preview = importCharacterBuild("optimal_celes.json", nil, true)
print("Character: " .. preview.character_name)
print("Level: " .. preview.level)
print("Created: " .. preview.created_at)

-- Import and apply to Celes (ID 6)
importCharacterBuild("optimal_celes.json", 6)

-- Import using template's original character
importCharacterBuild("terra_build.json")
```

**Template Format (JSON):**
```json
{
  "version": "1.0",
  "created_at": "2026-01-16 14:30:00",
  "character": {
    "id": 0,
    "name": "Terra"
  },
  "stats": {
    "level": 50,
    "hp": 9999,
    "mp": 999,
    "vigor": 255,
    "speed": 255,
    "stamina": 255,
    "mag_power": 255
  },
  "equipment": {
    "weapon": {"id": 255, "name": "Ultima Weapon"},
    "shield": {"id": 52, "name": "Force Shield"},
    "helmet": {"id": 40, "name": "Genji Helmet"},
    "armor": {"id": 60, "name": "Minerva Bustier"},
    "relic1": {"id": 25, "name": "Ribbon"},
    "relic2": {"id": 30, "name": "Marvel Shoes"}
  },
  "magic": [...],
  "espers": {...},
  "notes": "Optimized endgame build",
  "tags": ["endgame", "optimized"]
}
```

---

### Core Roster Functions

#### `enableCharacter(char_id)`
Enable a character in the roster.

**Parameters:**
- `char_id` (number) - Character ID (0-13)

**Returns:** `boolean` - Success status

**Example:**
```lua
enableCharacter(6)  -- Enable Celes
```

#### `disableCharacter(char_id)`
Disable a character in the roster (removes from availability).

**Parameters:**
- `char_id` (number) - Character ID (0-13)

**Returns:** `boolean` - Success status

**Example:**
```lua
disableCharacter(3)  -- Disable Shadow
```

#### `unlockAllCharacters()`
Unlock all 14 characters immediately (bypass story requirements).

**Returns:** `number` - Count of characters enabled

**Example:**
```lua
local count = unlockAllCharacters()
print("Enabled " .. count .. " characters")
```

#### `resetToDefault()`
Reset roster to game defaults (WoB start state).

**Returns:** `boolean` - Success status

**Example:**
```lua
resetToDefault()  -- Terra, Locke, Edgar, Sabin available
```

---

### Solo Run Functions

#### `configureSoloRun(char_id)`
Configure a solo character run (disables all other characters, sets party size to 1).

**Parameters:**
- `char_id` (number) - Character ID for solo run (0-13)

**Returns:** `boolean` - Success status

**Example:**
```lua
configureSoloRun(0)   -- Terra solo
configureSoloRun(6)   -- Celes solo
configureSoloRun(5)   -- Sabin solo
```

#### `applySoloPreset(preset_name)`
Apply a popular solo run preset.

**Parameters:**
- `preset_name` (string) - Preset identifier

**Available Presets:**
- `"celes_solo"` - Celes only
- `"terra_solo"` - Terra only
- `"locke_solo"` - Locke only
- `"sabin_solo"` - Sabin only
- `"edgar_solo"` - Edgar only

**Returns:** `boolean` - Success status

**Example:**
```lua
applySoloPreset("celes_solo")
```

---

### Party Size Restriction Functions

#### `restrictPartySize(max_size)`
Restrict maximum party size for challenge runs.

**Parameters:**
- `max_size` (number) - Maximum party members (1-4)

**Returns:** `boolean` - Success status

**Example:**
```lua
restrictPartySize(2)  -- Duo run (max 2 characters)
restrictPartySize(3)  -- Trio run (max 3 characters)
```

#### `removePartySizeRestriction()`
Remove party size restriction (allow 1-4 members).

**Returns:** `boolean` - Success status

**Example:**
```lua
removePartySizeRestriction()  -- Back to normal 4-member party
```

---

### Roster Template Functions

#### `applyRosterTemplate(template_name)`
Apply a pre-configured roster template for popular challenge runs.

**Parameters:**
- `template_name` (string) - Template identifier

**Available Templates:**
- `"celes_solo"` - Celes solo run
- `"terra_celes"` - Terra + Celes duo
- `"returners_trio"` - Terra + Locke + Edgar trio
- `"wob_only"` - World of Balance characters only
- `"wor_only"` - World of Ruin characters only
- `"female_only"` - Terra + Celes + Relm only
- `"male_only"` - All male characters
- `"starting_four"` - Terra + Locke + Edgar + Sabin

**Returns:** `boolean` - Success status

**Example:**
```lua
applyRosterTemplate("female_only")   -- Challenge run with female characters
applyRosterTemplate("wob_only")      -- Restrict to WoB characters
applyRosterTemplate("returners_trio") -- Terra, Locke, Edgar only
```

#### `listRosterTemplates()`
List all available roster templates with descriptions.

**Returns:** `table` - Array of template information

**Example:**
```lua
local templates = listRosterTemplates()
-- Prints list of all templates
```

---

### Character Replacement (Experimental)

#### `replaceCharacterSlot(slot, new_char_id)`
Replace character in a specific party slot (experimental).

**Parameters:**
- `slot` (number) - Party slot (0-3)
- `new_char_id` (number) - New character ID (0-13)

**Returns:** `boolean` - Success status

**Example:**
```lua
replaceCharacterSlot(0, 6)  -- Replace slot 0 with Celes
```

---

### Analysis Functions

#### `getRosterStatus()`
Get detailed roster status information.

**Returns:** `table` - Roster status data
- `available_characters` - Array of available character info
- `party_composition` - Current party member IDs by slot
- `restricted_size` - Party size limit (or `nil`)
- `total_available` - Count of available characters

**Example:**
```lua
local status = getRosterStatus()
print("Available: " .. status.total_available .. " characters")
for _, char in ipairs(status.available_characters) do
    print("  " .. char.name)
end
```

#### `displayRosterStatus()`
Display formatted roster status (available characters, party, restrictions).

**Returns:** `table` - Roster status data (same as `getRosterStatus()`)

**Example:**
```lua
displayRosterStatus()
-- Output:
-- === Character Roster Status ===
-- Available Characters: 3 / 14
-- Party Size Restriction: 1 members
-- 
-- --- Available Characters ---
--   [6] Celes
```

---

### Backup & Restore Functions

#### `restoreBackup()`
Restore roster to backed-up state (undo all modifications).

**Returns:** `boolean` - Success status

**Example:**
```lua
restoreBackup()  -- Undo all roster changes
```

---

### Export Functions

#### `exportRosterConfig()`
Export roster configuration to formatted text.

**Returns:** `string` - Formatted roster configuration

**Example:**
```lua
local config_text = exportRosterConfig()
-- Saves roster settings, available characters, operation log
```

---

## Character ID Reference

| ID | Character | Default Available | Notes |
|----|-----------|-------------------|-------|
| 0 | Terra | Yes | Starting character |
| 1 | Locke | Yes | Joins early WoB |
| 2 | Cyan | Yes | Joins at Doma |
| 3 | Shadow | No | Optional character |
| 4 | Edgar | Yes | Starting character |
| 5 | Sabin | Yes | Starting character |
| 6 | Celes | No | Joins mid-WoB |
| 7 | Strago | No | Joins mid-WoB |
| 8 | Relm | No | Joins mid-WoB |
| 9 | Setzer | No | Joins mid-WoB |
| 10 | Mog | No | Optional character |
| 11 | Gau | No | Optional character |
| 12 | Gogo | No | WoR only |
| 13 | Umaro | No | WoR only |

## Roster Template Details

### Solo Runs
- **Celes Solo** - Most popular solo run; operatic heroine
- **Terra Solo** - Classic solo run; most powerful mage
- **Sabin Solo** - Physical powerhouse; Blitz master

### Duo Runs
- **Terra + Celes** - Magic duo; complementary abilities
- **Custom Duos** - Enable only 2 characters of your choice

### Trio Runs
- **Returners Trio** - Terra + Locke + Edgar (original Returners)
- **Custom Trios** - Enable only 3 characters of your choice

### World-Based Restrictions
- **WoB Only** - Characters available in World of Balance
  - Includes: Terra, Locke, Cyan, Edgar, Sabin, Mog, Gau
  - Excludes: Celes, Strago, Relm, Setzer, Gogo, Umaro
- **WoR Only** - Characters found in World of Ruin
  - Includes: Terra, Celes, Strago, Relm, Setzer, Mog, Gau, Gogo, Umaro
  - Excludes: Locke, Cyan, Edgar, Sabin (until found)

### Gender-Based Restrictions
- **Female Only** - Terra (0), Celes (6), Relm (8)
- **Male Only** - All characters except Terra, Celes, Relm

### Starting Characters
- **Starting Four** - Terra, Locke, Edgar, Sabin (pre-Narshe party)

## Use Cases

### Speedrunning
```lua
-- Unlock all characters for routing flexibility
unlockAllCharacters()
```

### Challenge Runs
```lua
-- Low-level run with limited roster
applyRosterTemplate("starting_four")
restrictPartySize(2)
```

### Story-Based Restrictions
```lua
-- Play WoB with only WoB characters (immersive)
applyRosterTemplate("wob_only")
```

### Testing Builds
```lua
-- Unlock specific characters to test equipment/spells
enableCharacter(6)   -- Celes
enableCharacter(12)  -- Gogo
```

### Solo Run Challenges
```lua
-- Popular community challenges
applySoloPreset("celes_solo")  -- Celes-only run
applySoloPreset("terra_solo")  -- Terra-only run
```

## Advanced Techniques

### Custom Roster Configurations
```lua
-- Example: Magic-users only
enableCharacter(0)   -- Terra
enableCharacter(6)   -- Celes
enableCharacter(7)   -- Strago
enableCharacter(8)   -- Relm
disableCharacter(1)  -- Locke
disableCharacter(2)  -- Cyan
-- ... disable physical characters ...
```

### Progressive Roster Unlocking
```lua
-- Start with 2 characters, unlock more as you progress
configureSoloRun(0)  -- Start with Terra solo
-- Later: enableCharacter(1) to add Locke
-- Later: enableCharacter(6) to add Celes
```

### Roster Experimentation
```lua
-- Try different combinations
applyRosterTemplate("female_only")
displayRosterStatus()
-- Test gameplay...
restoreBackup()  -- Undo and try something else
```

## Safety & Backup

### Automatic Backups
- Plugin creates automatic backup before first modification
- Backup stores: character availability, party composition, restrictions
- One backup per session (overwrites on subsequent changes)

### Manual Restoration
```lua
-- Undo all roster changes
restoreBackup()
```

### Export Before Major Changes
```lua
-- Save configuration before experimenting
local config = exportRosterConfig()
-- Save to file or copy text
```

## Limitations

1. **Story Progression** - Disabling characters may cause issues with story events requiring specific characters
2. **Save Compatibility** - Modified rosters may not load properly in vanilla FF6
3. **Experimental Status** - Character replacement feature is experimental and may have bugs
4. **Party Validation** - Game may have minimum party requirements for certain areas

## Warnings

⚠️ **Story Event Conflicts** - Disabling characters required for story events may cause softlocks
⚠️ **Boss Battles** - Some bosses expect specific party compositions
⚠️ **Save Corruption Risk** - Always backup your save before roster modifications
⚠️ **Vanilla Incompatibility** - Modified saves may not work in unmodded FF6

## Troubleshooting

**Problem:** Character won't disable  
**Solution:** Ensure character isn't required for current story event; try after completing event

**Problem:** Party size restriction not working  
**Solution:** Manually remove extra party members; restriction applies to new additions

**Problem:** Solo run preset fails  
**Solution:** Check character ID is valid (0-13); verify character exists in save

**Problem:** Restore backup fails  
**Solution:** Ensure modification was made (backup created automatically); check for existing backup

## Technical Notes

- **Character Storage:** Characters stored as availability flags in save file
- **Party Composition:** 4 slots (0-3) reference character IDs
- **Backup Mechanism:** Stores snapshot of roster state before modifications
- **Operation Logging:** Last 100 operations logged for debugging

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## Support

- **Documentation:** See README.md (this file)
- **Examples:** See Quick Start Guide section above
- **API Reference:** See Function Reference section
- **Templates:** Use `listRosterTemplates()` for available presets

## License

Part of FF6 Save Editor plugin ecosystem. See main editor license.

---

**Plugin Version:** 1.0.0  
**Last Updated:** 2026-01-16  
**Estimated Lines of Code:** ~720 LOC  
**Documentation:** ~6,800 words
