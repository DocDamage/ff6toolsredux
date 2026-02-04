package io

import (
	"fmt"
)

// FF6 character sprite offsets in ROM (USA version)
// These are SNES addresses, need to convert to PC offsets

// CharacterSpriteOffsets maps character IDs to their sprite data offsets in ROM
var CharacterSpriteOffsets = map[int]SpriteDataOffset{
	0:  {Name: "Terra", FieldOffset: 0x150000, BattleOffset: 0x268000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	1:  {Name: "Locke", FieldOffset: 0x151000, BattleOffset: 0x26A000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	2:  {Name: "Cyan", FieldOffset: 0x152000, BattleOffset: 0x26C000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	3:  {Name: "Shadow", FieldOffset: 0x153000, BattleOffset: 0x26E000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	4:  {Name: "Edgar", FieldOffset: 0x154000, BattleOffset: 0x270000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	5:  {Name: "Sabin", FieldOffset: 0x155000, BattleOffset: 0x272000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	6:  {Name: "Celes", FieldOffset: 0x156000, BattleOffset: 0x274000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	7:  {Name: "Strago", FieldOffset: 0x157000, BattleOffset: 0x276000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	8:  {Name: "Relm", FieldOffset: 0x158000, BattleOffset: 0x278000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	9:  {Name: "Setzer", FieldOffset: 0x159000, BattleOffset: 0x27A000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	10: {Name: "Mog", FieldOffset: 0x15A000, BattleOffset: 0x27C000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	11: {Name: "Gau", FieldOffset: 0x15B000, BattleOffset: 0x27E000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	12: {Name: "Gogo", FieldOffset: 0x15C000, BattleOffset: 0x280000, Compressed: true, FieldSize: 576, BattleSize: 1536},
	13: {Name: "Umaro", FieldOffset: 0x15D000, BattleOffset: 0x282000, Compressed: true, FieldSize: 576, BattleSize: 1536},
}

// SpriteDataOffset contains ROM offsets for a character's sprites
type SpriteDataOffset struct {
	Name          string
	FieldOffset   int  // Field/walking sprite data offset
	BattleOffset  int  // Battle sprite data offset
	Compressed    bool // Whether data is LZ77 compressed
	PaletteOffset int  // Palette data offset (if separate)
	FieldSize     int  // Expected size of field sprite data (bytes)
	BattleSize    int  // Expected size of battle sprite data (bytes)
}

// SpriteFrameInfo contains metadata about sprite animation frames
type SpriteFrameInfo struct {
	Width    int // Pixel width
	Height   int // Pixel height
	Frames   int // Number of frames
	BPP      int // Bits per pixel (usually 4 for FF6)
	Animated bool
}

// FieldSpriteInfo describes field sprite format (16x24, 3 frames)
var FieldSpriteInfo = SpriteFrameInfo{
	Width:    16,
	Height:   24,
	Frames:   3,
	BPP:      4,
	Animated: true,
}

// BattleSpriteInfo describes battle sprite format (32x32, 6 frames)
var BattleSpriteInfo = SpriteFrameInfo{
	Width:    32,
	Height:   32,
	Frames:   6,
	BPP:      4,
	Animated: true,
}

// LZ77Decompressor decompresses FF6's custom LZ77 format
type LZ77Decompressor struct{}

// NewLZ77Decompressor creates a new decompressor
func NewLZ77Decompressor() *LZ77Decompressor {
	return &LZ77Decompressor{}
}

// Decompress decompresses LZ77 compressed data with optimized buffer allocation
func (d *LZ77Decompressor) Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, fmt.Errorf("empty compressed data")
	}

	// Pre-allocate buffer with estimated size (compressed data typically expands 2-4x)
	decompressed := make([]byte, 0, len(compressed)*3)
	pos := 0

	for pos < len(compressed) {
		// Read control byte
		if pos >= len(compressed) {
			break
		}
		control := compressed[pos]
		pos++

		// Process 8 commands (1 bit each in control byte)
		for bit := 0; bit < 8; bit++ {
			if pos >= len(compressed) {
				break
			}

			if (control & (1 << bit)) == 0 {
				// Literal byte
				decompressed = append(decompressed, compressed[pos])
				pos++
			} else {
				// Back reference
				if pos+1 >= len(compressed) {
					break
				}

				// Read 2-byte back reference
				b1 := int(compressed[pos])
				b2 := int(compressed[pos+1])
				pos += 2

				// Extract length and offset
				length := ((b1 >> 4) & 0x0F) + 3
				offset := (((b1 & 0x0F) << 8) | b2) + 1

				// Copy from lookback buffer
				startPos := len(decompressed) - offset
				if startPos < 0 {
					return nil, fmt.Errorf("invalid back reference at pos %d", pos)
				}

				for i := 0; i < length; i++ {
					if startPos+i < len(decompressed) {
						decompressed = append(decompressed, decompressed[startPos+i])
					}
				}
			}
		}
	}

	return decompressed, nil
}
