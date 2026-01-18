package cloud

import (
	"context"
	"io"
	"time"
)

// Provider defines the interface for cloud storage providers
type Provider interface {
	// Authentication
	// Authenticate performs OAuth2 authentication flow
	Authenticate(ctx context.Context) error

	// IsAuthenticated checks if the provider is authenticated
	IsAuthenticated() bool

	// Logout revokes the authentication token
	Logout(ctx context.Context) error

	// GetAuthURL returns the OAuth authentication URL for manual auth flow
	GetAuthURL(ctx context.Context) (string, error)

	// HandleAuthCallback processes the OAuth callback code
	HandleAuthCallback(ctx context.Context, code string, state string) error

	// File operations
	// Upload uploads a file to cloud storage
	Upload(ctx context.Context, filename string, reader io.Reader) (*FileMetadata, error)

	// UploadFile uploads a local file to cloud storage
	// localPath: path to local file
	// remotePath: destination in cloud storage
	UploadFile(ctx context.Context, localPath, remotePath string) (*FileMetadata, error)

	// Download downloads a file from cloud storage
	Download(ctx context.Context, fileID string) (io.ReadCloser, error)

	// DownloadFile downloads a cloud file to local path
	// fileID: cloud file ID
	// localPath: destination local path
	DownloadFile(ctx context.Context, fileID, localPath string) error

	// List lists all files in the cloud storage folder
	List(ctx context.Context, folder string) ([]*FileMetadata, error)

	// ListFolder lists files in a specific folder with optional recursion
	ListFolder(ctx context.Context, folderID string, recursive bool) ([]*FileMetadata, error)

	// Delete deletes a file from cloud storage
	Delete(ctx context.Context, fileID string) error

	// GetMetadata retrieves metadata for a file
	GetMetadata(ctx context.Context, fileID string) (*FileMetadata, error)

	// Folder operations
	// CreateFolder creates a new folder
	CreateFolder(ctx context.Context, name, parentID string) (string, error)

	// FindOrCreateFolder finds or creates a nested folder path
	FindOrCreateFolder(ctx context.Context, path string) (string, error)

	// Sync operations
	// SyncFolder synchronizes a local folder with cloud storage
	SyncFolder(ctx context.Context, localFolder, remoteFolder string, strategy ConflictResolution) ([]*Conflict, error)

	// GetName returns the provider name (e.g., "Google Drive", "Dropbox")
	GetName() string

	// GetStatus returns current sync status
	GetStatus() *SyncStatus

	// SetStatus updates sync status
	SetStatus(status *SyncStatus)

	// ValidateConnection tests connectivity and authentication
	ValidateConnection(ctx context.Context) (bool, string)

	// GetQuotaInfo returns storage quota info (used, total)
	GetQuotaInfo(ctx context.Context) (int64, int64, error)
}

// FileMetadata contains metadata about a cloud file
type FileMetadata struct {
	ID           string    // Unique file ID from provider
	Name         string    // File name
	Size         int64     // File size in bytes
	ModifiedTime time.Time // Last modified time
	Hash         string    // File hash (MD5 or SHA256)
	IsFolder     bool      // Whether this is a folder
	Path         string    // Full path in cloud storage
	Parents      []string  // Parent folder IDs
}

// ConflictResolution defines how to handle sync conflicts
type ConflictResolution int

const (
	// ConflictUseLocal keeps the local version
	ConflictUseLocal ConflictResolution = iota
	// ConflictUseRemote keeps the remote version
	ConflictUseRemote
	// ConflictCreateCopy creates a copy with timestamp suffix
	ConflictCreateCopy
	// ConflictPromptUser prompts the user to choose
	ConflictPromptUser
	// ConflictNewest keeps the newest version by timestamp
	ConflictNewest
)

// SyncConfig contains configuration for cloud sync
type SyncConfig struct {
	Provider           Provider
	Enabled            bool
	AutoSync           bool                // Auto-sync on save
	SyncInterval       time.Duration       // Interval for periodic sync
	ConflictResolution ConflictResolution  // How to handle conflicts
	EncryptionEnabled  bool                // Encrypt before upload
	EncryptionKey      []byte              // Encryption key (32 bytes for AES-256)
	FolderPath         string              // Remote folder path
	MaxRetries         int                 // Maximum retry attempts
	RetryDelay         time.Duration       // Delay between retries
	VerifyHashes       bool                // Verify file integrity via hashing
	CompressFiles      bool                // Compress files before upload
}

// SyncStatus represents the status of a sync operation
type SyncStatus struct {
	Provider          string        // Provider name
	InProgress        bool
	LastSync          time.Time
	NextSync          time.Time
	LastError         error
	FilesUploaded     int
	FilesDownloaded   int
	ConflictsFound    int
	Progress          float64       // 0.0 to 1.0
	CurrentFile       string        // Currently processing file
	IsAuthenticated   bool
	StorageUsed       int64
	StorageTotal      int64
}

// Conflict represents a sync conflict
type Conflict struct {
	FileName     string
	LocalTime    time.Time
	RemoteTime   time.Time
	LocalHash    string
	RemoteHash   string
	LocalSize    int64
	RemoteSize   int64
	LocalID      string  // Local file ID/path
	RemoteID     string  // Remote file ID
	Resolution   ConflictResolution
}
