package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done <id>",
		Short: "Mark a todo as done",
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
					todos[i].Done = true
					if err := storage.Save(todos); err != nil {
						ui.PrintError(err.Error())
						return err
					}
					ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("✓ Marked [%d] done", id))
					return nil
				}
			}
			return fmt.Errorf("todo #%d not found", id)
		},
	}
}
