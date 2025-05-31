package ui

import (
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/fsnotify/fsnotify"
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
	modeHelp
)

// FileChangedMsg is sent when the todo file is modified
type FileChangedMsg struct{}

// StatusMessageClearMsg is sent when a temporary status message should be cleared
type StatusMessageClearMsg struct{}

// FilterData represents filter information
type FilterData struct {
	name     string
	filterFn func(todotxt.TaskList) todotxt.TaskList
	count    int
}

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

// TaskCache はタスクデータのキャッシュを管理
type TaskCache struct {
	// 基本的なキャッシュ機能のみ残す（現在は使用していない）
}

// Model represents the main application state
type Model struct {
	filterList       SimpleList
	taskList         SimpleList
	textarea         textarea.Model
	todoFile         string
	tasks            todotxt.TaskList
	activePane       pane
	currentMode      mode
	editingTask      *todotxt.Task
	watcher          *fsnotify.Watcher
	width            int
	height           int
	filters          []FilterData
	filteredTasks    todotxt.TaskList
	deleteIndex      int // Index of task to delete in filteredTasks
	currentTheme     Theme
	appConfig        AppConfig     // Application configuration
	imeHelper        *IMEHelper    // Add IME helper for Japanese input
	helpContent      []HelpContent // Help content for key bindings
	helpScroll       int           // Current scroll position in help view
	statusMessage    string        // Temporary status message
	statusMessageEnd time.Time     // When status message should disappear
}
