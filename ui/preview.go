package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

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
	fgANSI := lipgloss.Color("255")
	if luminance > 0.5 {
		fgSuggestion = "black"
		fgANSI = lipgloss.Color("0")
	}

	hasAlpha := c.A < 0.999

	var rows []string
	previewRows := height - 1 // leave room for the fg hint line
	if previewRows < 1 {
		previewRows = 1
	}

	colorStyle := lipgloss.NewStyle().Background(lipgloss.Color(hexColor))
	checkerDarkStyle := lipgloss.NewStyle().Background(lipgloss.Color("236"))
	checkerLightStyle := lipgloss.NewStyle().Background(lipgloss.Color("250"))

	blockW := 2
	line := strings.Repeat(" ", width)

	for i := range previewRows {
		if hasAlpha {
			// Per-cell compositing: alternate color/checker blocks within each row.
			var cells []string
			col := 0
			for col < width {
				end := col + blockW
				if end > width {
					end = width
				}
				chunk := strings.Repeat(" ", end-col)
				// Alternate between color and checker based on position + row parity.
				block := (col/blockW + i) % 2
				if block == 0 {
					cells = append(cells, colorStyle.Render(chunk))
				} else {
					if (i/2)%2 == 0 {
						cells = append(cells, checkerDarkStyle.Render(chunk))
					} else {
						cells = append(cells, checkerLightStyle.Render(chunk))
					}
				}
				col = end
			}
			rows = append(rows, strings.Join(cells, ""))
		} else {
			rows = append(rows, colorStyle.Render(line))
		}
	}

	// Foreground hint line at bottom: styled with actual fg/bg colors.
	fgHint := fmt.Sprintf("Foreground: %s", fgSuggestion)
	if hasAlpha {
		fgHint += " (transparent)"
	}
	hintStyle := lipgloss.NewStyle().
		Foreground(fgANSI).
		Background(lipgloss.Color(hexColor)).
		Width(width)
	rows = append(rows, hintStyle.Render(fgHint))

	return strings.Join(rows, "\n")
}
