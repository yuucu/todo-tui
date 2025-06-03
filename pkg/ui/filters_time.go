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
		return tasks.Filter(func(task domain.Task, _ int) bool {
			// Skip deleted and completed tasks
			if task.IsDeleted() || task.IsCompleted() {
				return false
			}

			// Only include tasks that have due dates and are due today
			if !task.HasDueDate() {
				return false
			}

			// For Due Today filter, we still want to show overdue tasks from today
			// Only exclude tasks that are from previous days
			if task.IsOverdue(now) && !task.IsDueToday(now) {
				return false
			}

			return task.IsDueToday(now)
		})
	}
}

// getThisWeekFilterFn returns the filter function for this week tasks
func (m *Model) getThisWeekFilterFn() func(domain.Tasks) domain.Tasks {
	return func(tasks domain.Tasks) domain.Tasks {
		now := time.Now()
		return tasks.Filter(func(task domain.Task, _ int) bool {
			// Skip deleted and completed tasks
			if task.IsDeleted() || task.IsCompleted() {
				return false
			}

			// Only include tasks that have due dates and are due this week
			if !task.HasDueDate() {
				return false
			}

			// Skip overdue tasks - they should not appear in "This Week" filter
			if task.IsOverdue(now) {
				return false
			}

			return task.IsThisWeek(now)
		})
	}
}

// getOverdueFilterFn returns the filter function for overdue tasks
func (m *Model) getOverdueFilterFn() func(domain.Tasks) domain.Tasks {
	return func(tasks domain.Tasks) domain.Tasks {
		now := time.Now()
		return tasks.Filter(func(task domain.Task, _ int) bool {
			// Skip deleted and completed tasks
			if task.IsDeleted() || task.IsCompleted() {
				return false
			}

			// Only include tasks that have due dates and are overdue
			if !task.HasDueDate() {
				return false
			}

			return task.IsOverdue(now)
		})
	}
}
