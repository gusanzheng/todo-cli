package cmd_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"todo/cmd"
	"todo/internal/model"
)

func init() {
	color.NoColor = true
}

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := cmd.NewRootCmd()
	root.SilenceUsage = true
	root.SilenceErrors = true
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func setupTempStorage(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "todos.json")
	t.Setenv("TODO_FILE", p)
	return p
}

func TestAddCreatesAndOutputs(t *testing.T) {
	setupTempStorage(t)
	out, err := run(t, "add", "Buy milk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Added") || !strings.Contains(out, "Buy milk") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestAddEmptyTitleErrors(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "add", "")
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestAddIncrementsID(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "First")
	out, err := run(t, "add", "Second")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[2]") {
		t.Errorf("expected ID [2] in output, got %q", out)
	}
}

func TestAddPersistsToFile(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	var todos []model.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}
	if len(todos) != 1 || todos[0].Title != "Buy milk" {
		t.Errorf("unexpected file contents: %v", todos)
	}
}

func TestListEmpty(t *testing.T) {
	setupTempStorage(t)
	out, err := run(t, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No todos yet") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestListShowsTodos(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Buy milk")
	run(t, "add", "Write tests")
	out, err := run(t, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Buy milk") || !strings.Contains(out, "Write tests") {
		t.Errorf("expected both todos in output, got %q", out)
	}
}

func TestListPendingSymbol(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Buy milk")
	out, _ := run(t, "list")
	if !strings.Contains(out, "○") {
		t.Errorf("expected pending symbol ○, got %q", out)
	}
}

func TestDoneMarksTodo(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	out, err := run(t, "done", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Marked") {
		t.Errorf("expected 'Marked' in output, got %q", out)
	}
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
	if !todos[0].Done {
		t.Error("expected todo[0].Done to be true")
	}
}

func TestDoneNotFound(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "done", "99")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestDoneInvalidID(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "done", "abc")
	if err == nil {
		t.Error("expected error for non-integer ID")
	}
}

func TestDeleteRemovesTodo(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	out, err := run(t, "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Deleted") {
		t.Errorf("expected 'Deleted' in output, got %q", out)
	}
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
	if len(todos) != 0 {
		t.Errorf("expected 0 todos after delete, got %d", len(todos))
	}
}

func TestDeleteNotFound(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "delete", "99")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestDeleteInvalidID(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "delete", "abc")
	if err == nil {
		t.Error("expected error for non-integer ID")
	}
}

func TestUndoneMarksTodoPending(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	run(t, "done", "1")
	out, err := run(t, "undone", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Unmarked") {
		t.Errorf("expected 'Unmarked' in output, got %q", out)
	}
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
	if todos[0].Done {
		t.Error("expected todo[0].Done to be false after undone")
	}
}

func TestUndoneNotFound(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "undone", "99")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestUndoneInvalidID(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "undone", "abc")
	if err == nil {
		t.Error("expected error for non-integer ID")
	}
}

func TestAddSetsDateToToday(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
	today := time.Now().UTC().Format(model.DateFormat)
	if todos[0].Date != today {
		t.Errorf("expected Date %q, got %q", today, todos[0].Date)
	}
}

func TestDateSetsDate(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	out, err := run(t, "date", "1", "2026-04-20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2026-04-20") {
		t.Errorf("expected date in output, got %q", out)
	}
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
	if todos[0].Date != "2026-04-20" {
		t.Errorf("expected Date '2026-04-20', got %q", todos[0].Date)
	}
}

func TestDateInvalidFormat(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Buy milk")
	_, err := run(t, "date", "1", "14-04-2026")
	if err == nil {
		t.Error("expected error for invalid date format")
	}
}

func TestDateNotFound(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "date", "99", "2026-04-14")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestDateInvalidID(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "date", "abc", "2026-04-14")
	if err == nil {
		t.Error("expected error for non-integer ID")
	}
}

func TestListShowsDoneItems(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Buy milk")
	run(t, "done", "1")
	out, err := run(t, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Buy milk") {
		t.Errorf("expected done todo in list output, got %q", out)
	}
}

func TestListFilterDone(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Task A")
	run(t, "add", "Task B")
	run(t, "done", "1")
	out, err := run(t, "list", "done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Task A") {
		t.Errorf("expected done todo in output, got %q", out)
	}
	if strings.Contains(out, "Task B") {
		t.Errorf("expected undone todo excluded, got %q", out)
	}
}

func TestListFilterUndone(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Task A")
	run(t, "add", "Task B")
	run(t, "done", "1")
	out, err := run(t, "list", "undone")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "Task A") {
		t.Errorf("expected done todo excluded, got %q", out)
	}
	if !strings.Contains(out, "Task B") {
		t.Errorf("expected undone todo in output, got %q", out)
	}
}

func TestListFilterToday(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Today task")
	run(t, "add", "Old task")
	run(t, "date", "2", "2026-01-01")
	out, err := run(t, "list", "today")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Today task") {
		t.Errorf("expected today's todo in output, got %q", out)
	}
	if strings.Contains(out, "Old task") {
		t.Errorf("expected old todo excluded, got %q", out)
	}
}

func TestListFilterSpecificDate(t *testing.T) {
	setupTempStorage(t)
	run(t, "add", "Task A")
	run(t, "date", "1", "2026-04-14")
	run(t, "add", "Task B")
	run(t, "date", "2", "2026-04-15")
	out, err := run(t, "list", "2026-04-14")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Task A") {
		t.Errorf("expected Task A in output, got %q", out)
	}
	if strings.Contains(out, "Task B") {
		t.Errorf("expected Task B excluded, got %q", out)
	}
}

func TestListFilterInvalidArg(t *testing.T) {
	setupTempStorage(t)
	_, err := run(t, "list", "not-a-date")
	if err == nil {
		t.Error("expected error for invalid filter argument")
	}
}

func TestResetClearsAllTodos(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Task A")
	run(t, "add", "Task B")
	out, err := run(t, "reset", "--force")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2") || !strings.Contains(out, "cleared") {
		t.Errorf("expected confirmation with count, got %q", out)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("could not read storage file: %v", err)
	}
	var todos []model.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		t.Fatalf("could not parse storage file: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected 0 todos after reset, got %d", len(todos))
	}
}

func TestResetEmptyList(t *testing.T) {
	setupTempStorage(t)
	out, err := run(t, "reset", "--force")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Nothing to reset") {
		t.Errorf("expected 'Nothing to reset.' output, got %q", out)
	}
}
