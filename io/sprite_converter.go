package io

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"ffvi_editor/models"
)

// FF6SpriteConverter converts images to FF6 sprite format
type FF6SpriteConverter struct {
	dithering bool
}

// NewFF6SpriteConverter creates a new converter
func NewFF6SpriteConverter() *FF6SpriteConverter {
	return &FF6SpriteConverter{
		dithering: true,
	}
}

// ToFF6Format converts an image to FF6 sprite format
// Handles quantization to 16-color palette and tile-based encoding
func (c *FF6SpriteConverter) ToFF6Format(img image.Image, palette *models.Palette, spriteType models.SpriteType) *models.FF6Sprite {
	sprite := models.NewSprite("custom_sprite", "Custom Sprite", spriteType)
	sprite.Palette = palette

	// Get target dimensions
	expectedWidth, expectedHeight := spriteType.GetDimensions()

	// Resize/pad image to match expected dimensions
	resizedImg := c.fitImage(img, expectedWidth, expectedHeight)

	// Convert image to indexed color (palette indices)
	indexedData := c.imageToIndexed(resizedImg, palette)

	// Encode to FF6 tile format
	encodedData := c.encodeToTiles(indexedData, expectedWidth, expectedHeight)

	sprite.Data = encodedData
	sprite.Width = expectedWidth
	sprite.Height = expectedHeight
	sprite.IsCompressed = false

	return sprite
}

// fitImage resizes or pads image to target dimensions (optimized with draw.Draw)
func (c *FF6SpriteConverter) fitImage(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// Create new image with target size
	canvas := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Calculate offset to center the image
	offsetX := (targetWidth - srcWidth) / 2
	offsetY := (targetHeight - srcHeight) / 2

	// Ensure non-negative offsets
	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}

	// Use draw.Draw for efficient image copying (much faster than manual Set calls)
	dstRect := image.Rect(offsetX, offsetY, offsetX+srcWidth, offsetY+srcHeight)
	draw.Draw(canvas, dstRect, img, bounds.Min, draw.Src)

	return canvas
}

// imageToIndexed converts image to palette indices
func (c *FF6SpriteConverter) imageToIndexed(img image.Image, palette *models.Palette) []byte {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// 4 bits per pixel (16 colors) = 2 pixels per byte
	dataSize := (width * height) / 2
	if (width*height)%2 != 0 {
		dataSize++
	}

	indexedData := make([]byte, dataSize)

	pixelIdx := 0
	byteIdx := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get pixel color
			pixelColor := img.At(x, y)
			r, g, b, a := pixelColor.RGBA()

			// Convert to 8-bit
			r, g, b, a = r>>8, g>>8, b>>8, a>>8

			// Determine palette index
			var palIdx uint8
			if a < 128 {
				// Transparent pixel -> index 0
				palIdx = 0
			} else {
				// Find nearest color
				palIdx = uint8(findNearestColor(uint8(r), uint8(g), uint8(b), palette))
			}

			// Pack two pixels per byte (4 bits each)
			if pixelIdx%2 == 0 {
				// Low nibble
				indexedData[byteIdx] = (indexedData[byteIdx] & 0xF0) | (palIdx & 0x0F)
			} else {
				// High nibble
				indexedData[byteIdx] = (indexedData[byteIdx] & 0x0F) | ((palIdx & 0x0F) << 4)
				byteIdx++
			}

			pixelIdx++
		}
	}

	return indexedData
}

// encodeToTiles encodes pixel data into FF6 tile format
// FF6 uses 8x8 pixel tiles, with each tile being 32 bytes (4-bit indexed color)
func (c *FF6SpriteConverter) encodeToTiles(indexedData []byte, width, height int) []byte {
	var buffer bytes.Buffer

	// Calculate tile grid
	tilesWide := (width + 7) / 8
	tilesHigh := (height + 7) / 8

	// Process each tile in row-major order
	for tileRow := 0; tileRow < tilesHigh; tileRow++ {
		for tileCol := 0; tileCol < tilesWide; tileCol++ {
			// Encode this tile
			tileData := c.encodeTile(indexedData, width, height, tileCol, tileRow)
			buffer.Write(tileData)
		}
	}

	return buffer.Bytes()
}

// encodeTile encodes a single 8x8 tile
// Each row is 4 bytes (8 pixels × 4 bits = 32 bits = 4 bytes)
func (c *FF6SpriteConverter) encodeTile(indexedData []byte, imgWidth, imgHeight int, tileCol, tileRow int) []byte {
	tile := make([]byte, 32) // 8x8 pixels × 4 bits = 32 bytes

	pixelWidth := imgWidth
	bytesPerRow := (pixelWidth + 1) / 2

	for tileY := 0; tileY < 8; tileY++ {
		pixelY := tileRow*8 + tileY
		if pixelY >= imgHeight {
			// Pad with transparent pixels
			continue
		}

		// Calculate starting byte for this row
		rowStartByte := pixelY * bytesPerRow

		tileByteIdx := tileY * 4 // 4 bytes per row in tile

		for tileX := 0; tileX < 8; tileX++ {
			pixelX := tileCol*8 + tileX
			if pixelX >= pixelWidth {
				// Pad with transparent pixels (index 0)
				continue
			}

			// Get pixel index from indexed data
			pixelByteIdx := rowStartByte + (pixelX / 2)
			if pixelByteIdx >= len(indexedData) {
				continue
			}

			var pixelIdx uint8
			if pixelX%2 == 0 {
				// Low nibble
				pixelIdx = indexedData[pixelByteIdx] & 0x0F
			} else {
				// High nibble
				pixelIdx = (indexedData[pixelByteIdx] >> 4) & 0x0F
			}

			// Store in tile format
			tileBytesIdx := tileByteIdx + (tileX / 2)
			if tileX%2 == 0 {
				// Low nibble
				tile[tileBytesIdx] = (tile[tileBytesIdx] & 0xF0) | (pixelIdx & 0x0F)
			} else {
				// High nibble
				tile[tileBytesIdx] = (tile[tileBytesIdx] & 0x0F) | ((pixelIdx & 0x0F) << 4)
			}
		}
	}

	return tile
}

// DecodeFF6Sprite decodes FF6 sprite data back to an image
func (c *FF6SpriteConverter) DecodeFF6Sprite(sprite *models.FF6Sprite) (image.Image, error) {
	if sprite == nil {
		return nil, fmt.Errorf("sprite is nil")
	}

	if sprite.Data == nil || len(sprite.Data) == 0 {
		return nil, fmt.Errorf("sprite data is empty")
	}

	if sprite.Palette == nil {
		return nil, fmt.Errorf("sprite palette is nil")
	}

	// Create image buffer
	img := image.NewRGBA(image.Rect(0, 0, sprite.Width, sprite.Height))

	// Decode tiles
	tilesWide := (sprite.Width + 7) / 8
	tilesHigh := (sprite.Height + 7) / 8
	tileSize := 32 // 8x8 tile in 4-bit format = 32 bytes

	tileIdx := 0

	for tileRow := 0; tileRow < tilesHigh; tileRow++ {
		for tileCol := 0; tileCol < tilesWide; tileCol++ {
			tileDataStart := tileIdx * tileSize
			tileDataEnd := tileDataStart + tileSize

			if tileDataEnd > len(sprite.Data) {
				// Invalid sprite data
				return img, fmt.Errorf("sprite data truncated at tile %d", tileIdx)
			}

			tileData := sprite.Data[tileDataStart:tileDataEnd]

			// Decode tile pixels
			for tileY := 0; tileY < 8; tileY++ {
				pixelY := tileRow*8 + tileY
				if pixelY >= sprite.Height {
					continue
				}

				byteIdx := tileY * 4
				for tileX := 0; tileX < 8; tileX++ {
					pixelX := tileCol*8 + tileX
					if pixelX >= sprite.Width {
						continue
					}

					// Get pixel index
					bytesIdx := byteIdx + (tileX / 2)
					var pixelIdx uint8
					if tileX%2 == 0 {
						pixelIdx = tileData[bytesIdx] & 0x0F
					} else {
						pixelIdx = (tileData[bytesIdx] >> 4) & 0x0F
					}

					// Get color from palette
					if pixelIdx < 16 {
						paletteColor := sprite.Palette.Colors[pixelIdx]
						r, g, b := paletteColor.ToRGB888()

						// Transparent if index 0
						alpha := uint8(255)
						if pixelIdx == 0 {
							alpha = 0
						}

						img.Set(pixelX, pixelY, color.RGBA{R: r, G: g, B: b, A: alpha})
					}
				}
			}

			tileIdx++
		}
	}

	return img, nil
}

// ConvertPalette converts a palette to another format
func (c *FF6SpriteConverter) ConvertPalette(srcPalette *models.Palette, targetFormat string) (*models.Palette, error) {
	if srcPalette == nil {
		return nil, fmt.Errorf("source palette is nil")
	}

	// Create copy
	dst := srcPalette.Clone()
	dst.Modified = time.Now() // Update modification time

	// Could implement various transformations based on targetFormat
	// For now, return copy
	return dst, nil
}

// PadImage pads an image to tile boundaries (multiples of 8)
func (c *FF6SpriteConverter) PadImage(img image.Image, fillColor color.Color) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Round up to nearest multiple of 8
	paddedWidth := ((width + 7) / 8) * 8
	paddedHeight := ((height + 7) / 8) * 8

	if paddedWidth == width && paddedHeight == height {
		return img // Already padded
	}

	// Create padded image
	padded := image.NewRGBA(image.Rect(0, 0, paddedWidth, paddedHeight))

	// Fill with transparent by default
	transparentColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	for y := 0; y < paddedHeight; y++ {
		for x := 0; x < paddedWidth; x++ {
			padded.Set(x, y, transparentColor)
		}
	}

	// Copy original image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			padded.Set(x-bounds.Min.X, y-bounds.Min.Y, img.At(x, y))
		}
	}

	return padded
}
