package forms

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/models/search"
)

// SearchResult wraps the search result type
type SearchResult struct {
	ID      string
	Name    string
	Type    string // "character", "item", "spell", "esper"
	Details string
	Score   float64
}

// SearchDialog provides global search functionality
type SearchDialog struct {
	dialog          dialog.Dialog
	searchIndex     *search.Index
	resultsCallback func(*SearchResult)
}

// NewSearchDialog creates a new global search dialog
func NewSearchDialog(parent fyne.Window, index *search.Index, resultsCallback func(*SearchResult)) *SearchDialog {
	sd := &SearchDialog{
		searchIndex:     index,
		resultsCallback: resultsCallback,
	}

	// Search input
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search characters, items, spells, espers...")

	// Filter options
	filterGroup := container.NewVBox(
		widget.NewLabel("Filter by type:"),
	)

	filterCharacters := widget.NewCheck("Characters", func(b bool) {})
	filterItems := widget.NewCheck("Items", func(b bool) {})
	filterSpells := widget.NewCheck("Spells", func(b bool) {})
	filterEspers := widget.NewCheck("Espers", func(b bool) {})

	filterCharacters.Checked = true
	filterItems.Checked = true
	filterSpells.Checked = true
	filterEspers.Checked = true

	filterGroup.Add(filterCharacters)
	filterGroup.Add(filterItems)
	filterGroup.Add(filterSpells)
	filterGroup.Add(filterEspers)

	// Results list
	resultsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Result"),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	// Search results display
	resultsBox := container.NewVBox(
		widget.NewLabel("Results:"),
		resultsList,
	)

	// Search function
	performSearch := func(query string) {
		if query == "" {
			resultsList.Length = func() int { return 0 }
			resultsList.Refresh()
			return
		}

		// Get search results from index (using search results type)
		results := sd.searchIndex.Search(query)
		if results == nil {
			results = make([]search.SearchResult, 0)
		}

		// Filter results based on checkboxes
		filtered := make([]search.SearchResult, 0)
		for _, result := range results {
			switch result.Type {
			case search.ResultTypeCharacter:
				if filterCharacters.Checked {
					filtered = append(filtered, result)
				}
			case search.ResultTypeItem:
				if filterItems.Checked {
					filtered = append(filtered, result)
				}
			case search.ResultTypeSpell:
				if filterSpells.Checked {
					filtered = append(filtered, result)
				}
			case search.ResultTypeEsper:
				if filterEspers.Checked {
					filtered = append(filtered, result)
				}
			}
		}

		// Update results list
		resultsList.Length = func() int { return len(filtered) }
		resultsList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < len(filtered) {
				result := filtered[id]
				box := item.(*fyne.Container)
				if len(box.Objects) >= 2 {
					box.Objects[0].(*widget.Label).SetText(result.Name)
					box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("%s • %s", result.Type, result.Description))
				}
			}
		}

		resultsList.OnSelected = func(id widget.ListItemID) {
			if id < len(filtered) && sd.resultsCallback != nil {
				// Convert to SearchResult wrapper
				result := filtered[id]
				sr := &SearchResult{
					ID:      result.ID,
					Name:    result.Name,
					Type:    string(result.Type),
					Details: result.Description,
				}
				sd.resultsCallback(sr)
			}
		}

		resultsList.Refresh()
	}

	// Debounced search
	searchEntry.OnChanged = func(s string) {
		performSearch(s)
	}

	// Buttons
	closeBtn := widget.NewButton("Close", func() {
		sd.dialog.Hide()
	})

	buttons := container.NewHBox(closeBtn)

	// Main layout
	content := container.NewVBox(
		widget.NewCard("Search", "", container.NewVBox(
			searchEntry,
			filterGroup,
		)),
		widget.NewCard("Results", "", container.NewScroll(resultsBox)),
	)

	sd.dialog = dialog.NewCustom(
		"Search",
		"",
		container.NewBorder(nil, buttons, nil, nil, content),
		parent,
	)

	// Focus on search entry
	canvas := parent.Canvas()
	if canvas != nil {
		canvas.Focus(searchEntry)
	}

	return sd
}

// Show displays the search dialog
func (sd *SearchDialog) Show() {
	sd.dialog.Show()
}
