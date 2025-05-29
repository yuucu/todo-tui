package todo

import (
	"os"
	"path/filepath"
	"testing"

	todotxt "github.com/1set/todotxt"
)

func TestLoadAndSave(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "todo-tui-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	todoPath := filepath.Join(tmpDir, "test.todo.txt")

	// Test loading from non-existent file (should create empty file)
	list, err := Load(todoPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(list))
	}

	// Add some test tasks
	task1, err := todotxt.ParseTask("(A) Buy milk +grocery @home")
	if err != nil {
		t.Fatal(err)
	}
	task2, err := todotxt.ParseTask("Write tests +project @work")
	if err != nil {
		t.Fatal(err)
	}

	list.AddTask(task1)
	list.AddTask(task2)

	// Test saving
	if err := Save(list, todoPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Test loading again
	loadedList, err := Load(todoPath)
	if err != nil {
		t.Fatalf("Second load failed: %v", err)
	}

	if len(loadedList) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(loadedList))
	}

	// Check task content
	if loadedList[0].Todo != "Buy milk" {
		t.Errorf("First task content mismatch: %s", loadedList[0].Todo)
	}

	if loadedList[1].Todo != "Write tests" {
		t.Errorf("Second task content mismatch: %s", loadedList[1].Todo)
	}

	// Check project and context tags
	if len(loadedList[0].Projects) != 1 || loadedList[0].Projects[0] != "grocery" {
		t.Errorf("First task project mismatch: %v", loadedList[0].Projects)
	}

	if len(loadedList[0].Contexts) != 1 || loadedList[0].Contexts[0] != "home" {
		t.Errorf("First task context mismatch: %v", loadedList[0].Contexts)
	}
}

func TestLoadNonExistentDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "todo-tui-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test path with non-existent subdirectory
	todoPath := filepath.Join(tmpDir, "subdir", "test.todo.txt")

	list, err := Load(todoPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(list))
	}

	// Directory should be created
	if _, err := os.Stat(filepath.Dir(todoPath)); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}
