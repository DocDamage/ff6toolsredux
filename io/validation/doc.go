// Package validation provides save file validation.
//
// The validation package handles:
//   - Save file integrity checks
//   - Data consistency validation
//   - Error detection and reporting
//   - Auto-fix capabilities
//
// Validation levels:
//   - Strict: Reject any invalid data
//   - Normal: Allow minor issues, warn on major
//   - Permissive: Allow most issues, just log
//
// Example usage:
//
//	// Validate a save
//	validator := validation.NewValidator(validation.Normal)
//	results, err := validator.Validate(save)
//	if err != nil {
//	    return err
//	}
//
//	// Check for issues
//	if results.HasErrors() {
//	    for _, issue := range results.Issues {
//	        log.Printf("Issue: %s", issue.Description)
//	    }
//	}
package validation
