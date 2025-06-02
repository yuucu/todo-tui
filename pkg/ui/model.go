package ui

import (
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
	"github.com/yuucu/todotui/pkg/domain"
	"github.com/yuucu/todotui/pkg/logger"
	"github.com/yuucu/todotui/pkg/todo"
)

// watchFile watches for changes to the todo file
func (m *Model) watchFile() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			for {
				select {
				case event := <-m.watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						return TaskListChangedMsg{}
					}
				case err := <-m.watcher.Errors:
					logger.Error("File watcher error", "error", err)
				}
			}
		},
	)
}

// NewModel creates a new model with the given todo file and configuration
func NewModel(todoFile string, appConfig AppConfig) (*Model, error) {
	logger.Debug("Creating new model", "todo_file", todoFile, "theme", appConfig.Theme)

	// Load tasks from file
	taskList, err := todo.Load(todoFile)
	if err != nil {
		logger.Error("Failed to load tasks from file", "file", todoFile, "error", err)
		return nil, err
	}

	logger.Info("Loaded tasks from file", "file", todoFile, "task_count", len(taskList))

	// Get theme for the lists
	currentTheme := GetTheme(appConfig.Theme)

	model := &Model{
		filterList:       SimpleList{},
		taskList:         SimpleList{},
		todoFilePath:     todoFile,
		tasks:            domain.NewTasks(taskList),
		activePane:       paneFilter,
		viewMode:         ViewFilter,
		watcher:          nil,
		width:            DefaultTerminalWidth,
		height:           DefaultTerminalHeight,
		currentTheme:     &currentTheme,
		appConfig:        appConfig,
		editBuffer:       "",
		originalTask:     "",
		statusMessage:    "",
		statusMessageEnd: time.Now(),
		textarea:         nil,
		imeHelper:        nil,
		editingTask:      nil,
	}

	// Initialize enhanced list features
	model.filterList.SetTheme(&currentTheme)
	model.filterList.SetTaskList(false) // Filter list is not a task list

	model.taskList.SetTheme(&currentTheme)
	model.taskList.SetTaskList(true) // Task list requires special rendering

	// Initialize help content
	model.initializeHelpContent()

	// TODO: Initialize textarea properly
	// model.textarea.Placeholder = TaskInputPlaceholder
	// model.textarea.CharLimit = TextAreaCharLimit
	// model.textarea.SetWidth(DefaultTerminalWidth)
	// model.textarea.SetHeight(TextAreaHeight)

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
func (m *Model) saveAndRefresh() tea.Cmd {
	logger.Debug("Saving tasks to file", "file", m.todoFilePath, "task_count", m.tasks.Len())
	if err := todo.Save(m.tasks.ToTaskList(), m.todoFilePath); err != nil {
		logger.Error("Failed to save tasks to file", "file", m.todoFilePath, "error", err)
		return nil
	}
	logger.Debug("Tasks saved successfully", "file", m.todoFilePath)
	m.refreshLists()
	return nil
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
		if m.viewMode == ViewHelp {
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
				m.viewMode = ViewFilter
				m.helpScroll = 0 // Reset scroll position
				return m, nil
			}
		}

		// Handle input mode (add/edit)
		if m.viewMode == ViewAdd || m.viewMode == ViewEdit {
			switch msg.String() {
			case ctrlCKey, escKey:
				// Cancel input
				m.viewMode = ViewFilter
				m.editingTask = nil
				// TODO: m.textarea.SetValue("")
				return m, nil
			case enterKey, ctrlSKey:
				// Save task
				text := strings.TrimSpace(m.editBuffer)
				if text != "" {
					if m.viewMode == ViewAdd {
						// Create new task
						task, err := todotxt.ParseTask(text)
						if err == nil {
							taskList := m.tasks.ToTaskList()
							taskList = append(taskList, *task)
							m.tasks = domain.NewTasks(taskList)
						}
						// TODO: Add error handling for failed task parsing
					} else if m.viewMode == ViewEdit {
						// Update existing task
						if m.editingTask != nil {
							// Parse the edited text and update the task
							if newTask, err := todotxt.ParseTask(text); err == nil {
								*m.editingTask = *newTask
								// Find and update the task in the main list
								taskList := m.tasks.ToTaskList()
								for i := range taskList {
									if taskList[i].String() == m.originalTask {
										taskList[i] = *newTask
										break
									}
								}
								m.tasks = domain.NewTasks(taskList)
							}
						}
					}
					m.viewMode = ViewFilter
					m.editingTask = nil
					// TODO: m.textarea.SetValue("")
					return m, m.saveAndRefresh()
				}
				// If text is empty, just cancel the edit
				m.viewMode = ViewFilter
				m.editingTask = nil
				// TODO: m.textarea.SetValue("")
				return m, nil
			default:
				// Handle text input for editBuffer
				m.editBuffer += msg.String()
				return m, nil
			}
		}

		// Handle normal mode keys
		switch msg.String() {
		case helpKey:
			// Show help
			m.viewMode = ViewHelp
			m.helpScroll = 0 // Reset scroll position
			return m, nil
		case qKey, ctrlCKey:
			return m, tea.Quit
		case aKey:
			// Add new task
			m.viewMode = ViewAdd
			m.editBuffer = ""
			// TODO: m.textarea.Focus()
			return m, nil
		case eKey:
			// Edit selected task
			if m.activePane == paneTask {
				if m.taskList.selected < m.filteredTasks.Len() {
					m.viewMode = ViewEdit
					// Store the task content for editing
					selectedTask := m.filteredTasks.Get(m.taskList.selected)
					m.editBuffer = selectedTask.String()
					m.originalTask = selectedTask.String()
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
			// Toggle task completion
			if m.taskList.selected < m.filteredTasks.Len() {
				taskToToggle := m.filteredTasks.Get(m.taskList.selected)
				// Find the task in main tasks list and toggle completion using domain model
				index, task, found := m.findTaskInList(taskToToggle)
				if found {
					// Toggle completion directly on the domain task
					isCompleted := task.ToggleCompletion()
					// Update the task in the list
					taskList := m.tasks.ToTaskList()
					taskList[index] = *task.ToTodoTxtTask()
					m.tasks = domain.NewTasks(taskList)

					// Show status message
					if isCompleted {
						return m, tea.Batch(
							m.saveAndRefresh(),
							m.setStatusMessage("✅ Task completed", 2*time.Second),
						)
					} else {
						return m, tea.Batch(
							m.saveAndRefresh(),
							m.setStatusMessage("🔄 Task marked as incomplete", 2*time.Second),
						)
					}
				}
				return m, m.saveAndRefresh()
			}
		case dKey:
			if m.activePane == paneTask {
				// Delete task directly (only for non-deleted tasks)
				if m.taskList.selected < m.filteredTasks.Len() {
					// Check if current filter is "Deleted Tasks"
					if m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name != FilterDeletedTasks {
						// Soft delete with deleted_at field
						taskToDelete := m.filteredTasks.Get(m.taskList.selected)
						// Find the task in main tasks list and add deleted_at field using lo
						index, task, found := m.findTaskInList(taskToDelete)
						if found {
							// Add deleted_at field to mark as soft deleted
							currentDate := time.Now().Format(DateFormat)
							taskString := task.String()

							// Add deleted_at field to the task string
							if !strings.Contains(taskString, TaskFieldDeletedPrefix) {
								taskString += " " + TaskFieldDeletedPrefix + currentDate

								// Parse the modified task string back to update the task
								if newTask, err := todotxt.ParseTask(taskString); err == nil {
									m.updateTaskAtIndex(index, domain.NewTask(newTask))
								}
							}
						}
						return m, m.saveAndRefresh()
					}
				}
			}
		case pKey:
			if m.activePane == paneTask {
				// Toggle priority level
				if m.taskList.selected < m.filteredTasks.Len() {
					taskToUpdate := m.filteredTasks.Get(m.taskList.selected)
					// Find the task in main tasks list and cycle priority using lo
					index, task, found := m.findTaskInList(taskToUpdate)
					if found {
						todoTxtTask := task.ToTodoTxtTask()
						m.cyclePriority(todoTxtTask)
						// Update the task in the list
						taskList := m.tasks.ToTaskList()
						taskList[index] = *todoTxtTask
						m.tasks = domain.NewTasks(taskList)
					}
					return m, m.saveAndRefresh()
				}
			}
		case tKey:
			if m.activePane == paneTask {
				// Toggle due date to today
				if m.taskList.selected < m.filteredTasks.Len() {
					taskToUpdate := m.filteredTasks.Get(m.taskList.selected)
					// Find the task in main tasks list and toggle due date using lo
					index, task, found := m.findTaskInList(taskToUpdate)
					if found {
						todoTxtTask := task.ToTodoTxtTask()
						m.toggleDueToday(todoTxtTask)
						// Update the task in the list
						taskList := m.tasks.ToTaskList()
						taskList[index] = *todoTxtTask
						m.tasks = domain.NewTasks(taskList)
					}
					return m, m.saveAndRefresh()
				}
			}
		case rKey:
			if m.activePane == paneTask {
				// Restore deleted or completed task
				if m.taskList.selected < m.filteredTasks.Len() {
					currentFilter := ""
					if m.filterList.selected < len(m.filters) {
						currentFilter = m.filters[m.filterList.selected].name
					}

					taskToRestore := m.filteredTasks.Get(m.taskList.selected)

					// Handle deleted tasks restoration
					if currentFilter == FilterDeletedTasks {
						// Find the task in main tasks list and remove deleted_at field using lo
						index, task, found := m.findTaskInList(taskToRestore)
						if found {
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
									m.updateTaskAtIndex(index, domain.NewTask(newTask))
								}
							}
						}
						return m, m.saveAndRefresh()
					}

					// Handle completed tasks restoration
					if currentFilter == "Completed Tasks" {
						// Find the task in main tasks list and toggle completion using domain model
						index, task, found := m.findTaskInList(taskToRestore)
						if found {
							task.ToggleCompletion() // This will mark as incomplete
							// Update the task in the list
							taskList := m.tasks.ToTaskList()
							taskList[index] = *task.ToTodoTxtTask()
							m.tasks = domain.NewTasks(taskList)
						}
						return m, m.saveAndRefresh()
					}
				}
			}
		case yKey:
			if m.activePane == paneTask {
				// Copy task text to clipboard
				if m.taskList.selected < m.filteredTasks.Len() {
					taskToCopy := m.filteredTasks.Get(m.taskList.selected)
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

		// TODO: Also update textarea size if we're in add/edit mode
		// if m.width > TextAreaPadding {
		//     m.textarea.SetWidth(m.width - TextAreaPadding)
		// }

		// Force refresh of lists to apply new sizes and ensure content fits
		m.refreshLists()

		// Return nil to re-render without clearing screen
		return m, nil
	case TaskListChangedMsg:
		// Reload tasks from file
		if taskList, err := todo.Load(m.todoFilePath); err == nil {
			m.tasks = domain.NewTasks(taskList)
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
	now := time.Now()
	today := now.Format(DateFormat)

	// Get the current task string
	taskString := task.String()

	// Check if task is already due today using domain method
	domainTask := domain.NewTask(task)
	hasDueToday := domainTask.IsDueToday(now)

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

// findTaskInList finds a task in the main task list and returns its index and domain task
func (m *Model) findTaskInList(targetTask domain.Task) (int, domain.Task, bool) {
	for i := 0; i < m.tasks.Len(); i++ {
		task := m.tasks.Get(i)
		if task.String() == targetTask.String() {
			return i, task, true
		}
	}
	return -1, domain.Task{}, false
}

// updateTaskAtIndex updates a task at the given index using domain.Task
func (m *Model) updateTaskAtIndex(index int, newTask *domain.Task) {
	if index >= 0 && index < m.tasks.Len() {
		taskList := m.tasks.ToTaskList()
		taskList[index] = *newTask.ToTodoTxtTask()
		m.tasks = domain.NewTasks(taskList)
	}
}
