package editors

import (
	"ffvi_editor/io/pr"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SketchEditorDialog allows users to browse and manage sketches
type SketchEditorDialog struct {
	window     fyne.Window
	save       *pr.PR
	searchText *widget.Entry
	diffFilter *widget.Select
	sketchList *widget.List
	detailText *widget.RichText
	learnBtn   *widget.Button
	selected   int
}

// NewSketchEditorDialog creates a new sketch editor dialog
func NewSketchEditorDialog(window fyne.Window, save *pr.PR) fyne.CanvasObject {
	editor := &SketchEditorDialog{
		window:   window,
		save:     save,
		selected: -1,
	}

	return editor.buildUI()
}

// buildUI builds the complete UI
func (s *SketchEditorDialog) buildUI() fyne.CanvasObject {
	s.searchText = widget.NewEntry()
	s.searchText.SetPlaceHolder("Search sketches...")

	s.diffFilter = widget.NewSelect(
		[]string{"All", "Easy", "Normal", "Hard", "Very Hard"},
		func(str string) {},
	)
	s.diffFilter.SetSelected("All")

	s.sketchList = widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	s.detailText = widget.NewRichTextFromMarkdown("")
	s.learnBtn = widget.NewButton("Mark as Sketched", func() {})
	s.learnBtn.Disable()

	progressLabel := widget.NewLabel("Sketch Database")

	return container.NewBorder(
		progressLabel,
		nil,
		nil,
		nil,
		container.NewHBox(
			container.NewVBox(
				s.searchText,
				s.diffFilter,
				s.sketchList,
			),
			container.NewVBox(
				s.detailText,
				s.learnBtn,
			),
		),
	)
}
