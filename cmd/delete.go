package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a todo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid id: %s", args[0])
			}
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			for i, t := range todos {
				if t.ID == id {
					title := t.Title
					todos = append(todos[:i], todos[i+1:]...)
					if err := storage.Save(todos); err != nil {
						ui.PrintError(err.Error())
						return err
					}
					ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("✓ Deleted [%d] %s", id, title))
					return nil
				}
			}
			return fmt.Errorf("todo #%d not found", id)
		},
	}
}
