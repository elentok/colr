package color

import (
	"math"
	"testing"
)

func colorsEqual(a, b Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && math.Abs(a.A-b.A) < 0.01
}

var representativeColors = []struct {
	name  string
	color Color
}{
	{"black", Color{0, 0, 0, 1.0}},
	{"white", Color{255, 255, 255, 1.0}},
	{"red", Color{255, 0, 0, 1.0}},
	{"green", Color{0, 255, 0, 1.0}},
	{"blue", Color{0, 0, 255, 1.0}},
	{"gray", Color{128, 128, 128, 1.0}},
	{"semi-transparent red", Color{255, 0, 0, 0.5}},
}

func TestRGBToHSVToRGB(t *testing.T) {
	for _, tc := range representativeColors {
		t.Run(tc.name, func(t *testing.T) {
			hsv := RGBToHSV(tc.color)
			got := HSVToRGB(hsv, tc.color.A)
			if !colorsEqual(got, tc.color) {
				t.Errorf("round-trip RGB→HSV→RGB: got %+v, want %+v (via HSV %+v)", got, tc.color, hsv)
			}
		})
	}
}

func TestRGBToHSLToRGB(t *testing.T) {
	for _, tc := range representativeColors {
		t.Run(tc.name, func(t *testing.T) {
			hsl := RGBToHSL(tc.color)
			got := HSLToRGB(hsl, tc.color.A)
			if !colorsEqual(got, tc.color) {
				t.Errorf("round-trip RGB→HSL→RGB: got %+v, want %+v (via HSL %+v)", got, tc.color, hsl)
			}
		})
	}
}

func TestGrayscaleHueIsZero(t *testing.T) {
	grays := []Color{
		{0, 0, 0, 1.0},
		{128, 128, 128, 1.0},
		{255, 255, 255, 1.0},
	}
	for _, c := range grays {
		hsv := RGBToHSV(c)
		if hsv.H != 0 {
			t.Errorf("RGBToHSV(%+v): expected H=0 for gray, got H=%v", c, hsv.H)
		}
		hsl := RGBToHSL(c)
		if hsl.H != 0 {
			t.Errorf("RGBToHSL(%+v): expected H=0 for gray, got H=%v", c, hsl.H)
		}
	}
}

func TestKnownHSVValues(t *testing.T) {
	tests := []struct {
		color Color
		want  HSV
	}{
		{Color{255, 0, 0, 1.0}, HSV{H: 0, S: 1, V: 1}},
		{Color{0, 255, 0, 1.0}, HSV{H: 120, S: 1, V: 1}},
		{Color{0, 0, 255, 1.0}, HSV{H: 240, S: 1, V: 1}},
		{Color{0, 0, 0, 1.0}, HSV{H: 0, S: 0, V: 0}},
		{Color{255, 255, 255, 1.0}, HSV{H: 0, S: 0, V: 1}},
	}
	for _, tc := range tests {
		got := RGBToHSV(tc.color)
		if math.Abs(got.H-tc.want.H) > 0.5 || math.Abs(got.S-tc.want.S) > 0.01 || math.Abs(got.V-tc.want.V) > 0.01 {
			t.Errorf("RGBToHSV(%+v) = %+v, want %+v", tc.color, got, tc.want)
		}
	}
}

func TestKnownHSLValues(t *testing.T) {
	tests := []struct {
		color Color
		want  HSL
	}{
		{Color{255, 0, 0, 1.0}, HSL{H: 0, S: 1, L: 0.5}},
		{Color{0, 255, 0, 1.0}, HSL{H: 120, S: 1, L: 0.5}},
		{Color{0, 0, 255, 1.0}, HSL{H: 240, S: 1, L: 0.5}},
		{Color{0, 0, 0, 1.0}, HSL{H: 0, S: 0, L: 0}},
		{Color{255, 255, 255, 1.0}, HSL{H: 0, S: 0, L: 1}},
	}
	for _, tc := range tests {
		got := RGBToHSL(tc.color)
		if math.Abs(got.H-tc.want.H) > 0.5 || math.Abs(got.S-tc.want.S) > 0.01 || math.Abs(got.L-tc.want.L) > 0.01 {
			t.Errorf("RGBToHSL(%+v) = %+v, want %+v", tc.color, got, tc.want)
		}
	}
}

func TestHSVToRGBWrap(t *testing.T) {
	// H=360 should behave same as H=0 (red)
	got360 := HSVToRGB(HSV{H: 360, S: 1, V: 1}, 1.0)
	got0 := HSVToRGB(HSV{H: 0, S: 1, V: 1}, 1.0)
	if !colorsEqual(got360, got0) {
		t.Errorf("H=360 should equal H=0: got %+v vs %+v", got360, got0)
	}
}
