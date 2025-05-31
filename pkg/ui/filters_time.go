package ui

import (
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/samber/lo"
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
func (m *Model) getDueTodayFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
			domainTask := domain.NewTask(&task)
			if !domainTask.IsDueToday(now) {
				return false
			}
			// For completed tasks, only show them if they haven't moved to "Completed Tasks" yet
			if task.Completed {
				config := m.getCompletedTaskTransitionConfig()
				return !domainTask.ShouldMoveToCompleted(config)
			}
			return true
		})
	}
}

// getThisWeekFilterFn returns the filter function for this week tasks
func (m *Model) getThisWeekFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
			domainTask := domain.NewTask(&task)
			if !domainTask.IsThisWeek(now) {
				return false
			}
			// For completed tasks, only show them if they haven't moved to "Completed Tasks" yet
			if task.Completed {
				config := m.getCompletedTaskTransitionConfig()
				return !domainTask.ShouldMoveToCompleted(config)
			}
			return true
		})
	}
}

// getOverdueFilterFn returns the filter function for overdue tasks
func (m *Model) getOverdueFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
			domainTask := domain.NewTask(&task)
			if !domainTask.IsOverdue(now) {
				return false
			}
			// For completed tasks, only show them if they haven't moved to "Completed Tasks" yet
			if task.Completed {
				config := m.getCompletedTaskTransitionConfig()
				return !domainTask.ShouldMoveToCompleted(config)
			}
			return true
		})
	}
}
