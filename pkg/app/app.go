package app

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

// Config represents the parsed command line configuration
type Config struct {
	ConfigFile  string
	ThemeName   string
	TodoFile    string
	ShowVersion bool
	ShowHelp    bool
}

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
  -v, --version             Show version information
  -h, --help               Show this help message

For detailed documentation and keybindings, see: https://github.com/yuucu/todotui
`, os.Args[0])
}

// ParseFlags parses command line flags and returns configuration
func ParseFlags() (*Config, error) {
	// Define command line flags
	var (
		configFile  = flag.String("config", "", "Path to configuration file")
		themeName   = flag.String("theme", "", "Set color theme (catppuccin, nord, everforest-dark, everforest-light)")
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show this help message")
	)

	// Define short flag aliases
	flag.StringVar(configFile, "c", "", "Path to configuration file")
	flag.StringVar(themeName, "t", "", "Set color theme (catppuccin, nord, everforest-dark, everforest-light)")
	flag.BoolVar(showVersion, "v", false, "Show version information")
	flag.BoolVar(showHelp, "h", false, "Show this help message")

	// Set custom usage function
	flag.Usage = printUsage

	// Parse command line flags
	flag.Parse()

	// Get remaining non-flag arguments (todo file)
	args := flag.Args()
	var todoFile string
	if len(args) > 0 {
		todoFile = args[0]
	}
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments")
	}

	return &Config{
		ConfigFile:  *configFile,
		ThemeName:   *themeName,
		TodoFile:    todoFile,
		ShowVersion: *showVersion,
		ShowHelp:    *showHelp,
	}, nil
}

// RunBubbleTea runs the bubble tea application with given configuration
func RunBubbleTea(config *Config) error {
	// Setup IME environment for Japanese input support
	ui.SetupIMEEnvironment()

	// Load configuration
	appConfig := ui.LoadConfig(config.ConfigFile)

	// Override theme if specified via command line
	if config.ThemeName != "" {
		appConfig.Theme = config.ThemeName
	}

	// Initialize logging system
	var finalLogLevel string
	if appConfig.Logging.LogLevel != "" {
		finalLogLevel = appConfig.Logging.LogLevel
	} else {
		finalLogLevel = "WARN"
	}

	logConfig := logger.Config{
		Level:          parseLogLevel(finalLogLevel),
		OutputToStderr: true,
		AppName:        "todotui",
	}

	if err := logger.Init(logConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logger: %v\n", err)
	}

	logger.Info("todotui started", "version", version, "commit", commit)

	// Determine todo file path
	var finalTodoFile string
	if config.TodoFile != "" {
		finalTodoFile = config.TodoFile
		logger.Debug("TODO file specified via CLI", "file", finalTodoFile)
	} else if appConfig.DefaultTodoFile != "" {
		finalTodoFile = appConfig.DefaultTodoFile
		logger.Debug("TODO file specified via config", "file", finalTodoFile)
	} else {
		logger.Error("No todo file specified")
		return fmt.Errorf("no todo file specified. Use CLI argument or set default_todo_file in config")
	}

	// Expand ~ in path if present
	finalTodoFile = ui.ExpandHomePath(finalTodoFile)

	// Create model with configuration
	logger.Debug("Initializing model", "todo_file", finalTodoFile, "theme", appConfig.Theme)
	model, err := ui.NewModel(finalTodoFile, appConfig)
	if err != nil {
		logger.Error("Failed to initialize model", "error", err)
		return fmt.Errorf("failed to initialize model: %w", err)
	}
	defer model.Cleanup()

	// Start the Bubble Tea program
	logger.Info("Starting Bubble Tea program")
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("Failed to run Bubble Tea program", "error", err)
		return fmt.Errorf("failed to run Bubble Tea program: %w", err)
	}

	logger.Info("todotui exited")
	return nil
}

// Run is the main entry point
func Run() error {
	config, err := ParseFlags()
	if err != nil {
		return err
	}

	// Handle help and version flags
	if config.ShowHelp {
		flag.Usage()
		return nil
	}

	if config.ShowVersion {
		printVersion()
		return nil
	}

	return RunBubbleTea(config)
}
