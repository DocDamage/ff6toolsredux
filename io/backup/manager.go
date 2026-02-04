package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"ffvi_editor/models"
)

// Manager handles backup operations
type Manager struct {
	backupDir    string
	maxBackups   int
	autoBackup   bool
	mu           sync.RWMutex
	backups      map[string]*models.BackupMetadata
	metadataFile string
}

// NewManager creates a new backup manager
func NewManager(backupDir string, maxBackups int) (*Manager, error) {
	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	m := &Manager{
		backupDir:    backupDir,
		maxBackups:   maxBackups,
		autoBackup:   maxBackups > 0,
		backups:      make(map[string]*models.BackupMetadata),
		metadataFile: filepath.Join(backupDir, "backups.json"),
	}

	// Load existing backups metadata
	if err := m.loadMetadata(); err != nil {
		// If metadata file doesn't exist, that's OK
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return m, nil
}

// CreateBackup creates a new backup of save file data
func (m *Manager) CreateBackup(originalPath string, data []byte, description string) (*models.BackupMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.autoBackup {
		return nil, fmt.Errorf("auto-backup is disabled")
	}

	// Calculate hashes
	fileHash := models.CalculateHash(data)

	// Create metadata
	metadata := models.NewBackupMetadata(originalPath, int64(len(data)), fileHash, description)

	// Save backup file
	backupPath := filepath.Join(m.backupDir, metadata.ID+".bak")
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write backup file: %w", err)
	}

	// Store metadata
	m.backups[metadata.ID] = &metadata

	// Cleanup old backups
	if err := m.cleanupOldBackups(); err != nil {
		return nil, fmt.Errorf("failed to cleanup old backups: %w", err)
	}

	// Save metadata
	if err := m.saveMetadata(); err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return &metadata, nil
}

// ListBackups returns all backups sorted by timestamp (newest first)
func (m *Manager) ListBackups() []models.BackupListEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := make([]models.BackupListEntry, 0, len(m.backups))
	for _, meta := range m.backups {
		entries = append(entries, models.BackupListEntry{
			ID:          meta.ID,
			Timestamp:   meta.Timestamp,
			Description: meta.Description,
			FileSize:    meta.FileSize,
			TimeSince:   formatTimeSince(time.Since(meta.Timestamp)),
		})
	}

	// Sort by timestamp descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	return entries
}

// RestoreBackup restores a backup to its original path
func (m *Manager) RestoreBackup(backupID string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	meta, exists := m.backups[backupID]
	if !exists {
		return nil, fmt.Errorf("backup not found: %s", backupID)
	}

	// Read backup file
	backupPath := filepath.Join(m.backupDir, backupID+".bak")
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	// Verify hash
	if calculatedHash := models.CalculateHash(data); calculatedHash != meta.Hash {
		return nil, fmt.Errorf("backup integrity check failed: hash mismatch")
	}

	return data, nil
}

// DeleteBackup removes a backup
func (m *Manager) DeleteBackup(backupID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.backups[backupID]; !exists {
		return fmt.Errorf("backup not found: %s", backupID)
	}

	// Delete backup file
	backupPath := filepath.Join(m.backupDir, backupID+".bak")
	if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}

	// Remove metadata
	delete(m.backups, backupID)

	// Save updated metadata
	return m.saveMetadata()
}

// GetBackupMetadata returns metadata for a backup
func (m *Manager) GetBackupMetadata(backupID string) (*models.BackupMetadata, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	meta, exists := m.backups[backupID]
	if !exists {
		return nil, fmt.Errorf("backup not found: %s", backupID)
	}

	return meta, nil
}

// SetMaxBackups configures maximum number of backups to keep
func (m *Manager) SetMaxBackups(max int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.maxBackups = max
	m.autoBackup = max > 0
}

// GetMaxBackups returns configured maximum
func (m *Manager) GetMaxBackups() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.maxBackups
}

// BackupCount returns number of stored backups
func (m *Manager) BackupCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.backups)
}

// cleanupOldBackups removes oldest backups if count exceeds maximum
func (m *Manager) cleanupOldBackups() error {
	if m.maxBackups <= 0 || len(m.backups) <= m.maxBackups {
		return nil
	}

	// Sort by timestamp to find oldest
	backups := make([]*models.BackupMetadata, 0, len(m.backups))
	for _, meta := range m.backups {
		backups = append(backups, meta)
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.Before(backups[j].Timestamp)
	})

	// Delete oldest backups
	toDelete := len(backups) - m.maxBackups
	for i := 0; i < toDelete; i++ {
		backupPath := filepath.Join(m.backupDir, backups[i].ID+".bak")
		if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
			return err
		}
		delete(m.backups, backups[i].ID)
	}

	return nil
}

// loadMetadata loads backup metadata from file
func (m *Manager) loadMetadata() error {
	data, err := os.ReadFile(m.metadataFile)
	if err != nil {
		return err
	}

	var metadata map[string]models.BackupMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	m.backups = make(map[string]*models.BackupMetadata)
	for id, meta := range metadata {
		m.backups[id] = &models.BackupMetadata{
			ID:           meta.ID,
			Timestamp:    meta.Timestamp,
			OriginalPath: meta.OriginalPath,
			FileSize:     meta.FileSize,
			Hash:         meta.Hash,
			Description:  meta.Description,
			SaveGameHash: meta.SaveGameHash,
		}
	}

	return nil
}

// saveMetadata saves backup metadata to file
func (m *Manager) saveMetadata() error {
	data, err := json.MarshalIndent(m.backups, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(m.metadataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// formatTimeSince formats time duration for display
func formatTimeSince(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case d < 30*24*time.Hour:
		days := int(d.Hours()) / 24
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		weeks := int(d.Hours()) / (24 * 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}
}
