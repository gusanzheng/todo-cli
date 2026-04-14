package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"todo/internal/model"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <title>",
		Short: "Add a new todo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			if title == "" {
				return fmt.Errorf("title cannot be empty")
			}
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			todo := model.Todo{
				ID:        storage.NextID(todos),
				Title:     title,
				Done:      false,
				CreatedAt: time.Now().UTC(),
			}
			todos = append(todos, todo)
			if err := storage.Save(todos); err != nil {
				ui.PrintError(err.Error())
				return err
			}
			ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("✓ Added [%d] %s", todo.ID, todo.Title))
			return nil
		},
	}
}
