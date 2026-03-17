package color

// Color is the internal RGBA representation. Alpha is in [0.0, 1.0].
type Color struct {
	R uint8
	G uint8
	B uint8
	A float64 // 0.0 - 1.0
}

// HSV holds hue-saturation-value components.
// H: 0–360, S: 0–1, V: 0–1
type HSV struct {
	H float64
	S float64
	V float64
}

// HSL holds hue-saturation-lightness components.
// H: 0–360, S: 0–1, L: 0–1
type HSL struct {
	H float64
	S float64
	L float64
}
