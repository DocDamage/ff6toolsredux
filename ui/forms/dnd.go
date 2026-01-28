package forms

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// DragDropHandler manages drag and drop operations
type DragDropHandler struct {
	onFileDrop     func(filePath string) error
	onTemplateDrop func(templateID string) error
	onPresetDrop   func(presetID string) error
	onBackupDrop   func(backupFile string) error
}

// NewDragDropHandler creates a new drag and drop handler
func NewDragDropHandler() *DragDropHandler {
	return &DragDropHandler{}
}

// SetFileDropCallback sets the callback for file drops
func (dh *DragDropHandler) SetFileDropCallback(callback func(filePath string) error) {
	dh.onFileDrop = callback
}

// SetTemplateDropCallback sets the callback for template drops
func (dh *DragDropHandler) SetTemplateDropCallback(callback func(templateID string) error) {
	dh.onTemplateDrop = callback
}

// SetPresetDropCallback sets the callback for preset drops
func (dh *DragDropHandler) SetPresetDropCallback(callback func(presetID string) error) {
	dh.onPresetDrop = callback
}

// SetBackupDropCallback sets the callback for backup drops
func (dh *DragDropHandler) SetBackupDropCallback(callback func(backupFile string) error) {
	dh.onBackupDrop = callback
}

// HandleFilePath processes dropped file paths
func (dh *DragDropHandler) HandleFilePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("empty file path")
	}

	// Determine file type by extension
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".srm", ".sav":
		// Save file dropped
		if dh.onFileDrop != nil {
			return dh.onFileDrop(filePath)
		}
	case ".json":
		// Could be template or preset
		if dh.onTemplateDrop != nil {
			if err := dh.onTemplateDrop(filePath); err == nil {
				return nil
			}
		}
		if dh.onPresetDrop != nil {
			return dh.onPresetDrop(filePath)
		}
	case ".backup":
		// Backup file dropped
		if dh.onBackupDrop != nil {
			return dh.onBackupDrop(filePath)
		}
	}

	return fmt.Errorf("unsupported file type: %s", ext)
}

// HandleText processes dropped text (could be preset/template ID)
func (dh *DragDropHandler) HandleText(text string) error {
	if text == "" {
		return fmt.Errorf("empty text")
	}

	// Check if it looks like a template or preset ID
	if dh.onTemplateDrop != nil {
		if err := dh.onTemplateDrop(text); err == nil {
			return nil
		}
	}
	if dh.onPresetDrop != nil {
		if err := dh.onPresetDrop(text); err == nil {
			return nil
		}
	}

	return fmt.Errorf("unrecognized ID: %s", text)
}

// CreateDropZone creates a visual drop zone container for UI
func CreateDropZone(hint string, handler *DragDropHandler, callback func(string) error) *fyne.Container {
	label := widget.NewLabel(hint)

	zone := container.NewCenter(label)

	// Store the callback for reference
	if callback != nil && handler != nil {
		// In actual implementation, would hook into widget's drag/drop events
		// For now, this provides the visual structure
	}

	return zone
}

// GetSaveFileDropZoneHint returns hint text for save file drop zone
func GetSaveFileDropZoneHint() string {
	return "Drop save files (.srm, .sav) here to open"
}

// GetTemplateDropZoneHint returns hint text for template drop zone
func GetTemplateDropZoneHint() string {
	return "Drop template files (.json) here to apply"
}

// GetPresetDropZoneHint returns hint text for preset drop zone
func GetPresetDropZoneHint() string {
	return "Drop preset files (.json) here to load"
}

// GetBackupDropZoneHint returns hint text for backup drop zone
func GetBackupDropZoneHint() string {
	return "Drop backup files (.backup) here to restore"
}
