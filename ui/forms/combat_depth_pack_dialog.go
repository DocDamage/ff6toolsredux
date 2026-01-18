package forms

import (
	"context"
	"ffvi_editor/io/pr"
	"fmt"
	"image/color"

	"ffvi_editor/scripting"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// CombatDepthPackDialog provides a comprehensive interface for combat customization
type CombatDepthPackDialog struct {
	window    fyne.Window
	save      *pr.PR
	dialog    dialog.Dialog
	statusMsg *widget.Label
}

// NewCombatDepthPackDialog creates and shows the Combat Depth Pack dialog
func NewCombatDepthPackDialog(win fyne.Window, save *pr.PR) *CombatDepthPackDialog {
	cdpd := &CombatDepthPackDialog{
		window: win,
		save:   save,
	}
	cdpd.buildUI()
	cdpd.Show()
	return cdpd
}

// Show displays the dialog
func (cdpd *CombatDepthPackDialog) Show() {
	if cdpd.dialog != nil {
		cdpd.dialog.Show()
	}
}

// Hide hides the dialog
func (cdpd *CombatDepthPackDialog) Hide() {
	if cdpd.dialog != nil {
		cdpd.dialog.Hide()
	}
}

// buildUI constructs the complete dialog UI with modern styling
func (cdpd *CombatDepthPackDialog) buildUI() {
	// Status message display with modern styling
	cdpd.statusMsg = CreateStatusLabel("Ready to apply combat customizations")
	cdpd.statusMsg.Alignment = fyne.TextAlignCenter

	// ===== ENCOUNTER TUNER SECTION =====
	zoneEntry := widget.NewEntry()
	zoneEntry.SetPlaceHolder("e.g., Mt. Kolts")

	rateEntry := widget.NewEntry()
	rateEntry.SetText("1.0")
	rateEntry.SetPlaceHolder("0.5 - 2.0")

	eliteEntry := widget.NewEntry()
	eliteEntry.SetText("0.10")
	eliteEntry.SetPlaceHolder("0.0 - 1.0")

	applyEncounterBtn := CreateCardButton("⚙ Apply Tuning", func() {
		if zoneEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("zone cannot be empty"), cdpd.window)
			return
		}
		code := scripting.BuildEncounterScript(zoneEntry.Text, rateEntry.Text, eliteEntry.Text)
		_, err := scripting.RunSnippetWithSave(context.Background(), code, cdpd.save)
		if err != nil {
			dialog.ShowError(err, cdpd.window)
			return
		}
		cdpd.statusMsg.SetText(fmt.Sprintf("✓ Applied to %s", zoneEntry.Text))
		dialog.ShowInformation("Success", fmt.Sprintf("Zone: %s\nRate: %s\nElite: %s", zoneEntry.Text, rateEntry.Text, eliteEntry.Text), cdpd.window)
	})

	encounterCard := CreateModernCard(
		"🎲",
		"Dynamic Encounter Tuner",
		"Customize encounter rates and elite chances",
		container.NewVBox(
			CreateFormRow("Zone", zoneEntry),
			CreateFormRow("Spawn Rate", rateEntry),
			CreateFormRow("Elite Chance", eliteEntry),
		),
		applyEncounterBtn,
	)

	// ===== BOSS REMIX & AFFIXES SECTION =====
	affixEntry := widget.NewEntry()
	affixEntry.SetPlaceHolder("enraged, arcane_shield, ...")

	applyBossBtn := CreateCardButton("👹 Generate", func() {
		if affixEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("affixes cannot be empty"), cdpd.window)
			return
		}
		code := scripting.BuildBossScript(affixEntry.Text)
		_, err := scripting.RunSnippetWithSave(context.Background(), code, cdpd.save)
		if err != nil {
			dialog.ShowError(err, cdpd.window)
			return
		}
		cdpd.statusMsg.SetText("✓ Boss remix applied")
		dialog.ShowInformation("Success", fmt.Sprintf("Affixes applied:\n%s", affixEntry.Text), cdpd.window)
	})

	bossCard := CreateModernCard(
		"👹",
		"Boss Remix & Affixes",
		"Apply special effects to boss encounters",
		container.NewVBox(
			CreateFormRow("Affixes", affixEntry),
		),
		applyBossBtn,
	)

	// ===== AI COMPANION DIRECTOR SECTION =====
	profileEntry := widget.NewEntry()
	profileEntry.SetPlaceHolder("aggressive, support, balanced")

	riskEntry := widget.NewSelect([]string{"low", "normal", "high"}, func(s string) {})
	riskEntry.SetSelected("normal")

	applyCompanionBtn := CreateCardButton("🤖 Save", func() {
		if profileEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("profile cannot be empty"), cdpd.window)
			return
		}
		code := scripting.BuildCompanionScript(profileEntry.Text, riskEntry.Selected)
		_, err := scripting.RunSnippetWithSave(context.Background(), code, cdpd.save)
		if err != nil {
			dialog.ShowError(err, cdpd.window)
			return
		}
		cdpd.statusMsg.SetText(fmt.Sprintf("✓ Profile saved: %s", profileEntry.Text))
		dialog.ShowInformation("Success", fmt.Sprintf("Profile: %s\nRisk: %s", profileEntry.Text, riskEntry.Selected), cdpd.window)
	})

	companionCard := CreateModernCard(
		"🤖",
		"AI Companion Director",
		"Configure AI behavior profiles",
		container.NewVBox(
			CreateFormRow("Profile", profileEntry),
			CreateFormRow("Risk Level", riskEntry),
		),
		applyCompanionBtn,
	)

	// ===== SMOKE TESTS SECTION =====
	smokeBtn := CreateCardButton("🧪 Run Tests", func() {
		_, err := scripting.RunSnippetWithSave(context.Background(), scripting.BuildSmokeScript(), cdpd.save)
		if err != nil {
			dialog.ShowError(err, cdpd.window)
			return
		}
		cdpd.statusMsg.SetText("✓ Tests passed")
		dialog.ShowInformation("Success", "All smoke tests completed successfully", cdpd.window)
	})

	smokeCard := container.NewVBox(
		CreateDivider(),
		container.NewCenter(
			container.NewVBox(
				widget.NewLabel("Quality Assurance"),
				smokeBtn,
			),
		),
	)

	// ===== MAIN LAYOUT WITH MODERN DESIGN =====
	mainContent := container.NewVBox(
		// Header
		CreateHeader("⚔ Combat Configuration"),
		CreateStatusBar(cdpd.statusMsg),

		// Cards
		encounterCard,
		bossCard,
		companionCard,
		smokeCard,

		// Padding
		layout.NewSpacer(),
	)

	// Wrap in scroll for long content with dark background
	scrollContent := CreateScrollableContent(mainContent, 550, 650)

	// Modern button row
	closeBtn := widget.NewButton("Close", func() { cdpd.Hide() })
	closeBtn.Importance = widget.LowImportance

	buttonRow := CreateDialogButtonRow(closeBtn)

	// Create dark glass wrapper for entire dialog
	darkBg := canvas.NewRectangle(color.NRGBA{R: 18, G: 18, B: 24, A: 255})

	// Create dialog with border layout wrapped in dark background
	contentWithBg := container.NewMax(
		darkBg,
		container.NewBorder(
			nil,           // Top
			buttonRow,     // Bottom
			nil,           // Left
			nil,           // Right
			scrollContent, // Center
		),
	)

	cdpd.dialog = dialog.NewCustom(
		"Combat Depth Pack",
		"",
		contentWithBg,
		cdpd.window,
	)
	cdpd.dialog.Resize(fyne.NewSize(650, 750))
}
