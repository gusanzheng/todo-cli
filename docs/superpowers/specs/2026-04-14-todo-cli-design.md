# TODO CLI — Design Spec

**Date:** 2026-04-14  
**Status:** Approved  
**Stack:** Go, cobra, fatih/color  
**Storage:** Local file (`~/.todos.json`)

---

## Overview

A minimal CLI tool for developers to track TODO items locally. Four commands, colored terminal output, no server or sync required.

---

## Architecture

Option B — structured internal packages. Clean layer separation enables parallel development and isolated testing.

```
todo/
├── main.go                  # Entry point, cobra root command
├── cmd/
│   ├── add.go               # `todo add "buy milk"`
│   ├── list.go              # `todo list`
│   ├── done.go              # `todo done <id>`
│   └── delete.go            # `todo delete <id>`
├── internal/
│   ├── model/
│   │   └── todo.go          # Todo struct: ID, Title, Done, CreatedAt
│   ├── storage/
│   │   └── file.go          # Read/write ~/.todos.json
│   └── ui/
│       └── display.go       # Colored output with fatih/color
└── go.mod
```

**Data flow:**
```
CLI command → cmd/*.go → internal/storage → ~/.todos.json
                      ↘ internal/ui → terminal output
```

---

## Data Model

Storage file: `~/.todos.json`

```json
[
  { "id": 1, "title": "Buy milk", "done": false, "created_at": "2026-04-14T10:00:00Z" }
]
```

- IDs are stable, auto-incrementing integers — never reused after deletion
- `created_at` stored as RFC3339 UTC

---

## Commands

| Command | Usage | Success output |
|---|---|---|
| `todo add` | `todo add "Buy milk"` | `✓ Added [1] Buy milk` |
| `todo list` | `todo list` | Numbered list, color-coded by status |
| `todo done` | `todo done 1` | `✓ Marked [1] done` |
| `todo delete` | `todo delete 1` | `✓ Deleted [1] Buy milk` |

### `list` display

```
  1  ○  Buy milk
  2  ●  Write tests
  3  ○  Review PR
```

- `○` = pending (white)
- `●` = done (green + strikethrough)

### Error messages

| Scenario | Message |
|---|---|
| ID not found | `error: todo #99 not found` |
| Empty title | `error: title cannot be empty` |
| Corrupt storage | `error: failed to read todos, check ~/.todos.json` |

---

## Testing Strategy

| Layer | Approach |
|---|---|
| `internal/model` | Unit tests — struct fields, zero values |
| `internal/storage` | Unit tests — real temp files via `os.TempDir()`, covers read/write/corrupt/missing |
| `internal/ui` | Unit tests — `color.NoColor = true` to strip ANSI codes for assertions |
| `cmd/` | Integration tests via cobra test helpers, temp storage path injected via env var |

Coverage target: 80%+ on `storage` and `model`. UI display is lower priority.

No mocks — tests hit real temp files.

---

## Agent Team Roles

| Agent | Owns |
|---|---|
| PM | Requirements, spec, acceptance criteria |
| RD | Implementation (`cmd/`, `internal/`) |
| Reviewer | Code review, architecture feedback |
| Test | Test suite (`*_test.go`), coverage |
