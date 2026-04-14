package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"todo/internal/model"
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
			if _, err := time.Parse(model.DateFormat, dateStr); err != nil {
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
