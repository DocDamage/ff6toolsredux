package forms

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io"
	"ffvi_editor/models"
)

// AnimationExportDialog provides export options for animations
type AnimationExportDialog struct {
	dialog         *widget.PopUp
	animation      *models.AnimationData
	exporter       *io.AnimationExporter
	container      *fyne.Container
	formatSelect   *widget.Select
	scaleEntry     *widget.Entry
	ditherCheckbox *widget.Check
	bgColorEntry   *widget.Entry
	qualitySlider  *widget.Slider
	columnsEntry   *widget.Entry
	onClose        func()
}

// NewAnimationExportDialog creates a new animation export dialog
func NewAnimationExportDialog(animation *models.AnimationData) *AnimationExportDialog {
	if animation == nil || len(animation.Frames) == 0 {
		return nil
	}

	aed := &AnimationExportDialog{
		animation: animation,
	}

	// Create exporter
	aed.exporter = io.NewAnimationExporter(animation)
	if aed.exporter == nil {
		return nil
	}

	aed.buildUI()
	return aed
}

// buildUI constructs the dialog UI
func (aed *AnimationExportDialog) buildUI() {
	// Export format selection
	aed.formatSelect = widget.NewSelect(
		[]string{"GIF (Animated)", "PNG (Frames)", "JSON (Metadata)"},
		aed.onFormatChanged,
	)
	aed.formatSelect.SetSelected("GIF (Animated)")

	// Scale factor
	aed.scaleEntry = widget.NewEntry()
	aed.scaleEntry.SetText("1")
	aed.scaleEntry.SetPlaceHolder("1")

	// Dithering
	aed.ditherCheckbox = widget.NewCheck("Apply Dithering (GIF)", func(b bool) {})

	// Background color (for PNG frames)
	aed.bgColorEntry = widget.NewEntry()
	aed.bgColorEntry.SetText("000000")
	aed.bgColorEntry.SetPlaceHolder("RRGGBB hex")

	// Quality slider (for GIF)
	aed.qualitySlider = widget.NewSlider(1, 100)
	aed.qualitySlider.Value = 75
	qualityLabel := widget.NewLabel("Quality: 75")
	aed.qualitySlider.OnChanged = func(v float64) {
		qualityLabel.SetText(fmt.Sprintf("Quality: %d", int(v)))
	}

	// Sprite sheet columns
	aed.columnsEntry = widget.NewEntry()
	aed.columnsEntry.SetText("4")
	aed.columnsEntry.SetPlaceHolder("4")

	// Format options container
	gifOptionsBox := container.NewVBox(
		widget.NewLabel("GIF Options:"),
		container.NewVBox(
			widget.NewLabel("Quality:"),
			aed.qualitySlider,
			qualityLabel,
		),
		aed.ditherCheckbox,
	)

	pngOptionsBox := container.NewVBox(
		widget.NewLabel("PNG Options:"),
		container.NewVBox(
			widget.NewLabel("Background Color (Hex):"),
			aed.bgColorEntry,
		),
	)

	generalOptionsBox := container.NewVBox(
		widget.NewLabel("General Options:"),
		container.NewVBox(
			widget.NewLabel("Scale Factor:"),
			aed.scaleEntry,
		),
		container.NewVBox(
			widget.NewLabel("Sprite Sheet Columns:"),
			aed.columnsEntry,
		),
	)

	// Build layout
	aed.container = container.NewVBox(
		widget.NewLabel("Animation Export Options"),
		widget.NewCard("Export Format", "", container.NewVBox(
			aed.formatSelect,
		)),
		widget.NewCard("Format-Specific", "", container.NewVBox(
			gifOptionsBox,
			pngOptionsBox,
		)),
		widget.NewCard("General", "", generalOptionsBox),
	)
}

// Show displays the dialog
func (aed *AnimationExportDialog) Show(parent fyne.Window) {
	exportBtn := widget.NewButton("Export", func() {
		aed.onExport(parent)
	})

	dialog := container.NewVBox(
		aed.container,
		container.NewHBox(
			exportBtn,
			widget.NewButton("Close", func() {
				if aed.onClose != nil {
					aed.onClose()
				}
				aed.Hide()
			}),
		),
	)

	aed.dialog = widget.NewPopUp(dialog, parent.Canvas())
	aed.dialog.ShowAtPosition(fyne.NewPos(200, 200))
}

// Hide hides the dialog
func (aed *AnimationExportDialog) Hide() {
	if aed.dialog != nil {
		aed.dialog.Hide()
	}
}

// Event handlers

func (aed *AnimationExportDialog) onFormatChanged(value string) {
	// Format changed - could update UI
}

func (aed *AnimationExportDialog) onExport(parent fyne.Window) {
	format := aed.formatSelect.Selected

	// Show file/folder chooser based on format
	switch format {
	case "GIF (Animated)":
		aed.exportGIF(parent)
	case "PNG (Frames)":
		aed.exportPNG(parent)
	case "JSON (Metadata)":
		aed.exportJSON(parent)
	}
}

func (aed *AnimationExportDialog) exportGIF(parent fyne.Window) {
	fdialog := dialog.NewFileSave(
		func(uri fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			if uri == nil {
				return
			}

			defer uri.Close()

			// Get options
			scale, _ := strconv.Atoi(aed.scaleEntry.Text)
			if scale < 1 {
				scale = 1
			}

			quality := int(aed.qualitySlider.Value)
			dither := aed.ditherCheckbox.Checked

			opts := &io.AnimationExportOptions{
				Scale:   scale,
				Quality: quality,
				Dither:  dither,
			}

			// Export
			err = aed.exporter.ExportGIF(uri.URI().Path(), opts)
			if err != nil {
				dialog.ShowError(err, parent)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Animation exported to %s", uri.URI().Name()), parent)
			}
		},
		parent,
	)

	fdialog.SetFileName("animation.gif")
	fdialog.Show()
}

func (aed *AnimationExportDialog) exportPNG(parent fyne.Window) {
	fdialog := dialog.NewFolderOpen(
		func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			if uri == nil {
				return
			}

			// Get options
			scale, _ := strconv.Atoi(aed.scaleEntry.Text)
			if scale < 1 {
				scale = 1
			}

			opts := &io.AnimationExportOptions{
				Scale: scale,
			}

			// Export frames
			err = aed.exporter.ExportFramesPNG(uri.Path(), opts)
			if err != nil {
				dialog.ShowError(err, parent)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Frames exported to %s", uri.Name()), parent)
			}
		},
		parent,
	)

	fdialog.Show()
}

func (aed *AnimationExportDialog) exportJSON(parent fyne.Window) {
	fdialog := dialog.NewFileSave(
		func(uri fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			if uri == nil {
				return
			}

			defer uri.Close()

			// Export metadata
			err = aed.exporter.ExportJSON(uri.URI().Path())
			if err != nil {
				dialog.ShowError(err, parent)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Metadata exported to %s", uri.URI().Name()), parent)
			}
		},
		parent,
	)

	fdialog.SetFileName("animation.json")
	fdialog.Show()
}

// Helper methods

func (aed *AnimationExportDialog) getBackgroundColor() [3]uint8 {
	// Parse hex color
	hexStr := aed.bgColorEntry.Text
	if len(hexStr) != 6 {
		return [3]uint8{255, 255, 255}
	}

	var r, g, b uint8
	fmt.Sscanf(hexStr, "%02x%02x%02x", &r, &g, &b)
	return [3]uint8{r, g, b}
}

// SetOnClose sets the close callback
func (aed *AnimationExportDialog) SetOnClose(callback func()) {
	aed.onClose = callback
}
