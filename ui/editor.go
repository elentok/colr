package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/elentok/colr/color"
)

// RenderEditor renders the editor pane showing adjustable fields.
// lastHue is the preserved hue used when saturation is 0 (grayscale stability).
func RenderEditor(c color.Color, mode EditMode, selectedField int, lastHue float64, width int) string {
	var lines []string

	if mode == ModeHSV {
		lines = append(lines, ModeTitleStyle.Render("Edit Mode: HSV"))
		lines = append(lines, "")
		hsv := color.RGBToHSV(c)
		displayHue := hsv.H
		if hsv.S == 0 {
			displayHue = lastHue
		}
		opacity := math.Round(c.A * 100)
		fields := []struct {
			label   string
			value   string
			current float64
			min     float64
			max     float64
			idx     int
		}{
			{"Hue", fmt.Sprintf("%d°", int(math.Round(displayHue))), displayHue, 0, 360, FieldHue},
			{"Saturation", fmt.Sprintf("%d%%", int(math.Round(hsv.S*100))), hsv.S * 100, 0, 100, FieldSaturation},
			{"Value", fmt.Sprintf("%d%%", int(math.Round(hsv.V*100))), hsv.V * 100, 0, 100, FieldValue},
			{"Opacity", fmt.Sprintf("%d%%", int(opacity)), opacity, 0, 100, FieldOpacity},
		}
		for _, f := range fields {
			lines = append(lines, renderField(f.label, f.value, f.current, f.min, f.max, f.idx == selectedField, width))
		}
		lines = append(lines, "")
		lines = append(lines, ModeHintStyle.Render("Tab: switch to RGB"))
	} else {
		lines = append(lines, ModeTitleStyle.Render("Edit Mode: RGB"))
		lines = append(lines, "")
		opacity := math.Round(c.A * 100)
		fields := []struct {
			label   string
			value   string
			current float64
			min     float64
			max     float64
			idx     int
		}{
			{"Red", fmt.Sprintf("%d", c.R), float64(c.R), 0, 255, FieldRed},
			{"Green", fmt.Sprintf("%d", c.G), float64(c.G), 0, 255, FieldGreen},
			{"Blue", fmt.Sprintf("%d", c.B), float64(c.B), 0, 255, FieldBlue},
			{"Opacity", fmt.Sprintf("%d%%", int(opacity)), opacity, 0, 100, FieldOpacity},
		}
		for _, f := range fields {
			lines = append(lines, renderField(f.label, f.value, f.current, f.min, f.max, f.idx == selectedField, width))
		}
		lines = append(lines, "")
		lines = append(lines, ModeHintStyle.Render("Tab: switch to HSV"))
	}

	return strings.Join(lines, "\n")
}

// fieldPrefixWidth is the number of visual chars before the slider:
// 2 lead + 12 label + 1 sp + 4 value + 1 sp = 20
const fieldPrefixWidth = 20

func renderField(label, value string, current, min, max float64, selected bool, width int) string {
	sliderW := width - fieldPrefixWidth
	if sliderW < 3 {
		sliderW = 3
	}

	// Pad value to 4 chars for consistent alignment.
	paddedValue := fmt.Sprintf("%-4s", value)

	if selected {
		slider := renderSlider(current, min, max, sliderW, true)
		row := fmt.Sprintf("  %-12s %s %s", label, paddedValue, slider)
		// Pad to full width so the highlight background covers the whole line.
		rowRunes := []rune(row)
		if len(rowRunes) < width {
			row += strings.Repeat(" ", width-len(rowRunes))
		}
		return SelectedFieldStyle.Render(row)
	}

	slider := renderSlider(current, min, max, sliderW, false)
	labelStr := FieldLabelStyle.Render(fmt.Sprintf("  %-12s", label))
	valueStr := FieldValueStyle.Render(paddedValue)
	return labelStr + " " + valueStr + " " + slider
}

// renderSlider renders a text-based slider: [────●────]
func renderSlider(value, min, max float64, width int, selected bool) string {
	if width < 3 {
		return ""
	}
	innerW := width - 2 // subtract [ and ]
	if innerW < 1 {
		return "[]"
	}

	ratio := 0.0
	if max > min {
		ratio = (value - min) / (max - min)
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	knobPos := int(math.Round(ratio * float64(innerW-1)))
	track := strings.Repeat("─", knobPos) + "●" + strings.Repeat("─", innerW-1-knobPos)

	slider := "[" + track + "]"
	if !selected {
		return SliderStyle.Render(slider)
	}
	// For selected rows, return unstyled — SelectedFieldStyle colors the whole row.
	return slider
}
