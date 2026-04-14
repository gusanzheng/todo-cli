package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"todo/internal/model"
	"todo/internal/storage"
	"todo/internal/ui"
)

func newResetCmd() *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use:   "reset",
		Short: "Clear all todos",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			todos, err := storage.Load()
			if err != nil {
				ui.PrintError(err.Error())
				return err
			}
			n := len(todos)
			if n == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to reset.")
				return nil
			}
			if !force {
				fmt.Fprintf(os.Stderr, "Reset all %d todos? [y/N]: ", n)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				input := strings.TrimSpace(scanner.Text())
				if !strings.EqualFold(input, "y") && !strings.EqualFold(input, "yes") {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}
			if err := storage.Save([]model.Todo{}); err != nil {
				ui.PrintError(err.Error())
				return err
			}
			ui.PrintSuccess(cmd.OutOrStdout(), fmt.Sprintf("✓ Reset: all %d todos cleared.", n))
			return nil
		},
	}
	c.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	return c
}
