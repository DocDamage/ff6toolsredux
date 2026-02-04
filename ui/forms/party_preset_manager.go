package forms

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io/presets"
	"ffvi_editor/models"
)

// PartyPresetDialog manages party preset selection and creation
type PartyPresetDialog struct {
	dialog  dialog.Dialog
	manager *presets.Manager
	preset  *models.PartyPreset
	onApply func(*models.PartyPreset)
	onClose func()
}

// NewPartyPresetDialog creates a new party preset dialog
func NewPartyPresetDialog(parent fyne.Window, manager *presets.Manager, onApply func(*models.PartyPreset)) *PartyPresetDialog {
	ppd := &PartyPresetDialog{
		manager: manager,
		onApply: onApply,
	}

	// Create tabs
	presetsTab := ppd.createPresetsTab()
	favoritesTab := ppd.createFavoritesTab()
	newPresetTab := ppd.createNewPresetTab(parent)

	tabs := container.NewAppTabs(
		container.NewTabItem("Presets", presetsTab),
		container.NewTabItem("Favorites", favoritesTab),
		container.NewTabItem("New Preset", newPresetTab),
	)

	// Apply and Cancel buttons
	applyBtn := widget.NewButton("Apply", func() {
		if ppd.preset != nil {
			if ppd.onApply != nil {
				ppd.onApply(ppd.preset)
			}
			ppd.dialog.Hide()
		}
	})

	closeBtn := widget.NewButton("Close", func() {
		if ppd.onClose != nil {
			ppd.onClose()
		}
		ppd.dialog.Hide()
	})

	buttons := container.NewHBox(applyBtn, closeBtn)
	content := container.NewBorder(nil, buttons, nil, nil, tabs)

	ppd.dialog = dialog.NewCustom(
		"Party Presets",
		"",
		content,
		parent,
	)

	return ppd
}

// createPresetsTab creates the main presets tab
func (ppd *PartyPresetDialog) createPresetsTab() fyne.CanvasObject {
	// Get presets (they're stored in manager's map)
	allPresets := ppd.manager.ListPresets("")

	presetList := widget.NewList(
		func() int { return len(allPresets) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Preset Name")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < len(allPresets) {
				item.(*widget.Label).SetText(allPresets[id].Name)
			}
		},
	)

	presetList.OnSelected = func(id widget.ListItemID) {
		if id < len(allPresets) {
			ppd.preset = allPresets[id]
		}
	}

	return container.NewVBox(
		widget.NewLabel("Select a preset to apply:"),
		presetList,
	)
}

// createFavoritesTab creates the favorites tab
func (ppd *PartyPresetDialog) createFavoritesTab() fyne.CanvasObject {
	// Get all presets and filter for favorites
	allPresets := ppd.manager.ListPresets("")
	favorites := make([]*models.PartyPreset, 0)
	for _, p := range allPresets {
		if p.Favorite {
			favorites = append(favorites, p)
		}
	}

	if len(favorites) == 0 {
		return widget.NewLabel("No favorite presets yet.\nMark presets as favorites in the main Presets tab.")
	}

	favoriteList := widget.NewList(
		func() int { return len(favorites) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Favorite Preset")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < len(favorites) {
				item.(*widget.Label).SetText(fmt.Sprintf("★ %s", favorites[id].Name))
			}
		},
	)

	favoriteList.OnSelected = func(id widget.ListItemID) {
		if id < len(favorites) {
			ppd.preset = favorites[id]
		}
	}

	return container.NewVBox(
		widget.NewLabel("Quick access to your favorite presets:"),
		favoriteList,
	)
}

// createNewPresetTab creates the new preset tab
func (ppd *PartyPresetDialog) createNewPresetTab(parent fyne.Window) fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Preset name")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Description (optional)")

	// Character pickers (simplified)
	char1Picker := widget.NewSelect(
		[]string{"Select character 1", "Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"},
		func(s string) {},
	)
	char1Picker.SetSelected("Select character 1")

	char2Picker := widget.NewSelect(
		[]string{"Select character 2", "Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"},
		func(s string) {},
	)
	char2Picker.SetSelected("Select character 2")

	char3Picker := widget.NewSelect(
		[]string{"Select character 3", "Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"},
		func(s string) {},
	)
	char3Picker.SetSelected("Select character 3")

	char4Picker := widget.NewSelect(
		[]string{"Select character 4", "Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin", "Celes", "Strago", "Relm", "Setzer", "Mog", "Gau", "Gogo", "Umaro"},
		func(s string) {},
	)
	char4Picker.SetSelected("Select character 4")

	saveBtn := widget.NewButton("Save Preset", func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("preset name required"), parent)
			return
		}

		// Map character names to indices
		charToIndex := map[string]uint8{
			"Terra": 12, "Locke": 3, "Cyan": 1, "Shadow": 10,
			"Edgar": 2, "Sabin": 8, "Celes": 0, "Strago": 11,
			"Relm": 7, "Setzer": 9, "Mog": 6, "Gau": 4,
			"Gogo": 5, "Umaro": 13,
		}

		// Build party members array
		members := [4]uint8{255, 255, 255, 255} // 255 = empty slot
		pickers := []*widget.Select{char1Picker, char2Picker, char3Picker, char4Picker}

		for i, picker := range pickers {
			if idx, ok := charToIndex[picker.Selected]; ok {
				members[i] = idx
			}
		}

		// Check if at least one member is selected
		hasMember := false
		for _, m := range members {
			if m != 255 {
				hasMember = true
				break
			}
		}
		if !hasMember {
			dialog.ShowError(fmt.Errorf("select at least one party member"), parent)
			return
		}

		// Create and save the preset
		preset := &models.PartyPreset{
			Name:        nameEntry.Text,
			Description: descEntry.Text,
			Members:     members,
		}

		if err := ppd.manager.CreatePreset(preset); err != nil {
			dialog.ShowError(err, parent)
			return
		}

		dialog.ShowInformation("Success", fmt.Sprintf("Preset '%s' saved", nameEntry.Text), parent)
		nameEntry.SetText("")
		descEntry.SetText("")

		// Reset pickers
		char1Picker.SetSelected("Select character 1")
		char2Picker.SetSelected("Select character 2")
		char3Picker.SetSelected("Select character 3")
		char4Picker.SetSelected("Select character 4")
	})

	return container.NewVBox(
		widget.NewCard("Preset Details", "", container.NewVBox(
			nameEntry,
			descEntry,
		)),
		widget.NewCard("Party Members", "", container.NewVBox(
			char1Picker,
			char2Picker,
			char3Picker,
			char4Picker,
		)),
		saveBtn,
	)
}

// Show displays the party preset dialog
func (ppd *PartyPresetDialog) Show() {
	ppd.dialog.Show()
}

// SetOnClose sets the close callback
func (ppd *PartyPresetDialog) SetOnClose(callback func()) {
	ppd.onClose = callback
}
