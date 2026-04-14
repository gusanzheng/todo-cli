package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"todo/internal/model"
)

func Path() string {
	if p := os.Getenv("TODO_FILE"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".todos.json"
	}
	return filepath.Join(home, ".todos.json")
}

func Load() ([]model.Todo, error) {
	path := Path()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []model.Todo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read todos, check %s", path)
	}
	var todos []model.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, fmt.Errorf("failed to read todos, check %s", path)
	}
	return todos, nil
}

func Save(todos []model.Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0644)
}

func NextID(todos []model.Todo) int {
	max := 0
	for _, t := range todos {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}
