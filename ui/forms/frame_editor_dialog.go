package forms

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io"
	"ffvi_editor/models"
)

// FrameEditorDialog provides frame timing and sequence editing
type FrameEditorDialog struct {
	dialog           *widget.PopUp
	sequencer        *io.FrameSequencer
	container        *fyne.Container
	frameListBox     *widget.List
	frameDurationSp  *widget.Entry
	selectedFrameIdx int
	onClose          func()
}

// NewFrameEditorDialog creates a new frame editor dialog
func NewFrameEditorDialog(frames []*models.FF6Sprite) *FrameEditorDialog {
	if len(frames) == 0 {
		return nil
	}

	fed := &FrameEditorDialog{
		selectedFrameIdx: 0,
	}

	// Create sequencer
	fed.sequencer = io.NewFrameSequencer(frames)
	if fed.sequencer == nil {
		return nil
	}

	fed.buildUI()
	return fed
}

// buildUI constructs the dialog UI
func (fed *FrameEditorDialog) buildUI() {
	// Frame list
	fed.frameListBox = widget.NewList(
		func() int {
			return fed.sequencer.GetFrameCount()
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Frame "),
				widget.NewLabel(""),
				widget.NewLabel("Duration: "),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			duration := fed.sequencer.GetFrameDuration(id)

			labels := hbox.Objects
			frameLabel := labels[1].(*widget.Label)
			frameLabel.SetText(strconv.Itoa(id))

			durationLabel := labels[3].(*widget.Label)
			durationLabel.SetText(strconv.Itoa(duration) + "ms")
		},
	)

	fed.frameListBox.OnSelected = fed.onFrameSelected

	// Frame duration editor
	fed.frameDurationSp = widget.NewEntry()
	fed.frameDurationSp.SetPlaceHolder("100")
	fed.frameDurationSp.OnSubmitted = fed.onDurationSubmitted

	// Control buttons
	applyDurationBtn := widget.NewButton("Apply Duration", fed.onApplyDuration)
	bulkDurationBtn := widget.NewButton("Bulk Duration", fed.onBulkDuration)
	autoTimingBtn := widget.NewButton("Auto Timing", fed.onAutoTiming)

	// Sequence operations
	insertBtn := widget.NewButton("Insert Frame", fed.onInsertFrame)
	duplicateBtn := widget.NewButton("Duplicate Frame", fed.onDuplicateFrame)
	removeBtn := widget.NewButton("Remove Frame", fed.onRemoveFrame)
	moveUpBtn := widget.NewButton("Move Up", fed.onMoveFrameUp)
	moveDownBtn := widget.NewButton("Move Down", fed.onMoveFrameDown)

	// Info section
	infoLabel := widget.NewLabel("")
	fed.updateInfoLabel(infoLabel)

	// Build layout
	durationSection := container.NewVBox(
		widget.NewLabel("Frame Duration (ms):"),
		fed.frameDurationSp,
		container.NewHBox(applyDurationBtn, bulkDurationBtn),
		autoTimingBtn,
	)

	sequenceSection := container.NewVBox(
		widget.NewLabel("Sequence Operations:"),
		container.NewHBox(insertBtn, duplicateBtn, removeBtn),
		container.NewHBox(moveUpBtn, moveDownBtn),
	)

	fed.container = container.NewVBox(
		widget.NewLabel("Frame Sequence Editor"),
		widget.NewCard("Frames", "", container.NewVBox(
			container.NewBorder(nil, nil, widget.NewLabel("Frames:"), nil, fed.frameListBox),
			infoLabel,
		)),
		widget.NewCard("Timing", "", durationSection),
		widget.NewCard("Operations", "", sequenceSection),
	)
}

// Show displays the dialog
func (fed *FrameEditorDialog) Show(parent fyne.Window) {
	dialog := container.NewVBox(
		fed.container,
		container.NewHBox(
			widget.NewButton("Export to Animation", fed.onExportToAnimation),
			widget.NewButton("Close", func() {
				if fed.onClose != nil {
					fed.onClose()
				}
				fed.Hide()
			}),
		),
	)

	fed.dialog = widget.NewPopUp(dialog, parent.Canvas())
	fed.dialog.ShowAtPosition(fyne.NewPos(150, 150))
}

// Hide hides the dialog
func (fed *FrameEditorDialog) Hide() {
	if fed.dialog != nil {
		fed.dialog.Hide()
	}
}

// Event handlers

func (fed *FrameEditorDialog) onFrameSelected(id widget.ListItemID) {
	fed.selectedFrameIdx = id
	duration := fed.sequencer.GetFrameDuration(id)
	fed.frameDurationSp.SetText(strconv.Itoa(duration))
	fed.frameListBox.Refresh()
}

func (fed *FrameEditorDialog) onDurationSubmitted(s string) {
	fed.onApplyDuration()
}

func (fed *FrameEditorDialog) onApplyDuration() {
	duration, err := strconv.Atoi(fed.frameDurationSp.Text)
	if err != nil || duration <= 0 {
		return
	}

	fed.sequencer.SetFrameDuration(fed.selectedFrameIdx, duration)
	fed.frameListBox.Refresh()
}

func (fed *FrameEditorDialog) onBulkDuration() {
	duration, err := strconv.Atoi(fed.frameDurationSp.Text)
	if err != nil || duration <= 0 {
		return
	}

	fed.sequencer.SetAllFrameDurations(duration)
	fed.frameListBox.Refresh()
}

func (fed *FrameEditorDialog) onAutoTiming() {
	fed.sequencer.AutoTiming(100)
	fed.frameListBox.Refresh()
}

func (fed *FrameEditorDialog) onInsertFrame() {
	// Insert duplicate of current frame
	if fed.selectedFrameIdx < fed.sequencer.GetFrameCount() {
		fed.sequencer.DuplicateFrame(fed.selectedFrameIdx)
		fed.frameListBox.Refresh()
	}
}

func (fed *FrameEditorDialog) onDuplicateFrame() {
	if fed.selectedFrameIdx < fed.sequencer.GetFrameCount() {
		fed.sequencer.DuplicateFrame(fed.selectedFrameIdx)
		fed.frameListBox.Refresh()
	}
}

func (fed *FrameEditorDialog) onRemoveFrame() {
	if fed.sequencer.GetFrameCount() > 1 {
		fed.sequencer.RemoveFrame(fed.selectedFrameIdx)
		if fed.selectedFrameIdx >= fed.sequencer.GetFrameCount() {
			fed.selectedFrameIdx = fed.sequencer.GetFrameCount() - 1
		}
		fed.frameListBox.Refresh()
	}
}

func (fed *FrameEditorDialog) onMoveFrameUp() {
	if fed.selectedFrameIdx > 0 {
		fed.sequencer.MoveFrame(fed.selectedFrameIdx, fed.selectedFrameIdx-1)
		fed.selectedFrameIdx--
		fed.frameListBox.Refresh()
		fed.frameListBox.Select(fed.selectedFrameIdx)
	}
}

func (fed *FrameEditorDialog) onMoveFrameDown() {
	if fed.selectedFrameIdx < fed.sequencer.GetFrameCount()-1 {
		fed.sequencer.MoveFrame(fed.selectedFrameIdx, fed.selectedFrameIdx+1)
		fed.selectedFrameIdx++
		fed.frameListBox.Refresh()
		fed.frameListBox.Select(fed.selectedFrameIdx)
	}
}

func (fed *FrameEditorDialog) onExportToAnimation() {
	// Export from sequencer - returns AnimationData
	_ = fed.sequencer.ExportToAnimation(models.PlayContinuous)
	// In real implementation, would return this to parent
}

// Helper methods

func (fed *FrameEditorDialog) updateInfoLabel(label *widget.Label) {
	frameCount := fed.sequencer.GetFrameCount()
	totalDuration := fed.sequencer.GetTotalDuration()
	avgFPS := fed.sequencer.CalculateAverageFPS()

	label.SetText(fmt.Sprintf(
		"Frames: %d | Duration: %dms | Avg FPS: %.1f",
		frameCount, totalDuration, avgFPS,
	))
}

// GetSequencer returns the frame sequencer
func (fed *FrameEditorDialog) GetSequencer() *io.FrameSequencer {
	return fed.sequencer
}

// SetOnClose sets the close callback
func (fed *FrameEditorDialog) SetOnClose(callback func()) {
	fed.onClose = callback
}
