package domain

import (
	"testing"
	"time"

	todotxt "github.com/1set/todotxt"
)

func TestTask_ToggleCompletion(t *testing.T) {
	tests := []struct {
		name            string
		taskString      string
		expectCompleted bool
	}{
		{
			name:            "complete_incomplete_task",
			taskString:      "Buy groceries +shopping",
			expectCompleted: true,
		},
		{
			name:            "uncomplete_completed_task",
			taskString:      "x 2025-01-15 Completed task +project",
			expectCompleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask := NewTask(task)
			initialCompleted := domainTask.IsCompleted()

			result := domainTask.ToggleCompletion()

			// Check if the result matches expectation
			if result != tt.expectCompleted {
				t.Errorf("ToggleCompletion() returned %v, expected %v", result, tt.expectCompleted)
			}

			// Check if the completion state actually changed
			if domainTask.IsCompleted() == initialCompleted {
				t.Errorf("Task completion state did not change")
			}

			// Check completion date
			if domainTask.IsCompleted() && domainTask.GetCompletedDate().IsZero() {
				t.Errorf("Completed task should have completion date")
			}

			if !domainTask.IsCompleted() && !domainTask.GetCompletedDate().IsZero() {
				t.Errorf("Incomplete task should not have completion date")
			}
		})
	}
}

func TestTask_ShouldMoveToCompleted(t *testing.T) {
	// Create a completed task
	task, _ := todotxt.ParseTask("Buy groceries +shopping")
	domainTask := NewTask(task)
	domainTask.ToggleCompletion() // Complete the task

	tests := []struct {
		name             string
		config           CompletedTaskTransitionConfig
		completionOffset time.Duration // How much time ago the task was completed
		expected         bool
	}{
		{
			name: "immediate_move_delay_0",
			config: CompletedTaskTransitionConfig{
				DelayDays:      0,
				TransitionHour: 5,
			},
			completionOffset: 0,
			expected:         true,
		},
		{
			name: "not_yet_time_same_day",
			config: CompletedTaskTransitionConfig{
				DelayDays:      1,
				TransitionHour: 5,
			},
			completionOffset: 12 * time.Hour, // 12 hours ago, but same day
			expected:         false,
		},
		{
			name: "time_to_move_next_day",
			config: CompletedTaskTransitionConfig{
				DelayDays:      1,
				TransitionHour: 5,
			},
			completionOffset: 25 * time.Hour, // Over 1 day ago
			expected:         false,          // Still false because current time (04:41) < transition hour (05:00)
		},
		{
			name: "time_to_move_well_past_transition",
			config: CompletedTaskTransitionConfig{
				DelayDays:      1,
				TransitionHour: 3, // 3 AM transition time
			},
			completionOffset: 25 * time.Hour, // Over 1 day ago, and current time (04:41) > transition hour (03:00)
			expected:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set completion date to a specific time in the past
			completionTime := time.Now().Add(-tt.completionOffset)
			domainTask.task.CompletedDate = completionTime

			result := domainTask.ShouldMoveToCompleted(tt.config)

			if result != tt.expected {
				t.Errorf("ShouldMoveToCompleted() = %v, expected %v", result, tt.expected)
				t.Logf("Completion time: %v", completionTime)
				t.Logf("Current time: %v", time.Now())
				t.Logf("Config: %+v", tt.config)
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
			taskString: "Buy groceries +shopping",
			expected:   false,
		},
		{
			name:       "deleted_task",
			taskString: "Buy groceries +shopping deleted_at:2025-01-15",
			expected:   true,
		},
		{
			name:       "deleted_task_with_time",
			taskString: "Buy groceries +shopping deleted_at:2025-01-15T10:30:00",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask := NewTask(task)
			result := domainTask.IsDeleted()

			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestTask_IsOverdue(t *testing.T) {
	// テスト用の基準日時を設定（2025年5月31日 21:00）
	baseTime := time.Date(2025, 5, 31, 21, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		taskString  string
		now         time.Time
		expected    bool
		description string
	}{
		{
			name:        "overdue_task_yesterday",
			taskString:  "Test task due:2025-05-30",
			now:         baseTime,
			expected:    true,
			description: "昨日が期限のタスクはoverdue",
		},
		{
			name:        "overdue_task_week_ago",
			taskString:  "Test task due:2025-05-24",
			now:         baseTime,
			expected:    true,
			description: "1週間前が期限のタスクはoverdue",
		},
		{
			name:        "not_overdue_today",
			taskString:  "Test task due:2025-05-31",
			now:         baseTime,
			expected:    false,
			description: "今日が期限のタスクはoverdueではない",
		},
		{
			name:        "not_overdue_tomorrow",
			taskString:  "Test task due:2025-06-01",
			now:         baseTime,
			expected:    false,
			description: "明日が期限のタスクはoverdueではない",
		},
		{
			name:        "not_overdue_no_due_date",
			taskString:  "Test task without due date",
			now:         baseTime,
			expected:    false,
			description: "期限なしのタスクはoverdueではない",
		},
		{
			name:        "not_overdue_completed",
			taskString:  "x 2025-05-31 Test completed task due:2025-05-30",
			now:         baseTime,
			expected:    true,
			description: "完了済みタスクでもoverdueになる",
		},
		{
			name:        "not_overdue_deleted",
			taskString:  "Test deleted task due:2025-05-30 deleted_at:2025-05-31",
			now:         baseTime,
			expected:    false,
			description: "削除済みタスクはoverdueではない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask := NewTask(task)
			result := domainTask.IsOverdue(tt.now)
			if result != tt.expected {
				t.Errorf("IsOverdue() = %v, expected %v for %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestTask_IsDueToday(t *testing.T) {
	// テスト用の基準日時を設定（2025年5月31日 21:00）
	baseTime := time.Date(2025, 5, 31, 21, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		taskString  string
		now         time.Time
		expected    bool
		description string
	}{
		{
			name:        "due_today",
			taskString:  "Test task due:2025-05-31",
			now:         baseTime,
			expected:    true,
			description: "今日が期限のタスク",
		},
		{
			name:        "not_due_today_yesterday",
			taskString:  "Test task due:2025-05-30",
			now:         baseTime,
			expected:    false,
			description: "昨日が期限のタスクは今日ではない",
		},
		{
			name:        "not_due_today_tomorrow",
			taskString:  "Test task due:2025-06-01",
			now:         baseTime,
			expected:    false,
			description: "明日が期限のタスクは今日ではない",
		},
		{
			name:        "not_due_today_no_date",
			taskString:  "Test task without due date",
			now:         baseTime,
			expected:    false,
			description: "期限なしのタスクは今日ではない",
		},
		{
			name:        "not_due_today_completed",
			taskString:  "x 2025-05-31 Test completed task due:2025-05-31",
			now:         baseTime,
			expected:    true,
			description: "完了済みタスクでも今日になる",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask := NewTask(task)
			result := domainTask.IsDueToday(tt.now)
			if result != tt.expected {
				t.Errorf("IsDueToday() = %v, expected %v for %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestTask_IsThisWeek(t *testing.T) {
	// テスト用の基準日時を設定（2025年5月31日 土曜日 21:00）
	// この週は 2025-05-25(日) から 2025-05-31(土) まで
	baseTime := time.Date(2025, 5, 31, 21, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		taskString  string
		now         time.Time
		expected    bool
		description string
	}{
		{
			name:        "this_week_today",
			taskString:  "Test task due:2025-05-31",
			now:         baseTime,
			expected:    true,
			description: "今日（今週土曜日）が期限のタスク",
		},
		{
			name:        "this_week_sunday",
			taskString:  "Test task due:2025-05-25",
			now:         baseTime,
			expected:    true,
			description: "今週日曜日が期限のタスク",
		},
		{
			name:        "this_week_wednesday",
			taskString:  "Test task due:2025-05-28",
			now:         baseTime,
			expected:    true,
			description: "今週水曜日が期限のタスク",
		},
		{
			name:        "not_this_week_last_week",
			taskString:  "Test task due:2025-05-24",
			now:         baseTime,
			expected:    false,
			description: "先週土曜日が期限のタスクは今週ではない",
		},
		{
			name:        "not_this_week_next_week",
			taskString:  "Test task due:2025-06-01",
			now:         baseTime,
			expected:    false,
			description: "来週日曜日が期限のタスクは今週ではない",
		},
		{
			name:        "not_this_week_no_date",
			taskString:  "Test task without due date",
			now:         baseTime,
			expected:    false,
			description: "期限なしのタスクは今週ではない",
		},
		{
			name:        "not_this_week_completed",
			taskString:  "x 2025-05-31 Test completed task due:2025-05-31",
			now:         baseTime,
			expected:    true,
			description: "完了済みタスクでも今週になる",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask := NewTask(task)
			result := domainTask.IsThisWeek(tt.now)
			if result != tt.expected {
				t.Errorf("IsThisWeek() = %v, expected %v for %s", result, tt.expected, tt.description)
			}
		})
	}
}
