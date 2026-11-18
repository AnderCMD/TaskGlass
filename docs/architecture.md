# Architecture

```
cmd/taskglass        entry point: flags, single instance, logging, exe resources
internal/app         tray menu + main loop (ties everything together)
internal/config      schema, defaults, validation, hotkey parsing, atomic save   (portable)
internal/state       pure rule resolution: context -> appearance                  (portable)
internal/taskbar     applies appearances per taskbar, restores
internal/events      WinEvent hooks + global hotkey on a dedicated message-loop thread
internal/theme       reads the Windows light/dark mode
internal/winapi      the only package that calls Win32
internal/xaml        loads and drives the native TAP inside Explorer
internal/autostart   HKCU Run key
native/tap           C source of taskglass_tap.dll (runs inside explorer.exe)
assets/              embedded tray icons; assets/src holds the SVG/HTML sources
docs/                documentation, config schema, README images
scripts/             build-tap, gen-resources, render-assets
tools/genicon        packs rendered PNGs into .ico files
```

## How the effect works

Explorer's taskbar windows (`Shell_TrayWnd` for the primary monitor, `Shell_SecondaryTrayWnd` for the others) accept an accent policy through the undocumented `user32!SetWindowCompositionAttribute` (attribute 19). TaskGlass sets that policy (gradient, transparent gradient, blur-behind or acrylic with an `0xAABBGGRR` tint). Disabling the policy and sending `WM_DWMCOMPOSITIONCHANGED` makes Explorer repaint the stock look.

## Main loop

The loop in `internal/app` runs from a single goroutine and wakes on whichever comes first:

- **a window event** from `internal/events` (foreground change, show/hide/destroy, minimise/restore, move/resize of a top-level window). A WinEvent hook callback does nothing but a non-blocking send on a 1-slot channel; the loop then waits `debounce_ms` and drains the channel, so a burst of events costs one update;
- **the hotkey**, delivered by the same event thread (`RegisterHotKey` is thread-affine, so it lives there);
- **a kick** after a config change or tray click;
- **a 1 s housekeeping tick**, which hot-reloads the config (mtime check) and re-reads the Windows theme;
- **`poll_ms`**, a safety-net timer for missed events.

Each wake-up runs `taskbar.Manager.Tick`:

1. finds all taskbars (new monitors and Explorer restarts are handled for free; vanished handles are dropped; taskbars no longer selected by `taskbars: primary` are handed back),
2. only if a rule needs it, enumerates top-level windows once for monitors with a maximized, visible, non-cloaked, non-tool window, and reads the focused window's executable,
3. checks whether the foreground window is the Start/Search host (`Windows.UI.Core.CoreWindow` owned by `StartMenuExperienceHost` / `SearchHost` / `SearchApp`),
4. resolves each taskbar's appearance with `state.Resolve` (priority: master switch, Start menu, app rule, maximized, default) and resolves `auto` colours against the theme,
5. applies it only if the resolved policy changed or `refresh_ms` elapsed.

With nothing happening the process sleeps in the scheduler: measured at 0.0 ms of CPU over 8 s.

## Lifecycle and recovery

- `Run` returns after the tray quits, then `shutdown` runs synchronously and idempotently: the loop stops, the XAML fill is shown again, every accent is reset, hooks are removed. The tray library may return from `Run` before `onExit` finishes, so the process waits for `shutdown` itself.
- A panic inside an update is recovered and logged, the taskbar is restored, and the loop carries on. More than 5 panics in a minute and the app restores the taskbar and exits instead of flapping.
- `--quit` sets a named event the running instance waits on; Ctrl+C and SIGTERM do the same.
- If TaskGlass is killed, the TAP inside Explorer notices that the shared block vanished and clears the accent itself (see below). `--restore` is the manual fallback.
- An invalid config never replaces the running one; the error is logged once per edit.

## Removing the Windows 11 XAML fill

The Windows 11 taskbar is a XAML island that draws a grey fill (`BackgroundFill`) and a border (`BackgroundStroke`) over the DWM accent. `native/tap/tap.c` is a XAML Diagnostics *TAP* (Tool Assembly Provider) that hides them. The scheme mirrors the one used by TranslucentTB:

1. `internal/xaml` copies the DLL to a content-addressed path in the cache folder (so an updated DLL never collides with a copy Explorer still has loaded), installs a `WH_CALLWNDPROC` hook on the taskbar's UI thread with the DLL's exported `TaskGlassHookProc`, and nudges the thread with `WM_NULL`. Windows loads the DLL into Explorer on that thread.
2. In `DllMain` (only inside `explorer.exe`) the DLL starts a thread that calls `InitializeXamlDiagnosticsEx` on its own process, retrying for up to 30 s. Right after Explorer starts XAML is not ready and every endpoint fails with `ELEMENT_NOT_FOUND`; the endpoint is a small number (`VisualDiagConnection1`..) that does not grow with attempts, so each round tries the whole range of names. A fresh thread is used per call because XAML Diagnostics initialises once per thread.
3. XAML then creates the TAP object from the same DLL and calls `SetSite`; a separate thread calls `AdviseVisualTreeChange`.
4. On this Windows build the initial tree replay does not include the taskbar island, only live changes do. So on any live element (`Add`) the DLL climbs to the root of that element's tree with `VisualTreeHelper.GetParent`, walks it down, and keeps references to `Rectangle`s named `BackgroundFill` / `BackgroundStroke` inside a `TaskbarFrame`. If no natural change happens, TaskGlass shows a tiny off-screen window for 300 ms (`winapi.TaskbarAnchor`) to provoke one.
5. Opacity is set through `IUIElement::put_Opacity` (vtable slot read from `Windows.UI.Xaml.winmd`). XAML objects are only touched on the UI thread: from the visual tree callback, or via `CoreDispatcher.RunAsync`.

Control channel: the named shared memory `Local\TaskGlass.Xaml.State.v1` holds two int32s, `[0]` mode (TaskGlass to DLL: 0 leave alone, 1 hide) and `[1]` status (DLL to TaskGlass: 1 attached, 2 controlling the fill, negative = failed HRESULT). The DLL polls it every 150 ms and republishes status 2 to every new block, so a restarted TaskGlass learns that an earlier run's DLL is still working. If the block does not exist (TaskGlass exited or crashed) the mode is 0, so the stock taskbar returns by itself, and the DLL additionally clears the accent policy from inside Explorer.

The BSTRs that XAML Diagnostics hands to the visual tree callback are owned by the callback and are freed there. The DLL's log is capped at 1 MiB because it lives inside Explorer.

Debugging: create an empty `%LOCALAPPDATA%\TaskGlass\tap.dump` *before* the DLL is loaded and `tap.log` in the same folder records every visual tree change and scan. The DLL stays loaded until Explorer restarts.

## The tray overflow flyout

The same DLL controls one more XAML surface. The shared block has eight int32 slots: `[0]` taskbar fill, `[1]` status, `[2]` tray overflow flyout, `[7]` flyout tracked. Slot `[2]` holds `0` (leave alone) or `1 + opacity percent`; the DLL applies `opacity / 100` to the elements it tracks and `1.0` when told to leave them alone or when the block disappears. Slots `[3]` to `[6]` are unused and must stay `0`: copies of the DLL from unreleased builds can still be loaded in a running Explorer, and read them as the opacity of surfaces that no longer exist. Older DLLs read only slots 0 and 1, so the block stays compatible.

The flyout is a separate XAML island inside Explorer, created the first time the `^` button is opened. It has a `Border` named `OverflowFlyoutBackgroundBorder`. Its tree is not reachable from the taskbar's (both islands live on the same UI thread, but the diagnostics replay only contains one root), so the DLL keeps the last eight live elements it saw as anchors, climbs from each to its tree root and scans every distinct root (the scan runs 250 ms after a burst of events ends, and every 5 s while something wanted is still missing). `IXamlDiagnostics::HitTest` returns `E_INVALIDARG` on this build, so enumerating trees from coordinates is not an option.

If the island already exists when the DLL attaches, no event arrives for it by itself. `winapi.TrayProbe` provokes some: it adds a tray icon to a message-only window and removes it 400 ms later, which makes Explorer build and tear down a tray item and the tray's XAML emit events. The DLL sets slot `[7]` to 1 once it tracks the flyout, and `internal/xaml` stops probing (at most four attempts, six seconds apart, per Explorer process). Same-version restarts of TaskGlass need no probe because the DLL stays loaded and keeps its references.

## Appearance and assets

`assets/src` holds the logo and tray glyph SVGs and the banner / social preview pages. `scripts/render-assets.ps1` rasterises them with headless Edge and `tools/genicon` packs the `.ico` files; outputs are committed, so building never needs Edge. The exe icon, version info and manifest come from `cmd/taskglass/taskglass.rc` through `scripts/gen-resources.ps1`, which produces the committed `resource_windows_amd64.syso`.

## Design choices

- **Go without cgo**: `golang.org/x/sys/windows` plus lazy DLL procs; a single static binary. The only C is the helper DLL, which has to run inside Explorer.
- **Pure logic isolated** in `config` and `state` so it is tested on any OS; Win32 code is behind `_windows` build tags.
- **Events over polling**: WinEvent hooks wake the loop only when something relevant happens, with a slow poll as a safety net.
- **Fail soft**: invalid config keeps the previous settings; apply errors are logged and retried; if the TAP cannot load, accent effects still work.
- **Config file stays the source of truth**: tray choices are written to it, and hand edits are picked up live.
