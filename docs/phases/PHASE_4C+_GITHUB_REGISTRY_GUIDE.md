# GitHub Plugin Registry Setup Guide

**Date:** January 16, 2026  
**Status:** Phase 4C+ Setup Documentation

---

## Overview

The FF6 Save Editor Plugin Marketplace uses a GitHub-hosted registry to maintain the catalog of available plugins. This guide documents the setup process for the official registry.

---

## Repository Structure

### Main Registry Repository: `ff6-marketplace/registry`

```
ff6-marketplace/registry/
├── README.md                     # Registry overview and usage guide
├── plugins.json                  # Master plugin catalog
├── ratings.json                  # Community ratings (generated)
├── categories.json               # Plugin categories
│
├── SUBMISSION_GUIDE.md          # How to submit plugins
├── PLUGIN_TEMPLATE.md           # Template for plugin documentation
├── LICENSE                      # MIT license
│
├── api/
│   ├── v1/
│   │   ├── plugins.json         # API endpoint for plugin list
│   │   ├── plugins/{id}/        # Individual plugin endpoints
│   │   └── ratings/{id}/        # Rating endpoints
│   └── v2/                      # (Future version)
│
├── plugins/                     # Plugin storage
│   ├── speedrun-tracker/
│   │   ├── 1.0.0/
│   │   │   ├── plugin.lua       # Plugin source code
│   │   │   ├── metadata.json    # Plugin metadata
│   │   │   └── README.md        # Documentation
│   │   ├── 1.0.1/
│   │   └── 1.1.0/
│   │
│   ├── damage-calculator/
│   ├── stat-optimizer/
│   └── build-manager/
│
└── examples/                    # Example plugins
    ├── hello-world.lua
    └── template.lua
```

---

## Setup Steps

### 1. Create GitHub Repository

```bash
# Repository Details
Name: registry
Owner: ff6-marketplace
Description: Official plugin registry for FF6 Save Editor
Visibility: Public
License: MIT
```

**Initial Files:**
- README.md - Registry overview
- LICENSE - MIT license
- .gitignore - Ignore node_modules, build artifacts

### 2. Create Master Catalog (plugins.json)

**File:** `plugins.json`

```json
{
  "version": "1.0.0",
  "lastUpdated": "2026-01-16T00:00:00Z",
  "plugins": [
    {
      "id": "speedrun-tracker",
      "name": "Speedrun Tracker",
      "author": "FF6 Community",
      "description": "Track speedrun progress with real-time statistics and route guidance.",
      "version": "1.0.0",
      "releaseDate": "2026-01-16T00:00:00Z",
      "category": "Speedrun",
      "tags": ["speedrun", "tracker", "analytics"],
      "downloadCount": 125,
      "rating": 4.8,
      "ratingCount": 24,
      "homepageURL": "https://github.com/ff6-marketplace/speedrun-tracker",
      "downloadURL": "https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins/speedrun-tracker/1.0.0/plugin.lua",
      "checksumSHA256": "abc123...",
      "minEditorVersion": "4.1.0",
      "maxEditorVersion": "5.0.0",
      "permissions": ["read_save", "ui_display", "events"],
      "versions": [
        {
          "version": "1.0.0",
          "releaseDate": "2026-01-16T00:00:00Z",
          "downloadURL": "https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins/speedrun-tracker/1.0.0/plugin.lua",
          "checksumSHA256": "abc123...",
          "changelog": "Initial release"
        }
      ]
    },
    {
      "id": "damage-calculator",
      "name": "Damage Calculator",
      "author": "Math Enthusiast",
      "description": "Calculate precise damage values for all character abilities and spells.",
      "version": "1.2.0",
      "category": "Editor Tools",
      "tags": ["calculator", "damage", "analysis"],
      "downloadCount": 87,
      "rating": 4.6,
      "ratingCount": 18,
      "downloadURL": "https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins/damage-calculator/1.2.0/plugin.lua",
      "checksumSHA256": "def456...",
      "minEditorVersion": "4.1.0"
    }
  ]
}
```

### 3. Create Categories File (categories.json)

**File:** `categories.json`

```json
{
  "categories": [
    {
      "id": "editor-tools",
      "name": "Editor Tools",
      "description": "Tools and utilities for editing save data",
      "icon": "🛠️",
      "pluginCount": 12
    },
    {
      "id": "speedrun",
      "name": "Speedrun",
      "description": "Speedrun assistance and tracking",
      "icon": "⚡",
      "pluginCount": 8
    },
    {
      "id": "analytics",
      "name": "Analytics",
      "description": "Statistics and analysis tools",
      "icon": "📊",
      "pluginCount": 5
    },
    {
      "id": "automation",
      "name": "Automation",
      "description": "Automation and scripting helpers",
      "icon": "🤖",
      "pluginCount": 6
    },
    {
      "id": "utilities",
      "name": "Utilities",
      "description": "General utilities and helpers",
      "icon": "⚙️",
      "pluginCount": 10
    }
  ]
}
```

### 4. Plugin Submission Structure

Each plugin should follow this structure:

```
plugins/{plugin-id}/{version}/
├── plugin.lua          # Main plugin code
├── metadata.json       # Plugin metadata
└── README.md          # Plugin documentation
```

**Example: speedrun-tracker/1.0.0/metadata.json**

```json
{
  "id": "speedrun-tracker",
  "name": "Speedrun Tracker",
  "version": "1.0.0",
  "author": "FF6 Community",
  "description": "Track speedrun progress with real-time statistics and route guidance.",
  "category": "Speedrun",
  "tags": ["speedrun", "tracker", "analytics"],
  "permissions": ["read_save", "ui_display", "events"],
  "minEditorVersion": "4.1.0",
  "maxEditorVersion": "5.0.0",
  "changelog": "Initial release - Basic speedrun tracking with route suggestions",
  "maintainer": {
    "name": "FF6 Community",
    "email": "community@ff6-marketplace.local",
    "github": "ff6-community"
  },
  "dependencies": [],
  "keywords": ["speedrun", "tracker", "speedrunning"]
}
```

### 5. API Endpoints

The registry provides REST API endpoints for plugin discovery:

#### List All Plugins
```
GET /api/v1/plugins.json
```

Returns: Array of all plugins with metadata

#### Get Specific Plugin
```
GET /api/v1/plugins/{plugin-id}/metadata.json
```

#### Download Plugin
```
GET /api/v1/plugins/{plugin-id}/{version}/plugin.lua
```

#### Get Ratings
```
GET /api/v1/ratings/{plugin-id}/
```

---

## Initial Example Plugins

### 1. Speedrun Tracker

**File:** `plugins/speedrun-tracker/1.0.0/plugin.lua`

```lua
-- Speedrun Tracker Plugin
-- Tracks speedrun progress and provides route guidance

local Plugin = {}

function Plugin.OnLoad()
    print("Speedrun Tracker v1.0.0 loaded")
end

function Plugin.OnUnload()
    print("Speedrun Tracker unloaded")
end

function Plugin.GetStats()
    local stats = {
        playtime = 0,
        bosses_defeated = 0,
        treasures_found = 0
    }
    return stats
end

return Plugin
```

### 2. Damage Calculator

**File:** `plugins/damage-calculator/1.2.0/plugin.lua`

```lua
-- Damage Calculator Plugin
-- Calculate precise damage for all abilities

local Plugin = {}

function Plugin.CalculateDamage(attackerLevel, attackerStat, weaponPower, defenderDefense)
    local baseDamage = ((attackerLevel + attackerStat) / 2 + weaponPower) * 2
    local reducedDamage = baseDamage * (200 - defenderDefense) / 200
    return math.max(1, math.floor(reducedDamage))
end

function Plugin.OnLoad()
    print("Damage Calculator v1.2.0 loaded")
end

return Plugin
```

### 3. Stat Optimizer

**File:** `plugins/stat-optimizer/1.0.0/plugin.lua`

```lua
-- Stat Optimizer Plugin
-- Suggests optimal stat distributions

local Plugin = {}

function Plugin.SuggestBuild(character_class)
    local builds = {
        warrior = {hp=999, str=99, def=90, spd=50},
        mage = {hp=400, mag=99, spr=90, spd=80},
        hybrid = {hp=700, str=80, mag=80, def=75}
    }
    return builds[character_class] or builds.hybrid
end

function Plugin.OnLoad()
    print("Stat Optimizer loaded")
end

return Plugin
```

---

## Submission Process

### For Plugin Authors:

1. **Prepare your plugin** following the template structure
2. **Fork the registry repository** (optional, or submit via PR)
3. **Create plugin directory** following naming conventions
4. **Include metadata.json** with complete information
5. **Submit pull request** with:
   - Plugin code and metadata
   - Complete documentation
   - Changelog entry
   - Any required dependencies listed
6. **Await review** from maintainers
7. **Incorporate feedback** if requested
8. **Merge and release** to official registry

### Registry Review Criteria:

- ✅ Plugin follows code standards
- ✅ Metadata is complete and accurate
- ✅ Documentation is clear
- ✅ No malicious or harmful code
- ✅ Respects permission model
- ✅ Stable and tested
- ✅ Licensed appropriately

---

## Configuration URLs

Add these to your editor's marketplace settings:

**Primary Registry:**
```
https://raw.githubusercontent.com/ff6-marketplace/registry/main
```

**API Endpoints:**
- Plugins List: `https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins.json`
- Individual Plugin: `https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins/{id}/metadata.json`
- Plugin Code: `https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins/{id}/{version}/plugin.lua`
- Categories: `https://raw.githubusercontent.com/ff6-marketplace/registry/main/categories.json`

---

## Maintenance

### Registry Updates:

- **Daily:** Monitor new submissions and issues
- **Weekly:** Review plugin ratings and feedback
- **Monthly:** Update featured plugins and categories
- **Quarterly:** Archive old versions, publish stats

### Backup Strategy:

- GitHub automatically backs up all content
- Regular exports to local storage
- Version history available via git log
- Rollback capability for all changes

---

## Future Enhancements

1. **Plugin Signing:** Cryptographically sign plugins
2. **Automated Testing:** CI/CD pipeline for submissions
3. **Plugin Stats:** Downloads, ratings, usage analytics
4. **Beta Channel:** Pre-release plugins for testing
5. **Dependency Resolution:** Automatic dependency installation
6. **Auto-Updates:** Seamless plugin update delivery

---

## Support & Contact

**Issues:** Open on GitHub  
**Discussions:** GitHub Discussions forum  
**Email:** marketplace@ff6-community.local  
**Discord:** Community Discord server (link)

---

**Registry Setup Complete** ✅

The GitHub plugin registry is ready to accept community submissions and serve plugins to all FF6 Save Editor installations.
