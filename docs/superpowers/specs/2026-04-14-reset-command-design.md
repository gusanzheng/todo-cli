# Design: `todo reset` command

**Date:** 2026-04-14

## Overview

Add a `todo reset` command that clears all todos from storage. Because the operation is destructive and irreversible, it requires explicit confirmation before proceeding.

## Command signature

```
todo reset [--force]
```

- `--force` / `-f`: skip the interactive prompt (for scripting/non-interactive use)

## Behavior

1. **Empty list guard:** If there are no todos, print `Nothing to reset.` to stdout and exit 0. No prompt is shown.
2. **Without `--force`:** Print `Reset all N todos? [y/N]: ` to stderr. Read one line from stdin. If the input is `y` or `yes` (case-insensitive), proceed. Anything else (including bare Enter) prints `Aborted.` to stdout and exits 0.
3. **With `--force`:** Skip the prompt and proceed immediately.
4. **On confirm:** Save an empty slice via `storage.Save([]model.Todo{})`, then print `✓ Reset: all N todos cleared.` to stdout.

## Implementation

- New file `cmd/reset.go`, following the pattern of `cmd/delete.go`
- Registered in `cmd/root.go` with `root.AddCommand(newResetCmd())`
- Success output goes to `cmd.OutOrStdout()` (enables test capture)
- Confirmation prompt and `Aborted.` go to stderr / stdout respectively, consistent with other commands
- No new storage primitives needed; `storage.Save([]model.Todo{})` is sufficient

## Testing

Added to `cmd/cmd_test.go`:

| Test | Description |
|------|-------------|
| `TestResetClearsAllTodos` | `--force` clears all todos and outputs confirmation with count |
| `TestResetEmptyList` | `--force` on empty list prints "Nothing to reset." |

Interactive stdin tests are out of scope for this test suite (no tty in test environment). The `--force` path fully covers the core logic.

## Files changed

| File | Change |
|------|--------|
| `cmd/reset.go` | New — command definition |
| `cmd/root.go` | Register `newResetCmd()` |
| `cmd/cmd_test.go` | Add two tests |
