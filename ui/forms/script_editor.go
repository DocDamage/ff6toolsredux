package forms

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/scripting"
)

// ScriptEditorDialog provides a UI for writing and running Lua scripts
type ScriptEditorDialog struct {
	window fyne.Window
	vm     *scripting.VM
}

// NewScriptEditorDialog creates a new script editor dialog
func NewScriptEditorDialog(window fyne.Window) *ScriptEditorDialog {
	return &ScriptEditorDialog{
		window: window,
		vm:     scripting.NewVM(0),
	}
}

// Show displays the script editor dialog
func (s *ScriptEditorDialog) Show() {
	// Script text area
	scriptEntry := widget.NewMultiLineEntry()
	scriptEntry.SetPlaceHolder("-- Write your Lua script here\n-- Example:\nsetCharacterLevel(0, 99)\nprint('Terra set to level 99!')")
	scriptEntry.Wrapping = fyne.TextWrapWord
	scriptEntry.SetMinRowsVisible(15)

	// Built-in scripts dropdown
	builtInScripts := scripting.BuiltInScripts()
	scriptNames := make([]string, 0, len(builtInScripts))
	for name := range builtInScripts {
		scriptNames = append(scriptNames, name)
	}

	builtInSelect := widget.NewSelect(scriptNames, func(value string) {
		if script, ok := builtInScripts[value]; ok {
			scriptEntry.SetText(script)
		}
	})
	builtInSelect.PlaceHolder = "Load built-in script..."

	// Output area
	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetPlaceHolder("Script output will appear here...")
	outputEntry.Wrapping = fyne.TextWrapWord
	outputEntry.SetMinRowsVisible(5)
	outputEntry.Disable()

	// Run button
	runBtn := widget.NewButton("Run Script", func() {
		script := scriptEntry.Text
		if script == "" {
			dialog.ShowError(fmt.Errorf("please enter a script"), s.window)
			return
		}

		outputEntry.SetText("Executing script...\n")

		// Execute script
		ctx := context.Background()
		err := s.vm.Execute(ctx, script)
		if err != nil {
			outputEntry.SetText(fmt.Sprintf("Error: %v\n\nNote: Full Lua VM requires github.com/yuin/gopher-lua library", err))
			dialog.ShowError(err, s.window)
		} else {
			outputEntry.SetText("Script executed successfully!\n\nNote: Actual execution requires Lua VM integration.")
		}
	})

	// Stop button
	stopBtn := widget.NewButton("Stop", func() {
		// s.vm.Cancel() // TODO: Implement cancellation
		outputEntry.SetText(outputEntry.Text + "\nScript stopped by user.")
	})
	stopBtn.Disable()

	// Clear button
	clearBtn := widget.NewButton("Clear", func() {
		scriptEntry.SetText("")
		outputEntry.SetText("")
	})

	// Load/Save buttons
	loadBtn := widget.NewButton("Load Script", func() {
		dialog.ShowInformation("Load Script",
			"File picker for loading .lua scripts will be implemented here.",
			s.window)
	})

	saveBtn := widget.NewButton("Save Script", func() {
		dialog.ShowInformation("Save Script",
			"File picker for saving .lua scripts will be implemented here.",
			s.window)
	})

	// API reference button
	apiBtn := widget.NewButton("API Reference", func() {
		s.showAPIReference()
	})

	// Build layout
	toolbar := container.NewHBox(
		builtInSelect,
		widget.NewSeparator(),
		runBtn,
		stopBtn,
		clearBtn,
		widget.NewSeparator(),
		loadBtn,
		saveBtn,
		apiBtn,
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Lua Script Editor"),
			toolbar,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			widget.NewLabel("Output:"),
			outputEntry,
		),
		nil,
		nil,
		scriptEntry,
	)

	d := dialog.NewCustom("Script Editor", "Close", content, s.window)
	d.Resize(fyne.NewSize(800, 600))
	d.Show()
}

// showAPIReference shows the scripting API reference
func (s *ScriptEditorDialog) showAPIReference() {
	apiText := `
Lua Scripting API Reference

CHARACTER FUNCTIONS:
  getCharacter(id) - Get character data
  setCharacterLevel(id, level) - Set character level (1-99)
  setCharacterHP(id, hp) - Set character HP (1-9999)
  setCharacterMP(id, mp) - Set character MP (1-999)
  setCharacterStat(id, stat, value) - Set stat (vigor, speed, stamina, magic)

INVENTORY FUNCTIONS:
  getItemCount(itemId) - Get item count
  setItemCount(itemId, count) - Set item count (0-99)
  addItem(itemId, count) - Add items to inventory

PARTY FUNCTIONS:
  getPartyMembers() - Get current party
  setPartyMembers({id1, id2, ...}) - Set party members

MAGIC FUNCTIONS:
  hasMagic(charId, spellId) - Check if character has spell
  learnMagic(charId, spellId) - Teach spell to character

UTILITY FUNCTIONS:
  print(message) - Print to output
  log(message) - Log message
  min(a, b) - Get minimum value
  max(a, b) - Get maximum value
  clamp(value, min, max) - Clamp value to range

CHARACTER IDs:
  0=Terra, 1=Locke, 2=Cyan, 3=Shadow, 4=Edgar, 5=Sabin
  6=Celes, 7=Strago, 8=Relm, 9=Setzer, 10=Mog, 11=Gau
  12=Gogo, 13=Umaro
`

	apiEntry := widget.NewMultiLineEntry()
	apiEntry.SetText(apiText)
	apiEntry.Wrapping = fyne.TextWrapWord
	apiEntry.Disable()

	d := dialog.NewCustom("API Reference", "Close", apiEntry, s.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
