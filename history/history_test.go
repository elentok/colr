package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elentok/colr/color"
)

func TestRecordMovesDuplicateToFront(t *testing.T) {
	existing := []Entry{
		{RGB: "rgb(255 0 0)", HEX: "#FF0000", Name: "red"},
		{RGB: "rgb(0 0 255)", HEX: "#0000FF", Name: "blue"},
	}

	got := Record(existing, color.Color{R: 255, G: 0, B: 0, A: 1})
	if len(got) != 2 {
		t.Fatalf("Record returned %d entries, want 2", len(got))
	}
	if got[0].HEX != "#FF0000" {
		t.Fatalf("front entry hex = %q, want %q", got[0].HEX, "#FF0000")
	}
	if got[1].HEX != "#0000FF" {
		t.Fatalf("second entry hex = %q, want %q", got[1].HEX, "#0000FF")
	}
}

func TestRecordCapsHistory(t *testing.T) {
	existing := make([]Entry, 0, MaxEntries)
	for i := range MaxEntries {
		existing = append(existing, Entry{
			RGB:  color.FormatRGB(color.Color{R: uint8(i), G: 0, B: 0, A: 1}),
			HEX:  color.FormatHEX(color.Color{R: uint8(i), G: 0, B: 0, A: 1}),
			Name: color.NearestNamedColor(color.Color{R: uint8(i), G: 0, B: 0, A: 1}),
		})
	}

	got := Record(existing, color.Color{R: 200, G: 100, B: 50, A: 1})
	if len(got) != MaxEntries {
		t.Fatalf("Record returned %d entries, want %d", len(got), MaxEntries)
	}
	if got[0].HEX != "#C86432" {
		t.Fatalf("front entry hex = %q, want %q", got[0].HEX, "#C86432")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	want := []Entry{
		{RGB: "rgb(255 0 0)", HEX: "#FF0000", Name: "red"},
		{RGB: "rgb(0 255 0)", HEX: "#00FF00", Name: "lime"},
	}

	if err := Save(want); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("Load returned %d entries, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("entry %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestLoadMigratesLegacyEntries(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	path, err := filePath()
	if err != nil {
		t.Fatalf("filePath returned error: %v", err)
	}

	data := `[{"original":"#ff0000","color":{"R":255,"G":0,"B":0,"A":1}}]`
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Load returned %d entries, want 1", len(got))
	}
	if got[0].HEX != "#FF0000" || got[0].RGB != "rgb(255 0 0)" {
		t.Fatalf("migrated entry = %+v, want normalized red entry", got[0])
	}
}
