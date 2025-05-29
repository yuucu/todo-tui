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

// TaskCache はタスクデータのキャッシュを管理
type TaskCache struct {
	projects       []string
	contexts       []string
	deletedTasks   []int // deleted taskのインデックス
	completedTasks []int // completed taskのインデックス
	lastRefresh    time.Time
	taskVersion    int // タスクリストのバージョン管理
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
	appConfig     AppConfig  // Application configuration
	imeHelper     *IMEHelper // Add IME helper for Japanese input
	taskCache     *TaskCache // Add task cache for performance optimization
}
