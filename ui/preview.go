package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/elentok/colr/color"
)

// RenderPreview renders the color preview block.
func RenderPreview(c color.Color, width, height int) string {
	hexColor := fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)

	// Compute foreground suggestion via relative luminance.
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0
	luminance := 0.2126*r + 0.7152*g + 0.0722*b
	fgSuggestion := "white"
	if luminance > 0.5 {
		fgSuggestion = "black"
	}

	// Build the preview rows.
	// If alpha < 1.0, alternate rows between the color and a checkerboard hint.
	hasAlpha := c.A < 0.999

	var rows []string
	previewRows := height - 2 // leave room for the fg hint line
	if previewRows < 1 {
		previewRows = 1
	}

	colorStyle := lipgloss.NewStyle().Background(lipgloss.Color(hexColor))
	// A neutral dark/light checker background to contrast with the color.
	checkerDarkStyle := lipgloss.NewStyle().Background(lipgloss.Color("236"))
	checkerLightStyle := lipgloss.NewStyle().Background(lipgloss.Color("250"))

	line := strings.Repeat(" ", width)

	for i := range previewRows {
		if hasAlpha && i%2 == 1 {
			// Alternate row: split into checker blocks to hint at transparency.
			blockW := 2
			var cells []string
			col := 0
			for col < width {
				end := col + blockW
				if end > width {
					end = width
				}
				chunk := strings.Repeat(" ", end-col)
				block := (col/blockW + i) % 2
				if block == 0 {
					cells = append(cells, checkerDarkStyle.Render(chunk))
				} else {
					cells = append(cells, checkerLightStyle.Render(chunk))
				}
				col = end
			}
			rows = append(rows, strings.Join(cells, ""))
		} else {
			rows = append(rows, colorStyle.Render(line))
		}
	}

	// Foreground hint line at bottom of preview.
	fgHint := fmt.Sprintf("Foreground: %s", fgSuggestion)
	if hasAlpha {
		fgHint += "  (transparent)"
	}
	rows = append(rows, fgHint)

	return strings.Join(rows, "\n")
}
