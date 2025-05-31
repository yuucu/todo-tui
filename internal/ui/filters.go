package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
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
			return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
				return !task.Completed && !isTaskDeleted(task)
			})
		},
	}
	filters = append(filters, allTasksFilter)

	// Add "No Project" filter for tasks without any project tags
	noProjectFilter := FilterData{
		name: FilterNoProject,
		filterFn: func(tasks todotxt.TaskList) todotxt.TaskList {
			return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
				return !task.Completed && !isTaskDeleted(task) && len(task.Projects) == 0
			})
		},
	}
	filters = append(filters, noProjectFilter)

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
						return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
							return !task.Completed && !isTaskDeleted(task) &&
								lo.Contains(task.Projects, p)
						})
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
						return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
							return !task.Completed && !isTaskDeleted(task) &&
								lo.Contains(task.Contexts, c)
						})
					}
				}(context),
			})
		}
	}

	filters = append(filters, FilterData{
		name: FilterCompletedTasks,
		filterFn: func(tasks todotxt.TaskList) todotxt.TaskList {
			return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
				return task.Completed && !isTaskDeleted(task)
			})
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

	// Calculate counts and create display items using lo.Map
	items := lo.Map(filters, func(filter FilterData, i int) string {
		// Count is already calculated for time-based filters, calculate for others
		if filters[i].count == 0 {
			filtered := filter.filterFn(m.tasks)
			filters[i].count = len(filtered)
		}

		if strings.Contains(filter.name, "─") {
			// Header
			return filter.name
		}
		// Project/Context and Other filters
		return fmt.Sprintf("%s (%d)", filter.name, filters[i].count)
	})

	// Find the new index for the previously selected filter using lo.FindIndexOf
	newSelectedIndex := 0
	foundPreviousFilter := false

	_, index, found := lo.FindIndexOf(filters, func(filter FilterData) bool {
		return filter.name == currentFilterName
	})
	if found {
		newSelectedIndex = index
		foundPreviousFilter = true
	}

	// If the previous filter was not found (e.g., "Due Today" was removed),
	// switch to "All Tasks" filter
	if !foundPreviousFilter && currentFilterName != "" {
		// Check if it was a time-based filter that got removed
		timeBasedFilters := []string{"Due Today", "This Week", "Overdue"}
		wasTimeBasedFilter := lo.Contains(timeBasedFilters, currentFilterName)

		if wasTimeBasedFilter {
			// Find "All Tasks" filter and select it
			_, allTasksIndex, found := lo.FindIndexOf(filters, func(filter FilterData) bool {
				return filter.name == "All Tasks"
			})
			if found {
				newSelectedIndex = allTasksIndex
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
			// Default to all incomplete tasks (only for non-deleted task filters) using lo.Filter
			filteredTasks = lo.Filter(m.tasks, func(task todotxt.Task, _ int) bool {
				return !task.Completed && !isTaskDeleted(task)
			})
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
	// Get all incomplete, non-deleted tasks
	activeTasks := lo.Filter(m.tasks, func(task todotxt.Task, _ int) bool {
		return !task.Completed && !isTaskDeleted(task)
	})

	// Extract all projects from active tasks
	allProjects := lo.FlatMap(activeTasks, func(task todotxt.Task, _ int) []string {
		return task.Projects
	})

	// Get unique projects and sort
	uniqueProjects := lo.Uniq(allProjects)
	sort.Strings(uniqueProjects)
	return uniqueProjects
}

// getUniqueContexts returns sorted unique context names
func (m *Model) getUniqueContexts() []string {
	// Get all incomplete, non-deleted tasks
	activeTasks := lo.Filter(m.tasks, func(task todotxt.Task, _ int) bool {
		return !task.Completed && !isTaskDeleted(task)
	})

	// Extract all contexts from active tasks
	allContexts := lo.FlatMap(activeTasks, func(task todotxt.Task, _ int) []string {
		return task.Contexts
	})

	// Get unique contexts and sort
	uniqueContexts := lo.Uniq(allContexts)
	sort.Strings(uniqueContexts)
	return uniqueContexts
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

	// Task counts using lo.CountBy for more functional approach
	totalTasks := lo.CountBy(m.tasks, func(task todotxt.Task) bool {
		return !task.Completed
	})

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
