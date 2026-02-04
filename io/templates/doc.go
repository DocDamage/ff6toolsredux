// Package templates provides character template management.
//
// The templates package handles:
//   - Character build templates
//   - Template sharing and downloading
//   - Template application to characters
//
// Templates can include:
//   - Base stats and level
//   - Equipment loadouts
//   - Magic and ability sets
//   - Esper assignments
//
// Example usage:
//
//	// Create template manager
//	tm := templates.NewManager(templateDir)
//
//	// Apply template to character
//	if err := tm.ApplyTemplate(char, "Ultima Mage"); err != nil {
//	    return err
//	}
package templates
