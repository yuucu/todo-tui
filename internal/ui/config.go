package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// 設定ファイル用のディレクトリパーミッション
const defaultConfigDirMode = 0755

// AppConfig represents the complete application configuration
type AppConfig struct {
	// Theme settings
	Theme string `mapstructure:"theme"`

	// Priority levels configuration
	PriorityLevels []string `mapstructure:"priority_levels"`

	// Default todo file path
	DefaultTodoFile string `mapstructure:"default_todo_file"`

	// UI settings
	UI UIConfig `mapstructure:"ui"`
}

// UIConfig represents UI-specific configuration
type UIConfig struct {
	// Pane width ratio (left pane width / total width)
	LeftPaneRatio float64 `mapstructure:"left_pane_ratio"`

	// Minimum pane widths
	MinLeftPaneWidth  int `mapstructure:"min_left_pane_width"`
	MinRightPaneWidth int `mapstructure:"min_right_pane_width"`

	// Vertical padding/spacing settings
	VerticalPadding int `mapstructure:"vertical_padding"`
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
			VerticalPadding:   2,
		},
	}
}

// LoadConfigFromFile loads configuration from a file using Viper
// Supports YAML, JSON, TOML, HCL, envfile and Java properties config files
func LoadConfigFromFile(configPath string) (AppConfig, error) {
	config := DefaultAppConfig()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, fmt.Errorf("config file not found: %s", configPath)
	}

	// Setup Viper
	v := viper.New()

	// Set config file path
	v.SetConfigFile(configPath)

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into our struct
	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
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

// SaveConfigToFile saves configuration to a file using Viper
// Format is determined by file extension (yaml, json, toml, etc.)
func SaveConfigToFile(config AppConfig, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, defaultConfigDirMode); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Setup Viper
	v := viper.New()

	// Set the config data
	v.Set("theme", config.Theme)
	v.Set("priority_levels", config.PriorityLevels)
	v.Set("default_todo_file", config.DefaultTodoFile)
	v.Set("ui.left_pane_ratio", config.UI.LeftPaneRatio)
	v.Set("ui.min_left_pane_width", config.UI.MinLeftPaneWidth)
	v.Set("ui.min_right_pane_width", config.UI.MinRightPaneWidth)
	v.Set("ui.vertical_padding", config.UI.VerticalPadding)

	// Set config file path (Viper will determine format by extension)
	v.SetConfigFile(configPath)

	// Write the configuration file
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// CreateSampleConfig creates a sample configuration file
func CreateSampleConfig(configPath string) error {
	config := DefaultAppConfig()

	if err := SaveConfigToFile(config, configPath); err != nil {
		return fmt.Errorf("failed to create sample config: %w", err)
	}

	fmt.Printf("Sample configuration created at: %s\n", configPath)
	return nil
}
