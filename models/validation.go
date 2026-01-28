package models

// ValidationSeverity represents how serious a validation issue is
type ValidationSeverity string

const (
	// SeverityError prevents saving
	SeverityError ValidationSeverity = "error"
	// SeverityWarning warns but allows saving
	SeverityWarning ValidationSeverity = "warning"
	// SeverityInfo informational only
	SeverityInfo ValidationSeverity = "info"
)

// ValidationIssue represents a single validation problem
type ValidationIssue struct {
	Rule      string              `json:"rule"`
	Severity  ValidationSeverity  `json:"severity"`
	Message   string              `json:"message"`
	Target    string              `json:"target"`        // What was being validated
	TargetID  interface{}         `json:"targetID"`      // Character index, item ID, etc
	Fixable   bool                `json:"fixable"`
	FixAction string              `json:"fixAction"`     // Description of fix
}

// ValidationResult contains all validation results
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationIssue   `json:"errors"`
	Warnings []ValidationIssue   `json:"warnings"`
	Infomsgs []ValidationIssue   `json:"infoMessages"`
}

// ValidationMode determines how strict validation is
type ValidationMode string

const (
	// StrictMode prevents any illegal values
	StrictMode ValidationMode = "strict"
	// NormalMode (default) prevents critical errors but allows warnings
	NormalMode ValidationMode = "normal"
	// LenientMode allows everything but shows warnings
	LenientMode ValidationMode = "lenient"
)

// ValidationConfig holds validation settings
type ValidationConfig struct {
	Mode                ValidationMode
	RealTimeValidation  bool   // Validate as user edits
	PreSaveValidation   bool   // Validate before save
	MaxCharacterLevel   uint16 // Configurable limits
	MaxCharacterHP      uint16
	MaxCharacterMP      uint16
	MaxStatValue        uint16
	AutoFixSimpleIssues bool
}

// DefaultValidationConfig returns safe defaults
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		Mode:               NormalMode,
		RealTimeValidation: true,
		PreSaveValidation:  true,
		MaxCharacterLevel:  99,
		MaxCharacterHP:     9999,
		MaxCharacterMP:     9999,
		MaxStatValue:       255,
		AutoFixSimpleIssues: false,
	}
}

// HasErrors returns true if there are any error-level issues
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are any warning-level issues
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// AllIssues returns all issues combined
func (vr *ValidationResult) AllIssues() []ValidationIssue {
	result := make([]ValidationIssue, 0, len(vr.Errors)+len(vr.Warnings)+len(vr.Infomsgs))
	result = append(result, vr.Errors...)
	result = append(result, vr.Warnings...)
	result = append(result, vr.Infomsgs...)
	return result
}

// FixableIssues returns only issues that can be auto-fixed
func (vr *ValidationResult) FixableIssues() []ValidationIssue {
	result := make([]ValidationIssue, 0)
	for _, issue := range vr.AllIssues() {
		if issue.Fixable {
			result = append(result, issue)
		}
	}
	return result
}
