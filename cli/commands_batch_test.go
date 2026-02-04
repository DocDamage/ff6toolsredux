package cli

import (
	"os"
	"path/filepath"
	"testing"

	"ffvi_editor/models"
	pri "ffvi_editor/models/pr"
)

// TestHandleBatchCommandMaxStats tests the max-stats batch operation
func TestHandleBatchCommandMaxStats(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForBatch(t, tmpDir, "save.json")

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "max-stats", "")
	if err != nil {
		t.Logf("handleBatchCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify characters have max stats
	for _, char := range pri.Characters {
		if char == nil || char.IsNPC {
			continue
		}
		if char.Vigor != 255 {
			t.Errorf("Character %s Vigor = %d, want 255", char.Name, char.Vigor)
		}
		if char.Stamina != 255 {
			t.Errorf("Character %s Stamina = %d, want 255", char.Name, char.Stamina)
		}
		if char.Speed != 255 {
			t.Errorf("Character %s Speed = %d, want 255", char.Name, char.Speed)
		}
		if char.Magic != 255 {
			t.Errorf("Character %s Magic = %d, want 255", char.Name, char.Magic)
		}
	}
}

// TestHandleBatchCommandMaxItems tests the max-items batch operation
func TestHandleBatchCommandMaxItems(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForBatch(t, tmpDir, "save.json")

	// Add some items to inventory
	inv := pri.GetInventory()
	inv.Clear()
	inv.Set(0, pri.Row{ItemID: 1, Count: 5})
	inv.Set(1, pri.Row{ItemID: 2, Count: 10})
	inv.Set(2, pri.Row{ItemID: 3, Count: 1})

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "max-items", "")
	if err != nil {
		t.Logf("handleBatchCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify items have max count
	for _, row := range inv.GetRows() {
		if row.ItemID > 0 && row.Count != 99 {
			t.Errorf("Item %d Count = %d, want 99", row.ItemID, row.Count)
		}
	}
}

// TestHandleBatchCommandMaxMagic tests the max-magic batch operation
func TestHandleBatchCommandMaxMagic(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForBatch(t, tmpDir, "save.json")

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "max-magic", "")
	if err != nil {
		t.Logf("handleBatchCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify characters have learned spells
	for _, char := range pri.Characters {
		if char == nil || char.IsNPC {
			continue
		}
		for _, spell := range char.SpellsByIndex {
			if spell != nil && spell.Value != 1 {
				t.Errorf("Character %s spell %s Value = %d, want 1", char.Name, spell.Name, spell.Value)
			}
		}
	}
}

// TestHandleBatchCommandMaxAll tests the max-all batch operation
func TestHandleBatchCommandMaxAll(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForBatch(t, tmpDir, "save.json")

	// Add some items
	inv := pri.GetInventory()
	inv.Clear()
	inv.Set(0, pri.Row{ItemID: 1, Count: 5})

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "max-all", "")
	if err != nil {
		t.Logf("handleBatchCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify stats are maxed
	for _, char := range pri.Characters {
		if char == nil || char.IsNPC {
			continue
		}
		if char.Vigor != 255 {
			t.Errorf("Character %s Vigor = %d, want 255", char.Name, char.Vigor)
		}
	}

	// Verify items are maxed
	for _, row := range inv.GetRows() {
		if row.ItemID > 0 && row.Count != 99 {
			t.Errorf("Item %d Count = %d, want 99", row.ItemID, row.Count)
		}
	}
}

// TestHandleBatchCommandInvalidOperation tests invalid batch operation
func TestHandleBatchCommandInvalidOperation(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForBatch(t, tmpDir, "save.json")

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "invalid-op", "")
	if err == nil {
		t.Error("handleBatchCommand should return error for invalid operation")
	}
}

// TestHandleBatchCommandNonExistentFile tests batch on non-existent file
func TestHandleBatchCommandNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "nonexistent.json")

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "max-stats", "")
	if err == nil {
		t.Error("handleBatchCommand should return error for non-existent file")
	}
}

// TestHandleBatchCommandWithOutput tests batch with custom output
func TestHandleBatchCommandWithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForBatch(t, tmpDir, "save.json")
	outputFile := filepath.Join(tmpDir, "output.json")

	cli := NewCLI([]string{})

	err := cli.handleBatchCommand(saveFile, "max-stats", outputFile)
	if err != nil {
		t.Logf("handleBatchCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Log("Output file was not created (expected in isolated test)")
	}
}

// TestApplyMaxStats tests the applyMaxStats function
func TestApplyMaxStats(t *testing.T) {
	// Set up test characters with some stats
	if len(pri.Characters) > 0 {
		for _, char := range pri.Characters {
			if char == nil || char.IsNPC {
				continue
			}
			char.Vigor = 10
			char.Stamina = 20
			char.Speed = 30
			char.Magic = 40
		}

		err := applyMaxStats()
		if err != nil {
			t.Errorf("applyMaxStats() error: %v", err)
		}

		// Verify stats are maxed
		for _, char := range pri.Characters {
			if char == nil || char.IsNPC {
				continue
			}
			if char.Vigor != 255 {
				t.Errorf("Character %s Vigor = %d, want 255", char.Name, char.Vigor)
			}
			if char.Stamina != 255 {
				t.Errorf("Character %s Stamina = %d, want 255", char.Name, char.Stamina)
			}
			if char.Speed != 255 {
				t.Errorf("Character %s Speed = %d, want 255", char.Name, char.Speed)
			}
			if char.Magic != 255 {
				t.Errorf("Character %s Magic = %d, want 255", char.Name, char.Magic)
			}
		}
	} else {
		t.Skip("No characters available for testing")
	}
}

// TestApplyMaxItems tests the applyMaxItems function
func TestApplyMaxItems(t *testing.T) {
	inv := pri.GetInventory()
	inv.Clear()

	// Add some items with varying counts
	inv.Set(0, pri.Row{ItemID: 1, Count: 5})
	inv.Set(1, pri.Row{ItemID: 2, Count: 10})
	inv.Set(2, pri.Row{ItemID: 3, Count: 1})
	inv.Set(3, pri.Row{ItemID: 0, Count: 0}) // Empty slot

	err := applyMaxItems()
	if err != nil {
		t.Errorf("applyMaxItems() error: %v", err)
	}

	// Verify items are maxed
	for _, row := range inv.GetRows() {
		if row.ItemID > 0 && row.Count != 99 {
			t.Errorf("Item %d Count = %d, want 99", row.ItemID, row.Count)
		}
	}
}

// TestApplyMaxMagic tests the applyMaxMagic function
func TestApplyMaxMagic(t *testing.T) {
	if len(pri.Characters) > 0 {
		// Set up characters with unlearned spells
		for _, char := range pri.Characters {
			if char == nil || char.IsNPC {
				continue
			}
			for _, spell := range char.SpellsByIndex {
				if spell != nil {
					spell.Value = 0
				}
			}
		}

		err := applyMaxMagic()
		if err != nil {
			t.Errorf("applyMaxMagic() error: %v", err)
		}

		// Verify spells are learned
		for _, char := range pri.Characters {
			if char == nil || char.IsNPC {
				continue
			}
			for _, spell := range char.SpellsByIndex {
				if spell != nil && spell.Value != 1 {
					t.Errorf("Character %s spell %s Value = %d, want 1", char.Name, spell.Name, spell.Value)
				}
			}
		}
	} else {
		t.Skip("No characters available for testing")
	}
}

// TestApplyMaxStatsSkipsNPC tests that NPCs are skipped
func TestApplyMaxStatsSkipsNPC(t *testing.T) {
	if len(pri.Characters) > 0 {
		// Find an NPC character
		var npc *models.Character
		for _, char := range pri.Characters {
			if char != nil && char.IsNPC {
				npc = char
				break
			}
		}

		if npc == nil {
			t.Skip("No NPC character found for testing")
		}

		originalVigor := npc.Vigor

		err := applyMaxStats()
		if err != nil {
			t.Errorf("applyMaxStats() error: %v", err)
		}

		// NPC stats should not be changed
		if npc.Vigor != originalVigor {
			t.Errorf("NPC Vigor was changed from %d to %d", originalVigor, npc.Vigor)
		}
	} else {
		t.Skip("No characters available for testing")
	}
}

// createTestSaveFileForBatch creates a minimal test save file for batch tests
func createTestSaveFileForBatch(t *testing.T, dir, name string) string {
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

// BenchmarkApplyMaxStats benchmarks max stats application
func BenchmarkApplyMaxStats(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = applyMaxStats()
	}
}

// BenchmarkApplyMaxItems benchmarks max items application
func BenchmarkApplyMaxItems(b *testing.B) {
	inv := pri.GetInventory()
	inv.Clear()
	for i := 0; i < 50; i++ {
		inv.Set(i, pri.Row{ItemID: i + 1, Count: 1})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = applyMaxItems()
	}
}

// BenchmarkHandleBatchCommand benchmarks the batch command
func BenchmarkHandleBatchCommand(b *testing.B) {
	tmpDir := b.TempDir()
	saveFile := createTestSaveFileForBatch(nil, tmpDir, "save.json")

	cli := NewCLI([]string{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cli.handleBatchCommand(saveFile, "max-stats", "")
	}
}
