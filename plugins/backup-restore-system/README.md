# Backup & Restore System Plugin

**Version:** 1.1.0  
**Category:** Utility  
**Author:** FF6 Plugin Team

## Overview

The **Backup & Restore System** provides comprehensive backup, restore, versioning, and disaster recovery capabilities for FF6 save files. The plugin now includes a **Quick Backup** feature for one-click save backups with automatic version history.

**NEW in v1.1.0:** 🎉 **Quick Backup** - One-click save backups with instant restore!

## Features

### Quick Backup (NEW v1.1.0) ⭐
- **One-Click Backup** - Instant save file backup with single command
- **Backup History** - Track up to 50 recent backups automatically
- **One-Click Restore** - Restore any backup by name or ID
- **Auto-Backup** - Optional automatic backups before risky operations
- **Smart Management** - Automatic cleanup and storage management

### Core Backup Features
- **Full Backups** - Complete save file snapshots
- **Incremental Backups** - Only changed data since last backup
- **Differential Backups** - Changes since last full backup
- **Scheduled Backups** - Automatic backup scheduling

### Snapshot Management
- **Snapshots** - Point-in-time data capture
- **Snapshot Cloning** - Duplicate snapshots for testing
- **Export/Delete** - Manage snapshot lifecycle

### Recovery System
- **Restore Operations** - Recover from any backup
- **Integrity Verification** - Ensure backup validity
- **Point-in-Time Recovery** - Restore to specific timestamp
- **Test Restore** - Dry-run before actual restore

### Version Control
- **Versioning** - Tag and track save file versions
- **Rollback** - Revert to previous versions
- **Comparison** - Compare different versions
- **Release Tags** - Mark important milestones

## Installation

1. Copy the `backup-restore-system` folder to your `plugins` directory
2. Restart the FF6 Save Editor
3. Load the plugin via Plugin Menu → Backup & Restore System

## Quick Start Guide

### Quick Backup Usage (NEW! ⭐)

#### Example 1: Create Manual Backup
```lua
local BackupSystem = require("plugins/backup-restore-system/v1_0_core")
local QB = BackupSystem.QuickBackup

-- Create backup with custom name
local backup = QB.backupNow("C:/ff6saves/save01.ff6", "Before_Final_Battle")
print("Backup created: " .. backup.name)

-- Create backup with automatic naming
local backup = QB.backupNow("C:/ff6saves/save01.ff6")
-- Creates: "Manual_20260116_143022"
```

#### Example 2: View Backup History
```lua
-- Display formatted history
QB.displayBackupHistory()

-- Get backup history programmatically
local history = QB.getBackupHistory(10)
print("Total backups: " .. history.total_backups)

for i, backup in ipairs(history.backups) do
  print(string.format("%d. %s [%s]", i, backup.name, backup.created_date))
end
```

#### Example 3: Restore from Backup
```lua
-- Restore by backup name
local result = QB.restoreBackup("Before_Final_Battle")
if result.success then
  print("Restored successfully!")
else
  print("Error: " .. result.error)
end

-- Restore by backup ID
local result = QB.restoreBackup("QUICK_1737053422")
```

#### Example 4: Enable Auto-Backup
```lua
-- Enable automatic backups before risky operations
QB.setAutoBackup(true)

-- Now auto-backup before risky operation
local backup = QB.autoBackupBeforeOperation(
  "C:/ff6saves/save01.ff6",
  "max_all_stats"
)
-- Creates: "Auto_max_all_stats_20260116_143530"

-- Disable auto-backup
QB.setAutoBackup(false)
```

#### Example 5: Backup Management
```lua
-- Get info about specific backup
local info = QB.getBackupInfo("Before_Final_Battle")
if info.success then
  print("Backup: " .. info.backup.name)
  print("Created: " .. info.backup.created_date)
  print("Size: " .. info.backup.size_kb .. " KB")
end

-- Delete old backup
local result = QB.deleteBackup("old_backup_name")

-- Clean old backups (keep last 10)
local result = QB.cleanOldBackups(10)
print(string.format("Removed %d old backups", result.removed_count))
```

### Traditional Backup Features

#### Full Backup
```lua
local BackupSystem = require("plugins/backup-restore-system/v1_0_core")
local BE = BackupSystem.BackupEngine

local backup = BE.createFullBackup("Daily_Backup", save_data)
print("Backup created: " .. backup.backup_id)
```

#### Snapshot Management
```lua
local SM = BackupSystem.SnapshotManagement

-- Create snapshot
local snapshot = SM.createSnapshot("BeforeBoss", save_data)

-- List recent snapshots
local snapshots = SM.listSnapshots(10)
for _, snap in ipairs(snapshots.recent) do
  print(snap.name .. " - " .. snap.size_mb .. " MB")
end

-- Clone snapshot for testing
local clone = SM.cloneSnapshot(snapshot.snapshot_id, "Test_Clone")
```

#### Restore and Recovery
```lua
local RS = BackupSystem.RecoverySystem

-- Verify backup before restore
local verification = RS.verifyBackupIntegrity(backup_id)
if verification.integrity_score == 100.0 then
  -- Restore from backup
  local restore = RS.restoreFromBackup(backup_id, {full_restore = true})
  print("Restore status: " .. restore.status)
end

-- Point-in-time recovery
local recovery = RS.pointInTimeRecovery(target_timestamp)
```

#### Version Control
```lua
local VC = BackupSystem.VersionControl

-- Create version
local version = VC.createVersion("v1.5.0", save_data, "Feature complete")

-- List versions
local versions = VC.listVersions(5)
for _, ver in ipairs(versions.recent) do
  print(ver.label .. ": " .. ver.description)
end

-- Rollback to previous version
local rollback = VC.rollbackToVersion("v1.4.9")

-- Compare versions
local comparison = VC.compareVersions("v1.5.0", "v1.4.9")
print("Differences: " .. comparison.differences)
```

## Quick Backup API Reference

### `QuickBackup.backupNow(save_file_path, backup_name)`
Create instant backup of save file.

**Parameters:**
- `save_file_path` (string) - Path to save file to backup
- `backup_name` (string, optional) - Custom backup name (auto-generated if omitted)

**Returns:** `table` - Backup metadata with success status

### `QuickBackup.getBackupHistory(limit)`
Get list of recent backups.

**Parameters:**
- `limit` (number, optional) - Max backups to return (default: 10)

**Returns:** `table` - Backup history with total count and backup list

### `QuickBackup.displayBackupHistory()`
Display formatted backup history.

**Returns:** `string` - Formatted backup history display

### `QuickBackup.restoreBackup(backup_identifier)`
Restore save file from backup.

**Parameters:**
- `backup_identifier` (string) - Backup ID or name

**Returns:** `table` - Restore result with success status

### `QuickBackup.deleteBackup(backup_identifier)`
Delete backup from history.

**Parameters:**
- `backup_identifier` (string) - Backup ID or name

**Returns:** `table` - Deletion result with success status

### `QuickBackup.setAutoBackup(enabled)`
Enable/disable automatic backups.

**Parameters:**
- `enabled` (boolean) - true to enable, false to disable

**Returns:** `table` - Configuration result

### `QuickBackup.autoBackupBeforeOperation(save_file_path, operation_name)`
Create automatic backup before risky operation.

**Parameters:**
- `save_file_path` (string) - Path to save file
- `operation_name` (string) - Name of operation for backup naming

**Returns:** `table` - Backup metadata

### `QuickBackup.getBackupInfo(backup_identifier)`
Get detailed information about specific backup.

**Parameters:**
- `backup_identifier` (string) - Backup ID or name

**Returns:** `table` - Backup information

### `QuickBackup.cleanOldBackups(keep_count)`
Remove old backups, keep only recent ones.

**Parameters:**
- `keep_count` (number, optional) - Number of backups to keep (default: 10)

**Returns:** `table` - Cleanup result with counts

## Backup Strategy Recommendations

### For Casual Players
```lua
-- Simple approach: Manual backups before important events
QB.backupNow(save_file, "Before_Imperial_Base")
QB.backupNow(save_file, "Before_Kefka_Tower")
QB.backupNow(save_file, "Before_Final_Boss")
```

### For Experimenters
```lua
-- Enable auto-backup for safety
QB.setAutoBackup(true)

-- All risky operations now auto-backup
QB.autoBackupBeforeOperation(save_file, "max_stats_experiment")
QB.autoBackupBeforeOperation(save_file, "rare_item_duplication")

-- Restore if something goes wrong
QB.restoreBackup("Auto_max_stats_experiment_20260116_143530")
```

### For Power Users
```lua
-- Combination of manual and automatic
QB.setAutoBackup(true)

-- Manual milestone backups
QB.backupNow(save_file, "100_Percent_Complete")
QB.backupNow(save_file, "All_Characters_Level_99")

-- View history
QB.displayBackupHistory()

-- Clean old auto-backups, keep milestones
-- (Delete specific auto-backups manually)
```

## User Benefits

✅ **Never Lose Progress** - Backup before any risky action  
✅ **Experiment Safely** - Try anything, restore if it fails  
✅ **One-Click Operation** - Backup and restore in seconds  
✅ **Automatic Protection** - Auto-backup before risky edits  
✅ **Storage Management** - Automatic cleanup of old backups  
✅ **Peace of Mind** - Always have a recent backup available

## Technical Details

- **Plugin Version:** 1.1.0
- **Lines of Code:** ~1,150 (200 LOC added for Quick Backup)
- **Backup History:** Maintains last 50 backups automatically
- **Backup Speed:** <1 second for typical FF6 save (~128 KB)
- **Storage Efficient:** Backups are compressed and deduplicated
- **Safe Operations:** All operations validate inputs and handle errors

## Version History

- **v1.1.0** - Added Quick Backup feature (Quick Win #2)
- **v1.0.0** - Initial release with full backup/restore capabilities

## Support

For issues or feature requests, contact the FF6 Plugin Team.

## License

See main FF6 Save Editor license.
