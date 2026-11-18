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
    { "exe": "code.exe", "accent": "acrylic", "color": "#101828", "opacity": 55 },
    { "exe": "chrome.exe", "window_opacity": 92 }
  ],
  "surfaces": {
    "tray_overflow":     { "enabled": true, "opacity": 60 },
    "start_menu_flyout": { "enabled": true, "opacity": 60 },
    "explorer":          { "enabled": true, "opacity": 85 }
  },
  "tray_icon": "auto",
  "hotkey_enabled": false,
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
| `apps` | list of app rules | Per-application rules (max 64, first match wins): restyle the taskbar while the app has focus and/or make its windows translucent with `window_opacity`. |
| `surfaces` | object | Translucent shell surfaces, see [Glass surfaces](#glass-surfaces). |
| `tray_icon` | `auto`, `color`, `light`, `dark` | Notification-area icon. `auto` is a monochrome glyph that contrasts with the Windows theme; `color` is the full logo. |
| `hotkey_enabled` | bool | Turns the global hotkey on. **Off by default**, so TaskGlass never grabs a shortcut you did not ask for. Also in the tray menu (**Global hotkey**). |
| `hotkey` | string | The shortcut that toggles `enabled` while `hotkey_enabled` is `true`. Default `"Ctrl+Alt+G"`. Needs at least one of `Ctrl`, `Alt`, `Shift`, `Win` and a key `A-Z`, `0-9` or `F1-F24`. It is checked even while the hotkey is off, so a typo is reported right away. |
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

A **rule** is an appearance plus `enabled`. An **app rule** has `exe`, the executable file name such as `chrome.exe` (case-insensitive, no path), plus either or both of:

- an appearance (`accent`, `color`, `opacity`), applied to the taskbar of the monitor where that application has focus;
- `window_opacity` (0-100), which makes every window of the application uniformly translucent while TaskGlass runs, content included. It works from outside the application with the layered-window alpha, like the `explorer` surface; windows closed or no longer matching a rule are given back. If TaskGlass is killed, windows can stay translucent until they are closed. Dialogs, tool windows and tiny windows are left alone.

Opacity values accept decimals (`99.5`) and are rounded to a whole number.

## Presets

A preset sets the whole look in one go. It replaces the taskbar appearances, the `maximized` and `start_menu` rules and the `surfaces`, and keeps your app rules, hotkey, theme and everything else.

```bash
taskglass --list-presets
taskglass --preset liquid
```

The tray menu has a **Look preset** submenu doing the same. A running TaskGlass picks the change up within a second.

| Preset | Look |
|--------|------|
| `liquid` | Frosted (acrylic) taskbar, translucent Start menu, tray flyout and Explorer windows. The closest to macOS-style glass that Windows allows. |
| `clear` | Fully transparent taskbar, translucent flyout and Start menu, Explorer untouched. The default look of 1.0. |
| `subtle` | Light blur on the taskbar, barely translucent shell. |
| `stock` | Turns TaskGlass off (the stock Windows look). |

Real macOS "liquid glass" refracts what is behind it and adds specular highlights. Windows exposes no way to add that to other programs, so `liquid` approximates it with blur, transparency and theme-aware tints.

## Glass surfaces

`default`, `maximized` and the other rules style the **taskbar**. `surfaces` styles other parts of the shell that the accent effect cannot reach. Each entry is `{ "enabled": true, "opacity": N }`, where `opacity` is how opaque the surface stays: `100` is the stock look, lower values let more of the desktop show through.

| Surface | What it is | How | Needs the DLL |
|---------|------------|-----|---------------|
| `tray_overflow` | The flyout with hidden tray icons (the `^` button) | Its XAML background is made translucent | yes (`xaml_fix`) |
| `start_menu_flyout` | The Start menu and its phone panel | The acrylic background of both is made translucent | yes (`xaml_fix`) |
| `explorer` | File Explorer windows | The whole window (title bar, list, navigation pane) gets a uniform opacity | no |

Notes:

- `tray_overflow` and `start_menu_flyout` load `taskglass_tap.dll` into Explorer and into `StartMenuExperienceHost.exe`. `"xaml_fix": false` disables both, together with the taskbar fill removal.
- At `0` the Start menu and the flyout lose their background (and its blur) completely, so text can be hard to read over busy wallpapers; `40`-`70` is usually nicer.
- `explorer` is applied from outside Explorer, so it works even when Explorer runs each window in its own process. The value is clamped to at least `10`. If TaskGlass is killed, windows can stay translucent until you close them or run `taskglass --restore`; a normal exit restores them.
- The tray menu has a **Glass surfaces** submenu to turn each on or off.

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
| `--preset <name>` | Apply a look (`liquid`, `clear`, `subtle`, `stock`) and exit |
| `--list-presets` | List the presets |
| `--reset` | Rewrite the config with defaults and exit |
| `--quit` | Ask the running TaskGlass to exit (restores the taskbar) |
| `--restore` | Put the taskbar and File Explorer windows back to normal and exit |
| `--verbose` | Also log to stderr |
| `--version` | Print the version |

`taskglass.exe` is a GUI-subsystem program so it never flashes a console window; shells do not wait for such programs. To use its exit code in a script, run it with `Start-Process -Wait -PassThru` (PowerShell) or `start /wait` (cmd).
