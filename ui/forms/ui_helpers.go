package forms

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ============ MODERN CARD COMPONENTS ============

// CreateModernCard creates a styled card with icon, title, description, content and action button
func CreateModernCard(icon string, title string, description string, content fyne.CanvasObject, actionBtn *widget.Button) fyne.CanvasObject {
	// Header with icon and title
	headerLabel := widget.NewLabel(icon + " " + title)
	headerLabel.TextStyle.Bold = true
	headerLabel.Alignment = fyne.TextAlignLeading

	// Description
	descLabel := widget.NewLabel(description)
	descLabel.Alignment = fyne.TextAlignLeading

	// Content area
	contentBox := container.NewVBox(content)

	// Action button - ensure it's clickable
	actionBtn.Importance = widget.HighImportance

	// Combine into card with proper layout
	cardContent := container.NewVBox(
		headerLabel,
		descLabel,
		widget.NewSeparator(),
		contentBox,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			actionBtn,
		),
	)

	// Dark glass background that fills available space
	bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 35, B: 40, A: 220})

	// Pad the content slightly so it doesn't touch edges
	padded := container.NewPadded(cardContent)

	// Ensure background resizes to match content/container
	return container.NewMax(bg, padded)
}

// CreateHeader creates a modern header section with dark glass styling
func CreateHeader(title string) fyne.CanvasObject {
	headerLabel := widget.NewLabel(title)
	headerLabel.TextStyle.Bold = true
	headerLabel.Alignment = fyne.TextAlignCenter

	divider := CreateDivider()

	return container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			headerLabel,
			layout.NewSpacer(),
		),
		divider,
	)
}

// CreateStatusBar creates a modern status bar with dark glass styling
func CreateStatusBar(statusLabel *widget.Label) fyne.CanvasObject {
	// Dark glass background
	bg := canvas.NewRectangle(color.NRGBA{R: 25, G: 50, B: 35, A: 210})
	topDivider := canvas.NewLine(color.NRGBA{R: 100, G: 220, B: 140, A: 255})
	topDivider.StrokeWidth = 2
	bottomDivider := canvas.NewLine(color.NRGBA{R: 100, G: 220, B: 140, A: 255})
	bottomDivider.StrokeWidth = 2

	// Status text in light color
	statusText := canvas.NewText(statusLabel.Text, color.NRGBA{R: 200, G: 255, B: 200, A: 255})
	statusText.TextSize = 12

	content := container.NewVBox(
		topDivider,
		container.NewCenter(statusText),
		bottomDivider,
	)

	// Make sure the background fills the available space
	return container.NewMax(bg, content)
}

// CreateDivider creates a subtle divider line
func CreateDivider() fyne.CanvasObject {
	divider := canvas.NewLine(color.RGBA{200, 200, 200, 200})
	divider.StrokeWidth = 1
	return container.NewVBox(divider)
}

// ============ ORIGINAL HELPERS ============

// CreateStyledSection creates a visually distinct section with title and content
func CreateStyledSection(title string, content fyne.CanvasObject) fyne.CanvasObject {
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle.Bold = true
	titleLabel.Alignment = fyne.TextAlignLeading

	separator := canvas.NewLine(color.RGBA{150, 150, 150, 200})
	separator.StrokeWidth = 1

	return container.NewVBox(
		titleLabel,
		separator,
		container.NewVBox(content),
	)
}

// CreateCardButton creates a styled button with higher importance
func CreateCardButton(label string, callback func()) *widget.Button {
	btn := widget.NewButton(label, callback)
	btn.Importance = widget.HighImportance
	return btn
}

// CreateFormRow creates a horizontal form row with label and input
func CreateFormRow(label string, input fyne.CanvasObject) fyne.CanvasObject {
	return container.NewBorder(
		nil,
		nil,
		widget.NewLabel(label),
		nil,
		input,
	)
}

// CreateInfoBox creates a styled info/status box
func CreateInfoBox(title string, content string) fyne.CanvasObject {
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle.Bold = true

	contentLabel := widget.NewLabel(content)
	contentLabel.Alignment = fyne.TextAlignLeading

	bg := canvas.NewRectangle(color.RGBA{200, 220, 240, 40})

	return container.NewStack(
		bg,
		container.NewVBox(
			container.NewVBox(titleLabel, contentLabel),
		),
	)
}

// CreateToolbar creates a styled toolbar with buttons
func CreateToolbar(buttons ...*widget.Button) fyne.CanvasObject {
	items := []fyne.CanvasObject{}
	for _, btn := range buttons {
		items = append(items, btn)
	}
	items = append(items, layout.NewSpacer())
	return container.NewHBox(items...)
}

// CreateSectionWithForm creates a section with form items
func CreateSectionWithForm(title string, items ...*widget.FormItem) fyne.CanvasObject {
	form := widget.NewForm(items...)
	return CreateStyledSection(title, form)
}

// CreateDialogButtonRow creates a horizontal button row for dialog buttons
func CreateDialogButtonRow(buttons ...*widget.Button) fyne.CanvasObject {
	items := []fyne.CanvasObject{layout.NewSpacer()}
	for _, btn := range buttons {
		items = append(items, btn)
	}
	return container.NewHBox(items...)
}

// CreateScrollableContent creates scrollable content with proper sizing
func CreateScrollableContent(content fyne.CanvasObject, minWidth, minHeight float32) fyne.CanvasObject {
	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(minWidth, minHeight))
	return scroll
}

// CreateButtonGroup creates a group of buttons with consistent spacing
func CreateButtonGroup(buttons ...*widget.Button) fyne.CanvasObject {
	items := []fyne.CanvasObject{}
	for i, btn := range buttons {
		items = append(items, btn)
		if i < len(buttons)-1 {
			items = append(items, layout.NewSpacer())
		}
	}
	return container.NewVBox(items...)
}

// CreateDividedContent creates a layout with left and right sections
func CreateDividedContent(left, right fyne.CanvasObject, minLeftWidth float32) fyne.CanvasObject {
	leftScroll := container.NewVScroll(left)
	leftScroll.SetMinSize(fyne.NewSize(minLeftWidth, 500))

	rightScroll := container.NewVScroll(right)
	rightScroll.SetMinSize(fyne.NewSize(minLeftWidth, 500))

	return container.NewHBox(leftScroll, rightScroll)
}

// CreatePaddedBox creates a container with padding
func CreatePaddedBox(left, right, top, bottom float32, content fyne.CanvasObject) fyne.CanvasObject {
	spacers := container.NewVBox()

	if top > 0 {
		spacers.Add(layout.NewSpacer())
	}

	inner := container.NewHBox()
	if left > 0 {
		inner.Add(layout.NewSpacer())
	}
	inner.Add(content)
	if right > 0 {
		inner.Add(layout.NewSpacer())
	}

	spacers.Add(inner)
	if bottom > 0 {
		spacers.Add(layout.NewSpacer())
	}

	return spacers
}

// CreateStatusLabel creates a label for status messages
func CreateStatusLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Alignment = fyne.TextAlignLeading
	return label
}

// CreateIconButton creates a button with an icon (uses emoji for now)
func CreateIconButton(icon string, label string, callback func()) *widget.Button {
	return widget.NewButton(icon+" "+label, callback)
}

// WrapDialogWithDarkBackground wraps dialog content with a dark glass background
func WrapDialogWithDarkBackground(content fyne.CanvasObject) fyne.CanvasObject {
	darkBg := canvas.NewRectangle(color.NRGBA{R: 18, G: 18, B: 24, A: 255})
	return container.NewMax(darkBg, content)
}
