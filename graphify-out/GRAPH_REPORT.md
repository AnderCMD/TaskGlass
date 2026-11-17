# Graph Report - TaskGlass  (2026-11-17)

## Corpus Check
- 30 files · ~12,500 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 245 nodes · 463 edges · 23 communities (20 shown, 3 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `71608a43`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Config
- TaskGlass
- winapi_windows.go
- Resolve
- App
- Controller
- HRESULT
- TaskGlass
- tap.c
- Changelog
- logf_
- PULL_REQUEST_TEMPLATE.md
- ULONG
- github.com/AnderCMD/TaskGlass
- BOOL
- on_change
- TaskGlassHookProc
- CLAUDE.md

## God Nodes (most connected - your core abstractions)
1. `App` - 17 edges
2. `Controller` - 15 edges
3. `Config` - 14 edges
4. `logf_()` - 12 edges
5. `Default()` - 10 edges
6. `Context` - 10 edges
7. `on_change()` - 10 edges
8. `Appearance` - 9 edges
9. `Parse()` - 9 edges
10. `Load()` - 9 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `New()` --calls--> `Available()`  [INFERRED]
  internal/app/app_windows.go → internal/xaml/xaml_windows.go

## Import Cycles
- None detected.

## Communities (23 total, 3 thin omitted)

### Community 0 - "Config"
Cohesion: 0.14
Nodes (25): fatal(), Logger, main(), newLogger(), Accent, Appearance, Config, Rule (+17 more)

### Community 1 - "TaskGlass"
Cohesion: 0.08
Nodes (22): Branching, Contributing, Development, Knowledge graph (optional), Architecture, Design choices, How the effect works, Main loop (+14 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.12
Nodes (27): apply(), HWND, Logger, Time, New(), className(), cloaked(), enumTopLevel() (+19 more)

### Community 3 - "Resolve"
Cohesion: 0.16
Nodes (11): Attr, Handler, Resolve(), T, TestResolve(), TestResolveReturnsMatchingAppearance(), Level, Record (+3 more)

### Community 4 - "App"
Cohesion: 0.15
Nodes (10): App, CancelFunc, Logger, Time, New(), command(), Enabled(), Set() (+2 more)

### Community 5 - "Controller"
Cohesion: 0.13
Nodes (15): Duration, TaskbarAnchor(), WindowThreadPID(), Bool, CancelFunc, Handle, Logger, moduleLoaded() (+7 more)

### Community 6 - "HRESULT"
Cohesion: 0.23
Nodes (15): HRESULT, IObjectWithSite, IUnknown, C_QI(), DllCanUnloadNow(), DllGetClassObject(), F_Create(), F_QI() (+7 more)

### Community 7 - "TaskGlass"
Cohesion: 0.25
Nodes (8): Características, Contribuir, Cómo funciona, Instalación, Licencia, Limitaciones y seguridad, TaskGlass, Uso

### Community 8 - "tap.c"
Cohesion: 0.32
Nodes (12): DWORD, LONG, LPVOID, advise_thread(), any_tracked(), dispatch_apply(), init_call_thread(), install_thread() (+4 more)

### Community 9 - "Changelog"
Cohesion: 0.50
Nodes (3): Added, Changelog, [Unreleased]

### Community 10 - "logf_"
Cohesion: 0.28
Nodes (13): HSTRING, Inst, InstanceHandle, apply_all(), H_Invoke(), hs_eq(), logf_(), scan_from() (+5 more)

### Community 12 - "ULONG"
Cohesion: 0.19
Nodes (13): IClassFactory, IDispatchedHandler_, IVisualTreeServiceCallback, C_AddRef(), C_Release(), F_AddRef(), F_Lock(), F_Release() (+5 more)

### Community 18 - "BOOL"
Cohesion: 0.50
Nodes (5): HINSTANCE, BOOL, DllMain(), host_is_explorer(), wants()

### Community 19 - "on_change"
Cohesion: 0.70
Nodes (5): C_OnChange(), on_change(), ParentChildRelation, VisualElement, VisualMutationType

### Community 20 - "TaskGlassHookProc"
Cohesion: 0.50
Nodes (4): LPARAM, LRESULT, TaskGlassHookProc(), WPARAM

## Knowledge Gaps
- **33 isolated node(s):** `github.com/AnderCMD/TaskGlass`, `wndClassEx`, `accentPolicy`, `windowCompositionAttribData`, `What and why` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Config` connect `Config` to `winapi_windows.go`, `Resolve`, `App`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `Config`, `winapi_windows.go`, `Controller`?**
  _High betweenness centrality (0.062) - this node is a cross-community bridge._
- **Why does `Context` connect `Resolve` to `winapi_windows.go`, `App`, `Controller`?**
  _High betweenness centrality (0.059) - this node is a cross-community bridge._
- **What connects `github.com/AnderCMD/TaskGlass`, `wndClassEx`, `accentPolicy` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Config` be split into smaller, more focused modules?**
  _Cohesion score 0.1431451612903226 - nodes in this community are weakly interconnected._
- **Should `TaskGlass` be split into smaller, more focused modules?**
  _Cohesion score 0.082010582010582 - nodes in this community are weakly interconnected._
- **Should `winapi_windows.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12096774193548387 - nodes in this community are weakly interconnected._