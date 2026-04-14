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
		today := time.Now().UTC().Format(model.DateFormat)
		var out []model.Todo
		for _, t := range todos {
			if t.Date == today {
				out = append(out, t)
			}
		}
		return out, nil
	default:
		if _, err := time.Parse(model.DateFormat, args[0]); err != nil {
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
