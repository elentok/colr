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

	footer := ui.RenderFooter(totalW)
	isPortrait := usePortraitLayout(totalW, totalH)

	if isPortrait {
		previewPanelH := bodyH * 38 / 100
		if previewPanelH < 8 {
			previewPanelH = 8
		}
		if previewPanelH > bodyH-10 {
			previewPanelH = bodyH - 10
		}
		topSectionH := bodyH - previewPanelH
		if topSectionH < 8 {
			topSectionH = 8
			previewPanelH = bodyH - topSectionH
		}

		formatsPanelH := outputsPanelH
		editorPanelH := topSectionH - formatsPanelH
		if editorPanelH < 6 {
			editorPanelH = 6
			formatsPanelH = topSectionH - editorPanelH
		}

		editorContent := ui.RenderEditor(m.currentColor, m.editMode, m.selectedField, m.lastHue, totalW-2)
		editorPanel := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Width(totalW).
			Height(editorPanelH).
			Render(editorContent)

		formatsContent := ui.RenderOutputs(m.currentColor, totalW-2)
		formatsPanel := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Width(totalW).
			Height(formatsPanelH).
			Render(formatsContent)

		previewContent := ui.RenderPreview(m.originalColor, m.currentColor, totalW-2, previewPanelH-2, ui.PreviewSideBySide)
		previewPanel := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Width(totalW).
			Height(previewPanelH).
			Render(previewContent)

		body := lipgloss.JoinVertical(lipgloss.Left, editorPanel, formatsPanel, previewPanel)
		return lipgloss.JoinVertical(lipgloss.Left, headerPanel, body, footer)
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

	previewContent := ui.RenderPreview(m.originalColor, m.currentColor, rightW-2, bodyH-2, ui.PreviewStacked)
	previewPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(rightW).
		Height(bodyH).
		Render(previewContent)

	bodyRow := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, previewPanel)
	return lipgloss.JoinVertical(lipgloss.Left, headerPanel, bodyRow, footer)
}

func usePortraitLayout(width, height int) bool {
	// Terminal cells are visually taller than they are wide, so a raw width<height
	// comparison waits too long before switching to the portrait layout.
	return width <= height*2
}
