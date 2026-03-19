package color

import (
	"fmt"
	"math"
)

// FormatRGB returns the normalized CSS rgb() string.
// Alpha is omitted when 100%, otherwise shown as an integer percentage.
func FormatRGB(c Color) string {
	if c.A >= 1.0 {
		return fmt.Sprintf("rgb(%d %d %d)", c.R, c.G, c.B)
	}
	pct := int(math.Round(c.A * 100))
	return fmt.Sprintf("rgb(%d %d %d / %d%%)", c.R, c.G, c.B, pct)
}

// FormatHEX returns the normalized lowercase hex string.
// Alpha byte is omitted when alpha is 100%.
func FormatHEX(c Color) string {
	if c.A >= 1.0 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	alphaByte := uint8(math.Round(c.A * 255))
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, alphaByte)
}

// FormatHSL returns the normalized CSS hsl() string.
// Alpha is omitted when 100%, otherwise shown as an integer percentage.
func FormatHSL(c Color) string {
	hsl := RGBToHSL(c)
	h := int(math.Round(hsl.H))
	s := int(math.Round(hsl.S * 100))
	l := int(math.Round(hsl.L * 100))

	if c.A >= 1.0 {
		return fmt.Sprintf("hsl(%d %d%% %d%%)", h, s, l)
	}
	pct := int(math.Round(c.A * 100))
	return fmt.Sprintf("hsl(%d %d%% %d%% / %d%%)", h, s, l, pct)
}
