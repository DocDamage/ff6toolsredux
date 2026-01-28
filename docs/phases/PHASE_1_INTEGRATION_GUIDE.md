# Phase 1 Integration Guide

This document explains how to integrate Phase 1 components into your existing application.

## Overview

Phase 1 provides four main systems ready for integration:

1. **Backup Manager** - Save file backup and recovery
2. **Validation Panel** - Real-time save file validation  
3. **Theme Switcher** - Light/dark theme support
4. **Undo/Redo Controller** - Change history management

## Integration Steps

### 1. Backup Manager Integration

**File**: `ui/forms/dialogs/backup_manager.go`

#### In your main window:

```go
import "ffvi_editor/io/backup"
import "ffvi_editor/ui/forms/dialogs"

// In your window initialization
backupManager := backup.NewManager() // Creates ~/.ffvi-editor/backups/
backupDialog := dialogs.NewBackupManagerDialog(window, backupManager)

// Add to menu
backupMenuItem := fyne.NewMenuItem(
    "Manage Backups...",
    func() { backupDialog.Show() },
)
```

#### Usage:

```go
// Create backup on file save
data := pr.New() // Your save data
err := backupManager.CreateBackup(data.GetRawBytes(), "Manual save")

// Restore backup
err := backupManager.RestoreBackup(backupID)

// List backups
backups := backupManager.ListBackups()
for _, backup := range backups {
    fmt.Printf("%s - %d bytes\n", backup.Timestamp, backup.FileSize)
}
```

### 2. Validation Panel Integration

**File**: `ui/forms/validation_panel.go`

#### In your main window:

```go
import "ffvi_editor/io/validation"
import "ffvi_editor/ui/forms"

// In your window initialization
validator := validation.NewValidator()
validationPanel := forms.NewValidationPanel(validator)

// Add to main container
mainContainer := container.NewVBox(
    validationPanel.BuildPanel(),
    // ... other UI elements
)

// Set callback for when issues are fixed
validationPanel.SetOnIssueFixed(func() {
    // Refresh display or update data
})
```

#### Usage:

```go
// Validate on file load
loadedData := pr.New() // Load from file
validationPanel.ValidateSaveData(loadedData)

// Check validation status before save
isValid, errors, warnings := validationPanel.GetValidationStatus()
if !isValid {
    dialog.ShowWarning("Validation Failed", 
        fmt.Sprintf("%d errors found", errors), window)
    return
}
```

### 3. Theme Switcher Integration

**File**: `ui/theme_switcher.go`

#### In your main window:

```go
import "ffvi_editor/ui"

// In your app initialization
themeSwitcher := ui.NewThemeSwitcher(app.Preferences())
themeSwitcher.RegisterWindow(mainWindow)

// Add to menu
themeMenuItem := themeSwitcher.BuildMenuItems()
// Add to View menu or custom location
```

#### Usage:

```go
// Toggle theme with shortcut or button
button := widget.NewButton("Toggle Theme", themeSwitcher.ToggleTheme)

// Get current theme
currentTheme := themeSwitcher.GetCurrentThemeName()
fmt.Println(currentTheme) // "Light" or "Dark"
```

### 4. Undo/Redo Integration

**File**: `ui/undo_redo_controller.go`

#### In your main window:

```go
import "ffvi_editor/ui/state"
import "ffvi_editor/ui"
import "ffvi_editor/models"

// In your window initialization
undoStack := state.NewUndoStack(100) // Max 100 undo levels
undoRedoController := ui.NewUndoRedoController(undoStack)

undoMenuItem, redoMenuItem := undoRedoController.BuildMenuItems()
// Add to Edit menu

// Add status label to window
statusBox := container.NewHBox(
    undoRedoController.GetStatusLabel(),
)

// Set callbacks
undoRedoController.SetOnUndo(func() {
    // Handle undo - rebuild UI with previous state
    refreshDisplay()
})

undoRedoController.SetOnRedo(func() {
    // Handle redo - rebuild UI with next state
    refreshDisplay()
})
```

#### Usage:

```go
// Record a change when user edits something
change := &models.Change{
    Target:    "character_1",
    FieldName: "Level",
    OldValue:  50,
    NewValue:  51,
}
undoRedoController.RecordChange(change)

// Batch multiple related changes
batchID := undoRedoController.StartBatch("Character Edit")
// ... make multiple changes, each calls RecordChange()
undoRedoController.EndBatch(batchID)

// User presses Ctrl+Z
undoRedoController.executeUndo()
// All changes in batch are undone together
```

## Complete Integration Example

```go
package ui

import (
    "ffvi_editor/io/backup"
    "ffvi_editor/io/validation"
    "ffvi_editor/io/pr"
    "ffvi_editor/ui/forms"
    "ffvi_editor/ui/forms/dialogs"
    "ffvi_editor/ui/state"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
)

type AppState struct {
    Window       fyne.Window
    Backup       *backup.Manager
    Validator    *validation.Validator
    UndoRedo     *UndoRedoController
    ThemeSwitcher *ThemeSwitcher
    CurrentData  *pr.PR
}

func NewAppState() *AppState {
    myApp := app.New()
    mainWindow := myApp.NewWindow()
    
    state := &AppState{
        Window:        mainWindow,
        Backup:        backup.NewManager(),
        Validator:     validation.NewValidator(),
        UndoRedo:      NewUndoRedoController(state.NewUndoStack(100)),
        ThemeSwitcher: NewThemeSwitcher(myApp.Preferences()),
        CurrentData:   pr.New(),
    }
    
    state.setupUI()
    return state
}

func (s *AppState) setupUI() {
    // Setup menus
    editMenu := fyne.NewMenu("Edit",
        s.UndoRedo.BuildMenuItems()[0], // Undo
        s.UndoRedo.BuildMenuItems()[1], // Redo
    )
    
    viewMenu := fyne.NewMenu("View",
        s.ThemeSwitcher.BuildMenuItems(),
    )
    
    fileMenu := fyne.NewMenu("File",
        fyne.NewMenuItem("Manage Backups...", func() {
            backupDialog := dialogs.NewBackupManagerDialog(
                s.Window, s.Backup)
            backupDialog.Show()
        }),
    )
    
    mainMenu := fyne.NewMainMenu(fileMenu, editMenu, viewMenu)
    s.Window.SetMainMenu(mainMenu)
    
    // Setup validation panel
    validationPanel := forms.NewValidationPanel(s.Validator)
    validationPanel.SetOnIssueFixed(func() {
        // Refresh data after fix
    })
    
    // Setup layout
    content := container.NewVBox(
        validationPanel.BuildPanel(),
        // Add other UI elements here
    )
    
    s.Window.SetContent(content)
}

func (s *AppState) SaveFile(data *pr.PR) error {
    // Validate first
    s.Validator.Validate(data)
    
    // Create backup
    err := s.Backup.CreateBackup(data.GetRawBytes(), "Auto-save")
    if err != nil {
        return err
    }
    
    // Save file...
    // Clear undo history on successful save
    s.UndoRedo.Clear()
    
    return nil
}
```

## API Reference

### Backup Manager
```go
manager := backup.NewManager()
manager.CreateBackup(data []byte, description string) error
manager.RestoreBackup(id string) error
manager.DeleteBackup(id string) error
manager.ListBackups() []models.BackupMetadata
manager.SetMaxBackups(count int)
```

### Validator
```go
validator := validation.NewValidator()
result := validator.Validate(data *pr.PR) models.ValidationResult
fixed, err := validator.AutoFixIssues(data *pr.PR) (int, error)
validator.SetConfig(config models.ValidationConfig)
```

### Undo/Redo
```go
stack := state.NewUndoStack(maxDepth int)
stack.RecordChange(change models.Change)
batchID, err := stack.StartBatch(name string) (string, error)
err := stack.EndBatch(name string) error
changes, name, err := stack.PopUndo() ([]models.Change, string, error)
changes, name, err := stack.PopRedo() ([]models.Change, string, error)
```

### Theme
```go
switcher := NewThemeSwitcher(prefs fyne.Preferences)
switcher.RegisterWindow(window fyne.Window)
switcher.ToggleTheme()
switcher.SetTheme(theme *theme.Theme)
switcher.GetCurrentThemeName() string
```

## Data Structures

### Change
```go
type Change struct {
    ID        string        // Unique change ID
    Timestamp int64         // Unix nanoseconds
    Target    string        // What was changed
    FieldName string        // Which field
    OldValue  interface{}   // Previous value
    NewValue  interface{}   // New value
    Batch     bool          // Part of batch?
    BatchID   string        // Batch identifier
    BatchName string        // Batch description
}
```

### ValidationResult
```go
type ValidationResult struct {
    Valid    bool
    Errors   []ValidationIssue
    Warnings []ValidationIssue
    Infomsgs []ValidationIssue
}
```

## Testing the Integration

```go
// Test backup functionality
func TestBackupWorkflow(t *testing.T) {
    manager := backup.NewManager()
    originalData := []byte("save file data")
    
    // Create backup
    err := manager.CreateBackup(originalData, "test")
    assert.NoError(t, err)
    
    // List backups
    backups := manager.ListBackups()
    assert.Equal(t, 1, len(backups))
    
    // Restore backup
    err = manager.RestoreBackup(backups[0].ID)
    assert.NoError(t, err)
}

// Test undo/redo
func TestUndoRedo(t *testing.T) {
    stack := state.NewUndoStack(10)
    
    change := &models.Change{Target: "test", FieldName: "value"}
    stack.RecordChange(*change)
    
    assert.True(t, stack.CanUndo())
    assert.False(t, stack.CanRedo())
    
    stack.PopUndo()
    assert.False(t, stack.CanUndo())
    assert.True(t, stack.CanRedo())
}
```

## Common Patterns

### Auto-save with Backup
```go
func (s *AppState) OnFileModified() {
    // Record for undo
    s.UndoRedo.RecordChange(change)
    
    // Auto-backup every 5 minutes
    if time.Since(lastBackup) > 5*time.Minute {
        s.Backup.CreateBackup(s.CurrentData.GetRawBytes(), "Auto-save")
        lastBackup = time.Now()
    }
}
```

### Batch Edits
```go
func (s *AppState) EditCharacter(id int, updates map[string]interface{}) {
    batchID, _ := s.UndoRedo.StartBatch("Edit Character")
    
    for field, value := range updates {
        change := &models.Change{
            Target:    fmt.Sprintf("character_%d", id),
            FieldName: field,
            OldValue:  getOldValue(id, field),
            NewValue:  value,
        }
        s.UndoRedo.RecordChange(change)
    }
    
    s.UndoRedo.EndBatch(batchID)
}
```

### Pre-save Validation
```go
func (s *AppState) SaveFile() error {
    result := s.Validator.Validate(s.CurrentData)
    
    if !result.Valid {
        if len(result.Errors) > 0 {
            return fmt.Errorf("save failed: %d validation errors", 
                len(result.Errors))
        }
    }
    
    // Save...
    return nil
}
```

---

For questions about integration, refer to the sample code in `ui/window.go` and the existing editor implementations in `ui/forms/editors/`.
