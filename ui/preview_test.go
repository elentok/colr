package ui

import (
	"strings"
	"testing"

	"github.com/elentok/colr/color"
)

func TestRenderPreviewShowsOriginalAndEditedSections(t *testing.T) {
	rendered := RenderPreview(
		color.Color{R: 255, G: 0, B: 0, A: 1},
		color.Color{R: 0, G: 0, B: 255, A: 1},
		color.Color{R: 255, G: 255, B: 255, A: 1},
		24,
		12,
		PreviewStacked,
	)

	if !strings.Contains(rendered, "Original") {
		t.Fatal("rendered preview should include an Original section label")
	}
	if !strings.Contains(rendered, "Edited") {
		t.Fatal("rendered preview should include an Edited section label")
	}
	if strings.Count(rendered, "over-bg") != 2 {
		t.Fatalf("rendered preview should contain two over-bg hint lines, got %d", strings.Count(rendered, "over-bg"))
	}
}

func TestRenderPreviewSideBySideShowsOriginalAndEditedSections(t *testing.T) {
	rendered := RenderPreview(
		color.Color{R: 255, G: 0, B: 0, A: 1},
		color.Color{R: 0, G: 0, B: 255, A: 1},
		color.Color{R: 0, G: 0, B: 0, A: 1},
		40,
		8,
		PreviewSideBySide,
	)

	if !strings.Contains(rendered, "Original") {
		t.Fatal("side-by-side preview should include an Original section label")
	}
	if !strings.Contains(rendered, "Edited") {
		t.Fatal("side-by-side preview should include an Edited section label")
	}
	if strings.Count(rendered, "over-bg") != 2 {
		t.Fatalf("side-by-side preview should contain two over-bg hint lines, got %d", strings.Count(rendered, "over-bg"))
	}
}
