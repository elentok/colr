package ui

import (
	"strings"
	"testing"

	"github.com/elentok/colr/color"
)

func TestRenderOutputsShowsLowercaseHex(t *testing.T) {
	rendered := RenderOutputs(
		color.Color{R: 255, G: 128, B: 0, A: 1},
		color.Color{R: 255, G: 255, B: 255, A: 1},
		48,
	)

	if !strings.Contains(rendered, "#ff8000") {
		t.Fatalf("rendered outputs should contain lowercase hex, got %q", rendered)
	}
	if strings.Contains(rendered, "#FF8000") {
		t.Fatalf("rendered outputs should not contain uppercase hex, got %q", rendered)
	}
}

func TestRenderOutputsShowsOverBackgroundRow(t *testing.T) {
	rendered := RenderOutputs(
		color.Color{R: 255, G: 0, B: 0, A: 0.5},
		color.Color{R: 255, G: 255, B: 255, A: 1},
		52,
	)

	if !strings.Contains(rendered, "OVER") {
		t.Fatalf("rendered outputs should contain OVER row, got %q", rendered)
	}
	if !strings.Contains(rendered, "#ff8080 on white") {
		t.Fatalf("rendered outputs should show composited value, got %q", rendered)
	}
	if !strings.Contains(rendered, "[yo]") {
		t.Fatalf("rendered outputs should show yo key hint, got %q", rendered)
	}
}
