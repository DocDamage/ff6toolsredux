package io

import (
	"fmt"
	"os"
	"strings"
)

// ROMType identifies the FF6 ROM version
// Note: FF3 USA = FF6 (the US release was renumbered as FF3)
type ROMType int

const (
	ROMTypeUnknown ROMType = iota
	ROMTypeUSA             // Final Fantasy III (USA) - same as FF6
	ROMTypeJPN             // Final Fantasy VI (Japan)
	ROMTypeEUR             // Final Fantasy VI (Europe)
)

// SNES ROM header offsets
const (
	snesHeaderOffsetLoROM = 0x7FC0 // LoROM header location (FF6 uses LoROM)
	snesHeaderOffsetHiROM = 0xFFC0 // HiROM header location (for completeness)
	snesGameTitleOffset   = 0      // Game title at header + 0
	snesGameTitleLength   = 21     // Game title is 21 bytes
)

// Known FF6 game title patterns (case-insensitive substrings)
var ff6TitlePatterns = []string{
	"FINAL FANTASY",     // All versions contain this
	"FINAL FANTASY 6",   // Japan explicit
	"FINAL FANTASY VI",  // Europe explicit
	"FINAL FANTASY 3",   // USA version
	"FINAL FANTASY III", // USA version full
}

// ROMLoader handles loading and validating FF6 ROM files
type ROMLoader struct {
	path      string
	data      []byte
	romType   ROMType
	hasHeader bool
}

// NewROMLoader creates a new ROM loader
func NewROMLoader(path string) *ROMLoader {
	return &ROMLoader{
		path:    path,
		romType: ROMTypeUnknown,
	}
}

// TryLoadFromDefaultLocations attempts to load ROM from common default locations
func TryLoadFromDefaultLocations() *ROMLoader {
	// List of common paths to try
	defaultPaths := []string{
		"save_data/Final Fantasy III (USA) (Rev 1).sfc",
		"save_data/Final Fantasy III (USA) (Rev A).sfc",
		"save_data/ff3us.sfc",
		"save_data/Final Fantasy VI.sfc",
		"save_data/ff6.sfc",
		"save_data/ff6.smc",
		"snes/Final Fantasy III (USA) (Rev A).sfc",
		"snes/ff3us.sfc",
		"snes/Final Fantasy VI (USA).sfc",
		"snes/Final Fantasy VI.sfc",
		"Final Fantasy III (USA) (Rev A).sfc",
		"ff3us.sfc",
		"ff6.sfc",
		"ff6.smc",
	}

	for _, path := range defaultPaths {
		loader := NewROMLoader(path)
		if err := loader.Load(); err == nil {
			fmt.Printf("Found and loaded ROM from: %s\n", path)
			return loader
		}
	}

	// Return an empty loader if nothing found
	return NewROMLoader("")
}

// Load reads and validates the ROM file
func (r *ROMLoader) Load() error {
	// Read ROM file
	data, err := os.ReadFile(r.path)
	if err != nil {
		return fmt.Errorf("failed to read ROM: %w", err)
	}

	// Check if ROM has SMC header (512 bytes)
	if len(data)%1024 == 512 {
		r.hasHeader = true
		// Skip header
		data = data[512:]
	}

	// Validate ROM size (FF6 is 3MB or 4MB)
	if len(data) != 3*1024*1024 && len(data) != 4*1024*1024 {
		return fmt.Errorf("invalid ROM size: %d bytes (expected 3MB or 4MB)", len(data))
	}

	// Detect ROM type by reading SNES internal header
	r.romType = r.detectROMType(data)

	if r.romType == ROMTypeUnknown {
		return fmt.Errorf("not a valid Final Fantasy VI ROM (header validation failed)")
	}

	r.data = data
	return nil
}

// detectROMType identifies the ROM version by reading the SNES internal header
func (r *ROMLoader) detectROMType(data []byte) ROMType {
	// Try LoROM header first (FF6 uses LoROM)
	if title := r.readGameTitle(data, snesHeaderOffsetLoROM); title != "" {
		if r.isFF6Title(title) {
			return r.identifyRegion(title)
		}
	}

	// Try HiROM header as fallback
	if title := r.readGameTitle(data, snesHeaderOffsetHiROM); title != "" {
		if r.isFF6Title(title) {
			return r.identifyRegion(title)
		}
	}

	return ROMTypeUnknown
}

// readGameTitle extracts the game title from the SNES header
func (r *ROMLoader) readGameTitle(data []byte, headerOffset int) string {
	if len(data) < headerOffset+snesGameTitleLength {
		return ""
	}

	titleBytes := data[headerOffset : headerOffset+snesGameTitleLength]
	// Convert to string and trim null bytes and spaces
	title := string(titleBytes)
	title = strings.TrimRight(title, "\x00 ")
	return strings.ToUpper(title)
}

// isFF6Title checks if the title matches any known FF6 pattern
func (r *ROMLoader) isFF6Title(title string) bool {
	titleUpper := strings.ToUpper(title)
	for _, pattern := range ff6TitlePatterns {
		if strings.Contains(titleUpper, pattern) {
			return true
		}
	}
	return false
}

// identifyRegion determines the ROM region based on the title
func (r *ROMLoader) identifyRegion(title string) ROMType {
	titleUpper := strings.ToUpper(title)

	// USA uses "FINAL FANTASY III" or "FINAL FANTASY 3"
	if strings.Contains(titleUpper, "FINAL FANTASY III") ||
		strings.Contains(titleUpper, "FINAL FANTASY 3") {
		return ROMTypeUSA
	}

	// Japan uses "FINAL FANTASY 6" or just "FINAL FANTASY VI"
	if strings.Contains(titleUpper, "FINAL FANTASY 6") ||
		(strings.Contains(titleUpper, "FINAL FANTASY VI") && !strings.Contains(titleUpper, "III")) {
		return ROMTypeJPN
	}

	// If it contains FINAL FANTASY but isn't USA or JPN, assume EUR or return USA as default
	return ROMTypeEUR
}

// GetData returns the ROM data (without header if present)
func (r *ROMLoader) GetData() []byte {
	return r.data
}

// GetType returns the detected ROM type
func (r *ROMLoader) GetType() ROMType {
	return r.romType
}

// IsValid returns true if ROM loaded successfully
func (r *ROMLoader) IsValid() bool {
	return len(r.data) > 0
}

// ReadBytes reads bytes from ROM at the given offset
func (r *ROMLoader) ReadBytes(offset, length int) ([]byte, error) {
	if offset < 0 || offset+length > len(r.data) {
		return nil, fmt.Errorf("offset out of bounds: %d (ROM size: %d)", offset, len(r.data))
	}
	return r.data[offset : offset+length], nil
}
