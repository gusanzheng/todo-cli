# Changelog

## 2026-04-14 — Date support & list filtering
_git: `e91d6b2` → `8b84a3a`_

### Added
- **`todo date <id> <YYYY-MM-DD>`** — change the date on any todo item
- **`todo list today`** — show only items dated today
- **`todo list <YYYY-MM-DD>`** — show items for a specific date
- **`todo list done`** — show only completed items
- **`todo list undone`** — show only pending items
- New todos automatically get today's date (UTC) set on creation
- Date column now shown in `todo list` output between the status symbol and title

### Changed
- **`todo list`** (no args) now shows **all** items including done ones (previously only showed pending)

---

## Initial release — Core todo commands
_git: `fd48cbf` → `00c1234`_

### Added
- **`todo add "<title>"`** — create a new todo item
- **`todo list`** — list all pending todos
- **`todo done <id>`** — mark a todo as complete
- **`todo undone <id>`** — unmark a completed todo
- **`todo delete <id>`** — remove a todo permanently
- JSON storage at `~/.todos.json`; override path with `TODO_FILE` env var
- Color output: green for done items, plain for pending
