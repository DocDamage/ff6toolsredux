package io

import (
	"fmt"
	"time"

	"ffvi_editor/models"
)

// FrameSequencer manages frame timing and sequencing
type FrameSequencer struct {
	frames   []*models.FF6Sprite
	timings  []int
	metadata SequenceMetadata
}

// SequenceMetadata stores sequence information
type SequenceMetadata struct {
	Name          string
	Description   string
	DefaultSpeed  float32
	TotalFrames   int
	TotalDuration int64
	Author        string
	Created       time.Time
	Modified      time.Time
}

// NewFrameSequencer creates a new frame sequencer
func NewFrameSequencer(frames []*models.FF6Sprite) *FrameSequencer {
	if len(frames) == 0 {
		return nil
	}

	// Initialize default timings (100ms per frame = 10 FPS)
	timings := make([]int, len(frames))
	totalDuration := int64(0)
	for i := range timings {
		timings[i] = 100
		totalDuration += 100
	}

	return &FrameSequencer{
		frames:  frames,
		timings: timings,
		metadata: SequenceMetadata{
			Name:          "Sequence",
			DefaultSpeed:  1.0,
			TotalFrames:   len(frames),
			TotalDuration: totalDuration,
			Created:       time.Now(),
			Modified:      time.Now(),
		},
	}
}

// SetFrameDuration sets duration for specific frame
func (fs *FrameSequencer) SetFrameDuration(frameIdx int, durationMs int) error {
	if frameIdx < 0 || frameIdx >= len(fs.timings) {
		return fmt.Errorf("invalid frame index: %d", frameIdx)
	}

	if durationMs <= 0 {
		return fmt.Errorf("duration must be positive, got %d", durationMs)
	}

	oldDuration := fs.timings[frameIdx]
	fs.timings[frameIdx] = durationMs
	fs.metadata.TotalDuration += int64(durationMs - oldDuration)
	fs.metadata.Modified = time.Now()

	return nil
}

// GetFrameDuration returns duration for specific frame
func (fs *FrameSequencer) GetFrameDuration(frameIdx int) int {
	if frameIdx < 0 || frameIdx >= len(fs.timings) {
		return 0
	}
	return fs.timings[frameIdx]
}

// SetAllFrameDurations sets same duration for all frames
func (fs *FrameSequencer) SetAllFrameDurations(durationMs int) error {
	if durationMs <= 0 {
		return fmt.Errorf("duration must be positive, got %d", durationMs)
	}

	for i := range fs.timings {
		fs.timings[i] = durationMs
	}

	fs.metadata.TotalDuration = int64(durationMs) * int64(len(fs.frames))
	fs.metadata.Modified = time.Now()

	return nil
}

// AutoTiming automatically sets frame durations based on base speed
func (fs *FrameSequencer) AutoTiming(baseSpeedMs int) error {
	if baseSpeedMs <= 0 {
		return fmt.Errorf("base speed must be positive, got %d", baseSpeedMs)
	}

	// Set all frames to base speed with slight variation
	for i := range fs.timings {
		variation := (i % 3) - 1 // -1, 0, or 1
		duration := baseSpeedMs + (variation * baseSpeedMs / 10)
		if duration < 10 {
			duration = 10
		}
		fs.timings[i] = duration
	}

	// Recalculate total
	total := int64(0)
	for _, timing := range fs.timings {
		total += int64(timing)
	}
	fs.metadata.TotalDuration = total
	fs.metadata.Modified = time.Now()

	return nil
}

// GetTotalDuration returns total animation duration
func (fs *FrameSequencer) GetTotalDuration() int64 {
	if len(fs.timings) == 0 {
		return 0
	}

	total := int64(0)
	for _, timing := range fs.timings {
		total += int64(timing)
	}
	return total
}

// GetFrameAtTime returns frame index at given time
func (fs *FrameSequencer) GetFrameAtTime(timeMs int64) int {
	if len(fs.frames) == 0 {
		return 0
	}

	currentTime := int64(0)
	for i, duration := range fs.timings {
		if currentTime+int64(duration) > timeMs {
			return i
		}
		currentTime += int64(duration)
	}

	return len(fs.frames) - 1
}

// InsertFrame inserts a frame at position
func (fs *FrameSequencer) InsertFrame(position int, frame *models.FF6Sprite) error {
	if position < 0 || position > len(fs.frames) {
		return fmt.Errorf("invalid position: %d", position)
	}

	if frame == nil {
		return fmt.Errorf("frame is nil")
	}

	// Insert frame
	newFrames := make([]*models.FF6Sprite, len(fs.frames)+1)
	copy(newFrames[:position], fs.frames[:position])
	newFrames[position] = frame
	copy(newFrames[position+1:], fs.frames[position:])
	fs.frames = newFrames

	// Insert timing
	newTimings := make([]int, len(fs.timings)+1)
	copy(newTimings[:position], fs.timings[:position])
	newTimings[position] = 100
	copy(newTimings[position+1:], fs.timings[position:])
	fs.timings = newTimings

	fs.metadata.TotalFrames = len(fs.frames)
	fs.metadata.TotalDuration += 100
	fs.metadata.Modified = time.Now()

	return nil
}

// RemoveFrame removes frame at index
func (fs *FrameSequencer) RemoveFrame(frameIdx int) error {
	if frameIdx < 0 || frameIdx >= len(fs.frames) {
		return fmt.Errorf("invalid frame index: %d", frameIdx)
	}

	if len(fs.frames) <= 1 {
		return fmt.Errorf("cannot remove last frame")
	}

	removedDuration := fs.timings[frameIdx]
	fs.frames = append(fs.frames[:frameIdx], fs.frames[frameIdx+1:]...)
	fs.timings = append(fs.timings[:frameIdx], fs.timings[frameIdx+1:]...)

	fs.metadata.TotalFrames = len(fs.frames)
	fs.metadata.TotalDuration -= int64(removedDuration)
	fs.metadata.Modified = time.Now()

	return nil
}

// DuplicateFrame duplicates frame at index
func (fs *FrameSequencer) DuplicateFrame(frameIdx int) error {
	if frameIdx < 0 || frameIdx >= len(fs.frames) {
		return fmt.Errorf("invalid frame index: %d", frameIdx)
	}

	return fs.InsertFrame(frameIdx+1, fs.frames[frameIdx])
}

// MoveFrame moves frame from one position to another
func (fs *FrameSequencer) MoveFrame(fromIdx int, toIdx int) error {
	if fromIdx < 0 || fromIdx >= len(fs.frames) {
		return fmt.Errorf("invalid from index: %d", fromIdx)
	}

	if toIdx < 0 || toIdx >= len(fs.frames) {
		return fmt.Errorf("invalid to index: %d", toIdx)
	}

	if fromIdx == toIdx {
		return nil
	}

	// Save and remove
	frame := fs.frames[fromIdx]
	timing := fs.timings[fromIdx]

	fs.frames = append(fs.frames[:fromIdx], fs.frames[fromIdx+1:]...)
	fs.timings = append(fs.timings[:fromIdx], fs.timings[fromIdx+1:]...)

	// Insert at new position
	newFrames := make([]*models.FF6Sprite, len(fs.frames)+1)
	newTimings := make([]int, len(fs.timings)+1)

	copy(newFrames[:toIdx], fs.frames[:toIdx])
	newFrames[toIdx] = frame
	copy(newFrames[toIdx+1:], fs.frames[toIdx:])

	copy(newTimings[:toIdx], fs.timings[:toIdx])
	newTimings[toIdx] = timing
	copy(newTimings[toIdx+1:], fs.timings[toIdx:])

	fs.frames = newFrames
	fs.timings = newTimings
	fs.metadata.Modified = time.Now()

	return nil
}

// Validate checks sequence for consistency
func (fs *FrameSequencer) Validate() error {
	if len(fs.frames) == 0 {
		return fmt.Errorf("sequence has no frames")
	}

	if len(fs.frames) != len(fs.timings) {
		return fmt.Errorf("frame count (%d) does not match timing count (%d)",
			len(fs.frames), len(fs.timings))
	}

	for i, frame := range fs.frames {
		if frame == nil {
			return fmt.Errorf("frame %d is nil", i)
		}
	}

	for i, timing := range fs.timings {
		if timing <= 0 {
			return fmt.Errorf("frame %d has invalid timing: %d", i, timing)
		}
	}

	return nil
}

// GetFrames returns copy of frames array
func (fs *FrameSequencer) GetFrames() []*models.FF6Sprite {
	frames := make([]*models.FF6Sprite, len(fs.frames))
	copy(frames, fs.frames)
	return frames
}

// GetTimings returns copy of timings array
func (fs *FrameSequencer) GetTimings() []int {
	timings := make([]int, len(fs.timings))
	copy(timings, fs.timings)
	return timings
}

// GetMetadata returns sequence metadata
func (fs *FrameSequencer) GetMetadata() SequenceMetadata {
	return fs.metadata
}

// SetMetadata sets sequence metadata
func (fs *FrameSequencer) SetMetadata(meta SequenceMetadata) {
	fs.metadata = meta
	fs.metadata.Modified = time.Now()
}

// GetFrameCount returns number of frames
func (fs *FrameSequencer) GetFrameCount() int {
	return len(fs.frames)
}

// CalculateAverageFPS calculates average frames per second
func (fs *FrameSequencer) CalculateAverageFPS() float32 {
	if len(fs.frames) == 0 {
		return 0
	}

	totalDuration := fs.GetTotalDuration()
	if totalDuration == 0 {
		return 0
	}

	return float32(len(fs.frames)) * 1000.0 / float32(totalDuration)
}

// ExportToAnimation creates AnimationData from sequence
func (fs *FrameSequencer) ExportToAnimation(playbackMode models.PlaybackMode) *models.AnimationData {
	anim := models.NewAnimationData(fs.frames)
	if anim == nil {
		return nil
	}

	copy(anim.FrameTimings, fs.timings)
	anim.PlaybackMode = playbackMode
	anim.Metadata = models.AnimationMetadata{
		Name:          fs.metadata.Name,
		Description:   fs.metadata.Description,
		Author:        fs.metadata.Author,
		Created:       fs.metadata.Created,
		Modified:      fs.metadata.Modified,
		TotalDuration: fs.metadata.TotalDuration,
		FrameCount:    len(fs.frames),
	}

	return anim
}
