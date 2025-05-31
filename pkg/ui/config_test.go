package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHomePath(t *testing.T) {
	// ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory, skipping test")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "home_directory_only",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "home_with_subdirectory",
			input:    "~/Documents/todo.txt",
			expected: filepath.Join(homeDir, "Documents/todo.txt"),
		},
		{
			name:     "absolute_path_unchanged",
			input:    "/usr/local/bin/todo.txt",
			expected: "/usr/local/bin/todo.txt",
		},
		{
			name:     "relative_path_unchanged",
			input:    "todo.txt",
			expected: "todo.txt",
		},
		{
			name:     "relative_path_with_dot",
			input:    "./todo.txt",
			expected: "./todo.txt",
		},
		{
			name:     "tilde_in_middle_unchanged",
			input:    "/path/~user/file.txt",
			expected: "/path/~user/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandHomePath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandHomePath(%s) = %s, expected %s",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultAppConfig(t *testing.T) {
	config := DefaultAppConfig()

	// テーマのデフォルト値
	if config.Theme != "catppuccin" {
		t.Errorf("Default theme = %s, expected catppuccin", config.Theme)
	}

	// 優先度レベルのデフォルト値
	expectedPriorities := []string{"", "A", "B", "C", "D"}
	if len(config.PriorityLevels) != len(expectedPriorities) {
		t.Errorf("Default priority levels length = %d, expected %d",
			len(config.PriorityLevels), len(expectedPriorities))
	}

	for i, expected := range expectedPriorities {
		if i >= len(config.PriorityLevels) || config.PriorityLevels[i] != expected {
			t.Errorf("Default priority level[%d] = %s, expected %s",
				i, config.PriorityLevels[i], expected)
		}
	}

	// デフォルトTodoファイルは空
	if config.DefaultTodoFile != "" {
		t.Errorf("Default todo file = %s, expected empty string", config.DefaultTodoFile)
	}

	// UI設定のデフォルト値
	if config.UI.LeftPaneRatio != 0.33 {
		t.Errorf("Default left pane ratio = %f, expected 0.33", config.UI.LeftPaneRatio)
	}

	if config.UI.MinLeftPaneWidth != 18 {
		t.Errorf("Default min left pane width = %d, expected 18", config.UI.MinLeftPaneWidth)
	}

	if config.UI.MinRightPaneWidth != 28 {
		t.Errorf("Default min right pane width = %d, expected 28", config.UI.MinRightPaneWidth)
	}

	if config.UI.VerticalPadding != 2 {
		t.Errorf("Default vertical padding = %d, expected 2", config.UI.VerticalPadding)
	}
}

func TestValidateAndFixConfig(t *testing.T) {
	t.Run("valid config unchanged", func(t *testing.T) {
		config := AppConfig{
			Theme:          "catppuccin",
			PriorityLevels: []string{"", "A", "B", "C", "D"},
			UI: UIConfig{
				LeftPaneRatio:     0.4,
				MinLeftPaneWidth:  20,
				MinRightPaneWidth: 30,
				VerticalPadding:   3,
			},
			Logging: LoggingConfig{
				MaxLogDays: 15,
			},
		}

		result := validateAndFixConfig(config)

		if result.Theme != "catppuccin" {
			t.Errorf("Theme = %s, expected catppuccin", result.Theme)
		}
		if result.UI.LeftPaneRatio != 0.4 {
			t.Errorf("LeftPaneRatio = %f, expected 0.4", result.UI.LeftPaneRatio)
		}
	})

	t.Run("invalid theme fixed", func(t *testing.T) {
		config := AppConfig{
			Theme: "invalid-theme",
		}

		result := validateAndFixConfig(config)

		if result.Theme != "catppuccin" {
			t.Errorf("Invalid theme should fallback to catppuccin, got %s", result.Theme)
		}
	})

	t.Run("everforest themes are valid", func(t *testing.T) {
		testCases := []string{"everforest-dark", "everforest-light"}

		for _, theme := range testCases {
			config := AppConfig{
				Theme: theme,
			}

			result := validateAndFixConfig(config)

			if result.Theme != theme {
				t.Errorf("Theme %s should be valid, but got %s", theme, result.Theme)
			}
		}
	})

	t.Run("invalid ratio fixed", func(t *testing.T) {
		config := AppConfig{
			UI: UIConfig{
				LeftPaneRatio: 1.5, // Invalid - too high
			},
		}

		result := validateAndFixConfig(config)

		if result.UI.LeftPaneRatio != 0.33 {
			t.Errorf("Invalid ratio should be fixed to 0.33, got %f", result.UI.LeftPaneRatio)
		}
	})

	t.Run("negative min width fixed", func(t *testing.T) {
		config := AppConfig{
			UI: UIConfig{
				MinLeftPaneWidth:  -5,
				MinRightPaneWidth: -10,
			},
		}

		result := validateAndFixConfig(config)

		if result.UI.MinLeftPaneWidth != 18 {
			t.Errorf("Negative MinLeftPaneWidth should be fixed to 18, got %d", result.UI.MinLeftPaneWidth)
		}
		if result.UI.MinRightPaneWidth != 28 {
			t.Errorf("Negative MinRightPaneWidth should be fixed to 28, got %d", result.UI.MinRightPaneWidth)
		}
	})
}

func TestValidateAndFixConfigPriorityLevels(t *testing.T) {
	tests := []struct {
		name            string
		inputPriorities []string
		expectedFirst   string
		expectedLength  int
		description     string
	}{
		{
			name:            "valid_priorities_unchanged",
			inputPriorities: []string{"", "A", "B", "C"},
			expectedFirst:   "",
			expectedLength:  4,
			description:     "有効な優先度リストは変更されない",
		},
		{
			name:            "empty_priorities_filled_with_defaults",
			inputPriorities: []string{},
			expectedFirst:   "",
			expectedLength:  5, // デフォルトの優先度リスト
			description:     "空の優先度リストはデフォルトで埋められる",
		},
		{
			name:            "missing_empty_priority_added",
			inputPriorities: []string{"A", "B", "C"},
			expectedFirst:   "", // 先頭に空文字列が追加される
			expectedLength:  4,
			description:     "空の優先度が先頭に追加される",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AppConfig{
				Theme:          "catppuccin",
				PriorityLevels: tt.inputPriorities,
				UI:             UIConfig{LeftPaneRatio: 0.33},
			}

			result := validateAndFixConfig(config)

			if len(result.PriorityLevels) != tt.expectedLength {
				t.Errorf("validateAndFixConfig() priority levels length = %d, expected %d for %s",
					len(result.PriorityLevels), tt.expectedLength, tt.description)
			}

			if len(result.PriorityLevels) > 0 && result.PriorityLevels[0] != tt.expectedFirst {
				t.Errorf("validateAndFixConfig() first priority = %s, expected %s for %s",
					result.PriorityLevels[0], tt.expectedFirst, tt.description)
			}
		})
	}
}

func TestLoadConfigFromFileNotFound(t *testing.T) {
	// 存在しないファイルパス
	nonExistentPath := "/path/that/does/not/exist/config.yaml"

	config, err := LoadConfigFromFile(nonExistentPath)

	// エラーが返されるべき
	if err == nil {
		t.Error("LoadConfigFromFile should return error for non-existent file")
	}

	// デフォルト設定が返されるべき
	defaultConfig := DefaultAppConfig()
	if config.Theme != defaultConfig.Theme {
		t.Errorf("LoadConfigFromFile returned config with theme %s, expected default %s",
			config.Theme, defaultConfig.Theme)
	}
}

func TestLoadConfigFromFileValidYAML(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "todotui-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 有効なYAML設定ファイルを作成
	configContent := `theme: nord
priority_levels: ["", "A", "B"]
default_todo_file: "~/todo.txt"
ui:
  left_pane_ratio: 0.4
  min_left_pane_width: 25
  min_right_pane_width: 35
  vertical_padding: 3
`

	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// 設定を読み込み
	config, err := LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	// 設定値を確認
	if config.Theme != "nord" {
		t.Errorf("Theme = %s, expected nord", config.Theme)
	}

	if config.DefaultTodoFile != "~/todo.txt" {
		t.Errorf("DefaultTodoFile = %s, expected ~/todo.txt", config.DefaultTodoFile)
	}

	if config.UI.LeftPaneRatio != 0.4 {
		t.Errorf("LeftPaneRatio = %f, expected 0.4", config.UI.LeftPaneRatio)
	}

	if config.UI.MinLeftPaneWidth != 25 {
		t.Errorf("MinLeftPaneWidth = %d, expected 25", config.UI.MinLeftPaneWidth)
	}

	expectedPriorities := []string{"", "A", "B"}
	if len(config.PriorityLevels) != len(expectedPriorities) {
		t.Errorf("PriorityLevels length = %d, expected %d",
			len(config.PriorityLevels), len(expectedPriorities))
	}
}

func TestLoadConfigFromFileInvalidYAML(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "todotui-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 無効なYAML設定ファイルを作成
	invalidContent := `theme: nord
priority_levels: [
  - invalid yaml structure
ui:
  left_pane_ratio: not_a_number
`

	configPath := filepath.Join(tmpDir, "invalid_config.yaml")
	err = os.WriteFile(configPath, []byte(invalidContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// 設定の読み込みでエラーが返されるべき
	_, err = LoadConfigFromFile(configPath)
	if err == nil {
		t.Error("LoadConfigFromFile should return error for invalid YAML")
	}
}

func TestLoadConfigWithoutPath(t *testing.T) {
	// 設定パスなしで設定を読み込み
	config := LoadConfig("")

	// デフォルト設定が返されるべき
	defaultConfig := DefaultAppConfig()
	if config.Theme != defaultConfig.Theme {
		t.Errorf("LoadConfig without path returned theme %s, expected %s",
			config.Theme, defaultConfig.Theme)
	}

	if config.UI.LeftPaneRatio != defaultConfig.UI.LeftPaneRatio {
		t.Errorf("LoadConfig without path returned ratio %f, expected %f",
			config.UI.LeftPaneRatio, defaultConfig.UI.LeftPaneRatio)
	}
}

// 統合テスト: 完全な設定ワークフロー
func TestConfigIntegration(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "todotui-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// カスタム設定を作成
	customConfig := AppConfig{
		Theme:           "nord",
		PriorityLevels:  []string{"", "HIGH", "MEDIUM", "LOW"},
		DefaultTodoFile: "~/custom-todo.txt",
		UI: UIConfig{
			LeftPaneRatio:     0.25,
			MinLeftPaneWidth:  15,
			MinRightPaneWidth: 40,
			VerticalPadding:   4,
		},
	}

	configPath := filepath.Join(tmpDir, "test_config.yaml")

	// 設定を保存
	err = SaveConfigToFile(customConfig, configPath)
	if err != nil {
		t.Fatalf("SaveConfigToFile failed: %v", err)
	}

	// 保存された設定を読み込み
	loadedConfig, err := LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	// バリデーションと修正を適用
	validatedConfig := validateAndFixConfig(loadedConfig)

	// 値が正しく保存・読み込みされたことを確認
	if validatedConfig.Theme != customConfig.Theme {
		t.Errorf("Integrated config theme = %s, expected %s",
			validatedConfig.Theme, customConfig.Theme)
	}

	// パス展開が行われるため、展開された結果と比較
	expectedTodoFile := ExpandHomePath(customConfig.DefaultTodoFile)
	if validatedConfig.DefaultTodoFile != expectedTodoFile {
		t.Errorf("Integrated config todo file = %s, expected %s",
			validatedConfig.DefaultTodoFile, expectedTodoFile)
	}

	if len(validatedConfig.PriorityLevels) != len(customConfig.PriorityLevels) {
		t.Errorf("Integrated config priority levels length = %d, expected %d",
			len(validatedConfig.PriorityLevels), len(customConfig.PriorityLevels))
	}

	// UI設定の確認
	if validatedConfig.UI.LeftPaneRatio != customConfig.UI.LeftPaneRatio {
		t.Errorf("Integrated config left pane ratio = %f, expected %f",
			validatedConfig.UI.LeftPaneRatio, customConfig.UI.LeftPaneRatio)
	}
}
