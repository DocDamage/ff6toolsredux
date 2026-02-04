package state

import (
	"fmt"
	"sync"
	"time"

	"ffvi_editor/models"
)

// UndoStack manages undo and redo operations
type UndoStack struct {
	undoStack       []models.Change
	redoStack       []models.Change
	maxDepth        int
	batchMode       bool
	currentBatchID  string
	currentBatchOps []models.Change
	mu              sync.Mutex
}

// NewUndoStack creates a new undo/redo stack
func NewUndoStack(maxDepth int) *UndoStack {
	if maxDepth <= 0 {
		maxDepth = 100 // Default
	}
	return &UndoStack{
		undoStack: make([]models.Change, 0, maxDepth),
		redoStack: make([]models.Change, 0, maxDepth),
		maxDepth:  maxDepth,
	}
}

// RecordChange adds a change to the undo stack
func (us *UndoStack) RecordChange(change models.Change) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if !models.ValidateChange(change) {
		return fmt.Errorf("invalid change: %+v", change)
	}

	// If in batch mode, accumulate changes
	if us.batchMode {
		change.Batch = true
		change.BatchID = us.currentBatchID
		us.currentBatchOps = append(us.currentBatchOps, change)
		return nil
	}

	// Add change to undo stack
	us.undoStack = append(us.undoStack, change)
	if len(us.undoStack) > us.maxDepth {
		us.undoStack = us.undoStack[len(us.undoStack)-us.maxDepth:]
	}

	// Clear redo stack (can't redo after new change)
	us.redoStack = make([]models.Change, 0)

	return nil
}

// StartBatch begins a batch of operations
func (us *UndoStack) StartBatch(name string) (string, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	if us.batchMode {
		return "", fmt.Errorf("batch already in progress")
	}

	us.batchMode = true
	us.currentBatchID = generateBatchID()
	us.currentBatchOps = make([]models.Change, 0)

	return us.currentBatchID, nil
}

// EndBatch completes a batch and records all changes as single undo action
func (us *UndoStack) EndBatch(name string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if !us.batchMode {
		return fmt.Errorf("no batch in progress")
	}

	// Add all batch operations to undo stack
	for _, change := range us.currentBatchOps {
		change.BatchName = name
		us.undoStack = append(us.undoStack, change)
	}

	// Maintain max depth
	if len(us.undoStack) > us.maxDepth {
		us.undoStack = us.undoStack[len(us.undoStack)-us.maxDepth:]
	}

	// Clear redo stack
	us.redoStack = make([]models.Change, 0)

	// Reset batch state
	us.batchMode = false
	us.currentBatchID = ""
	us.currentBatchOps = make([]models.Change, 0)

	return nil
}

// CancelBatch cancels current batch without recording
func (us *UndoStack) CancelBatch() error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if !us.batchMode {
		return fmt.Errorf("no batch in progress")
	}

	us.batchMode = false
	us.currentBatchID = ""
	us.currentBatchOps = make([]models.Change, 0)

	return nil
}

// CanUndo returns whether undo is available
func (us *UndoStack) CanUndo() bool {
	us.mu.Lock()
	defer us.mu.Unlock()
	return len(us.undoStack) > 0
}

// CanRedo returns whether redo is available
func (us *UndoStack) CanRedo() bool {
	us.mu.Lock()
	defer us.mu.Unlock()
	return len(us.redoStack) > 0
}

// PopUndo returns the next change to undo (doesn't actually undo)
func (us *UndoStack) PopUndo() ([]models.Change, string, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	if len(us.undoStack) == 0 {
		return nil, "", fmt.Errorf("nothing to undo")
	}

	// Find all changes in this batch (if batch) or single change
	lastChange := us.undoStack[len(us.undoStack)-1]
	batchName := lastChange.BatchName
	batchID := lastChange.BatchID

	if lastChange.Batch && batchID != "" {
		// Find all changes with this batch ID
		var batchChanges []models.Change
		for i := len(us.undoStack) - 1; i >= 0; i-- {
			if us.undoStack[i].BatchID == batchID {
				batchChanges = append([]models.Change{us.undoStack[i]}, batchChanges...)
			} else {
				break
			}
		}

		// Remove batch changes from undo stack
		us.undoStack = us.undoStack[:len(us.undoStack)-len(batchChanges)]

		// Add to redo stack
		us.redoStack = append(us.redoStack, batchChanges...)

		return batchChanges, batchName, nil
	}

	// Single change
	change := us.undoStack[len(us.undoStack)-1]
	us.undoStack = us.undoStack[:len(us.undoStack)-1]
	us.redoStack = append(us.redoStack, change)

	return []models.Change{change}, change.FieldName, nil
}

// PopRedo returns the next change to redo
func (us *UndoStack) PopRedo() ([]models.Change, string, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	if len(us.redoStack) == 0 {
		return nil, "", fmt.Errorf("nothing to redo")
	}

	// Similar logic to PopUndo but in reverse
	lastChange := us.redoStack[len(us.redoStack)-1]
	batchID := lastChange.BatchID

	if lastChange.Batch && batchID != "" {
		// Find all changes with this batch ID
		var batchChanges []models.Change
		for i := len(us.redoStack) - 1; i >= 0; i-- {
			if us.redoStack[i].BatchID == batchID {
				batchChanges = append([]models.Change{us.redoStack[i]}, batchChanges...)
			} else {
				break
			}
		}

		// Remove from redo stack
		us.redoStack = us.redoStack[:len(us.redoStack)-len(batchChanges)]

		// Add to undo stack
		us.undoStack = append(us.undoStack, batchChanges...)

		return batchChanges, lastChange.BatchName, nil
	}

	// Single change
	change := us.redoStack[len(us.redoStack)-1]
	us.redoStack = us.redoStack[:len(us.redoStack)-1]
	us.undoStack = append(us.undoStack, change)

	return []models.Change{change}, change.FieldName, nil
}

// GetUndoCount returns number of undo operations available
func (us *UndoStack) GetUndoCount() int {
	us.mu.Lock()
	defer us.mu.Unlock()
	return len(us.undoStack)
}

// GetRedoCount returns number of redo operations available
func (us *UndoStack) GetRedoCount() int {
	us.mu.Lock()
	defer us.mu.Unlock()
	return len(us.redoStack)
}

// GetUndoPreview returns what the next undo would do
func (us *UndoStack) GetUndoPreview() string {
	us.mu.Lock()
	defer us.mu.Unlock()

	if len(us.undoStack) == 0 {
		return ""
	}

	lastChange := us.undoStack[len(us.undoStack)-1]
	if lastChange.BatchName != "" {
		return fmt.Sprintf("Undo: %s", lastChange.BatchName)
	}
	return fmt.Sprintf("Undo: Set %s to %v", lastChange.FieldName, lastChange.OldValue)
}

// GetRedoPreview returns what the next redo would do
func (us *UndoStack) GetRedoPreview() string {
	us.mu.Lock()
	defer us.mu.Unlock()

	if len(us.redoStack) == 0 {
		return ""
	}

	lastChange := us.redoStack[len(us.redoStack)-1]
	if lastChange.BatchName != "" {
		return fmt.Sprintf("Redo: %s", lastChange.BatchName)
	}
	return fmt.Sprintf("Redo: Set %s to %v", lastChange.FieldName, lastChange.NewValue)
}

// Clear resets both undo and redo stacks
func (us *UndoStack) Clear() {
	us.mu.Lock()
	defer us.mu.Unlock()

	us.undoStack = make([]models.Change, 0, us.maxDepth)
	us.redoStack = make([]models.Change, 0, us.maxDepth)
	us.batchMode = false
	us.currentBatchID = ""
	us.currentBatchOps = make([]models.Change, 0)
}

// Undo performs an undo operation and returns the changes
func (us *UndoStack) Undo() ([]models.Change, error) {
	changes, _, err := us.PopUndo()
	return changes, err
}

// Redo performs a redo operation and returns the changes
func (us *UndoStack) Redo() ([]models.Change, error) {
	changes, _, err := us.PopRedo()
	return changes, err
}

// UndoDepth returns the number of undo operations available
func (us *UndoStack) UndoDepth() int {
	return us.GetUndoCount()
}

// RedoDepth returns the number of redo operations available
func (us *UndoStack) RedoDepth() int {
	return us.GetRedoCount()
}

// generateBatchID creates a unique batch ID
func generateBatchID() string {
	return fmt.Sprintf("batch_%d", time.Now().UnixNano())
}
