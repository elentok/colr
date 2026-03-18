# Migration: Bubbletea v1 → v2.0.2 (+ Lipgloss v2)

## Overview

Bubbletea v2 changes the import path to `charm.land/bubbletea/v2` and makes several breaking API changes. Lipgloss v2 (`charm.land/lipgloss/v2`) ships alongside it with its own breaking changes. This plan covers every required change in the colr codebase.

## Impact Summary

| Area | Files affected | Severity |
|------|---------------|----------|
| Import paths | All files importing bubbletea or lipgloss | Mechanical |
| `View()` return type | `app/view.go` | Moderate — new `tea.View` struct |
| `tea.WithAltScreen()` removal | `main.go` | Trivial — moves into `View()` |
| `tea.KeyMsg` → `tea.KeyPressMsg` | `app/update.go` | Low — one type switch |
| `tea.WindowSizeMsg` (verify) | `app/update.go` | Verify — may be renamed |
| `tea.Tick` (verify) | `app/copy.go` | Verify — confirm API unchanged |
| lipgloss color API | `app/view.go`, `ui/styles.go`, `ui/preview.go`, `ui/outputs.go`, `ui/help.go` | Moderate — `lipgloss.Color()` may change |
| Test type assertions | `app/copy_test.go`, `app/edit_test.go` | Low — `tea.Model` assertion may change |

## Step-by-step Plan

### Step 1: Update go.mod

- [ ] `go get charm.land/bubbletea/v2@v2.0.2`
- [ ] `go get charm.land/lipgloss/v2` (latest compatible version)
- [ ] Remove old `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss` requires
- [ ] `go mod tidy`

### Step 2: Rewrite import paths (all files)

Every file that imports bubbletea or lipgloss needs its import path updated:

| Old | New |
|-----|-----|
| `tea "github.com/charmbracelet/bubbletea"` | `tea "charm.land/bubbletea/v2"` |
| `"github.com/charmbracelet/lipgloss"` | `"charm.land/lipgloss/v2"` |

**Files to update:**

bubbletea imports:
- `main.go`
- `app/model.go`
- `app/update.go`
- `app/copy.go`

lipgloss imports:
- `app/view.go`
- `ui/styles.go`
- `ui/preview.go`
- `ui/outputs.go`
- `ui/help.go`
- `ui/editor.go` (uses `SliderStyle` from styles.go — same package, no direct import; verify)

### Step 3: `View()` return type (`app/view.go`)

**Before:**
```go
func (m Model) View() string {
    ...
    return frame
}
```

**After:**
```go
func (m Model) View() tea.View {
    ...
    v := tea.NewView(frame)
    v.AltScreen = true
    return v
}
```

- The `View()` method must return `tea.View` instead of `string`
- `tea.WithAltScreen()` is removed from `main.go` — the alt screen flag moves here as `v.AltScreen = true`
- The help overlay branch also needs to return `tea.View`:
  ```go
  if m.showHelp {
      v := tea.NewView(ui.RenderHelp(m.width, m.height))
      v.AltScreen = true
      return v
  }
  ```

### Step 4: Remove `tea.WithAltScreen()` from `main.go`

**Before:**
```go
p := tea.NewProgram(model, tea.WithAltScreen())
```

**After:**
```go
p := tea.NewProgram(model)
```

The alt screen is now declaratively set in `View()` (see Step 3).

### Step 5: `tea.KeyMsg` → `tea.KeyPressMsg` (`app/update.go`)

**Before:**
```go
case tea.KeyMsg:
    return handleKeyMsg(m, msg.String())
```

**After:**
```go
case tea.KeyPressMsg:
    return handleKeyMsg(m, msg.String())
```

The `msg.String()` method still exists on `tea.KeyPressMsg`, so `handleKeyMsg` and all downstream key handlers (`handleEditKey`, `handleKey`) remain unchanged — they operate on `string` values, not on the message type directly.

### Step 6: Verify `tea.WindowSizeMsg` (`app/update.go`)

Check whether `tea.WindowSizeMsg` is renamed in v2. If so, update the type switch case. The fields `.Width` and `.Height` are expected to remain the same.

### Step 7: Verify `tea.Tick` (`app/copy.go`)

```go
func scheduleToast(duration time.Duration) tea.Cmd {
    return tea.Tick(duration, func(time.Time) tea.Msg {
        return ClearToastMsg{}
    })
}
```

Confirm `tea.Tick` still exists and has the same signature. If renamed, update the call.

### Step 8: Lipgloss color API (`ui/styles.go`, `ui/preview.go`, `ui/outputs.go`, `ui/help.go`, `app/view.go`)

In lipgloss v1, colors are created with:
```go
lipgloss.Color("39")
lipgloss.Color("#FF0000")
```

In lipgloss v2, `lipgloss.Color()` still works but returns a different type. Verify:
- [ ] `lipgloss.Color("39")` (ANSI 256) still compiles and works
- [ ] `lipgloss.Color("#FF0000")` (hex string) still compiles and works
- [ ] `lipgloss.NewStyle()` API is unchanged
- [ ] `lipgloss.NormalBorder()`, `lipgloss.RoundedBorder()` still exist
- [ ] `lipgloss.JoinHorizontal()`, `lipgloss.JoinVertical()` still exist
- [ ] `lipgloss.Place()` still exists (used in `ui/help.go`)
- [ ] `.Width()`, `.Height()`, `.Border()`, `.BorderForeground()`, `.Foreground()`, `.Background()` chainable methods still work
- [ ] `.Bold()`, `.Italic()`, `.Render()` still work

If the color constructor changes, update all call sites. These are spread across:
- `ui/styles.go` — 16 `lipgloss.Color()` calls (all style definitions)
- `ui/preview.go` — 4 `lipgloss.Color()` calls (hex color, checker backgrounds, fg hint)
- `ui/outputs.go` — 1 `lipgloss.Color()` call
- `ui/help.go` — uses styles from `ui/styles.go`, plus `helpBorderStyle`
- `app/view.go` — 2 `lipgloss.Color()` calls (panel border, outer frame border)

### Step 9: Verify `tea.Model` interface (`app/model.go`)

The `tea.Model` interface in v2 requires:
- `Init() tea.Cmd` — unchanged
- `Update(tea.Msg) (tea.Model, tea.Cmd)` — unchanged
- `View() tea.View` — **changed** (handled in Step 3)

Confirm `Model` satisfies the new interface after the `View()` change.

### Step 10: Update tests (`app/copy_test.go`, `app/edit_test.go`)

Tests use `m.Update(ClearToastMsg{})` and assert `result.(Model)`. The `Update` signature is unchanged (`(tea.Model, tea.Cmd)`), so the type assertion should still work.

However, verify:
- [ ] `tea.Msg` is still the base message interface
- [ ] No test imports need updating (tests don't import bubbletea directly — they call `handleKeyMsg` which takes strings)

The only test file that calls `m.Update()` directly is `copy_test.go:TestClearToastMsg`. This should work unchanged since `Update` still returns `(tea.Model, tea.Cmd)`.

### Step 11: Build and test

- [ ] `go build ./...` — fix any compilation errors
- [ ] `go test ./...` — all tests pass
- [ ] `go vet ./...` — no warnings
- [ ] Manual test: run the app, verify all keybindings, alt screen, visual rendering

### Step 12: Cleanup

- [ ] Remove any `// indirect` deps from go.sum that are no longer needed (`go mod tidy`)
- [ ] Verify no old `github.com/charmbracelet/bubbletea` or `github.com/charmbracelet/lipgloss` references remain in the codebase

## Risk areas

1. **Lipgloss color API** — The biggest unknown. If `lipgloss.Color("39")` no longer works the same way, every style definition needs updating. This affects visual appearance of the entire app.

2. **`tea.Tick`** — If the timer API changed, the toast system breaks. Low risk — `tea.Tick` is fundamental and likely preserved.

3. **Border rendering** — `injectBorderTitle` in `app/view.go` relies on the literal `╭──────` string being present in lipgloss's rendered output. If lipgloss v2 changes how borders are rendered (different characters, different ANSI wrapping), this breaks silently.

4. **`lipgloss.Place()`** — Used in the help overlay. Verify it exists in v2.

## Estimated scope

~10 files changed, mostly mechanical import path updates. The only structural change is `View() string` → `View() tea.View` and moving `AltScreen` from program option to view field. No logic changes needed.
