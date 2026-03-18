package app

import (
	"github.com/elentok/colr/color"
	"github.com/elentok/colr/ui"
)

// applyAdjust adjusts the currently selected field by delta steps.
// delta should be +1 (increase) or -1 (decrease).
// large selects the larger step size defined in the spec.
func applyAdjust(m Model, delta int, large bool) Model {
	if m.editMode == ui.ModeHSV {
		return applyHSVAdjust(m, delta, large)
	}
	return applyRGBAdjust(m, delta, large)
}

func applyHSVAdjust(m Model, delta int, large bool) Model {
	if m.selectedField == ui.FieldOpacity {
		step := 1
		if large {
			step = 5
		}
		newA := m.currentColor.A + float64(delta*step)/100.0
		m.currentColor.A = color.ClampFloat(newA, 0, 1)
		return m
	}

	hsv := color.RGBToHSV(m.currentColor)
	// Restore remembered hue when achromatic to prevent snap-to-zero.
	if hsv.S == 0 {
		hsv.H = m.lastHue
	}

	switch m.selectedField {
	case ui.FieldHue:
		step := 1.0
		if large {
			step = 10
		}
		hsv.H = color.WrapFloat(hsv.H+float64(delta)*step, 0, 360)
		m.lastHue = hsv.H

	case ui.FieldSaturation:
		step := 1.0
		if large {
			step = 5
		}
		hsv.S = color.ClampFloat(hsv.S+float64(delta)*step/100.0, 0, 1)
		if hsv.S > 0 {
			m.lastHue = hsv.H
		}

	case ui.FieldValue:
		step := 1.0
		if large {
			step = 5
		}
		hsv.V = color.ClampFloat(hsv.V+float64(delta)*step/100.0, 0, 1)
	}

	m.currentColor = color.HSVToRGB(hsv, m.currentColor.A)
	return m
}

func applyRGBAdjust(m Model, delta int, large bool) Model {
	switch m.selectedField {
	case ui.FieldOpacity:
		step := 1
		if large {
			step = 5
		}
		newA := m.currentColor.A + float64(delta*step)/100.0
		m.currentColor.A = color.ClampFloat(newA, 0, 1)

	case ui.FieldRed:
		step := 1
		if large {
			step = 10
		}
		m.currentColor.R = color.ClampUint8(int(m.currentColor.R) + delta*step)

	case ui.FieldGreen:
		step := 1
		if large {
			step = 10
		}
		m.currentColor.G = color.ClampUint8(int(m.currentColor.G) + delta*step)

	case ui.FieldBlue:
		step := 1
		if large {
			step = 10
		}
		m.currentColor.B = color.ClampUint8(int(m.currentColor.B) + delta*step)
	}
	return m
}

// applyReset restores the current color to the original clipboard color.
func applyReset(m Model) Model {
	m.currentColor = m.originalColor
	hsv := color.RGBToHSV(m.originalColor)
	m.lastHue = hsv.H
	return m
}
