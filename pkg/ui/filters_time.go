package ui

import (
	"time"

	"github.com/yuucu/todotui/pkg/domain"
)

// getTimeBasedFilters returns all time-based filters
func (m *Model) getTimeBasedFilters() []FilterData {
	var filters []FilterData

	// Add time-based filters only if they have tasks
	if filter := m.addFilterIfNotEmpty("Due Today", m.getDueTodayFilterFn()); filter != nil {
		filters = append(filters, *filter)
	}

	if filter := m.addFilterIfNotEmpty("This Week", m.getThisWeekFilterFn()); filter != nil {
		filters = append(filters, *filter)
	}

	if filter := m.addFilterIfNotEmpty("Overdue", m.getOverdueFilterFn()); filter != nil {
		filters = append(filters, *filter)
	}

	return filters
}

// getDueTodayFilterFn returns the filter function for due today tasks
func (m *Model) getDueTodayFilterFn() func(domain.Tasks) domain.Tasks {
	return func(tasks domain.Tasks) domain.Tasks {
		now := time.Now()
		config := m.getCompletedTaskTransitionConfig()
		return tasks.Filter(func(task domain.Task, _ int) bool {
			// Skip deleted tasks
			if task.IsDeleted() {
				return false
			}

			// Only include tasks that have due dates and are due today
			if !task.HasDueDate() {
				return false
			}

			if task.IsDueToday(now) {
				// For incomplete tasks, always include
				if !task.IsCompleted() {
					return true
				}
				// For completed tasks, include if they haven't been removed from original filters yet
				return !task.ShouldRemoveFromOriginalFilters(config, now)
			}
			return false
		})
	}
}

// getThisWeekFilterFn returns the filter function for this week tasks
func (m *Model) getThisWeekFilterFn() func(domain.Tasks) domain.Tasks {
	return func(tasks domain.Tasks) domain.Tasks {
		now := time.Now()
		config := m.getCompletedTaskTransitionConfig()
		return tasks.Filter(func(task domain.Task, _ int) bool {
			// Skip deleted tasks
			if task.IsDeleted() {
				return false
			}

			// Only include tasks that have due dates and are due this week
			if !task.HasDueDate() {
				return false
			}

			if task.IsThisWeek(now) {
				// For incomplete tasks, always include
				if !task.IsCompleted() {
					return true
				}
				// For completed tasks, include if they haven't been removed from original filters yet
				return !task.ShouldRemoveFromOriginalFilters(config, now)
			}
			return false
		})
	}
}

// getOverdueFilterFn returns the filter function for overdue tasks
func (m *Model) getOverdueFilterFn() func(domain.Tasks) domain.Tasks {
	return func(tasks domain.Tasks) domain.Tasks {
		now := time.Now()
		config := m.getCompletedTaskTransitionConfig()
		return tasks.Filter(func(task domain.Task, _ int) bool {
			// Skip deleted tasks
			if task.IsDeleted() {
				return false
			}

			// Only include tasks that have due dates and are overdue
			if !task.HasDueDate() {
				return false
			}

			if task.IsOverdue(now) {
				// For incomplete tasks, always include
				if !task.IsCompleted() {
					return true
				}
				// For completed tasks, include if they haven't been removed from original filters yet
				return !task.ShouldRemoveFromOriginalFilters(config, now)
			}
			return false
		})
	}
}
