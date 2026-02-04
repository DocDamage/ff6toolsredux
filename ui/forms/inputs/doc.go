// Package inputs provides specialized input widgets.
//
// The inputs package contains:
//   - Integer entry fields
//   - Float entry fields
//   - Labeled inputs
//   - Validation input fields
//
// These are reusable input components with
// built-in validation and formatting.
//
// Example usage:
//
//	// Create integer entry
//	entry := inputs.NewIntEntry()
//	entry.SetValue(100)
//	entry.SetRange(0, 9999)
//
//	// Create labeled input
//	labeled := inputs.NewLabeledInput("HP:", entry)
//	labeled.SetHelpText("Current HP value")
package inputs
