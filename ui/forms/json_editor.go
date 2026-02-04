package forms

import (
	"encoding/json"
	"fmt"
	"io"

	ioJson "ffvi_editor/io/json"
	ipr "ffvi_editor/io/pr"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// JSONExportImportDialog handles JSON export/import of save data
type JSONExportImportDialog struct {
	dialog       dialog.Dialog
	prData       *ipr.PR
	window       fyne.Window
	formatSelect *widget.Select
	onExport     func(string) error
	onImport     func([]byte) error
}

// NewJSONExportImportDialog creates a new JSON export/import dialog
func NewJSONExportImportDialog(prData *ipr.PR, window fyne.Window) *JSONExportImportDialog {
	return &JSONExportImportDialog{
		prData: prData,
		window: window,
		formatSelect: widget.NewSelect(
			[]string{"Full Save", "Characters Only", "Inventory", "Party", "Magic", "Espers", "Equipment"},
			func(s string) {},
		),
	}
}

// Show displays the JSON export/import dialog
func (j *JSONExportImportDialog) Show() {
	j.formatSelect.SetSelected("Full Save")

	// Export section
	exportLabel := widget.NewLabel("Export Format:")

	exportDesc := widget.NewLabel("Select what data to include in the export:")

	exportBtn := widget.NewButton("Export as JSON", func() {
		format := j.getExportFormat(j.formatSelect.Selected)

		// Create exporter
		exporter := ioJson.NewExporter(j.prData)

		// Export based on format
		jsonBytes, err := exporter.ExportToJSON(format)
		if err != nil {
			dialog.ShowError(fmt.Errorf("export failed: %w", err), j.window)
			return
		}

		// Trigger file save dialog
		saveDialog := dialog.NewFileSave(func(reader fyne.URIWriteCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			_, err = reader.Write(jsonBytes)
			if err != nil {
				dialog.ShowError(err, j.window)
				return
			}

			dialog.ShowInformation("Export Complete", "Save data exported successfully", j.window)
		}, j.window)

		saveDialog.SetFileName("ff6_export.json")
		saveDialog.Show()
	})

	exportBox := container.NewVBox(
		exportLabel,
		j.formatSelect,
		exportDesc,
		exportBtn,
	)

	// Import section
	importLabel := widget.NewLabel("Import JSON:")

	importDesc := widget.NewLabel("Select a JSON file to import data from:")

	importBtn := widget.NewButton("Import from JSON", func() {
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			// Read file
			jsonBytes, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to read file: %w", err), j.window)
				return
			}

			// Validate JSON structure
			var export ioJson.SaveExport
			if err := json.Unmarshal(jsonBytes, &export); err != nil {
				dialog.ShowError(fmt.Errorf("invalid JSON format: %w", err), j.window)
				return
			}

			// Import the data with actual importer
			importer := ioJson.NewImporter(j.prData)
			if err := importer.ImportFromJSON(jsonBytes, ioJson.FormatFull); err != nil {
				dialog.ShowError(fmt.Errorf("import failed: %w", err), j.window)
				return
			}

			dialog.ShowInformation("Import Complete", fmt.Sprintf("Imported %d characters", len(export.Characters)), j.window)
		}, j.window)

		openDialog.Show()
	})

	importBox := container.NewVBox(
		importLabel,
		importDesc,
		importBtn,
	)

	// Preview area
	previewLabel := widget.NewLabel("Format: Full Save - Contains all character and equipment data")

	updatePreview := func() {
		format := j.formatSelect.Selected
		previewLabel.SetText(fmt.Sprintf("Format: %s - Selective export", format))
	}

	j.formatSelect.OnChanged = func(s string) {
		updatePreview()
	}

	updatePreview()

	// Buttons
	closeBtn := widget.NewButton("Close", func() {
		j.dialog.Hide()
	})

	buttons := container.NewHBox(closeBtn)

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Export Save Data"),
		exportBox,
		widget.NewLabel(""),
		widget.NewLabel("Import Save Data"),
		importBox,
		widget.NewLabel(""),
		previewLabel,
		buttons,
	)

	scrollContainer := container.NewScroll(content)

	j.dialog = dialog.NewCustom("JSON Export/Import", "Close", scrollContainer, j.window)
	j.dialog.Resize(fyne.NewSize(500, 600))
	j.dialog.Show()
}

// getExportFormat converts UI format string to export format enum
func (j *JSONExportImportDialog) getExportFormat(formatName string) ioJson.ExportFormat {
	switch formatName {
	case "Characters Only":
		return ioJson.FormatCharacters
	case "Inventory":
		return ioJson.FormatInventory
	case "Party":
		return ioJson.FormatParty
	case "Magic":
		return ioJson.FormatMagic
	case "Espers":
		return ioJson.FormatEspers
	case "Equipment":
		return ioJson.FormatEquipment
	default:
		return ioJson.FormatFull
	}
}
