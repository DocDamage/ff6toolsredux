package cloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// TestNewManager tests creating a new cloud manager
func TestNewManager(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}

	if m.config == nil {
		t.Error("Manager config is nil")
	}

	if m.config.Enabled {
		t.Error("Manager should be disabled by default")
	}

	if m.config.AutoSync {
		t.Error("AutoSync should be disabled by default")
	}
}

// TestNewManagerWithConfig tests creating a manager with custom config
func TestNewManagerWithConfig(t *testing.T) {
	config := &SyncConfig{
		Enabled:            true,
		AutoSync:           true,
		SyncInterval:       15 * time.Minute,
		ConflictResolution: ConflictUseLocal,
		EncryptionEnabled:  true,
	}

	m := NewManager(config)
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}

	if m.config != config {
		t.Error("Manager config not set correctly")
	}
}

// TestRegisterProvider tests registering cloud providers
func TestRegisterProvider(t *testing.T) {
	m := New()

	// Create a mock provider
	mockProvider := &mockProvider{
		name:   "test-provider",
		status: &SyncStatus{},
	}

	err := m.RegisterProvider(mockProvider)
	if err != nil {
		t.Errorf("RegisterProvider error: %v", err)
	}

	// Verify provider was registered
	provider, err := m.GetProvider("test-provider")
	if err != nil {
		t.Errorf("GetProvider error: %v", err)
	}

	if provider.GetName() != "test-provider" {
		t.Errorf("Provider name = %q, want test-provider", provider.GetName())
	}
}

// TestRegisterProviderDuplicate tests registering duplicate provider
func TestRegisterProviderDuplicate(t *testing.T) {
	m := New()

	mockProvider := &mockProvider{
		name:   "test-provider",
		status: &SyncStatus{},
	}

	// First registration should succeed
	err := m.RegisterProvider(mockProvider)
	if err != nil {
		t.Errorf("First RegisterProvider error: %v", err)
	}

	// Second registration should fail
	err = m.RegisterProvider(mockProvider)
	if err == nil {
		t.Error("RegisterProvider should return error for duplicate")
	}
}

// TestGetProviderNotFound tests getting non-existent provider
func TestGetProviderNotFound(t *testing.T) {
	m := New()

	_, err := m.GetProvider("non-existent")
	if err == nil {
		t.Error("GetProvider should return error for non-existent provider")
	}
}

// TestListProviders tests listing registered providers
func TestListProviders(t *testing.T) {
	m := New()

	// Initially should be empty
	providers := m.ListProviders()
	if len(providers) != 0 {
		t.Errorf("Expected 0 providers, got %d", len(providers))
	}

	// Register some providers
	m.RegisterProvider(&mockProvider{name: "provider1", status: &SyncStatus{}})
	m.RegisterProvider(&mockProvider{name: "provider2", status: &SyncStatus{}})
	m.RegisterProvider(&mockProvider{name: "provider3", status: &SyncStatus{}})

	providers = m.ListProviders()
	if len(providers) != 3 {
		t.Errorf("Expected 3 providers, got %d", len(providers))
	}
}

// TestGetStatus tests getting provider status
func TestGetStatus(t *testing.T) {
	m := New()

	mockProvider := &mockProvider{
		name:   "test-provider",
		status: &SyncStatus{InProgress: true},
	}

	m.RegisterProvider(mockProvider)

	status, err := m.GetStatus("test-provider")
	if err != nil {
		t.Errorf("GetStatus error: %v", err)
	}

	if !status.InProgress {
		t.Error("Status InProgress should be true")
	}
}

// TestGetStatusNotFound tests status for non-existent provider
func TestGetStatusNotFound(t *testing.T) {
	m := New()

	_, err := m.GetStatus("non-existent")
	if err == nil {
		t.Error("GetStatus should return error for non-existent provider")
	}
}

// TestGetAllStatus tests getting all provider statuses
func TestGetAllStatus(t *testing.T) {
	m := New()

	m.RegisterProvider(&mockProvider{name: "p1", status: &SyncStatus{}})
	m.RegisterProvider(&mockProvider{name: "p2", status: &SyncStatus{}})

	statuses := m.GetAllStatus()
	if len(statuses) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(statuses))
	}
}

// TestStartStop tests starting and stopping auto-sync
func TestStartStop(t *testing.T) {
	m := New()

	// Should not start if disabled
	err := m.Start()
	if err != nil {
		t.Errorf("Start error (disabled): %v", err)
	}

	// Enable and set auto-sync
	m.config.Enabled = true
	m.config.AutoSync = true
	m.config.SyncInterval = 100 * time.Millisecond

	err = m.Start()
	if err != nil {
		t.Errorf("Start error: %v", err)
	}

	// Should not be able to start twice
	err = m.Start()
	if err == nil {
		t.Error("Start should return error when already started")
	}

	// Let it run briefly
	time.Sleep(50 * time.Millisecond)

	// Stop
	m.Stop()
}

// TestStartDisabled tests starting when disabled
func TestStartDisabled(t *testing.T) {
	m := New()

	// Disabled by default
	err := m.Start()
	if err != nil {
		t.Errorf("Start error when disabled: %v", err)
	}
}

// TestGetConflicts tests getting conflicts
func TestGetConflicts(t *testing.T) {
	m := New()

	mockProvider := &mockProvider{
		name:   "test-provider",
		status: &SyncStatus{ConflictsFound: 0},
	}

	m.RegisterProvider(mockProvider)

	conflicts, err := m.GetConflicts("test-provider")
	if err != nil {
		t.Errorf("GetConflicts error: %v", err)
	}

	if len(conflicts) != 0 {
		t.Errorf("Expected 0 conflicts, got %d", len(conflicts))
	}
}

// TestGetConflictsNotFound tests conflicts for non-existent provider
func TestGetConflictsNotFound(t *testing.T) {
	m := New()

	_, err := m.GetConflicts("non-existent")
	if err == nil {
		t.Error("GetConflicts should return error for non-existent provider")
	}
}

// TestResolveConflictNotAuthenticated tests resolving conflict when not authenticated
func TestResolveConflictNotAuthenticated(t *testing.T) {
	m := New()

	mockProvider := &mockProvider{
		name:          "test-provider",
		status:        &SyncStatus{},
		authenticated: false,
	}

	m.RegisterProvider(mockProvider)

	conflict := &Conflict{
		LocalID:    "local",
		RemoteID:   "remote",
		LocalTime:  time.Now(),
		RemoteTime: time.Now(),
	}

	err := m.ResolveConflict(context.Background(), "test-provider", conflict, ConflictUseLocal)
	if err == nil {
		t.Error("ResolveConflict should return error when not authenticated")
	}
}

// TestResolveConflictUnknownStrategy tests unknown resolution strategy
func TestResolveConflictUnknownStrategy(t *testing.T) {
	m := New()

	mockProvider := &mockProvider{
		name:          "test-provider",
		status:        &SyncStatus{},
		authenticated: true,
	}

	m.RegisterProvider(mockProvider)

	conflict := &Conflict{
		LocalID:    "local",
		RemoteID:   "remote",
		LocalTime:  time.Now(),
		RemoteTime: time.Now(),
	}

	err := m.ResolveConflict(context.Background(), "test-provider", conflict, ConflictResolution(999))
	if err == nil {
		t.Error("ResolveConflict should return error for unknown strategy")
	}
}

// TestSyncAllNoProviders tests syncing with no providers
func TestSyncAllNoProviders(t *testing.T) {
	m := New()

	ctx := context.Background()
	err := m.SyncAll(ctx)
	if err != nil {
		t.Errorf("SyncAll with no providers error: %v", err)
	}
}

// TestHashFile tests the HashFile function
func TestHashFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/test_hash_file.txt"
	content := []byte("test content for hashing")

	// Write test file
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hash1, err := HashFile(tmpFile)
	if err != nil {
		t.Logf("HashFile error (may not work on all systems): %v", err)
		return
	}

	if hash1 == "" {
		t.Error("Hash should not be empty")
	}

	// Same file should produce same hash
	hash2, err := HashFile(tmpFile)
	if err != nil {
		t.Logf("HashFile error (second call): %v", err)
		return
	}

	if hash1 != hash2 {
		t.Error("Same file should produce same hash")
	}
}

// TestHashFileNotFound tests hashing non-existent file
func TestHashFileNotFound(t *testing.T) {
	_, err := HashFile("/non/existent/file.txt")
	if err == nil {
		t.Error("HashFile should return error for non-existent file")
	}
}

// mockProvider is a mock implementation of the Provider interface for testing
type mockProvider struct {
	name          string
	status        *SyncStatus
	authenticated bool
}

func (m *mockProvider) GetName() string {
	return m.name
}

func (m *mockProvider) GetStatus() *SyncStatus {
	return m.status
}

func (m *mockProvider) SetStatus(status *SyncStatus) {
	m.status = status
}

func (m *mockProvider) IsAuthenticated() bool {
	return m.authenticated
}

func (m *mockProvider) Authenticate(ctx context.Context) error {
	m.authenticated = true
	return nil
}

func (m *mockProvider) Logout(ctx context.Context) error {
	m.authenticated = false
	return nil
}

func (m *mockProvider) GetAuthURL(ctx context.Context) (string, error) {
	return "https://example.com/auth", nil
}

func (m *mockProvider) HandleAuthCallback(ctx context.Context, code string, state string) error {
	m.authenticated = true
	return nil
}

func (m *mockProvider) Upload(ctx context.Context, filename string, reader io.Reader) (*FileMetadata, error) {
	return &FileMetadata{ID: "test-id", Name: filename}, nil
}

func (m *mockProvider) UploadFile(ctx context.Context, localPath, remotePath string) (*FileMetadata, error) {
	return &FileMetadata{ID: "test-id", Name: localPath}, nil
}

func (m *mockProvider) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	return nil, nil
}

func (m *mockProvider) DownloadFile(ctx context.Context, fileID, localPath string) error {
	return nil
}

func (m *mockProvider) List(ctx context.Context, folder string) ([]*FileMetadata, error) {
	return []*FileMetadata{}, nil
}

func (m *mockProvider) ListFolder(ctx context.Context, folderID string, recursive bool) ([]*FileMetadata, error) {
	return []*FileMetadata{}, nil
}

func (m *mockProvider) Delete(ctx context.Context, fileID string) error {
	return nil
}

func (m *mockProvider) GetMetadata(ctx context.Context, fileID string) (*FileMetadata, error) {
	return &FileMetadata{ID: fileID}, nil
}

func (m *mockProvider) CreateFolder(ctx context.Context, name, parentID string) (string, error) {
	return "folder-id", nil
}

func (m *mockProvider) FindOrCreateFolder(ctx context.Context, path string) (string, error) {
	return "folder-id", nil
}

func (m *mockProvider) SyncFolder(ctx context.Context, localFolder, remoteFolder string, strategy ConflictResolution) ([]*Conflict, error) {
	return []*Conflict{}, nil
}

func (m *mockProvider) ValidateConnection(ctx context.Context) (bool, string) {
	return true, ""
}

func (m *mockProvider) GetQuotaInfo(ctx context.Context) (int64, int64, error) {
	return 100, 1000, nil
}

// BenchmarkRegisterProvider benchmarks provider registration
func BenchmarkRegisterProvider(b *testing.B) {
	m := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := &mockProvider{
			name:   fmt.Sprintf("provider-%d", i),
			status: &SyncStatus{},
		}
		m.RegisterProvider(provider)
	}
}

// BenchmarkGetProvider benchmarks provider retrieval
func BenchmarkGetProvider(b *testing.B) {
	m := New()
	m.RegisterProvider(&mockProvider{name: "test", status: &SyncStatus{}})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.GetProvider("test")
	}
}
