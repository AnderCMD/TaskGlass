# Graph Report - TaskGlass  (2026-11-20)

## Corpus Check
- 50 files · ~97,994 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 609 nodes · 1055 edges · 53 communities (48 shown, 5 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 87 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c8f79e7d`
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
- theme
- PULL_REQUEST_TEMPLATE.md
- definitions
- Config
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
- xaml_fix
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
- main.go
- enabled
- hotkey
- enum
- anchor_windows.go
- Server
- exe
- autostart_windows.go
- tray_icon
- language
- properties

## God Nodes (most connected - your core abstractions)
1. `App` - 38 edges
2. `T()` - 38 edges
3. `Config` - 35 edges
4. `Controller` - 21 edges
5. `Default()` - 19 edges
6. `Parse()` - 17 edges
7. `Server` - 17 edges
8. `main()` - 13 edges
9. `logf_()` - 12 edges
10. `Context` - 11 edges

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

## Communities (53 total, 5 thin omitted)

### Community 0 - "main"
Cohesion: 0.26
Nodes (10): attachConsole(), fatal(), Logger, main(), newLogger(), signalEvent(), signalQuit(), versionString() (+2 more)

### Community 1 - "README.md"
Cohesion: 0.18
Nodes (10): Before you start, Contributing, Features, How it works, Install, License, Limitations and safety, Roadmap (+2 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.12
Nodes (28): TestToPolicy(), HWND, Logger, Time, New(), toPolicy(), className(), cloaked() (+20 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.06
Nodes (27): App, host, menuItems, Options, panicGuard, Int32, check(), Bool (+19 more)

### Community 5 - "Controller"
Cohesion: 0.07
Nodes (25): Attr, Handler, GUID, Duration, TrayProbe(), Handle, SetCallWndProcHook(), Available() (+17 more)

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
Cohesion: 0.11
Nodes (18): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, [1.1.1] - 2026-11-18, Added, Added, Added, Changed (+10 more)

### Community 10 - "theme"
Cohesion: 0.24
Nodes (10): theme, default, description, enum, type, enum, auto, color (+2 more)

### Community 12 - "definitions"
Cohesion: 0.20
Nodes (10): description, pattern, type, definitions, color, opacity, description, maximum (+2 more)

### Community 13 - "Config"
Cohesion: 0.09
Nodes (16): Accent, Appearance, AppRule, Config, Percent, Rule, Surface, Surfaces (+8 more)

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

### Community 28 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 29 - "T"
Cohesion: 0.09
Nodes (44): Lang, Default(), Load(), Parse(), TestAppRuleMatching(), TestAutoColorFollowsTheme(), TestDefaultIsValid(), TestExplicitFalseSticks() (+36 more)

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

### Community 41 - "main.go"
Cohesion: 0.83
Nodes (3): entry, main(), writeICO()

### Community 42 - "enabled"
Cohesion: 0.25
Nodes (8): surface, default, description, type, enabled, additionalProperties, properties, type

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 44 - "enum"
Cohesion: 0.22
Nodes (9): description, enum, type, accent, acrylic, blur, clear, normal (+1 more)

### Community 45 - "anchor_windows.go"
Cohesion: 0.40
Nodes (5): Duration, TaskbarAnchor(), msg, point, wndClassEx

### Community 46 - "Server"
Cohesion: 0.20
Nodes (12): HandlerFunc, Handler, loopbackHost(), New(), writeErr(), writeJSON(), Listener, Mutex (+4 more)

### Community 47 - "exe"
Cohesion: 0.40
Nodes (5): description, minLength, pattern, type, exe

### Community 48 - "autostart_windows.go"
Cohesion: 0.38
Nodes (8): clearDisabled(), command(), Enabled(), Repair(), sameCommand(), Set(), TestCommandIsPlainQuotedPath(), TestSameCommandRecognisesLegacyDoubledBackslashes()

### Community 49 - "tray_icon"
Cohesion: 0.50
Nodes (4): tray_icon, default, description, type

### Community 50 - "language"
Cohesion: 0.29
Nodes (7): default, description, enum, type, language, en, es

### Community 51 - "properties"
Cohesion: 0.50
Nodes (4): rule, additionalProperties, properties, type

## Knowledge Gaps
- **162 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+157 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `main`, `winapi_windows.go`, `Controller`, `Config`, `Server`, `Watcher`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Why does `Config` connect `Config` to `winapi_windows.go`, `App`, `T`, `Server`?**
  _High betweenness centrality (0.072) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `hotkey_enabled`, `items`, `config.schema.json`, `default`, `enabled`, `hotkey`, `theme`, `tray_icon`, `language`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`, `xaml_fix`?**
  _High betweenness centrality (0.050) - this node is a cross-community bridge._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _162 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._
- **Should `App` be split into smaller, more focused modules?**
  _Cohesion score 0.061955965181771634 - nodes in this community are weakly interconnected._
- **Should `Controller` be split into smaller, more focused modules?**
  _Cohesion score 0.07342995169082125 - nodes in this community are weakly interconnected._