package cmd

import (
	"github.com/spf13/cobra"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		RunE: func(cmd *cobra.Command, args []string) error {
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			ui.PrintList(cmd.OutOrStdout(), todos)
			return nil
		},
	}
}
