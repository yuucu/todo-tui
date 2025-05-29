package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppConfig represents the complete application configuration
type AppConfig struct {
	// Theme settings
	Theme string `json:"theme,omitempty"`

	// Priority levels configuration
	PriorityLevels []string `json:"priority_levels,omitempty"`

	// Default todo file path
	DefaultTodoFile string `json:"default_todo_file,omitempty"`

	// UI settings
	UI UIConfig `json:"ui,omitempty"`
}

// UIConfig represents UI-specific configuration
type UIConfig struct {
	// Pane width ratio (left pane width / total width)
	LeftPaneRatio float64 `json:"left_pane_ratio,omitempty"`

	// Minimum pane widths
	MinLeftPaneWidth  int `json:"min_left_pane_width,omitempty"`
	MinRightPaneWidth int `json:"min_right_pane_width,omitempty"`

	// Vertical padding/spacing settings
	VerticalPadding int `json:"vertical_padding,omitempty"`
}

// DefaultAppConfig returns the default application configuration
func DefaultAppConfig() AppConfig {
	homeDir, _ := os.UserHomeDir()
	defaultTodoFile := filepath.Join(homeDir, "todo.txt")

	return AppConfig{
		Theme:           "catppuccin",
		PriorityLevels:  []string{"", "A", "B", "C", "D"},
		DefaultTodoFile: defaultTodoFile,
		UI: UIConfig{
			LeftPaneRatio:     0.33,
			MinLeftPaneWidth:  18,
			MinRightPaneWidth: 28,
			VerticalPadding:   2, // Default vertical padding for spacing optimization
		},
	}
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(configPath string) (AppConfig, error) {
	config := DefaultAppConfig()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, fmt.Errorf("config file not found: %s", configPath)
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and set defaults for missing fields
	config = validateAndFixConfig(config)

	return config, nil
}

// LoadConfig loads configuration from file if specified, otherwise from environment variables
func LoadConfig(configPath string) AppConfig {
	var config AppConfig

	if configPath != "" {
		var err error
		config, err = LoadConfigFromFile(configPath)
		if err != nil {
			// If config file loading fails, print warning and use defaults
			fmt.Fprintf(os.Stderr, "Warning: %v\nUsing default configuration.\n", err)
			config = DefaultAppConfig()
		}
	} else {
		config = DefaultAppConfig()
	}

	// Override with environment variables if set
	config = applyEnvironmentOverrides(config)

	return config
}

// applyEnvironmentOverrides applies environment variable overrides
func applyEnvironmentOverrides(config AppConfig) AppConfig {
	// Theme override
	if theme := os.Getenv("TODO_TUI_THEME"); theme != "" {
		config.Theme = theme
	}

	// Priority levels override
	if envPriorities := os.Getenv("TODO_TUI_PRIORITY_LEVELS"); envPriorities != "" {
		levels := strings.Split(envPriorities, ",")
		for i, level := range levels {
			levels[i] = strings.TrimSpace(level)
		}
		// Ensure empty string is first (no priority)
		if len(levels) > 0 && levels[0] != "" {
			levels = append([]string{""}, levels...)
		}
		config.PriorityLevels = levels
	}

	return config
}

// validateAndFixConfig validates the configuration and sets defaults for invalid values
func validateAndFixConfig(config AppConfig) AppConfig {
	// Validate theme
	validThemes := GetAvailableThemes()
	isValidTheme := false
	for _, theme := range validThemes {
		if config.Theme == theme {
			isValidTheme = true
			break
		}
	}
	if !isValidTheme {
		config.Theme = "catppuccin"
	}

	// Validate priority levels
	if len(config.PriorityLevels) == 0 {
		config.PriorityLevels = []string{"", "A", "B", "C", "D"}
	}

	// Ensure empty string is first (no priority)
	if config.PriorityLevels[0] != "" {
		config.PriorityLevels = append([]string{""}, config.PriorityLevels...)
	}

	// Validate UI settings
	if config.UI.LeftPaneRatio <= 0 || config.UI.LeftPaneRatio >= 1 {
		config.UI.LeftPaneRatio = 0.33
	}

	if config.UI.MinLeftPaneWidth <= 0 {
		config.UI.MinLeftPaneWidth = 18
	}

	if config.UI.MinRightPaneWidth <= 0 {
		config.UI.MinRightPaneWidth = 28
	}

	// Validate vertical padding
	if config.UI.VerticalPadding < 1 {
		config.UI.VerticalPadding = 2
	}

	// Set default todo file if not specified
	if config.DefaultTodoFile == "" {
		homeDir, _ := os.UserHomeDir()
		config.DefaultTodoFile = filepath.Join(homeDir, "todo.txt")
	}

	return config
}

// SaveConfigToFile saves configuration to a JSON file
func SaveConfigToFile(config AppConfig, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// CreateSampleConfig creates a sample configuration file with comments
func CreateSampleConfig(configPath string) error {
	config := DefaultAppConfig()

	// Create the configuration file with sample values
	if err := SaveConfigToFile(config, configPath); err != nil {
		return err
	}

	fmt.Printf("Sample configuration file created at: %s\n", configPath)
	fmt.Println("You can edit this file to customize your settings.")

	return nil
}
