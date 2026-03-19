package ui

import (
	"charm.land/lipgloss/v2"
)

var helpTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("39"))

var helpSectionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("245")).
	Bold(true)

var helpKeyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("39")).
	Width(18)

var helpDescStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255"))

var helpBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("39")).
	Padding(1, 2)

// RenderHelp returns a centered help overlay box.
func RenderHelp(width, height int) string {
	rows := []struct{ key, desc string }{} // typed below per section
	_ = rows

	content := helpTitleStyle.Render("colr help") + "\n\n"

	content += helpSectionStyle.Render("Movement") + "\n"
	content += helpRow("j/k  ↑/↓", "move between fields")
	content += helpRow("h/l  ←/→", "adjust field value")
	content += helpRow("H/L", "adjust by larger step")
	content += helpRow("g/G", "first / last field")
	content += "\n"

	content += helpSectionStyle.Render("Modes") + "\n"
	content += helpRow("tab", "toggle HSV / RGB")
	content += helpRow("1", "HSV mode")
	content += helpRow("2", "RGB mode")
	content += "\n"

	content += helpSectionStyle.Render("Copy") + "\n"
	content += helpRow("yr / yy", "copy RGB")
	content += helpRow("yx", "copy HEX")
	content += helpRow("yh", "copy HSL")
	content += helpRow("yn", "copy CSS name")
	content += "\n"

	content += helpSectionStyle.Render("History") + "\n"
	content += helpRow("p", "open color history")
	content += helpRow("j/k  ↑/↓", "move within history")
	content += helpRow("enter", "load selected color")
	content += helpRow("Esc / q / p", "close history")
	content += "\n"

	content += helpSectionStyle.Render("Other") + "\n"
	content += helpRow("R", "reset color")
	content += helpRow("?  Esc  q", "close this help")
	content += helpRow("ctrl+c", "quit")

	box := helpBorderStyle.Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

func helpRow(key, desc string) string {
	return "  " + helpKeyStyle.Render(key) + helpDescStyle.Render(desc) + "\n"
}
