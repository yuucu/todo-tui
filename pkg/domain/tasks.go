package domain

import (
	"sort"

	todotxt "github.com/1set/todotxt"
	"github.com/samber/lo"
)

// Tasks represents a collection of tasks with domain methods
type Tasks []Task

// TaskFilter represents a function that filters tasks
type TaskFilter func(Tasks) Tasks

// NewTasks creates a new Tasks instance from a TaskList
func NewTasks(taskList todotxt.TaskList) Tasks {
	tasks := make(Tasks, 0, len(taskList))
	for _, task := range taskList {
		domainTask, err := NewTask(&task)
		if err != nil {
			// Skip invalid tasks, but this should rarely happen
			continue
		}
		tasks = append(tasks, *domainTask)
	}
	return tasks
}

// ToTaskList converts Tasks back to todotxt.TaskList
func (t Tasks) ToTaskList() todotxt.TaskList {
	taskList := make(todotxt.TaskList, len(t))
	for i, task := range t {
		taskList[i] = *task.task
	}
	return taskList
}

// Len returns the number of tasks
func (t Tasks) Len() int {
	return len(t)
}

// Get returns the task at the specified index
// Panics if index is out of bounds - use SafeGet for safe access
func (t Tasks) Get(index int) Task {
	if index < 0 || index >= len(t) {
		panic("index out of range")
	}
	return t[index]
}

// SafeGet returns the task at the specified index and a boolean indicating whether the index was valid
func (t Tasks) SafeGet(index int) (Task, bool) {
	if index < 0 || index >= len(t) {
		return Task{}, false
	}
	return t[index], true
}

// Filter applies a filter function and returns a new Tasks instance
func (t Tasks) Filter(filterFn func(Task, int) bool) Tasks {
	var filtered Tasks
	for i, task := range t {
		if filterFn(task, i) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// FilterByStatus filters tasks by completion status
func (t Tasks) FilterByStatus(completed bool) Tasks {
	return t.Filter(func(task Task, _ int) bool {
		return task.IsCompleted() == completed
	})
}

// FilterActive returns only active (incomplete, non-deleted) tasks
func (t Tasks) FilterActive() Tasks {
	return t.Filter(func(task Task, _ int) bool {
		return !task.IsCompleted() && !task.IsDeleted()
	})
}

// FilterDeleted returns only deleted tasks
func (t Tasks) FilterDeleted() Tasks {
	return t.Filter(func(task Task, _ int) bool {
		return task.IsDeleted()
	})
}

// FilterByProject returns tasks that belong to the specified project
func (t Tasks) FilterByProject(project string) Tasks {
	return t.Filter(func(task Task, _ int) bool {
		return lo.Contains(task.Projects(), project)
	})
}

// FilterByContext returns tasks that belong to the specified context
func (t Tasks) FilterByContext(context string) Tasks {
	return t.Filter(func(task Task, _ int) bool {
		return lo.Contains(task.Contexts(), context)
	})
}

// FilterWithoutProjects returns tasks that have no projects
func (t Tasks) FilterWithoutProjects() Tasks {
	return t.Filter(func(task Task, _ int) bool {
		return len(task.Projects()) == 0
	})
}

// SortByCompletionStatus sorts tasks by completion status.
// Incomplete tasks are placed at the top, completed and deleted tasks at the bottom.
// The original order is preserved within each group (stable sort).
// Returns a new Tasks instance without modifying the original.
func (t Tasks) SortByCompletionStatus() Tasks {
	// Make a copy to avoid modifying the original slice
	sortedTasks := make(Tasks, len(t))
	copy(sortedTasks, t)

	// ソート: 完了したタスクを下の方に移動
	// 完了していないタスクを最初に、完了したタスクをその後に配置
	sort.SliceStable(sortedTasks, func(i, j int) bool {
		taskI := sortedTasks[i]
		taskJ := sortedTasks[j]

		// 削除されたタスクは完了タスクと同様に扱う
		isCompletedI := taskI.IsCompleted() || taskI.IsDeleted()
		isCompletedJ := taskJ.IsCompleted() || taskJ.IsDeleted()

		// 完了状態が異なる場合：未完了を上に、完了を下に
		if isCompletedI != isCompletedJ {
			return !isCompletedI // 未完了（false）が上に来る
		}

		// 両方とも同じ完了状態の場合：元の順序を保持（安定ソート）
		return false
	})

	return sortedTasks
}

// SortTasksByCompletionStatus sorts tasks by completion status.
// Incomplete tasks are placed at the top, completed and deleted tasks at the bottom.
// The original order is preserved within each group (stable sort).
//
// Deprecated: Use Tasks.SortByCompletionStatus() instead.
func SortTasksByCompletionStatus(tasks todotxt.TaskList) todotxt.TaskList {
	t := NewTasks(tasks)
	return t.SortByCompletionStatus().ToTaskList()
}
