# Milestone 3 — Interactive Editing

## Goal

Make the editor pane interactive: navigate fields, adjust values, switch between HSV/RGB modes, and reset to the original color. All edits update the preview and output strings immediately.

## Existing Code Touched

- `app/update.go` — currently only handles `q`/`ctrl+c` and window resize; will gain all keybinding logic
- `app/model.go` — may need a `lastHue` field for grayscale hue stability
- `ui/editor.go` — already renders fields; no changes expected (rendering is read-only)

## Design Decisions

### Source of truth

`currentColor` (RGBA) remains the single source of truth.

- **HSV editing:** convert `currentColor` → HSV, apply delta, convert back to RGBA, store in `currentColor`.
- **RGB editing:** mutate `currentColor` fields directly.
- **Opacity:** mutate `currentColor.A` directly in both modes.

### Grayscale hue stability

When saturation is 0, converting RGB→HSV always yields H=0 regardless of the original hue. This would cause the hue slider to "snap to 0" while the user is adjusting saturation near zero.

Solution: store a `lastHue float64` on the model. When converting to HSV for display or editing:
- if the converted S > 0, use the converted H and update `lastHue`
- if S == 0, substitute `lastHue` for H

This keeps the hue slider stable when the color is achromatic.

### Hue wrapping

Hue wraps around: decrementing below 0 → 359, incrementing past 359 → 0. All other channels clamp.

## Implementation Steps

### Step 1: Add `lastHue` to model (`app/model.go`)

- [ ] Add `lastHue float64` field to `Model`
- [ ] In `NewModel`, initialize `lastHue` from the parsed color's HSV hue

### Step 2: Add a clamp/wrap helper (`color/clamp.go`)

- [ ] `ClampInt(v, min, max int) int` — clamp to [min, max]
- [ ] `ClampFloat(v, min, max float64) float64` — clamp to [min, max]
- [ ] `WrapFloat(v, min, max float64) float64` — wrap around (for hue: 0–360)

### Step 3: Add editing logic (`app/edit.go`)

Separate file to keep `update.go` focused on message dispatch.

- [ ] `adjustField(m *Model, delta int, large bool)` — adjust the currently selected field
  - Reads `m.editMode` and `m.selectedField` to determine which field
  - Computes the actual step (small vs large, per the spec table)
  - For HSV fields: convert `currentColor` → HSV, apply delta, convert back
  - For RGB fields: mutate directly
  - For opacity: mutate `m.currentColor.A` directly
  - Clamps/wraps appropriately
  - Updates `lastHue` when hue changes with S > 0

#### Step sizes (from spec)

**HSV mode:**
| Field      | Small | Large |
|------------|-------|-------|
| Hue        | 1     | 10    |
| Saturation | 1     | 5     |
| Value      | 1     | 5     |
| Opacity    | 1     | 5     |

**RGB mode:**
| Field   | Small | Large |
|---------|-------|-------|
| Red     | 1     | 10    |
| Green   | 1     | 10    |
| Blue    | 1     | 10    |
| Opacity | 1     | 5     |

### Step 4: Keybindings in update (`app/update.go`)

- [ ] **Movement:**
  - `j` / `down` → next field (clamp at FieldCount-1)
  - `k` / `up` → previous field (clamp at 0)
  - `g` → first field (0)
  - `G` → last field (FieldCount-1)
- [ ] **Adjust:**
  - `h` / `left` / `-` → decrease by small step
  - `l` / `right` / `+` → increase by small step
  - `H` → decrease by large step
  - `L` → increase by large step
- [ ] **Mode switching:**
  - `tab` → toggle HSV ↔ RGB
  - `1` → switch to HSV
  - `2` → switch to RGB
  - On mode switch: preserve `selectedField` index (both modes have 4 fields), preserve `currentColor` unchanged
- [ ] **Reset:**
  - `R` → reset `currentColor` to `originalColor`, reset `lastHue` from original

### Step 5: Update editor view for hue stability (`ui/editor.go`)

- [ ] `RenderEditor` gains a `lastHue float64` parameter
- [ ] In HSV mode: if converted S == 0, display `lastHue` instead of 0 for the hue field
- [ ] This is a display concern only — the model already holds the correct `lastHue`

### Step 6: Wire lastHue into View (`app/view.go`)

- [ ] Pass `m.lastHue` to `ui.RenderEditor`

### Step 7: Tests for editing logic (`app/edit_test.go`)

- [ ] **HSV hue adjustment:** start at H=0, increment → H=1; at H=359, increment → H=0 (wrap)
- [ ] **HSV saturation/value clamping:** at 100, increment stays at 100; at 0, decrement stays at 0
- [ ] **RGB channel clamping:** at 255, increment stays at 255; at 0, decrement stays at 0
- [ ] **Opacity clamping:** at 100%, increment stays; at 0%, decrement stays
- [ ] **Large steps:** verify correct step sizes per spec table
- [ ] **Mode switch preserves color:** switch HSV→RGB→HSV, color unchanged
- [ ] **Reset:** edit a color, then reset → matches original
- [ ] **Grayscale hue stability:** set S=0, lastHue should be preserved on adjustment

### Step 8: Tests for clamp/wrap helpers (`color/clamp_test.go`)

- [ ] `ClampInt` basic cases
- [ ] `WrapFloat` for hue: 361 → 1, -1 → 359, 360 → 0
- [ ] Edge cases: exact boundary values

### Step 9: Integration verification

- [ ] `go build ./...` clean
- [ ] `go test ./...` all pass
- [ ] Manual test: launch with `#FF8000`, navigate fields, adjust values, verify preview updates
- [ ] Manual test: switch modes, verify color preserved
- [ ] Manual test: reset after edits, verify exact original restored
- [ ] Manual test: grayscale color — adjust saturation to 0 and back, hue shouldn't jump

## Scope Boundaries

**In scope:** all keybindings for movement/adjustment/mode/reset, clamping/wrapping, hue stability, editing tests.

**Deferred to milestone 4:** copy commands (`yr`/`yx`/`yh`), help overlay, toast messages.

## Acceptance Criteria

- [ ] Adjusting fields updates preview and output strings immediately
- [ ] Switching modes (tab/1/2) preserves the current color
- [ ] Opacity editing works in both modes
- [ ] `R` resets to exactly the original parsed color
- [ ] Hue is stable for grayscale colors (no jitter when S=0)
- [ ] Hue wraps at 0/360
- [ ] All other channels clamp to their valid ranges
- [ ] All tests pass
