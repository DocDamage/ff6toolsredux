package share

import (
	"fmt"
	"strings"
	"time"
)

// ShareCode represents a shareable build code
type ShareCode struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	BuildType   string    `json:"buildType"` // "character", "party", "full"
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Data        string    `json:"data"` // Encoded build data
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	Downloads   int       `json:"downloads"`
}

// CharacterBuild represents an exported character for sharing
type CharacterBuild struct {
	Name      string
	Level     uint8
	HP        uint16
	MP        uint8
	Stats     CharacterStats
	Equipment CharacterEquipment
	Magic     []string
	Commands  []string
}

// CharacterStats represents character stats for sharing
type CharacterStats struct {
	Vigor    uint8
	Speed    uint8
	Stamina  uint8
	MagicPwr uint8
	Defense  uint8
	MagicDef uint8
}

// CharacterEquipment represents character equipment for sharing
type CharacterEquipment struct {
	Weapon string
	Armor  string
	Shield string
	Helmet string
	Relic1 string
	Relic2 string
}

// PartyBuild represents an exported party for sharing
type PartyBuild struct {
	Members    []string
	Characters []CharacterBuild
}

// CodeGenerator generates shareable codes for builds
type CodeGenerator struct {
	baseURL string
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		baseURL: "ffvi.build",
	}
}

// GenerateCharacterCode generates a shareable code for a character build
// character should be a CharacterBuild struct
func (g *CodeGenerator) GenerateCharacterCode(build *CharacterBuild, maxCodeLength int) (string, error) {
	if build == nil {
		return "", fmt.Errorf("build is nil")
	}

	// Encode build to compact format
	encoded := encodeCharacterBuild(*build)

	// Generate code prefix and suffix
	prefix := strings.ToUpper(build.Name[:min(4, len(build.Name))])
	code := fmt.Sprintf("%s-BLD-%s", prefix, encoded)

	// Truncate if needed
	if maxCodeLength > 0 && len(code) > maxCodeLength {
		// Keep prefix and truncate encoded part
		maxEncoded := maxCodeLength - len(prefix) - 6 // Account for "-BLD-"
		if maxEncoded > 10 {
			code = fmt.Sprintf("%s-BLD-%s", prefix, encoded[:maxEncoded])
		}
	}

	return code, nil
}

// GeneratePartyCode generates a shareable code for a party
func (g *CodeGenerator) GeneratePartyCode(party *PartyBuild) (string, error) {
	if party == nil || len(party.Members) == 0 {
		return "", fmt.Errorf("party is empty")
	}

	// Encode party to compact format
	encoded := encodePartyBuild(*party)

	// Generate code
	code := fmt.Sprintf("PTY-%s", encoded)

	return code, nil
}

// DecodeCharacterCode decodes a shareable character code
func (g *CodeGenerator) DecodeCharacterCode(code string) (*CharacterBuild, error) {
	// Extract encoded part (remove prefix)
	parts := strings.Split(code, "-")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid code format: %s", code)
	}

	encoded := strings.Join(parts[2:], "-")

	// Decode build
	build, err := decodeCharacterBuild(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode character build: %w", err)
	}

	return build, nil
}

// DecodePartyCode decodes a shareable party code
func (g *CodeGenerator) DecodePartyCode(code string) (*PartyBuild, error) {
	// Extract encoded part
	if !strings.HasPrefix(code, "PTY-") {
		return nil, fmt.Errorf("invalid party code format: %s", code)
	}

	encoded := strings.TrimPrefix(code, "PTY-")

	// Decode build
	build, err := decodePartyBuild(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode party build: %w", err)
	}

	return build, nil
}

// Encoding functions using simple base32-like encoding

// encodeCharacterBuild encodes a character build to compact format
func encodeCharacterBuild(build CharacterBuild) string {
	// Create encoding: name(4) + level(1) + stats(6) + equip(6 short names)
	// Using base32 alphabet for URL safety
	encoded := ""

	// Encode name (first 4 chars, uppercase)
	name := strings.ToUpper(build.Name)
	if len(name) > 4 {
		name = name[:4]
	}
	encoded += name

	// Encode level (2 hex digits)
	encoded += fmt.Sprintf("%02X", build.Level)

	// Encode stats (each stat is 1 hex digit, 0-9, A-F = 0-15)
	encoded += fmt.Sprintf("%X%X%X%X%X%X",
		build.Stats.Vigor/17,   // 0-5
		build.Stats.Speed/17,
		build.Stats.Stamina/17,
		build.Stats.MagicPwr/17,
		build.Stats.Defense/17,
		build.Stats.MagicDef/17,
	)

	// Encode equipment (simplified: just use first char of each piece)
	equip := ""
	if build.Equipment.Weapon != "" {
		equip += string(build.Equipment.Weapon[0])
	} else {
		equip += "-"
	}
	if build.Equipment.Armor != "" {
		equip += string(build.Equipment.Armor[0])
	} else {
		equip += "-"
	}
	if build.Equipment.Helmet != "" {
		equip += string(build.Equipment.Helmet[0])
	} else {
		equip += "-"
	}
	if build.Equipment.Relic1 != "" {
		equip += string(build.Equipment.Relic1[0])
	} else {
		equip += "-"
	}

	encoded += equip

	return encoded
}

// decodeCharacterBuild decodes a character build from encoded format
func decodeCharacterBuild(encoded string) (*CharacterBuild, error) {
	if len(encoded) < 16 {
		return nil, fmt.Errorf("encoded data too short")
	}

	build := &CharacterBuild{}

	// Decode name
	build.Name = encoded[:4]

	// Decode level
	var level int
	_, err := fmt.Sscanf(encoded[4:6], "%X", &level)
	if err != nil {
		return nil, fmt.Errorf("failed to decode level: %w", err)
	}
	build.Level = uint8(level)

	// Decode stats
	var v, s, st, m, d, md int
	_, err = fmt.Sscanf(encoded[6:12], "%X%X%X%X%X%X", &v, &s, &st, &m, &d, &md)
	if err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	build.Stats.Vigor = uint8(v * 17)
	build.Stats.Speed = uint8(s * 17)
	build.Stats.Stamina = uint8(st * 17)
	build.Stats.MagicPwr = uint8(m * 17)
	build.Stats.Defense = uint8(d * 17)
	build.Stats.MagicDef = uint8(md * 17)

	return build, nil
}

// encodePartyBuild encodes a party build to compact format
func encodePartyBuild(build PartyBuild) string {
	encoded := ""

	// Encode each character name (first 3 chars)
	for i, char := range build.Characters {
		if i < 4 { // Max 4 party members
			name := strings.ToUpper(char.Name)
			if len(name) > 3 {
				name = name[:3]
			}
			encoded += name
			encoded += fmt.Sprintf("%02X", char.Level)
		}
	}

	return encoded
}

// decodePartyBuild decodes a party build from encoded format
func decodePartyBuild(encoded string) (*PartyBuild, error) {
	if len(encoded) < 5 {
		return nil, fmt.Errorf("encoded data too short")
	}

	build := &PartyBuild{
		Members:    make([]string, 0),
		Characters: make([]CharacterBuild, 0),
	}

	// Decode each character
	for i := 0; i < len(encoded); i += 5 {
		if i+5 > len(encoded) {
			break
		}

		// Extract name (3 chars) and level (2 hex digits)
		name := encoded[i : i+3]
		build.Members = append(build.Members, name)

		var level int
		_, err := fmt.Sscanf(encoded[i+3:i+5], "%X", &level)
		if err != nil {
			continue
		}

		charBuild := CharacterBuild{
			Name:  name,
			Level: uint8(level),
		}
		build.Characters = append(build.Characters, charBuild)
	}

	return build, nil
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
