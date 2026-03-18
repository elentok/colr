package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/elentok/colr/color"
	"github.com/elentok/colr/ui"
)

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	// Width: outer border takes 1 col per side = 2 total.
	outerW := m.width
	innerW := outerW - 2

	// Height budget (all rows that aren't the split panels):
	//   outer border:    2  (top + bottom)
	//   header:          3
	//   outputs panel:   5  (3 content + 2 border)
	//   footer:          1
	//   ─────────────────────
	//   fixed total:    11
	//
	// The split panels (including their NormalBorder) fill the rest.
	const fixedH = 11
	splitPanelH := m.height - fixedH // total panel height including border
	if splitPanelH < 4 {
		splitPanelH = 4
	}
	splitContentH := splitPanelH - 2 // content height passed to Height()

	// Split widths: preview ~40%, editor ~60%.
	previewW := innerW * 40 / 100
	editorW := innerW - previewW

	// Render each region.
	header := ui.RenderHeader(
		m.originalClip,
		color.FormatRGB(m.currentColor),
		m.toastMessage,
		innerW,
	)
	preview := ui.RenderPreview(m.currentColor, previewW-2, splitContentH)
	editor := ui.RenderEditor(m.currentColor, m.editMode, m.selectedField, m.lastHue, editorW-2)
	// outputsPanel content width = innerW-2 (NormalBorder takes 1 per side).
	outputs := ui.RenderOutputs(m.currentColor, innerW-2)
	footer := ui.RenderFooter(innerW)

	// Assemble panels.
	previewPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(previewW - 2).
		Height(splitContentH).
		Render(preview)

	editorPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(editorW - 2).
		Height(splitContentH).
		Render(editor)

	splitRow := lipgloss.JoinHorizontal(lipgloss.Top, previewPanel, editorPanel)

	outputsPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(innerW - 2).
		Render(outputs)

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		splitRow,
		outputsPanel,
		footer,
	)

	// Outer frame: Height(m.height-2) so the border fills exactly m.height rows.
	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Width(innerW).
		Height(m.height - 2).
		Render(body)
	// Inject "colr" into the top border line after rendering.
	frame = injectBorderTitle(frame)

	if m.showHelp {
		return ui.RenderHelp(m.width, m.height)
	}

	return frame
}

// injectBorderTitle replaces "╭──────" with "╭ colr " in the rendered frame,
// embedding the app name in the top border line without changing visual width.
func injectBorderTitle(s string) string {
	return strings.Replace(s, "╭──────", "╭ colr ", 1)
}
