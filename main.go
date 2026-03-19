package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/elentok/colr/app"
	"github.com/elentok/colr/clipboard"
	"github.com/elentok/colr/color"
	"github.com/elentok/colr/history"
)

var errNoColorAvailable = errors.New("input does not contain a valid CSS color and history is empty")

func resolveInput(args []string, readClipboard func() (string, error)) (string, error) {
	if len(args) > 1 {
		return strings.Join(args[1:], " "), nil
	}

	return readClipboard()
}

func resolveStartupColor(inputText string, historyEntries []history.Entry) (string, color.Color, string, error) {
	if c, err := color.FindFirst(inputText); err == nil {
		return inputText, c, "", nil
	}

	if len(historyEntries) > 0 {
		return historyEntries[0].Original, historyEntries[0].Color, "Input had no color; loaded last history color", nil
	}

	return "", color.Color{}, "", errNoColorAvailable
}

func main() {
	inputText, err := resolveInput(os.Args, clipboard.Read)
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: failed to read clipboard")
		os.Exit(1)
	}

	historyEntries, err := history.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: failed to load history")
		os.Exit(1)
	}

	startupText, c, startupToast, err := resolveStartupColor(inputText, historyEntries)
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: input does not contain a valid CSS color")
		fmt.Fprintln(os.Stderr, "colr: no history fallback is available")
		fmt.Fprintln(os.Stderr, "Supported formats: RGB, RGBA, HEX, HSL")
		os.Exit(1)
	}

	model := app.NewModel(startupText, c, historyEntries).WithToast(startupToast)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr:", err)
		os.Exit(1)
	}

	m, ok := finalModel.(app.Model)
	if !ok {
		fmt.Fprintln(os.Stderr, "colr: unexpected final model type")
		os.Exit(1)
	}

	if err := history.Save(m.HistoryEntriesForSave()); err != nil {
		fmt.Fprintln(os.Stderr, "colr: failed to save history")
		os.Exit(1)
	}
}
