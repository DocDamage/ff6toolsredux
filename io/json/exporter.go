package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	ipr "ffvi_editor/io/pr"
)

// ExportFormat determines which parts of the save to export
type ExportFormat string

const (
	FormatFull       ExportFormat = "full"
	FormatCharacters ExportFormat = "characters"
	FormatInventory  ExportFormat = "inventory"
	FormatParty      ExportFormat = "party"
	FormatMagic      ExportFormat = "magic"
	FormatEspers     ExportFormat = "espers"
	FormatEquipment  ExportFormat = "equipment"
)

// SaveExport represents exported save data in human-readable JSON format
type SaveExport struct {
	Format     string                     `json:"format"`
	Version    string                     `json:"version"`
	Characters []CharacterExport          `json:"characters,omitempty"`
	Party      *PartyExport               `json:"party,omitempty"`
	Inventory  *InventoryExport           `json:"inventory,omitempty"`
	Equipment  map[string]EquipmentExport `json:"equipment,omitempty"`
	Magic      map[string]MagicExport     `json:"magic,omitempty"`
	Espers     *EsperExport               `json:"espers,omitempty"`
	Metadata   ExportMetadata             `json:"metadata"`
}

// ExportMetadata contains metadata about the export
type ExportMetadata struct {
	ExportedAt string `json:"exportedAt"`
	Format     string `json:"format"`
	Note       string `json:"note,omitempty"`
}

// CharacterExport represents a character for export
type CharacterExport struct {
	Name     string     `json:"name"`
	Level    uint8      `json:"level"`
	HP       uint16     `json:"hp"`
	MaxHP    uint16     `json:"maxHp"`
	MP       uint8      `json:"mp"`
	MaxMP    uint8      `json:"maxMp"`
	Exp      uint32     `json:"exp"`
	Stats    StatExport `json:"stats"`
	Status   []string   `json:"status,omitempty"`
	Commands []string   `json:"commands,omitempty"`
	Espers   []string   `json:"espers,omitempty"`
	Relic1   string     `json:"relic1,omitempty"`
	Relic2   string     `json:"relic2,omitempty"`
	RowState string     `json:"rowState,omitempty"`
}

// StatExport represents character stats for export
type StatExport struct {
	Vigor    uint8 `json:"vigor"`
	Speed    uint8 `json:"speed"`
	Stamina  uint8 `json:"stamina"`
	MagicPwr uint8 `json:"magicPwr"`
	Defense  uint8 `json:"defense"`
	MagicDef uint8 `json:"magicDef"`
}

// PartyExport represents party composition for export
type PartyExport struct {
	Members []string `json:"members"`
}

// InventoryExport represents inventory for export
type InventoryExport struct {
	Items []ItemExport `json:"items"`
}

// ItemExport represents a single item for export
type ItemExport struct {
	Name     string `json:"name"`
	Quantity uint8  `json:"quantity"`
}

// EquipmentExport represents equipment for export
type EquipmentExport struct {
	Character string `json:"character"`
	Weapon    string `json:"weapon,omitempty"`
	Armor     string `json:"armor,omitempty"`
	Shield    string `json:"shield,omitempty"`
	Helmet    string `json:"helmet,omitempty"`
	Relic1    string `json:"relic1,omitempty"`
	Relic2    string `json:"relic2,omitempty"`
}

// MagicExport represents magic learned for export
type MagicExport struct {
	Character string   `json:"character"`
	Spells    []string `json:"spells"`
}

// EsperExport represents esper information for export
type EsperExport struct {
	Unlocked []string          `json:"unlocked"`
	Equipped map[string]string `json:"equipped"`
}

// Exporter handles exporting save data to JSON format
type Exporter struct {
	prData *ipr.PR
}

// NewExporter creates a new exporter for the given save data
func NewExporter(pr *ipr.PR) *Exporter {
	return &Exporter{
		prData: pr,
	}
}

// ExportToJSON exports save data in specified format to JSON bytes
func (e *Exporter) ExportToJSON(format ExportFormat) ([]byte, error) {
	export := &SaveExport{
		Format:  string(format),
		Version: "1.0",
		Metadata: ExportMetadata{
			ExportedAt: time.Now().Format(time.RFC3339),
			Format:     string(format),
		},
	}

	switch format {
	case FormatFull:
		if err := e.populateCharacters(export); err != nil {
			return nil, fmt.Errorf("failed to export characters: %w", err)
		}
		if err := e.populateParty(export); err != nil {
			return nil, fmt.Errorf("failed to export party: %w", err)
		}
		if err := e.populateInventory(export); err != nil {
			return nil, fmt.Errorf("failed to export inventory: %w", err)
		}
		if err := e.populateEquipment(export); err != nil {
			return nil, fmt.Errorf("failed to export equipment: %w", err)
		}
		if err := e.populateMagic(export); err != nil {
			return nil, fmt.Errorf("failed to export magic: %w", err)
		}
		if err := e.populateEspers(export); err != nil {
			return nil, fmt.Errorf("failed to export espers: %w", err)
		}

	case FormatCharacters:
		if err := e.populateCharacters(export); err != nil {
			return nil, fmt.Errorf("failed to export characters: %w", err)
		}

	case FormatInventory:
		if err := e.populateInventory(export); err != nil {
			return nil, fmt.Errorf("failed to export inventory: %w", err)
		}

	case FormatParty:
		if err := e.populateParty(export); err != nil {
			return nil, fmt.Errorf("failed to export party: %w", err)
		}

	case FormatMagic:
		if err := e.populateMagic(export); err != nil {
			return nil, fmt.Errorf("failed to export magic: %w", err)
		}

	case FormatEspers:
		if err := e.populateEspers(export); err != nil {
			return nil, fmt.Errorf("failed to export espers: %w", err)
		}

	case FormatEquipment:
		if err := e.populateEquipment(export); err != nil {
			return nil, fmt.Errorf("failed to export equipment: %w", err)
		}

	default:
		return nil, fmt.Errorf("unknown export format: %s", format)
	}

	// Marshal with pretty indentation
	jsonBytes, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonBytes, nil
}

// ExportToFile exports save data to a JSON file
func (e *Exporter) ExportToFile(format ExportFormat, filePath string) error {
	jsonBytes, err := e.ExportToJSON(format)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
			return fmt.Errorf("failed to write JSON file: %w", err)
		}
	}

	return os.WriteFile(filePath, jsonBytes, 0644)
}

// populateCharacters adds character data to export
func (e *Exporter) populateCharacters(export *SaveExport) error {
	export.Characters = make([]CharacterExport, 0)

	// Note: Full implementation would require parsing the PR OrderedMap structure
	// For now, this is a placeholder for schema definition
	// In production, this would extract character data from pr.Characters OrderedMaps

	return nil
}

// populateParty adds party data to export
func (e *Exporter) populateParty(export *SaveExport) error {
	partyExport := &PartyExport{
		Members: make([]string, 0),
	}
	export.Party = partyExport
	return nil
}

// populateInventory adds inventory data to export
func (e *Exporter) populateInventory(export *SaveExport) error {
	inventoryExport := &InventoryExport{
		Items: make([]ItemExport, 0),
	}
	export.Inventory = inventoryExport
	return nil
}

// populateEquipment adds equipment data to export
func (e *Exporter) populateEquipment(export *SaveExport) error {
	export.Equipment = make(map[string]EquipmentExport)
	return nil
}

// populateMagic adds magic data to export
func (e *Exporter) populateMagic(export *SaveExport) error {
	export.Magic = make(map[string]MagicExport)
	return nil
}

// populateEspers adds esper data to export
func (e *Exporter) populateEspers(export *SaveExport) error {
	export.Espers = &EsperExport{
		Unlocked: make([]string, 0),
		Equipped: make(map[string]string),
	}
	return nil
}
