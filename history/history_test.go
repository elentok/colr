package history

import (
	"testing"

	"github.com/elentok/colr/color"
)

func TestRecordMovesDuplicateToFront(t *testing.T) {
	existing := []Entry{
		{Original: "old red", Color: color.Color{R: 255, G: 0, B: 0, A: 1}},
		{Original: "blue", Color: color.Color{R: 0, G: 0, B: 255, A: 1}},
	}

	got := Record(existing, "new red", color.Color{R: 255, G: 0, B: 0, A: 1})
	if len(got) != 2 {
		t.Fatalf("Record returned %d entries, want 2", len(got))
	}
	if got[0].Original != "new red" {
		t.Fatalf("front entry original = %q, want %q", got[0].Original, "new red")
	}
	if got[1].Original != "blue" {
		t.Fatalf("second entry original = %q, want %q", got[1].Original, "blue")
	}
}

func TestRecordCapsHistory(t *testing.T) {
	existing := make([]Entry, 0, MaxEntries)
	for i := range MaxEntries {
		existing = append(existing, Entry{
			Original: "entry",
			Color:    color.Color{R: uint8(i), G: 0, B: 0, A: 1},
		})
	}

	got := Record(existing, "new", color.Color{R: 200, G: 100, B: 50, A: 1})
	if len(got) != MaxEntries {
		t.Fatalf("Record returned %d entries, want %d", len(got), MaxEntries)
	}
	if got[0].Original != "new" {
		t.Fatalf("front entry original = %q, want %q", got[0].Original, "new")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	want := []Entry{
		{Original: "#ff0000", Color: color.Color{R: 255, G: 0, B: 0, A: 1}},
		{Original: "rgb(0 255 0)", Color: color.Color{R: 0, G: 255, B: 0, A: 1}},
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
		if got[i].Original != want[i].Original || got[i].Color != want[i].Color {
			t.Fatalf("entry %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}
