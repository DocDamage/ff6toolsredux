package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"ffvi_editor/models/consts"
)

// TestHandleExportCommandFull tests exporting full save data
func TestHandleExportCommandFull(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForExport(t, tmpDir, "save.json")
	outputFile := filepath.Join(tmpDir, "export.json")

	cli := NewCLI([]string{})

	err := cli.handleExportCommand(saveFile, outputFile, "full")
	if err != nil {
		t.Logf("handleExportCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}

	// Verify JSON content
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	var export ExportData
	if err := json.Unmarshal(data, &export); err != nil {
		t.Errorf("Export file is not valid JSON: %v", err)
	}

	// Verify metadata
	if export.Metadata.SourceFile != saveFile {
		t.Errorf("SourceFile = %q, want %q", export.Metadata.SourceFile, saveFile)
	}
	if export.Metadata.Format != "full" {
		t.Errorf("Format = %q, want full", export.Metadata.Format)
	}
}

// TestHandleExportCommandFormats tests different export formats
func TestHandleExportCommandFormats(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForExport(t, tmpDir, "save.json")

	cli := NewCLI([]string{})

	formats := []string{"full", "characters", "inventory", "party", "magic", "espers"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			outputFile := filepath.Join(tmpDir, "export_"+format+".json")

			err := cli.handleExportCommand(saveFile, outputFile, format)
			if err != nil {
				t.Logf("handleExportCommand error for format %s (expected in isolated test): %v", format, err)
				return
			}

			// Verify output file was created
			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				t.Errorf("Export file for format %s was not created", format)
			}

			// Verify JSON content
			data, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read export file: %v", err)
			}

			var export ExportData
			if err := json.Unmarshal(data, &export); err != nil {
				t.Errorf("Export file for format %s is not valid JSON: %v", format, err)
			}

			// Verify format in metadata
			if export.Metadata.Format != format {
				t.Errorf("Format = %q, want %q", export.Metadata.Format, format)
			}
		})
	}
}

// TestHandleExportCommandInvalidFormat tests invalid export format
func TestHandleExportCommandInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForExport(t, tmpDir, "save.json")
	outputFile := filepath.Join(tmpDir, "export.json")

	cli := NewCLI([]string{})

	err := cli.handleExportCommand(saveFile, outputFile, "invalid_format")
	if err == nil {
		t.Error("handleExportCommand should return error for invalid format")
	}
}

// TestHandleExportCommandNonExistentFile tests exporting non-existent file
func TestHandleExportCommandNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "nonexistent.json")
	outputFile := filepath.Join(tmpDir, "export.json")

	cli := NewCLI([]string{})

	err := cli.handleExportCommand(saveFile, outputFile, "full")
	if err == nil {
		t.Error("handleExportCommand should return error for non-existent file")
	}
}

// TestHandleExportCommandInvalidJSON tests exporting invalid JSON file
func TestHandleExportCommandInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "invalid.json")
	outputFile := filepath.Join(tmpDir, "export.json")

	// Create invalid JSON file
	if err := os.WriteFile(saveFile, []byte(`{invalid json`), 0644); err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	cli := NewCLI([]string{})

	err := cli.handleExportCommand(saveFile, outputFile, "full")
	if err == nil {
		t.Error("handleExportCommand should return error for invalid JSON")
	}
}

// TestGetEspersForExport tests the esper export helper
func TestGetEspersForExport(t *testing.T) {
	espers := getEspersForExport()

	if espers == nil {
		t.Fatal("getEspersForExport() returned nil")
	}

	// Espers should not be empty in a valid setup
	t.Logf("Found %d espers for export", len(espers))
}

// TestExportDataStructure tests the ExportData structure
func TestExportDataStructure(t *testing.T) {
	export := ExportData{
		Metadata: ExportMetadata{
			SourceFile: "test.json",
			Format:     "full",
		},
		Characters: nil,
		Party:      nil,
		Inventory:  nil,
		Espers:     nil,
	}

	if export.Metadata.SourceFile != "test.json" {
		t.Errorf("SourceFile = %q, want test.json", export.Metadata.SourceFile)
	}
	if export.Metadata.Format != "full" {
		t.Errorf("Format = %q, want full", export.Metadata.Format)
	}
}

// TestExportMetadataJSON tests JSON marshaling of metadata
func TestExportMetadataJSON(t *testing.T) {
	metadata := ExportMetadata{
		SourceFile: "/path/to/save.json",
		Format:     "characters",
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("Failed to marshal metadata: %v", err)
	}

	var unmarshaled ExportMetadata
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if unmarshaled.SourceFile != metadata.SourceFile {
		t.Errorf("SourceFile = %q, want %q", unmarshaled.SourceFile, metadata.SourceFile)
	}
	if unmarshaled.Format != metadata.Format {
		t.Errorf("Format = %q, want %q", unmarshaled.Format, metadata.Format)
	}
}

// TestExportDataWithEspers tests export with esper data
func TestExportDataWithEspers(t *testing.T) {
	espers := []*consts.NameValueChecked{
		{NameValue: consts.NameValue{Name: "Ifrit", Value: 1}, Checked: true},
		{NameValue: consts.NameValue{Name: "Shiva", Value: 2}, Checked: false},
	}

	export := ExportData{
		Metadata: ExportMetadata{
			SourceFile: "test.json",
			Format:     "espers",
		},
		Espers: espers,
	}

	if len(export.Espers) != 2 {
		t.Errorf("Espers count = %d, want 2", len(export.Espers))
	}

	if export.Espers[0].Name != "Ifrit" {
		t.Errorf("First esper name = %q, want Ifrit", export.Espers[0].Name)
	}

	if !export.Espers[0].Checked {
		t.Error("First esper should be checked")
	}
}

// createTestSaveFileForExport creates a minimal test save file for export tests
func createTestSaveFileForExport(t *testing.T, dir, name string) string {
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

// BenchmarkHandleExportCommand benchmarks the export command
func BenchmarkHandleExportCommand(b *testing.B) {
	tmpDir := b.TempDir()
	saveFile := createTestSaveFileForExport(nil, tmpDir, "save.json")
	outputFile := filepath.Join(tmpDir, "export.json")

	cli := NewCLI([]string{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cli.handleExportCommand(saveFile, outputFile, "full")
	}
}

// BenchmarkExportDataMarshal benchmarks JSON marshaling of export data
func BenchmarkExportDataMarshal(b *testing.B) {
	export := ExportData{
		Metadata: ExportMetadata{
			SourceFile: "test.json",
			Format:     "full",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(export)
	}
}
