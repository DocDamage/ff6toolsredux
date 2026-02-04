package forms

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io/pr"
	"ffvi_editor/io/validation"
	"ffvi_editor/models"
)

// ValidationPanel displays validation results and allows fixing issues
type ValidationPanel struct {
	validator        *validation.Validator
	currentData      *pr.PR
	resultsContainer fyne.CanvasObject
	summaryLabel     *widget.Label
	modeSelect       *widget.Select
	autoFixBtn       *widget.Button
	issuesList       *widget.List
	lastResult       models.ValidationResult
	onIssueFixed     func() // Callback when issue is fixed
}

// NewValidationPanel creates a new validation panel
func NewValidationPanel(validator *validation.Validator) *ValidationPanel {
	vp := &ValidationPanel{
		validator:    validator,
		summaryLabel: widget.NewLabel("Ready to validate"),
		issuesList: widget.NewList(
			func() int { return 0 },
			func() fyne.CanvasObject {
				return widget.NewLabel("Issue")
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {},
		),
	}

	vp.modeSelect = widget.NewSelect(
		[]string{"Strict", "Normal", "Lenient"},
		vp.onModeChanged,
	)
	vp.modeSelect.SetSelected("Normal")

	vp.autoFixBtn = widget.NewButton("Auto-Fix All Issues", vp.onAutoFixAll)
	vp.autoFixBtn.Disable()

	vp.resultsContainer = widget.NewLabel("")

	return vp
}

// BuildPanel builds the validation panel widget
func (vp *ValidationPanel) BuildPanel() fyne.CanvasObject {
	return container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Validation Mode:"),
			vp.modeSelect,
		),
		vp.summaryLabel,
		vp.resultsContainer,
		vp.autoFixBtn,
	)
}

// ValidateSaveData validates the provided save data
func (vp *ValidationPanel) ValidateSaveData(data *pr.PR) {
	vp.currentData = data

	// Run validation
	result := vp.validator.Validate(data)
	vp.lastResult = result

	// Update display
	vp.updateDisplay(result)
}

// updateDisplay updates the panel display with validation results
func (vp *ValidationPanel) updateDisplay(result models.ValidationResult) {
	// Update summary
	summaryText := fmt.Sprintf(
		"✓ Valid: %d errors, %d warnings, %d info",
		len(result.Errors), len(result.Warnings), len(result.Infomsgs),
	)
	if !result.Valid {
		summaryText = fmt.Sprintf(
			"✗ Invalid: %d errors, %d warnings, %d info",
			len(result.Errors), len(result.Warnings), len(result.Infomsgs),
		)
	}
	vp.summaryLabel.SetText(summaryText)

	// Enable/disable auto-fix button
	if len(result.Errors) > 0 || len(result.Warnings) > 0 {
		vp.autoFixBtn.Enable()
	} else {
		vp.autoFixBtn.Disable()
	}

	// Build results display
	vp.resultsContainer = vp.buildResultsDisplay(result)
}

// buildResultsDisplay builds the results container
func (vp *ValidationPanel) buildResultsDisplay(result models.ValidationResult) fyne.CanvasObject {
	boxes := make([]fyne.CanvasObject, 0)

	// Errors section
	if len(result.Errors) > 0 {
		boxes = append(boxes, vp.buildIssueSection("Errors", result.Errors, color.RGBA{200, 50, 50, 255}))
	}

	// Warnings section
	if len(result.Warnings) > 0 {
		boxes = append(boxes, vp.buildIssueSection("Warnings", result.Warnings, color.RGBA{200, 150, 50, 255}))
	}

	// Info section
	if len(result.Infomsgs) > 0 {
		boxes = append(boxes, vp.buildIssueSection("Information", result.Infomsgs, color.RGBA{50, 150, 200, 255}))
	}

	if len(boxes) == 0 {
		return widget.NewLabel("No validation issues found!")
	}

	return container.NewVBox(boxes...)
}

// buildIssueSection builds a section for issues of a specific severity
func (vp *ValidationPanel) buildIssueSection(title string, issues []models.ValidationIssue, titleColor color.Color) fyne.CanvasObject {
	titleText := canvas.NewText(title, titleColor)
	titleText.TextSize = 14
	titleText.TextStyle.Bold = true

	issueItems := make([]fyne.CanvasObject, len(issues))
	for i, issue := range issues {
		issueItems[i] = vp.buildIssueItem(&issue)
	}

	return container.NewVBox(
		titleText,
		container.NewVBox(issueItems...),
	)
}

// buildIssueItem builds a single issue item
func (vp *ValidationPanel) buildIssueItem(issue *models.ValidationIssue) fyne.CanvasObject {
	issueText := fmt.Sprintf(
		"%s: %s (Target: %s)",
		issue.Rule, issue.Message, issue.Target,
	)

	label := widget.NewLabel(issueText)
	label.Wrapping = fyne.TextWrapWord

	if issue.Fixable {
		fixBtn := widget.NewButton("Fix", func() {
			vp.onFixIssue(issue)
		})
		return container.NewHBox(label, fixBtn)
	}

	return label
}

// onModeChanged handles validation mode change
func (vp *ValidationPanel) onModeChanged(mode string) {
	var configMode models.ValidationMode

	switch mode {
	case "Strict":
		configMode = models.StrictMode
	case "Normal":
		configMode = models.NormalMode
	case "Lenient":
		configMode = models.LenientMode
	}

	config := vp.validator.GetConfig()
	config.Mode = configMode
	vp.validator.SetConfig(config)

	// Re-validate with new mode
	if vp.currentData != nil {
		vp.ValidateSaveData(vp.currentData)
	}
}

// onAutoFixAll handles auto-fix all button click
func (vp *ValidationPanel) onAutoFixAll() {
	if vp.currentData == nil {
		return
	}

	fixed, err := vp.validator.AutoFixIssues(vp.currentData)
	if err != nil {
		return
	}

	if fixed > 0 {
		// Re-validate to show results
		vp.ValidateSaveData(vp.currentData)

		// Call callback
		if vp.onIssueFixed != nil {
			vp.onIssueFixed()
		}
	}
}

// onFixIssue handles fixing a single issue
func (vp *ValidationPanel) onFixIssue(issue *models.ValidationIssue) {
	if vp.currentData == nil {
		return
	}

	// Note: This would need specific fix implementations
	// For now, we'll just re-validate
	if vp.onIssueFixed != nil {
		vp.onIssueFixed()
	}

	// Re-validate
	vp.ValidateSaveData(vp.currentData)
}

// SetOnIssueFixed sets the callback for when issue is fixed
func (vp *ValidationPanel) SetOnIssueFixed(callback func()) {
	vp.onIssueFixed = callback
}

// GetValidationStatus returns current validation status
func (vp *ValidationPanel) GetValidationStatus() (isValid bool, errors int, warnings int) {
	return vp.lastResult.Valid, len(vp.lastResult.Errors), len(vp.lastResult.Warnings)
}

// HasErrors checks if validation has errors
func (vp *ValidationPanel) HasErrors() bool {
	return len(vp.lastResult.Errors) > 0
}

// HasWarnings checks if validation has warnings
func (vp *ValidationPanel) HasWarnings() bool {
	return len(vp.lastResult.Warnings) > 0
}
