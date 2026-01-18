package forms

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// DialogConfig holds configuration for consistent dialog styling
type DialogConfig struct {
	Title      string
	Width      float32
	Height     float32
	MinWidth   float32
	MinHeight  float32
	ShowScroll bool
}

// DefaultDialogConfig returns standard dialog sizing
var DefaultDialogConfig = DialogConfig{
	Width:      600,
	Height:     500,
	MinWidth:   400,
	MinHeight:  300,
	ShowScroll: true,
}

// LargeDialogConfig returns sizing for content-heavy dialogs
var LargeDialogConfig = DialogConfig{
	Width:      800,
	Height:     700,
	MinWidth:   600,
	MinHeight:  500,
	ShowScroll: true,
}

// CompactDialogConfig returns sizing for simple dialogs
var CompactDialogConfig = DialogConfig{
	Width:      400,
	Height:     300,
	MinWidth:   300,
	MinHeight:  200,
	ShowScroll: false,
}

// ValidationFeedback provides real-time input validation
type ValidationFeedback struct {
	Entry    *widget.Entry
	Feedback *widget.Label
}

// NewValidationFeedback creates a text entry with validation feedback
func NewValidationFeedback(placeholder string, validator func(string) (bool, string)) *ValidationFeedback {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)

	feedback := widget.NewLabel("")

	// Setup validation on change
	entry.OnChanged = func(s string) {
		valid, msg := validator(s)
		if valid {
			feedback.SetText("✓ " + msg)
		} else {
			feedback.SetText("✗ " + msg)
		}
	}

	return &ValidationFeedback{
		Entry:    entry,
		Feedback: feedback,
	}
}

// GetText returns the entry text
func (vf *ValidationFeedback) GetText() string {
	return vf.Entry.Text
}

// SetText sets the entry text
func (vf *ValidationFeedback) SetText(text string) {
	vf.Entry.SetText(text)
}

// FormHelper assists in building consistent forms
type FormHelper struct {
	Items []*widget.FormItem
}

// NewFormHelper creates a new form helper
func NewFormHelper() *FormHelper {
	return &FormHelper{}
}

// AddEntry adds a simple text entry
func (fh *FormHelper) AddEntry(label, placeholder, defaultValue string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entry.SetText(defaultValue)
	fh.Items = append(fh.Items, widget.NewFormItem(label, entry))
	return entry
}

// AddSelect adds a dropdown select
func (fh *FormHelper) AddSelect(label string, options []string, defaultOption string) *widget.Select {
	selectWidget := widget.NewSelect(options, func(s string) {})
	selectWidget.SetSelected(defaultOption)
	fh.Items = append(fh.Items, widget.NewFormItem(label, selectWidget))
	return selectWidget
}

// AddCheck adds a checkbox
func (fh *FormHelper) AddCheck(label string, defaultValue bool) *widget.Check {
	check := widget.NewCheck(label, func(b bool) {})
	check.Checked = defaultValue
	fh.Items = append(fh.Items, widget.NewFormItem("", check))
	return check
}

// Build creates the form
func (fh *FormHelper) Build() *widget.Form {
	return widget.NewForm(fh.Items...)
}

// SuccessMessage displays a styled success message
func SuccessMessage(title string, message string) fyne.CanvasObject {
	titleLabel := widget.NewLabel("✓ " + title)
	titleLabel.TextStyle.Bold = true

	msgLabel := widget.NewLabel(message)
	msgLabel.Alignment = fyne.TextAlignLeading

	return CreateInfoBox(title, message)
}

// ErrorMessage displays a styled error message
func ErrorMessage(title string, message string) fyne.CanvasObject {
	titleLabel := widget.NewLabel("✗ " + title)
	titleLabel.TextStyle.Bold = true

	msgLabel := widget.NewLabel(message)
	msgLabel.Alignment = fyne.TextAlignLeading

	return CreateInfoBox(title, message)
}

// WarningMessage displays a styled warning message
func WarningMessage(title string, message string) fyne.CanvasObject {
	titleLabel := widget.NewLabel("⚠ " + title)
	titleLabel.TextStyle.Bold = true

	msgLabel := widget.NewLabel(message)
	msgLabel.Alignment = fyne.TextAlignLeading

	return CreateInfoBox(title, message)
}
