package cmd_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
