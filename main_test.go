package main

import (
	"errors"
	"testing"
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
