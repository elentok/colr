package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/elentok/colr/app"
	"github.com/elentok/colr/clipboard"
	"github.com/elentok/colr/color"
)

func resolveInput(args []string, readClipboard func() (string, error)) (string, error) {
	if len(args) > 1 {
		return strings.Join(args[1:], " "), nil
	}

	return readClipboard()
}

func main() {
	inputText, err := resolveInput(os.Args, clipboard.Read)
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: failed to read clipboard")
		os.Exit(1)
	}

	c, err := color.FindFirst(inputText)
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: input does not contain a valid CSS color")
		fmt.Fprintln(os.Stderr, "Supported formats: RGB, RGBA, HEX, HSL")
		os.Exit(1)
	}

	model := app.NewModel(inputText, c)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "colr:", err)
		os.Exit(1)
	}
}
