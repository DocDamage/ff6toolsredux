# Equipment Optimizer Plugin

**Version:** 1.2.0  
**Category:** Utility  
**Author:** FF6 Plugin Team

## Overview

The **Equipment Optimizer** finds and compares the best gear loadouts for every character. v1.2.0 adds Phase 11 integrations: analytics-assisted optimization, dashboards, automation hooks, import/export, and loadout guide generation.

## Whats New
- **v1.2.0 (Phase 11 Tier 1):** Analytics-assisted optimization, dashboards, automation, import/export, loadout guides
- **v1.1.0 (Quick Win #3):** Equipment comparison dashboard with stat diffs and synergies

## Key Features

### Optimization
- `optimizeEquipment(char_id, goal)`  Core optimizer (offense/defense/balanced/magic)
- `optimizeEquipmentLoadout(char_id, goal)`  Analytics-enhanced optimization
- `predictEquipmentSynergies(loadout)`, `detectEquipmentOutliers(loadouts)`, `recommendSpecializedGear(scenario)`
- `analyzeResistancePatterns(enemy_profile)` for defensive gearing

### Comparison & Visualization
- `compareLoadouts(loadouts)`, `displayComparisonDashboard(comparison)`
- `generateEquipmentComparison(loadouts)`, `createLoadoutDashboard(loadouts)`
- `visualizeArmorCoverage(loadout)` gauge, `exportLoadoutGuide(loadouts, format, output_path)`

### Automation
- `autoEquipOptimal(char_id)` executes optimization + workflow
- `setupEquipmentRules(rule_set)`, `scheduleEquipmentReview(cron_expr)`, `triggerEquipmentAlert(message)`

### Import/Export & Sync
- `exportEquipmentTemplate(loadout, format, path)`, `importEquipmentTemplate(path, format)`
- `batchApplyLoadouts(loadouts)`, `syncEquipmentWithTeam(team_ids)`

## Quick Start
```lua
-- Optimize and compare
local current = getCurrentLoadout(0)
local offense = optimizeEquipmentLoadout(0, "offense")
local balanced = optimizeEquipmentLoadout(0, "balanced")

local comparison = compareLoadouts({current, offense, balanced})
displayComparisonDashboard(comparison)

-- Visual chart
local chart = generateEquipmentComparison({current, offense, balanced})

-- Automation
autoEquipOptimal(0)
scheduleEquipmentReview("0 9 * * *")

-- Export template
exportEquipmentTemplate(offense, "json", "terra_offense.json")
```

## Dependencies
- Phase 11: Advanced Analytics Engine, Data Visualization Suite, Import/Export Manager, Automation Framework
- Functions fail gracefully if dependencies are missing.

## Commands (console)
- Optimization: `optimizeEquipment`, `optimizeEquipmentLoadout`
- Comparison: `compareLoadouts`, `displayComparisonDashboard`, `generateEquipmentComparison`
- Automation: `autoEquipOptimal`, `setupEquipmentRules`, `scheduleEquipmentReview`
- Import/Export: `exportEquipmentTemplate`, `importEquipmentTemplate`

## Safety
- Display and analytics helpers are non-destructive
- Import/export writes only to provided paths
- Automation hooks are no-ops if Automation Framework is absent
