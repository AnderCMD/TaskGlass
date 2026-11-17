# Changelog

Format based on [Keep a Changelog](https://keepachangelog.com/); versioning follows [SemVer](https://semver.org/).

## [Unreleased]

### Added
- Initial release: clear / blur / acrylic / opaque / normal taskbar styles.
- Dynamic rules for maximized windows and the Start menu.
- `taskglass_tap.dll`: removes the Windows 11 XAML taskbar fill and border for true transparency (`xaml_fix` option), loaded into Explorer via the XAML Diagnostics API.
- Multi-monitor support, hot-reloaded JSON config, tray menu, autostart.
- `--restore`, `--reset`, `--config`, `--verbose`, `--version` flags.
