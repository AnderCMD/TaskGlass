# Changelog

Format based on [Keep a Changelog](https://keepachangelog.com/); versioning follows [SemVer](https://semver.org/).

## [Unreleased]

### Changed
- **The settings window is the web-style one again** (the 1.2.x window, shown in an Edge app window; the native Win32 window from 1.3.0 is removed). It keeps the 1.3.0 improvements: Blur works on Windows 11, and the preview shows how Blur really looks there.

### Added
- The settings window now has the TaskGlass icon (title bar and taskbar) through a favicon served by the local settings page.

## [1.3.0] - 2026-11-20

### Fixed
- **Blur rendered solid black on Windows 11.** The classic blur-behind accent is drawn black by Windows 11 (measured on build 26300), whatever the tint. On Windows 11 Blur now uses the acrylic backdrop, and its strength lightens or darkens the glass (that backdrop takes its brightness from the tint colour and ignores the alpha). Windows 10 keeps the classic blur.

### Changed
- **The settings window is now native Windows** (real Win32 controls with the TaskGlass icon) instead of a web page, so there is no Edge window or local server any more. The live preview is drawn in-app with real blur, and shows the maximized and per-app looks and the hidden-icons flyout. The maximized rule has its own tab.

## [1.2.1] - 2026-11-20

### Fixed
- **Settings window showed "Cannot read properties of null".** With no per-app rules, 1.2.0 saved `"apps": null` to `config.json`, which the window could not read. `null` now loads as an empty list, is never written again, and the window tolerates it.

## [1.2.0] - 2026-11-20

### Fixed
- **Start with Windows did not start TaskGlass.** The Run-key value was written with Go's `%q`, which doubled every backslash (`"C:\Dir\taskglass.exe"`), and Windows could not launch it at sign-in. The path is now written in plain quotes. An entry written by an older version is repaired automatically the next time TaskGlass runs; turning "Start with Windows" on also clears a "disabled in Task Manager" mark.

### Added
- **Settings window** with a live preview of the taskbar (default, maximized and per-app looks) and of the hidden-icons flyout. Edit everything there: style, tint colour and strength, theme, language, tray icon, hotkey, per-app rules, flyout opacity, start with Windows and advanced timings. Open it with a left click on the tray icon, "Settings…" in the menu, `taskglass --settings` or by launching TaskGlass again. It is a local page (loopback only, protected by a per-run token) shown in an Edge app window.
- **English and Spanish** everywhere: tray menu, settings window, installer and license. New `language` setting (`auto`, `en`, `es`).
- **Installer**: asks for the taskbar style, language and extras on a first install and writes the initial `config.json`; shows the license in the installer language (`LICENSE.es.txt`, an unofficial translation, alongside the binding English `LICENSE`); adds a "TaskGlass Settings" Start-menu entry.
- Quick presets (pure glass, frosted, acrylic tint, solid).

### Changed
- Blur no longer renders opaque black when its tint is 0 (same floor Clear and Acrylic already had).

## [1.1.1] - 2026-11-18

### Changed
- **License**: TaskGlass is now released under the Apache License 2.0 (was MIT). The `NOTICE` file ships in the zip and the installer. Versions up to 1.1.0 stay available under MIT.

## [1.1.0] - 2026-11-18

### Added
- **Website** (GitHub Pages, `site/`): bilingual landing page with installer and portable downloads, install steps and FAQ for non-technical users.
- **Translucent tray flyout**: `surfaces.tray_overflow` makes the flyout with the hidden tray icons (the `^` button) translucent, with its own opacity, enabled by default, and a tray menu item. The helper DLL finds the flyout even when it already existed before TaskGlass started; for that TaskGlass briefly adds and removes a tray icon, and stops as soon as the DLL reports it is tracked.
- **Windows installer** (`TaskGlass-<version>-Setup.exe`, Inno Setup): per-user, no administrator rights, English and Spanish, optional start with Windows and desktop shortcut. Uninstalling asks the running TaskGlass to quit, restores the taskbar and offers to delete settings. `scripts/build-installer.ps1` builds it; the release workflow publishes it with checksums for both files and CI builds it on every run.
- Opacity values accept decimals (`99.5`) and are rounded; `null` and out-of-range values are rejected with a clear message.

### Changed
- The tray icon's right-click menu follows the Windows light/dark mode (or the `theme` you pin in the config) and switches live.
- Documented that *Show accent color on Start and taskbar* must be off in Windows, because it overrides the taskbar effect.
- **Defaults**: the taskbar stays fully transparent with a maximized window too (`maximized.enabled` is now `false`), and the tray flyout keeps 10 % opacity, just enough not to lose the icons.
- **Removed** the `start_menu` rule. An old config that still has the key keeps loading; the key is ignored and dropped on the next save.
- TaskGlass has a single purpose again: the taskbar and the tray flyout. Experiments with other surfaces (Start menu, File Explorer windows, per-application window opacity, look presets) were removed before release. A config that still has `surfaces.start_menu_flyout` or `surfaces.explorer` keeps loading; the keys are ignored and dropped on the next save.
- GitHub Actions bumped: `checkout` 7, `setup-go` 6, `upload-artifact` 7.
- The helper DLL finds XAML elements in every tree it has seen an event from, rescans after bursts of events and keeps looking while something wanted is missing.
- The debug dump (`tap.dump`) can be switched on and off while the DLL is loaded.

### Known limitations
- Context menus, notifications, calendar, quick settings, the Start menu and File Explorer windows are not styled.

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
