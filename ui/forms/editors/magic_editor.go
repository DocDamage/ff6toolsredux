package editors

import (
	"fmt"

	"ffvi_editor/io/pr"
	"ffvi_editor/models/game"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MagicEditorDialog allows users to browse and learn magic spells
type MagicEditorDialog struct {
	window     fyne.Window
	save       *pr.PR
	searchText *widget.Entry
	schoolSel  *widget.Select
	magicList  *widget.List
	detailText *widget.RichText
	learnBtn   *widget.Button
	magics     []*game.MagicEntry
	selected   int
}

// NewMagicEditorDialog creates a new magic editor dialog
func NewMagicEditorDialog(window fyne.Window, save *pr.PR) fyne.CanvasObject {
	magics := make([]*game.MagicEntry, 0)

	editor := &MagicEditorDialog{
		window:   window,
		save:     save,
		magics:   magics,
		selected: -1,
	}

	return editor.buildUI()
}

// buildUI builds the complete UI
func (m *MagicEditorDialog) buildUI() fyne.CanvasObject {
	m.searchText = widget.NewEntry()
	m.searchText.SetPlaceHolder("Search spells...")
	m.searchText.OnChanged = func(s string) {
		m.updateMagicList()
	}

	schools := []string{"All Schools", "Black", "White", "Blue", "Red", "Special"}
	m.schoolSel = widget.NewSelect(schools, func(s string) {
		m.updateMagicList()
	})
	m.schoolSel.SetSelected("All Schools")

	m.magicList = widget.NewList(
		func() int { return len(m.magics) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(m.magics) {
				spell := m.magics[id]
				label.SetText(fmt.Sprintf("%s - Cost: %d MP", spell.Name, spell.Cost))
			}
		},
	)
	m.magicList.OnSelected = func(id widget.ListItemID) {
		m.selected = id
		if id < len(m.magics) {
			m.updateDetailView(m.magics[id])
		}
	}

	m.detailText = widget.NewRichTextFromMarkdown("")
	m.learnBtn = widget.NewButton("Learn Spell", func() {})
	m.learnBtn.Disable()

	progressLabel := widget.NewLabel("Magic Browser")

	return container.NewBorder(
		progressLabel,
		nil,
		nil,
		nil,
		container.NewHBox(
			container.NewVBox(
				m.searchText,
				m.schoolSel,
				m.magicList,
			),
			container.NewVBox(
				m.detailText,
				m.learnBtn,
			),
		),
	)
}

func (m *MagicEditorDialog) updateMagicList() {
}

func (m *MagicEditorDialog) updateDetailView(spell *game.MagicEntry) {
	text := fmt.Sprintf("## %s\nCost: %d MP\nSchool: %s", spell.Name, spell.Cost, spell.School)
	m.detailText.ParseMarkdown(text)
}

func (m *MagicEditorDialog) updateProgressLabel(label *widget.Label) {
	label.SetText("Magic Browser")
}
