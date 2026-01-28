# 🎨 PHASE 13: SPRITE EDITOR BACKEND API REFERENCE

## Quick Start

### Import a Sprite
```go
package main

import (
	"log"
	"ffvi_editor/io"
	"ffvi_editor/models"
)

func ImportSprite(imagePath string) {
	// Create importer
	importer := io.NewSpriteImporter()
	
	// Setup import options
	opts := models.NewSpriteImportOptions()
	opts.SourcePath = imagePath
	opts.TargetType = models.SpriteTypeCharacter
	opts.MaxColors = 16
	opts.DitherMethod = "floyd-steinberg"
	
	// Import
	result := importer.Import(opts)
	
	// Check result
	if result.Success {
		sprite := result.Sprite
		log.Printf("✓ Imported sprite: %s", sprite.Name)
		log.Printf("  Dimensions: %dx%d", sprite.Width, sprite.Height)
		log.Printf("  Frames: %d", sprite.Frames)
	} else {
		for _, err := range result.Errors {
			log.Printf("✗ %s", err.Message)
		}
	}
}
```

---

## Core Types

### FF6Sprite
```go
type FF6Sprite struct {
	ID          string           // Unique identifier
	Name        string           // Display name
	Type        SpriteType       // Character, Battle, Portrait, etc.
	Description string           // User notes
	Data        []byte           // Raw sprite data
	Palette     *Palette         // 16-color palette
	Frames      int              // Number of animation frames
	FrameRate   int              // FPS
	FrameDurations []int         // Duration per frame (ms)
	Width       int              // Pixels
	Height      int              // Pixels
	IsCompressed bool            // LZ77 compressed
	Author      string           // Creator
	License     string           // Usage rights
	Tags        []string         // Metadata tags
	CreatedDate time.Time        // Creation timestamp
	ModifiedDate time.Time       // Last edit
}
```

### Palette
```go
type Palette struct {
	Colors   [16]RGB555  // 16 colors, index 0 = transparent
	Name     string      // Palette name
	Author   string      // Creator
	Tags     []string    // Metadata
	License  string      // Usage rights
	Created  time.Time   // Creation date
	Modified time.Time   // Last modified
	Notes    string      // User notes
}
```

### RGB555
```go
type RGB555 struct {
	R uint8  // 0-31 (5 bits)
	G uint8  // 0-31 (5 bits)
	B uint8  // 0-31 (5 bits)
}

// Methods
func (c RGB555) ToRGB888() (r, g, b uint8)
func (c RGB555) ToColor() color.Color
```

### SpriteType
```go
const (
	SpriteTypeCharacter  // 16x24 (standing sprite)
	SpriteTypeBattle     // 32x32 (battle sprite)
	SpriteTypePortrait   // 48x48 (menu portrait)
	SpriteTypeNPC        // 16x24 (NPC sprite)
	SpriteTypeEnemy      // 32x32 (enemy sprite)
	SpriteTypeOverworld  // 16x16 (map sprite)
)

// Methods
func (st SpriteType) String() string
func (st SpriteType) GetDimensions() (width, height int)
```

---

## Image Import

### ImageDecoder
```go
decoder := io.NewImageDecoder()
img, format, err := decoder.Decode("/path/to/image.png")
// Supported: PNG, GIF, BMP, JPEG
// Returns: image.Image, format string ("png", "gif", etc), error
```

### PaletteExtractor
```go
extractor := io.NewPaletteExtractor()
palette, err := extractor.Extract(img, 16)
// Returns best 16 colors from image
// Optimized for large images (samples every Nth pixel)
```

### ColorQuantizer
```go
quantizer := io.NewColorQuantizer("floyd-steinberg")
// Methods: "floyd-steinberg", "bayer", "none"

quantizedImg, palette, err := quantizer.Quantize(img, 16)
// Reduces image to 16-color palette with dithering
```

---

## Sprite Conversion

### FF6SpriteConverter
```go
converter := io.NewFF6SpriteConverter()

// Convert image to FF6 format
sprite := converter.ToFF6Format(img, palette, models.SpriteTypeCharacter)

// Decode FF6 sprite back to image
img, err := converter.DecodeFF6Sprite(sprite)

// Pad image to tile boundaries (multiples of 8)
paddedImg := converter.PadImage(img, nil)

// Fit image to sprite dimensions
fitted := converter.fitImage(img, 16, 24)
```

---

## Sprite Validation

### SpriteValidator
```go
validator := io.NewSpriteValidator()
result := validator.Validate(sprite)

if result.IsValid {
	log.Println("✓ Sprite is valid")
} else {
	for _, err := range result.Errors {
		if err.Level == "error" {
			log.Printf("✗ %s", err.Message)
		} else {
			log.Printf("⚠ %s", err.Message)
		}
	}
}
```

**Checks:**
- Dimensions match sprite type
- Dimensions are multiples of 8
- Palette has 16 colors
- All RGB values are 5-bit (0-31)
- Data size is reasonable
- Frame count > 0

---

## Import Pipeline

### SpriteImporter
```go
importer := io.NewSpriteImporter()

// Single import
opts := models.NewSpriteImportOptions()
opts.SourcePath = "/path/to/sprite.png"
opts.MaxColors = 16
opts.DitherMethod = "floyd-steinberg"
opts.AutoPadding = true
opts.AutoDetectType = true

result := importer.Import(opts)

type ImportResult struct {
	Sprite      *FF6Sprite         // Resulting sprite
	Success     bool               // Import succeeded
	Errors      []ValidationError  // Error messages
	Warnings    []string           // Warning messages
	StepsFailed []string           // Which step failed
	Duration    time.Duration      // Time taken
}
```

### Batch Import
```go
paths := []string{"/img1.png", "/img2.png", "/img3.png"}

// Simple batch
results := importer.ImportMultiple(paths, opts)

// Advanced batch with options
batchOpts := &io.BatchImportOptions{
	BaseOptions:     opts,
	StopOnError:     false,
	ReusePalette:    true,  // Share palette across imports
	CreateZipArchive: false,
}

results := importer.BatchImport(paths, batchOpts)
```

---

## Palette Editing

### PaletteEditor
```go
editor := io.NewPaletteEditor(sprite.Palette)

// Set individual color
editor.SetColor(0, models.RGB555{R: 31, G: 0, B: 0}) // Red

// Get color
color, _ := editor.GetColor(0)

// Generate harmonious palette
editor.GenerateHarmony(baseColor, "triadic")
// Schemes: "complementary", "triadic", "analogous", 
//          "monochromatic", "split-complementary", "tetradic"

// Apply transformations
editor.ApplyTransform("brighten", 0.3)
editor.ApplyTransform("saturate", 0.5)
// Transforms: "brighten", "darken", "saturate", "desaturate",
//             "shift-hue", "invert", "grayscale", "sepia"

// Transform a color range
editor.ApplyTransformToRange("darken", 0.2, 8, 15)

// Swap colors
editor.SwapColors(0, 15)

// Rotate colors
editor.RotateColors(0, 7, 1)  // Rotate first 8 colors left

// Fill with gradient
editor.GradientFill(0, 7, darkColor, lightColor)

// Get result
palette := editor.GetPalette()
```

### ColorHarmonizer
```go
harmonizer := io.NewColorHarmonizer()

palette := harmonizer.Generate(baseColor, "triadic")
// Generates full 16-color palette based on harmony scheme
```

### ColorTransformer
```go
transformer := io.NewColorTransformer()

newColor := transformer.Apply(color, "brighten", 0.5)
// Transforms: "brighten", "darken", "saturate", "desaturate",
//             "shift-hue", "invert", "grayscale", "sepia"
// Amount: 0.0-1.0
```

---

## Export

### SpriteExporter
```go
exporter := io.NewSpriteExporter()

// Export to PNG
err := exporter.ExportToPNG(sprite, "/output/sprite.png")

// Export animated GIF (future)
// err := exporter.ExportToGIF(sprite, "/output/sprite.gif", 12) // FPS
```

---

## Undo/Redo

### SpriteHistory
```go
history := models.NewSpriteHistory(50)  // 50-level undo

// Save state
history.Push(sprite)

// Undo
prevSprite := history.Undo()

// Redo
nextSprite := history.Redo()

// Check availability
if history.CanUndo() {
	// Undo button enabled
}

if history.CanRedo() {
	// Redo button enabled
}
```

---

## Color Conversion

### RGB555 Methods
```go
rgb555 := models.RGB555{R: 31, G: 16, B: 0}

// To standard RGB
r, g, b := rgb555.ToRGB888()

// To Go color
goColor := rgb555.ToColor()

// From RGB888
rgb555 = models.FromRGB888(255, 128, 0)
```

---

## Common Patterns

### Import and Apply to Character
```go
importer := io.NewSpriteImporter()
opts := models.NewSpriteImportOptions()
opts.SourcePath = "terra_custom.png"
opts.TargetType = models.SpriteTypeCharacter

result := importer.Import(opts)
if result.Success {
	sprite := result.Sprite
	sprite.Name = "Terra - Fire Outfit"
	sprite.Author = "CustomArtist"
	sprite.Tags = []string{"character", "terra", "custom"}
	
	// Apply to game (next phase)
	// saveFile.SetCharacterSprite(0, sprite) // Terra
}
```

### Quick Palette Swap
```go
// Get current sprite palette
currentPalette := sprite.Palette

// Create editor
editor := io.NewPaletteEditor(currentPalette)

// Apply preset transformation
editor.ApplyTransform("shift-hue", 0.5)  // Shift colors

// Update sprite
sprite.Palette = editor.GetPalette()
sprite.ModifiedDate = time.Now()
```

### Batch Process with Status
```go
importer := io.NewSpriteImporter()
imageFiles := []string{...}

for i, file := range imageFiles {
	opts := models.NewSpriteImportOptions()
	opts.SourcePath = file
	
	result := importer.Import(opts)
	
	progressPct := (i + 1) * 100 / len(imageFiles)
	if result.Success {
		log.Printf("[%d%%] ✓ %s", progressPct, filepath.Base(file))
	} else {
		log.Printf("[%d%%] ✗ %s", progressPct, filepath.Base(file))
	}
}
```

---

## Error Handling

### Import Errors
```go
result := importer.Import(opts)

if !result.Success {
	// Errors (hard failures)
	for _, err := range result.Errors {
		if err.Level == "error" {
			log.Printf("Failed: %s", err.Message)
		}
	}
	
	// Warnings (soft warnings)
	for _, warning := range result.Warnings {
		log.Printf("Warning: %s", warning)
	}
	
	// What failed
	for _, step := range result.StepsFailed {
		log.Printf("Failed at: %s", step)
	}
	
	// How long it took
	log.Printf("Duration: %v", result.Duration)
}
```

### Validation Errors
```go
validator := io.NewSpriteValidator()
result := validator.Validate(sprite)

for _, err := range result.Errors {
	display := io.PrintValidationResult(result)
	// Shows formatted error/warning list
}
```

---

## Performance Notes

- **Import time:** ~50-200ms depending on image size
- **Palette extraction:** Samples large images (>1MP) for speed
- **Color quantization:** Floyd-Steinberg takes ~100-300ms
- **Memory:** ~1-2MB for typical 16x24 sprites

---

## Future Extensions

The API is designed for easy expansion:

1. **Compression:** Add LZ77 compression to `FF6SpriteConverter`
2. **Animation:** Frame duration arrays already in place
3. **Editing:** Pixel-level editing via image manipulation
4. **Batch operations:** Already supported in pipeline
5. **Cloud storage:** Add metadata for sync (author, license, tags)

---

*API Reference - Phase 13 Foundation Complete*
