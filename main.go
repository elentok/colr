package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/elentok/colr/app"
	"github.com/elentok/colr/clipboard"
	"github.com/elentok/colr/color"
)

func main() {
	clipText, err := clipboard.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: failed to read clipboard")
		os.Exit(1)
	}

	c, err := color.Parse(clipText)
	if err != nil {
		fmt.Fprintln(os.Stderr, "colr: clipboard does not contain a valid CSS color")
		fmt.Fprintln(os.Stderr, "Supported formats: RGB, RGBA, HEX, HSL")
		os.Exit(1)
	}

	model := app.NewModel(clipText, c)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "colr:", err)
		os.Exit(1)
	}
}
