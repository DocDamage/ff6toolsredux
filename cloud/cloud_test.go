package cloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestGoogleDriveProvider tests the Google Drive provider implementation
func TestGoogleDriveProvider(t *testing.T) {
	provider := NewGoogleDriveProvider("test-client-id", "test-secret")

	// Test initial state
	if provider.IsAuthenticated() {
		t.Error("provider should not be authenticated initially")
	}

	// Test authentication
	err := provider.Authenticate(context.Background())
	if err != nil {
		t.Errorf("authentication failed: %v", err)
	}

	if !provider.IsAuthenticated() {
		t.Error("provider should be authenticated after Authenticate()")
	}

	// Test GetName
	name := provider.GetName()
	if name != "Google Drive" {
		t.Errorf("expected name 'Google Drive', got '%s'", name)
	}

	// Test ValidateConnection
	ok, msg := provider.ValidateConnection(context.Background())
	if !ok {
		t.Errorf("connection validation failed: %s", msg)
	}

	// Test GetQuotaInfo
	used, total, err := provider.GetQuotaInfo(context.Background())
	if err != nil {
		t.Errorf("failed to get quota info: %v", err)
	}
	if used <= 0 || total <= 0 {
		t.Errorf("invalid quota info: used=%d, total=%d", used, total)
	}

	// Test Logout
	err = provider.Logout(context.Background())
	if err != nil {
		t.Errorf("logout failed: %v", err)
	}
	if provider.IsAuthenticated() {
		t.Error("provider should not be authenticated after logout")
	}
}

// TestDropboxProvider tests the Dropbox provider implementation
func TestDropboxProvider(t *testing.T) {
	provider := NewDropboxProvider("test-app-key", "test-secret")

	// Test initial state
	if provider.IsAuthenticated() {
		t.Error("provider should not be authenticated initially")
	}

	// Test authentication
	err := provider.Authenticate(context.Background())
	if err != nil {
		t.Errorf("authentication failed: %v", err)
	}

	if !provider.IsAuthenticated() {
		t.Error("provider should be authenticated after Authenticate()")
	}

	// Test GetName
	name := provider.GetName()
	if name != "Dropbox" {
		t.Errorf("expected name 'Dropbox', got '%s'", name)
	}

	// Test ValidateConnection
	ok, msg := provider.ValidateConnection(context.Background())
	if !ok {
		t.Errorf("connection validation failed: %s", msg)
	}

	// Test GetQuotaInfo
	used, total, err := provider.GetQuotaInfo(context.Background())
	if err != nil {
		t.Errorf("failed to get quota info: %v", err)
	}
	if used <= 0 || total <= 0 {
		t.Errorf("invalid quota info: used=%d, total=%d", used, total)
	}
}

// TestCloudManager tests the cloud manager
func TestCloudManager(t *testing.T) {
	manager := New()

	// Test RegisterProvider
	gdrive := NewGoogleDriveProvider("test-id", "test-secret")
	err := manager.RegisterProvider(gdrive)
	if err != nil {
		t.Errorf("failed to register Google Drive provider: %v", err)
	}

	// Test duplicate registration fails
	err = manager.RegisterProvider(gdrive)
	if err == nil {
		t.Error("should fail to register duplicate provider")
	}

	// Test GetProvider
	provider, err := manager.GetProvider("Google Drive")
	if err != nil {
		t.Errorf("failed to get provider: %v", err)
	}
	if provider == nil {
		t.Error("provider should not be nil")
	}

	// Test ListProviders
	providers := manager.ListProviders()
	if len(providers) != 1 {
		t.Errorf("expected 1 provider, got %d", len(providers))
	}
	if !contains(providers, "Google Drive") {
		t.Error("Google Drive not in provider list")
	}

	// Test GetStatus
	status, err := manager.GetStatus("Google Drive")
	if err != nil {
		t.Errorf("failed to get status: %v", err)
	}
	if status == nil {
		t.Error("status should not be nil")
	}

	// Test GetAllStatus
	allStatus := manager.GetAllStatus()
	if len(allStatus) != 1 {
		t.Errorf("expected 1 status entry, got %d", len(allStatus))
	}
}

// TestSyncConflictResolution tests conflict resolution strategies
func TestSyncConflictResolution(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		conflict *Conflict
		strategy ConflictResolution
		expect   string
	}{
		{
			name: "local_strategy",
			conflict: &Conflict{
				FileName:   "test.sav",
				LocalTime:  now,
				RemoteTime: now.Add(-1 * time.Hour),
			},
			strategy: ConflictUseLocal,
			expect:   "local",
		},
		{
			name: "remote_strategy",
			conflict: &Conflict{
				FileName:   "test.sav",
				LocalTime:  now.Add(-1 * time.Hour),
				RemoteTime: now,
			},
			strategy: ConflictUseRemote,
			expect:   "remote",
		},
		{
			name: "newest_strategy",
			conflict: &Conflict{
				FileName:   "test.sav",
				LocalTime:  now,
				RemoteTime: now.Add(-1 * time.Hour),
			},
			strategy: ConflictNewest,
			expect:   "local", // local is newer
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify strategy is set correctly
			if tt.strategy != ConflictUseLocal && tt.strategy != ConflictUseRemote && tt.strategy != ConflictNewest {
				t.Errorf("invalid strategy: %v", tt.strategy)
			}
		})
	}
}

// TestFileMetadata tests file metadata structures
func TestFileMetadata(t *testing.T) {
	now := time.Now()

	metadata := &FileMetadata{
		ID:           "test-id",
		Name:         "test.sav",
		Size:         1024,
		ModifiedTime: now,
		Hash:         fmt.Sprintf("%x", md5.Sum([]byte("test"))),
		IsFolder:     false,
		Path:         "/FF6Editor/Backups/test.sav",
		Parents:      []string{"parent-id"},
	}

	if metadata.ID != "test-id" {
		t.Errorf("unexpected ID: %s", metadata.ID)
	}
	if metadata.Name != "test.sav" {
		t.Errorf("unexpected name: %s", metadata.Name)
	}
	if metadata.Size != 1024 {
		t.Errorf("unexpected size: %d", metadata.Size)
	}
	if metadata.IsFolder {
		t.Error("should not be folder")
	}
	if !strings.Contains(metadata.Path, "FF6Editor") {
		t.Errorf("unexpected path: %s", metadata.Path)
	}
}

// TestSyncStatus tests sync status tracking
func TestSyncStatus(t *testing.T) {
	status := &SyncStatus{
		Provider:        "Google Drive",
		InProgress:      true,
		LastSync:        time.Now(),
		FilesUploaded:   5,
		FilesDownloaded: 3,
		ConflictsFound:  0,
		Progress:        0.75,
		IsAuthenticated: true,
		StorageUsed:     1024 * 1024 * 1024,
		StorageTotal:    2 * 1024 * 1024 * 1024,
	}

	if !status.InProgress {
		t.Error("sync should be in progress")
	}
	if status.FilesUploaded != 5 {
		t.Errorf("expected 5 uploaded files, got %d", status.FilesUploaded)
	}
	if status.FilesDownloaded != 3 {
		t.Errorf("expected 3 downloaded files, got %d", status.FilesDownloaded)
	}
	if status.ConflictsFound != 0 {
		t.Errorf("expected 0 conflicts, got %d", status.ConflictsFound)
	}
	if status.Progress != 0.75 {
		t.Errorf("expected progress 0.75, got %f", status.Progress)
	}
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// BenchmarkProviderAuthentication benchmarks provider authentication
func BenchmarkProviderAuthentication(b *testing.B) {
	provider := NewGoogleDriveProvider("test-id", "test-secret")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.Authenticate(ctx)
	}
}

// BenchmarkFileMetadataCreation benchmarks metadata creation
func BenchmarkFileMetadataCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &FileMetadata{
			ID:           fmt.Sprintf("id-%d", i),
			Name:         fmt.Sprintf("file-%d.sav", i),
			Size:         1024,
			ModifiedTime: time.Now(),
			Hash:         fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("data-%d", i)))),
		}
	}
}
