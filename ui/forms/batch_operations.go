package forms

import (
	"fmt"

	ipr "ffvi_editor/io/pr"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// BatchOperationsDialog displays batch operations UI
type BatchOperationsDialog struct {
	dialog          dialog.Dialog
	prData          *ipr.PR
	window          fyne.Window
	selectedOps     []string
	previewText     *widget.Label
	executeCallback func([]string) error
}

// NewBatchOperationsDialog creates a new batch operations dialog from PR data
func NewBatchOperationsDialog(prData *ipr.PR, window fyne.Window) *BatchOperationsDialog {
	bod := &BatchOperationsDialog{
		prData:      prData,
		window:      window,
		selectedOps: make([]string, 0),
	}

	bod.previewText = widget.NewLabel("Select operations to see preview...")

	bod.Show()
	return bod
}

// NewBatchOperationsDialogWithCallback creates batch operations with custom callback
func NewBatchOperationsDialogWithCallback(parent fyne.Window, executeCallback func([]string) error) *BatchOperationsDialog {
	bod := &BatchOperationsDialog{
		window:          parent,
		selectedOps:     make([]string, 0),
		executeCallback: executeCallback,
	}

	bod.previewText = widget.NewLabel("Select operations to see preview...")

	// Character operations
	charOps := container.NewVBox(
		widget.NewCard("Character Operations", "", container.NewVBox(
			bod.createCheckOperation("Max All Stats", "maxStats", "Maximize HP, MP, and all stats for all characters"),
			bod.createCheckOperation("Max HP/MP", "maxHPMP", "Maximize HP and MP for all characters"),
			bod.createCheckOperation("Level to 99", "maxLevel", "Set all characters to level 99"),
			bod.createCheckOperation("Max Experience", "maxExp", "Give all characters max experience"),
		)),
	)

	// Inventory operations
	invOps := container.NewVBox(
		widget.NewCard("Inventory Operations", "", container.NewVBox(
			bod.createCheckOperation("Add 99 Consumables", "maxConsumables", "Add 99 of each consumable item"),
			bod.createCheckOperation("Clear Inventory", "clearInv", "Remove all items from inventory"),
			bod.createCheckOperation("Max All Items", "maxItems", "Set all items to max quantity"),
		)),
	)

	// Magic/Esper operations
	magicOps := container.NewVBox(
		widget.NewCard("Magic & Esper Operations", "", container.NewVBox(
			bod.createCheckOperation("Learn All Magic", "learnMagic", "Unlock all spells for all characters"),
			bod.createCheckOperation("Unlock All Espers", "unlockEspers", "Unlock all espers/summons"),
			bod.createCheckOperation("Max Esper Stats", "maxEsperStats", "Maximize all esper statistics"),
		)),
	)

	// Equipment operations
	equipOps := container.NewVBox(
		widget.NewCard("Equipment Operations", "", container.NewVBox(
			bod.createCheckOperation("Equip Best Gear", "bestGear", "Equip optimal equipment for all characters"),
			bod.createCheckOperation("Equip All Items", "allItems", "Equip all available items"),
		)),
	)

	// Scrollable operations list
	scrollOps := container.NewScroll(
		container.NewVBox(charOps, invOps, magicOps, equipOps),
	)

	// Preview panel
	previewCard := widget.NewCard(
		"Preview",
		"What will change",
		container.NewScroll(bod.previewText),
	)

	// Buttons
	applyBtn := widget.NewButton("Apply", func() {
		if len(bod.selectedOps) == 0 {
			dialog.ShowError(fmt.Errorf("select at least one operation"), bod.window)
			return
		}

		if bod.executeCallback != nil {
			if err := bod.executeCallback(bod.selectedOps); err != nil {
				dialog.ShowError(err, bod.window)
			}
		}
		bod.dialog.Hide()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		bod.dialog.Hide()
	})

	buttons := container.NewHBox(applyBtn, cancelBtn)

	// Main layout
	content := container.NewBorder(nil, buttons, nil, nil,
		container.NewHSplit(scrollOps, previewCard),
	)

	bod.dialog = dialog.NewCustom(
		"Batch Operations",
		"",
		content,
		bod.window,
	)

	return bod
}

// createCheckOperation creates a checkbox with description
func (bod *BatchOperationsDialog) createCheckOperation(label, opID, description string) fyne.CanvasObject {
	check := widget.NewCheck(label, func(b bool) {
		if b {
			// Add to selected operations
			bod.selectedOps = append(bod.selectedOps, opID)
		} else {
			// Remove from selected operations
			for i, op := range bod.selectedOps {
				if op == opID {
					bod.selectedOps = append(bod.selectedOps[:i], bod.selectedOps[i+1:]...)
					break
				}
			}
		}
		bod.updatePreview()
	})

	return container.NewVBox(
		check,
		widget.NewLabel("  "+description),
	)
}

// updatePreview updates the preview text
func (bod *BatchOperationsDialog) updatePreview() {
	preview := "Selected Operations:\n\n"

	if len(bod.selectedOps) == 0 {
		preview = "Select operations to see preview..."
		bod.previewText.SetText(preview)
		return
	}

	operationDescriptions := map[string]string{
		"maxStats":       "- Maximize all character stats (HP, MP, STR, DEF, etc.)\n",
		"maxHPMP":        "- Set HP to 9999 and MP to 999 for all characters\n",
		"maxLevel":       "- Level all characters to 99\n",
		"maxExp":         "- Give all characters max experience\n",
		"maxConsumables": "- Add 99 of each consumable item to inventory\n",
		"clearInv":       "- Remove all items from inventory\n",
		"maxItems":       "- Set all items to 99 quantity\n",
		"learnMagic":     "- All characters learn all available spells\n",
		"unlockEspers":   "- Unlock all espers/summons for all characters\n",
		"maxEsperStats":  "- Maximize all esper statistics\n",
		"bestGear":       "- Equip optimal equipment based on character stats\n",
		"allItems":       "- Equip all available items on characters\n",
	}

	for _, op := range bod.selectedOps {
		if desc, exists := operationDescriptions[op]; exists {
			preview += desc
		}
	}
	preview += fmt.Sprintf("\nTotal changes: %d operations selected\n", len(bod.selectedOps))
	preview += "Note: This action will create a single undo point.\n"

	bod.previewText.SetText(preview)
}

// Show displays the batch operations dialog
func (bod *BatchOperationsDialog) Show() {
	bod.dialog.Show()
}
