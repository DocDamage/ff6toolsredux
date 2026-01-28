package theme

import "image/color"

// Theme defines colors and styling for the UI
type Theme struct {
	Name string

	// Primary colors
	PrimaryColor   color.Color
	SecondaryColor color.Color
	AccentColor    color.Color
	ErrorColor     color.Color
	WarningColor   color.Color
	SuccessColor   color.Color

	// Background colors
	BackgroundColor color.Color
	SurfaceColor    color.Color
	CardColor       color.Color
	HoverColor      color.Color
	FocusColor      color.Color
	DisabledColor   color.Color

	// Text colors
	TextColor       color.Color
	TextSecondary   color.Color
	TextDisabled    color.Color
	TextOnPrimary   color.Color
	TextOnSecondary color.Color

	// Border colors
	BorderColor color.Color
	BorderLight color.Color
	BorderDark  color.Color

	// Component specific
	InputBackground  color.Color
	InputBorder      color.Color
	InputText        color.Color
	InputPlaceholder color.Color
	ButtonBackground color.Color
	ButtonText       color.Color

	// Sizing
	TextSize      int
	TextSizeSmall int
	TextSizeLarge int
	Padding       int
	Spacing       int
	CornerRadius  int
	BorderWidth   int

	// Display
	IsDark bool
}

// CurrentTheme holds the active theme
var CurrentTheme *Theme

// SetTheme updates the current theme
func SetTheme(t *Theme) {
	CurrentTheme = t
}

// GetTheme returns the current theme
func GetTheme() *Theme {
	if CurrentTheme == nil {
		CurrentTheme = GetLightTheme()
	}
	return CurrentTheme
}

// ToggleTheme switches between light and dark themes
func ToggleTheme() {
	if CurrentTheme != nil && CurrentTheme.IsDark {
		SetTheme(GetLightTheme())
	} else {
		SetTheme(GetDarkTheme())
	}
}

// IsDarkTheme checks if current theme is dark
func IsDarkTheme() bool {
	theme := GetTheme()
	return theme.IsDark
}
