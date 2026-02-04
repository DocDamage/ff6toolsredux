package cli

import (
	"encoding/json"
	"fmt"

	"ffvi_editor/io/pr"
	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

// LoadSaveFile loads a save file from the specified path
func (c *CLI) LoadSaveFile(filepath string) (*pr.PR, error) {
	p := pr.New()
	if err := p.Load(filepath, 0); err != nil {
		return nil, fmt.Errorf("failed to load save file: %w", err)
	}
	return p, nil
}

// SaveSaveFile saves a save file to the specified path
func (c *CLI) SaveSaveFile(save *pr.PR, filepath string) error {
	if err := save.Save(0, filepath, 0); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	fmt.Printf("Successfully saved to: %s\n", filepath)
	return nil
}

// ExportData represents the structure for CLI export
type ExportData struct {
	Metadata   ExportMetadata             `json:"metadata"`
	Characters []*models.Character        `json:"characters,omitempty"`
	Party      *pri.Party                 `json:"party,omitempty"`
	Inventory  *pri.Inventory             `json:"inventory,omitempty"`
	Espers     []*consts.NameValueChecked `json:"espers,omitempty"`
}

// ExportMetadata contains metadata about the export
type ExportMetadata struct {
	SourceFile string `json:"sourceFile"`
	Format     string `json:"format"`
}

// validationIssue represents a single validation problem
type validationIssue struct {
	severity string // "error", "warning", "info"
	field    string
	message  string
	fixable  bool
}

// Helper functions for working with OrderedMap
func getIntFromMap(m *jo.OrderedMap, key string) (int, error) {
	i, ok := m.GetValue(key)
	if !ok {
		return 0, fmt.Errorf("key not found: %s", key)
	}
	switch v := i.(type) {
	case json.Number:
		val, err := v.Int64()
		return int(val), err
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("expected number, got %T", i)
	}
}

func unmarshalFromMap(from *jo.OrderedMap, key string, m *jo.OrderedMap) error {
	i, ok := from.GetValue(key)
	if !ok || i == nil {
		return fmt.Errorf("unable to find %s", key)
	}
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("unable to unmarshal %s: expected string", key)
	}
	return m.UnmarshalJSON([]byte(s))
}
