package plugins

import "context"

// ShowDialog shows a dialog message
func (a *APIImpl) ShowDialog(ctx context.Context, title, message string) error {
	if !a.HasPermission(CommonPermissions.UIDisplay) {
		return ErrInsufficientPermissions
	}
	return a.showDialogFn(title, message)
}

// ShowConfirm shows a confirmation dialog
func (a *APIImpl) ShowConfirm(ctx context.Context, title, message string) (bool, error) {
	if !a.HasPermission(CommonPermissions.UIDisplay) {
		return false, ErrInsufficientPermissions
	}
	return a.showConfirmFn(title, message), nil
}

// ShowInput shows an input dialog
func (a *APIImpl) ShowInput(ctx context.Context, prompt string) (string, error) {
	if !a.HasPermission(CommonPermissions.UIDisplay) {
		return "", ErrInsufficientPermissions
	}
	return a.showInputFn(prompt)
}

// Log writes a log message
func (a *APIImpl) Log(ctx context.Context, level string, message string) error {
	a.logger(level, message)
	return nil
}
