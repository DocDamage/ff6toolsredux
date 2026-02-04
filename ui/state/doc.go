// Package state provides UI state management utilities.
//
// The state package manages application state that persists during
// the UI session, including:
//   - Undo/redo stacks for editor operations
//   - Batch operation grouping
//   - Change tracking and history management
//
// This package enables users to undo and redo changes made to save files,
// providing a safe editing experience with the ability to revert mistakes.
//
// # Usage
//
//	stack := state.NewUndoStack(100) // Max 100 undo levels
//	stack.Push(change)
//	stack.Undo()
//	stack.Redo()
//
// The undo stack is thread-safe and can be used concurrently from
// multiple UI goroutines.
package state
