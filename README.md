# TaskGlass

[![CI](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml/badge.svg)](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A tiny, dependency-light tray app that makes the Windows taskbar **transparent, blurred or acrylic**, so only the icons remain. Written in Go, no cgo, no runtime, a single ~3 MB `.exe`.

[Leer en español](README.es.md)

## Features

- **Clear / Blur / Acrylic / Opaque / Normal** taskbar styles
- **Dynamic rules**: become opaque when a window is maximized, restyle while the Start menu is open
- **Multi-monitor**: every taskbar is handled independently, including ones added later
- **Survives Explorer restarts** and silently re-applies the effect
- **Hot-reloaded JSON config**; a typo never breaks a running session (the previous settings are kept and the error is logged)
- Tray menu, optional start-with-Windows, single instance
- Restores the stock taskbar on exit, and offers `--restore` for crash recovery

## Install

Download `taskglass.exe` from the [Releases](https://github.com/AnderCMD/TaskGlass/releases) page and run it. It lives in the notification area.

Build from source (Go 1.26+, Windows):

```bash
go build -ldflags "-H=windowsgui -s -w" -o dist/taskglass.exe ./cmd/taskglass
```

## Usage

Right-click the tray icon to pick a style, toggle the dynamic rules, enable autostart, or open the config file.

| Flag | Purpose |
|------|---------|
| `--config <path>` | Use a custom config file |
| `--reset` | Rewrite the config with defaults and exit |
| `--restore` | Put the taskbar back to normal and exit (after a crash) |
| `--verbose` | Also log to stderr |
| `--version` | Print the version |

Files: config `%APPDATA%\TaskGlass\config.json`, log `%LOCALAPPDATA%\TaskGlass\taskglass.log`.

See [docs/configuration.md](docs/configuration.md) for every option and [docs/architecture.md](docs/architecture.md) for how it works.

## Limitations

TaskGlass uses the undocumented `SetWindowCompositionAttribute` API, the same technique as TranslucentTB. Microsoft can change its behaviour in any Windows update. Windows 10 and 11 are the targets; exact visuals (especially "clear" on the Windows 11 XAML taskbar) can differ between builds. Please [report](../../issues) your Windows build if something looks off.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Security issues: [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
