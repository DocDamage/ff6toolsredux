package theme

import "image/color"

// GetLightTheme returns the light theme configuration
func GetLightTheme() *Theme {
	return &Theme{
		Name: "Light",

		// Primary colors - blue-ish tones
		PrimaryColor:   color.RGBA{0x1F, 0x5A, 0xD8, 0xFF}, // Deep Blue
		SecondaryColor: color.RGBA{0x1B, 0x8A, 0x3E, 0xFF}, // Deep Green
		AccentColor:    color.RGBA{0xF5, 0x8C, 0x00, 0xFF}, // Amber
		ErrorColor:     color.RGBA{0xD3, 0x28, 0x14, 0xFF}, // Dark Red
		WarningColor:   color.RGBA{0xF5, 0x7C, 0x00, 0xFF}, // Orange
		SuccessColor:   color.RGBA{0x2E, 0x7D, 0x32, 0xFF}, // Forest Green

		// Background colors - light palette
		BackgroundColor: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White
		SurfaceColor:    color.RGBA{0xF5, 0xF5, 0xF5, 0xFF}, // Light gray
		CardColor:       color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // Card background
		HoverColor:      color.RGBA{0xED, 0xED, 0xED, 0xFF}, // Hover state
		FocusColor:      color.RGBA{0x1F, 0x5A, 0xD8, 0x1F}, // Blue with transparency
		DisabledColor:   color.RGBA{0xE0, 0xE0, 0xE0, 0xFF}, // Disabled state

		// Text colors - dark on light
		TextColor:       color.RGBA{0x21, 0x21, 0x21, 0xFF}, // Almost black
		TextSecondary:   color.RGBA{0x61, 0x61, 0x61, 0xFF}, // Medium gray
		TextDisabled:    color.RGBA{0x99, 0x99, 0x99, 0xFF}, // Light gray
		TextOnPrimary:   color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White on blue
		TextOnSecondary: color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White on green

		// Border colors
		BorderColor: color.RGBA{0xCC, 0xCC, 0xCC, 0xFF}, // Medium gray
		BorderLight: color.RGBA{0xE0, 0xE0, 0xE0, 0xFF}, // Lighter border
		BorderDark:  color.RGBA{0x99, 0x99, 0x99, 0xFF}, // Darker border

		// Component specific
		InputBackground:  color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // White input
		InputBorder:      color.RGBA{0xBD, 0xBD, 0xBD, 0xFF}, // Gray border
		InputText:        color.RGBA{0x21, 0x21, 0x21, 0xFF}, // Dark text
		InputPlaceholder: color.RGBA{0x99, 0x99, 0x99, 0xFF}, // Gray placeholder
		ButtonBackground: color.RGBA{0x1F, 0x5A, 0xD8, 0xFF}, // Primary blue
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
		IsDark: false,
	}
}
