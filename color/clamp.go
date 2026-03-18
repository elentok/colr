package color

import "math"

// ClampInt clamps v to [min, max].
func ClampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ClampFloat clamps v to [min, max].
func ClampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// WrapFloat wraps v to [min, max). Used for hue wraparound.
// e.g. WrapFloat(360, 0, 360) == 0, WrapFloat(-1, 0, 360) == 359.
func WrapFloat(v, min, max float64) float64 {
	r := max - min
	if r <= 0 {
		return min
	}
	v = math.Mod(v-min, r)
	if v < 0 {
		v += r
	}
	return v + min
}

// ClampUint8 clamps an int to [0, 255] and returns uint8.
func ClampUint8(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
