# Phase 4G: GitHub Registry & Example Plugins - Implementation Plan

**Date Created:** January 16, 2026  
**Phase Status:** In Progress  
**Estimated Duration:** 8-12 hours  
**Priority:** High

---

## Overview

Phase 4G establishes the official plugin registry infrastructure on GitHub and creates production-quality example plugins that demonstrate the full capabilities of the plugin system. This phase bridges the gap between the completed marketplace backend (Phase 4C), the enhanced Plugin API (Phases 4C+, 4D, 4E), and the Marketplace UI (Phase 4F).

---

## Goals

1. **GitHub Registry Setup**
   - Create official `ff6-save-editor/plugin-registry` repository
   - Design and implement `plugins.json` catalog schema
   - Establish plugin submission and review workflow
   - Set up automated validation and testing

2. **Example Plugin Development**
   - Create 3-5 production-quality example plugins
   - Demonstrate all major API capabilities
   - Provide templates for community plugin developers
   - Include comprehensive documentation for each plugin

3. **End-to-End Testing**
   - Verify marketplace discovery from live GitHub registry
   - Test plugin installation workflow end-to-end
   - Validate plugin execution and permissions
   - Confirm bidirectional sync between Plugin Manager and Marketplace

4. **Developer Documentation**
   - Plugin submission guidelines
   - Registry structure documentation
   - Best practices for plugin development
   - Testing and validation procedures

---

## Phase 4G Deliverables

### 1. GitHub Registry Structure

#### Repository: `ff6-save-editor/plugin-registry`

```
plugin-registry/
├── README.md                           # Registry overview and usage
├── CONTRIBUTING.md                     # Plugin submission guidelines
├── PLUGIN_SPECIFICATION.md             # Plugin format specification
├── plugins.json                        # Main plugin catalog (auto-generated)
├── schema/
│   └── plugin-schema.json              # JSON schema for plugin metadata
├── plugins/
│   ├── stats-display/
│   │   ├── plugin.lua                  # Plugin code
│   │   ├── metadata.json               # Plugin metadata
│   │   ├── README.md                   # Plugin documentation
│   │   ├── screenshot.png              # Plugin screenshot
│   │   └── checksum.sha256             # SHA256 checksum
│   ├── item-manager/
│   │   ├── plugin.lua
│   │   ├── metadata.json
│   │   ├── README.md
│   │   ├── screenshot.png
│   │   └── checksum.sha256
│   ├── party-optimizer/
│   │   ├── plugin.lua
│   │   ├── metadata.json
│   │   ├── README.md
│   │   ├── screenshot.png
│   │   └── checksum.sha256
│   └── ...
├── scripts/
│   ├── validate-plugin.py              # Plugin validation script
│   ├── generate-catalog.py             # Catalog generation script
│   └── update-checksums.py             # Checksum generation
└── .github/
    └── workflows/
        ├── validate-plugins.yml        # CI validation workflow
        └── generate-catalog.yml        # Auto-generate plugins.json
```

#### plugins.json Schema

```json
{
  "version": "1.0.0",
  "generated": "2026-01-16T12:00:00Z",
  "plugins": [
    {
      "id": "stats-display",
      "name": "Character Stats Display",
      "author": "FF6 Editor Team",
      "version": "1.0.0",
      "description": "Display comprehensive character stats including HP, MP, equipment, and spells",
      "category": "utility",
      "tags": ["stats", "display", "character", "analysis"],
      "permissions": ["read_save", "ui_display"],
      "downloadUrl": "https://raw.githubusercontent.com/ff6-save-editor/plugin-registry/main/plugins/stats-display/plugin.lua",
      "homepage": "https://github.com/ff6-save-editor/plugin-registry/tree/main/plugins/stats-display",
      "checksum": "sha256:a1b2c3d4e5f6...",
      "size": 4096,
      "rating": 4.8,
      "downloads": 1250,
      "created": "2026-01-16T12:00:00Z",
      "updated": "2026-01-16T12:00:00Z",
      "minEditorVersion": "3.4.0",
      "changelog": "https://github.com/ff6-save-editor/plugin-registry/blob/main/plugins/stats-display/CHANGELOG.md"
    }
  ]
}
```

### 2. Example Plugins

#### Plugin 1: Character Stats Display
**Category:** Utility  
**Permissions:** read_save, ui_display  
**Description:** Comprehensive character information display with formatted output

**Features:**
- Display all character stats (HP, MP, Level, Experience)
- Show all 6 equipment slots with item names
- List learned spells with proficiency percentages
- Show command assignments
- Format output in a user-friendly dialog

**Use Case:** Quick character overview without navigating multiple editor tabs

#### Plugin 2: Item Manager
**Category:** Automation  
**Permissions:** read_save, write_save, ui_display  
**Description:** Advanced inventory management with batch operations

**Features:**
- Add/remove items in bulk
- Set item quantities (e.g., "Max all consumables")
- Filter items by category
- Export/import inventory lists
- Duplicate rare items

**Use Case:** Quickly configure inventory for speedruns or testing

#### Plugin 3: Party Optimizer
**Category:** Analysis  
**Permissions:** read_save, write_save, ui_display  
**Description:** Automatically select and configure optimal party composition

**Features:**
- Analyze all characters' stats and equipment
- Recommend best party for different scenarios (magic-heavy, physical, balanced)
- Auto-equip best gear from inventory
- Calculate party-wide stat totals
- Suggest esper assignments for optimal growth

**Use Case:** Find the best party setup for different game scenarios

#### Plugin 4: Spell Coverage Analyzer (Optional)
**Category:** Analysis  
**Permissions:** read_save, ui_display  
**Description:** Analyze party-wide spell coverage and identify gaps

**Features:**
- List all learned spells across party
- Identify missing critical spells
- Show spell coverage by school (Black, White, Blue, etc.)
- Recommend characters to learn specific spells
- Calculate total spell diversity

**Use Case:** Ensure party has comprehensive spell coverage

#### Plugin 5: Batch Level Up (Optional)
**Category:** Automation  
**Permissions:** read_save, write_save, ui_display  
**Description:** Level up characters with optimal stat growth

**Features:**
- Level up to target level
- Calculate optimal esper assignments per level
- Apply stat bonuses automatically
- Show projected stat totals
- Undo support

**Use Case:** Quickly level characters with optimal stat distribution

### 3. Plugin Metadata Specification

Each plugin must include a `metadata.json` file:

```json
{
  "id": "unique-plugin-id",
  "name": "Human Readable Name",
  "version": "1.0.0",
  "author": "Author Name",
  "contact": "author@example.com",
  "description": "Detailed plugin description",
  "longDescription": "Extended description with features and use cases",
  "category": "utility|automation|analysis|enhancement",
  "tags": ["tag1", "tag2", "tag3"],
  "permissions": ["read_save", "write_save", "ui_display", "events"],
  "minEditorVersion": "3.4.0",
  "homepage": "https://github.com/...",
  "repository": "https://github.com/...",
  "license": "MIT",
  "created": "2026-01-16T12:00:00Z",
  "updated": "2026-01-16T12:00:00Z",
  "dependencies": [],
  "screenshots": ["screenshot1.png", "screenshot2.png"],
  "documentation": "README.md",
  "changelog": "CHANGELOG.md"
}
```

### 4. Plugin Submission Workflow

#### For Plugin Authors:

1. **Develop Plugin:**
   - Follow Plugin API documentation
   - Test thoroughly with example save files
   - Include comprehensive error handling
   - Add detailed inline comments

2. **Prepare Submission:**
   - Create plugin directory structure
   - Write `metadata.json`
   - Write README.md with usage instructions
   - Create CHANGELOG.md
   - Add screenshot(s)
   - Generate SHA256 checksum

3. **Submit to Registry:**
   - Fork `plugin-registry` repository
   - Add plugin to `plugins/` directory
   - Run validation script: `python scripts/validate-plugin.py plugins/your-plugin`
   - Create pull request with plugin submission template
   - Wait for automated CI checks to pass
   - Address review feedback

4. **Plugin Review:**
   - Automated validation (syntax, metadata, checksum)
   - Security review (no forbidden API calls)
   - Code quality review
   - Testing by maintainers
   - Approval and merge

5. **Publication:**
   - CI automatically regenerates `plugins.json`
   - Plugin appears in marketplace within 5 minutes
   - Author receives notification

### 5. Registry Automation

#### GitHub Actions Workflows:

**validate-plugins.yml:**
- Trigger: On pull request to main branch
- Steps:
  1. Checkout repository
  2. Validate plugin structure
  3. Check metadata.json schema
  4. Verify checksums
  5. Run Lua syntax check
  6. Test plugin loads without errors
  7. Check for forbidden API calls
  8. Comment PR with results

**generate-catalog.yml:**
- Trigger: On push to main branch
- Steps:
  1. Checkout repository
  2. Scan all plugins in `plugins/` directory
  3. Generate `plugins.json` catalog
  4. Commit and push updated catalog
  5. Trigger marketplace cache refresh

### 6. Testing Strategy

#### Unit Tests:
- Each plugin includes test cases
- Test data validation
- Test error handling
- Test permissions

#### Integration Tests:
- Load plugin from registry
- Install through marketplace UI
- Execute plugin in editor
- Verify save file modifications
- Test uninstall

#### End-to-End Tests:
1. Start with clean editor installation
2. Open marketplace browser
3. Search for test plugin
4. Install plugin
5. Run plugin
6. Verify results
7. Check plugin logs
8. Uninstall plugin
9. Verify clean removal

---

## Implementation Tasks

### Task 1: Repository Setup (2 hours)
- [ ] Create GitHub repository `ff6-save-editor/plugin-registry`
- [ ] Set up repository structure
- [ ] Write README.md with registry overview
- [ ] Write CONTRIBUTING.md with submission guidelines
- [ ] Create PLUGIN_SPECIFICATION.md
- [ ] Design JSON schema for plugin metadata
- [ ] Set up branch protection rules
- [ ] Configure issue templates

### Task 2: Automation Scripts (2 hours)
- [ ] Write `validate-plugin.py`
  - Validate directory structure
  - Validate metadata.json against schema
  - Check Lua syntax
  - Verify checksums
  - Detect forbidden patterns
- [ ] Write `generate-catalog.py`
  - Scan plugins directory
  - Aggregate metadata
  - Generate plugins.json
  - Update timestamp
- [ ] Write `update-checksums.py`
  - Generate SHA256 checksums
  - Update checksum files
- [ ] Set up GitHub Actions workflows

### Task 3: Example Plugin Development (3-4 hours)
- [ ] **Stats Display Plugin** (1 hour)
  - Implement character stats extraction
  - Format output in readable layout
  - Create dialog with tabbed interface
  - Write documentation
  - Create screenshot
- [ ] **Item Manager Plugin** (1 hour)
  - Implement inventory read/write
  - Add batch operations
  - Create item filter UI
  - Write documentation
  - Create screenshot
- [ ] **Party Optimizer Plugin** (1.5 hours)
  - Implement party analysis
  - Calculate optimal compositions
  - Auto-equip functionality
  - Write documentation
  - Create screenshot

### Task 4: Documentation (1.5 hours)
- [ ] Write plugin development guide
- [ ] Create plugin submission checklist
- [ ] Write testing guidelines
- [ ] Create plugin template with boilerplate
- [ ] Document API best practices
- [ ] Create troubleshooting guide

### Task 5: Testing & Validation (1.5 hours)
- [ ] Test plugin discovery from registry
- [ ] Test marketplace installation workflow
- [ ] Verify plugin execution
- [ ] Test permission enforcement
- [ ] Verify checksum validation
- [ ] Test auto-update detection
- [ ] End-to-end marketplace integration test

### Task 6: Documentation & Finalization (1 hour)
- [ ] Update FEATURE_ROADMAP_DETAILED.md
- [ ] Create PHASE_4G_COMPLETION_REPORT.md
- [ ] Update main README.md
- [ ] Create plugin developer quick start guide
- [ ] Record demo video (optional)

---

## Success Criteria

1. ✅ GitHub registry repository is live and accessible
2. ✅ plugins.json catalog is auto-generated and valid
3. ✅ At least 3 example plugins are published
4. ✅ All plugins have complete documentation
5. ✅ Automated validation workflow is functional
6. ✅ Marketplace UI successfully discovers plugins
7. ✅ Plugin installation works end-to-end
8. ✅ Plugin execution is verified
9. ✅ Plugin submission workflow is documented
10. ✅ Zero compilation errors after integration

---

## Risk Assessment

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| GitHub API rate limits | Medium | Low | Use caching, fallback to manual refresh |
| Plugin validation complexity | Medium | Medium | Start with basic validation, iterate |
| Community adoption | Low | Medium | Provide excellent examples and docs |
| Security vulnerabilities | High | Low | Comprehensive sandboxing already in place |
| Performance issues | Medium | Low | Limit plugin size, enforce timeouts |

---

## Dependencies

- Phase 4C: Marketplace Backend (Complete) ✅
- Phase 4C+: Enhanced Plugin API (Complete) ✅
- Phase 4D: Plugin Write Operations (Complete) ✅
- Phase 4E: Inventory & Party Operations (Complete) ✅
- Phase 4F: Marketplace UI (Complete) ✅
- GitHub account for registry hosting
- Python 3.8+ for automation scripts
- GitHub Actions for CI/CD

---

## Timeline

**Total Estimated Time:** 8-12 hours

- **Day 1 (4 hours):** Repository setup + automation scripts
- **Day 2 (4 hours):** Example plugin development
- **Day 3 (2-4 hours):** Testing, documentation, finalization

---

## Next Phase Preview: Phase 4H (CLI Tools)

After Phase 4G completion, the next focus will be:

1. **Command-Line Interface**
   - Headless save file operations
   - Batch processing scripts
   - CI/CD integration
   - Automated testing support

2. **Scripting Support**
   - Shell script examples
   - PowerShell module
   - Python wrapper library
   - API documentation

---

## Notes

- This phase focuses on creating production-quality infrastructure
- Example plugins should demonstrate best practices
- Documentation is critical for community adoption
- Registry automation reduces maintenance burden
- End-to-end testing ensures marketplace reliability

---

## References

- [PHASE_4C_COMPLETION_REPORT.md](PHASE_4C_COMPLETION_REPORT.md) - Marketplace backend
- [PHASE_4B_PLUGIN_GUIDE.md](PHASE_4B_PLUGIN_GUIDE.md) - Plugin development guide
- [PHASE_4B_API_REFERENCE.md](PHASE_4B_API_REFERENCE.md) - Complete API reference
- [PHASE_4F_COMPLETION_REPORT.md](PHASE_4F_COMPLETION_REPORT.md) - Marketplace UI integration
