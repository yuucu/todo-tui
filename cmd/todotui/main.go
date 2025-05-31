package main

import (
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
  TODO_FILE    Path to todo.txt file (default: ~/todo.txt)

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

	var todoFile string
	var themeName string
	var configFile string

	// Parse command line arguments
	args := os.Args[1:]
	i := 0
	for i < len(args) {
		switch args[i] {
		case "--config", "-c":
			if i+1 >= len(args) {
				fmt.Printf("Error: --config requires a value\n")
				printUsage()
				os.Exit(1)
			}
			configFile = args[i+1]
			i += 2
		case "--theme", "-t":
			if i+1 >= len(args) {
				fmt.Printf("Error: --theme requires a value\n")
				printUsage()
				os.Exit(1)
			}
			themeName = args[i+1]
			i += 2
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		case "--version", "-v":
			printVersion()
			os.Exit(0)
		default:
			if todoFile == "" {
				todoFile = args[i]
			} else {
				fmt.Printf("Error: unexpected argument %s\n", args[i])
				printUsage()
				os.Exit(1)
			}
			i++
		}
	}

	// Load configuration
	appConfig := ui.LoadConfig(configFile)

	// Override theme if specified via command line
	if themeName != "" {
		appConfig.Theme = themeName
	}

	// Determine todo file path
	if todoFile == "" {
		// Use from config if not specified
		todoFile = appConfig.DefaultTodoFile
	}

	// Create model with configuration
	model, err := ui.NewModel(todoFile, appConfig)
	if err != nil {
		fmt.Printf("Error initializing: %v\n", err)
		os.Exit(1)
	}
	defer model.Cleanup()

	// Start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
