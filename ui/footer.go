package ui

// RenderFooter renders the key hint bar.
func RenderFooter(width int) string {
	sep := FooterSepStyle.Render(" • ")
	hints := []string{
		"hjkl move",
		"h/l adjust",
		"H/L larger step",
		"tab mode",
		"yr/yx/yh copy",
		"R reset",
		"? help",
		"q quit",
	}
	line := ""
	for i, h := range hints {
		if i > 0 {
			line += sep
		}
		line += FooterStyle.Render(h)
	}
	return line
}
