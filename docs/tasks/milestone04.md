# Milestone 4 — Clipboard Copy, Help Overlay, Toast Messages

## Goal

Add the `yr`/`yx`/`yh` copy commands with clipboard write, a help overlay toggled by `?`, and timed toast messages for success/error feedback.

## Existing Code Touched

- `app/model.go` — add `pendingY bool` for the two-key `y` prefix
- `app/update.go` — `y` prefix handling, `?` toggle, help-overlay key interception
- `app/view.go` — render help overlay on top when `showHelp` is true
- `ui/help.go` — new: help overlay renderer
- `ui/header.go` — already reads `toastMsg`, no changes needed

## Design Decisions

### Two-key `y` prefix

The `yr`, `yx`, `yh` commands are two-key sequences. When `y` is pressed:
- Set `pendingY = true` on the model
- On next key:
  - `r` → copy RGB, clear `pendingY`
  - `x` → copy HEX, clear `pendingY`
  - `h` → copy HSL, clear `pendingY`
  - anything else → clear `pendingY`, ignore (no action)

### Help overlay is modal

When `showHelp` is true:
- Render a centered overlay on top of the normal view
- Only `?`, `Esc` dismiss the overlay
- `q` **inside help** closes help (does not quit the app) — spec says "q still quits app globally only if not in help overlay"
- All other keys are ignored while help is showing

### Toast timer

Use `tea.Tick` to schedule a `ClearToastMsg` after a delay:
- 1 second for success messages
- 2 seconds for error messages

The `ClearToastMsg` type already exists on the model. Wire up the tick command.

## Implementation Steps

### Step 1: Add `pendingY` to model (`app/model.go`)

- [ ] Add `pendingY bool` field to `Model`

### Step 2: Toast timer command (`app/model.go`)

- [ ] `toastCmd(msg string, duration time.Duration) tea.Cmd` — sets `toastMessage` on model and returns a `tea.Tick` that sends `ClearToastMsg` after `duration`
- [ ] Actually: set toast message directly, return the tick cmd from Update

### Step 3: Copy commands (`app/copy.go`)

- [ ] `applyCopy(m Model, format string) (Model, tea.Cmd)` — copies a format to clipboard, returns updated model + toast tick
  - `"rgb"` → `color.FormatRGB(m.currentColor)` → clipboard write
  - `"hex"` → `color.FormatHEX(m.currentColor)` → clipboard write
  - `"hsl"` → `color.FormatHSL(m.currentColor)` → clipboard write
- [ ] On success: set `toastMessage = "Copied RGB to clipboard"`, return 1s tick
- [ ] On error: set `toastMessage = "Failed to copy RGB to clipboard"`, return 2s tick
- [ ] Uses `clipboard.Write()` from our clipboard wrapper

### Step 4: Update key handling (`app/update.go`)

- [ ] **Help overlay intercept:** when `showHelp` is true, only handle:
  - `?` → close help
  - `Esc` → close help
  - `q` → close help (NOT quit app)
  - `ctrl+c` → quit (always works)
  - All other keys → ignore
- [ ] **`?` key:** toggle `showHelp`
- [ ] **`y` prefix:** when `pendingY` is false and `y` is pressed, set `pendingY = true`
- [ ] **`y` continuation:** when `pendingY` is true, dispatch:
  - `r` → `applyCopy(m, "rgb")`
  - `x` → `applyCopy(m, "hex")`
  - `h` → `applyCopy(m, "hsl")`
  - anything else → clear `pendingY` only
- [ ] Note: since `Update` needs to return a `tea.Cmd` from copy, refactor to return `(Model, tea.Cmd)` from the key handler path

### Step 5: Help overlay renderer (`ui/help.go`)

- [ ] `RenderHelp(width, height int) string`
- [ ] Content (from spec):
  ```
  colr help

  Movement
    j/k or ↑/↓    move between fields
    h/l or ←/→    adjust field
    H/L            adjust by larger step
    g/G            first/last field

  Modes
    tab            toggle HSV/RGB
    1              HSV mode
    2              RGB mode

  Copy
    yr             copy RGB
    yx             copy HEX
    yh             copy HSL

  Other
    R              reset color
    ?              toggle this help
    q              quit
  ```
- [ ] Render in a bordered, centered box using Lip Gloss
- [ ] Box should be sized to content, not fill the terminal

### Step 6: View composition for help overlay (`app/view.go`)

- [ ] When `m.showHelp` is true: render the help overlay centered on top of the normal view
- [ ] Use `lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, helpBox)` to center the help overlay within the full terminal area
- [ ] The normal view should still render underneath (dimmed or not — keep it simple for now)

### Step 7: Tests

- [ ] **Copy command test (`app/copy_test.go`):**
  - Verify `applyCopy` sets correct toast message for each format
  - Verify returned command is non-nil (tick scheduled)
- [ ] **Help overlay toggle test:**
  - `?` toggles `showHelp` on/off
  - `Esc` closes help
  - `q` in help closes help (doesn't quit)
- [ ] **`y` prefix test:**
  - `y` then `r` → triggers copy
  - `y` then `x` → triggers copy
  - `y` then unknown key → clears pendingY, no action
  - `y` alone → sets pendingY
- [ ] **Toast clear test:**
  - Verify `ClearToastMsg` clears `toastMessage`

### Step 8: Integration verification

- [ ] `go build ./...` clean
- [ ] `go test ./...` all pass
- [ ] Manual: `yr` copies RGB to clipboard, toast appears and fades
- [ ] Manual: `yx` copies HEX to clipboard
- [ ] Manual: `yh` copies HSL to clipboard
- [ ] Manual: `?` shows help overlay, `Esc` closes it
- [ ] Manual: `q` inside help closes help, doesn't quit

## Scope Boundaries

**In scope:** `yr`/`yx`/`yh` copy, help overlay, toast with timers, `y` prefix handling.

**Deferred to milestone 5:** styling polish, `yy`/Enter for selected output row, checkerboard refinement.

## Acceptance Criteria

- [ ] `yr` copies normalized RGB to clipboard
- [ ] `yx` copies normalized HEX to clipboard
- [ ] `yh` copies normalized HSL to clipboard
- [ ] Success toast shows for ~1 second
- [ ] Failed clipboard write shows error toast for ~2 seconds (app stays running)
- [ ] `?` toggles help overlay
- [ ] Help is dismissible with `?`, `Esc`, or `q`
- [ ] `q` inside help does NOT quit the app
- [ ] `ctrl+c` always quits regardless of help state
- [ ] `y` followed by an invalid key is a no-op
- [ ] All tests pass
