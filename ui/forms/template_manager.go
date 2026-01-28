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
	dialog   dialog.Dialog
	manager  *templates.Manager
	template *tmpl.CharacterTemplate
}

// NewTemplateManagerDialog creates a new template manager dialog
func NewTemplateManagerDialog(parent fyne.Window, templateManager *templates.Manager) *TemplateManagerDialog {
	tmd := &TemplateManagerDialog{
		manager: templateManager,
	}

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
				tmd.createDetailsPanel(),
			),
		),
	)

	// Create buttons
	applyBtn := widget.NewButton("Apply Template", func() {
		if tmd.template != nil {
			// TODO: Show apply dialog with mode selection
			dialog.ShowInformation("Apply Template", fmt.Sprintf("Applying template: %s", tmd.template.Name), parent)
		}
	})

	editBtn := widget.NewButton("Edit Template", func() {
		if tmd.template != nil {
			dialog.ShowInformation("Edit Template", fmt.Sprintf("Editing template: %s", tmd.template.Name), parent)
		}
	})

	deleteBtn := widget.NewButton("Delete Template", func() {
		if tmd.template != nil {
			dialog.ShowConfirm(
				"Delete Template",
				fmt.Sprintf("Delete template '%s'?", tmd.template.Name),
				func(b bool) {
					if b {
						// TODO: Implement template deletion
						dialog.ShowInformation("Deleted", fmt.Sprintf("Template '%s' deleted", tmd.template.Name), parent)
					}
				},
				parent,
			)
		}
	})

	exportBtn := widget.NewButton("Export", func() {
		if tmd.template != nil {
			dialog.ShowInformation("Export Template", fmt.Sprintf("Exporting template: %s", tmd.template.Name), parent)
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
	list := widget.NewList(
		func() int {
			// TODO: Get templates from manager
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Name")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// TODO: Update template list items
		},
	)

	return list
}

// createDetailsPanel creates the template details panel
func (tmd *TemplateManagerDialog) createDetailsPanel() fyne.CanvasObject {
	if tmd.template == nil {
		return widget.NewLabel("Select a template to view details")
	}

	return container.NewVBox(
		widget.NewCard("Name", "", widget.NewLabel(tmd.template.Name)),
		widget.NewCard("Description", "", widget.NewLabel(tmd.template.Description)),
		widget.NewCard("Tags", "", widget.NewLabel(fmt.Sprintf("%v", tmd.template.Tags))),
		widget.NewCard("Created", "", widget.NewLabel(tmd.template.CreatedAt.String())),
		widget.NewCard("Favorite", "", widget.NewCheck("Favorite", nil)),
	)
}

// Show displays the template manager dialog
func (tmd *TemplateManagerDialog) Show() {
	tmd.dialog.Show()
}
