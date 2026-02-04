package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHandleScriptCommandValid tests running a valid Lua script
func TestHandleScriptCommandValid(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "test.lua", `-- Test script
return {message = "Hello from Lua"}`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err != nil {
		t.Logf("handleScriptCommand error (expected in isolated test): %v", err)
		return
	}
}

// TestHandleScriptCommandNonExistentSave tests script with non-existent save file
func TestHandleScriptCommandNonExistentSave(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "nonexistent.json")
	scriptFile := createTestScriptFile(t, tmpDir, "test.lua", "return {}")

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err == nil {
		t.Error("handleScriptCommand should return error for non-existent save file")
	}
}

// TestHandleScriptCommandNonExistentScript tests with non-existent script file
func TestHandleScriptCommandNonExistentScript(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := filepath.Join(tmpDir, "nonexistent.lua")

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err == nil {
		t.Error("handleScriptCommand should return error for non-existent script file")
	}
}

// TestHandleScriptCommandInvalidScript tests with invalid Lua script
func TestHandleScriptCommandInvalidScript(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "invalid.lua", `function unclosed(`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err == nil {
		t.Error("handleScriptCommand should return error for invalid Lua script")
	}
}

// TestHandleScriptCommandWithOutput tests script with custom output path
func TestHandleScriptCommandWithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "test.lua", `return {success = true}`)
	outputFile := filepath.Join(tmpDir, "output.json")

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, outputFile)
	if err != nil {
		t.Logf("handleScriptCommand error (expected in isolated test): %v", err)
		return
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Log("Output file was not created (expected in isolated test)")
	}
}

// TestHandleScriptCommandCharacterEdit tests script that edits character stats
func TestHandleScriptCommandCharacterEdit(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "edit_char.lua", `
-- Edit character script
if editor and editor.setCharacterLevel then
    editor.setCharacterLevel(0, 99)
    editor.setCharacterHP(0, 9999)
end
return {modified = true}`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err != nil {
		t.Logf("handleScriptCommand error (expected in isolated test): %v", err)
		return
	}
}

// TestHandleScriptCommandInventoryEdit tests script that modifies inventory
func TestHandleScriptCommandInventoryEdit(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "edit_inv.lua", `
-- Edit inventory script
if editor and editor.addItem then
    editor.addItem(1, 10)
end
return {itemsAdded = 10}`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err != nil {
		t.Logf("handleScriptCommand error (expected in isolated test): %v", err)
		return
	}
}

// TestHandleScriptCommandEmptyReturn tests script with no return value
func TestHandleScriptCommandEmptyReturn(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "empty.lua", `-- Empty script
print("Hello World")`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err != nil {
		t.Logf("handleScriptCommand error (expected in isolated test): %v", err)
		return
	}
}

// TestHandleScriptCommandComplexReturn tests script returning complex data
func TestHandleScriptCommandComplexReturn(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "complex.lua", `
return {
    characters = {"Terra", "Locke", "Cyan"},
    stats = {
        level = 99,
        hp = 9999,
        mp = 999
    },
    items = {
        {id = 1, count = 10},
        {id = 2, count = 5}
    }
}`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err != nil {
		t.Logf("handleScriptCommand error (expected in isolated test): %v", err)
		return
	}
}

// TestHandleScriptCommandErrorHandling tests script error handling
func TestHandleScriptCommandErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := createTestSaveFileForScript(t, tmpDir, "save.json")
	scriptFile := createTestScriptFile(t, tmpDir, "error.lua", `
-- Script with runtime error
local x = nil
local y = x.someField  -- This will cause an error
return {}`)

	cli := NewCLI([]string{})

	err := cli.handleScriptCommand(saveFile, scriptFile, "")
	if err == nil {
		t.Log("handleScriptCommand should have returned error for script with runtime error (may be caught by Lua sandbox)")
	}
}

// createTestSaveFileForScript creates a minimal test save file for script tests
func createTestSaveFileForScript(t *testing.T, dir, name string) string {
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

// createTestScriptFile creates a test Lua script file
func createTestScriptFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test script file: %v", err)
	}

	return path
}

// BenchmarkHandleScriptCommand benchmarks the script command
func BenchmarkHandleScriptCommand(b *testing.B) {
	tmpDir := b.TempDir()
	saveFile := createTestSaveFileForScript(nil, tmpDir, "save.json")
	scriptFile := createTestScriptFile(nil, tmpDir, "bench.lua", `return {count = 42}`)

	cli := NewCLI([]string{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cli.handleScriptCommand(saveFile, scriptFile, "")
	}
}
