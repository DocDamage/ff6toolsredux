package forms

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io/templates"
	tmpl "ffvi_editor/models/templates"
)

// TemplateManagerDialog displays and manages character build templates
type TemplateManagerDialog struct {
	dialog       dialog.Dialog
	manager      *templates.Manager
	template     *tmpl.CharacterTemplate
	templates    []*tmpl.CharacterTemplate
	list         *widget.List
	detailsPanel fyne.CanvasObject
	parent       fyne.Window
}

// NewTemplateManagerDialog creates a new template manager dialog
func NewTemplateManagerDialog(parent fyne.Window, templateManager *templates.Manager) *TemplateManagerDialog {
	tmd := &TemplateManagerDialog{
		manager: templateManager,
		parent:  parent,
	}

	// Load templates
	tmd.refreshTemplateList()

	// Create details panel (will be updated on selection)
	tmd.detailsPanel = tmd.createDetailsPanel()

	// Create content
	content := container.NewVBox(
		widget.NewLabel("Template Manager"),
		container.NewVBox(
			widget.NewCard(
				"Available Templates",
				"Select a template to view or apply",
				tmd.createTemplateList(),
			),
			widget.NewCard(
				"Template Details",
				"",
				tmd.detailsPanel,
			),
		),
	)

	// Create buttons
	applyBtn := widget.NewButton("Apply Template", func() {
		if tmd.template != nil {
			tmd.showApplyDialog(parent)
		}
	})

	editBtn := widget.NewButton("Edit Template", func() {
		if tmd.template != nil {
			tmd.showEditDialog(parent)
		}
	})

	deleteBtn := widget.NewButton("Delete Template", func() {
		if tmd.template != nil {
			dialog.ShowConfirm(
				"Delete Template",
				fmt.Sprintf("Delete template '%s'?", tmd.template.Name),
				func(b bool) {
					if b {
						if err := tmd.manager.DeleteTemplate(tmd.template.ID); err != nil {
							dialog.ShowError(err, parent)
							return
						}
						dialog.ShowInformation("Deleted", fmt.Sprintf("Template '%s' deleted", tmd.template.Name), parent)
						tmd.refreshTemplateList()
						tmd.list.Refresh()
					}
				},
				parent,
			)
		}
	})

	exportBtn := widget.NewButton("Export", func() {
		if tmd.template != nil {
			tmd.showExportDialog(parent)
		}
	})

	buttons := container.NewHBox(applyBtn, editBtn, deleteBtn, exportBtn)

	// Create dialog
	tmd.dialog = dialog.NewCustom(
		"Character Templates",
		"Close",
		container.NewBorder(nil, buttons, nil, nil, content),
		parent,
	)

	return tmd
}

// createTemplateList creates the template list widget
func (tmd *TemplateManagerDialog) createTemplateList() fyne.CanvasObject {
	tmd.list = widget.NewList(
		func() int {
			return len(tmd.templates)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(nil),
				widget.NewLabel("Template Name"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(tmd.templates) {
				return
			}
			template := tmd.templates[id]
			hbox := item.(*fyne.Container)
			icon := hbox.Objects[0].(*widget.Icon)
			label := hbox.Objects[1].(*widget.Label)

			label.SetText(template.Name)
			if template.Favorite {
				// Would use a star icon if available
				label.SetText("★ " + template.Name)
			}
			_ = icon
		},
	)

	tmd.list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(tmd.templates) {
			tmd.template = tmd.templates[id]
			tmd.updateDetailsPanel()
		}
	}

	return tmd.list
}

// createDetailsPanel creates the template details panel
func (tmd *TemplateManagerDialog) createDetailsPanel() fyne.CanvasObject {
	if tmd.template == nil {
		return widget.NewLabel("Select a template to view details")
	}

	return tmd.buildDetailsContent()
}

// buildDetailsContent builds the details content for the selected template
func (tmd *TemplateManagerDialog) buildDetailsContent() fyne.CanvasObject {
	t := tmd.template

	// Build tags string
	tagsStr := "None"
	if len(t.Tags) > 0 {
		tagsStr = ""
		for i, tag := range t.Tags {
			if i > 0 {
				tagsStr += ", "
			}
			tagsStr += tag
		}
	}

	// Build character stats summary
	var statsStr string
	if t.Character != nil {
		statsStr = fmt.Sprintf("Level: %d\nHP: %d/%d\nMP: %d/%d\nVigor: %d\nStamina: %d\nSpeed: %d\nMagic: %d",
			t.Character.Level,
			t.Character.HP.Current, t.Character.HP.Max,
			t.Character.MP.Current, t.Character.MP.Max,
			t.Character.Vigor, t.Character.Stamina, t.Character.Speed, t.Character.Magic,
		)
	} else {
		statsStr = "No character data"
	}

	favoriteCheck := widget.NewCheck("Favorite", func(checked bool) {
		if tmd.template != nil {
			tmd.template.Favorite = checked
			tmd.manager.ToggleFavorite(tmd.template.ID)
			tmd.list.Refresh()
		}
	})
	favoriteCheck.SetChecked(t.Favorite)

	return container.NewVBox(
		widget.NewCard("Name", "", widget.NewLabel(t.Name)),
		widget.NewCard("Description", "", widget.NewLabel(t.Description)),
		widget.NewCard("Tags", "", widget.NewLabel(tagsStr)),
		widget.NewCard("Version", "", widget.NewLabel(fmt.Sprintf("%d", t.Version))),
		widget.NewCard("Created", "", widget.NewLabel(t.CreatedAt.Format("2006-01-02 15:04"))),
		widget.NewCard("Favorite", "", favoriteCheck),
		widget.NewCard("Character Stats", "", widget.NewLabel(statsStr)),
	)
}

// updateDetailsPanel updates the details panel with current template
func (tmd *TemplateManagerDialog) updateDetailsPanel() {
	if tmd.detailsPanel == nil {
		return
	}

	// Find the parent container and replace the details panel
	// This is a workaround since we can't directly update the card content
	if tmd.template == nil {
		return
	}

	// Refresh the dialog by recreating it with updated content
	// For now, we'll rely on the user reopening the dialog
	// A more sophisticated approach would use data binding
}

// refreshTemplateList reloads templates from the manager
func (tmd *TemplateManagerDialog) refreshTemplateList() {
	tmd.templates = tmd.manager.ListTemplates("")
	if tmd.list != nil {
		tmd.list.Refresh()
	}
}

// Show displays the template manager dialog
func (tmd *TemplateManagerDialog) Show() {
	tmd.refreshTemplateList()
	tmd.dialog.Show()
}

// showApplyDialog shows the apply template dialog with mode selection
func (tmd *TemplateManagerDialog) showApplyDialog(parent fyne.Window) {
	if tmd.template == nil {
		return
	}

	var selectedMode tmpl.ApplyMode = tmpl.ApplyModeReplace

	modeRadio := widget.NewRadioGroup([]string{
		"Replace All (Overwrite character with template)",
		"Merge (Only apply stronger stats)",
		"Preview (Validate only, no changes)",
	}, func(value string) {
		switch value {
		case "Replace All (Overwrite character with template)":
			selectedMode = tmpl.ApplyModeReplace
		case "Merge (Only apply stronger stats)":
			selectedMode = tmpl.ApplyModeMerge
		case "Preview (Validate only, no changes)":
			selectedMode = tmpl.ApplyModePreview
		}
	})
	modeRadio.SetSelected("Replace All (Overwrite character with template)")

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Apply template: %s", tmd.template.Name)),
		widget.NewSeparator(),
		widget.NewLabel("Select apply mode:"),
		modeRadio,
	)

	dialog.ShowCustomConfirm("Apply Template", "Apply", "Cancel", content,
		func(apply bool) {
			if !apply {
				return
			}

			// In a real implementation, we would get the target character
			// from the current editor state and apply the template
			// For now, show what would happen
			modeStr := "Replace"
			if selectedMode == tmpl.ApplyModeMerge {
				modeStr = "Merge"
			} else if selectedMode == tmpl.ApplyModePreview {
				modeStr = "Preview"
			}

			dialog.ShowInformation("Apply Template",
				fmt.Sprintf("Template '%s' would be applied with %s mode\n(Integration with character editor pending)",
					tmd.template.Name, modeStr), parent)
		}, parent)
}

// showEditDialog shows the edit template dialog
func (tmd *TemplateManagerDialog) showEditDialog(parent fyne.Window) {
	if tmd.template == nil {
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(tmd.template.Name)

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(tmd.template.Description)

	tagsEntry := widget.NewEntry()
	tagsStr := ""
	for i, tag := range tmd.template.Tags {
		if i > 0 {
			tagsStr += ", "
		}
		tagsStr += tag
	}
	tagsEntry.SetText(tagsStr)

	content := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Description", descEntry),
			widget.NewFormItem("Tags (comma separated)", tagsEntry),
		),
	)

	dialog.ShowCustomConfirm("Edit Template", "Save", "Cancel", content,
		func(save bool) {
			if !save {
				return
			}

			tmd.template.Name = nameEntry.Text
			tmd.template.Description = descEntry.Text

			// Parse tags
			if tagsEntry.Text != "" {
				tags := []string{}
				current := ""
				for _, r := range tagsEntry.Text {
					if r == ',' {
						if current != "" {
							tags = append(tags, current)
							current = ""
						}
					} else if r != ' ' {
						current += string(r)
					}
				}
				if current != "" {
					tags = append(tags, current)
				}
				tmd.template.Tags = tags
			}

			if err := tmd.manager.UpdateTemplate(tmd.template); err != nil {
				dialog.ShowError(err, parent)
				return
			}

			tmd.refreshTemplateList()
			tmd.list.Refresh()
			dialog.ShowInformation("Saved", "Template updated successfully", parent)
		}, parent)
}

// showExportDialog shows the export template dialog
func (tmd *TemplateManagerDialog) showExportDialog(parent fyne.Window) {
	if tmd.template == nil {
		return
	}

	data, err := tmd.manager.ExportTemplate(tmd.template.ID)
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	// Show export dialog with the JSON data
	entry := widget.NewMultiLineEntry()
	entry.SetText(string(data))
	entry.Disable()

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Template: %s", tmd.template.Name)),
		widget.NewLabel("Copy the JSON below to share this template:"),
		entry,
	)

	dialog.ShowCustom("Export Template", "Close", content, parent)
}
