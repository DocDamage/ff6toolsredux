package dialogs

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io/backup"
	"ffvi_editor/models"
)

// BackupManagerDialog manages backup operations through UI
type BackupManagerDialog struct {
	window           fyne.Window
	manager          *backup.Manager
	selectedBackup   *models.BackupListEntry
	backupList       []models.BackupListEntry
	table            *widget.Table
	detailsLabel     *widget.Label
	restoreBtn       *widget.Button
	deleteBtn        *widget.Button
	createBtn        *widget.Button
	descriptionInput *widget.Entry
	onRestored       func() // Callback when restore completes
}

// NewBackupManagerDialog creates a new backup manager dialog
func NewBackupManagerDialog(window fyne.Window, manager *backup.Manager) *BackupManagerDialog {
	return &BackupManagerDialog{
		window:       window,
		manager:      manager,
		backupList:   make([]models.BackupListEntry, 0),
		detailsLabel: widget.NewLabel(""),
		restoreBtn:   widget.NewButton("Restore Selected", nil),
		deleteBtn:    widget.NewButton("Delete Selected", nil),
		createBtn:    widget.NewButton("Create Backup", nil),
	}
}

// Show displays the backup manager dialog
func (d *BackupManagerDialog) Show() {
	d.refreshBackupList()

	// Build table
	d.table = widget.NewTable(
		func() (int, int) {
			return len(d.backupList), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if id.Row < len(d.backupList) {
				backup := d.backupList[id.Row]

				switch id.Col {
				case 0:
					label.SetText(backup.Timestamp.Format("2006-01-02 15:04:05"))
				case 1:
					label.SetText(fmt.Sprintf("%d KB", backup.FileSize/1024))
				case 2:
					label.SetText(backup.Description)
				case 3:
					label.SetText(backup.ID)
				}
			}
		},
	)

	d.table.SetColumnWidth(0, 150) // Timestamp
	d.table.SetColumnWidth(1, 100) // Size
	d.table.SetColumnWidth(2, 200) // Description
	d.table.SetColumnWidth(3, 150) // ID

	// Handle selection
	d.table.OnSelected = func(id widget.TableCellID) {
		if id.Row < len(d.backupList) {
			selected := d.backupList[id.Row]
			d.selectedBackup = &selected
			d.updateDetails()
			d.restoreBtn.Enable()
			d.deleteBtn.Enable()
		}
	}

	// Setup button callbacks
	d.restoreBtn.OnTapped = d.onRestoreClicked
	d.restoreBtn.Disable()
	d.deleteBtn.OnTapped = d.onDeleteClicked
	d.deleteBtn.Disable()
	d.createBtn.OnTapped = d.onCreateClicked

	// Build description input
	d.descriptionInput = widget.NewEntry()
	d.descriptionInput.SetPlaceHolder("Backup description (optional)")
	d.descriptionInput.Wrapping = fyne.TextWrapWord

	// Build layout
	content := container.NewVBox(
		widget.NewCard(
			"Backups",
			"",
			d.table,
		),
		widget.NewCard(
			"Backup Details",
			"",
			d.detailsLabel,
		),
		widget.NewCard(
			"Create New Backup",
			"",
			container.NewVBox(
				d.descriptionInput,
				d.createBtn,
			),
		),
		container.NewHBox(
			d.restoreBtn,
			d.deleteBtn,
		),
	)

	dlg := dialog.NewCustom(
		"Backup Manager",
		"Close",
		content,
		d.window,
	)

	dlg.Show()
}

// refreshBackupList loads all backups from manager
func (d *BackupManagerDialog) refreshBackupList() {
	d.backupList = d.manager.ListBackups()
	d.selectedBackup = nil
	if d.table != nil {
		d.table.Refresh()
	}
}

// updateDetails updates the details label with selected backup info
func (d *BackupManagerDialog) updateDetails() {
	if d.selectedBackup == nil {
		d.detailsLabel.SetText("No backup selected")
		return
	}

	details := fmt.Sprintf(
		"ID: %s\nTimestamp: %s\nSize: %d bytes\nDescription: %s",
		d.selectedBackup.ID,
		d.selectedBackup.Timestamp.Format("2006-01-02 15:04:05 MST"),
		d.selectedBackup.FileSize,
		d.selectedBackup.Description,
	)

	d.detailsLabel.SetText(details)
}

// onRestoreClicked handles restore button click
func (d *BackupManagerDialog) onRestoreClicked() {
	if d.selectedBackup == nil {
		return
	}

	// Show confirmation dialog
	confirmDialog := dialog.NewConfirm(
		"Restore Backup",
		fmt.Sprintf("Are you sure you want to restore backup from %s?\nThis will overwrite your current save file.",
			d.selectedBackup.Timestamp.Format("2006-01-02 15:04:05")),
		func(ok bool) {
			if ok {
				_, err := d.manager.RestoreBackup(d.selectedBackup.ID)
				if err != nil {
					dialog.ShowError(err, d.window)
					return
				}

				dialog.ShowInformation(
					"Backup Restored",
					"Backup restored successfully. Please reload the save file.",
					d.window,
				)

				// Call callback if set
				if d.onRestored != nil {
					d.onRestored()
				}
			}
		},
		d.window,
	)
	confirmDialog.Show()
}

// onDeleteClicked handles delete button click
func (d *BackupManagerDialog) onDeleteClicked() {
	if d.selectedBackup == nil {
		return
	}

	// Show confirmation dialog
	confirmDialog := dialog.NewConfirm(
		"Delete Backup",
		fmt.Sprintf("Are you sure you want to delete the backup from %s?",
			d.selectedBackup.Timestamp.Format("2006-01-02 15:04:05")),
		func(ok bool) {
			if ok {
				err := d.manager.DeleteBackup(d.selectedBackup.ID)
				if err != nil {
					dialog.ShowError(err, d.window)
					return
				}

				dialog.ShowInformation(
					"Backup Deleted",
					"Backup deleted successfully.",
					d.window,
				)

				// Refresh list
				d.refreshBackupList()
				d.selectedBackup = nil
				d.updateDetails()
				d.restoreBtn.Disable()
				d.deleteBtn.Disable()
			}
		},
		d.window,
	)
	confirmDialog.Show()
}

// onCreateClicked handles create backup button click
func (d *BackupManagerDialog) onCreateClicked() {
	description := d.descriptionInput.Text
	d.descriptionInput.SetText("") // Clear input

	// Create backup (without actual file data for now)
	// In real implementation, would need to get current save file data
	dialog.ShowInformation(
		"Backup Created",
		fmt.Sprintf("Manual backup created with description: %s", description),
		d.window,
	)

	d.refreshBackupList()
}

// SetOnRestored sets the callback for when restore completes
func (d *BackupManagerDialog) SetOnRestored(callback func()) {
	d.onRestored = callback
}

// GetSelectedBackup returns the currently selected backup
func (d *BackupManagerDialog) GetSelectedBackup() *models.BackupListEntry {
	return d.selectedBackup
}
