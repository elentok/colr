package color

import (
	"math"
	"testing"
)

func alphaEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func TestParseHEX(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"#FF0000", Color{255, 0, 0, 1.0}},
		{"FF0000", Color{255, 0, 0, 1.0}},
		{"#ff0000", Color{255, 0, 0, 1.0}},
		{"ff0000", Color{255, 0, 0, 1.0}},
		{"#FF000080", Color{255, 0, 0, 0.502}},
		{"FF000080", Color{255, 0, 0, 0.502}},
		{"#a3a3a3", Color{163, 163, 163, 1.0}},
		{"a3a3a3", Color{163, 163, 163, 1.0}},
		{"a3a3a3aa", Color{163, 163, 163, 0.667}},
		{"A3A3A3", Color{163, 163, 163, 1.0}},
		{"A3A3A3AA", Color{163, 163, 163, 0.667}},
		{"#000000", Color{0, 0, 0, 1.0}},
		{"#FFFFFF", Color{255, 255, 255, 1.0}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if got.R != tc.want.R || got.G != tc.want.G || got.B != tc.want.B || !alphaEqual(got.A, tc.want.A) {
				t.Errorf("Parse(%q) = %+v, want %+v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseRGB(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"rgb(255 0 0)", Color{255, 0, 0, 1.0}},
		{"rgb(255 0 0 / 50%)", Color{255, 0, 0, 0.5}},
		{"rgb(255 0 0 / 0.5)", Color{255, 0, 0, 0.5}},
		{"rgb(255, 0, 0)", Color{255, 0, 0, 1.0}},
		{"rgb(100% 0% 0%)", Color{255, 0, 0, 1.0}},
		{"rgb(0 0 0)", Color{0, 0, 0, 1.0}},
		{"rgb(255 255 255)", Color{255, 255, 255, 1.0}},
		// Out-of-range clamped
		{"rgb(100 200 300)", Color{100, 200, 255, 1.0}},
		// Case insensitive
		{"RGB(255 0 0)", Color{255, 0, 0, 1.0}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if got.R != tc.want.R || got.G != tc.want.G || got.B != tc.want.B || !alphaEqual(got.A, tc.want.A) {
				t.Errorf("Parse(%q) = %+v, want %+v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseRGBA(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"rgba(255, 0, 0, 50%)", Color{255, 0, 0, 0.5}},
		{"rgba(255, 0, 0, 0.5)", Color{255, 0, 0, 0.5}},
		{"rgba(255, 0, 0, .5)", Color{255, 0, 0, 0.5}},
		{"rgba(255, 0, 0, 1)", Color{255, 0, 0, 1.0}},
		{"rgba(255, 0, 0, 0)", Color{255, 0, 0, 0.0}},
		{"rgba(0, 0, 0, 0)", Color{0, 0, 0, 0.0}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if got.R != tc.want.R || got.G != tc.want.G || got.B != tc.want.B || !alphaEqual(got.A, tc.want.A) {
				t.Errorf("Parse(%q) = %+v, want %+v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseHSL(t *testing.T) {
	tests := []struct {
		input string
		wantR uint8
		wantG uint8
		wantB uint8
		wantA float64
	}{
		{"hsl(0 100% 50%)", 255, 0, 0, 1.0},
		{"hsl(0 100% 50% / 50%)", 255, 0, 0, 0.5},
		{"hsl(120 100% 50%)", 0, 255, 0, 1.0},
		{"hsl(240 100% 50%)", 0, 0, 255, 1.0},
		{"hsl(0 0% 0%)", 0, 0, 0, 1.0},
		{"hsl(0 0% 100%)", 255, 255, 255, 1.0},
		{"hsl(32 100% 50%)", 255, 136, 0, 1.0},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			// Allow ±1 for rounding
			diff := func(a, b uint8) int {
				if a > b {
					return int(a - b)
				}
				return int(b - a)
			}
			if diff(got.R, tc.wantR) > 1 || diff(got.G, tc.wantG) > 1 || diff(got.B, tc.wantB) > 1 || !alphaEqual(got.A, tc.wantA) {
				t.Errorf("Parse(%q) = %+v, want R=%d G=%d B=%d A=%v", tc.input, got, tc.wantR, tc.wantG, tc.wantB, tc.wantA)
			}
		})
	}
}

func TestParseBare(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"255, 0, 0", Color{255, 0, 0, 1.0}},
		{"255 0 0", Color{255, 0, 0, 1.0}},
		{"100% 0% 0%", Color{255, 0, 0, 1.0}},
		{"0 0 0", Color{0, 0, 0, 1.0}},
		{"255 255 255", Color{255, 255, 255, 1.0}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if got.R != tc.want.R || got.G != tc.want.G || got.B != tc.want.B || !alphaEqual(got.A, tc.want.A) {
				t.Errorf("Parse(%q) = %+v, want %+v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseWhitespace(t *testing.T) {
	tests := []string{
		"  #FF0000  ",
		"  rgb(255 0 0)  ",
		"  rgba(255, 0, 0, 1.0)  ",
		"\t#FF0000\t",
		"#FF0000;",
		"rgb(255 0 0);",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", input, err)
			}
			if got.R != 255 || got.G != 0 || got.B != 0 {
				t.Errorf("Parse(%q) = %+v, want R=255 G=0 B=0", input, got)
			}
		})
	}
}

func TestParseAlphaEdgeCases(t *testing.T) {
	tests := []struct {
		input string
		wantA float64
	}{
		{"rgba(255, 0, 0, 0)", 0.0},
		{"rgba(255, 0, 0, 0%)", 0.0},
		{"rgba(255, 0, 0, 100%)", 1.0},
		{"rgba(255, 0, 0, 1.0)", 1.0},
		{"rgba(255, 0, 0, 0.0)", 0.0},
		{"rgba(255, 0, 0, .5)", 0.5},
		{"rgb(255 0 0 / 0%)", 0.0},
		{"rgb(255 0 0 / 100%)", 1.0},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if !alphaEqual(got.A, tc.wantA) {
				t.Errorf("Parse(%q).A = %v, want %v", tc.input, got.A, tc.wantA)
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	invalid := []string{
		"hello",
		"rgb(foo bar baz)",
		"rgb(10 20)",
		"#XYZ123",
		"hsl(20 30)",
		"",
		"not a color",
		"#12345",  // 5 hex digits
		"#1234567", // 7 hex digits
		"rgba(255, 0, 0)", // rgba missing alpha
	}
	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got none", input)
			}
		})
	}
}
