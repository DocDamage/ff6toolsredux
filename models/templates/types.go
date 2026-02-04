package templates

import "errors"

// ApplyMode defines how a template should be applied to a character
type ApplyMode int

const (
	ApplyModeReplace ApplyMode = iota // Replace all fields
	ApplyModeMerge                    // Merge only stronger stats
	ApplyModePreview                  // Validation only
)

var (
	ErrInvalidTemplate  = errors.New("invalid template: missing character data")
	ErrTemplateNotFound = errors.New("template not found")
	ErrDuplicateName    = errors.New("template name already exists")
	ErrIOError          = errors.New("failed to read/write template file")
)
