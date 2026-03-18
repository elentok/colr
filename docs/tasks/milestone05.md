# Milestone 5 — Polish

## Goal

Make the TUI look and feel polished. Visual consistency, better preview rendering, slider-style field display, and the optional `yy`/`Enter` copy for the selected output row.

## Areas of Work

### 1. Preview pane improvements

**Current state:** alternating rows of solid color and checkerboard blocks (2-char wide blocks in two grays). Foreground suggestion is plain text at the bottom.

**Changes:**
- Improve checkerboard: composite the actual color on top of light/dark cells using Lip Gloss background, so the preview shows the color at its actual opacity visually blended against a check pattern, not just alternating solid/checker rows
- Style the foreground suggestion text with the actual suggested color (white on dark preview, black on light preview)
- Make the "Foreground: white" and "(transparent)" lines styled, not raw text

### 2. Editor pane — slider visualization

**Current state:** fields show label and value, e.g. `Hue  0°`. The spec mockup shows sliders: `Hue 0° [────●────]`.

**Changes:**
- Add a text-based slider bar after each field value: `[────●────]`
- Slider width proportional to available editor pane width
- Knob position (`●`) reflects the current value relative to the field's range
- Selected field's slider should use the highlight color; unselected sliders dimmed

### 3. Header — app name in title border

**Current state:** header is 3 plain text lines. The spec mockup shows `┌ colr ──...` in the outer border.

**Changes:**
- The outer rounded border's top line should include "colr" as a title — use Lip Gloss border title if available, otherwise embed it manually in the header
- Keep the 3-line content (Clipboard / Normalized / Status) as-is

### 4. Output pane — selectable rows with yy/Enter

**Current state:** output rows are display-only. Copy keys `yr`/`yx`/`yh` work globally.

**Changes:**
- Add `selectedOutput int` (0=RGB, 1=HEX, 2=HSL) to the model — but only active when a mode for it is desired. Keep it simple: always track a selected output row.
- Highlight the selected output row (subtle — dimmer than the editor field highlight)
- `yy` copies the selected output row
- `Enter` copies the selected output row
- Navigation: tab between editor and output pane? **No** — keep it simple. The output rows cycle via some key. Let's use the existing `j`/`k` behavior: when `selectedField` goes past the last editor field (Opacity), cursor moves into the output rows. When it goes above the first output row, it returns to the editor.

  Actually this adds complexity to field indexing. **Simpler approach:** `yy` and `Enter` copy the selected output row. Use `J`/`K` (shifted) or just skip output row navigation entirely and always default the selected output to RGB. Simplest: **don't add output navigation**. Just add `yy` → copies RGB (the first/default output). That matches the spec's "nice to have" without over-engineering.

  **Final decision:** Track `selectedOutput` (0–2). No special navigation for it — it defaults to 0 (RGB). `yy` copies the row at `selectedOutput`. `Enter` copies the row at `selectedOutput`. Not adding j/k navigation into output rows. This is the minimal useful version of the feature.

### 5. Footer — style refinement

**Current state:** plain dimmed text with bullet separators.

**Changes:**
- Highlight key names differently from descriptions (keys in accent color, descriptions in dim)
- Ensure footer content fits or degrades gracefully on narrow terminals

### 6. General style cleanup

- Consistent color palette across all regions (accent color for interactive elements, dim for chrome)
- Ensure the app looks good on both dark and light terminal themes (test with dark; spec says high contrast)
- Remove any unused styles from `ui/styles.go`

## Implementation Steps

### Step 1: Improve checkerboard preview (`ui/preview.go`)

- [ ] Replace alternating-row approach with per-cell compositing: for each cell, alternate between two checker backgrounds, then blend the color on top at its alpha
- [ ] Since true alpha compositing isn't possible in terminal cells, approximate by: when alpha < 1.0, show every other block in the color and every other block in a checker background — gives a visual hint of transparency
- [ ] Style the "Foreground: white/black" hint using the actual fg color against the bg color
- [ ] Add "(transparent)" indicator only when alpha < 100%

### Step 2: Add sliders to editor fields (`ui/editor.go`)

- [ ] `renderSlider(value, min, max float64, width int, selected bool) string`
  - Generates `[────●────]` where `●` position is proportional to `(value - min) / (max - min)`
  - Track width = `width - 2` (for `[` and `]`)
  - Filled portion uses `─`, knob uses `●`
  - Selected: accent color for knob; unselected: dim
- [ ] Append slider after value in each field row
- [ ] Slider width: allocate remaining space after label + value + spacing
- [ ] Field ranges:
  - Hue: 0–360
  - Saturation: 0–100
  - Value: 0–100
  - Opacity: 0–100
  - Red/Green/Blue: 0–255

### Step 3: Header title in outer border (`app/view.go`)

- [ ] Check if Lip Gloss v1.x supports border title — if yes, use `.BorderTop(true).SetString(" colr ")` or similar
- [ ] If not supported, manually render: replace the first line of the outer border string to inject " colr " after the `╭` character
- [ ] Keep it minimal: `╭ colr ─────...╮`

### Step 4: Selected output row + yy/Enter (`app/model.go`, `app/update.go`, `ui/outputs.go`)

- [ ] Add `selectedOutput int` to Model (default 0 = RGB)
- [ ] `yy` → `applyCopyByIndex(m, m.selectedOutput)`
- [ ] `Enter` → same
- [ ] `applyCopyByIndex` maps 0→"rgb", 1→"hex", 2→"hsl"
- [ ] `RenderOutputs` gains `selectedOutput int` param; highlight that row with a subtle style

### Step 5: Footer key style (`ui/footer.go`)

- [ ] Render key names (e.g., "hjkl", "yr/yx/yh") in accent color
- [ ] Render descriptions (e.g., "move", "copy") in dim
- [ ] Adjust separator bullet style

### Step 6: General style audit (`ui/styles.go`)

- [ ] Review all styles for consistency
- [ ] Remove any unused styles
- [ ] Ensure accent color (currently "39" = bright cyan) is used consistently for interactive elements

### Step 7: Verification

- [ ] `go build ./...` clean
- [ ] `go test ./...` all pass
- [ ] Manual test: launch with various colors, verify visual coherence
- [ ] Manual test: sliders update with value changes
- [ ] Manual test: `yy` copies the correct output format
- [ ] Manual test: Enter copies the correct output format
- [ ] Manual test: help overlay still renders correctly
- [ ] Manual test: preview checkerboard looks reasonable for semi-transparent colors

## Acceptance Criteria

- [ ] App looks coherent and polished
- [ ] Color preview has improved transparency hint
- [ ] Editor fields have slider visualization
- [ ] Outer border shows "colr" title
- [ ] `yy` and `Enter` copy the selected output row
- [ ] Footer keys are visually distinct from descriptions
- [ ] Important state (selected field, edit mode, toast) is visually obvious
- [ ] All existing tests still pass
- [ ] No visual regressions (help overlay, outputs alignment, height fill)
