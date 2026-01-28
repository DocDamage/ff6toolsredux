package models

import "time"

// AnimationData represents a complete sprite animation
type AnimationData struct {
	Frames       []*FF6Sprite      // Animation frames
	FrameTimings []int             // Duration per frame in milliseconds
	Metadata     AnimationMetadata // Animation metadata
	PlaybackMode PlaybackMode      // How to play (once, loop, pingpong)
	DefaultSpeed float32           // 1.0 = 60 FPS (16.67ms per frame)
}

// AnimationMetadata stores animation information
type AnimationMetadata struct {
	Name          string    // Animation name
	Description   string    // Animation description
	Author        string    // Creator name
	Created       time.Time // Creation timestamp
	Modified      time.Time // Last modified timestamp
	TotalDuration int64     // Total animation duration in milliseconds
	FrameCount    int       // Number of frames
}

// PlaybackMode defines animation playback behavior
type PlaybackMode int

const (
	// PlayOnce plays animation once then stops
	PlayOnce PlaybackMode = iota
	// PlayContinuous loops animation continuously
	PlayContinuous
	// PlayPingPong plays forward then backward repeatedly
	PlayPingPong
)

// String returns the string representation of PlaybackMode
func (p PlaybackMode) String() string {
	switch p {
	case PlayOnce:
		return "Once"
	case PlayContinuous:
		return "Loop"
	case PlayPingPong:
		return "PingPong"
	default:
		return "Unknown"
	}
}

// PlaybackState represents current animation playback state
type PlaybackState int

const (
	// StateStopped animation is stopped
	StateStopped PlaybackState = iota
	// StatePlaying animation is playing
	StatePlaying
	// StatePaused animation is paused
	StatePaused
)

// String returns the string representation of PlaybackState
func (s PlaybackState) String() string {
	switch s {
	case StateStopped:
		return "Stopped"
	case StatePlaying:
		return "Playing"
	case StatePaused:
		return "Paused"
	default:
		return "Unknown"
	}
}

// FrameInfo contains per-frame animation data
type FrameInfo struct {
	SpriteFrame *FF6Sprite // The sprite frame
	DurationMs  int        // How long to display this frame
	Index       int        // Frame index in animation
	DisplayTime int64      // Time offset where this frame displays (ms)
}

// NewAnimationData creates a new animation from sprite frames
func NewAnimationData(frames []*FF6Sprite) *AnimationData {
	if len(frames) == 0 {
		return nil
	}

	// Default timing: 100ms per frame (10 FPS)
	timings := make([]int, len(frames))
	totalDuration := int64(0)
	for i := range timings {
		timings[i] = 100
		totalDuration += 100
	}

	return &AnimationData{
		Frames:       frames,
		FrameTimings: timings,
		Metadata: AnimationMetadata{
			Name:          "Animation",
			Description:   "",
			Author:        "",
			Created:       time.Now(),
			Modified:      time.Now(),
			TotalDuration: totalDuration,
			FrameCount:    len(frames),
		},
		PlaybackMode: PlayContinuous,
		DefaultSpeed: 1.0,
	}
}

// SetFrameTiming sets the duration for a specific frame
func (a *AnimationData) SetFrameTiming(frameIdx int, durationMs int) error {
	if frameIdx < 0 || frameIdx >= len(a.FrameTimings) {
		return ErrInvalidFrameIndex
	}

	if durationMs <= 0 {
		return ErrInvalidDuration
	}

	oldDuration := a.FrameTimings[frameIdx]
	a.FrameTimings[frameIdx] = durationMs

	// Update total duration
	a.Metadata.TotalDuration += int64(durationMs - oldDuration)
	a.Metadata.Modified = time.Now()

	return nil
}

// GetFrameAtTime returns the frame index at a given time in the animation
func (a *AnimationData) GetFrameAtTime(timeMs int64) int {
	if len(a.Frames) == 0 {
		return 0
	}

	// Calculate total animation duration
	totalDuration := int64(0)
	for _, duration := range a.FrameTimings {
		totalDuration += int64(duration)
	}

	// Handle based on playback mode
	adjustedTime := timeMs
	switch a.PlaybackMode {
	case PlayOnce:
		if adjustedTime >= totalDuration {
			return len(a.Frames) - 1
		}

	case PlayContinuous:
		adjustedTime = timeMs % totalDuration
		if adjustedTime < 0 {
			adjustedTime += totalDuration
		}

	case PlayPingPong:
		fullCycleDuration := totalDuration * 2
		cycleTime := timeMs % fullCycleDuration
		if cycleTime < 0 {
			cycleTime += fullCycleDuration
		}

		if cycleTime >= totalDuration {
			// Reverse direction
			adjustedTime = totalDuration*2 - cycleTime - 1
		} else {
			adjustedTime = cycleTime
		}
	}

	// Find frame at adjusted time
	currentTime := int64(0)
	for i, duration := range a.FrameTimings {
		if currentTime+int64(duration) > adjustedTime {
			return i
		}
		currentTime += int64(duration)
	}

	return len(a.Frames) - 1
}

// GetTotalDuration returns the total animation duration in milliseconds
func (a *AnimationData) GetTotalDuration() int64 {
	if len(a.FrameTimings) == 0 {
		return 0
	}

	total := int64(0)
	for _, duration := range a.FrameTimings {
		total += int64(duration)
	}
	return total
}

// SetAllFrameDurations sets the same duration for all frames
func (a *AnimationData) SetAllFrameDurations(durationMs int) error {
	if durationMs <= 0 {
		return ErrInvalidDuration
	}

	for i := range a.FrameTimings {
		a.FrameTimings[i] = durationMs
	}

	a.Metadata.TotalDuration = int64(durationMs) * int64(len(a.Frames))
	a.Metadata.Modified = time.Now()

	return nil
}

// Validate checks animation data for consistency
func (a *AnimationData) Validate() error {
	if len(a.Frames) == 0 {
		return ErrNoFrames
	}

	if len(a.Frames) != len(a.FrameTimings) {
		return ErrMismatchedFramesAndTimings
	}

	for i, duration := range a.FrameTimings {
		if duration <= 0 {
			return ErrInvalidFrameTiming
		}

		if a.Frames[i] == nil {
			return ErrNilFrame
		}
	}

	return nil
}

// Copy creates a deep copy of the animation data
func (a *AnimationData) Copy() *AnimationData {
	copied := &AnimationData{
		Frames:       make([]*FF6Sprite, len(a.Frames)),
		FrameTimings: make([]int, len(a.FrameTimings)),
		Metadata:     a.Metadata,
		PlaybackMode: a.PlaybackMode,
		DefaultSpeed: a.DefaultSpeed,
	}

	copy(copied.Frames, a.Frames)
	copy(copied.FrameTimings, a.FrameTimings)

	return copied
}

// ExportFrameInfo returns detailed frame information for export
func (a *AnimationData) ExportFrameInfo() []FrameInfo {
	info := make([]FrameInfo, len(a.Frames))

	displayTime := int64(0)
	for i, frame := range a.Frames {
		info[i] = FrameInfo{
			SpriteFrame: frame,
			DurationMs:  a.FrameTimings[i],
			Index:       i,
			DisplayTime: displayTime,
		}
		displayTime += int64(a.FrameTimings[i])
	}

	return info
}

// animError is a simple error type for animation operations
type animError struct {
	message string
}

// Error implements the error interface
func (e *animError) Error() string {
	return e.message
}

// AnimationError types
var (
	ErrInvalidFrameIndex          = &animError{"invalid frame index"}
	ErrInvalidDuration            = &animError{"invalid duration (must be > 0)"}
	ErrNoFrames                   = &animError{"animation has no frames"}
	ErrMismatchedFramesAndTimings = &animError{"frame count does not match timing count"}
	ErrInvalidFrameTiming         = &animError{"invalid frame timing"}
	ErrNilFrame                   = &animError{"frame is nil"}
)

// NewError creates a new animation error
func NewError(msg string) error {
	return &animError{message: msg}
}
