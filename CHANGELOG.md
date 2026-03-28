# Changelog

## [0.1.1] - 2026-03-28

- Replaced checkerboard alpha preview with a background-aware composited swatch
- Preview background toggle with `b` (white/black)
- `yo` keybinding to copy the composited over-background HEX value
- Output panel now shows both raw values and an `OVER` row for the current background context

## [0.1.0] - 2026-03-20

### Added

- Explicit history save with `s`, so the current edited color can be written to history without quitting
- Startup fallback to the most recent history color when clipboard or CLI input contains no color

### Changed

- History entries now store normalized RGB, HEX, and CSS name values instead of clipboard-origin text
- HEX output is now normalized to lowercase consistently across the UI, history, and copied HEX values
- The preview compares original versus edited colors, with responsive behavior:
  - landscape terminals show a right sidebar with stacked original/edited previews
  - portrait terminals move the preview to the bottom and show original/edited side-by-side
- Main layout updated to use separate bordered panels for the header, editor, formats, and preview

## [0.0.2] - 2026-03-19

### Added

- CSS named color display: a NAME row in the outputs panel shows the nearest CSS named color (`red`, `~salmon`, etc.)
- `yn` keybinding to copy the CSS color name to the clipboard
- History selector (`p`) to browse and reload recently used colors (stored in `XDG_STATE_HOME/colr/history.json`, last 100 entries)
- Command-line color argument: pass a color directly as `colr "#ff8000"` or `colr 255 128 0`
- Multi-color input: `colr "button=rgb(255 128 0); border=#112233"` picks the first valid color found

## [0.0.1] - initial release

- Interactive TUI color editor with HSV and RGB edit modes
- Live preview swatch
- Copy to clipboard in RGB, HEX, and HSL formats (`yr`/`yy`, `yx`, `yh`)
- Clipboard input on launch
