package clipboard

import (
	"strings"

	"github.com/atotto/clipboard"
)

// Read reads text from the system clipboard and trims surrounding whitespace.
func Read() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// Write writes text to the system clipboard.
func Write(text string) error {
	return clipboard.WriteAll(text)
}
