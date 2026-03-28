package color

import "testing"

func TestCompositeOver(t *testing.T) {
	tests := []struct {
		name string
		fg   Color
		bg   Color
		want Color
	}{
		{
			name: "opaque foreground replaces background",
			fg:   Color{R: 255, G: 0, B: 0, A: 1},
			bg:   Color{R: 0, G: 0, B: 255, A: 1},
			want: Color{R: 255, G: 0, B: 0, A: 1},
		},
		{
			name: "transparent foreground keeps background",
			fg:   Color{R: 255, G: 0, B: 0, A: 0},
			bg:   Color{R: 12, G: 34, B: 56, A: 1},
			want: Color{R: 12, G: 34, B: 56, A: 1},
		},
		{
			name: "half alpha blends channels",
			fg:   Color{R: 255, G: 0, B: 0, A: 0.5},
			bg:   Color{R: 0, G: 0, B: 0, A: 1},
			want: Color{R: 128, G: 0, B: 0, A: 1},
		},
		{
			name: "alpha is clamped",
			fg:   Color{R: 255, G: 255, B: 255, A: 2},
			bg:   Color{R: 10, G: 20, B: 30, A: 1},
			want: Color{R: 255, G: 255, B: 255, A: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CompositeOver(tc.fg, tc.bg)
			if got != tc.want {
				t.Fatalf("CompositeOver(%+v, %+v) = %+v, want %+v", tc.fg, tc.bg, got, tc.want)
			}
		})
	}
}
