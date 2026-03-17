# Milestone 1 — Parsing and Formatting Library

## Goal

Implement the pure color logic layer: parsing CSS color strings, internal RGBA model, color space conversions, and normalized output formatting. No TUI code — this is the foundation everything else builds on.

## Dependencies

- Go 1.26
- Module path: `github.com/elentok/colr`
- `github.com/atotto/clipboard` (for clipboard wrapper only)
- No other external dependencies for this milestone

## Package Structure

```
colr/
  go.mod
  main.go              # placeholder, just enough to compile

  color/
    color.go           # Color struct (RGBA), HSV struct, HSL struct
    parse.go           # parse clipboard string -> Color
    parse_test.go      # comprehensive parser tests
    format.go          # Color -> normalized RGB/HEX/HSL strings
    format_test.go     # formatting tests
    convert.go         # RGB<->HSV, RGB<->HSL conversions
    convert_test.go    # conversion tests

  clipboard/
    clipboard.go       # thin wrapper around atotto/clipboard
```

## Implementation Steps

### Step 1: Project scaffolding

- [ ] `go mod init github.com/elentok/colr`
- [ ] Create directory structure: `color/`, `clipboard/`
- [ ] Minimal `main.go` (just `package main` + `func main()`)
- [ ] `go mod tidy`

### Step 2: Data types (`color/color.go`)

- [ ] `Color` struct: `R uint8`, `G uint8`, `B uint8`, `A float64` (0.0–1.0)
- [ ] `HSV` struct: `H float64` (0–360), `S float64` (0–1), `V float64` (0–1)
- [ ] `HSL` struct: `H float64` (0–360), `S float64` (0–1), `L float64` (0–1)

### Step 3: Conversions (`color/convert.go`)

- [ ] `RGBToHSV(Color) HSV`
- [ ] `HSVToRGB(HSV, alpha float64) Color`
- [ ] `RGBToHSL(Color) HSL`
- [ ] `HSLToRGB(HSL, alpha float64) Color`
- [ ] Handle grayscale edge case (S=0 → H=0, deterministic)

### Step 4: Conversion tests (`color/convert_test.go`)

- [ ] Round-trip tests: RGB→HSV→RGB for representative colors
- [ ] Round-trip tests: RGB→HSL→RGB for representative colors
- [ ] Representative colors: black, white, red, green, blue, gray, semi-transparent red
- [ ] Grayscale stability: converting gray should yield H=0
- [ ] Edge cases: pure black (V=0), pure white (V=1,S=0)

### Step 5: Parser (`color/parse.go`)

- [ ] `Parse(input string) (Color, error)` — main entry point
- [ ] Trim whitespace and semicolons from input
- [ ] Try each sub-parser in order, return first match
- [ ] Sub-parsers:
  - [ ] **HEX**: `#RRGGBB`, `#RRGGBBAA`, `RRGGBB`, `RRGGBBAA` (case-insensitive)
  - [ ] **rgb(...)**: space-separated `rgb(R G B)`, with optional alpha `rgb(R G B / A)`
  - [ ] **rgb(...)**: comma-separated `rgb(R, G, B)`
  - [ ] **rgb(...)**: percentage values `rgb(R% G% B%)`
  - [ ] **rgba(...)**: comma-separated with alpha `rgba(R, G, B, A)`
  - [ ] **hsl(...)**: `hsl(H S% L%)`, with optional alpha `hsl(H S% L% / A)`
  - [ ] **Bare values**: `R, G, B` and `R G B` and `R% G% B%`
- [ ] Alpha parsing rules:
  - [ ] `50%` → 0.5
  - [ ] `0.5` → 0.5
  - [ ] `.5` → 0.5
  - [ ] No alpha specified → 1.0
- [ ] RGB percentage: `100%` → 255, `0%` → 0
- [ ] Reject invalid inputs with descriptive errors

### Step 6: Parser tests (`color/parse_test.go`)

- [ ] **HEX valid**: `#FF0000`, `FF0000`, `#ff0000`, `ff0000`, `#FF000080`, `FF000080`, `#a3a3a3`, `a3a3a3aa`
- [ ] **RGB valid**: `rgb(255 0 0)`, `rgb(255 0 0 / 50%)`, `rgb(255 0 0 / 0.5)`, `rgb(255, 0, 0)`, `rgb(100% 0% 0%)`
- [ ] **RGBA valid**: `rgba(255, 0, 0, 50%)`, `rgba(255, 0, 0, 0.5)`, `rgba(255, 0, 0, .5)`
- [ ] **HSL valid**: `hsl(0 100% 50%)`, `hsl(0 100% 50% / 50%)`
- [ ] **Bare valid**: `255, 0, 0`, `255 0 0`, `100% 0% 0%`
- [ ] **Invalid**: `hello`, `rgb(foo bar baz)`, `rgb(10 20)`, `#XYZ123`, `hsl(20 30)`, empty string
- [ ] **Whitespace**: `  #FF0000  `, `  rgb(255 0 0)  `
- [ ] **Alpha edge cases**: alpha 0%, alpha 100%, `.5`, `0.0`, `1.0`

### Step 7: Formatter (`color/format.go`)

- [ ] `FormatRGB(Color) string` — e.g. `rgb(255 0 0)` or `rgb(255 0 0 / 50%)`
- [ ] `FormatHEX(Color) string` — e.g. `#FF0000` or `#FF000080`
- [ ] `FormatHSL(Color) string` — e.g. `hsl(0 100% 50%)` or `hsl(0 100% 50% / 50%)`
- [ ] Alpha rules:
  - [ ] Omit alpha when A == 1.0
  - [ ] Display as integer percentage when A < 1.0
- [ ] Rounding: RGB channels nearest int, HSL H/S/L nearest int, alpha nearest int %

### Step 8: Formatter tests (`color/format_test.go`)

- [ ] Opaque red → `rgb(255 0 0)`, `#FF0000`, `hsl(0 100% 50%)`
- [ ] Semi-transparent red (A=0.5) → `rgb(255 0 0 / 50%)`, `#FF000080`, `hsl(0 100% 50% / 50%)`
- [ ] Black → `rgb(0 0 0)`, `#000000`, `hsl(0 0% 0%)`
- [ ] White → `rgb(255 255 255)`, `#FFFFFF`, `hsl(0 0% 100%)`
- [ ] Gray → appropriate values
- [ ] Alpha 0% → included in output
- [ ] Alpha 100% → omitted from output

### Step 9: Clipboard wrapper (`clipboard/clipboard.go`)

- [ ] `Read() (string, error)` — wraps `atotto/clipboard.ReadAll()`, trims whitespace
- [ ] `Write(text string) error` — wraps `atotto/clipboard.WriteAll()`

### Step 10: Integration smoke test

- [ ] Verify: `clipboard.Read()` → `color.Parse()` → `color.FormatRGB/HEX/HSL()` pipeline works
- [ ] Run full test suite, all green

## Acceptance Criteria

- [ ] Can parse all listed valid input examples from the spec
- [ ] Rejects all listed invalid examples with errors
- [ ] Outputs normalized RGB / HEX / HSL strings correctly
- [ ] Alpha handled correctly throughout (parsing, conversion, formatting)
- [ ] RGB↔HSV and RGB↔HSL round-trips are stable for representative colors
- [ ] All tests pass
