package selections

import (
	"ffvi_editor/io/validation"
	"ffvi_editor/ui/forms"
	"ffvi_editor/ui/forms/editors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type (
	Editor struct {
		widget.BaseWidget
		tabs            *container.AppTabs
		validationPanel *forms.ValidationPanel
		onTabChanged    func(title string)
	}
)

func NewEditor() *Editor {
	s := &Editor{}
	s.ExtendBaseWidget(s)
	return s
}

// SetOnTabChanged registers a callback for tab selection changes
func (s *Editor) SetOnTabChanged(cb func(title string)) {
	s.onTabChanged = cb
}

// GetValidationPanel returns the validation panel for refresh
func (s *Editor) GetValidationPanel() *forms.ValidationPanel {
	return s.validationPanel
}

func (s *Editor) CreateRenderer() fyne.WidgetRenderer {
	validator := validation.NewValidator()
	s.validationPanel = forms.NewValidationPanel(validator)

	s.tabs = container.NewAppTabs(
		container.NewTabItem("Validation", s.validationPanel.BuildPanel()),
		container.NewTabItem("Characters", NewCharacters()),
		container.NewTabItem("Inventory", NewInventory()),
		container.NewTabItem("Skills", editors.NewSkills()),
		container.NewTabItem("Espers", editors.NewEsper()),
		container.NewTabItem("Party", editors.NewParty()),
		container.NewTabItem("Map", editors.NewMapData()),
		container.NewTabItem("Veldt", editors.NewVeldt()),
	)

	// Emit tab change events
	s.tabs.OnSelected = func(tab *container.TabItem) {
		if s.onTabChanged != nil && tab != nil {
			s.onTabChanged(tab.Text)
		}
	}

	return widget.NewSimpleRenderer(s.tabs)
}
