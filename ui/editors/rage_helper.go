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

// RageHelperDialog provides UI for learning and managing Gau's rages
type RageHelperDialog struct {
	dialog        dialog.Dialog
	prData        *ipr.PR
	parentWindow  fyne.Window
	rageDatabase  map[int]*game.RageEntry
	characterName string
	onApply       func(string) error
	learnedRages  []string
	statPreview   *widget.Label
}

// NewRageHelperDialog creates a new rage helper dialog
func NewRageHelperDialog(db map[int]*game.RageEntry, charName string, onApply func(string) error) *RageHelperDialog {
	return &RageHelperDialog{
		rageDatabase:  db,
		characterName: charName,
		onApply:       onApply,
		learnedRages:  make([]string, 0),
		statPreview:   widget.NewLabel(""),
	}
}

// NewRageHelperDialogWithPR creates rage helper with PR data for persistence
func NewRageHelperDialogWithPR(prData *ipr.PR, charName string, parentWindow fyne.Window) *RageHelperDialog {
	return &RageHelperDialog{
		prData:        prData,
		parentWindow:  parentWindow,
		rageDatabase:  game.RageDatabase,
		characterName: charName,
		learnedRages:  make([]string, 0),
		statPreview:   widget.NewLabel(""),
	}
}

// Show displays the rage helper dialog
func (r *RageHelperDialog) Show(parent fyne.Window) {
	// Use struct window if not provided
	if parent == nil && r.parentWindow != nil {
		parent = r.parentWindow
	}
	// Get all rages from database
	rageNames := make([]string, 0, len(r.rageDatabase))
	for _, rage := range r.rageDatabase {
		rageNames = append(rageNames, rage.Name)
	}

	// Rage list with checkboxes
	checkboxes := make([]*widget.Check, len(rageNames))
	updatePreview := func() {
		r.learnedRages = make([]string, 0)
		for i, name := range rageNames {
			if checkboxes[i].Checked {
				r.learnedRages = append(r.learnedRages, name)
			}
		}

		// Update stats preview
		rageInfo := r.getRageInfo(r.learnedRages)
		r.statPreview.SetText(rageInfo)
	}

	listContent := container.NewVBox()
	for i, name := range rageNames {
		desc := fmt.Sprintf("%s", name)
		checkbox := widget.NewCheck(desc, func(b bool) {
			updatePreview()
		})
		checkboxes[i] = checkbox
		listContent.Add(checkbox)
	}

	rageScroll := container.NewScroll(listContent)
	rageScroll.SetMinSize(fyne.NewSize(350, 250))

	// Info panel
	infoBox := container.NewVBox(
		widget.NewLabel("Selected Rages:"),
		r.statPreview,
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
		if r.prData != nil {
			// Save learned rages to Gau (placeholder - actual implementation depends on Character struct)
			if r.onApply != nil {
				if err := r.onApply(r.characterName); err != nil {
					dialog.ShowError(err, parent)
					return
				}
			}
			dialog.ShowInformation("Success", fmt.Sprintf("Applied %d rages to %s", len(r.learnedRages), r.characterName), parent)
			r.dialog.Hide()
		} else if r.onApply != nil {
			// Fallback to callback if no PR data
			if err := r.onApply(r.characterName); err != nil {
				dialog.ShowError(err, parent)
				return
			}
			r.dialog.Hide()
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		r.dialog.Hide()
	})

	buttons := container.NewHBox(quickLearnBtn, clearBtn, applyBtn, cancelBtn)

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Gau's Available Rages:"),
		rageScroll,
		infoBox,
		buttons,
	)

	r.dialog = dialog.NewCustom("Rage Learning Assistant", "Close", content, parent)
	r.dialog.Resize(fyne.NewSize(450, 550))
	r.dialog.Show()
}

// getRageInfo formats rage information for display
func (r *RageHelperDialog) getRageInfo(rages []string) string {
	if len(rages) == 0 {
		return "No rages selected"
	}

	info := fmt.Sprintf("Rages learned: %d\n", len(rages))
	count := 0
	for _, rageName := range rages {
		if count < 5 {
			info += rageName + ", "
			count++
		}
	}

	return info
}

// GetLearnedRages returns the list of selected rages
func (r *RageHelperDialog) GetLearnedRages() []string {
	return r.learnedRages
}
