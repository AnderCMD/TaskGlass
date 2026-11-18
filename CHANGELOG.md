# Changelog

Format based on [Keep a Changelog](https://keepachangelog.com/); versioning follows [SemVer](https://semver.org/).

## [1.1.0] - 2026-11-18

### Added
- **Windows installer** (`TaskGlass-<version>-Setup.exe`, Inno Setup): per-user, no administrator rights, English and Spanish, optional start with Windows and desktop shortcut. Uninstalling asks the running TaskGlass to quit, restores the taskbar and offers to delete settings. `scripts/build-installer.ps1` builds it; the release workflow publishes it with checksums for both files and CI builds it on every run.
- **Glass surfaces** (`surfaces` in the config, **Glass surfaces** in the tray menu), each with its own opacity and enabled by default:
  - `tray_overflow`: the flyout with hidden tray icons.
  - `start_menu_flyout`: the Start menu and its phone panel (the helper DLL is now also loaded into `StartMenuExperienceHost.exe`, hooked through its UI threads because it owns no window while Start is closed).
- **Looks**: `taskglass --preset liquid|clear|subtle|stock`, `--list-presets` and a **Look preset** tray submenu. `liquid` frosts the taskbar and makes the flyout and Start menu translucent; presets keep your app rules, hotkey and theme.
- **Window opacity per application**: `window_opacity` in an app rule makes every window of that application uniformly translucent. A rule may use only `window_opacity` and no taskbar accent.
- Opacity values accept decimals (`99.5`) and are rounded; `null` and out-of-range values are rejected with a clear message.
- The shared-memory block grew to 8 slots (still readable by older DLLs). Groundwork for `ShellExperienceHost` (notification centre) is in the DLL and controller but not exposed: see `docs/architecture.md`.

### Changed
- GitHub Actions bumped: `checkout` 7, `setup-go` 6, `upload-artifact` 7.
- The helper DLL finds XAML elements in every tree it has seen an event from, rescans after bursts of events and keeps looking while something wanted is missing.
- The tray overflow flyout is found even when it already existed before TaskGlass started or before an updated DLL attached: TaskGlass briefly adds and removes a tray icon to make the tray emit the events the DLL needs, and stops as soon as the DLL reports it is tracked.
- The debug dump (`tap.dump`) can be switched on and off while the DLL is loaded.

### Known limitations
- Context menus, notifications, calendar and quick settings are not styled yet.

## [1.0.1] - 2026-11-17

### Changed
- The global hotkey now has an explicit switch, `hotkey_enabled`, which is **off by default**; the shortcut itself (`hotkey`, default `Ctrl+Alt+G`) is preconfigured so enabling it is one change. It can also be toggled from the tray menu. The shortcut is validated even while disabled.
- Migration: a 1.0.0 config that set `hotkey` without `hotkey_enabled` no longer registers the shortcut until you add `"hotkey_enabled": true`.

## [1.0.0] - 2026-11-17

First stable release.

### Added
- Taskbar styles `clear`, `blur`, `acrylic`, `opaque` and `normal`.
- **True transparency on Windows 11**: `taskglass_tap.dll` removes the taskbar's own XAML fill and border (`xaml_fix` option). It is loaded into Explorer through the XAML Diagnostics API and reverts by itself if TaskGlass exits or crashes.
- **Light and dark mode**: `theme` (`auto`, `light`, `dark`), `auto` colours that follow Windows live, and a monochrome tray icon that contrasts with the theme (`tray_icon`).
- **Rules**: `maximized`, `start_menu`, per-application `apps`, with a documented priority; `taskbars` to style every monitor or only the primary.
- **Master switch** `enabled`, tray menu with style and theme submenus, and a global `hotkey` to toggle the effect.
- Config: `$schema` and a JSON Schema for editor completion, strict validation with clear messages, `--validate`, hot reload that keeps the previous settings on errors, atomic saves.
- Command line: `--validate`, `--quit`, `--restore`, `--reset`, `--config`, `--verbose`, `--version`.
- Logo, banner, social preview and generated icons; exe icon, version info and manifest.
- Documentation in English and Spanish, architecture notes, contribution guide, CI and release workflows.

### Changed
- The engine is event driven (WinEvent hooks) instead of polling: about 0 % CPU when idle. `poll_ms` is now only a safety net (default 500 ms) and `debounce_ms` is new.
- Window scans are skipped entirely when no rule needs them.

### Fixed
- The process could exit before the shutdown cleanup finished, leaving the accent applied; shutdown is now synchronous and idempotent.
- After an Explorer restart the XAML fill removal could fail to reconnect; endpoint names are now retried correctly and a half-loaded DLL is no longer mistaken for a working one.
- After a forced kill, a leftover accent tint stayed on the taskbar; the DLL now clears it from inside Explorer.
- Internal panics are recovered instead of killing the app.

### Known limitations
- Verified on Windows 11 build 26300 with one monitor and auto-hide. Windows 10 and multi-monitor setups are untested.
- File Explorer windows and the Start menu flyout are not styled.
