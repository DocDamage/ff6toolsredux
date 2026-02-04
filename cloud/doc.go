// Package cloud provides cloud synchronization capabilities for save files.
//
// Supported Providers:
//   - Google Drive
//   - Dropbox
//
// Features:
//   - Automatic backup to cloud storage
//   - Conflict resolution strategies
//   - Encryption support
//   - Bandwidth management
//   - Sync status tracking
//
// The cloud manager coordinates multiple providers and handles:
//   - Authentication flows
//   - File upload/download
//   - Conflict detection and resolution
//   - Quota monitoring
//
// Usage:
//
//	manager := cloud.New()
//	gdrive := cloud.NewGoogleDriveProvider(clientID, secret)
//	manager.RegisterProvider(gdrive)
//	status, _ := manager.GetStatus("Google Drive")
package cloud
