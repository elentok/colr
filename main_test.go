package main

import (
	"errors"
	"testing"

	"github.com/elentok/colr/color"
	"github.com/elentok/colr/history"
)

func TestResolveInputPrefersArgument(t *testing.T) {
	t.Helper()

	called := false
	got, err := resolveInput([]string{"colr", "#ff0000"}, func() (string, error) {
		called = true
		return "", nil
	})
	if err != nil {
		t.Fatalf("resolveInput returned error: %v", err)
	}
	if got != "#ff0000" {
		t.Fatalf("resolveInput returned %q, want %q", got, "#ff0000")
	}
	if called {
		t.Fatal("resolveInput should not read the clipboard when an argument is provided")
	}
}

func TestResolveInputJoinsMultipleArguments(t *testing.T) {
	t.Helper()

	got, err := resolveInput([]string{"colr", "255", "0", "0"}, func() (string, error) {
		t.Fatal("resolveInput should not read the clipboard when arguments are provided")
		return "", nil
	})
	if err != nil {
		t.Fatalf("resolveInput returned error: %v", err)
	}
	if got != "255 0 0" {
		t.Fatalf("resolveInput returned %q, want %q", got, "255 0 0")
	}
}

func TestResolveInputFallsBackToClipboard(t *testing.T) {
	t.Helper()

	got, err := resolveInput([]string{"colr"}, func() (string, error) {
		return "rgb(255 0 0)", nil
	})
	if err != nil {
		t.Fatalf("resolveInput returned error: %v", err)
	}
	if got != "rgb(255 0 0)" {
		t.Fatalf("resolveInput returned %q, want %q", got, "rgb(255 0 0)")
	}
}

func TestResolveInputReturnsClipboardError(t *testing.T) {
	t.Helper()

	wantErr := errors.New("clipboard unavailable")
	_, err := resolveInput([]string{"colr"}, func() (string, error) {
		return "", wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("resolveInput error = %v, want %v", err, wantErr)
	}
}

func TestResolveStartupColorUsesParsedInput(t *testing.T) {
	inputText, got, toast, err := resolveStartupColor("button=#ff0000", []history.Entry{
		{RGB: "rgb(0 255 0)", HEX: "#00FF00", Name: "lime"},
	})
	if err != nil {
		t.Fatalf("resolveStartupColor returned error: %v", err)
	}
	if inputText != "button=#ff0000" {
		t.Fatalf("resolveStartupColor inputText = %q, want %q", inputText, "button=#ff0000")
	}
	if got != (color.Color{R: 255, G: 0, B: 0, A: 1}) {
		t.Fatalf("resolveStartupColor color = %+v, want red", got)
	}
	if toast != "" {
		t.Fatalf("resolveStartupColor toast = %q, want empty", toast)
	}
}

func TestResolveStartupColorFallsBackToHistory(t *testing.T) {
	entry := history.Entry{
		RGB:  "rgb(0 255 0)",
		HEX:  "#00FF00",
		Name: "lime",
	}

	inputText, got, toast, err := resolveStartupColor("not a color", []history.Entry{entry})
	if err != nil {
		t.Fatalf("resolveStartupColor returned error: %v", err)
	}
	if inputText != entry.RGB {
		t.Fatalf("resolveStartupColor inputText = %q, want %q", inputText, entry.RGB)
	}
	if got != (color.Color{R: 0, G: 255, B: 0, A: 1}) {
		t.Fatalf("resolveStartupColor color = %+v, want lime", got)
	}
	if toast == "" {
		t.Fatal("resolveStartupColor should return a toast message for history fallback")
	}
}

func TestResolveStartupColorErrorsWithoutHistoryFallback(t *testing.T) {
	_, _, _, err := resolveStartupColor("not a color", nil)
	if !errors.Is(err, errNoColorAvailable) {
		t.Fatalf("resolveStartupColor error = %v, want %v", err, errNoColorAvailable)
	}
}
