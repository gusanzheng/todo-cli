# todo-cli

A command-line todo tracker with date support.

## Install

```bash
go install .
```

## Usage

```bash
todo add "task title"        # add a new item (date defaults to today)
todo list                    # all items, with date
todo list today              # today's items only
todo list 2026-04-14         # items for a specific date
todo list done               # done items only
todo list undone             # pending items only
todo done <id>               # mark an item done
todo undone <id>             # mark an item undone
todo date <id> 2026-04-14    # change the date on an item
todo delete <id>             # delete an item
```

Data is stored in `~/.todos.json`. Override with the `TODO_FILE` environment variable.
