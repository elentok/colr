# Changelog

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
