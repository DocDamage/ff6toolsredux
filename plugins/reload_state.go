package plugins

import (
	"encoding/json"
	"fmt"
	"time"
)

// PluginState represents a snapshot of a plugin's state
type PluginState struct {
	PluginID       string
	Version        string
	Timestamp      time.Time
	Config         interface{}
	Enabled        bool
	Data           map[string]interface{}
	LastExecutions []ExecutionRecord
}

// Serialize converts the plugin state to JSON
func (s *PluginState) Serialize() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// Deserialize loads plugin state from JSON
func (s *PluginState) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}

// Validate checks if the state is valid
func (s *PluginState) Validate() error {
	if s.PluginID == "" {
		return fmt.Errorf("plugin ID is empty")
	}
	if s.Version == "" {
		return fmt.Errorf("version is empty")
	}
	if s.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is zero")
	}
	return nil
}

// Clone creates a deep copy of the plugin state
func (s *PluginState) Clone() (*PluginState, error) {
	data, err := s.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize state: %w", err)
	}

	clone := &PluginState{}
	if err := clone.Deserialize(data); err != nil {
		return nil, fmt.Errorf("failed to deserialize clone: %w", err)
	}

	return clone, nil
}

// Merge merges another state into this one
func (s *PluginState) Merge(other *PluginState) error {
	if other == nil {
		return fmt.Errorf("cannot merge nil state")
	}

	if s.PluginID != other.PluginID {
		return fmt.Errorf("cannot merge states from different plugins")
	}

	// Merge data maps
	if other.Data != nil {
		if s.Data == nil {
			s.Data = make(map[string]interface{})
		}
		for k, v := range other.Data {
			s.Data[k] = v
		}
	}

	// Update config if provided
	if other.Config != nil {
		s.Config = other.Config
	}

	// Keep the newer timestamp
	if other.Timestamp.After(s.Timestamp) {
		s.Timestamp = other.Timestamp
	}

	return nil
}

// Diff compares this state with another and returns differences
func (s *PluginState) Diff(other *PluginState) *StateDiff {
	diff := &StateDiff{
		PluginID:      s.PluginID,
		OldVersion:    s.Version,
		NewVersion:    other.Version,
		Changes:       make([]StateChange, 0),
		ConfigChanged: false,
	}

	// Check version change
	if s.Version != other.Version {
		diff.Changes = append(diff.Changes, StateChange{
			Field:    "Version",
			OldValue: s.Version,
			NewValue: other.Version,
		})
	}

	// Check enabled status
	if s.Enabled != other.Enabled {
		diff.Changes = append(diff.Changes, StateChange{
			Field:    "Enabled",
			OldValue: s.Enabled,
			NewValue: other.Enabled,
		})
	}

	// Check config (simple comparison)
	if s.Config != other.Config {
		diff.ConfigChanged = true
		diff.Changes = append(diff.Changes, StateChange{
			Field:    "Config",
			OldValue: s.Config,
			NewValue: other.Config,
		})
	}

	// Check data differences
	for key, oldVal := range s.Data {
		if newVal, exists := other.Data[key]; exists {
			if oldVal != newVal {
				diff.Changes = append(diff.Changes, StateChange{
					Field:    fmt.Sprintf("Data.%s", key),
					OldValue: oldVal,
					NewValue: newVal,
				})
			}
		} else {
			diff.Changes = append(diff.Changes, StateChange{
				Field:    fmt.Sprintf("Data.%s", key),
				OldValue: oldVal,
				NewValue: nil,
			})
		}
	}

	// Check for new data keys
	for key, newVal := range other.Data {
		if _, exists := s.Data[key]; !exists {
			diff.Changes = append(diff.Changes, StateChange{
				Field:    fmt.Sprintf("Data.%s", key),
				OldValue: nil,
				NewValue: newVal,
			})
		}
	}

	return diff
}

// StateDiff represents differences between two plugin states
type StateDiff struct {
	PluginID      string
	OldVersion    string
	NewVersion    string
	Changes       []StateChange
	ConfigChanged bool
}

// StateChange represents a single state change
type StateChange struct {
	Field    string
	OldValue interface{}
	NewValue interface{}
}

// HasChanges returns true if there are any changes
func (d *StateDiff) HasChanges() bool {
	return len(d.Changes) > 0
}

// Summary returns a human-readable summary of changes
func (d *StateDiff) Summary() string {
	if !d.HasChanges() {
		return "No changes detected"
	}

	summary := fmt.Sprintf("%d change(s) detected:\n", len(d.Changes))
	for _, change := range d.Changes {
		summary += fmt.Sprintf("  - %s: %v → %v\n", change.Field, change.OldValue, change.NewValue)
	}
	return summary
}

// StateSnapshot provides utilities for managing plugin state snapshots
type StateSnapshot struct {
	States    map[string]*PluginState
	Timestamp time.Time
	Metadata  map[string]string
}

// NewStateSnapshot creates a new state snapshot
func NewStateSnapshot() *StateSnapshot {
	return &StateSnapshot{
		States:    make(map[string]*PluginState),
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// AddState adds a plugin state to the snapshot
func (s *StateSnapshot) AddState(state *PluginState) error {
	if state == nil {
		return fmt.Errorf("cannot add nil state")
	}

	if err := state.Validate(); err != nil {
		return fmt.Errorf("invalid state: %w", err)
	}

	s.States[state.PluginID] = state
	return nil
}

// GetState retrieves a plugin state from the snapshot
func (s *StateSnapshot) GetState(pluginID string) (*PluginState, bool) {
	state, ok := s.States[pluginID]
	return state, ok
}

// RemoveState removes a plugin state from the snapshot
func (s *StateSnapshot) RemoveState(pluginID string) {
	delete(s.States, pluginID)
}

// Serialize converts the entire snapshot to JSON
func (s *StateSnapshot) Serialize() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// Deserialize loads a snapshot from JSON
func (s *StateSnapshot) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}

// Size returns the number of states in the snapshot
func (s *StateSnapshot) Size() int {
	return len(s.States)
}

// PluginIDs returns a list of all plugin IDs in the snapshot
func (s *StateSnapshot) PluginIDs() []string {
	ids := make([]string, 0, len(s.States))
	for id := range s.States {
		ids = append(ids, id)
	}
	return ids
}
