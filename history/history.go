package history

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elentok/colr/color"
)

const MaxEntries = 100

// Entry is a persisted color history record, ordered newest-first.
type Entry struct {
	Original string      `json:"original"`
	Color    color.Color `json:"color"`
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
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	if len(entries) > MaxEntries {
		entries = entries[:MaxEntries]
	}

	return entries, nil
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
func Record(entries []Entry, original string, c color.Color) []Entry {
	recorded := make([]Entry, 0, len(entries)+1)
	recorded = append(recorded, Entry{
		Original: original,
		Color:    c,
	})

	newKey := color.FormatHEX(c)
	for _, entry := range entries {
		if color.FormatHEX(entry.Color) == newKey {
			continue
		}
		recorded = append(recorded, entry)
		if len(recorded) == MaxEntries {
			break
		}
	}

	return recorded
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
