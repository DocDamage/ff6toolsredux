package io

import (
	"fmt"
	"math"
	"time"

	"ffvi_editor/models"
)

// PaletteEditor manages 16-color FF6 palettes with advanced editing features
type PaletteEditor struct {
	palette     *models.Palette
	history     *models.SpriteHistory
	harmonizer  *ColorHarmonizer
	transformer *ColorTransformer
}

// NewPaletteEditor creates a new palette editor
func NewPaletteEditor(palette *models.Palette) *PaletteEditor {
	if palette == nil {
		palette = models.NewPalette("Default")
	}

	return &PaletteEditor{
		palette:     palette.Clone(),
		history:     nil, // Would be initialized with sprite history
		harmonizer:  NewColorHarmonizer(),
		transformer: NewColorTransformer(),
	}
}

// SetColor sets a color at the given index
func (pe *PaletteEditor) SetColor(index int, color models.RGB555) error {
	if index < 0 || index >= 16 {
		return fmt.Errorf("palette index out of range: %d", index)
	}

	// Save snapshot for undo
	pe.saveSnapshot()
	pe.palette.Colors[index] = color
	pe.palette.Modified = time.Now()

	return nil
}

// GetColor gets a color at the given index
func (pe *PaletteEditor) GetColor(index int) (models.RGB555, error) {
	if index < 0 || index >= 16 {
		return models.RGB555{}, fmt.Errorf("palette index out of range: %d", index)
	}
	return pe.palette.Colors[index], nil
}

// GenerateHarmony generates a harmonious palette based on a base color
func (pe *PaletteEditor) GenerateHarmony(baseColor models.RGB555, scheme string) error {
	pe.saveSnapshot()

	newPalette := pe.harmonizer.Generate(baseColor, scheme)
	if newPalette == nil {
		return fmt.Errorf("failed to generate harmony palette for scheme: %s", scheme)
	}

	pe.palette.Colors = newPalette.Colors
	pe.palette.Modified = time.Now()

	return nil
}

// ApplyTransform applies a color transformation to all palette colors
func (pe *PaletteEditor) ApplyTransform(transform string, amount float64) error {
	pe.saveSnapshot()

	for i := 0; i < 16; i++ {
		newColor := pe.transformer.Apply(pe.palette.Colors[i], transform, amount)
		pe.palette.Colors[i] = newColor
	}

	pe.palette.Modified = time.Now()
	return nil
}

// ApplyTransformToRange applies a transformation to a color range
func (pe *PaletteEditor) ApplyTransformToRange(transform string, amount float64, startIdx, endIdx int) error {
	if startIdx < 0 || endIdx >= 16 || startIdx > endIdx {
		return fmt.Errorf("invalid color range: %d-%d", startIdx, endIdx)
	}

	pe.saveSnapshot()

	for i := startIdx; i <= endIdx; i++ {
		newColor := pe.transformer.Apply(pe.palette.Colors[i], transform, amount)
		pe.palette.Colors[i] = newColor
	}

	pe.palette.Modified = time.Now()
	return nil
}

// SwapColors swaps two palette colors
func (pe *PaletteEditor) SwapColors(idx1, idx2 int) error {
	if idx1 < 0 || idx1 >= 16 || idx2 < 0 || idx2 >= 16 {
		return fmt.Errorf("palette index out of range")
	}

	if idx1 == idx2 {
		return nil
	}

	pe.saveSnapshot()
	pe.palette.Colors[idx1], pe.palette.Colors[idx2] = pe.palette.Colors[idx2], pe.palette.Colors[idx1]
	pe.palette.Modified = time.Now()

	return nil
}

// RotateColors rotates palette colors (shifts left)
func (pe *PaletteEditor) RotateColors(startIdx, endIdx, steps int) error {
	if startIdx < 0 || endIdx >= 16 || startIdx > endIdx {
		return fmt.Errorf("invalid color range: %d-%d", startIdx, endIdx)
	}

	pe.saveSnapshot()

	count := endIdx - startIdx + 1
	steps = ((steps % count) + count) % count // Normalize steps

	if steps == 0 {
		return nil
	}

	// Create temporary copy
	temp := make([]models.RGB555, count)
	copy(temp, pe.palette.Colors[startIdx:endIdx+1])

	// Rotate
	for i := 0; i < count; i++ {
		newPos := (i + steps) % count
		pe.palette.Colors[startIdx+newPos] = temp[i]
	}

	pe.palette.Modified = time.Now()
	return nil
}

// GradientFill fills a range with a smooth color gradient
func (pe *PaletteEditor) GradientFill(startIdx, endIdx int, startColor, endColor models.RGB555) error {
	if startIdx < 0 || endIdx >= 16 || startIdx > endIdx {
		return fmt.Errorf("invalid color range: %d-%d", startIdx, endIdx)
	}

	pe.saveSnapshot()

	count := endIdx - startIdx + 1

	for i := 0; i < count; i++ {
		t := float64(i) / float64(count-1)

		// Interpolate between start and end color
		r := uint8(float64(startColor.R)*(1-t) + float64(endColor.R)*t)
		g := uint8(float64(startColor.G)*(1-t) + float64(endColor.G)*t)
		b := uint8(float64(startColor.B)*(1-t) + float64(endColor.B)*t)

		pe.palette.Colors[startIdx+i] = models.RGB555{R: r, G: g, B: b}
	}

	pe.palette.Modified = time.Now()
	return nil
}

// GetPalette returns the current palette
func (pe *PaletteEditor) GetPalette() *models.Palette {
	return pe.palette.Clone()
}

// SetPalette replaces the entire palette
func (pe *PaletteEditor) SetPalette(newPalette *models.Palette) error {
	if newPalette == nil {
		return fmt.Errorf("palette is nil")
	}

	pe.saveSnapshot()
	pe.palette = newPalette.Clone()
	pe.palette.Modified = time.Now()

	return nil
}

// saveSnapshot saves current palette state for undo
func (pe *PaletteEditor) saveSnapshot() {
	// TODO: Implement proper undo/redo with history
}

// ColorHarmonizer generates harmonious color palettes
type ColorHarmonizer struct{}

// NewColorHarmonizer creates a new color harmonizer
func NewColorHarmonizer() *ColorHarmonizer {
	return &ColorHarmonizer{}
}

// HarmonyScheme constants
const (
	HarmonyComplementary      = "complementary"
	HarmonyTriadic            = "triadic"
	HarmonyAnalogous          = "analogous"
	HarmonyMonochromatic      = "monochromatic"
	HarmonySplitComplementary = "split-complementary"
	HarmonyTetradic           = "tetradic"
)

// Generate creates a harmonious palette based on a base color and scheme
func (ch *ColorHarmonizer) Generate(baseColor models.RGB555, scheme string) *models.Palette {
	palette := models.NewPalette("Harmony")

	// Convert to HSL for better color manipulation
	h, s, l := ch.rgb555ToHSL(baseColor)

	switch scheme {
	case HarmonyComplementary:
		ch.generateComplementary(palette, h, s, l)
	case HarmonyTriadic:
		ch.generateTriadic(palette, h, s, l)
	case HarmonyAnalogous:
		ch.generateAnalogous(palette, h, s, l)
	case HarmonyMonochromatic:
		ch.generateMonochromatic(palette, h, s, l)
	case HarmonySplitComplementary:
		ch.generateSplitComplementary(palette, h, s, l)
	case HarmonyTetradic:
		ch.generateTetradic(palette, h, s, l)
	default:
		// Default to monochromatic
		ch.generateMonochromatic(palette, h, s, l)
	}

	return palette
}

// generateComplementary creates complementary color scheme
func (ch *ColorHarmonizer) generateComplementary(palette *models.Palette, h, s, l float64) {
	// Base color + opposite (180 degrees)
	colors := [][3]float64{
		{h, s, l},
		{fmod(h+180, 360), s, l},
		{fmod(h+180, 360), s, l * 0.8},
		{fmod(h+180, 360), s, l * 0.6},
		{h, s * 0.8, l * 0.8},
		{h, s * 0.6, l * 0.6},
		{fmod(h+180, 360), s * 0.8, l * 0.8},
		{fmod(h+180, 360), s * 0.6, l * 0.6},
	}

	for i, hsl := range colors {
		if i < 16 {
			palette.Colors[i] = ch.hslToRGB555(hsl[0], hsl[1], hsl[2])
		}
	}
}

// generateTriadic creates triadic color scheme
func (ch *ColorHarmonizer) generateTriadic(palette *models.Palette, h, s, l float64) {
	// Three colors equally spaced (120 degrees)
	colors := [][3]float64{
		{h, s, l},
		{fmod(h+120, 360), s, l},
		{fmod(h+240, 360), s, l},
		{fmod(h+120, 360), s, l * 0.8},
		{fmod(h+240, 360), s, l * 0.8},
		{h, s * 0.8, l * 0.8},
		{fmod(h+120, 360), s * 0.6, l * 0.6},
		{fmod(h+240, 360), s * 0.6, l * 0.6},
	}

	for i, hsl := range colors {
		if i < 16 {
			palette.Colors[i] = ch.hslToRGB555(hsl[0], hsl[1], hsl[2])
		}
	}
}

// generateAnalogous creates analogous color scheme
func (ch *ColorHarmonizer) generateAnalogous(palette *models.Palette, h, s, l float64) {
	// Colors adjacent on color wheel (±30 degrees)
	colors := [][3]float64{
		{fmod(h-30, 360), s, l},
		{h, s, l},
		{fmod(h+30, 360), s, l},
		{fmod(h-30, 360), s, l * 0.8},
		{fmod(h+30, 360), s, l * 0.8},
		{h, s * 0.8, l * 0.8},
		{fmod(h-30, 360), s * 0.6, l * 0.6},
		{fmod(h+30, 360), s * 0.6, l * 0.6},
	}

	for i, hsl := range colors {
		if i < 16 {
			palette.Colors[i] = ch.hslToRGB555(hsl[0], hsl[1], hsl[2])
		}
	}
}

// generateMonochromatic creates monochromatic color scheme
func (ch *ColorHarmonizer) generateMonochromatic(palette *models.Palette, h, s, l float64) {
	// Same hue, different lightness levels
	for i := 0; i < 16; i++ {
		lightness := l * (float64(i)/15.0*0.9 + 0.1) // 10% to 100% of original lightness
		palette.Colors[i] = ch.hslToRGB555(h, s, lightness)
	}
}

// generateSplitComplementary creates split-complementary scheme
func (ch *ColorHarmonizer) generateSplitComplementary(palette *models.Palette, h, s, l float64) {
	// Base + two adjacent to complement
	colors := [][3]float64{
		{h, s, l},
		{fmod(h+150, 360), s, l},
		{fmod(h+210, 360), s, l},
		{fmod(h+150, 360), s, l * 0.8},
		{fmod(h+210, 360), s, l * 0.8},
		{h, s * 0.8, l * 0.8},
		{fmod(h+150, 360), s * 0.6, l * 0.6},
		{fmod(h+210, 360), s * 0.6, l * 0.6},
	}

	for i, hsl := range colors {
		if i < 16 {
			palette.Colors[i] = ch.hslToRGB555(hsl[0], hsl[1], hsl[2])
		}
	}
}

// generateTetradic creates tetradic (square) color scheme
func (ch *ColorHarmonizer) generateTetradic(palette *models.Palette, h, s, l float64) {
	// Four colors equally spaced (90 degrees)
	colors := [][3]float64{
		{h, s, l},
		{fmod(h+90, 360), s, l},
		{fmod(h+180, 360), s, l},
		{fmod(h+270, 360), s, l},
		{fmod(h+45, 360), s, l * 0.8},
		{fmod(h+135, 360), s, l * 0.8},
		{fmod(h+225, 360), s, l * 0.8},
		{fmod(h+315, 360), s, l * 0.8},
	}

	for i, hsl := range colors {
		if i < 16 {
			palette.Colors[i] = ch.hslToRGB555(hsl[0], hsl[1], hsl[2])
		}
	}
}

// rgb555ToHSL converts 5-bit RGB to HSL
func (ch *ColorHarmonizer) rgb555ToHSL(rgb models.RGB555) (h, s, l float64) {
	r := float64(rgb.R) / 31.0
	g := float64(rgb.G) / 31.0
	b := float64(rgb.B) / 31.0

	maxc := math.Max(math.Max(r, g), b)
	minc := math.Min(math.Min(r, g), b)
	l = (maxc + minc) / 2.0

	if minc == maxc {
		h = 0
		s = 0
	} else {
		if l <= 0.5 {
			s = (maxc - minc) / (maxc + minc)
		} else {
			s = (maxc - minc) / (2.0 - maxc - minc)
		}

		rc := (maxc - r) / (maxc - minc)
		gc := (maxc - g) / (maxc - minc)
		bc := (maxc - b) / (maxc - minc)

		if r == maxc {
			h = bc - gc
		} else if g == maxc {
			h = 2.0 + rc - bc
		} else {
			h = 4.0 + gc - rc
		}

		h *= 60.0
		if h < 0 {
			h += 360.0
		}
	}

	return
}

// hslToRGB555 converts HSL to 5-bit RGB
func (ch *ColorHarmonizer) hslToRGB555(h, s, l float64) models.RGB555 {
	var r, g, b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		var t1, t2 float64
		if l < 0.5 {
			t2 = l * (1.0 + s)
		} else {
			t2 = l + s - l*s
		}
		t1 = 2.0*l - t2

		hNorm := h / 360.0
		r = ch.hueToRGB(t1, t2, fmod(hNorm+1.0/3.0, 1.0))
		g = ch.hueToRGB(t1, t2, fmod(hNorm, 1.0))
		b = ch.hueToRGB(t1, t2, fmod(hNorm-1.0/3.0, 1.0))
	}

	// Clamp to 0-1
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	// Convert to 5-bit
	return models.RGB555{
		R: uint8(r * 31),
		G: uint8(g * 31),
		B: uint8(b * 31),
	}
}

// hueToRGB helper for HSL conversion
func (ch *ColorHarmonizer) hueToRGB(t1, t2, hue float64) float64 {
	if hue < 0 {
		hue += 1.0
	}
	if hue > 1 {
		hue -= 1.0
	}
	if hue < 1.0/6.0 {
		return t1 + (t2-t1)*6.0*hue
	}
	if hue < 1.0/2.0 {
		return t2
	}
	if hue < 2.0/3.0 {
		return t1 + (t2-t1)*(2.0/3.0-hue)*6.0
	}
	return t1
}

// ColorTransformer applies transformations to colors
type ColorTransformer struct{}

// NewColorTransformer creates a new transformer
func NewColorTransformer() *ColorTransformer {
	return &ColorTransformer{}
}

// Apply applies a transformation to a color
func (ct *ColorTransformer) Apply(color models.RGB555, transform string, amount float64) models.RGB555 {
	harmonizer := &ColorHarmonizer{}
	h, s, l := harmonizer.rgb555ToHSL(color)

	// Clamp amount to 0-1
	amount = math.Max(0, math.Min(1, amount))

	switch transform {
	case "brighten":
		l = math.Min(1, l+amount*0.3)
	case "darken":
		l = math.Max(0, l-amount*0.3)
	case "saturate":
		s = math.Min(1, s+amount*0.3)
	case "desaturate":
		s = math.Max(0, s-amount*0.3)
	case "shift-hue":
		h = fmod(h+amount*360, 360)
	case "invert":
		r := float64(color.R) / 31.0
		g := float64(color.G) / 31.0
		b := float64(color.B) / 31.0
		return models.RGB555{
			R: uint8((1 - r) * 31),
			G: uint8((1 - g) * 31),
			B: uint8((1 - b) * 31),
		}
	case "grayscale":
		gray := uint8((float64(color.R) + float64(color.G) + float64(color.B)) / 3)
		return models.RGB555{R: gray, G: gray, B: gray}
	case "sepia":
		r := float64(color.R) / 31.0
		g := float64(color.G) / 31.0
		b := float64(color.B) / 31.0
		return models.RGB555{
			R: uint8(math.Min(31, (r*0.393+g*0.769+b*0.189)) * 31),
			G: uint8(math.Min(31, (r*0.349+g*0.686+b*0.168)) * 31),
			B: uint8(math.Min(31, (r*0.272+g*0.534+b*0.131)) * 31),
		}
	}

	return harmonizer.hslToRGB555(h, s, l)
}

// Helper function: fmod for floating point modulo
func fmod(x, y float64) float64 {
	return math.Mod(x, y)
}
