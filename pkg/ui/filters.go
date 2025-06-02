package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"github.com/yuucu/todotui/pkg/domain"
)

// getCompletedTaskTransitionConfig converts UI config to domain config
func (m *Model) getCompletedTaskTransitionConfig() domain.CompletedTaskTransitionConfig {
	return domain.CompletedTaskTransitionConfig{
		DelayDays:      m.appConfig.UI.CompletedTaskTransition.DelayDays,
		TransitionHour: m.appConfig.UI.CompletedTaskTransition.TransitionHour,
	}
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
		filterFn: func(tasks domain.Tasks) domain.Tasks {
			return tasks.Filter(func(task domain.Task, _ int) bool {
				if task.IsDeleted() {
					return false
				}
				// Show incomplete tasks OR completed tasks that haven't moved to "Completed Tasks" yet
				if !task.IsCompleted() {
					return true
				}
				// For completed tasks, check if they should move
				config := m.getCompletedTaskTransitionConfig()
				return !task.ShouldMoveToCompleted(config, time.Now())
			})
		},
	}
	filters = append(filters, allTasksFilter)

	// Add "No Project" filter for tasks without any project tags
	noProjectFilter := FilterData{
		name: FilterNoProject,
		filterFn: func(tasks domain.Tasks) domain.Tasks {
			return tasks.Filter(func(task domain.Task, _ int) bool {
				if task.IsDeleted() || len(task.Projects()) > 0 {
					return false
				}
				// Show incomplete tasks OR completed tasks that haven't moved to "Completed Tasks" yet
				if !task.IsCompleted() {
					return true
				}
				config := m.getCompletedTaskTransitionConfig()
				return !task.ShouldMoveToCompleted(config, time.Now())
			})
		},
	}
	filters = append(filters, noProjectFilter)

	// Add project filters if any exist
	projects := m.getUniqueProjects()
	if len(projects) > 0 {
		filters = append(filters, FilterData{
			name: FilterHeaderProjects,
			filterFn: func(tasks domain.Tasks) domain.Tasks {
				return domain.Tasks{}
			},
		})
		for _, project := range projects {
			filters = append(filters, FilterData{
				name: "  +" + project,
				filterFn: func(p string) func(domain.Tasks) domain.Tasks {
					return func(tasks domain.Tasks) domain.Tasks {
						return tasks.Filter(func(task domain.Task, _ int) bool {
							if task.IsDeleted() || !lo.Contains(task.Projects(), p) {
								return false
							}
							// Show incomplete tasks OR completed tasks that haven't moved to "Completed Tasks" yet
							if !task.IsCompleted() {
								return true
							}
							config := m.getCompletedTaskTransitionConfig()
							return !task.ShouldMoveToCompleted(config, time.Now())
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
			filterFn: func(tasks domain.Tasks) domain.Tasks {
				return domain.Tasks{}
			},
		})
		for _, context := range contexts {
			filters = append(filters, FilterData{
				name: "  @" + context,
				filterFn: func(c string) func(domain.Tasks) domain.Tasks {
					return func(tasks domain.Tasks) domain.Tasks {
						return tasks.Filter(func(task domain.Task, _ int) bool {
							if task.IsDeleted() || !lo.Contains(task.Contexts(), c) {
								return false
							}
							// Show incomplete tasks OR completed tasks that haven't moved to "Completed Tasks" yet
							if !task.IsCompleted() {
								return true
							}
							config := m.getCompletedTaskTransitionConfig()
							return !task.ShouldMoveToCompleted(config, time.Now())
						})
					}
				}(context),
			})
		}
	}

	filters = append(filters, FilterData{
		name: FilterCompletedTasks,
		filterFn: func(tasks domain.Tasks) domain.Tasks {
			return tasks.Filter(func(task domain.Task, _ int) bool {
				// Only show completed tasks that should be moved to the "Completed Tasks" filter
				// based on the transition configuration
				config := m.getCompletedTaskTransitionConfig()
				return task.IsCompleted() && !task.IsDeleted() &&
					task.ShouldMoveToCompleted(config, time.Now())
			})
		},
	})

	// Deleted Tasks filter
	deletedTasks := m.tasks.FilterDeleted()
	if deletedTasks.Len() > 0 {
		filters = append(filters, FilterData{
			name: FilterDeletedTasks,
			filterFn: func(tasks domain.Tasks) domain.Tasks {
				return tasks.FilterDeleted()
			},
			count: deletedTasks.Len(),
		})
	}

	// Calculate counts and create display items using lo.Map
	items := lo.Map(filters, func(filter FilterData, i int) string {
		// Count is already calculated for time-based filters, calculate for others
		if filters[i].count == 0 {
			filtered := filter.filterFn(m.tasks)
			filters[i].count = filtered.Len()
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

	// Restore or update the selected filter position after SetItems
	// Use preserve scroll method to avoid unwanted scrolling
	if newSelectedIndex < len(items) {
		m.filterList.SetSelectedIndexPreserveScroll(newSelectedIndex)
	}
}

// refreshTaskList updates the task list based on current filter
func (m *Model) refreshTaskList() {
	var filteredTasks domain.Tasks

	if m.filterList.selected < len(m.filters) {
		filter := m.filters[m.filterList.selected]
		if !strings.Contains(filter.name, "─") { // Skip headers
			filteredTasks = filter.filterFn(m.tasks)
		}
	}

	// Only show default tasks if we're not viewing a specific filter that returned 0 results
	// Exception: Don't show default tasks for "Deleted Tasks" and "No Project" filters
	if filteredTasks.Len() == 0 {
		// Check if current filter is "Deleted Tasks" or "No Project"
		isDeletedTasksFilter := m.filterList.selected < len(m.filters) &&
			m.filters[m.filterList.selected].name == FilterDeletedTasks
		isNoProjectFilter := m.filterList.selected < len(m.filters) &&
			m.filters[m.filterList.selected].name == FilterNoProject

		if !isDeletedTasksFilter && !isNoProjectFilter {
			// Default to all incomplete tasks AND completed tasks that haven't moved yet
			// (only for non-deleted and non-no-project task filters) using Filter
			filteredTasks = m.tasks.Filter(func(task domain.Task, _ int) bool {
				if task.IsDeleted() {
					return false
				}
				// Show incomplete tasks OR completed tasks that haven't moved to "Completed Tasks" yet
				if !task.IsCompleted() {
					return true
				}
				// For completed tasks, check if they should move
				config := m.getCompletedTaskTransitionConfig()
				return !task.ShouldMoveToCompleted(config, time.Now())
			})
		}
	}

	// ソート処理をdomainの新しいTasks型に委譲
	sortedTasks := filteredTasks.SortByCompletionStatus()
	filteredTasks = sortedTasks

	// Convert tasks to display strings and track completion status
	var items []string
	var completedItems []bool
	var checkboxColors []lipgloss.Color

	taskList := filteredTasks.ToTaskList()
	for i := range taskList {
		task := &taskList[i]

		// Track completion status for the enhanced list display
		// Consider both completed tasks and deleted tasks as "completed" for UI purposes
		var isTaskCompleted bool
		if domainTask, err := domain.NewTask(task); err == nil {
			isTaskCompleted = task.Completed || domainTask.IsDeleted()
		} else {
			isTaskCompleted = task.Completed // Fallback if domain task creation fails
		}
		completedItems = append(completedItems, isTaskCompleted)

		// Calculate checkbox color based on due date for incomplete tasks
		var checkboxColor lipgloss.Color
		if !isTaskCompleted && task.HasDueDate() {
			now := time.Now()
			if domainTask, err := domain.NewTask(task); err == nil {
				if domainTask.IsOverdue(now) {
					checkboxColor = m.currentTheme.Danger // Overdue - red
				} else if domainTask.IsDueToday(now) {
					checkboxColor = m.currentTheme.Warning // Due today - yellow
				} else {
					checkboxColor = m.currentTheme.Success // Future - green
				}
			} else {
				checkboxColor = m.currentTheme.Success // Default to green if domain task creation fails
			}
		} else {
			// For completed tasks or tasks without due date, use default muted color
			checkboxColor = m.currentTheme.TextMuted
		}
		checkboxColors = append(checkboxColors, checkboxColor)

		// Build task display string - only plain text, styling will be done in renderTaskItem
		display := task.Todo

		// For both completed and active tasks, keep plain text and let renderTaskItem handle all styling
		if task.HasPriority() {
			display = fmt.Sprintf("(%s) ", task.Priority) + display
		}

		var tags []string
		for _, project := range task.Projects {
			tags = append(tags, "+"+project)
		}
		for _, context := range task.Contexts {
			tags = append(tags, "@"+context)
		}
		if task.HasDueDate() {
			dueDate := task.DueDate.Format(DateFormat)
			tags = append(tags, TaskFieldDuePrefix+dueDate)
		}

		if len(tags) > 0 {
			display += " " + strings.Join(tags, " ")
		}

		items = append(items, display)
	}

	m.filteredTasks = filteredTasks
	m.taskList.SetItems(items)
	// Set completion status for enhanced display
	m.taskList.SetCompletedItems(completedItems)
	// Set checkbox colors to match due date colors
	m.taskList.SetCheckboxColors(checkboxColors)
}

// getUniqueProjects returns sorted unique project names
func (m *Model) getUniqueProjects() []string {
	// Get all incomplete, non-deleted tasks
	activeTasks := m.tasks.FilterActive()

	// Extract all projects from active tasks
	var allProjects []string
	for _, task := range activeTasks {
		allProjects = append(allProjects, task.Projects()...)
	}

	// Get unique projects and sort
	uniqueProjects := lo.Uniq(allProjects)
	sort.Strings(uniqueProjects)
	return uniqueProjects
}

// getUniqueContexts returns sorted unique context names
func (m *Model) getUniqueContexts() []string {
	// Get all incomplete, non-deleted tasks
	activeTasks := m.tasks.FilterActive()

	// Extract all contexts from active tasks
	var allContexts []string
	for _, task := range activeTasks {
		allContexts = append(allContexts, task.Contexts()...)
	}

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

	// Task counts
	totalTasks := m.tasks.FilterActive().Len()
	filteredCount := m.filteredTasks.Len()

	// Icons and info
	return fmt.Sprintf("🏷️  %s │ 📋 %d/%d │ 🕐 %s",
		currentFilter, filteredCount, totalTasks, now)
}

// Helper function to create filters with count check
func (m *Model) addFilterIfNotEmpty(name string, filterFn func(domain.Tasks) domain.Tasks) *FilterData {
	filtered := filterFn(m.tasks)
	if filtered.Len() > 0 {
		return &FilterData{
			name:     name,
			filterFn: filterFn,
			count:    filtered.Len(),
		}
	}
	return nil
}

func FilterActiveTask(tasks domain.Tasks) domain.Tasks {
	return tasks.Filter(func(task domain.Task, _ int) bool {
		return !task.IsCompleted() && !task.IsDeleted()
	})
}

// SortTasksByCompletionStatus sorts tasks by completion status
func SortTasksByCompletionStatus(tasks todotxt.TaskList, currentFilter string) todotxt.TaskList {
	// Convert each task only once and handle errors
	domainTasks := make([]struct {
		original    todotxt.Task
		domain      *domain.Task
		isCompleted bool
	}, len(tasks))

	for i, task := range tasks {
		domainTasks[i].original = task
		if domainTask, err := domain.NewTask(&task); err == nil {
			domainTasks[i].domain = domainTask
			domainTasks[i].isCompleted = task.Completed || domainTask.IsDeleted()
		} else {
			// Skip invalid tasks or treat as incomplete
			domainTasks[i].domain = nil
			domainTasks[i].isCompleted = task.Completed
		}
	}

	// ソート: 完了したタスクを下の方に移動
	sort.SliceStable(domainTasks, func(i, j int) bool {
		// 完了状態が異なる場合：未完了を上に、完了を下に
		if domainTasks[i].isCompleted != domainTasks[j].isCompleted {
			return !domainTasks[i].isCompleted // 未完了（false）が上に来る
		}
		// 両方とも同じ完了状態の場合：元の順序を保持（安定ソート）
		return false
	})

	// Extract original tasks in sorted order
	result := make(todotxt.TaskList, len(tasks))
	for i, taskInfo := range domainTasks {
		result[i] = taskInfo.original
	}

	return result
}
