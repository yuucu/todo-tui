package ui

import "testing"

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name      string
		themeName string
		expected  bool // whether theme should exist
	}{
		{
			name:      "catppuccin theme exists",
			themeName: "catppuccin",
			expected:  true,
		},
		{
			name:      "nord theme exists",
			themeName: "nord",
			expected:  true,
		},
		{
			name:      "everforest-dark theme exists",
			themeName: "everforest-dark",
			expected:  true,
		},
		{
			name:      "everforest-light theme exists",
			themeName: "everforest-light",
			expected:  true,
		},
		{
			name:      "nonexistent theme fallbacks to catppuccin",
			themeName: "nonexistent",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := GetTheme(tt.themeName)

			// Check that theme is not empty
			if theme.Text == "" {
				t.Errorf("GetTheme(%s) returned empty theme", tt.themeName)
			}

			// For nonexistent themes, should return catppuccin theme
			if !tt.expected {
				expectedTheme := GetTheme("catppuccin")
				if theme.Text != expectedTheme.Text {
					t.Errorf("GetTheme(%s) should fallback to catppuccin theme", tt.themeName)
				}
			}
		})
	}
}

func TestEverforestThemeColors(t *testing.T) {
	darkTheme := GetTheme("everforest-dark")
	lightTheme := GetTheme("everforest-light")

	// Test that Everforest themes have different colors
	if darkTheme.Background == lightTheme.Background {
		t.Error("Everforest dark and light themes should have different background colors")
	}

	if darkTheme.Text == lightTheme.Text {
		t.Error("Everforest dark and light themes should have different text colors")
	}

	// Test that all required colors are defined (not empty)
	testThemeCompleteness := func(t *testing.T, theme Theme, themeName string) {
		if theme.PriorityHigh == "" {
			t.Errorf("%s theme missing PriorityHigh color", themeName)
		}
		if theme.Primary == "" {
			t.Errorf("%s theme missing Primary color", themeName)
		}
		if theme.Text == "" {
			t.Errorf("%s theme missing Text color", themeName)
		}
		if theme.Background == "" {
			t.Errorf("%s theme missing Background color", themeName)
		}
	}

	t.Run("everforest-dark completeness", func(t *testing.T) {
		testThemeCompleteness(t, darkTheme, "everforest-dark")
	})

	t.Run("everforest-light completeness", func(t *testing.T) {
		testThemeCompleteness(t, lightTheme, "everforest-light")
	})
}
