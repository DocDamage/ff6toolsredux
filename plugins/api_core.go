package plugins

import ioPR "ffvi_editor/io/pr"

// NewAPIImpl creates a new API implementation
func NewAPIImpl(prData *ioPR.PR, permissions []string) *APIImpl {
	api := &APIImpl{
		prData:      prData,
		hooks:       make(map[string][]func(interface{}) error),
		settings:    make(map[string]interface{}),
		permissions: make(map[string]bool),
	}

	// Set permissions
	for _, perm := range permissions {
		api.permissions[perm] = true
	}

	// Set default logger
	api.logger = func(level, msg string) {
		// Default: silent
	}

	// Set default dialog functions (no-op)
	api.showDialogFn = func(title, message string) error {
		return nil
	}
	api.showConfirmFn = func(title, message string) bool {
		return false
	}
	api.showInputFn = func(prompt string) (string, error) {
		return "", nil
	}

	return api
}

// SetLogger sets the logging function
func (a *APIImpl) SetLogger(logger func(level, msg string)) {
	if logger != nil {
		a.logger = logger
	}
}

// SetDialogFunctions sets UI dialog functions
func (a *APIImpl) SetDialogFunctions(
	showDialog func(title, message string) error,
	showConfirm func(title, message string) bool,
	showInput func(prompt string) (string, error),
) {
	if showDialog != nil {
		a.showDialogFn = showDialog
	}
	if showConfirm != nil {
		a.showConfirmFn = showConfirm
	}
	if showInput != nil {
		a.showInputFn = showInput
	}
}

// HasPermission checks if the API has a specific permission
func (a *APIImpl) HasPermission(permission string) bool {
	return a.permissions[permission]
}

// GetSetting retrieves a setting value
func (a *APIImpl) GetSetting(key string) interface{} {
	return a.settings[key]
}

// SetSetting stores a setting value
func (a *APIImpl) SetSetting(key string, value interface{}) error {
	a.settings[key] = value
	return nil
}
