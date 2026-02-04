// Package speedrun provides speedrun configuration and timing.
//
// The speedrun package handles:
//   - Speedrun category configurations
//   - Split timing
//   - Route optimization
//   - Validation rules
//
// Supported categories:
//   - Any%
//   - 100%
//   - Glitchless
//   - Specific character routes
//
// Example usage:
//
//	// Load speedrun config
//	cfg := speedrun.LoadConfig("any_percent")
//
//	// Apply to save
//	if err := cfg.ApplyToSave(save); err != nil {
//	    return err
//	}
package speedrun
