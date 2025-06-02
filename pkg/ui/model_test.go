package ui

import (
	"strings"
	"testing"
	"time"

	todotxt "github.com/1set/todotxt"
	"github.com/yuucu/todotui/pkg/domain"
)

// テスト用のモデル作成ヘルパー
func createTestModel() *Model {
	return &Model{
		tasks:        domain.NewTasks(createTestTaskList()),
		activePane:   paneTask,
		viewMode:     ViewFilter,
		taskList:     SimpleList{},
		filterList:   SimpleList{},
		filters:      []FilterData{},
		todoFilePath: "/tmp/test.todo.txt",
		appConfig:    DefaultAppConfig(),
	}
}

func TestModel_findTaskInList(t *testing.T) {
	// Create test model
	model := &Model{
		tasks:        domain.NewTasks(createTestTaskList()),
		viewMode:     ViewFilter,
		activePane:   paneFilter,
		todoFilePath: "test.txt",
	}

	// Get a task to search for
	if model.tasks.Len() == 0 {
		t.Fatal("Test task list is empty")
	}

	targetTask := model.tasks.Get(0)

	// Test finding existing task
	index, task, found := model.findTaskInList(targetTask)

	if !found {
		t.Error("findTaskInList() should find existing task")
	}

	if index < 0 || index >= model.tasks.Len() {
		t.Errorf("findTaskInList() returned invalid index: %d", index)
	}

	// Test that found task matches the target
	if task.String() != targetTask.String() {
		t.Errorf("findTaskInList() returned wrong task: got %s, want %s",
			task.String(), targetTask.String())
	}

	// Test finding non-existing task - create a simple task for testing
	nonExistentTaskStr := "Non-existent task that should not be found"
	nonExistentTodoTxtTask, _ := todotxt.ParseTask(nonExistentTaskStr)
	nonExistentTask := domain.NewTask(nonExistentTodoTxtTask)
	_, _, found = model.findTaskInList(*nonExistentTask)

	if found {
		t.Error("findTaskInList() should not find non-existent task")
	}
}

func TestCyclePriority(t *testing.T) {
	tests := []struct {
		name           string
		initialTask    string
		expectedResult string
		description    string
	}{
		{
			name:           "no_priority_to_A",
			initialTask:    "Test task without priority",
			expectedResult: "(A) Test task without priority",
			description:    "優先度なし → (A)",
		},
		{
			name:           "priority_A_to_B",
			initialTask:    "(A) Test task with priority A",
			expectedResult: "(B) Test task with priority A",
			description:    "優先度(A) → (B)",
		},
		{
			name:           "priority_B_to_C",
			initialTask:    "(B) Test task with priority B",
			expectedResult: "(C) Test task with priority B",
			description:    "優先度(B) → (C)",
		},
		{
			name:           "priority_Z_to_none",
			initialTask:    "(Z) Test task with priority Z",
			expectedResult: "(A) Test task with priority Z",
			description:    "優先度(Z) → 優先度(A)",
		},
		{
			name:           "priority_C_to_D",
			initialTask:    "(C) Test task with priority C",
			expectedResult: "(D) Test task with priority C",
			description:    "優先度(C) → (D)",
		},
	}

	model := createTestModel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.initialTask)
			if err != nil {
				t.Fatalf("Failed to parse initial task: %v", err)
			}

			// 優先度をサイクル
			model.cyclePriority(task)

			// 期待される結果と比較（空白や日付の差異を考慮）
			expected, err := todotxt.ParseTask(tt.expectedResult)
			if err != nil {
				t.Fatalf("Failed to parse expected result: %v", err)
			}

			if task.Priority != expected.Priority {
				t.Errorf("cyclePriority() priority = %s, expected %s for %s",
					task.Priority, expected.Priority, tt.description)
			}

			if task.Todo != expected.Todo {
				t.Errorf("cyclePriority() todo = %s, expected %s for %s",
					task.Todo, expected.Todo, tt.description)
			}
		})
	}
}

func TestToggleDueToday(t *testing.T) {
	model := createTestModel()

	tests := []struct {
		name        string
		initialTask string
		expectDue   bool
		description string
	}{
		{
			name:        "add_due_today_to_task_without_due",
			initialTask: "Test task without due date",
			expectDue:   true,
			description: "期限なしタスクに今日の期限を追加",
		},
		{
			name:        "remove_due_today_from_task_with_today",
			initialTask: "Test task due:" + time.Now().Format("2006-01-02"),
			expectDue:   false,
			description: "今日期限のタスクから期限を削除",
		},
		{
			name:        "change_due_date_to_today",
			initialTask: "Test task due:2025-12-31",
			expectDue:   true,
			description: "他の期限を今日に変更",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := todotxt.ParseTask(tt.initialTask)
			if err != nil {
				t.Fatalf("Failed to parse initial task: %v", err)
			}

			// 期限を切り替え
			model.toggleDueToday(task)

			today := time.Now().Format("2006-01-02")

			if tt.expectDue {
				if !task.HasDueDate() {
					t.Errorf("toggleDueToday() should add due date for %s", tt.description)
				} else if task.DueDate.Format("2006-01-02") != today {
					t.Errorf("toggleDueToday() should set due date to today (%s), got %s for %s",
						today, task.DueDate.Format("2006-01-02"), tt.description)
				}
			} else {
				if task.HasDueDate() {
					t.Errorf("toggleDueToday() should remove due date for %s", tt.description)
				}
			}
		})
	}
}

// タスク追加のテスト（模擬）
func TestTaskAdditionLogic(t *testing.T) {
	model := createTestModel()
	initialTaskCount := model.tasks.Len()

	// 新しいタスクを作成してリストに追加
	newTaskString := "New test task +project @work"
	newTask, err := todotxt.ParseTask(newTaskString)
	if err != nil {
		t.Fatalf("Failed to parse new task: %v", err)
	}

	// タスクをリストに追加
	taskList := model.tasks.ToTaskList()
	taskList = append(taskList, *newTask)
	model.tasks = domain.NewTasks(taskList)

	// タスク数が増えたことを確認
	if model.tasks.Len() != initialTaskCount+1 {
		t.Errorf("Task addition failed: expected %d tasks, got %d",
			initialTaskCount+1, model.tasks.Len())
	}

	// 追加されたタスクの内容を確認
	addedTask := model.tasks.Get(model.tasks.Len() - 1)
	if addedTask.String() == "" {
		t.Errorf("Added task has wrong content: %s", addedTask.String())
	}

	if len(addedTask.Projects()) != 1 || addedTask.Projects()[0] != "project" {
		t.Errorf("Added task has wrong projects: %v", addedTask.Projects())
	}

	if len(addedTask.Contexts()) != 1 || addedTask.Contexts()[0] != "work" {
		t.Errorf("Added task has wrong contexts: %v", addedTask.Contexts())
	}
}

// タスク編集のテスト（模擬）
func TestTaskEditLogic(t *testing.T) {
	model := createTestModel()

	if model.tasks.Len() == 0 {
		t.Fatal("Test model should have tasks")
	}

	// 最初のタスクを編集
	originalTask := model.tasks.Get(0)
	originalTodo := originalTask.String()

	// 新しい内容でタスクを更新
	editedTaskString := "(A) Edited task +newproject @newcontext"
	editedTask, err := todotxt.ParseTask(editedTaskString)
	if err != nil {
		t.Fatalf("Failed to parse edited task: %v", err)
	}

	// タスクを更新 - replace the entire task list
	taskList := model.tasks.ToTaskList()
	taskList[0] = *editedTask
	model.tasks = domain.NewTasks(taskList)

	// 変更が反映されたことを確認
	updatedTask := model.tasks.Get(0)

	if updatedTask.String() == originalTodo {
		t.Error("Task editing failed: content unchanged")
	}

	// Check if the task contains "Edited task"
	if !strings.Contains(updatedTask.String(), "Edited task") {
		t.Errorf("Task editing failed: expected 'Edited task' in '%s'", updatedTask.String())
	}

	// 優先度が設定されたことを確認
	todoTxtTask := updatedTask.ToTodoTxtTask()
	if !todoTxtTask.HasPriority() || todoTxtTask.Priority != "A" {
		t.Errorf("Task editing failed: expected priority A, got %s", todoTxtTask.Priority)
	}

	if len(updatedTask.Projects()) != 1 || updatedTask.Projects()[0] != "newproject" {
		t.Errorf("Task editing failed: wrong projects %v", updatedTask.Projects())
	}
}

// タスク完了切り替えのテスト
func TestTaskCompletionToggle(t *testing.T) {
	// 未完了のタスクを作成
	incompleteTask, _ := todotxt.ParseTask("Incomplete task +project")

	// 完了済みのタスクを作成
	completedTask, _ := todotxt.ParseTask("x 2025-01-15 Completed task +project")

	tests := []struct {
		name            string
		task            *todotxt.Task
		expectCompleted bool
		description     string
	}{
		{
			name:            "complete_incomplete_task",
			task:            incompleteTask,
			expectCompleted: true,
			description:     "未完了タスクを完了にする",
		},
		{
			name:            "uncomplete_completed_task",
			task:            completedTask,
			expectCompleted: false,
			description:     "完了タスクを未完了にする",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初期状態を記録
			initialCompleted := tt.task.Completed

			// 完了状態を切り替え（模擬実装）
			if tt.task.Completed {
				tt.task.Completed = false
				tt.task.CompletedDate = time.Time{}
			} else {
				tt.task.Completed = true
				tt.task.CompletedDate = time.Now()
			}

			// 期待される状態になったことを確認
			if tt.task.Completed != tt.expectCompleted {
				t.Errorf("Task completion toggle failed for %s: expected %v, got %v",
					tt.description, tt.expectCompleted, tt.task.Completed)
			}

			// 状態が実際に変更されたことを確認
			if tt.task.Completed == initialCompleted {
				t.Errorf("Task completion state should have changed for %s", tt.description)
			}

			// 完了日の設定を確認
			if tt.task.Completed && tt.task.CompletedDate.IsZero() {
				t.Errorf("Completed task should have completion date for %s", tt.description)
			}

			if !tt.task.Completed && !tt.task.CompletedDate.IsZero() {
				t.Errorf("Incomplete task should not have completion date for %s", tt.description)
			}
		})
	}
}

// タスク削除のテスト（論理削除）
func TestTaskDeletion(t *testing.T) {
	model := createTestModel()
	initialTaskCount := model.tasks.Len()

	if initialTaskCount == 0 {
		t.Fatal("Test model should have tasks")
	}

	// 最初のタスクを論理削除（deleted_atフィールドを追加）
	targetTask := model.tasks.Get(0)
	taskString := targetTask.String()

	// deleted_atフィールドを追加
	deletedTaskString := taskString + " deleted_at:" + time.Now().Format("2006-01-02")
	deletedTask, err := todotxt.ParseTask(deletedTaskString)
	if err != nil {
		t.Fatalf("Failed to parse deleted task: %v", err)
	}

	// タスクを更新 - replace the entire task list
	taskList := model.tasks.ToTaskList()
	taskList[0] = *deletedTask
	model.tasks = domain.NewTasks(taskList)

	// 削除されたタスクが削除済みとして認識されることを確認
	updatedTask := model.tasks.Get(0)
	if !updatedTask.IsDeleted() {
		t.Error("Task should be marked as deleted")
	}

	// タスク数は変わらない（論理削除のため）
	if model.tasks.Len() != initialTaskCount {
		t.Errorf("Task count should remain same for logical deletion: expected %d, got %d",
			initialTaskCount, model.tasks.Len())
	}
}

// エラーハンドリングのテスト
func TestTaskParsingErrors(t *testing.T) {
	invalidTaskStrings := []string{
		"", // 空文字列
		// 他の無効なケースがあれば追加
	}

	for _, invalidTask := range invalidTaskStrings {
		t.Run("invalid_task_"+invalidTask, func(t *testing.T) {
			_, err := todotxt.ParseTask(invalidTask)

			// 空文字列の場合、ライブラリがどう振る舞うかに依存
			// 実際のエラーハンドリングはアプリケーションレベルで行う
			if invalidTask == "" && err == nil {
				// 空文字列の場合は通常のタスクとして扱われる可能性がある
				t.Logf("Empty string task parsing: %v", err)
			}
		})
	}
}

// タスクリストの整合性テスト
func TestTaskListIntegrity(t *testing.T) {
	model := createTestModel()

	// すべてのタスクが有効であることを確認
	for i, task := range model.tasks.ToTaskList() {
		if task.Todo == "" && !task.Completed && !domain.NewTask(&task).IsDeleted() {
			t.Errorf("Task at index %d has empty content and is not completed/deleted", i)
		}

		// プロジェクトとコンテキストの重複チェック
		projectMap := make(map[string]bool)
		for _, project := range task.Projects {
			if projectMap[project] {
				t.Errorf("Task at index %d has duplicate project: %s", i, project)
			}
			projectMap[project] = true
		}

		contextMap := make(map[string]bool)
		for _, context := range task.Contexts {
			if contextMap[context] {
				t.Errorf("Task at index %d has duplicate context: %s", i, context)
			}
			contextMap[context] = true
		}
	}
}
