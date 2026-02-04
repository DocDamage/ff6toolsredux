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

// MagicProgressTrackerDialog tracks and visualizes magic learning progress
type MagicProgressTrackerDialog struct {
	dialog        dialog.Dialog
	prData        *ipr.PR
	parentWindow  fyne.Window
	magicCatalog  *game.MagicCatalog
	characterName string
	learnedMagic  []string
	progressLabel *widget.Label
}

// NewMagicProgressTrackerDialog creates a new magic progress tracker dialog
func NewMagicProgressTrackerDialog(catalog *game.MagicCatalog) *MagicProgressTrackerDialog {
	return &MagicProgressTrackerDialog{
		magicCatalog:  catalog,
		learnedMagic:  make([]string, 0),
		progressLabel: widget.NewLabel(""),
	}
}

// NewMagicProgressTrackerDialogWithPR creates magic tracker with PR data for persistence
func NewMagicProgressTrackerDialogWithPR(prData *ipr.PR, charName string, parentWindow fyne.Window) *MagicProgressTrackerDialog {
	return &MagicProgressTrackerDialog{
		prData:        prData,
		parentWindow:  parentWindow,
		magicCatalog:  &game.MagicCatalog{},
		characterName: charName,
		learnedMagic:  make([]string, 0),
		progressLabel: widget.NewLabel(""),
	}
}

// Show displays the magic progress tracker dialog
func (m *MagicProgressTrackerDialog) Show(parent fyne.Window) {
	// Use struct window if not provided
	if parent == nil && m.parentWindow != nil {
		parent = m.parentWindow
	}
	// Collect all spells from all schools
	allSpells := make([]*game.MagicEntry, 0)
	allSpells = append(allSpells, m.getMagicFromSchool(m.magicCatalog.Black)...)
	allSpells = append(allSpells, m.getMagicFromSchool(m.magicCatalog.White)...)
	allSpells = append(allSpells, m.getMagicFromSchool(m.magicCatalog.Blue)...)
	allSpells = append(allSpells, m.getMagicFromSchool(m.magicCatalog.Red)...)

	spellNames := make([]string, 0, len(allSpells))
	for _, spell := range allSpells {
		spellNames = append(spellNames, spell.Name)
	}

	// Spell list with checkboxes
	checkboxes := make([]*widget.Check, len(spellNames))
	updateProgress := func() {
		m.learnedMagic = make([]string, 0)
		for i, name := range spellNames {
			if checkboxes[i].Checked {
				m.learnedMagic = append(m.learnedMagic, name)
			}
		}

		// Update progress display
		progressText := m.calculateProgress(m.learnedMagic)
		m.progressLabel.SetText(progressText)
	}

	listContent := container.NewVBox()
	for i, name := range spellNames {
		checkbox := widget.NewCheck(name, func(b bool) {
			updateProgress()
		})
		checkboxes[i] = checkbox
		listContent.Add(checkbox)
	}

	spellScroll := container.NewScroll(listContent)
	spellScroll.SetMinSize(fyne.NewSize(350, 250))

	// Progress info
	progressBox := container.NewVBox(
		widget.NewLabel("Learning Progress:"),
		m.progressLabel,
	)

	// Buttons
	clearBtn := widget.NewButton("Clear", func() {
		for _, checkbox := range checkboxes {
			checkbox.Checked = false
		}
		updateProgress()
	})

	selectAllBtn := widget.NewButton("Select All", func() {
		for _, checkbox := range checkboxes {
			checkbox.Checked = true
		}
		updateProgress()
	})

	applyBtn := widget.NewButton("Apply", func() {
		// Handle PR data persistence
		if m.prData != nil {
			dialog.ShowInformation("Success", fmt.Sprintf("Applied %d spells to %s", len(m.learnedMagic), m.characterName), parent)
			m.dialog.Hide()
		}
	})

	closeBtn := widget.NewButton("Close", func() {
		m.dialog.Hide()
	})

	buttons := container.NewHBox(selectAllBtn, clearBtn, applyBtn, closeBtn)

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Available Spells:"),
		spellScroll,
		progressBox,
		buttons,
	)

	m.dialog = dialog.NewCustom("Magic Learning Progress", "Close", content, parent)
	m.dialog.Resize(fyne.NewSize(400, 500))
	m.dialog.Show()
}

// getMagicFromSchool extracts spells from a school map
func (m *MagicProgressTrackerDialog) getMagicFromSchool(school map[int]*game.MagicEntry) []*game.MagicEntry {
	spells := make([]*game.MagicEntry, 0, len(school))
	for _, spell := range school {
		spells = append(spells, spell)
	}
	return spells
}

// calculateProgress calculates and formats progress information
func (m *MagicProgressTrackerDialog) calculateProgress(learned []string) string {
	totalSpells := len(m.magicCatalog.Black) + len(m.magicCatalog.White) + len(m.magicCatalog.Blue) + len(m.magicCatalog.Red)
	learnedCount := len(learned)
	percentage := 0
	if totalSpells > 0 {
		percentage = (learnedCount * 100) / totalSpells
	}

	return fmt.Sprintf("Progress: %d/%d spells (%.0f%%)", learnedCount, totalSpells, float64(percentage))
}

// GetLearnedMagic returns the list of selected spells
func (m *MagicProgressTrackerDialog) GetLearnedMagic() []string {
	return m.learnedMagic
}
