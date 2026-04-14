package ui_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"todo/internal/model"
	"todo/internal/ui"
)

func init() {
	color.NoColor = true
}

func TestPrintListEmpty(t *testing.T) {
	var buf bytes.Buffer
	ui.PrintList(&buf, []model.Todo{})
	if !strings.Contains(buf.String(), "No todos yet") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestPrintListPending(t *testing.T) {
	var buf bytes.Buffer
	todos := []model.Todo{
		{ID: 1, Title: "Buy milk", Done: false, CreatedAt: time.Now()},
	}
	ui.PrintList(&buf, todos)
	out := buf.String()
	if !strings.Contains(out, "1") || !strings.Contains(out, "Buy milk") || !strings.Contains(out, "○") {
		t.Errorf("unexpected output for pending todo: %q", out)
	}
}

func TestPrintListDone(t *testing.T) {
	var buf bytes.Buffer
	todos := []model.Todo{
		{ID: 2, Title: "Write tests", Done: true, CreatedAt: time.Now()},
	}
	ui.PrintList(&buf, todos)
	out := buf.String()
	if !strings.Contains(out, "2") || !strings.Contains(out, "Write tests") || !strings.Contains(out, "●") {
		t.Errorf("unexpected output for done todo: %q", out)
	}
}

func TestPrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	ui.PrintSuccess(&buf, "✓ Added [1] Buy milk")
	if !strings.Contains(buf.String(), "✓ Added [1] Buy milk") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestPrintListShowsDate(t *testing.T) {
	var buf bytes.Buffer
	todos := []model.Todo{
		{ID: 1, Title: "Buy milk", Done: false, Date: "2026-04-14"},
	}
	ui.PrintList(&buf, todos)
	if !strings.Contains(buf.String(), "2026-04-14") {
		t.Errorf("expected date in output, got %q", buf.String())
	}
}
