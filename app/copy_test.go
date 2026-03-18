package app

import (
	"strings"
	"testing"

	"github.com/elentok/colr/color"
)

func TestApplyCopyToastMessages(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	tests := []struct {
		format      string
		wantContain string
	}{
		{"rgb", "RGB"},
		{"hex", "HEX"},
		{"hsl", "HSL"},
	}

	for _, tc := range tests {
		t.Run(tc.format, func(t *testing.T) {
			updated, cmd := applyCopy(m, tc.format)
			// Toast message should mention the format label.
			if !strings.Contains(updated.toastMessage, tc.wantContain) {
				t.Errorf("toastMessage %q does not contain %q", updated.toastMessage, tc.wantContain)
			}
			// A tick command should be scheduled.
			if cmd == nil {
				t.Error("expected a tick cmd, got nil")
			}
		})
	}
}

func TestApplyCopyReturnsCmdForTimer(t *testing.T) {
	m := newTestModel(color.Color{R: 0, G: 128, B: 255, A: 0.5})
	_, cmd := applyCopy(m, "hex")
	if cmd == nil {
		t.Error("expected non-nil cmd from applyCopy")
	}
}

func TestClearToastMsg(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.toastMessage = "Copied RGB to clipboard"

	// Simulate ClearToastMsg being processed.
	result, _ := m.Update(ClearToastMsg{})
	updated := result.(Model)
	if updated.toastMessage != "" {
		t.Errorf("expected empty toast after ClearToastMsg, got %q", updated.toastMessage)
	}
}

func TestYPrefixCopyRGB(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	// Press y
	result, _ := handleKeyMsg(m, "y")
	updated := result.(Model)
	if !updated.pendingY {
		t.Error("expected pendingY=true after 'y'")
	}

	// Press r
	result2, cmd := handleKeyMsg(updated, "r")
	updated2 := result2.(Model)
	if updated2.pendingY {
		t.Error("expected pendingY=false after 'r'")
	}
	if !strings.Contains(updated2.toastMessage, "RGB") {
		t.Errorf("expected RGB toast, got %q", updated2.toastMessage)
	}
	if cmd == nil {
		t.Error("expected tick cmd after copy")
	}
}

func TestYPrefixCopyHEX(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.pendingY = true

	result, _ := handleKeyMsg(m, "x")
	updated := result.(Model)
	if !strings.Contains(updated.toastMessage, "HEX") {
		t.Errorf("expected HEX toast, got %q", updated.toastMessage)
	}
}

func TestYPrefixCopyHSL(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.pendingY = true

	result, _ := handleKeyMsg(m, "h")
	updated := result.(Model)
	if !strings.Contains(updated.toastMessage, "HSL") {
		t.Errorf("expected HSL toast, got %q", updated.toastMessage)
	}
}

func TestYPrefixUnknownKeyNoop(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.pendingY = true

	result, cmd := handleKeyMsg(m, "z")
	updated := result.(Model)
	if updated.pendingY {
		t.Error("expected pendingY cleared after unknown key")
	}
	if updated.toastMessage != "" {
		t.Errorf("expected no toast after unknown key, got %q", updated.toastMessage)
	}
	if cmd != nil {
		t.Error("expected nil cmd after unknown key")
	}
}

func TestHelpToggle(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	// Open help
	result, _ := handleKeyMsg(m, "?")
	updated := result.(Model)
	if !updated.showHelp {
		t.Error("expected showHelp=true after '?'")
	}

	// Close help with ?
	result2, _ := handleKeyMsg(updated, "?")
	updated2 := result2.(Model)
	if updated2.showHelp {
		t.Error("expected showHelp=false after second '?'")
	}
}

func TestHelpDismissEsc(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.showHelp = true

	result, _ := handleKeyMsg(m, "esc")
	updated := result.(Model)
	if updated.showHelp {
		t.Error("expected showHelp=false after Esc")
	}
}

func TestHelpQClosesHelpNotApp(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.showHelp = true

	result, cmd := handleKeyMsg(m, "q")
	updated := result.(Model)
	if updated.showHelp {
		t.Error("expected showHelp=false after q in help")
	}
	if cmd != nil {
		t.Error("expected no quit cmd when q pressed in help overlay")
	}
}

func TestYYCopiesRGB(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	result, _ := handleKeyMsg(m, "y")
	updated := result.(Model)
	result2, cmd := handleKeyMsg(updated, "y")
	updated2 := result2.(Model)

	if !strings.Contains(updated2.toastMessage, "RGB") {
		t.Errorf("expected RGB toast from yy, got %q", updated2.toastMessage)
	}
	if cmd == nil {
		t.Error("expected tick cmd from yy")
	}
}

func TestHelpBlocksOtherKeys(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	m.showHelp = true
	original := m.selectedField

	// j should NOT move field when help is open
	result, _ := handleKeyMsg(m, "j")
	updated := result.(Model)
	if updated.selectedField != original {
		t.Errorf("expected field unchanged while help open, got %d", updated.selectedField)
	}
}
