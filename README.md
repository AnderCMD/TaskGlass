# TaskGlass

[![CI](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml/badge.svg)](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A tiny tray app that makes the Windows taskbar **fully transparent, blurred or acrylic**, so only the icons remain. Written in Go (no cgo) plus a small native helper DLL, no runtime to install.

[Leer en español](README.es.md)

## Features

- **Clear / Blur / Acrylic / Opaque / Normal** taskbar styles, including truly transparent on the Windows 11 XAML taskbar
- **Dynamic rules**: become opaque when a window is maximized, restyle while the Start menu is open
- **Multi-monitor**: every taskbar is handled independently, including ones added later
- **Survives Explorer restarts** and silently re-applies the effect
- **Hot-reloaded JSON config**; a typo never breaks a running session (the previous settings are kept and the error is logged)
- Tray menu, optional start-with-Windows, single instance
- Restores the stock taskbar on exit, and offers `--restore` for crash recovery

## Install

Download the zip from the [Releases](https://github.com/AnderCMD/TaskGlass/releases) page, extract it and run `taskglass.exe`. Keep `taskglass_tap.dll` in the same folder (see [How it works](#how-it-works)). It lives in the notification area.

Build from source (Windows, Go 1.26+, and a MinGW-w64 compiler for the helper DLL, e.g. `winget install MartinStorsjo.LLVM-MinGW.UCRT`):

```bash
go build -ldflags "-H=windowsgui -s -w" -o dist/taskglass.exe ./cmd/taskglass
powershell -File scripts/build-tap.ps1   # builds dist/taskglass_tap.dll
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

## How it works

Windows exposes an (undocumented) accent API, `SetWindowCompositionAttribute`, that TaskGlass uses to blur or tint the taskbar. On Windows 11 the taskbar additionally paints its own semi-transparent grey XAML fill and border on top, so the accent alone can never be fully clear.

To remove that fill, `taskglass_tap.dll` is loaded into `explorer.exe` through the XAML Diagnostics API, finds the taskbar's `BackgroundFill` / `BackgroundStroke` rectangles and sets their opacity to 0 while TaskGlass runs. If TaskGlass exits or crashes, the DLL notices and the stock look comes back by itself. Details in [docs/architecture.md](docs/architecture.md).

## Limitations and safety

- This relies on undocumented Windows internals. A Windows update can break it; the app then falls back to accent-only effects (tints and blur still work) and logs why.
- Loading code into Explorer is what security tools call "injection" and some antivirus products may flag it. The DLL is small and its source is in [native/tap](native/tap); you can build it yourself. Set `"xaml_fix": false` in the config to never touch Explorer.
- Verified on Windows 11 build 26300 with a single monitor and auto-hide enabled. Windows 10 (no XAML taskbar) only needs the accent API and is untested; multi-monitor setups are untested. Please [report](../../issues) your Windows build if something looks off.
- If the taskbar ever looks wrong after a crash: `taskglass.exe --restore`, or restart Explorer.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Security issues: [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
