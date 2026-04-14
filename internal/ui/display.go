package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"todo/internal/model"
)

func PrintList(w io.Writer, todos []model.Todo) {
	if len(todos) == 0 {
		fmt.Fprintln(w, `No todos yet. Add one with: todo add "your task"`)
		return
	}
	green := color.New(color.FgGreen)
	for _, t := range todos {
		if t.Done {
			green.Fprintf(w, "  %d  ●  %s  %s\n", t.ID, t.Date, t.Title)
		} else {
			fmt.Fprintf(w, "  %d  ○  %s  %s\n", t.ID, t.Date, t.Title)
		}
	}
}

func PrintSuccess(w io.Writer, msg string) {
	color.New(color.FgGreen).Fprintln(w, msg)
}

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
}
