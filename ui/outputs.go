package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/elentok/colr/color"
)

// RenderOutputs renders the RGB/HEX/HSL output rows.
func RenderOutputs(c color.Color, width int) string {
	rgb := color.FormatRGB(c)
	hex := color.FormatHEX(c)
	hsl := color.FormatHSL(c)

	name := color.NearestNamedColor(c)

	rows := []struct {
		label string
		value string
		key   string
	}{
		{"RGB", rgb, "[yr]"},
		{"HEX", hex, "[yx]"},
		{"HSL", hsl, "[yh]"},
		{"NAME", name, "[yn]"},
	}

	// Line structure: "  " + label(6) + "  " + value(fill) + "  " + key(4)
	// Total = 2 + 6 + 2 + valueW + 2 + 4 = 16 + valueW = width
	const labelW = 6
	const keyW = 4
	const padding = 2 + 2 + 2 // leading + label/value gap + value/key gap
	valueW := width - labelW - keyW - padding
	if valueW < 4 {
		valueW = 4
	}

	var lines []string
	for _, r := range rows {
		label := OutputLabelStyle.Width(labelW).Render(r.label)
		value := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Width(valueW).
			Render(r.value)
		key := OutputKeyStyle.Render(r.key)
		lines = append(lines, "  "+label+"  "+value+"  "+key)
	}
	return strings.Join(lines, "\n")
}
