package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/elentok/colr/color"
	"github.com/elentok/colr/history"
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

func TestPOpensHistoryOverlay(t *testing.T) {
	m := newHistoryTestModel(color.Color{R: 255, G: 0, B: 0, A: 1}, []history.Entry{
		{Original: "#00ff00", Color: color.Color{R: 0, G: 255, B: 0, A: 1}},
	})

	result, _ := handleKeyMsg(m, "p")
	updated := result.(Model)
	if !updated.showHistory {
		t.Error("expected showHistory=true after 'p'")
	}
}

func TestHistoryOverlayLoadsSelectedEntry(t *testing.T) {
	entry := history.Entry{Original: "#00ff00", Color: color.Color{R: 0, G: 255, B: 0, A: 1}}
	m := newHistoryTestModel(color.Color{R: 255, G: 0, B: 0, A: 1}, []history.Entry{entry})
	m.showHistory = true

	result, _ := handleKeyMsg(m, "enter")
	updated := result.(Model)
	if updated.showHistory {
		t.Error("expected showHistory=false after loading a history color")
	}
	if updated.currentColor != entry.Color {
		t.Errorf("currentColor = %+v, want %+v", updated.currentColor, entry.Color)
	}
	if updated.originalClip != entry.Original {
		t.Errorf("originalClip = %q, want %q", updated.originalClip, entry.Original)
	}
}

func TestHistoryOverlayConsumesEditingKeys(t *testing.T) {
	m := newHistoryTestModel(color.Color{R: 255, G: 0, B: 0, A: 1}, []history.Entry{
		{Original: "#00ff00", Color: color.Color{R: 0, G: 255, B: 0, A: 1}},
		{Original: "#0000ff", Color: color.Color{R: 0, G: 0, B: 255, A: 1}},
	})
	m.showHistory = true
	m.selectedField = FieldOpacity

	result, _ := handleKeyMsg(m, "j")
	updated := result.(Model)
	if updated.historyIndex != 1 {
		t.Errorf("historyIndex = %d, want 1", updated.historyIndex)
	}
	if updated.selectedField != FieldOpacity {
		t.Errorf("selectedField changed while history open: got %d want %d", updated.selectedField, FieldOpacity)
	}
}

func TestSPressesStartsHistorySave(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	result, cmd := handleKeyMsg(m, "s")
	updated := result.(Model)
	if cmd == nil {
		t.Fatal("expected save history command from 's'")
	}
	if updated.pendingY {
		t.Fatal("expected pendingY cleared when saving history")
	}
}

func TestSaveHistoryMsgUpdatesEntriesAndToast(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})
	entries := []history.Entry{
		{Original: "rgb(255 0 0)", Color: color.Color{R: 255, G: 0, B: 0, A: 1}},
	}

	result, _ := m.Update(SaveHistoryMsg{entries: entries})
	updated := result.(Model)
	if len(updated.historyEntries) != 1 {
		t.Fatalf("historyEntries length = %d, want 1", len(updated.historyEntries))
	}
	if updated.toastMessage != "Saved color to history" {
		t.Fatalf("toastMessage = %q, want %q", updated.toastMessage, "Saved color to history")
	}
}

func TestSaveHistoryMsgFailureSetsErrorToast(t *testing.T) {
	m := newTestModel(color.Color{R: 255, G: 0, B: 0, A: 1})

	result, _ := m.Update(SaveHistoryMsg{err: errors.New("write failed")})
	updated := result.(Model)
	if updated.toastMessage != "Failed to save color to history" {
		t.Fatalf("toastMessage = %q, want %q", updated.toastMessage, "Failed to save color to history")
	}
}
