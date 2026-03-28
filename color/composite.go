package color

import "math"

// CompositeOver returns the opaque color produced by drawing fg over bg.
func CompositeOver(fg, bg Color) Color {
	a := ClampFloat(fg.A, 0, 1)
	inverse := 1 - a

	r := int(math.Round(float64(fg.R)*a + float64(bg.R)*inverse))
	g := int(math.Round(float64(fg.G)*a + float64(bg.G)*inverse))
	b := int(math.Round(float64(fg.B)*a + float64(bg.B)*inverse))

	return Color{
		R: ClampUint8(r),
		G: ClampUint8(g),
		B: ClampUint8(b),
		A: 1,
	}
}
