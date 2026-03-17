package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/elentok/colr/color"
)

// ClearToastMsg is sent when a toast timer expires.
type ClearToastMsg struct{}

// Model is the Bubble Tea application model.
type Model struct {
	originalClip  string
	originalColor color.Color
	currentColor  color.Color
	editMode      EditMode
	selectedField int
	toastMessage  string
	toastExpiry   time.Time
	width         int
	height        int
	showHelp      bool
}

// NewModel creates a new Model from parsed clipboard input.
func NewModel(clipText string, c color.Color) Model {
	return Model{
		originalClip:  clipText,
		originalColor: c,
		currentColor:  c,
		editMode:      ModeHSV,
		selectedField: 0,
		width:         80,
		height:        24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
