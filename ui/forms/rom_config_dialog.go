package forms

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// ROMConfigDialog allows the user to configure the ROM path
type ROMConfigDialog struct {
	window         fyne.Window
	romPathEntry   *widget.Entry
	onSaveCallback func(romPath string)
}

// NewROMConfigDialog creates a new ROM configuration dialog
func NewROMConfigDialog(window fyne.Window, currentROMPath string) *ROMConfigDialog {
	d := &ROMConfigDialog{
		window:       window,
		romPathEntry: widget.NewEntry(),
	}

	d.romPathEntry.SetText(currentROMPath)
	d.romPathEntry.SetPlaceHolder("Path to Final Fantasy VI ROM file...")

	return d
}

// OnSave sets the callback for when the user saves the ROM path
func (d *ROMConfigDialog) OnSave(callback func(romPath string)) {
	d.onSaveCallback = callback
}

// Show displays the ROM configuration dialog
func (d *ROMConfigDialog) Show() {
	browseButton := widget.NewButton("Browse...", func() {
		fileDialog := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, d.window)
				return
			}
			if uc == nil {
				return
			}
			defer uc.Close()

			path := uc.URI().Path()
			d.romPathEntry.SetText(path)
		}, d.window)

		// Set file filter for ROM files
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".smc", ".sfc", ".fig"}))
		fileDialog.Show()
	})

	instructionsLabel := widget.NewLabel(
		"Select the Final Fantasy VI ROM file (released as FF3 in USA).\n" +
			"This is required to load character sprites for animation playback and editing.\n\n" +
			"Supported ROM formats: .smc, .sfc, .fig\n" +
			"Expected ROM size: 3-4 MB")
	instructionsLabel.Wrapping = fyne.TextWrapWord

	pathContainer := container.NewBorder(nil, nil, nil, browseButton, d.romPathEntry)

	content := container.NewVBox(
		instructionsLabel,
		widget.NewSeparator(),
		widget.NewLabel("ROM Path:"),
		pathContainer,
	)

	confirmDialog := dialog.NewCustomConfirm(
		"ROM Configuration",
		"Save",
		"Cancel",
		content,
		func(confirmed bool) {
			if confirmed && d.onSaveCallback != nil {
				romPath := d.romPathEntry.Text
				// Validate that the path exists
				if romPath != "" {
					// Clean the path
					romPath = filepath.Clean(romPath)
					d.onSaveCallback(romPath)
					dialog.ShowInformation(
						"ROM Configuration Saved",
						fmt.Sprintf("ROM path configured:\n%s\n\nRestart the application to apply changes.", romPath),
						d.window,
					)
				} else {
					d.onSaveCallback("")
					dialog.ShowInformation(
						"ROM Configuration Cleared",
						"ROM path has been cleared. Sprites will not be loaded from ROM.",
						d.window,
					)
				}
			}
		},
		d.window,
	)

	confirmDialog.Resize(fyne.NewSize(600, 300))
	confirmDialog.Show()
}
