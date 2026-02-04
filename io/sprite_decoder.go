package io

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	"ffvi_editor/models"
)

// ImageDecoder handles decoding of various image formats
type ImageDecoder struct {
	supportedFormats map[string]bool
}

// NewImageDecoder creates a new image decoder
func NewImageDecoder() *ImageDecoder {
	return &ImageDecoder{
		supportedFormats: map[string]bool{
			".png":  true,
			".bmp":  true,
			".gif":  true,
			".jpg":  true,
			".jpeg": true,
		},
	}
}

// Decode reads and decodes an image file
// Returns the image, format string, and any error
func (d *ImageDecoder) Decode(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Detect format by file extension
	ext := filepath.Ext(filePath)
	if ext == "" {
		return nil, "", fmt.Errorf("file has no extension")
	}

	switch ext {
	case ".png":
		img, err := png.Decode(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode PNG: %w", err)
		}
		return img, "png", nil

	case ".gif":
		img, err := gif.Decode(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode GIF: %w", err)
		}
		return img, "gif", nil

	case ".bmp":
		img, _, err := image.Decode(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode BMP: %w", err)
		}
		return img, "bmp", nil

	case ".jpg", ".jpeg":
		img, _, err := image.Decode(file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode JPEG: %w", err)
		}
		return img, "jpeg", nil

	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}
}

// PaletteExtractor extracts a 16-color palette from an image
type PaletteExtractor struct {
	maxColors int
}

// NewPaletteExtractor creates a new palette extractor
func NewPaletteExtractor() *PaletteExtractor {
	return &PaletteExtractor{
		maxColors: 16,
	}
}

// ColorInfo tracks color frequency in image
type ColorInfo struct {
	RGBA   color.Color
	Count  int
	RGB555 models.RGB555
	Index  int
}

// Extract analyzes an image and extracts a 16-color palette
func (pe *PaletteExtractor) Extract(img image.Image, maxColors int) (*models.Palette, error) {
	if maxColors > 256 {
		maxColors = 256
	}
	if maxColors < 2 {
		maxColors = 2
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Sample every Nth pixel for large images to speed up extraction
	sampleRate := 1
	if width*height > 1000000 { // > 1 megapixel
		sampleRate = 4
	} else if width*height > 100000 {
		sampleRate = 2
	}

	// Count color frequencies
	colorCounts := make(map[uint32]*ColorInfo)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += sampleRate {
		for x := bounds.Min.X; x < bounds.Max.X; x += sampleRate {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()

			// Normalize to 8-bit
			r, g, b, a = r>>8, g>>8, b>>8, a>>8

			// Skip transparent pixels (alpha < 128)
			if a < 128 {
				continue
			}

			// Convert to uint32 for map key (R, G, B, A)
			key := (uint32(r) << 24) | (uint32(g) << 16) | (uint32(b) << 8) | uint32(a)

			if info, exists := colorCounts[key]; exists {
				info.Count++
			} else {
				colorCounts[key] = &ColorInfo{
					RGBA:   color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
					Count:  1,
					RGB555: models.FromRGB888(uint8(r), uint8(g), uint8(b)),
				}
			}
		}
	}

	// Sort colors by frequency
	colors := make([]*ColorInfo, 0, len(colorCounts))
	for _, info := range colorCounts {
		colors = append(colors, info)
	}

	sort.Slice(colors, func(i, j int) bool {
		return colors[i].Count > colors[j].Count
	})

	// Take top N colors
	if len(colors) > maxColors {
		colors = colors[:maxColors]
	}

	// Create palette
	palette := models.NewPalette("Extracted")
	for i, info := range colors {
		if i < 16 {
			palette.Colors[i] = info.RGB555
		}
	}

	// If we have fewer than 16 colors, fill with black
	for i := len(colors); i < 16; i++ {
		palette.Colors[i] = models.RGB555{R: 0, G: 0, B: 0}
	}

	return palette, nil
}

// ColorQuantizer reduces image colors to a target palette
type ColorQuantizer struct {
	ditherMethod string // "floyd-steinberg", "bayer", "none"
}

// NewColorQuantizer creates a new quantizer
func NewColorQuantizer(ditherMethod string) *ColorQuantizer {
	if ditherMethod == "" {
		ditherMethod = "floyd-steinberg"
	}
	return &ColorQuantizer{
		ditherMethod: ditherMethod,
	}
}

// Quantize reduces image to target number of colors
func (cq *ColorQuantizer) Quantize(img image.Image, maxColors int) (image.Image, *models.Palette, error) {
	// Extract palette
	extractor := NewPaletteExtractor()
	palette, err := extractor.Extract(img, maxColors)
	if err != nil {
		return nil, nil, err
	}

	// Apply dithering if needed
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)
	draw.Draw(output, bounds, img, bounds.Min, draw.Src)

	// Apply dithering
	switch cq.ditherMethod {
	case "floyd-steinberg":
		cq.applyFloydSteinbergDithering(output, palette)
	case "bayer":
		cq.applyBayerDithering(output, palette)
	case "none":
		// No dithering, just quantize
		cq.quantizeNoised(output, palette)
	}

	return output, palette, nil
}

// quantizeNoised maps image colors to nearest palette colors
func (cq *ColorQuantizer) quantizeNoised(img *image.RGBA, palette *models.Palette) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8

			// Find nearest palette color
			idx := findNearestColor(uint8(r), uint8(g), uint8(b), palette)
			paletteColor := palette.Colors[idx].ToColor()
			img.Set(x, y, paletteColor)
		}
	}
}

// applyFloydSteinbergDithering applies Floyd-Steinberg dithering
func (cq *ColorQuantizer) applyFloydSteinbergDithering(img *image.RGBA, palette *models.Palette) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Pre-allocate error buffers for R, G, B channels (outside loop for efficiency)
	errR := make([][]float64, height+1)
	errG := make([][]float64, height+1)
	errB := make([][]float64, height+1)

	// Pre-allocate all rows at once instead of in a loop
	errorLine := make([]float64, (height+1)*(width+1))
	for i := range errR {
		errR[i] = errorLine[i*(width+1) : (i+1)*(width+1)]
		errG[i] = make([]float64, width+1)
		errB[i] = make([]float64, width+1)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			idx := y - bounds.Min.Y
			idxX := x - bounds.Min.X

			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			r, g, b = r>>8, g>>8, b>>8

			// Apply accumulated error
			r8 := uint8(max(0, min(255, int(r)+int(errR[idx][idxX]))))
			g8 := uint8(max(0, min(255, int(g)+int(errG[idx][idxX]))))
			b8 := uint8(max(0, min(255, int(b)+int(errB[idx][idxX]))))

			// Find nearest palette color
			paletteIdx := findNearestColor(r8, g8, b8, palette)
			paletteColor := palette.Colors[paletteIdx]
			pr, pg, pb := paletteColor.ToRGB888()

			img.Set(x, y, color.RGBA{R: pr, G: pg, B: pb, A: 255})

			// Calculate and distribute error
			errRVal := float64(r8) - float64(pr)
			errGVal := float64(g8) - float64(pg)
			errBVal := float64(b8) - float64(pb)

			// Distribute error to neighboring pixels
			if idxX+1 < width {
				errR[idx][idxX+1] += errRVal * 7.0 / 16.0
				errG[idx][idxX+1] += errGVal * 7.0 / 16.0
				errB[idx][idxX+1] += errBVal * 7.0 / 16.0
			}

			if idx+1 < height {
				if idxX > 0 {
					errR[idx+1][idxX-1] += errRVal * 3.0 / 16.0
					errG[idx+1][idxX-1] += errGVal * 3.0 / 16.0
					errB[idx+1][idxX-1] += errBVal * 3.0 / 16.0
				}

				errR[idx+1][idxX] += errRVal * 5.0 / 16.0
				errG[idx+1][idxX] += errGVal * 5.0 / 16.0
				errB[idx+1][idxX] += errBVal * 5.0 / 16.0

				if idxX+1 < width {
					errR[idx+1][idxX+1] += errRVal * 1.0 / 16.0
					errG[idx+1][idxX+1] += errGVal * 1.0 / 16.0
					errB[idx+1][idxX+1] += errBVal * 1.0 / 16.0
				}
			}
		}
	}
}

// applyBayerDithering applies Bayer matrix dithering
func (cq *ColorQuantizer) applyBayerDithering(img *image.RGBA, palette *models.Palette) {
	// Simplified Bayer dithering - could be enhanced
	bounds := img.Bounds()

	bayerMatrix := [4][4]int{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	}

	threshold := 255 / 16

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Get dither value from Bayer matrix
			ditherVal := int(bayerMatrix[y%4][x%4]) * threshold / 15

			// Apply dither threshold
			if int(r8) > ditherVal {
				r8 = uint8(min(255, int(r8)+ditherVal))
			}
			if int(g8) > ditherVal {
				g8 = uint8(min(255, int(g8)+ditherVal))
			}
			if int(b8) > ditherVal {
				b8 = uint8(min(255, int(b8)+ditherVal))
			}

			// Find nearest palette color
			paletteIdx := findNearestColor(r8, g8, b8, palette)
			paletteColor := palette.Colors[paletteIdx]
			pr, pg, pb := paletteColor.ToRGB888()

			img.Set(x, y, color.RGBA{R: pr, G: pg, B: pb, A: 255})
		}
	}
}

// findNearestColor finds the closest palette color to the given RGB value
func findNearestColor(r, g, b uint8, palette *models.Palette) int {
	minDist := float64(999999)
	bestIdx := 0

	for i := 0; i < 16; i++ {
		palColor := palette.Colors[i]
		pr, pg, pb := palColor.ToRGB888()

		// Euclidean distance in color space
		dist := float64((int(r)-int(pr))*(int(r)-int(pr))) +
			float64((int(g)-int(pg))*(int(g)-int(pg))) +
			float64((int(b)-int(pb))*(int(b)-int(pb)))

		if dist < minDist {
			minDist = dist
			bestIdx = i
		}
	}

	return bestIdx
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
