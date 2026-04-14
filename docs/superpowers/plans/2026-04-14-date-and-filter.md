# Date & List Filtering Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `Date` field to todos (defaulting to today), a `todo date` command to change it, and filtering options to `todo list`.

**Architecture:** `Date` is stored as a plain `YYYY-MM-DD` string on the `Todo` model. All filtering logic lives in `cmd/list.go` as a private helper. Display shows the date column inline. No new packages needed.

**Tech Stack:** Go, Cobra, fatih/color (already in use)

---

### Task 1: Add `Date string` field to the Todo model

**Files:**
- Modify: `internal/model/todo.go`
- Modify: `internal/model/todo_test.go`

- [ ] **Step 1: Write the failing test**

Add to the bottom of `internal/model/todo_test.go`:

```go
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
```

- [ ] **Step 2: Verify it fails**

```
go test ./internal/model/ -run TestTodoDateField
```

Expected: `FAIL` — `unknown field 'Date'`

- [ ] **Step 3: Add the field**

Replace the `Todo` struct in `internal/model/todo.go`:

```go
package model

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	Date      string    `json:"date"`
}
```

- [ ] **Step 4: Verify all tests pass**

```
go test ./...
```

Expected: all `ok`

- [ ] **Step 5: Commit**

```bash
git add internal/model/todo.go internal/model/todo_test.go
git commit -m "feat: add Date field to Todo model"
```

---

### Task 2: Set today's date when adding a todo

**Files:**
- Modify: `cmd/add.go`
- Modify: `cmd/cmd_test.go`

- [ ] **Step 1: Write the failing test**

Add to the bottom of `cmd/cmd_test.go`. Also add `"time"` to the imports:

```go
func TestAddSetsDateToToday(t *testing.T) {
	p := setupTempStorage(t)
	run(t, "add", "Buy milk")
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
	today := time.Now().UTC().Format("2006-01-02")
	if todos[0].Date != today {
		t.Errorf("expected Date %q, got %q", today, todos[0].Date)
	}
}
```

- [ ] **Step 2: Verify it fails**

```
go test ./cmd/ -run TestAddSetsDateToToday
```

Expected: `FAIL` — Date is `""`, not today's date

- [ ] **Step 3: Set Date in the add command**

In `cmd/add.go`, update the todo literal inside `RunE`:

```go
todo := model.Todo{
    ID:        storage.NextID(todos),
    Title:     title,
    Done:      false,
    CreatedAt: time.Now().UTC(),
    Date:      time.Now().UTC().Format("2006-01-02"),
}
```

- [ ] **Step 4: Verify all tests pass**

```
go test ./...
```

Expected: all `ok`

- [ ] **Step 5: Commit**

```bash
git add cmd/add.go cmd/cmd_test.go
git commit -m "feat: set today's date on todo add"
```

---

### Task 3: Show date column in `todo list` display

**Files:**
- Modify: `internal/ui/display.go`
- Modify: `internal/ui/display_test.go`

- [ ] **Step 1: Write the failing test**

Add to the bottom of `internal/ui/display_test.go`:

```go
func TestPrintListShowsDate(t *testing.T) {
	var buf bytes.Buffer
	todos := []model.Todo{
		{ID: 1, Title: "Buy milk", Done: false, Date: "2026-04-14"},
	}
	ui.PrintList(&buf, todos)
	if !strings.Contains(buf.String(), "2026-04-14") {
		t.Errorf("expected date in output, got %q", buf.String())
	}
}
```

- [ ] **Step 2: Verify it fails**

```
go test ./internal/ui/ -run TestPrintListShowsDate
```

Expected: `FAIL` — date not in output

- [ ] **Step 3: Update PrintList to include the date column**

Replace `PrintList` in `internal/ui/display.go`:

```go
func PrintList(w io.Writer, todos []model.Todo) {
	if len(todos) == 0 {
		fmt.Fprintln(w, `No todos yet. Add one with: todo add "your task"`)
		return
	}
	green := color.New(color.FgGreen)
	for _, t := range todos {
		if t.Done {
			green.Fprintf(w, "  %d  ●  %s  %s\n", t.ID, t.Date, t.Title)
		} else {
			fmt.Fprintf(w, "  %d  ○  %s  %s\n", t.ID, t.Date, t.Title)
		}
	}
}
```

- [ ] **Step 4: Verify all tests pass**

```
go test ./...
```

Expected: all `ok`

- [ ] **Step 5: Commit**

```bash
git add internal/ui/display.go internal/ui/display_test.go
git commit -m "feat: show date column in list display"
```

---

### Task 4: Add `todo date <id> <YYYY-MM-DD>` command

**Files:**
- Create: `cmd/date.go`
- Modify: `cmd/root.go`
- Modify: `cmd/cmd_test.go`

- [ ] **Step 1: Write the failing tests**

Add to the bottom of `cmd/cmd_test.go`:

```go
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
```

- [ ] **Step 2: Verify they fail**

```
go test ./cmd/ -run TestDate
```

Expected: `FAIL` — `unknown command "date"`

- [ ] **Step 3: Create `cmd/date.go`**

```go
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newDateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "date <id> <YYYY-MM-DD>",
		Short: "Set the date of a todo",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid id: %s", args[0])
			}
			dateStr := args[1]
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return fmt.Errorf("invalid date format: use YYYY-MM-DD")
			}
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			for i, t := range todos {
				if t.ID == id {
					todos[i].Date = dateStr
					if err := storage.Save(todos); err != nil {
						ui.PrintError(err.Error())
						return err
					}
					ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("✓ Set [%d] date to %s", id, dateStr))
					return nil
				}
			}
			return fmt.Errorf("todo #%d not found", id)
		},
	}
}
```

- [ ] **Step 4: Register the command in `cmd/root.go`**

Add `root.AddCommand(newDateCmd())` to `NewRootCmd`, after `newUndoneCmd()`:

```go
root.AddCommand(newAddCmd())
root.AddCommand(newListCmd())
root.AddCommand(newDoneCmd())
root.AddCommand(newUndoneCmd())
root.AddCommand(newDateCmd())
root.AddCommand(newDeleteCmd())
```

- [ ] **Step 5: Verify all tests pass**

```
go test ./...
```

Expected: all `ok`

- [ ] **Step 6: Commit**

```bash
git add cmd/date.go cmd/root.go cmd/cmd_test.go
git commit -m "feat: add 'date' command to set todo date"
```

---

### Task 5: Add filtering to `todo list`

Supported filters: `done`, `undone`, `today`, `YYYY-MM-DD`. No arg shows all items (including done ones — changed from current behavior).

**Files:**
- Modify: `cmd/list.go`
- Modify: `cmd/cmd_test.go`

- [ ] **Step 1: Write the failing tests**

Add to the bottom of `cmd/cmd_test.go`:

```go
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
```

- [ ] **Step 2: Verify they fail**

```
go test ./cmd/ -run "TestListFilter|TestListShowsDone"
```

Expected: `FAIL` — filters not supported, unknown arg errors

- [ ] **Step 3: Rewrite `cmd/list.go`**

```go
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"todo/internal/model"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [today|done|undone|YYYY-MM-DD]",
		Short: "List todos (optionally filtered)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			filtered, err := filterTodos(todos, args)
			if err != nil {
				return err
			}
			ui.PrintList(cmd.OutOrStdout(), filtered)
			return nil
		},
	}
}

func filterTodos(todos []model.Todo, args []string) ([]model.Todo, error) {
	if len(args) == 0 {
		return todos, nil
	}
	switch args[0] {
	case "done":
		var out []model.Todo
		for _, t := range todos {
			if t.Done {
				out = append(out, t)
			}
		}
		return out, nil
	case "undone":
		var out []model.Todo
		for _, t := range todos {
			if !t.Done {
				out = append(out, t)
			}
		}
		return out, nil
	case "today":
		today := time.Now().UTC().Format("2006-01-02")
		var out []model.Todo
		for _, t := range todos {
			if t.Date == today {
				out = append(out, t)
			}
		}
		return out, nil
	default:
		if _, err := time.Parse("2006-01-02", args[0]); err != nil {
			return nil, fmt.Errorf("unknown filter %q: use done, undone, today, or YYYY-MM-DD", args[0])
		}
		var out []model.Todo
		for _, t := range todos {
			if t.Date == args[0] {
				out = append(out, t)
			}
		}
		return out, nil
	}
}
```

- [ ] **Step 4: Verify all tests pass**

```
go test ./...
```

Expected: all `ok`

- [ ] **Step 5: Rebuild the binary**

```
go build -o todo .
```

- [ ] **Step 6: Commit**

```bash
git add cmd/list.go cmd/cmd_test.go
git commit -m "feat: add date filtering to list command"
```
