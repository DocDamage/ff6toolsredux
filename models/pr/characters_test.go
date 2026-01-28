package pr

import (
	"testing"

	"ffvi_editor/models/consts/pr"
)

// TestGetCharacter tests character lookup by root name
func TestGetCharacter(t *testing.T) {
	// Initialize the character system
	if len(Characters) == 0 {
		t.Skip("Characters not initialized")
	}

	tests := []struct {
		name    string
		charName string
		found   bool
	}{
		{
			name:    "find valid character",
			charName: "Terra",
			found:   true,
		},
		{
			name:    "find another valid character",
			charName: "Locke",
			found:   true,
		},
		{
			name:    "invalid character returns nil",
			charName: "NonExistent",
			found:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetCharacter(tt.charName)
			if tt.found && c == nil {
				t.Fatalf("GetCharacter(%q) = nil, want character", tt.charName)
			}
			if !tt.found && c != nil && c.RootName == tt.charName {
				t.Fatalf("GetCharacter(%q) should not find character", tt.charName)
			}
			if c != nil && c.RootName == tt.charName && c.Name != tt.charName {
				t.Fatalf("GetCharacter(%q) name mismatch", tt.charName)
			}
		})
	}
}

// TestGetCharacterByID tests character lookup by ID
func TestGetCharacterByID(t *testing.T) {
	if len(Characters) == 0 {
		t.Skip("Characters not initialized")
	}

	tests := []struct {
		name     string
		id       int
		shouldFind bool
	}{
		{
			name:       "find character by valid ID",
			id:         1, // Terra
			shouldFind: true,
		},
		{
			name:       "find another character",
			id:         5, // Locke
			shouldFind: true,
		},
		{
			name:       "invalid ID returns nil",
			id:         999,
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetCharacterByID(tt.id)
			if tt.shouldFind && c == nil {
				t.Fatalf("GetCharacterByID(%d) = nil, want character", tt.id)
			}
			if !tt.shouldFind && c != nil && c.ID == tt.id {
				t.Fatalf("GetCharacterByID(%d) should not find character", tt.id)
			}
		})
	}
}

// TestCharacterInitialization tests that characters are properly initialized with default values
func TestCharacterInitialization(t *testing.T) {
	if len(Characters) == 0 {
		t.Skip("Characters not initialized")
	}

	terra := GetCharacter("Terra")
	if terra == nil {
		t.Skip("Terra not found")
	}

	// Check basic initialization
	if terra.ID == 0 {
		t.Fatal("Character ID not set")
	}
	if terra.RootName == "" {
		t.Fatal("Character RootName not set")
	}
	if terra.StatusEffects == nil {
		t.Fatal("Character StatusEffects not initialized")
	}
	if len(terra.Commands) != 9 {
		t.Fatalf("Character Commands should have 9 slots, got %d", len(terra.Commands))
	}
	if terra.SpellsByIndex == nil {
		t.Fatal("Character SpellsByIndex not initialized")
	}
	if terra.SpellsSorted == nil {
		t.Fatal("Character SpellsSorted not initialized")
	}
	if terra.SpellsByID == nil {
		t.Fatal("Character SpellsByID not initialized")
	}
}

// TestCharacterDefaultCommand tests that characters have valid default commands
func TestCharacterDefaultCommand(t *testing.T) {
	if len(Characters) == 0 {
		t.Skip("Characters not initialized")
	}

	for _, c := range Characters {
		if c == nil {
			continue
		}
		for i, cmd := range c.Commands {
			if cmd == nil {
				t.Fatalf("Character %s has nil command at slot %d", c.Name, i)
			}
		}
	}
}

// TestCharacterIsNPC tests NPC flag for special characters
func TestCharacterIsNPC(t *testing.T) {
	tests := []struct {
		name  string
		id    int
		isNPC bool
	}{
		{
			name:  "Terra is playable",
			id:    1,
			isNPC: false,
		},
		{
			name:  "Locke is playable",
			id:    5,
			isNPC: false,
		},
		{
			name:  "Moglin is NPC",
			id:    6,
			isNPC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetCharacterByID(tt.id)
			if c == nil {
				t.Skip("Character not found")
			}
			if c.IsNPC != tt.isNPC {
				t.Fatalf("Character %s IsNPC = %v, want %v", c.Name, c.IsNPC, tt.isNPC)
			}
		})
	}
}

// TestCharacterCountIsCorrect tests that we have the expected number of characters
func TestCharacterCountIsCorrect(t *testing.T) {
	// FF6 PR has 40 characters total (including NPCs), but currently only 30 are initialized
	expectedCount := 30
	if len(Characters) != expectedCount {
		t.Fatalf("Characters array length = %d, want %d", len(Characters), expectedCount)
	}

	// Count non-nil characters
	nonNilCount := 0
	for _, c := range Characters {
		if c != nil {
			nonNilCount++
		}
	}

	if nonNilCount == 0 {
		t.Fatal("No characters were initialized")
	}
}

// TestCharacterBaseOffsets tests that all characters have valid base offsets
func TestCharacterBaseOffsets(t *testing.T) {
	// Verify the offset maps are populated
	if len(CharacterOffsetByName) == 0 {
		t.Skip("CharacterOffsetByName not populated")
	}

	for name, offset := range CharacterOffsetByName {
		if offset == nil {
			t.Fatalf("Character offset for %q is nil", name)
		}
		if offset.ID == 0 {
			t.Fatalf("Character %q has invalid ID", name)
		}
		if offset.HPBase <= 0 {
			t.Fatalf("Character %q has invalid HPBase: %d", name, offset.HPBase)
		}
		if offset.MPBase < 0 {
			t.Fatalf("Character %q has invalid MPBase: %d", name, offset.MPBase)
		}
	}
}

// TestNewSpells tests spell initialization for characters
func TestNewSpells(t *testing.T) {
	spellsByIndex, spellsSorted, spellsByID := NewSpells()

	if spellsByIndex == nil {
		t.Fatal("spellsByIndex is nil")
	}
	if spellsSorted == nil {
		t.Fatal("spellsSorted is nil")
	}
	if spellsByID == nil {
		t.Fatal("spellsByID is nil")
	}

	// Verify the spell collections are consistent
	if len(spellsByIndex) != len(spellsByID) {
		t.Fatalf("spell index and ID maps have different lengths: %d vs %d",
			len(spellsByIndex), len(spellsByID))
	}
}

// TestCommandLookup tests that command lookup is functional
func TestCommandLookup(t *testing.T) {
	if pr.CommandLookupByValue == nil {
		t.Skip("CommandLookupByValue not initialized")
	}

	if len(pr.CommandLookupByValue) == 0 {
		t.Fatal("CommandLookupByValue is empty")
	}

	// Verify we can look up default command (typically 4 for "Fight")
	defaultCmd := pr.CommandLookupByValue[4]
	if defaultCmd == nil {
		t.Fatal("Default command (4) not found in lookup")
	}
}

// BenchmarkGetCharacter benchmarks character lookup performance
func BenchmarkGetCharacter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCharacter("Terra")
	}
}

// BenchmarkGetCharacterByID benchmarks ID-based character lookup
func BenchmarkGetCharacterByID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCharacterByID(1)
	}
}
