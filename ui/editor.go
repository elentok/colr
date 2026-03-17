package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/elentok/colr/color"
)

// RenderEditor renders the editor pane showing adjustable fields.
func RenderEditor(c color.Color, mode EditMode, selectedField int, width int) string {
	var lines []string

	if mode == ModeHSV {
		lines = append(lines, ModeTitleStyle.Render("Edit Mode: HSV"))
		lines = append(lines, "")
		hsv := color.RGBToHSV(c)
		opacity := int(math.Round(c.A * 100))
		fields := []struct {
			label string
			value string
			idx   int
		}{
			{"Hue", fmt.Sprintf("%d°", int(math.Round(hsv.H))), FieldHue},
			{"Saturation", fmt.Sprintf("%d%%", int(math.Round(hsv.S*100))), FieldSaturation},
			{"Value", fmt.Sprintf("%d%%", int(math.Round(hsv.V*100))), FieldValue},
			{"Opacity", fmt.Sprintf("%d%%", opacity), FieldOpacity},
		}
		for _, f := range fields {
			lines = append(lines, renderField(f.label, f.value, f.idx == selectedField, width))
		}
		lines = append(lines, "")
		lines = append(lines, ModeHintStyle.Render("Tab: switch to RGB"))
	} else {
		lines = append(lines, ModeTitleStyle.Render("Edit Mode: RGB"))
		lines = append(lines, "")
		opacity := int(math.Round(c.A * 100))
		fields := []struct {
			label string
			value string
			idx   int
		}{
			{"Red", fmt.Sprintf("%d", c.R), FieldRed},
			{"Green", fmt.Sprintf("%d", c.G), FieldGreen},
			{"Blue", fmt.Sprintf("%d", c.B), FieldBlue},
			{"Opacity", fmt.Sprintf("%d%%", opacity), FieldOpacity},
		}
		for _, f := range fields {
			lines = append(lines, renderField(f.label, f.value, f.idx == selectedField, width))
		}
		lines = append(lines, "")
		lines = append(lines, ModeHintStyle.Render("Tab: switch to HSV"))
	}

	return strings.Join(lines, "\n")
}

func renderField(label, value string, selected bool, width int) string {
	if selected {
		row := fmt.Sprintf("  %-12s %s", label, value)
		if len(row) < width {
			row += strings.Repeat(" ", width-len(row))
		}
		return SelectedFieldStyle.Render(row)
	}
	return FieldLabelStyle.Render(fmt.Sprintf("  %-12s", label)) + " " + FieldValueStyle.Render(value)
}
