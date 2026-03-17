package color

import "math"

// RGBToHSV converts a Color to HSV. For achromatic colors (S=0) H is 0.
func RGBToHSV(c Color) HSV {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	v := max
	s := 0.0
	if max != 0 {
		s = delta / max
	}

	h := 0.0
	if delta != 0 {
		switch max {
		case r:
			h = (g - b) / delta
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/delta + 2
		case b:
			h = (r-g)/delta + 4
		}
		h *= 60
	}

	return HSV{H: h, S: s, V: v}
}

// HSVToRGB converts HSV + alpha to Color.
func HSVToRGB(hsv HSV, alpha float64) Color {
	h, s, v := hsv.H, hsv.S, hsv.V

	if s == 0 {
		val := uint8(math.Round(v * 255))
		return Color{R: val, G: val, B: val, A: alpha}
	}

	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	sector := h / 60
	i := math.Floor(sector)
	f := sector - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var r, g, b float64
	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return Color{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: alpha,
	}
}

// RGBToHSL converts a Color to HSL. For achromatic colors (S=0) H is 0.
func RGBToHSL(c Color) HSL {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	l := (max + min) / 2

	s := 0.0
	if delta != 0 {
		s = delta / (1 - math.Abs(2*l-1))
	}

	h := 0.0
	if delta != 0 {
		switch max {
		case r:
			h = (g - b) / delta
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/delta + 2
		case b:
			h = (r-g)/delta + 4
		}
		h *= 60
	}

	return HSL{H: h, S: s, L: l}
}

// HSLToRGB converts HSL + alpha to Color.
func HSLToRGB(hsl HSL, alpha float64) Color {
	h, s, l := hsl.H, hsl.S, hsl.L

	if s == 0 {
		val := uint8(math.Round(l * 255))
		return Color{R: val, G: val, B: val, A: alpha}
	}

	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return Color{
		R: uint8(math.Round((r + m) * 255)),
		G: uint8(math.Round((g + m) * 255)),
		B: uint8(math.Round((b + m) * 255)),
		A: alpha,
	}
}
