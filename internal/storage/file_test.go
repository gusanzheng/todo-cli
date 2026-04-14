package storage_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"todo/internal/model"
	"todo/internal/storage"
)

func setupTempStorage(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "todos.json")
	t.Setenv("TODO_FILE", p)
	return p
}

func TestLoadFileMissing(t *testing.T) {
	setupTempStorage(t)
	todos, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected 0 todos, got %d", len(todos))
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTempStorage(t)
	todos := []model.Todo{
		{ID: 1, Title: "Buy milk", Done: false, CreatedAt: time.Now().UTC().Truncate(time.Second)},
	}
	if err := storage.Save(todos); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(loaded))
	}
	if loaded[0].Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", loaded[0].Title)
	}
	if loaded[0].Done != false {
		t.Error("expected Done false")
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "todos.json")
	os.WriteFile(p, []byte("{not valid json"), 0644)
	t.Setenv("TODO_FILE", p)
	_, err := storage.Load()
	if err == nil {
		t.Error("expected error for corrupt JSON, got nil")
	}
}

func TestNextIDEmpty(t *testing.T) {
	if got := storage.NextID(nil); got != 1 {
		t.Errorf("expected NextID=1 for nil slice, got %d", got)
	}
}

func TestNextIDMax(t *testing.T) {
	todos := []model.Todo{{ID: 3}, {ID: 1}, {ID: 2}}
	if got := storage.NextID(todos); got != 4 {
		t.Errorf("expected NextID=4, got %d", got)
	}
}
