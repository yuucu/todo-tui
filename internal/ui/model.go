package ui

import (
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
	"github.com/yuucu/todotui/internal/todo"
	"github.com/yuucu/todotui/pkg/logger"
)

// watchFile returns a command that watches for file changes
func (m *Model) watchFile() tea.Cmd {
	return func() tea.Msg {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				return FileChangedMsg{}
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return nil
			}
			// Handle error if needed
			_ = err
		}
		return nil
	}
}

// NewModel creates a new model instance
func NewModel(todoFile string, appConfig AppConfig) (*Model, error) {
	logger.Debug("Creating new model", "todo_file", todoFile, "theme", appConfig.Theme)

	// Load tasks from file
	taskList, err := todo.Load(todoFile)
	if err != nil {
		logger.Error("Failed to load tasks from file", "file", todoFile, "error", err)
		return nil, err
	}

	logger.Info("Loaded tasks from file", "file", todoFile, "task_count", len(taskList))

	model := &Model{
		filterList:   SimpleList{},
		taskList:     SimpleList{},
		textarea:     textarea.New(),
		todoFile:     todoFile,
		tasks:        taskList,
		activePane:   paneFilter,
		currentMode:  modeView,
		watcher:      nil,
		width:        DefaultTerminalWidth,
		height:       DefaultTerminalHeight,
		currentTheme: GetTheme(appConfig.Theme),
		appConfig:    appConfig,
		imeHelper:    NewIMEHelper(),
	}

	// Initialize help content
	model.initializeHelpContent()

	// Initialize textarea
	model.textarea.Placeholder = TaskInputPlaceholder
	model.textarea.CharLimit = TextAreaCharLimit
	model.textarea.SetWidth(DefaultTerminalWidth)
	model.textarea.SetHeight(TextAreaHeight)

	// Initialize file watcher
	logger.Debug("Initializing file watcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("Failed to create file watcher", "error", err)
		return nil, err
	}
	model.watcher = watcher

	// Watch the todo file
	err = watcher.Add(todoFile)
	if err != nil {
		logger.Error("Failed to watch todo file", "file", todoFile, "error", err)
		watcher.Close()
		return nil, err
	}

	logger.Debug("File watcher initialized successfully", "file", todoFile)
	return model, nil
}

// saveAndRefresh saves the task list and refreshes the UI
func (m *Model) saveAndRefresh() {
	logger.Debug("Saving tasks to file", "file", m.todoFile, "task_count", len(m.tasks))
	if err := todo.Save(m.tasks, m.todoFile); err != nil {
		logger.Error("Failed to save tasks to file", "file", m.todoFile, "error", err)
		return
	}
	logger.Info("Tasks saved successfully", "file", m.todoFile)
	m.refreshLists()
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Initialize pane sizes first
	m.updatePaneSizes()
	m.refreshLists()
	return m.watchFile()
}

// Update handles key input and state changes
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle help mode with scrolling support
		if m.currentMode == modeHelp {
			switch msg.String() {
			case jKey, downKey:
				// Scroll down
				// Estimate total lines (rough calculation for responsive scrolling)
				estimatedLines := len(m.helpContent)*8 + 20 // Approximate lines per category + header/footer
				maxScroll := max(0, estimatedLines-(m.height-8))
				if m.helpScroll < maxScroll {
					m.helpScroll++
				}
				return m, nil
			case kKey, upKey:
				// Scroll up
				if m.helpScroll > 0 {
					m.helpScroll--
				}
				return m, nil
			case gKey:
				// Go to top
				m.helpScroll = 0
				return m, nil
			case GKey:
				// Go to bottom
				estimatedLines := len(m.helpContent)*8 + 20
				maxScroll := max(0, estimatedLines-(m.height-8))
				m.helpScroll = maxScroll
				return m, nil
			default:
				// Any other key exits help
				m.currentMode = modeView
				m.helpScroll = 0 // Reset scroll position
				return m, nil
			}
		}

		// Handle input mode (add/edit)
		if m.currentMode == modeAdd || m.currentMode == modeEdit {
			switch msg.String() {
			case ctrlCKey, escKey:
				// Cancel input
				m.currentMode = modeView
				m.editingTask = nil
				m.textarea.SetValue("")
				return m, nil
			case enterKey, ctrlSKey:
				// Save task
				text := strings.TrimSpace(m.textarea.Value())
				if text != "" {
					if m.currentMode == modeAdd {
						// Create new task
						task, err := todotxt.ParseTask(text)
						if err == nil {
							m.tasks = append(m.tasks, *task)
						}
						// TODO: Add error handling for failed task parsing
					} else if m.currentMode == modeEdit && m.editingTask != nil {
						// Update existing task
						newTask, err := todotxt.ParseTask(text)
						if err == nil {
							*m.editingTask = *newTask
						}
						// TODO: Add error handling for failed task parsing
					}
					m.currentMode = modeView
					m.editingTask = nil
					m.textarea.SetValue("")
					m.saveAndRefresh()
					return m, nil
				}
				// If text is empty, just cancel the edit
				m.currentMode = modeView
				m.editingTask = nil
				m.textarea.SetValue("")
				return m, nil
			default:
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			}
		}

		// Handle normal mode keys
		switch msg.String() {
		case helpKey:
			// Show help
			m.currentMode = modeHelp
			m.helpScroll = 0 // Reset scroll position
			return m, nil
		case qKey, ctrlCKey:
			return m, tea.Quit
		case aKey:
			// Add new task
			m.currentMode = modeAdd
			m.textarea.SetValue("")
			m.textarea.Focus()
			return m, nil
		case eKey:
			// Edit selected task
			if m.activePane == paneTask {
				if m.taskList.selected < len(m.filteredTasks) {
					m.currentMode = modeEdit
					// Store the reference to the original task, not the filtered one
					selectedTask := m.filteredTasks[m.taskList.selected]
					// Find the original task in m.tasks using lo
					_, task := m.findTaskInList(selectedTask)
					if task != nil {
						m.editingTask = task // Point to the original task
						m.textarea.SetValue(m.editingTask.String())
						m.textarea.Focus()
					}
					return m, nil
				}
			}
		case tabKey:
			// Switch between panes
			if m.activePane == paneFilter {
				m.activePane = paneTask
			} else {
				m.activePane = paneFilter
			}
			return m, nil
		case hKey:
			// Move to left pane (filter)
			m.activePane = paneFilter
			return m, nil
		case lKey:
			// Move to right pane (task)
			m.activePane = paneTask
			return m, nil
		case enterKey:
			if m.activePane == paneFilter {
				// Filter selection changed, refresh task list
				m.refreshTaskList()
				// Move to right pane (task)
				m.activePane = paneTask
				return m, nil
			}
			// Complete task
			if m.taskList.selected < len(m.filteredTasks) {
				taskToComplete := m.filteredTasks[m.taskList.selected]
				// Find the task in main tasks list and mark as completed using lo
				_, task := m.findTaskInList(taskToComplete)
				if task != nil {
					task.Complete()
				}
				m.saveAndRefresh()
				return m, nil
			}
		case dKey:
			if m.activePane == paneTask {
				// Delete task directly (only for non-deleted tasks)
				if m.taskList.selected < len(m.filteredTasks) {
					// Check if current filter is "Deleted Tasks"
					if m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name != FilterDeletedTasks {
						// Soft delete with deleted_at field
						taskToDelete := m.filteredTasks[m.taskList.selected]
						// Find the task in main tasks list and add deleted_at field using lo
						index, task := m.findTaskInList(taskToDelete)
						if task != nil {
							// Add deleted_at field to mark as soft deleted
							currentDate := time.Now().Format(DateFormat)
							taskString := task.String()

							// Add deleted_at field to the task string
							if !strings.Contains(taskString, TaskFieldDeletedPrefix) {
								taskString += " " + TaskFieldDeletedPrefix + currentDate

								// Parse the modified task string back to update the task
								if newTask, err := todotxt.ParseTask(taskString); err == nil {
									m.tasks[index] = *newTask
								}
							}
						}
						m.saveAndRefresh()
						return m, nil
					}
				}
			}
		case pKey:
			if m.activePane == paneTask {
				// Toggle priority level
				if m.taskList.selected < len(m.filteredTasks) {
					taskToUpdate := m.filteredTasks[m.taskList.selected]
					// Find the task in main tasks list and cycle priority using lo
					_, task := m.findTaskInList(taskToUpdate)
					if task != nil {
						m.cyclePriority(task)
					}
					m.saveAndRefresh()
					return m, nil
				}
			}
		case tKey:
			if m.activePane == paneTask {
				// Toggle due date to today
				if m.taskList.selected < len(m.filteredTasks) {
					taskToUpdate := m.filteredTasks[m.taskList.selected]
					// Find the task in main tasks list and toggle due date using lo
					_, task := m.findTaskInList(taskToUpdate)
					if task != nil {
						m.toggleDueToday(task)
					}
					m.saveAndRefresh()
					return m, nil
				}
			}
		case rKey:
			if m.activePane == paneTask {
				// Restore deleted or completed task
				if m.taskList.selected < len(m.filteredTasks) {
					currentFilter := ""
					if m.filterList.selected < len(m.filters) {
						currentFilter = m.filters[m.filterList.selected].name
					}

					taskToRestore := m.filteredTasks[m.taskList.selected]

					// Handle deleted tasks restoration
					if currentFilter == FilterDeletedTasks {
						// Find the task in main tasks list and remove deleted_at field using lo
						index, task := m.findTaskInList(taskToRestore)
						if task != nil {
							// Remove deleted_at field to restore the task
							taskString := task.String()

							// Remove deleted_at field from the task string using lo.Filter
							if strings.Contains(taskString, "deleted_at:") {
								parts := strings.Fields(taskString)
								cleanParts := lo.Filter(parts, func(part string, _ int) bool {
									return !strings.HasPrefix(part, "deleted_at:")
								})
								taskString = strings.Join(cleanParts, " ")

								// Parse the modified task string back to update the task
								if newTask, err := todotxt.ParseTask(taskString); err == nil {
									m.tasks[index] = *newTask
								}
							}
						}
						m.saveAndRefresh()
						return m, nil
					}

					// Handle completed tasks restoration
					if currentFilter == "Completed Tasks" {
						// Find the task in main tasks list and mark as incomplete using lo
						_, task := m.findTaskInList(taskToRestore)
						if task != nil {
							task.Completed = false
							task.CompletedDate = time.Time{} // Clear completion date
						}
						m.saveAndRefresh()
						return m, nil
					}
				}
			}
		case yKey:
			if m.activePane == paneTask {
				// Copy task text to clipboard
				if m.taskList.selected < len(m.filteredTasks) {
					taskToCopy := m.filteredTasks[m.taskList.selected]
					taskText := taskToCopy.String()

					// Copy to clipboard
					if err := clipboard.WriteAll(taskText); err == nil {
						// Show success message for 2 seconds
						return m, m.setStatusMessage("📋 Task copied to clipboard", 2*time.Second)
					}
					// Show error message for 3 seconds if copy failed
					return m, m.setStatusMessage("❌ Failed to copy task", 3*time.Second)
				}
			}
		case jKey, downKey:
			if m.activePane == paneFilter {
				m.filterList.MoveDown()
				m.refreshTaskList()
			} else {
				m.taskList.MoveDown()
			}
			return m, nil
		case kKey, upKey:
			if m.activePane == paneFilter {
				m.filterList.MoveUp()
				m.refreshTaskList()
			} else {
				m.taskList.MoveUp()
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update pane sizes with new terminal dimensions
		m.updatePaneSizes()

		// Also update textarea size if we're in add/edit mode
		if m.width > TextAreaPadding {
			m.textarea.SetWidth(m.width - TextAreaPadding)
		}

		// Force refresh of lists to apply new sizes and ensure content fits
		m.refreshLists()

		// Return nil to re-render without clearing screen
		return m, nil
	case FileChangedMsg:
		// Reload tasks from file
		if taskList, err := todo.Load(m.todoFile); err == nil {
			m.tasks = taskList
			m.refreshLists()
		}
		// Continue watching
		return m, m.watchFile()
	case StatusMessageClearMsg:
		// Clear status message if it has expired
		if time.Now().After(m.statusMessageEnd) {
			m.statusMessage = ""
		}
		return m, nil
	}

	return m, nil
}

// cyclePriority cycles through priority levels based on configuration
func (m *Model) cyclePriority(task *todotxt.Task) {
	currentPriority := ""
	if task.HasPriority() {
		currentPriority = task.Priority
	}

	// Find current priority index in configuration using lo.FindIndexOf
	_, currentIndex, found := lo.FindIndexOf(m.appConfig.PriorityLevels, func(priority string) bool {
		return priority == currentPriority
	})
	if !found {
		currentIndex = 0
	}

	// Move to next priority level (cycle around)
	nextIndex := (currentIndex + 1) % len(m.appConfig.PriorityLevels)
	nextPriority := m.appConfig.PriorityLevels[nextIndex]

	// Set the new priority
	if nextPriority == "" {
		task.Priority = ""
	} else {
		task.Priority = nextPriority
	}
}

// toggleDueToday toggles the due date of a task to today or removes it if already set to today
func (m *Model) toggleDueToday(task *todotxt.Task) {
	today := time.Now().Format(DateFormat)

	// Get the current task string
	taskString := task.String()

	// Check if task already has a due date
	hasDueToday := false
	if task.HasDueDate() && task.DueDate.Format(DateFormat) == today {
		hasDueToday = true
	}

	var newTaskString string

	if hasDueToday {
		// Remove due date - remove due:YYYY-MM-DD from task string using lo.Filter
		parts := strings.Fields(taskString)
		newParts := lo.Filter(parts, func(part string, _ int) bool {
			return !strings.HasPrefix(part, TaskFieldDuePrefix)
		})
		newTaskString = strings.Join(newParts, " ")
	} else {
		// Add or update due date
		if task.HasDueDate() {
			// Replace existing due date using lo.Map
			parts := strings.Fields(taskString)
			newParts := lo.Map(parts, func(part string, _ int) string {
				if strings.HasPrefix(part, TaskFieldDuePrefix) {
					return TaskFieldDuePrefix + today
				}
				return part
			})
			newTaskString = strings.Join(newParts, " ")
		} else {
			// Add new due date
			newTaskString = taskString + " " + TaskFieldDuePrefix + today
		}
	}

	// Parse the new task string and update the task
	if newTask, err := todotxt.ParseTask(newTaskString); err == nil {
		*task = *newTask
	}
}

// Cleanup closes the file watcher
func (m *Model) Cleanup() {
	if m.watcher != nil {
		m.watcher.Close()
	}
}

// setStatusMessage sets a temporary status message with auto-clear timer
func (m *Model) setStatusMessage(message string, duration time.Duration) tea.Cmd {
	m.statusMessage = message
	m.statusMessageEnd = time.Now().Add(duration)

	// Return a command that clears the status message after the duration
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return StatusMessageClearMsg{}
	})
}

// findTaskInList finds a task in the main task list and returns its index and pointer
func (m *Model) findTaskInList(targetTask todotxt.Task) (int, *todotxt.Task) {
	_, index, found := lo.FindIndexOf(m.tasks, func(task todotxt.Task) bool {
		return task.String() == targetTask.String()
	})
	if found {
		return index, &m.tasks[index]
	}
	return -1, nil
}
