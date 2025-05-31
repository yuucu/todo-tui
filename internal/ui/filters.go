package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/charmbracelet/lipgloss"
)

// よく使用される文字列定数
const (
	enterKeyStr = "enter"
)

// isTaskDeleted checks if a task has the deleted_at field
func isTaskDeleted(task todotxt.Task) bool {
	taskString := task.String()
	return strings.Contains(taskString, TaskFieldDeletedPrefix)
}

// refreshLists updates both filter and task lists
func (m *Model) refreshLists() {
	m.refreshFilterList()
	m.refreshTaskList()
}

// refreshFilterList builds the filter list with projects and due dates
func (m *Model) refreshFilterList() {
	// Remember currently selected filter name for restoration
	var currentFilterName string
	if m.filterList.selected < len(m.filters) {
		currentFilterName = m.filters[m.filterList.selected].name
	}

	var filters []FilterData

	// Add time-based filters
	filters = append(filters, m.getTimeBasedFilters()...)

	// Always add "All Tasks" filter
	allTasksFilter := FilterData{
		name: FilterAllTasks,
		filterFn: func(tasks todotxt.TaskList) todotxt.TaskList {
			var result todotxt.TaskList
			for _, task := range tasks {
				if !task.Completed && !isTaskDeleted(task) {
					result = append(result, task)
				}
			}
			return result
		},
	}
	filters = append(filters, allTasksFilter)

	// Add project filters if any exist
	projects := m.getUniqueProjects()
	if len(projects) > 0 {
		filters = append(filters, FilterData{
			name: FilterHeaderProjects,
			filterFn: func(tasks todotxt.TaskList) todotxt.TaskList {
				return todotxt.TaskList{}
			},
		})
		for _, project := range projects {
			filters = append(filters, FilterData{
				name: "  +" + project,
				filterFn: func(p string) func(todotxt.TaskList) todotxt.TaskList {
					return func(tasks todotxt.TaskList) todotxt.TaskList {
						var result todotxt.TaskList
						for _, task := range tasks {
							if !task.Completed && !isTaskDeleted(task) {
								for _, taskProject := range task.Projects {
									if taskProject == p {
										result = append(result, task)
										break
									}
								}
							}
						}
						return result
					}
				}(project),
			})
		}
	}

	// Add context filters if any exist
	contexts := m.getUniqueContexts()
	if len(contexts) > 0 {
		filters = append(filters, FilterData{
			name: FilterHeaderContexts,
			filterFn: func(tasks todotxt.TaskList) todotxt.TaskList {
				return todotxt.TaskList{}
			},
		})
		for _, context := range contexts {
			filters = append(filters, FilterData{
				name: "  @" + context,
				filterFn: func(c string) func(todotxt.TaskList) todotxt.TaskList {
					return func(tasks todotxt.TaskList) todotxt.TaskList {
						var result todotxt.TaskList
						for _, task := range tasks {
							if !task.Completed && !isTaskDeleted(task) {
								for _, taskContext := range task.Contexts {
									if taskContext == c {
										result = append(result, task)
										break
									}
								}
							}
						}
						return result
					}
				}(context),
			})
		}
	}

	filters = append(filters, FilterData{
		name: FilterCompletedTasks,
		filterFn: func(tasks todotxt.TaskList) todotxt.TaskList {
			var result todotxt.TaskList
			for _, task := range tasks {
				if task.Completed && !isTaskDeleted(task) {
					result = append(result, task)
				}
			}
			return result
		},
	})

	// Deleted Tasks filter
	deletedTasks := m.tasks.Filter(isTaskDeleted)
	if len(deletedTasks) > 0 {
		filters = append(filters, FilterData{
			name: FilterDeletedTasks,
			filterFn: func(list todotxt.TaskList) todotxt.TaskList {
				return list.Filter(isTaskDeleted)
			},
			count: len(deletedTasks),
		})
	}

	// Calculate counts and create display items
	var items []string
	for i := range filters {
		// Count is already calculated for time-based filters, calculate for others
		if filters[i].count == 0 {
			filtered := filters[i].filterFn(m.tasks)
			filters[i].count = len(filtered)
		}

		if strings.Contains(filters[i].name, "─") {
			// Header
			items = append(items, filters[i].name)
		} else {
			// Project/Context and Other filters
			items = append(items, fmt.Sprintf("%s (%d)", filters[i].name, filters[i].count))
		}
	}

	// Find the new index for the previously selected filter
	newSelectedIndex := 0
	foundPreviousFilter := false

	// First, try to find the exact same filter name
	for i, filter := range filters {
		if filter.name == currentFilterName {
			newSelectedIndex = i
			foundPreviousFilter = true
			break
		}
	}

	// If the previous filter was not found (e.g., "Due Today" was removed),
	// switch to "All Tasks" filter
	if !foundPreviousFilter && currentFilterName != "" {
		// Check if it was a time-based filter that got removed
		timeBasedFilters := []string{"Due Today", "This Week", "Overdue"}
		wasTimeBasedFilter := false
		for _, tbf := range timeBasedFilters {
			if currentFilterName == tbf {
				wasTimeBasedFilter = true
				break
			}
		}

		if wasTimeBasedFilter {
			// Find "All Tasks" filter and select it
			for i, filter := range filters {
				if filter.name == "All Tasks" {
					newSelectedIndex = i
					break
				}
			}
		}
	}

	m.filters = filters
	m.filterList.SetItems(items)

	// Restore or update the selected filter position
	if newSelectedIndex < len(items) {
		m.filterList.selected = newSelectedIndex
	}
}

// refreshTaskList updates the task list based on current filter
func (m *Model) refreshTaskList() {
	var filteredTasks todotxt.TaskList

	if m.filterList.selected < len(m.filters) {
		filter := m.filters[m.filterList.selected]
		if !strings.Contains(filter.name, "─") { // Skip headers
			filteredTasks = filter.filterFn(m.tasks)
		}
	}

	// Only show default tasks if we're not viewing a specific filter that returned 0 results
	// Exception: Don't show default tasks for "Deleted Tasks" filter
	if len(filteredTasks) == 0 {
		// Check if current filter is "Deleted Tasks"
		isDeletedTasksFilter := m.filterList.selected < len(m.filters) &&
			m.filters[m.filterList.selected].name == FilterDeletedTasks

		if !isDeletedTasksFilter {
			// Default to all incomplete tasks (only for non-deleted task filters)
			for _, task := range m.tasks {
				if !task.Completed && !isTaskDeleted(task) {
					filteredTasks = append(filteredTasks, task)
				}
			}
		}
	}

	// Convert tasks to display strings
	var items []string
	for i := range filteredTasks {
		task := &filteredTasks[i]

		// Build task display string
		display := task.Todo

		// Add priority with color
		if task.HasPriority() {
			priorityStyle := lipgloss.NewStyle().Bold(true)
			switch task.Priority {
			case "A":
				priorityStyle = priorityStyle.Foreground(m.currentTheme.PriorityHigh)
			case "B":
				priorityStyle = priorityStyle.Foreground(m.currentTheme.PriorityMedium)
			case "C":
				priorityStyle = priorityStyle.Foreground(m.currentTheme.PriorityLow)
			case "D":
				priorityStyle = priorityStyle.Foreground(m.currentTheme.PriorityLowest)
			default:
				priorityStyle = priorityStyle.Foreground(m.currentTheme.PriorityDefault)
			}
			display = priorityStyle.Render(fmt.Sprintf("(%s) ", task.Priority)) + display
		}

		// Add contexts and projects
		var tags []string
		for _, project := range task.Projects {
			tags = append(tags, lipgloss.NewStyle().
				Foreground(m.currentTheme.Secondary).
				Render("+"+project))
		}
		for _, context := range task.Contexts {
			tags = append(tags, lipgloss.NewStyle().
				Foreground(m.currentTheme.Primary).
				Render("@"+context))
		}

		// Add due date
		if task.HasDueDate() {
			dueStyle := lipgloss.NewStyle()
			dueDate := task.DueDate.Format(DateFormat)
			now := time.Now()
			today := now.Format(DateFormat)

			if dueDate < today {
				dueStyle = dueStyle.Foreground(m.currentTheme.Danger) // Overdue
			} else if dueDate == today {
				dueStyle = dueStyle.Foreground(m.currentTheme.Warning) // Due today
			} else {
				dueStyle = dueStyle.Foreground(m.currentTheme.Success) // Future
			}

			tags = append(tags, dueStyle.Render(TaskFieldDuePrefix+dueDate))
		}

		if len(tags) > 0 {
			display += " " + strings.Join(tags, " ")
		}

		items = append(items, display)
	}

	m.filteredTasks = filteredTasks
	m.taskList.SetItems(items)
}

// getUniqueProjects returns sorted unique project names
func (m *Model) getUniqueProjects() []string {
	projectMap := make(map[string]bool)
	for _, task := range m.tasks {
		if !task.Completed && !isTaskDeleted(task) {
			for _, project := range task.Projects {
				projectMap[project] = true
			}
		}
	}

	var projects []string
	for project := range projectMap {
		projects = append(projects, project)
	}
	sort.Strings(projects)
	return projects
}

// getUniqueContexts returns sorted unique context names
func (m *Model) getUniqueContexts() []string {
	contextMap := make(map[string]bool)
	for _, task := range m.tasks {
		if !task.Completed && !isTaskDeleted(task) {
			for _, context := range task.Contexts {
				contextMap[context] = true
			}
		}
	}

	var contexts []string
	for context := range contextMap {
		contexts = append(contexts, context)
	}
	sort.Strings(contexts)
	return contexts
}

// getStatusInfo returns status information for display
func (m *Model) getStatusInfo() string {
	// If there's an active status message, show it instead with appropriate color
	if m.statusMessage != "" && time.Now().Before(m.statusMessageEnd) {
		// Apply color based on message type
		messageStyle := lipgloss.NewStyle()
		if strings.Contains(m.statusMessage, "📋") || strings.Contains(m.statusMessage, "✅") {
			// Success message - green
			messageStyle = messageStyle.Foreground(m.currentTheme.Success)
		} else if strings.Contains(m.statusMessage, "❌") {
			// Error message - red
			messageStyle = messageStyle.Foreground(m.currentTheme.Danger)
		} else {
			// Default - normal text color
			messageStyle = messageStyle.Foreground(m.currentTheme.Text)
		}
		return messageStyle.Render(m.statusMessage)
	}

	// Current time
	now := time.Now().Format(TimeFormat)

	// Current filter name
	var currentFilter string
	if m.filterList.selected < len(m.filters) {
		filterName := m.filters[m.filterList.selected].name
		// Clean up filter name for display
		if strings.HasPrefix(filterName, "  +") || strings.HasPrefix(filterName, "  @") {
			currentFilter = strings.TrimSpace(filterName)
		} else if !strings.Contains(filterName, "─") {
			currentFilter = strings.Split(filterName, " (")[0]
		} else {
			currentFilter = "All"
		}
	} else {
		currentFilter = "All"
	}

	// Task counts
	totalTasks := 0
	for _, task := range m.tasks {
		if !task.Completed {
			totalTasks++
		}
	}

	filteredCount := len(m.filteredTasks)

	// Icons and info
	return fmt.Sprintf("🏷️  %s │ 📋 %d/%d │ 🕐 %s",
		currentFilter, filteredCount, totalTasks, now)
}

// Helper function to create filters with count check
func (m *Model) addFilterIfNotEmpty(name string, filterFn func(todotxt.TaskList) todotxt.TaskList) *FilterData {
	filtered := filterFn(m.tasks)
	if len(filtered) > 0 {
		return &FilterData{
			name:     name,
			filterFn: filterFn,
			count:    len(filtered),
		}
	}
	return nil
}
