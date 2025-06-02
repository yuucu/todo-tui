package ui

import (
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/fsnotify/fsnotify"
	"github.com/yuucu/todotui/pkg/domain"
)

// ViewMode represents the current view mode
type ViewMode int

const (
	ViewFilter ViewMode = iota
	ViewTask
	ViewHelp
	ViewEdit
	ViewAdd
)

// Pane represents which pane is active
type Pane int

const (
	paneFilter Pane = iota
	paneTask
)

// FilterData holds information about a filter
type FilterData struct {
	name     string
	filterFn func(domain.Tasks) domain.Tasks
	count    int
}

// StatusMessageClearMsg is a message to clear the status message
type StatusMessageClearMsg struct{}

// TaskListChangedMsg is sent when the task list file changes
type TaskListChangedMsg struct{}

// HelpContent represents help information for key bindings
type HelpContent struct {
	Category string
	Items    []HelpItem
}

// HelpItem represents a single help item
type HelpItem struct {
	Key         string
	Description string
}

// Model represents the main application model
type Model struct {
	appConfig        AppConfig
	tasks            domain.Tasks
	filterList       SimpleList
	taskList         SimpleList
	filters          []FilterData
	filteredTasks    domain.Tasks
	activePane       Pane
	viewMode         ViewMode
	width            int
	height           int
	currentTheme     *Theme
	todoFilePath     string
	statusMessage    string
	statusMessageEnd time.Time
	editBuffer       string
	originalTask     string
	watcher          *fsnotify.Watcher
	helpContent      []HelpContent // Help content for key bindings
	helpScroll       int           // Current scroll position in help view
	textarea         interface{}   // Placeholder for textarea
	imeHelper        interface{}   // Placeholder for imeHelper
	editingTask      *todotxt.Task // Currently editing task
}
