package ui

import (
	"testing"

	todotxt "github.com/1set/todotxt"
	"github.com/yuucu/todotui/pkg/domain"
)

// テスト用のタスクリスト作成ヘルパー
func createTestTaskList() todotxt.TaskList {
	tasks := []string{
		"Buy milk +grocery @home",            // アクティブタスク、プロジェクト・コンテキスト付き
		"Write tests +project @work",         // アクティブタスク、プロジェクト・コンテキスト付き
		"Read book",                          // アクティブタスク、タグなし
		"Call mom @phone",                    // アクティブタスク、コンテキスト付き
		"x 2024-01-15 Completed task +other", // 完了済みタスク
		"Deleted task deleted_at:2024-01-15", // 削除済みタスク
	}

	var taskList todotxt.TaskList
	for _, taskStr := range tasks {
		if task, err := todotxt.ParseTask(taskStr); err == nil {
			taskList = append(taskList, *task)
		}
	}

	return taskList
}

func TestIsTaskDeleted(t *testing.T) {
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

			domainTask, err := domain.NewTask(task)
			if err != nil {
				t.Fatalf("Failed to create domain task: %v", err)
			}

			result := domainTask.IsDeleted()
			if result != tt.expected {
				t.Errorf("domain.NewTask(&).IsDeleted() = %v, expected %v for task: %s",
					result, tt.expected, tt.taskString)
			}
		})
	}
}

func TestGetUniqueProjects(t *testing.T) {
	// モデルを作成（最小限の設定）
	model := &Model{
		tasks: domain.NewTasks(createTestTaskList()),
	}

	projects := model.getUniqueProjects()

	// 期待されるプロジェクト（ソート済み）- アクティブタスクのみ
	expected := []string{"grocery", "project"}

	if len(projects) != len(expected) {
		t.Errorf("getUniqueProjects() returned %d projects, expected %d",
			len(projects), len(expected))
		t.Logf("Actual projects: %v", projects)
		t.Logf("Expected projects: %v", expected)
	}

	for i, project := range projects {
		if i >= len(expected) || project != expected[i] {
			t.Errorf("getUniqueProjects()[%d] = %s, expected %s",
				i, project, expected[i])
		}
	}
}

func TestGetUniqueContexts(t *testing.T) {
	// モデルを作成（最小限の設定）
	model := &Model{
		tasks: domain.NewTasks(createTestTaskList()),
	}

	contexts := model.getUniqueContexts()

	// 期待されるコンテキスト（ソート済み）- アクティブタスクのみ
	expected := []string{"home", "phone", "work"}

	if len(contexts) != len(expected) {
		t.Errorf("getUniqueContexts() returned %d contexts, expected %d",
			len(contexts), len(expected))
		t.Logf("Actual contexts: %v", contexts)
		t.Logf("Expected contexts: %v", expected)
	}

	for i, context := range contexts {
		if i >= len(expected) || context != expected[i] {
			t.Errorf("getUniqueContexts()[%d] = %s, expected %s",
				i, context, expected[i])
		}
	}
}

func TestProjectFilter(t *testing.T) {
	tasks := createTestTaskList()

	// プロジェクト "grocery" でフィルタ
	groceryFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed {
				domainTask, err := domain.NewTask(&task)
				if err != nil {
					continue // Skip invalid tasks
				}
				if !domainTask.IsDeleted() {
					for _, project := range task.Projects {
						if project == "grocery" {
							filtered = append(filtered, task)
							break
						}
					}
				}
			}
		}
		return filtered
	}

	filtered := groceryFilter(tasks)

	// "grocery" プロジェクトを持つアクティブなタスクが1つ期待される
	if len(filtered) != 1 {
		t.Errorf("Project filter for 'grocery' returned %d tasks, expected 1", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].Todo != "Buy milk" {
		t.Errorf("Project filter returned wrong task: %s", filtered[0].Todo)
	}
}

func TestContextFilter(t *testing.T) {
	tasks := createTestTaskList()

	// コンテキスト "@work" でフィルタ
	workFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed {
				domainTask, err := domain.NewTask(&task)
				if err != nil {
					continue // Skip invalid tasks
				}
				if !domainTask.IsDeleted() {
					for _, context := range task.Contexts {
						if context == "work" {
							filtered = append(filtered, task)
							break
						}
					}
				}
			}
		}
		return filtered
	}

	filtered := workFilter(tasks)

	// "@work" コンテキストを持つアクティブなタスクが1つ期待される
	if len(filtered) != 1 {
		t.Errorf("Context filter for '@work' returned %d tasks, expected 1", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].Todo != "Write tests" {
		t.Errorf("Context filter returned wrong task: %s", filtered[0].Todo)
	}
}

func TestAllTasksFilter(t *testing.T) {
	tasks := createTestTaskList()

	// 全タスクフィルタ（完了済み・削除済みを除く）
	allTasksFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed {
				domainTask, err := domain.NewTask(&task)
				if err != nil {
					continue // Skip invalid tasks
				}
				if !domainTask.IsDeleted() {
					filtered = append(filtered, task)
				}
			}
		}
		return filtered
	}

	filtered := allTasksFilter(tasks)

	// アクティブなタスクが4つ期待される（完了済み・削除済みを除く）
	expected := 4
	if len(filtered) != expected {
		t.Errorf("All tasks filter returned %d tasks, expected %d", len(filtered), expected)
	}

	// 完了済みタスクと削除済みタスクが含まれていないことを確認
	for _, task := range filtered {
		if task.Completed {
			t.Errorf("All tasks filter included completed task: %s", task.Todo)
		}
		domainTask, err := domain.NewTask(&task)
		if err != nil {
			t.Errorf("Failed to create domain task for validation: %v", err)
			continue
		}
		if domainTask.IsDeleted() {
			t.Errorf("All tasks filter included deleted task: %s", task.Todo)
		}
	}
}

func TestNoProjectFilter(t *testing.T) {
	tasks := createTestTaskList()

	// プロジェクトなしフィルタ
	noProjectFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed && len(task.Projects) == 0 {
				domainTask, err := domain.NewTask(&task)
				if err != nil {
					continue // Skip invalid tasks
				}
				if !domainTask.IsDeleted() {
					filtered = append(filtered, task)
				}
			}
		}
		return filtered
	}

	filtered := noProjectFilter(tasks)

	// プロジェクトなしのタスクが2つ期待される: "Read book", "Call mom @phone"
	expected := 2
	if len(filtered) != expected {
		t.Errorf("No project filter returned %d tasks, expected %d", len(filtered), expected)
		for i, task := range filtered {
			t.Logf("Filtered task %d: %s (projects: %v)", i, task.Todo, task.Projects)
		}
	}

	// すべてのタスクがプロジェクトを持たないことを確認
	for _, task := range filtered {
		if len(task.Projects) > 0 {
			t.Errorf("No project filter included task with projects: %s (projects: %v)",
				task.Todo, task.Projects)
		}
	}
}

func TestCompletedTasksFilter(t *testing.T) {
	tasks := createTestTaskList()

	// 完了済みタスクフィルタ
	completedFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			if task.Completed {
				domainTask, err := domain.NewTask(&task)
				if err != nil {
					continue // Skip invalid tasks
				}
				if !domainTask.IsDeleted() {
					filtered = append(filtered, task)
				}
			}
		}
		return filtered
	}

	filtered := completedFilter(tasks)

	// 完了済みタスクが1つ期待される
	expected := 1
	if len(filtered) != expected {
		t.Errorf("Completed tasks filter returned %d tasks, expected %d", len(filtered), expected)
	}

	// すべてのタスクが完了済みであることを確認
	for _, task := range filtered {
		if !task.Completed {
			t.Errorf("Completed tasks filter included non-completed task: %s", task.Todo)
		}
		domainTask, err := domain.NewTask(&task)
		if err != nil {
			t.Errorf("Failed to create domain task for validation: %v", err)
			continue
		}
		if domainTask.IsDeleted() {
			t.Errorf("Completed tasks filter included deleted task: %s", task.Todo)
		}
	}
}

func TestDeletedTasksFilter(t *testing.T) {
	tasks := createTestTaskList()

	// 削除済みタスクフィルタ
	deletedFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			domainTask, err := domain.NewTask(&task)
			if err != nil {
				continue // Skip invalid tasks
			}
			if domainTask.IsDeleted() {
				filtered = append(filtered, task)
			}
		}
		return filtered
	}

	filtered := deletedFilter(tasks)

	// 削除済みタスクが1つ期待される
	expected := 1
	if len(filtered) != expected {
		t.Errorf("Deleted tasks filter returned %d tasks, expected %d", len(filtered), expected)
	}

	// すべてのタスクが削除済みであることを確認
	for _, task := range filtered {
		domainTask, err := domain.NewTask(&task)
		if err != nil {
			t.Errorf("Failed to create domain task for validation: %v", err)
			continue
		}
		if !domainTask.IsDeleted() {
			t.Errorf("Deleted tasks filter included non-deleted task: %s", task.Todo)
		}
	}
}

// エッジケースのテスト
func TestFilterWithEmptyTaskList(t *testing.T) {
	emptyTasks := todotxt.TaskList{}

	// 空のタスクリストに対するフィルタ
	allTasksFilter := func(tasks todotxt.TaskList) todotxt.TaskList {
		var filtered todotxt.TaskList
		for _, task := range tasks {
			if !task.Completed {
				domainTask, err := domain.NewTask(&task)
				if err != nil {
					continue // Skip invalid tasks
				}
				if !domainTask.IsDeleted() {
					filtered = append(filtered, task)
				}
			}
		}
		return filtered
	}

	filtered := allTasksFilter(emptyTasks)

	if len(filtered) != 0 {
		t.Errorf("Filter on empty task list returned %d tasks, expected 0", len(filtered))
	}
}

func TestFilterIntegration(t *testing.T) {
	// 複数のフィルタが正しく動作することを確認する統合テスト
	tasks := createTestTaskList()

	// モデルを作成
	model := &Model{
		tasks: domain.NewTasks(tasks),
	}

	// プロジェクトリストの取得
	projects := model.getUniqueProjects()
	if len(projects) == 0 {
		t.Error("No projects found in test task list")
	}

	// コンテキストリストの取得
	contexts := model.getUniqueContexts()
	if len(contexts) == 0 {
		t.Error("No contexts found in test task list")
	}

	// 各プロジェクトに対するフィルタが動作することを確認
	for _, project := range projects {
		projectFilter := func(p string) func(todotxt.TaskList) todotxt.TaskList {
			return func(tasks todotxt.TaskList) todotxt.TaskList {
				var filtered todotxt.TaskList
				for _, task := range tasks {
					if !task.Completed {
						domainTask, err := domain.NewTask(&task)
						if err != nil {
							continue // Skip invalid tasks
						}
						if !domainTask.IsDeleted() {
							for _, taskProject := range task.Projects {
								if taskProject == p {
									filtered = append(filtered, task)
									break
								}
							}
						}
					}
				}
				return filtered
			}
		}(project)

		filtered := projectFilter(tasks)

		// 各フィルタが少なくとも何らかの結果を返すことを確認
		if len(filtered) == 0 {
			t.Logf("Warning: Project filter for '%s' returned no tasks", project)
		}

		// フィルタされたタスクが条件を満たすことを確認
		for _, task := range filtered {
			found := false
			for _, taskProject := range task.Projects {
				if taskProject == project {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Project filter for '%s' returned task without that project: %s",
					project, task.Todo)
			}
		}
	}
}
