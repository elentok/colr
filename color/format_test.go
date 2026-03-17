package color

import (
	"testing"
)

func TestFormatRGB(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{Color{255, 0, 0, 1.0}, "rgb(255 0 0)"},
		{Color{255, 0, 0, 0.5}, "rgb(255 0 0 / 50%)"},
		{Color{0, 0, 0, 1.0}, "rgb(0 0 0)"},
		{Color{255, 255, 255, 1.0}, "rgb(255 255 255)"},
		{Color{0, 0, 0, 0.0}, "rgb(0 0 0 / 0%)"},
		{Color{163, 163, 163, 1.0}, "rgb(163 163 163)"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			got := FormatRGB(tc.color)
			if got != tc.want {
				t.Errorf("FormatRGB(%+v) = %q, want %q", tc.color, got, tc.want)
			}
		})
	}
}

func TestFormatHEX(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{Color{255, 0, 0, 1.0}, "#FF0000"},
		{Color{255, 0, 0, 0.5}, "#FF000080"},
		{Color{0, 0, 0, 1.0}, "#000000"},
		{Color{255, 255, 255, 1.0}, "#FFFFFF"},
		{Color{0, 0, 0, 0.0}, "#00000000"},
		{Color{163, 163, 163, 1.0}, "#A3A3A3"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			got := FormatHEX(tc.color)
			if got != tc.want {
				t.Errorf("FormatHEX(%+v) = %q, want %q", tc.color, got, tc.want)
			}
		})
	}
}

func TestFormatHSL(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{Color{255, 0, 0, 1.0}, "hsl(0 100% 50%)"},
		{Color{255, 0, 0, 0.5}, "hsl(0 100% 50% / 50%)"},
		{Color{0, 0, 0, 1.0}, "hsl(0 0% 0%)"},
		{Color{255, 255, 255, 1.0}, "hsl(0 0% 100%)"},
		{Color{0, 0, 0, 0.0}, "hsl(0 0% 0% / 0%)"},
		{Color{0, 255, 0, 1.0}, "hsl(120 100% 50%)"},
		{Color{0, 0, 255, 1.0}, "hsl(240 100% 50%)"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			got := FormatHSL(tc.color)
			if got != tc.want {
				t.Errorf("FormatHSL(%+v) = %q, want %q", tc.color, got, tc.want)
			}
		})
	}
}

func TestParseRoundTrip(t *testing.T) {
	inputs := []string{
		"#FF0000",
		"#000000",
		"#FFFFFF",
		"#00FF00",
		"#0000FF",
		"#808080",
		"#FF000080",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			c, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", input, err)
			}
			// Re-format to HEX and parse again — should be stable
			hex := FormatHEX(c)
			c2, err := Parse(hex)
			if err != nil {
				t.Fatalf("Parse(FormatHEX(%q) = %q): %v", input, hex, err)
			}
			if c.R != c2.R || c.G != c2.G || c.B != c2.B || !alphaEqual(c.A, c2.A) {
				t.Errorf("round-trip %q → %q not stable: %+v vs %+v", input, hex, c, c2)
			}
		})
	}
}
