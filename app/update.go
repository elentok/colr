package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/elentok/colr/ui"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			m = handleKey(m, msg.String())
		}

	case ClearToastMsg:
		m.toastMessage = ""
	}

	return m, nil
}

func handleKey(m Model, key string) Model {
	switch key {
	// Movement
	case "j", "down":
		m.selectedField = clampField(m.selectedField + 1)
	case "k", "up":
		m.selectedField = clampField(m.selectedField - 1)
	case "g":
		m.selectedField = 0
	case "G":
		m.selectedField = ui.FieldCount - 1

	// Adjust — small step
	case "h", "left", "-":
		m = applyAdjust(m, -1, false)
	case "l", "right", "*":
		m = applyAdjust(m, +1, false)

	// Adjust — large step
	case "H":
		m = applyAdjust(m, -1, true)
	case "L":
		m = applyAdjust(m, +1, true)

	// Mode switching
	case "tab":
		if m.editMode == ModeHSV {
			m.editMode = ModeRGB
		} else {
			m.editMode = ModeHSV
		}
	case "1":
		m.editMode = ModeHSV
	case "2":
		m.editMode = ModeRGB

	// Reset
	case "R":
		m = applyReset(m)
	}

	return m
}

func clampField(f int) int {
	if f < 0 {
		return 0
	}
	if f >= ui.FieldCount {
		return ui.FieldCount - 1
	}
	return f
}
