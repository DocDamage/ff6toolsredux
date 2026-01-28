package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/models"
	"ffvi_editor/ui/state"
)

// UndoRedoController manages undo/redo operations through UI
type UndoRedoController struct {
	undoStack    *state.UndoStack
	undoMenuItem *fyne.MenuItem
	redoMenuItem *fyne.MenuItem
	statusLabel  *widget.Label
	onUndo       func() // Callback when undo completes
	onRedo       func() // Callback when redo completes
}

// NewUndoRedoController creates a new undo/redo controller
func NewUndoRedoController(undoStack *state.UndoStack) *UndoRedoController {
	urc := &UndoRedoController{
		undoStack:   undoStack,
		statusLabel: widget.NewLabel(""),
	}

	// Build menu items
	urc.undoMenuItem = fyne.NewMenuItem(
		"Undo",
		urc.executeUndo,
	)
	urc.undoMenuItem.Disabled = true

	urc.redoMenuItem = fyne.NewMenuItem(
		"Redo",
		urc.executeRedo,
	)
	urc.redoMenuItem.Disabled = true

	return urc
}

// BuildMenuItems returns the undo/redo menu items
func (urc *UndoRedoController) BuildMenuItems() (*fyne.MenuItem, *fyne.MenuItem) {
	return urc.undoMenuItem, urc.redoMenuItem
}

// GetStatusLabel returns the status label widget
func (urc *UndoRedoController) GetStatusLabel() *widget.Label {
	return urc.statusLabel
}

// RecordChange records a change in the undo stack
func (urc *UndoRedoController) RecordChange(change *models.Change) {
	urc.undoStack.RecordChange(*change)
	urc.updateMenuState()
}

// StartBatch starts a batch of changes
func (urc *UndoRedoController) StartBatch(name string) string {
	id, _ := urc.undoStack.StartBatch(name)
	return id
}

// EndBatch ends the current batch
func (urc *UndoRedoController) EndBatch(name string) {
	_ = urc.undoStack.EndBatch(name)
	urc.updateMenuState()
}

// CancelBatch cancels the current batch without saving
func (urc *UndoRedoController) CancelBatch() {
	urc.undoStack.CancelBatch()
}

// executeUndo performs an undo operation
func (urc *UndoRedoController) executeUndo() {
	if !urc.undoStack.CanUndo() {
		return
	}

	changes, _, _ := urc.undoStack.PopUndo()
	if len(changes) == 0 {
		return
	}

	// Call callback
	if urc.onUndo != nil {
		urc.onUndo()
	}

	urc.updateMenuState()
}

// executeRedo performs a redo operation
func (urc *UndoRedoController) executeRedo() {
	if !urc.undoStack.CanRedo() {
		return
	}

	changes, _, _ := urc.undoStack.PopRedo()
	if len(changes) == 0 {
		return
	}

	// Call callback
	if urc.onRedo != nil {
		urc.onRedo()
	}

	urc.updateMenuState()
}

// updateMenuState updates the enabled/disabled state of menu items
func (urc *UndoRedoController) updateMenuState() {
	// Update undo menu item
	if urc.undoStack.CanUndo() {
		urc.undoMenuItem.Disabled = false
		undoCount := urc.undoStack.GetUndoCount()
		preview := urc.undoStack.GetUndoPreview()
		urc.undoMenuItem.Label = fmt.Sprintf("Undo (%d) - %s", undoCount, preview)
	} else {
		urc.undoMenuItem.Disabled = true
		urc.undoMenuItem.Label = "Undo"
	}

	// Update redo menu item
	if urc.undoStack.CanRedo() {
		urc.redoMenuItem.Disabled = false
		redoCount := urc.undoStack.GetRedoCount()
		preview := urc.undoStack.GetRedoPreview()
		urc.redoMenuItem.Label = fmt.Sprintf("Redo (%d) - %s", redoCount, preview)
	} else {
		urc.redoMenuItem.Disabled = true
		urc.redoMenuItem.Label = "Redo"
	}

	// Update status label
	urc.updateStatusLabel()
}

// updateStatusLabel updates the status label display
func (urc *UndoRedoController) updateStatusLabel() {
	undo := urc.undoStack.GetUndoCount()
	redo := urc.undoStack.GetRedoCount()

	statusText := fmt.Sprintf("Undo: %d | Redo: %d", undo, redo)
	urc.statusLabel.SetText(statusText)
}

// HandleKeyboardShortcut handles keyboard shortcuts for undo/redo
func (urc *UndoRedoController) HandleKeyboardShortcut(key *fyne.KeyEvent) {
	// Note: Fyne v2 keyboard handling is different
	// This would need to be integrated with app-level key bindings
	switch key.Name {
	case fyne.KeyZ:
		if urc.undoStack.CanUndo() {
			urc.executeUndo()
		}
	case fyne.KeyY:
		if urc.undoStack.CanRedo() {
			urc.executeRedo()
		}
	}
}

// SetOnUndo sets the callback for undo operations
func (urc *UndoRedoController) SetOnUndo(callback func()) {
	urc.onUndo = callback
}

// SetOnRedo sets the callback for redo operations
func (urc *UndoRedoController) SetOnRedo(callback func()) {
	urc.onRedo = callback
}

// Clear clears all undo/redo history
func (urc *UndoRedoController) Clear() {
	urc.undoStack.Clear()
	urc.updateMenuState()
}

// GetUndoCount returns the number of undoable actions
func (urc *UndoRedoController) GetUndoCount() int {
	return urc.undoStack.GetUndoCount()
}

// GetRedoCount returns the number of redoable actions
func (urc *UndoRedoController) GetRedoCount() int {
	return urc.undoStack.GetRedoCount()
}

// GetUndoPreview returns the preview text for undo
func (urc *UndoRedoController) GetUndoPreview() string {
	return urc.undoStack.GetUndoPreview()
}

// GetRedoPreview returns the preview text for redo
func (urc *UndoRedoController) GetRedoPreview() string {
	return urc.undoStack.GetRedoPreview()
}

// ShowHistory shows the undo/redo history dialog
func (urc *UndoRedoController) ShowHistory(window fyne.Window) {
	historyText := fmt.Sprintf(
		"Undo Stack: %d items\nRedo Stack: %d items\n\nUndo Preview: %s\nRedo Preview: %s",
		urc.undoStack.GetUndoCount(),
		urc.undoStack.GetRedoCount(),
		urc.undoStack.GetUndoPreview(),
		urc.undoStack.GetRedoPreview(),
	)

	dialog.ShowInformation(
		"Undo/Redo History",
		historyText,
		window,
	)
}
