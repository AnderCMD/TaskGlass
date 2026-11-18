# Graph Report - TaskGlass  (2026-11-18)

## Corpus Check
- 41 files · ~55,064 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 517 nodes · 868 edges · 38 communities (33 shown, 5 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 65 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ea5f9415`
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
- enabled
- config.schema.json
- [1.0.0] - 2026-11-17
- tray_icon
- PULL_REQUEST_TEMPLATE.md
- Context
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
- theme
- main.go
- hotkey
- autostart_windows.go

## God Nodes (most connected - your core abstractions)
1. `App` - 32 edges
2. `Config` - 25 edges
3. `Controller` - 21 edges
4. `Parse()` - 16 edges
5. `Default()` - 15 edges
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
- `main()` --calls--> `Parse()`  [INFERRED]
  tools/genicon/main.go → internal/config/config.go

## Import Cycles
- None detected.

## Communities (38 total, 5 thin omitted)

### Community 0 - "main"
Cohesion: 0.31
Nodes (8): attachConsole(), fatal(), Logger, main(), newLogger(), signalQuit(), versionString(), watchForQuit()

### Community 1 - "README.md"
Cohesion: 0.05
Nodes (42): Branching, Contributing, Development, Knowledge graph (optional), Appearance and assets, Architecture, Design choices, How the effect works (+34 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.11
Nodes (28): T, TestToPolicy(), HWND, Logger, Time, New(), toPolicy(), className() (+20 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.09
Nodes (24): App, menuItems, Options, panicGuard, Int32, check(), Bool, CancelFunc (+16 more)

### Community 5 - "Controller"
Cohesion: 0.08
Nodes (25): GUID, Duration, TaskbarAnchor(), Duration, TrayProbe(), NotifyNull(), Unhook(), WindowThreadPID() (+17 more)

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

### Community 10 - "tray_icon"
Cohesion: 0.24
Nodes (10): tray_icon, enum, default, description, enum, type, auto, color (+2 more)

### Community 12 - "Context"
Cohesion: 0.15
Nodes (15): Attr, Handler, NeedsWindowScan(), Resolve(), T, TestDisabledResolvesToNormal(), TestFirstMatchingAppRuleWins(), TestNeedsWindowScan() (+7 more)

### Community 13 - "TaskGlassHookProc"
Cohesion: 0.50
Nodes (4): LPARAM, LRESULT, TaskGlassHookProc(), WPARAM

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

### Community 29 - "Config"
Cohesion: 0.09
Nodes (38): Accent, Appearance, AppRule, Config, Percent, Rule, Surface, Surfaces (+30 more)

### Community 40 - "theme"
Cohesion: 0.50
Nodes (4): theme, default, description, type

### Community 41 - "main.go"
Cohesion: 0.83
Nodes (3): entry, main(), writeICO()

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 46 - "autostart_windows.go"
Cohesion: 0.83
Nodes (3): command(), Enabled(), Set()

## Knowledge Gaps
- **149 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+144 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `main`, `winapi_windows.go`, `Controller`, `Watcher`, `Config`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `hotkey_enabled`, `enabled`, `config.schema.json`, `theme`, `tray_icon`, `hotkey`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`, `xaml_fix`?**
  _High betweenness centrality (0.062) - this node is a cross-community bridge._
- **Why does `Config` connect `Config` to `winapi_windows.go`, `App`, `Context`?**
  _High betweenness centrality (0.060) - this node is a cross-community bridge._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _149 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1140819964349376 - nodes in this community are weakly interconnected._
- **Should `App` be split into smaller, more focused modules?**
  _Cohesion score 0.08973172987974098 - nodes in this community are weakly interconnected._