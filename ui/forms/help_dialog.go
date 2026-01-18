package forms

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/docs"
)

// HelpDialog displays help documentation
type HelpDialog struct {
	window     fyne.Window
	helpSystem *docs.HelpSystem
}

// NewHelpDialog creates a new help dialog
func NewHelpDialog(window fyne.Window) *HelpDialog {
	return &HelpDialog{
		window:     window,
		helpSystem: docs.NewHelpSystem(),
	}
}

// Show displays the help dialog
func (h *HelpDialog) Show() {
	// Search bar
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search help topics...")

	searchBtn := widget.NewButton("Search", func() {
		query := searchEntry.Text
		if query == "" {
			return
		}
		results := h.helpSystem.Search(query)
		h.showSearchResults(results)
	})

	// Category selector
	categories := h.helpSystem.GetCategories()
	categorySelect := widget.NewSelect(categories, func(value string) {
		topics := h.helpSystem.GetTopicsByCategory(value)
		// Would update topic list here
		_ = topics
	})
	categorySelect.PlaceHolder = "Filter by category..."

	// Topic list
	allTopics := h.helpSystem.GetAllTopics()
	topicList := widget.NewList(
		func() int { return len(allTopics) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Topic")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < len(allTopics) {
				label := item.(*widget.Label)
				label.SetText(allTopics[id].Title)
			}
		},
	)

	topicList.OnSelected = func(id widget.ListItemID) {
		if id < len(allTopics) {
			h.showTopic(allTopics[id])
		}
		topicList.UnselectAll()
	}

	// Quick links
	quickLinksLabel := widget.NewLabel("Quick Links:")
	quickStartBtn := widget.NewButton("Quick Start Guide", func() {
		if topic, err := h.helpSystem.GetTopic("quickstart"); err == nil {
			h.showTopic(topic)
		}
	})
	interfaceBtn := widget.NewButton("Interface Overview", func() {
		if topic, err := h.helpSystem.GetTopic("interface"); err == nil {
			h.showTopic(topic)
		}
	})
	troubleshootBtn := widget.NewButton("Troubleshooting", func() {
		if topic, err := h.helpSystem.GetTopic("validation_errors"); err == nil {
			h.showTopic(topic)
		}
	})
	shortcutsBtn := widget.NewButton("Keyboard Shortcuts", func() {
		if topic, err := h.helpSystem.GetTopic("keyboard_shortcuts"); err == nil {
			h.showTopic(topic)
		}
	})

	quickLinks := container.NewVBox(
		quickLinksLabel,
		quickStartBtn,
		interfaceBtn,
		troubleshootBtn,
		shortcutsBtn,
	)

	// Build layout
	searchBar := container.NewBorder(
		nil, nil, nil,
		searchBtn,
		searchEntry,
	)

	sidebar := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Help Topics"),
			categorySelect,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			quickLinks,
		),
		nil,
		nil,
		topicList,
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Help & Documentation"),
			widget.NewSeparator(),
			searchBar,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		sidebar,
	)

	d := dialog.NewCustom("Help", "Close", WrapDialogWithDarkBackground(content), h.window)
	d.Resize(fyne.NewSize(700, 500))
	d.Show()
}

// showTopic displays a specific help topic
func (h *HelpDialog) showTopic(topic *docs.HelpTopic) {
	// Content area
	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetText(topic.Content)
	contentEntry.Wrapping = fyne.TextWrapWord
	contentEntry.Disable()

	// Related topics
	var relatedButtons []fyne.CanvasObject
	if len(topic.RelatedTopics) > 0 {
		relatedButtons = append(relatedButtons, widget.NewLabel("Related Topics:"))
		for _, relatedID := range topic.RelatedTopics {
			if relatedTopic, err := h.helpSystem.GetTopic(relatedID); err == nil {
				// Capture in loop
				rt := relatedTopic
				btn := widget.NewButton(rt.Title, func() {
					h.showTopic(rt)
				})
				relatedButtons = append(relatedButtons, btn)
			}
		}
	}

	var content fyne.CanvasObject
	if len(relatedButtons) > 0 {
		content = container.NewBorder(
			nil,
			container.NewVBox(
				widget.NewSeparator(),
				container.NewVBox(relatedButtons...),
			),
			nil,
			nil,
			contentEntry,
		)
	} else {
		content = contentEntry
	}

	d := dialog.NewCustom(topic.Title, "Close", content, h.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}

// showSearchResults displays search results
func (h *HelpDialog) showSearchResults(results []*docs.HelpTopic) {
	if len(results) == 0 {
		dialog.ShowInformation("No Results", "No help topics found matching your search.", h.window)
		return
	}

	resultList := widget.NewList(
		func() int { return len(results) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Topic Title"),
				widget.NewLabel("Category"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < len(results) {
				box := item.(*fyne.Container)
				titleLabel := box.Objects[0].(*widget.Label)
				categoryLabel := box.Objects[1].(*widget.Label)

				titleLabel.SetText(results[id].Title)
				categoryLabel.SetText("Category: " + results[id].Category)
			}
		},
	)

	resultList.OnSelected = func(id widget.ListItemID) {
		if id < len(results) {
			h.showTopic(results[id])
		}
		resultList.UnselectAll()
	}

	content := container.NewBorder(
		widget.NewLabel("Search Results"),
		nil,
		nil,
		nil,
		resultList,
	)

	d := dialog.NewCustom("Search Results", "Close", content, h.window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}

// ShowQuickHelp shows a quick help dialog for a specific feature
func ShowQuickHelp(window fyne.Window, topicID string) {
	helpSystem := docs.NewHelpSystem()
	topic, err := helpSystem.GetTopic(topicID)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetText(topic.Content)
	contentEntry.Wrapping = fyne.TextWrapWord
	contentEntry.Disable()

	d := dialog.NewCustom(topic.Title, "Close", contentEntry, window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
