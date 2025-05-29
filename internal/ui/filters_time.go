package ui

import (
	"time"

	todotxt "github.com/1set/todotxt"
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
func (m *Model) getDueTodayFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		today := time.Now().Format("2006-01-02")
		var result todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed && !isTaskDeleted(task) && task.HasDueDate() && task.DueDate.Format("2006-01-02") == today {
				result = append(result, task)
			}
		}
		return result
	}
}

// getThisWeekFilterFn returns the filter function for this week tasks
func (m *Model) getThisWeekFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		weekStart := now.AddDate(0, 0, -int(now.Weekday()))
		weekEnd := weekStart.AddDate(0, 0, 7)
		var result todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed && !isTaskDeleted(task) && task.HasDueDate() &&
				!task.DueDate.Before(weekStart) && task.DueDate.Before(weekEnd) {
				result = append(result, task)
			}
		}
		return result
	}
}

// getOverdueFilterFn returns the filter function for overdue tasks
func (m *Model) getOverdueFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		var result todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed && !isTaskDeleted(task) && task.HasDueDate() && task.DueDate.Before(now) {
				result = append(result, task)
			}
		}
		return result
	}
}
