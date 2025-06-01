package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Config はログの設定を定義
type Config struct {
	Level          slog.Level
	OutputToStderr bool   // stderrへの出力を制御
	AppName        string // アプリケーション名（ログディレクトリの決定に使用）
}

// Logger はアプリケーション全体で使用するロガー
var globalLogger *slog.Logger

// CLIアプリケーション用のログディレクトリを取得
func getLogDirectory(appName string) (string, error) {
	var logDir string

	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Logs/{appName}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		logDir = filepath.Join(homeDir, "Library", "Logs", appName)
	case "linux":
		// Linux: ~/.local/share/{appName}/logs
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		logDir = filepath.Join(homeDir, ".local", "share", appName, "logs")
	case "windows":
		// Windows: %APPDATA%/{appName}/logs
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		logDir = filepath.Join(appData, appName, "logs")
	default:
		// その他: ~/.{appName}/logs
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		logDir = filepath.Join(homeDir, "."+appName, "logs")
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("ログディレクトリの作成に失敗: %w", err)
	}

	return logDir, nil
}

// Init はログシステムを初期化
func Init(config Config) error {
	var writers []io.Writer

	// OutputToStderrが有効な場合のみstderrに出力
	if config.OutputToStderr {
		writers = append(writers, os.Stderr)
	}

	// 出力先が設定されていない場合はエラー
	if len(writers) == 0 {
		return fmt.Errorf("ログの出力先が設定されていません（OutputToStderr を有効にしてください）")
	}

	// マルチライターを作成
	multiWriter := io.MultiWriter(writers...)

	// ログレベルを設定
	level := slog.LevelInfo
	if config.Level != 0 {
		level = config.Level
	}

	// CLI向けの人に見やすいハンドラーを作成
	handlerOpts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 時刻フォーマットをより読みやすく
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("15:04:05"))
			}
			// ソースの情報を短くする（ファイル名:行番号のみ）
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line))
			}
			return a
		},
		AddSource: level == slog.LevelDebug, // デバッグ時のみソース情報を表示
	}

	handler := slog.NewTextHandler(multiWriter, handlerOpts)
	globalLogger = slog.New(handler)

	// slogのデフォルトロガーも設定
	slog.SetDefault(globalLogger)

	return nil
}

// GetLogger はグローバルロガーを取得
func GetLogger() *slog.Logger {
	return globalLogger
}

// レベル別のヘルパー関数
func Debug(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.Error(msg, args...)
	}
}

// CleanupOldLogs は古いログファイルを削除する
func CleanupOldLogs(appName string, maxDays int) error {
	logDir, err := getLogDirectory(appName)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -maxDays)

	return filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ログファイルかチェック
		if !info.IsDir() && filepath.Ext(path) == ".log" {
			if info.ModTime().Before(cutoff) {
				Debug("古いログファイルを削除", "file", path)
				return os.Remove(path)
			}
		}

		return nil
	})
}
