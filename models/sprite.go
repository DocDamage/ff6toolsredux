package models

import (
	"image"
	"image/color"
	"time"
)

// SpriteType defines the classification of sprites in FF6
type SpriteType int

const (
	// Character sprites (16x24 pixels, walking sprites)
	SpriteTypeCharacter SpriteType = iota
	// Battle sprites (32x32 pixels)
	SpriteTypeBattle
	// Portrait sprites (48x48 pixels, menu/dialogue)
	SpriteTypePortrait
	// NPC sprites (16x24 pixels, world map NPCs)
	SpriteTypeNPC
	// Enemy sprites (32x32 pixels, battle enemies)
	SpriteTypeEnemy
	// Overworld sprite (16x16 pixels, map sprites)
	SpriteTypeOverworld
)

// String returns the name of the sprite type
func (st SpriteType) String() string {
	switch st {
	case SpriteTypeCharacter:
		return "Character"
	case SpriteTypeBattle:
		return "Battle"
	case SpriteTypePortrait:
		return "Portrait"
	case SpriteTypeNPC:
		return "NPC"
	case SpriteTypeEnemy:
		return "Enemy"
	case SpriteTypeOverworld:
		return "Overworld"
	default:
		return "Unknown"
	}
}

// GetDimensions returns width and height for this sprite type
func (st SpriteType) GetDimensions() (width, height int) {
	switch st {
	case SpriteTypeCharacter, SpriteTypeNPC:
		return 16, 24
	case SpriteTypeBattle, SpriteTypeEnemy:
		return 32, 32
	case SpriteTypePortrait:
		return 48, 48
	case SpriteTypeOverworld:
		return 16, 16
	default:
		return 0, 0
	}
}

// RGB555 represents a 5-bit RGB color (FF6 format: 0-31 per channel)
// SNES uses 5-bit RGB, not standard 8-bit RGB
type RGB555 struct {
	R uint8 // 0-31 (5 bits)
	G uint8 // 0-31 (5 bits)
	B uint8 // 0-31 (5 bits)
}

// ToRGB888 converts 5-bit RGB to standard 8-bit RGB
func (c RGB555) ToRGB888() (r, g, b uint8) {
	// Scale from 0-31 to 0-255
	r = uint8((int(c.R) * 255) / 31)
	g = uint8((int(c.G) * 255) / 31)
	b = uint8((int(c.B) * 255) / 31)
	return
}

// FromRGB888 converts standard 8-bit RGB to 5-bit RGB
func FromRGB888(r, g, b uint8) RGB555 {
	// Scale from 0-255 to 0-31
	return RGB555{
		R: uint8((int(r) * 31) / 255),
		G: uint8((int(g) * 31) / 255),
		B: uint8((int(b) * 31) / 255),
	}
}

// ToColor converts to Go's color.Color interface
func (c RGB555) ToColor() color.Color {
	r, g, b := c.ToRGB888()
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// Palette represents a 16-color FF6 sprite palette
type Palette struct {
	// 16 colors, index 0 is transparent
	Colors   [16]RGB555
	Name     string
	Author   string
	Tags     []string
	License  string
	Created  time.Time
	Modified time.Time
	Notes    string
}

// NewPalette creates a new palette with default black colors
func NewPalette(name string) *Palette {
	return &Palette{
		Name:    name,
		Created: time.Now(),
		Colors:  [16]RGB555{}, // All black by default
	}
}

// Clone creates a deep copy of the palette
func (p *Palette) Clone() *Palette {
	cloned := &Palette{
		Name:     p.Name,
		Author:   p.Author,
		License:  p.License,
		Created:  p.Created,
		Modified: p.Modified,
		Notes:    p.Notes,
		Colors:   p.Colors, // Arrays are copied by value
	}
	cloned.Tags = make([]string, len(p.Tags))
	copy(cloned.Tags, p.Tags)
	return cloned
}

// FF6Sprite represents a sprite that can be applied to FF6 save files
type FF6Sprite struct {
	// Basic properties
	ID          string // Unique identifier (e.g., "terra_custom_001")
	Name        string // Display name (e.g., "Fire Terra")
	Type        SpriteType
	Description string

	// Sprite data
	Data      []byte   // Compressed sprite data
	Palette   *Palette // 16-color palette
	Frames    int      // Number of animation frames
	FrameRate int      // Frames per second

	// Animation metadata
	FrameDurations []int // Duration of each frame in ms

	// Source tracking
	SourceFile   string
	ImportedFrom string // PNG, GIF, Aseprite, etc.
	ImportDate   time.Time

	// Validation & metadata
	Width        int
	Height       int
	IsCompressed bool
	Checksum     string // For integrity validation

	// Attribution
	Author       string
	License      string
	Tags         []string
	CreatedDate  time.Time
	ModifiedDate time.Time
}

// NewSprite creates a new sprite
func NewSprite(id, name string, spriteType SpriteType) *FF6Sprite {
	w, h := spriteType.GetDimensions()
	return &FF6Sprite{
		ID:           id,
		Name:         name,
		Type:         spriteType,
		Width:        w,
		Height:       h,
		Frames:       1,
		FrameRate:    12,
		Palette:      NewPalette("Default"),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
}

// Clone creates a deep copy of the sprite
func (s *FF6Sprite) Clone() *FF6Sprite {
	cloned := &FF6Sprite{
		ID:           s.ID,
		Name:         s.Name,
		Type:         s.Type,
		Description:  s.Description,
		Data:         make([]byte, len(s.Data)),
		Frames:       s.Frames,
		FrameRate:    s.FrameRate,
		SourceFile:   s.SourceFile,
		ImportedFrom: s.ImportedFrom,
		ImportDate:   s.ImportDate,
		Width:        s.Width,
		Height:       s.Height,
		IsCompressed: s.IsCompressed,
		Checksum:     s.Checksum,
		Author:       s.Author,
		License:      s.License,
		CreatedDate:  s.CreatedDate,
		ModifiedDate: s.ModifiedDate,
	}

	copy(cloned.Data, s.Data)

	if s.Palette != nil {
		cloned.Palette = s.Palette.Clone()
	}

	cloned.Tags = make([]string, len(s.Tags))
	copy(cloned.Tags, s.Tags)

	cloned.FrameDurations = make([]int, len(s.FrameDurations))
	copy(cloned.FrameDurations, s.FrameDurations)

	return cloned
}

// GetFrameCount returns the number of frames
func (s *FF6Sprite) GetFrameCount() int {
	if s.Frames <= 0 {
		return 1
	}
	return s.Frames
}

// GetFrameSize returns the size in bytes for one frame
func (s *FF6Sprite) GetFrameSize() int {
	if s.Frames <= 0 || len(s.Data) == 0 {
		return 0
	}
	return len(s.Data) / s.Frames
}

// GetExpectedDataSize calculates expected uncompressed data size
// Based on tile count (8x8 tiles with 4-bit color = 32 bytes per tile)
func (s *FF6Sprite) GetExpectedDataSize() int {
	tilesWide := (s.Width + 7) / 8  // Round up to nearest tile
	tilesHigh := (s.Height + 7) / 8 // Round up to nearest tile
	bytesPerTile := 32              // 64 pixels × 4-bit = 32 bytes
	return tilesWide * tilesHigh * bytesPerTile * s.Frames
}

// SpriteFrame represents a single animation frame
type SpriteFrame struct {
	Index    int          // Frame number (0-based)
	Data     []byte       // Raw frame data
	Duration int          // Duration in milliseconds
	Image    *image.Image // Decoded image for preview
}

// PaletteSnapshot represents a palette state for undo/redo
type PaletteSnapshot struct {
	Timestamp time.Time
	Palette   *Palette
}

// SpriteHistory tracks modifications for undo/redo
type SpriteHistory struct {
	Snapshots    []*FF6Sprite
	CurrentIdx   int
	MaxSnapshots int
}

// NewSpriteHistory creates a new history tracker
func NewSpriteHistory(maxSnapshots int) *SpriteHistory {
	return &SpriteHistory{
		Snapshots:    make([]*FF6Sprite, 0, maxSnapshots),
		MaxSnapshots: maxSnapshots,
		CurrentIdx:   -1,
	}
}

// Push adds a sprite snapshot to history
func (h *SpriteHistory) Push(sprite *FF6Sprite) {
	// Remove any redo history after current position
	if h.CurrentIdx >= 0 && h.CurrentIdx < len(h.Snapshots)-1 {
		h.Snapshots = h.Snapshots[:h.CurrentIdx+1]
	}

	h.Snapshots = append(h.Snapshots, sprite.Clone())
	h.CurrentIdx++

	// Enforce max snapshots
	if len(h.Snapshots) > h.MaxSnapshots {
		h.Snapshots = h.Snapshots[len(h.Snapshots)-h.MaxSnapshots:]
		h.CurrentIdx = len(h.Snapshots) - 1
	}
}

// Undo returns the previous sprite state
func (h *SpriteHistory) Undo() *FF6Sprite {
	if h.CurrentIdx <= 0 {
		return nil
	}
	h.CurrentIdx--
	return h.Snapshots[h.CurrentIdx].Clone()
}

// Redo returns the next sprite state
func (h *SpriteHistory) Redo() *FF6Sprite {
	if h.CurrentIdx >= len(h.Snapshots)-1 {
		return nil
	}
	h.CurrentIdx++
	return h.Snapshots[h.CurrentIdx].Clone()
}

// CanUndo checks if undo is available
func (h *SpriteHistory) CanUndo() bool {
	return h.CurrentIdx > 0
}

// CanRedo checks if redo is available
func (h *SpriteHistory) CanRedo() bool {
	return h.CurrentIdx < len(h.Snapshots)-1
}

// SpriteImportOptions configures sprite import behavior
type SpriteImportOptions struct {
	// Quantization options
	MaxColors       int    // Reduce to this many colors (default 16)
	DitherMethod    string // "floyd-steinberg", "bayer", "none" (default "floyd-steinberg")
	PreservePalette bool   // Keep original palette if possible

	// Processing options
	RemoveBackground bool // Try to detect and remove background
	AutoDetectType   bool // Auto-detect sprite type from dimensions
	AutoPadding      bool // Pad to nearest tile boundary

	// Import source
	SourceFormat string // "png", "gif", "bmp", "jpg"
	SourcePath   string

	// Target settings
	TargetType    SpriteType
	CharacterID   int      // For character sprites
	TargetPalette *Palette // Use existing palette
}

// NewSpriteImportOptions creates default import options
func NewSpriteImportOptions() *SpriteImportOptions {
	return &SpriteImportOptions{
		MaxColors:        16,
		DitherMethod:     "floyd-steinberg",
		PreservePalette:  false,
		RemoveBackground: false,
		AutoDetectType:   true,
		AutoPadding:      true,
	}
}
