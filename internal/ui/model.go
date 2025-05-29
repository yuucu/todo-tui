package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	todotxt "github.com/1set/todotxt"
	"github.com/yuucu/todo-tui/internal/todo"
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
func NewModel(todoFile string) (*Model, error) {
	// Load tasks from file
	taskList, err := todo.Load(todoFile)
	if err != nil {
		return nil, err
	}

	model := &Model{
		todoFile:     todoFile,
		tasks:        taskList,
		activePane:   paneFilter,
		currentMode:  modeView,
		textarea:     textarea.New(),
		deleteIndex:  -1,
		currentTheme: GetTheme(),
		imeHelper:    NewIMEHelper(),
	}

	// Initialize textarea
	model.textarea.Placeholder = "タスクの説明を入力してください (例: '電話 @母 +home due:2025-01-15')"
	model.textarea.CharLimit = 0
	model.textarea.SetWidth(80)
	model.textarea.SetHeight(3)

	// Initialize file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	model.watcher = watcher

	// Watch the todo file
	err = watcher.Add(todoFile)
	if err != nil {
		watcher.Close()
		return nil, err
	}

	return model, nil
}

// saveAndRefresh saves the task list and refreshes the UI
func (m *Model) saveAndRefresh() tea.Cmd {
	if err := todo.Save(m.tasks, m.todoFile); err != nil {
		// TODO: Handle error properly
		return nil
	}
	m.refreshLists()
	return nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	m.refreshLists()
	return m.watchFile()
}

// Update handles key input and state changes
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle delete confirmation mode
		if m.currentMode == modeDeleteConfirm {
			switch msg.String() {
			case "y", "Y":
				// Confirm deletion - use soft delete with deleted_at field
				if m.deleteIndex >= 0 && m.deleteIndex < len(m.filteredTasks) {
					taskToDelete := m.filteredTasks[m.deleteIndex]
					// Find the task in main tasks list and add deleted_at field
					for i := 0; i < len(m.tasks); i++ {
						if m.tasks[i].String() == taskToDelete.String() {
							// Add deleted_at field to mark as soft deleted
							currentDate := time.Now().Format("2006-01-02")
							taskString := m.tasks[i].String()
							
							// Add deleted_at field to the task string
							if !strings.Contains(taskString, "deleted_at:") {
								taskString += " deleted_at:" + currentDate
								
								// Parse the modified task string back to update the task
								if newTask, err := todotxt.ParseTask(taskString); err == nil {
									m.tasks[i] = *newTask
								}
							}
							break
						}
					}
				}
				m.currentMode = modeView
				m.deleteIndex = -1
				return m, m.saveAndRefresh()
			case "n", "N", "esc", "ctrl+c":
				// Cancel deletion
				m.currentMode = modeView
				m.deleteIndex = -1
				return m, nil
			}
			return m, nil
		}
		
		// Handle input mode (add/edit)
		if m.currentMode == modeAdd || m.currentMode == modeEdit {
			switch msg.String() {
			case "ctrl+c", "esc":
				// Cancel input
				m.currentMode = modeView
				m.editingTask = nil
				m.textarea.SetValue("")
				return m, nil
			case "enter", "ctrl+s":
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
					return m, m.saveAndRefresh()
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			// Add new task
			m.currentMode = modeAdd
			m.textarea.SetValue("")
			m.textarea.Focus()
			return m, nil
		case "e":
			// Edit selected task
			if m.activePane == paneTask {
				if m.taskList.selected < len(m.filteredTasks) {
					m.currentMode = modeEdit
					// Store the reference to the original task, not the filtered one
					selectedTask := m.filteredTasks[m.taskList.selected]
					// Find the original task in m.tasks
					for i := 0; i < len(m.tasks); i++ {
						if m.tasks[i].String() == selectedTask.String() {
							m.editingTask = &m.tasks[i] // Point to the original task
							break
						}
					}
					if m.editingTask != nil {
						m.textarea.SetValue(m.editingTask.String())
						m.textarea.Focus()
					}
					return m, nil
				}
			}
		case "tab":
			// Switch between panes
			if m.activePane == paneFilter {
				m.activePane = paneTask
			} else {
				m.activePane = paneFilter
			}
			return m, nil
		case "h":
			// Move to left pane (filter)
			m.activePane = paneFilter
			return m, nil
		case "l":
			// Move to right pane (task)
			m.activePane = paneTask
			return m, nil
		case "enter":
			if m.activePane == paneFilter {
				// Filter selection changed, refresh task list
				m.refreshTaskList()
				// Move to right pane (task)
				m.activePane = paneTask
				return m, nil
			} else {
				// Complete task
				if m.taskList.selected < len(m.filteredTasks) {
					taskToComplete := m.filteredTasks[m.taskList.selected]
					// Find the task in main tasks list and mark as completed
					for i := 0; i < len(m.tasks); i++ {
						if m.tasks[i].String() == taskToComplete.String() {
							m.tasks[i].Complete()
							break
						}
					}
					return m, m.saveAndRefresh()
				}
			}
		case "d":
			if m.activePane == paneTask {
				// Delete task (only for non-deleted tasks)
				if m.taskList.selected < len(m.filteredTasks) {
					// Check if current filter is "Deleted Tasks"
					if m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name != "Deleted Tasks" {
						m.currentMode = modeDeleteConfirm
						m.deleteIndex = m.taskList.selected
						return m, nil
					}
				}
			}
		case "r":
			if m.activePane == paneTask {
				// Restore deleted or completed task
				if m.taskList.selected < len(m.filteredTasks) {
					currentFilter := ""
					if m.filterList.selected < len(m.filters) {
						currentFilter = m.filters[m.filterList.selected].name
					}
					
					taskToRestore := m.filteredTasks[m.taskList.selected]
					
					// Handle deleted tasks restoration
					if currentFilter == "Deleted Tasks" {
						// Find the task in main tasks list and remove deleted_at field
						for i := 0; i < len(m.tasks); i++ {
							if m.tasks[i].String() == taskToRestore.String() {
								// Remove deleted_at field to restore the task
								taskString := m.tasks[i].String()
								
								// Remove deleted_at field from the task string
								if strings.Contains(taskString, "deleted_at:") {
									// Simple approach: split and rejoin without deleted_at parts
									parts := strings.Fields(taskString)
									var cleanParts []string
									for _, part := range parts {
										if !strings.HasPrefix(part, "deleted_at:") {
											cleanParts = append(cleanParts, part)
										}
									}
									taskString = strings.Join(cleanParts, " ")
									
									// Parse the modified task string back to update the task
									if newTask, err := todotxt.ParseTask(taskString); err == nil {
										m.tasks[i] = *newTask
									}
								}
								break
							}
						}
						return m, m.saveAndRefresh()
					}
					
					// Handle completed tasks restoration
					if currentFilter == "Completed Tasks" {
						// Find the task in main tasks list and mark as incomplete
						for i := 0; i < len(m.tasks); i++ {
							if m.tasks[i].String() == taskToRestore.String() {
								m.tasks[i].Completed = false
								m.tasks[i].CompletedDate = time.Time{} // Clear completion date
								break
							}
						}
						return m, m.saveAndRefresh()
					}
				}
			}
		case "j", "down":
			if m.activePane == paneFilter {
				m.filterList.MoveDown()
				m.refreshTaskList()
			} else {
				m.taskList.MoveDown()
			}
			return m, nil
		case "k", "up":
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
		m.updatePaneSizes() // This is in view.go
		// Also update textarea size
		if m.width > 10 {
			m.textarea.SetWidth(m.width - 10)
		}
	case FileChangedMsg:
		// Reload tasks from file
		if taskList, err := todo.Load(m.todoFile); err == nil {
			m.tasks = taskList
			m.refreshLists()
		}
		// Continue watching
		return m, m.watchFile()
	}

	return m, nil
}

// Cleanup closes the file watcher
func (m *Model) Cleanup() {
	if m.watcher != nil {
		m.watcher.Close()
	}
} 