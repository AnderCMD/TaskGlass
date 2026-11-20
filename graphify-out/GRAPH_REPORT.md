# Graph Report - TaskGlass  (2026-11-20)

## Corpus Check
- 59 files · ~109,426 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 812 nodes · 1582 edges · 63 communities (54 shown, 9 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 149 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fdcc0463`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main
- README.md
- winapi_windows.go
- preview.go
- Config
- Controller
- tap.c
- properties
- config.schema.json
- Changelog
- theme
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
- configuration.md
- default
- window
- enabled
- Client
- enum
- App
- Server
- .InstallUpdate
- autostart_windows.go
- New
- language
- host
- Lang
- fakeHost
- updates
- menuItems
- opacity
- tray_icon
- xaml_fix
- SystemIsLight
- window

## God Nodes (most connected - your core abstractions)
1. `T()` - 54 edges
2. `App` - 42 edges
3. `Config` - 40 edges
4. `window` - 35 edges
5. `Default()` - 24 edges
6. `Controller` - 21 edges
7. `Parse()` - 20 edges
8. `Server` - 20 edges
9. `Render()` - 18 edges
10. `Appearance` - 15 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `watchForQuit()` --calls--> `Notify()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/settings/window_windows.go

## Import Cycles
- None detected.

## Communities (63 total, 9 thin omitted)

### Community 0 - "main"
Cohesion: 0.24
Nodes (11): App, attachConsole(), fatal(), Logger, main(), newLogger(), signalEvent(), signalQuit() (+3 more)

### Community 1 - "README.md"
Cohesion: 0.18
Nodes (10): Before you start, Contributing, Features, How it works, Install, License, Limitations and safety, Roadmap (+2 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.10
Nodes (33): blurPolicy(), TestToPolicy(), HWND, Logger, Time, lighten(), New(), toPolicy() (+25 more)

### Community 3 - "preview.go"
Cohesion: 0.19
Nodes (30): at(), blur(), boxH(), boxV(), clamp8(), Color(), drawFlyout(), drawTaskbarContents() (+22 more)

### Community 5 - "Controller"
Cohesion: 0.06
Nodes (30): Attr, Handler, GUID, Duration, TaskbarAnchor(), Duration, TrayProbe(), Handle (+22 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (72): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+64 more)

### Community 7 - "properties"
Cohesion: 0.20
Nodes (11): $ref, properties, $ref, description, minLength, pattern, type, properties (+3 more)

### Community 8 - "config.schema.json"
Cohesion: 0.29
Nodes (6): additionalProperties, description, $id, $schema, title, type

### Community 9 - "Changelog"
Cohesion: 0.07
Nodes (29): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, [1.1.1] - 2026-11-18, [1.2.0] - 2026-11-20, [1.2.1] - 2026-11-20, [1.3.0] - 2026-11-20, [1.3.1] - 2026-11-20 (+21 more)

### Community 10 - "theme"
Cohesion: 0.24
Nodes (10): theme, default, description, enum, type, enum, auto, color (+2 more)

### Community 12 - "definitions"
Cohesion: 0.29
Nodes (7): type, description, pattern, type, definitions, appearance, color

### Community 13 - "config.go"
Cohesion: 0.13
Nodes (13): Accent, Appearance, AppRule, Percent, Rule, Surface, SurfaceMaximized, Surfaces (+5 more)

### Community 18 - "properties"
Cohesion: 0.14
Nodes (14): default, description, default, description, type, type, properties, hotkey (+6 more)

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
Cohesion: 0.29
Nodes (7): rule, default, $ref, opacity, additionalProperties, properties, type

### Community 29 - "T"
Cohesion: 0.09
Nodes (47): entry, Default(), Load(), Parse(), ParseHotkey(), Save(), TestAppRuleMatching(), TestAutoColorFollowsTheme() (+39 more)

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
Nodes (9): Appearance, Command line, Configuration, Light and dark mode, Options, Tray flyout, Updates, Which rule wins (+1 more)

### Community 38 - "Contributing"
Cohesion: 0.40
Nodes (4): Branching, Contributing, Development, Knowledge graph (optional)

### Community 39 - "configuration.md"
Cohesion: 0.33
Nodes (4): Reporting a vulnerability, Scope, Security Policy, Supported versions

### Community 40 - "default"
Cohesion: 0.40
Nodes (5): additionalProperties, description, properties, type, default

### Community 41 - "window"
Cohesion: 0.06
Nodes (42): Bitmap, CheckBox, ComboBox, Composite, ImageView, SystemLocale(), systemLocale(), Container (+34 more)

### Community 42 - "enabled"
Cohesion: 0.15
Nodes (14): surface, default, description, type, additionalProperties, description, properties, $ref (+6 more)

### Community 43 - "Client"
Cohesion: 0.17
Nodes (17): updater, Mutex, githubOnly(), New(), Newer(), parseSums(), parseVersion(), fake() (+9 more)

### Community 44 - "enum"
Cohesion: 0.22
Nodes (9): description, enum, type, accent, acrylic, blur, clear, normal (+1 more)

### Community 45 - "App"
Cohesion: 0.14
Nodes (8): Int32, App, Bool, CancelFunc, Once, ApplyMenuTheme(), RWMutex, Value

### Community 46 - "Server"
Cohesion: 0.19
Nodes (12): HandlerFunc, favicon(), Handler, Mutex, loopbackHost(), writeErr(), writeJSON(), Listener (+4 more)

### Community 47 - ".InstallUpdate"
Cohesion: 0.23
Nodes (7): App, Ask(), CopyToClipboard(), edgePath(), Notify(), Show(), RunInstaller()

### Community 48 - "autostart_windows.go"
Cohesion: 0.38
Nodes (8): clearDisabled(), command(), Enabled(), Repair(), sameCommand(), Set(), TestCommandIsPlainQuotedPath(), TestSameCommandRecognisesLegacyDoubledBackslashes()

### Community 49 - "New"
Cohesion: 0.22
Nodes (9): Options, panicGuard, Duration, Logger, Time, New(), newPanicGuard(), TestPanicGuardAllowsBurstsUpToMax() (+1 more)

### Community 50 - "language"
Cohesion: 0.29
Nodes (7): default, description, enum, type, language, en, es

### Community 53 - "Lang"
Cohesion: 0.25
Nodes (8): Lang, Keys(), Resolve(), Table(), TestEveryKeyIsTranslated(), TestFallback(), TestResolve(), stateResp

### Community 55 - "updates"
Cohesion: 0.29
Nodes (7): updates, default, description, enum, type, notify, off

### Community 56 - "menuItems"
Cohesion: 0.40
Nodes (3): menuItems, Theme, MenuItem

### Community 57 - "opacity"
Cohesion: 0.40
Nodes (5): opacity, description, maximum, minimum, type

### Community 58 - "tray_icon"
Cohesion: 0.50
Nodes (4): tray_icon, default, description, type

### Community 59 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 60 - "SystemIsLight"
Cohesion: 0.83
Nodes (3): AppsAreLight(), readBool(), SystemIsLight()

## Knowledge Gaps
- **178 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+173 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Config` connect `Config` to `winapi_windows.go`, `window`, `config.go`, `App`, `Server`, `New`, `host`, `Lang`, `fakeHost`, `menuItems`, `T`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `winapi_windows.go`, `Config`, `Controller`, `window`, `Client`, `Server`, `New`, `host`, `Watcher`, `menuItems`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Why does `T()` connect `T` to `winapi_windows.go`, `preview.go`, `Config`, `Controller`, `window`, `Client`, `App`, `autostart_windows.go`, `New`, `Lang`?**
  _High betweenness centrality (0.064) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `T()` (e.g. with `.syncMenu()` and `.tr()`) actually correct?**
  _`T()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 20 inferred relationships involving `Default()` (e.g. with `main()` and `TestDefaultIsValid()`) actually correct?**
  _`Default()` has 20 INFERRED edges - model-reasoned connections that need verification._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _178 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0975609756097561 - nodes in this community are weakly interconnected._