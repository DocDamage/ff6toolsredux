package io

import (
	"fmt"

	"ffvi_editor/models"
)

// ValidationError represents a sprite validation issue
type ValidationError struct {
	Level   string // "error", "warning"
	Message string
}

// SpriteValidator validates FF6 sprite compliance
type SpriteValidator struct {
	strictMode bool
}

// NewSpriteValidator creates a new validator
func NewSpriteValidator() *SpriteValidator {
	return &SpriteValidator{
		strictMode: false,
	}
}

// ValidationResult holds validation results
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// Validate checks if a sprite meets FF6 requirements
func (sv *SpriteValidator) Validate(sprite *models.FF6Sprite) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  make([]ValidationError, 0),
	}

	// Check basic properties
	sv.validateBasicProperties(sprite, result)

	// Check dimensions
	sv.validateDimensions(sprite, result)

	// Check palette
	sv.validatePalette(sprite, result)

	// Check data
	sv.validateData(sprite, result)

	// Determine overall validity
	for _, err := range result.Errors {
		if err.Level == "error" {
			result.IsValid = false
			break
		}
	}

	return result
}

// validateBasicProperties checks sprite ID, name, etc.
func (sv *SpriteValidator) validateBasicProperties(sprite *models.FF6Sprite, result *ValidationResult) {
	if sprite == nil {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Sprite is nil",
		})
		return
	}

	if sprite.Name == "" {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "warning",
			Message: "Sprite name is empty",
		})
	}

	if sprite.Type < 0 || sprite.Type > 5 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("Invalid sprite type: %d", sprite.Type),
		})
	}

	if sprite.Frames <= 0 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "warning",
			Message: "Sprite has no frames",
		})
	}
}

// validateDimensions checks sprite dimensions
func (sv *SpriteValidator) validateDimensions(sprite *models.FF6Sprite, result *ValidationResult) {
	expectedWidth, expectedHeight := sprite.Type.GetDimensions()

	// Check if dimensions match expected for type
	if sprite.Width != expectedWidth || sprite.Height != expectedHeight {
		result.Errors = append(result.Errors, ValidationError{
			Level: "error",
			Message: fmt.Sprintf("Dimensions %dx%d don't match expected %dx%d for %s sprite",
				sprite.Width, sprite.Height, expectedWidth, expectedHeight, sprite.Type.String()),
		})
	}

	// Check dimensions are multiples of 8 (tile-based)
	if sprite.Width%8 != 0 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("Width %d is not a multiple of 8 (tile-based)", sprite.Width),
		})
	}

	if sprite.Height%8 != 0 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("Height %d is not a multiple of 8 (tile-based)", sprite.Height),
		})
	}

	// Check reasonable size limits
	maxTiles := 64 // 64 tiles max (~512x512 at 8x8 tiles)
	tilesWide := (sprite.Width + 7) / 8
	tilesHigh := (sprite.Height + 7) / 8
	totalTiles := tilesWide * tilesHigh

	if totalTiles > maxTiles {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("Sprite exceeds maximum size: %d tiles (max %d)", totalTiles, maxTiles),
		})
	}

	// Check minimum size
	if sprite.Width < 8 || sprite.Height < 8 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Sprite too small (minimum 8x8)",
		})
	}
}

// validatePalette checks palette validity
func (sv *SpriteValidator) validatePalette(sprite *models.FF6Sprite, result *ValidationResult) {
	if sprite.Palette == nil {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Sprite has no palette",
		})
		return
	}

	// Check palette has 16 colors
	colorCount := 0
	for i := 0; i < 16; i++ {
		// Count non-black colors (non-zero)
		color := sprite.Palette.Colors[i]
		if color.R > 0 || color.G > 0 || color.B > 0 {
			colorCount++
		}
	}

	if colorCount == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "warning",
			Message: "Palette appears to be all black (likely needs adjustment)",
		})
	}

	// Verify 5-bit RGB values
	for i := 0; i < 16; i++ {
		color := sprite.Palette.Colors[i]
		if color.R > 31 {
			result.Errors = append(result.Errors, ValidationError{
				Level:   "error",
				Message: fmt.Sprintf("Palette color %d has invalid R value: %d (max 31)", i, color.R),
			})
		}
		if color.G > 31 {
			result.Errors = append(result.Errors, ValidationError{
				Level:   "error",
				Message: fmt.Sprintf("Palette color %d has invalid G value: %d (max 31)", i, color.G),
			})
		}
		if color.B > 31 {
			result.Errors = append(result.Errors, ValidationError{
				Level:   "error",
				Message: fmt.Sprintf("Palette color %d has invalid B value: %d (max 31)", i, color.B),
			})
		}
	}
}

// validateData checks sprite data validity
func (sv *SpriteValidator) validateData(sprite *models.FF6Sprite, result *ValidationResult) {
	if sprite.Data == nil {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Sprite data is nil",
		})
		return
	}

	if len(sprite.Data) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Sprite data is empty",
		})
		return
	}

	// Calculate expected data size
	expectedSize := sprite.GetExpectedDataSize()
	actualSize := len(sprite.Data)

	if !sprite.IsCompressed {
		if actualSize < expectedSize {
			result.Errors = append(result.Errors, ValidationError{
				Level: "warning",
				Message: fmt.Sprintf("Sprite data size %d is smaller than expected %d bytes (may be incomplete)",
					actualSize, expectedSize),
			})
		} else if actualSize > expectedSize*2 { // Allow some overhead
			result.Errors = append(result.Errors, ValidationError{
				Level: "warning",
				Message: fmt.Sprintf("Sprite data size %d is much larger than expected %d bytes",
					actualSize, expectedSize),
			})
		}
	} else {
		// For compressed data, check reasonable compression ratio
		compressionRatio := float64(actualSize) / float64(expectedSize)
		if compressionRatio > 0.8 {
			result.Errors = append(result.Errors, ValidationError{
				Level:   "warning",
				Message: "Compressed sprite data appears to have poor compression ratio",
			})
		}
	}

	// Check for reasonable frame count
	if sprite.Frames > 32 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "warning",
			Message: fmt.Sprintf("Sprite has unusually many frames: %d", sprite.Frames),
		})
	}

	if sprite.Frames < 1 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Sprite has no frames",
		})
	}
}

// ValidateImportOptions validates import options
func ValidateImportOptions(opts *models.SpriteImportOptions) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  make([]ValidationError, 0),
	}

	if opts == nil {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Import options are nil",
		})
		result.IsValid = false
		return result
	}

	if opts.MaxColors < 2 || opts.MaxColors > 256 {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("MaxColors must be between 2 and 256, got %d", opts.MaxColors),
		})
	}

	validDithers := map[string]bool{
		"floyd-steinberg": true,
		"bayer":           true,
		"none":            true,
	}

	if !validDithers[opts.DitherMethod] {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "warning",
			Message: fmt.Sprintf("Unknown dither method: %s, using floyd-steinberg", opts.DitherMethod),
		})
	}

	if opts.SourcePath == "" {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: "Source path not specified",
		})
	}

	// Determine overall validity
	for _, err := range result.Errors {
		if err.Level == "error" {
			result.IsValid = false
			break
		}
	}

	return result
}

// PrintValidationResult formats validation results for display
func PrintValidationResult(result *ValidationResult) string {
	if result.IsValid && len(result.Errors) == 0 {
		return "✓ Sprite is valid"
	}

	output := ""
	if result.IsValid {
		output = "✓ Sprite is valid with warnings:\n"
	} else {
		output = "✗ Sprite validation failed:\n"
	}

	for _, err := range result.Errors {
		prefix := "⚠"
		if err.Level == "error" {
			prefix = "✗"
		}
		output += fmt.Sprintf("%s %s\n", prefix, err.Message)
	}

	return output
}
