package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/elentok/colr/color"
	"github.com/elentok/colr/ui"
)

func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m Model) render() string {
	if m.width == 0 {
		return ""
	}

	if m.showHelp {
		return ui.RenderHelp(m.width, m.height)
	}

	// In lipgloss v2, Width() is the total outer width (border included).
	// outerW = full terminal width; innerW = content width inside outer border.
	outerW := m.width
	innerW := outerW - 2 // content width inside the outer rounded border

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

	// Split widths (total including each panel's NormalBorder):
	// previewW + editorW = innerW; each panel's Width() is its total outer width.
	previewW := innerW * 40 / 100
	editorW := innerW - previewW

	// Render each region (pass content widths = total - 2 for border).
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
	// Width() = total outer width (lipgloss v2 includes borders in the count).
	previewPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(previewW).
		Height(splitContentH + 2).
		Render(preview)

	editorPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(editorW).
		Height(splitContentH + 2).
		Render(editor)

	splitRow := lipgloss.JoinHorizontal(lipgloss.Top, previewPanel, editorPanel)

	outputsPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(innerW).
		Render(outputs)

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		splitRow,
		outputsPanel,
		footer,
	)

	// Outer frame: Width/Height = total outer dimensions.
	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Width(outerW).
		Height(m.height).
		Render(body)
	// Inject "colr" into the top border line after rendering.
	frame = injectBorderTitle(frame)

	return frame
}

// injectBorderTitle replaces "╭──────" with "╭ colr " in the rendered frame,
// embedding the app name in the top border line without changing visual width.
func injectBorderTitle(s string) string {
	return strings.Replace(s, "╭──────", "╭ colr ", 1)
}
