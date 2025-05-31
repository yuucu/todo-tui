package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/yuucu/todotui/pkg/logger"
)

// 設定ファイル用のディレクトリパーミッション
const defaultConfigDirMode = 0755

// ExpandHomePath expands ~ in file paths to the user's home directory
func ExpandHomePath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path // Return original path if we can't get home directory
	}

	if path == "~" {
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path // Return original path if it doesn't match expected patterns
}

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

	// Logging settings
	Logging LoggingConfig `mapstructure:"logging"`
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

	// Checkbox style for task display
	CheckboxStyle string `mapstructure:"checkbox_style"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	// Enable debug logging
	EnableDebug bool `mapstructure:"enable_debug"`

	// Custom log file path (optional, ファイル出力は常に有効)
	LogFilePath string `mapstructure:"log_file_path"`

	// Maximum days to keep log files
	MaxLogDays int `mapstructure:"max_log_days"`
}

// DefaultAppConfig returns the default application configuration
func DefaultAppConfig() AppConfig {
	return AppConfig{
		Theme:           "catppuccin",
		PriorityLevels:  []string{"", "A", "B", "C", "D"},
		DefaultTodoFile: "", // No default file, must be specified explicitly
		UI: UIConfig{
			LeftPaneRatio:     0.33,
			MinLeftPaneWidth:  18,
			MinRightPaneWidth: 28,
			VerticalPadding:   2,
			CheckboxStyle:     DefaultCheckboxStyle,
		},
		Logging: LoggingConfig{
			EnableDebug: false,
			LogFilePath: "", // 空の場合はデフォルトパスを使用
			MaxLogDays:  30, // 30日間ログを保持
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

	// Note: validateAndFixConfig will be called by LoadConfig
	return config, nil
}

// LoadConfig loads configuration from file if specified, otherwise uses default configuration
func LoadConfig(configPath string) AppConfig {
	var config AppConfig

	if configPath != "" {
		// Use specified config file path
		var err error
		config, err = LoadConfigFromFile(configPath)
		if err != nil {
			// If config file loading fails, print warning and use defaults
			// ログシステムが初期化されているかチェックして両方に出力
			if logger.GetLogger() != nil {
				logger.Warn("Config file loading failed, using defaults", "error", err)
			}
			fmt.Fprintf(os.Stderr, "Warning: %v\nUsing default configuration.\n", err)
			config = DefaultAppConfig()
		} else {
			// Validate and fix config, which includes path expansion
			config = validateAndFixConfig(config)
		}
	} else {
		// Try to load from default locations if no config path is specified
		configPath = findDefaultConfigFile()
		if configPath != "" {
			var err error
			config, err = LoadConfigFromFile(configPath)
			if err != nil {
				// If config file loading fails, print warning and use defaults
				// ログシステムが初期化されているかチェックして両方に出力
				if logger.GetLogger() != nil {
					logger.Warn("Config file loading failed from default location, using defaults", "path", configPath, "error", err)
				}
				fmt.Fprintf(os.Stderr, "Warning: Failed to load config from %s: %v\nUsing default configuration.\n", configPath, err)
				config = DefaultAppConfig()
			} else {
				// Validate and fix config, which includes path expansion
				config = validateAndFixConfig(config)
			}
		} else {
			config = DefaultAppConfig()
		}
	}

	return config
}

// findDefaultConfigFile searches for config files in common locations
// Returns the path to the first config file found, or empty string if none found
func findDefaultConfigFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Check for config.yaml in ~/.config/todotui/ directory
	configPath := filepath.Join(homeDir, ".config", "todotui", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	return ""
}

// validateAndFixConfig validates the configuration and sets defaults for invalid values
func validateAndFixConfig(config AppConfig) AppConfig {
	// Validate theme
	validThemes := []string{"catppuccin", "nord", "everforest-dark", "everforest-light"}
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

	// Validate checkbox style
	validCheckboxStyles := []string{CheckboxStyleCircle, CheckboxStyleSquare, CheckboxStyleCheck, CheckboxStyleDiamond, CheckboxStyleStar}
	isValidCheckboxStyle := false
	for _, style := range validCheckboxStyles {
		if config.UI.CheckboxStyle == style {
			isValidCheckboxStyle = true
			break
		}
	}
	if !isValidCheckboxStyle {
		config.UI.CheckboxStyle = DefaultCheckboxStyle
	}

	// Validate logging settings
	if config.Logging.MaxLogDays <= 0 {
		config.Logging.MaxLogDays = 30
	}

	// Expand ~ in path if present (only if path is specified)
	if config.DefaultTodoFile != "" {
		config.DefaultTodoFile = ExpandHomePath(config.DefaultTodoFile)
	}

	// Expand ~ in log file path if present
	if config.Logging.LogFilePath != "" {
		config.Logging.LogFilePath = ExpandHomePath(config.Logging.LogFilePath)
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
	v.Set("ui.checkbox_style", config.UI.CheckboxStyle)

	// Set logging configuration
	v.Set("logging.enable_debug", config.Logging.EnableDebug)
	v.Set("logging.log_file_path", config.Logging.LogFilePath)
	v.Set("logging.max_log_days", config.Logging.MaxLogDays)

	// Set config file path (Viper will determine format by extension)
	v.SetConfigFile(configPath)

	// Write the configuration file
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
