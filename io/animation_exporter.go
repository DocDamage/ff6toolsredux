package io

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"

	"ffvi_editor/models"
)

// AnimationExporter handles export of animations in various formats
type AnimationExporter struct {
	animation *models.AnimationData
	sequencer *FrameSequencer
}

// AnimationExportOptions configures export behavior
type AnimationExportOptions struct {
	Quality       int  // 1-100 for GIF compression
	Scale         int  // Scale factor for output
	Dither        bool // Apply dithering
	BackgroundRGB [3]uint8
}

// NewAnimationExporter creates a new animation exporter
func NewAnimationExporter(animation *models.AnimationData) *AnimationExporter {
	if animation == nil {
		return nil
	}

	return &AnimationExporter{
		animation: animation,
	}
}

// decodeSpriteToRGBA decodes a tile-encoded FF6 sprite into RGBA for export.
func decodeSpriteToRGBA(sprite *models.FF6Sprite) (*image.RGBA, error) {
	if sprite == nil {
		return nil, fmt.Errorf("sprite is nil")
	}

	conv := NewFF6SpriteConverter()
	img, err := conv.DecodeFF6Sprite(sprite)
	if err != nil {
		return nil, err
	}

	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}

	b := img.Bounds()
	out := image.NewRGBA(b)
	draw.Draw(out, b, img, b.Min, draw.Src)
	return out, nil
}

// ExportGIF exports animation as animated GIF
func (ae *AnimationExporter) ExportGIF(filePath string, opts *AnimationExportOptions) error {
	if ae.animation == nil || len(ae.animation.Frames) == 0 {
		return fmt.Errorf("animation is empty")
	}

	if opts == nil {
		opts = &AnimationExportOptions{
			Quality: 75,
			Scale:   1,
		}
	}

	if opts.Scale < 1 {
		opts.Scale = 1
	}

	// Create output images
	var images []*image.Paletted
	var delays []int

	for i, frame := range ae.animation.Frames {
		if frame == nil {
			continue
		}

		// Decode sprite to RGBA
		rgba, err := decodeSpriteToRGBA(frame)
		if err != nil {
			return fmt.Errorf("failed to decode frame %d: %w", i, err)
		}

		img := rgbaToPaletted(rgba, frame.Palette)

		// Scale if needed
		if opts.Scale > 1 {
			img = scaleImagePaletted(img, opts.Scale)
		}

		images = append(images, img)

		// Get frame delay
		if i < len(ae.animation.FrameTimings) {
			delay := ae.animation.FrameTimings[i] / 10 // ms to centiseconds
			if delay < 1 {
				delay = 1
			}
			delays = append(delays, delay)
		}
	}

	if len(images) == 0 {
		return fmt.Errorf("no valid frames to export")
	}

	// Create and save GIF
	g := &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: -1,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return gif.EncodeAll(file, g)
}

// ExportJSON exports animation metadata as JSON
func (ae *AnimationExporter) ExportJSON(filePath string) error {
	if ae.animation == nil {
		return fmt.Errorf("animation is nil")
	}

	exportData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":          ae.animation.Metadata.Name,
			"description":   ae.animation.Metadata.Description,
			"author":        ae.animation.Metadata.Author,
			"totalDuration": ae.animation.Metadata.TotalDuration,
			"frameCount":    ae.animation.Metadata.FrameCount,
		},
		"settings": map[string]interface{}{
			"playbackMode": ae.animation.PlaybackMode.String(),
			"defaultSpeed": ae.animation.DefaultSpeed,
		},
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// ExportFramesPNG exports individual frames as PNG files
func (ae *AnimationExporter) ExportFramesPNG(outputDir string, opts *AnimationExportOptions) error {
	if ae.animation == nil || len(ae.animation.Frames) == 0 {
		return fmt.Errorf("animation is empty")
	}

	if opts == nil {
		opts = &AnimationExportOptions{
			Scale: 1,
		}
	}

	if opts.Scale < 1 {
		opts.Scale = 1
	}

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for i, frame := range ae.animation.Frames {
		if frame == nil {
			continue
		}

		rgba, err := decodeSpriteToRGBA(frame)
		if err != nil {
			return fmt.Errorf("failed to decode frame %d: %w", i, err)
		}
		if opts.Scale > 1 {
			rgba = scaleImageRGBA(rgba, opts.Scale)
		}

		fileName := fmt.Sprintf("frame_%03d.png", i)
		filePath := filepath.Join(outputDir, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", fileName, err)
		}
		defer file.Close()

		err = png.Encode(file, rgba)
		if err != nil {
			return fmt.Errorf("failed to encode frame %d: %w", i, err)
		}
	}

	return nil
}

// Helper functions

func rgbaToPaletted(src *image.RGBA, palette *models.Palette) *image.Paletted {
	if src == nil || palette == nil {
		return nil
	}

	bounds := src.Bounds()
	out := image.NewPaletted(bounds, convertPalette(palette))
	draw.FloydSteinberg.Draw(out, bounds, src, bounds.Min)
	return out
}

func convertPalette(palette *models.Palette) color.Palette {
	colors := make([]color.Color, 16)
	for i := 0; i < 16; i++ {
		c := palette.Colors[i]
		r, g, b := c.ToRGB888()
		colors[i] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return colors
}

func scaleImagePaletted(img *image.Paletted, scale int) *image.Paletted {
	if scale <= 1 {
		return img
	}

	bounds := img.Bounds()
	newBounds := image.Rect(0, 0, bounds.Dx()*scale, bounds.Dy()*scale)
	scaled := image.NewPaletted(newBounds, img.Palette)

	for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
		for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
			srcX := bounds.Min.X + (x-newBounds.Min.X)/scale
			srcY := bounds.Min.Y + (y-newBounds.Min.Y)/scale
			if srcX < bounds.Max.X && srcY < bounds.Max.Y {
				scaled.SetColorIndex(x, y, img.ColorIndexAt(srcX, srcY))
			}
		}
	}

	return scaled
}

func palettedToRGBA(img *image.Paletted) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

func scaleImageRGBA(img *image.RGBA, scale int) *image.RGBA {
	if scale <= 1 {
		return img
	}

	bounds := img.Bounds()
	newBounds := image.Rect(0, 0, bounds.Dx()*scale, bounds.Dy()*scale)
	scaled := image.NewRGBA(newBounds)

	for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
		for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
			srcX := bounds.Min.X + (x-newBounds.Min.X)/scale
			srcY := bounds.Min.Y + (y-newBounds.Min.Y)/scale
			if srcX < bounds.Max.X && srcY < bounds.Max.Y {
				scaled.Set(x, y, img.At(srcX, srcY))
			}
		}
	}

	return scaled
}
