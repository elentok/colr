# Milestone 2 — Minimal TUI Shell

## Goal

Get a Bubble Tea TUI that boots from clipboard input, renders the full 5-region layout (header, preview, editor, outputs, footer), and displays all values read-only. No keyboard editing yet — that's milestone 3. The focus here is startup flow, layout, and rendering.

## Dependencies

- Bubble Tea (`github.com/charmbracelet/bubbletea`)
- Lip Gloss (`github.com/charmbracelet/lipgloss`)
- Existing `color/` and `clipboard/` packages from milestone 1

## Package Structure (new files)

```
app/
  model.go      # Model struct, Init, initial state
  update.go     # Update loop (q to quit only for now)
  view.go       # View — composes all 5 regions
  types.go      # EditMode enum, field index constants

ui/
  styles.go     # Lip Gloss styles: borders, colors, dimensions
  header.go     # Header region renderer
  preview.go    # Color preview block renderer
  editor.go     # Editor pane renderer (read-only display for now)
  outputs.go    # Output pane renderer (RGB/HEX/HSL rows)
  footer.go     # Footer key hints renderer

main.go         # Updated: clipboard read → parse → launch TUI or exit
```

## Implementation Steps

### Step 1: Dependencies

- [ ] `go get github.com/charmbracelet/bubbletea`
- [ ] `go get github.com/charmbracelet/lipgloss`

### Step 2: App types (`app/types.go`)

- [ ] `EditMode` type (int): `ModeHSV`, `ModeRGB`
- [ ] `Field` type (int): field index constants per mode
  - HSV: `FieldHue`, `FieldSaturation`, `FieldValue`, `FieldOpacity`
  - RGB: `FieldRed`, `FieldGreen`, `FieldBlue`, `FieldOpacity`
- [ ] Constants for field count per mode (both are 4)

### Step 3: App model (`app/model.go`)

- [ ] `Model` struct containing:
  - `originalClip string` — raw clipboard text
  - `originalColor color.Color` — parsed original
  - `currentColor color.Color` — working copy (same as original for now)
  - `editMode EditMode` — HSV or RGB, default HSV
  - `selectedField int` — 0-based index of focused field
  - `toastMessage string` — status/toast text
  - `toastExpiry time.Time` — when toast disappears
  - `width, height int` — terminal dimensions
  - `showHelp bool` — help overlay toggle
- [ ] `NewModel(clipText string, c color.Color) Model` constructor
- [ ] `Init() tea.Cmd` — return `tea.WindowSize` or nil

### Step 4: App update (`app/update.go`)

- [ ] Handle `tea.KeyMsg`:
  - `q` / `ctrl+c` → `tea.Quit`
  - (Other keys will be added in milestone 3)
- [ ] Handle `tea.WindowSizeMsg` → store width/height

### Step 5: Lip Gloss styles (`ui/styles.go`)

- [ ] Define a `Styles` struct or package-level vars for:
  - Border style for the outer frame
  - Header text style
  - Preview block style
  - Editor label / value styles
  - Selected field highlight style (defined now, used in M3)
  - Output row styles
  - Footer style (dimmed)
  - Separator/divider style
- [ ] Keep styles minimal — high contrast, clean

### Step 6: Header renderer (`ui/header.go`)

- [ ] `RenderHeader(clipText, normalizedText, toastMsg string, width int) string`
- [ ] Display:
  - `Clipboard: <original text>`
  - `Normalized: <normalized rgb string>`
  - Toast message (if non-empty) or `Status: OK`
- [ ] Respect terminal width for wrapping/truncation

### Step 7: Preview pane renderer (`ui/preview.go`)

- [ ] `RenderPreview(c color.Color, width, height int) string`
- [ ] Render a solid block of color using Lip Gloss background styling
  - Use `lipgloss.NewStyle().Background(lipgloss.Color(hex))` to fill cells
- [ ] Fill block with spaces styled with the color's background
- [ ] If alpha < 100%: render a simple checkerboard pattern mixing the color with a neutral background to hint at transparency
- [ ] Show foreground suggestion text at bottom of preview:
  - Compute relative luminance: `L = 0.2126*R + 0.7152*G + 0.0722*B` (normalized 0-1)
  - If L > 0.5 → suggest "black", else "white"

### Step 8: Editor pane renderer (`ui/editor.go`)

- [ ] `RenderEditor(c color.Color, mode EditMode, selectedField int, width int) string`
- [ ] Display mode header: `Edit Mode: HSV` or `Edit Mode: RGB`
- [ ] In HSV mode, show:
  - `Hue        <value>°`
  - `Saturation <value>%`
  - `Value      <value>%`
  - `Opacity    <value>%`
  - (Values computed via `color.RGBToHSV()`)
- [ ] In RGB mode, show:
  - `Red        <value>`
  - `Green      <value>`
  - `Blue       <value>`
  - `Opacity    <value>%`
- [ ] Highlight the selected field row (use selected style)
- [ ] Show `Tab: switch to <other mode>` hint at bottom

### Step 9: Output pane renderer (`ui/outputs.go`)

- [ ] `RenderOutputs(c color.Color, width int) string`
- [ ] Three rows:
  - `RGB  rgb(R G B)          [yr] copy`
  - `HEX  #RRGGBB             [yx]`
  - `HSL  hsl(H S% L%)        [yh]`
- [ ] Use `color.FormatRGB()`, `color.FormatHEX()`, `color.FormatHSL()`

### Step 10: Footer renderer (`ui/footer.go`)

- [ ] `RenderFooter(width int) string`
- [ ] Single line of compact key hints, dimmed:
  - `hjkl move • h/l adjust • H/L larger step • tab mode • yr/yx/yh copy • q quit`

### Step 11: Compose view (`app/view.go`)

- [ ] `View() string` — assemble all 5 regions:
  1. Header (full width)
  2. Horizontal split: Preview (left ~40%) | Editor (right ~60%)
  3. Outputs (full width)
  4. Footer (full width)
- [ ] Use `lipgloss.JoinHorizontal` for preview/editor split
- [ ] Use `lipgloss.JoinVertical` for stacking regions
- [ ] Add borders between regions (using Lip Gloss border styles)
- [ ] Respect `model.width` / `model.height` for sizing

### Step 12: Main entry point (`main.go`)

- [ ] Read clipboard via `clipboard.Read()`
- [ ] If clipboard read fails:
  - Print `colr: failed to read clipboard` to stderr
  - `os.Exit(1)`
- [ ] Parse with `color.Parse(clipText)`
- [ ] If parse fails:
  - Print `colr: clipboard does not contain a valid CSS color` to stderr
  - Print `Supported formats: RGB, RGBA, HEX, HSL` to stderr
  - `os.Exit(1)`
- [ ] Construct `app.NewModel(clipText, parsedColor)`
- [ ] Run `tea.NewProgram(model, tea.WithAltScreen()).Run()`
- [ ] Handle program error → print and exit

### Step 13: Build and manual test

- [ ] `go build -o colr .`
- [ ] Test with valid clipboard: copy `#FF0000` → run `./colr` → TUI renders
- [ ] Test with invalid clipboard: copy `hello` → run `./colr` → error message, exits
- [ ] Test with semi-transparent: copy `rgba(255,0,0,.5)` → preview shows, alpha indicated
- [ ] Verify `q` quits cleanly
- [ ] Verify layout adapts to terminal resize

## Scope Boundaries

**In scope:** rendering all 5 regions, startup/error flow, `q` to quit, window resize.

**Deferred to milestone 3:** all keyboard editing (hjkl, tab, value adjustment, reset).

**Deferred to milestone 4:** copy commands (yr/yx/yh), help overlay content, toast timers.

## Acceptance Criteria

- [ ] App launches with valid clipboard color, displays full layout
- [ ] Color preview renders with correct background color
- [ ] Header shows original clipboard text and normalized form
- [ ] Editor pane shows HSV fields with correct values (read-only)
- [ ] Output pane shows correct RGB, HEX, HSL strings
- [ ] Footer shows key hints
- [ ] Invalid clipboard input exits with error message (no TUI)
- [ ] Failed clipboard read exits with error message
- [ ] `q` and `ctrl+c` quit cleanly
- [ ] Terminal resize updates layout dimensions
