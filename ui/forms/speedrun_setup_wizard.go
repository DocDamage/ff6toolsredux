package forms

import (
	"fmt"

	ipr "ffvi_editor/io/pr"
	"ffvi_editor/models/speedrun"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SpeedrunSetupWizard creates a dialog for setting up speedrun configurations
type SpeedrunSetupWizard struct {
	dialog       dialog.Dialog
	prData       *ipr.PR
	window       fyne.Window
	configs      map[string]*speedrun.SpeedrunConfig
	selectedCfg  *speedrun.SpeedrunConfig
	onApply      func(*speedrun.SpeedrunConfig) error
	categoryList *widget.List
	configList   *widget.List
}

// NewSpeedrunSetupWizard creates a new speedrun setup wizard from PR data
func NewSpeedrunSetupWizard(prData *ipr.PR, window fyne.Window) *SpeedrunSetupWizard {
	wizard := &SpeedrunSetupWizard{
		prData:  prData,
		window:  window,
		configs: speedrun.RegisteredConfigs,
	}
	wizard.Show()
	return wizard
}

// NewSpeedrunSetupWizardWithCallback creates a speedrun wizard with custom callback
func NewSpeedrunSetupWizardWithCallback(configs map[string]*speedrun.SpeedrunConfig, onApply func(*speedrun.SpeedrunConfig) error) *SpeedrunSetupWizard {
	wizard := &SpeedrunSetupWizard{
		configs: configs,
		onApply: onApply,
	}

	return wizard
}

// Show displays the speedrun setup wizard dialog
func (w *SpeedrunSetupWizard) Show() {
	parent := w.window
	if parent == nil {
		// If no window set, don't show
		return
	}
	// Categories
	categories := []string{"any%", "100%", "low_level", "solo", "glitchless", "pacifist"}

	// Category list
	w.categoryList = widget.NewList(
		func() int { return len(categories) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(categories[i])
		},
	)

	// Config list
	configIDs := make([]string, 0, len(w.configs))
	for id := range w.configs {
		configIDs = append(configIDs, id)
	}

	w.configList = widget.NewList(
		func() int { return len(configIDs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			if i < len(configIDs) {
				cfg := w.configs[configIDs[i]]
				obj.(*widget.Label).SetText(fmt.Sprintf("%s - %s", cfg.Name, cfg.Category))
			}
		},
	)

	w.configList.OnSelected = func(id widget.ListItemID) {
		if id < len(configIDs) {
			w.selectedCfg = w.configs[configIDs[id]]
		}
	}

	// Description
	descLabel := widget.NewLabel("")

	// Details panel - reserved for future use

	// Stats display
	statsLabel := widget.NewLabel("")

	updateDetails := func() {
		if w.selectedCfg != nil {
			descLabel.SetText(w.selectedCfg.Description)
			statsLabel.SetText(fmt.Sprintf("Category: %s\nCharacters: %d\nLevel Cap: %d\nSequence Breaks: %v",
				w.selectedCfg.Category,
				len(w.selectedCfg.Characters),
				w.selectedCfg.LevelCap,
				w.selectedCfg.AllowSequenceBreaks,
			))
		}
	}

	w.configList.OnSelected = func(id widget.ListItemID) {
		if id < len(configIDs) {
			w.selectedCfg = w.configs[configIDs[id]]
			updateDetails()
		}
	}

	// Two-column layout
	categoryPanel := container.NewBorder(
		widget.NewLabel("Categories"),
		nil,
		nil,
		nil,
		w.categoryList,
	)

	detailsPanel := container.NewVBox(
		widget.NewLabel("Available Configs"),
		w.configList,
	)

	infoPanel := container.NewVBox(
		widget.NewLabel("Description:"),
		descLabel,
		widget.NewLabel("Configuration:"),
		statsLabel,
	)

	// Buttons
	applyBtn := widget.NewButton("Apply Configuration", func() {
		if w.selectedCfg != nil {
			// Handle PR data persistence
			if w.prData != nil {
				err := speedrun.ApplyConfigToSave(w.selectedCfg)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to apply config: %w", err), parent)
					return
				}
				dialog.ShowInformation("Success", fmt.Sprintf("Applied %s configuration", w.selectedCfg.Name), parent)
				w.dialog.Hide()
			} else if w.onApply != nil {
				// Fallback to callback if no PR data
				if err := w.onApply(w.selectedCfg); err != nil {
					dialog.ShowError(err, parent)
					return
				}
				w.dialog.Hide()
			}
		}
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		w.dialog.Hide()
	})

	buttons := container.NewHBox(applyBtn, cancelBtn)

	// Main layout
	content := container.NewVBox(
		categoryPanel,
		detailsPanel,
		infoPanel,
		buttons,
	)

	w.dialog = dialog.NewCustom("Speedrun Setup Wizard", "Close", content, parent)
	w.dialog.Resize(fyne.NewSize(600, 500))
	w.dialog.Show()
}

// GetSelectedConfig returns the currently selected configuration
func (w *SpeedrunSetupWizard) GetSelectedConfig() *speedrun.SpeedrunConfig {
	return w.selectedCfg
}
