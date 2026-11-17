# Graph Report - TaskGlass  (2026-11-17)

## Corpus Check
- 24 files · ~6,392 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 140 nodes · 232 edges · 18 communities (15 shown, 3 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 30 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Config
- README.md
- winapi_windows.go
- Resolve
- App
- Manager
- TaskGlass
- TaskGlass
- main
- Changelog
- autostart_windows.go
- PULL_REQUEST_TEMPLATE.md
- TestToPolicy
- github.com/AnderCMD/TaskGlass

## God Nodes (most connected - your core abstractions)
1. `App` - 16 edges
2. `Config` - 14 edges
3. `Appearance` - 9 edges
4. `Default()` - 9 edges
5. `Load()` - 9 edges
6. `Parse()` - 8 edges
7. `Save()` - 8 edges
8. `Resolve()` - 8 edges
9. `Manager` - 8 edges
10. `main()` - 7 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Default()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Load()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Parse()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `main()` --calls--> `Save()`  [INFERRED]
  cmd/taskglass/main_windows.go → internal/config/config.go
- `TestResolve()` --calls--> `Default()`  [INFERRED]
  internal/state/state_test.go → internal/config/config.go

## Import Cycles
- None detected.

## Communities (18 total, 3 thin omitted)

### Community 0 - "Config"
Cohesion: 0.20
Nodes (18): Accent, Appearance, Config, Rule, Default(), Load(), Parse(), ParseColor() (+10 more)

### Community 1 - "README.md"
Cohesion: 0.13
Nodes (13): Branching, Contributing, Development, Knowledge graph (optional), Architecture, Design choices, How the effect works, Main loop (+5 more)

### Community 2 - "winapi_windows.go"
Cohesion: 0.22
Nodes (16): apply(), className(), cloaked(), enumTopLevel(), HWND, IsWindow(), MaximizedMonitors(), Monitor() (+8 more)

### Community 3 - "Resolve"
Cohesion: 0.16
Nodes (11): Attr, Handler, Resolve(), T, TestResolve(), TestResolveReturnsMatchingAppearance(), Level, Record (+3 more)

### Community 4 - "App"
Cohesion: 0.22
Nodes (6): App, CancelFunc, Logger, Time, New(), RWMutex

### Community 5 - "Manager"
Cohesion: 0.43
Nodes (6): HWND, Logger, Time, New(), applied, Manager

### Community 6 - "TaskGlass"
Cohesion: 0.29
Nodes (7): Contributing, Features, Install, License, Limitations, TaskGlass, Usage

### Community 7 - "TaskGlass"
Cohesion: 0.29
Nodes (7): Características, Contribuir, Instalación, Licencia, Limitaciones, TaskGlass, Uso

### Community 8 - "main"
Cohesion: 0.60
Nodes (4): fatal(), Logger, main(), newLogger()

### Community 9 - "Changelog"
Cohesion: 0.50
Nodes (3): Added, Changelog, [Unreleased]

### Community 10 - "autostart_windows.go"
Cohesion: 0.83
Nodes (3): command(), Enabled(), Set()

## Knowledge Gaps
- **28 isolated node(s):** `github.com/AnderCMD/TaskGlass`, `accentPolicy`, `windowCompositionAttribData`, `What and why`, `Checklist` (+23 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Config` connect `Config` to `winapi_windows.go`, `Resolve`, `App`?**
  _High betweenness centrality (0.134) - this node is a cross-community bridge._
- **Why does `Resolve()` connect `Resolve` to `Config`, `winapi_windows.go`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `App` connect `App` to `Config`, `Manager`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `Default()` (e.g. with `main()` and `TestDefaultIsValid()`) actually correct?**
  _`Default()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 4 inferred relationships involving `Load()` (e.g. with `.reload()` and `main()`) actually correct?**
  _`Load()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/AnderCMD/TaskGlass`, `accentPolicy`, `windowCompositionAttribData` to the rest of the system?**
  _28 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `README.md` be split into smaller, more focused modules?**
  _Cohesion score 0.1286549707602339 - nodes in this community are weakly interconnected._