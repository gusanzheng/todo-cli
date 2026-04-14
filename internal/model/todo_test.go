package model_test

import (
	"testing"
	"time"

	"todo/internal/model"
)

func TestTodoFields(t *testing.T) {
	now := time.Now()
	todo := model.Todo{
		ID:        1,
		Title:     "Buy milk",
		Done:      false,
		CreatedAt: now,
	}
	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}
	if todo.Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", todo.Title)
	}
	if todo.Done != false {
		t.Error("expected Done to be false")
	}
	if !todo.CreatedAt.Equal(now) {
		t.Error("expected CreatedAt to match")
	}
}

func TestTodoZeroValue(t *testing.T) {
	var todo model.Todo
	if todo.ID != 0 {
		t.Errorf("expected zero ID, got %d", todo.ID)
	}
	if todo.Title != "" {
		t.Errorf("expected empty title, got %q", todo.Title)
	}
	if todo.Done != false {
		t.Error("expected Done to be false")
	}
}

func TestTodoDateField(t *testing.T) {
	todo := model.Todo{
		ID:    1,
		Title: "Buy milk",
		Date:  "2026-04-14",
	}
	if todo.Date != "2026-04-14" {
		t.Errorf("expected Date '2026-04-14', got %q", todo.Date)
	}
}
