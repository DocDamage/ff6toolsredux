package theme

import "image/color"

// GetDarkTheme returns the dark theme configuration
func GetDarkTheme() *Theme {
	return &Theme{
		Name: "Dark",

		// Primary colors - blue-ish tones
		PrimaryColor:   color.RGBA{0x42, 0x85, 0xF4, 0xFF}, // Google Blue
		SecondaryColor: color.RGBA{0x34, 0xA8, 0x53, 0xFF}, // Green
		AccentColor:    color.RGBA{0xFB, 0xBC, 0x04, 0xFF}, // Amber
		ErrorColor:     color.RGBA{0xEA, 0x4B, 0x35, 0xFF}, // Red
		WarningColor:   color.RGBA{0xF5, 0x7C, 0x00, 0xFF}, // Orange
		SuccessColor:   color.RGBA{0x0F, 0x9D, 0x58, 0xFF}, // Light Green

		// Background colors - dark palette
		BackgroundColor: color.RGBA{0x12, 0x12, 0x12, 0xFF}, // Almost black
		SurfaceColor:    color.RGBA{0x1E, 0x1E, 0x1E, 0xFF}, // Slightly lighter
		CardColor:       color.RGBA{0x28, 0x28, 0x28, 0xFF}, // Card background
		HoverColor:      color.RGBA{0x2F, 0x2F, 0x2F, 0xFF}, // Hover state
		FocusColor:      color.RGBA{0x42, 0x85, 0xF4, 0x4D}, // Blue with transparency
		DisabledColor:   color.RGBA{0x42, 0x42, 0x42, 0xFF}, // Disabled state

		// Text colors - light on dark
		TextColor:       color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White
		TextSecondary:   color.RGBA{0xB0, 0xB0, 0xB0, 0xFF}, // Light gray
		TextDisabled:    color.RGBA{0x66, 0x66, 0x66, 0xFF}, // Dark gray
		TextOnPrimary:   color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White on blue
		TextOnSecondary: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White on green

		// Border colors
		BorderColor: color.RGBA{0x3F, 0x3F, 0x3F, 0xFF}, // Medium gray
		BorderLight: color.RGBA{0x4A, 0x4A, 0x4A, 0xFF}, // Lighter border
		BorderDark:  color.RGBA{0x2A, 0x2A, 0x2A, 0xFF}, // Darker border

		// Component specific
		InputBackground:  color.RGBA{0x2E, 0x2E, 0x2E, 0xFF}, // Input background
		InputBorder:      color.RGBA{0x3F, 0x3F, 0x3F, 0xFF}, // Input border
		InputText:        color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // Input text
		InputPlaceholder: color.RGBA{0x76, 0x76, 0x76, 0xFF}, // Placeholder
		ButtonBackground: color.RGBA{0x42, 0x85, 0xF4, 0xFF}, // Primary blue
		ButtonText:       color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White text

		// Sizing
		TextSize:      14,
		TextSizeSmall: 12,
		TextSizeLarge: 16,
		Padding:       10,
		Spacing:       8,
		CornerRadius:  4,
		BorderWidth:   1,

		// Display
		IsDark: true,
	}
}
