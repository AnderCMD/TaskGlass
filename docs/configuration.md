# Configuration

File: `%APPDATA%\TaskGlass\config.json`, created with defaults on first run. Open it from the tray menu (**Open config file**). Changes are picked up within about a second, no restart needed. If the file is invalid the running settings are kept and the error is written to the log. You can also check it first:

```bash
taskglass --validate
```

Every key is optional: a missing key keeps its default. Unknown keys are rejected, so typos are caught. Add the `$schema` line (already in generated files) to get completion and validation in VS Code and other editors; the schema lives in [config.schema.json](config.schema.json).

Changes made from the tray menu are saved to this same file.

```json
{
  "$schema": "https://raw.githubusercontent.com/AnderCMD/TaskGlass/main/docs/config.schema.json",
  "version": 1,
  "enabled": true,
  "theme": "auto",
  "taskbars": "all",
  "default":    { "accent": "clear",  "color": "auto", "opacity": 0 },
  "maximized":  { "enabled": true,  "accent": "opaque",  "color": "auto", "opacity": 100 },
  "start_menu": { "enabled": false, "accent": "acrylic", "color": "auto", "opacity": 60 },
  "apps": [
    { "exe": "code.exe", "accent": "acrylic", "color": "#101828", "opacity": 55 }
  ],
  "tray_icon": "auto",
  "hotkey": "Ctrl+Alt+G",
  "poll_ms": 500,
  "debounce_ms": 25,
  "refresh_ms": 1000,
  "xaml_fix": true,
  "log_level": "info"
}
```

## Options

| Key | Values | Meaning |
|-----|--------|---------|
| `enabled` | bool | Master switch. `false` shows the stock taskbar. |
| `theme` | `auto`, `light`, `dark` | Which colour mode drives `"auto"` colours and the tray icon. `auto` follows the Windows *system* theme (the one the taskbar uses). |
| `taskbars` | `all`, `primary` | Style every monitor's taskbar or only the primary one. |
| `default` | appearance | Used when no other rule matches. |
| `maximized` | rule | Applies on a monitor that has a maximized window. |
| `start_menu` | rule | Applies while the Start menu or Search has focus. |
| `apps` | list of app rules | Applies while the named application has focus (max 64, first match wins). |
| `tray_icon` | `auto`, `color`, `light`, `dark` | Notification-area icon. `auto` is a monochrome glyph that contrasts with the Windows theme; `color` is the full logo. |
| `hotkey` | string | Global shortcut toggling `enabled`, for example `"Ctrl+Alt+G"`. Needs at least one of `Ctrl`, `Alt`, `Shift`, `Win` and a key `A-Z`, `0-9` or `F1-F24`. Empty (default) disables it. |
| `poll_ms` | 50-10000 | Fallback sampling interval. Window events normally trigger updates instantly; this only guards against missed events. |
| `debounce_ms` | 0-1000 | Coalesces bursts of window events into one update. |
| `refresh_ms` | 100-60000 | Re-apply the effect this often even if nothing changed, to survive Explorer silently resetting it. |
| `xaml_fix` | bool | Remove Windows 11's own XAML taskbar fill (needs `taskglass_tap.dll` next to the exe). Without it `clear` cannot be fully transparent. `false` never touches Explorer. |
| `log_level` | `debug`, `info`, `warn`, `error` | Log verbosity. |

### Appearance

| Key | Values | Meaning |
|-----|--------|---------|
| `accent` | `clear`, `blur`, `acrylic`, `opaque`, `normal` | `clear`: fully transparent. `blur`: blur with a tint. `acrylic`: blur plus noise. `opaque`: solid colour. `normal`: stock taskbar. |
| `color` | `#RRGGBB` or `auto` | Tint colour. `auto` follows the theme: `#1C1C1C` in dark mode, `#F3F3F3` in light mode (matching the stock taskbar). |
| `opacity` | 0-100 | Tint opacity in percent. |

A **rule** is an appearance plus `enabled`. An **app rule** is an appearance plus `exe`, the executable file name such as `chrome.exe` (case-insensitive, no path). An app rule applies to the taskbar of the monitor where that application has focus.

## Which rule wins

1. `enabled: false` → stock taskbar.
2. `start_menu` (when enabled and the Start menu or Search has focus).
3. The first matching entry in `apps`.
4. `maximized` (when enabled and the taskbar's monitor has a maximized window).
5. `default`.

## Light and dark mode

Colours set to `"auto"` and the `auto` tray icon follow Windows: switch Windows between light and dark and TaskGlass adapts within a second, no restart. Set `"theme": "light"` or `"dark"` to pin the behaviour regardless of Windows. A fixed colour such as `"#102030"` is never changed by the theme.

## Command line

| Flag | Purpose |
|------|---------|
| `--config <path>` | Use a custom config file |
| `--validate` | Check the config and exit; non-zero on errors |
| `--reset` | Rewrite the config with defaults and exit |
| `--quit` | Ask the running TaskGlass to exit (restores the taskbar) |
| `--restore` | Put the taskbar back to normal and exit |
| `--verbose` | Also log to stderr |
| `--version` | Print the version |

`taskglass.exe` is a GUI-subsystem program so it never flashes a console window; shells do not wait for such programs. To use its exit code in a script, run it with `Start-Process -Wait -PassThru` (PowerShell) or `start /wait` (cmd).
