package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuucu/todotui/internal/ui"
)

func printUsage() {
	fmt.Printf(`Usage: %s [OPTIONS] [TODO_FILE]

A beautiful terminal todo.txt manager with vim-like keybindings.

Arguments:
  TODO_FILE    Path to todo.txt file (default: ~/todo.txt or from config file)

Options:
  -c, --config CONFIG  Path to configuration file (JSON format)
  -t, --theme THEME    Set color theme (catppuccin, nord)
  --list-themes        List available themes
  --init-config        Create a sample configuration file
  -h, --help          Show this help message

Environment Variables:
  TODO_TUI_THEME         Set default theme (same as --theme)
  TODO_TUI_PRIORITY_LEVELS  Set priority levels (example: "A,B,C")

Examples:
  %s                           # Use default ~/todo.txt with catppuccin theme
  %s my-tasks.txt              # Use custom file
  %s -c ~/.config/todotui/config.json  # Use configuration file
  %s -t nord                   # Use nord theme
  %s --theme catppuccin ~/todo.txt  # Use catppuccin theme with custom file
  %s --init-config             # Create sample configuration file

Keybindings:
  Tab/h/l      Switch between panes
  j/k          Navigate up/down
  a            Add new task
  e            Edit selected task
  Enter        Select filter & move to tasks (left pane) / Complete task (right pane)
  d            Delete task (with confirmation)
  r            Restore deleted/completed task
  p            Toggle task priority
  q            Quit

Priority Configuration:
  Set TODO_TUI_PRIORITY_LEVELS environment variable to customize priority levels.
  Example: TODO_TUI_PRIORITY_LEVELS="A,B,C" for A→B→C cycle
  Default: A,B,C,D
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	// Setup IME environment for Japanese input support
	ui.SetupIMEEnvironment()

	var todoFile string
	var themeName string
	var configFile string
	var initConfig bool

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
		case "--list-themes":
			fmt.Println("Available themes:")
			for _, theme := range ui.GetAvailableThemes() {
				fmt.Printf("  %s\n", theme)
			}
			os.Exit(0)
		case "--init-config":
			initConfig = true
			i++
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

	// Handle init config
	if initConfig {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		defaultConfigPath := filepath.Join(homeDir, ".config", "todotui", "config.json")
		if configFile != "" {
			defaultConfigPath = configFile
		}

		if err := ui.CreateSampleConfig(defaultConfigPath); err != nil {
			fmt.Printf("Error creating sample config: %v\n", err)
			os.Exit(1)
		}
		return
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
