package forms

import (
	"fmt"
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/io"
	"ffvi_editor/models"
)

// AnimationPlayerDialog provides animation playback controls and preview
type AnimationPlayerDialog struct {
	dialog     *widget.PopUp
	animation  *models.AnimationData
	controller *io.AnimationController
	container  *fyne.Container

	previewWidget fyne.CanvasObject
	previewImage  *canvas.Image
	frameImages   []image.Image

	playButton       *widget.Button
	pauseButton      *widget.Button
	stopButton       *widget.Button
	frameSlider      *widget.Slider
	frameLabel       *widget.Label
	speedSlider      *widget.Slider
	speedLabel       *widget.Label
	loopModeSelect   *widget.Select
	fpsLabel         *widget.Label
	timeLabel        *widget.Label
	autoPlayCheckbox *widget.Check

	tickCount     int64
	isAutoPlaying bool
	onClose       func()
}

// NewAnimationPlayerDialog creates a new animation player dialog
func NewAnimationPlayerDialog(animation *models.AnimationData) *AnimationPlayerDialog {
	if animation == nil || len(animation.Frames) == 0 {
		return nil
	}

	apd := &AnimationPlayerDialog{
		animation:     animation,
		isAutoPlaying: false,
	}

	apd.frameImages = decodeFrameImages(animation)

	var err error
	apd.controller, err = io.NewAnimationController(animation)
	if err != nil {
		return nil
	}

	apd.buildUI()
	return apd
}

// buildUI constructs the dialog UI
func (apd *AnimationPlayerDialog) buildUI() {
	apd.previewWidget = apd.createPreviewWidget()

	apd.playButton = widget.NewButton("▶ Play", apd.onPlay)
	apd.pauseButton = widget.NewButton("⏸ Pause", apd.onPause)
	apd.pauseButton.Disable()
	apd.stopButton = widget.NewButton("⏹ Stop", apd.onStop)

	apd.frameSlider = widget.NewSlider(0, float64(len(apd.animation.Frames)-1))
	apd.frameSlider.OnChanged = apd.onFrameSliderChanged

	apd.frameLabel = widget.NewLabel("Frame: 0 / " + strconv.Itoa(len(apd.animation.Frames)))

	apd.speedSlider = widget.NewSlider(0.5, 2.0)
	if apd.animation.DefaultSpeed > 0 {
		apd.speedSlider.Value = float64(apd.animation.DefaultSpeed)
	} else {
		apd.speedSlider.Value = 1.0
	}
	apd.speedSlider.OnChanged = apd.onSpeedChanged
	apd.speedLabel = widget.NewLabel(fmt.Sprintf("Speed: %.2fx", apd.speedSlider.Value))

	apd.loopModeSelect = widget.NewSelect(
		[]string{"Once", "Loop", "PingPong"},
		apd.onLoopModeChanged,
	)
	switch apd.animation.PlaybackMode {
	case models.PlayOnce:
		apd.loopModeSelect.SetSelected("Once")
	case models.PlayPingPong:
		apd.loopModeSelect.SetSelected("PingPong")
	default:
		apd.loopModeSelect.SetSelected("Loop")
	}

	apd.fpsLabel = widget.NewLabel(fmt.Sprintf("FPS: %.1f", apd.calculateFPS()))
	apd.timeLabel = widget.NewLabel("Time: 0ms")

	apd.autoPlayCheckbox = widget.NewCheck("Auto-play on open", apd.onAutoPlayToggled)
	apd.autoPlayCheckbox.SetChecked(apd.isAutoPlaying)

	controlsContainer := container.NewVBox(
		container.NewHBox(apd.playButton, apd.pauseButton, apd.stopButton),
		container.NewVBox(
			widget.NewLabel("Frame Navigation:"),
			apd.frameSlider,
			apd.frameLabel,
		),
		container.NewVBox(
			widget.NewLabel("Playback Speed:"),
			apd.speedSlider,
			apd.speedLabel,
		),
		container.NewVBox(
			widget.NewLabel("Loop Mode:"),
			apd.loopModeSelect,
		),
		container.NewHBox(apd.fpsLabel, apd.timeLabel),
		apd.autoPlayCheckbox,
	)

	apd.container = container.NewVBox(
		apd.previewWidget,
		controlsContainer,
	)
}

// createPreviewWidget creates a preview widget showing current frame
func (apd *AnimationPlayerDialog) createPreviewWidget() fyne.CanvasObject {
	var img image.Image = image.NewRGBA(image.Rect(0, 0, 1, 1))
	if len(apd.frameImages) > 0 {
		img = apd.frameImages[0]
	}

	apd.previewImage = canvas.NewImageFromImage(img)
	apd.previewImage.FillMode = canvas.ImageFillContain
	apd.previewImage.SetMinSize(fyne.NewSize(128, 128))

	return container.NewVBox(
		widget.NewLabel("Animation Preview"),
		apd.previewImage,
	)
}

// Show displays the dialog
func (apd *AnimationPlayerDialog) Show(parent fyne.Window) {
	dialog := container.NewVBox(
		apd.container,
		container.NewHBox(
			widget.NewButton("Close", func() {
				if apd.onClose != nil {
					apd.onClose()
				}
				apd.Hide()
			}),
		),
	)

	apd.dialog = widget.NewPopUp(dialog, parent.Canvas())
	apd.dialog.ShowAtPosition(fyne.NewPos(150, 150))

	if apd.isAutoPlaying {
		apd.onPlay()
	}
}

// Hide hides the dialog
func (apd *AnimationPlayerDialog) Hide() {
	if apd.dialog != nil {
		apd.dialog.Hide()
	}
	apd.onStop()
}

// Tick updates animation playback
func (apd *AnimationPlayerDialog) Tick(deltaTimeMs int64) {
	if !apd.controller.IsPlaying() {
		return
	}

	frameIdx, changed := apd.controller.Update(deltaTimeMs)

	if changed {
		apd.frameSlider.Value = float64(frameIdx)
		apd.frameLabel.SetText(fmt.Sprintf("Frame: %d / %d", frameIdx, len(apd.animation.Frames)))
		apd.frameLabel.Refresh()
		apd.updatePreviewFrame(frameIdx)
	}

	apd.timeLabel.SetText(fmt.Sprintf("Time: %dms", apd.controller.GetCurrentTime()))
	apd.timeLabel.Refresh()
}

// Event handlers

func (apd *AnimationPlayerDialog) onPlay() {
	apd.controller.Play()
	apd.playButton.Disable()
	apd.pauseButton.Enable()
}

func (apd *AnimationPlayerDialog) onPause() {
	apd.controller.Pause()
	apd.playButton.Enable()
	apd.pauseButton.Disable()
}

func (apd *AnimationPlayerDialog) onStop() {
	apd.controller.Stop()
	apd.playButton.Enable()
	apd.pauseButton.Disable()

	apd.frameSlider.Value = 0
	apd.frameLabel.SetText(fmt.Sprintf("Frame: 0 / %d", len(apd.animation.Frames)))
	apd.frameLabel.Refresh()
	apd.timeLabel.SetText("Time: 0ms")
	apd.timeLabel.Refresh()
	apd.updatePreviewFrame(0)
}

func (apd *AnimationPlayerDialog) onFrameSliderChanged(value float64) {
	frameIdx := int(value)
	if err := apd.controller.JumpToFrame(frameIdx); err != nil {
		return
	}

	apd.frameLabel.SetText(fmt.Sprintf("Frame: %d / %d", frameIdx, len(apd.animation.Frames)))
	apd.frameLabel.Refresh()
	apd.timeLabel.SetText(fmt.Sprintf("Time: %dms", apd.controller.GetCurrentTime()))
	apd.timeLabel.Refresh()
	apd.updatePreviewFrame(frameIdx)
}

func (apd *AnimationPlayerDialog) onSpeedChanged(value float64) {
	apd.controller.SetPlaybackSpeed(float32(value))
	apd.animation.DefaultSpeed = float32(value)
	apd.speedLabel.SetText(fmt.Sprintf("Speed: %.2fx", value))
	apd.speedLabel.Refresh()
}

func (apd *AnimationPlayerDialog) onLoopModeChanged(value string) {
	var mode models.PlaybackMode
	switch value {
	case "Once":
		mode = models.PlayOnce
	case "PingPong":
		mode = models.PlayPingPong
	default:
		mode = models.PlayContinuous
	}
	apd.animation.PlaybackMode = mode
	apd.controller.SetLoopMode(mode)
}

func (apd *AnimationPlayerDialog) onAutoPlayToggled(value bool) {
	apd.isAutoPlaying = value
}

// Helper methods

func (apd *AnimationPlayerDialog) updatePreviewFrame(idx int) {
	if apd.previewImage == nil || len(apd.frameImages) == 0 {
		return
	}
	if idx < 0 || idx >= len(apd.frameImages) {
		return
	}

	apd.previewImage.Image = apd.frameImages[idx]
	apd.previewImage.Refresh()
}

func decodeFrameImages(animation *models.AnimationData) []image.Image {
	if animation == nil || len(animation.Frames) == 0 {
		return nil
	}

	conv := io.NewFF6SpriteConverter()
	frames := make([]image.Image, 0, len(animation.Frames))
	for _, frame := range animation.Frames {
		img, err := conv.DecodeFF6Sprite(frame)
		if err != nil {
			continue
		}
		frames = append(frames, img)
	}
	return frames
}

func (apd *AnimationPlayerDialog) calculateFPS() float32 {
	if len(apd.animation.Frames) == 0 {
		return 0
	}

	total := int64(0)
	for _, timing := range apd.animation.FrameTimings {
		total += int64(timing)
	}

	if total <= 0 {
		return 0
	}

	return float32(len(apd.animation.Frames)) * 1000.0 / float32(total)
}

// GetController returns the animation controller
func (apd *AnimationPlayerDialog) GetController() *io.AnimationController {
	return apd.controller
}

// SetOnClose sets the close callback
func (apd *AnimationPlayerDialog) SetOnClose(callback func()) {
	apd.onClose = callback
}
