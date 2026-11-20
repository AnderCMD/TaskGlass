# Graph Report - TaskGlass  (2026-11-20)

## Corpus Check
- 50 files · ~98,812 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 622 nodes · 1081 edges · 53 communities (47 shown, 6 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 91 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cc5eb08c`
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
- config.go
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
- fakeHost
- enabled
- hotkey
- enum
- theme
- Server
- exe
- autostart_windows.go
- xaml_fix
- language
- main.go

## God Nodes (most connected - your core abstractions)
1. `T()` - 40 edges
2. `App` - 38 edges
3. `Config` - 35 edges
4. `Controller` - 21 edges
5. `Default()` - 20 edges
6. `Parse()` - 18 edges
7. `Server` - 17 edges
8. `main()` - 13 edges
9. `logf_()` - 12 edges
10. `host` - 11 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  tools/genicon/main.go → internal/config/config.go

## Import Cycles
- None detected.

## Communities (53 total, 6 thin omitted)

### Community 0 - "main"
Cohesion: 0.26
Nodes (10): attachConsole(), fatal(), Logger, main(), newLogger(), signalEvent(), signalQuit(), versionString() (+2 more)

### Community 1 - "README.md"
Cohesion: 0.18
Nodes (10): Before you start, Contributing, Features, How it works, Install, License, Limitations and safety, Roadmap (+2 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.11
Nodes (30): blurPolicy(), TestToPolicy(), HWND, Logger, Time, lighten(), New(), toPolicy() (+22 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.06
Nodes (27): App, host, Options, panicGuard, Config, Int32, check(), Bool (+19 more)

### Community 5 - "Controller"
Cohesion: 0.06
Nodes (29): Attr, Handler, GUID, Duration, TaskbarAnchor(), Duration, TrayProbe(), Handle (+21 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (72): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+64 more)

### Community 7 - "opacity"
Cohesion: 0.25
Nodes (11): $ref, properties, type, $ref, properties, appearance, properties, $ref (+3 more)

### Community 8 - "config.schema.json"
Cohesion: 0.29
Nodes (6): additionalProperties, description, $id, $schema, title, type

### Community 9 - "Changelog"
Cohesion: 0.08
Nodes (23): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, [1.1.1] - 2026-11-18, [1.2.0] - 2026-11-20, [1.2.1] - 2026-11-20, [1.3.0] - 2026-11-20, Added (+15 more)

### Community 10 - "tray_icon"
Cohesion: 0.24
Nodes (10): tray_icon, enum, default, description, enum, type, auto, color (+2 more)

### Community 12 - "definitions"
Cohesion: 0.22
Nodes (9): description, pattern, type, definitions, color, surface, additionalProperties, properties (+1 more)

### Community 13 - "config.go"
Cohesion: 0.12
Nodes (14): menuItems, Accent, Appearance, AppRule, Percent, Rule, Surface, Surfaces (+6 more)

### Community 18 - "properties"
Cohesion: 0.22
Nodes (9): description, $ref, properties, maximized, $schema, version, type, const (+1 more)

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
Cohesion: 0.08
Nodes (48): Lang, Default(), Load(), Parse(), ParseHotkey(), Save(), TestAppRuleMatching(), TestAutoColorFollowsTheme() (+40 more)

### Community 30 - "Architecture"
Cohesion: 0.18
Nodes (10): Appearance and assets, Architecture, Design choices, How the effect works, Installer and release, Lifecycle and recovery, Main loop, Removing the Windows 11 XAML fill (+2 more)

### Community 34 - "README.es.md"
Cohesion: 0.18
Nodes (10): Antes de empezar, Características, Contribuir, Cómo funciona, Hoja de ruta, Instalación, Licencia, Limitaciones y seguridad (+2 more)

### Community 36 - "items"
Cohesion: 0.22
Nodes (9): description, items, maxItems, type, additionalProperties, required, type, apps (+1 more)

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

### Community 42 - "enabled"
Cohesion: 0.25
Nodes (8): rule, default, description, type, enabled, additionalProperties, properties, type

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 44 - "enum"
Cohesion: 0.22
Nodes (9): description, enum, type, accent, acrylic, blur, clear, normal (+1 more)

### Community 45 - "theme"
Cohesion: 0.50
Nodes (4): theme, default, description, type

### Community 46 - "Server"
Cohesion: 0.19
Nodes (13): HandlerFunc, favicon(), Handler, loopbackHost(), New(), writeErr(), writeJSON(), Listener (+5 more)

### Community 47 - "exe"
Cohesion: 0.40
Nodes (5): description, minLength, pattern, type, exe

### Community 48 - "autostart_windows.go"
Cohesion: 0.38
Nodes (8): clearDisabled(), command(), Enabled(), Repair(), sameCommand(), Set(), TestCommandIsPlainQuotedPath(), TestSameCommandRecognisesLegacyDoubledBackslashes()

### Community 49 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 50 - "language"
Cohesion: 0.29
Nodes (7): default, description, enum, type, language, en, es

### Community 51 - "main.go"
Cohesion: 0.83
Nodes (3): entry, main(), writeICO()

## Knowledge Gaps
- **165 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+160 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `main`, `winapi_windows.go`, `Controller`, `config.go`, `Server`, `Watcher`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **Why does `Config` connect `App` to `winapi_windows.go`, `fakeHost`, `config.go`, `Server`, `T`?**
  _High betweenness centrality (0.073) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `hotkey_enabled`, `items`, `config.schema.json`, `default`, `enabled`, `hotkey`, `tray_icon`, `theme`, `xaml_fix`, `language`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Are the 16 inferred relationships involving `Default()` (e.g. with `main()` and `TestDefaultIsValid()`) actually correct?**
  _`Default()` has 16 INFERRED edges - model-reasoned connections that need verification._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _165 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10810810810810811 - nodes in this community are weakly interconnected._
- **Should `App` be split into smaller, more focused modules?**
  _Cohesion score 0.06448412698412699 - nodes in this community are weakly interconnected._