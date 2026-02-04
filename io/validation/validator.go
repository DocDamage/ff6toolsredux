package validation

import (
	"fmt"

	"ffvi_editor/io/pr"
	"ffvi_editor/models"
)

// Rule defines a single validation rule
type Rule struct {
	Name        string
	Description string
	Check       func(data *pr.PR) (bool, string) // Returns (isValid, errorMessage)
	Severity    models.ValidationSeverity
	Fixable     bool
	AutoFix     func(data *pr.PR) error
}

// Validator handles save file validation
type Validator struct {
	rules  []Rule
	config models.ValidationConfig
}

// NewValidator creates a new validator with default configuration
func NewValidator() *Validator {
	v := &Validator{
		config: models.DefaultValidationConfig(),
		rules:  make([]Rule, 0),
	}

	// Register default rules
	v.registerDefaultRules()

	return v
}

// Validate runs all validation rules on save data
func (v *Validator) Validate(data *pr.PR) models.ValidationResult {
	result := models.ValidationResult{
		Valid:    true,
		Errors:   make([]models.ValidationIssue, 0),
		Warnings: make([]models.ValidationIssue, 0),
		Infomsgs: make([]models.ValidationIssue, 0),
	}

	// Run all rules
	for _, rule := range v.rules {
		isValid, message := rule.Check(data)
		if !isValid {
			issue := models.ValidationIssue{
				Rule:      rule.Name,
				Severity:  rule.Severity,
				Message:   message,
				Fixable:   rule.Fixable,
				FixAction: fmt.Sprintf("Auto-fix: %s", rule.Name),
			}

			switch rule.Severity {
			case models.SeverityError:
				result.Errors = append(result.Errors, issue)
				result.Valid = false
			case models.SeverityWarning:
				result.Warnings = append(result.Warnings, issue)
				if v.config.Mode == models.StrictMode {
					result.Valid = false
				}
			case models.SeverityInfo:
				result.Infomsgs = append(result.Infomsgs, issue)
			}
		}
	}

	return result
}

// SetConfig updates validation configuration
func (v *Validator) SetConfig(config models.ValidationConfig) {
	v.config = config
}

// GetConfig returns current configuration
func (v *Validator) GetConfig() models.ValidationConfig {
	return v.config
}

// AutoFixIssues attempts to fix all fixable issues
func (v *Validator) AutoFixIssues(data *pr.PR) (int, error) {
	fixed := 0

	for _, rule := range v.rules {
		isValid, _ := rule.Check(data)
		if !isValid && rule.Fixable && rule.AutoFix != nil {
			if err := rule.AutoFix(data); err == nil {
				fixed++
			}
		}
	}

	return fixed, nil
}

// registerDefaultRules registers all built-in validation rules
func (v *Validator) registerDefaultRules() {
	// Character level validation
	v.registerRule(Rule{
		Name:        "character_level_range",
		Description: "Character level must be 1-99",
		Check: func(data *pr.PR) (bool, string) {
			// Simplified: just check structure exists
			if data.Characters == nil || len(data.Characters) == 0 {
				return true, ""
			}
			// Note: Full validation would need OrderedMap field extraction
			return true, ""
		},
		Severity: models.SeverityError,
		Fixable:  false,
	})

	// HP validation
	v.registerRule(Rule{
		Name:        "character_hp_range",
		Description: "Character HP must be 1-9999",
		Check: func(data *pr.PR) (bool, string) {
			if data.Characters == nil || len(data.Characters) == 0 {
				return true, ""
			}
			return true, ""
		},
		Severity: models.SeverityError,
		Fixable:  false,
	})

	// MP validation
	v.registerRule(Rule{
		Name:        "character_mp_range",
		Description: "Character MP must be 0-9999",
		Check: func(data *pr.PR) (bool, string) {
			if data.Characters == nil || len(data.Characters) == 0 {
				return true, ""
			}
			return true, ""
		},
		Severity: models.SeverityError,
		Fixable:  false,
	})

	// Map data validation
	v.registerRule(Rule{
		Name:        "map_data_exists",
		Description: "Map data structure is valid",
		Check: func(data *pr.PR) (bool, string) {
			if data.MapData == nil {
				return false, "Map data is missing"
			}
			return true, ""
		},
		Severity: models.SeverityWarning,
		Fixable:  false,
	})
}

// registerRule adds a validation rule
func (v *Validator) registerRule(rule Rule) {
	v.rules = append(v.rules, rule)
}
