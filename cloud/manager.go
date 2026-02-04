package cloud

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Manager handles cloud sync operations
type Manager struct {
	providers  map[string]Provider // provider name -> provider instance
	config     *SyncConfig
	status     map[string]*SyncStatus // provider -> status
	mu         sync.RWMutex
	stopCh     chan struct{}
	syncTicker *time.Ticker
	logger     *log.Logger
}

// New creates a new cloud sync manager with default config
func New() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
		config: &SyncConfig{
			Enabled:            false,
			AutoSync:           false,
			SyncInterval:       30 * time.Minute,
			ConflictResolution: ConflictPromptUser,
			EncryptionEnabled:  true,
			EncryptionKey:      nil,
			FolderPath:         "",
		},
		status: make(map[string]*SyncStatus),
		stopCh: make(chan struct{}),
		logger: log.New(log.Writer(), "[cloud] ", log.LstdFlags),
	}
}

// NewManager creates a new cloud sync manager with configuration
func NewManager(config *SyncConfig) *Manager {
	return &Manager{
		providers: make(map[string]Provider),
		config:    config,
		status:    make(map[string]*SyncStatus),
		stopCh:    make(chan struct{}),
		logger:    log.New(log.Writer(), "[cloud] ", log.LstdFlags),
	}
}

// RegisterProvider registers a cloud storage provider
func (m *Manager) RegisterProvider(provider Provider) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := provider.GetName()
	if _, exists := m.providers[name]; exists {
		return fmt.Errorf("provider '%s' already registered", name)
	}

	m.providers[name] = provider
	m.status[name] = provider.GetStatus()
	return nil
}

// GetProvider retrieves a registered provider by name
func (m *Manager) GetProvider(name string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not registered", name)
	}

	return provider, nil
}

// ListProviders returns names of all registered providers
func (m *Manager) ListProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.providers))
	for name := range m.providers {
		names = append(names, name)
	}
	return names
}

// Start begins automatic sync if enabled
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Enabled || !m.config.AutoSync {
		return nil
	}

	if m.syncTicker != nil {
		return fmt.Errorf("sync already started")
	}

	m.syncTicker = time.NewTicker(m.config.SyncInterval)

	// Start periodic sync in background
	go m.syncLoop()

	return nil
}

// Stop stops automatic sync
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.syncTicker != nil {
		m.syncTicker.Stop()
		m.syncTicker = nil
	}

	close(m.stopCh)
}

// syncLoop runs periodic sync operations
func (m *Manager) syncLoop() {
	for {
		select {
		case <-m.syncTicker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			if err := m.SyncAll(ctx); err != nil {
				m.logger.Printf("[cloud] periodic sync failed: %v", err)
			}
			cancel()
		case <-m.stopCh:
			return
		}
	}
}

// Sync performs a full sync operation with the configured provider
func (m *Manager) Sync(ctx context.Context, providerName string) error {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return err
	}

	status := provider.GetStatus()
	if status.InProgress {
		return fmt.Errorf("sync already in progress for provider '%s'", providerName)
	}

	status.InProgress = true
	status.LastSync = time.Now()
	provider.SetStatus(status)

	defer func() {
		status.InProgress = false
		provider.SetStatus(status)
	}()

	// Check authentication
	if !provider.IsAuthenticated() {
		err := provider.Authenticate(ctx)
		if err != nil {
			status.LastError = fmt.Errorf("authentication failed: %w", err)
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Validate connection
	ok, msg := provider.ValidateConnection(ctx)
	if !ok {
		status.LastError = errors.New(msg)
		return fmt.Errorf("connection validation failed: %s", msg)
	}

	// Get storage quota
	used, total, err := provider.GetQuotaInfo(ctx)
	if err == nil {
		status.StorageUsed = used
		status.StorageTotal = total
	}

	// Perform sync
	if m.config.FolderPath != "" {
		_, err := provider.SyncFolder(ctx, m.config.FolderPath, m.config.FolderPath, m.config.ConflictResolution)
		if err != nil {
			status.LastError = err
			return fmt.Errorf("folder sync failed: %w", err)
		}
	}

	return nil
}

// SyncAll performs sync with all enabled providers
func (m *Manager) SyncAll(ctx context.Context) error {
	m.mu.RLock()
	providers := make([]Provider, 0, len(m.providers))
	for _, p := range m.providers {
		providers = append(providers, p)
	}
	m.mu.RUnlock()

	var errs []error
	for _, provider := range providers {
		if err := m.Sync(ctx, provider.GetName()); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", provider.GetName(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("sync errors: %v", errs)
	}

	return nil
}

// UploadFile uploads a single file to cloud storage via specified provider
func (m *Manager) UploadFile(ctx context.Context, providerName, localPath, remotePath string) error {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return err
	}

	if !provider.IsAuthenticated() {
		return fmt.Errorf("provider not authenticated")
	}

	// Open local file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Encrypt if enabled
	var reader io.Reader = file
	if m.config.EncryptionEnabled && m.config.EncryptionKey != nil {
		encrypted, err := m.encryptReader(file)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
		reader = encrypted
	}

	// Upload to provider
	_, err = provider.Upload(ctx, filepath.Base(remotePath), reader)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	return nil
}

// DownloadFile downloads a single file from cloud storage via specified provider
func (m *Manager) DownloadFile(ctx context.Context, providerName, fileID, localPath string) error {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return err
	}

	if !provider.IsAuthenticated() {
		return fmt.Errorf("provider not authenticated")
	}

	// Download from provider
	reader, err := provider.Download(ctx, fileID)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer reader.Close()

	// Decrypt if enabled
	var dataReader io.Reader = reader
	if m.config.EncryptionEnabled && m.config.EncryptionKey != nil {
		decrypted, err := m.decryptReader(reader)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		dataReader = decrypted
	}

	// Write to local file
	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, dataReader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetStatus returns the sync status for a specific provider
func (m *Manager) GetStatus(providerName string) (*SyncStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, exists := m.status[providerName]
	if !exists {
		return nil, fmt.Errorf("no status for provider '%s'", providerName)
	}

	return status, nil
}

// GetAllStatus returns status for all providers
func (m *Manager) GetAllStatus() map[string]*SyncStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*SyncStatus)
	for name, status := range m.status {
		statusCopy := *status
		result[name] = &statusCopy
	}
	return result
}

// GetConflicts returns all detected conflicts for a provider
func (m *Manager) GetConflicts(providerName string) ([]*Conflict, error) {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	status := provider.GetStatus()
	if status.ConflictsFound == 0 {
		return []*Conflict{}, nil
	}

	// In real implementation, would fetch conflicts from provider
	return []*Conflict{}, nil
}

// ResolveConflict resolves a specific conflict
func (m *Manager) ResolveConflict(ctx context.Context, providerName string, conflict *Conflict, resolution ConflictResolution) error {
	provider, err := m.GetProvider(providerName)
	if err != nil {
		return err
	}

	if !provider.IsAuthenticated() {
		return fmt.Errorf("provider not authenticated")
	}

	// Handle different resolution strategies
	switch resolution {
	case ConflictUseLocal:
		// Delete remote version
		return provider.Delete(ctx, conflict.RemoteID)

	case ConflictUseRemote:
		// Delete local version
		return os.Remove(conflict.LocalID)

	case ConflictNewest:
		if conflict.LocalTime.After(conflict.RemoteTime) {
			return provider.Delete(ctx, conflict.RemoteID)
		}
		return os.Remove(conflict.LocalID)

	case ConflictCreateCopy:
		// Keep both, rename remote
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		_ = fmt.Sprintf("%s.%s", conflict.RemoteID, timestamp)
		// This would require provider-specific renaming logic
		return nil

	default:
		return fmt.Errorf("unknown conflict resolution strategy: %v", resolution)
	}
}

// encryptReader encrypts data from a reader using AES-256-GCM
func (m *Manager) encryptReader(reader io.Reader) (io.Reader, error) {
	// Read all data (for simplicity; streaming encryption would be better for large files)
	plaintext, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Create cipher
	block, err := aes.NewCipher(m.config.EncryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Return as reader
	return &readSeeker{data: ciphertext}, nil
}

// decryptReader decrypts data from a reader using AES-256-GCM
func (m *Manager) decryptReader(reader io.Reader) (io.Reader, error) {
	// Read all data
	ciphertext, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Create cipher
	block, err := aes.NewCipher(m.config.EncryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	// Return as reader
	return &readSeeker{data: plaintext}, nil
}

// HashFile computes SHA256 hash of a file
func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// readSeeker wraps a byte slice as an io.Reader
type readSeeker struct {
	data []byte
	pos  int
}

func (r *readSeeker) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
