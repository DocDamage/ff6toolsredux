# Phase 1 Remaining Work - UI Integration (40%)

## Quick Reference

| Component | Lines | Status | Priority | Duration |
|-----------|-------|--------|----------|----------|
| Backup Manager Dialog | 250 | Not Started | HIGH | 2 days |
| Validation UI Panel | 200 | Not Started | HIGH | 2 days |
| Theme Switcher | 100 | Not Started | MEDIUM | 1 day |
| Undo/Redo Keyboard Shortcuts | 150 | Not Started | MEDIUM | 1 day |
| Integration Testing | N/A | Not Started | HIGH | 2 days |
| **TOTAL** | **700** | | | **~8 days** |

## Detailed Component Breakdown

### 1. Backup Manager Dialog (250 lines)
**File**: `ui/forms/dialogs/backup_manager.go`
**Timeline**: 2 days

#### Features to Implement:
```go
type BackupManagerDialog struct {
    window *fyne.Window
    manager *backup.Manager
    selectedBackup *models.BackupMetadata
    backupList fyne.CanvasObject
    table *widget.Table
}
```

#### UI Elements:
- **Backup List Table**:
  - Columns: Timestamp, Size, Description, Actions
  - Double-click to restore
  - Right-click context menu

- **Action Buttons**:
  - Restore Selected (enabled when backup selected)
  - Delete Selected (with confirmation)
  - Create Manual Backup (with description input)
  - Refresh List

- **Backup Details Panel**:
  - Show hash for verification
  - Display full timestamp
  - Show file size

#### Code Structure:
```go
// Dialog creation
func NewBackupManagerDialog(window *fyne.Window, manager *backup.Manager) *BackupManagerDialog

// Event handlers
func (d *BackupManagerDialog) OnRestoreClicked()
func (d *BackupManagerDialog) OnDeleteClicked()
func (d *BackupManagerDialog) OnCreateClicked()

// UI builders
func (d *BackupManagerDialog) buildBackupTable() fyne.CanvasObject
func (d *BackupManagerDialog) buildActionButtons() fyne.CanvasObject
func (d *BackupManagerDialog) buildDetailsPanel() fyne.CanvasObject

// Data refresh
func (d *BackupManagerDialog) RefreshBackupList()
func (d *BackupManagerDialog) ShowRestoreConfirmation(backup *models.BackupMetadata)
```

#### Integration Points:
- Hook into file menu: "Manage Backups..."
- Call `manager.RestoreBackup()` on restore
- Call `manager.DeleteBackup()` on delete
- Emit event to reload save file after restore

---

### 2. Validation UI Panel (200 lines)
**File**: `ui/forms/validation_panel.go`
**Timeline**: 2 days

#### Features to Implement:
```go
type ValidationPanel struct {
    validator *validation.Validator
    currentData *pr.PR
    resultsDisplay fyne.CanvasObject
    autoFixButton *widget.Button
    strictnessSelect *widget.Select
}
```

#### UI Elements:
- **Validation Mode Selector**:
  - Dropdown: Strict / Normal / Lenient
  - Description of each mode

- **Results Display**:
  - **Errors** section (if any)
  - **Warnings** section (if any)
  - **Info** section (if any)
  - Color-coded by severity

- **Issue List**:
  - Rule name
  - Message
  - Target (character, field)
  - Fixable indicator
  - Quick-fix buttons for fixable issues

- **Summary**:
  - "X errors, Y warnings, Z info"
  - Validation status badge
  - Auto-fix all button

#### Code Structure:
```go
// Panel creation
func NewValidationPanel(validator *validation.Validator) *ValidationPanel

// Validation execution
func (p *ValidationPanel) ValidateSaveData(data *pr.PR)
func (p *ValidationPanel) ClearResults()

// UI builders
func (p *ValidationPanel) buildModeSelector() fyne.CanvasObject
func (p *ValidationPanel) buildResultsDisplay(result models.ValidationResult) fyne.CanvasObject
func (p *ValidationPanel) buildSummary(result models.ValidationResult) fyne.CanvasObject
func (p *ValidationPanel) buildIssueList(issues []models.ValidationIssue) fyne.CanvasObject

// Event handlers
func (p *ValidationPanel) OnStrictnessChanged(mode string)
func (p *ValidationPanel) OnAutoFixAll()
func (p *ValidationPanel) OnFixIssue(issue *models.ValidationIssue)
```

#### Integration Points:
- Add as collapsible panel in main window
- Call on file load: `panel.ValidateSaveData(loadedData)`
- Call on pre-save: Warn if errors in Strict mode
- Display in status bar: "Validation: OK / 3 warnings / 1 error"

---

### 3. Theme Switcher (100 lines)
**File**: `ui/theme_switcher.go`
**Timeline**: 1 day

#### Features to Implement:
```go
type ThemeSwitcher struct {
    currentTheme *theme.Theme
    preferences fyne.Preferences
}
```

#### UI Elements:
- **Menu Item**: "View" → "Theme" → Toggle
- **Keyboard Shortcut**: Ctrl+T
- **Status Bar Indicator**: Show current theme

#### Code Structure:
```go
// Switcher creation
func NewThemeSwitcher(prefs fyne.Preferences) *ThemeSwitcher

// Theme management
func (ts *ThemeSwitcher) ToggleTheme()
func (ts *ThemeSwitcher) SetTheme(t *theme.Theme)
func (ts *ThemeSwitcher) LoadSavedTheme()
func (ts *ThemeSwitcher) SaveThemePreference()

// Event handler
func (ts *ThemeSwitcher) OnToggleRequested()
```

#### Integration Points:
- Load saved preference on startup
- Apply theme to all open windows
- Persist preference to app settings
- Update menu item text to show current theme

---

### 4. Undo/Redo System Integration (150 lines)
**File**: `ui/undo_redo_integration.go`
**Timeline**: 1 day

#### Features to Implement:
```go
type UndoRedoController struct {
    undoStack *state.UndoStack
    menuItems struct {
        Undo *fyne.MenuItem
        Redo *fyne.MenuItem
    }
    statusLabel *widget.Label
}
```

#### Menu Items:
- "Edit" → "Undo" (Ctrl+Z)
- "Edit" → "Redo" (Ctrl+Y)
- Both show action preview in tooltip
- Both show enabled/disabled state

#### Code Structure:
```go
// Controller creation
func NewUndoRedoController(undoStack *state.UndoStack) *UndoRedoController

// Menu setup
func (uc *UndoRedoController) BuildMenuItems() (*fyne.MenuItem, *fyne.MenuItem)

// Actions
func (uc *UndoRedoController) ExecuteUndo()
func (uc *UndoRedoController) ExecuteRedo()

// State updates
func (uc *UndoRedoController) OnChangeRecorded()
func (uc *UndoRedoController) UpdateMenuState()
func (uc *UndoRedoController) UpdateStatusLabel()

// Keyboard shortcuts
func (uc *UndoRedoController) HandleKeyboardShortcut(key *fyne.KeyEvent)
```

#### Integration Points:
- Record changes from all editors
- Hook into character/equipment/inventory editors
- Disable redo stack on new changes
- Show preview on menu hover
- Visual feedback for undo/redo state

---

## Integration Testing Checklist

### Phase 1 Complete When:
- [ ] All 4 UI components built and functional
- [ ] Backup manager successfully creates and restores
- [ ] Validation catches and fixes issues
- [ ] Theme switching works without restart
- [ ] Undo/redo works across all editors
- [ ] Keyboard shortcuts respond correctly
- [ ] No memory leaks or goroutine leaks
- [ ] Comprehensive testing document completed
- [ ] Zero compile errors
- [ ] No runtime panics

### Test Scenarios:

#### Backup Testing:
1. Create backup manually
2. Modify save file
3. Restore from backup
4. Verify file matches backup
5. Delete backup
6. Auto-cleanup when > max

#### Validation Testing:
1. Create character with invalid level
2. Run validation - should flag error
3. Enable auto-fix - should correct
4. Change to strict mode - should show warning as error
5. Create inventory overflow
6. Validation catches it

#### Theme Testing:
1. Start app - load light theme (default)
2. Toggle to dark theme
3. Verify all colors applied
4. Restart app - verify dark theme persists
5. Toggle back to light
6. Verify all colors applied correctly

#### Undo/Redo Testing:
1. Edit character level
2. Press Ctrl+Z - should revert
3. Press Ctrl+Y - should restore
4. Batch edit equipment
5. Single Ctrl+Z undoes entire batch
6. Rapid edits work correctly
7. Redo stack clears on new change

---

## Code Style Guidelines

For consistency with existing codebase:

```go
// Error handling
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// Naming: exported functions are PascalCase
func (d *Dialog) ShowDialog() {}

// Private functions are camelCase
func (d *Dialog) updateDisplay() {}

// Receiver variables: short names (d, p, c for dialog/panel/controller)
// Use full names for clarity in package-level functions

// Comments: Before exported functions
// ShowDialog displays the backup manager dialog
func (d *Dialog) ShowDialog() {}
```

---

## Development Priority

**Week 1-2 (Days 1-8):**
1. Backup Manager Dialog (Days 1-2)
2. Validation UI Panel (Days 3-4)
3. Theme Switcher (Days 5)
4. Undo/Redo Integration (Days 6-7)
5. Testing & Fixes (Day 8)

**Week 2-3 (Days 9-21):**
- Phase 2 features (UI improvements, advanced editors)
- Phase 3 features (Character optimizer, party system)
- Phase 4 features (Search, statistics, documentation)

---

## Success Metrics

Phase 1 will be considered complete when:
- ✅ All 10 files compile without warnings
- ✅ All 4 UI dialogs functional and tested
- ✅ Backup system working end-to-end
- ✅ Validation preventing invalid saves
- ✅ Themes applying correctly
- ✅ Undo/redo responsive and reliable
- ✅ No performance regressions
- ✅ Documentation complete

**Estimated Completion**: Week 2 (Days 8-14)

---

## Resources & Dependencies

- Existing code: `ui/window.go`, `ui/forms/`, existing editors
- Libraries: `fyne.io/fyne/v2`, `models/`, `io/backup/`, `io/validation/`, `ui/state/`, `ui/theme/`
- Test data: Use existing save file in repository
- References: Fyne documentation, Go best practices

---

**Document Version**: 1.0
**Created**: 2026-01-15
**Status**: Ready for Implementation
