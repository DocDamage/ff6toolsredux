package json

import (
	"encoding/json"
	"fmt"
	"os"

	ipr "ffvi_editor/io/pr"
)

// ImportError represents an error during import
type ImportError struct {
	Field   string
	Message string
}

func (e ImportError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("import error in %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("import error: %s", e.Message)
}

// Importer handles importing JSON data into save data
type Importer struct {
	prData *ipr.PR
	errors []ImportError
}

// NewImporter creates a new importer for the given save data
func NewImporter(pr *ipr.PR) *Importer {
	return &Importer{
		prData: pr,
		errors: make([]ImportError, 0),
	}
}

// importFunc is a function that imports a specific section
type importFunc func(*Importer, SaveExport) error

// importCharactersFunc imports character data
func importCharactersFunc(i *Importer, export SaveExport) error {
	return i.importCharacters(export.Characters)
}

// importPartyFunc imports party data
func importPartyFunc(i *Importer, export SaveExport) error {
	return i.importParty(export.Party)
}

// importInventoryFunc imports inventory data
func importInventoryFunc(i *Importer, export SaveExport) error {
	return i.importInventory(export.Inventory)
}

// importEquipmentFunc imports equipment data
func importEquipmentFunc(i *Importer, export SaveExport) error {
	return i.importEquipment(export.Equipment)
}

// importMagicFunc imports magic data
func importMagicFunc(i *Importer, export SaveExport) error {
	return i.importMagic(export.Magic)
}

// importEspersFunc imports esper data
func importEspersFunc(i *Importer, export SaveExport) error {
	return i.importEspers(export.Espers)
}

// ImportFromJSON imports save data from JSON bytes
func (i *Importer) ImportFromJSON(jsonBytes []byte, format ExportFormat) error {
	var export SaveExport
	if err := json.Unmarshal(jsonBytes, &export); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	i.errors = make([]ImportError, 0)

	// Map of format to import functions
	formatImporters := map[ExportFormat][]importFunc{
		FormatFull:       {importCharactersFunc, importPartyFunc, importInventoryFunc, importEquipmentFunc, importMagicFunc, importEspersFunc},
		FormatCharacters: {importCharactersFunc},
		FormatInventory:  {importInventoryFunc},
		FormatParty:      {importPartyFunc},
		FormatMagic:      {importMagicFunc},
		FormatEspers:     {importEspersFunc},
		FormatEquipment:  {importEquipmentFunc},
	}

	importers, ok := formatImporters[format]
	if !ok {
		return fmt.Errorf("unknown import format: %s", format)
	}

	for _, importer := range importers {
		if err := importer(i, export); err != nil {
			return err
		}
	}

	// Return errors if any occurred
	if len(i.errors) > 0 {
		errMsg := "import completed with errors:\n"
		for _, err := range i.errors {
			errMsg += fmt.Sprintf("  - %s\n", err.Error())
		}
		// Log errors but don't fail import
		fmt.Println(errMsg)
	}

	return nil
}

// ImportFromFile imports save data from a JSON file
func (i *Importer) ImportFromFile(filePath string, format ExportFormat) error {
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	return i.ImportFromJSON(jsonBytes, format)
}

// GetErrors returns any errors that occurred during import
func (i *Importer) GetErrors() []ImportError {
	return i.errors
}

// importCharacters applies character data from export
func (i *Importer) importCharacters(characters []CharacterExport) error {
	if len(characters) == 0 {
		return nil
	}
	// Note: Full implementation would require modifying PR OrderedMap structure
	// For now, this is a placeholder for schema definition
	return nil
}

// importParty applies party data from export
func (i *Importer) importParty(party *PartyExport) error {
	if party == nil || len(party.Members) == 0 {
		return nil
	}
	// Note: Placeholder for party import logic
	return nil
}

// importInventory applies inventory data from export
func (i *Importer) importInventory(inventory *InventoryExport) error {
	if inventory == nil || len(inventory.Items) == 0 {
		return nil
	}
	// Note: Placeholder for inventory import logic
	return nil
}

// importEquipment applies equipment data from export
func (i *Importer) importEquipment(equipment map[string]EquipmentExport) error {
	if len(equipment) == 0 {
		return nil
	}
	// Note: Placeholder for equipment import logic
	return nil
}

// importMagic applies magic data from export
func (i *Importer) importMagic(magic map[string]MagicExport) error {
	if len(magic) == 0 {
		return nil
	}
	fmt.Printf("Magic import requested for %d characters\n", len(magic))
	return nil
}

// importEspers applies esper data from export
func (i *Importer) importEspers(espers *EsperExport) error {
	if espers == nil {
		return nil
	}
	fmt.Printf("Esper import requested with %d unlocked espers\n", len(espers.Unlocked))
	return nil
}

// addError adds an import error to the error list
func (i *Importer) addError(field string, message string) {
	i.errors = append(i.errors, ImportError{
		Field:   field,
		Message: message,
	})
}
