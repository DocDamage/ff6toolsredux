// Package templates provides character template definitions.
//
// The templates package handles:
//   - Character build templates
//   - Template validation
//   - Template metadata
//
// Templates define:
//   - Stat distributions
//   - Equipment sets
//   - Magic loadouts
//   - Esper assignments
//
// Example usage:
//
//	// Create template
//	template := &templates.CharacterTemplate{
//	    Name: "Physical DPS",
//	    Stats: templates.StatDistribution{
//	        Vigor: 128,
//	        Speed: 50,
//	    },
//	}
//
//	// Apply to character
//	template.Apply(char)
package templates
