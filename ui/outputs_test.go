package ui

import (
	"strings"
	"testing"

	"github.com/elentok/colr/color"
)

func TestRenderOutputsShowsLowercaseHex(t *testing.T) {
	rendered := RenderOutputs(color.Color{R: 255, G: 128, B: 0, A: 1}, 48)

	if !strings.Contains(rendered, "#ff8000") {
		t.Fatalf("rendered outputs should contain lowercase hex, got %q", rendered)
	}
	if strings.Contains(rendered, "#FF8000") {
		t.Fatalf("rendered outputs should not contain uppercase hex, got %q", rendered)
	}
}
