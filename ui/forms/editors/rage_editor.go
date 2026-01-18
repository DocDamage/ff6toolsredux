package editors

import (
	"fmt"

	"ffvi_editor/io/pr"
	"ffvi_editor/models/game"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RageEditorDialog allows users to browse and learn Gau rages
type RageEditorDialog struct {
	window     fyne.Window
	save       *pr.PR
	searchText *widget.Entry
	filterDiff *widget.Select
	rageList   *widget.List
	detailText *widget.RichText
	learnBtn   *widget.Button
	rages      []*game.RageEntry
	selected   int
}

// NewRageEditorDialog creates a new rage editor dialog
func NewRageEditorDialog(window fyne.Window, save *pr.PR) fyne.CanvasObject {
	rages := make([]*game.RageEntry, 0)
	for _, rage := range game.RageDatabase {
		rages = append(rages, rage)
	}

	editor := &RageEditorDialog{
		window:   window,
		save:     save,
		rages:    rages,
		selected: -1,
	}

	return editor.buildUI()
}

// buildUI builds the complete UI
func (r *RageEditorDialog) buildUI() fyne.CanvasObject {
	r.searchText = widget.NewEntry()
	r.searchText.SetPlaceHolder("Search rages...")
	r.searchText.OnChanged = func(s string) {
		r.updateRageList()
	}

	r.filterDiff = widget.NewSelect(
		[]string{"All Difficulties", "Easy", "Normal", "Hard", "Very Hard"},
		func(s string) {
			r.updateRageList()
		},
	)
	r.filterDiff.SetSelected("All Difficulties")

	r.rageList = widget.NewList(
		func() int { return len(r.rages) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Rage Name")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(r.rages) {
				rage := r.rages[id]
				label.SetText(fmt.Sprintf("%s (%d) - %s", rage.Name, rage.ID, rage.Enemy))
			}
		},
	)
	r.rageList.OnSelected = func(id widget.ListItemID) {
		r.selected = id
		if id < len(r.rages) {
			r.updateDetailView(r.rages[id])
		}
	}

	r.detailText = widget.NewRichTextFromMarkdown("")

	r.learnBtn = widget.NewButton("Learn Rage", func() {
		if r.selected >= 0 && r.selected < len(r.rages) {
			// Learn rage functionality would be implemented here
			r.learnBtn.Disable()
		}
	})
	r.learnBtn.Disable()

	progressLabel := widget.NewLabel("Rage Database Browser")

	searchContainer := container.NewVBox(
		widget.NewLabel("Search and Filter:"),
		r.searchText,
		r.filterDiff,
	)

	contentContainer := container.NewHBox(
		container.NewVBox(
			searchContainer,
			container.NewBorder(nil, nil, nil, nil, r.rageList),
		),
		container.NewVBox(
			widget.NewLabel("Rage Details:"),
			r.detailText,
			r.learnBtn,
		),
	)

	return container.NewBorder(
		progressLabel,
		nil,
		nil,
		nil,
		contentContainer,
	)
}

// updateRageList updates the visible rages based on filters
func (r *RageEditorDialog) updateRageList() {
	searchTerm := r.searchText.Text
	difficulty := r.filterDiff.Selected

	filtered := make([]*game.RageEntry, 0)
	for _, rage := range game.RageDatabase {
		// Filter by search term
		if searchTerm != "" && !stringContains(rage.Name, searchTerm) &&
			!stringContains(rage.Enemy, searchTerm) {
			continue
		}

		// Filter by difficulty
		if difficulty != "All Difficulties" && difficultyToString(rage.Difficulty) != difficulty {
			continue
		}

		filtered = append(filtered, rage)
	}

	r.rages = filtered
	r.rageList.Refresh()
}

// difficultyToString converts difficulty int to string
func difficultyToString(d int) string {
	switch d {
	case 1:
		return "Easy"
	case 2:
		return "Normal"
	case 3:
		return "Hard"
	case 4:
		return "Very Hard"
	default:
		return "Normal"
	}
}

// updateDetailView updates the detail text
func (r *RageEditorDialog) updateDetailView(rage *game.RageEntry) {
	detail := fmt.Sprintf(`
## %s (ID: %d)

**Enemy:** %s
**Location:** %s
**Level:** %d
**Difficulty:** %s

**Stats:**
- HP: %d
- Speed: %d
- Attack: %s

**Element:** %s
**Status Effects:** %v

**Notes:** %s
`,
		rage.Name, rage.ID,
		rage.Enemy, rage.Location, rage.Level, difficultyToString(rage.Difficulty),
		rage.HP, rage.Speed, rage.Attack,
		rage.ElementType, rage.Status,
		rage.Notes,
	)

	r.detailText.ParseMarkdown(detail)
	r.learnBtn.Enable()
}

// stringContains checks if str contains substr (case-insensitive)
func stringContains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(str[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// toLower converts a byte to lowercase
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}
