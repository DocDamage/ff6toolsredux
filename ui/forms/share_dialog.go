package forms

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ipr "ffvi_editor/io/pr"
	"ffvi_editor/models/share"
)

// ShareDialog handles sharing builds via share codes
type ShareDialog struct {
	window     fyne.Window
	charBuild  *share.CharacterBuild
	partyBuild *share.PartyBuild
	shareType  string // "character", "party"
}

// NewShareDialog creates a new share dialog for a character
func NewShareDialog(window fyne.Window, charBuild *share.CharacterBuild) *ShareDialog {
	return &ShareDialog{
		window:    window,
		charBuild: charBuild,
		shareType: "character",
	}
}

// NewShareDialogFromPR creates a share dialog from PR data
func NewShareDialogFromPR(window fyne.Window, prData *ipr.PR) *ShareDialog {
	// Build character build from first character in party
	// This is a simplified version - in production would extract specific character
	return &ShareDialog{
		window:    window,
		charBuild: &share.CharacterBuild{}, // Would populate from prData
		shareType: "character",
	}
}

// NewShareDialogParty creates a new share dialog for party
func NewShareDialogParty(window fyne.Window, partyBuild *share.PartyBuild) *ShareDialog {
	return &ShareDialog{
		window:     window,
		partyBuild: partyBuild,
		shareType:  "party",
	}
}

// Show displays the share dialog
func (sd *ShareDialog) Show() {
	content := container.NewVBox()

	// Title
	title := canvas.NewText("Share Build", color.White)
	title.TextSize = 18
	content.Add(title)

	// Info label
	infoLabel := widget.NewLabel("Generate a shareable code for this build:")
	content.Add(infoLabel)

	// Share code display
	codeContainer := container.NewVBox()

	// Generate code button
	generateBtn := widget.NewButton("Generate Share Code", func() {
		sd.generateAndDisplay(codeContainer)
	})
	content.Add(generateBtn)

	content.Add(codeContainer)

	// Add separator
	content.Add(widget.NewSeparator())

	// Import code section
	content.Add(widget.NewLabel("Or paste a code to import a build:"))

	// Code input
	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("TERRA-BLD-XXXXX or PTY-XXXX")
	content.Add(codeEntry)

	// Import button
	importBtn := widget.NewButton("Import Build", func() {
		if codeEntry.Text != "" {
			sd.importBuild(codeEntry.Text)
		} else {
			dialog.ShowError(fmt.Errorf("please enter a share code"), sd.window)
		}
	})
	content.Add(importBtn)

	// Add buttons
	content.Add(widget.NewSeparator())

	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Close", func() {
			// Dialog will close automatically
		}),
	)
	content.Add(buttonContainer)

	// Create and show dialog
	d := dialog.NewCustom("Share Build", "Close", content, sd.window)
	d.Show()
}

// generateAndDisplay generates a share code and displays it
func (sd *ShareDialog) generateAndDisplay(displayContainer *fyne.Container) {
	generator := share.NewCodeGenerator()

	var code string
	var err error

	if sd.shareType == "character" && sd.charBuild != nil {
		code, err = generator.GenerateCharacterCode(sd.charBuild, 0)
	} else if sd.shareType == "party" && sd.partyBuild != nil {
		code, err = generator.GeneratePartyCode(sd.partyBuild)
	} else {
		dialog.ShowError(fmt.Errorf("no valid build to share"), sd.window)
		return
	}

	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to generate share code: %w", err), sd.window)
		return
	}

	// Clear container
	displayContainer.Objects = make([]fyne.CanvasObject, 0)

	// Display code in box with copy button
	codeText := canvas.NewText(code, color.White)
	codeText.TextSize = 16

	copyBtn := widget.NewButton("Copy to Clipboard", func() {
		sd.window.Clipboard().SetContent(code)
		dialog.ShowInformation("Copied", "Share code copied to clipboard!", sd.window)
	})

	qrBtn := widget.NewButton("Show QR Code", func() {
		sd.showQRCode(code)
	})

	displayContainer.Add(widget.NewLabel("Your share code:"))
	displayContainer.Add(codeText)
	displayContainer.Add(container.NewHBox(
		layout.NewSpacer(),
		copyBtn,
		qrBtn,
		layout.NewSpacer(),
	))

	// Add note about code
	if sd.shareType == "character" {
		displayContainer.Add(widget.NewLabel("Share this code to let others import your character build."))
	} else {
		displayContainer.Add(widget.NewLabel("Share this code to let others import your party composition."))
	}

	displayContainer.Refresh()
}

// importBuild imports a build from a share code
func (sd *ShareDialog) importBuild(code string) {
	generator := share.NewCodeGenerator()

	// Try to determine code type
	if sd.shareType == "party" || (len(code) > 3 && code[0:3] == "PTY") {
		// Import party
		partyBuild, err := generator.DecodePartyCode(code)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to decode party code: %w", err), sd.window)
			return
		}

		dialog.ShowInformation("Party Imported",
			fmt.Sprintf("Imported party with %d members.\nPlease manually assign characters.",
				len(partyBuild.Characters)), sd.window)
	} else {
		// Import character
		charBuild, err := generator.DecodeCharacterCode(code)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to decode character code: %w", err), sd.window)
			return
		}

		dialog.ShowInformation("Character Build",
			fmt.Sprintf("Build: %s\nLevel: %d\nVigour: %d, Speed: %d, Stamina: %d",
				charBuild.Name, charBuild.Level, charBuild.Stats.Vigor,
				charBuild.Stats.Speed, charBuild.Stats.Stamina), sd.window)
	}
}

// showQRCode shows a QR code for the share code
func (sd *ShareDialog) showQRCode(code string) {
	// Note: Full QR code generation would require an external library
	// For now, show a dialog with the code in large text

	qrContent := container.NewVBox()

	title := canvas.NewText("Share Code (QR Format)", color.White)
	title.TextSize = 16
	qrContent.Add(title)

	// Display code in large text (ASCII art QR-like format)
	codeText := canvas.NewText(code, color.White)
	codeText.TextSize = 20
	codeText.TextStyle.Monospace = true
	qrContent.Add(codeText)

	qrContent.Add(widget.NewLabel("Scan this with your device or share this code directly."))

	copyBtn := widget.NewButton("Copy Code", func() {
		sd.window.Clipboard().SetContent(code)
		dialog.ShowInformation("Copied", "Code copied to clipboard!", sd.window)
	})

	qrContent.Add(copyBtn)

	d := dialog.NewCustom("QR Code", "Close", qrContent, sd.window)
	d.Show()
}
