package models

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"
)

// BackupMetadata holds information about a backup
type BackupMetadata struct {
	ID           string    `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	OriginalPath string    `json:"originalPath"`
	FileSize     int64     `json:"fileSize"`
	Hash         string    `json:"hash"`
	Description  string    `json:"description"`
	SaveGameHash string    `json:"saveGameHash"` // Hash of save file contents
}

// Backup represents a complete backup with metadata and file data
type Backup struct {
	Metadata BackupMetadata
	Data     []byte
}

// CalculateHash computes SHA256 hash of data
func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// CalculateHashFromReader computes SHA256 hash from reader
func CalculateHashFromReader(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// NewBackupMetadata creates a new backup metadata entry
func NewBackupMetadata(originalPath string, fileSize int64, fileHash string, description string) BackupMetadata {
	return BackupMetadata{
		ID:           generateID(),
		Timestamp:    time.Now(),
		OriginalPath: originalPath,
		FileSize:     fileSize,
		Hash:         fileHash,
		Description:  description,
		SaveGameHash: fileHash,
	}
}

// generateID creates a unique backup ID
func generateID() string {
	return time.Now().Format("20060102_150405_") + generateRandomSuffix()
}

// generateRandomSuffix creates a random suffix for uniqueness
func generateRandomSuffix() string {
	// Simple implementation - in production use crypto/rand
	return time.Now().Format("9999")
}

// BackupListEntry represents a backup in the backup list UI
type BackupListEntry struct {
	ID          string
	Timestamp   time.Time
	Description string
	FileSize    int64
	TimeSince   string
}
