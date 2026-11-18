# Graph Report - TaskGlass  (2026-11-18)

## Corpus Check
- 44 files · ~58,599 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 563 nodes · 989 edges · 38 communities (33 shown, 5 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 96 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `04b14ce5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Default
- README.md
- winapi_windows.go
- hotkey_enabled
- App
- Controller
- tap.c
- properties
- config.schema.json
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
- main
- render-assets.ps1
- apply-baseline.sh
- main.go
- anchor_windows.go
- hotkey
- tray_icon

## God Nodes (most connected - your core abstractions)
1. `App` - 30 edges
2. `Config` - 27 edges
3. `Default()` - 21 edges
4. `Controller` - 20 edges
5. `Parse()` - 17 edges
6. `main()` - 12 edges
7. `logf_()` - 12 edges
8. `Appearance` - 11 edges
9. `Resolve()` - 11 edges
10. `Context` - 10 edges

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

### Community 0 - "Default"
Cohesion: 0.08
Nodes (47): Accent, Appearance, AppRule, Percent, Rule, Surface, Surfaces, TrayIcon (+39 more)

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
Cohesion: 0.08
Nodes (28): App, menuItems, Options, panicGuard, Config, Taskbars, Theme, check() (+20 more)

### Community 5 - "Controller"
Cohesion: 0.08
Nodes (28): Attr, Handler, Handle, PostThreadNull(), ProcessGUIThreads(), SetCallWndProcHook(), SetGetMessageHook(), Unhook() (+20 more)

### Community 6 - "tap.c"
Cohesion: 0.07
Nodes (70): DWORD, HINSTANCE, HRESULT, HSTRING, IClassFactory, IDispatchedHandler_, Inst, InstanceHandle (+62 more)

### Community 7 - "properties"
Cohesion: 0.05
Nodes (46): $ref, properties, type, description, items, maxItems, type, $ref (+38 more)

### Community 8 - "config.schema.json"
Cohesion: 0.08
Nodes (25): description, enum, type, additionalProperties, description, pattern, type, definitions (+17 more)

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
Cohesion: 0.18
Nodes (11): start_menu_flyout, surfaces, tray_overflow, description, $ref, additionalProperties, description, properties (+3 more)

### Community 28 - "xaml_fix"
Cohesion: 0.50
Nodes (4): xaml_fix, default, description, type

### Community 29 - "main"
Cohesion: 0.16
Nodes (17): attachConsole(), fatal(), Logger, main(), newLogger(), signalQuit(), versionString(), watchForQuit() (+9 more)

### Community 34 - "main.go"
Cohesion: 0.83
Nodes (3): entry, main(), writeICO()

### Community 41 - "anchor_windows.go"
Cohesion: 0.40
Nodes (5): Duration, TaskbarAnchor(), msg, point, wndClassEx

### Community 43 - "hotkey"
Cohesion: 0.50
Nodes (4): default, description, type, hotkey

### Community 44 - "tray_icon"
Cohesion: 0.50
Nodes (4): tray_icon, default, description, type

## Knowledge Gaps
- **156 isolated node(s):** `$schema`, `$id`, `title`, `description`, `type` (+151 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Config` connect `App` to `Default`, `winapi_windows.go`, `Manager`, `main`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `main`, `Watcher`, `Controller`?**
  _High betweenness centrality (0.061) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `hotkey_enabled`, `properties`, `config.schema.json`, `theme`, `hotkey`, `tray_icon`, `enum`, `taskbars`, `debounce_ms`, `poll_ms`, `refresh_ms`, `surfaces`, `xaml_fix`?**
  _High betweenness centrality (0.060) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `Default()` (e.g. with `main()` and `TestDefaultIsValid()`) actually correct?**
  _`Default()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `$schema`, `$id`, `title` to the rest of the system?**
  _156 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Default` be split into smaller, more focused modules?**
  _Cohesion score 0.07650273224043716 - nodes in this community are weakly interconnected._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.04734299516908213 - nodes in this community are weakly interconnected._