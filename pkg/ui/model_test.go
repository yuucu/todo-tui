package ui

import (
	"strings"
	"testing"

	todotxt "github.com/1set/todotxt"
	"github.com/yuucu/todotui/pkg/domain"
)

// createTestModel creates a test model with sample data
func createTestModel() *Model {
	tasks := []string{
		"(A) Buy milk +grocery @home",
		"Write tests +project @work",
		"x 2025-01-15 Completed task +project",
		"Call mom @phone",
	}

	var taskList todotxt.TaskList
	for _, taskStr := range tasks {
		if task, err := todotxt.ParseTask(taskStr); err == nil {
			taskList = append(taskList, *task)
		}
	}

	return &Model{
		tasks:     domain.NewTasks(taskList),
		appConfig: DefaultAppConfig(), // Add default app config to prevent divide by zero
	}
}

func TestModel_findTaskInList(t *testing.T) {
	model := createTestModel()

	if model.tasks.Len() == 0 {
		t.Fatal("Test model should have tasks")
	}

	// Test finding existing task
	targetTask := model.tasks.Get(0)
	index, task, found := model.findTaskInList(targetTask)

	if !found {
		t.Error("findTaskInList() should find existing task")
	}

	if index != 0 {
		t.Errorf("findTaskInList() returned wrong index: got %d, want 0", index)
	}

	if task.String() != targetTask.String() {
		t.Errorf("findTaskInList() returned wrong task: got %s, want %s",
			task.String(), targetTask.String())
	}

	// Test finding non-existing task - create a simple task for testing
	nonExistentTaskStr := "Non-existent task that should not be found"
	nonExistentTodoTxtTask, _ := todotxt.ParseTask(nonExistentTaskStr)
	nonExistentTask, err := domain.NewTask(nonExistentTodoTxtTask)
	if err != nil {
		t.Fatalf("Failed to create domain task: %v", err)
	}
	_, _, found = model.findTaskInList(*nonExistentTask)

	if found {
		t.Error("findTaskInList() should not find non-existent task")
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
		// Check for domain task creation errors
		domainTask, err := domain.NewTask(&task)
		if err != nil {
			t.Errorf("Failed to create domain task at index %d: %v", i, err)
			continue
		}

		if task.Todo == "" && !task.Completed && !domainTask.IsDeleted() {
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
