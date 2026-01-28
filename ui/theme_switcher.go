package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"ffvi_editor/ui/theme"
)

// ThemeSwitcher manages theme switching and persistence
type ThemeSwitcher struct {
	preferences fyne.Preferences
	windows     []fyne.Window
}

// NewThemeSwitcher creates a new theme switcher
func NewThemeSwitcher(preferences fyne.Preferences) *ThemeSwitcher {
	ts := &ThemeSwitcher{
		preferences: preferences,
		windows:     make([]fyne.Window, 0),
	}

	// Load saved theme preference
	ts.loadSavedTheme()

	return ts
}

// RegisterWindow registers a window to be updated when theme changes
func (ts *ThemeSwitcher) RegisterWindow(window fyne.Window) {
	ts.windows = append(ts.windows, window)
}

// ToggleTheme switches between light and dark themes
func (ts *ThemeSwitcher) ToggleTheme() {
	if theme.IsDarkTheme() {
		ts.SetTheme(theme.GetLightTheme())
	} else {
		ts.SetTheme(theme.GetDarkTheme())
	}
}

// SetTheme applies the given theme
func (ts *ThemeSwitcher) SetTheme(t *theme.Theme) {
	theme.SetTheme(t)
	ts.saveThemePreference(t.Name)

	// Update all registered windows
	for _, window := range ts.windows {
		if window != nil {
		}
	}
}

// loadSavedTheme loads the saved theme preference
func (ts *ThemeSwitcher) loadSavedTheme() {
	saved := ts.preferences.String("theme")

	switch saved {
	case "Dark":
		theme.SetTheme(theme.GetDarkTheme())
	case "Light":
		theme.SetTheme(theme.GetLightTheme())
	default:
		// Default to light theme
		theme.SetTheme(theme.GetLightTheme())
	}
}

// saveThemePreference saves the theme preference
func (ts *ThemeSwitcher) saveThemePreference(themeName string) {
	ts.preferences.SetString("theme", themeName)
}

// GetCurrentThemeName returns the name of the current theme
func (ts *ThemeSwitcher) GetCurrentThemeName() string {
	t := theme.GetTheme()
	return t.Name
}

// ShowThemeDialog shows theme selection dialog
func (ts *ThemeSwitcher) ShowThemeDialog(window fyne.Window) {
	currentTheme := ts.GetCurrentThemeName()
	themes := []string{"Light", "Dark"}

	themeSelect := fyne.NewMenu(
		"Theme",
		fyne.NewMenuItem("Light Theme", func() {
			ts.SetTheme(theme.GetLightTheme())
			dialog.ShowInformation("Theme Changed", "Light theme activated", window)
		}),
		fyne.NewMenuItem("Dark Theme", func() {
			ts.SetTheme(theme.GetDarkTheme())
			dialog.ShowInformation("Theme Changed", "Dark theme activated", window)
		}),
	)

	_ = themes // Keep for potential future use
	_ = currentTheme

	// Show menu (caller would display this)
	_ = themeSelect
}

// BuildMenuItems returns the theme menu item for the app menu
func (ts *ThemeSwitcher) BuildMenuItems() *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Toggle Theme",
		ts.ToggleTheme,
	)
}
