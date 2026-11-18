# Graph Report - TaskGlass  (2026-11-17)

## Corpus Check
- 41 files · ~56,312 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 533 nodes · 913 edges · 45 communities (40 shown, 5 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 73 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `85d49618`
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
- opacity
- definitions
- [1.0.0] - 2026-11-17
- theme
- PULL_REQUEST_TEMPLATE.md
- main
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
- autostart_windows.go
- multi
- render-assets.ps1
- apply-baseline.sh
- main.go
- enum
- items
- enabled
- config.schema.json
- exe
- default
- properties
- hotkey
- tray_icon

## God Nodes (most connected - your core abstractions)
1. `App` - 30 edges
2. `Config` - 23 edges
3. `Controller` - 20 edges
4. `Default()` - 15 edges
5. `Parse()` - 14 edges
6. `main()` - 12 edges
7. `logf_()` - 12 edges
8. `Context` - 10 edges
9. `Resolve()` - 10 edges
10. `Appearance` - 9 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `RestoreAll()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/winglass/winglass_windows.go

## Import Cycles
- None detected.

## Communities (45 total, 5 thin omitted)

### Community 0 - "Config"
Cohesion: 0.08
Nodes (43): Accent, Appearance, AppRule, Config, Rule, Surface, Surfaces, Taskbars (+35 more)

### Community 1 - "README.md"
Cohesion: 0.05
Nodes (39): Branching, Contributing, Development, Knowledge graph (optional), Appearance and assets, Architecture, Design choices, How the effect works (+31 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.08
Nodes (44): T, TestToPolicy(), HWND, Logger, Time, New(), toPolicy(), className() (+36 more)

### Community 3 - "hotkey_enabled"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey_enabled

### Community 4 - "App"
Cohesion: 0.10
Nodes (22): App, menuItems, Options, panicGuard, check(), Bool, CancelFunc, Duration (+14 more)

### Community 5 - "Controller"
Cohesion: 0.09
Nodes (27): Duration, TaskbarAnchor(), Handle, PostThreadNull(), ProcessGUIThreads(), SetCallWndProcHook(), SetGetMessageHook(), Unhook() (+19 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (69): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+61 more)

### Community 7 - "opacity"
Cohesion: 0.25
Nodes (11): $ref, properties, type, $ref, properties, appearance, properties, $ref (+3 more)

### Community 8 - "definitions"
Cohesion: 0.20
Nodes (10): description, pattern, type, definitions, color, opacity, description, maximum (+2 more)

### Community 9 - "[1.0.0] - 2026-11-17"
Cohesion: 0.15
Nodes (12): [1.0.0] - 2026-11-17, [1.0.1] - 2026-11-17, [1.1.0] - 2026-11-18, Added, Added, Changed, Changed, Changed (+4 more)

### Community 10 - "theme"
Cohesion: 0.24
Nodes (10): theme, default, description, enum, type, enum, auto, color (+2 more)

### Community 12 - "main"
Cohesion: 0.31
Nodes (8): attachConsole(), fatal(), Logger, main(), newLogger(), signalQuit(), versionString(), watchForQuit()

### Community 13 - "TaskGlassGetMsgProc"
Cohesion: 0.60
Nodes (5): LPARAM, LRESULT, TaskGlassGetMsgProc(), TaskGlassHookProc(), WPARAM

### Community 18 - "properties"
Cohesion: 0.17
Nodes (12): description, $ref, properties, maximized, $schema, start_menu, version, type (+4 more)

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
Cohesion: 0.14
Nodes (14): description, $ref, explorer, start_menu_flyout, surfaces, tray_overflow, description, $ref (+6 more)

### Community 28 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 29 - "autostart_windows.go"
Cohesion: 0.83
Nodes (3): command(), Enabled(), Set()

### Community 30 - "multi"
Cohesion: 0.24
Nodes (5): Attr, Handler, Level, Record, multi

### Community 34 - "main.go"
Cohesion: 0.83
Nodes (3): entry, main(), writeICO()

### Community 36 - "enum"
Cohesion: 0.22
Nodes (9): description, enum, type, accent, acrylic, blur, clear, normal (+1 more)

### Community 37 - "items"
Cohesion: 0.22
Nodes (9): description, items, maxItems, type, additionalProperties, required, type, apps (+1 more)

### Community 38 - "enabled"
Cohesion: 0.25
Nodes (8): surface, default, description, type, enabled, additionalProperties, properties, type

### Community 39 - "config.schema.json"
Cohesion: 0.29
Nodes (6): additionalProperties, description, $id, $schema, title, type

### Community 40 - "exe"
Cohesion: 0.40
Nodes (5): description, minLength, pattern, type, exe

### Community 41 - "default"
Cohesion: 0.50
Nodes (4): additionalProperties, description, type, default

### Community 42 - "properties"
Cohesion: 0.50
Nodes (4): rule, additionalProperties, properties, type

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 44 - "tray_icon"
Cohesion: 0.50
Nodes (4): tray_icon, default, description, type

## Knowledge Gaps
- **152 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+147 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `properties` connect `properties` to `hotkey_enabled`, `items`, `enabled`, `config.schema.json`, `default`, `theme`, `hotkey`, `tray_icon`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`, `xaml_fix`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `Config`, `Watcher`, `main`, `Controller`?**
  _High betweenness centrality (0.057) - this node is a cross-community bridge._
- **Why does `Config` connect `Config` to `winapi_windows.go`, `App`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _152 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Config` be split into smaller, more focused modules?**
  _Cohesion score 0.08333333333333333 - nodes in this community are weakly interconnected._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.048484848484848485 - nodes in this community are weakly interconnected._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0803633822501747 - nodes in this community are weakly interconnected._