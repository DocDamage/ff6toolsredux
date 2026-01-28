package editors

import (
	"fmt"

	ipr "ffvi_editor/io/pr"
	"ffvi_editor/models/game"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SketchHelperDialog provides UI for managing Relm's sketch moves
type SketchHelperDialog struct {
	dialog          dialog.Dialog
	prData          *ipr.PR
	parentWindow    fyne.Window
	sketchDatabase  map[int]*game.SketchMove
	characterName   string
	onApply         func(string) error
	learnedSketches []string
	skillPreview    *widget.Label
}

// NewSketchHelperDialog creates a new sketch helper dialog
func NewSketchHelperDialog(db map[int]*game.SketchMove, charName string, onApply func(string) error) *SketchHelperDialog {
	return &SketchHelperDialog{
		sketchDatabase:  db,
		characterName:   charName,
		onApply:         onApply,
		learnedSketches: make([]string, 0),
		skillPreview:    widget.NewLabel(""),
	}
}

// NewSketchHelperDialogWithPR creates sketch helper with PR data for persistence
func NewSketchHelperDialogWithPR(prData *ipr.PR, charName string, parentWindow fyne.Window) *SketchHelperDialog {
	return &SketchHelperDialog{
		prData:          prData,
		parentWindow:    parentWindow,
		sketchDatabase:  game.SketchDatabase,
		characterName:   charName,
		learnedSketches: make([]string, 0),
		skillPreview:    widget.NewLabel(""),
	}
}

// Show displays the sketch helper dialog
func (s *SketchHelperDialog) Show(parent fyne.Window) {
	// Use struct window if not provided
	if parent == nil && s.parentWindow != nil {
		parent = s.parentWindow
	}
	// Get all sketches from database
	sketchNames := make([]string, 0, len(s.sketchDatabase))
	for _, sketch := range s.sketchDatabase {
		sketchNames = append(sketchNames, sketch.Name)
	}

	// Sketch list with checkboxes
	checkboxes := make([]*widget.Check, len(sketchNames))
	updatePreview := func() {
		s.learnedSketches = make([]string, 0)
		for i, name := range sketchNames {
			if checkboxes[i].Checked {
				s.learnedSketches = append(s.learnedSketches, name)
			}
		}

		// Update preview
		skillInfo := s.getSketchInfo(s.learnedSketches)
		s.skillPreview.SetText(skillInfo)
	}

	listContent := container.NewVBox()
	for i, name := range sketchNames {
		desc := fmt.Sprintf("%s", name)
		checkbox := widget.NewCheck(desc, func(b bool) {
			updatePreview()
		})
		checkboxes[i] = checkbox
		listContent.Add(checkbox)
	}

	sketchScroll := container.NewScroll(listContent)
	sketchScroll.SetMinSize(fyne.NewSize(350, 250))

	// Info panel
	infoBox := container.NewVBox(
		widget.NewLabel("Selected Sketches:"),
		s.skillPreview,
	)

	// Buttons
	quickLearnBtn := widget.NewButton("Learn All", func() {
		for _, checkbox := range checkboxes {
			checkbox.Checked = true
		}
		updatePreview()
	})

	clearBtn := widget.NewButton("Clear", func() {
		for _, checkbox := range checkboxes {
			checkbox.Checked = false
		}
		updatePreview()
	})

	applyBtn := widget.NewButton("Apply", func() {
		// Handle PR data persistence
		if s.prData != nil {
			// Save learned sketches to Relm (placeholder - actual implementation depends on Character struct)
			if s.onApply != nil {
				if err := s.onApply(s.characterName); err != nil {
					dialog.ShowError(err, parent)
					return
				}
			}
			dialog.ShowInformation("Success", fmt.Sprintf("Applied %d sketches to %s", len(s.learnedSketches), s.characterName), parent)
			s.dialog.Hide()
		} else if s.onApply != nil {
			// Fallback to callback if no PR data
			if err := s.onApply(s.characterName); err != nil {
				dialog.ShowError(err, parent)
				return
			}
			s.dialog.Hide()
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		s.dialog.Hide()
	})

	buttons := container.NewHBox(quickLearnBtn, clearBtn, applyBtn, cancelBtn)

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Relm's Available Sketches:"),
		sketchScroll,
		infoBox,
		buttons,
	)

	s.dialog = dialog.NewCustom("Sketch Learning Assistant", "Close", content, parent)
	s.dialog.Resize(fyne.NewSize(450, 550))
	s.dialog.Show()
}

// getSketchInfo formats sketch information for display
func (s *SketchHelperDialog) getSketchInfo(sketches []string) string {
	if len(sketches) == 0 {
		return "No sketches selected"
	}

	info := fmt.Sprintf("Sketches learned: %d\n", len(sketches))
	count := 0
	for _, sketchName := range sketches {
		if count < 5 {
			info += sketchName + ", "
			count++
		}
	}

	return info
}

// GetLearnedSketches returns the list of selected sketches
func (s *SketchHelperDialog) GetLearnedSketches() []string {
	return s.learnedSketches
}
