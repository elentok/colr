# Initial prompt (for planning with ChatGPT)

I want to create a TUI named "colr" for handling color values, specifically for CSS color formats.

Definition: **colr** is a terminal color inspector and tuner for CSS color values.

## Support color formats

These are the supported input formats, for each one I marked the one to use for displaying:

- RGB:
  - `rgb(100 200 300)` (for display - without opacity)
  - `rgb(100 200 300 / 50%)` (for display - with opacity)
  - `rgb(100 200 300 / 0.5)`
  - `rgb(100, 200, 300)`
  - `rgb(100% 20% 30%)`
  - `100, 200, 300`
  - `100 200 300`
  - `100% 20% 30%`
- RGBA: (use the RGB format above for displaying)
  - `rgba(100, 200, 300, 50%)`
  - `rgba(100, 200, 300, 0.5)`
  - `rgba(100, 200, 300, .5)`
- HEX: (with or without '#' prefix)
  - `a3a3a3`
  - `A3A3A3`
  - `a3a3a3aa`
  - `A3A3A3AA`
- HSL:
  - `hsl(32 100% 50%)`
  - `hsl(32 100% 50% / 50%)`

## Use cases

There are two primary use cases:

1. Convert a given format to another format
2. Tweak a given color and copy the new value

When running "colr" read the clipboard text:

- Try to parse it as a color value
  - If it's an invalid value, exit with an error message
- Render a block in this color (so I can preview what it looks like)
- Allow modifying the color by changing either RGB values or HSV values
- Allow modifying the opacity (display it as percentage)
- Show them in RGB, HEX and HSL and allow me copy those formats to clipboard

## Guidelines

- The TUI should use vim bindings:
  - HJKL for movement in addition to regular arrows
  - single key or chained keys to perform actions (e.g. "d" to delete a worktree)
- Write using Go using the BubbleTea framework with lipgloss for formatting and
  bubbles for components

## UI Design

```
┌ colr ───────────────────────────────────────────────────────────────┐
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
