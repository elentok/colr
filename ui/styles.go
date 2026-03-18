package ui

import "charm.land/lipgloss/v2"

var (
	// Borders
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	// Header
	HeaderLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Width(12)

	HeaderValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	ToastOKStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("40"))

	ToastErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	// Editor
	FieldLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(14)

	FieldValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	SelectedFieldStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("39")).
				Bold(true)

	ModeTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	ModeHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true)

	// Outputs
	OutputLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Width(6)

	OutputValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	OutputKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	// Slider
	SliderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	FooterKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	FooterSepStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("237"))
)
