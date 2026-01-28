# 🎬 PHASE 13 - WEEK 3 ANIMATION SYSTEM PLAN

**Status:** Planning Phase  
**Duration:** Week 3 (January 20-24, 2026)  
**Focus:** Animation management, playback, and export

---

## 📋 Overview

Build complete animation system on top of Week 1-2 sprite foundation:
- Frame management and sequencing
- Playback controls with timing
- Animation export formats
- Professional animation editor UI

---

## 🏗️ Architecture Design

### Component Hierarchy

```
┌─────────────────────────────────────────┐
│ Animation Menu (ui/window.go)           │
├─────────────────────────────────────────┤
│ Animation Player Dialog (UI)            │
│ - Playback controls                     │
│ - Frame display                         │
│ - Speed/loop controls                   │
├─────────────────────────────────────────┤
│ Frame Editor Dialog (UI)                │
│ - Duration editor                       │
│ - Frame sequencing                      │
│ - Frame preview grid                    │
├─────────────────────────────────────────┤
│ Animation Backend (io/)                 │
│ - AnimationController                   │
│ - FrameSequencer                        │
│ - AnimationExporter                     │
├─────────────────────────────────────────┤
│ Animation Models (models/)              │
│ - AnimationData struct                  │
│ - FrameInfo struct                      │
│ - PlaybackState enum                    │
└─────────────────────────────────────────┘
```

### Data Flow

```
Sprite with Frames
    ↓
Animation Backend (Frame management)
    ↓
Playback Controller (Timing logic)
    ↓
Animation Player UI (Display + Controls)
    ↓
Frame Editor UI (Sequencing)
    ↓
Animation Exporter (Format output)
```

---

## 📦 Backend Components (io/)

### 1. Animation Controller (`io/animation_controller.go`)

**Purpose:** Core animation playback engine

```go
type AnimationController struct {
    frames              []*models.FF6Sprite
    frameTimings        []int                    // ms per frame
    currentFrameIdx     int
    currentTime         int64                    // milliseconds
    isPlaying           bool
    playbackSpeed       float32                  // 0.5x to 2.0x
    loopMode            LoopMode                // Once, Loop, PingPong
    onFrameChange       func(int)
    onPlaybackStateChange func(bool)
}

type LoopMode int
const (
    LoopOnce    LoopMode = iota
    LoopContinuous
    LoopPingPong
)

// Methods
func NewAnimationController(frames []*models.FF6Sprite) *AnimationController
func (a *AnimationController) Start() error
func (a *AnimationController) Stop()
func (a *AnimationController) Pause()
func (a *AnimationController) Resume()
func (a *AnimationController) Update(deltaTimeMs int64) (currentFrameIdx int, changed bool)
func (a *AnimationController) SetFrameTiming(frameIdx, durationMs int) error
func (a *AnimationController) SetPlaybackSpeed(speed float32)
func (a *AnimationController) SetLoopMode(mode LoopMode)
func (a *AnimationController) GetFrameTimings() []int
func (a *AnimationController) GetTotalDuration() int64  // milliseconds
func (a *AnimationController) JumpToFrame(frameIdx int) error
```

**Key Features:**
- Precise frame timing (millisecond accuracy)
- Variable playback speed (0.5x to 2.0x)
- 3 loop modes (Once, Continuous, PingPong)
- Event callbacks for frame changes
- Non-blocking update loop

---

### 2. Frame Sequencer (`io/frame_sequencer.go`)

**Purpose:** Manage frame timing and sequencing

```go
type FrameSequencer struct {
    frames           []*models.FF6Sprite
    timings          []int              // Duration per frame in ms
    metadata         *SequenceMetadata
}

type SequenceMetadata struct {
    Name            string
    Description     string
    DefaultSpeed    float32            // 1.0 = 60 FPS
    TotalFrames     int
    TotalDuration   int64              // milliseconds
    Author          string
    Created         time.Time
    Modified        time.Time
}

// Methods
func NewFrameSequencer(frames []*models.FF6Sprite) *FrameSequencer
func (f *FrameSequencer) SetFrameDuration(frameIdx, durationMs int) error
func (f *FrameSequencer) GetFrameDuration(frameIdx int) int
func (f *FrameSequencer) SetAllFrameDurations(durationMs int)
func (f *FrameSequencer) AutoTiming(baseSpeedMs int)                    // Smart defaults
func (f *FrameSequencer) GetTotalDuration() int64
func (f *FrameSequencer) GetFrameAtTime(timeMs int64) int              // Frame for time
func (f *FrameSequencer) InsertFrame(position int, frame *models.FF6Sprite)
func (f *FrameSequencer) RemoveFrame(frameIdx int) error
func (f *FrameSequencer) DuplicateFrame(frameIdx int)
func (f *FrameSequencer) Validate() error
```

**Key Features:**
- Per-frame duration control
- Automatic timing generation
- Frame insertion/removal
- Validation and error checking

---

### 3. Animation Exporter (`io/animation_exporter.go`)

**Purpose:** Export animations to various formats

```go
type AnimationExporter struct {
    sequencer *FrameSequencer
}

type ExportFormat int
const (
    FormatGIF ExportFormat = iota
    FormatWebP
    FormatMP4
    FormatJSON
)

// Methods
func NewAnimationExporter(seq *FrameSequencer) *AnimationExporter
func (e *AnimationExporter) ExportToGIF(filename string, opts *GIFOptions) error
func (e *AnimationExporter) ExportToJSON(filename string) error          // Sprite sheet
func (e *AnimationExporter) ExportToWebP(filename string, opts *WebPOptions) error
func (e *AnimationExporter) ExportFrames(outputDir string) error         // Individual PNGs

type GIFOptions struct {
    Quality     int         // 1-256 (default: 256)
    Looping     bool        // Default: true
    BackgroundColor color.Color
    PixelScale  int         // Upscale factor
}

type WebPOptions struct {
    Quality     float32     // 0.0-1.0 (default: 0.8)
    Looping     bool
    PixelScale  int
}
```

**Supported Exports:**
- GIF (animated, web-ready)
- WebP (modern, efficient)
- PNG frames (for video editors)
- JSON metadata (for sprite sheets)

---

## 🎨 UI Components (ui/forms/)

### 1. Animation Player Dialog (`ui/forms/animation_player_dialog.go`)

**Purpose:** Playback interface with controls

**Features:**
- Play/Pause/Stop buttons
- Frame slider
- Speed control (0.5x - 2.0x)
- Loop mode selector (Once, Loop, PingPong)
- Frame counter (e.g., "Frame 5/24")
- Animation preview pane
- FPS display
- Timeline scrubber

**UI Layout:**
```
┌────────────────────────────────┐
│ Animation Player               │
├────────────────────────────────┤
│ ┌──────────────────────────┐   │
│ │  Frame 5/24              │   │
│ │  ┌────────────────────┐  │   │
│ │  │  Sprite Preview    │  │   │
│ │  │                    │  │   │
│ │  └────────────────────┘  │   │
│ │  FPS: 60  Duration: 1.2s  │  │
│ └──────────────────────────┘   │
│ ▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░   │
│ [◀◀] [▶] [║] [▶▶]              │
│ Speed: [━━━●────]  1.0x        │
│ Loop: [Once ▼] ☑ Auto-advance  │
├────────────────────────────────┤
│ [Export...] [Edit Frames...] [Close]
└────────────────────────────────┘
```

**Key Methods:**
```go
type AnimationPlayerDialog struct {
    window              fyne.Window
    dialog              dialog.Dialog
    controller          *io.AnimationController
    playBtn             *widget.Button
    pauseBtn            *widget.Button
    stopBtn             *widget.Button
    frameSlider         *widget.Slider
    speedSlider         *widget.Slider
    loopSelect          *widget.Select
    frameLabel          *widget.Label
    spritePreview       *SpritePreviewWidget
    timeline            fyne.CanvasObject
    fpsLabel            *widget.Label
}

func NewAnimationPlayerDialog(window fyne.Window, sprite *models.FF6Sprite) *AnimationPlayerDialog
func (a *AnimationPlayerDialog) Play()
func (a *AnimationPlayerDialog) Pause()
func (a *AnimationPlayerDialog) Stop()
func (a *AnimationPlayerDialog) SetSpeed(speed float32)
func (a *AnimationPlayerDialog) SetLoopMode(mode string)
func (a *AnimationPlayerDialog) Show()
```

---

### 2. Frame Editor Dialog (`ui/forms/frame_editor_dialog.go`)

**Purpose:** Edit frame timings and sequences

**Features:**
- Frame grid with thumbnails
- Duration input per frame
- Drag-to-reorder frames
- Duplicate/delete frame buttons
- Auto-timing option
- Preview changes
- Batch duration editor

**UI Layout:**
```
┌──────────────────────────────────┐
│ Frame Editor                     │
├──────────────────────────────────┤
│ Frame 1      Duration: [100] ms  │
│ [←] [▲] [▼] [→] [+] [-] [Dup]   │
│                                  │
│ ┌──────────────────────────────┐ │
│ │ Frame Thumbnails:            │ │
│ │ ┌──┬──┬──┬──┬──┐            │ │
│ │ │1 │2 │3 │4 │5 │            │ │
│ │ ├──┼──┼──┼──┼──┤            │ │
│ │ │6 │7 │8 │9 │10│            │ │
│ │ └──┴──┴──┴──┴──┘            │ │
│ └──────────────────────────────┘ │
│                                  │
│ Batch Timing:                    │
│ Set all frames to: [100] ms      │
│ [Apply] [Auto-calc]              │
│                                  │
│ Total Duration: 1.2 seconds      │
├──────────────────────────────────┤
│ [Revert] [Close] [Apply]        │
└──────────────────────────────────┘
```

**Key Methods:**
```go
type FrameEditorDialog struct {
    window           fyne.Window
    dialog           dialog.Dialog
    sequencer        *io.FrameSequencer
    frameList        *fyne.Container
    selectedFrameIdx int
    durationEntry    *widget.Entry
    previewFrame     *SpritePreviewWidget
    timingTable      *container.Container
}

func NewFrameEditorDialog(window fyne.Window, frames []*models.FF6Sprite) *FrameEditorDialog
func (f *FrameEditorDialog) SelectFrame(idx int)
func (f *FrameEditorDialog) SetFrameDuration(frameIdx, ms int)
func (f *FrameEditorDialog) SetAllDurations(ms int)
func (f *FrameEditorDialog) AutoCalculateTiming()
func (f *FrameEditorDialog) Show()
```

---

### 3. Animation Export Dialog (`ui/forms/animation_export_dialog.go`)

**Purpose:** Export animations to various formats

**Features:**
- Format selector (GIF, WebP, PNG frames, JSON)
- Quality/compression options
- Preview before export
- Looping toggle
- Scale factor control
- Save location picker

**UI Layout:**
```
┌────────────────────────────────┐
│ Export Animation               │
├────────────────────────────────┤
│ Format: [GIF ▼]                │
│                                │
│ Export Options:                │
│ Quality: [━━━●────] 80%        │
│ Scale: [━━━○────] 2x           │
│ ☑ Looping                      │
│ ☑ Optimize                     │
│                                │
│ Output: [____________]         │
│         [Browse...]            │
│                                │
│ Preview: (first frame)         │
│ Estimated size: 245 KB         │
├────────────────────────────────┤
│ [Cancel] [Export]              │
└────────────────────────────────┘
```

**Formats Supported:**
- **GIF** - Animated, web-ready
- **WebP** - Modern, efficient
- **PNG** - Individual frames
- **JSON** - Metadata + sprite sheet reference

---

## 🔄 Integration Points

### 1. Menu Integration (ui/window.go)

```go
fyne.NewMenuItem("Animation Player...", func() {
    if sprite != nil && len(sprite.Frames) > 1 {
        dialog := forms.NewAnimationPlayerDialog(g.window, sprite)
        dialog.Show()
    } else {
        dialog.ShowError(fmt.Errorf("no multi-frame sprite"), g.window)
    }
}),

fyne.NewMenuItem("Edit Frames...", func() {
    if sprite != nil && len(sprite.Frames) > 0 {
        dialog := forms.NewFrameEditorDialog(g.window, sprite.Frames)
        dialog.Show()
    } else {
        dialog.ShowError(fmt.Errorf("no sprite frames"), g.window)
    }
}),

fyne.NewMenuItem("Export Animation...", func() {
    if sprite != nil && len(sprite.Frames) > 1 {
        dialog := forms.NewAnimationExportDialog(g.window, sprite)
        dialog.Show()
    } else {
        dialog.ShowError(fmt.Errorf("no animation to export"), g.window)
    }
}),
```

### 2. Sprite Editor Integration

From existing sprite import:
```go
// After successful import
if len(sprite.Frames) > 1 {
    // Offer to open animation player
    dialog.ShowInformation("Multi-frame sprite", 
        "This sprite has animation frames!", window)
    // Auto-open animation player
}
```

---

## 📊 Implementation Plan

### Phase 3.1: Backend (2-3 days)

**Day 1:**
- [ ] Create `models/animation.go` (AnimationData structures)
- [ ] Create `io/animation_controller.go` (Playback engine)
- [ ] Unit tests for playback timing

**Day 2:**
- [ ] Create `io/frame_sequencer.go` (Frame management)
- [ ] Create `io/animation_exporter.go` (GIF/WebP export)
- [ ] Integration tests

**Day 3:**
- [ ] Optimize performance
- [ ] Add error handling
- [ ] Documentation

### Phase 3.2: UI Layer (2-3 days)

**Day 1:**
- [ ] Create `animation_player_dialog.go`
- [ ] Implement playback controls
- [ ] Real-time preview

**Day 2:**
- [ ] Create `frame_editor_dialog.go`
- [ ] Implement frame timing UI
- [ ] Drag-to-reorder support

**Day 3:**
- [ ] Create `animation_export_dialog.go`
- [ ] Integrate menu items
- [ ] Polish UI/UX

---

## 📦 Dependencies

### New Models (models/animation.go)
```go
type AnimationData struct {
    Frames          []*FF6Sprite
    FrameTimings    []int              // ms per frame
    Metadata        AnimationMetadata
    PlaybackMode    PlaybackMode
}

type AnimationMetadata struct {
    Name            string
    Description     string
    Author          string
    Created         time.Time
    Modified        time.Time
    TotalDuration   int64
}

type PlaybackMode int
const (
    PlayOnce PlaybackMode = iota
    PlayLoop
    PlayPingPong
)
```

### External Dependencies
- **Gifimage** - Already in Fyne for GIF support (stdlib)
- **WebP** - May need library (or skip WebP for MVP)
- **Image encoding** - Go stdlib covers PNG/GIF

---

## ✅ Success Criteria

**Week 3 Completion (by Friday):**
- [ ] Animation playback controller working
- [ ] Frame sequencer functional
- [ ] Animation player UI displaying
- [ ] Frame editor working
- [ ] GIF export functional
- [ ] Menu items integrated
- [ ] 500+ lines UI code
- [ ] 400+ lines backend code
- [ ] All code compiles
- [ ] Comprehensive documentation

**Performance Targets:**
- Playback: 60 FPS smooth
- Export: <5 seconds for typical animation
- Memory: <50MB per animation

---

## 📚 Documentation Deliverables

1. **PHASE_13_WEEK3_ANIMATION_COMPLETE.md** - Implementation summary
2. **PHASE_13_ANIMATION_API.md** - Animation API reference
3. **PHASE_13_ANIMATION_USER_GUIDE.md** - How to use animation features

---

## 🎯 Quick Reference

**Key Classes:**
- `AnimationController` - Playback engine
- `FrameSequencer` - Frame management
- `AnimationExporter` - Export formats
- `AnimationPlayerDialog` - Player UI
- `FrameEditorDialog` - Editor UI

**Key Methods:**
- `controller.Update()` - Main update loop (call each frame)
- `sequencer.SetFrameDuration()` - Set timing
- `exporter.ExportToGIF()` - Export animation

**Menu Items:**
- Tools → Animation Player...
- Tools → Edit Frames...
- Tools → Export Animation...

---

## 🚀 Ready to Build?

All Week 1-2 components are ready to support animation system.

**Starting with:** Backend foundation (AnimationController + FrameSequencer)

Proceed? ✅
