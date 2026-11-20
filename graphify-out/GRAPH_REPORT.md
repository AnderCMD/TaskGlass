# Graph Report - TaskGlass  (2026-11-20)

## Corpus Check
- 49 files · ~99,000 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 679 nodes · 1289 edges · 49 communities (44 shown, 5 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 112 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `05ca89d4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main
- README.md
- winapi_windows.go
- hotkey_enabled
- App
- Controller
- tap.c
- opacity
- config.schema.json
- Changelog
- tray_icon
- PULL_REQUEST_TEMPLATE.md
- definitions
- preview.go
- github.com/AnderCMD/TaskGlass
- properties
- Watcher
- enum
- CLAUDE.md
- taskbars
- debounce_ms
- poll_ms
- refresh_ms
- surfaces
- opacity
- T
- Architecture
- render-assets.ps1
- apply-baseline.sh
- README.es.md
- items
- Configuration
- Contributing
- Security Policy
- default
- apps
- enabled
- hotkey
- enum
- theme
- window
- autostart_windows.go
- language

## God Nodes (most connected - your core abstractions)
1. `T()` - 41 edges
2. `App` - 38 edges
3. `window` - 35 edges
4. `Config` - 33 edges
5. `Controller` - 21 edges
6. `Render()` - 18 edges
7. `Parse()` - 17 edges
8. `Default()` - 16 edges
9. `Appearance` - 15 edges
10. `window` - 15 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `watchForQuit()` --references--> `App`  [EXTRACTED]
  cmd/taskglass/main_windows.go → internal/app/app_windows.go

## Import Cycles
- None detected.

## Communities (49 total, 5 thin omitted)

### Community 0 - "main"
Cohesion: 0.26
Nodes (10): attachConsole(), fatal(), Logger, main(), newLogger(), signalEvent(), signalQuit(), versionString() (+2 more)

### Community 1 - "README.md"
Cohesion: 0.18
Nodes (10): Before you start, Contributing, Features, How it works, Install, License, Limitations and safety, Roadmap (+2 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.11
Nodes (31): blurPolicy(), TestToPolicy(), HWND, Logger, Time, lighten(), New(), toPolicy() (+23 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.07
Nodes (25): App, host, menuItems, Options, panicGuard, Int32, check(), Bool (+17 more)

### Community 5 - "Controller"
Cohesion: 0.07
Nodes (28): Attr, GUID, Handler, Duration, TaskbarAnchor(), Duration, TrayProbe(), Unhook() (+20 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (72): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+64 more)

### Community 7 - "opacity"
Cohesion: 0.19
Nodes (14): $ref, properties, type, $ref, properties, appearance, rule, $ref (+6 more)

### Community 8 - "config.schema.json"
Cohesion: 0.29
Nodes (6): additionalProperties, description, $id, $schema, title, type

### Community 9 - "Changelog"
Cohesion: 0.08
Nodes (23): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, [1.1.1] - 2026-11-18, [1.2.0] - 2026-11-20, [1.2.1] - 2026-11-20, Added, Added (+15 more)

### Community 10 - "tray_icon"
Cohesion: 0.24
Nodes (10): tray_icon, enum, default, description, enum, type, auto, color (+2 more)

### Community 12 - "definitions"
Cohesion: 0.25
Nodes (8): description, type, description, pattern, type, definitions, accent, color

### Community 13 - "preview.go"
Cohesion: 0.19
Nodes (30): at(), blur(), boxH(), boxV(), clamp8(), Color(), drawFlyout(), drawTaskbarContents() (+22 more)

### Community 18 - "properties"
Cohesion: 0.15
Nodes (13): description, $ref, properties, maximized, $schema, version, xaml_fix, type (+5 more)

### Community 19 - "Watcher"
Cohesion: 0.33
Nodes (5): msg, point, Watcher, Logger, Start()

### Community 20 - "enum"
Cohesion: 0.25
Nodes (8): default, enum, type, log_level, debug, error, info, warn

### Community 23 - "taskbars"
Cohesion: 0.29
Nodes (7): taskbars, default, description, enum, type, all, primary

### Community 24 - "debounce_ms"
Cohesion: 0.33
Nodes (6): default, description, maximum, minimum, type, debounce_ms

### Community 25 - "poll_ms"
Cohesion: 0.33
Nodes (6): default, description, maximum, minimum, type, poll_ms

### Community 26 - "refresh_ms"
Cohesion: 0.33
Nodes (6): refresh_ms, default, description, maximum, minimum, type

### Community 27 - "surfaces"
Cohesion: 0.25
Nodes (8): surfaces, tray_overflow, additionalProperties, description, properties, type, description, $ref

### Community 28 - "opacity"
Cohesion: 0.40
Nodes (5): opacity, description, maximum, minimum, type

### Community 29 - "T"
Cohesion: 0.06
Nodes (57): Accent, Appearance, AppRule, Config, Percent, Rule, Surface, Surfaces (+49 more)

### Community 30 - "Architecture"
Cohesion: 0.18
Nodes (10): Appearance and assets, Architecture, Design choices, How the effect works, Installer and release, Lifecycle and recovery, Main loop, Removing the Windows 11 XAML fill (+2 more)

### Community 34 - "README.es.md"
Cohesion: 0.18
Nodes (10): Antes de empezar, Características, Contribuir, Cómo funciona, Hoja de ruta, Instalación, Licencia, Limitaciones y seguridad (+2 more)

### Community 36 - "items"
Cohesion: 0.18
Nodes (11): items, description, minLength, pattern, type, additionalProperties, properties, required (+3 more)

### Community 37 - "Configuration"
Cohesion: 0.22
Nodes (8): Appearance, Command line, Configuration, Light and dark mode, Options, Tray flyout, Which rule wins, Windows setting that overrides the bar

### Community 38 - "Contributing"
Cohesion: 0.40
Nodes (4): Branching, Contributing, Development, Knowledge graph (optional)

### Community 39 - "Security Policy"
Cohesion: 0.40
Nodes (4): Reporting a vulnerability, Scope, Security Policy, Supported versions

### Community 40 - "default"
Cohesion: 0.50
Nodes (4): additionalProperties, description, type, default

### Community 41 - "apps"
Cohesion: 0.50
Nodes (4): description, maxItems, type, apps

### Community 42 - "enabled"
Cohesion: 0.25
Nodes (8): surface, default, description, type, enabled, additionalProperties, properties, type

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 44 - "enum"
Cohesion: 0.33
Nodes (6): enum, acrylic, blur, clear, normal, opaque

### Community 45 - "theme"
Cohesion: 0.50
Nodes (4): theme, default, description, type

### Community 46 - "window"
Cohesion: 0.06
Nodes (41): Bitmap, CheckBox, ComboBox, Composite, ImageView, Container, window, unsafeSizeofChooseColor() (+33 more)

### Community 48 - "autostart_windows.go"
Cohesion: 0.38
Nodes (8): clearDisabled(), command(), Enabled(), Repair(), sameCommand(), Set(), TestCommandIsPlainQuotedPath(), TestSameCommandRecognisesLegacyDoubledBackslashes()

### Community 50 - "language"
Cohesion: 0.29
Nodes (7): default, description, enum, type, language, en, es

## Knowledge Gaps
- **165 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+160 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `main`, `winapi_windows.go`, `Controller`, `window`, `Watcher`, `T`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Why does `Config` connect `T` to `winapi_windows.go`, `App`, `window`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Why does `T()` connect `T` to `winapi_windows.go`, `App`, `Controller`, `preview.go`, `window`, `autostart_windows.go`?**
  _High betweenness centrality (0.052) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `T()` (e.g. with `.syncMenu()` and `.t()`) actually correct?**
  _`T()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _165 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10810810810810811 - nodes in this community are weakly interconnected._
- **Should `App` be split into smaller, more focused modules?**
  _Cohesion score 0.06779661016949153 - nodes in this community are weakly interconnected._