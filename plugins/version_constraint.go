package plugins

import (
	"fmt"
	"strconv"
	"strings"
)

// VersionConstraint represents a version requirement
type VersionConstraint struct {
	Original string   // Original constraint string (e.g., "^1.5.0")
	Min      *Version // Minimum version (inclusive)
	Max      *Version // Maximum version (exclusive)
	Operator string   // Operator: "==", ">=", "<=", ">", "<", "~", "^", "*"
	Exact    *Version // For exact matches
}

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string // Pre-release identifier (e.g., "alpha", "beta.1")
	Build string // Build metadata
}

// ParseVersion parses a semantic version string
func ParseVersion(s string) (*Version, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Remove leading 'v' if present
	s = strings.TrimPrefix(s, "v")

	// Split on '+' for build metadata
	parts := strings.SplitN(s, "+", 2)
	version := &Version{}
	if len(parts) == 2 {
		version.Build = parts[1]
	}

	// Split on '-' for pre-release
	parts = strings.SplitN(parts[0], "-", 2)
	if len(parts) == 2 {
		version.Pre = parts[1]
	}

	// Parse major.minor.patch
	nums := strings.Split(parts[0], ".")
	if len(nums) < 1 || len(nums) > 3 {
		return nil, fmt.Errorf("invalid version format: %s", s)
	}

	var err error
	version.Major, err = strconv.Atoi(nums[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}

	if len(nums) >= 2 {
		version.Minor, err = strconv.Atoi(nums[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %w", err)
		}
	}

	if len(nums) >= 3 {
		version.Patch, err = strconv.Atoi(nums[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %w", err)
		}
	}

	return version, nil
}

// String returns the version as a string
func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}

// Compare compares two versions
// Returns: -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Compare pre-release versions
	if v.Pre != other.Pre {
		// Version without pre-release is greater
		if v.Pre == "" {
			return 1
		}
		if other.Pre == "" {
			return -1
		}
		// Lexicographic comparison of pre-release
		if v.Pre < other.Pre {
			return -1
		}
		return 1
	}

	return 0
}

// LessThan returns true if v < other
func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}

// GreaterThan returns true if v > other
func (v *Version) GreaterThan(other *Version) bool {
	return v.Compare(other) > 0
}

// Equal returns true if v == other
func (v *Version) Equal(other *Version) bool {
	return v.Compare(other) == 0
}

// ParseConstraint parses a version constraint string
func ParseConstraint(s string) (*VersionConstraint, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "*" {
		return &VersionConstraint{
			Original: s,
			Operator: "*",
		}, nil
	}

	constraint := &VersionConstraint{Original: s}

	// Handle operators
	if strings.HasPrefix(s, "^") {
		// Caret: ^1.2.3 allows >=1.2.3 <2.0.0
		constraint.Operator = "^"
		version, err := ParseVersion(s[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid caret constraint: %w", err)
		}
		constraint.Min = version
		constraint.Max = &Version{Major: version.Major + 1, Minor: 0, Patch: 0}
		return constraint, nil

	} else if strings.HasPrefix(s, "~") {
		// Tilde: ~1.2.3 allows >=1.2.3 <1.3.0
		constraint.Operator = "~"
		version, err := ParseVersion(s[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid tilde constraint: %w", err)
		}
		constraint.Min = version
		constraint.Max = &Version{Major: version.Major, Minor: version.Minor + 1, Patch: 0}
		return constraint, nil

	} else if strings.HasPrefix(s, ">=") {
		constraint.Operator = ">="
		version, err := ParseVersion(s[2:])
		if err != nil {
			return nil, fmt.Errorf("invalid >= constraint: %w", err)
		}
		constraint.Min = version
		return constraint, nil

	} else if strings.HasPrefix(s, "<=") {
		constraint.Operator = "<="
		version, err := ParseVersion(s[2:])
		if err != nil {
			return nil, fmt.Errorf("invalid <= constraint: %w", err)
		}
		constraint.Max = version
		return constraint, nil

	} else if strings.HasPrefix(s, ">") {
		constraint.Operator = ">"
		version, err := ParseVersion(s[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid > constraint: %w", err)
		}
		// Greater than means min is next patch
		constraint.Min = &Version{
			Major: version.Major,
			Minor: version.Minor,
			Patch: version.Patch + 1,
		}
		return constraint, nil

	} else if strings.HasPrefix(s, "<") {
		constraint.Operator = "<"
		version, err := ParseVersion(s[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid < constraint: %w", err)
		}
		constraint.Max = version
		return constraint, nil

	} else if strings.HasPrefix(s, "==") || strings.HasPrefix(s, "=") {
		// Exact match
		constraint.Operator = "=="
		versionStr := strings.TrimPrefix(s, "==")
		versionStr = strings.TrimPrefix(versionStr, "=")
		version, err := ParseVersion(versionStr)
		if err != nil {
			return nil, fmt.Errorf("invalid exact constraint: %w", err)
		}
		constraint.Exact = version
		return constraint, nil

	} else {
		// No operator means exact match
		constraint.Operator = "=="
		version, err := ParseVersion(s)
		if err != nil {
			return nil, fmt.Errorf("invalid version: %w", err)
		}
		constraint.Exact = version
		return constraint, nil
	}
}

// Satisfies checks if a version satisfies the constraint
func (c *VersionConstraint) Satisfies(version *Version) bool {
	if c.Operator == "*" {
		return true
	}

	if c.Exact != nil {
		return version.Equal(c.Exact)
	}

	// Check minimum
	if c.Min != nil {
		if c.Operator == ">=" || c.Operator == "^" || c.Operator == "~" {
			if version.LessThan(c.Min) {
				return false
			}
		} else if c.Operator == ">" {
			if !version.GreaterThan(c.Min) {
				return false
			}
		}
	}

	// Check maximum
	if c.Max != nil {
		if c.Operator == "<=" {
			if version.GreaterThan(c.Max) {
				return false
			}
		} else if c.Operator == "<" || c.Operator == "^" || c.Operator == "~" {
			if !version.LessThan(c.Max) {
				return false
			}
		}
	}

	return true
}

// String returns the constraint as a string
func (c *VersionConstraint) String() string {
	return c.Original
}

// ParseConstraintRange parses a range constraint like ">=1.0.0,<2.0.0"
func ParseConstraintRange(s string) ([]*VersionConstraint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty constraint range")
	}

	parts := strings.Split(s, ",")
	constraints := make([]*VersionConstraint, 0, len(parts))

	for _, part := range parts {
		constraint, err := ParseConstraint(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid constraint '%s': %w", part, err)
		}
		constraints = append(constraints, constraint)
	}

	return constraints, nil
}

// SatisfiesAll checks if a version satisfies all constraints in a range
func SatisfiesAll(version *Version, constraints []*VersionConstraint) bool {
	for _, constraint := range constraints {
		if !constraint.Satisfies(version) {
			return false
		}
	}
	return true
}

// FindCompatibleVersion finds the highest version that satisfies all constraints
func FindCompatibleVersion(versions []*Version, constraints []*VersionConstraint) *Version {
	var best *Version

	for _, version := range versions {
		if SatisfiesAll(version, constraints) {
			if best == nil || version.GreaterThan(best) {
				best = version
			}
		}
	}

	return best
}
