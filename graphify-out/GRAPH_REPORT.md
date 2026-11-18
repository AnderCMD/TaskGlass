# Graph Report - TaskGlass  (2026-11-18)

## Corpus Check
- 45 files · ~59,231 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 570 nodes · 999 edges · 48 communities (43 shown, 5 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 97 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `40c82ec0`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Config
- README.md
- winapi_windows.go
- hotkey_enabled
- App
- Controller
- tap.c
- properties
- definitions
- [1.0.0] - 2026-11-17
- theme
- PULL_REQUEST_TEMPLATE.md
- Manager
- TaskGlassGetMsgProc
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
- config_test.go
- multi
- render-assets.ps1
- apply-baseline.sh
- items
- enabled
- config.schema.json
- enum
- opacity
- exe
- anchor_windows.go
- window_opacity
- hotkey
- tray_icon
- properties
- autostart_windows.go
- start_menu

## God Nodes (most connected - your core abstractions)
1. `App` - 30 edges
2. `Config` - 27 edges
3. `Default()` - 21 edges
4. `Controller` - 21 edges
5. `Parse()` - 17 edges
6. `main()` - 12 edges
7. `logf_()` - 12 edges
8. `Appearance` - 11 edges
9. `Context` - 11 edges
10. `Resolve()` - 11 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `ApplyPreset()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/presets.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go

## Import Cycles
- None detected.

## Communities (48 total, 5 thin omitted)

### Community 0 - "Config"
Cohesion: 0.09
Nodes (34): Accent, Appearance, AppRule, Config, Percent, Preset, Rule, Surface (+26 more)

### Community 1 - "README.md"
Cohesion: 0.05
Nodes (40): Branching, Contributing, Development, Knowledge graph (optional), Appearance and assets, Architecture, Design choices, How the effect works (+32 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.11
Nodes (36): T, TestToPolicy(), HWND, Logger, Time, New(), toPolicy(), AppWindows() (+28 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.10
Nodes (22): App, menuItems, Options, panicGuard, check(), Bool, CancelFunc, Duration (+14 more)

### Community 5 - "Controller"
Cohesion: 0.08
Nodes (28): GUID, Duration, TrayProbe(), Handle, PostThreadNull(), ProcessGUIThreads(), SetCallWndProcHook(), SetGetMessageHook() (+20 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (70): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+62 more)

### Community 7 - "properties"
Cohesion: 0.17
Nodes (15): $ref, properties, type, $ref, additionalProperties, description, properties, type (+7 more)

### Community 8 - "definitions"
Cohesion: 0.25
Nodes (8): description, type, description, pattern, type, definitions, accent, color

### Community 9 - "[1.0.0] - 2026-11-17"
Cohesion: 0.15
Nodes (12): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, Added, Added, Changed, Changed, Changed (+4 more)

### Community 10 - "theme"
Cohesion: 0.24
Nodes (10): theme, default, description, enum, type, enum, auto, color (+2 more)

### Community 12 - "Manager"
Cohesion: 0.19
Nodes (12): T, pct(), TestWanted(), Alpha(), HWND, Logger, New(), T (+4 more)

### Community 13 - "TaskGlassGetMsgProc"
Cohesion: 0.60
Nodes (5): LPARAM, LRESULT, TaskGlassGetMsgProc(), TaskGlassHookProc(), WPARAM

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
Cohesion: 0.18
Nodes (11): start_menu_flyout, surfaces, tray_overflow, description, $ref, additionalProperties, description, properties (+3 more)

### Community 28 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 29 - "config_test.go"
Cohesion: 0.11
Nodes (36): attachConsole(), fatal(), Logger, main(), newLogger(), signalQuit(), versionString(), watchForQuit() (+28 more)

### Community 30 - "multi"
Cohesion: 0.24
Nodes (5): Attr, Handler, Level, Record, multi

### Community 34 - "items"
Cohesion: 0.22
Nodes (9): description, items, maxItems, type, additionalProperties, required, type, apps (+1 more)

### Community 36 - "enabled"
Cohesion: 0.25
Nodes (8): surface, default, description, type, enabled, additionalProperties, properties, type

### Community 37 - "config.schema.json"
Cohesion: 0.29
Nodes (6): additionalProperties, description, $id, $schema, title, type

### Community 38 - "enum"
Cohesion: 0.33
Nodes (6): enum, acrylic, blur, clear, normal, opaque

### Community 39 - "opacity"
Cohesion: 0.40
Nodes (5): opacity, description, maximum, minimum, type

### Community 40 - "exe"
Cohesion: 0.40
Nodes (5): description, minLength, pattern, type, exe

### Community 41 - "anchor_windows.go"
Cohesion: 0.40
Nodes (5): Duration, TaskbarAnchor(), msg, point, wndClassEx

### Community 42 - "window_opacity"
Cohesion: 0.40
Nodes (5): window_opacity, description, maximum, minimum, type

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 44 - "tray_icon"
Cohesion: 0.50
Nodes (4): tray_icon, default, description, type

### Community 45 - "properties"
Cohesion: 0.50
Nodes (4): rule, additionalProperties, properties, type

### Community 46 - "autostart_windows.go"
Cohesion: 0.83
Nodes (3): command(), Enabled(), Set()

### Community 47 - "start_menu"
Cohesion: 0.67
Nodes (3): start_menu, description, $ref

## Knowledge Gaps
- **156 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+151 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Config` connect `Config` to `Manager`, `winapi_windows.go`, `App`, `config_test.go`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `Config`, `config_test.go`, `Watcher`, `Controller`?**
  _High betweenness centrality (0.063) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `items`, `hotkey_enabled`, `enabled`, `config.schema.json`, `properties`, `theme`, `hotkey`, `tray_icon`, `start_menu`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`, `xaml_fix`?**
  _High betweenness centrality (0.059) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `Default()` (e.g. with `main()` and `TestDefaultIsValid()`) actually correct?**
  _`Default()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _156 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Config` be split into smaller, more focused modules?**
  _Cohesion score 0.08549019607843138 - nodes in this community are weakly interconnected._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.04734299516908213 - nodes in this community are weakly interconnected._