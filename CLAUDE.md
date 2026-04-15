# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...                          # build
go run . <command>                      # run without installing
go test ./...                           # run all tests
go test ./cmd/ -run TestAddCreatesAndOutputs  # run a single test
```

Install and use:
```bash
go install .
todo add "task title"
todo list                   # all items (including done), with date
todo list today             # today's items only
todo list 2026-04-14        # items for a specific date
todo list done              # done items only
todo list undone            # pending items only
todo done <id>
todo undone <id>
todo date <id> 2026-04-14   # change the date on an item
todo delete <id>
```

## Feature Completion Checklist

When a feature is finished, you must:

1. Update `docs/CHANGELOG.md` with the new feature or change.
2. Build and install the binary:
   ```bash
   go build ./... && go install .
   ```

## Architecture

This is a Go CLI todo tracker (module `todo`) built with [Cobra](https://github.com/spf13/cobra) and [fatih/color](https://github.com/fatih/color).

**Layer structure:**

- `cmd/` — Cobra command definitions (`add`, `list`, `done`, `undone`, `date`, `delete`) wired together in `root.go`. Each command calls `storage.Load`, mutates the slice, calls `storage.Save`, then prints via `ui`. Integration tests live here as `cmd_test.go` (external `cmd_test` package).
- `internal/storage/` — JSON persistence. Reads/writes `~/.todos.json` by default; the `TODO_FILE` env var overrides the path. Tests set `TODO_FILE` to a temp file via `t.Setenv`.
- `internal/model/` — `Todo` struct (ID, Title, Done, CreatedAt, Date). `model.DateFormat = "2006-01-02"` is the single source of truth for the date format used across all commands.
- `internal/ui/` — Terminal output helpers. `PrintList` writes to an `io.Writer` (enabling test capture); `PrintError` always writes to stderr. Tests disable color globally with `color.NoColor = true` in `init()`.

**ID generation:** `storage.NextID` returns `max(existing IDs) + 1`; IDs are never reused after deletion.

**Testing pattern:** Commands accept `io.Writer` output via `cmd.OutOrStdout()`, so tests capture output by calling `root.SetOut(buf)` rather than capturing stdout.
