package color

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var hexRe = regexp.MustCompile(`^#?([0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
var hexSearchRe = regexp.MustCompile(`#?[0-9a-fA-F]{6}(?:[0-9a-fA-F]{2})?`)
var rgbSearchRe = regexp.MustCompile(`(?i)rgba?\([^)]*\)`)
var hslSearchRe = regexp.MustCompile(`(?i)hsla?\([^)]*\)`)
var bareSearchRe = regexp.MustCompile(`(?:\d+(?:\.\d+)?%?\s*,\s*){2}\d+(?:\.\d+)?%?|(?:\d+(?:\.\d+)?%?\s+){2}\d+(?:\.\d+)?%?`)

type candidateMatch struct {
	start int
	text  string
}

// Parse parses a CSS color string into a Color.
// Trims surrounding whitespace and trailing semicolons.
func Parse(input string) (Color, error) {
	input = strings.TrimSpace(input)
	input = strings.TrimRight(input, ";")
	input = strings.TrimSpace(input)

	if input == "" {
		return Color{}, fmt.Errorf("empty input")
	}

	for _, p := range []func(string) (Color, error){
		parseHEX,
		parseRGBFunc,
		parseHSLFunc,
		parseBare,
	} {
		if c, err := p(input); err == nil {
			return c, nil
		}
	}

	return Color{}, fmt.Errorf("unrecognized color format")
}

// FindFirst parses the first valid color found anywhere in the input text.
func FindFirst(input string) (Color, error) {
	if c, err := Parse(input); err == nil {
		return c, nil
	}

	candidates := collectCandidates(input)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].start != candidates[j].start {
			return candidates[i].start < candidates[j].start
		}
		return len(candidates[i].text) > len(candidates[j].text)
	})

	for _, candidate := range candidates {
		if c, err := Parse(candidate.text); err == nil {
			return c, nil
		}
	}

	return Color{}, fmt.Errorf("unrecognized color format")
}

func collectCandidates(input string) []candidateMatch {
	candidates := make([]candidateMatch, 0)

	for _, idx := range rgbSearchRe.FindAllStringIndex(input, -1) {
		candidates = append(candidates, candidateMatch{
			start: idx[0],
			text:  input[idx[0]:idx[1]],
		})
	}

	for _, idx := range hslSearchRe.FindAllStringIndex(input, -1) {
		candidates = append(candidates, candidateMatch{
			start: idx[0],
			text:  input[idx[0]:idx[1]],
		})
	}

	for _, idx := range bareSearchRe.FindAllStringIndex(input, -1) {
		if !isBareBoundary(input, idx[0], idx[1]) {
			continue
		}
		candidates = append(candidates, candidateMatch{
			start: idx[0],
			text:  input[idx[0]:idx[1]],
		})
	}

	for _, idx := range hexSearchRe.FindAllStringIndex(input, -1) {
		if !isHexBoundary(input, idx[0], idx[1]) {
			continue
		}
		candidates = append(candidates, candidateMatch{
			start: idx[0],
			text:  input[idx[0]:idx[1]],
		})
	}

	return candidates
}

func isHexBoundary(input string, start, end int) bool {
	if start > 0 {
		prev := input[start-1]
		if isHexDigit(prev) || prev == '#' {
			return false
		}
	}

	if end < len(input) && isHexDigit(input[end]) {
		return false
	}

	return true
}

func isBareBoundary(input string, start, end int) bool {
	if start > 0 && isBareTokenChar(input[start-1]) {
		return false
	}

	if end < len(input) && isBareTokenChar(input[end]) {
		return false
	}

	return true
}

func isHexDigit(b byte) bool {
	return ('0' <= b && b <= '9') || ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F')
}

func isBareTokenChar(b byte) bool {
	return ('0' <= b && b <= '9') || ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == '.' || b == '%'
}

// parseHEX handles #RRGGBB, RRGGBB, #RRGGBBAA, RRGGBBAA (case-insensitive).
func parseHEX(s string) (Color, error) {
	if !hexRe.MatchString(s) {
		return Color{}, fmt.Errorf("not a hex color")
	}

	hex := strings.TrimPrefix(s, "#")
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	a := 1.0
	if len(hex) == 8 {
		av, _ := strconv.ParseUint(hex[6:8], 16, 8)
		a = float64(av) / 255.0
	}

	return Color{R: uint8(r), G: uint8(g), B: uint8(b), A: a}, nil
}

// parseRGBFunc handles rgb(...) and rgba(...).
func parseRGBFunc(s string) (Color, error) {
	lower := strings.ToLower(s)
	isRGBA := strings.HasPrefix(lower, "rgba(")
	isRGB := strings.HasPrefix(lower, "rgb(")

	if !isRGB && !isRGBA {
		return Color{}, fmt.Errorf("not rgb/rgba")
	}

	if !strings.HasSuffix(strings.TrimSpace(lower), ")") {
		return Color{}, fmt.Errorf("missing closing parenthesis")
	}

	start := strings.Index(s, "(")
	end := strings.LastIndex(s, ")")
	content := strings.TrimSpace(s[start+1 : end])

	if strings.Contains(content, ",") {
		return parseRGBComma(content, isRGBA)
	}
	return parseRGBSpace(content)
}

// parseRGBComma parses comma-separated rgb/rgba: rgb(R, G, B) or rgba(R, G, B, A).
func parseRGBComma(content string, isRGBA bool) (Color, error) {
	parts := strings.Split(content, ",")

	alpha := 1.0
	if isRGBA {
		if len(parts) != 4 {
			return Color{}, fmt.Errorf("rgba requires 4 components, got %d", len(parts))
		}
		var err error
		alpha, err = parseAlpha(strings.TrimSpace(parts[3]))
		if err != nil {
			return Color{}, fmt.Errorf("invalid alpha: %w", err)
		}
		parts = parts[:3]
	} else {
		if len(parts) != 3 {
			return Color{}, fmt.Errorf("rgb requires 3 components, got %d", len(parts))
		}
	}

	return parseThreeChannels(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), alpha)
}

// parseRGBSpace parses space-separated rgb: rgb(R G B) or rgb(R G B / A).
func parseRGBSpace(content string) (Color, error) {
	alpha := 1.0

	if idx := strings.Index(content, "/"); idx != -1 {
		var err error
		alpha, err = parseAlpha(strings.TrimSpace(content[idx+1:]))
		if err != nil {
			return Color{}, fmt.Errorf("invalid alpha: %w", err)
		}
		content = strings.TrimSpace(content[:idx])
	}

	parts := strings.Fields(content)
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("rgb requires 3 components, got %d", len(parts))
	}

	return parseThreeChannels(parts[0], parts[1], parts[2], alpha)
}

// parseHSLFunc handles hsl(...) and hsla(...).
func parseHSLFunc(s string) (Color, error) {
	lower := strings.ToLower(s)
	if !strings.HasPrefix(lower, "hsl(") && !strings.HasPrefix(lower, "hsla(") {
		return Color{}, fmt.Errorf("not hsl/hsla")
	}

	if !strings.HasSuffix(strings.TrimSpace(lower), ")") {
		return Color{}, fmt.Errorf("missing closing parenthesis")
	}

	start := strings.Index(s, "(")
	end := strings.LastIndex(s, ")")
	content := strings.TrimSpace(s[start+1 : end])

	alpha := 1.0
	if idx := strings.Index(content, "/"); idx != -1 {
		var err error
		alpha, err = parseAlpha(strings.TrimSpace(content[idx+1:]))
		if err != nil {
			return Color{}, fmt.Errorf("invalid alpha: %w", err)
		}
		content = strings.TrimSpace(content[:idx])
	}

	parts := strings.Fields(content)
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("hsl requires 3 components, got %d", len(parts))
	}

	h, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hue %q: %w", parts[0], err)
	}

	sv, err := parsePercent(parts[1])
	if err != nil {
		return Color{}, fmt.Errorf("invalid saturation %q: %w", parts[1], err)
	}

	lv, err := parsePercent(parts[2])
	if err != nil {
		return Color{}, fmt.Errorf("invalid lightness %q: %w", parts[2], err)
	}

	return HSLToRGB(HSL{H: h, S: sv / 100.0, L: lv / 100.0}, alpha), nil
}

// parseBare handles bare values: "R G B", "R, G, B", "R% G% B%".
func parseBare(s string) (Color, error) {
	if strings.Contains(s, ",") {
		parts := strings.Split(s, ",")
		if len(parts) != 3 {
			return Color{}, fmt.Errorf("bare: expected 3 comma-separated values")
		}
		return parseThreeChannels(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), 1.0)
	}

	parts := strings.Fields(s)
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("bare: expected 3 space-separated values")
	}
	return parseThreeChannels(parts[0], parts[1], parts[2], 1.0)
}

// parseThreeChannels parses three RGB channel strings with a given alpha.
func parseThreeChannels(rs, gs, bs string, alpha float64) (Color, error) {
	r, err := parseChannel(rs)
	if err != nil {
		return Color{}, fmt.Errorf("invalid red channel: %w", err)
	}
	g, err := parseChannel(gs)
	if err != nil {
		return Color{}, fmt.Errorf("invalid green channel: %w", err)
	}
	b, err := parseChannel(bs)
	if err != nil {
		return Color{}, fmt.Errorf("invalid blue channel: %w", err)
	}
	return Color{R: r, G: g, B: b, A: alpha}, nil
}

// parseChannel parses a single RGB channel value (integer 0-255 or percentage).
// Values outside range are clamped.
func parseChannel(s string) (uint8, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid percentage %q: %w", s, err)
		}
		return clampUint8(math.Round(v * 255.0 / 100.0)), nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q: %w", s, err)
	}
	return clampUint8(math.Round(v)), nil
}

// parseAlpha parses an alpha value from a percentage (50%) or decimal (0.5, .5).
// Result is clamped to [0.0, 1.0].
func parseAlpha(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid alpha percentage %q: %w", s, err)
		}
		return math.Max(0, math.Min(1, v/100.0)), nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid alpha %q: %w", s, err)
	}
	return math.Max(0, math.Min(1, v)), nil
}

// parsePercent parses a percentage string like "50%" and returns 50.0.
// The % suffix is optional.
func parsePercent(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid percentage %q: %w", s, err)
	}
	return v, nil
}

// clampUint8 rounds and clamps a float64 to uint8 range.
func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
