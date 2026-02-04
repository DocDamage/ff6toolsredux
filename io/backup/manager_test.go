package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"ffvi_editor/models"
)

// setupTestManager creates a test manager with a temporary directory
func setupTestManager(t *testing.T) (*Manager, string) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	manager, err := NewManager(tmpDir, 10)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create manager: %v", err)
	}

	return manager, tmpDir
}

// cleanup removes the temporary directory
func cleanup(tmpDir string) {
	os.RemoveAll(tmpDir)
}

// TestNewManager tests manager creation
func TestNewManager(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	if manager.backupDir != tmpDir {
		t.Errorf("backupDir = %s, want %s", manager.backupDir, tmpDir)
	}

	if manager.maxBackups != 10 {
		t.Errorf("maxBackups = %d, want 10", manager.maxBackups)
	}

	if !manager.autoBackup {
		t.Error("autoBackup should be true when maxBackups > 0")
	}

	if manager.backups == nil {
		t.Error("backups map is nil")
	}
}

// TestNewManagerCreatesDirectory tests that manager creates backup directory
func TestNewManagerCreatesDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	parentDir := tmpDir
	backupDir := filepath.Join(tmpDir, "backups", "nested")
	defer os.RemoveAll(parentDir)

	_, err = NewManager(backupDir, 10)
	if err != nil {
		t.Fatalf("NewManager() returned error: %v", err)
	}

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		t.Error("NewManager() did not create backup directory")
	}
}

// TestCreateBackup tests backup creation
func TestCreateBackup(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	data := []byte("test save data")
	originalPath := "/path/to/save.sav"
	description := "Test backup"

	meta, err := manager.CreateBackup(originalPath, data, description)
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	if meta == nil {
		t.Fatal("CreateBackup() returned nil metadata")
	}

	if meta.OriginalPath != originalPath {
		t.Errorf("OriginalPath = %s, want %s", meta.OriginalPath, originalPath)
	}

	if meta.Description != description {
		t.Errorf("Description = %s, want %s", meta.Description, description)
	}

	// Verify backup file was created
	backupPath := filepath.Join(tmpDir, meta.ID+".bak")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Backup file was not created")
	}
}

// TestCreateBackupDisabled tests backup creation when disabled
func TestCreateBackupDisabled(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager, err := NewManager(tmpDir, 0) // maxBackups = 0 disables auto-backup
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	data := []byte("test save data")
	_, err = manager.CreateBackup("/path", data, "desc")
	if err == nil {
		t.Error("CreateBackup() should return error when auto-backup is disabled")
	}
}

// TestCreateBackupVerifiesHash tests that backup verifies file hash
func TestCreateBackupVerifiesHash(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	data := []byte("test save data for hash")
	meta, err := manager.CreateBackup("/path", data, "desc")
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	expectedHash := models.CalculateHash(data)
	if meta.Hash != expectedHash {
		t.Errorf("Hash mismatch: got %s, want %s", meta.Hash, expectedHash)
	}
}

// TestRestoreBackup tests backup restoration
func TestRestoreBackup(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	originalData := []byte("original save data")
	meta, err := manager.CreateBackup("/path", originalData, "desc")
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	restoredData, err := manager.RestoreBackup(meta.ID)
	if err != nil {
		t.Fatalf("RestoreBackup() returned error: %v", err)
	}

	if string(restoredData) != string(originalData) {
		t.Error("Restored data does not match original")
	}
}

// TestRestoreBackupNotFound tests restoring non-existent backup
func TestRestoreBackupNotFound(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	_, err := manager.RestoreBackup("non_existent_id")
	if err == nil {
		t.Error("RestoreBackup() should return error for non-existent backup")
	}
}

// TestRestoreBackupIntegrityCheck tests backup integrity verification
func TestRestoreBackupIntegrityCheck(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	originalData := []byte("original save data")
	meta, err := manager.CreateBackup("/path", originalData, "desc")
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	// Corrupt the backup file
	backupPath := filepath.Join(tmpDir, meta.ID+".bak")
	corruptedData := []byte("corrupted data")
	if err := os.WriteFile(backupPath, corruptedData, 0644); err != nil {
		t.Fatalf("Failed to corrupt backup file: %v", err)
	}

	_, err = manager.RestoreBackup(meta.ID)
	if err == nil {
		t.Error("RestoreBackup() should return error for corrupted backup")
	}
}

// TestDeleteBackup tests backup deletion
func TestDeleteBackup(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	meta, err := manager.CreateBackup("/path", []byte("data"), "desc")
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	// Verify backup exists
	backupPath := filepath.Join(tmpDir, meta.ID+".bak")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatal("Backup file was not created")
	}

	// Delete backup
	err = manager.DeleteBackup(meta.ID)
	if err != nil {
		t.Fatalf("DeleteBackup() returned error: %v", err)
	}

	// Verify backup is deleted
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Error("Backup file was not deleted")
	}

	// Verify metadata is removed
	if _, exists := manager.backups[meta.ID]; exists {
		t.Error("Backup metadata was not removed")
	}
}

// TestDeleteBackupNotFound tests deleting non-existent backup
func TestDeleteBackupNotFound(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	err := manager.DeleteBackup("non_existent_id")
	if err == nil {
		t.Error("DeleteBackup() should return error for non-existent backup")
	}
}

// TestGetBackupMetadata tests getting backup metadata
func TestGetBackupMetadata(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	meta, err := manager.CreateBackup("/path", []byte("data"), "desc")
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	retrievedMeta, err := manager.GetBackupMetadata(meta.ID)
	if err != nil {
		t.Fatalf("GetBackupMetadata() returned error: %v", err)
	}

	if retrievedMeta.ID != meta.ID {
		t.Errorf("Retrieved metadata ID = %s, want %s", retrievedMeta.ID, meta.ID)
	}
}

// TestGetBackupMetadataNotFound tests getting non-existent backup metadata
func TestGetBackupMetadataNotFound(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	_, err := manager.GetBackupMetadata("non_existent_id")
	if err == nil {
		t.Error("GetBackupMetadata() should return error for non-existent backup")
	}
}

// TestSetMaxBackups tests setting max backups
func TestSetMaxBackups(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager, err := NewManager(tmpDir, 10)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	manager.SetMaxBackups(5)
	if manager.GetMaxBackups() != 5 {
		t.Errorf("maxBackups = %d, want 5", manager.GetMaxBackups())
	}

	// Setting to 0 should disable auto-backup
	manager.SetMaxBackups(0)
	if manager.GetMaxBackups() != 0 {
		t.Errorf("maxBackups = %d, want 0", manager.GetMaxBackups())
	}
	if manager.autoBackup {
		t.Error("autoBackup should be false when maxBackups = 0")
	}
}

// TestBackupCount tests backup count
func TestBackupCount(t *testing.T) {
	manager, tmpDir := setupTestManager(t)
	defer cleanup(tmpDir)

	if manager.BackupCount() != 0 {
		t.Errorf("Initial backup count = %d, want 0", manager.BackupCount())
	}

	_, err := manager.CreateBackup("/path", []byte("data1"), "desc")
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}
	time.Sleep(20 * time.Millisecond)

	if manager.BackupCount() != 1 {
		t.Errorf("Backup count = %d, want 1", manager.BackupCount())
	}
}

// TestCleanupOldBackups tests old backup cleanup
func TestCleanupOldBackups(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager, err := NewManager(tmpDir, 2) // Max 2 backups
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create 3 backups with delays to ensure unique IDs
	backupIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		data := []byte("data" + string(rune(i)))
		meta, err := manager.CreateBackup("/path", data, "desc")
		if err != nil {
			t.Fatalf("Failed to create backup: %v", err)
		}
		backupIDs[i] = meta.ID
		time.Sleep(25 * time.Millisecond)
	}

	// Should only have 2 backups (the newest ones)
	if manager.BackupCount() > 2 {
		t.Errorf("Backup count = %d, want at most 2", manager.BackupCount())
	}
}

// TestFormatTimeSince tests time formatting
func TestFormatTimeSince(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "just now"},
		{1 * time.Minute, "1 minute ago"},
		{5 * time.Minute, "5 minutes ago"},
		{1 * time.Hour, "1 hour ago"},
		{3 * time.Hour, "3 hours ago"},
		{1 * 24 * time.Hour, "1 day ago"},
		{5 * 24 * time.Hour, "5 days ago"},
		{14 * 24 * time.Hour, "14 days ago"},
	}

	for _, tt := range tests {
		result := formatTimeSince(tt.duration)
		if result != tt.expected {
			t.Errorf("formatTimeSince(%v) = %s, want %s", tt.duration, result, tt.expected)
		}
	}
}

// TestLoadMetadata tests metadata loading
func TestLoadMetadata(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create manager and add backup
	manager1, err := NewManager(tmpDir, 10)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	meta, err := manager1.CreateBackup("/path", []byte("data"), "desc")
	if err != nil {
		t.Fatalf("CreateBackup() returned error: %v", err)
	}

	// Create new manager pointing to same directory (simulating restart)
	manager2, err := NewManager(tmpDir, 10)
	if err != nil {
		t.Fatalf("Failed to create second manager: %v", err)
	}

	// Should be able to get the backup from new manager
	retrievedMeta, err := manager2.GetBackupMetadata(meta.ID)
	if err != nil {
		t.Fatalf("GetBackupMetadata() returned error: %v", err)
	}

	if retrievedMeta.ID != meta.ID {
		t.Errorf("Retrieved metadata ID = %s, want %s", retrievedMeta.ID, meta.ID)
	}
}

// BenchmarkCreateBackup benchmarks backup creation
func BenchmarkCreateBackup(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "backup_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager, err := NewManager(tmpDir, 100)
	if err != nil {
		b.Fatalf("Failed to create manager: %v", err)
	}

	data := make([]byte, 1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := manager.CreateBackup("/path", data, "desc")
		if err != nil {
			b.Fatalf("CreateBackup() returned error: %v", err)
		}
	}
}

// BenchmarkRestoreBackup benchmarks backup restoration
func BenchmarkRestoreBackup(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "backup_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager, err := NewManager(tmpDir, 100)
	if err != nil {
		b.Fatalf("Failed to create manager: %v", err)
	}

	data := make([]byte, 1024)
	meta, err := manager.CreateBackup("/path", data, "desc")
	if err != nil {
		b.Fatalf("CreateBackup() returned error: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := manager.RestoreBackup(meta.ID)
		if err != nil {
			b.Fatalf("RestoreBackup() returned error: %v", err)
		}
	}
}
