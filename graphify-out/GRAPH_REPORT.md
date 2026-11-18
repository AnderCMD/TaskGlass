# Graph Report - TaskGlass  (2026-11-17)

## Corpus Check
- 38 files · ~52,003 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 472 nodes · 790 edges · 36 communities (31 shown, 5 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 57 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `928fb309`
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
- items
- config.schema.json
- [1.0.0] - 2026-11-17
- theme
- PULL_REQUEST_TEMPLATE.md
- main
- TaskGlassHookProc
- github.com/AnderCMD/TaskGlass
- properties
- Watcher
- enum
- CLAUDE.md
- taskbars
- debounce_ms
- poll_ms
- refresh_ms
- version
- xaml_fix
- autostart_windows.go
- start_menu
- render-assets.ps1
- apply-baseline.sh
- main.go

## God Nodes (most connected - your core abstractions)
1. `App` - 30 edges
2. `Config` - 22 edges
3. `Controller` - 15 edges
4. `Default()` - 14 edges
5. `Parse()` - 13 edges
6. `logf_()` - 12 edges
7. `main()` - 11 edges
8. `Context` - 10 edges
9. `Resolve()` - 10 edges
10. `on_change()` - 10 edges

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

## Communities (36 total, 5 thin omitted)

### Community 0 - "Config"
Cohesion: 0.09
Nodes (40): Accent, Appearance, AppRule, Config, Rule, Taskbars, Theme, TrayIcon (+32 more)

### Community 1 - "README.md"
Cohesion: 0.05
Nodes (37): Branching, Contributing, Development, Knowledge graph (optional), Appearance and assets, Architecture, Design choices, How the effect works (+29 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.12
Nodes (29): T, TestToPolicy(), HWND, Logger, Time, New(), toPolicy(), className() (+21 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.11
Nodes (21): App, menuItems, Options, panicGuard, check(), Bool, CancelFunc, Duration (+13 more)

### Community 5 - "Controller"
Cohesion: 0.08
Nodes (25): Attr, Handler, Duration, TaskbarAnchor(), Handle, SetCallWndProcHook(), Unhook(), Available() (+17 more)

### Community 6 - "tap.c"
Cohesion: 0.08
Nodes (66): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+58 more)

### Community 7 - "items"
Cohesion: 0.07
Nodes (35): $ref, properties, description, items, maxItems, type, $ref, additionalProperties (+27 more)

### Community 8 - "config.schema.json"
Cohesion: 0.07
Nodes (27): description, enum, type, additionalProperties, type, description, pattern, type (+19 more)

### Community 9 - "[1.0.0] - 2026-11-17"
Cohesion: 0.17
Nodes (11): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.0.2] - 2026-11-17, Added, Added, Changed, Changed, Changed (+3 more)

### Community 10 - "theme"
Cohesion: 0.16
Nodes (14): theme, tray_icon, default, description, enum, type, default, description (+6 more)

### Community 12 - "main"
Cohesion: 0.31
Nodes (8): attachConsole(), fatal(), Logger, main(), newLogger(), signalQuit(), versionString(), watchForQuit()

### Community 13 - "TaskGlassHookProc"
Cohesion: 0.50
Nodes (4): LPARAM, LRESULT, TaskGlassHookProc(), WPARAM

### Community 18 - "properties"
Cohesion: 0.20
Nodes (10): default, description, type, description, $ref, properties, hotkey, maximized (+2 more)

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

### Community 27 - "version"
Cohesion: 0.67
Nodes (3): version, const, description

### Community 28 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 29 - "autostart_windows.go"
Cohesion: 0.83
Nodes (3): command(), Enabled(), Set()

### Community 30 - "start_menu"
Cohesion: 0.67
Nodes (3): start_menu, description, $ref

### Community 34 - "main.go"
Cohesion: 0.83
Nodes (3): entry, main(), writeICO()

## Knowledge Gaps
- **138 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+133 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `properties` connect `properties` to `hotkey_enabled`, `items`, `config.schema.json`, `theme`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `version`, `xaml_fix`, `start_menu`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `Config`, `winapi_windows.go`, `Controller`, `main`, `Watcher`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **Why does `Config` connect `Config` to `winapi_windows.go`, `App`?**
  _High betweenness centrality (0.053) - this node is a cross-community bridge._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _138 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Config` be split into smaller, more focused modules?**
  _Cohesion score 0.0899854862119013 - nodes in this community are weakly interconnected._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.05094130675526024 - nodes in this community are weakly interconnected._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._