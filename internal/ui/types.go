package ui

import (
	todotxt "github.com/1set/todotxt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/fsnotify/fsnotify"
	"os"
	"strings"
)

type pane int
type mode int

const (
	paneFilter pane = iota
	paneTask
)

const (
	modeView mode = iota
	modeAdd
	modeEdit
	modeDeleteConfirm
)

// FileChangedMsg is sent when the todo file is modified
type FileChangedMsg struct{}

// FilterData represents filter information
type FilterData struct {
	name     string
	filterFn func(todotxt.TaskList) todotxt.TaskList
	count    int
}

// Config represents application configuration
type Config struct {
	// Priority levels in order (empty string means no priority)
	PriorityLevels []string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	// Try to read priority levels from environment variable
	envPriorities := os.Getenv("TODO_TUI_PRIORITY_LEVELS")
	if envPriorities != "" {
		// Parse comma-separated priority levels
		levels := strings.Split(envPriorities, ",")
		for i, level := range levels {
			levels[i] = strings.TrimSpace(level)
		}
		// Ensure empty string is first (no priority)
		if len(levels) > 0 && levels[0] != "" {
			levels = append([]string{""}, levels...)
		}
		return Config{
			PriorityLevels: levels,
		}
	}
	
	return Config{
		PriorityLevels: []string{"", "A", "B", "C", "D"}, // "" means no priority, then A->B->C->D
	}
}

// Model represents the main application state
type Model struct {
	filterList    SimpleList
	taskList      SimpleList
	textarea      textarea.Model
	todoFile      string
	tasks         todotxt.TaskList
	activePane    pane
	currentMode   mode
	editingTask   *todotxt.Task
	watcher       *fsnotify.Watcher
	width         int
	height        int
	filters       []FilterData
	filteredTasks todotxt.TaskList
	deleteIndex   int // Index of task to delete in filteredTasks
	currentTheme  Theme
	config        Config // Application configuration
	imeHelper     *IMEHelper // Add IME helper for Japanese input
} 