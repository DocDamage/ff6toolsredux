package forms

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/settings"
)

// PreferencesDialog manages application preferences
type PreferencesDialog struct {
	window          fyne.Window
	settingsManager *settings.Manager
}

// NewPreferencesDialog creates a new preferences dialog
func NewPreferencesDialog(window fyne.Window, settingsManager *settings.Manager) *PreferencesDialog {
	return &PreferencesDialog{
		window:          window,
		settingsManager: settingsManager,
	}
}

// Show displays the preferences dialog
func (p *PreferencesDialog) Show() {
	currentSettings := p.settingsManager.Get()

	// Create tabs for different preference categories
	tabs := container.NewAppTabs(
		container.NewTabItem("General", p.buildGeneralTab(currentSettings)),
		container.NewTabItem("Editor", p.buildEditorTab(currentSettings)),
		container.NewTabItem("Backup", p.buildBackupTab(currentSettings)),
		container.NewTabItem("Validation", p.buildValidationTab(currentSettings)),
		container.NewTabItem("Shortcuts", p.buildShortcutsTab(currentSettings)),
		container.NewTabItem("Advanced", p.buildAdvancedTab(currentSettings)),
	)

	// Save button
	saveBtn := widget.NewButton("Save", func() {
		if err := p.settingsManager.Save(); err != nil {
			dialog.ShowError(fmt.Errorf("failed to save settings: %w", err), p.window)
		} else {
			dialog.ShowInformation("Success", "Settings saved successfully!", p.window)
		}
	})

	// Reset button
	resetBtn := widget.NewButton("Reset to Defaults", func() {
		dialog.ShowConfirm("Reset Settings",
			"Are you sure you want to reset all settings to defaults?",
			func(confirm bool) {
				if confirm {
					p.settingsManager.Reset()
					dialog.ShowInformation("Reset", "Settings reset to defaults. Please reopen preferences.", p.window)
				}
			}, p.window)
	})

	content := container.NewBorder(
		nil,
		container.NewHBox(resetBtn, saveBtn),
		nil,
		nil,
		tabs,
	)

	d := dialog.NewCustom("Preferences", "Close", WrapDialogWithDarkBackground(content), p.window)
	d.Resize(fyne.NewSize(700, 500))
	d.Show()
}

// buildGeneralTab builds the general preferences tab
func (p *PreferencesDialog) buildGeneralTab(settings *settings.Settings) fyne.CanvasObject {
	themeSelect := widget.NewSelect([]string{"Dark", "Light"}, func(value string) {
		settings.Theme = value
	})
	if settings.Theme == "dark" {
		themeSelect.SetSelected("Dark")
	} else {
		themeSelect.SetSelected("Light")
	}

	languageSelect := widget.NewSelect([]string{"English", "Japanese"}, func(value string) {
		if value == "English" {
			settings.Language = "en"
		} else {
			settings.Language = "ja"
		}
	})
	languageSelect.SetSelected("English")

	autoSaveCheck := widget.NewCheck("Enable auto-save", func(checked bool) {
		settings.AutoSave = checked
	})
	autoSaveCheck.SetChecked(settings.AutoSave)

	autoSaveDelayEntry := widget.NewEntry()
	autoSaveDelayEntry.SetText(fmt.Sprintf("%d", settings.AutoSaveDelay))
	autoSaveDelayEntry.SetPlaceHolder("Seconds")

	return container.NewVBox(
		widget.NewLabel("General Settings"),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Theme", themeSelect),
			widget.NewFormItem("Language", languageSelect),
		),
		autoSaveCheck,
		widget.NewForm(
			widget.NewFormItem("Auto-save delay", autoSaveDelayEntry),
		),
	)
}

// buildEditorTab builds the editor preferences tab
func (p *PreferencesDialog) buildEditorTab(settings *settings.Settings) fyne.CanvasObject {
	confirmCloseCheck := widget.NewCheck("Confirm on close", func(checked bool) {
		settings.ConfirmOnClose = checked
	})
	confirmCloseCheck.SetChecked(settings.ConfirmOnClose)

	showUndoCheck := widget.NewCheck("Show undo history", func(checked bool) {
		settings.ShowUndoHistory = checked
	})
	showUndoCheck.SetChecked(settings.ShowUndoHistory)

	dragDropCheck := widget.NewCheck("Enable drag & drop", func(checked bool) {
		settings.EnableDragDrop = checked
	})
	dragDropCheck.SetChecked(settings.EnableDragDrop)

	tooltipsCheck := widget.NewCheck("Show tooltips", func(checked bool) {
		settings.ShowTooltips = checked
	})
	tooltipsCheck.SetChecked(settings.ShowTooltips)

	maxUndoEntry := widget.NewEntry()
	maxUndoEntry.SetText(fmt.Sprintf("%d", settings.MaxUndoSteps))
	maxUndoEntry.SetPlaceHolder("Steps")

	fontSizeEntry := widget.NewEntry()
	fontSizeEntry.SetText(fmt.Sprintf("%d", settings.EditorFontSize))
	fontSizeEntry.SetPlaceHolder("Size")

	return container.NewVBox(
		widget.NewLabel("Editor Settings"),
		widget.NewSeparator(),
		confirmCloseCheck,
		showUndoCheck,
		dragDropCheck,
		tooltipsCheck,
		widget.NewForm(
			widget.NewFormItem("Max undo steps", maxUndoEntry),
			widget.NewFormItem("Editor font size", fontSizeEntry),
		),
	)
}

// buildBackupTab builds the backup preferences tab
func (p *PreferencesDialog) buildBackupTab(settings *settings.Settings) fyne.CanvasObject {
	autoBackupCheck := widget.NewCheck("Enable auto-backup", func(checked bool) {
		settings.AutoBackup = checked
	})
	autoBackupCheck.SetChecked(settings.AutoBackup)

	backupsToKeepEntry := widget.NewEntry()
	backupsToKeepEntry.SetText(fmt.Sprintf("%d", settings.BackupsToKeep))
	backupsToKeepEntry.SetPlaceHolder("Number")

	backupLocationEntry := widget.NewEntry()
	backupLocationEntry.SetText(settings.BackupLocation)
	backupLocationEntry.SetPlaceHolder("Leave empty for default location")

	browseBtn := widget.NewButton("Browse...", func() {
		dialog.ShowInformation("Browse", "Directory picker would open here.", p.window)
	})

	return container.NewVBox(
		widget.NewLabel("Backup Settings"),
		widget.NewSeparator(),
		autoBackupCheck,
		widget.NewForm(
			widget.NewFormItem("Backups to keep", backupsToKeepEntry),
			widget.NewFormItem("Backup location", container.NewBorder(nil, nil, nil, browseBtn, backupLocationEntry)),
		),
	)
}

// buildValidationTab builds the validation preferences tab
func (p *PreferencesDialog) buildValidationTab(settings *settings.Settings) fyne.CanvasObject {
	validationLevelSelect := widget.NewSelect(
		[]string{"Strict", "Normal", "Permissive"},
		func(value string) {
			settings.ValidationLevel = value
		},
	)
	validationLevelSelect.SetSelected(settings.ValidationLevel)

	autoFixCheck := widget.NewCheck("Auto-fix issues", func(checked bool) {
		settings.AutoFix = checked
	})
	autoFixCheck.SetChecked(settings.AutoFix)

	warnOnSaveCheck := widget.NewCheck("Warn before saving invalid data", func(checked bool) {
		settings.WarnOnSave = checked
	})
	warnOnSaveCheck.SetChecked(settings.WarnOnSave)

	return container.NewVBox(
		widget.NewLabel("Validation Settings"),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Validation level", validationLevelSelect),
		),
		autoFixCheck,
		warnOnSaveCheck,
		widget.NewLabel("\nValidation Levels:"),
		widget.NewLabel("• Strict: Enforce all rules strictly"),
		widget.NewLabel("• Normal: Standard validation (recommended)"),
		widget.NewLabel("• Permissive: Allow more flexibility"),
	)
}

// buildShortcutsTab builds the shortcuts preferences tab
func (p *PreferencesDialog) buildShortcutsTab(settings *settings.Settings) fyne.CanvasObject {
	shortcutList := widget.NewList(
		func() int { return len(settings.Shortcuts) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Action"),
				widget.NewButton("Ctrl+X", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// Would populate with actual shortcuts
		},
	)

	resetShortcutsBtn := widget.NewButton("Reset All Shortcuts", func() {
		dialog.ShowConfirm("Reset Shortcuts",
			"Reset all keyboard shortcuts to defaults?",
			func(confirm bool) {
				if confirm {
					// Reset shortcuts to defaults
					settings.Shortcuts = map[string]string{
						"save":        "Ctrl+S",
						"undo":        "Ctrl+Z",
						"redo":        "Ctrl+Y",
						"find":        "Ctrl+F",
						"palette":     "Ctrl+Shift+P",
						"close":       "Ctrl+W",
						"new":         "Ctrl+N",
						"open":        "Ctrl+O",
						"backup":      "Ctrl+B",
						"validate":    "Ctrl+Shift+V",
						"preferences": "Ctrl+,",
					}
					dialog.ShowInformation("Reset", "Shortcuts reset to defaults.", p.window)
				}
			}, p.window)
	})

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Keyboard Shortcuts"),
			widget.NewSeparator(),
			widget.NewLabel("Click a shortcut to edit it"),
		),
		resetShortcutsBtn,
		nil,
		nil,
		shortcutList,
	)
}

// buildAdvancedTab builds the advanced preferences tab
func (p *PreferencesDialog) buildAdvancedTab(settings *settings.Settings) fyne.CanvasObject {
	pluginsCheck := widget.NewCheck("Enable plugins", func(checked bool) {
		settings.EnablePlugins = checked
	})
	pluginsCheck.SetChecked(settings.EnablePlugins)

	scriptingCheck := widget.NewCheck("Enable scripting", func(checked bool) {
		settings.EnableScripting = checked
	})
	scriptingCheck.SetChecked(settings.EnableScripting)

	achievementsCheck := widget.NewCheck("Enable achievements", func(checked bool) {
		settings.EnableAchievements = checked
	})
	achievementsCheck.SetChecked(settings.EnableAchievements)

	debugCheck := widget.NewCheck("Show debug info", func(checked bool) {
		settings.ShowDebugInfo = checked
	})
	debugCheck.SetChecked(settings.ShowDebugInfo)

	logLevelSelect := widget.NewSelect(
		[]string{"Debug", "Info", "Warn", "Error"},
		func(value string) {
			settings.LogLevel = value
		},
	)
	logLevelSelect.SetSelected(settings.LogLevel)

	return container.NewVBox(
		widget.NewLabel("Advanced Settings"),
		widget.NewSeparator(),
		pluginsCheck,
		scriptingCheck,
		achievementsCheck,
		debugCheck,
		widget.NewForm(
			widget.NewFormItem("Log level", logLevelSelect),
		),
		widget.NewSeparator(),
		widget.NewLabel("⚠ Warning: Changing these settings may affect stability."),
	)
}
