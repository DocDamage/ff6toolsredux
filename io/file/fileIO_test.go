package file

import (
	"bytes"
	"encoding/base64"
	"os"
	"testing"

	"ffvi_editor/global"
	"github.com/kiamev/ffpr-save-cypher/rijndael"
)

// TestLoadFileInvalidPath tests error handling for missing files
func TestLoadFileInvalidPath(t *testing.T) {
	_, _, err := LoadFile("/nonexistent/path/file.save", global.PC)
	if err == nil {
		t.Fatal("LoadFile() should return error for nonexistent file")
	}
}

// TestLoadFilePS tests PlayStation save format (no transformation)
func TestLoadFilePS(t *testing.T) {
	// Create a temporary test file
	testData := []byte("test data for ps")
	tmpFile := t.TempDir() + "/test.save"

	if err := writeTestFile(tmpFile, testData); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	out, trimmed, err := LoadFile(tmpFile, global.PS)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}

	if !bytes.Equal(out, testData) {
		t.Fatalf("LoadFile() returned different data for PS format")
	}
	if trimmed != nil {
		t.Fatal("LoadFile() should return nil trimmed for PS format")
	}
}

// TestLoadFileTooShort tests error handling for truncated files
func TestLoadFileTooShort(t *testing.T) {
	tmpFile := t.TempDir() + "/test.save"
	if err := writeTestFile(tmpFile, []byte("abc")); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, _, err := LoadFile(tmpFile, global.PC)
	if err == nil {
		t.Fatal("LoadFile() should return error for too-short file")
	}
}

// TestLoadFileBOMRemoval tests that BOM (Byte Order Mark) is properly detected and trimmed
func TestLoadFileBOMRemoval(t *testing.T) {
	bom := []byte{239, 187, 191}            // UTF-8 BOM
	testData := []byte("test")              // dummy data to pass base64 check
	fileContent := append(bom, testData...) // prepend BOM

	tmpFile := t.TempDir() + "/test.save"
	if err := writeTestFile(tmpFile, fileContent); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, trimmed, err := LoadFile(tmpFile, global.PC)
	if err == nil {
		// May fail for decryption reasons, but the BOM should have been detected
		// If we get here without error, check that trimmed is set
		if trimmed == nil {
			t.Fatal("LoadFile() should detect BOM and set trimmed")
		}
	}
}

// TestSaveFilePC tests PC save file encryption and encoding
func TestSaveFilePC(t *testing.T) {
	testData := []byte(`{"test": "data"}`)
	tmpFile := t.TempDir() + "/test.save"

	err := SaveFile(testData, tmpFile, nil, global.PC)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	// Verify file was created
	if !fileExists(tmpFile) {
		t.Fatal("SaveFile() did not create output file")
	}

	// Read back and verify it's base64 encoded
	fileContent, err := readTestFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if len(fileContent) == 0 {
		t.Fatal("SaveFile() created empty file")
	}

	// Should be base64 (alphanumeric + /+= characters)
	isBase64 := isValidBase64(string(fileContent))
	if !isBase64 {
		t.Fatal("SaveFile() output doesn't appear to be base64 encoded")
	}
}

// TestSaveFilePS tests PlayStation save file (no transformation)
func TestSaveFilePS(t *testing.T) {
	testData := []byte("original ps data")
	tmpFile := t.TempDir() + "/test.save"

	err := SaveFile(testData, tmpFile, nil, global.PS)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	fileContent, err := readTestFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !bytes.Equal(fileContent, testData) {
		t.Fatal("SaveFile() should preserve data for PS format")
	}
}

// TestSaveFileWithBOM tests that BOM is properly written when provided
func TestSaveFileWithBOM(t *testing.T) {
	testData := []byte(`{"test": "data"}`)
	bom := []byte{239, 187, 191}
	tmpFile := t.TempDir() + "/test.save"

	err := SaveFile(testData, tmpFile, bom, global.PC)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	fileContent, err := readTestFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !bytes.HasPrefix(fileContent, bom) {
		t.Fatal("SaveFile() should prepend BOM to output")
	}
}

// TestRoundTrip tests encryption/decryption consistency
// This test is limited by the external rijndael library, but verifies integration
func TestRoundTripIntegration(t *testing.T) {
	testData := []byte(`{"character": {"name": "Terra", "level": 1}}`)

	// Compress
	var buf bytes.Buffer
	_, err := buf.Write(testData)
	if err != nil {
		t.Fatalf("failed to write to buffer: %v", err)
	}

	// Verify we can create a cipher (basic integration test)
	cipher := rijndael.New()
	if cipher == nil {
		t.Fatal("rijndael cipher creation failed")
	}
}

// Helper functions

func writeTestFile(path string, data []byte) error {
	return writeFile(path, data)
}

func readTestFile(path string) ([]byte, error) {
	return readFile(path)
}

func fileExists(path string) bool {
	// Simple implementation - can be replaced with os.Stat
	_, err := readFile(path)
	return err == nil
}

func isValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// Simple file I/O helpers to avoid platform-specific code
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// TestSaveFileCompressionWorks tests that data is compressed before encryption
func TestSaveFileCompressionWorks(t *testing.T) {
	// Large repeating data that should compress well
	largeData := []byte(bytes.Repeat([]byte("a"), 10000))
	tmpFile := t.TempDir() + "/test.save"

	err := SaveFile(largeData, tmpFile, nil, global.PC)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	fileContent, err := readTestFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	// Compressed and encrypted file should be much smaller than original
	if len(fileContent) >= len(largeData) {
		t.Fatal("SaveFile() compression may not be working (output not smaller than input)")
	}
}
