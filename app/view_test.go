package app

import "testing"

func TestUsePortraitLayoutUsesVisualAspectRatio(t *testing.T) {
	tests := []struct {
		width  int
		height int
		want   bool
		name   string
	}{
		{name: "very wide stays landscape", width: 140, height: 40, want: false},
		{name: "visually narrow terminal becomes portrait", width: 90, height: 50, want: true},
		{name: "square-ish terminal becomes portrait", width: 80, height: 50, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := usePortraitLayout(tc.width, tc.height)
			if got != tc.want {
				t.Fatalf("usePortraitLayout(%d, %d) = %v, want %v", tc.width, tc.height, got, tc.want)
			}
		})
	}
}
