package forms

import (
	"fmt"
	"image"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ffvi_editor/io"
	"ffvi_editor/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// SpriteImportDialog provides a complete sprite import UI with preview and options
type SpriteImportDialog struct {
	window fyne.Window
	dialog dialog.Dialog

	// File selection
	fileEntry *widget.Entry
	browseBtn *widget.Button

	// Import options
	ditheringSelect    *widget.Select
	spriteTypeSelect   *widget.Select
	paddingCheckbox    *widget.Check
	autoTypeCheckbox   *widget.Check
	maxColorsEntry     *widget.Entry
	colorQualitySlider *widget.Slider

	// Preview
	previewContainer *fyne.Container
	previewImage     *canvas.Image
	previewLabel     *widget.Label
	statusLabel      *widget.Label

	// Buttons
	importBtn *widget.Button
	cancelBtn *widget.Button

	// State
	selectedFile    string
	importedSprite  *models.FF6Sprite
	onImportSuccess func(*models.FF6Sprite)
	onImportError   func(error)

	// Backend
	importer *io.SpriteImporter
}

// NewSpriteImportDialog creates a new sprite import dialog
func NewSpriteImportDialog(window fyne.Window) *SpriteImportDialog {
	s := &SpriteImportDialog{
		window:   window,
		importer: io.NewSpriteImporter(),
	}

	s.buildUI()
	return s
}

// buildUI constructs the complete dialog UI
func (s *SpriteImportDialog) buildUI() {
	// File selection section
	s.fileEntry = widget.NewEntry()
	s.fileEntry.SetPlaceHolder("Select image file...")
	s.fileEntry.OnChanged = func(str string) {
		s.selectedFile = str
	}

	s.browseBtn = widget.NewButton("Browse...", func() {
		s.showFileDialog()
	})

	fileSection := container.NewBorder(nil, nil, s.browseBtn, nil, s.fileEntry)

	// Import options section
	s.ditheringSelect = widget.NewSelect(
		[]string{"None", "Floyd-Steinberg", "Bayer"},
		func(val string) {},
	)
	s.ditheringSelect.SetSelected("Floyd-Steinberg")

	s.spriteTypeSelect = widget.NewSelect(
		[]string{"Auto-detect", "Character", "Battle", "Portrait", "NPC", "Enemy", "Overworld"},
		func(val string) {},
	)
	s.spriteTypeSelect.SetSelected("Auto-detect")

	s.autoTypeCheckbox = widget.NewCheck("Auto-detect sprite type", func(checked bool) {
		if checked {
			s.spriteTypeSelect.Disable()
			s.spriteTypeSelect.SetSelected("Auto-detect")
		} else {
			s.spriteTypeSelect.Enable()
		}
	})
	s.autoTypeCheckbox.SetChecked(true)

	s.paddingCheckbox = widget.NewCheck("Pad to tile boundaries", func(bool) {})
	s.paddingCheckbox.SetChecked(true)

	s.maxColorsEntry = widget.NewEntry()
	s.maxColorsEntry.SetText("16")
	s.maxColorsEntry.SetPlaceHolder("Max colors")

	s.colorQualitySlider = widget.NewSlider(0.5, 2.0)
	s.colorQualitySlider.Value = 1.0
	s.colorQualitySlider.Step = 0.1

	optionsForm := container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("Dithering:"), nil, s.ditheringSelect),
		container.NewBorder(nil, nil, widget.NewLabel("Sprite Type:"), nil, s.spriteTypeSelect),
		container.NewBorder(nil, nil, widget.NewLabel("Max Colors:"), nil, s.maxColorsEntry),
		s.autoTypeCheckbox,
		s.paddingCheckbox,
		widget.NewSeparator(),
		widget.NewLabel("Color Quality:"),
		s.colorQualitySlider,
	)

	// Preview section
	s.previewLabel = widget.NewLabel("Preview")
	s.previewLabel.Alignment = fyne.TextAlignCenter
	s.previewImage = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 128, 128)))
	s.previewImage.FillMode = canvas.ImageFillContain
	s.previewImage.SetMinSize(fyne.NewSize(256, 256))

	s.previewContainer = container.NewVBox(
		s.previewLabel,
		s.previewImage,
	)

	s.statusLabel = widget.NewLabel("Ready to import")
	s.statusLabel.Alignment = fyne.TextAlignCenter

	// Main layout: options on left, preview on right
	contentBox := container.NewHBox(
		container.NewVBox(
			widget.NewLabel("Import Options"),
			optionsForm,
			widget.NewSeparator(),
			s.statusLabel,
		),
		s.previewContainer,
	)

	// Buttons
	s.importBtn = widget.NewButton("Import", func() {
		s.performImport()
	})
	s.importBtn.Importance = widget.HighImportance

	s.cancelBtn = widget.NewButton("Cancel", func() {
		s.dialog.Hide()
	})

	buttons := container.NewHBox(
		layout.NewSpacer(),
		s.cancelBtn,
		s.importBtn,
	)

	// Main dialog content
	mainContent := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Import Sprite Image"),
			widget.NewSeparator(),
			fileSection,
			widget.NewSeparator(),
		),
		buttons,
		nil,
		nil,
		contentBox,
	)

	s.dialog = dialog.NewCustom(
		"Import Sprite",
		"Cancel",
		mainContent,
		s.window,
	)
	s.dialog.Resize(fyne.NewSize(900, 600))
}

// Show displays the import dialog
func (s *SpriteImportDialog) Show() {
	s.dialog.Show()
}

// Hide hides the import dialog
func (s *SpriteImportDialog) Hide() {
	s.dialog.Hide()
}

// OnImportSuccess sets the callback for successful import
func (s *SpriteImportDialog) OnImportSuccess(fn func(*models.FF6Sprite)) {
	s.onImportSuccess = fn
}

// OnImportError sets the callback for import errors
func (s *SpriteImportDialog) OnImportError(fn func(error)) {
	s.onImportError = fn
}

// showFileDialog opens the file browser
func (s *SpriteImportDialog) showFileDialog() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, s.window)
			return
		}

		if reader == nil {
			return
		}

		path := reader.URI().Path()
		reader.Close()

		s.fileEntry.SetText(path)
		s.selectedFile = path

		// Load and preview the image
		go s.updatePreview()
	}, s.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".gif", ".bmp", ".jpg", ".jpeg"}))

	fd.Show()
}

// updatePreview loads and displays a preview of the selected image
func (s *SpriteImportDialog) updatePreview() {
	if s.selectedFile == "" {
		return
	}

	s.statusLabel.SetText("Loading preview...")

	decoder := io.NewImageDecoder()
	img, _, err := decoder.Decode(s.selectedFile)
	if err != nil {
		s.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	s.previewImage = canvas.NewImageFromImage(img)
	s.previewImage.FillMode = canvas.ImageFillContain
	s.previewImage.SetMinSize(fyne.NewSize(256, 256))

	bounds := img.Bounds()
	s.previewLabel.SetText(fmt.Sprintf("Preview (%dx%d)", bounds.Dx(), bounds.Dy()))

	s.previewContainer.RemoveAll()
	s.previewContainer.Add(s.previewLabel)
	s.previewContainer.Add(s.previewImage)
	s.previewContainer.Refresh()

	s.statusLabel.SetText("Ready to import")
}

// performImport executes the sprite import
func (s *SpriteImportDialog) performImport() {
	if s.selectedFile == "" {
		dialog.ShowError(fmt.Errorf("please select an image file"), s.window)
		return
	}

	s.statusLabel.SetText("Importing... please wait")
	s.importBtn.Disable()
	s.browseBtn.Disable()

	go func() {
		defer func() {
			s.importBtn.Enable()
			s.browseBtn.Enable()
		}()

		opts := models.NewSpriteImportOptions()
		opts.DitherMethod = s.getDitheringMethod()
		opts.AutoPadding = s.paddingCheckbox.Checked
		opts.AutoDetectType = s.autoTypeCheckbox.Checked

		if maxColors, err := strconv.Atoi(strings.TrimSpace(s.maxColorsEntry.Text)); err == nil && maxColors > 0 {
			opts.MaxColors = maxColors
		}

		if !opts.AutoDetectType {
			opts.TargetType = s.getSelectedSpriteType()
		}

		opts.SourcePath = s.selectedFile
		opts.SourceFormat = strings.TrimPrefix(strings.ToLower(filepath.Ext(s.selectedFile)), ".")

		result := s.importer.Import(opts)

		if !result.Success {
			errMsg := "import failed"
			if len(result.Errors) > 0 {
				errMsg = result.Errors[0].Message
			}
			err := fmt.Errorf(errMsg)
			s.statusLabel.SetText(fmt.Sprintf("Import failed: %v", err))
			if s.onImportError != nil {
				s.onImportError(err)
			}
			dialog.ShowError(err, s.window)
			return
		}

		s.importedSprite = result.Sprite
		s.statusLabel.SetText(fmt.Sprintf("✓ Import successful (%dms)", result.Duration.Milliseconds()))

		if s.onImportSuccess != nil {
			s.onImportSuccess(result.Sprite)
		}

		time.Sleep(500 * time.Millisecond)
		s.dialog.Hide()
	}()
}

// getDitheringMethod converts the UI selection to backend enum
func (s *SpriteImportDialog) getDitheringMethod() string {
	switch s.ditheringSelect.Selected {
	case "Floyd-Steinberg":
		return "floyd-steinberg"
	case "Bayer":
		return "bayer"
	default:
		return "none"
	}
}

// GetImportedSprite returns the last imported sprite
func (s *SpriteImportDialog) GetImportedSprite() *models.FF6Sprite {
	return s.importedSprite
}

func (s *SpriteImportDialog) getSelectedSpriteType() models.SpriteType {
	switch s.spriteTypeSelect.Selected {
	case "Character":
		return models.SpriteTypeCharacter
	case "Battle":
		return models.SpriteTypeBattle
	case "Portrait":
		return models.SpriteTypePortrait
	case "NPC":
		return models.SpriteTypeNPC
	case "Enemy":
		return models.SpriteTypeEnemy
	case "Overworld":
		return models.SpriteTypeOverworld
	default:
		return models.SpriteTypeCharacter
	}
}
