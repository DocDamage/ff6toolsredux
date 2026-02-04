// Package backup provides save file backup management.
//
// The backup package handles:
//   - Creating timestamped backups of save files
//   - Managing backup retention policies
//   - Listing and restoring from backups
//   - Automatic backup on save
//
// Features:
//   - Configurable number of backups to keep
//   - Custom backup location
//   - Backup metadata (timestamp, original file, etc.)
//
// Example usage:
//
//	// Create a backup manager
//	bm := backup.NewManager(backupDir, maxBackups)
//
//	// Create a backup
//	if err := bm.CreateBackup(saveFilePath); err != nil {
//	    log.Printf("Backup failed: %v", err)
//	}
//
//	// List backups
//	backups, err := bm.ListBackups()
//	if err != nil {
//	    return err
//	}
package backup
