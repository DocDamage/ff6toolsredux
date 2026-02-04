package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	prconsts "ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"
)

// handleImportCommand imports JSON data into a save file
// Imports character stats, inventory, party composition, magic, and espers
// from a previously exported JSON file into the target save file
func (c *CLI) handleImportCommand(file, input, format string, backup bool) error {
	// 1. Create backup if requested
	if backup {
		backupPath := fmt.Sprintf("%s.%s.backup", file, time.Now().Format("20060102_150405"))
		if err := c.handleBackupCommand(file, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// 2. Load the target save file to be modified
	save, err := c.LoadSaveFile(file)
	if err != nil {
		return fmt.Errorf("failed to load save file: %w", err)
	}

	// 3. Read and parse the input JSON
	jsonData, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	var importData ExportData
	if err := json.Unmarshal(jsonData, &importData); err != nil {
		return fmt.Errorf("failed to parse import JSON: %w", err)
	}

	// 4. Apply imported data based on format
	switch format {
	case "full":
		if err := importCharacters(importData.Characters); err != nil {
			return fmt.Errorf("failed to import characters: %w", err)
		}
		if err := importParty(importData.Party); err != nil {
			return fmt.Errorf("failed to import party: %w", err)
		}
		if err := importInventory(importData.Inventory); err != nil {
			return fmt.Errorf("failed to import inventory: %w", err)
		}
		if err := importEspers(importData.Espers); err != nil {
			return fmt.Errorf("failed to import espers: %w", err)
		}
		fmt.Println("Imported: characters, party, inventory, espers, magic")
	case "characters":
		if err := importCharacters(importData.Characters); err != nil {
			return fmt.Errorf("failed to import characters: %w", err)
		}
		fmt.Println("Imported: characters")
	case "inventory":
		if err := importInventory(importData.Inventory); err != nil {
			return fmt.Errorf("failed to import inventory: %w", err)
		}
		fmt.Println("Imported: inventory")
	case "party":
		if err := importParty(importData.Party); err != nil {
			return fmt.Errorf("failed to import party: %w", err)
		}
		fmt.Println("Imported: party")
	case "magic":
		if err := importMagic(importData.Characters); err != nil {
			return fmt.Errorf("failed to import magic: %w", err)
		}
		fmt.Println("Imported: magic")
	case "espers":
		if err := importEspers(importData.Espers); err != nil {
			return fmt.Errorf("failed to import espers: %w", err)
		}
		fmt.Println("Imported: espers")
	default:
		return fmt.Errorf("unknown import format: %s (valid: full, characters, inventory, party, magic, espers)", format)
	}

	// 5. Save the modified file
	if err := c.SaveSaveFile(save, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("Successfully imported data from %s into: %s\n", input, file)
	return nil
}

// importCharacters imports character stats, levels, and HP/MP
func importCharacters(characters []*models.Character) error {
	if characters == nil || len(characters) == 0 {
		return nil
	}
	for _, importChar := range characters {
		if importChar == nil {
			continue
		}
		// Find matching character by ID
		char := pri.GetCharacterByID(importChar.ID)
		if char == nil {
			continue
		}
		// Import basic stats
		if importChar.Name != "" {
			char.Name = importChar.Name
		}
		if importChar.Level > 0 {
			char.Level = importChar.Level
		}
		if importChar.Exp >= 0 {
			char.Exp = importChar.Exp
		}
		// Import HP/MP
		if importChar.HP.Current > 0 {
			char.HP.Current = importChar.HP.Current
		}
		if importChar.HP.Max > 0 {
			char.HP.Max = importChar.HP.Max
		}
		if importChar.MP.Current > 0 {
			char.MP.Current = importChar.MP.Current
		}
		if importChar.MP.Max > 0 {
			char.MP.Max = importChar.MP.Max
		}
		// Import stats
		if importChar.Vigor > 0 {
			char.Vigor = importChar.Vigor
		}
		if importChar.Stamina > 0 {
			char.Stamina = importChar.Stamina
		}
		if importChar.Speed > 0 {
			char.Speed = importChar.Speed
		}
		if importChar.Magic > 0 {
			char.Magic = importChar.Magic
		}
		// Import enabled status
		char.IsEnabled = importChar.IsEnabled
		// Import equipment
		char.Equipment = importChar.Equipment
		// Import spells
		if importChar.SpellsByID != nil {
			for spellID, spell := range importChar.SpellsByID {
				if existingSpell, ok := char.SpellsByID[spellID]; ok && spell != nil {
					existingSpell.Value = spell.Value
				}
			}
		}
	}
	return nil
}

// importParty imports party composition
func importParty(partyData *pri.Party) error {
	if partyData == nil {
		return nil
	}
	party := pri.GetParty()
	// Import party members
	for i, member := range partyData.Members {
		if i >= 4 {
			break
		}
		if member != nil {
			party.Members[i] = member
		}
	}
	return nil
}

// importInventory imports inventory items
func importInventory(inventoryData *pri.Inventory) error {
	if inventoryData == nil || inventoryData.Rows == nil {
		return nil
	}
	inv := pri.GetInventory()
	// Clear existing inventory
	inv.Clear()
	// Import items
	for i, row := range inventoryData.Rows {
		if i >= inv.Size {
			break
		}
		if row != nil && row.ItemID > 0 && row.Count > 0 {
			inv.Rows[i].ItemID = row.ItemID
			inv.Rows[i].Count = row.Count
		}
	}
	return nil
}

// importMagic imports spell learn status for all characters
func importMagic(characters []*models.Character) error {
	if characters == nil || len(characters) == 0 {
		return nil
	}
	for _, importChar := range characters {
		if importChar == nil || importChar.SpellsByID == nil {
			continue
		}
		// Find matching character by ID
		char := pri.GetCharacterByID(importChar.ID)
		if char == nil {
			continue
		}
		// Import spells
		for spellID, spell := range importChar.SpellsByID {
			if existingSpell, ok := char.SpellsByID[spellID]; ok && spell != nil {
				existingSpell.Value = spell.Value
			}
		}
	}
	return nil
}

// importEspers imports esper unlock status
func importEspers(espers []*consts.NameValueChecked) error {
	if espers == nil || len(espers) == 0 {
		return nil
	}
	for _, esper := range espers {
		if esper == nil {
			continue
		}
		// Find matching esper by value
		if existingEsper, ok := prconsts.EspersByValue[esper.Value]; ok {
			existingEsper.Checked = esper.Checked
		}
	}
	return nil
}
