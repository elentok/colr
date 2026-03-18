package color

import (
	"testing"
)

func TestClampInt(t *testing.T) {
	tests := []struct{ v, min, max, want int }{
		{5, 0, 10, 5},
		{-1, 0, 10, 0},
		{11, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
	}
	for _, tc := range tests {
		got := ClampInt(tc.v, tc.min, tc.max)
		if got != tc.want {
			t.Errorf("ClampInt(%d, %d, %d) = %d, want %d", tc.v, tc.min, tc.max, got, tc.want)
		}
	}
}

func TestClampFloat(t *testing.T) {
	tests := []struct{ v, min, max, want float64 }{
		{0.5, 0, 1, 0.5},
		{-0.1, 0, 1, 0},
		{1.1, 0, 1, 1},
		{0, 0, 1, 0},
		{1, 0, 1, 1},
	}
	for _, tc := range tests {
		got := ClampFloat(tc.v, tc.min, tc.max)
		if got != tc.want {
			t.Errorf("ClampFloat(%v, %v, %v) = %v, want %v", tc.v, tc.min, tc.max, got, tc.want)
		}
	}
}

func TestWrapFloat(t *testing.T) {
	tests := []struct {
		v, min, max, want float64
	}{
		{180, 0, 360, 180},  // middle stays
		{0, 0, 360, 0},      // at min stays
		{360, 0, 360, 0},    // at max wraps to min
		{361, 0, 360, 1},    // just over max wraps
		{-1, 0, 360, 359},   // just under min wraps
		{-10, 0, 360, 350},  // further under min wraps
		{720, 0, 360, 0},    // double wrap
		{359, 0, 360, 359},  // just under max stays
	}
	for _, tc := range tests {
		got := WrapFloat(tc.v, tc.min, tc.max)
		if got != tc.want {
			t.Errorf("WrapFloat(%v, %v, %v) = %v, want %v", tc.v, tc.min, tc.max, got, tc.want)
		}
	}
}

func TestClampUint8(t *testing.T) {
	tests := []struct {
		v    int
		want uint8
	}{
		{128, 128},
		{0, 0},
		{255, 255},
		{-1, 0},
		{256, 255},
		{1000, 255},
	}
	for _, tc := range tests {
		got := ClampUint8(tc.v)
		if got != tc.want {
			t.Errorf("ClampUint8(%d) = %d, want %d", tc.v, got, tc.want)
		}
	}
}
