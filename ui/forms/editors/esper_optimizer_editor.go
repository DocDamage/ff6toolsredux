package editors

import (
	"ffvi_editor/io/pr"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// EsperOptimizerDialog helps optimize esper sequences
type EsperOptimizerDialog struct {
	window     fyne.Window
	save       *pr.PR
	targetStat *widget.Select
	charSel    *widget.Select
	esperList  *widget.List
	detailText *widget.RichText
	applyBtn   *widget.Button
	selected   int
}

// NewEsperOptimizerDialog creates a new esper optimizer dialog
func NewEsperOptimizerDialog(window fyne.Window, save *pr.PR) fyne.CanvasObject {
	editor := &EsperOptimizerDialog{
		window:   window,
		save:     save,
		selected: -1,
	}

	return editor.buildUI()
}

// buildUI builds the complete UI
func (e *EsperOptimizerDialog) buildUI() fyne.CanvasObject {
	chars := []string{
		"Terra", "Locke", "Cyan", "Shadow", "Edgar", "Sabin",
		"Celes", "Strago", "Relm", "Setzer", "Mog", "Gau",
		"Gogo", "Umaro", "Kefka", "Esper",
	}
	e.charSel = widget.NewSelect(chars, func(s string) {
	})
	e.charSel.SetSelected("Terra")

	stats := []string{"Balanced", "Vigor", "Speed", "Stamina", "Magic Power", "Defense", "Magic Defense"}
	e.targetStat = widget.NewSelect(stats, func(s string) {
	})
	e.targetStat.SetSelected("Balanced")

	e.esperList = widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	e.detailText = widget.NewRichTextFromMarkdown("")
	e.applyBtn = widget.NewButton("Use This Sequence", func() {})
	e.applyBtn.Disable()

	infoLabel := widget.NewLabel("Esper Optimizer")

	return container.NewBorder(
		infoLabel,
		nil,
		nil,
		nil,
		container.NewHBox(
			container.NewVBox(
				widget.NewLabel("Character:"),
				e.charSel,
				widget.NewLabel("Target:"),
				e.targetStat,
				e.esperList,
			),
			container.NewVBox(
				e.detailText,
				e.applyBtn,
			),
		),
	)
}
