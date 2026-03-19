package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/elentok/colr/color"
)

const MaxEntries = 100

// Entry is a persisted color history record, ordered newest-first.
type Entry struct {
	RGB  string `json:"rgb"`
	HEX  string `json:"hex"`
	Name string `json:"name"`
}

// Load reads the stored history file. Missing files return an empty history.
func Load() ([]Entry, error) {
	path, err := filePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err == nil {
		if normalizedEntriesValid(entries) {
			return clampEntries(entries), nil
		}
	}

	var legacyEntries []legacyEntry
	if err := json.Unmarshal(data, &legacyEntries); err != nil {
		return nil, err
	}

	entries = make([]Entry, 0, len(legacyEntries))
	for _, legacy := range legacyEntries {
		entry, err := NewEntry(legacy.Color)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return clampEntries(entries), nil
}

// Save writes the full history file, ensuring the parent directory exists.
func Save(entries []Entry) error {
	path, err := filePath()
	if err != nil {
		return err
	}

	if len(entries) > MaxEntries {
		entries = entries[:MaxEntries]
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

// Record inserts a color at the front of the history, removing older duplicates.
func Record(entries []Entry, c color.Color) []Entry {
	entry, err := NewEntry(c)
	if err != nil {
		return clampEntries(entries)
	}

	recorded := make([]Entry, 0, len(entries)+1)
	recorded = append(recorded, entry)

	for _, existing := range entries {
		if existing.HEX == entry.HEX {
			continue
		}
		recorded = append(recorded, existing)
		if len(recorded) == MaxEntries {
			break
		}
	}

	return recorded
}

// NewEntry builds a normalized history entry from a color.
func NewEntry(c color.Color) (Entry, error) {
	return Entry{
		RGB:  color.FormatRGB(c),
		HEX:  color.FormatHEX(c),
		Name: color.NearestNamedColor(c),
	}, nil
}

// Color reconstructs the saved color from the normalized hex form.
func (e Entry) Color() (color.Color, error) {
	if e.HEX == "" {
		return color.Color{}, fmt.Errorf("history entry is missing hex")
	}
	return color.Parse(e.HEX)
}

type legacyEntry struct {
	Color color.Color `json:"color"`
}

func clampEntries(entries []Entry) []Entry {
	if len(entries) > MaxEntries {
		return entries[:MaxEntries]
	}
	return entries
}

func normalizedEntriesValid(entries []Entry) bool {
	for _, entry := range entries {
		if entry.RGB == "" || entry.HEX == "" || entry.Name == "" {
			return false
		}
	}
	return true
}

func filePath() (string, error) {
	stateDir, err := stateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(stateDir, "colr", "history.json"), nil
}

func stateDir() (string, error) {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".local", "state"), nil
}
