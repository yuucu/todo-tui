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
			expected:   false,
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

			result := domainTask.IsThisWeek(now)
			if result != tt.expected {
				t.Errorf("IsThisWeek() = %v, expected %v for task: %s",
					result, tt.expected, tt.taskString)
			}
		})
	}
}
