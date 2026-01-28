# 🎨 PHASE 13: SPRITE EDITOR & ASSET MANAGER - COMPLETE SPECIFICATION

**Version:** 2.0 (Enhanced)  
**Date:** January 17, 2026  
**Status:** Ready for Implementation  
**Author:** Development Team  
**Reviewers:** Product, Engineering, Community

---

## 📋 Table of Contents

1. [Executive Summary](#executive-summary)
2. [Strategic Vision](#strategic-vision)
3. [Quick Wins (Week 1)](#quick-wins-week-1)
4. [FF6 Technical Specifications](#ff6-technical-specifications)
5. [Core Features (Tier 1)](#core-features-tier-1)
6. [Advanced Features (Tier 2)](#advanced-features-tier-2)
7. [Community Features (Tier 3)](#community-features-tier-3)
8. [Import/Export Compatibility](#import-export-compatibility)
9. [Performance Optimization](#performance-optimization)
10. [Keyboard Shortcuts](#keyboard-shortcuts)
11. [Batch Operations](#batch-operations)
12. [Sprite Validation Rules](#sprite-validation-rules)
13. [Tutorial System](#tutorial-system)
14. [Cloud Backup System](#cloud-backup-system)
15. [Phase 11 Integration](#phase-11-integration)
16. [Technical Architecture](#technical-architecture)
17. [Implementation Timeline](#implementation-timeline)
18. [Testing Strategy](#testing-strategy)
19. [Risk Mitigation](#risk-mitigation)
20. [Success Metrics](#success-metrics)
21. [User Personas](#user-personas)
22. [Deployment Plan](#deployment-plan)
23. [ROI Analysis](#roi-analysis)
24. [Mobile/Touch Roadmap](#mobile-touch-roadmap)
25. [Next Steps](#next-steps)

---

## Executive Summary

**Plugin Name:** FF6 Sprite Studio  
**Version:** 1.0.0  
**Target:** Phase 13 (Post Phase 12.3 completion)  
**Effort:** 5 weeks, 2 developers  
**Total LOC:** ~10,700 lines (8,500 backend + 2,200 UI)  
**Risk:** Medium (graphics manipulation complexity)  
**Impact:** VERY HIGH (unlocks visual customization ecosystem)  
**ROI:** 480x in Year 1  
**Break-Even:** < 1 week

### What This Delivers

Transform FF6 Save Editor into a **complete visual customization platform** by adding:
- ✅ Professional-grade sprite editing with pixel-perfect canvas
- ✅ Advanced palette management (16-color FF6 format)
- ✅ Multi-frame animation tools with preview
- ✅ Import from PNG/GIF/Aseprite and popular tools
- ✅ Export to ROM patches, texture packs, community hub
- ✅ AI-assisted sprite generation and enhancement
- ✅ Community asset library with 10,000+ sprites
- ✅ Real-time collaboration and version control

### Key Benefits

**For Users:**
- 🎨 Design custom character sprites in minutes (vs hours in external tools)
- 🎨 Import fan art and mods seamlessly
- 🎨 Share creations with community instantly
- 🎨 No need to learn complex ROM hacking tools

**For Community:**
- 🌟 Vibrant sprite ecosystem (10,000+ custom sprites Year 1)
- 🌟 Lower barrier to entry for sprite artists
- 🌟 Enables sprite-based mods and texture packs
- 🌟 Differentiates FF6 editor from all competitors

**For Business:**
- 💰 480x ROI in Year 1 (user time savings)
- 💰 Enables Phase 14+ content marketplace
- 💰 Premium sprite packs monetization opportunity
- 💰 Unassailable competitive moat

---

## Strategic Vision

### The Big Picture

**Current State:**
- Users must use external tools (Aseprite, GraphicsGale) to create sprites
- Manual hex editing to apply sprites to save files
- No easy way to share custom sprites
- High barrier to entry for visual customization

**Future State (Phase 13):**
- One-click sprite import from popular tools
- Professional pixel art editor built-in
- Community sprite library with instant download
- AI-assisted sprite generation
- Seamless integration with save editing workflow

**Ultimate Vision (Phase 14+):**
- Content marketplace for premium sprite packs
- Mobile companion app with touch editing
- Real-time multiplayer sprite jam sessions
- Procedural sprite generation with ML
- Cross-game sprite compatibility (FF4, FF5, etc.)

### Competitive Landscape

| Feature | FF6 Editor (Current) | FF6 Editor (Phase 13) | Competitors |
|---------|---------------------|---------------------|-------------|
| Sprite Editing | ❌ None | ✅ Full Suite | ⚠️ Basic |
| Import Tools | ❌ Manual | ✅ Aseprite/etc | ❌ None |
| Community Library | ❌ None | ✅ 10,000+ sprites | ⚠️ Forums |
| Animation Tools | ❌ None | ✅ Professional | ❌ None |
| AI Generation | ❌ None | ✅ Style transfer | ❌ None |
| ROM Patching | ⚠️ Manual | ✅ One-click | ⚠️ Manual |

**Competitive Advantage:** Phase 13 creates an **unassailable moat** - no other save editor has integrated sprite editing.

---

## Quick Wins (Week 1)

Following the Phase 11+ pattern, we start with **3 Quick Wins** to validate the approach before building the full editor.

### Quick Win #1: Basic Sprite Import (3 days)

**What:** Import PNG sprites and apply to characters
**Why First:** Validates core import pipeline with minimal risk
**Effort:** 3 days, 1 developer
**LOC:** ~400 lines

**Features:**
- ✅ Import single PNG file (16x24 or 32x32)
- ✅ Auto-extract 16-color palette
- ✅ Map to FF6 format
- ✅ Apply to one character (Terra for testing)
- ✅ Preview before/after
- ✅ Undo/rollback

**Technical Stack:**
```go
type BasicImporter struct {
    decoder      *png.Decoder
    palettizer   *SimplePaletteExtractor  // Reduce to 16 colors
    converter    *FF6FormatConverter      // To 5-bit RGB
    validator    *BasicValidator          // Size/color checks
}

func ImportSprite(path string) (*FF6Sprite, error) {
    // 1. Decode PNG
    img, err := png.Decode(file)
    
    // 2. Extract palette (16 colors)
    palette := ExtractPalette(img, 16)
    
    // 3. Convert to FF6 format
    sprite := ConvertToFF6Format(img, palette)
    
    // 4. Validate
    if err := ValidateSprite(sprite); err != nil {
        return nil, err
    }
    
    return sprite, nil
}
```

**Success Criteria:**
- ✅ Import 10 test sprites successfully
- ✅ No data corruption in save files
- ✅ Import completes in <500ms
- ✅ Preview accuracy 100%

**Deliverable:** Working import dialog with "Import Sprite" button

---

### Quick Win #2: Palette Swap Tool (2 days)

**What:** Recolor existing sprites with preset palettes
**Why Second:** Instant visual customization without complex editing
**Effort:** 2 days, 1 developer
**LOC:** ~300 lines

**Features:**
- ✅ 10 preset palettes (fire, ice, poison, gold, etc.)
- ✅ Apply palette to any character
- ✅ Real-time preview
- ✅ Save as template
- ✅ Undo/rollback

**Preset Palettes:**
```go
var PresetPalettes = map[string]Palette{
    "Fire Theme":    {Red/Orange/Yellow hues},
    "Ice Theme":     {Blue/Cyan/White hues},
    "Poison Theme":  {Purple/Green/Black hues},
    "Gold Theme":    {Yellow/Orange/Brown metallic},
    "Shadow Theme":  {Black/Dark Gray/Purple},
    "Forest Theme":  {Green/Brown/Earth tones},
    "Ocean Theme":   {Blue/Teal/White aquatic},
    "Sunset Theme":  {Orange/Pink/Purple warm},
    "Monochrome":    {Grayscale 16 shades},
    "Rainbow":       {Full spectrum},
}
```

**UI Workflow:**
1. User selects character (e.g., Locke)
2. Opens "Palette Swap" dialog
3. Clicks palette thumbnail (e.g., "Fire Theme")
4. Preview shows Locke with fire colors
5. Click "Apply" → character updated in save
6. Undo available from toolbar

**Success Criteria:**
- ✅ All 10 presets look visually appealing
- ✅ Preview updates in <100ms
- ✅ No palette corruption
- ✅ Works on all 14 playable characters

**Deliverable:** Palette swap dialog with 10 working presets

---

### Quick Win #3: Sprite Preview System (2 days)

**What:** Preview sprites in authentic FF6 game context
**Why Third:** Critical feedback loop for sprite quality
**Effort:** 2 days, 1 developer
**LOC:** ~350 lines

**Features:**
- ✅ Preview character standing sprite
- ✅ Preview walk cycle animation (3 frames)
- ✅ Preview battle stance (6 frames)
- ✅ Side-by-side comparison (before/after)
- ✅ Authentic FF6 background rendering
- ✅ Zoom levels (1x, 2x, 4x)

**Preview Contexts:**
```go
type PreviewContext int
const (
    ContextStanding PreviewContext = iota  // Static, front-facing
    ContextWalking                          // 3-frame walk animation
    ContextBattle                           // 6-frame battle animations
    ContextPortrait                         // 48x48 menu portrait
    ContextWorldMap                         // 16x16 overworld sprite
)

func RenderPreview(sprite *FF6Sprite, context PreviewContext) *image.RGBA {
    // Render sprite in authentic FF6 style
    // With pixel-perfect scaling
    // Against appropriate background
}
```

**UI Layout:**
```
┌─────────────────────────────────────────┐
│ Sprite Preview                          │
├─────────────────────────────────────────┤
│  Before           │  After              │
│  ┌─────────┐     │  ┌─────────┐        │
│  │ Original│     │  │ Modified│        │
│  │  Terra  │     │  │  Terra  │        │
│  └─────────┘     │  └─────────┘        │
│                  │                      │
│  Context: [Standing ▼]  Zoom: [2x ▼]   │
│  ☑ Animate  ☐ Show Grid  ☐ Onion Skin │
│                                         │
│  [◀ Previous Frame] [▶ Play] [▶| Next] │
└─────────────────────────────────────────┘
```

**Success Criteria:**
- ✅ Pixel-perfect rendering (no artifacts)
- ✅ Animations play at correct speed (60 FPS)
- ✅ Zoom maintains sharp pixels
- ✅ Preview loads in <200ms

**Deliverable:** Preview panel integrated into sprite editor dialog

---

### Quick Wins Summary

**Week 1 Deliverables:**
- ✓ Day 1-3: Basic import working (400 LOC)
- ✓ Day 4-5: Palette swap working (300 LOC)
- ✓ Day 6-7: Preview system working (350 LOC)
- **Total:** 1,050 LOC, fully functional prototype

**Validation Checkpoint:**
After Week 1, assess:
- ✅ Import quality (is palette extraction good enough?)
- ✅ Performance (are operations fast enough?)
- ✅ User feedback (is workflow intuitive?)
- ✅ Technical feasibility (any blockers?)

**Decision Point:**
- 🟢 **GO:** Proceed to full editor (Weeks 2-5)
- 🟡 **PIVOT:** Adjust approach based on learnings
- 🔴 **STOP:** Halt if fundamental issues discovered

**Expected Outcome:** 🟢 GO (95% confidence)

---


## FF6 Technical Specifications

### SNES Hardware Constraints

Understanding FF6's technical limitations is **critical** for sprite editing:

```
FF6 SPRITE FORMAT SPECIFICATIONS:
├── Platform: Super Nintendo Entertainment System (SNES)
├── PPU: Picture Processing Unit (5C77/5C78 chips)
├── Color Depth: 4-bit indexed (16 colors per palette)
├── Palette Format: 5-bit RGB (0-31 per channel, 32 levels)
├── Resolution: 256x224 pixels (NTSC), 256x239 (PAL)
├── Sprite Sizes: 8x8, 16x16, 32x32, 64x64 pixels
├── Max Sprites: 128 sprites (hardware limit)
├── Max On-Screen: 32 sprites per scanline
├── Tile-Based: All graphics use 8x8 tile blocks
├── Compression: Custom LZ77 variant
└── Memory: 8KB sprite RAM, 512B palette RAM
```

### Character Sprite Specifications

**Standing/Walking Sprites:**
- Dimensions: 16x24 pixels (2 tiles wide × 3 tiles tall)
- Frames: 3 frames per direction (standing, walk1, walk2)
- Directions: 4 directions (up, down, left, right)
- Total: 12 frames per character
- Palette: 16 colors (index 0 = transparent)
- Size: ~1.5KB per character (compressed)

**Battle Sprites:**
- Dimensions: 32x32 pixels (4 tiles × 4 tiles)
- Frames: 6 frames (idle, attack, magic, damage, victory, dead)
- Palette: 16 colors (shared with standing sprite)
- Size: ~4KB per character (compressed)
- Special: Some attacks use additional sprite overlays

**Portrait Sprites:**
- Dimensions: 48x48 pixels (menu/dialogue)
- Frames: Usually 2 (normal, talking)
- Palette: 16 colors (can differ from battle palette)
- Size: ~1KB per portrait

### Palette Format (5-bit RGB)

FF6 uses **5-bit RGB color** (not standard 8-bit!):

```go
type RGB555 struct {
    R uint8  // 0-31 (5 bits)
    G uint8  // 0-31 (5 bits)  
    B uint8  // 0-31 (5 bits)
}

// Convert to standard RGB888
func (c RGB555) ToRGB888() (r, g, b uint8) {
    r = uint8((c.R * 255) / 31)
    g = uint8((c.G * 255) / 31)
    b = uint8((c.B * 255) / 31)
    return
}

// Convert from standard RGB888
func FromRGB888(r, g, b uint8) RGB555 {
    return RGB555{
        R: uint8((int(r) * 31) / 255),
        G: uint8((int(g) * 31) / 255),
        B: uint8((int(b) * 31) / 255),
    }
}
```

**Important:** Color precision loss occurs (256 levels → 32 levels per channel)

### Tile-Based Graphics

All FF6 graphics are **tile-based** (8x8 blocks):

```
16x24 Character Sprite Structure:
┌────┬────┐
│ T0 │ T1 │  ← Top row (8x8 + 8x8)
├────┼────┤
│ T2 │ T3 │  ← Middle row
├────┼────┤
│ T4 │ T5 │  ← Bottom row
└────┴────┘

Each tile: 64 pixels, 4-bit indexed (32 bytes)
Total: 6 tiles × 32 bytes = 192 bytes per frame
```

### Compression Format

FF6 uses **custom LZ77 compression**:
- Window size: 4KB lookback buffer
- Match length: 3-18 bytes
- Compression ratio: ~40-60% (sprites compress well due to large transparent areas)
- Decompression: Fast (hardware-accelerated on SNES)

### Memory Layout

```
SNES Sprite Memory Map:
├── 0x7E2000-0x7E3FFF: Sprite graphics (8KB)
├── 0x7E4000-0x7E41FF: Sprite palettes (512B)
├── 0x7E4200-0x7E43FF: Sprite attributes (512B)
└── 0x7E4400-0x7E4FFF: Compression workspace (3KB)
```

### Validation Requirements

Sprites MUST meet these criteria to work in FF6:

**✅ Valid Sprite Checklist:**
- [ ] Dimensions are multiples of 8 (tile-based)
- [ ] Maximum 16 colors (including transparency)
- [ ] Color index 0 is transparent
- [ ] Fits within size limits (16x24, 32x32, etc.)
- [ ] Uses 5-bit RGB values (0-31 per channel)
- [ ] Compressed size ≤ 8KB
- [ ] No more than 128 tiles total

**❌ Common Import Errors:**
- Odd dimensions (e.g., 17x25) → fails
- >16 colors → requires quantization
- No transparent color → sprites overlap incorrectly
- >64x64 pixels → exceeds hardware limit
- 8-bit RGB → precision loss warning

---

## Core Features (Tier 1: Foundation)

### 1. Sprite Import & Export System

**What:** Import external sprite graphics into save files  
**LOC:** ~600 lines  
**Time:** 3 days  
**Risk:** Low (standard image processing)

#### Import Pipeline Architecture

```go
type SpriteImporter struct {
    decoder      ImageDecoder          // PNG/GIF/BMP support
    palettizer   PaletteExtractor      // Extract 16-color palette
    quantizer    ColorQuantizer        // Reduce colors if >16
    converter    FF6SpriteConverter    // Convert to FF6 format
    compressor   LZ77Compressor        // Compress sprite data
    validator    SpriteValidator       // Verify FF6 compliance
}

// Import pipeline
func (s *SpriteImporter) Import(path string, spriteType SpriteType) (*FF6Sprite, error) {
    // Step 1: Decode image file
    img, format, err := s.decoder.Decode(path)
    if err != nil {
        return nil, fmt.Errorf("decode failed: %w", err)
    }
    
    // Step 2: Extract palette
    palette, err := s.palettizer.Extract(img, 16)
    if err != nil {
        // Image has >16 colors, quantize
        img, palette = s.quantizer.Reduce(img, 16)
    }
    
    // Step 3: Convert to FF6 format
    sprite := s.converter.ToFF6Format(img, palette, spriteType)
    
    // Step 4: Compress
    sprite.Data = s.compressor.Compress(sprite.Data)
    
    // Step 5: Validate
    if err := s.validator.Validate(sprite); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    return sprite, nil
}
```

#### Supported Image Formats

| Format | Support | Notes |
|--------|---------|-------|
| PNG | ✅ Full | Recommended (lossless) |
| GIF | ✅ Full | Supports animation frames |
| BMP | ✅ Full | Uncompressed |
| JPEG | ⚠️ Partial | Lossy, not recommended |
| WEBP | ⚠️ Partial | Requires external decoder |
| TIFF | ❌ None | Rarely used for pixel art |

#### Export Formats

```go
type ExportFormat string
const (
    ExportPNG        ExportFormat = "png"         // Standard PNG
    ExportGIF        ExportFormat = "gif"         // Animated GIF
    ExportSpriteSheet ExportFormat = "spritesheet" // All frames in grid
    ExportJSON       ExportFormat = "json"        // Metadata + base64
    ExportROMPatch   ExportFormat = "ips"         // IPS patch file
    ExportTexturePack ExportFormat = "texpack"    // Emulator texture pack
)

func (s *SpriteExporter) Export(sprite *FF6Sprite, format ExportFormat, path string) error {
    switch format {
    case ExportPNG:
        return s.exportPNG(sprite, path)
    case ExportGIF:
        return s.exportAnimatedGIF(sprite, path)
    case ExportSpriteSheet:
        return s.exportSpriteSheet(sprite, path)
    case ExportROMPatch:
        return s.exportIPSPatch(sprite, path)
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}
```

#### Features

**Import Features:**
- ✅ Drag & drop file import
- ✅ Batch import multiple sprites
- ✅ Auto-detect sprite type (16x24 vs 32x32)
- ✅ Color quantization (>16 colors → 16 colors)
- ✅ Dithering options (Floyd-Steinberg, Bayer)
- ✅ Preview before applying
- ✅ Import history (last 20 imports)

**Export Features:**
- ✅ Export to PNG (lossless)
- ✅ Export animated GIF (all frames)
- ✅ Export sprite sheet (grid layout)
- ✅ Export ROM patch (IPS format)
- ✅ Export texture pack (emulator-compatible)
- ✅ Batch export all sprites
- ✅ Export with metadata (JSON)

---

### 2. Advanced Palette Editor

**What:** Professional 16-color palette management  
**LOC:** ~900 lines  
**Time:** 4 days  
**Risk:** Low

#### Palette Editor Architecture

```go
type PaletteEditor struct {
    palette      Palette              // 16 colors, 5-bit RGB
    history      []PaletteSnapshot    // Undo/redo (50 levels)
    harmonizer   ColorHarmonizer      // Generate harmonious colors
    transformer  ColorTransformer     // Batch transformations
    analyzer     ColorAnalyzer        // Contrast, accessibility checks
}

type Palette struct {
    Colors   [16]RGB555
    Name     string
    Tags     []string
    Author   string
    License  string
    Created  time.Time
    Modified time.Time
}

// Core operations
func (p *PaletteEditor) SetColor(index int, color RGB555) {
    if index < 0 || index >= 16 {
        return
    }
    p.saveSnapshot()  // For undo
    p.palette.Colors[index] = color
}

func (p *PaletteEditor) GenerateHarmony(baseColor RGB555, scheme HarmonyScheme) Palette {
    // Generate harmonious palette based on color theory
    switch scheme {
    case Complementary:
        return p.harmonizer.Complementary(baseColor)
    case Triadic:
        return p.harmonizer.Triadic(baseColor)
    case Analogous:
        return p.harmonizer.Analogous(baseColor)
    case Monochromatic:
        return p.harmonizer.Monochromatic(baseColor)
    }
}
```

#### Color Harmony Schemes

**1. Complementary** (opposite on color wheel)
- Base color + opposite hue
- High contrast, vibrant
- Example: Blue + Orange

**2. Triadic** (3 colors, evenly spaced)
- Base + 120° + 240° on wheel
- Balanced, colorful
- Example: Red + Yellow + Blue

**3. Analogous** (adjacent colors)
- Base + neighbors (±30°)
- Harmonious, subtle
- Example: Blue + Blue-Green + Green

**4. Monochromatic** (same hue, varying lightness)
- Base color with different brightnesses
- Cohesive, elegant
- Example: Light Blue → Dark Blue

**5. Split-Complementary**
- Base + two adjacent to complement
- Less tension than complementary
- Example: Blue + Red-Orange + Yellow-Orange

#### Color Transformations

```go
type ColorTransform int
const (
    TransformBrighten ColorTransform = iota
    TransformDarken
    TransformSaturate
    TransformDesaturate
    TransformShiftHue
    TransformInvert
    TransformGrayscale
    TransformSepia
)

func (p *PaletteEditor) ApplyTransform(transform ColorTransform, amount float64) {
    p.saveSnapshot()
    for i := range p.palette.Colors {
        p.palette.Colors[i] = p.transformer.Apply(
            p.palette.Colors[i], 
            transform, 
            amount,
        )
    }
}
```

#### UI Components

```
┌────────────────────────────────────────────────────┐
│ Palette Editor                           [×]       │
├────────────────────────────────────────────────────┤
│ Current Palette: "Fire Theme"                      │
│                                                    │
│ ┌──┬──┬──┬──┬──┬──┬──┬──┐                         │
│ │00│01│02│03│04│05│06│07│  ← Color swatches      │
│ └──┴──┴──┴──┴──┴──┴──┴──┘                         │
│ ┌──┬──┬──┬──┬──┬──┬──┬──┐                         │
│ │08│09│10│11│12│13│14│15│                         │
│ └──┴──┴──┴──┴──┴──┴──┴──┘                         │
│                                                    │
│ Selected: Color 03                                 │
│ ┌─────────────────┐                               │
│ │    Color Picker │  RGB: (31, 20, 5)             │
│ │   ┌───────────┐ │  HEX: #F8A028                 │
│ │   │  [Color]  │ │  HSL: (35°, 85%, 60%)         │
│ │   └───────────┘ │                               │
│ │   R: [████░░░░] │  31  (5-bit)                  │
│ │   G: [███░░░░░] │  20                           │
│ │   B: [█░░░░░░░] │   5                           │
│ └─────────────────┘                               │
│                                                    │
│ Harmony: [Complementary ▼] [Generate]             │
│ Transform: [Brighten ▼] Amount: [░░░░] [Apply]    │
│                                                    │
│ [Import Palette] [Export Palette] [Templates...] │
│                                                    │
│ Preview on sprite: ☑ Real-time                    │
│ ┌────────────────────────────────────────────┐    │
│ │  [Before]         [After]                  │    │
│ │  Terra sprite     Terra sprite (recolored) │    │
│ └────────────────────────────────────────────┘    │
│                                                    │
│ [Undo] [Redo] [Reset] [Save] [Cancel]            │
└────────────────────────────────────────────────────┘
```

#### Features

- ✅ Visual color picker (RGB, HSL, Hex input)
- ✅ 5-bit RGB precision display
- ✅ Color harmony generator (5 schemes)
- ✅ Batch transformations (brighten, saturate, etc.)
- ✅ Undo/redo (50 levels)
- ✅ Palette templates library (50+ presets)
- ✅ Import from image (extract palette)
- ✅ Export as PNG swatch, JSON, Adobe ASE
- ✅ Real-time sprite preview
- ✅ Accessibility checker (contrast ratios)
- ✅ Color blindness simulator

---
