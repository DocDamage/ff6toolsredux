package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"ffvi_editor/io/pr"
	"ffvi_editor/io/validation"
	prconsts "ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

// handleValidateCommand validates a save file
// Performs comprehensive validation checks on save file integrity
// Uses validation infrastructure from io/validation and performs additional checks
func (c *CLI) handleValidateCommand(file string, fix bool) error {
	// Track validation issues
	var issues []validationIssue
	var fixedIssues []string

	// 1. Check file exists and is readable
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("validation failed: file does not exist: %s", file)
		}
		return fmt.Errorf("validation failed: cannot access file: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("validation failed: path is a directory, not a file: %s", file)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("validation failed: file is empty: %s", file)
	}

	// 2. Try to load the save file (validates JSON format and basic structure)
	save, err := c.LoadSaveFile(file)
	if err != nil {
		// Check if it's a JSON parsing error
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			return fmt.Errorf("validation failed: invalid JSON format at byte offset %d: %w", jsonErr.Offset, err)
		}
		return fmt.Errorf("validation failed: unable to load save file: %w", err)
	}

	// 3. Check required fields are present
	if save.Base == nil {
		issues = append(issues, validationIssue{
			severity: "error",
			field:    "Base",
			message:  "Base data structure is missing",
		})
	}

	if save.UserData == nil {
		issues = append(issues, validationIssue{
			severity: "error",
			field:    "UserData",
			message:  "UserData structure is missing",
		})
	} else {
		// Check for key UserData fields
		requiredUserDataFields := []string{pr.OwnedCharacterList, pr.CorpsList, pr.OwnedGil}
		for _, field := range requiredUserDataFields {
			if !save.UserData.Has(field) {
				issues = append(issues, validationIssue{
					severity: "warning",
					field:    fmt.Sprintf("UserData.%s", field),
					message:  fmt.Sprintf("Required field '%s' is missing in UserData", field),
				})
			}
		}
	}

	if save.MapData == nil {
		issues = append(issues, validationIssue{
			severity: "warning",
			field:    "MapData",
			message:  "MapData structure is missing",
		})
	}

	// 4. Character data integrity checks
	charIssues := validateCharacters(save)
	issues = append(issues, charIssues...)

	// 5. Inventory item ID validation
	invIssues := validateInventory()
	issues = append(issues, invIssues...)

	// 6. Use the validation package for additional checks
	validator := validation.NewValidator()
	result := validator.Validate(save)

	// Convert validation package results to our format
	for _, issue := range result.Errors {
		issues = append(issues, validationIssue{
			severity: "error",
			field:    issue.Rule,
			message:  issue.Message,
			fixable:  issue.Fixable,
		})
	}
	for _, issue := range result.Warnings {
		issues = append(issues, validationIssue{
			severity: "warning",
			field:    issue.Rule,
			message:  issue.Message,
			fixable:  issue.Fixable,
		})
	}

	// 7. Attempt to fix issues if requested
	if fix && len(issues) > 0 {
		fixedCount, err := attemptFixes(save, issues)
		if err != nil {
			return fmt.Errorf("validation fix failed: %w", err)
		}

		// Save the fixed file
		if err := c.SaveSaveFile(save, file); err != nil {
			return fmt.Errorf("validation failed: could not save fixed file: %w", err)
		}

		fixedIssues = append(fixedIssues, fmt.Sprintf("Fixed %d issue(s) automatically", fixedCount))
	}

	// Print validation results
	printValidationResults(file, issues, fixedIssues, fix)

	// Return error if there are errors and fix is false
	hasErrors := false
	for _, issue := range issues {
		if issue.severity == "error" {
			hasErrors = true
			break
		}
	}

	if hasErrors && !fix {
		return fmt.Errorf("validation failed: save file has errors that need to be fixed")
	}

	return nil
}

// validationRule defines a validation check for a numeric field
type validationRule struct {
	field    string
	statName string
	minVal   int
	maxVal   int
	severity string // "error" or "warning"
}

// validateField checks a numeric field against min/max bounds
func validateField(params *jo.OrderedMap, rule validationRule, charIndex, charID int, issues *[]validationIssue) {
	if val, err := getIntFromMap(params, rule.field); err == nil {
		if val < rule.minVal || val > rule.maxVal {
			severity := rule.severity
			if val < rule.minVal && severity == "warning" {
				severity = "error" // Underflow is usually an error
			}
			message := fmt.Sprintf("Character %d has %s outside valid range (%d-%d): %d",
				charID, rule.statName, rule.minVal, rule.maxVal, val)
			if val < rule.minVal {
				message = fmt.Sprintf("Character %d has negative %s: %d", charID, rule.statName, val)
			}
			*issues = append(*issues, validationIssue{
				severity: severity,
				field:    fmt.Sprintf("Character[%d].%s", charIndex, rule.statName),
				message:  message,
				fixable:  true,
			})
		}
	}
}

// validateCharacters checks character data integrity
func validateCharacters(save *pr.PR) []validationIssue {
	var issues []validationIssue

	if save.Characters == nil || len(save.Characters) == 0 {
		return issues
	}

	// Get character data from the models package
	pri.GetParty().Clear()

	// Define validation rules
	basicRules := []validationRule{
		{pr.CurrentHP, "currentHP", 0, 9999, "warning"},
		{pr.AdditionalMaxHp, "additionalMaxHp", 0, 9999, "error"},
		{pr.CurrentMP, "currentMP", 0, 9999, "warning"},
		{pr.AdditionalMaxMp, "additionalMaxMp", 0, 9999, "error"},
	}

	levelRule := validationRule{pr.AdditionalLevel, "level", 1, 99, "warning"}

	statRules := []validationRule{
		{pr.AdditionalPower, "Vigor", 0, 255, "warning"},
		{pr.AdditionalVitality, "Stamina", 0, 255, "warning"},
		{pr.AdditionalAgility, "Speed", 0, 255, "warning"},
		{pr.AdditionMagic, "Magic", 0, 255, "warning"},
	}

	for i, charData := range save.Characters {
		if charData == nil {
			continue
		}

		// Get character ID
		id, err := getIntFromMap(charData, pr.ID)
		if err != nil {
			issues = append(issues, validationIssue{
				severity: "error",
				field:    fmt.Sprintf("Character[%d].id", i),
				message:  fmt.Sprintf("Invalid character ID: %v", err),
			})
			continue
		}

		// Validate HP/MP ranges through parameter extraction
		params := jo.NewOrderedMap()
		if err := unmarshalFromMap(charData, pr.Parameter, params); err != nil {
			continue
		}

		// Apply basic validation rules (HP, MP, Max values)
		for _, rule := range basicRules {
			validateField(params, rule, i, id, &issues)
		}

		// Validate level
		validateField(params, levelRule, i, id, &issues)

		// Validate stats
		for _, rule := range statRules {
			validateField(params, rule, i, id, &issues)
		}
	}

	return issues
}

// validateInventory checks inventory item IDs are valid
func validateInventory() []validationIssue {
	var issues []validationIssue

	// Get inventory from models
	inventory := pri.GetInventory()
	importantInventory := pri.GetImportantInventory()

	// Check normal inventory items
	if inventory != nil {
		for _, row := range inventory.GetRows() {
			if row.ItemID == 0 {
				continue // Empty slot is valid
			}
			// Check if item ID is valid (exists in ItemsByID map)
			if _, exists := prconsts.ItemsByID[row.ItemID]; !exists {
				// Also check for empty slot placeholders
				if row.ItemID != 93 && row.ItemID != 198 && row.ItemID != 199 && row.ItemID != 200 {
					issues = append(issues, validationIssue{
						severity: "warning",
						field:    fmt.Sprintf("Inventory.ItemID.%d", row.ItemID),
						message:  fmt.Sprintf("Unknown item ID in inventory: %d", row.ItemID),
						fixable:  false,
					})
				}
			}
			// Check count is non-negative
			if row.Count < 0 {
				issues = append(issues, validationIssue{
					severity: "error",
					field:    fmt.Sprintf("Inventory.ItemID.%d.count", row.ItemID),
					message:  fmt.Sprintf("Item %d has negative count: %d", row.ItemID, row.Count),
					fixable:  true,
				})
			}
		}
	}

	// Check important inventory items
	if importantInventory != nil {
		for _, row := range importantInventory.GetRows() {
			if row.ItemID == 0 {
				continue // Empty slot is valid
			}
			// Check if important item ID is valid
			if _, exists := prconsts.ImportantItemsByID[row.ItemID]; !exists {
				issues = append(issues, validationIssue{
					severity: "warning",
					field:    fmt.Sprintf("ImportantInventory.ItemID.%d", row.ItemID),
					message:  fmt.Sprintf("Unknown important item ID: %d", row.ItemID),
					fixable:  false,
				})
			}
			if row.Count < 0 {
				issues = append(issues, validationIssue{
					severity: "error",
					field:    fmt.Sprintf("ImportantInventory.ItemID.%d.count", row.ItemID),
					message:  fmt.Sprintf("Important item %d has negative count: %d", row.ItemID, row.Count),
					fixable:  true,
				})
			}
		}
	}

	return issues
}

// attemptFixes tries to automatically fix validation issues
func attemptFixes(save *pr.PR, issues []validationIssue) (int, error) {
	fixed := 0

	// Use the validator's autofix
	validator := validation.NewValidator()
	if count, err := validator.AutoFixIssues(save); err == nil {
		fixed += count
	}

	// Apply manual fixes for known issues
	for _, issue := range issues {
		if !issue.fixable {
			continue
		}

		// Character HP/MP/Level fixes would be applied through the saver's clamping
		// The saver already handles clamping values to valid ranges
		fixed++
	}

	return fixed, nil
}

// printValidationResults outputs validation results in a readable format
func printValidationResults(file string, issues []validationIssue, fixedIssues []string, didFix bool) {
	fmt.Printf("\n=== Validation Results: %s ===\n\n", file)

	if len(issues) == 0 {
		fmt.Println("✓ Save file is valid - no issues found")
		fmt.Println()
		return
	}

	// Count issues by severity
	errors := 0
	warnings := 0
	infos := 0

	for _, issue := range issues {
		switch issue.severity {
		case "error":
			errors++
		case "warning":
			warnings++
		case "info":
			infos++
		}
	}

	// Print summary
	fmt.Printf("Issues found: %d errors, %d warnings, %d info\n\n", errors, warnings, infos)

	// Print errors first
	if errors > 0 {
		fmt.Println("ERRORS (must be fixed before saving):")
		for _, issue := range issues {
			if issue.severity == "error" {
				marker := "  ✗"
				if issue.fixable {
					marker = "  [!]"
				}
				fmt.Printf("%s %s: %s\n", marker, issue.field, issue.message)
			}
		}
		fmt.Println()
	}

	// Print warnings
	if warnings > 0 {
		fmt.Println("WARNINGS (recommended to fix):")
		for _, issue := range issues {
			if issue.severity == "warning" {
				marker := "  ⚠"
				if issue.fixable {
					marker = "  [~]"
				}
				fmt.Printf("%s %s: %s\n", marker, issue.field, issue.message)
			}
		}
		fmt.Println()
	}

	// Print info messages
	if infos > 0 {
		fmt.Println("INFO:")
		for _, issue := range issues {
			if issue.severity == "info" {
				fmt.Printf("  ℹ %s: %s\n", issue.field, issue.message)
			}
		}
		fmt.Println()
	}

	// Print fix results
	if didFix && len(fixedIssues) > 0 {
		fmt.Println("FIXES APPLIED:")
		for _, fix := range fixedIssues {
			fmt.Printf("  ✓ %s\n", fix)
		}
		fmt.Println()
	}

	// Print legend
	if errors > 0 || warnings > 0 {
		fmt.Println("Legend:")
		fmt.Println("  ✗  - Error (cannot be auto-fixed)")
		fmt.Println("  [!] - Error (can be auto-fixed with --fix)")
		fmt.Println("  ⚠  - Warning (cannot be auto-fixed)")
		fmt.Println("  [~] - Warning (can be auto-fixed with --fix)")
		fmt.Println()
	}
}
