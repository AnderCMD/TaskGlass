# Graph Report - TaskGlass  (2026-11-18)

## Corpus Check
- 40 files · ~54,812 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 512 nodes · 858 edges · 36 communities (31 shown, 5 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 63 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `04fa29ac`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Context
- README.md
- winapi_windows.go
- hotkey_enabled
- App
- Controller
- tap.c
- enabled
- config.schema.json
- [1.0.0] - 2026-11-17
- theme
- PULL_REQUEST_TEMPLATE.md
- version
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
- surfaces
- xaml_fix
- Config
- render-assets.ps1
- apply-baseline.sh
- hotkey
- autostart_windows.go

## God Nodes (most connected - your core abstractions)
1. `App` - 30 edges
2. `Config` - 23 edges
3. `Controller` - 21 edges
4. `Default()` - 15 edges
5. `Parse()` - 15 edges
6. `logf_()` - 12 edges
7. `main()` - 11 edges
8. `Context` - 11 edges
9. `Appearance` - 10 edges
10. `Resolve()` - 10 edges

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

## Communities (36 total, 5 thin omitted)

### Community 0 - "Context"
Cohesion: 0.15
Nodes (15): Attr, Handler, NeedsWindowScan(), Resolve(), T, TestDisabledResolvesToNormal(), TestFirstMatchingAppRuleWins(), TestNeedsWindowScan() (+7 more)

### Community 1 - "README.md"
Cohesion: 0.05
Nodes (39): Branching, Contributing, Development, Knowledge graph (optional), Appearance and assets, Architecture, Design choices, How the effect works (+31 more)

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
Cohesion: 0.07
Nodes (26): GUID, Duration, TaskbarAnchor(), Duration, TrayProbe(), Handle, SetCallWndProcHook(), Unhook() (+18 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (68): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+60 more)

### Community 7 - "enabled"
Cohesion: 0.06
Nodes (41): $ref, properties, type, description, items, maxItems, type, $ref (+33 more)

### Community 8 - "config.schema.json"
Cohesion: 0.08
Nodes (25): description, enum, type, additionalProperties, description, pattern, type, definitions (+17 more)

### Community 9 - "[1.0.0] - 2026-11-17"
Cohesion: 0.15
Nodes (12): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, Added, Added, Changed, Changed, Changed (+4 more)

### Community 10 - "theme"
Cohesion: 0.16
Nodes (14): theme, tray_icon, default, description, enum, type, default, description (+6 more)

### Community 12 - "version"
Cohesion: 0.67
Nodes (3): version, const, description

### Community 13 - "TaskGlassHookProc"
Cohesion: 0.50
Nodes (4): LPARAM, LRESULT, TaskGlassHookProc(), WPARAM

### Community 18 - "properties"
Cohesion: 0.22
Nodes (9): description, $ref, properties, maximized, $schema, start_menu, type, description (+1 more)

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

### Community 29 - "Config"
Cohesion: 0.07
Nodes (48): attachConsole(), fatal(), Logger, main(), newLogger(), signalQuit(), versionString(), watchForQuit() (+40 more)

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 46 - "autostart_windows.go"
Cohesion: 0.83
Nodes (3): command(), Enabled(), Set()

## Knowledge Gaps
- **148 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+143 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `Controller`, `winapi_windows.go`, `Watcher`, `Config`?**
  _High betweenness centrality (0.068) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `hotkey_enabled`, `enabled`, `config.schema.json`, `theme`, `hotkey`, `version`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`, `xaml_fix`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Why does `Config` connect `Config` to `Context`, `winapi_windows.go`, `App`?**
  _High betweenness centrality (0.056) - this node is a cross-community bridge._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _148 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Context` be split into smaller, more focused modules?**
  _Cohesion score 0.1471861471861472 - nodes in this community are weakly interconnected._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.048484848484848485 - nodes in this community are weakly interconnected._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._