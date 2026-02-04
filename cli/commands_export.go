package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"ffvi_editor/models/consts"
	prconsts "ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"
)

// handleExportCommand exports save data to JSON in various formats
func (c *CLI) handleExportCommand(file, output, format string) error {
	// Load the save file
	_, err := c.LoadSaveFile(file)
	if err != nil {
		return fmt.Errorf("failed to load save file: %w", err)
	}

	// Create export data structure
	exportData := &ExportData{
		Metadata: ExportMetadata{
			SourceFile: file,
			Format:     format,
		},
	}

	// Populate based on format
	switch format {
	case "full":
		exportData.Characters = pri.Characters
		exportData.Party = pri.GetParty()
		exportData.Inventory = pri.GetInventory()
		exportData.Espers = getEspersForExport()
	case "characters":
		exportData.Characters = pri.Characters
	case "inventory":
		exportData.Inventory = pri.GetInventory()
	case "party":
		exportData.Party = pri.GetParty()
	case "magic":
		exportData.Characters = pri.Characters
	case "espers":
		exportData.Espers = getEspersForExport()
	default:
		return fmt.Errorf("unknown export format: %s (valid: full, characters, inventory, party, magic, espers)", format)
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	// Write to output file
	if err := os.WriteFile(output, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	fmt.Printf("Successfully exported save data to: %s (format: %s)\n", output, format)
	return nil
}

// getEspersForExport returns the list of espers
func getEspersForExport() []*consts.NameValueChecked {
	return prconsts.Espers
}
