package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme
type Theme struct {
	// Priority colors
	PriorityHigh    lipgloss.Color // A
	PriorityMedium  lipgloss.Color // B
	PriorityLow     lipgloss.Color // C
	PriorityLowest  lipgloss.Color // D
	PriorityDefault lipgloss.Color // Other priorities

	// UI element colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Danger    lipgloss.Color

	// Text colors
	Text       lipgloss.Color
	TextMuted  lipgloss.Color
	TextSubtle lipgloss.Color

	// Background colors
	Background   lipgloss.Color
	Surface      lipgloss.Color
	SurfaceLight lipgloss.Color

	// Border colors
	BorderActive   lipgloss.Color
	BorderInactive lipgloss.Color

	// Selection colors
	SelectionBg lipgloss.Color
	SelectionFg lipgloss.Color
}

// Available themes
var themes = map[string]Theme{
	"catppuccin": {
		PriorityHigh:    lipgloss.Color("#f38ba8"), // Red
		PriorityMedium:  lipgloss.Color("#f9e2af"), // Yellow
		PriorityLow:     lipgloss.Color("#89b4fa"), // Blue
		PriorityLowest:  lipgloss.Color("#fab387"), // Peach
		PriorityDefault: lipgloss.Color("#fab387"), // Peach

		Primary:   lipgloss.Color("#89b4fa"), // Blue
		Secondary: lipgloss.Color("#74c7ec"), // Sapphire
		Success:   lipgloss.Color("#a6e3a1"), // Green
		Warning:   lipgloss.Color("#f9e2af"), // Yellow
		Danger:    lipgloss.Color("#f38ba8"), // Red

		Text:       lipgloss.Color("#cdd6f4"), // Text
		TextMuted:  lipgloss.Color("#bac2de"), // Subtext1
		TextSubtle: lipgloss.Color("#a6adc8"), // Subtext0

		Background:   lipgloss.Color("#1e1e2e"), // Base
		Surface:      lipgloss.Color("#313244"), // Surface0
		SurfaceLight: lipgloss.Color("#45475a"), // Surface1

		BorderActive:   lipgloss.Color("#89b4fa"), // Blue
		BorderInactive: lipgloss.Color("#6c7086"), // Overlay0

		SelectionBg: lipgloss.Color("#45475a"), // Surface1
		SelectionFg: lipgloss.Color("#cdd6f4"), // Text
	},
	"nord": {
		PriorityHigh:    lipgloss.Color("#BF616A"), // Nord11 - Red
		PriorityMedium:  lipgloss.Color("#EBCB8B"), // Nord13 - Yellow
		PriorityLow:     lipgloss.Color("#5E81AC"), // Nord10 - Blue
		PriorityLowest:  lipgloss.Color("#D08770"), // Nord12 - Orange
		PriorityDefault: lipgloss.Color("#D08770"), // Nord12 - Orange

		Primary:   lipgloss.Color("#5E81AC"), // Nord10 - Blue
		Secondary: lipgloss.Color("#88C0D0"), // Nord8 - Cyan
		Success:   lipgloss.Color("#A3BE8C"), // Nord14 - Green
		Warning:   lipgloss.Color("#EBCB8B"), // Nord13 - Yellow
		Danger:    lipgloss.Color("#BF616A"), // Nord11 - Red

		Text:       lipgloss.Color("#ECEFF4"), // Nord6 - White
		TextMuted:  lipgloss.Color("#E5E9F0"), // Nord5 - Off-white
		TextSubtle: lipgloss.Color("#D8DEE9"), // Nord4 - Light gray

		Background:   lipgloss.Color("#2E3440"), // Nord0 - Dark
		Surface:      lipgloss.Color("#3B4252"), // Nord1 - Dark gray
		SurfaceLight: lipgloss.Color("#434C5E"), // Nord2 - Medium gray

		BorderActive:   lipgloss.Color("#5E81AC"), // Nord10 - Blue
		BorderInactive: lipgloss.Color("#4C566A"), // Nord3 - Gray

		SelectionBg: lipgloss.Color("#434C5E"), // Nord2
		SelectionFg: lipgloss.Color("#ECEFF4"), // Nord6
	},
	"everforest-dark": {
		PriorityHigh:    lipgloss.Color("#e67e80"), // Red
		PriorityMedium:  lipgloss.Color("#dbbc7f"), // Yellow
		PriorityLow:     lipgloss.Color("#7fbbb3"), // Aqua
		PriorityLowest:  lipgloss.Color("#d699b6"), // Purple
		PriorityDefault: lipgloss.Color("#a7c080"), // Green

		Primary:   lipgloss.Color("#a7c080"), // Green
		Secondary: lipgloss.Color("#7fbbb3"), // Aqua
		Success:   lipgloss.Color("#a7c080"), // Green
		Warning:   lipgloss.Color("#dbbc7f"), // Yellow
		Danger:    lipgloss.Color("#e67e80"), // Red

		Text:       lipgloss.Color("#d3c6aa"), // Foreground
		TextMuted:  lipgloss.Color("#9da9a0"), // Grey1
		TextSubtle: lipgloss.Color("#859289"), // Grey0

		Background:   lipgloss.Color("#2d353b"), // Background
		Surface:      lipgloss.Color("#3d484d"), // Background Dark
		SurfaceLight: lipgloss.Color("#475258"), // Background Light

		BorderActive:   lipgloss.Color("#a7c080"), // Green
		BorderInactive: lipgloss.Color("#543a48"), // StatusLine

		SelectionBg: lipgloss.Color("#475258"), // Background Light
		SelectionFg: lipgloss.Color("#d3c6aa"), // Foreground
	},
	"everforest-light": {
		PriorityHigh:    lipgloss.Color("#f85552"), // Red
		PriorityMedium:  lipgloss.Color("#dfa000"), // Yellow
		PriorityLow:     lipgloss.Color("#35a77c"), // Aqua
		PriorityLowest:  lipgloss.Color("#df69ba"), // Purple
		PriorityDefault: lipgloss.Color("#8da101"), // Green

		Primary:   lipgloss.Color("#8da101"), // Green
		Secondary: lipgloss.Color("#35a77c"), // Aqua
		Success:   lipgloss.Color("#8da101"), // Green
		Warning:   lipgloss.Color("#dfa000"), // Yellow
		Danger:    lipgloss.Color("#f85552"), // Red

		Text:       lipgloss.Color("#5c6a72"), // Foreground
		TextMuted:  lipgloss.Color("#829181"), // Grey1
		TextSubtle: lipgloss.Color("#a6b0a0"), // Grey0

		Background:   lipgloss.Color("#fdf6e3"), // Background
		Surface:      lipgloss.Color("#f4f0d9"), // Background Dark
		SurfaceLight: lipgloss.Color("#efebd4"), // Background Light

		BorderActive:   lipgloss.Color("#8da101"), // Green
		BorderInactive: lipgloss.Color("#f0f2d4"), // StatusLine

		SelectionBg: lipgloss.Color("#efebd4"), // Background Light
		SelectionFg: lipgloss.Color("#5c6a72"), // Foreground
	},
}

// GetTheme returns the theme based on theme name
func GetTheme(themeName string) Theme {
	if theme, exists := themes[themeName]; exists {
		return theme
	}

	// Fallback to catppuccin if theme not found
	return themes["catppuccin"]
}
