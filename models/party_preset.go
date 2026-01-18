package models

import (
	"errors"
	"time"
)

// PartyPreset represents a saved party configuration
type PartyPreset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Members     [4]uint8  `json:"members"`      // Character indices (0-13)
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []string  `json:"tags"`
	Favorite    bool      `json:"favorite"`
	Version     int       `json:"version"`
}

// Validate checks if the preset is valid
func (p *PartyPreset) Validate() error {
	if p.Name == "" {
		return errors.New("preset name cannot be empty")
	}

	for i, member := range p.Members {
		if member > 13 {
			return errors.New("invalid character index in party member " + string(rune(i)))
		}
	}

	return nil
}

// IsComplete checks if all party slots are filled
func (p *PartyPreset) IsComplete() bool {
	for _, member := range p.Members {
		if member == 0 && p.Members[0] != 0 { // Skip first member if it's 0
			return false
		}
	}
	return true
}
