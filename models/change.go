package models

import (
	"time"
)

// Change represents a single edit operation that can be undone/redone
type Change struct {
	ID        string        `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Target    string        `json:"target"`        // "Character", "Inventory", etc
	FieldName string        `json:"fieldName"`     // e.g., "HP", "Level"
	OldValue  interface{}   `json:"oldValue"`
	NewValue  interface{}   `json:"newValue"`
	Batch     bool          `json:"batch"`         // Part of batch operation?
	BatchID   string        `json:"batchID"`       // Groups multiple changes
	BatchName string        `json:"batchName"`     // User-friendly batch name
}

// ChangeGroup represents multiple changes that should be undone together
type ChangeGroup struct {
	ID      string
	Name    string
	Changes []Change
	Time    time.Time
}

// NewChange creates a new change record
func NewChange(target, fieldName string, oldValue, newValue interface{}) Change {
	return Change{
		ID:        generateChangeID(),
		Timestamp: time.Now(),
		Target:    target,
		FieldName: fieldName,
		OldValue:  oldValue,
		NewValue:  newValue,
		Batch:     false,
		BatchID:   "",
	}
}

// NewBatchChange creates a new change that's part of a batch
func NewBatchChange(batchID, batchName, target, fieldName string, oldValue, newValue interface{}) Change {
	change := NewChange(target, fieldName, oldValue, newValue)
	change.Batch = true
	change.BatchID = batchID
	change.BatchName = batchName
	return change
}

// generateChangeID creates a unique change ID
func generateChangeID() string {
	return time.Now().Format("20060102150405000") // Nanosecond precision
}

// ValidateChange checks if a change is valid
func ValidateChange(c Change) bool {
	if c.Target == "" || c.FieldName == "" {
		return false
	}
	if c.OldValue == nil || c.NewValue == nil {
		return false
	}
	return true
}
