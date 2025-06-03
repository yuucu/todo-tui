package domain

import (
	"errors"
	"strings"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/samber/lo"
)

// ===============================
// Task Field Constants
// ===============================

// Task field names
const (
	TaskFieldDeleted = "deleted_at"
	TaskFieldDue     = "due"
)

// Task field prefixes (with colon)
const (
	TaskFieldDeletedPrefix = TaskFieldDeleted + ":"
	TaskFieldDuePrefix     = TaskFieldDue + ":"
)

// ===============================
// Date Format Constants
// ===============================

const (
	// Go standard date format (ISO 8601)
	DateFormat = "2006-01-02"
)

// Task represents the domain logic for a task
type Task struct {
	task *todotxt.Task
}

// NewTask creates a new Task domain instance
func NewTask(task *todotxt.Task) (*Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}
	return &Task{
		task: task,
	}, nil
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

// ShouldRemoveFromOriginalFilters checks if a completed task should be removed from
// original filters (Due Today, All Tasks, etc.) based on the transition configuration
// This is separate from ShouldMoveToCompleted to allow completed tasks to appear
// in both Completed Tasks and original filters during the delay period
func (t *Task) ShouldRemoveFromOriginalFilters(config CompletedTaskTransitionConfig, now time.Time) bool {
	if !t.task.Completed || t.task.CompletedDate.IsZero() {
		return false
	}

	// If delay is 0, remove immediately
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

	// If the target time has passed, the task should be removed from original filters
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
// For completed tasks, this will return true if the task was originally due today
// and hasn't been removed from original filters yet
func (t *Task) IsDueToday(now time.Time) bool {
	if !t.task.HasDueDate() || t.IsDeleted() {
		return false
	}

	today := now.Format("2006-01-02")
	taskDateStr := t.task.DueDate.Format("2006-01-02")
	return taskDateStr == today
}

// IsDueTodayForCompleted checks if a completed task is due today and should still
// be shown in the Due Today filter (within the delay period)
func (t *Task) IsDueTodayForCompleted(config CompletedTaskTransitionConfig, now time.Time) bool {
	if !t.task.Completed || !t.task.HasDueDate() || t.IsDeleted() {
		return false
	}

	// Check if it's due today
	today := now.Format("2006-01-02")
	taskDateStr := t.task.DueDate.Format("2006-01-02")
	isDueToday := taskDateStr == today

	// If it's due today, check if it should still be shown in original filters
	return isDueToday && !t.ShouldRemoveFromOriginalFilters(config, now)
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

// IsThisWeekForCompleted checks if a completed task is due this week and should still
// be shown in the This Week filter (within the delay period)
func (t *Task) IsThisWeekForCompleted(config CompletedTaskTransitionConfig, now time.Time) bool {
	if !t.task.Completed || !t.task.HasDueDate() || t.IsDeleted() {
		return false
	}

	// Check if it's due this week
	if !t.IsThisWeek(now) {
		return false
	}

	// If it's due this week, check if it should still be shown in original filters
	return !t.ShouldRemoveFromOriginalFilters(config, now)
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

// Projects returns the projects associated with the task
func (t *Task) Projects() []string {
	return t.task.Projects
}

// Contexts returns the contexts associated with the task
func (t *Task) Contexts() []string {
	return t.task.Contexts
}

// HasDueDate returns true if the task has a due date
func (t *Task) HasDueDate() bool {
	return t.task.HasDueDate()
}

// IsDueThisWeek returns true if the task is due this week
func (t *Task) IsDueThisWeek(now time.Time) bool {
	if !t.HasDueDate() {
		return false
	}

	// Get the start of this week (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	dueDate := t.task.DueDate.Truncate(24 * time.Hour)
	return !dueDate.Before(startOfWeek) && dueDate.Before(endOfWeek)
}

// String returns the string representation of the task
func (t *Task) String() string {
	return t.task.String()
}

// ToTodoTxtTask returns the underlying todotxt.Task
func (t *Task) ToTodoTxtTask() *todotxt.Task {
	return t.task
}

// CyclePriority cycles through priority levels based on the provided configuration
func (t *Task) CyclePriority(priorityLevels []string) error {
	// Safety check: ensure PriorityLevels is not empty
	if len(priorityLevels) == 0 {
		return errors.New("priority levels cannot be empty")
	}

	currentPriority := ""
	if t.task.HasPriority() {
		currentPriority = t.task.Priority
	}

	// Find current priority index in configuration using lo.FindIndexOf
	_, currentIndex, found := lo.FindIndexOf(priorityLevels, func(priority string) bool {
		return priority == currentPriority
	})
	if !found {
		currentIndex = -1 // Use -1 to indicate not found, so next index becomes 0
	}

	// Move to next priority level (cycle around)
	nextIndex := (currentIndex + 1) % len(priorityLevels)
	nextPriority := priorityLevels[nextIndex]

	// Set the new priority
	if nextPriority == "" {
		t.task.Priority = ""
	} else {
		t.task.Priority = nextPriority
	}

	return nil
}

// ToggleDueToday toggles the due date of a task to today or removes it if already set to today
func (t *Task) ToggleDueToday(now time.Time) error {
	today := now.Format(DateFormat)

	// Get the current task string
	taskString := t.task.String()

	// Check if task is already due today
	hasDueToday := t.IsDueToday(now)

	var newTaskString string

	if hasDueToday {
		// Remove due date - remove due:YYYY-MM-DD from task string using lo.Filter
		parts := strings.Fields(taskString)
		newParts := lo.Filter(parts, func(part string, _ int) bool {
			return !strings.HasPrefix(part, TaskFieldDuePrefix)
		})
		newTaskString = strings.Join(newParts, " ")
	} else {
		// Add or update due date
		if t.task.HasDueDate() {
			// Replace existing due date using lo.Map
			parts := strings.Fields(taskString)
			newParts := lo.Map(parts, func(part string, _ int) string {
				if strings.HasPrefix(part, TaskFieldDuePrefix) {
					return TaskFieldDuePrefix + today
				}
				return part
			})
			newTaskString = strings.Join(newParts, " ")
		} else {
			// Add new due date
			newTaskString = taskString + " " + TaskFieldDuePrefix + today
		}
	}

	// Parse the new task string and update the task
	newTask, err := todotxt.ParseTask(newTaskString)
	if err != nil {
		return err
	}
	*t.task = *newTask
	return nil
}

// SoftDelete marks the task as deleted by adding a deleted_at field
func (t *Task) SoftDelete(now time.Time) error {
	taskString := t.task.String()

	// Check if already deleted
	if strings.Contains(taskString, TaskFieldDeletedPrefix) {
		return nil // Already deleted, no action needed
	}

	// Add deleted_at field to mark as soft deleted
	currentDate := now.Format(DateFormat)
	newTaskString := taskString + " " + TaskFieldDeletedPrefix + currentDate

	// Parse the modified task string back to update the task
	newTask, err := todotxt.ParseTask(newTaskString)
	if err != nil {
		return err
	}
	*t.task = *newTask
	return nil
}

// RestoreFromDeleted removes the deleted_at field to restore the task
func (t *Task) RestoreFromDeleted() error {
	taskString := t.task.String()

	// Check if task is deleted
	if !strings.Contains(taskString, TaskFieldDeletedPrefix) {
		return nil // Not deleted, no action needed
	}

	// Remove deleted_at field from the task string using lo.Filter
	parts := strings.Fields(taskString)
	cleanParts := lo.Filter(parts, func(part string, _ int) bool {
		return !strings.HasPrefix(part, TaskFieldDeletedPrefix)
	})
	newTaskString := strings.Join(cleanParts, " ")

	// Parse the modified task string back to update the task
	newTask, err := todotxt.ParseTask(newTaskString)
	if err != nil {
		return err
	}
	*t.task = *newTask
	return nil
}

// GetPriority returns the current priority of the task
func (t *Task) GetPriority() string {
	if t.task.HasPriority() {
		return t.task.Priority
	}
	return ""
}

// HasPriority returns true if the task has a priority
func (t *Task) HasPriority() bool {
	return t.task.HasPriority()
}

// GetDueDate returns the due date of the task
func (t *Task) GetDueDate() time.Time {
	return t.task.DueDate
}
