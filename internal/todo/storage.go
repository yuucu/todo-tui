package todo

import (
	"os"
	"path/filepath"

	todotxt "github.com/1set/todotxt"
)

// ディレクトリ作成時のデフォルトパーミッション
const defaultDirMode = 0755

// Load reads a todo.txt file and returns a TaskList
func Load(path string) (todotxt.TaskList, error) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, defaultDirMode); err != nil {
		return nil, err
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	return todotxt.LoadFromPath(path)
}

// Save writes a TaskList to a todo.txt file
func Save(list todotxt.TaskList, path string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, defaultDirMode); err != nil {
		return err
	}

	return list.WriteToPath(path)
}
