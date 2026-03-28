package app

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/elentok/colr/clipboard"
	"github.com/elentok/colr/color"
)

// applyCopy writes the requested format to the clipboard and sets a toast.
// format must be "rgb", "hex", "hsl", "name", or "over".
func applyCopy(m Model, format string) (Model, tea.Cmd) {
	var text, label string
	switch format {
	case "rgb":
		text = color.FormatRGB(m.currentColor)
		label = "RGB"
	case "hex":
		text = color.FormatHEX(m.currentColor)
		label = "HEX"
	case "hsl":
		text = color.FormatHSL(m.currentColor)
		label = "HSL"
	case "name":
		text = strings.TrimPrefix(color.NearestNamedColor(m.currentColor), "~")
		label = "Name"
	case "over":
		bg := previewBackgroundColor(m.previewDarkBG)
		text = color.FormatHEX(color.CompositeOver(m.currentColor, bg))
		label = "Over-bg"
	}

	err := clipboard.Write(text)
	if err != nil {
		m.toastMessage = fmt.Sprintf("Failed to copy %s to clipboard", label)
		return m, scheduleToast(2 * time.Second)
	}

	m.toastMessage = fmt.Sprintf("Copied %s to clipboard", label)
	return m, scheduleToast(1 * time.Second)
}

// scheduleToast returns a Cmd that sends ClearToastMsg after duration.
func scheduleToast(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return ClearToastMsg{}
	})
}
