package domain

import (
	"testing"
	"time"

	todotxt "github.com/1set/todotxt"
)

// createTestTask creates a test task with the given text and optional completion status
func createTestTask(text string, completed bool) todotxt.Task {
	task, err := todotxt.ParseTask(text)
	if err != nil {
		panic(err)
	}

	if completed {
		task.Completed = true
		task.CompletedDate = time.Now()
	}

	return *task
}

// createDeletedTestTask creates a test task marked as deleted
func createDeletedTestTask(text string) todotxt.Task {
	task, err := todotxt.ParseTask(text + " deleted_at:2025-01-15")
	if err != nil {
		panic(err)
	}
	return *task
}

// TestTasks_Basic tests basic Tasks functionality
func TestTasks_Basic(t *testing.T) {
	taskList := todotxt.TaskList{
		createTestTask("First task", false),
		createTestTask("Second task", true),
	}

	tasks := NewTasks(taskList)

	// Test Len
	if tasks.Len() != 2 {
		t.Errorf("Expected length 2, got %d", tasks.Len())
	}

	// Test Get
	firstTask := tasks.Get(0)
	if firstTask == nil {
		t.Fatal("Expected task at index 0, got nil")
	}
	if firstTask.Todo != "First task" {
		t.Errorf("Expected 'First task', got '%s'", firstTask.Todo)
	}

	// Test Get out of bounds
	invalidTask := tasks.Get(5)
	if invalidTask != nil {
		t.Errorf("Expected nil for out of bounds index, got %v", invalidTask)
	}

	// Test ToTaskList
	convertedBack := tasks.ToTaskList()
	if len(convertedBack) != 2 {
		t.Errorf("Expected converted list length 2, got %d", len(convertedBack))
	}
}

// TestTasks_SortByCompletionStatus tests the new Tasks method
func TestTasks_SortByCompletionStatus_EmptyList(t *testing.T) {
	tasks := NewTasks(todotxt.TaskList{})
	result := tasks.SortByCompletionStatus()

	if result.Len() != 0 {
		t.Errorf("Expected empty list, got %d tasks", result.Len())
	}
}

func TestTasks_SortByCompletionStatus_OnlyIncompleteTasks(t *testing.T) {
	taskList := todotxt.TaskList{
		createTestTask("First incomplete task +project", false),
		createTestTask("Second incomplete task @context", false),
		createTestTask("Third incomplete task", false),
	}

	tasks := NewTasks(taskList)
	result := tasks.SortByCompletionStatus()

	// Should maintain original order for incomplete tasks
	if result.Len() != 3 {
		t.Errorf("Expected 3 tasks, got %d", result.Len())
	}

	for i := 0; i < result.Len(); i++ {
		task := result.Get(i)
		if task.Completed {
			t.Errorf("Task %d should be incomplete, but is completed", i)
		}
	}

	// Check original order is preserved (use the actual Todo field content)
	expectedTodos := []string{
		"First incomplete task", // Projects/contexts are parsed separately
		"Second incomplete task",
		"Third incomplete task",
	}

	for i, expected := range expectedTodos {
		task := result.Get(i)
		if task.Todo != expected {
			t.Errorf("Task %d: expected '%s', got '%s'", i, expected, task.Todo)
		}
	}

	// Verify projects and contexts were parsed correctly
	firstTask := result.Get(0)
	if len(firstTask.Projects) != 1 || firstTask.Projects[0] != "project" {
		t.Errorf("First task should have project 'project', got %v", firstTask.Projects)
	}

	secondTask := result.Get(1)
	if len(secondTask.Contexts) != 1 || secondTask.Contexts[0] != "context" {
		t.Errorf("Second task should have context 'context', got %v", secondTask.Contexts)
	}
}

func TestTasks_SortByCompletionStatus_MixedTasks(t *testing.T) {
	taskList := todotxt.TaskList{
		createTestTask("First incomplete task", false),
		createTestTask("First completed task", true),
		createTestTask("Second incomplete task", false),
		createTestTask("Second completed task", true),
		createTestTask("Third incomplete task", false),
	}

	tasks := NewTasks(taskList)
	result := tasks.SortByCompletionStatus()

	if result.Len() != 5 {
		t.Errorf("Expected 5 tasks, got %d", result.Len())
	}

	// First three tasks should be incomplete
	for i := 0; i < 3; i++ {
		task := result.Get(i)
		if task.Completed {
			t.Errorf("Task %d should be incomplete, but is completed", i)
		}
	}

	// Last two tasks should be completed
	for i := 3; i < 5; i++ {
		task := result.Get(i)
		if !task.Completed {
			t.Errorf("Task %d should be completed, but is incomplete", i)
		}
	}

	// Check that original order is preserved within each group
	expectedIncompleteOrder := []string{
		"First incomplete task",
		"Second incomplete task",
		"Third incomplete task",
	}

	for i, expected := range expectedIncompleteOrder {
		task := result.Get(i)
		if task.Todo != expected {
			t.Errorf("Incomplete task %d: expected '%s', got '%s'",
				i, expected, task.Todo)
		}
	}

	expectedCompletedOrder := []string{
		"First completed task",
		"Second completed task",
	}

	for i, expected := range expectedCompletedOrder {
		task := result.Get(i + 3) // offset by 3 incomplete tasks
		if task.Todo != expected {
			t.Errorf("Completed task %d: expected '%s', got '%s'",
				i, expected, task.Todo)
		}
	}
}

func TestTasks_SortByCompletionStatus_DoesNotModifyOriginal(t *testing.T) {
	originalTaskList := todotxt.TaskList{
		createTestTask("Incomplete task", false),
		createTestTask("Completed task", true),
	}

	tasks := NewTasks(originalTaskList)

	// Make a copy for comparison
	originalCopy := make(todotxt.TaskList, len(originalTaskList))
	copy(originalCopy, originalTaskList)

	result := tasks.SortByCompletionStatus()

	// Original should be unchanged
	if tasks.Len() != len(originalCopy) {
		t.Errorf("Original tasks length changed")
	}

	for i := 0; i < tasks.Len(); i++ {
		originalTask := tasks.Get(i)
		if originalTask.String() != originalCopy[i].String() {
			t.Errorf("Original task %d was modified", i)
		}
	}

	// Result should be different order
	firstResult := result.Get(0)
	secondResult := result.Get(1)
	if firstResult.Completed {
		t.Errorf("First task in result should be incomplete")
	}
	if !secondResult.Completed {
		t.Errorf("Second task in result should be completed")
	}
}

// 既存のテスト（後方互換性のため）

func TestSortTasksByCompletionStatus_EmptyList(t *testing.T) {
	tasks := todotxt.TaskList{}
	result := SortTasksByCompletionStatus(tasks)

	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d tasks", len(result))
	}
}

func TestSortTasksByCompletionStatus_OnlyIncompleteTasks(t *testing.T) {
	tasks := todotxt.TaskList{
		createTestTask("First incomplete task +project", false),
		createTestTask("Second incomplete task @context", false),
		createTestTask("Third incomplete task", false),
	}

	result := SortTasksByCompletionStatus(tasks)

	// Should maintain original order for incomplete tasks
	if len(result) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(result))
	}

	for i, task := range result {
		if task.Completed {
			t.Errorf("Task %d should be incomplete, but is completed", i)
		}
	}

	// Check original order is preserved (use the actual Todo field content)
	expectedTodos := []string{
		"First incomplete task", // Projects/contexts are parsed separately
		"Second incomplete task",
		"Third incomplete task",
	}

	for i, expected := range expectedTodos {
		if result[i].Todo != expected {
			t.Errorf("Task %d: expected '%s', got '%s'", i, expected, result[i].Todo)
		}
	}

	// Verify projects and contexts were parsed correctly
	if len(result[0].Projects) != 1 || result[0].Projects[0] != "project" {
		t.Errorf("First task should have project 'project', got %v", result[0].Projects)
	}

	if len(result[1].Contexts) != 1 || result[1].Contexts[0] != "context" {
		t.Errorf("Second task should have context 'context', got %v", result[1].Contexts)
	}
}

func TestSortTasksByCompletionStatus_OnlyCompletedTasks(t *testing.T) {
	tasks := todotxt.TaskList{
		createTestTask("First completed task", true),
		createTestTask("Second completed task", true),
		createTestTask("Third completed task", true),
	}

	result := SortTasksByCompletionStatus(tasks)

	// Should maintain original order for completed tasks
	if len(result) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(result))
	}

	for i, task := range result {
		if !task.Completed {
			t.Errorf("Task %d should be completed, but is incomplete", i)
		}
	}
}

func TestSortTasksByCompletionStatus_MixedTasks(t *testing.T) {
	tasks := todotxt.TaskList{
		createTestTask("First incomplete task", false),
		createTestTask("First completed task", true),
		createTestTask("Second incomplete task", false),
		createTestTask("Second completed task", true),
		createTestTask("Third incomplete task", false),
	}

	result := SortTasksByCompletionStatus(tasks)

	if len(result) != 5 {
		t.Errorf("Expected 5 tasks, got %d", len(result))
	}

	// First three tasks should be incomplete
	incompleteTasks := result[:3]
	for i, task := range incompleteTasks {
		if task.Completed {
			t.Errorf("Task %d should be incomplete, but is completed", i)
		}
	}

	// Last two tasks should be completed
	completedTasks := result[3:]
	for i, task := range completedTasks {
		if !task.Completed {
			t.Errorf("Completed task %d should be completed, but is incomplete", i)
		}
	}

	// Check that original order is preserved within each group
	expectedIncompleteOrder := []string{
		"First incomplete task",
		"Second incomplete task",
		"Third incomplete task",
	}

	for i, expected := range expectedIncompleteOrder {
		if incompleteTasks[i].Todo != expected {
			t.Errorf("Incomplete task %d: expected '%s', got '%s'",
				i, expected, incompleteTasks[i].Todo)
		}
	}

	expectedCompletedOrder := []string{
		"First completed task",
		"Second completed task",
	}

	for i, expected := range expectedCompletedOrder {
		if completedTasks[i].Todo != expected {
			t.Errorf("Completed task %d: expected '%s', got '%s'",
				i, expected, completedTasks[i].Todo)
		}
	}
}

func TestSortTasksByCompletionStatus_WithDeletedTasks(t *testing.T) {
	tasks := todotxt.TaskList{
		createTestTask("Incomplete task", false),
		createDeletedTestTask("Deleted task"),
		createTestTask("Completed task", true),
		createTestTask("Another incomplete task", false),
	}

	result := SortTasksByCompletionStatus(tasks)

	if len(result) != 4 {
		t.Errorf("Expected 4 tasks, got %d", len(result))
	}

	// First two tasks should be incomplete
	incompleteTasks := result[:2]
	for i, task := range incompleteTasks {
		if task.Completed {
			t.Errorf("Task %d should be incomplete, but is completed", i)
		}
	}

	// Last two tasks should be completed or deleted (treated as completed)
	completedOrDeletedTasks := result[2:]
	for i, task := range completedOrDeletedTasks {
		// Check if it's either completed or has deleted_at field
		hasDeletedField := task.String() != "" &&
			(task.Completed || len(task.String()) > len(task.Todo))

		if !hasDeletedField && !task.Completed {
			t.Errorf("Task %d should be completed or deleted, but appears to be active", i)
		}
	}
}

func TestSortTasksByCompletionStatus_DoesNotModifyOriginal(t *testing.T) {
	originalTasks := todotxt.TaskList{
		createTestTask("Incomplete task", false),
		createTestTask("Completed task", true),
	}

	// Make a copy for comparison
	originalCopy := make(todotxt.TaskList, len(originalTasks))
	copy(originalCopy, originalTasks)

	result := SortTasksByCompletionStatus(originalTasks)

	// Original should be unchanged
	if len(originalTasks) != len(originalCopy) {
		t.Errorf("Original slice length changed")
	}

	for i, task := range originalTasks {
		if task.String() != originalCopy[i].String() {
			t.Errorf("Original task %d was modified", i)
		}
	}

	// Result should be different order
	if result[0].Completed {
		t.Errorf("First task in result should be incomplete")
	}
	if !result[1].Completed {
		t.Errorf("Second task in result should be completed")
	}
}

func TestSortTasksByCompletionStatus_StableSort(t *testing.T) {
	// Create multiple tasks with same completion status to test stability
	tasks := todotxt.TaskList{
		createTestTask("Alpha incomplete", false),
		createTestTask("Beta incomplete", false),
		createTestTask("Alpha completed", true),
		createTestTask("Beta completed", true),
		createTestTask("Gamma incomplete", false),
		createTestTask("Gamma completed", true),
	}

	result := SortTasksByCompletionStatus(tasks)

	// Check that incomplete tasks maintain their relative order
	incompleteTasks := result[:3]
	expectedIncompleteOrder := []string{
		"Alpha incomplete",
		"Beta incomplete",
		"Gamma incomplete",
	}

	for i, expected := range expectedIncompleteOrder {
		if incompleteTasks[i].Todo != expected {
			t.Errorf("Incomplete task order not stable: expected '%s', got '%s'",
				expected, incompleteTasks[i].Todo)
		}
	}

	// Check that completed tasks maintain their relative order
	completedTasks := result[3:]
	expectedCompletedOrder := []string{
		"Alpha completed",
		"Beta completed",
		"Gamma completed",
	}

	for i, expected := range expectedCompletedOrder {
		if completedTasks[i].Todo != expected {
			t.Errorf("Completed task order not stable: expected '%s', got '%s'",
				expected, completedTasks[i].Todo)
		}
	}
}
