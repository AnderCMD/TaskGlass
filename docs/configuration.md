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
  "maximized":  { "enabled": false, "accent": "opaque",  "color": "auto", "opacity": 100 },
  "apps": [
    { "exe": "code.exe", "accent": "acrylic", "color": "#101828", "opacity": 55 }
  ],
  "surfaces": {
    "tray_overflow": { "enabled": true, "opacity": 10 }
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
| `updates` | `notify`, `auto`, `off` | `notify` (default) checks GitHub once a day and asks before installing a new version; `auto` installs without asking; `off` never connects. See [Updates](#updates). |
| `language` | `auto`, `en`, `es` | Language of the tray menu and settings window. `auto` follows the Windows display language. |
| `theme` | `auto`, `light`, `dark` | Which colour mode drives `"auto"` colours, the tray icon and the tray menu (the menu follows the Windows *apps* mode when `auto`). `auto` follows the Windows *system* theme (the one the taskbar uses). |
| `taskbars` | `all`, `primary` | Style every monitor's taskbar or only the primary one. |
| `default` | appearance | Used when no other rule matches. |
| `maximized` | rule | Applies on a monitor that has a maximized window. Off by default: the bar stays as transparent as without a maximized window. |
| `apps` | list of app rules | Applies while the named application has focus (max 64, first match wins). |
| `surfaces` | object | Currently only `tray_overflow`, see [Tray flyout](#tray-flyout). |
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

A **rule** is an appearance plus `enabled`. An **app rule** is an appearance plus `exe`, the executable file name such as `chrome.exe` (case-insensitive, no path). An app rule applies to the taskbar of the monitor where that application has focus.

Opacity values accept decimals (`99.5`) and are rounded to a whole number.

## Tray flyout

`surfaces.tray_overflow` makes the flyout with the hidden tray icons (the `^` button) translucent: `{ "enabled": true, "opacity": N }`, where `opacity` is how opaque its background stays. `100` is the stock look and `0` removes the background completely (the icons then float over the desktop). It needs `taskglass_tap.dll` and `xaml_fix`, and a tray menu item turns it on or off.

Why the DLL: the flyout draws its own XAML background, which the accent effect does not reach, so TaskGlass loads its helper DLL into Explorer to change it, the same way it removes the taskbar's own fill.

`surfaces.tray_overflow.maximized` is a second opacity for the flyout, `{ "enabled": false, "opacity": 85 }`: while a window is maximized or fullscreen (a bright window, a game) on any monitor the flyout uses it instead of `opacity`, for example to keep the panel more solid there.

## Updates

With `"updates": "notify"` or `"auto"`, TaskGlass asks GitHub (`api.github.com`) for the latest release once a day, starting a minute and a half after it launches, and from the tray menu (**Check for updates…**) or the settings window. This is the only network connection TaskGlass makes; `"off"` disables it. Development builds never update.

When you accept (or in `auto` mode) TaskGlass downloads `TaskGlass-<version>-Setup.exe` from the release, checks its SHA-256 against the release's `SHA256SUMS.txt`, discards it if it does not match, and runs it silently; the installer closes TaskGlass, replaces it and starts the new version. Downloads and redirects are limited to GitHub hosts. The files are not code-signed, so the checksum protects against corruption and tampering in transit, not against a compromised GitHub account. If security software stops the installer from starting, the verified file stays in `%TEMP%\TaskGlass-update` and its path is copied to the clipboard.

## Windows setting that overrides the bar

Keep *Show accent color on Start and taskbar* (Settings > Personalization > Colors) turned **off**. When it is on, Windows paints the taskbar with the accent color and overrides whatever TaskGlass sets, whatever the rules say.

## Which rule wins

1. `enabled: false` → stock taskbar.
2. The first matching entry in `apps`.
3. `maximized` (when enabled and the taskbar's monitor has a maximized window).
4. `default`.

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
