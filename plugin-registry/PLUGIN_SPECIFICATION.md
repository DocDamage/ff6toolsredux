# Plugin Specification v1.0

**Version:** 1.0.0  
**Last Updated:** January 16, 2026  
**Status:** Stable

---

## Overview

This document defines the technical specification for FF6 Save Editor plugins. All plugins submitted to the official registry must comply with this specification.

---

## Table of Contents

1. [Plugin Structure](#plugin-structure)
2. [File Requirements](#file-requirements)
3. [Metadata Format](#metadata-format)
4. [Code Requirements](#code-requirements)
5. [Permissions](#permissions)
6. [API Contract](#api-contract)
7. [Security Requirements](#security-requirements)
8. [Versioning](#versioning)
9. [Validation Rules](#validation-rules)

---

## Plugin Structure

Every plugin must follow this directory structure:

```
plugins/plugin-id/
├── plugin.lua          # Main plugin code (required)
├── metadata.json       # Plugin metadata (required)
├── README.md           # Plugin documentation (required)
├── CHANGELOG.md        # Version history (required)
├── LICENSE             # License file (optional, defaults to MIT)
├── screenshot.png      # Screenshot (recommended, max 2MB)
└── checksum.sha256     # SHA256 checksum (auto-generated)
```

### File Size Limits

- **plugin.lua:** Maximum 10MB
- **metadata.json:** Maximum 100KB
- **README.md:** Maximum 1MB
- **CHANGELOG.md:** Maximum 500KB
- **screenshot.png:** Maximum 2MB, recommended 800x600px

---

## File Requirements

### 1. plugin.lua

**Required:** Yes  
**Format:** UTF-8 encoded Lua script  
**Extension:** `.lua` (lowercase)

**Structure:**

```lua
-- ============================================================================
-- Plugin Metadata (Required in comments)
-- ============================================================================
-- @id: unique-plugin-id
-- @name: Human Readable Name
-- @version: 1.0.0
-- @author: Author Name
-- @description: Brief description
-- @permissions: read_save, ui_display

-- ============================================================================
-- Main Entry Point (Required)
-- ============================================================================
function main(api)
    -- Plugin implementation
    return true  -- Return true on success, false on failure
end
```

**Requirements:**
- Must define a `main(api)` function
- Must return boolean (true = success, false = failure)
- Must include metadata comments at top of file
- Must use proper error handling (pcall recommended)

### 2. metadata.json

**Required:** Yes  
**Format:** JSON (RFC 8259 compliant)  
**Validation:** Must conform to [plugin-schema.json](schema/plugin-schema.json)

**Complete Schema:**

```json
{
  "id": "string (required)",
  "name": "string (required)",
  "version": "string (required, semver format)",
  "author": "string (required)",
  "contact": "string (required, email or URL)",
  "description": "string (required, max 200 chars)",
  "longDescription": "string (required, max 2000 chars)",
  "category": "string (required, enum)",
  "tags": ["array of strings (required, max 5)"],
  "permissions": ["array of strings (required)"],
  "minEditorVersion": "string (required, semver format)",
  "homepage": "string (required, URL)",
  "repository": "string (optional, URL)",
  "license": "string (optional, SPDX identifier)",
  "dependencies": ["array of strings (optional)"],
  "screenshots": ["array of strings (optional)"],
  "documentation": "string (optional, default: README.md)",
  "changelog": "string (optional, default: CHANGELOG.md)"
}
```

**Field Specifications:**

#### id (required)
- **Type:** String
- **Format:** Lowercase alphanumeric with hyphens only
- **Pattern:** `^[a-z0-9]+(-[a-z0-9]+)*$`
- **Min Length:** 3 characters
- **Max Length:** 50 characters
- **Example:** `"stats-display"`, `"item-manager"`
- **Must be unique** across all plugins in registry

#### name (required)
- **Type:** String
- **Min Length:** 3 characters
- **Max Length:** 50 characters
- **Example:** `"Character Stats Display"`
- **Must be human-readable**

#### version (required)
- **Type:** String
- **Format:** Semantic Versioning (MAJOR.MINOR.PATCH)
- **Pattern:** `^[0-9]+\.[0-9]+\.[0-9]+$`
- **Example:** `"1.0.0"`, `"2.1.3"`

#### author (required)
- **Type:** String
- **Min Length:** 2 characters
- **Max Length:** 100 characters
- **Example:** `"John Doe"`, `"FF6 Editor Team"`

#### contact (required)
- **Type:** String
- **Format:** Email address or URL
- **Example:** `"john@example.com"`, `"https://github.com/johndoe"`

#### description (required)
- **Type:** String
- **Min Length:** 20 characters
- **Max Length:** 200 characters
- **Example:** `"Display comprehensive character stats including HP, MP, equipment, and spells"`
- **Shown in marketplace listing**

#### longDescription (required)
- **Type:** String
- **Min Length:** 50 characters
- **Max Length:** 2000 characters
- **Example:** Long-form description with features, use cases, etc.
- **Shown in plugin details page**

#### category (required)
- **Type:** String (enum)
- **Valid Values:**
  - `"utility"` - General-purpose tools
  - `"automation"` - Batch operations, automated tasks
  - `"analysis"` - Stats analysis, optimization
  - `"enhancement"` - UI improvements, quality-of-life
  - `"debug"` - Development and debugging tools

#### tags (required)
- **Type:** Array of strings
- **Min Items:** 1
- **Max Items:** 5
- **Item Length:** 2-20 characters each
- **Format:** Lowercase, alphanumeric with hyphens
- **Example:** `["stats", "character", "display", "analysis"]`

#### permissions (required)
- **Type:** Array of strings
- **Min Items:** 1
- **Valid Values:**
  - `"read_save"` - Read save data
  - `"write_save"` - Modify save data
  - `"ui_display"` - Show dialogs to user
  - `"events"` - Register event hooks
- **Example:** `["read_save", "ui_display"]`

#### minEditorVersion (required)
- **Type:** String
- **Format:** Semantic Versioning
- **Example:** `"3.4.0"`
- **Must be >= 3.4.0**

#### homepage (required)
- **Type:** String
- **Format:** Valid HTTP/HTTPS URL
- **Example:** `"https://github.com/user/plugin"`

#### repository (optional)
- **Type:** String
- **Format:** Valid HTTP/HTTPS URL
- **Example:** `"https://github.com/user/plugin"`

#### license (optional)
- **Type:** String
- **Format:** SPDX License Identifier
- **Example:** `"MIT"`, `"GPL-3.0"`, `"Apache-2.0"`
- **Default:** `"MIT"` if not specified

#### dependencies (optional)
- **Type:** Array of strings
- **Format:** Array of plugin IDs
- **Example:** `["dependency-plugin-id"]`
- **Note:** Dependencies must be available in registry

#### screenshots (optional)
- **Type:** Array of strings
- **Format:** Array of filenames
- **Example:** `["screenshot.png", "screenshot2.png"]`
- **Files must exist in plugin directory**

#### documentation (optional)
- **Type:** String
- **Format:** Filename
- **Default:** `"README.md"`

#### changelog (optional)
- **Type:** String
- **Format:** Filename
- **Default:** `"CHANGELOG.md"`

### 3. README.md

**Required:** Yes  
**Format:** Markdown (CommonMark)

**Required Sections:**
1. **Plugin Name** (H1 heading)
2. **Description** - What the plugin does
3. **Features** - List of features
4. **Installation** - How to install
5. **Usage** - How to use the plugin
6. **Permissions** - What permissions are required and why
7. **Compatibility** - Editor version requirements
8. **Changelog** - Link to CHANGELOG.md
9. **License** - License information

**Optional Sections:**
- Screenshots
- Configuration
- Known Issues
- FAQ
- Support/Contact

### 4. CHANGELOG.md

**Required:** Yes  
**Format:** Markdown (Keep a Changelog format)

**Structure:**

```markdown
# Changelog

## [1.0.0] - 2026-01-16

### Added
- Initial release
- Feature 1
- Feature 2

### Changed
- Change 1

### Fixed
- Bug fix 1

### Security
- Security improvement 1
```

**Categories:**
- **Added:** New features
- **Changed:** Changes in existing functionality
- **Deprecated:** Soon-to-be removed features
- **Removed:** Removed features
- **Fixed:** Bug fixes
- **Security:** Security fixes

### 5. LICENSE (optional)

**Required:** No (defaults to MIT)  
**Format:** Plain text license file

If provided, must be a valid open-source license (MIT, GPL, Apache, BSD, etc.)

### 6. checksum.sha256

**Required:** Yes  
**Format:** Plain text SHA256 hash  
**Generated automatically** by validation script

```
a1b2c3d4e5f6789...
```

---

## Code Requirements

### Lua Version

Plugins execute in **Lua 5.1** environment with gopher-lua VM.

### Allowed Standard Library Modules

- **table:** Table manipulation functions
- **string:** String manipulation functions
- **math:** Mathematical functions
- **utf8:** UTF-8 string handling

### Blocked Modules/Functions

The following are **forbidden** and will cause validation failure:

- **io:** File I/O operations
- **os:** Operating system interaction
- **debug:** Debug library
- **package:** Module loading
- **require:** Dynamic loading
- **load/loadfile/loadstring:** Code execution
- **dofile:** File execution
- **setfenv/getfenv:** Environment manipulation
- **module:** Module definition

### Main Function Signature

```lua
function main(api)
    -- api: Plugin API object (see API Contract)
    -- Returns: boolean (true = success, false = failure)
end
```

### Error Handling

Plugins should use `pcall` for error handling:

```lua
local success, result = pcall(function()
    -- Plugin logic here
    return true
end)

if not success then
    api:Log("error", "Plugin error: " .. tostring(result))
    return false
end

return result
```

### Performance Requirements

- **Execution Time:** Maximum 30 seconds
- **Memory Usage:** Maximum 50MB
- **API Calls:** No enforced limit, but should be reasonable

### Code Style

While not strictly enforced, follow these guidelines:
- Use 2 or 4 spaces for indentation (consistent)
- Comment complex logic
- Use descriptive variable names
- Group related code with comments
- Include usage examples in comments

---

## Permissions

Plugins must declare all required permissions in both:
1. Metadata comments in plugin.lua
2. permissions array in metadata.json

### Permission Descriptions

#### read_save
- **Purpose:** Read save data
- **Grants Access To:**
  - Character data (stats, equipment, spells)
  - Inventory items
  - Party composition
  - Equipment items
  - Map/world state
- **Risk Level:** Low
- **Required For:** Any plugin that needs to read save data

#### write_save
- **Purpose:** Modify save data
- **Grants Access To:**
  - Character modifications
  - Inventory changes
  - Party changes
  - Equipment changes
- **Risk Level:** Medium
- **Required For:** Plugins that modify save data
- **Implies:** read_save (automatically granted)

#### ui_display
- **Purpose:** Show dialogs and prompts
- **Grants Access To:**
  - ShowDialog() - Display information
  - ShowConfirm() - Ask yes/no questions
  - ShowInput() - Request user input
- **Risk Level:** Low
- **Required For:** Any plugin with user interface

#### events
- **Purpose:** Register event hooks
- **Grants Access To:**
  - RegisterHook() - Register for events
  - FireEvent() - Trigger custom events
- **Risk Level:** Low
- **Required For:** Plugins that respond to editor events

### Permission Validation

During plugin execution, the API enforces permissions:

```lua
-- This will fail if plugin doesn't have write_save permission
api:SetCharacter(char)  -- Throws error if permission denied
```

---

## API Contract

### Available API Methods

Plugins interact with the editor through the Plugin API object passed to `main(api)`.

#### Character Methods

```lua
-- Get character by index (0-15)
local char = api:GetCharacter(index)

-- Set character data
api:SetCharacter(character)

-- Find characters matching predicate
local chars = api:FindCharacter(function(char)
    return char.Level > 50
end)
```

#### Inventory Methods

```lua
-- Get inventory
local inventory = api:GetInventory()

-- Set inventory
api:SetInventory(inventory)
```

#### Party Methods

```lua
-- Get party composition
local party = api:GetParty()

-- Set party composition
api:SetParty(party)
```

#### Equipment Methods

```lua
-- Get equipment
local equipment = api:GetEquipment()

-- Set equipment
api:SetEquipment(equipment)
```

#### UI Methods

```lua
-- Show information dialog
api:ShowDialog(title, message)

-- Show confirmation dialog
local confirmed = api:ShowConfirm(title, message)

-- Show input dialog
local value = api:ShowInput(title, prompt, defaultValue)
```

#### Logging Methods

```lua
-- Log message with level
api:Log(level, message)
-- Levels: "debug", "info", "warn", "error"
```

#### Settings Methods

```lua
-- Get plugin setting
local value = api:GetSetting(key)

-- Set plugin setting
api:SetSetting(key, value)
```

#### Permission Methods

```lua
-- Check if plugin has permission
local hasPermission = api:HasPermission("write_save")
```

#### Event Methods

```lua
-- Register event hook
api:RegisterHook(hookType, callback)

-- Fire custom event
api:FireEvent(eventName, data)
```

### Data Structures

#### Character Object

```lua
{
    ID = 0,              -- Character ID (0-15)
    Name = "Terra",      -- Character name
    Enabled = true,      -- Is character enabled
    Level = 50,          -- Character level
    Experience = 123456, -- Experience points
    CurrentHP = 850,     -- Current HP
    MaxHP = 1000,        -- Maximum HP
    CurrentMP = 180,     -- Current MP
    MaxMP = 200,         -- Maximum MP
    Vigor = 35,          -- Physical attack
    Stamina = 30,        -- Physical defense
    Speed = 40,          -- Agility/turn frequency
    Magic = 45,          -- Magical power
    Weapon = 100,        -- Weapon item ID
    Shield = 101,        -- Shield item ID
    Armor = 102,         -- Armor item ID
    Helmet = 103,        -- Helmet item ID
    Relic1 = 104,        -- Relic 1 item ID
    Relic2 = 105,        -- Relic 2 item ID
    Spells = {           -- Learned spells (spell ID → proficiency)
        [1] = 50,
        [2] = 100
    },
    Commands = {         -- Command list
        { Name = "Fight", Value = 0 },
        { Name = "Magic", Value = 17 }
    }
}
```

#### Inventory Object

```lua
{
    Items = {
        { ItemID = 1, Count = 10 },
        { ItemID = 2, Count = 5 }
    }
}
```

#### Party Object

```lua
{
    Members = { 0, 1, 2, 3 }  -- Character IDs (0-15)
}
```

Complete API documentation: [PHASE_4B_API_REFERENCE.md](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_API_REFERENCE.md)

---

## Security Requirements

### Sandboxing

All plugins run in a sandboxed Lua environment:
- No file system access
- No network access
- No OS command execution
- No arbitrary code loading
- Restricted standard library

### Forbidden Patterns

Plugins will fail validation if they contain:
- Attempts to access forbidden modules
- Attempts to bypass sandbox
- Obfuscated code
- Malicious patterns
- Code that attempts to access plugin system internals

### Safe Coding Practices

✅ **Do:**
- Validate all user input
- Handle errors gracefully
- Use appropriate permissions
- Document security considerations
- Test with malformed data

❌ **Don't:**
- Assume data is valid
- Ignore error conditions
- Request unnecessary permissions
- Hardcode sensitive data
- Attempt to bypass security

---

## Versioning

Plugins must follow [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR:** Incompatible API changes or breaking changes
- **MINOR:** New functionality, backward compatible
- **PATCH:** Bug fixes, backward compatible

Examples:
- `1.0.0` - Initial release
- `1.1.0` - Added new feature
- `1.1.1` - Fixed bug
- `2.0.0` - Breaking changes

### Version Updates

When updating a plugin:
1. Update version in metadata.json
2. Update version in plugin.lua metadata comments
3. Add entry to CHANGELOG.md
4. Regenerate checksum
5. Submit pull request

---

## Validation Rules

The validation script (`scripts/validate-plugin.py`) enforces these rules:

### Structure Validation

- ✅ Plugin directory exists
- ✅ Plugin ID matches directory name
- ✅ All required files present
- ✅ File sizes within limits

### Metadata Validation

- ✅ Valid JSON format
- ✅ Conforms to JSON schema
- ✅ All required fields present
- ✅ Field values match patterns/constraints
- ✅ Plugin ID is unique

### Code Validation

- ✅ Valid Lua syntax
- ✅ Defines main(api) function
- ✅ No forbidden modules/functions
- ✅ Metadata comments present
- ✅ Metadata matches metadata.json

### Documentation Validation

- ✅ README.md exists and readable
- ✅ CHANGELOG.md exists and readable
- ✅ Required sections present
- ✅ Markdown format valid

### Checksum Validation

- ✅ Checksum file exists
- ✅ SHA256 format correct
- ✅ Checksum matches plugin.lua

### Security Validation

- ✅ No forbidden patterns
- ✅ No obfuscation attempts
- ✅ Permission usage appropriate
- ✅ No security vulnerabilities

---

## Compliance

Plugins that don't comply with this specification will be:
- Rejected during CI validation
- Flagged during manual review
- Not published to the registry

---

## Updates to Specification

This specification may be updated. Plugin authors will be notified of:
- **Breaking Changes:** Require plugin updates
- **New Features:** Optional enhancements
- **Clarifications:** Documentation improvements

Version history: See [CHANGELOG.md](../CHANGELOG.md)

---

## Support

Questions about this specification:
- **GitHub Discussions:** Ask in plugin-registry discussions
- **Discord:** #plugin-development channel
- **Email:** plugins@ff6editor.dev

---

## References

- [Plugin API Reference](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_API_REFERENCE.md)
- [Plugin Development Guide](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_PLUGIN_GUIDE.md)
- [Plugin Examples](https://github.com/ff6-save-editor/editor/blob/main/PHASE_4B_PLUGIN_EXAMPLES.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [JSON Schema](schema/plugin-schema.json)
