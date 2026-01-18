package io

import (
	"ffvi_editor/models"
	"fmt"
	"sync"
	"time"
)

// ROMSpriteExtractor extracts sprite data from FF6 ROM files
type ROMSpriteExtractor struct {
	rom               *ROMLoader
	decompressor      *LZ77Decompressor
	mu                sync.RWMutex
	paletteCache      map[int]*models.Palette   // Cached palettes (all 14 characters)
	battleSpriteCache map[int]*models.FF6Sprite // Cached battle sprites
	fieldSpriteCache  map[int]*models.FF6Sprite // Cached field sprites
	cachesLoadedOnce  bool                      // Flag to prevent redundant batch loads
}

// NewROMSpriteExtractor creates a new sprite extractor
func NewROMSpriteExtractor(romPath string) (*ROMSpriteExtractor, error) {
	var rom *ROMLoader

	// If path provided, try to load it
	if romPath != "" {
		rom = NewROMLoader(romPath)
		if err := rom.Load(); err != nil {
			return nil, fmt.Errorf("failed to load ROM from %s: %w", romPath, err)
		}
	} else {
		// Try to find ROM from default locations
		rom = TryLoadFromDefaultLocations()
		if !rom.IsValid() {
			return nil, fmt.Errorf("no ROM path provided and no ROM found in default locations")
		}
	}

	return &ROMSpriteExtractor{
		rom:               rom,
		decompressor:      NewLZ77Decompressor(),
		paletteCache:      make(map[int]*models.Palette),
		battleSpriteCache: make(map[int]*models.FF6Sprite),
		fieldSpriteCache:  make(map[int]*models.FF6Sprite),
	}, nil
}

// ExtractCharacterSprite extracts a character's field sprite from ROM
func (e *ROMSpriteExtractor) ExtractCharacterSprite(characterID int) (*models.FF6Sprite, error) {
	// Get sprite offset for this character
	offset, ok := CharacterSpriteOffsets[characterID]
	if !ok {
		return nil, fmt.Errorf("unknown character ID: %d", characterID)
	}

	// Read compressed sprite data from ROM
	// Typical field sprite is ~1-2KB compressed
	compressedData, err := e.rom.ReadBytes(offset.FieldOffset, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to read sprite data: %w", err)
	}

	// Decompress if needed
	var spriteData []byte
	if offset.Compressed {
		spriteData, err = e.decompressor.Decompress(compressedData)
		if err != nil {
			// If decompression fails, try to use raw data
			spriteData = compressedData
		}
	} else {
		spriteData = compressedData
	}

	// Field sprites are 16x24 pixels, 3 frames
	// Each frame is 16*24/2 = 192 bytes (4bpp = 0.5 bytes per pixel)
	expectedSize := 16 * 24 * 3 / 2 // 576 bytes for 3 frames

	// If we got more data, truncate to expected size
	if len(spriteData) > expectedSize {
		spriteData = spriteData[:expectedSize]
	}

	// If we got less data, pad with zeros
	if len(spriteData) < expectedSize {
		padding := make([]byte, expectedSize-len(spriteData))
		spriteData = append(spriteData, padding...)
	}

	// Extract palette (usually stored near sprite data or use default)
	palette := e.extractDefaultPalette(characterID)

	// Create FF6Sprite
	sprite := &models.FF6Sprite{
		Width:   16,
		Height:  24,
		Frames:  3,
		Data:    spriteData,
		Palette: palette,
	}

	return sprite, nil
}

// extractDefaultPalette returns a default palette for a character
func (e *ROMSpriteExtractor) extractDefaultPalette(characterID int) *models.Palette {
	// Create default palette with some basic colors
	// In a full implementation, this would read from ROM palette data
	colors := [16]models.RGB555{
		{R: 0, G: 0, B: 0},    // 0: Transparent/Black
		{R: 31, G: 31, B: 31}, // 1: White
		{R: 25, G: 20, B: 15}, // 2: Skin tone
		{R: 15, G: 10, B: 5},  // 3: Shadow
		{R: 20, G: 15, B: 10}, // 4: Light skin
		{R: 10, G: 8, B: 25},  // 5: Blue/purple (clothing)
		{R: 25, G: 10, B: 5},  // 6: Red/brown (hair)
		{R: 30, G: 25, B: 5},  // 7: Yellow/blonde
		{R: 5, G: 15, B: 5},   // 8: Green
		{R: 20, G: 5, B: 5},   // 9: Dark red
		{R: 15, G: 15, B: 20}, // 10: Light blue
		{R: 25, G: 25, B: 15}, // 11: Beige
		{R: 10, G: 10, B: 10}, // 12: Dark gray
		{R: 20, G: 20, B: 20}, // 13: Medium gray
		{R: 5, G: 5, B: 20},   // 14: Dark blue
		{R: 28, G: 28, B: 28}, // 15: Light gray
	}

	return &models.Palette{
		Colors: colors,
		Name:   fmt.Sprintf("Character %d Default", characterID),
	}
}

// IsROMLoaded returns true if a valid ROM is loaded
func (e *ROMSpriteExtractor) IsROMLoaded() bool {
	return e.rom != nil && e.rom.IsValid()
}

// ExtractBattleSprite extracts a character's battle sprite from ROM
// Battle sprites are 32x32 pixels with 6 frames (idle, attack, magic, damage, victory, dead)
func (e *ROMSpriteExtractor) ExtractBattleSprite(characterID int) (*models.FF6Sprite, error) {
	// Get sprite offset for this character
	offset, ok := CharacterSpriteOffsets[characterID]
	if !ok {
		return nil, fmt.Errorf("unknown character ID: %d", characterID)
	}

	// Read compressed battle sprite data from ROM
	// Battle sprites are typically 3-8KB compressed
	compressedData, err := e.rom.ReadBytes(offset.BattleOffset, 8192)
	if err != nil {
		return nil, fmt.Errorf("failed to read battle sprite data: %w", err)
	}

	// Decompress if needed
	var spriteData []byte
	if offset.Compressed {
		spriteData, err = e.decompressor.Decompress(compressedData)
		if err != nil {
			// If decompression fails, try to use raw data
			spriteData = compressedData
		}
	} else {
		spriteData = compressedData
	}

	// Battle sprites are 32x32 pixels, 6 frames
	// Each frame: 32*32/2 = 512 bytes (4bpp = 0.5 bytes per pixel)
	// Total: 512 * 6 = 3072 bytes for 6 frames
	expectedSize := 32 * 32 * 6 / 2 // 3072 bytes for 6 frames

	// If we got more data, truncate to expected size
	if len(spriteData) > expectedSize {
		spriteData = spriteData[:expectedSize]
	}

	// If we got less data, pad with zeros
	if len(spriteData) < expectedSize {
		padding := make([]byte, expectedSize-len(spriteData))
		spriteData = append(spriteData, padding...)
	}

	// Extract palette for battle sprite
	palette, err := e.ExtractCharacterPalette(characterID)
	if err != nil {
		// Fall back to default if extraction fails
		palette = e.extractDefaultPalette(characterID)
	}

	// Create FF6Sprite for battle sprite
	sprite := &models.FF6Sprite{
		Width:   32,
		Height:  32,
		Frames:  6,
		Data:    spriteData,
		Palette: palette,
		Type:    models.SpriteTypeBattle,
	}

	return sprite, nil
}

// ExtractCharacterPalette extracts the 16-color palette for a character from ROM
// FF6 stores palettes in a separate location from sprite data
// Palette locations typically start at 0x2C0000 in USA version
func (e *ROMSpriteExtractor) ExtractCharacterPalette(characterID int) (*models.Palette, error) {
	if characterID < 0 || characterID > 13 {
		return nil, fmt.Errorf("invalid character ID: %d", characterID)
	}

	// FF6 character palettes are stored in sequence
	// Each palette is 16 colors × 2 bytes (5-bit RGB) = 32 bytes
	paletteBaseOffset := 0x2C0000 // Character palette base in ROM
	paletteSize := 32             // 16 colors × 2 bytes
	paletteOffset := paletteBaseOffset + (characterID * paletteSize)

	// Read palette data from ROM
	paletteData, err := e.rom.ReadBytes(paletteOffset, paletteSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read palette data: %w", err)
	}

	if len(paletteData) < paletteSize {
		return nil, fmt.Errorf("insufficient palette data read: got %d, expected %d", len(paletteData), paletteSize)
	}

	// Parse palette data (5-bit RGB, little-endian 16-bit values)
	var colors [16]models.RGB555
	for i := 0; i < 16; i++ {
		// Each color is 2 bytes, little-endian
		// Format: GGGRRRRR XBBBBBGG (where X is unused)
		low := paletteData[i*2]
		high := paletteData[i*2+1]

		// Extract 5-bit components
		r := uint8(low & 0x1F)
		g := uint8(((low >> 5) & 0x07) | ((high & 0x03) << 3))
		b := uint8((high >> 2) & 0x1F)

		colors[i] = models.RGB555{R: r, G: g, B: b}
	}

	character, ok := CharacterSpriteOffsets[characterID]
	if !ok {
		character.Name = fmt.Sprintf("Character %d", characterID)
	}

	palette := &models.Palette{
		Colors:   colors,
		Name:     fmt.Sprintf("%s Palette", character.Name),
		Created:  time.Now(),
		Modified: time.Now(),
	}

	return palette, nil
}

// ExtractAllCharacterPalettes extracts palettes for all 14 playable characters
func (e *ROMSpriteExtractor) ExtractAllCharacterPalettes() (map[int]*models.Palette, error) {
	palettes := make(map[int]*models.Palette)
	errors := make(map[int]error)

	for i := 0; i < 14; i++ {
		pal, err := e.ExtractCharacterPalette(i)
		if err != nil {
			errors[i] = err
			// Use default palette if extraction fails
			palettes[i] = e.extractDefaultPalette(i)
		} else {
			palettes[i] = pal
		}
	}

	// If all extractions failed, return error
	if len(errors) == 14 {
		return palettes, fmt.Errorf("failed to extract any character palettes from ROM")
	}

	return palettes, nil
}

// LoadAllPalettesCached loads all 14 character palettes into cache concurrently for fast access
// This should be called once on startup in a background goroutine
func (e *ROMSpriteExtractor) LoadAllPalettesCached() error {
	e.mu.Lock()
	if e.cachesLoadedOnce {
		e.mu.Unlock()
		return nil // Already loaded
	}
	e.mu.Unlock()

	// Load all 14 palettes concurrently using worker pool pattern
	const maxWorkers = 6
	results := make(chan struct {
		id  int
		pal *models.Palette
		err error
	}, 14)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxWorkers)

	for i := 0; i < 14; i++ {
		wg.Add(1)
		go func(charID int) {
			defer wg.Done()

			// Acquire semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			pal, err := e.ExtractCharacterPalette(charID)
			results <- struct {
				id  int
				pal *models.Palette
				err error
			}{charID, pal, err}
		}(i)
	}

	// Wait for all to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	palettes := make(map[int]*models.Palette)
	errorCount := 0

	for result := range results {
		if result.err != nil {
			errorCount++
			palettes[result.id] = e.extractDefaultPalette(result.id)
		} else {
			palettes[result.id] = result.pal
		}
	}

	e.mu.Lock()
	e.paletteCache = palettes
	e.cachesLoadedOnce = true
	e.mu.Unlock()

	if errorCount == 14 {
		return fmt.Errorf("failed to extract any character palettes from ROM")
	}

	fmt.Printf("✅ All character palettes loaded into cache (%d successful, %d fallback)\n", 14-errorCount, errorCount)
	return nil
}

// ExtractCharacterPalette returns palette from cache if available, otherwise extracts and caches it
func (e *ROMSpriteExtractor) ExtractCharacterPaletteWithCache(characterID int) (*models.Palette, error) {
	e.mu.RLock()
	if pal, ok := e.paletteCache[characterID]; ok {
		e.mu.RUnlock()
		return pal, nil
	}
	e.mu.RUnlock()

	// Not in cache, extract it
	pal, err := e.ExtractCharacterPalette(characterID)
	if err != nil {
		return nil, err
	}

	// Cache for future use
	e.mu.Lock()
	e.paletteCache[characterID] = pal
	e.mu.Unlock()

	return pal, nil
}
