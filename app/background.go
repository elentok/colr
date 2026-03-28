package app

import "github.com/elentok/colr/color"

func previewBackgroundColor(useDark bool) color.Color {
	if useDark {
		return color.Color{R: 0, G: 0, B: 0, A: 1}
	}

	return color.Color{R: 255, G: 255, B: 255, A: 1}
}

func previewBackgroundName(useDark bool) string {
	if useDark {
		return "black"
	}

	return "white"
}
