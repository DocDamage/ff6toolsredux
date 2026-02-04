// Package share provides share code generation and parsing.
//
// The share package handles:
//   - Generating compact share codes
//   - Parsing share codes
//   - Code validation
//
// Share codes allow users to share:
//   - Character builds
//   - Party presets
//   - Configuration settings
//
// Example usage:
//
//	// Generate share code
//	code, err := share.GenerateCode(data)
//	if err != nil {
//	    return err
//	}
//
//	// Parse share code
//	data, err := share.ParseCode(code)
//	if err != nil {
//	    return err
//	}
package share
