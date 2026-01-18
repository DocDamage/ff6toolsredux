package io

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"ffvi_editor/models"
)

// SpriteImporter handles the complete sprite import workflow
type SpriteImporter struct {
	decoder     *ImageDecoder
	extractor   *PaletteExtractor
	quantizer   *ColorQuantizer
	converter   *FF6SpriteConverter
	validator   *SpriteValidator
	lastErrors  []ValidationError
	importCount int
}

// NewSpriteImporter creates a new sprite importer
func NewSpriteImporter() *SpriteImporter {
	return &SpriteImporter{
		decoder:     NewImageDecoder(),
		extractor:   NewPaletteExtractor(),
		quantizer:   NewColorQuantizer("floyd-steinberg"),
		converter:   NewFF6SpriteConverter(),
		validator:   NewSpriteValidator(),
		lastErrors:  make([]ValidationError, 0),
		importCount: 0,
	}
}

// ImportResult contains the result of an import operation
type ImportResult struct {
	Sprite      *models.FF6Sprite
	Success     bool
	Errors      []ValidationError
	Warnings    []string
	StepsFailed []string
	Duration    time.Duration
}

// Import performs the complete import workflow
func (si *SpriteImporter) Import(opts *models.SpriteImportOptions) *ImportResult {
	startTime := time.Now()
	result := &ImportResult{
		Errors:      make([]ValidationError, 0),
		Warnings:    make([]string, 0),
		StepsFailed: make([]string, 0),
		Success:     true,
	}

	// Step 1: Validate options
	optsValidation := ValidateImportOptions(opts)
	if !optsValidation.IsValid {
		for _, err := range optsValidation.Errors {
			if err.Level == "error" {
				result.Errors = append(result.Errors, err)
				result.StepsFailed = append(result.StepsFailed, "Validate Options")
				result.Success = false
				result.Duration = time.Since(startTime)
				return result
			}
		}
		result.Warnings = append(result.Warnings, "Import options have warnings")
	}

	// Step 2: Decode image
	img, format, err := si.decoder.Decode(opts.SourcePath)
	if err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("Failed to decode image: %v", err),
		})
		result.StepsFailed = append(result.StepsFailed, "Decode Image")
		result.Success = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Step 3: Extract palette
	palette, err := si.extractor.Extract(img, opts.MaxColors)
	if err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Level:   "error",
			Message: fmt.Sprintf("Failed to extract palette: %v", err),
		})
		result.StepsFailed = append(result.StepsFailed, "Extract Palette")
		result.Success = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Step 4: Apply quantization if needed
	if opts.MaxColors < 16 {
		var quantErr error
		img, palette, quantErr = si.quantizer.Quantize(img, opts.MaxColors)
		if quantErr != nil {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Color quantization produced warning: %v", quantErr))
		}
	}

	// Step 5: Use target palette if provided
	if opts.TargetPalette != nil {
		palette = opts.TargetPalette
	}

	// Step 6: Pad image to tile boundaries if enabled
	if opts.AutoPadding {
		img = si.converter.PadImage(img, nil)
	}

	// Step 7: Detect sprite type if enabled
	spriteType := opts.TargetType
	if opts.AutoDetectType {
		bounds := img.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y
		spriteType = si.detectSpriteType(width, height)
	}

	// Step 8: Convert to FF6 format
	sprite := si.converter.ToFF6Format(img, palette, spriteType)

	// Set source information
	sprite.SourceFile = opts.SourcePath
	sprite.ImportedFrom = format
	sprite.ImportDate = time.Now()
	sprite.Author = "Imported"
	sprite.CreatedDate = time.Now()

	// Generate ID
	si.importCount++
	sprite.ID = fmt.Sprintf("imported_%s_%d", format, si.importCount)

	// Step 9: Validate sprite
	valResult := si.validator.Validate(sprite)
	if !valResult.IsValid {
		for _, err := range valResult.Errors {
			result.Errors = append(result.Errors, err)
		}
		result.StepsFailed = append(result.StepsFailed, "Validate Sprite")
		result.Success = false
	} else {
		for _, err := range valResult.Errors {
			if err.Level == "warning" {
				result.Warnings = append(result.Warnings, err.Message)
			}
		}
	}

	result.Sprite = sprite
	result.Duration = time.Since(startTime)
	si.lastErrors = result.Errors

	return result
}

// detectSpriteType guesses sprite type from dimensions
func (si *SpriteImporter) detectSpriteType(width, height int) models.SpriteType {
	// Match against known sprite sizes
	if width == 16 && height == 24 {
		return models.SpriteTypeCharacter
	}
	if width == 32 && height == 32 {
		return models.SpriteTypeBattle
	}
	if width == 48 && height == 48 {
		return models.SpriteTypePortrait
	}
	if width == 16 && height == 16 {
		return models.SpriteTypeOverworld
	}

	// Try to guess based on aspect ratio
	aspectRatio := float64(width) / float64(height)

	if aspectRatio < 0.8 {
		// Taller than wide - likely character
		return models.SpriteTypeCharacter
	} else if aspectRatio > 1.2 {
		// Wider than tall - likely overworld
		return models.SpriteTypeOverworld
	}

	// Default to battle sprite
	return models.SpriteTypeBattle
}

// ImportMultiple imports multiple image files
func (si *SpriteImporter) ImportMultiple(paths []string, opts *models.SpriteImportOptions) []*ImportResult {
	results := make([]*ImportResult, 0, len(paths))

	for _, path := range paths {
		// Create options copy with updated path
		optsCopy := *opts
		optsCopy.SourcePath = path
		optsCopy.SourceFormat = filepath.Ext(path)

		result := si.Import(&optsCopy)
		results = append(results, result)
	}

	return results
}

// GetLastErrors returns the last import errors
func (si *SpriteImporter) GetLastErrors() []ValidationError {
	return si.lastErrors
}

// BatchImportOptions configures batch import behavior
type BatchImportOptions struct {
	BaseOptions      *models.SpriteImportOptions
	StopOnError      bool
	CreatePaletteFor int // Create separate palette every N sprites (0 = use same)
	ReusePalette     bool
	OutputDirectory  string
	CreateZipArchive bool
}

// BatchImport performs batch import with options
func (si *SpriteImporter) BatchImport(paths []string, batchOpts *BatchImportOptions) []*ImportResult {
	if batchOpts == nil {
		batchOpts = &BatchImportOptions{
			BaseOptions:      models.NewSpriteImportOptions(),
			StopOnError:      false,
			ReusePalette:     false,
			CreateZipArchive: false,
		}
	}

	results := make([]*ImportResult, 0, len(paths))
	var lastPalette *models.Palette

	for i, path := range paths {
		// Create options copy
		optsCopy := *batchOpts.BaseOptions
		optsCopy.SourcePath = path
		optsCopy.SourceFormat = filepath.Ext(path)

		// Reuse palette if configured
		if batchOpts.ReusePalette && lastPalette != nil {
			optsCopy.TargetPalette = lastPalette
		}

		// Create new palette every N sprites
		if batchOpts.CreatePaletteFor > 0 && i%batchOpts.CreatePaletteFor == 0 {
			optsCopy.TargetPalette = nil
		}

		result := si.Import(&optsCopy)
		results = append(results, result)

		// Save palette for reuse
		if result.Success && result.Sprite != nil {
			lastPalette = result.Sprite.Palette
		}

		// Stop on error if configured
		if batchOpts.StopOnError && !result.Success {
			break
		}
	}

	return results
}

// SpriteExporter handles sprite export operations
type SpriteExporter struct {
	converter *FF6SpriteConverter
}

// NewSpriteExporter creates a new sprite exporter
func NewSpriteExporter() *SpriteExporter {
	return &SpriteExporter{
		converter: NewFF6SpriteConverter(),
	}
}

// ExportToPNG exports sprite as PNG image
func (se *SpriteExporter) ExportToPNG(sprite *models.FF6Sprite, outputPath string) error {
	if sprite == nil {
		return fmt.Errorf("sprite is nil")
	}

	// Decode sprite back to image
	img, err := se.converter.DecodeFF6Sprite(sprite)
	if err != nil {
		return fmt.Errorf("failed to decode sprite: %w", err)
	}

	// Save as PNG
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}
