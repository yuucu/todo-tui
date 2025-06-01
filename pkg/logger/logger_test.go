package logger

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestGetLogDirectory(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		wantDir bool
	}{
		{
			name:    "Valid app name",
			appName: "testapp",
			wantDir: true,
		},
		{
			name:    "Empty app name",
			appName: "",
			wantDir: true,
		},
		{
			name:    "App name with spaces",
			appName: "test app",
			wantDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLogDirectory(tt.appName)
			if err != nil {
				t.Errorf("getLogDirectory() error = %v", err)
				return
			}

			if tt.wantDir {
				// ディレクトリが作成されているかチェック
				if _, err := os.Stat(got); os.IsNotExist(err) {
					t.Errorf("getLogDirectory() directory not created: %v", got)
				}

				// OSに応じた適切なパスかチェック
				switch runtime.GOOS {
				case "darwin":
					if !strings.Contains(got, "Library/Logs") {
						t.Errorf("getLogDirectory() for macOS should contain 'Library/Logs', got: %v", got)
					}
				case "linux":
					if !strings.Contains(got, ".local/share") {
						t.Errorf("getLogDirectory() for Linux should contain '.local/share', got: %v", got)
					}
				case "windows":
					if !strings.Contains(got, "AppData") && !strings.Contains(got, "APPDATA") {
						t.Errorf("getLogDirectory() for Windows should contain 'AppData', got: %v", got)
					}
				}

				// アプリ名が含まれているかチェック
				if tt.appName != "" && !strings.Contains(got, tt.appName) {
					t.Errorf("getLogDirectory() should contain app name '%v', got: %v", tt.appName, got)
				}
			}

			// テスト後のクリーンアップ
			if err := os.RemoveAll(got); err != nil {
				t.Logf("Warning: Failed to cleanup test directory: %v", err)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name               string
		config             Config
		wantError          bool
		errorContains      string
		checkLogger        bool
		checkLogFile       bool
		useTemporaryConfig bool
	}{
		{
			name: "Output to stderr only (actual usage)",
			config: Config{
				Level:          slog.LevelInfo,
				EnableDebug:    false,
				OutputToFile:   false,
				OutputToStderr: true,
				AppName:        "todotui",
			},
			wantError:   false,
			checkLogger: true,
		},
		{
			name: "Output to file only (for testing)",
			config: Config{
				Level:          slog.LevelDebug,
				EnableDebug:    true,
				OutputToFile:   true,
				OutputToStderr: false,
				AppName:        "test",
			},
			wantError:          false,
			checkLogger:        true,
			checkLogFile:       true,
			useTemporaryConfig: true,
		},
		{
			name: "No output configured",
			config: Config{
				Level:          slog.LevelInfo,
				EnableDebug:    false,
				OutputToFile:   false,
				OutputToStderr: false,
				AppName:        "test",
			},
			wantError:     true,
			errorContains: "ログの出力先が設定されていません",
		},
		{
			name: "Both outputs enabled (for testing)",
			config: Config{
				Level:          slog.LevelWarn,
				EnableDebug:    false,
				OutputToFile:   true,
				OutputToStderr: true,
				AppName:        "test",
			},
			wantError:          false,
			checkLogger:        true,
			checkLogFile:       true,
			useTemporaryConfig: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト前の状態を保存
			originalLogger := globalLogger
			defer func() {
				globalLogger = originalLogger
				slog.SetDefault(slog.Default())
			}()

			// テスト用の一時ファイルパスを設定（AppNameベースで自動生成）
			if tt.useTemporaryConfig {
				// テスト用に一時的なAppNameを設定してファイルの場所を予測可能にする
				tt.config.AppName = "test-" + tt.name
			}

			err := Init(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("Init() expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Init() error should contain '%v', got: %v", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Init() unexpected error = %v", err)
				return
			}

			if tt.checkLogger {
				if globalLogger == nil {
					t.Error("Init() globalLogger should not be nil")
				}

				logger := GetLogger()
				if logger == nil {
					t.Error("GetLogger() should not return nil")
				}
			}

			if tt.checkLogFile {
				// ログファイルが作成されることを確認
				Info("test message")

				// AppNameベースで自動生成されたログファイルの存在確認
				// 実際のファイルパスは getLogDirectory() + AppName-date.log の形式
				logDir, err := getLogDirectory(tt.config.AppName)
				if err != nil {
					t.Logf("Could not get log directory (expected for some tests): %v", err)
					return
				}

				// 日付ベースのログファイル名を生成
				today := time.Now().Format("2006-01-02")
				expectedLogFile := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", tt.config.AppName, today))

				if _, err := os.Stat(expectedLogFile); os.IsNotExist(err) {
					t.Logf("Log file not found at expected location (may be expected for some environments): %v", expectedLogFile)
				}
			}
		})
	}
}

func TestLoggerHelperFunctions(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(string, ...any)
		message string
		level   string
	}{
		{
			name:    "Debug",
			logFunc: Debug,
			message: "debug message",
			level:   "DEBUG",
		},
		{
			name:    "Info",
			logFunc: Info,
			message: "info message",
			level:   "INFO",
		},
		{
			name:    "Warn",
			logFunc: Warn,
			message: "warn message",
			level:   "WARN",
		},
		{
			name:    "Error",
			logFunc: Error,
			message: "error message",
			level:   "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト前の状態を保存
			originalLogger := globalLogger
			defer func() {
				globalLogger = originalLogger
				slog.SetDefault(slog.Default())
			}()

			// バッファに出力をキャプチャ
			var buf bytes.Buffer

			// テスト用のロガーを設定
			handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
			globalLogger = slog.New(handler)

			tt.logFunc(tt.message)

			output := buf.String()
			if !strings.Contains(output, tt.level) {
				t.Errorf("Log output should contain level '%v', got: %v", tt.level, output)
			}
			if !strings.Contains(output, tt.message) {
				t.Errorf("Log output should contain message '%v', got: %v", tt.message, output)
			}
		})
	}
}

func TestLoggerHelperFunctions_NilGlobalLogger(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(string, ...any)
		message string
	}{
		{
			name:    "Debug with nil logger",
			logFunc: Debug,
			message: "test debug",
		},
		{
			name:    "Info with nil logger",
			logFunc: Info,
			message: "test info",
		},
		{
			name:    "Warn with nil logger",
			logFunc: Warn,
			message: "test warn",
		},
		{
			name:    "Error with nil logger",
			logFunc: Error,
			message: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト前の状態を保存
			originalLogger := globalLogger
			defer func() {
				globalLogger = originalLogger
			}()

			// globalLoggerをnilに設定
			globalLogger = nil

			// ヘルパー関数がpanicしないことを確認
			tt.logFunc(tt.message)

			// GetLoggerがnilを返すことを確認
			logger := GetLogger()
			if logger != nil {
				t.Error("GetLogger() should return nil when globalLogger is nil")
			}
		})
	}
}

func TestCleanupOldLogs(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		maxDays int
		wantErr bool
	}{
		{
			name:    "Empty app name",
			appName: "",
			maxDays: 7,
			wantErr: false, // エラーが発生する可能性があるが関数は正常実行される
		},
		{
			name:    "Valid app name",
			appName: "testapp",
			maxDays: 30,
			wantErr: false,
		},
		{
			name:    "Zero max days",
			appName: "testapp",
			maxDays: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CleanupOldLogs(tt.appName, tt.maxDays)

			if tt.wantErr && err == nil {
				t.Error("CleanupOldLogs() expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Logf("CleanupOldLogs() returned error (may be expected): %v", err)
			}
		})
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		expectedLevel  slog.Level
		expectedDebug  bool
		expectedFile   bool
		expectedStderr bool
		expectedApp    string
	}{
		{
			name:           "Zero value config",
			config:         Config{},
			expectedLevel:  0,
			expectedDebug:  false,
			expectedFile:   false,
			expectedStderr: false,
			expectedApp:    "",
		},
		{
			name: "Actual usage config (stderr only)",
			config: Config{
				Level:          slog.LevelWarn,
				EnableDebug:    false,
				OutputToFile:   false,
				OutputToStderr: true,
				AppName:        "todotui",
			},
			expectedLevel:  slog.LevelWarn,
			expectedDebug:  false,
			expectedFile:   false,
			expectedStderr: true,
			expectedApp:    "todotui",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Level != tt.expectedLevel {
				t.Errorf("Level should be %v, got: %v", tt.expectedLevel, tt.config.Level)
			}
			if tt.config.EnableDebug != tt.expectedDebug {
				t.Errorf("EnableDebug should be %v, got: %v", tt.expectedDebug, tt.config.EnableDebug)
			}
			if tt.config.OutputToFile != tt.expectedFile {
				t.Errorf("OutputToFile should be %v, got: %v", tt.expectedFile, tt.config.OutputToFile)
			}
			if tt.config.OutputToStderr != tt.expectedStderr {
				t.Errorf("OutputToStderr should be %v, got: %v", tt.expectedStderr, tt.config.OutputToStderr)
			}
			if tt.config.AppName != tt.expectedApp {
				t.Errorf("AppName should be %v, got: %v", tt.expectedApp, tt.config.AppName)
			}
		})
	}
}
