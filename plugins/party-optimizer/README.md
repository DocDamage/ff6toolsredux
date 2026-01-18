# Party Optimizer Plugin

**Version:** 1.2.0  
**Category:** Utility  
**Author:** FF6 Plugin Team

## Overview
Party Optimizer recommends and applies optimal party compositions. v1.2.0 adds Phase 11 integrations: analytics-assisted recommendations, dashboards, automation hooks, import/export, backup, and cross-plugin sync.

## Whats New
- **v1.2.0 (Phase 11 Tier 1):** Analytics + visualization + automation + import/export + backup/sync
- **v1.1.0:** Equipment integration, growth prediction, esper optimization, boss strategies (see v1_1_upgrades.lua)

## Key Features
- **Optimization:** `optimizePartyComposition(party, goal)`, `recommendPartyForScenario(info)`
- **Automation:** `autoConfigureParty(party, info)`
- **Visualization:** `visualizePartySynergy(party)`, `generatePartyReport(party)`
- **Import/Export:** `exportPartyTemplate(party, format, path)`, `importPartyTemplate(path, format)`
- **Backup/Sync:** `snapshotParty(party, label)`, `restoreParty(id)`, `syncPartyData(targets)`

## Quick Start
```lua
local party = {
  {name = "Terra", level = 45, synergy = 88},
  {name = "Celes", level = 44, synergy = 85},
  {name = "Locke", level = 43, synergy = 80},
  {name = "Edgar", level = 42, synergy = 78}
}

local rec = optimizePartyComposition(party, "balanced")
visualizePartySynergy(party)
autoConfigureParty(party, {boss = "Kefka", difficulty = 90})
exportPartyTemplate(party, "json", "party_balanced.json")
snapshotParty(party, "pre-kefka")
```

## Dependencies
Phase 11 plugins (lazy, optional): Advanced Analytics Engine, Data Visualization Suite, Import/Export Manager, Automation Framework, Integration Hub, Backup & Restore System. Functions no-op with warnings when missing.

## Safety
Display/analytics helpers are non-destructive; import/export writes only to provided paths; automation is skipped if the Automation Framework is absent.
