package color

import (
	"strings"
	"testing"
)

func TestNearestNamedColor_exact(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{Color{R: 255, G: 0, B: 0}, "red"},
		{Color{R: 0, G: 0, B: 255}, "blue"},
		{Color{R: 0, G: 128, B: 0}, "green"},
		{Color{R: 255, G: 255, B: 255}, "white"},
		{Color{R: 0, G: 0, B: 0}, "black"},
		{Color{R: 255, G: 99, B: 71}, "tomato"},
	}
	for _, tt := range tests {
		got := NearestNamedColor(tt.color)
		if got != tt.want {
			t.Errorf("NearestNamedColor(%v) = %q, want %q", tt.color, got, tt.want)
		}
	}
}

func TestNearestNamedColor_approximate(t *testing.T) {
	// A color that is not an exact CSS named color should get a ~ prefix.
	c := Color{R: 255, G: 96, B: 0} // close to tomato (255,99,71) or darkorange (255,140,0)
	got := NearestNamedColor(c)
	if !strings.HasPrefix(got, "~") {
		t.Errorf("NearestNamedColor(%v) = %q, want ~ prefix for approximate match", c, got)
	}
}

func TestNearestNamedColor_alphaIgnored(t *testing.T) {
	// Same RGB with different alpha should still match exactly.
	c := Color{R: 255, G: 0, B: 0, A: 0.5}
	got := NearestNamedColor(c)
	if got != "red" {
		t.Errorf("NearestNamedColor(%v) = %q, want \"red\"", c, got)
	}
}
