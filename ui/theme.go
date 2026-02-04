package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ThemeType defines available theme options
type ThemeType string

const (
	ThemeDark      ThemeType = "dark"
	ThemeLight     ThemeType = "light"
	ThemeCyberpunk ThemeType = "cyberpunk"
	ThemeOcean     ThemeType = "ocean"
	ThemeForest    ThemeType = "forest"
	ThemeAmethyst  ThemeType = "amethyst"
)

// FF6Theme provides a custom theme inspired by Final Fantasy VI aesthetics
type FF6Theme struct {
	variant   fyne.ThemeVariant
	themeType ThemeType
}

// NewFF6Theme creates a new FF6-inspired theme
func NewFF6Theme(variant fyne.ThemeVariant) fyne.Theme {
	return &FF6Theme{variant: variant, themeType: ThemeDark}
}

// NewFF6ThemeWithType creates a theme with a specific type
func NewFF6ThemeWithType(themeType ThemeType) fyne.Theme {
	return &FF6Theme{variant: theme.VariantDark, themeType: themeType}
}

func (t *FF6Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Route to specific theme handler
	switch t.themeType {
	case ThemeCyberpunk:
		return t.cyberpunkColor(name)
	case ThemeOcean:
		return t.oceanColor(name)
	case ThemeForest:
		return t.forestColor(name)
	case ThemeAmethyst:
		return t.amethystColor(name)
	case ThemeLight:
		return t.lightColor(name)
	default:
		return t.darkColor(name)
	}
}

func (t *FF6Theme) darkColor(name fyne.ThemeColorName) color.Color {
	// Dark blue-black theme (original)
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 103, G: 58, B: 183, A: 255} // Deep purple
	case theme.ColorNameButton:
		return color.NRGBA{R: 33, G: 150, B: 243, A: 255} // Blue
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green
	case theme.ColorNameFocus:
		return color.NRGBA{R: 255, G: 193, B: 7, A: 255} // Gold
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 152, B: 0, A: 255} // Orange
	case theme.ColorNameError:
		return color.NRGBA{R: 244, G: 67, B: 54, A: 255} // Red
	case theme.ColorNameBackground:
		return color.NRGBA{R: 18, G: 18, B: 24, A: 255} // Dark blue-black
	case theme.ColorNameForeground:
		return color.NRGBA{R: 240, G: 240, B: 240, A: 255} // Light text
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 150}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 100, G: 100, B: 100, A: 255}
	default:
		return theme.DarkTheme().Color(name, theme.VariantDark)
	}
}

func (t *FF6Theme) lightColor(name fyne.ThemeColorName) color.Color {
	// Light theme
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 250, G: 250, B: 253, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 33, G: 33, B: 33, A: 255}
	default:
		return theme.LightTheme().Color(name, theme.VariantLight)
	}
}

func (t *FF6Theme) cyberpunkColor(name fyne.ThemeColorName) color.Color {
	// Cyberpunk: Neon pink, cyan, purple on dark background
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 13, G: 13, B: 20, A: 255} // Deep black-blue
	case theme.ColorNameForeground:
		return color.NRGBA{R: 255, G: 0, B: 255, A: 255} // Neon magenta text
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 255, B: 255, A: 255} // Cyan
	case theme.ColorNameButton:
		return color.NRGBA{R: 255, G: 20, B: 147, A: 255} // Deep pink
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 0, G: 255, B: 255, A: 255} // Cyan
	case theme.ColorNameFocus:
		return color.NRGBA{R: 255, G: 0, B: 255, A: 255} // Magenta
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0, G: 255, B: 127, A: 255} // Spring green
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Orange
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Bright red
	default:
		return color.NRGBA{R: 30, G: 30, B: 40, A: 255}
	}
}

func (t *FF6Theme) oceanColor(name fyne.ThemeColorName) color.Color {
	// Ocean: Deep blues, cyans, teals
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 15, G: 30, B: 50, A: 255} // Deep ocean blue
	case theme.ColorNameForeground:
		return color.NRGBA{R: 200, G: 240, B: 255, A: 255} // Light blue text
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 150, B: 200, A: 255} // Teal
	case theme.ColorNameButton:
		return color.NRGBA{R: 30, G: 180, B: 220, A: 255} // Bright cyan
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 100, G: 220, B: 255, A: 255} // Cyan
	case theme.ColorNameFocus:
		return color.NRGBA{R: 255, G: 200, B: 87, A: 255} // Golden
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 100, G: 255, B: 150, A: 255} // Light green
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Orange
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 100, B: 100, A: 255} // Light red
	default:
		return color.NRGBA{R: 30, G: 50, B: 70, A: 255}
	}
}

func (t *FF6Theme) forestColor(name fyne.ThemeColorName) color.Color {
	// Forest: Greens, earth tones
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 20, G: 35, B: 20, A: 255} // Deep forest green
	case theme.ColorNameForeground:
		return color.NRGBA{R: 220, G: 240, B: 200, A: 255} // Light green text
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green
	case theme.ColorNameButton:
		return color.NRGBA{R: 100, G: 200, B: 100, A: 255} // Bright green
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 144, G: 238, B: 144, A: 255} // Light green
	case theme.ColorNameFocus:
		return color.NRGBA{R: 255, G: 215, B: 0, A: 255} // Gold
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 100, G: 255, B: 100, A: 255} // Bright green
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Orange
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 100, B: 100, A: 255} // Light red
	default:
		return color.NRGBA{R: 35, G: 60, B: 35, A: 255}
	}
}

func (t *FF6Theme) amethystColor(name fyne.ThemeColorName) color.Color {
	// Amethyst: Purple, violet, lavender
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 25, G: 15, B: 40, A: 255} // Deep purple
	case theme.ColorNameForeground:
		return color.NRGBA{R: 230, G: 210, B: 255, A: 255} // Light lavender text
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 180, G: 100, B: 220, A: 255} // Purple
	case theme.ColorNameButton:
		return color.NRGBA{R: 200, G: 120, B: 240, A: 255} // Bright purple
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 220, G: 160, B: 255, A: 255} // Light purple
	case theme.ColorNameFocus:
		return color.NRGBA{R: 255, G: 215, B: 0, A: 255} // Gold
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 144, G: 238, B: 144, A: 255} // Light green
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 165, B: 0, A: 255} // Orange
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 100, B: 150, A: 255} // Light pink-red
	default:
		return color.NRGBA{R: 50, G: 30, B: 70, A: 255}
	}
}

func (t *FF6Theme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 8
	case theme.SizeNameSeparatorThickness:
		return 2
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInputBorder:
		return 2
	}
	return theme.DefaultTheme().Size(name)
}

func (t *FF6Theme) Font(n fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(n)
}

func (t *FF6Theme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
