# Changelog

Format based on [Keep a Changelog](https://keepachangelog.com/); versioning follows [SemVer](https://semver.org/).

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
