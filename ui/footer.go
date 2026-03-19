package ui

// RenderFooter renders the key hint bar with keys in accent color.
func RenderFooter(width int) string {
	sep := FooterSepStyle.Render(" • ")
	// Each hint: key styled in accent color, description in dim.
	hints := []struct{ key, desc string }{
		{"hjkl", "move"},
		{"h/l", "adjust"},
		{"H/L", "larger step"},
		{"tab", "mode"},
		{"p", "history"},
		{"yr/yx/yh/yn/yy", "copy"},
		{"R", "reset"},
		{"?", "help"},
		{"q", "quit"},
	}
	line := ""
	for i, h := range hints {
		if i > 0 {
			line += sep
		}
		line += FooterKeyStyle.Render(h.key) + FooterStyle.Render(" "+h.desc)
	}
	return line
}
