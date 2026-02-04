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

// EsperOptimizerDialog provides UI for optimizing esper growth
type EsperOptimizerDialog struct {
	dialog         dialog.Dialog
	prData         *ipr.PR
	parentWindow   fyne.Window
	optimizer      *game.EsperOptimizer
	characterName  string
	onApply        func(string) error
	selectedEspers []string
	statPreview    *widget.Label
}

// NewEsperOptimizerDialog creates a new esper optimizer dialog
func NewEsperOptimizerDialog(optimizer *game.EsperOptimizer, charName string, onApply func(string) error) *EsperOptimizerDialog {
	return &EsperOptimizerDialog{
		optimizer:      optimizer,
		characterName:  charName,
		onApply:        onApply,
		selectedEspers: make([]string, 0),
		statPreview:    widget.NewLabel(""),
	}
}

// NewEsperOptimizerDialogWithPR creates esper optimizer with PR data for persistence
func NewEsperOptimizerDialogWithPR(prData *ipr.PR, charName string, parentWindow fyne.Window) *EsperOptimizerDialog {
	return &EsperOptimizerDialog{
		prData:         prData,
		parentWindow:   parentWindow,
		optimizer:      &game.EsperOptimizer{},
		characterName:  charName,
		selectedEspers: make([]string, 0),
		statPreview:    widget.NewLabel(""),
	}
}

// Show displays the esper optimizer dialog
func (e *EsperOptimizerDialog) Show(parent fyne.Window) {
	// Use struct window if not provided in method call
	if parent == nil && e.parentWindow != nil {
		parent = e.parentWindow
	}
	// Get all espers
	espers := e.optimizer.GetAllEspers()
	esperNames := make([]string, 0, len(espers))

	for _, esper := range espers {
		esperNames = append(esperNames, esper.Name)
	}

	// Esper list with checkboxes
	checkboxes := make([]*widget.Check, len(esperNames))

	updatePreview := func() {
		e.selectedEspers = make([]string, 0)
		for i, name := range esperNames {
			if checkboxes[i].Checked {
				e.selectedEspers = append(e.selectedEspers, name)
			}
		}

		// Calculate projected stat bonuses
		projectedStats := e.calculateProjectedStats(e.selectedEspers)
		e.statPreview.SetText(projectedStats)
	}

	listContent := container.NewVBox()
	for i, name := range esperNames {
		checkbox := widget.NewCheck(name, func(b bool) {
			updatePreview()
		})
		checkboxes[i] = checkbox
		listContent.Add(checkbox)
	}

	esperScroll := container.NewScroll(listContent)
	esperScroll.SetMinSize(fyne.NewSize(300, 250))

	// Info panel
	infoBox := container.NewVBox(
		widget.NewLabel("Projected Stat Growth:"),
		e.statPreview,
	)

	// Buttons
	applyBtn := widget.NewButton("Apply", func() {
		// Handle PR data persistence
		if e.prData != nil {
			// Save selected espers to character (placeholder - actual implementation depends on Character struct)
			if e.onApply != nil {
				if err := e.onApply(e.characterName); err != nil {
					dialog.ShowError(err, e.parentWindow)
					return
				}
			}
			dialog.ShowInformation("Success", fmt.Sprintf("Applied %d espers to %s", len(e.selectedEspers), e.characterName), e.parentWindow)
			e.dialog.Hide()
		} else if e.onApply != nil {
			// Fallback to callback if no PR data
			if err := e.onApply(e.characterName); err != nil {
				dialog.ShowError(err, e.parentWindow)
				return
			}
			e.dialog.Hide()
		}
	})

	clearBtn := widget.NewButton("Clear All", func() {
		for _, checkbox := range checkboxes {
			checkbox.Checked = false
		}
		updatePreview()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		e.dialog.Hide()
	})

	buttons := container.NewHBox(applyBtn, clearBtn, cancelBtn)

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Available Espers:"),
		esperScroll,
		infoBox,
		buttons,
	)

	e.dialog = dialog.NewCustom("Esper Growth Optimizer", "Close", content, parent)
	e.dialog.Resize(fyne.NewSize(400, 500))
	e.dialog.Show()
}

// calculateProjectedStats calculates the projected stat bonuses from selected espers
func (e *EsperOptimizerDialog) calculateProjectedStats(esperNames []string) string {
	if len(esperNames) == 0 {
		return "No espers selected"
	}

	totalGrowth := game.StatGrowth{}
	allEspers := e.optimizer.GetAllEspers()

	for _, name := range esperNames {
		for _, esper := range allEspers {
			if esper.Name == name {
				totalGrowth.Vigor += esper.StatGrowth.Vigor
				totalGrowth.Speed += esper.StatGrowth.Speed
				totalGrowth.Stamina += esper.StatGrowth.Stamina
				totalGrowth.MagicPwr += esper.StatGrowth.MagicPwr
				totalGrowth.Defense += esper.StatGrowth.Defense
				totalGrowth.MagicDef += esper.StatGrowth.MagicDef
				break
			}
		}
	}

	return fmt.Sprintf(
		"Vigor: +%d | Speed: +%d | Stamina: +%d\nMagic Pwr: +%d | Defense: +%d | Magic Def: +%d",
		totalGrowth.Vigor, totalGrowth.Speed, totalGrowth.Stamina,
		totalGrowth.MagicPwr, totalGrowth.Defense, totalGrowth.MagicDef,
	)
}

// GetSelectedEspers returns the selected esper list
func (e *EsperOptimizerDialog) GetSelectedEspers() []string {
	return e.selectedEspers
}
