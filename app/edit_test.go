package app

import (
	"math"
	"testing"

	"github.com/elentok/colr/color"
	"github.com/elentok/colr/history"
	"github.com/elentok/colr/ui"
)

func newTestModel(c color.Color) Model {
	return NewModel("test", c, nil)
}

func newHistoryTestModel(c color.Color, entries []history.Entry) Model {
	return NewModel("test", c, entries)
}

func TestAdjustHSVHue(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1}) // red, H=0
	m.selectedField = ui.FieldHue

	m = applyAdjust(m, +1, false)
	hsv := color.RGBToHSV(m.currentColor)
	if math.Round(hsv.H) != 1 {
		t.Errorf("expected H=1 after +1, got H=%v", hsv.H)
	}

	m = applyAdjust(m, -1, false)
	hsv = color.RGBToHSV(m.currentColor)
	if math.Round(hsv.H) != 0 {
		t.Errorf("expected H=0 after -1, got H=%v", hsv.H)
	}
}

func TestAdjustHSVHueLargeStep(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.selectedField = ui.FieldHue

	m = applyAdjust(m, +1, true)
	hsv := color.RGBToHSV(m.currentColor)
	if math.Round(hsv.H) != 10 {
		t.Errorf("expected H=10 after large +1, got H=%v", hsv.H)
	}
}

func TestAdjustHSVHueWrapForward(t *testing.T) {
	// Set lastHue to 359 directly and increment — should wrap to 0.
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.selectedField = ui.FieldHue
	m.lastHue = 359
	// Make color achromatic so applyHSVAdjust uses lastHue.
	m.currentColor = color.Color{R: 128, G: 128, B: 128, A: 1}

	m = applyAdjust(m, +1, false)
	if m.lastHue != 0 {
		t.Errorf("expected lastHue=0 after wrap from 359, got %v", m.lastHue)
	}
}

func TestAdjustHSVHueWrapBackward(t *testing.T) {
	// Set lastHue to 0 and decrement — should wrap to 359.
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.selectedField = ui.FieldHue
	m.lastHue = 0
	m.currentColor = color.Color{R: 128, G: 128, B: 128, A: 1}

	m = applyAdjust(m, -1, false)
	if m.lastHue != 359 {
		t.Errorf("expected lastHue=359 after wrap from 0, got %v", m.lastHue)
	}
}

func TestAdjustHSVSaturationClamp(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1}) // S=1
	m.selectedField = ui.FieldSaturation

	// Clamp at max
	for range 5 {
		m = applyAdjust(m, +1, false)
	}
	hsv := color.RGBToHSV(m.currentColor)
	if hsv.S > 1.0 {
		t.Errorf("saturation exceeded 1.0: %v", hsv.S)
	}

	// Drive to 0
	m2 := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m2.selectedField = ui.FieldSaturation
	for range 110 {
		m2 = applyAdjust(m2, -1, false)
	}
	hsv2 := color.RGBToHSV(m2.currentColor)
	if hsv2.S < 0 {
		t.Errorf("saturation went below 0: %v", hsv2.S)
	}
}

func TestAdjustHSVValueClamp(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1}) // V=1
	m.selectedField = ui.FieldValue

	for range 10 {
		m = applyAdjust(m, +1, false)
	}
	hsv := color.RGBToHSV(m.currentColor)
	if hsv.V > 1.0 {
		t.Errorf("value exceeded 1.0: %v", hsv.V)
	}

	m2 := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m2.selectedField = ui.FieldValue
	for range 110 {
		m2 = applyAdjust(m2, -1, false)
	}
	hsv2 := color.RGBToHSV(m2.currentColor)
	if hsv2.V < 0 {
		t.Errorf("value went below 0: %v", hsv2.V)
	}
}

func TestAdjustRGBRedClamp(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.editMode = ModeRGB
	m.selectedField = ui.FieldRed

	m = applyAdjust(m, +1, false)
	if m.currentColor.R != 255 {
		t.Errorf("expected R=255 (clamped), got %d", m.currentColor.R)
	}

	m2 := newTestModel(color.Color{R: 0, G: 0, B: 0, A: 1})
	m2.editMode = ModeRGB
	m2.selectedField = ui.FieldRed
	m2 = applyAdjust(m2, -1, false)
	if m2.currentColor.R != 0 {
		t.Errorf("expected R=0 (clamped), got %d", m2.currentColor.R)
	}
}

func TestAdjustRGBLargeStep(t *testing.T) {
	m := newTestModel(color.Color{R: 100, G: 100, B: 100, A: 1})
	m.editMode = ModeRGB
	m.selectedField = ui.FieldGreen

	m = applyAdjust(m, +1, true)
	if m.currentColor.G != 110 {
		t.Errorf("expected G=110 after large +1, got %d", m.currentColor.G)
	}

	m = applyAdjust(m, -1, true)
	if m.currentColor.G != 100 {
		t.Errorf("expected G=100 after large -1, got %d", m.currentColor.G)
	}
}

func TestAdjustOpacityHSV(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1.0})
	m.selectedField = ui.FieldOpacity

	m = applyAdjust(m, +1, false) // should clamp at 1.0
	if m.currentColor.A > 1.0 {
		t.Errorf("opacity exceeded 1.0: %v", m.currentColor.A)
	}

	m2 := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 0.0})
	m2.selectedField = ui.FieldOpacity
	m2 = applyAdjust(m2, -1, false) // should clamp at 0.0
	if m2.currentColor.A < 0.0 {
		t.Errorf("opacity went below 0: %v", m2.currentColor.A)
	}
}

func TestAdjustOpacityLargeStep(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 0.5})
	m.selectedField = ui.FieldOpacity

	m = applyAdjust(m, +1, true) // +5%
	if math.Abs(m.currentColor.A-0.55) > 0.001 {
		t.Errorf("expected A=0.55, got %v", m.currentColor.A)
	}
}

func TestModeSwitchPreservesColor(t *testing.T) {
	original := color.Color{R: 123, G: 45, B: 67, A: 0.8}
	m := newTestModel(original)

	// Switch to RGB and back
	m = handleKey(m, "tab") // → RGB
	m = handleKey(m, "tab") // → HSV

	if m.currentColor != original {
		t.Errorf("mode switch changed color: got %+v, want %+v", m.currentColor, original)
	}
}

func TestModeSwitchKeys(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	m = handleKey(m, "2")
	if m.editMode != ModeRGB {
		t.Errorf("expected ModeRGB after '2', got %v", m.editMode)
	}

	m = handleKey(m, "1")
	if m.editMode != ModeHSV {
		t.Errorf("expected ModeHSV after '1', got %v", m.editMode)
	}

	m = handleKey(m, "tab")
	if m.editMode != ModeRGB {
		t.Errorf("expected ModeRGB after tab, got %v", m.editMode)
	}
}

func TestReset(t *testing.T) {
	original := color.Color{R: 255, G: 128, B: 0, A: 0.75}
	m := newTestModel(original)

	// Apply several edits
	m.editMode = ModeRGB
	m.selectedField = ui.FieldRed
	for range 20 {
		m = applyAdjust(m, +1, false)
	}
	m.selectedField = ui.FieldBlue
	for range 10 {
		m = applyAdjust(m, -1, false)
	}

	// Reset should restore exactly
	m = applyReset(m)
	if m.currentColor != original {
		t.Errorf("reset: got %+v, want %+v", m.currentColor, original)
	}
}

func TestFieldNavigation(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.selectedField = 0

	m = handleKey(m, "j")
	if m.selectedField != 1 {
		t.Errorf("j: expected field 1, got %d", m.selectedField)
	}

	m = handleKey(m, "k")
	if m.selectedField != 0 {
		t.Errorf("k: expected field 0, got %d", m.selectedField)
	}

	// Clamp at top
	m = handleKey(m, "k")
	if m.selectedField != 0 {
		t.Errorf("k at 0: expected field 0 (clamped), got %d", m.selectedField)
	}

	m = handleKey(m, "G")
	if m.selectedField != ui.FieldCount-1 {
		t.Errorf("G: expected last field %d, got %d", ui.FieldCount-1, m.selectedField)
	}

	// Clamp at bottom
	m = handleKey(m, "j")
	if m.selectedField != ui.FieldCount-1 {
		t.Errorf("j at last: expected %d (clamped), got %d", ui.FieldCount-1, m.selectedField)
	}

	m = handleKey(m, "g")
	if m.selectedField != 0 {
		t.Errorf("g: expected field 0, got %d", m.selectedField)
	}
}

func TestGrayscaleHueStability(t *testing.T) {
	// Start with red (H=0), drive saturation to 0, lastHue should be preserved.
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	initialHue := m.lastHue // should be ~0 for red

	// Drive saturation to 0
	m.selectedField = ui.FieldSaturation
	for range 110 {
		m = applyAdjust(m, -1, false)
	}
	hsv := color.RGBToHSV(m.currentColor)
	if hsv.S != 0 {
		t.Fatalf("expected S=0 after driving down, got %v", hsv.S)
	}

	// lastHue should still be the original hue, not snapped to 0.
	// (For red starting at H=0, lastHue starts at 0, so we test with a different color.)
	_ = initialHue

	// Use a non-zero starting hue (green, H=120).
	m2 := newTestModel(color.Color{R: 0, G: 255, B: 0, A: 1}) // H=120
	if math.Abs(m2.lastHue-120) > 1 {
		t.Fatalf("expected lastHue≈120 for green, got %v", m2.lastHue)
	}

	m2.selectedField = ui.FieldSaturation
	for range 110 {
		m2 = applyAdjust(m2, -1, false)
	}
	// lastHue should still be ~120 even though color is now gray.
	if math.Abs(m2.lastHue-120) > 1 {
		t.Errorf("lastHue changed to %v after driving S to 0 (expected ~120)", m2.lastHue)
	}
}
