package ui

import (
	"testing"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/yuucu/todotui/pkg/domain"
)

func TestIsOverdue(t *testing.T) {
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
			name:        "overdue_yesterday",
			taskString:  "Test task due:2025-05-30",
			now:         baseTime,
			expected:    true,
			description: "昨日が期限のタスクは期限切れ",
		},
		{
			name:        "overdue_week_ago",
			taskString:  "Test task due:2025-05-24",
			now:         baseTime,
			expected:    true,
			description: "1週間前が期限のタスクは期限切れ",
		},
		{
			name:        "not_overdue_today",
			taskString:  "Test task due:2025-05-31",
			now:         baseTime,
			expected:    false,
			description: "今日が期限のタスクは期限切れではない",
		},
		{
			name:        "not_overdue_tomorrow",
			taskString:  "Test task due:2025-06-01",
			now:         baseTime,
			expected:    false,
			description: "明日が期限のタスクは期限切れではない",
		},
		{
			name:        "not_overdue_no_due_date",
			taskString:  "Test task without due date",
			now:         baseTime,
			expected:    false,
			description: "期限なしのタスクは期限切れではない",
		},
		{
			name:        "not_overdue_completed",
			taskString:  "x 2025-05-31 Test completed task due:2025-05-30",
			now:         baseTime,
			expected:    false,
			description: "完了済みタスクは期限切れとしない",
		},
		{
			name:        "not_overdue_deleted",
			taskString:  "Test deleted task due:2025-05-30 deleted_at:2025-05-31",
			now:         baseTime,
			expected:    false,
			description: "削除済みタスクは期限切れとしない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := domain.NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}
			result := domainTask.IsOverdue(tt.now)
			if result != tt.expected {
				t.Errorf("IsOverdue() = %v, expected %v for %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestIsDueToday(t *testing.T) {
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
			expected:    false,
			description: "完了済みタスクは今日期限として判定しない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := domain.NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}
			result := domainTask.IsDueToday(tt.now)
			if result != tt.expected {
				t.Errorf("IsDueToday() = %v, expected %v for %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestIsThisWeek(t *testing.T) {
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
			expected:    false,
			description: "完了済みタスクは今週として判定しない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := domain.NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}
			result := domainTask.IsThisWeek(tt.now)
			if result != tt.expected {
				t.Errorf("IsThisWeek() = %v, expected %v for %s", result, tt.expected, tt.description)
			}
		})
	}
}

// 境界値テスト - 日付の境界での動作を確認
func TestOverdueBoundaryConditions(t *testing.T) {
	// 2025年5月31日の23:59:59
	endOfDay := time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC)
	// 2025年6月1日の00:00:00
	startOfNextDay := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		taskString string
		now        time.Time
		expected   bool
	}{
		{
			name:       "overdue_at_end_of_day",
			taskString: "Test task due:2025-05-30",
			now:        endOfDay,
			expected:   true,
		},
		{
			name:       "not_overdue_today_at_end_of_day",
			taskString: "Test task due:2025-05-31",
			now:        endOfDay,
			expected:   false,
		},
		{
			name:       "overdue_at_start_of_next_day",
			taskString: "Test task due:2025-05-31",
			now:        startOfNextDay,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.taskString)
			if err != nil {
				t.Fatalf("Failed to parse task: %v", err)
			}

			domainTask, err := domain.NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}
			result := domainTask.IsOverdue(tt.now)
			if result != tt.expected {
				t.Errorf("IsOverdue() = %v, expected %v at %v", result, tt.expected, tt.now)
			}
		})
	}
}
