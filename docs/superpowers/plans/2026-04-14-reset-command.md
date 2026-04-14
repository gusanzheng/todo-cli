# Reset Command Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `todo reset` command that clears all todos, with interactive confirmation and a `--force` flag.

**Architecture:** New `cmd/reset.go` follows the exact pattern of existing commands — load from storage, guard, optionally prompt, save empty slice, print result. The `--force` flag bypasses stdin for scripting and tests.

**Tech Stack:** Go, Cobra, fatih/color, standard `bufio`/`os`/`strings` for stdin reading.

---

### Task 1: Write the two failing tests

**Files:**
- Modify: `cmd/cmd_test.go`

- [ ] **Step 1: Append both tests to `cmd/cmd_test.go`**

Add these two functions at the end of the file:

```go
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
	data, _ := os.ReadFile(p)
	var todos []model.Todo
	json.Unmarshal(data, &todos)
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/ -run "TestResetClearsAllTodos|TestResetEmptyList" -v
```

Expected: both tests FAIL with `unknown command "reset"` or similar.

---

### Task 2: Implement `cmd/reset.go`

**Files:**
- Create: `cmd/reset.go`

- [ ] **Step 3: Create `cmd/reset.go` with this content**

```go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"todo/internal/model"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newResetCmd() *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use:   "reset",
		Short: "Clear all todos",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			n := len(todos)
			if n == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to reset.")
				return nil
			}
			if !force {
				fmt.Fprintf(os.Stderr, "Reset all %d todos? [y/N]: ", n)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				input := strings.TrimSpace(scanner.Text())
				if !strings.EqualFold(input, "y") && !strings.EqualFold(input, "yes") {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}
			if err := storage.Save([]model.Todo{}); err != nil {
				ui.PrintError(err.Error())
				return err
			}
			ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("✓ Reset: all %d todos cleared.", n))
			return nil
		},
	}
	c.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	return c
}
```

---

### Task 3: Register the command and verify

**Files:**
- Modify: `cmd/root.go`

- [ ] **Step 4: Register `newResetCmd()` in `root.go`**

In `NewRootCmd()`, add after the existing `AddCommand` calls:

```go
root.AddCommand(newResetCmd())
```

The full `NewRootCmd` body should look like:

```go
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "todo",
		Short: "A minimal CLI todo tracker",
	}
	root.AddCommand(newAddCmd())
	root.AddCommand(newListCmd())
	root.AddCommand(newDoneCmd())
	root.AddCommand(newUndoneCmd())
	root.AddCommand(newDateCmd())
	root.AddCommand(newDeleteCmd())
	root.AddCommand(newResetCmd())
	return root
}
```

- [ ] **Step 5: Run the two new tests to verify they pass**

```bash
go test ./cmd/ -run "TestResetClearsAllTodos|TestResetEmptyList" -v
```

Expected: both PASS.

- [ ] **Step 6: Run the full test suite to check for regressions**

```bash
go test ./...
```

Expected: all tests PASS.

- [ ] **Step 7: Build to verify no compile errors**

```bash
go build ./...
```

Expected: exits 0 with no output.

- [ ] **Step 8: Commit**

```bash
git add cmd/reset.go cmd/root.go cmd/cmd_test.go
git commit -m "feat: add reset command with confirmation prompt and --force flag"
```
