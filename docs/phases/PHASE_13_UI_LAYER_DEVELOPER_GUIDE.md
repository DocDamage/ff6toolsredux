# 🔧 PHASE 13 UI LAYER - DEVELOPER INTEGRATION GUIDE

**Purpose:** Guide for developers adding new sprite features  
**Audience:** Backend and frontend developers  
**Updated:** January 17, 2026

---

## 📚 Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Component Reference](#component-reference)
3. [Adding UI Components](#adding-ui-components)
4. [Connecting to Backend](#connecting-to-backend)
5. [Common Patterns](#common-patterns)
6. [Testing](#testing)
7. [Performance Tips](#performance-tips)

---

## 🏗️ Architecture Overview

### Layer Structure
```
┌─────────────────────────────┐
│  Menu & Window Management   │  ui/window.go
├─────────────────────────────┤
│    Dialog Components        │  ui/forms/*.go
├─────────────────────────────┤
│   Fyne UI Framework (v2)    │
├─────────────────────────────┤
│   Backend Services          │  io/sprite_*.go
├─────────────────────────────┤
│   Data Models               │  models/sprite.go
└─────────────────────────────┘
```

### File Organization
```
ui/
├── window.go                          # Main window + menu
└── forms/
    ├── sprite_import_dialog.go        # Import UI
    ├── palette_editor_dialog.go       # Palette UI
    ├── sprite_preview_widget.go       # Preview widget
    ├── ... (other dialogs)
    └── ... (existing components)

models/
└── sprite.go                          # Sprite data types

io/
├── sprite_decoder.go                  # Image loading
├── sprite_converter.go                # FF6 format
├── sprite_validator.go                # Validation
├── sprite_importer.go                 # Import pipeline
└── palette_editor.go                  # Palette ops
```

---

## 📖 Component Reference

### Sprite Import Dialog

**Location:** `ui/forms/sprite_import_dialog.go`

**Constructor:**
```go
dialog := forms.NewSpriteImportDialog(window fyne.Window) *SpriteImportDialog
```

**Public Methods:**
```go
dialog.Show()                                    // Display
dialog.Hide()                                    // Hide
dialog.OnImportSuccess(func(*models.FF6Sprite)) // Success callback
dialog.OnImportError(func(error))               // Error callback
sprite := dialog.GetImportedSprite()            // Get result
```

**State:**
```go
selectedFile   string                  // Currently selected file
importedSprite *models.FF6Sprite      // Last import result
```

**Extension Points:**
- Override `performImport()` to customize workflow
- Add new dithering methods in `getDitheringMethod()`
- Modify preview in `updatePreview()`

---

### Palette Editor Dialog

**Location:** `ui/forms/palette_editor_dialog.go`

**Constructor:**
```go
dialog := forms.NewPaletteEditorDialog(
    window fyne.Window,
    palette *models.Palette,
) *PaletteEditorDialog
```

**Public Methods:**
```go
dialog.Show()                              // Display
dialog.Hide()                              // Hide
dialog.OnApply(func(*models.Palette))     // Apply callback
dialog.Revert()                            // Restore original
```

**State:**
```go
palette          *models.Palette   // Edited palette
originalPalette  *models.Palette   // Backup
selectedColorIdx int               // Current color
editor           *io.PaletteEditor // Backend
```

**Extension Points:**
- Add new harmony schemes in `generateHarmony()`
- Add new transformations in `applyTransform()`
- Customize color editor in `buildUI()`

---

### Sprite Preview Widget

**Location:** `ui/forms/sprite_preview_widget.go`

**Constructor:**
```go
widget := forms.NewSpritePreviewWidget(sprite *models.FF6Sprite) *SpritePreviewWidget
```

**Public Methods:**
```go
widget.SetSprite(sprite *models.FF6Sprite)  // Update sprite
widget.SetFrame(frameIdx int)               // Go to frame
frameIdx := widget.GetCurrentFrame() int    // Current frame
```

**State:**
```go
sprite       *models.FF6Sprite  // Current sprite
currentFrame int                // Current frame index
scale        float32            // Display scale
showGrid     bool               // Grid visibility
```

**Extension Points:**
- Customize tile grid in `drawGrid()`
- Add overlay information in `renderSprite()`
- Extend frame controls in `CreateRenderer()`

---

## 🎨 Adding UI Components

### Example: Add Animation Controls

**File:** `ui/forms/animation_controls_dialog.go`

```go
package forms

import (
    "ffvi_editor/models"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// AnimationControlsDialog provides animation playback UI
type AnimationControlsDialog struct {
    window         fyne.Window
    sprite         *models.FF6Sprite
    currentFrame   int
    isPlaying      bool
    frameDelay     int
    playBtn        *widget.Button
    stopBtn        *widget.Button
    delaySpinner   *widget.Entry
    onFrameChange  func(int)
}

// NewAnimationControlsDialog creates animation controls
func NewAnimationControlsDialog(
    window fyne.Window,
    sprite *models.FF6Sprite,
) *AnimationControlsDialog {
    a := &AnimationControlsDialog{
        window:       window,
        sprite:       sprite,
        currentFrame: 0,
        frameDelay:   100,
    }
    a.buildUI()
    return a
}

// buildUI constructs the dialog
func (a *AnimationControlsDialog) buildUI() {
    // Implementation details
}

// Play starts animation
func (a *AnimationControlsDialog) Play() {
    a.isPlaying = true
    a.playBtn.Disable()
    a.stopBtn.Enable()
    go a.animationLoop()
}

// Stop halts animation
func (a *AnimationControlsDialog) Stop() {
    a.isPlaying = false
    a.playBtn.Enable()
    a.stopBtn.Disable()
}

// animationLoop handles frame advancement
func (a *AnimationControlsDialog) animationLoop() {
    for a.isPlaying {
        time.Sleep(time.Duration(a.frameDelay) * time.Millisecond)
        a.currentFrame++
        if a.currentFrame >= len(a.sprite.Frames) {
            a.currentFrame = 0
        }
        if a.onFrameChange != nil {
            a.onFrameChange(a.currentFrame)
        }
    }
}

// OnFrameChange sets callback for frame changes
func (a *AnimationControlsDialog) OnFrameChange(fn func(int)) {
    a.onFrameChange = fn
}
```

### Integration Steps:

1. **Create dialog file** in `ui/forms/`
2. **Implement standard methods** (Show, Hide, Close)
3. **Use backend services** from `io/` package
4. **Add to menu** in `ui/window.go`
5. **Test with sample data**

---

## 🔌 Connecting to Backend

### Using Sprite Importer

```go
import "ffvi_editor/io"

// Create importer
importer := &io.SpriteImporter{}

// Configure options
opts := &io.SpriteImportOptions{
    DitheringMethod:   "floyd-steinberg",
    MaxColors:         16,
    PadToTileBoundary: true,
    AutoDetectType:    true,
}

// Perform import
result := importer.Import("path/to/image.png", opts)

if result.Success {
    sprite := result.Sprite
    // Use sprite
} else {
    err := result.Error
    // Handle error
}
```

### Using Palette Editor Backend

```go
import "ffvi_editor/io"

// Create palette editor
editor := &io.PaletteEditor{}

// Generate harmony
colors := editor.GenerateHarmony(baseColor, "complementary")

// Apply transformation
transformedColor := editor.Transform(color, "brighten")
```

### Using Validators

```go
import "ffvi_editor/io"

// Create validator
validator := &io.SpriteValidator{}

// Validate sprite
result := validator.Validate(sprite)

if !result.IsValid {
    for _, err := range result.Errors {
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

## 🎯 Common Patterns

### Pattern 1: Dialog with Callbacks

```go
type MyDialog struct {
    window   fyne.Window
    dialog   dialog.Dialog
    onApply  func(result interface{})
    onCancel func()
}

func NewMyDialog(window fyne.Window) *MyDialog {
    d := &MyDialog{window: window}
    d.buildUI()
    return d
}

func (d *MyDialog) Show() {
    d.dialog.Show()
}

func (d *MyDialog) OnApply(fn func(interface{})) {
    d.onApply = fn
}

func (d *MyDialog) OnCancel(fn func()) {
    d.onCancel = fn
}
```

### Pattern 2: Real-time Preview

```go
// In buildUI()
entry := widget.NewEntry()
entry.OnChanged = func(s string) {
    go d.updatePreview(s)
}

// Async update
func (d *MyDialog) updatePreview(input string) {
    // Process in background
    result := d.process(input)
    
    // Update UI on main thread
    fyne.CurrentApp().Window().Canvas().Refresh()
}
```

### Pattern 3: State Management

```go
// Save original state
originalState := d.sprite.Copy()

// Modify state
d.sprite.ApplyChanges()

// Optionally revert
if cancelled {
    d.sprite = originalState
}
```

### Pattern 4: Error Handling

```go
if err != nil {
    dialog.ShowError(
        fmt.Errorf("operation failed: %w", err),
        d.window,
    )
    return
}

dialog.ShowInformation(
    "Success",
    "Operation completed successfully",
    d.window,
)
```

### Pattern 5: Non-blocking Operations

```go
d.importBtn.Disable()

go func() {
    defer d.importBtn.Enable()
    
    // Long-running operation
    result := d.performLongOperation()
    
    // Update UI when done
    d.statusLabel.SetText("Complete!")
}()
```

---

## 🧪 Testing

### Unit Test Template

```go
package forms

import (
    "testing"
    "ffvi_editor/models"
    "fyne.io/fyne/v2/app"
)

func TestNewSpriteImportDialog(t *testing.T) {
    myApp := app.New()
    window := myApp.NewWindow()
    
    dialog := NewSpriteImportDialog(window)
    if dialog == nil {
        t.Fatal("expected dialog, got nil")
    }
}

func TestImportFile(t *testing.T) {
    myApp := app.New()
    window := myApp.NewWindow()
    
    dialog := NewSpriteImportDialog(window)
    
    // Set test file
    dialog.fileEntry.SetText("testdata/sprite.png")
    
    // Verify state
    if dialog.selectedFile != "testdata/sprite.png" {
        t.Errorf("expected file path, got %s", dialog.selectedFile)
    }
}

func TestPaletteEditorHarmony(t *testing.T) {
    palette := &models.Palette{
        Colors: make([]models.RGB555, 16),
    }
    
    myApp := app.New()
    window := myApp.NewWindow()
    
    editor := NewPaletteEditorDialog(window, palette)
    editor.harmonySelect.SetSelected("Complementary")
    editor.generateHarmony()
    
    // Verify colors were generated
    if editor.palette.Colors[0] == editor.palette.Colors[1] {
        t.Error("harmony colors should differ")
    }
}
```

### Integration Test Template

```go
func TestFullImportWorkflow(t *testing.T) {
    // Setup
    myApp := app.New()
    window := myApp.NewWindow()
    
    // Create dialog
    dialog := NewSpriteImportDialog(window)
    
    // Verify initial state
    if dialog.importedSprite != nil {
        t.Error("should have no sprite initially")
    }
    
    // Setup callback
    var importedSprite *models.FF6Sprite
    dialog.OnImportSuccess(func(sprite *models.FF6Sprite) {
        importedSprite = sprite
    })
    
    // Simulate import
    dialog.fileEntry.SetText("testdata/character.png")
    dialog.performImport()
    
    // Verify callback
    if importedSprite == nil {
        t.Error("callback should have been called")
    }
}
```

---

## ⚡ Performance Tips

### 1. Lazy Load Dialogs
```go
// Bad - loads all on startup
dialogs := map[string]dialog.Dialog{
    "import": forms.NewSpriteImportDialog(window),
    "palette": forms.NewPaletteEditorDialog(window, p),
}

// Good - load on demand
func openImportDialog() {
    d := forms.NewSpriteImportDialog(window)
    d.Show()
}
```

### 2. Cache Preview Images
```go
// Bad - rerender every time
func updatePreview() {
    img := renderSprite()  // Expensive
    canvas.SetImage(img)
}

// Good - cache result
func updatePreview() {
    if cachedImage == nil {
        cachedImage = renderSprite()
    }
    canvas.SetImage(cachedImage)
}
```

### 3. Use Goroutines for I/O
```go
// Bad - blocks UI
sprite := importer.Import(file, opts)

// Good - non-blocking
go func() {
    sprite := importer.Import(file, opts)
    d.handleImportResult(sprite)
}()
```

### 4. Debounce Rapid Updates
```go
// Bad - updates every keystroke
entry.OnChanged = func(s string) {
    d.updatePreview(s)
}

// Good - debounce
var debounceTimer *time.Timer
entry.OnChanged = func(s string) {
    if debounceTimer != nil {
        debounceTimer.Stop()
    }
    debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
        d.updatePreview(s)
    })
}
```

### 5. Limit Palette Grid Size
```go
// Bad - render all 256 colors
for i := 0; i < 256; i++ {
    grid.Add(createColorBox(colors[i]))
}

// Good - only show 16 (FF6 limitation)
for i := 0; i < 16; i++ {
    grid.Add(createColorBox(colors[i]))
}
```

---

## 📋 Checklist for New Components

- [ ] Create file in `ui/forms/`
- [ ] Implement `NewXxxDialog()` constructor
- [ ] Implement `Show()` method
- [ ] Implement `Hide()` method
- [ ] Add callback support (OnApply, OnCancel, etc.)
- [ ] Implement `buildUI()` to create layout
- [ ] Add state management
- [ ] Use backend services for logic
- [ ] Implement error handling
- [ ] Add to menu in `ui/window.go`
- [ ] Write unit tests
- [ ] Test with sample data
- [ ] Profile performance
- [ ] Document public API
- [ ] Add inline comments
- [ ] Format with `go fmt`

---

## 🚀 Quick Start - Add Your Component

### Step 1: Create File
```bash
touch ui/forms/my_feature_dialog.go
```

### Step 2: Scaffold Code
```go
package forms

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/dialog"
)

type MyFeatureDialog struct {
    window fyne.Window
    dialog dialog.Dialog
}

func NewMyFeatureDialog(window fyne.Window) *MyFeatureDialog {
    d := &MyFeatureDialog{window: window}
    d.buildUI()
    return d
}

func (d *MyFeatureDialog) buildUI() {
    // TODO: Build UI
}

func (d *MyFeatureDialog) Show() {
    d.dialog.Show()
}
```

### Step 3: Build UI
```go
func (d *MyFeatureDialog) buildUI() {
    content := container.NewVBox(
        widget.NewLabel("My Feature"),
        // Add widgets here
    )
    
    d.dialog = dialog.NewCustom(
        "My Feature",
        "Close",
        content,
        d.window,
    )
    d.dialog.Resize(fyne.NewSize(400, 300))
}
```

### Step 4: Add to Menu
```go
// In ui/window.go
fyne.NewMenuItem("My Feature...", func() {
    d := forms.NewMyFeatureDialog(g.window)
    d.Show()
}),
```

### Step 5: Test
```bash
go test ./ui/forms -v
```

---

## 📞 Support Resources

**Documentation:**
- [Fyne Documentation](https://developer.fyne.io/api/container)
- [Backend API Reference](PHASE_13_API_REFERENCE.md)
- [Integration Guide](PHASE_13_INTEGRATION_GUIDE.md)

**Code Examples:**
- `ui/forms/plugin_manager.go` - Complex dialog example
- `ui/forms/palette_editor_dialog.go` - Color editing pattern
- `ui/forms/sprite_import_dialog.go` - Import workflow pattern

**Contact:**
- Review existing dialogs for patterns
- Check `models/` for data structures
- See `io/` for backend services

---

## ✨ Success Criteria

Your new component is ready when:
- ✅ Code compiles with zero errors
- ✅ UI displays correctly
- ✅ Callbacks work as expected
- ✅ Errors handled gracefully
- ✅ Performance acceptable (<200ms)
- ✅ Tests pass
- ✅ Menu item shows and works
- ✅ Documentation complete

---

**Happy coding! 🎉**
