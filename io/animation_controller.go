package io

import (
	"fmt"
	"time"

	"ffvi_editor/models"
)

// AnimationController manages animation playback with precise timing
type AnimationController struct {
	animation             *models.AnimationData
	currentFrameIdx       int
	currentTimeMs         int64
	isPlaying             bool
	isPaused              bool
	playbackSpeed         float32
	lastUpdateTimeNs      int64
	directionForward      bool
	onFrameChange         func(int)
	onPlaybackStateChange func(bool)
}

// NewAnimationController creates a new animation playback controller
func NewAnimationController(animation *models.AnimationData) (*AnimationController, error) {
	if animation == nil {
		return nil, fmt.Errorf("animation is nil")
	}

	if err := animation.Validate(); err != nil {
		return nil, err
	}

	return &AnimationController{
		animation:        animation,
		currentFrameIdx:  0,
		currentTimeMs:    0,
		isPlaying:        false,
		isPaused:         false,
		playbackSpeed:    1.0,
		directionForward: true,
	}, nil
}

// Play starts or resumes animation playback
func (ac *AnimationController) Play() error {
	if ac.animation == nil {
		return fmt.Errorf("no animation loaded")
	}

	ac.isPlaying = true
	ac.isPaused = false
	ac.lastUpdateTimeNs = time.Now().UnixNano()

	if ac.onPlaybackStateChange != nil {
		ac.onPlaybackStateChange(true)
	}

	return nil
}

// Pause pauses animation playback
func (ac *AnimationController) Pause() {
	if ac.isPlaying {
		ac.isPlaying = false
		ac.isPaused = true

		if ac.onPlaybackStateChange != nil {
			ac.onPlaybackStateChange(false)
		}
	}
}

// Stop stops animation and resets to first frame
func (ac *AnimationController) Stop() {
	ac.isPlaying = false
	ac.isPaused = false
	ac.currentFrameIdx = 0
	ac.currentTimeMs = 0
	ac.directionForward = true

	if ac.onPlaybackStateChange != nil {
		ac.onPlaybackStateChange(false)
	}

	if ac.onFrameChange != nil {
		ac.onFrameChange(0)
	}
}

// Resume resumes from pause
func (ac *AnimationController) Resume() error {
	if ac.isPaused {
		return ac.Play()
	}
	return nil
}

// Update updates animation playback state
func (ac *AnimationController) Update(deltaTimeMs int64) (int, bool) {
	if !ac.isPlaying || ac.animation == nil || len(ac.animation.Frames) == 0 {
		return ac.currentFrameIdx, false
	}

	// Apply speed factor
	scaledDelta := int64(float32(deltaTimeMs) * ac.playbackSpeed)
	ac.currentTimeMs += scaledDelta

	// Get total duration
	totalDuration := ac.getTotalDuration()
	if totalDuration <= 0 {
		return 0, false
	}

	// Calculate frame based on playback mode
	frameChanged := false
	oldFrame := ac.currentFrameIdx

	switch ac.animation.PlaybackMode {
	case models.PlayOnce:
		if ac.currentTimeMs >= totalDuration {
			ac.currentTimeMs = totalDuration
			ac.currentFrameIdx = len(ac.animation.Frames) - 1
			ac.isPlaying = false
			if ac.onPlaybackStateChange != nil {
				ac.onPlaybackStateChange(false)
			}
		} else {
			ac.currentFrameIdx = ac.animation.GetFrameAtTime(ac.currentTimeMs)
		}

	case models.PlayContinuous:
		ac.currentTimeMs = ac.currentTimeMs % totalDuration
		ac.currentFrameIdx = ac.animation.GetFrameAtTime(ac.currentTimeMs)

	case models.PlayPingPong:
		fullCycleDuration := totalDuration * 2
		cycleTime := ac.currentTimeMs % fullCycleDuration
		if cycleTime < totalDuration {
			ac.directionForward = true
			ac.currentFrameIdx = ac.animation.GetFrameAtTime(cycleTime)
		} else {
			ac.directionForward = false
			adjustedTime := totalDuration*2 - cycleTime - 1
			ac.currentFrameIdx = ac.animation.GetFrameAtTime(adjustedTime)
		}
	}

	if ac.currentFrameIdx != oldFrame {
		frameChanged = true
		if ac.onFrameChange != nil {
			ac.onFrameChange(ac.currentFrameIdx)
		}
	}

	return ac.currentFrameIdx, frameChanged
}

// SetPlaybackSpeed sets animation playback speed (0.5x to 2.0x)
func (ac *AnimationController) SetPlaybackSpeed(speed float32) error {
	if speed < 0.5 || speed > 2.0 {
		return fmt.Errorf("playback speed must be between 0.5 and 2.0, got %.2f", speed)
	}
	ac.playbackSpeed = speed
	return nil
}

// SetLoopMode sets animation loop mode
func (ac *AnimationController) SetLoopMode(mode models.PlaybackMode) {
	ac.animation.PlaybackMode = mode
}

// JumpToFrame jumps to specific frame index
func (ac *AnimationController) JumpToFrame(frameIdx int) error {
	if frameIdx < 0 || frameIdx >= len(ac.animation.Frames) {
		return fmt.Errorf("invalid frame index: %d", frameIdx)
	}

	ac.currentFrameIdx = frameIdx
	ac.currentTimeMs = ac.getTimeAtFrame(frameIdx)

	if ac.onFrameChange != nil {
		ac.onFrameChange(frameIdx)
	}

	return nil
}

// JumpToTime jumps to specific time in animation
func (ac *AnimationController) JumpToTime(timeMs int64) error {
	if timeMs < 0 {
		return fmt.Errorf("time cannot be negative")
	}

	totalDuration := ac.getTotalDuration()
	if timeMs > totalDuration {
		timeMs = totalDuration
	}

	ac.currentTimeMs = timeMs
	ac.currentFrameIdx = ac.animation.GetFrameAtTime(timeMs)

	if ac.onFrameChange != nil {
		ac.onFrameChange(ac.currentFrameIdx)
	}

	return nil
}

// GetCurrentFrame returns current frame index
func (ac *AnimationController) GetCurrentFrame() int {
	return ac.currentFrameIdx
}

// GetCurrentSprite returns current sprite frame
func (ac *AnimationController) GetCurrentSprite() *models.FF6Sprite {
	if ac.animation == nil || ac.currentFrameIdx < 0 || ac.currentFrameIdx >= len(ac.animation.Frames) {
		return nil
	}
	return ac.animation.Frames[ac.currentFrameIdx]
}

// GetCurrentTime returns current time in milliseconds
func (ac *AnimationController) GetCurrentTime() int64 {
	return ac.currentTimeMs
}

// CalculateFPS calculates current frames per second based on timings
func (ac *AnimationController) CalculateFPS() float32 {
	if ac.animation == nil || len(ac.animation.Frames) == 0 {
		return 0
	}

	totalDuration := ac.getTotalDuration()
	if totalDuration <= 0 {
		return 0
	}

	return float32(len(ac.animation.Frames)) * 1000.0 / float32(totalDuration)
}

// OnFrameChange sets the frame change callback
func (ac *AnimationController) OnFrameChange(callback func(int)) {
	ac.onFrameChange = callback
}

// OnPlaybackStateChange sets the playback state change callback
func (ac *AnimationController) OnPlaybackStateChange(callback func(bool)) {
	ac.onPlaybackStateChange = callback
}

// IsPlaying returns true if animation is currently playing
func (ac *AnimationController) IsPlaying() bool {
	return ac.isPlaying
}

// IsPaused returns true if animation is paused
func (ac *AnimationController) IsPaused() bool {
	return ac.isPaused
}

// Helper methods

func (ac *AnimationController) getTotalDuration() int64 {
	if ac.animation == nil || len(ac.animation.FrameTimings) == 0 {
		return 0
	}

	total := int64(0)
	for _, timing := range ac.animation.FrameTimings {
		total += int64(timing)
	}
	return total
}

func (ac *AnimationController) getTimeAtFrame(frameIdx int) int64 {
	if ac.animation == nil || frameIdx < 0 || frameIdx >= len(ac.animation.Frames) {
		return 0
	}

	time := int64(0)
	for i := 0; i < frameIdx && i < len(ac.animation.FrameTimings); i++ {
		time += int64(ac.animation.FrameTimings[i])
	}
	return time
}
