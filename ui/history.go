package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/elentok/colr/history"
)

var historyTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("39"))

var historyMetaStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("245"))

var historyBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("39")).
	Padding(1, 2)

var historySelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("39"))

func RenderHistory(width, height int, entries []history.Entry, selected int) string {
	if len(entries) == 0 {
		box := historyBorderStyle.Render(historyTitleStyle.Render("history") + "\n\nNo history yet.")
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
	}

	if selected < 0 {
		selected = 0
	}
	if selected >= len(entries) {
		selected = len(entries) - 1
	}

	visibleRows := height - 8
	if visibleRows < 4 {
		visibleRows = 4
	}
	if visibleRows > len(entries) {
		visibleRows = len(entries)
	}

	start := selected - visibleRows/2
	if start < 0 {
		start = 0
	}
	if maxStart := len(entries) - visibleRows; start > maxStart {
		start = maxStart
	}
	end := start + visibleRows
	if end > len(entries) {
		end = len(entries)
	}

	contentWidth := width - 10
	if contentWidth < 48 {
		contentWidth = 48
	}

	swatchW := 6
	rgbW := 22
	hexW := 12
	nameW := 16
	metaPadding := 14

	var lines []string
	for i := start; i < end; i++ {
		entry := entries[i]
		indicator := "  "
		if i == selected {
			indicator = "› "
		}

		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(entry.HEX)).
			Width(swatchW).
			Render(strings.Repeat(" ", swatchW))

		rgbCell := OutputValueStyle.Width(rgbW).Render(entry.RGB)
		hexCell := OutputValueStyle.Width(hexW).Render(entry.HEX)
		nameCell := historyMetaStyle.Width(nameW).Render(entry.Name)
		row := fmt.Sprintf("%s%s  %s  %s  %s", indicator, swatch, rgbCell, hexCell, nameCell)
		if i == selected {
			row = historySelectedStyle.Width(contentWidth - metaPadding).Render(row)
		}
		lines = append(lines, row)
	}

	title := historyTitleStyle.Render("history")
	meta := historyMetaStyle.Render(fmt.Sprintf("%d colors  •  enter load  •  esc close", len(entries)))
	body := title + "\n" + meta + "\n\n" + strings.Join(lines, "\n")
	box := historyBorderStyle.Width(contentWidth + 4).Render(body)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}
