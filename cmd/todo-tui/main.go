package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuucu/todo-tui/internal/ui"
)

func printUsage() {
	fmt.Printf(`Usage: %s [OPTIONS] [TODO_FILE]

A beautiful terminal todo.txt manager with vim-like keybindings.

Arguments:
  TODO_FILE    Path to todo.txt file (default: ~/todo.txt)

Options:
  -t, --theme THEME    Set color theme (catppuccin, nord)
  --list-themes        List available themes
  -h, --help          Show this help message

Environment Variables:
  TODO_TUI_THEME         Set default theme (same as --theme)

Examples:
  %s                           # Use default ~/todo.txt with catppuccin theme
  %s my-tasks.txt              # Use custom file
  %s -t nord                   # Use nord theme
  %s --theme catppuccin ~/todo.txt  # Use catppuccin theme with custom file

Keybindings:
  Tab/h/l      Switch between panes
  j/k          Navigate up/down
  a            Add new task
  e            Edit selected task
  Enter        Select filter & move to tasks (left pane) / Complete task (right pane)
  d            Delete task (with confirmation)
  r            Restore deleted/completed task
  q            Quit

Priority Configuration:
  Set TODO_TUI_PRIORITY_LEVELS environment variable to customize priority levels.
  Example: TODO_TUI_PRIORITY_LEVELS="A,B,C" for A→B→C cycle
  Default: A,B,C,D
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	// Setup IME environment for Japanese input support
	ui.SetupIMEEnvironment()
	
	var todoFile string
	var themeName string

	// Parse command line arguments
	args := os.Args[1:]
	i := 0
	for i < len(args) {
		switch args[i] {
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

	// Set theme environment variable if specified
	if themeName != "" {
		os.Setenv("TODO_TUI_THEME", themeName)
	}

	// Default todo file path
	if todoFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		todoFile = filepath.Join(homeDir, "todo.txt")
	}

	// Create model
	model, err := ui.NewModel(todoFile)
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
