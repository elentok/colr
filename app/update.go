package app

import (
	tea "charm.land/bubbletea/v2"

	"github.com/elentok/colr/ui"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		return handleKeyMsg(m, msg.String())

	case ClearToastMsg:
		m.toastMessage = ""

	case SaveHistoryMsg:
		if msg.err != nil {
			m.toastMessage = "Failed to save color to history"
			return m, nil
		}
		m.historyEntries = msg.entries
		m.historyIndex = 0
		m.toastMessage = "Saved color to history"
		return m, nil
	}

	return m, nil
}

func handleKeyMsg(m Model, key string) (tea.Model, tea.Cmd) {
	// ctrl+c always quits.
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	// When help overlay is visible, only dismiss keys are handled.
	if m.showHelp {
		switch key {
		case "?", "esc", "q":
			m.showHelp = false
		}
		return m, nil
	}

	if m.showHistory {
		return handleHistoryKey(m, key), nil
	}

	// Normal key handling.
	switch key {
	case "q":
		return m, tea.Quit

	case "?":
		m.showHelp = true

	case "p":
		m.pendingY = false
		if len(m.historyEntries) == 0 {
			m.toastMessage = "No history yet"
			return m, nil
		}
		if m.historyIndex >= len(m.historyEntries) {
			m.historyIndex = len(m.historyEntries) - 1
		}
		m.showHistory = true

	case "s":
		m.pendingY = false
		return applySaveHistory(m)

	default:
		var cmd tea.Cmd
		m, cmd = handleEditKey(m, key)
		return m, cmd
	}

	return m, nil
}

func handleEditKey(m Model, key string) (Model, tea.Cmd) {
	// Handle pending y-prefix for copy commands.
	if m.pendingY {
		m.pendingY = false
		switch key {
		case "r":
			return applyCopy(m, "rgb")
		case "x":
			return applyCopy(m, "hex")
		case "h":
			return applyCopy(m, "hsl")
		case "n":
			return applyCopy(m, "name")
		case "y":
			return applyCopy(m, "rgb")
		}
		// Unknown key after y — no action, pendingY already cleared.
		return m, nil
	}

	// y starts the copy prefix.
	if key == "y" {
		m.pendingY = true
		return m, nil
	}

	return handleKey(m, key), nil
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

func handleHistoryKey(m Model, key string) Model {
	switch key {
	case "p", "esc", "q":
		m.showHistory = false
	case "j", "down":
		if m.historyIndex < len(m.historyEntries)-1 {
			m.historyIndex++
		}
	case "k", "up":
		if m.historyIndex > 0 {
			m.historyIndex--
		}
	case "g":
		m.historyIndex = 0
	case "G":
		if len(m.historyEntries) > 0 {
			m.historyIndex = len(m.historyEntries) - 1
		}
	case "enter":
		if len(m.historyEntries) > 0 {
			m = applyHistoryEntry(m, m.historyEntries[m.historyIndex])
			m.toastMessage = "Loaded color from history"
		}
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
