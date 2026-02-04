package editors

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	weaponsTextBox      fyne.CanvasObject
	shieldsTextBox      fyne.CanvasObject
	helmetTextBox       fyne.CanvasObject
	armorTextBox        fyne.CanvasObject
	relic1TextBox       fyne.CanvasObject
	relic2TextBox       fyne.CanvasObject
	itemsTextBox        fyne.CanvasObject
	allEquipmentTextBox fyne.CanvasObject
	importItemsTextBox  fyne.CanvasObject
	mapsTextBox         fyne.CanvasObject
)

func CreateTextBoxes() {
	if weaponsTextBox == nil {
		weaponsTextBox = container.NewVScroll(widget.NewRichTextWithText(weaponsText))
		shieldsTextBox = container.NewVScroll(widget.NewRichTextWithText(shieldsText))
		helmetTextBox = container.NewVScroll(widget.NewRichTextWithText(helmetText))
		armorTextBox = container.NewVScroll(widget.NewRichTextWithText(armorText))
		relic1TextBox = container.NewVScroll(widget.NewRichTextWithText(relicText1))
		relic2TextBox = container.NewVScroll(widget.NewRichTextWithText(relicText2Header + relicText2))
		itemsTextBox = container.NewVScroll(widget.NewRichTextWithText(itemsText))
		allEquipmentTextBox = container.NewVScroll(widget.NewRichTextWithText(weaponsText + shieldsText + helmetText + armorText + relicText1 + relicText2))
		importItemsTextBox = container.NewVScroll(widget.NewRichTextWithText(importantItemsText))
		mapsTextBox = container.NewVScroll(widget.NewRichTextWithText(mapText))
	}
}
