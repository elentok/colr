package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/elentok/colr/color"
)

type PreviewLayout int

const (
	PreviewStacked PreviewLayout = iota
	PreviewSideBySide
)

// RenderPreview renders the color preview block.
func RenderPreview(original, current, previewBG color.Color, width, height int, layout PreviewLayout) string {
	if layout == PreviewSideBySide {
		return renderPreviewSideBySide(original, current, previewBG, width, height)
	}

	return renderPreviewStacked(original, current, previewBG, width, height)
}

func renderPreviewStacked(original, current, previewBG color.Color, width, height int) string {
	if height < 6 {
		height = 6
	}

	topH := height / 2
	bottomH := height - topH
	if topH < 3 {
		topH = 3
		bottomH = height - topH
	}
	if bottomH < 3 {
		bottomH = 3
		topH = height - bottomH
	}

	top := renderPreviewSection("Original", original, previewBG, width, topH)
	bottom := renderPreviewSection("Edited", current, previewBG, width, bottomH)
	return strings.Join([]string{top, bottom}, "\n")
}

func renderPreviewSideBySide(original, current, previewBG color.Color, width, height int) string {
	if width < 20 {
		width = 20
	}
	if height < 4 {
		height = 4
	}

	leftW := width / 2
	rightW := width - leftW
	if leftW < 10 {
		leftW = 10
		rightW = width - leftW
	}
	if rightW < 10 {
		rightW = 10
		leftW = width - rightW
	}

	left := renderPreviewSection("Original", original, previewBG, leftW, height)
	right := renderPreviewSection("Edited", current, previewBG, rightW, height)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderPreviewSection(title string, c, previewBG color.Color, width, height int) string {
	rawHex := color.FormatHEX(c)
	overBG := color.CompositeOver(c, previewBG)
	overHex := color.FormatHEX(overBG)
	overBGHex := fmt.Sprintf("#%02x%02x%02x", overBG.R, overBG.G, overBG.B)
	previewBGHex := fmt.Sprintf("#%02x%02x%02x", previewBG.R, previewBG.G, previewBG.B)

	// Compute foreground suggestion via relative luminance for composited color.
	r := float64(overBG.R) / 255.0
	g := float64(overBG.G) / 255.0
	b := float64(overBG.B) / 255.0
	luminance := 0.2126*r + 0.7152*g + 0.0722*b
	fgANSI := lipgloss.Color("255")
	if luminance > 0.5 {
		fgANSI = lipgloss.Color("0")
	}

	bgName := "white"
	if previewBG.R == 0 && previewBG.G == 0 && previewBG.B == 0 {
		bgName = "black"
	}

	var rows []string
	previewRows := height - 2 // leave room for title + fg hint line
	if previewRows < 1 {
		previewRows = 1
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("238")).
		Bold(true).
		Width(width)
	rows = append(rows, titleStyle.Render(title))

	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(previewBGHex))
	overStyle := lipgloss.NewStyle().Background(lipgloss.Color(overBGHex))
	innerW := width - 2
	if innerW < 0 {
		innerW = 0
	}

	for range previewRows {
		switch {
		case width <= 0:
			rows = append(rows, "")
		case width == 1:
			rows = append(rows, bgStyle.Render(" "))
		default:
			row := bgStyle.Render(" ") + overStyle.Render(strings.Repeat(" ", innerW)) + bgStyle.Render(" ")
			rows = append(rows, row)
		}
	}

	// Bottom line shows raw and composited values.
	fgHint := previewInfoLine(width, rawHex, overHex, bgName)
	hintStyle := lipgloss.NewStyle().
		Foreground(fgANSI).
		Background(lipgloss.Color(overBGHex)).
		Width(width)
	rows = append(rows, hintStyle.Render(fgHint))

	return strings.Join(rows, "\n")
}

func previewInfoLine(width int, rawHex, overHex, bgName string) string {
	switch {
	case width >= 44:
		return fmt.Sprintf("raw %s  over-bg %s  bg %s", rawHex, overHex, bgName)
	case width >= 20:
		return fmt.Sprintf("over-bg %s  bg %s", overHex, bgName)
	case width >= 16:
		return fmt.Sprintf("over %s", overHex)
	default:
		return overHex
	}
}
