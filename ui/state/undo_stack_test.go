package state

import (
	"ffvi_editor/models"
	"testing"
)

func TestUndoStack_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		changes    []models.Change
		doUndo     int
		doRedo     int
		expectUndo int
		expectRedo int
	}{
		{
			name: "Single change, undo/redo",
			changes: []models.Change{
				models.NewChange("Character", "HP", 100, 200),
			},
			doUndo: 1, doRedo: 1, expectUndo: 1, expectRedo: 0,
		},
		{
			name: "Multiple changes, undo all, redo all",
			changes: []models.Change{
				models.NewChange("Character", "HP", 100, 200),
				models.NewChange("Character", "MP", 50, 75),
			},
			doUndo: 2, doRedo: 2, expectUndo: 2, expectRedo: 0,
		},
		{
			name: "Undo more than stack",
			changes: []models.Change{
				models.NewChange("Character", "HP", 100, 200),
			},
			doUndo: 2, doRedo: 0, expectUndo: 0, expectRedo: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewUndoStack(10)
			for _, c := range tt.changes {
				if err := us.RecordChange(c); err != nil {
					t.Fatalf("RecordChange failed: %v", err)
				}
			}
			for i := 0; i < tt.doUndo; i++ {
				_, _ = us.Undo()
			}
			for i := 0; i < tt.doRedo; i++ {
				_, _ = us.Redo()
			}
			if got := us.UndoDepth(); got != tt.expectUndo {
				t.Errorf("UndoDepth = %d, want %d", got, tt.expectUndo)
			}
			if got := us.RedoDepth(); got != tt.expectRedo {
				t.Errorf("RedoDepth = %d, want %d", got, tt.expectRedo)
			}
		})
	}
}

func TestUndoStack_Property_Idempotence(t *testing.T) {
	us := NewUndoStack(5)
	c := models.NewChange("Character", "HP", 100, 200)
	_ = us.RecordChange(c)
	for i := 0; i < 10; i++ {
		_, _ = us.Undo()
	}
	for i := 0; i < 10; i++ {
		_, _ = us.Redo()
	}
	if us.UndoDepth() != 1 {
		t.Errorf("UndoDepth should be 1 after excessive undo/redo, got %d", us.UndoDepth())
	}
}
