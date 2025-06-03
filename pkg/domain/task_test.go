package domain

import (
	"testing"
	"time"

	todotxt "github.com/1set/todotxt"
)

func TestNewTask_Error(t *testing.T) {
	task, err := NewTask(nil)
	if err == nil {
		t.Error("Expected error when creating task with nil input")
	}
	if task != nil {
		t.Error("Expected nil task when error occurs")
	}
	if err.Error() != "task cannot be nil" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestNewTask_Success(t *testing.T) {
	todoTxtTask, _ := todotxt.ParseTask("Test task")
	task, err := NewTask(todoTxtTask)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if task == nil {
		t.Error("Expected valid task, got nil")
	}
	if task.String() != "Test task" {
		t.Errorf("Expected 'Test task', got %s", task.String())
	}
}

func TestTask_ToggleCompletion(t *testing.T) {
	tests := []struct {
		name            string
		taskString      string
		initialComplete bool
		expectedResult  bool
	}{
		{
			name:            "complete_incomplete_task",
			taskString:      "Test incomplete task",
			initialComplete: false,
			expectedResult:  true,
		},
		{
			name:            "uncomplete_completed_task",
			taskString:      "x 2025-01-15 Test completed task",
			initialComplete: true,
			expectedResult:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			if domainTask.IsCompleted() != tt.initialComplete {
				t.Errorf("Initial completion state mismatch: expected %v, got %v",
					tt.initialComplete, domainTask.IsCompleted())
			}

			result := domainTask.ToggleCompletion()

			if result != tt.expectedResult {
				t.Errorf("ToggleCompletion() = %v, expected %v", result, tt.expectedResult)
			}

			if domainTask.IsCompleted() != tt.expectedResult {
				t.Errorf("Task completion state after toggle: expected %v, got %v",
					tt.expectedResult, domainTask.IsCompleted())
			}
		})
	}
}

func TestTask_ShouldMoveToCompleted(t *testing.T) {
	tests := []struct {
		name           string
		taskString     string
		config         CompletedTaskTransitionConfig
		testTime       time.Time
		expectedResult bool
	}{
		{
			name:       "immediate_move_delay_0",
			taskString: "x 2025-01-15 Completed task",
			config: CompletedTaskTransitionConfig{
				DelayDays:      0,
				TransitionHour: 9,
			},
			testTime:       time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name:       "not_yet_time_same_day",
			taskString: "x 2025-01-15 Completed task",
			config: CompletedTaskTransitionConfig{
				DelayDays:      1,
				TransitionHour: 9,
			},
			testTime:       time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			expectedResult: false,
		},
		{
			name:       "time_to_move_next_day_before_transition",
			taskString: "x 2025-01-15 Completed task",
			config: CompletedTaskTransitionConfig{
				DelayDays:      1,
				TransitionHour: 9,
			},
			testTime:       time.Date(2025, 1, 16, 8, 59, 0, 0, time.UTC),
			expectedResult: false,
		},
		{
			name:       "time_to_move_next_day_after_transition",
			taskString: "x 2025-01-15 Completed task",
			config: CompletedTaskTransitionConfig{
				DelayDays:      1,
				TransitionHour: 9,
			},
			testTime:       time.Date(2025, 1, 16, 9, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			result := domainTask.ShouldMoveToCompleted(tt.config, tt.testTime)

			if result != tt.expectedResult {
				t.Errorf("ShouldMoveToCompleted() = %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}

func TestTask_IsDeleted(t *testing.T) {
	tests := []struct {
		name       string
		taskString string
		expected   bool
	}{
		{
			name:       "normal_task",
			taskString: "Buy milk +grocery @home",
			expected:   false,
		},
		{
			name:       "deleted_task",
			taskString: "Deleted task deleted_at:2025-01-15",
			expected:   true,
		},
		{
			name:       "deleted_task_with_time",
			taskString: "Another deleted task deleted_at:2025-01-15T10:30:00",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			result := domainTask.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v for task: %s",
					result, tt.expected, tt.taskString)
			}
		})
	}
}

func TestTask_IsOverdue(t *testing.T) {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	weekAgo := now.AddDate(0, 0, -7).Format("2006-01-02")
	today := now.Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")

	tests := []struct {
		name       string
		taskString string
		expected   bool
	}{
		{
			name:       "overdue_task_yesterday",
			taskString: "Overdue task due:" + yesterday,
			expected:   true,
		},
		{
			name:       "overdue_task_week_ago",
			taskString: "Very overdue task due:" + weekAgo,
			expected:   true,
		},
		{
			name:       "not_overdue_today",
			taskString: "Task due today due:" + today,
			expected:   false,
		},
		{
			name:       "not_overdue_tomorrow",
			taskString: "Future task due:" + tomorrow,
			expected:   false,
		},
		{
			name:       "not_overdue_no_due_date",
			taskString: "Task without due date",
			expected:   false,
		},
		{
			name:       "not_overdue_completed",
			taskString: "x " + today + " Completed task due:" + yesterday,
			expected:   true,
		},
		{
			name:       "not_overdue_deleted",
			taskString: "Deleted task due:" + yesterday + " deleted_at:" + today,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			result := domainTask.IsOverdue(now)
			if result != tt.expected {
				t.Errorf("IsOverdue() = %v, expected %v for task: %s",
					result, tt.expected, tt.taskString)
			}
		})
	}
}

func TestTask_IsDueToday(t *testing.T) {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	today := now.Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")

	tests := []struct {
		name       string
		taskString string
		expected   bool
	}{
		{
			name:       "due_today",
			taskString: "Task due today due:" + today,
			expected:   true,
		},
		{
			name:       "not_due_today_yesterday",
			taskString: "Task due yesterday due:" + yesterday,
			expected:   false,
		},
		{
			name:       "not_due_today_tomorrow",
			taskString: "Task due tomorrow due:" + tomorrow,
			expected:   false,
		},
		{
			name:       "not_due_today_no_date",
			taskString: "Task without due date",
			expected:   false,
		},
		{
			name:       "not_due_today_completed",
			taskString: "x " + today + " Completed task due:" + today,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			result := domainTask.IsDueToday(now)
			if result != tt.expected {
				t.Errorf("IsDueToday() = %v, expected %v for task: %s",
					result, tt.expected, tt.taskString)
			}
		})
	}
}

func TestTask_IsThisWeek(t *testing.T) {
	// Wednesday, January 15, 2025 (week: Jan 12-18, 2025)
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// This week dates
	sunday := time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	wednesday := now.Format("2006-01-02")

	// Other week dates
	lastWeek := time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC).Format("2006-01-02") // Saturday before
	nextWeek := time.Date(2025, 1, 19, 0, 0, 0, 0, time.UTC).Format("2006-01-02") // Sunday after

	tests := []struct {
		name       string
		taskString string
		expected   bool
	}{
		{
			name:       "this_week_today",
			taskString: "Task due today due:" + wednesday,
			expected:   true,
		},
		{
			name:       "this_week_sunday",
			taskString: "Task due this Sunday due:" + sunday,
			expected:   true,
		},
		{
			name:       "this_week_wednesday",
			taskString: "Task due this Wednesday due:" + wednesday,
			expected:   true,
		},
		{
			name:       "not_this_week_last_week",
			taskString: "Task due last week due:" + lastWeek,
			expected:   false,
		},
		{
			name:       "not_this_week_next_week",
			taskString: "Task due next week due:" + nextWeek,
			expected:   false,
		},
		{
			name:       "not_this_week_no_date",
			taskString: "Task without due date",
			expected:   false,
		},
		{
			name:       "not_this_week_completed",
			taskString: "x " + wednesday + " Completed task due:" + wednesday,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			result := domainTask.IsThisWeek(now)
			if result != tt.expected {
				t.Errorf("IsThisWeek() = %v, expected %v for task: %s",
					result, tt.expected, tt.taskString)
			}
		})
	}
}

func TestTask_CyclePriority(t *testing.T) {
	tests := []struct {
		name           string
		taskText       string
		priorityLevels []string
		expectedSteps  []string
		shouldError    bool
	}{
		{
			name:           "Basic priority cycling",
			taskText:       "Test task",
			priorityLevels: []string{"", "A", "B", "C"},
			expectedSteps:  []string{"(A) Test task", "(B) Test task", "(C) Test task", "Test task"},
			shouldError:    false,
		},
		{
			name:           "Cycling with existing priority A",
			taskText:       "(A) Test task",
			priorityLevels: []string{"", "A", "B", "C"},
			expectedSteps:  []string{"(B) Test task", "(C) Test task", "Test task", "(A) Test task"},
			shouldError:    false,
		},
		{
			name:           "Cycling with priority not in levels",
			taskText:       "(Z) Test task",
			priorityLevels: []string{"", "A", "B", "C"},
			expectedSteps:  []string{"Test task", "(A) Test task", "(B) Test task", "(C) Test task"},
			shouldError:    false,
		},
		{
			name:           "Single priority level",
			taskText:       "Test task",
			priorityLevels: []string{"A"},
			expectedSteps:  []string{"(A) Test task", "(A) Test task"},
			shouldError:    false,
		},
		{
			name:           "Empty priority levels",
			taskText:       "Test task",
			priorityLevels: []string{},
			expectedSteps:  nil,
			shouldError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			if tt.shouldError {
				err := task.CyclePriority(tt.priorityLevels)
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			// Test cycling through all expected steps
			for i, expected := range tt.expectedSteps {
				err := task.CyclePriority(tt.priorityLevels)
				if err != nil {
					t.Errorf("Step %d: Unexpected error: %v", i, err)
					continue
				}

				actual := task.String()
				if actual != expected {
					t.Errorf("Step %d: Expected %q, got %q", i, expected, actual)
				}
			}
		})
	}
}

func TestTask_ToggleDueToday(t *testing.T) {
	testTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		taskText string
		expected string
	}{
		{
			name:     "Add due today to task without due date",
			taskText: "Test task",
			expected: "Test task due:2025-01-15",
		},
		{
			name:     "Remove due today from task due today",
			taskText: "Test task due:2025-01-15",
			expected: "Test task",
		},
		{
			name:     "Change due tomorrow to due today",
			taskText: "Test task due:2025-01-16",
			expected: "Test task due:2025-01-15",
		},
		{
			name:     "Add due today to task with priority",
			taskText: "(A) Test task",
			expected: "(A) Test task due:2025-01-15",
		},
		{
			name:     "Remove due today from prioritized task",
			taskText: "(A) Test task due:2025-01-15",
			expected: "(A) Test task",
		},
		{
			name:     "Toggle with project and context",
			taskText: "Test task +project @context",
			expected: "Test task @context +project due:2025-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			err = task.ToggleDueToday(testTime)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			actual := task.String()
			if actual != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestTask_ToggleDueToday_DoubleToggle(t *testing.T) {
	testTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	original := "Test task"

	todoTask, err := todotxt.ParseTask(original)
	if err != nil {
		t.Fatalf("Failed to parse task: %v", err)
	}

	task, err := NewTask(todoTask)
	if err != nil {
		t.Fatalf("Failed to create domain task: %v", err)
	}

	// First toggle: add due today
	err = task.ToggleDueToday(testTime)
	if err != nil {
		t.Fatalf("First toggle failed: %v", err)
	}

	expected1 := "Test task due:2025-01-15"
	if task.String() != expected1 {
		t.Errorf("After first toggle: expected %q, got %q", expected1, task.String())
	}

	// Second toggle: remove due today
	err = task.ToggleDueToday(testTime)
	if err != nil {
		t.Fatalf("Second toggle failed: %v", err)
	}

	if task.String() != original {
		t.Errorf("After second toggle: expected %q, got %q", original, task.String())
	}
}

func TestTask_SoftDelete(t *testing.T) {
	testTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		taskText string
		expected string
	}{
		{
			name:     "Delete simple task",
			taskText: "Test task",
			expected: "Test task deleted_at:2025-01-15",
		},
		{
			name:     "Delete task with priority",
			taskText: "(A) Test task",
			expected: "(A) Test task deleted_at:2025-01-15",
		},
		{
			name:     "Delete task with project and context",
			taskText: "Test task +project @context",
			expected: "Test task @context +project deleted_at:2025-01-15",
		},
		{
			name:     "Delete already deleted task (no change)",
			taskText: "Test task deleted_at:2025-01-10",
			expected: "Test task deleted_at:2025-01-10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			err = task.SoftDelete(testTime)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			actual := task.String()
			if actual != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, actual)
			}

			// Verify IsDeleted returns true
			if !task.IsDeleted() {
				t.Errorf("Task should be marked as deleted")
			}
		})
	}
}

func TestTask_RestoreFromDeleted(t *testing.T) {
	tests := []struct {
		name     string
		taskText string
		expected string
	}{
		{
			name:     "Restore simple deleted task",
			taskText: "Test task deleted_at:2025-01-15",
			expected: "Test task",
		},
		{
			name:     "Restore deleted task with priority",
			taskText: "(A) Test task deleted_at:2025-01-15",
			expected: "(A) Test task",
		},
		{
			name:     "Restore deleted task with project and context",
			taskText: "Test task +project @context deleted_at:2025-01-15",
			expected: "Test task @context +project",
		},
		{
			name:     "Restore task that's not deleted (no change)",
			taskText: "Test task",
			expected: "Test task",
		},
		{
			name:     "Restore task with multiple fields including deleted_at",
			taskText: "Test task +project @context due:2025-01-20 deleted_at:2025-01-15",
			expected: "Test task @context +project due:2025-01-20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			err = task.RestoreFromDeleted()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			actual := task.String()
			if actual != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, actual)
			}

			// Verify IsDeleted returns false
			if task.IsDeleted() {
				t.Errorf("Task should not be marked as deleted after restoration")
			}
		})
	}
}

func TestTask_GetPriority(t *testing.T) {
	tests := []struct {
		name     string
		taskText string
		expected string
	}{
		{
			name:     "Task with priority A",
			taskText: "(A) Test task",
			expected: "A",
		},
		{
			name:     "Task with priority Z",
			taskText: "(Z) Test task",
			expected: "Z",
		},
		{
			name:     "Task without priority",
			taskText: "Test task",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			actual := task.GetPriority()
			if actual != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestTask_HasPriority(t *testing.T) {
	tests := []struct {
		name     string
		taskText string
		expected bool
	}{
		{
			name:     "Task with priority",
			taskText: "(A) Test task",
			expected: true,
		},
		{
			name:     "Task without priority",
			taskText: "Test task",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			actual := task.HasPriority()
			if actual != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestTask_GetDueDate(t *testing.T) {
	tests := []struct {
		name     string
		taskText string
		hasDate  bool
	}{
		{
			name:     "Task with due date",
			taskText: "Test task due:2025-01-15",
			hasDate:  true,
		},
		{
			name:     "Task without due date",
			taskText: "Test task",
			hasDate:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoTask, err := todotxt.ParseTask(tt.taskText)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			task, err := NewTask(todoTask)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			actual := task.GetDueDate()
			if tt.hasDate {
				if actual.IsZero() {
					t.Errorf("Expected a valid date, got zero time")
				}
				// Check if the date is 2025-01-15, regardless of timezone
				if actual.Year() != 2025 || actual.Month() != 1 || actual.Day() != 15 {
					t.Errorf("Expected date 2025-01-15, got %v", actual)
				}
			} else {
				if !actual.IsZero() {
					t.Errorf("Expected zero time, got %v", actual)
				}
			}
		})
	}
}

// Integration test for soft delete and restore workflow
func TestTask_SoftDeleteAndRestore_Integration(t *testing.T) {
	testTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	original := "(A) Important task +project @context due:2025-01-20"

	todoTask, err := todotxt.ParseTask(original)
	if err != nil {
		t.Fatalf("Failed to parse task: %v", err)
	}

	task, err := NewTask(todoTask)
	if err != nil {
		t.Fatalf("Failed to create domain task: %v", err)
	}

	// Initial state - not deleted
	if task.IsDeleted() {
		t.Errorf("Task should not be deleted initially")
	}

	// Soft delete
	err = task.SoftDelete(testTime)
	if err != nil {
		t.Fatalf("Failed to soft delete: %v", err)
	}

	// Should be marked as deleted
	if !task.IsDeleted() {
		t.Errorf("Task should be marked as deleted")
	}

	expectedDeleted := "(A) Important task @context +project deleted_at:2025-01-15 due:2025-01-20"
	if task.String() != expectedDeleted {
		t.Errorf("Expected %q, got %q", expectedDeleted, task.String())
	}

	// Restore from deleted
	err = task.RestoreFromDeleted()
	if err != nil {
		t.Fatalf("Failed to restore: %v", err)
	}

	// Should not be marked as deleted
	if task.IsDeleted() {
		t.Errorf("Task should not be marked as deleted after restoration")
	}

	expectedRestored := "(A) Important task @context +project due:2025-01-20"
	if task.String() != expectedRestored {
		t.Errorf("Expected %q, got %q", expectedRestored, task.String())
	}
}
