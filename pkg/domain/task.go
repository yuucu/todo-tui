package domain

import (
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
)

// Task represents the domain logic for a task
type Task struct {
	task *todotxt.Task
}

// NewTask creates a new Task domain instance
func NewTask(task *todotxt.Task) *Task {
	return &Task{
		task: task,
	}
}

// ToggleCompletion toggles the completion status of the task
// Returns true if the task is now completed, false if it's now incomplete
func (t *Task) ToggleCompletion() bool {
	if t.task.Completed {
		// Mark as incomplete
		t.task.Completed = false
		t.task.CompletedDate = time.Time{}
		return false
	} else {
		// Mark as completed
		t.task.Completed = true
		t.task.CompletedDate = time.Now()
		return true
	}
}

// IsCompleted returns whether the task is completed
func (t *Task) IsCompleted() bool {
	return t.task.Completed
}

// GetCompletedDate returns the completion date of the task
func (t *Task) GetCompletedDate() time.Time {
	return t.task.CompletedDate
}

// ShouldMoveToCompleted checks if a completed task should be moved to "Completed Tasks" filter
// based on the completion date and transition configuration
func (t *Task) ShouldMoveToCompleted(config CompletedTaskTransitionConfig, now time.Time) bool {
	if !t.task.Completed || t.task.CompletedDate.IsZero() {
		return false
	}

	// If delay is 0, move immediately
	if config.DelayDays == 0 {
		return true
	}

	completedDate := t.task.CompletedDate

	// Calculate the target transition date/time
	// Add the delay days to the completion date
	targetDate := completedDate.AddDate(0, 0, config.DelayDays)

	// Set the transition time to the specified hour on the target date
	targetDateTime := time.Date(
		targetDate.Year(),
		targetDate.Month(),
		targetDate.Day(),
		config.TransitionHour,
		0, 0, 0,
		now.Location(), // Use current timezone
	)

	// If the target time has passed, the task should be moved
	return now.After(targetDateTime) || now.Equal(targetDateTime)
}

// IsDeleted checks if the task has been soft deleted
func (t *Task) IsDeleted() bool {
	taskString := t.task.String()
	return containsDeletedPrefix(taskString)
}

// IsOverdue checks if a task is overdue based on the current date
func (t *Task) IsOverdue(now time.Time) bool {
	if !t.task.HasDueDate() || t.IsDeleted() {
		return false
	}

	today := now.Format("2006-01-02")
	taskDateStr := t.task.DueDate.Format("2006-01-02")
	return taskDateStr < today
}

// IsDueToday checks if a task is due today
func (t *Task) IsDueToday(now time.Time) bool {
	if !t.task.HasDueDate() || t.IsDeleted() {
		return false
	}

	today := now.Format("2006-01-02")
	taskDateStr := t.task.DueDate.Format("2006-01-02")
	return taskDateStr == today
}

// IsThisWeek checks if a task is due this week
func (t *Task) IsThisWeek(now time.Time) bool {
	if !t.task.HasDueDate() || t.IsDeleted() {
		return false
	}

	// 週の開始を日曜日として計算（Goでは日曜日が0）
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	// 週の終了を土曜日の翌日（日曜日）として計算
	weekEnd := weekStart.AddDate(0, 0, 7)

	// 日付レベルでの比較
	taskDate := t.task.DueDate.Format("2006-01-02")
	weekStartStr := weekStart.Format("2006-01-02")
	weekEndStr := weekEnd.Format("2006-01-02")

	// タスクの期限日が週の範囲内かチェック（開始日含む、終了日除く）
	return taskDate >= weekStartStr && taskDate < weekEndStr
}

// CompletedTaskTransitionConfig represents settings for when completed tasks move to "Completed Tasks" filter
type CompletedTaskTransitionConfig struct {
	// Number of days to wait before moving completed tasks to "Completed Tasks" filter
	DelayDays int

	// Time of day (24-hour format) when the transition should occur (0-23)
	TransitionHour int
}

// containsDeletedPrefix checks if a task string contains the deleted_at field
func containsDeletedPrefix(taskString string) bool {
	// Check if the task string contains the "deleted_at:" prefix
	return strings.Contains(taskString, "deleted_at:")
}
