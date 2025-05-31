package ui

import (
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/samber/lo"
)

// 時間関連の定数
const (
	daysInWeek = 7
)

// isOverdue checks if a task is overdue based on the current date
func isOverdue(task todotxt.Task, now time.Time) bool {
	if !task.HasDueDate() || task.Completed || isTaskDeleted(task) {
		return false
	}

	today := now.Format("2006-01-02")
	taskDateStr := task.DueDate.Format("2006-01-02")
	return taskDateStr < today
}

// isDueToday checks if a task is due today
func isDueToday(task todotxt.Task, now time.Time) bool {
	if !task.HasDueDate() || task.Completed || isTaskDeleted(task) {
		return false
	}

	today := now.Format("2006-01-02")
	taskDateStr := task.DueDate.Format("2006-01-02")
	return taskDateStr == today
}

// isThisWeek checks if a task is due this week
func isThisWeek(task todotxt.Task, now time.Time) bool {
	if !task.HasDueDate() || task.Completed || isTaskDeleted(task) {
		return false
	}

	// 週の開始を日曜日として計算（Goでは日曜日が0）
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	// 週の終了を土曜日の翌日（日曜日）として計算
	weekEnd := weekStart.AddDate(0, 0, daysInWeek)

	// 日付レベルでの比較
	taskDate := task.DueDate.Format("2006-01-02")
	weekStartStr := weekStart.Format("2006-01-02")
	weekEndStr := weekEnd.Format("2006-01-02")

	// タスクの期限日が週の範囲内かチェック（開始日含む、終了日除く）
	return taskDate >= weekStartStr && taskDate < weekEndStr
}

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
			return isDueToday(task, now)
		})
	}
}

// getThisWeekFilterFn returns the filter function for this week tasks
func (m *Model) getThisWeekFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
			return isThisWeek(task, now)
		})
	}
}

// getOverdueFilterFn returns the filter function for overdue tasks
func (m *Model) getOverdueFilterFn() func(todotxt.TaskList) todotxt.TaskList {
	return func(tasks todotxt.TaskList) todotxt.TaskList {
		now := time.Now()
		return lo.Filter(tasks, func(task todotxt.Task, _ int) bool {
			return isOverdue(task, now)
		})
	}
}
