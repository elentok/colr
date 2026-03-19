package app

import (
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

	if m.showHistory {
		return ui.RenderHistory(m.width, m.height, m.historyEntries, m.historyIndex)
	}

	totalW := m.width
	totalH := m.height

	const (
		headerPanelH  = 5
		outputsPanelH = 6
		footerH       = 1
	)

	bodyH := totalH - headerPanelH - footerH
	if bodyH < 8 {
		bodyH = 8
	}

	rightW := totalW * 38 / 100
	if rightW < 24 {
		rightW = 24
	}
	if rightW > totalW-24 {
		rightW = totalW - 24
	}
	leftW := totalW - rightW
	if leftW < 24 {
		leftW = 24
		rightW = totalW - leftW
	}

	editorPanelH := bodyH - outputsPanelH
	if editorPanelH < 8 {
		editorPanelH = 8
	}
	if outputsPanelH > bodyH-editorPanelH {
		editorPanelH = bodyH - outputsPanelH
	}
	if editorPanelH < 4 {
		editorPanelH = 4
	}
	formatsPanelH := bodyH - editorPanelH
	if formatsPanelH < 6 {
		formatsPanelH = 6
		editorPanelH = bodyH - formatsPanelH
	}
	if editorPanelH < 4 {
		editorPanelH = 4
		formatsPanelH = bodyH - editorPanelH
	}

	headerContent := ui.RenderHeader(
		m.originalClip,
		color.FormatRGB(m.currentColor),
		m.toastMessage,
		totalW-2,
	)
	headerPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(totalW).
		Height(headerPanelH).
		Render(headerContent)

	editorContent := ui.RenderEditor(m.currentColor, m.editMode, m.selectedField, m.lastHue, leftW-2)
	editorPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(leftW).
		Height(editorPanelH).
		Render(editorContent)

	formatsContent := ui.RenderOutputs(m.currentColor, leftW-2)
	formatsPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(leftW).
		Height(formatsPanelH).
		Render(formatsContent)

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, editorPanel, formatsPanel)

	previewContent := ui.RenderPreview(m.currentColor, rightW-2, bodyH-2)
	previewPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(rightW).
		Height(bodyH).
		Render(previewContent)

	bodyRow := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, previewPanel)
	footer := ui.RenderFooter(totalW)

	return lipgloss.JoinVertical(lipgloss.Left, headerPanel, bodyRow, footer)
}
