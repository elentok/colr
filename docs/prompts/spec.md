# colr — TUI CSS Color Inspector and Editor

## Overview

`colr` is a terminal-based color inspector and editor for CSS color values.

It reads a color value from the clipboard, parses it, displays a visual preview,
and allows the user to tweak the color interactively. The tool also converts the
color into several CSS formats and allows copying them to the clipboard.

The interface is implemented as a TUI using:

- Go
- Bubble Tea
- Lip Gloss
- Bubbles (optional for reusable components)

The UI uses Vim-style keybindings and is designed for fast keyboard-driven workflows.

---

# Goals

Primary goals:

1. Convert color formats quickly
2. Visually inspect colors
3. Interactively tweak colors
4. Copy CSS-ready formats

Secondary goals:

- Provide intuitive HSV editing
- Maintain CSS-compliant output
- Work entirely with keyboard input

---

# Supported Color Formats

## Accepted Input Formats

### RGB

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

### RGBA

```
rgba(100, 200, 300, 50%)
rgba(100, 200, 300, 0.5)
rgba(100, 200, 300, .5)
```

### HEX

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

### HSL

```
hsl(32 100% 50%)
hsl(32 100% 50% / 50%)
```

---

## Display Output Formats

The application displays the following normalized formats:

### RGB

Without alpha:

```
rgb(255 0 0)
```

With alpha:

```
rgb(255 0 0 / 50%)
```

### HEX

```
#FF0000
#FF000080
```

### HSL

```
hsl(0 100% 50%)
hsl(0 100% 50% / 50%)
```

Output rules:

- RGB uses space-separated syntax
- Alpha is displayed as percentage
- HEX output is uppercase
- `rgba()` and `hsla()` are never emitted

---

# Application Startup Flow

1. Read clipboard contents
2. Attempt to parse a color
3. If parsing fails:
   - Print error
   - Exit with code 1
4. If parsing succeeds:
   - Normalize color internally
   - Launch the TUI

Error example:

```
colr: clipboard does not contain a valid CSS color
Supported formats: RGB, RGBA, HEX, HSL
```

---

# Internal Color Model

Internally the application stores:

```
type Color struct {
    R uint8
    G uint8
    B uint8
    A float64

    H float64
    S float64
    V float64

    HslH float64
    HslS float64
    HslL float64
}
```

Source of truth is **RGBA**.

HSV and HSL values are derived whenever RGB changes.

---

# UI Layout

The interface is divided into three main zones:

```
┌ colr ────────────────────────────────────────────────────────────────┐
│ Clipboard: rgb(255 0 0 / 50%)                                       │
│ Parsed as: RGB with alpha                                           │
│ Status: OK                                                          │
├───────────────────────────────┬─────────────────────────────────────┤
│                               │ Edit Mode: HSV                      │
│         COLOR PREVIEW         │                                     │
│                               │  Hue         0°        [────●────]  │
│      ███████████████████      │  Saturation  100%      [────────●]  │
│      ███████████████████      │  Value       100%      [────────●]  │
│      ███████████████████      │  Opacity     50%       [────●────]  │
│                               │                                     │
│  Foreground suggestion: white │  Alt mode: RGB   (press tab)        │
│  Alpha preview: checkerboard  │                                     │
├───────────────────────────────┴─────────────────────────────────────┤
│ RGB  rgb(255 0 0 / 50%)                                  [y r] copy │
│ HEX  #FF000080                                              [y x]    │
│ HSL  hsl(0 100% 50% / 50%)                                 [y h]    │
├──────────────────────────────────────────────────────────────────────┤
│ hjkl/←↓↑→ move • +/- adjust • shift for larger step • tab mode      │
│ y{r,x,h} copy • R reset • q quit • ? help                           │
└──────────────────────────────────────────────────────────────────────┘
```

---

# UI Sections

## Header

Displays:

- Original clipboard value
- Parsed format
- Status

Example:

```
Clipboard: rgba(255, 0, 0, .5)
Normalized: rgb(255 0 0 / 50%)
```

---

## Color Preview Panel

Displays a large block showing the current color.

Behavior:

- If alpha < 100%, show checkerboard background
- Render using terminal background colors
- Display text contrast recommendation

Example:

```
Foreground suggestion: white
```

---

## Editor Panel

Allows modifying color values.

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
Saturation: 0–100%
Value: 0–100%
Opacity: 0–100%
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
Opacity: 0–100%
```

Mode switching:

```
tab
1 → HSV
2 → RGB
```

---

## Output Panel

Displays normalized formats.

```
RGB
HEX
HSL
```

Each row can be copied.

Example:

```
RGB  rgb(255 0 0 / 50%)
HEX  #FF000080
HSL  hsl(0 100% 50% / 50%)
```

---

# Keybindings

## Movement

```
j / ↓   move down
k / ↑   move up
h / ←   decrease value
l / →   increase value
g       first field
G       last field
```

---

## Adjustment

```
h/l adjust by small step
H/L adjust by large step
+   increase
-   decrease
```

Step sizes:

HSV mode:

```
Hue: 1 / 10
Saturation: 1% / 5%
Value: 1% / 5%
Opacity: 1% / 5%
```

RGB mode:

```
RGB: 1 / 10
Opacity: 1% / 5%
```

---

## Copy Commands

Copy formats directly:

```
yr → copy RGB
yx → copy HEX
yh → copy HSL
```

Alternative:

```
yy → copy selected output row
Enter → copy selected output row
```

---

## Mode Switching

```
tab → toggle HSV/RGB
1 → HSV mode
2 → RGB mode
```

---

## Utility

```
R → reset to original color
q → quit
? → help
```

---

# Toast Notifications

Temporary messages appear at the bottom.

Examples:

```
Copied HEX to clipboard
Copied RGB to clipboard
```

Duration:

```
~1 second
```

---

# Help Overlay

Opened with:

```
?
```

Example:

```
colr help

Movement
  j/k or ↑/↓    move between fields
  h/l or ←/→    adjust field
  H/L           adjust by larger step
  tab           switch edit mode

Modes
  1             HSV mode
  2             RGB mode

Copy
  yr            yank RGB
  yx            yank HEX
  yh            yank HSL
  yy            yank selected output

Other
  R             reset color
  q             quit
```

---

# Bubble Tea Architecture

## Model

The main model stores:

```
type Model struct {
    width  int
    height int

    original Color
    color    Color

    editMode EditMode
    focus    FocusArea

    selectedField int
    selectedOutput int

    statusMessage string
}
```

---

## Update Loop

Handles:

- key events
- color updates
- mode switching
- copy commands
- reset

---

## View Rendering

The view renders:

```
Header
Preview panel
Editor panel
Output panel
Footer help
Toast
```

Lip Gloss is used for layout and styling.

---

# Suggested Project Structure

```
colr/
  main.go

  app/
    model.go
    update.go
    view.go
    keymap.go

  color/
    parse.go
    convert.go
    format.go

  ui/
    preview.go
    editor.go
    outputs.go
    header.go
    help.go

  clipboard/
    clipboard.go
```

---

# Future Enhancements

Possible future improvements:

- OKLCH display
- Color palette generation
- Named CSS color support
- Terminal theme preview
- Gradient preview
- Direct input editing
- Mouse support
