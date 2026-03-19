package app

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/elentok/colr/color"
	"github.com/elentok/colr/history"
)

// ClearToastMsg is sent when a toast timer expires.
type ClearToastMsg struct{}

// Model is the Bubble Tea application model.
type Model struct {
	originalClip   string
	originalColor  color.Color
	currentColor   color.Color
	editMode       EditMode
	selectedField  int
	lastHue        float64 // preserved hue for grayscale stability
	toastMessage   string
	toastExpiry    time.Time
	width          int
	height         int
	showHelp       bool
	showHistory    bool
	pendingY       bool
	historyEntries []history.Entry
	historyIndex   int
}

// NewModel creates a new Model from parsed clipboard input.
func NewModel(clipText string, c color.Color, historyEntries []history.Entry) Model {
	hsv := color.RGBToHSV(c)
	return Model{
		originalClip:   clipText,
		originalColor:  c,
		currentColor:   c,
		editMode:       ModeHSV,
		selectedField:  0,
		lastHue:        hsv.H,
		width:          80,
		height:         24,
		historyEntries: append([]history.Entry(nil), historyEntries...),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) WithToast(msg string) Model {
	m.toastMessage = msg
	return m
}

func (m Model) HistoryEntriesForSave() []history.Entry {
	return history.Record(m.historyEntries, m.originalClip, m.currentColor)
}

func applyHistoryEntry(m Model, entry history.Entry) Model {
	hsv := color.RGBToHSV(entry.Color)
	m.originalClip = entry.Original
	m.originalColor = entry.Color
	m.currentColor = entry.Color
	m.lastHue = hsv.H
	m.showHistory = false
	m.pendingY = false
	return m
}
