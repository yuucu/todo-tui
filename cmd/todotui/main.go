package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuucu/todotui/internal/ui"
)

// Build information variables (set by goreleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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
  -c, --config CONFIG  Path to configuration file
  -t, --theme THEME    Set color theme (catppuccin, nord)
  -v, --version        Show version information
  -h, --help          Show this help message

For detailed documentation and keybindings, see: https://github.com/yuucu/todotui
`, os.Args[0])
}

func main() {
	// Setup IME environment for Japanese input support
	ui.SetupIMEEnvironment()

	// Define command line flags
	var (
		configFile  = flag.String("config", "", "Path to configuration file")
		themeName   = flag.String("theme", "", "Set color theme (catppuccin, nord)")
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show this help message")
	)

	// Define short flag aliases (reuse same help text)
	flag.StringVar(configFile, "c", "", "Path to configuration file")
	flag.StringVar(themeName, "t", "", "Set color theme (catppuccin, nord)")
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

	// Determine todo file path with priority: CLI argument > config file > error
	var finalTodoFile string
	if todoFile != "" {
		// Priority 1: CLI argument specified
		finalTodoFile = todoFile
	} else if appConfig.DefaultTodoFile != "" {
		// Priority 2: Config file specified (already expanded in LoadConfig)
		finalTodoFile = appConfig.DefaultTodoFile
	} else {
		// No todo file specified
		fmt.Fprintf(os.Stderr, "Error: No todo file specified. Use CLI argument or set default_todo_file in config.\n")
		fmt.Fprintf(os.Stderr, "Example: %s ~/todo.txt\n", os.Args[0])
		os.Exit(1)
	}

	// Expand ~ in path if present (for CLI arguments)
	finalTodoFile = ui.ExpandHomePath(finalTodoFile)

	// Create model with configuration
	model, err := ui.NewModel(finalTodoFile, appConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %v\n", err)
		os.Exit(1)
	}
	defer model.Cleanup()

	// Start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
