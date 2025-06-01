package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuucu/todotui/pkg/logger"
	"github.com/yuucu/todotui/pkg/ui"
)

// Build information variables (set by goreleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelWarn // デフォルトは警告レベル
	}
}

func printVersion() {
	fmt.Printf("todotui %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built: %s\n", date)
}

func printUsage() {
	fmt.Printf(`Usage: %s [OPTIONS] [TODO_FILE]

A terminal todo.txt manager with vim-like keybindings.

Arguments:
  TODO_FILE    Path to todo.txt file (required unless set in config)

Options:
  -c, --config CONFIG       Path to configuration file
  -t, --theme THEME         Set color theme (catppuccin, nord, everforest-dark, everforest-light)
  -d, --debug               Enable debug level logging (default: warning level)
  -v, --version             Show version information
  -h, --help               Show this help message

For detailed documentation and keybindings, see: https://github.com/yuucu/todotui
`, os.Args[0])
}

func main() {
	// Setup IME environment for Japanese input support
	ui.SetupIMEEnvironment()

	// Define command line flags
	var (
		configFile  = flag.String("config", "", "Path to configuration file")
		themeName   = flag.String("theme", "", "Set color theme (catppuccin, nord, everforest-dark, everforest-light)")
		enableDebug = flag.Bool("debug", false, "Enable debug logging")
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show this help message")
	)

	// Define short flag aliases (reuse same help text)
	flag.StringVar(configFile, "c", "", "Path to configuration file")
	flag.StringVar(themeName, "t", "", "Set color theme (catppuccin, nord, everforest-dark, everforest-light)")
	flag.BoolVar(enableDebug, "d", false, "Enable debug logging")
	flag.BoolVar(showVersion, "v", false, "Show version information")
	flag.BoolVar(showHelp, "h", false, "Show this help message")

	// Set custom usage function
	flag.Usage = printUsage

	// Parse command line flags
	flag.Parse()

	// Handle help flag first (convention)
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Handle version flag
	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// Get remaining non-flag arguments (todo file)
	args := flag.Args()
	var todoFile string
	if len(args) > 0 {
		todoFile = args[0]
	}
	if len(args) > 1 {
		// ログ初期化前なので、stderrに直接出力
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n")
		flag.Usage()
		os.Exit(1)
	}

	// Load configuration
	appConfig := ui.LoadConfig(*configFile)

	// Override theme if specified via command line
	if *themeName != "" {
		appConfig.Theme = *themeName
	}

	// Override debug logging if specified via command line
	if *enableDebug {
		appConfig.Logging.EnableDebug = true
	}

	// Initialize logging system (first initialization without UI channel)
	var finalLogLevel string
	if *enableDebug {
		finalLogLevel = "DEBUG"
	} else if appConfig.Logging.LogLevel != "" {
		finalLogLevel = appConfig.Logging.LogLevel
	} else {
		finalLogLevel = "WARN"
	}

	logConfig := logger.Config{
		Level:          parseLogLevel(finalLogLevel),
		EnableDebug:    parseLogLevel(finalLogLevel) == slog.LevelDebug,
		OutputToFile:   false, // ファイル出力は無効
		OutputToStderr: true,  // 標準エラー出力のみ
		LogFilePath:    "",    // ファイルパスは不要
		AppName:        "todotui",
	}

	if err := logger.Init(logConfig); err != nil {
		// ログ初期化に失敗した場合のみstderrに出力
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logger: %v\n", err)
	}

	// ログシステム開始メッセージ
	logger.Info("todotui started", "version", version, "commit", commit)

	// Determine todo file path with priority: CLI argument > config file > error
	var finalTodoFile string
	if todoFile != "" {
		// Priority 1: CLI argument specified
		finalTodoFile = todoFile
		logger.Debug("TODO file specified via CLI", "file", finalTodoFile)
	} else if appConfig.DefaultTodoFile != "" {
		// Priority 2: Config file specified (already expanded in LoadConfig)
		finalTodoFile = appConfig.DefaultTodoFile
		logger.Debug("TODO file specified via config", "file", finalTodoFile)
	} else {
		// No todo file specified
		logger.Error("No todo file specified")
		fmt.Fprintf(os.Stderr, "Error: No todo file specified. Use CLI argument or set default_todo_file in config.\n")
		fmt.Fprintf(os.Stderr, "Example: %s ~/todo.txt\n", os.Args[0])
		os.Exit(1)
	}

	// Expand ~ in path if present (for CLI arguments)
	finalTodoFile = ui.ExpandHomePath(finalTodoFile)

	// Create model with configuration
	logger.Debug("Initializing model", "todo_file", finalTodoFile, "theme", appConfig.Theme)
	model, err := ui.NewModel(finalTodoFile, appConfig)
	if err != nil {
		logger.Error("Failed to initialize model", "error", err)
		fmt.Fprintf(os.Stderr, "Error initializing: %v\n", err)
		os.Exit(1)
	}
	defer model.Cleanup()

	// Start the Bubble Tea program
	logger.Info("Starting Bubble Tea program")
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("Failed to run Bubble Tea program", "error", err)
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	logger.Info("todotui exited")
}
