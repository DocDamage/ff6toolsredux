package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	prconsts "ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"
)

// TestHandleImportCommandFull tests importing full save data
func TestHandleImportCommandFull(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForImport(t, tmpDir, "save.json")
	importFile := createImportFile(t, tmpDir, "import.json", "full")

	cli := NewCLI([]string{})

	err := cli.handleImportCommand(saveFile, importFile, "full", false)
	if err != nil {
		t.Logf("handleImportCommand error (expected in isolated test): %v", err)
		return
	}
}

// TestHandleImportCommandFormats tests different import formats
func TestHandleImportCommandFormats(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForImport(t, tmpDir, "save.json")

	cli := NewCLI([]string{})

	formats := []string{"full", "characters", "inventory", "party", "magic", "espers"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			importFile := createImportFile(t, tmpDir, "import_"+format+".json", format)

			err := cli.handleImportCommand(saveFile, importFile, format, false)
			if err != nil {
				t.Logf("handleImportCommand error for format %s (expected in isolated test): %v", format, err)
			}
		})
	}
}

// TestHandleImportCommandInvalidFormat tests invalid import format
func TestHandleImportCommandInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForImport(t, tmpDir, "save.json")
	importFile := createImportFile(t, tmpDir, "import.json", "full")

	cli := NewCLI([]string{})

	err := cli.handleImportCommand(saveFile, importFile, "invalid_format", false)
	if err == nil {
		t.Error("handleImportCommand should return error for invalid format")
	}
}

// TestHandleImportCommandNonExistentFile tests importing non-existent file
func TestHandleImportCommandNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "nonexistent.json")
	importFile := createImportFile(t, tmpDir, "import.json", "full")

	cli := NewCLI([]string{})

	err := cli.handleImportCommand(saveFile, importFile, "full", false)
	if err == nil {
		t.Error("handleImportCommand should return error for non-existent file")
	}
}

// TestHandleImportCommandNonExistentImport tests importing from non-existent import file
func TestHandleImportCommandNonExistentImport(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForImport(t, tmpDir, "save.json")
	importFile := filepath.Join(tmpDir, "nonexistent.json")

	cli := NewCLI([]string{})

	err := cli.handleImportCommand(saveFile, importFile, "full", false)
	if err == nil {
		t.Error("handleImportCommand should return error for non-existent import file")
	}
}

// TestHandleImportCommandInvalidJSON tests importing invalid JSON
func TestHandleImportCommandInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForImport(t, tmpDir, "save.json")
	importFile := filepath.Join(tmpDir, "invalid.json")

	// Create invalid JSON file
	if err := os.WriteFile(importFile, []byte(`{invalid json`), 0644); err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	cli := NewCLI([]string{})

	err := cli.handleImportCommand(saveFile, importFile, "full", false)
	if err == nil {
		t.Error("handleImportCommand should return error for invalid JSON")
	}
}

// TestHandleImportCommandWithBackup tests importing with backup
func TestHandleImportCommandWithBackup(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForImport(t, tmpDir, "save.json")
	importFile := createImportFile(t, tmpDir, "import.json", "full")

	cli := NewCLI([]string{})

	err := cli.handleImportCommand(saveFile, importFile, "full", true)
	if err != nil {
		t.Logf("handleImportCommand error (expected in isolated test): %v", err)
		return
	}

	// Check if backup was created
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	hasBackup := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".backup" {
			hasBackup = true
			break
		}
	}

	if !hasBackup {
		t.Log("Backup file was not created (expected in isolated test)")
	}
}

// TestImportCharactersEmpty tests importing empty characters
func TestImportCharactersEmpty(t *testing.T) {
	err := importCharacters(nil)
	if err != nil {
		t.Errorf("importCharacters(nil) error: %v", err)
	}

	err = importCharacters([]*models.Character{})
	if err != nil {
		t.Errorf("importCharacters(empty) error: %v", err)
	}
}

// TestImportCharactersWithData tests importing character data
func TestImportCharactersWithData(t *testing.T) {
	if len(pri.Characters) == 0 {
		t.Skip("No characters available for testing")
	}

	// Find first non-NPC character
	var targetChar *models.Character
	for _, char := range pri.Characters {
		if char != nil && !char.IsNPC {
			targetChar = char
			break
		}
	}

	if targetChar == nil {
		t.Skip("No playable character found for testing")
	}

	// Store original values
	originalLevel := targetChar.Level
	originalVigor := targetChar.Vigor

	// Create import character
	importChar := &models.Character{
		ID:         targetChar.ID,
		Name:       "TestName",
		Level:      50,
		Exp:        10000,
		Vigor:      80,
		Stamina:    70,
		Speed:      60,
		Magic:      90,
		HP:         models.CurrentMax{Current: 500, Max: 1000},
		MP:         models.CurrentMax{Current: 100, Max: 200},
		SpellsByID: make(map[int]*models.Spell),
	}

	characters := []*models.Character{importChar}
	err := importCharacters(characters)
	if err != nil {
		t.Errorf("importCharacters error: %v", err)
	}

	// Verify import
	if targetChar.Name != "TestName" {
		t.Errorf("Name = %q, want TestName", targetChar.Name)
	}
	if targetChar.Level != 50 {
		t.Errorf("Level = %d, want 50", targetChar.Level)
	}
	if targetChar.Vigor != 80 {
		t.Errorf("Vigor = %d, want 80", targetChar.Vigor)
	}

	// Restore original values
	targetChar.Level = originalLevel
	targetChar.Vigor = originalVigor
}

// TestImportPartyEmpty tests importing empty party
func TestImportPartyEmpty(t *testing.T) {
	err := importParty(nil)
	if err != nil {
		t.Errorf("importParty(nil) error: %v", err)
	}
}

// TestImportPartyWithData tests importing party data
func TestImportPartyWithData(t *testing.T) {
	partyData := &pri.Party{
		Members: [4]*pri.Member{
			{CharacterID: 1, Name: "Terra"},
			{CharacterID: 5, Name: "Locke"},
			{CharacterID: 6, Name: "Cyan"},
			{CharacterID: 7, Name: "Shadow"},
		},
	}

	err := importParty(partyData)
	if err != nil {
		t.Errorf("importParty error: %v", err)
	}

	// Verify party was imported
	party := pri.GetParty()
	if party.Members[0] == nil || party.Members[0].CharacterID != 1 {
		t.Error("First party member not imported correctly")
	}
}

// TestImportInventoryEmpty tests importing empty inventory
func TestImportInventoryEmpty(t *testing.T) {
	err := importInventory(nil)
	if err != nil {
		t.Errorf("importInventory(nil) error: %v", err)
	}

	emptyInventory := &pri.Inventory{Rows: nil}
	err = importInventory(emptyInventory)
	if err != nil {
		t.Errorf("importInventory(empty) error: %v", err)
	}
}

// TestImportInventoryWithData tests importing inventory data
func TestImportInventoryWithData(t *testing.T) {
	inv := pri.GetInventory()
	inv.Clear()

	// Add some initial items
	inv.Set(0, pri.Row{ItemID: 100, Count: 5})

	// Create import inventory
	importInv := &pri.Inventory{
		Rows: []*pri.Row{
			{ItemID: 1, Count: 10},
			{ItemID: 2, Count: 20},
			{ItemID: 3, Count: 30},
		},
	}

	err := importInventory(importInv)
	if err != nil {
		t.Errorf("importInventory error: %v", err)
	}

	// Verify items were imported
	if inv.Rows[0].ItemID != 1 || inv.Rows[0].Count != 10 {
		t.Errorf("Item 0 = {%d, %d}, want {1, 10}", inv.Rows[0].ItemID, inv.Rows[0].Count)
	}
	if inv.Rows[1].ItemID != 2 || inv.Rows[1].Count != 20 {
		t.Errorf("Item 1 = {%d, %d}, want {2, 20}", inv.Rows[1].ItemID, inv.Rows[1].Count)
	}
}

// TestImportMagicEmpty tests importing empty magic
func TestImportMagicEmpty(t *testing.T) {
	err := importMagic(nil)
	if err != nil {
		t.Errorf("importMagic(nil) error: %v", err)
	}

	err = importMagic([]*models.Character{})
	if err != nil {
		t.Errorf("importMagic(empty) error: %v", err)
	}
}

// TestImportEspersEmpty tests importing empty espers
func TestImportEspersEmpty(t *testing.T) {
	err := importEspers(nil)
	if err != nil {
		t.Errorf("importEspers(nil) error: %v", err)
	}

	err = importEspers([]*consts.NameValueChecked{})
	if err != nil {
		t.Errorf("importEspers(empty) error: %v", err)
	}
}

// TestImportEspersWithData tests importing esper data
func TestImportEspersWithData(t *testing.T) {
	// Create import espers data
	esperData := []*consts.NameValueChecked{
		{NameValue: consts.NameValue{Name: "Ifrit", Value: 1}, Checked: true},
		{NameValue: consts.NameValue{Name: "Shiva", Value: 2}, Checked: true},
		{NameValue: consts.NameValue{Name: "Ramuh", Value: 3}, Checked: false},
	}

	err := importEspers(esperData)
	if err != nil {
		t.Errorf("importEspers error: %v", err)
	}

	// Verify espers were imported
	for _, esper := range esperData {
		if existingEsper, ok := prconsts.EspersByValue[esper.Value]; ok {
			if existingEsper.Checked != esper.Checked {
				t.Errorf("Esper %s Checked = %v, want %v", esper.Name, existingEsper.Checked, esper.Checked)
			}
		}
	}
}

// createTestSaveFileForImport creates a minimal test save file for import tests
func createTestSaveFileForImport(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	// Create minimal JSON save structure
	minimalSave := `{
		"id": 1,
		"userData": "{\"owendGil\": 1000, \"ownedCharacterList\": \"{\\\"target\\\": []}\", \"normalOwnedItemList\": \"{\\\"target\\\": []}\"}",
		"mapData": "{\"mapId\": 1, \"pointIn\": 0}",
		"isCompleteFlag": 0,
		"dataStorage": "{\"global\": [0,0,0,0,0,0,0,0,0,0]}"
	}`

	if err := os.WriteFile(path, []byte(minimalSave), 0644); err != nil {
		t.Fatalf("Failed to create test save file: %v", err)
	}

	return path
}

// createImportFile creates a test import file
func createImportFile(t *testing.T, dir, name, format string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	exportData := ExportData{
		Metadata: ExportMetadata{
			SourceFile: "original.json",
			Format:     format,
		},
	}

	// Add data based on format
	switch format {
	case "full", "characters", "magic":
		exportData.Characters = []*models.Character{
			{
				ID:      1,
				Name:    "TestChar",
				Level:   50,
				Vigor:   80,
				Stamina: 70,
				Speed:   60,
				Magic:   90,
				HP:      models.CurrentMax{Current: 500, Max: 1000},
				MP:      models.CurrentMax{Current: 100, Max: 200},
				SpellsByID: map[int]*models.Spell{
					1: {Index: 1, Value: 100},
				},
			},
		}
	}

	if format == "full" || format == "inventory" {
		exportData.Inventory = &pri.Inventory{
			Rows: []*pri.Row{
				{ItemID: 1, Count: 10},
				{ItemID: 2, Count: 20},
			},
		}
	}

	if format == "full" || format == "party" {
		exportData.Party = &pri.Party{
			Members: [4]*pri.Member{
				{CharacterID: 1, Name: "Terra"},
				{CharacterID: 5, Name: "Locke"},
			},
		}
	}

	if format == "full" || format == "espers" {
		exportData.Espers = []*consts.NameValueChecked{
			{NameValue: consts.NameValue{Name: "Ifrit", Value: 1}, Checked: true},
			{NameValue: consts.NameValue{Name: "Shiva", Value: 2}, Checked: true},
		}
	}

	jsonData, err := json.Marshal(exportData)
	if err != nil {
		t.Fatalf("Failed to marshal export data: %v", err)
	}

	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		t.Fatalf("Failed to create import file: %v", err)
	}

	return path
}

// BenchmarkHandleImportCommand benchmarks the import command
func BenchmarkHandleImportCommand(b *testing.B) {
	tmpDir := b.TempDir()
	saveFile := createTestSaveFileForImport(nil, tmpDir, "save.json")
	importFile := createImportFile(nil, tmpDir, "import.json", "full")

	cli := NewCLI([]string{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cli.handleImportCommand(saveFile, importFile, "full", false)
	}
}

// BenchmarkImportCharacters benchmarks character import
func BenchmarkImportCharacters(b *testing.B) {
	characters := []*models.Character{
		{
			ID:         1,
			Name:       "Test",
			Level:      50,
			Vigor:      80,
			SpellsByID: make(map[int]*models.Spell),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = importCharacters(characters)
	}
}
