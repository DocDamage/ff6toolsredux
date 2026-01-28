package forms

import (
	"fmt"

	"ffvi_editor/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// PaletteViewerDialog displays a read-only palette extracted from ROM
type PaletteViewerDialog struct {
	window   fyne.Window
	dialog   dialog.Dialog
	palette  *models.Palette
	charName string

	paletteGrid *fyne.Container
	colorBoxes  []*canvas.Rectangle
	infoLabel   *widget.Label
}

// NewPaletteViewerDialog creates a new palette viewer dialog
func NewPaletteViewerDialog(parentWindow fyne.Window, palette *models.Palette, characterName string) *PaletteViewerDialog {
	pvd := &PaletteViewerDialog{
		window:   parentWindow,
		palette:  palette,
		charName: characterName,
	}

	// Create color boxes (16 colors in a 4x4 grid)
	pvd.paletteGrid = container.New(layout.NewGridLayout(4))
	pvd.colorBoxes = make([]*canvas.Rectangle, 16)

	for i := 0; i < 16; i++ {
		idx := i
		color := palette.Colors[i].ToColor()
		rect := canvas.NewRectangle(color)
		rect.SetMinSize(fyne.NewSize(60, 60))

		// Create button to show color info
		btn := widget.NewButton(fmt.Sprintf("%d", i), func() {
			pvd.showColorInfo(idx)
		})

		// Create a container with the color rectangle and button on top
		colorBox := container.New(layout.NewCenterLayout(), rect, btn)

		pvd.colorBoxes[i] = rect
		pvd.paletteGrid.Add(colorBox)
	}

	// Info label for selected color
	pvd.infoLabel = widget.NewLabel("Click on a color to see details")

	// Build content
	content := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel(""),
		widget.NewLabel(fmt.Sprintf("%s - ROM Palette", characterName)),
		widget.NewLabel(""),
		pvd.paletteGrid,
		widget.NewLabel(""),
		pvd.infoLabel,
	)

	scrollContainer := container.New(
		layout.NewVBoxLayout(),
		content,
	)

	// Create close button
	closeBtn := widget.NewButton("Close", func() {
		if pvd.dialog != nil {
			pvd.dialog.Hide()
		}
	})

	// Create export button
	exportBtn := widget.NewButton("Export Palette", func() {
		dialog.NewFileSave(func(writeCloser fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, pvd.window)
				return
			}
			if writeCloser == nil {
				return
			}
			defer writeCloser.Close()

			// Export as JSON or other format
			data := fmt.Sprintf(`{
  "name": "%s",
  "character": "%s",
  "colors": [
`, palette.Name, characterName)

			for i := 0; i < 16; i++ {
				c := palette.Colors[i]
				data += fmt.Sprintf(`    {"index": %d, "r": %d, "g": %d, "b": %d}`, i, c.R, c.G, c.B)
				if i < 15 {
					data += ",\n"
				} else {
					data += "\n"
				}
			}
			data += "  ]\n}\n"

			writeCloser.Write([]byte(data))
			dialog.ShowInformation("Export", "Palette exported successfully", pvd.window)
		}, pvd.window).Show()
	})

	buttons := container.New(layout.NewHBoxLayout(), exportBtn, layout.NewSpacer(), closeBtn)

	mainContent := container.New(
		layout.NewVBoxLayout(),
		scrollContainer,
		buttons,
	)

	// Create dialog
	pvd.dialog = dialog.NewCustom(
		fmt.Sprintf("Palette Viewer - %s", characterName),
		"",
		mainContent,
		pvd.window,
	)

	return pvd
}

// showColorInfo displays information about a selected color
func (pvd *PaletteViewerDialog) showColorInfo(index int) {
	c := pvd.palette.Colors[index]
	r888, g888, b888 := c.ToRGB888()
	hex := fmt.Sprintf("%02X%02X%02X", r888, g888, b888)
	info := fmt.Sprintf("Color %d: RGB(5-bit)=(%d,%d,%d) RGB(8-bit)=(%d,%d,%d) #%s",
		index, c.R, c.G, c.B, r888, g888, b888, hex)
	pvd.infoLabel.SetText(info)
}

// Show displays the palette viewer dialog
func (pvd *PaletteViewerDialog) Show() {
	if pvd.dialog != nil {
		pvd.dialog.Show()
	}
}
