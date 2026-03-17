# colr — Claude Code Implementation Spec

## Project Summary

Build a terminal UI application named `colr` in Go.

`colr` reads a CSS color value from the clipboard, parses it, shows a live preview,
allows interactive editing, and lets the user copy the color in normalized CSS formats.

This is a keyboard-driven TUI built with:

- Go
- Bubble Tea
- Lip Gloss
- Bubbles (optional / only where useful)

The application should feel fast, minimal, and Vim-friendly.

---

# Product Definition

## Main Use Cases

### 1. Convert a color from one format to another

User copies a color string such as:

```
rgba(255, 0, 0, .5)
```

Then runs:

```
colr
```

The app parses the clipboard value and displays normalized output formats:

```
RGB rgb(255 0 0 / 50%)
HEX #FF000080
HSL hsl(0 100% 50% / 50%)
```

The user can copy any format to the clipboard.

---

### 2. Tweak a color interactively

User copies a color string and runs `colr`.

The app opens in a split layout with:

- a large live preview
- editable channels
- normalized outputs

The user adjusts the color via keyboard, sees the preview update live, and copies the desired format.

---

# Scope

## In Scope

- Read clipboard on startup
- Parse supported CSS color formats
- Normalize into internal color representation
- Display a split preview/editor TUI
- Support editing in:
  - HSV mode
  - RGB mode
- Support opacity editing
- Show normalized RGB, HEX, HSL output strings
- Copy RGB / HEX / HSL to clipboard
- Vim-style keybindings
- Fast failure on invalid clipboard input

## Out of Scope for v1

- Mouse support
- Arbitrary text editing inside the app
- Named CSS colors
- OKLCH / LAB / LCH / HWB
- Saving palettes
- Reading input from command-line args
- Reading from stdin
- Multi-color workflows
- Palette picker UI

---

# Supported Input Formats

## RGB Inputs

```
rgb(100 200 300)
rgb(100 200 300 / 50%)
rgb(100 200 300 / 0.5)
rgb(100, 200, 300)
rgb(100% 20% 30%)
100, 200, 300
100 200 300
100% 20% 30%
```

## RGBA Inputs

```
rgba(100, 200, 300, 50%)
rgba(100, 200, 300, 0.5)
rgba(100, 200, 300, .5)
```

## HEX Inputs

Accept both with and without `#`:

```
a3a3a3
A3A3A3
a3a3a3aa
A3A3A3AA
#a3a3a3
#A3A3A3
#a3a3a3aa
#A3A3A3AA
```

## HSL Inputs

```
hsl(32 100% 50%)
hsl(32 100% 50% / 50%)
```

---

# Output Formats

The app must always display normalized output in these three formats:

## RGB Output

Without alpha:

```
rgb(255 0 0)
```

With alpha:

```
rgb(255 0 0 / 50%)
```

Rules:

- use `rgb(...)`, never `rgba(...)`
- use space-separated syntax
- alpha is displayed as percentage
- RGB channels are integers

## HEX Output

Without alpha:

```
#FF0000
```

With alpha:

```
#FF000080
```

Rules:

- uppercase
- include `#`
- include alpha only when opacity is not 100%

## HSL Output

Without alpha:

```
hsl(0 100% 50%)
```

With alpha:

```
hsl(0 100% 50% / 50%)
```

Rules:

- use `hsl(...)`, never `hsla(...)`
- use space-separated syntax
- alpha is displayed as percentage
- H is integer degrees
- S and L are integer percentages

---

# Startup Behavior

## Normal startup

1. Read clipboard text
2. Trim surrounding whitespace
3. Parse color
4. If valid:
   - store original parsed color
   - compute derived formats
   - launch TUI

## Invalid startup

If clipboard does not contain a supported color value:

- print a short error message to stderr
- exit with non-zero status
- do not open the TUI

Example:

```
colr: clipboard does not contain a valid CSS color
Supported formats: RGB, RGBA, HEX, HSL
```

---

# Internal Data Model

Use RGBA as the source of truth.

Recommended structure:

```
type Color struct {
R uint8
G uint8
B uint8
A float64 // 0.0 - 1.0
}
```

Derived values should be computed as needed:

- HSV
- HSL
- formatted RGB string
- formatted HEX string
- formatted HSL string

Recommended helper structures:

```
type HSV struct {
H float64 // 0-360
S float64 // 0-1
V float64 // 0-1
}

type HSL struct {
H float64 // 0-360
S float64 // 0-1
L float64 // 0-1
}
```

Do not treat HSV/HSL as separate sources of truth.

When editing in HSV mode:

- convert RGBA -> HSV
- apply edits in HSV
- convert HSV -> RGB
- preserve alpha separately

When editing in RGB mode:

- mutate RGB directly
- preserve alpha separately

---

# UI Layout

Use a split editor/preview layout.

```
┌ colr ────────────────────────────────────────────────────────────────┐
│ Clipboard: rgba(255, 0, 0, .5) │
│ Normalized: rgb(255 0 0 / 50%) │
│ Status: OK │
├───────────────────────────────┬─────────────────────────────────────┤
│ │ Edit Mode: HSV │
│ COLOR PREVIEW │ │
│ │ Hue 0° [────●────] │
│ ███████████████████ │ Saturation 100% [────────●] │
│ ███████████████████ │ Value 100% [────────●] │
│ ███████████████████ │ Opacity 50% [────●────] │
│ │ │
│ Foreground suggestion: white │ Tab: switch to RGB │
│ Alpha preview: checkerboard │ │
├───────────────────────────────┴─────────────────────────────────────┤
│ RGB rgb(255 0 0 / 50%) [y r] copy │
│ HEX #FF000080 [y x] │
│ HSL hsl(0 100% 50% / 50%) [y h] │
├──────────────────────────────────────────────────────────────────────┤
│ hjkl/←↓↑→ move • +/- adjust • H/L larger step • q quit • ? help │
└──────────────────────────────────────────────────────────────────────┘
```

---

# UI Regions

## 1. Header

Must display:

- app name
- original clipboard text
- normalized color string
- status / toast area

Example:

```
Clipboard: rgba(255, 0, 0, .5)
Normalized: rgb(255 0 0 / 50%)
```

---

## 2. Preview Pane

Must display:

- large solid color preview
- checkerboard indication when alpha < 100%
- suggested foreground text color:
  - `black`
  - `white`

Foreground suggestion can be based on luminance threshold.

---

## 3. Editor Pane

Two editing modes:

### HSV Mode (default)

Fields:

```
Hue
Saturation
Value
Opacity
```

Ranges:

```
Hue: 0–360
Saturation: 0–100
Value: 0–100
Opacity: 0–100
```

### RGB Mode

Fields:

```
Red
Green
Blue
Opacity
```

Ranges:

```
Red: 0–255
Green: 0–255
Blue: 0–255
Opacity: 0–100
```

One field is selected at a time.

Selected field must be visually highlighted.

Adjustments must update preview and output strings immediately.

---

## 4. Output Pane

Must show these rows:

- RGB
- HEX
- HSL

Example:

```
RGB rgb(255 0 0 / 50%)
HEX #FF000080
HSL hsl(0 100% 50% / 50%)
```

These rows are display-first; copy commands work globally and do not require focus.

Optional v1 behavior:

- allow selecting output row
- `yy` copies selected output
- `Enter` copies selected output

This is nice to have, not required for MVP.

---

## 5. Footer

Must show compact key hints.

Example:

```
hjkl move • h/l adjust • H/L larger step • tab mode • yr/yx/yh copy • q quit
```

---

# Keybindings

## Global

```
q quit
? toggle help overlay
R reset to original clipboard color
tab toggle edit mode
1 switch to HSV mode
2 switch to RGB mode
```

## Movement

```
j / ↓ move to next editable field
k / ↑ move to previous editable field
g jump to first field
G jump to last field
```

## Adjust Current Field

```
h / ← decrease by small step
l / → increase by small step
H decrease by large step
L increase by large step

-       decrease by small step

*       increase by small step
```

## Copy Commands

```
yr copy normalized RGB
yx copy normalized HEX
yh copy normalized HSL
```

Optional:

```
yy copy selected output row
Enter copy selected output row
```

---

# Value Adjustment Rules

## HSV Mode

Small / large steps:

```
Hue: 1 / 10
Saturation: 1 / 5
Value: 1 / 5
Opacity: 1 / 5
```

## RGB Mode

Small / large steps:

```
Red: 1 / 10
Green: 1 / 10
Blue: 1 / 10
Opacity: 1 / 5
```

## Clamping

All edits must clamp safely:

- RGB channels: 0–255
- Hue: wrap or clamp to 0–360
- Saturation: 0–100
- Value: 0–100
- Opacity: 0–100

Recommended behavior for hue:

- wrap around rather than hard clamp
- e.g. decrementing below 0 goes to 359
- incrementing above 360 wraps to 0 or 1 depending on representation

Use a consistent rule throughout the app.

---

# Parsing Rules

## General

- trim leading/trailing whitespace
- accept uppercase or lowercase function names if practical
- ignore case for hex digits
- accept both `#AABBCC` and `AABBCC`
- accept alpha as:
  - percentage
  - decimal in 0–1 range
  - decimal without leading zero where supported by listed formats

## Numeric interpretation

### RGB integer values

```
rgb(255 0 0)
```

Interpret channels as 0–255.

### RGB percentage values

```
rgb(100% 20% 30%)
```

Convert percentages to 0–255.

### Alpha percentage

```
50%
```

Convert to 0.5 internally.

### Alpha decimal

```
0.5
.5
```

Convert directly to internal alpha.

## Invalid inputs

Reject clearly invalid values.

Examples to reject:

```
hello
rgb(foo bar baz)
rgb(10 20)
#XYZ123
hsl(20 30)
```

---

# Formatting Rules

## RGB formatting

If alpha == 100%:

```
rgb(R G B)
```

Otherwise:

```
rgb(R G B / A%)
```

Examples:

```
rgb(255 0 0)
rgb(255 0 0 / 50%)
```

## HEX formatting

If alpha == 100%:

```
#RRGGBB
```

Otherwise:

```
#RRGGBBAA
```

Examples:

```
#FF0000
#FF000080
```

## HSL formatting

If alpha == 100%:

```
hsl(H S% L%)
```

Otherwise:

```
hsl(H S% L% / A%)
```

Examples:

```
hsl(0 100% 50%)
hsl(0 100% 50% / 50%)
```

Rounding rules should be consistent and deterministic.

Recommended:

- RGB channels: nearest integer
- HSL H: nearest integer
- HSL S/L: nearest integer percentage
- alpha output: nearest integer percentage

---

# Clipboard Behavior

## On startup

Read text from system clipboard.

## On copy command

Write the selected normalized output string to the system clipboard.

## Error handling

If clipboard read fails on startup:

- print a clear error
- exit non-zero

If clipboard write fails during TUI use:

- keep app running
- show temporary error toast

Example:

```
Failed to copy HEX to clipboard
```

---

# Help Overlay

Pressing `?` toggles a help overlay.

Suggested content:

```
colr help

Movement
j/k or ↑/↓ move between fields
h/l or ←/→ adjust field
H/L adjust by larger step
g/G first/last field

Modes
tab toggle HSV/RGB
1 HSV mode
2 RGB mode

Copy
yr copy RGB
yx copy HEX
yh copy HSL

Other
R reset color
q quit
```

Overlay can be modal and dismiss on:

- `?`
- `Esc`
- `q` only if desired, but do not accidentally quit the app from help unless intentional

Recommended:

- `?` and `Esc` close help
- `q` still quits app globally only if not in help overlay

---

# Toast / Status Messages

Use short-lived messages for successful copy actions and non-fatal errors.

Examples:

```
Copied RGB to clipboard
Copied HEX to clipboard
Copied HSL to clipboard
Failed to copy HEX to clipboard
```

Recommended duration:

- 1 second for success
- 1.5–2 seconds for errors

---

# Visual Design Requirements

## General

- minimal
- clean
- keyboard-first
- high contrast
- good use of spacing
- no excessive decoration

## Lip Gloss usage

Use Lip Gloss for:

- panel borders
- selected field highlighting
- preview area layout
- footer/status styling

## Preview rendering

Use terminal background coloring to render the color preview.

If alpha < 100%, try to communicate transparency visually with a checkerboard-like pattern behind or around the preview.

Exact rendering can be approximate; do not block implementation on perfect alpha compositing in terminal cells.

---

# Suggested Package Structure

```
colr/
main.go

app/
model.go
update.go
view.go
keymap.go
types.go

color/
parse.go
format.go
convert.go
clamp.go

clipboard/
clipboard.go

ui/
header.go
preview.go
editor.go
outputs.go
help.go
styles.go
```

This structure is guidance, not a strict requirement.

---

# Implementation Milestones

## Milestone 1 — Parsing and formatting library

Implement:

- clipboard text input wrapper
- parser for supported input formats
- internal RGBA color model
- RGB / HEX / HSL formatting
- RGB <-> HSV conversion
- RGB <-> HSL conversion

### Acceptance criteria

- can parse all listed valid examples
- rejects clearly invalid examples
- outputs normalized RGB / HEX / HSL strings correctly
- alpha is handled correctly

### Suggested tests

- unit tests for parsing
- unit tests for formatting
- unit tests for conversions
- edge-case tests for alpha and percentages

---

## Milestone 2 — Minimal TUI shell

Implement:

- Bubble Tea app boot
- startup clipboard read
- invalid clipboard error path
- basic layout with header / preview / editor / outputs / footer
- render normalized values

### Acceptance criteria

- app launches with valid clipboard color
- preview displays
- outputs display
- invalid clipboard exits cleanly without entering TUI

---

## Milestone 3 — Interactive editing

Implement:

- HSV mode
- RGB mode
- field selection
- value adjustment
- live preview updates
- reset behavior

### Acceptance criteria

- adjusting fields updates preview immediately
- switching modes preserves current color
- opacity editing works
- reset returns exactly to original parsed color

---

## Milestone 4 — Clipboard copy actions and help

Implement:

- `yr`, `yx`, `yh`
- help overlay
- toast messages
- non-fatal clipboard copy errors

### Acceptance criteria

- copy commands place normalized output on clipboard
- success and failure messages appear
- help overlay is readable and dismissible

---

## Milestone 5 — Polish

Implement:

- better styles
- selected field highlight
- contrast hint
- checkerboard transparency hint
- optional selected output row / `yy`

### Acceptance criteria

- app feels coherent and polished
- important state is visually obvious
- controls are discoverable via footer/help

---

# Test Plan

## Unit tests

Must cover:

- parse RGB variants
- parse RGBA variants
- parse HEX variants
- parse HSL variants
- invalid inputs
- formatting behavior
- alpha handling
- conversion correctness for representative colors

Representative colors:

```
#000000
#FFFFFF
#FF0000
#00FF00
#0000FF
#808080
#FF000080
```

## Integration-ish tests

At minimum, logic tests should verify:

- startup parse flow
- editing mutations
- mode switching behavior
- reset behavior
- copy string generation

If full Bubble Tea interaction tests are too heavy, prioritize model/update tests over snapshot-heavy UI tests.

---

# Edge Cases

Must handle these safely:

- lowercase / uppercase hex input
- hex with alpha
- alpha 0%
- alpha 100%
- percentage RGB values
- decimal alpha values like `.5`
- whitespace around clipboard content
- hue wraparound
- clamping on repeated edits
- reset after multiple edits
- fully transparent colors
- grayscale colors where hue may be unstable

For grayscale colors in HSV/HSL conversions:

- use a stable, deterministic hue strategy
- do not let the UI jitter unpredictably

Recommended:

- preserve previous hue where reasonable during edits
- otherwise default to 0

---

# Non-Functional Requirements

- startup should feel instant
- keypress response should feel immediate
- code should be organized and testable
- parser and formatter logic should be independent from TUI code
- no unnecessary abstraction
- prioritize clarity over cleverness

---

# Delivery Requirements for Claude Code

Please implement this project incrementally.

Recommended execution order:

1. parsing + formatting package
2. conversion helpers
3. minimal Bubble Tea app shell
4. editor interactions
5. clipboard copy commands
6. styling and polish
7. tests

When implementing:

- keep pure color logic separate from UI logic
- write tests for parser/formatter/conversion code
- avoid overengineering
- prefer straightforward readable code
- keep key handling explicit and easy to follow

---

# Definition of Done

The project is done when:

- `colr` reads a supported color from clipboard
- launches the TUI successfully
- shows a live preview
- allows editing in HSV and RGB modes
- allows opacity editing
- shows normalized RGB / HEX / HSL output
- supports copying RGB / HEX / HSL to clipboard
- handles invalid clipboard input gracefully
- has a help overlay
- has basic tests for parser/formatter/conversion logic
- feels usable from keyboard only

---

# Nice-to-Have After v1

Possible future additions:

- named CSS colors
- OKLCH display
- direct text input mode
- CLI arg input
- stdin input
- palette generation
- theme preview
- gradient preview
- export multiple formats at once
